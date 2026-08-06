package wce

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
)

type EmitterDef struct {
	folders                  []string
	TagIndex                 int
	Name                     string
	Texture                  string
	RelativeToBone           int32
	NoOcclusion              int32
	AdditiveBlending         int32
	ScaleWithActor           int32
	SpriteOrientation        int32
	StickToActor             int32
	DefaultLifeSpan          float32
	ParticleLifeSpan         float32
	ParticlesAtCreation      int32
	ParticlesAtInterval      int32
	IntervalsPerSecond       float32
	SpawnDelay               float32
	FadeInTime               float32
	FadeOutTime              float32
	ScaleInTime              float32
	ScaleOutTime             float32
	ReductionDistance        float32
	MaxAlpha                 float32
	SpawnShape               int32
	ShapeRadius              float32
	ShapeRadiusMinor         float32
	ShapeHeight              float32
	ShapeOffset              [3]float32
	ShapeTilt                [2]float32
	ParticleWidthMin         float32
	ParticleZBias            float32
	TintStart                [3]uint32
	TintEnd                  [3]uint32
	SpeedMin                 [3]float32
	SpeedMax                 [3]float32
	Acceleration             [3]float32
	OutwardSpeedMin          float32
	OutwardSpeedMax          float32
	OutwardSpeedAcceleration float32
	OrbitalSpeedMin          float32
	OrbitalSpeedMax          float32
	OrbitalSpeedAcceleration float32
	ScalarGravity            float32
	WindSpeed                float32
	AnimationFrames          int32
	AnimationRate            float32
	ParticleSpinRate         float32
	OldParticleType          int32
	OldFlags                 uint32
	OldSize                  uint32
	Gravity                  [3]float32
	BBMin                    [3]float32
	BBMax                    [3]float32
	SpawnScale               float32
	Alpha                    float32
	RandomRotation           int32
	ParticleOrientation      int32
	ParticleHeightMin        float32
	ParticleHeightMax        float32
	ParticleWidthMax         float32
	ParticleSpinRateMax      float32
	ProportionalSizeScaling  int32
	HeightSquashTime         float32
	WidthSquashTime          float32
	AllowCenterPassThrough   int32
	ScaleEmitterWithActor    int32
}

func (e *EmitterDef) Definition() string {
	return "EQGEMITTERDEF"
}

