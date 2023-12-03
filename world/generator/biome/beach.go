package biome

import (
	"mycraft/world/block"
	"mycraft/world/chunk"
)

type Beach struct {
}

func NewBeach() *Beach {
	b := new(Beach)

	return b
}

func (b *Beach) FillGround(chunk *chunk.Chunk, ground, x, z float32) {
	for y := float32(0); y <= ground; y++ {
		chunk.SetBlockAtF(x, y, z, b.getBlockAt(ground, x, y, z))
	}
}

func (b *Beach) getBlockAt(ground, x, y, z float32) uint8 {
	if y < ground-10 {
		return block.TypeStone
	}

	return block.TypeSand
}
