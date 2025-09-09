package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDmRGBTrackDef is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type WldFragDmRGBTrackDef struct {
	nameRef    int32
	Flags      uint32 // usually contains 1
	RGBAFrames [][][4]uint8
	Sleep      uint32 // usually contains 200
	Data4      uint32 // usually contains 0
	//RGBAs      [][4]uint8
}

func (e *WldFragDmRGBTrackDef) FragCode() int {
	return FragCodeDmRGBTrackDef
}

func (e *WldFragDmRGBTrackDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.Int32(e.nameRef)
	if e.Flags == 0 {
		e.Flags = 1
	}
	enc.Uint32(e.Flags)
	if len(e.RGBAFrames) < 1 {
		return fmt.Errorf("no frames found")
	}
	enc.Uint32(uint32(len(e.RGBAFrames[0])))
	enc.Uint32(uint32(len(e.RGBAFrames)))
	enc.Uint32(e.Sleep)
	enc.Uint32(e.Data4)

	for _, frame := range e.RGBAFrames {
		for _, rgba := range frame {
			enc.Uint8(rgba[0])
			enc.Uint8(rgba[1])
			enc.Uint8(rgba[2])
			enc.Uint8(rgba[3])
		}
	}

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrackDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	numRGBA := dec.Uint32()
	frameCount := dec.Uint32()
	e.Sleep = dec.Uint32()
	e.Data4 = dec.Uint32()
	if e.Data4 != 0 {
		fmt.Printf("Data4 (NumVertices) on rgbtrack is not 0 (%d), tell xack you found this!\n", e.Data4)
	}
	e.RGBAFrames = make([][][4]uint8, frameCount)
	for i := range e.RGBAFrames {
		e.RGBAFrames[i] = make([][4]uint8, numRGBA)
		for j := range e.RGBAFrames[i] {
			e.RGBAFrames[i][j][0] = dec.Uint8()
			e.RGBAFrames[i][j][1] = dec.Uint8()
			e.RGBAFrames[i][j][2] = dec.Uint8()
			e.RGBAFrames[i][j][3] = dec.Uint8()
		}
	}
	return nil
}

func (e *WldFragDmRGBTrackDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragDmRGBTrackDef) SetNameRef(id int32) {
	e.nameRef = id
}
