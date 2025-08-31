// helper/rendermethod.go
package helper

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

/*
Bit layout (LSB -> MSB)
bits  0-1: Drawstyle
bits  2-4: Lighting
bits  5-6: Shading
bit     7: Masked Transparency Toggle (TRANS)
bits  8-15: Texture
bits 16-19: Alpha Blend Opacity (0..15)  => shown as % = field/16 * 100
bit    20: Additive Toggle
bits 21-23: UnknownA
bit    24: Alpha Blend Toggle (BLEND)
bits 25-30: UnknownB
bit    31: Userdefined Toggle
*/

const (
	bfDrawstyleShift = 0
	bfDrawstyleMask  = 0b11

	bfLightingShift = 2
	bfLightingMask  = 0b111

	bfShadingShift = 5
	bfShadingMask  = 0b11

	bfMaskedShift = 7
	bfMaskedMask  = 0b1

	bfTextureShift = 8
	bfTextureMask  = 0xFF

	bfAlphaShift = 16
	bfAlphaMask  = 0xF

	bfAdditiveShift = 20
	bfAdditiveMask  = 0x1

	bfDynamicShift = 21
	bfDynamicMask  = 0x1

	bfAlphaToggleShift = 24
	bfAlphaToggleMask  = 0x1

	bfPrelitShift = 30
	bfPrelitMask  = 0x1

	bfUserDefinedShift = 31
	bfUserDefinedMask  = 0x1
)

// ——— Public API ————————————————————————————————————————————————————————

func RenderMethodStr(v uint32) string {
	// TRANSPARENT if completely zero.
	if v == 0 {
		return "TRANSPARENT"
	}
	// USERDEFINED_ if bit 31 is set.
	if (v>>bfUserDefinedShift)&bfUserDefinedMask == 1 {
		// “the rest of the Bits written as one integer value +1”
		rest := (v & 0x7FFFFFFF) + 1
		return fmt.Sprintf("USERDEFINED_%d", rest)
	}

	var b strings.Builder

	// Masked transparency (“TRANS”) – written first if present.
	if ((v >> bfMaskedShift) & bfMaskedMask) == 1 {
		b.WriteString("TRANS")
	}

	// Drawstyle (with the SOLIDFILL ↔ TEXTURE special case)
	draw := (v >> bfDrawstyleShift) & bfDrawstyleMask
	tex := (v >> bfTextureShift) & bfTextureMask

	if draw == 3 /* SOLIDFILL */ && tex > 0 {
		// Texture takes the place of SOLIDFILL
		b.WriteString(fmt.Sprintf("TEXTURE%d", tex))
	} else {
		b.WriteString(drawstyleName(draw))
		// NOTE: Per your spec, Texture is only written “instead of SOLIDFILL”.
		// If you ever want TEXTURE# always appended when >0, move this out.
	}

	// Lighting
	light := (v >> bfLightingShift) & bfLightingMask
	b.WriteString(lightingName(light))

	// Shading
	shade := (v >> bfShadingShift) & bfShadingMask
	b.WriteString(shadingName(shade))

	// Additive
	if ((v >> bfAdditiveShift) & bfAdditiveMask) == 1 {
		b.WriteString("ADDITIVE")
	}

	// DYNAMIC (bit 21)
	if ((v >> bfDynamicShift) & bfDynamicMask) == 1 {
		b.WriteString("DYNAMIC")
	}

	// PRELIT (bit 30)
	if ((v >> bfPrelitShift) & bfPrelitMask) == 1 {
		b.WriteString("PRELIT")
	}

	// Alpha Blend toggle
	alphaToggle := ((v >> bfAlphaToggleShift) & bfAlphaToggleMask) == 1
	if alphaToggle {
		b.WriteString("BLEND")
	}

	// Alpha Blend Opacity
	alpha := (v >> bfAlphaShift) & bfAlphaMask
	if alphaToggle || alpha > 0 {
		percent := float64(alpha) / 16.0 * 100.0
		// Always show one decimal place (matches your examples like 50.0%)
		b.WriteString(fmt.Sprintf("OPACITY%.1f%%", percent))
	}

	return b.String()
}

