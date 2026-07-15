package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Writes an EQAL (.anl) file.
func (anl *Anl) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.StringFixed("EQAL", 4)
	enc.Uint32(anl.Version)

	// Build string block and record offsets.
	offsets := make([]uint32, len(anl.Animations))

	var stringBlock []byte

	for i, name := range anl.Animations {
		offsets[i] = uint32(len(stringBlock))
		stringBlock = append(stringBlock, []byte(name)...)
		stringBlock = append(stringBlock, 0)
	}

	defaultOffset := uint32(len(stringBlock))
	stringBlock = append(stringBlock, []byte(anl.DefaultAnimation)...)
	stringBlock = append(stringBlock, 0)

	// Header
	enc.Uint32(uint32(len(stringBlock)))
	enc.Uint32(uint32(len(anl.Animations)))

	// String block
	enc.Bytes(stringBlock)

	// Animation offsets
	for _, off := range offsets {
		enc.Uint32(off)
	}

	// Default animation offset
	enc.Uint32(defaultOffset)

	if err := enc.Error(); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}