func (e *EmitterDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		if err := token.SetWriter(folder); err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "%s \"%d\"\n", e.Definition(), e.TagIndex)
		fmt.Fprintf(w, "\tNAME \"%s\"\n", e.Name)
		fmt.Fprintf(w, "\tTEXTURE \"%s\"\n", e.Texture)
		fmt.Fprintf(w, "\tRELATIVETOBONE %d\n", e.RelativeToBone)
		fmt.Fprintf(w, "\tNOOCCLUSION %d\n", e.NoOcclusion)
		fmt.Fprintf(w, "\tADDITIVEBLENDING %d\n", e.AdditiveBlending)
		fmt.Fprintf(w, "\tSCALEWITHACTOR %d\n", e.ScaleWithActor)
		fmt.Fprintf(w, "\tSPRITEORIENTATION %d\n", e.SpriteOrientation)
		fmt.Fprintf(w, "\tSTICKTOACTOR %d\n", e.StickToActor)
		fmt.Fprintf(w, "\tDEFAULTLIFESPAN %0.8e\n", e.DefaultLifeSpan)
		fmt.Fprintf(w, "\tPARTICLELIFESPAN %0.8e\n", e.ParticleLifeSpan)
		fmt.Fprintf(w, "\tPARTICLESATCREATION %d\n", e.ParticlesAtCreation)
		fmt.Fprintf(w, "\tPARTICLESATINTERVAL %d\n", e.ParticlesAtInterval)
		fmt.Fprintf(w, "\tINTERVALSPERSECOND %0.8e\n", e.IntervalsPerSecond)
		fmt.Fprintf(w, "\tSPAWNDELAY %0.8e\n", e.SpawnDelay)
		fmt.Fprintf(w, "\tFADEINTIME %0.8e\n", e.FadeInTime)
		fmt.Fprintf(w, "\tFADEOUTTIME %0.8e\n", e.FadeOutTime)
		fmt.Fprintf(w, "\tSCALEINTIME %0.8e\n", e.ScaleInTime)
		fmt.Fprintf(w, "\tSCALEOUTTIME %0.8e\n", e.ScaleOutTime)
		fmt.Fprintf(w, "\tREDUCTIONDISTANCE %0.8e\n", e.ReductionDistance)
		fmt.Fprintf(w, "\tMAXALPHA %0.8e\n", e.MaxAlpha)
		fmt.Fprintf(w, "\tSPAWNSHAPE %d\n", e.SpawnShape)
		fmt.Fprintf(w, "\tSHAPERADIUS %0.8e\n", e.ShapeRadius)
		fmt.Fprintf(w, "\tSHAPERADIUSMINOR %0.8e\n", e.ShapeRadiusMinor)
		fmt.Fprintf(w, "\tSHAPEHEIGHT %0.8e\n", e.ShapeHeight)
		fmt.Fprintf(w, "\tSHAPEOFFSET %0.8e %0.8e %0.8e\n", e.ShapeOffset[0], e.ShapeOffset[1], e.ShapeOffset[2])
		fmt.Fprintf(w, "\tSHAPETILT %0.8e %0.8e\n", e.ShapeTilt[0], e.ShapeTilt[1])
		fmt.Fprintf(w, "\tPARTICLEWIDTHMIN %0.8e\n", e.ParticleWidthMin)
		fmt.Fprintf(w, "\tPARTICLEZBIAS %0.8e\n", e.ParticleZBias)
		fmt.Fprintf(w, "\tTINTSTART %d %d %d\n", e.TintStart[0], e.TintStart[1], e.TintStart[2])
		fmt.Fprintf(w, "\tTINTEND %d %d %d\n", e.TintEnd[0], e.TintEnd[1], e.TintEnd[2])
		fmt.Fprintf(w, "\tSPEEDMIN %0.8e %0.8e %0.8e\n", e.SpeedMin[0], e.SpeedMin[1], e.SpeedMin[2])
		fmt.Fprintf(w, "\tSPEEDMAX %0.8e %0.8e %0.8e\n", e.SpeedMax[0], e.SpeedMax[1], e.SpeedMax[2])
		fmt.Fprintf(w, "\tACCELERATION %0.8e %0.8e %0.8e\n", e.Acceleration[0], e.Acceleration[1], e.Acceleration[2])
		fmt.Fprintf(w, "\tOUTWARDSPEEDMIN %0.8e\n", e.OutwardSpeedMin)
		fmt.Fprintf(w, "\tOUTWARDSPEEDMAX %0.8e\n", e.OutwardSpeedMax)
		fmt.Fprintf(w, "\tOUTWARDSPEEDACCELERATION %0.8e\n", e.OutwardSpeedAcceleration)
		fmt.Fprintf(w, "\tORBITALSPEEDMIN %0.8e\n", e.OrbitalSpeedMin)
		fmt.Fprintf(w, "\tORBITALSPEEDMAX %0.8e\n", e.OrbitalSpeedMax)
		fmt.Fprintf(w, "\tORBITALSPEEDACCELERATION %0.8e\n", e.OrbitalSpeedAcceleration)
		fmt.Fprintf(w, "\tSCALARGRAVITY %0.8e\n", e.ScalarGravity)
		fmt.Fprintf(w, "\tWINDSPEED %0.8e\n", e.WindSpeed)
		fmt.Fprintf(w, "\tANIMATIONFRAMES %d\n", e.AnimationFrames)
		fmt.Fprintf(w, "\tANIMATIONRATE %0.8e\n", e.AnimationRate)
		fmt.Fprintf(w, "\tPARTICLESPINRATE %0.8e\n", e.ParticleSpinRate)
		fmt.Fprintf(w, "\tOLDPARTICLETYPE %d\n", e.OldParticleType)
		fmt.Fprintf(w, "\tOLDFLAGS %d\n", e.OldFlags)
		fmt.Fprintf(w, "\tOLDSIZE %d\n", e.OldSize)
		fmt.Fprintf(w, "\tGRAVITY %0.8e %0.8e %0.8e\n", e.Gravity[0], e.Gravity[1], e.Gravity[2])
		fmt.Fprintf(w, "\tBBMIN %0.8e %0.8e %0.8e\n", e.BBMin[0], e.BBMin[1], e.BBMin[2])
		fmt.Fprintf(w, "\tBBMAX %0.8e %0.8e %0.8e\n", e.BBMax[0], e.BBMax[1], e.BBMax[2])
		fmt.Fprintf(w, "\tSPAWNSCALE %0.8e\n", e.SpawnScale)
		fmt.Fprintf(w, "\tALPHA %0.8e\n", e.Alpha)
		fmt.Fprintf(w, "\tRANDOMROTATION %d\n", e.RandomRotation)
		fmt.Fprintf(w, "\tPARTICLEORIENTATION %d\n", e.ParticleOrientation)
		fmt.Fprintf(w, "\tPARTICLEHEIGHTMIN %0.8e\n", e.ParticleHeightMin)
		fmt.Fprintf(w, "\tPARTICLEHEIGHTMAX %0.8e\n", e.ParticleHeightMax)
		fmt.Fprintf(w, "\tPARTICLEWIDTHMAX %0.8e\n", e.ParticleWidthMax)
		fmt.Fprintf(w, "\tPARTICLESPINRATEMAX %0.8e\n", e.ParticleSpinRateMax)
		fmt.Fprintf(w, "\tPROPORTIONALSIZESCALING %d\n", e.ProportionalSizeScaling)
		fmt.Fprintf(w, "\tHEIGHTSQUASHTIME %0.8e\n", e.HeightSquashTime)
		fmt.Fprintf(w, "\tWIDTHSQUASHTIME %0.8e\n", e.WidthSquashTime)
		fmt.Fprintf(w, "\tALLOWCENTERPASSTHROUGH %d\n", e.AllowCenterPassThrough)
		fmt.Fprintf(w, "\tSCALEEMITTERWITHACTOR %d\n", e.ScaleEmitterWithActor)
		fmt.Fprintln(w)
	}
	return nil
}

