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

	fileName := filepath.Base(path)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	q.Wld = wce.New(baseName)

	switch strings.ToLower(fileName) {
	case "spells.eff":
		eff := &raw.EffOld{}
		err = eff.Read(f)
		if err != nil {
			return fmt.Errorf("read old eff: %w", err)
		}
		eff.SetFileName(fileName)

		err = q.Wld.ReadEffRaw(eff)
		if err != nil {
			return fmt.Errorf("convert old eff: %w", err)
		}
	case "spellsnew.eff":
		eff := &raw.EffNew{}
		err = eff.Read(f)
		if err != nil {
			return fmt.Errorf("read new eff: %w", err)
		}
		eff.SetFileName(fileName)

		err = q.Wld.ReadEffNewRaw(eff)
		if err != nil {
			return fmt.Errorf("convert new eff: %w", err)
		}
	default:
		return fmt.Errorf("unsupported eff file %s: expected spells.eff or spellsnew.eff", fileName)
	}

	return nil
}
