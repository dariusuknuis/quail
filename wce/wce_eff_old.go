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
	Label       string
	ExtraEffect [12]ExtraEffect
	EffectMode  int32
	SoundRef    int32
	Sub         [3]EffectOldSub
	UnknownDW   [51]uint32
	UnknownF32  [12]float32
}

type EffectOldSub struct {
	PrimarySprite string
	DagIndex      int32
	EffectType    int32
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

type ExtraEffect struct {
	Sprite              string
	ColorBGR            [3]uint8
	AnimSpeedMultiplier float32
	AngleRangeA         int16
	AngleRangeB         int16
	Radius              float32
	EffectType          int16
	Scale               float32
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
			fmt.Fprintf(w, "\t\tEFFECTMODE %d\n", b.EffectMode)
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

			for j := 0; j < 12; j++ {
				ee := b.ExtraEffect[j]
				fmt.Fprintf(w, "\t\tEXTRAEFFECT // %d\n", j)
				fmt.Fprintf(w, "\t\t\tSPRITE \"%s\"\n", ee.Sprite)
				fmt.Fprintf(w, "\t\t\tRGB %d %d %d \n", ee.ColorBGR[2], ee.ColorBGR[1], ee.ColorBGR[0])
				fmt.Fprintf(w, "\t\t\tANIMSPEEDMULTIPLIER %0.8e\n", ee.AnimSpeedMultiplier)
				fmt.Fprintf(w, "\t\t\tANGLERANGEA %d \n", ee.AngleRangeA)
				fmt.Fprintf(w, "\t\t\tANGLERANGEB %d \n", ee.AngleRangeB)
				fmt.Fprintf(w, "\t\t\tRADIUS %0.8e\n", ee.Radius)
				fmt.Fprintf(w, "\t\t\tEFFECTTYPE %d\n", ee.EffectType)
				fmt.Fprintf(w, "\t\t\tSCALE %0.8e\n", ee.Scale)
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

		records, err = token.ReadProperty("EFFECTMODE", 1)
		if err != nil {
			return err
		}
		err = parse(&dst.EffectMode, records[1])
		if err != nil {
			return fmt.Errorf("%s effect mode: %w", which, err)
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
			_, err := token.ReadProperty("SUBEFFECT", 0)
			if err != nil {
				return fmt.Errorf("%s sub %d: %w", which, s, err)
			}
			sub := &dst.Sub[s]

			records, err = token.ReadProperty("BLITSPRITE", 1)
			if err != nil {
				return fmt.Errorf("%s sub %d blit sprite: %w", which, s, err)
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

		for s := 0; s < 12; s++ {
			_, err := token.ReadProperty("EXTRAEFFECT", 0)
			if err != nil {
				return fmt.Errorf("%s extra effect %d: %w", which, s, err)
			}
			ee := &dst.ExtraEffect[s]

			records, err = token.ReadProperty("SPRITE", 1)
			if err != nil {
				return fmt.Errorf("%s extra effect %d sprite: %w", which, s, err)
			}
			ee.Sprite = records[1]

			records, err = token.ReadProperty("RGB", 3)
			if err != nil {
				return err
			}
			err = parse(&ee.ColorBGR, records[1:]...)
			if err != nil {
				return fmt.Errorf("%s extra effect %d rgb: %w", which, s, err)
			}

			records, err = token.ReadProperty("ANIMSPEEDMULTIPLIER", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.AnimSpeedMultiplier, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d animation speed multiplier: %w", which, s, err)
			}

			records, err = token.ReadProperty("ANGLERANGEA", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.AngleRangeA, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d angle range A: %w", which, s, err)
			}

			records, err = token.ReadProperty("ANGLERANGEB", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.AngleRangeB, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d angle range B: %w", which, s, err)
			}

			records, err = token.ReadProperty("RADIUS", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.Radius, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d radius: %w", which, s, err)
			}

			records, err = token.ReadProperty("EFFECTTYPE", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.EffectType, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d effecttype: %w", which, s, err)
			}

			records, err = token.ReadProperty("SCALE", 1)
			if err != nil {
				return err
			}
			err = parse(&ee.Scale, records[1])
			if err != nil {
				return fmt.Errorf("%s extra effect %d scale: %w", which, s, err)
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
		dstB.EffectMode = src.EffectMode
		dstB.SoundRef = src.SoundRef

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

		for i := 0; i < len(dstB.ExtraEffect) && i < len(src.ExtraEffect); i++ {
			dstB.ExtraEffect[i].Blit = src.ExtraEffect[i].Sprite
			dstB.ExtraEffect[i].ColorBGR = src.ExtraEffect[i].ColorBGR
			dstB.ExtraEffect[i].AnimSpeedMultiplier = src.ExtraEffect[i].AnimSpeedMultiplier
			dstB.ExtraEffect[i].AngleRangeA = src.ExtraEffect[i].AngleRangeA
			dstB.ExtraEffect[i].AngleRangeB = src.ExtraEffect[i].AngleRangeB
			dstB.ExtraEffect[i].Radius = src.ExtraEffect[i].Radius
			dstB.ExtraEffect[i].EffectType = src.ExtraEffect[i].EffectType
			dstB.ExtraEffect[i].Scale = src.ExtraEffect[i].Scale
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
		dstB.EffectMode = srcB.EffectMode
		dstB.SoundRef = srcB.SoundRef

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

		for i := 0; i < 12; i++ {
			s := srcB.ExtraEffect[i]
			dstB.ExtraEffect[i] = ExtraEffect{
				Sprite:              s.Blit,
				ColorBGR:            s.ColorBGR,
				AnimSpeedMultiplier: s.AnimSpeedMultiplier,
				AngleRangeA:         s.AngleRangeA,
				AngleRangeB:         s.AngleRangeB,
				Radius:              s.Radius,
				EffectType:          s.EffectType,
				Scale:               s.Scale,
			}
		}
	}

	cpBlock(&src.Source, &e.Source)
	cpBlock(&src.SourceToTarget, &e.SourceToTarget)
	cpBlock(&src.Target, &e.Target)

	return nil
}
