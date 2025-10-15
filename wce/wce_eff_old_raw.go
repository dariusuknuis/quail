// wce_eff_old_raw.go
package wce

import (
	"fmt"

	"github.com/xackery/quail/raw"
)

// NOTE: This file intentionally ONLY contains raw <-> in-memory conversions
// for the "old" EFF format. ASCII read/write stays in wce_eff_old.go (as requested).

// -------------------------
// EffOldDef <-> raw.EffOld
// -------------------------

// ToRaw copies the ASCII/in-memory EffOldDef records into raw.EffOld.
// Patterned similar to your WLD ToRaw, but this format isn't fragment-based.
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