func RenderMethodInt(s string) uint32 {
	// Fast paths
	if s == "" || s == "TRANSPARENT" {
		return 0
	}
	if strings.HasPrefix(s, "USERDEFINED_") {
		// value = 0x80000000 | ((N-1) & 0x7FFFFFFF)
		nStr := strings.TrimPrefix(s, "USERDEFINED_")
		if n, err := strconv.ParseUint(nStr, 10, 32); err == nil {
			return 0x80000000 | (uint32(n-1) & 0x7FFFFFFF)
		}
		// If parse fails, fall through to best-effort parse.
	}

	var v uint32
	rest := s

	// TRANS
	if strings.HasPrefix(rest, "TRANS") {
		v |= 1 << bfMaskedShift
		rest = strings.TrimPrefix(rest, "TRANS")
	}

	// Drawstyle / Texture special case
	if strings.HasPrefix(rest, "TEXTURE") {
		num := takeNumber(&rest, "TEXTURE")
		v |= 3 // SOLIDFILL code in drawstyle field
		v |= (num & bfTextureMask) << bfTextureShift
	} else {
		switch {
		case strings.HasPrefix(rest, "DRAW0"):
			v |= 0
			rest = strings.TrimPrefix(rest, "DRAW0")
		case strings.HasPrefix(rest, "DRAW1"):
			v |= 1
			rest = strings.TrimPrefix(rest, "DRAW1")
		case strings.HasPrefix(rest, "WIREFRAME"):
			v |= 2
			rest = strings.TrimPrefix(rest, "WIREFRAME")
		case strings.HasPrefix(rest, "SOLIDFILL"):
			v |= 3
			rest = strings.TrimPrefix(rest, "SOLIDFILL")
		}
		// (Per spec, only replace SOLIDFILL with TEXTURE when Texture>0.
		// We do not append a separate TEXTURE# here if present in the string
		// after SOLIDFILL; adapt if you later decide otherwise.)
	}

	// Lighting
	for name, code := range lightingCodes {
		if strings.HasPrefix(rest, name) {
			v |= code << bfLightingShift
			rest = strings.TrimPrefix(rest, name)
			break
		}
	}

	// Shading
	for name, code := range shadingCodes {
		if strings.HasPrefix(rest, name) {
			v |= code << bfShadingShift
			rest = strings.TrimPrefix(rest, name)
			break
		}
	}

	// ADDITIVE
	if strings.HasPrefix(rest, "ADDITIVE") {
		v |= 1 << bfAdditiveShift
		rest = strings.TrimPrefix(rest, "ADDITIVE")
	}

	// DYNAMIC
	if strings.HasPrefix(rest, "DYNAMIC") {
		v |= 1 << bfDynamicShift
		rest = strings.TrimPrefix(rest, "DYNAMIC")
	}

	// PRELIT
	if strings.HasPrefix(rest, "PRELIT") {
		v |= 1 << bfPrelitShift
		rest = strings.TrimPrefix(rest, "PRELIT")
	}

	// BLEND
	if strings.HasPrefix(rest, "BLEND") {
		v |= 1 << bfAlphaToggleShift
		rest = strings.TrimPrefix(rest, "BLEND")
	}

	// OPACITYxx.x%
	if strings.HasPrefix(rest, "OPACITY") {
		re := regexp.MustCompile(`^OPACITY([\d.]+)%`)
		if m := re.FindStringSubmatch(rest); len(m) == 2 {
			// round to nearest 1/16 step
			percent, _ := strconv.ParseFloat(m[1], 64)
			field := uint32(math.Round((percent / 100.0) * 16.0))
			if field > 0xF {
				field = 0xF
			}
			v |= field << bfAlphaShift
			rest = rest[len(m[0]):]
		}
	}

	return v
}

// ——— helpers ————————————————————————————————————————————————————————

func drawstyleName(code uint32) string {
	switch code {
	case 0:
		return "DRAW0"
	case 1:
		return "DRAW1"
	case 2:
		return "WIREFRAME"
	case 3:
		return "SOLIDFILL"
	default:
		return "DRAW0"
	}
}

var lightingNames = [...]string{
	"ZEROINTENSITY",
	"LIGHT1",
	"CONSTANT",
	"LIGHT3",
	"AMBIENT",
	"SCALEDAMBIENT",
	"LIGHT6",
	"LIGHT7",
}
var lightingCodes = map[string]uint32{
	"ZEROINTENSITY": 0,
	"LIGHT1":        1,
	"CONSTANT":      2,
	"LIGHT3":        3,
	"AMBIENT":       4,
	"SCALEDAMBIENT": 5,
	"LIGHT6":        6,
	"LIGHT7":        7,
}

func lightingName(code uint32) string {
	if code < uint32(len(lightingNames)) {
		return lightingNames[code]
	}
	return "ZEROINTENSITY"
}

var shadingNames = [...]string{
	"SHADE0",
	"SHADE1",
	"GOURAUD1",
	"GOURAUD2",
}
var shadingCodes = map[string]uint32{
	"SHADE0":   0,
	"SHADE1":   1,
	"GOURAUD1": 2,
	"GOURAUD2": 3,
}

func shadingName(code uint32) string {
	if code < uint32(len(shadingNames)) {
		return shadingNames[code]
	}
	return "SHADE0"
}

func takeNumber(rest *string, prefix string) uint32 {
	t := strings.TrimPrefix(*rest, prefix)
	i := 0
	for i < len(t) && t[i] >= '0' && t[i] <= '9' {
		i++
	}
	num := uint32(0)
	if i > 0 {
		n, _ := strconv.ParseUint(t[:i], 10, 32)
		num = uint32(n)
	}
	*rest = t[i:]
	return num
}