func (e *EmitterDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("NAME", 1)
	if err != nil {
		return err
	}
	e.Name = records[1]

	records, err = token.ReadProperty("TEXTURE", 1)
	if err != nil {
		return err
	}
	e.Texture = records[1]

	read := func(name string, count int, dst interface{}) error {
		records, err := token.ReadProperty(name, count)
		if err != nil {
			return err
		}
		if err := parse(dst, records[1:]...); err != nil {
			return fmt.Errorf("%s: %w", strings.ToLower(name), err)
		}
		return nil
	}

	properties := []struct {
		name  string
		count int
		dst   interface{}
	}{
		{"RELATIVETOBONE", 1, &e.RelativeToBone},
		{"NOOCCLUSION", 1, &e.NoOcclusion},
		{"ADDITIVEBLENDING", 1, &e.AdditiveBlending},
		{"SCALEWITHACTOR", 1, &e.ScaleWithActor},
		{"SPRITEORIENTATION", 1, &e.SpriteOrientation},
		{"STICKTOACTOR", 1, &e.StickToActor},
		{"DEFAULTLIFESPAN", 1, &e.DefaultLifeSpan},
		{"PARTICLELIFESPAN", 1, &e.ParticleLifeSpan},
		{"PARTICLESATCREATION", 1, &e.ParticlesAtCreation},
		{"PARTICLESATINTERVAL", 1, &e.ParticlesAtInterval},
		{"INTERVALSPERSECOND", 1, &e.IntervalsPerSecond},
		{"SPAWNDELAY", 1, &e.SpawnDelay},
		{"FADEINTIME", 1, &e.FadeInTime},
		{"FADEOUTTIME", 1, &e.FadeOutTime},
		{"SCALEINTIME", 1, &e.ScaleInTime},
		{"SCALEOUTTIME", 1, &e.ScaleOutTime},
		{"REDUCTIONDISTANCE", 1, &e.ReductionDistance},
		{"MAXALPHA", 1, &e.MaxAlpha},
		{"SPAWNSHAPE", 1, &e.SpawnShape},
		{"SHAPERADIUS", 1, &e.ShapeRadius},
		{"SHAPERADIUSMINOR", 1, &e.ShapeRadiusMinor},
		{"SHAPEHEIGHT", 1, &e.ShapeHeight},
		{"SHAPEOFFSET", 3, &e.ShapeOffset},
		{"SHAPETILT", 2, &e.ShapeTilt},
		{"PARTICLEWIDTHMIN", 1, &e.ParticleWidthMin},
		{"PARTICLEZBIAS", 1, &e.ParticleZBias},
		{"TINTSTART", 3, &e.TintStart},
		{"TINTEND", 3, &e.TintEnd},
		{"SPEEDMIN", 3, &e.SpeedMin},
		{"SPEEDMAX", 3, &e.SpeedMax},
		{"ACCELERATION", 3, &e.Acceleration},
		{"OUTWARDSPEEDMIN", 1, &e.OutwardSpeedMin},
		{"OUTWARDSPEEDMAX", 1, &e.OutwardSpeedMax},
		{"OUTWARDSPEEDACCELERATION", 1, &e.OutwardSpeedAcceleration},
		{"ORBITALSPEEDMIN", 1, &e.OrbitalSpeedMin},
		{"ORBITALSPEEDMAX", 1, &e.OrbitalSpeedMax},
		{"ORBITALSPEEDACCELERATION", 1, &e.OrbitalSpeedAcceleration},
		{"SCALARGRAVITY", 1, &e.ScalarGravity},
		{"WINDSPEED", 1, &e.WindSpeed},
		{"ANIMATIONFRAMES", 1, &e.AnimationFrames},
		{"ANIMATIONRATE", 1, &e.AnimationRate},
		{"PARTICLESPINRATE", 1, &e.ParticleSpinRate},
		{"OLDPARTICLETYPE", 1, &e.OldParticleType},
		{"OLDFLAGS", 1, &e.OldFlags},
		{"OLDSIZE", 1, &e.OldSize},
		{"GRAVITY", 3, &e.Gravity},
		{"BBMIN", 3, &e.BBMin},
		{"BBMAX", 3, &e.BBMax},
		{"SPAWNSCALE", 1, &e.SpawnScale},
		{"ALPHA", 1, &e.Alpha},
		{"RANDOMROTATION", 1, &e.RandomRotation},
		{"PARTICLEORIENTATION", 1, &e.ParticleOrientation},
		{"PARTICLEHEIGHTMIN", 1, &e.ParticleHeightMin},
		{"PARTICLEHEIGHTMAX", 1, &e.ParticleHeightMax},
		{"PARTICLEWIDTHMAX", 1, &e.ParticleWidthMax},
		{"PARTICLESPINRATEMAX", 1, &e.ParticleSpinRateMax},
		{"PROPORTIONALSIZESCALING", 1, &e.ProportionalSizeScaling},
		{"HEIGHTSQUASHTIME", 1, &e.HeightSquashTime},
		{"WIDTHSQUASHTIME", 1, &e.WidthSquashTime},
		{"ALLOWCENTERPASSTHROUGH", 1, &e.AllowCenterPassThrough},
		{"SCALEEMITTERWITHACTOR", 1, &e.ScaleEmitterWithActor},
	}

	for _, property := range properties {
		if err := read(property.name, property.count, property.dst); err != nil {
			return err
		}
	}
	return nil
}

