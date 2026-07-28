package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Write writes a complete LIT file.
func (lit *Lit) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	// File identifier
	enc.StringFixed("EQGP", 4)

	// Number of four-byte lighting entries
	enc.Uint32(uint32(len(lit.Entries)))

	for _, entry := range lit.Entries {
		enc.Uint8(entry[0])
		enc.Uint8(entry[1])
		enc.Uint8(entry[2])
		enc.Uint8(entry[3])
	}

	if err := enc.Error(); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
