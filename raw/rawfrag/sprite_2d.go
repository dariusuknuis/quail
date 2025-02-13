package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSprite2D is Sprite2D in libeq, Two-Dimensional Object Reference in openzone, 2DSPRITE (ref) in wld, Fragment07 in lantern
type WldFragSprite2D struct {
	NameRef       int32  `yaml:"name_ref"`
	TwoDSpriteRef int32  `yaml:"two_d_sprite_ref"`
	Flags         uint32 `yaml:"flags"`
}

func (e *WldFragSprite2D) FragCode() int {
	return FragCodeSprite2D
}

func (e *WldFragSprite2D) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.TwoDSpriteRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite2D) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TwoDSpriteRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
