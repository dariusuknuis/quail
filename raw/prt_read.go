package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Prt is a Particle Render
type Prt struct {
	MetaFileName string
	Version      uint32
	Entries      []*PrtEntry
}

// PrtEntry is  ParticleRender entry
type PrtEntry struct {
	EmitterID     int32
	ParticlePoint string
	ParticleType  int32
	AnimNumber    int32
	AnimVariation int32
	RandomAnim    int32
	StartTime     int32
	Lifespan      int32
	Ground        int32
	PlayWithMat   int32
	Sporadic      int32
	ColdEmitterID int32
}

// Identity returns the type of the struct
func (prt *Prt) Identity() string {
	return "prt"
}

func (prt *Prt) String() string {
	out := ""
	out += fmt.Sprintf("metafilename: %s\n", prt.MetaFileName)
	out += fmt.Sprintf("version: %d\n", prt.Version)
	out += fmt.Sprintf("entries: %d\n", len(prt.Entries))

	for i, entry := range prt.Entries {
		out += fmt.Sprintf(
			"  %d: point=%s emitter=%d coldEmitter=%d "+
				"type=%d animation=%d variation=%d randomVariation=%d "+
				"start=%d lifespan=%d ground=%d playWithMat=%d sporadic=%d\n",
			i,
			entry.ParticlePoint,
			entry.EmitterID,
			entry.ColdEmitterID,
			entry.ParticleType,
			entry.AnimNumber,
			entry.AnimVariation,
			entry.RandomAnim,
			entry.StartTime,
			entry.Lifespan,
			entry.Ground,
			entry.PlayWithMat,
			entry.Sporadic,
		)
	}

	return out
}

// Read reads a PRT file
func (prt *Prt) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "PTCL" {
		return fmt.Errorf("invalid header %q, wanted PTCL", header)
	}

	particleCount := dec.Uint32()
	prt.Version = dec.Uint32()

	if prt.Version < 1 || prt.Version > 5 {
		return fmt.Errorf(
			"unsupported PRT version %d, wanted version 1 through 5",
			prt.Version,
		)
	}

	prt.Entries = nil

	for i := uint32(0); i < particleCount; i++ {
		entry := &PrtEntry{}

		// Present in every version.
		entry.EmitterID = dec.Int32()

		// Fixed char particlePointName[64].
		pointData := dec.Bytes(64)

		nullIndex := bytes.IndexByte(pointData, 0)
		if nullIndex < 0 {
			return fmt.Errorf(
				"particle %d: particle-point name is not null-terminated",
				i,
			)
		}

		entry.ParticlePoint = string(pointData[:nullIndex])

		// Present in every version.
		entry.ParticleType = dec.Int32()
		entry.AnimNumber = dec.Int32()
		entry.AnimVariation = dec.Int32()
		entry.RandomAnim = dec.Int32()
		entry.StartTime = dec.Int32()
		entry.Lifespan = dec.Int32()

		// Added in version 2.
		if prt.Version >= 2 {
			entry.Ground = dec.Int32()
		}

		// Added in version 3.
		if prt.Version >= 3 {
			entry.PlayWithMat = dec.Int32()
		} else {
			entry.PlayWithMat = -1
		}

		// Added in version 4.
		if prt.Version >= 4 {
			entry.Sporadic = dec.Int32()
		}

		// Added in version 5 at the end of the record.
		if prt.Version >= 5 {
			entry.ColdEmitterID = dec.Int32()
		} else {
			entry.ColdEmitterID = entry.EmitterID
		}

		prt.Entries = append(prt.Entries, entry)
	}

	if err := dec.Error(); err != nil {
		return fmt.Errorf("read PRT version %d: %w", prt.Version, err)
	}

	// Check for trailing data
	pos := dec.Pos()

	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek to end: %w", err)
	}

	if pos > endPos {
		return fmt.Errorf(
			"read past end of PRT file: position %d, size %d",
			pos,
			endPos,
		)
	}

	if pos < endPos {
		if _, err := r.Seek(pos, io.SeekStart); err != nil {
			return fmt.Errorf("seek to trailing data: %w", err)
		}

		remaining, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read trailing data: %w", err)
		}

		// Some files contain a trailing zero DWORD. Accept zero
		// padding, but reject unexplained nonzero data.
		if len(bytes.Trim(remaining, "\x00")) != 0 {
			fmt.Printf(
				"remaining PRT bytes:\n%s\n",
				hex.Dump(remaining),
			)

			return fmt.Errorf(
				"%d non-padding bytes remain in PRT file",
				len(remaining),
			)
		}
	}

	return nil
}

// SetFileName sets the name of the file
func (prt *Prt) SetFileName(name string) {
	prt.MetaFileName = name
}

// FileName returns the name of the file
func (prt *Prt) FileName() string {
	return prt.MetaFileName
}
