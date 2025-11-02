package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (eff *EffOld) Write(w io.Writer) error {
	if len(eff.Records) != 256 {
		return fmt.Errorf("expected 256 records, got %d", len(eff.Records))
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)

	writeStr32 := func(s string) {
		var b [str32Len]byte
		copy(b[:], []byte(s))
		enc.Bytes(b[:])
	}

	for i := 0; i < 256; i++ {
		rec := eff.Records[i]

		// header (2 DWORDs)
		enc.Uint32(rec.Header[0])
		enc.Uint32(rec.Header[1])

		// three blocks in order: Source, SourceToTarget, Target
		blocks := []*EffOldBlock{
			&rec.Source,
			&rec.SourceToTarget,
			&rec.Target,
		}
		for bi := 0; bi < 3; bi++ {
			b := blocks[bi]

			// 3x primary sprite strings
			writeStr32(b.SubEffect[0].Blit)
			writeStr32(b.SubEffect[1].Blit)
			writeStr32(b.SubEffect[2].Blit)

			// label
			writeStr32(b.Label)

			// 3x DagIndex
			enc.Int32(b.SubEffect[0].DagIndex)
			enc.Int32(b.SubEffect[1].DagIndex)
			enc.Int32(b.SubEffect[2].DagIndex)

			// 3x EffectType
			enc.Int32(b.SubEffect[0].EffectType)
			enc.Int32(b.SubEffect[1].EffectType)
			enc.Int32(b.SubEffect[2].EffectType)

			// 12 extra sprites
			for j := 0; j < 12; j++ {
				writeStr32(b.ExtraEffect[j].Blit)
			}

			// shared params
			enc.Int32(b.EffectMode)
			enc.Int32(b.SoundRef)

			// colors BGRA (4 bytes each, matches your read)
			for s := 0; s < 3; s++ {
				enc.Uint8(b.SubEffect[s].ColorBGRA[2])
				enc.Uint8(b.SubEffect[s].ColorBGRA[1])
				enc.Uint8(b.SubEffect[s].ColorBGRA[0])
				enc.Uint8(b.SubEffect[s].ColorBGRA[3])
			}

			// gravity
			enc.Float32(b.SubEffect[0].Gravity)
			enc.Float32(b.SubEffect[1].Gravity)
			enc.Float32(b.SubEffect[2].Gravity)

			// spawn normal vec3 per sub
			for s := 0; s < 3; s++ {
				enc.Float32(b.SubEffect[s].SpawnNormal[0])
				enc.Float32(b.SubEffect[s].SpawnNormal[1])
				enc.Float32(b.SubEffect[s].SpawnNormal[2])
			}

			// spawn radius
			enc.Float32(b.SubEffect[0].SpawnRadius)
			enc.Float32(b.SubEffect[1].SpawnRadius)
			enc.Float32(b.SubEffect[2].SpawnRadius)

			// spawn angle
			enc.Float32(b.SubEffect[0].SpawnAngle)
			enc.Float32(b.SubEffect[1].SpawnAngle)
			enc.Float32(b.SubEffect[2].SpawnAngle)

			// lifespan (DWORD)
			enc.Uint32(b.SubEffect[0].Lifespan)
			enc.Uint32(b.SubEffect[1].Lifespan)
			enc.Uint32(b.SubEffect[2].Lifespan)

			// spawn velocity
			enc.Float32(b.SubEffect[0].SpawnVelocity)
			enc.Float32(b.SubEffect[1].SpawnVelocity)
			enc.Float32(b.SubEffect[2].SpawnVelocity)

			// spawn rate (DWORD)
			enc.Uint32(b.SubEffect[0].SpawnRate)
			enc.Uint32(b.SubEffect[1].SpawnRate)
			enc.Uint32(b.SubEffect[2].SpawnRate)

			// spawn scale
			enc.Float32(b.SubEffect[0].SpawnScale)
			enc.Float32(b.SubEffect[1].SpawnScale)
			enc.Float32(b.SubEffect[2].SpawnScale)

			for j := 0; j < 12; j++ {
				enc.Uint8(b.ExtraEffect[j].ColorBGR[2])
				enc.Uint8(b.ExtraEffect[j].ColorBGR[1])
				enc.Uint8(b.ExtraEffect[j].ColorBGR[0])
			}

			for j := 0; j < 12; j++ {
				enc.Float32(b.ExtraEffect[j].AnimSpeedMultiplier)
			}

			for j := 0; j < 12; j++ {
				enc.Int16(b.ExtraEffect[j].AngleRangeA)
			}

			for j := 0; j < 12; j++ {
				enc.Int16(b.ExtraEffect[j].AngleRangeB)
			}

			for j := 0; j < 12; j++ {
				enc.Float32(b.ExtraEffect[j].Radius)
			}

			for j := 0; j < 12; j++ {
				enc.Int16(b.ExtraEffect[j].EffectType)
			}

			for j := 0; j < 12; j++ {
				enc.Float32(b.ExtraEffect[j].Scale)
			}
		}
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
