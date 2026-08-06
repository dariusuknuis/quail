package quail

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/raw"
)

func (q *Quail) EddWrite(path string) error {
	if q.Wld == nil {
		return fmt.Errorf("no wce loaded")
	}

	edd := &raw.Edd{}
	edd.SetFileName(filepath.Base(path))

	err := q.Wld.WriteEddRaw(edd)
	if err != nil {
		return fmt.Errorf("convert edd: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create edd: %w", err)
	}
	defer f.Close()

	err = edd.Write(f)
	if err != nil {
		return fmt.Errorf("write edd: %w", err)
	}

	return nil
}
