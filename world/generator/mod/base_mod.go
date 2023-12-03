package mod

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/chunk"
)

type modBlock struct {
	position *math32.Vector3
	block    uint8
}

type BaseMod struct {
	modBlocks   []modBlock
	bounds      *math32.Vector4
	boundsDirty bool
}

func NewBaseMod() *BaseMod {
	b := new(BaseMod)
	b.modBlocks = []modBlock{}
	b.bounds = &math32.Vector4{X: 0, Y: 0, Z: 0, W: 0}
	b.boundsDirty = false

	return b
}

func (b *BaseMod) Apply(chunk *chunk.Chunk) {
	outOfBounds := make([]modBlock, 0)

	for idx, block := range b.modBlocks {
		if chunk.InBounds(block.position.X, block.position.Y, block.position.Z) {
			chunk.SetBlockAt(int(block.position.X), int(block.position.Y), int(block.position.Z), block.block)
			continue
		}
		outOfBounds = append(outOfBounds, b.modBlocks[idx])
	}

	b.modBlocks = []modBlock{}
}

func (b *BaseMod) GetBounds() *math32.Vector4 {
	if b.boundsDirty {
		b.bounds = &math32.Vector4{X: 0, Y: 0, Z: 0, W: 0}
		for _, block := range b.modBlocks {
			if block.position.X < b.bounds.X {
				b.bounds.X = block.position.X
			}
			if block.position.Z < b.bounds.Y {
				b.bounds.Y = block.position.Z
			}
			if block.position.X > b.bounds.W {
				b.bounds.W = block.position.X
			}
			if block.position.Z > b.bounds.W {
				b.bounds.W = block.position.Z
			}
		}
		b.boundsDirty = false
	}
	return b.bounds
}

func (b *BaseMod) AddBlock(x, y, z float32, block uint8) {
	b.modBlocks = append(b.modBlocks, modBlock{
		position: &math32.Vector3{
			X: x,
			Y: y,
			Z: z,
		},
		block: block,
	})
	b.boundsDirty = true
}

func (b *BaseMod) IsEmpty() bool {
	return len(b.modBlocks) == 0
}
