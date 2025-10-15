package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

const (
	EffOldRecordCount = 256
	EffOldRecordSize  = 0xAE8      // 2792
	EffOldHeaderSize  = 0x08       // 8
	EffOldBlockSize   = 0x3A0      // 928
	EffOldFileSize    = 698 * 1024 // 714,752 bytes
	str32Len          = 0x20
)

// EffOld models classic spells.eff (256 records, each with 3 sub-effect blocks)
type EffOld struct {
	MetaFileName string
	Records      []*EffOldRecord
}

type EffOldRecord struct {
	Header         [EffOldHeaderSize]byte
	Source         EffOldBlock
	SourceToTarget EffOldBlock
	Target         EffOldBlock
}

// EffOldBlock is one 0x3A0 sub-effect block containing 3 sub-effects + shared fields
type EffOldBlock struct {
	// Per-sub-effect triplet (index 0..2)
	Sub [3]EffSubEffect

	// Shared across the whole block
	Label        string     // "Source" / "Target" (observed)
	ExtraSprites [12]string // additional blitsprite references
	UnknownParam uint32
	SoundRef     uint32

	// Trailer / still-shared (not clearly per-sub)
	UnknownDW  [51]uint32
	UnknownF32 [12]float32
}

// EffSubEffect groups the fields that refer to the same logical sub-effect
type EffSubEffect struct {
	PrimarySprite string // first 3 strings (one per sub-effect)
	DagIndex      uint32 //1=head, 2=right hand, 3=left hand
	//4=right foot, 5=left foot, Other=chest
	EffectType    uint32
	ColorBGRA     uint32
	Gravity       float32
	SpawnNormal   [3]float32
	SpawnRadius   float32
	SpawnAngle    float32
	Lifespan      uint32
	SpawnVelocity float32
	SpawnRate     uint32
	SpawnScale    float32
}

func (eff *EffOld) Identity() string {
	return "eff_old"
}

func (eff *EffOld) String() string {
	return fmt.Sprintf("filename: %s\nrecords: %d (expected %d)\nsize(bytes): %d (expected %d)",
		eff.MetaFileName, len(eff.Records), EffOldRecordCount, len(eff.Records)*EffOldRecordSize, EffOldFileSize)
}

func (eff *EffOld) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	// sanity on size
	start := dec.Pos()
	end, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if _, err := r.Seek(start, io.SeekStart); err != nil {
		return fmt.Errorf("seek start: %w", err)
	}
	if end != EffOldFileSize {
		return fmt.Errorf("unexpected file size: %d (expected %d)", end, EffOldFileSize)
	}

	eff.Records = make([]*EffOldRecord, 0, EffOldRecordCount)
	for i := 0; i < EffOldRecordCount; i++ {
		rec := &EffOldRecord{}

		if _, err := io.ReadFull(r, rec.Header[:]); err != nil {
			return fmt.Errorf("record %d: header: %w", i, err)
		}
		if err := rec.Source.read(dec); err != nil {
			return fmt.Errorf("record %d: source: %w", i, err)
		}
		if err := rec.SourceToTarget.read(dec); err != nil {
			return fmt.Errorf("record %d: source_to_target: %w", i, err)
		}
		if err := rec.Target.read(dec); err != nil {
			return fmt.Errorf("record %d: target: %w", i, err)
		}

		e.Records = append(e.Records, rec)
	}

	// trailing sanity
	pos := dec.Pos()
	if pos < end {
		rem := make([]byte, end-pos)
		if _, err := io.ReadFull(r, rem); err == nil && len(rem) > 0 {
			fmt.Printf("remaining bytes:\n%s\n", hex.Dump(rem))
		}
		return fmt.Errorf("%d bytes remaining (%d total)", end-pos, end)
	}
	if pos > end {
		return fmt.Errorf("read past end of file: pos=%d end=%d", pos, end)
	}
	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}
	return nil
}

func (b *EffOldBlock) read(dec *encdec.Decoder) error {
	blkStart := dec.Pos()

	readStr32 := func() string {
		raw := dec.Bytes(str32Len)
		if dec.Error() != nil {
			return ""
		}
		if i := bytes.IndexByte(raw, 0); i >= 0 {
			raw = raw[:i]
		}
		return string(raw)
	}

	// 3 x 0x20 primary sprite strings (per sub-effect)
	for i := 0; i < 3; i++ {
		b.Sub[i].PrimarySprite = readStr32()
	}

	// 1 x 0x20 label string
	b.Label = readStr32()

	// 3 DagIndex (DWORD) -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].DagIndex = dec.Uint32()
	}
	// 3 EffectType (DWORD) -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].EffectType = dec.Uint32()
	}

	// 12 x 0x20 extra blitsprite strings (shared)
	for i := 0; i < 12; i++ {
		b.ExtraSprites[i] = readStr32()
	}

	// unknown + sound (shared)
	b.UnknownParam = dec.Uint32()
	b.SoundRef = dec.Uint32()

	// 3 BGRA colors -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].ColorBGRA = dec.Uint32()
	}
	// 3 Gravity -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].Gravity = dec.Float32()
	}
	// 3 SpawnNormals (each vec3) -> per sub-effect
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b.Sub[i].SpawnNormal[j] = dec.Float32()
		}
	}
	// 3 SpawnRadii -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].SpawnRadius = dec.Float32()
	}
	// 3 SpawnAngle -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].SpawnAngle = dec.Float32()
	}
	// 3 Lifespan (DWORD) -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].Lifespan = dec.Uint32()
	}
	// 3 SpawnVelocity -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].SpawnVelocity = dec.Float32()
	}
	// 3 SpawnRate (DWORD) -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].SpawnRate = dec.Uint32()
	}
	// 3 SpawnScale -> per sub-effect
	for i := 0; i < 3; i++ {
		b.Sub[i].SpawnScale = dec.Float32()
	}

	// 51 unknown DWORDs (shared)
	for i := 0; i < 51; i++ {
		b.UnknownDW[i] = dec.Uint32()
	}
	// 12 unknown floats (shared)
	for i := 0; i < 12; i++ {
		b.UnknownF32[i] = dec.Float32()
	}

	// exact size check of the block
	blkEnd := dec.Pos()
	if blkEnd-blkStart != EffOldBlockSize {
		return fmt.Errorf("eff block size mismatch: read %d, expected %d", blkEnd-blkStart, EffOldBlockSize)
	}
	return dec.Error()
}

// SetFileName sets the file name
func (eff *EffOld) SetFileName(name string) {
	eff.MetaFileName = name
}

// FileName returns the file name
func (eff *EffOld) FileName() string {
	return eff.MetaFileName
}
