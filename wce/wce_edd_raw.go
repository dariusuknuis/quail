// wce_edd_raw.go
package wce

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/raw"
)

func (w *Wce) ReadEddRaw(src *raw.Edd) error {
	if src == nil {
		return fmt.Errorf("src is nil")
	}

	w.reset()
	w.FileName = src.FileName()
	w.EmitterDefs = w.EmitterDefs[:0]

	folder := strings.TrimSuffix(
		strings.ToLower(filepath.Base(src.FileName())),
		filepath.Ext(src.FileName()),
	)
	if folder == "" {
		folder = "emitters"
	}

	if w.WorldDef == nil {
		w.WorldDef = &WorldDef{folders: []string{folder}}
	}

	for i, rr := range src.Entries {
		if rr == nil {
			return fmt.Errorf("record %d is nil", i)
		}
		def := &EmitterDef{
			folders:  []string{folder},
			TagIndex: i,
		}
		if err := def.FromRaw(w, rr); err != nil {
			return fmt.Errorf("record %d: %w", i, err)
		}
		w.EmitterDefs = append(w.EmitterDefs, def)
	}
	return nil
}

func (w *Wce) WriteEddRaw(dst *raw.Edd) error {
	if dst == nil {
		return fmt.Errorf("dst is nil")
	}

	dst.Version = "110"
	dst.Entries = dst.Entries[:0]

	for i, def := range w.EmitterDefs {
		if def == nil {
			return fmt.Errorf("emitter %d is nil", i)
		}
		if def.TagIndex != i {
			return fmt.Errorf("emitter %d has index %d", i, def.TagIndex)
		}
		rec := &raw.EddEntry{}
		if err := def.ToRaw(w, rec); err != nil {
			return fmt.Errorf("emitter %d: %w", i, err)
		}
		dst.Entries = append(dst.Entries, rec)
	}
	return nil
}