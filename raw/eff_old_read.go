package raw

import (
	"bytes"
	"encoding/binary"
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
	Header         [2]uint32
	Source         EffOldBlock
	SourceToTarget EffOldBlock
	Target         EffOldBlock
}

// EffOldBlock is one 0x3A0 sub-effect block containing 3 sub-effects + shared fields
type EffOldBlock struct {
	// Per-sub-effect triplet (index 0..2)
	SubEffect [3]EffSubEffect

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
	Blit     string // first 3 strings (one per sub-effect)
	DagIndex uint32 //1=head, 2=right hand, 3=left hand
	//4=right foot, 5=left foot, Other=chest
	EffectType    uint32
	ColorBGRA     [4]uint8
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
	out := fmt.Sprintf("EffOld: %s,", eff.MetaFileName)
	out += fmt.Sprintf("records: %d,", len(eff.Records))
	return out
}

func (eff *EffOld) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	readStr32 := func() string {
		raw := dec.Bytes(str32Len)
		if n := bytes.IndexByte(raw, 0); n >= 0 {
			raw = raw[:n]
		}
		return string(raw)
	}
	for i := 0; i < 256; i++ {
		effRec := &EffOldRecord{}
		effRec.Header[0] = dec.Uint32()
		effRec.Header[1] = dec.Uint32()
		blocks := []*EffOldBlock{
			&effRec.Source,
			&effRec.SourceToTarget,
			&effRec.Target,
		}
		for b := 0; b < 3; b++ {
			block := blocks[b]
			block.SubEffect[0].Blit = readStr32()
			block.SubEffect[1].Blit = readStr32()
			block.SubEffect[2].Blit = readStr32()
			block.Label = readStr32()
			block.SubEffect[0].DagIndex = dec.Uint32()
			block.SubEffect[1].DagIndex = dec.Uint32()
			block.SubEffect[2].DagIndex = dec.Uint32()
			block.SubEffect[0].EffectType = dec.Uint32()
			block.SubEffect[1].EffectType = dec.Uint32()
			block.SubEffect[2].EffectType = dec.Uint32()
			for j := 0; j < 12; j++ {
				block.ExtraSprites[j] = readStr32()
			}
			block.UnknownParam = dec.Uint32()
			block.SoundRef = dec.Uint32()
			block.SubEffect[0].ColorBGRA[0] = dec.Uint8()
			block.SubEffect[0].ColorBGRA[1] = dec.Uint8()
			block.SubEffect[0].ColorBGRA[2] = dec.Uint8()
			block.SubEffect[0].ColorBGRA[3] = dec.Uint8()
			block.SubEffect[1].ColorBGRA[0] = dec.Uint8()
			block.SubEffect[1].ColorBGRA[1] = dec.Uint8()
			block.SubEffect[1].ColorBGRA[2] = dec.Uint8()
			block.SubEffect[1].ColorBGRA[3] = dec.Uint8()
			block.SubEffect[2].ColorBGRA[0] = dec.Uint8()
			block.SubEffect[2].ColorBGRA[1] = dec.Uint8()
			block.SubEffect[2].ColorBGRA[2] = dec.Uint8()
			block.SubEffect[2].ColorBGRA[3] = dec.Uint8()
			block.SubEffect[0].Gravity = dec.Float32()
			block.SubEffect[1].Gravity = dec.Float32()
			block.SubEffect[2].Gravity = dec.Float32()
			block.SubEffect[0].SpawnNormal[0] = dec.Float32()
			block.SubEffect[0].SpawnNormal[1] = dec.Float32()
			block.SubEffect[0].SpawnNormal[2] = dec.Float32()
			block.SubEffect[1].SpawnNormal[0] = dec.Float32()
			block.SubEffect[1].SpawnNormal[1] = dec.Float32()
			block.SubEffect[1].SpawnNormal[2] = dec.Float32()
			block.SubEffect[2].SpawnNormal[0] = dec.Float32()
			block.SubEffect[2].SpawnNormal[1] = dec.Float32()
			block.SubEffect[2].SpawnNormal[2] = dec.Float32()
			block.SubEffect[0].SpawnRadius = dec.Float32()
			block.SubEffect[1].SpawnRadius = dec.Float32()
			block.SubEffect[2].SpawnRadius = dec.Float32()
			block.SubEffect[0].SpawnAngle = dec.Float32()
			block.SubEffect[1].SpawnAngle = dec.Float32()
			block.SubEffect[2].SpawnAngle = dec.Float32()
			block.SubEffect[0].Lifespan = dec.Uint32()
			block.SubEffect[1].Lifespan = dec.Uint32()
			block.SubEffect[2].Lifespan = dec.Uint32()
			block.SubEffect[0].SpawnVelocity = dec.Float32()
			block.SubEffect[1].SpawnVelocity = dec.Float32()
			block.SubEffect[2].SpawnVelocity = dec.Float32()
			block.SubEffect[0].SpawnRate = dec.Uint32()
			block.SubEffect[1].SpawnRate = dec.Uint32()
			block.SubEffect[2].SpawnRate = dec.Uint32()
			block.SubEffect[0].SpawnScale = dec.Float32()
			block.SubEffect[1].SpawnScale = dec.Float32()
			block.SubEffect[2].SpawnScale = dec.Float32()
			for k := 0; k < 51; k++ {
				block.UnknownDW[k] = dec.Uint32()
			}
			for l := 0; l < 12; l++ {
				block.UnknownF32[l] = dec.Float32()
			}
		}
		eff.Records = append(eff.Records, effRec)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// SetFileName sets the file name
func (eff *EffOld) SetFileName(name string) {
	eff.MetaFileName = name
}

// FileName returns the file name
func (eff *EffOld) FileName() string {
	return eff.MetaFileName
}
