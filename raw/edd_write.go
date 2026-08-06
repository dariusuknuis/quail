package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (edd *Edd) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	version := edd.Version
	if version == "" {
		version = "110"
	}
	if len(version) > 3 {
		return fmt.Errorf("EDD version %q is too long", version)
	}

	var versionData [4]byte
	copy(versionData[:], []byte(version))

	enc.Bytes([]byte{'E', 'D', 'D', 0})
	enc.Bytes(versionData[:])

	writeString := func(value string, size int) error {
		if len(value) >= size {
			return fmt.Errorf(
				"string %q is too long: maximum is %d bytes",
				value,
				size-1,
			)
		}

		data := make([]byte, size)
		copy(data, []byte(value))
		enc.Bytes(data)
		return nil
	}

	writeFloat3 := func(value [3]float32) {
		for i := 0; i < 3; i++ {
			enc.Float32(value[i])
		}
	}

	writeFloat2 := func(value [2]float32) {
		for i := 0; i < 2; i++ {
			enc.Float32(value[i])
		}
	}

	for i, entry := range edd.Entries {
		if entry == nil {
			return fmt.Errorf("entry %d is nil", i)
		}

		err := writeString(entry.EmitterName, eddNameLen)
		if err != nil {
			return fmt.Errorf("entry %d emitter name: %w", i, err)
		}

		err = writeString(entry.TextureName, eddTextureLen)
		if err != nil {
			return fmt.Errorf("entry %d texture name: %w", i, err)
		}

		enc.Int32(entry.RelativeToBone)
		enc.Int32(entry.NoOcclusion)
		enc.Int32(entry.AdditiveBlending)
		enc.Int32(entry.ScaleWithActor)
		enc.Int32(entry.SpriteOrientation)
		enc.Int32(entry.StickToActor)
		enc.Float32(entry.DefaultLifeSpan)
		enc.Float32(entry.ParticleLifeSpan)
		enc.Int32(entry.ParticlesAtCreation)
		enc.Int32(entry.ParticlesAtInterval)
		enc.Float32(entry.IntervalsPerSecond)
		enc.Float32(entry.SpawnDelay)
		enc.Float32(entry.FadeInTime)
		enc.Float32(entry.FadeOutTime)
		enc.Float32(entry.ScaleInTime)
		enc.Float32(entry.ScaleOutTime)
		enc.Float32(entry.ReductionDistance)
		enc.Float32(entry.MaxAlpha)
		enc.Int32(entry.SpawnShape)
		enc.Float32(entry.ShapeRadius)
		enc.Float32(entry.ShapeRadiusMinor)
		enc.Float32(entry.ShapeHeight)
		writeFloat3(entry.ShapeOffset)
		writeFloat2(entry.ShapeTilt)
		enc.Float32(entry.ParticleWidthMin)
		enc.Float32(entry.ParticleZBias)
		enc.Uint32(entry.TintStartR)
		enc.Uint32(entry.TintStartG)
		enc.Uint32(entry.TintStartB)
		enc.Uint32(entry.TintEndR)
		enc.Uint32(entry.TintEndG)
		enc.Uint32(entry.TintEndB)
		writeFloat3(entry.SpeedMin)
		writeFloat3(entry.SpeedMax)
		writeFloat3(entry.Acceleration)
		enc.Float32(entry.OutwardSpeedMin)
		enc.Float32(entry.OutwardSpeedMax)
		enc.Float32(entry.OutwardSpeedAcceleration)
		enc.Float32(entry.OrbitalSpeedMin)
		enc.Float32(entry.OrbitalSpeedMax)
		enc.Float32(entry.OrbitalSpeedAcceleration)
		enc.Float32(entry.ScalarGravity)
		enc.Float32(entry.WindSpeed)
		enc.Int32(entry.AnimationFrames)
		enc.Float32(entry.AnimationRate)
		enc.Float32(entry.ParticleSpinRate)
		enc.Int32(entry.OldParticleType)
		enc.Uint32(entry.OldFlags)
		enc.Uint32(entry.OldSize)
		writeFloat3(entry.Gravity)
		writeFloat3(entry.BBMin)
		writeFloat3(entry.BBMax)
		enc.Float32(entry.SpawnScale)
		enc.Float32(entry.Alpha)
		enc.Int32(entry.RandomRotation)
		enc.Int32(entry.ParticleOrientation)
		enc.Float32(entry.ParticleHeightMin)
		enc.Float32(entry.ParticleHeightMax)
		enc.Float32(entry.ParticleWidthMax)
		enc.Float32(entry.ParticleSpinRateMax)
		enc.Int32(entry.ProportionalSizeScaling)
		enc.Float32(entry.HeightSquashTime)
		enc.Float32(entry.WidthSquashTime)
		enc.Int32(entry.AllowCenterPassThrough)
		enc.Int32(entry.ScaleEmitterWithActor)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}