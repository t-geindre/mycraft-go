package world

import (
	"github.com/g3n/engine/math32"
	"mycraft/block"
)

type DemoWorld struct {
	Center            math32.Vector2
	RenderingDistance float32 // block
	Blocks            map[*block.Block]math32.Vector3
	BlockRepository   *block.Repository
	LatestUpdate      math32.Vector3
	Initialized       bool
	AddMeshChan       chan []*block.Block
	RemoveMeshChan    chan []*block.Block
	PositionChan      chan math32.Vector3
	ChanPackSize      int
}

func NewDemoWorld(
	renderingDistance float32,
	repository *block.Repository,
	addMeshChannel chan []*block.Block,
	removeMeshChannel chan []*block.Block,
	positionChannel chan math32.Vector3,
	chanPackSize int,
) *DemoWorld {
	demo := DemoWorld{
		Center:            math32.Vector2{X: 0, Y: 0},
		RenderingDistance: renderingDistance,
		Blocks:            make(map[*block.Block]math32.Vector3),
		BlockRepository:   repository,
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
	var toAdd []*block.Block
	for x := pos.X - dw.RenderingDistance; x <= pos.X+dw.RenderingDistance; x++ {
		for z := pos.Z - dw.RenderingDistance; z <= pos.Z+dw.RenderingDistance; z++ {
			meshPos := math32.Vector3{X: x, Y: -2, Z: z}
			if !dw.hasMeshAt(meshPos) {
				b := dw.createBlockAt(meshPos)
				dw.Blocks[b] = meshPos
				toAdd = append(toAdd, b)

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

func (dw *DemoWorld) createBlockAt(pos math32.Vector3) *block.Block {
	b := dw.BlockRepository.Get("green_grass")
	b.SetPosition(pos)

	return b
}

func (dw *DemoWorld) hasMeshAt(pos math32.Vector3) bool {
	for _, meshPos := range dw.Blocks {
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
	var toRemove []*block.Block
	for b, meshPos := range dw.Blocks {
		dist := math32.Max(math32.Abs(meshPos.X-pos.X), math32.Abs(meshPos.Z-pos.Z))

		if dist > dw.RenderingDistance {
			toRemove = append(toRemove, b)
			delete(dw.Blocks, b)

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
