package world

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"mycraft/block"
)

type DemoWorld struct {
	Center            math32.Vector2
	RenderingDistance float32 // block
	ContainerNode     *core.Node
	Meshes            map[*graphic.Mesh]math32.Vector3
	Blocks            *block.Repository
	LatestUpdate      math32.Vector3
	Initialized       bool
}

func NewDemoWorld(container *core.Node, renderingDistance float32, repository *block.Repository) *DemoWorld {
	demo := DemoWorld{
		Center:            math32.Vector2{X: 0, Y: 0},
		ContainerNode:     container,
		RenderingDistance: renderingDistance,
		Meshes:            make(map[*graphic.Mesh]math32.Vector3),
		Blocks:            repository,
		LatestUpdate:      math32.Vector3{X: 0, Y: 0, Z: 0},
		Initialized:       false,
	}

	return &demo
}

func (dw *DemoWorld) Update(pos math32.Vector3) {
	pos = dw.getWorldPosition(pos)

	if dw.Initialized && dw.LatestUpdate.X == pos.X && dw.LatestUpdate.Z == pos.Z {
		return
	}

	dw.Initialized = true
	dw.LatestUpdate = pos

	dw.cleanTooFar(pos)
	dw.populate(pos)
}

func (dw *DemoWorld) populate(pos math32.Vector3) {
	for x := pos.X - dw.RenderingDistance; x <= pos.X+dw.RenderingDistance; x++ {
		for z := pos.Z - dw.RenderingDistance; z <= pos.Z+dw.RenderingDistance; z++ {
			meshPos := math32.Vector3{X: x, Y: -2, Z: z}
			if !dw.hasMeshAt(meshPos) {
				dw.populateMeshAt(meshPos)
			}
		}
	}
}

func (dw *DemoWorld) populateMeshAt(pos math32.Vector3) {
	mesh := dw.Blocks.Get("green_grass").CreateMesh()
	mesh.SetPositionVec(&pos)
	dw.ContainerNode.Add(mesh)
	dw.Meshes[mesh] = pos
}

func (dw *DemoWorld) hasMeshAt(pos math32.Vector3) bool {
	for _, meshPos := range dw.Meshes {
		if meshPos.X == pos.X && meshPos.Z == pos.Z {
			return true
		}
	}

	return false
}

func (dw *DemoWorld) getWorldPosition(pos math32.Vector3) math32.Vector3 {
	return math32.Vector3{
		X: float32(int(pos.X)),
		Y: float32(int(pos.Y)),
		Z: float32(int(pos.Z)),
	}
}

func (dw *DemoWorld) cleanTooFar(pos math32.Vector3) {
	for mesh, meshPos := range dw.Meshes {
		dist := math32.Max(math32.Abs(meshPos.X-pos.X), math32.Abs(meshPos.Z-pos.Z))

		if dist > dw.RenderingDistance {
			dw.ContainerNode.Remove(mesh)
			delete(dw.Meshes, mesh)
		}
	}
}
