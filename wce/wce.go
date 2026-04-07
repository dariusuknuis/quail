// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wce

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/qfs"
	"github.com/xackery/quail/raw"
)

var AsciiVersion = "v0.0.1"

// Wce is a struct representing a Wce file
type Wce struct {
	FileSystem             qfs.QFS
	isVariationMaterial    bool   // set true while writing or reading variations
	lastReadFolder         string // used during wce parsing to remember context
	isObj                  bool   // true when a _obj suffix is found in path
	isChr                  bool   // true when a _chr suffix is found in path
	maxMaterialHeads       map[string]int
	maxMaterialTextures    map[string]int
	indexedTags            map[string]int32 // used when parsing to keep track of indexes
	fragToIndexedTags      map[int32]string // used for reverse lookup of above
	FileName               string
	WorldDef               *WorldDef
	GlobalAmbientLightDef  *GlobalAmbientLightDef
	Version                uint32
	ActorDefs              []*ActorDef
	ActorInsts             []*ActorInst
	AmbientLights          []*AmbientLight
	BlitSpriteDefs         []*BlitSpriteDef
	DefaultPalette         *DefaultPalette
	UserDatas              []*UserData
	DMSpriteDef2s          []*DMSpriteDef2
	DMSpriteDefs           []*DMSpriteDef
	DMTrackDef2s           []*DMTrackDef2
	HierarchicalSpriteDefs []*HierarchicalSpriteDef
	LightDefs              []*LightDef
	MaterialDefs           []*MaterialDef
	MaterialPalettes       []*MaterialPalette
	ParticleCloudDefs      []*ParticleCloudDef
	PointLights            []*PointLight
	DirectionalLights      []*DirectionalLight
	PolyhedronDefs         []*PolyhedronDefinition
	Regions                []*Region
	RGBTrackDefs           []*RGBTrackDef
	SimpleSpriteDefs       []*SimpleSpriteDef
	Sprite2DDefs           []*Sprite2DDef
	Sprite3DDefs           []*Sprite3DDef
	SphereListDefs         []*SphereListDef
	TrackDefs              []*TrackDef
	TrackInstances         []*TrackInstance
	variationMaterialDefs  map[string][]*MaterialDef
	WorldTrees             []*WorldTree
	Zones                  []*Zone
	AniDefs                []*EqgAniDef
	MdsDefs                []*EqgMdsDef
	ModDefs                []*EqgModDef
	TerDefs                []*EqgTerDef
	LayDefs                []*EqgLayDef
	PtsDefs                []*EqgParticlePointDef
	PrtDefs                []*EqgParticleRenderDef
	LodDefs                []*EqgLodDef
	ZonDefs                []*EqgZonDef
	EffectOlds             []*EffectOld
}

type WldDefinitioner interface {
	Definition() string
	ToRaw(src *Wce, dst *raw.Wld) (int32, error)
	Write(token *AsciiWriteToken) error
}

func New(filename string) *Wce {
	isObj := strings.Contains(filename, "_obj")
	isChr := strings.Contains(filename, "_chr")

	return &Wce{
		FileSystem:            qfs.OSFS{},
		FileName:              filename,
		isObj:                 isObj,
		isChr:                 isChr,
		maxMaterialHeads:      make(map[string]int),
		maxMaterialTextures:   make(map[string]int),
		variationMaterialDefs: make(map[string][]*MaterialDef),
		WorldDef:              &WorldDef{folders: []string{"world"}},
	}
}

