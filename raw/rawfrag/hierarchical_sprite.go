package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragHierarchicalSprite is HierarchicalSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHICALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragHierarchicalSprite struct {
	NameRef               int32
	HierarchicalSpriteRef int32
	Param                 uint32
}

func (e *WldFragHierarchicalSprite) FragCode() int {
	return FragCodeHierarchicalSprite
}

func (e *WldFragHierarchicalSprite) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.HierarchicalSpriteRef)
	enc.Uint32(e.Param)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchicalSprite) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.HierarchicalSpriteRef = dec.Int32()
	e.Param = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
