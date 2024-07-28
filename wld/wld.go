// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
)

var AsciiVersion = "v0.0.1"

// Wld is a struct representing a Wld file
type Wld struct {
	FileName               string
	GlobalAmbientLight     string
	Version                uint32
	SimpleSpriteDefs       []*SimpleSpriteDef
	MaterialDefs           []*MaterialDef
	MaterialPalettes       []*MaterialPalette
	DMSpriteDefs           []*DMSpriteDef
	DMSpriteInsts          []*DMSprite
	DMSpriteDef2s          []*DMSpriteDef2
	ActorDefs              []*ActorDef
	ActorInsts             []*ActorInst
	LightDefs              []*LightDef
	PointLights            []*PointLight
	Sprite3DDefs           []*Sprite3DDef
	TrackInstances         []*TrackInstance
	TrackDefs              []*TrackDef
	HierarchicalSpriteDefs []*HierarchicalSpriteDef
	PolyhedronDefs         []*PolyhedronDefinition
	WorldTrees             []*WorldTree
	Regions                []*Region
	AmbientLights          []*AmbientLight
	Zones                  []*Zone
}

// ReadAscii reads the ascii file at path
func (wld *Wld) ReadAscii(path string) error {

	asciiReader, err := LoadAsciiFile(path, wld)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	fmt.Println(asciiReader.TotalLineCountRead(), "total lines parsed for", filepath.Base(path))
	return nil
}

func (wld *Wld) WriteAscii(path string, isDir bool) error {
	var err error
	//var err error

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	rootPath := path + "/_root.wce"
	if !isDir {
		rootPath = path + "/" + wld.FileName
	}

	var w *os.File
	rootBuf, err := os.Create(rootPath)
	if err != nil {
		return err
	}
	defer rootBuf.Close()
	writeAsciiHeader(rootBuf)

	// now we can write

	w = rootBuf
	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]

		baseTag := strings.ToLower(strings.TrimSuffix(strings.ToUpper(dmSprite.Tag), "_DMSPRITEDEF"))
		dmBuf, err := os.Create(path + "/" + baseTag + ".mod")
		if err != nil {
			return err
		}
		defer dmBuf.Close()
		writeAsciiHeader(dmBuf)

		w = dmBuf

		if dmSprite.MaterialPaletteTag != "" {
			isMaterialPaletteFound := false
			for _, materialPal := range wld.MaterialPalettes {
				if materialPal.Tag != dmSprite.MaterialPaletteTag {
					continue
				}

				for _, materialTag := range materialPal.Materials {
					isMaterialDefFound := false
					for _, materialDef := range wld.MaterialDefs {
						if materialDef.Tag != materialTag {
							continue
						}

						if materialDef.SimpleSpriteInstTag != "" {
							isSimpleSpriteFound := false
							for _, simpleSprite := range wld.SimpleSpriteDefs {
								if simpleSprite.Tag != materialDef.SimpleSpriteInstTag {
									continue
								}
								isSimpleSpriteFound = true
								err = simpleSprite.Write(w)
								if err != nil {
									return fmt.Errorf("simple sprite %s: %w", simpleSprite.Tag, err)
								}
								break
							}
							if !isSimpleSpriteFound {
								return fmt.Errorf("simple sprite %s not found", materialDef.SimpleSpriteInstTag)
							}
						}

						isMaterialDefFound = true
						err = materialDef.Write(w)
						if err != nil {
							return fmt.Errorf("material %s: %w", materialDef.Tag, err)
						}
						break
					}
					if !isMaterialDefFound {
						return fmt.Errorf("dmsprite %s materialdef %s not found", dmSprite.Tag, materialTag)
					}
				}

				isMaterialPaletteFound = true
				err = materialPal.Write(w)
				if err != nil {
					return fmt.Errorf("material palette %s: %w", materialPal.Tag, err)
				}
				break
			}
			if !isMaterialPaletteFound {
				return fmt.Errorf("material palette %s not found", dmSprite.MaterialPaletteTag)
			}
		}

		err = dmSprite.Write(w)
		if err != nil {
			return fmt.Errorf("dm sprite def %s: %w", dmSprite.Tag, err)
		}
		fmt.Fprintf(rootBuf, "INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag))

		tracksWritten := map[string]bool{}
		var aniBuf *os.File
		for _, hierarchySprite := range wld.HierarchicalSpriteDefs {
			isFound := false

			if hierarchySprite.DMSpriteTag == dmSprite.Tag {
				isFound = true
			}
			if !isFound {
				for _, skin := range hierarchySprite.AttachedSkins {
					if skin.DMSpriteTag != dmSprite.Tag {
						continue
					}
					isFound = true
					break
				}
			}
			if !isFound {
				for _, dag := range hierarchySprite.Dags {
					if dag.SpriteTag == dmSprite.Tag {
						isFound = true
						break
					}
				}
			}
			if !isFound {
				continue
			}

			err = hierarchySprite.Write(w)
			if err != nil {
				return fmt.Errorf("hierarchical sprite %s: %w", hierarchySprite.Tag, err)
			}

			for _, dag := range hierarchySprite.Dags {
				if dag.Track == "" {
					continue
				}
				if tracksWritten[dag.Track] {
					continue
				}

				isTrackFound := false
				for _, track := range wld.TrackInstances {
					if track.Tag != dag.Track {
						continue
					}

					if aniBuf == nil {
						aniBuf, err = os.Create(path + "/" + baseTag + ".ani")
						if err != nil {
							return err
						}
						defer aniBuf.Close()
						writeAsciiHeader(aniBuf)
					}
					isTrackDefFound := false
					for _, trackDef := range wld.TrackDefs {
						if trackDef.Tag != track.DefinitionTag {
							continue
						}
						isTrackDefFound = true

						if tracksWritten[trackDef.Tag] {
							break
						}

						err = trackDef.Write(aniBuf)
						if err != nil {
							return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
						}

						tracksWritten[trackDef.Tag] = true
						break
					}
					if !isTrackDefFound {
						return fmt.Errorf("hierarchy %s track %s definition not found", hierarchySprite.Tag, track.DefinitionTag)
					}

					isTrackFound = true

					tracksWritten[dag.Track] = true

					err = track.Write(aniBuf)
					if err != nil {
						return fmt.Errorf("track %s: %w", track.Tag, err)
					}
				}
				if !isTrackFound {
					return fmt.Errorf("hierarchy %s track %s not found", hierarchySprite.Tag, dag.Track)
				}
			}

			break
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.PolyhedronDefs); i++ {
		polyhedron := wld.PolyhedronDefs[i]
		err = polyhedron.Write(w)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", polyhedron.Tag, err)
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.HierarchicalSpriteDefs); i++ {
		hierarchicalSprite := wld.HierarchicalSpriteDefs[i]
		err = hierarchicalSprite.Write(w)
		if err != nil {
			return fmt.Errorf("hierarchical sprite def %s: %w", hierarchicalSprite.Tag, err)
		}
	}

	return nil
}

func writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n\n", common.Version)
}

func (wld *Wld) Write(w io.Writer) error {
	var err error

	raw.NameClear()

	out := &raw.Wld{
		MetaFileName: wld.FileName,
		Version:      wld.Version,
		Fragments:    []model.FragmentReadWriter{},
	}

	/*
		for i := 0; i < len(wld.SimpleSpriteDefs); i++ {
			sprite := wld.SimpleSpriteDefs[i]

			if sprite.fragID == 0 { // if sprite hasn't been inesrted yet
				bitmapRefs := []uint32{}

				for j := 0; j < len(sprite.Bitmaps); j++ {
					bitmap := sprite.Bitmaps[j]
					bmInfo := wld.bitmapByTag(bitmap)
					if bmInfo == nil {
						return fmt.Errorf("spriteInstance %s refers sprite %s which refers to bitmap %s which does not exist", sprite.Tag, sprite.Tag, bitmap)
					}
					if bmInfo.fragID > 0 {
						bitmapRefs = append(bitmapRefs, bmInfo.fragID)
						continue
					}

					nameRef := raw.NameAdd(bmInfo.Tag)
					out.Fragments = append(out.Fragments, &rawfrag.WldFragBMInfo{
						NameRef:      nameRef,
						TextureNames: bmInfo.Textures,
					})
					bitmapRefs = append(bitmapRefs, uint32(len(out.Fragments)))
				}

				nameRef := raw.NameAdd(sprite.Tag)
				out.Fragments = append(out.Fragments, &rawfrag.WldFragSimpleSpriteDef{
					NameRef:      nameRef,
					Flags:        sprite.Flags,
					CurrentFrame: sprite.CurrentFrame,
					Sleep:        sprite.Sleep,
					BitmapRefs:   bitmapRefs,
				})
			}

			nameRef := raw.NameAdd(sprite.Tag)
			out.Fragments = append(out.Fragments, &rawfrag.WldFragSimpleSprite{
				NameRef:   nameRef,
				SpriteRef: int16(sprite.fragID),
				Flags:     sprite.Flags,
			})
		}
	*/
	err = out.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func checkSharedAssets(wld *Wld) map[string]int {
	sharedMap := map[string]int{}

	// sharing pass
	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]
		sharedAdd(dmSprite.MaterialPaletteTag, sharedMap)
		sharedAdd(dmSprite.DmTrackTag, sharedMap)
		sharedAdd(dmSprite.PolyhedronTag, sharedMap)
	}

	for i := 0; i < len(wld.MaterialPalettes); i++ {
		palette := wld.MaterialPalettes[i]
		for _, material := range palette.Materials {
			sharedAdd(material, sharedMap)
		}
	}

	return sharedMap
}

func sharedAdd(tag string, sharedMap map[string]int) {
	if tag == "" {
		return
	}
	if sharedMap[tag] == 0 {
		sharedMap[tag] = 1
		return
	}
	sharedMap[tag]++
}
