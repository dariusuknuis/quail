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
	Header         [8]byte
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
	ColorBGRA     uint32
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
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		// de-dupe per Tag (same pattern used elsewhere)
		if token.TagIsWritten(tag) {
			continue
		}
		token.TagSetIsWritten(tag)

		fmt.Fprintf(w, "%s \"%s\"\n", effectOldDefinition(), tag)
		fmt.Fprintf(w, "\tNUMRECORDS %d\n", len(records))

		for i, rec := range records {
			fmt.Fprintf(w, "\t\tRECORD // %d\n", i)
			// Header as 8 bytes (decimal)
			fmt.Fprintf(w, "\t\t\tHEADER %d %d %d %d %d %d %d %d\n",
				rec.Header[0], rec.Header[1], rec.Header[2], rec.Header[3],
				rec.Header[4], rec.Header[5], rec.Header[6], rec.Header[7])

			writeBlock := func(name string, b *EffectOldBlock) {
				fmt.Fprintf(w, "\t\t\tBLOCK \"%s\"\n", name)
				fmt.Fprintf(w, "\t\t\t\tLABEL \"%s\"\n", b.Label)

				// Emit declared count (we’ll just print all 12; empty ones are "")
				fmt.Fprintf(w, "\t\t\t\tEXTRASPRITES %d\n", len(b.ExtraSprites))
				for _, s := range b.ExtraSprites {
					fmt.Fprintf(w, "\t\t\t\t\tSPRITE \"%s\"\n", s)
				}

				fmt.Fprintf(w, "\t\t\t\tUNKNOWNPARAM %d\n", b.UnknownParam)
				fmt.Fprintf(w, "\t\t\t\tSOUNDREF %d\n", b.SoundRef)

				for j := 0; j < 3; j++ {
					sub := b.Sub[j]
					fmt.Fprintf(w, "\t\t\t\tSUB // %d\n", j)
					fmt.Fprintf(w, "\t\t\t\t\tPRIMARYSPRITE \"%s\"\n", sub.PrimarySprite)
					fmt.Fprintf(w, "\t\t\t\t\tDAGINDEX %d\n", sub.DagIndex)
					fmt.Fprintf(w, "\t\t\t\t\tEFFECTTYPE %d\n", sub.EffectType)
					fmt.Fprintf(w, "\t\t\t\t\tCOLORBGRA %d\n", sub.ColorBGRA)
					fmt.Fprintf(w, "\t\t\t\t\tGRAVITY %0.8e\n", sub.Gravity)
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNNORMAL %0.8e %0.8e %0.8e\n",
						sub.SpawnNormal[0], sub.SpawnNormal[1], sub.SpawnNormal[2])
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNRADIUS %0.8e\n", sub.SpawnRadius)
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNANGLE %0.8e\n", sub.SpawnAngle)
					fmt.Fprintf(w, "\t\t\t\t\tLIFESPAN %d\n", sub.Lifespan)
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNVELOCITY %0.8e\n", sub.SpawnVelocity)
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNRATE %d\n", sub.SpawnRate)
					fmt.Fprintf(w, "\t\t\t\t\tSPAWNSCALE %0.8e\n", sub.SpawnScale)
				}

				// --- Unknown DWORDs ---
				for i, v := range b.UnknownDW {
					fmt.Fprintf(w, "\t\t\t\tUNKNOWNDW_%d %d\n", i+1, v)
				}

				// --- Unknown Float32s ---
				for i, v := range b.UnknownF32 {
					fmt.Fprintf(w, "\t\t\t\tUNKNOWNF32_%d %0.8e\n", i+1, v)
				}
			}

			writeBlock("Source", &rec.Source)
			writeBlock("SourceToTarget", &rec.SourceToTarget)
			writeBlock("Target", &rec.Target)
		}

		fmt.Fprintln(w)
	}
	return nil
}

