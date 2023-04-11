package mod

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model/geo"
)

// BlenderImport imports a blender structure to MOD
func (e *MOD) BlenderImport(dir string) error {
	e.version = 1
	e.materials = []*geo.Material{}
	e.vertices = []*geo.Vertex{}
	e.triangles = []*geo.Triangle{}
	e.bones = []*geo.Bone{}
	e.particleRenders = []*geo.ParticleRender{}
	e.particlePoints = []*geo.ParticlePoint{}
	path := dir

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat _%s: %w", e.Name(), err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	var r *os.File
	var scanner *bufio.Scanner
	var lineNumber int

	curPath := fmt.Sprintf("%s/info.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "=")
			if len(parts) < 2 {
				return fmt.Errorf("invalid version.txt (expected 2 records) line %d: %s", lineNumber, line)
			}
			switch parts[0] {
			case "version":
				e.version = helper.AtoU32(parts[1])
			default:
				return fmt.Errorf("invalid info.txt line %d: %s", lineNumber, line)
			}
		}
	}

	curPath = fmt.Sprintf("%s/material.txt", path)
	if helper.IsFile(curPath) {

		r, err := os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner := bufio.NewScanner(r)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 3 {
				return fmt.Errorf("invalid material.txt (expected 3 records) line %d: %s", lineNumber, line)
			}
			material := &geo.Material{
				Name:       parts[0],
				Flag:       helper.AtoU32(parts[1]),
				ShaderName: parts[2],
			}
			e.materials = append(e.materials, material)
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/material_property.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 4 {
				return fmt.Errorf("invalid material_property.txt (expected 4 records) line %d: %s", lineNumber, line)
			}
			isFound := false
			for _, material := range e.materials {
				if material.Name != parts[0] {
					continue
				}
				isFound = true
				material.Properties = append(material.Properties, &geo.Property{
					Name:     parts[1],
					Value:    parts[2],
					Category: helper.AtoU32(parts[3]),
				})
				// TODO: validate value/category of material property
			}
			if !isFound {
				return fmt.Errorf("material_property.txt material not found: %s", parts[0])
			}
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/particle_point.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 5 {
				return fmt.Errorf("invalid particle_point.txt (expected 5 records) line %d: %s", lineNumber, line)
			}
			e.particlePoints = append(e.particlePoints, &geo.ParticlePoint{
				Name:        parts[0],
				Bone:        parts[1],
				Translation: geo.AtoVector3(parts[2]),
				Rotation:    geo.AtoVector3(parts[3]),
				Scale:       geo.AtoVector3(parts[4]),
			})
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/particle_render.txt", path)
	if helper.IsFile(curPath) {

		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 7 {
				return fmt.Errorf("invalid particle_render.txt (expected 7 records) line %d: %s", lineNumber, line)
			}
			e.particleRenders = append(e.particleRenders, &geo.ParticleRender{
				Duration:      helper.AtoU32(parts[0]),
				ID:            helper.AtoU32(parts[1]),
				ID2:           helper.AtoU32(parts[2]),
				ParticlePoint: parts[3],
				//UnknownA:        helper.AtoU32(parts[4]),
				UnknownB: helper.AtoU32(parts[5]),
				//UnknownFFFFFFFF: helper.AtoU32(parts[6]),
			})
			return fmt.Errorf("todo: blender import fix for particles")
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/triangle.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 3 {
				return fmt.Errorf("invalid triangle.txt (expected 3 records) line %d: %s", lineNumber, line)
			}
			e.triangles = append(e.triangles, &geo.Triangle{
				Index:        geo.AtoUIndex3(parts[0]),
				Flag:         helper.AtoU32(parts[1]),
				MaterialName: parts[2],
			})
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/vertex.txt", path)
	if helper.IsFile(curPath) {

		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 5 {
				return fmt.Errorf("invalid vertex.txt (expected 5 records) line %d: %s", lineNumber, line)
			}
			vert := &geo.Vertex{
				Position: geo.AtoVector3(parts[0]),
				Normal:   geo.AtoVector3(parts[1]),
				Uv:       geo.AtoVector2(parts[2]),
				Uv2:      geo.AtoVector2(parts[3]),
				Tint:     geo.AtoRGBA(parts[4]),
			}
			vert.Position = &geo.Vector3{X: vert.Position.Y, Y: -vert.Position.X, Z: vert.Position.Z}
			vert.Normal = &geo.Vector3{X: vert.Normal.Y, Y: -vert.Normal.X, Z: vert.Normal.Z}

			e.vertices = append(e.vertices, vert)
		}
		r.Close()
	}

	curPath = fmt.Sprintf("%s/bone.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 7 {

				return fmt.Errorf("invalid bone.txt (expected 7 records) line %d: %s", lineNumber, line)
			}
			e.bones = append(e.bones, &geo.Bone{
				Name:          parts[0],
				ChildIndex:    helper.AtoI32(parts[1]),
				ChildrenCount: helper.AtoU32(parts[2]),
				Next:          helper.AtoI32(parts[3]),
				Pivot:         geo.AtoVector3(parts[4]),
				Rotation:      geo.AtoQuad4(parts[5]),
				Scale:         geo.AtoVector3(parts[6]),
			})
		}
		r.Close()
	}

	return nil
}
