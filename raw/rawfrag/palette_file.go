package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragDefaultPaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type WldFragDefaultPaletteFile struct {
	NameLength uint16
	FileName   string
}

func (e *WldFragDefaultPaletteFile) FragCode() int {
	return FragCodeDefaultPaletteFile
}

func (e *WldFragDefaultPaletteFile) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	start := enc.Pos()
	encodedStr := helper.WriteStringHash(e.FileName + "\x00")
	enc.Uint16(uint16(len(encodedStr)))
	enc.String(string(encodedStr))

	diff := enc.Pos() - start
	paddingSize := (4 - diff%4) % 4
	enc.Bytes(make([]byte, paddingSize))

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameLength = dec.Uint16()
	decodedStr := helper.ReadStringHash((dec.Bytes(int(e.NameLength))))
	e.FileName = strings.TrimRight(decodedStr, "\x00")
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) NameRef() int32 {
	return 0
}
