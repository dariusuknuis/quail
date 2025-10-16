// quail/eff_read.go
package quail

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

func (q *Quail) EffRead(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open eff: %w", err)
	}
	defer f.Close()

	eff := &raw.EffOld{}
	if err := eff.Read(f); err != nil {
		return fmt.Errorf("read eff: %w", err)
	}
	eff.SetFileName(filepath.Base(path))

	w := wce.New(eff.FileName())
	if err := w.ReadEffRaw(eff); err != nil {
		return fmt.Errorf("convert eff: %w", err)
	}

	q.Wld = w // or append to q.Wce if you want
	return nil
}