func baseTag(tag string) string {
	// detect ".###" at end
	if len(tag) > 4 && tag[len(tag)-4] == '.' {
		suffix := tag[len(tag)-3:]
		if isNumeric(suffix) {
			return tag[:len(tag)-4]
		}
	}
	return tag
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ByTag returns a instance by tag
func (wce *Wce) ByTag(tag string) WldDefinitioner {
	if tag == "" {
		return nil
	}
	if strings.HasSuffix(baseTag(tag), "_SPRITE") || strings.HasPrefix(baseTag(tag), "I_") {
		for _, sprite := range wce.SimpleSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
		for _, sprite := range wce.BlitSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_PCD") {
		for _, cloud := range wce.ParticleCloudDefs {
			if cloud.Tag == tag {
				return cloud
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_SPB") {
		for _, sprite := range wce.BlitSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_MDF") {
		for _, material := range wce.MaterialDefs {
			if material.Tag == tag {
				return material
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_MP") {
		for _, palette := range wce.MaterialPalettes {
			if palette.Tag == tag {
				return palette
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_DMSPRITEDEF") {
		for _, sprite := range wce.DMSpriteDef2s {
			if sprite.Tag == tag {
				return sprite
			}
		}
		for _, sprite := range wce.DMSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_DMTRACKDEF") {
		for _, track := range wce.DMTrackDef2s {
			if track.Tag == tag {
				return track
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_LIGHTDEF") {
		for _, light := range wce.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}
	if strings.HasSuffix(baseTag(tag), "_LDEF") {
		for _, light := range wce.LightDefs {
			if light.Tag == tag {
				return light
			}
		}
	}

	if strings.HasSuffix(baseTag(tag), "_TRACKDEF") {
		for _, track := range wce.TrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	if strings.HasSuffix(baseTag(tag), "_HS_DEF") {
		for _, sprite := range wce.HierarchicalSpriteDefs {
			if sprite.Tag == tag {
				return sprite
			}
		}
	}

	if strings.HasSuffix(baseTag(tag), "_POLYHDEF") {
		for _, polyhedron := range wce.PolyhedronDefs {
			if polyhedron.Tag == tag {
				return polyhedron
			}
		}
	}

	if strings.HasSuffix(baseTag(tag), "_DMT") {
		for _, track := range wce.RGBTrackDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	if strings.HasSuffix(baseTag(tag), "_SPHRLDEF") {
		for _, track := range wce.SphereListDefs {
			if track.Tag == tag {
				return track
			}
		}
	}

	for _, sprite := range wce.Sprite3DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}
	for _, region := range wce.Regions {
		if region.Tag == tag {
			return region
		}
	}

	for _, actor := range wce.ActorDefs {
		if actor.Tag == tag {
			return actor
		}
	}

	for _, track := range wce.TrackInstances {
		if track.Tag == tag {
			return track
		}
	}

	for _, sprite := range wce.Sprite2DDefs {
		if sprite.Tag == tag {
			return sprite
		}
	}

	// for _, sprite := range wce.SimpleSpriteDefs {
	// 	if sprite.Tag == tag {
	// 		return sprite
	// 	}
	// 	if strings.HasSuffix(sprite.Tag, "_SPRITE") && !strings.HasSuffix(tag, "_SPRITE") {
	// 		if sprite.Tag == tag+"_SPRITE" {
	// 			return sprite
	// 		}
	// 	}
	// }
	return nil
}

// ByTagWithIndex returns a instance by tag with index included
// func (wce *Wce) ByTagWithIndex(tag string, index int) WldDefinitioner {
// 	if tag == "" {
// 		return nil
// 	}

// 	if strings.HasSuffix(tag, "_DMSPRITEDEF") {
// 		for _, dmsprite := range wce.DMSpriteDef2s {
// 			if dmsprite.Tag == tag && dmsprite.TagIndex == index {
// 				return dmsprite
// 			}
// 		}
// 		for _, dmsprite := range wce.DMSpriteDefs {
// 			if dmsprite.Tag == tag && dmsprite.TagIndex == index {
// 				return dmsprite
// 			}
// 		}
// 	}

// 	if strings.HasSuffix(tag, "_TRACK") {
// 		for _, track := range wce.TrackInstances {
// 			if track.Tag == tag && track.TagIndex == index {
// 				return track
// 			}
// 		}
// 	}

// 	if strings.HasSuffix(tag, "_TRACKDEF") {
// 		for _, track := range wce.TrackDefs {
// 			if track.Tag == tag && track.TagIndex == index {
// 				return track
// 			}
// 		}
// 	}

// 	if strings.HasSuffix(tag, "_MDF") {
// 		for _, material := range wce.MaterialDefs {
// 			if material.Tag == tag && material.TagIndex == index {
// 				return material
// 			}
// 		}
// 	}

// 	if strings.HasSuffix(tag, "_SPRITE") {
// 		for _, sprite := range wce.SimpleSpriteDefs {
// 			if sprite.Tag == tag && sprite.TagIndex == index {
// 				return sprite
// 			}
// 		}
// 	}

// 	if strings.HasSuffix(tag, "_PCD") {
// 		for _, sprite := range wce.ParticleCloudDefs {
// 			if sprite.Tag == tag && sprite.TagIndex == index {
// 				return sprite
// 			}
// 		}
// 	}

// 	return nil
// }

// NextTagIndex returns the next available index for a tag
func (wce *Wce) NextIndexedTag(tag string, fragID int32) string {
	if tag == "" {
		return ""
	}

	// If base doesn't exist → use it
	if _, exists := wce.indexedTags[tag]; !exists {
		wce.indexedTags[tag] = fragID
		wce.fragToIndexedTags[fragID] = tag
		return tag
	}

	// Otherwise, find next available .###
	for i := 1; ; i++ {
		newTag := fmt.Sprintf("%s.%03d", tag, i)
		if _, exists := wce.indexedTags[newTag]; !exists {
			wce.indexedTags[newTag] = fragID
			wce.fragToIndexedTags[fragID] = newTag
			return newTag
		}
	}
}

func (wce *Wce) reset() {
	wce.GlobalAmbientLightDef = nil
	wce.lastReadFolder = ""
	wce.indexedTags = make(map[string]int32)
	wce.fragToIndexedTags = make(map[int32]string)
	wce.SimpleSpriteDefs = []*SimpleSpriteDef{}
	wce.MaterialDefs = []*MaterialDef{}
	wce.variationMaterialDefs = make(map[string][]*MaterialDef)
	wce.MaterialPalettes = []*MaterialPalette{}
	wce.DefaultPalette = nil
	wce.UserDatas = []*UserData{}
	wce.DMSpriteDefs = []*DMSpriteDef{}
	wce.DMSpriteDef2s = []*DMSpriteDef2{}
	wce.ActorDefs = []*ActorDef{}
	wce.ActorInsts = []*ActorInst{}
	wce.LightDefs = []*LightDef{}
	wce.PointLights = []*PointLight{}
	wce.DirectionalLights = []*DirectionalLight{}
	wce.Sprite3DDefs = []*Sprite3DDef{}
	wce.SphereListDefs = []*SphereListDef{}
	wce.TrackInstances = []*TrackInstance{}
	wce.TrackDefs = []*TrackDef{}
	wce.HierarchicalSpriteDefs = []*HierarchicalSpriteDef{}
	wce.PolyhedronDefs = []*PolyhedronDefinition{}
	wce.WorldTrees = []*WorldTree{}
	wce.Regions = []*Region{}
	wce.AmbientLights = []*AmbientLight{}
	wce.Zones = []*Zone{}
	wce.RGBTrackDefs = []*RGBTrackDef{}
	wce.ParticleCloudDefs = []*ParticleCloudDef{}
	wce.Sprite2DDefs = []*Sprite2DDef{}
	wce.MdsDefs = []*EqgMdsDef{}
	wce.ModDefs = []*EqgModDef{}
	wce.TerDefs = []*EqgTerDef{}
	wce.LayDefs = []*EqgLayDef{}
	wce.PtsDefs = []*EqgParticlePointDef{}
	wce.PrtDefs = []*EqgParticleRenderDef{}
	wce.LodDefs = []*EqgLodDef{}
	wce.ZonDefs = []*EqgZonDef{}
	wce.EffectOlds = []*EffectOld{}
}
