package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Lit struct {
	MetaFileName string
	Entries      [][4]uint8
}

// Identity returns the type of the struct.
func (lit *Lit) Identity() string {
	return "lit"
}

// Read reads a LIT file.
func (lit *Lit) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGP" {
		return fmt.Errorf("invalid header %s, wanted EQGP", header)
	}

	lightCount := dec.Uint32()

	lit.Entries = make([][4]uint8, 0, lightCount)

	for i := uint32(0); i < lightCount; i++ {
		entry := [4]uint8{
			dec.Uint8(),
			dec.Uint8(),
			dec.Uint8(),
			dec.Uint8(),
		}

		lit.Entries = append(lit.Entries, entry)
	}

	pos := dec.Pos()

	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}

	if pos < endPos {
		remaining := dec.Bytes(int(endPos - pos))

		// Some files may contain one trailing zero DWORD.
		if !bytes.Equal(remaining, []byte{0x00, 0x00, 0x00, 0x00}) {
			fmt.Printf("remaining bytes: %s\n", hex.Dump(remaining))

			return fmt.Errorf(
				"%d bytes remaining (%d total)",
				endPos-pos,
				endPos,
			)
		}
	}

	if pos > endPos {
		return fmt.Errorf("read past end of file")
	}

	if err := dec.Error(); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	return nil
}

// SetFileName sets the name of the file.
func (lit *Lit) SetFileName(name string) {
	lit.MetaFileName = name
}

// FileName returns the name of the file.
func (lit *Lit) FileName() string {
	return lit.MetaFileName
}
