// wce_eff_old_raw.go
package wce

import (
	"fmt"

	"github.com/xackery/quail/raw"
)

// ReadEffRaw converts a parsed raw.EffOld into wce.EffectOlds for ASCII output.
func (w *Wce) ReadEffRaw(src *raw.EffOld) error {
	if src == nil {
		return fmt.Errorf("src is nil")
	}

	w.reset()
	w.FileName = src.FileName()
	w.EffectOlds = w.EffectOlds[:0]

	// ensure WorldDef exists so WriteAscii doesn't bail out
	if w.WorldDef == nil {
		w.WorldDef = &WorldDef{folders: []string{"spells"}}
	}

	for i, rr := range src.Records {
		if rr == nil {
			return fmt.Errorf("record %d is nil", i)
		}
		def := &EffectOld{
			folders:  []string{"spells"},
			TagIndex: i,
		}
		// Your EffectOld.FromRaw should accept *raw.EffOldRecord
		if err := def.FromRaw(w, rr); err != nil {
			return fmt.Errorf("record %d: %w", i, err)
		}
		w.EffectOlds = append(w.EffectOlds, def)
	}
	return nil
}

func (w *Wce) WriteEffRaw(dst *raw.EffOld) error {
	if dst == nil {
		return fmt.Errorf("dst is nil")
	}
	dst.Records = dst.Records[:0]

	for i, def := range w.EffectOlds {
		if def == nil {
			return fmt.Errorf("effect %d is nil", i)
		}
		rec := &raw.EffOldRecord{}
		if err := def.ToRaw(w, rec); err != nil {
			return fmt.Errorf("effect %d: %w", i, err)
		}
		dst.Records = append(dst.Records, rec)
	}
	return nil
}