func (e *EmitterDef) ToRaw(_ *Wce, dst *raw.EddEntry) error {
	if e == nil || dst == nil {
		return fmt.Errorf("nil receiver or destination")
	}

	dst.EmitterName = e.Name
	dst.TextureName = e.Texture
	dst.RelativeToBone = e.RelativeToBone
	dst.NoOcclusion = e.NoOcclusion
	dst.AdditiveBlending = e.AdditiveBlending
	dst.ScaleWithActor = e.ScaleWithActor
	dst.SpriteOrientation = e.SpriteOrientation
	dst.StickToActor = e.StickToActor
	dst.DefaultLifeSpan = e.DefaultLifeSpan
	dst.ParticleLifeSpan = e.ParticleLifeSpan
	dst.ParticlesAtCreation = e.ParticlesAtCreation
	dst.ParticlesAtInterval = e.ParticlesAtInterval
	dst.IntervalsPerSecond = e.IntervalsPerSecond
	dst.SpawnDelay = e.SpawnDelay
	dst.FadeInTime = e.FadeInTime
	dst.FadeOutTime = e.FadeOutTime
	dst.ScaleInTime = e.ScaleInTime
	dst.ScaleOutTime = e.ScaleOutTime
	dst.ReductionDistance = e.ReductionDistance
	dst.MaxAlpha = e.MaxAlpha
	dst.SpawnShape = e.SpawnShape
	dst.ShapeRadius = e.ShapeRadius
	dst.ShapeRadiusMinor = e.ShapeRadiusMinor
	dst.ShapeHeight = e.ShapeHeight
	dst.ShapeOffset = e.ShapeOffset
	dst.ShapeTilt = e.ShapeTilt
	dst.ParticleWidthMin = e.ParticleWidthMin
	dst.ParticleZBias = e.ParticleZBias
	dst.TintStartR = e.TintStart[0]
	dst.TintStartG = e.TintStart[1]
	dst.TintStartB = e.TintStart[2]
	dst.TintEndR = e.TintEnd[0]
	dst.TintEndG = e.TintEnd[1]
	dst.TintEndB = e.TintEnd[2]
	dst.SpeedMin = e.SpeedMin
	dst.SpeedMax = e.SpeedMax
	dst.Acceleration = e.Acceleration
	dst.OutwardSpeedMin = e.OutwardSpeedMin
	dst.OutwardSpeedMax = e.OutwardSpeedMax
	dst.OutwardSpeedAcceleration = e.OutwardSpeedAcceleration
	dst.OrbitalSpeedMin = e.OrbitalSpeedMin
	dst.OrbitalSpeedMax = e.OrbitalSpeedMax
	dst.OrbitalSpeedAcceleration = e.OrbitalSpeedAcceleration
	dst.ScalarGravity = e.ScalarGravity
	dst.WindSpeed = e.WindSpeed
	dst.AnimationFrames = e.AnimationFrames
	dst.AnimationRate = e.AnimationRate
	dst.ParticleSpinRate = e.ParticleSpinRate
	dst.OldParticleType = e.OldParticleType
	dst.OldFlags = e.OldFlags
	dst.OldSize = e.OldSize
	dst.Gravity = e.Gravity
	dst.BBMin = e.BBMin
	dst.BBMax = e.BBMax
	dst.SpawnScale = e.SpawnScale
	dst.Alpha = e.Alpha
	dst.RandomRotation = e.RandomRotation
	dst.ParticleOrientation = e.ParticleOrientation
	dst.ParticleHeightMin = e.ParticleHeightMin
	dst.ParticleHeightMax = e.ParticleHeightMax
	dst.ParticleWidthMax = e.ParticleWidthMax
	dst.ParticleSpinRateMax = e.ParticleSpinRateMax
	dst.ProportionalSizeScaling = e.ProportionalSizeScaling
	dst.HeightSquashTime = e.HeightSquashTime
	dst.WidthSquashTime = e.WidthSquashTime
	dst.AllowCenterPassThrough = e.AllowCenterPassThrough
	dst.ScaleEmitterWithActor = e.ScaleEmitterWithActor
	return nil
}

