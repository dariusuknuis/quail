package wce

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
)

// Ensure Wce has:  EffectOlds []*EffectOld
// type Wce struct { ... EffectOlds []*EffectOld ... }

// ===============================
// In-memory EFF
// ===============================

type EffectOld struct {
	folders        []string
	TagIndex       int
	Header         [2]uint32
	Source         EffectOldBlock
	SourceToTarget EffectOldBlock
	Target         EffectOldBlock
}

type EffectOldBlock struct {
	Label        string
	ExtraSprites [12]string
	UnknownParam uint32
	SoundRef     uint32
	Sub          [3]EffectOldSub
	UnknownDW    [51]uint32
	UnknownF32   [12]float32
}

type EffectOldSub struct {
	PrimarySprite string
	DagIndex      uint32
	EffectType    uint32
	ColorBGRA     [4]uint8
	Gravity       float32
	SpawnNormal   [3]float32
	SpawnRadius   float32
	SpawnAngle    float32
	Lifespan      uint32
	SpawnVelocity float32
	SpawnRate     uint32
	SpawnScale    float32
}

// ===============================
// ASCII I/O (no raw here)
// ===============================

// Definition token for this ASCII block.
// Match the style of your other DEF headers; adjust if you prefer another label.
func (e *EffectOld) Definition() string {
	return "CLASSICEFFECT"
}

// WriteEffOld writes the provided records using your ASCII token system.
// You pass the folders and the tag you want this emitted under.
func (e *EffectOld) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		if err := token.SetWriter(folder); err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		// CLASSICEFFECT "<index>"
		fmt.Fprintf(w, "%s \"%d\"\n", e.Definition(), e.TagIndex)

		// Header as two DWORD ints (little-endian)
		fmt.Fprintf(w, "\tHEADER %d %d\n", e.Header[0], e.Header[1])

		writeBlock := func(name string, b *EffectOldBlock) {
			fmt.Fprintf(w, "\tBLOCK \"%s\"\n", name)
			fmt.Fprintf(w, "\t\tLABEL \"%s\"\n", b.Label)

			// Keep sprites as a count + lines (same style you used)
			fmt.Fprintf(w, "\t\tEXTRASPRITES %d\n", len(b.ExtraSprites))
			for _, s := range b.ExtraSprites {
				fmt.Fprintf(w, "\t\t\tSPRITE \"%s\"\n", s)
			}

			fmt.Fprintf(w, "\t\tUNKNOWNPARAM %d\n", b.UnknownParam)
			fmt.Fprintf(w, "\t\tSOUNDREF %d\n", b.SoundRef)

			// Sub effects: multiline, exactly as you prefer
			for j := 0; j < 3; j++ {
				sub := b.Sub[j]
				fmt.Fprintf(w, "\t\tSUBEFFECT // %d\n", j)
				fmt.Fprintf(w, "\t\t\tBLITSPRITE \"%s\"\n", sub.PrimarySprite)
				fmt.Fprintf(w, "\t\t\tDAGINDEX %d\n", sub.DagIndex)
				fmt.Fprintf(w, "\t\t\tEFFECTTYPE %d\n", sub.EffectType)
				fmt.Fprintf(w, "\t\t\tRGBA %d %d %d %d \n", sub.ColorBGRA[2], sub.ColorBGRA[1], sub.ColorBGRA[0], sub.ColorBGRA[3])
				fmt.Fprintf(w, "\t\t\tGRAVITY %0.8e\n", sub.Gravity)
				fmt.Fprintf(w, "\t\t\tSPAWNNORMAL %0.8e %0.8e %0.8e\n",
					sub.SpawnNormal[0], sub.SpawnNormal[1], sub.SpawnNormal[2])
				fmt.Fprintf(w, "\t\t\tSPAWNRADIUS %0.8e\n", sub.SpawnRadius)
				fmt.Fprintf(w, "\t\t\tSPAWNANGLE %0.8e\n", sub.SpawnAngle)
				fmt.Fprintf(w, "\t\t\tLIFESPAN %d\n", sub.Lifespan)
				fmt.Fprintf(w, "\t\t\tSPAWNVELOCITY %0.8e\n", sub.SpawnVelocity)
				fmt.Fprintf(w, "\t\t\tSPAWNRATE %d\n", sub.SpawnRate)
				fmt.Fprintf(w, "\t\t\tSPAWNSCALE %0.8e\n", sub.SpawnScale)
			}

			// Fixed 51 DWORDs, numbered with no zero padding
			for i := 0; i < 51; i++ {
				fmt.Fprintf(w, "\t\tUNKNOWNDW_%d %d\n", i+1, b.UnknownDW[i])
			}
			// Fixed 12 floats, numbered with no zero padding
			for i := 0; i < 12; i++ {
				fmt.Fprintf(w, "\t\tUNKNOWNF32_%d %0.8e\n", i+1, b.UnknownF32[i])
			}
		}

		writeBlock("Source", &e.Source)
		writeBlock("SourceToTarget", &e.SourceToTarget)
		writeBlock("Target", &e.Target)

		fmt.Fprintln(w)
	}
	return nil
}

