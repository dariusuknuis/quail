// quail/eff_write.go
package quail

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/raw"
)

func (q *Quail) EffWrite(path string) error {
	if q.Wld == nil {
		return fmt.Errorf("no wld loaded")
	}

	if len(q.Wld.EffectOlds) > 0 && len(q.Wld.EffectNews) > 0 {
		return fmt.Errorf("cannot write eff: both old and new effect definitions are present")
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create eff: %w", err)
	}
	defer f.Close()

	if len(q.Wld.EffectNews) > 0 {
		eff := &raw.EffNew{}
		eff.SetFileName(filepath.Base(path))

		err = q.Wld.WriteEffNewRaw(eff)
		if err != nil {
			return fmt.Errorf("convert new eff: %w", err)
		}

		err = eff.Write(f)
		if err != nil {
			return fmt.Errorf("write new eff: %w", err)
		}

		return nil
	}

	if len(q.Wld.EffectOlds) > 0 {
		eff := &raw.EffOld{}
		eff.SetFileName(filepath.Base(path))

		err = q.Wld.WriteEffRaw(eff)
		if err != nil {
			return fmt.Errorf("convert old eff: %w", err)
		}

		err = eff.Write(f)
		if err != nil {
			return fmt.Errorf("write old eff: %w", err)
		}

		return nil
	}

	return fmt.Errorf("cannot write eff: no old or new effect definitions found")
}
