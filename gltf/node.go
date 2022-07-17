package gltf

import (
	"fmt"

	"github.com/qmuntal/gltf"
)

type nodeEntry struct {
	index *uint32
	node  *gltf.Node
}

func (e *GLTF) NodeAdd(node *gltf.Node) {
	e.doc.Nodes = append(e.doc.Nodes, node)
}

func (e *GLTF) NodeSetAttributes(name string, translation [3]float32, rotation [4]float32, scale [3]float32) error {
	for _, node := range e.doc.Nodes {
		if node.Name != name {
			continue
		}
		node.Translation = translation
		node.Rotation = rotation
		node.Scale = scale
		return nil
	}
	return fmt.Errorf("node '%s' not found", name)
}