// ReadEffOld parses the next EFFECTOLDDEF block and returns its records.
// It does not clear or touch any Wce; you can append these wherever you like.
func (e *EffectOld) Read(token *AsciiReadToken) error {
	// The definition line "CLASSICEFFECT <index>" has already been consumed by the dispatcher.
	// We just read the properties that follow it.

	records, err := token.ReadProperty("HEADER", 2)
	if err != nil {
		return err
	}
	err = parse(&e.Header, records[1:]...)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	readBlock := func(which string, dst *EffectOldBlock) error {
		records, err := token.ReadProperty("BLOCK", 1)
		if err != nil {
			return fmt.Errorf("%s block open: %w", which, err)
		}
		got := strings.Trim(records[1], "\"")
		if !strings.EqualFold(got, which) {
			return fmt.Errorf("expected BLOCK %q, got %q", which, got)
		}

		records, err = token.ReadProperty("LABEL", 1)
		if err != nil {
			return fmt.Errorf("%s label: %w", which, err)
		}
		dst.Label = records[1]

		records, err = token.ReadProperty("EXTRASPRITES", 1)
		if err != nil {
			return fmt.Errorf("%s extrasprites: %w", which, err)
		}
		cnt := 0
		err = parse(&cnt, records[1])
		if err != nil {
			return fmt.Errorf("%s extrasprites cnt: %w", which, err)
		}
		for j := 0; j < cnt && j < len(dst.ExtraSprites); j++ {
			records, err = token.ReadProperty("SPRITE", 1)
			if err != nil {
				return fmt.Errorf("%s sprite %d: %w", which, j, err)
			}
			dst.ExtraSprites[j] = records[1]
		}

		records, err = token.ReadProperty("UNKNOWNPARAM", 1)
		if err != nil {
			return err
		}
		err = parse(&dst.UnknownParam, records[1])
		if err != nil {
			return fmt.Errorf("%s unknownparam: %w", which, err)
		}

		records, err = token.ReadProperty("SOUNDREF", 1)
		if err != nil {
			return err
		}
		err = parse(&dst.SoundRef, records[1])
		if err != nil {
			return fmt.Errorf("%s soundref: %w", which, err)
		}

		for s := 0; s < 3; s++ {
			_, err := token.ReadProperty("SUB", 0)
			if err != nil {
				return fmt.Errorf("%s sub %d: %w", which, s, err)
			}
			sub := &dst.Sub[s]

			records, err = token.ReadProperty("PRIMARYSPRITE", 1)
			if err != nil {
				return fmt.Errorf("%s sub %d primarysprite: %w", which, s, err)
			}
			sub.PrimarySprite = records[1]

			records, err = token.ReadProperty("DAGINDEX", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.DagIndex, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d dagindex: %w", which, s, err)
			}

			records, err = token.ReadProperty("EFFECTTYPE", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.EffectType, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d effecttype: %w", which, s, err)
			}

			records, err = token.ReadProperty("RGBA", 4)
			if err != nil {
				return err
			}
			err = parse(&sub.ColorBGRA, records[1:]...)
			if err != nil {
				return fmt.Errorf("%s sub %d rgba: %w", which, s, err)
			}

			records, err = token.ReadProperty("GRAVITY", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.Gravity, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d gravity: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNNORMAL", 3)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnNormal, records[1:]...)
			if err != nil {
				return fmt.Errorf("%s sub %d spawnnormal: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNRADIUS", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnRadius, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d spawnradius: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNANGLE", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnAngle, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d spawnangle: %w", which, s, err)
			}

			records, err = token.ReadProperty("LIFESPAN", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.Lifespan, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d lifespan: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNVELOCITY", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnVelocity, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d spawnvelocity: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNRATE", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnRate, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d spawnrate: %w", which, s, err)
			}

			records, err = token.ReadProperty("SPAWNSCALE", 1)
			if err != nil {
				return err
			}
			err = parse(&sub.SpawnScale, records[1])
			if err != nil {
				return fmt.Errorf("%s sub %d spawnscale: %w", which, s, err)
			}
		}

		// Fixed 51 DWORDs, read numbered keys
		for i := 0; i < 51; i++ {
			key := fmt.Sprintf("UNKNOWNDW_%d", i+1)
			records, err = token.ReadProperty(key, 1)
			if err != nil {
				return err
			}
			err = parse(&dst.UnknownDW[i], records[1])
			if err != nil {
				return fmt.Errorf("%s %s parse: %w", which, key, err)
			}
		}

		// Fixed 12 floats, read numbered keys
		for i := 0; i < 12; i++ {
			key := fmt.Sprintf("UNKNOWNF32_%d", i+1)
			records, err = token.ReadProperty(key, 1)
			if err != nil {
				return fmt.Errorf("%s %s: %w", which, key, err)
			}
			err = parse(&dst.UnknownF32[i], records[1])
			if err != nil {
				return fmt.Errorf("%s %s parse: %w", which, key, err)
			}
		}

		return nil
	}

	err = readBlock("Source", &e.Source)
	if err != nil {
		return err
	}

	err = readBlock("SourceToTarget", &e.SourceToTarget)
	if err != nil {
		return err
	}

	err = readBlock("Target", &e.Target)
	if err != nil {
		return err
	}

	return nil
}

func (e *EffectOld) ToRaw(wce *Wce, dst *raw.EffOldRecord) error {
	if e == nil || dst == nil {
		return fmt.Errorf("nil receiver or destination")
	}

	// Copy header
	dst.Header = e.Header

	// Inline copy for each of the three blocks
	copyBlock := func(src *EffectOldBlock, dstB *raw.EffOldBlock) {
		dstB.Label = src.Label
		dstB.UnknownParam = src.UnknownParam
		dstB.SoundRef = src.SoundRef

		for i := 0; i < len(dstB.ExtraSprites) && i < len(src.ExtraSprites); i++ {
			dstB.ExtraSprites[i] = src.ExtraSprites[i]
		}

		for i := 0; i < 3; i++ {
			s := src.Sub[i]
			dstB.SubEffect[i].Blit = s.PrimarySprite
			dstB.SubEffect[i].DagIndex = s.DagIndex
			dstB.SubEffect[i].EffectType = s.EffectType
			dstB.SubEffect[i].ColorBGRA = s.ColorBGRA
			dstB.SubEffect[i].Gravity = s.Gravity
			dstB.SubEffect[i].SpawnNormal = s.SpawnNormal
			dstB.SubEffect[i].SpawnRadius = s.SpawnRadius
			dstB.SubEffect[i].SpawnAngle = s.SpawnAngle
			dstB.SubEffect[i].Lifespan = s.Lifespan
			dstB.SubEffect[i].SpawnVelocity = s.SpawnVelocity
			dstB.SubEffect[i].SpawnRate = s.SpawnRate
			dstB.SubEffect[i].SpawnScale = s.SpawnScale
		}

		for i := 0; i < len(dstB.UnknownDW) && i < len(src.UnknownDW); i++ {
			dstB.UnknownDW[i] = src.UnknownDW[i]
		}
		for i := 0; i < len(dstB.UnknownF32) && i < len(src.UnknownF32); i++ {
			dstB.UnknownF32[i] = src.UnknownF32[i]
		}
	}

	copyBlock(&e.Source, &dst.Source)
	copyBlock(&e.SourceToTarget, &dst.SourceToTarget)
	copyBlock(&e.Target, &dst.Target)

	return nil
}

// FromRaw fills the ASCII/in-memory EffOldDef from raw.EffOld.
// FromRaw populates an EffectOld from a raw.EffOldRecord.
func (e *EffectOld) FromRaw(_ *Wce, src *raw.EffOldRecord) error {
	if e == nil || src == nil {
		return fmt.Errorf("nil receiver or source")
	}

	e.Header = src.Header

	cpBlock := func(srcB *raw.EffOldBlock, dstB *EffectOldBlock) {
		dstB.Label = srcB.Label
		dstB.UnknownParam = srcB.UnknownParam
		dstB.SoundRef = srcB.SoundRef

		for i := 0; i < len(dstB.ExtraSprites) && i < len(srcB.ExtraSprites); i++ {
			dstB.ExtraSprites[i] = srcB.ExtraSprites[i]
		}

		for i := 0; i < 3; i++ {
			s := srcB.SubEffect[i]
			dstB.Sub[i] = EffectOldSub{
				PrimarySprite: s.Blit,
				DagIndex:      s.DagIndex,
				EffectType:    s.EffectType,
				ColorBGRA:     s.ColorBGRA,
				Gravity:       s.Gravity,
				SpawnNormal:   s.SpawnNormal,
				SpawnRadius:   s.SpawnRadius,
				SpawnAngle:    s.SpawnAngle,
				Lifespan:      s.Lifespan,
				SpawnVelocity: s.SpawnVelocity,
				SpawnRate:     s.SpawnRate,
				SpawnScale:    s.SpawnScale,
			}
		}

		for i := 0; i < len(dstB.UnknownDW) && i < len(srcB.UnknownDW); i++ {
			dstB.UnknownDW[i] = srcB.UnknownDW[i]
		}
		for i := 0; i < len(dstB.UnknownF32) && i < len(srcB.UnknownF32); i++ {
			dstB.UnknownF32[i] = srcB.UnknownF32[i]
		}
	}

	cpBlock(&src.Source, &e.Source)
	cpBlock(&src.SourceToTarget, &e.SourceToTarget)
	cpBlock(&src.Target, &e.Target)

	return nil
}
