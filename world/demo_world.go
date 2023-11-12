package world

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"mycraft/block"
)

type DemoWorld struct {
	Center            math32.Vector2
	RenderingDistance float32 // block
	Meshes            map[*graphic.Mesh]math32.Vector3
	Blocks            *block.Repository
	LatestUpdate      math32.Vector3
	Initialized       bool
	AddMeshChan       chan []*graphic.Mesh
	RemoveMeshChan    chan []*graphic.Mesh
	PositionChan      chan math32.Vector3
	ChanPackSize      int
}

func NewDemoWorld(
	renderingDistance float32,
	repository *block.Repository,
	addMeshChannel chan []*graphic.Mesh,
	removeMeshChannel chan []*graphic.Mesh,
	positionChannel chan math32.Vector3,
	chanPackSize int,
) *DemoWorld {
	demo := DemoWorld{
		Center:            math32.Vector2{X: 0, Y: 0},
		RenderingDistance: renderingDistance,
		Meshes:            make(map[*graphic.Mesh]math32.Vector3),
		Blocks:            repository,
		LatestUpdate:      math32.Vector3{X: 0, Y: 0, Z: 0},
		Initialized:       false,
		AddMeshChan:       addMeshChannel,
		RemoveMeshChan:    removeMeshChannel,
		PositionChan:      positionChannel,
		ChanPackSize:      chanPackSize,
	}

	return &demo
}

func (dw *DemoWorld) Run() {
	for {
		pos, ok := <-dw.PositionChan

		if !ok {
			return
		}

		pos = dw.getWorldPosition(pos)

		if !dw.needsUpdate(pos) {
			continue
		}

		dw.populate(pos)
		dw.cleanTooFar(pos)
	}
}

func (dw *DemoWorld) needsUpdate(pos math32.Vector3) bool {
	needsUpdate := false
	if !dw.Initialized || dw.LatestUpdate.X != pos.X || dw.LatestUpdate.Z != pos.Z {
		needsUpdate = true
	}

	dw.Initialized = true
	dw.LatestUpdate = pos

	return needsUpdate
}

func (dw *DemoWorld) populate(pos math32.Vector3) {
	var toAdd []*graphic.Mesh
	for x := pos.X - dw.RenderingDistance; x <= pos.X+dw.RenderingDistance; x++ {
		for z := pos.Z - dw.RenderingDistance; z <= pos.Z+dw.RenderingDistance; z++ {
			meshPos := math32.Vector3{X: x, Y: -2, Z: z}
			if !dw.hasMeshAt(meshPos) {
				mesh := dw.createMeshAt(meshPos)
				dw.Meshes[mesh] = meshPos
				toAdd = append(toAdd, mesh)

				if len(toAdd) > dw.ChanPackSize {
					dw.AddMeshChan <- toAdd
					toAdd = nil
				}
			}
		}
	}

	if len(toAdd) > 0 {
		dw.AddMeshChan <- toAdd
		toAdd = nil
	}
}

func (dw *DemoWorld) createMeshAt(pos math32.Vector3) *graphic.Mesh {
	mesh := dw.Blocks.Get("green_grass").CreateMesh()
	mesh.SetPositionVec(&pos)

	return mesh
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
	var toRemove []*graphic.Mesh
	for mesh, meshPos := range dw.Meshes {
		dist := math32.Max(math32.Abs(meshPos.X-pos.X), math32.Abs(meshPos.Z-pos.Z))

		if dist > dw.RenderingDistance {
			toRemove = append(toRemove, mesh)
			delete(dw.Meshes, mesh)

			if len(toRemove) > dw.ChanPackSize {
				dw.RemoveMeshChan <- toRemove
				toRemove = nil
			}
		}
	}

	if len(toRemove) > 0 {
		dw.RemoveMeshChan <- toRemove
		toRemove = nil
	}
}
