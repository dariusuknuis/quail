package mod

import (
	"fmt"
	"os"

	"github.com/xackery/quail/obj"
)

func (e *MOD) ObjImport(objPath string, mtlPath string, matPath string) error {
	var err error
	rm, err := os.Open(mtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()
	req := &obj.ObjRequest{
		ObjPath:    objPath,
		MtlPath:    mtlPath,
		MattxtPath: matPath,
	}
	err = obj.Import(req)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}
	e.materials = req.Obj.Materials
	e.triangles = req.Obj.Triangles
	e.vertices = req.Obj.Vertices

	return nil
}
