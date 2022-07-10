package obj

import (
	"fmt"
	"os"
)

func objExport(req *ObjRequest) error {

	if req == nil {
		return fmt.Errorf("request is nil")
	}
	if req.Data == nil {
		return fmt.Errorf("request object is nil")
	}
	if req.Data.Name == "" {
		req.Data.Name = "unnamed"
	}
	w, err := os.Create(req.ObjPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.WriteString("# exported by quail\n\n")
	if err != nil {
		return fmt.Errorf("export header: %w", err)
	}

	_, err = w.WriteString(fmt.Sprintf("mtllib %s.mtl\no %s\n", req.Data.Name, req.Data.Name))
	if err != nil {
		return fmt.Errorf("mtllib: %w", err)
	}

	for _, e := range req.Data.Vertices {
		fmt.Println(e)
		_, err = w.WriteString(fmt.Sprintf("v %0.6f %0.6f %0.6f\n", e.Position.X, e.Position.Y, e.Position.Z))
		if err != nil {
			return fmt.Errorf("export pos: %w", err)
		}
	}
	for _, e := range req.Data.Vertices {
		_, err = w.WriteString(fmt.Sprintf("vt %0.6f %0.6f\n", e.Uv.X, e.Uv.Y))
		if err != nil {
			return fmt.Errorf("export uv: %w", err)
		}
	}
	for _, e := range req.Data.Vertices {
		_, err = w.WriteString(fmt.Sprintf("vn %0.6f %0.6f %0.6f\n", e.Normal.X, e.Normal.Y, e.Normal.Z))
		if err != nil {
			return fmt.Errorf("export normal: %w", err)
		}
	}

	lastMaterial := ""
	group := 0
	for _, e := range req.Data.Triangles {
		if lastMaterial != e.MaterialName {
			lastMaterial = e.MaterialName
			_, err = w.WriteString(fmt.Sprintf("usemtl %s\ns off\ng piece%d\n", e.MaterialName, group))
			if err != nil {
				return fmt.Errorf("usemtl: %w", err)
			}
			group++
		}
		_, err = w.WriteString(fmt.Sprintf("f %d/%d/%d %d/%d/%d %d/%d/%d\n", int(e.Index[0]+1), int(e.Index[0]+1), int(e.Index[0]+1), int(e.Index[1]+1), int(e.Index[1]+1), int(e.Index[1]+1), int(e.Index[2]+1), int(e.Index[2]+1), int(e.Index[2]+1)))
		if err != nil {
			return fmt.Errorf("f: %w", err)
		}
	}

	return nil
}