// ReadEffOld parses the next EFFECTOLDDEF block and returns its records.
// It does not clear or touch any Wce; you can append these wherever you like.
func ReadEffOld(token *AsciiReadToken) ([]*EffectOld, error) {
	// Some of your other readers assume the header (like EQG*DEF "Tag") is already consumed
	// by the caller. Here we only rely on the body properties.
	records, err := token.ReadProperty("NUMRECORDS", 1)
	if err != nil {
		return nil, err
	}
	n := 0
	if err := parse(&n, records[1]); err != nil {
		return nil, fmt.Errorf("numrecords: %w", err)
	}
	out := make([]*EffectOld, 0, n)

	for i := 0; i < n; i++ {
		if _, err := token.ReadProperty("RECORD", 0); err != nil {
			return nil, fmt.Errorf("record %d: %w", i, err)
		}

		rec := &EffectOld{}

		hdr, err := token.ReadProperty("HEADER", 8)
		if err != nil {
			return nil, fmt.Errorf("record %d header: %w", i, err)
		}
		for k := 0; k < 8; k++ {
			if err := parse(&rec.Header[k], hdr[k+1]); err != nil {
				return nil, fmt.Errorf("record %d header[%d]: %w", i, k, err)
			}
		}

		readBlock := func(which string, dst *EffectOldBlock) error {
			props, err := token.ReadProperty("BLOCK", 1)
			if err != nil {
				return fmt.Errorf("%s block open: %w", which, err)
			}
			got := strings.Trim(props[1], "\"")
			if !strings.EqualFold(got, which) {
				return fmt.Errorf("expected BLOCK %q, got %q", which, got)
			}

			props, err = token.ReadProperty("LABEL", 1)
			if err != nil {
				return fmt.Errorf("%s label: %w", which, err)
			}
			dst.Label = props[1]

			props, err = token.ReadProperty("EXTRASPRITES", 1)
			if err != nil {
				return fmt.Errorf("%s extrasprites: %w", which, err)
			}
			cnt := 0
			if err := parse(&cnt, props[1]); err != nil {
				return fmt.Errorf("%s extrasprites cnt: %w", which, err)
			}
			for j := 0; j < cnt && j < len(dst.ExtraSprites); j++ {
				props, err = token.ReadProperty("SPRITE", 1)
				if err != nil {
					return fmt.Errorf("%s sprite %d: %w", which, j, err)
				}
				dst.ExtraSprites[j] = props[1]
			}

			if props, err = token.ReadProperty("UNKNOWNPARAM", 1); err != nil {
				return fmt.Errorf("%s unknownparam: %w", which, err)
			}
			if err := parse(&dst.UnknownParam, props[1]); err != nil {
				return fmt.Errorf("%s unknownparam: %w", which, err)
			}

			if props, err = token.ReadProperty("SOUNDREF", 1); err != nil {
				return fmt.Errorf("%s soundref: %w", which, err)
			}
			if err := parse(&dst.SoundRef, props[1]); err != nil {
				return fmt.Errorf("%s soundref: %w", which, err)
			}

			for s := 0; s < 3; s++ {
				if _, err := token.ReadProperty("SUB", 0); err != nil {
					return fmt.Errorf("%s sub %d: %w", which, s, err)
				}
				sub := &dst.Sub[s]

				if props, err = token.ReadProperty("PRIMARYSPRITE", 1); err != nil {
					return fmt.Errorf("%s sub %d primarysprite: %w", which, s, err)
				}
				sub.PrimarySprite = props[1]

				if props, err = token.ReadProperty("DAGINDEX", 1); err != nil {
					return fmt.Errorf("%s sub %d dagindex: %w", which, s, err)
				}
				if err := parse(&sub.DagIndex, props[1]); err != nil {
					return fmt.Errorf("%s sub %d dagindex: %w", which, s, err)
				}

				if props, err = token.ReadProperty("EFFECTTYPE", 1); err != nil {
					return fmt.Errorf("%s sub %d effecttype: %w", which, s, err)
				}
				if err := parse(&sub.EffectType, props[1]); err != nil {
					return fmt.Errorf("%s sub %d effecttype: %w", which, s, err)
				}

				if props, err = token.ReadProperty("COLORBGRA", 1); err != nil {
					return fmt.Errorf("%s sub %d colorbgra: %w", which, s, err)
				}
				if err := parse(&sub.ColorBGRA, props[1]); err != nil {
					return fmt.Errorf("%s sub %d colorbgra: %w", which, s, err)
				}

				if props, err = token.ReadProperty("GRAVITY", 1); err != nil {
					return fmt.Errorf("%s sub %d gravity: %w", which, s, err)
				}
				if err := parse(&sub.Gravity, props[1]); err != nil {
					return fmt.Errorf("%s sub %d gravity: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNNORMAL", 3); err != nil {
					return fmt.Errorf("%s sub %d spawnnormal: %w", which, s, err)
				}
				if err := parse(&sub.SpawnNormal, props[1:]...); err != nil {
					return fmt.Errorf("%s sub %d spawnnormal: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNRADIUS", 1); err != nil {
					return fmt.Errorf("%s sub %d spawnradius: %w", which, s, err)
				}
				if err := parse(&sub.SpawnRadius, props[1]); err != nil {
					return fmt.Errorf("%s sub %d spawnradius: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNANGLE", 1); err != nil {
					return fmt.Errorf("%s sub %d spawnangle: %w", which, s, err)
				}
				if err := parse(&sub.SpawnAngle, props[1]); err != nil {
					return fmt.Errorf("%s sub %d spawnangle: %w", which, s, err)
				}

				if props, err = token.ReadProperty("LIFESPAN", 1); err != nil {
					return fmt.Errorf("%s sub %d lifespan: %w", which, s, err)
				}
				if err := parse(&sub.Lifespan, props[1]); err != nil {
					return fmt.Errorf("%s sub %d lifespan: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNVELOCITY", 1); err != nil {
					return fmt.Errorf("%s sub %d spawnvelocity: %w", which, s, err)
				}
				if err := parse(&sub.SpawnVelocity, props[1]); err != nil {
					return fmt.Errorf("%s sub %d spawnvelocity: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNRATE", 1); err != nil {
					return fmt.Errorf("%s sub %d spawnrate: %w", which, s, err)
				}
				if err := parse(&sub.SpawnRate, props[1]); err != nil {
					return fmt.Errorf("%s sub %d spawnrate: %w", which, s, err)
				}

				if props, err = token.ReadProperty("SPAWNSCALE", 1); err != nil {
					return fmt.Errorf("%s sub %d spawnscale: %w", which, s, err)
				}
				if err := parse(&sub.SpawnScale, props[1]); err != nil {
					return fmt.Errorf("%s sub %d spawnscale: %w", which, s, err)
				}
			}

			for i := 0; i < 51; i++ {
				key := fmt.Sprintf("UNKNOWNDW_%d", i+1)
				props, err := token.ReadProperty(key, 1)
				if err != nil {
					return fmt.Errorf("%s dw %d: %w", which, i, err)
				}
				if err := parse(&dst.UnknownDW[i], props[1]); err != nil {
					return fmt.Errorf("%s dw %d: %w", which, i, err)
				}
			}

			for i := 0; i < 12; i++ {
				key := fmt.Sprintf("UNKNOWNF32_%d", i+1)
				props, err := token.ReadProperty(key, 1)
				if err != nil {
					return fmt.Errorf("%s f32 %d: %w", which, i, err)
				}
				if err := parse(&dst.UnknownF32[i], props[1]); err != nil {
					return fmt.Errorf("%s f32 %d: %w", which, i, err)
				}
			}

			return nil
		}

		if err := readBlock("Source", &rec.Source); err != nil {
			return nil, fmt.Errorf("record %d: %w", i, err)
		}
		if err := readBlock("SourceToTarget", &rec.SourceToTarget); err != nil {
			return nil, fmt.Errorf("record %d: %w", i, err)
		}
		if err := readBlock("Target", &rec.Target); err != nil {
			return nil, fmt.Errorf("record %d: %w", i, err)
		}

		out = append(out, rec)
	}

	return out, nil
}

func (e *EffectOld) ToRaw(_ *Wce, dst *raw.EffOld) error {
	if e == nil || dst == nil {
		return fmt.Errorf("nil receiver or destination")
	}

	dst.Records = dst.Records[:0]

	for i, rec := range e.Records {
		if rec == nil {
			return fmt.Errorf("record %d is nil", i)
		}

		rRec := &raw.EffOld{}
		// Header (8 bytes)
		rRec.Header = rec.Header

		// Blocks
		copyBlock := func(src *EffectOldBlock) (*raw.EffOldBlock, error) {
			if src == nil {
				return nil, fmt.Errorf("nil block")
			}
			dstB := &raw.EffOldBlock{
				Label:        src.Label,
				UnknownParam: src.UnknownParam,
				SoundRef:     src.SoundRef,
			}

			// ExtraSprites (fixed 12)
			for i := 0; i < len(dstB.ExtraSprites) && i < len(src.ExtraSprites); i++ {
				dstB.ExtraSprites[i] = src.ExtraSprites[i]
			}

			// Sub (fixed 3)
			for i := 0; i < len(dstB.Sub) && i < len(src.Sub); i++ {
				s := src.Sub[i]
				dstB.Sub[i] = raw.EffSubEffect{
					PrimarySprite: s.PrimarySprite,
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

			// UnknownDW (fixed 51)
			for i := 0; i < len(dstB.UnknownDW) && i < len(src.UnknownDW); i++ {
				dstB.UnknownDW[i] = src.UnknownDW[i]
			}

			// UnknownF32 (fixed 12)
			for i := 0; i < len(dstB.UnknownF32) && i < len(src.UnknownF32); i++ {
				dstB.UnknownF32[i] = src.UnknownF32[i]
			}

			return dstB, nil
		}

		var err error
		if rRec.Source, err = copyBlock(&rec.Source); err != nil {
			return fmt.Errorf("record %d source: %w", i, err)
		}
		if rRec.SourceToTarget, err = copyBlock(&rec.SourceToTarget); err != nil {
			return fmt.Errorf("record %d sourceToTarget: %w", i, err)
		}
		if rRec.Target, err = copyBlock(&rec.Target); err != nil {
			return fmt.Errorf("record %d target: %w", i, err)
		}

		dst.Records = append(dst.Records, rRec)
	}

	return nil
}

// FromRaw fills the ASCII/in-memory EffOldDef from raw.EffOld.
func (e *EffectOld) FromRaw(_ *Wce, src *raw.EffOld) error {
	if e == nil || src == nil {
		return fmt.Errorf("nil receiver or source")
	}

	e.Records = e.Records[:0]

	for i, rec := range src.Records {
		if rec == nil {
			return fmt.Errorf("raw record %d is nil", i)
		}

		dstRec := &EffectOld{
			Header: rec.Header,
		}

		cpBlock := func(srcB *raw.EffOldBlock, dstB *EffectOldBlock) error {
			if srcB == nil {
				return fmt.Errorf("nil raw block")
			}
			dstB.Label = srcB.Label
			dstB.UnknownParam = srcB.UnknownParam
			dstB.SoundRef = srcB.SoundRef

			for i := 0; i < len(dstB.ExtraSprites) && i < len(srcB.ExtraSprites); i++ {
				dstB.ExtraSprites[i] = srcB.ExtraSprites[i]
			}
			for i := 0; i < len(dstB.Sub) && i < len(srcB.Sub); i++ {
				s := srcB.Sub[i]
				dstB.Sub[i] = EffectOldSub{
					PrimarySprite: s.PrimarySprite,
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
			return nil
		}

		if err := cpBlock(rec.Source, &dstRec.Source); err != nil {
			return fmt.Errorf("raw record %d source: %w", i, err)
		}
		if err := cpBlock(rec.SourceToTarget, &dstRec.SourceToTarget); err != nil {
			return fmt.Errorf("raw record %d sourceToTarget: %w", i, err)
		}
		if err := cpBlock(rec.Target, &dstRec.Target); err != nil {
			return fmt.Errorf("raw record %d target: %w", i, err)
		}

		e.Records = append(e.Records, dstRec)
	}

	return nil
}

// -----------------------------------------
// Optional helpers for Wce.EffectOlds table
// -----------------------------------------
//
// If you prefer to convert directly to/from Wce.EffectOlds instead of via
// EffOldDef.Records, these helpers keep the same mapping logic.

func (w *Wce) EffOldsToRaw(dst *raw.EffOld) error {
	if w == nil || dst == nil {
		return fmt.Errorf("nil receiver or destination")
	}
	def := &EffectOld{Records: w.EffectOlds}
	return def.ToRaw(w, dst)
}

func (w *Wce) EffOldsFromRaw(src *raw.EffOld) error {
	if w == nil || src == nil {
		return fmt.Errorf("nil receiver or source")
	}
	def := &EffectOld{}
	if err := def.FromRaw(w, src); err != nil {
		return err
	}
	w.EffectOlds = def.Records
	return nil
}

// ===============================
// Optional Wce helpers
// ===============================

// Write all effects currently on w.EffectOlds.
func (w *Wce) WriteEffOld(token *AsciiWriteToken, folders []string, tag string) error {
	return WriteEffOld(token, folders, tag, w.EffectOlds)
}

// Read a block and append to w.EffectOlds.
func (w *Wce) ReadEffOld(token *AsciiReadToken) error {
	recs, err := ReadEffOld(token)
	if err != nil {
		return err
	}
	w.EffectOlds = append(w.EffectOlds, recs...)
	return nil
}
