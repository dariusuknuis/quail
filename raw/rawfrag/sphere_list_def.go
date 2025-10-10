package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSphereListDef is SphereListDef in libeq, empty in openzone, SPHERELISTDEFINITION in wld
type WldFragSphereListDef struct {
	nameRef int32
	Flags   uint32
	Radius  float32
	Scale   float32
	Spheres [][4]float32
}

func (e *WldFragSphereListDef) FragCode() int {
	return FragCodeSphereListDef
}

func (e *WldFragSphereListDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Spheres)))
	enc.Float32(e.Radius)
	if e.Flags&0x01 == 0x01 {
		enc.Float32(e.Scale)
	}
	for _, sphere := range e.Spheres {
		enc.Float32(sphere[0])
		enc.Float32(sphere[1])
		enc.Float32(sphere[2])
		enc.Float32(sphere[3])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSphereListDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	sphereCount := dec.Uint32()
	e.Radius = dec.Float32()
	if e.Flags&0x01 == 0x01 {
		e.Scale = dec.Float32()
	}
	for i := 0; i < int(sphereCount); i++ {
		var sphere [4]float32
		sphere[0] = dec.Float32()
		sphere[1] = dec.Float32()
		sphere[2] = dec.Float32()
		sphere[3] = dec.Float32()
		e.Spheres = append(e.Spheres, sphere)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragSphereListDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragSphereListDef) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
