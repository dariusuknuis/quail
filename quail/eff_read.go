package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	err = eff.Read(f)
	if err != nil {
		return fmt.Errorf("read eff: %w", err)
	}
	eff.SetFileName(filepath.Base(path))

	baseName := strings.TrimSuffix(filepath.Base(path), ".eff")
	q.Wld = wce.New(baseName)

	// No WorldDef manipulation here — ReadEffRaw will ensure it exists with the correct folder.
	err = q.Wld.ReadEffRaw(eff)
	if err != nil {
		return fmt.Errorf("convert eff: %w", err)
	}

	return nil
}
