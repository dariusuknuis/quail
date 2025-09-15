package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDirectionalLight is DirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLight struct {
	nameRef  int32
	LightRef int32
	Flags    uint32
	Normal   [3]float32
	Regions  []uint32
}

func (e *WldFragDirectionalLight) FragCode() int {
	return FragCodeDirectionalLight
}

func (e *WldFragDirectionalLight) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)

	enc.Float32(e.Normal[0])
	enc.Float32(e.Normal[1])
	enc.Float32(e.Normal[2])
	enc.Uint32(uint32(len(e.Regions)))
	for _, id := range e.Regions {
		enc.Uint32(id)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDirectionalLight) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.LightRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.Normal[0] = dec.Float32()
	e.Normal[1] = dec.Float32()
	e.Normal[2] = dec.Float32()
	n := dec.Uint32()
	e.Regions = make([]uint32, n)
	for i := uint32(0); i < n; i++ {
		e.Regions[i] = dec.Uint32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDirectionalLight) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragDirectionalLight) SetNameRef(id int32) {
	e.nameRef = id
}
