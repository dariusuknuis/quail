// quail/eff_read.go
package quail

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/raw"
)

func (q *Quail) EffRead(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open eff: %w", err)
	}
	defer f.Close()

	reader, err := raw.Read(".eff", f)
	if err != nil {
		return fmt.Errorf("eff read: %w", err)
	}
	reader.SetFileName(filepath.Base(path))

	if err := q.RawRead(reader); err != nil {
		return fmt.Errorf("q rawRead eff: %w", err)
	}
	return nil
}
