// quail/eff_write.go
package quail

import (
	"fmt"
	"os"

	"github.com/xackery/quail/raw"
)

func (q *Quail) EffWrite(path string) error {
	if q.Wld == nil {
		return fmt.Errorf("no wld found")
	}
	if len(q.Wld.EffectOlds) == 0 {
		return fmt.Errorf("no effects to write")
	}

	dst := &raw.EffOld{}
	dst.SetFileName(q.Wld.FileName)

	if err := q.Wld.WriteEffRaw(dst); err != nil {
		return fmt.Errorf("write eff raw: %w", err)
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create eff: %w", err)
	}
	defer w.Close()

	if err := dst.Write(w); err != nil {
		return fmt.Errorf("encode eff: %w", err)
	}

	fmt.Printf("Wrote %s (%d records)\n", path, len(dst.Records))
	return nil
}
