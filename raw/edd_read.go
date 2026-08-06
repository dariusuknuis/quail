package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

const (
	EddHeaderSize  = 0x08
	EddRecordSize  = 0x1A0
	eddNameLen     = 0x40
	eddTextureLen  = 0x20
)

// Edd contains particle definitions used by prt.
// Examples are actoremittersnew.edd, environmentemittersnew.edd,
// and spellsnew.edd.
type Edd struct {
	MetaFileName string
	Version      string
	Entries      []*EddEntry
}

type EddEntry struct {
	EmitterName                 string
	TextureName                 string
	RelativeToBone              int32
	NoOcclusion                 int32
	AdditiveBlending            int32
	ScaleWithActor              int32
	SpriteOrientation           int32
	StickToActor                int32
	DefaultLifeSpan             float32
	ParticleLifeSpan            float32
	ParticlesAtCreation         int32
	ParticlesAtInterval         int32
	IntervalsPerSecond          float32
	SpawnDelay                  float32
	FadeInTime                  float32
	FadeOutTime                 float32
	ScaleInTime                 float32
	ScaleOutTime                float32
	ReductionDistance           float32
	MaxAlpha                    float32
	SpawnShape                  int32
	ShapeRadius                 float32
	ShapeRadiusMinor            float32
	ShapeHeight                 float32
	ShapeOffset                 [3]float32
	ShapeTilt                   [2]float32
	ParticleWidthMin            float32
	ParticleZBias               float32
	TintStartR                  uint32
	TintStartG                  uint32
	TintStartB                  uint32
	TintEndR                    uint32
	TintEndG                    uint32
	TintEndB                    uint32
	SpeedMin                    [3]float32
	SpeedMax                    [3]float32
	Acceleration                [3]float32
	OutwardSpeedMin             float32
	OutwardSpeedMax             float32
	OutwardSpeedAcceleration    float32
	OrbitalSpeedMin             float32
	OrbitalSpeedMax             float32
	OrbitalSpeedAcceleration    float32
	ScalarGravity               float32
	WindSpeed                   float32
	AnimationFrames             int32
	AnimationRate               float32
	ParticleSpinRate            float32
	OldParticleType             int32
	OldFlags                    uint32
	OldSize                     uint32
	Gravity                     [3]float32
	BBMin                       [3]float32
	BBMax                       [3]float32
	SpawnScale                  float32
	Alpha                       float32
	RandomRotation              int32
	ParticleOrientation         int32
	ParticleHeightMin           float32
	ParticleHeightMax           float32
	ParticleWidthMax            float32
	ParticleSpinRateMax         float32
	ProportionalSizeScaling     int32
	HeightSquashTime            float32
	WidthSquashTime             float32
	AllowCenterPassThrough      int32
	ScaleEmitterWithActor       int32
}

// Identity returns the type of the struct
func (edd *Edd) Identity() string {
	return "edd"
}

