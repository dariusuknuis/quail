package raw

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
)

// Write writes a PRT file.
func (prt *Prt) Write(w io.Writer) error {
	if prt.Version < 1 || prt.Version > 5 {
		return fmt.Errorf(
			"unsupported PRT version %d, wanted version 1 through 5",
			prt.Version,
		)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.String("PTCL")
	enc.Uint32(uint32(len(prt.Entries)))
	enc.Uint32(prt.Version)

	for i, entry := range prt.Entries {
		if entry == nil {
			return fmt.Errorf("particle %d is nil", i)
		}

		// -----------------------------------------------------
		// Version 1+
		// -----------------------------------------------------

		enc.Int32(entry.EmitterID)

		// ParticlePoint is stored as char[64]. The name may use
		// at most 63 bytes because it requires a null terminator.
		if len(entry.ParticlePoint) > 63 {
			return fmt.Errorf(
				"particle %d: particle-point name %q is %d bytes; "+
					"maximum is 63 bytes",
				i,
				entry.ParticlePoint,
				len(entry.ParticlePoint),
			)
		}

		if strings.IndexByte(entry.ParticlePoint, 0) >= 0 {
			return fmt.Errorf(
				"particle %d: particle-point name contains a null byte",
				i,
			)
		}

		pointData := make([]byte, 64)
		copy(pointData, entry.ParticlePoint)
		enc.Bytes(pointData)

		enc.Int32(entry.ParticleType)
		enc.Int32(entry.AnimNumber)
		enc.Int32(entry.AnimVariation)
		enc.Int32(entry.RandomAnim)
		enc.Int32(entry.StartTime)
		enc.Int32(entry.Lifespan)

		// Fields added by later versions
		if prt.Version >= 2 {
			enc.Int32(entry.Ground)
		}

		if prt.Version >= 3 {
			enc.Int32(entry.PlayWithMat)
		}

		if prt.Version >= 4 {
			enc.Int32(entry.Sporadic)
		}

		// Version 5 adds ColdEmitterID at the end of the record.
		if prt.Version >= 5 {
			enc.Int32(entry.ColdEmitterID)
		}
	}

	if err := enc.Error(); err != nil {
		return fmt.Errorf("write PRT version %d: %w", prt.Version, err)
	}

	return nil
}
