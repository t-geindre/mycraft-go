package world

import (
	"github.com/g3n/engine/math32"
	"mycraft/block"
)

type DemoWorld struct {
	Center            math32.Vector2
	RenderingDistance float32 // block
	Blocks            map[*block.Block]*math32.Vector3
	BlockRepository   *block.Repository
	LatestUpdate      math32.Vector3
	Initialized       bool
	AddBlockChan      chan []*block.Block
	RemoveBlockChan   chan []*block.Block
	PositionChan      chan math32.Vector3
	ChanPackSize      int
}

func NewDemoWorld(
	renderingDistance float32,
	addBlockChannel chan []*block.Block,
	removeBlockChannel chan []*block.Block,
	positionChannel chan math32.Vector3,
	chanPackSize int,
) *DemoWorld {
	demo := DemoWorld{
		Center:            math32.Vector2{X: 0, Y: 0},
		RenderingDistance: renderingDistance,
		Blocks:            make(map[*block.Block]*math32.Vector3),
		BlockRepository:   block.GetRepository(),
		LatestUpdate:      math32.Vector3{X: 0, Y: 0, Z: 0},
		Initialized:       false,
		AddBlockChan:      addBlockChannel,
		RemoveBlockChan:   removeBlockChannel,
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
			blockPos := &math32.Vector3{X: x, Y: -2, Z: z}
			if blockPos.DistanceTo(&pos) > dw.RenderingDistance {
				continue
			}
			if !dw.hasMeshAt(blockPos) {
				bl := dw.BlockRepository.Get("green_grass")
				bl.SetPositionVec(blockPos)

				dw.Blocks[bl] = blockPos
				toAdd = append(toAdd, bl)

				flPos := &math32.Vector3{X: x, Y: -1, Z: z}
				fl := dw.BlockRepository.Get("dandelion")
				fl.SetPositionVec(flPos)

				dw.Blocks[fl] = blockPos
				toAdd = append(toAdd, fl)

				if len(toAdd) > dw.ChanPackSize {
					dw.AddBlockChan <- toAdd
					toAdd = nil
				}
			}
		}
	}

	if len(toAdd) > 0 {
		dw.AddBlockChan <- toAdd
		toAdd = nil
	}
}

func (dw *DemoWorld) hasMeshAt(pos *math32.Vector3) bool {
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
	for bl, meshPos := range dw.Blocks {
		dist := math32.Max(math32.Abs(meshPos.X-pos.X), math32.Abs(meshPos.Z-pos.Z))

		if dist > dw.RenderingDistance {
			toRemove = append(toRemove, bl)
			delete(dw.Blocks, bl)

			if len(toRemove) > dw.ChanPackSize {
				dw.RemoveBlockChan <- toRemove
				toRemove = nil
			}
		}
	}

	if len(toRemove) > 0 {
		dw.RemoveBlockChan <- toRemove
		toRemove = nil
	}
}