func (edd *Edd) Read(r io.ReadSeeker) error {
	fileSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if fileSize < EddHeaderSize {
		return fmt.Errorf("invalid EDD size %d", fileSize)
	}
	if (fileSize-EddHeaderSize)%EddRecordSize != 0 {
		return fmt.Errorf(
			"invalid EDD size %d: data size is not divisible by record size %d",
			fileSize,
			EddRecordSize,
		)
	}
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek start: %w", err)
	}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.Bytes(4)
	if !bytes.Equal(header, []byte{'E', 'D', 'D', 0}) {
		return fmt.Errorf("invalid header %q, wanted EDD", header)
	}

	versionData := dec.Bytes(4)
	if n := bytes.IndexByte(versionData, 0); n >= 0 {
		versionData = versionData[:n]
	}
	edd.Version = string(versionData)

	recordCount := int((fileSize - EddHeaderSize) / EddRecordSize)
	edd.Entries = nil

	readFloat3 := func(value *[3]float32) {
		for i := 0; i < 3; i++ {
			value[i] = dec.Float32()
		}
	}

	readFloat2 := func(value *[2]float32) {
		for i := 0; i < 2; i++ {
			value[i] = dec.Float32()
		}
	}

	for i := 0; i < recordCount; i++ {
		entry := &EddEntry{}

		nameData := dec.Bytes(eddNameLen)
		if n := bytes.IndexByte(nameData, 0); n >= 0 {
			nameData = nameData[:n]
		}
		entry.EmitterName = string(nameData)

		textureData := dec.Bytes(eddTextureLen)
		if n := bytes.IndexByte(textureData, 0); n >= 0 {
			textureData = textureData[:n]
		}
		entry.TextureName = string(textureData)

		entry.RelativeToBone = dec.Int32()
		entry.NoOcclusion = dec.Int32()
		entry.AdditiveBlending = dec.Int32()
		entry.ScaleWithActor = dec.Int32()
		entry.SpriteOrientation = dec.Int32()
		entry.StickToActor = dec.Int32()
		entry.DefaultLifeSpan = dec.Float32()
		entry.ParticleLifeSpan = dec.Float32()
		entry.ParticlesAtCreation = dec.Int32()
		entry.ParticlesAtInterval = dec.Int32()
		entry.IntervalsPerSecond = dec.Float32()
		entry.SpawnDelay = dec.Float32()
		entry.FadeInTime = dec.Float32()
		entry.FadeOutTime = dec.Float32()
		entry.ScaleInTime = dec.Float32()
		entry.ScaleOutTime = dec.Float32()
		entry.ReductionDistance = dec.Float32()
		entry.MaxAlpha = dec.Float32()
		entry.SpawnShape = dec.Int32()
		entry.ShapeRadius = dec.Float32()
		entry.ShapeRadiusMinor = dec.Float32()
		entry.ShapeHeight = dec.Float32()
		readFloat3(&entry.ShapeOffset)
		readFloat2(&entry.ShapeTilt)
		entry.ParticleWidthMin = dec.Float32()
		entry.ParticleZBias = dec.Float32()
		entry.TintStartR = dec.Uint32()
		entry.TintStartG = dec.Uint32()
		entry.TintStartB = dec.Uint32()
		entry.TintEndR = dec.Uint32()
		entry.TintEndG = dec.Uint32()
		entry.TintEndB = dec.Uint32()
		readFloat3(&entry.SpeedMin)
		readFloat3(&entry.SpeedMax)
		readFloat3(&entry.Acceleration)
		entry.OutwardSpeedMin = dec.Float32()
		entry.OutwardSpeedMax = dec.Float32()
		entry.OutwardSpeedAcceleration = dec.Float32()
		entry.OrbitalSpeedMin = dec.Float32()
		entry.OrbitalSpeedMax = dec.Float32()
		entry.OrbitalSpeedAcceleration = dec.Float32()
		entry.ScalarGravity = dec.Float32()
		entry.WindSpeed = dec.Float32()
		entry.AnimationFrames = dec.Int32()
		entry.AnimationRate = dec.Float32()
		entry.ParticleSpinRate = dec.Float32()
		entry.OldParticleType = dec.Int32()
		entry.OldFlags = dec.Uint32()
		entry.OldSize = dec.Uint32()
		readFloat3(&entry.Gravity)
		readFloat3(&entry.BBMin)
		readFloat3(&entry.BBMax)
		entry.SpawnScale = dec.Float32()
		entry.Alpha = dec.Float32()
		entry.RandomRotation = dec.Int32()
		entry.ParticleOrientation = dec.Int32()
		entry.ParticleHeightMin = dec.Float32()
		entry.ParticleHeightMax = dec.Float32()
		entry.ParticleWidthMax = dec.Float32()
		entry.ParticleSpinRateMax = dec.Float32()
		entry.ProportionalSizeScaling = dec.Int32()
		entry.HeightSquashTime = dec.Float32()
		entry.WidthSquashTime = dec.Float32()
		entry.AllowCenterPassThrough = dec.Int32()
		entry.ScaleEmitterWithActor = dec.Int32()

		edd.Entries = append(edd.Entries, entry)
	}

	err = dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if dec.Pos() != fileSize {
		return fmt.Errorf(
			"read ended at %d, expected %d",
			dec.Pos(),
			fileSize,
		)
	}

	return nil
}

// SetFileName sets the name of the file
func (edd *Edd) SetFileName(name string) {
	edd.MetaFileName = name
}

// FileName returns the name of the file
func (edd *Edd) FileName() string {
	return edd.MetaFileName
}