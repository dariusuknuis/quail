package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSphereList is SphereList in libeq, empty in openzone, SPHERELIST (ref) in wld
type WldFragSphereList struct {
	nameRef          int32
	SphereListDefRef int32
	Flags            uint32
	ScaleFactor      float32
}

func (e *WldFragSphereList) FragCode() int {
	return FragCodeSphereList
}

func (e *WldFragSphereList) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.SphereListDefRef)
	enc.Uint32(e.Flags)
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.ScaleFactor)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSphereList) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.SphereListDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.ScaleFactor = dec.Float32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSphereList) NameRef() int32 {
	return e.nameRef
}
