package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragUserData is empty in libeq, empty in openzone, USERDATA in wld
type WldFragUserData struct {
	Length uint32
	Data   string
}

func (e *WldFragUserData) FragCode() int {
	return FragCodeUserData
}

func (e *WldFragUserData) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	start := enc.Pos()
	encodedStr := helper.WriteStringHash(e.Data + "\x00")
	enc.Uint32(uint32(len(encodedStr)))
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

func (e *WldFragUserData) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.Length = dec.Uint32()
	decodedStr := helper.ReadStringHash((dec.Bytes(int(e.Length))))
	e.Data = strings.TrimRight(decodedStr, "\x00")
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragUserData) NameRef() int32 {
	return 0
}
