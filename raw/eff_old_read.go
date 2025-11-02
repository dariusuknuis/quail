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
	Label       string             // "Source" / "Target" (observed)
	ExtraEffect [12]EffExtraEffect // additional blitsprite references
	EffectMode  int32
	SoundRef    int32

	// Trailer / still-shared (not clearly per-sub)
	UnknownDW  [51]uint32
	UnknownF32 [12]float32
}

// EffSubEffect groups the fields that refer to the same logical sub-effect
type EffSubEffect struct {
	Blit     string // first 3 strings (one per sub-effect)
	DagIndex int32  //1=head, 2=right hand, 3=left hand
	//4=right foot, 5=left foot, Other=chest
	EffectType    int32
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

type EffExtraEffect struct {
	Blit                string // first 3 strings (one per sub-effect)
	ColorBGR            [3]uint8
	AnimSpeedMultiplier float32
	AngleRangeA         int16
	AngleRangeB         int16
	Radius              float32
	EffectType          int16
	Scale               float32
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
			for j := 0; j < 3; j++ {
				block.SubEffect[j].Blit = readStr32()
			}

			block.Label = readStr32()
			for j := 0; j < 3; j++ {
				block.SubEffect[j].DagIndex = dec.Int32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].EffectType = dec.Int32()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].Blit = readStr32()
			}

			block.EffectMode = dec.Int32()
			block.SoundRef = dec.Int32()
			for j := 0; j < 3; j++ {
				for k := 0; k < 4; k++ {
					block.SubEffect[j].ColorBGRA[k] = dec.Uint8()
				}
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].Gravity = dec.Float32()
			}

			for j := 0; j < 3; j++ {
				for k := 0; k < 3; k++ {
					block.SubEffect[j].SpawnNormal[k] = dec.Float32()
				}
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].SpawnRadius = dec.Float32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].SpawnAngle = dec.Float32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].Lifespan = dec.Uint32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].SpawnVelocity = dec.Float32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].SpawnRate = dec.Uint32()
			}

			for j := 0; j < 3; j++ {
				block.SubEffect[j].SpawnScale = dec.Float32()
			}

			for j := 0; j < 12; j++ {
				for k := 0; k < 3; k++ {
					block.ExtraEffect[j].ColorBGR[k] = dec.Uint8()
				}
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].AnimSpeedMultiplier = dec.Float32()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].AngleRangeA = dec.Int16()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].AngleRangeB = dec.Int16()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].Radius = dec.Float32()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].EffectType = dec.Int16()
			}

			for j := 0; j < 12; j++ {
				block.ExtraEffect[j].Scale = dec.Float32()
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