func (e *EmitterDef) FromRaw(_ *Wce, src *raw.EddEntry) error {
	if e == nil || src == nil {
		return fmt.Errorf("nil receiver or source")
	}

	e.Name = src.EmitterName
	e.Texture = src.TextureName
	e.RelativeToBone = src.RelativeToBone
	e.NoOcclusion = src.NoOcclusion
	e.AdditiveBlending = src.AdditiveBlending
	e.ScaleWithActor = src.ScaleWithActor
	e.SpriteOrientation = src.SpriteOrientation
	e.StickToActor = src.StickToActor
	e.DefaultLifeSpan = src.DefaultLifeSpan
	e.ParticleLifeSpan = src.ParticleLifeSpan
	e.ParticlesAtCreation = src.ParticlesAtCreation
	e.ParticlesAtInterval = src.ParticlesAtInterval
	e.IntervalsPerSecond = src.IntervalsPerSecond
	e.SpawnDelay = src.SpawnDelay
	e.FadeInTime = src.FadeInTime
	e.FadeOutTime = src.FadeOutTime
	e.ScaleInTime = src.ScaleInTime
	e.ScaleOutTime = src.ScaleOutTime
	e.ReductionDistance = src.ReductionDistance
	e.MaxAlpha = src.MaxAlpha
	e.SpawnShape = src.SpawnShape
	e.ShapeRadius = src.ShapeRadius
	e.ShapeRadiusMinor = src.ShapeRadiusMinor
	e.ShapeHeight = src.ShapeHeight
	e.ShapeOffset = src.ShapeOffset
	e.ShapeTilt = src.ShapeTilt
	e.ParticleWidthMin = src.ParticleWidthMin
	e.ParticleZBias = src.ParticleZBias
	e.TintStart = [3]uint32{src.TintStartR, src.TintStartG, src.TintStartB}
	e.TintEnd = [3]uint32{src.TintEndR, src.TintEndG, src.TintEndB}
	e.SpeedMin = src.SpeedMin
	e.SpeedMax = src.SpeedMax
	e.Acceleration = src.Acceleration
	e.OutwardSpeedMin = src.OutwardSpeedMin
	e.OutwardSpeedMax = src.OutwardSpeedMax
	e.OutwardSpeedAcceleration = src.OutwardSpeedAcceleration
	e.OrbitalSpeedMin = src.OrbitalSpeedMin
	e.OrbitalSpeedMax = src.OrbitalSpeedMax
	e.OrbitalSpeedAcceleration = src.OrbitalSpeedAcceleration
	e.ScalarGravity = src.ScalarGravity
	e.WindSpeed = src.WindSpeed
	e.AnimationFrames = src.AnimationFrames
	e.AnimationRate = src.AnimationRate
	e.ParticleSpinRate = src.ParticleSpinRate
	e.OldParticleType = src.OldParticleType
	e.OldFlags = src.OldFlags
	e.OldSize = src.OldSize
	e.Gravity = src.Gravity
	e.BBMin = src.BBMin
	e.BBMax = src.BBMax
	e.SpawnScale = src.SpawnScale
	e.Alpha = src.Alpha
	e.RandomRotation = src.RandomRotation
	e.ParticleOrientation = src.ParticleOrientation
	e.ParticleHeightMin = src.ParticleHeightMin
	e.ParticleHeightMax = src.ParticleHeightMax
	e.ParticleWidthMax = src.ParticleWidthMax
	e.ParticleSpinRateMax = src.ParticleSpinRateMax
	e.ProportionalSizeScaling = src.ProportionalSizeScaling
	e.HeightSquashTime = src.HeightSquashTime
	e.WidthSquashTime = src.WidthSquashTime
	e.AllowCenterPassThrough = src.AllowCenterPassThrough
	e.ScaleEmitterWithActor = src.ScaleEmitterWithActor
	return nil
}
