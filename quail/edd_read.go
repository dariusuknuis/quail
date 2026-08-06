package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

func (q *Quail) EddRead(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open edd: %w", err)
	}
	defer f.Close()

	edd := &raw.Edd{}
	err = edd.Read(f)
	if err != nil {
		return fmt.Errorf("read edd: %w", err)
	}
	edd.SetFileName(filepath.Base(path))

	baseName := strings.TrimSuffix(
		filepath.Base(path),
		filepath.Ext(path),
	)
	q.Wld = wce.New(baseName)

	err = q.Wld.ReadEddRaw(edd)
	if err != nil {
		return fmt.Errorf("convert edd: %w", err)
	}

	return nil
}