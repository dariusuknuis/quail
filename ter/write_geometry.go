package ter

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/quail/helper"
)

func (e *TER) writeGeometry() ([]byte, []byte, error) {
	var err error

	names := []string{}
	nameBuf := bytes.NewBuffer(nil)
	dataBuf := bytes.NewBuffer(nil)
	nameID := -1
	// materials

	for materialID, o := range e.materials {
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(materialID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material id %s: %w", o.name, err)
		}

		nameID = -1
		for i, name := range names {
			if name == o.name {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			names = append(names, o.name)
			nameID = len(names) - 1
			err = helper.WriteString(nameBuf, o.name)
			if err != nil {
				return nil, nil, fmt.Errorf("writestring material %s: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material name id %s: %w", o.name, err)
		}

		nameID = -1
		for i, name := range names {
			if name == o.shaderName {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			names = append(names, o.shaderName)
			nameID = len(names) - 1
			err = helper.WriteString(nameBuf, o.shaderName)
			if err != nil {
				return nil, nil, fmt.Errorf("writestring material %s shader: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write material shader id %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(o.properties)))
		if err != nil {
			return nil, nil, fmt.Errorf("write material property count %s: %w", o.name, err)
		}

		for propertyID, p := range o.properties {
			nameID = -1
			for i, name := range names {
				if name == p.name {
					nameID = i
					break
				}
			}
			if nameID == -1 {
				names = append(names, p.name)
				nameID = len(names) - 1
				err = helper.WriteString(nameBuf, p.name)
				if err != nil {
					return nil, nil, fmt.Errorf("writestring material %s property %s: %w", o.name, p.name, err)
				}
			}

			err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
			if err != nil {
				return nil, nil, fmt.Errorf("write material %s property %s id %d: %w", o.name, p.name, propertyID, err)
			}

			if p.typeValue == 0 {
				err = binary.Write(dataBuf, binary.LittleEndian, p.floatValue)
				if err != nil {
					return nil, nil, fmt.Errorf("write material %s property %s id %d value (float): %w", o.name, p.name, propertyID, err)
				}
			} else {
				err = binary.Write(dataBuf, binary.LittleEndian, p.intValue)
				if err != nil {
					return nil, nil, fmt.Errorf("write material %s property %s id %d value (int): %w", o.name, p.name, propertyID, err)
				}
			}
		}
	}

	// verts
	for i, o := range e.vertices {
		err = binary.Write(dataBuf, binary.LittleEndian, o.position)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d position: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.position2)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d position2: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.uv)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d uv: %w", i, err)
		}
	}

	// triangles
	for i, o := range e.triangles {
		nameID = -1
		for i, name := range names {
			if name == o.materialName {
				nameID = i
				break
			}
		}
		if nameID == -1 {
			return nil, nil, fmt.Errorf("triangle %d refers to material %s, which is not declared", i, o.materialName)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.index)
		if err != nil {
			return nil, nil, fmt.Errorf("write triangle %d index: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(nameID))
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d position2: %w", i, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.flag)
		if err != nil {
			return nil, nil, fmt.Errorf("write vertex %d flag: %w", i, err)
		}
	}
	return nameBuf.Bytes(), dataBuf.Bytes(), nil
}
