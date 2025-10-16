// wce_eff_old_raw.go
package wce

import (
	"fmt"

	"github.com/xackery/quail/raw"
)

// ReadEffRaw converts a parsed raw.EffOld into wce.EffectOlds for ASCII output.
func (wce *Wce) ReadEffRaw(src *raw.EffOld) error {
	if src == nil {
		return fmt.Errorf("src is nil")
	}

	wce.reset()
	wce.FileName = src.FileName()
	wce.EffectOlds = wce.EffectOlds[:0]

	for i, rr := range src.Records {
		if rr == nil {
			return fmt.Errorf("record %d is nil", i)
		}
		def := &EffectOld{
			folders:  []string{"spells"}, // ensures it writes to spells/spells.wce
			TagIndex: i,
		}
		err := def.FromRaw(wce, rr)
		if err != nil { // NOTE: FromRaw must accept *raw.EffOldRecord
			return fmt.Errorf("record %d: %w", i, err)
		}
		wce.EffectOlds = append(wce.EffectOlds, def)
	}
	return nil
}
