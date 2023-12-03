package biome

import (
	"mycraft/world/block"
	"mycraft/world/chunk"
)

type Plains struct {
}

func NewPlains() *Plains {
	p := new(Plains)

	return p
}

func (p *Plains) FillGround(chunk *chunk.Chunk, ground, x, z float32) {
	for y := ground; y > 0; y-- {
		chunk.SetBlockAtF(x, y, z, p.getBlockAt(ground, x, y, z))
	}
}

func (p *Plains) getBlockAt(ground, x, y, z float32) uint8 {
	if y == ground {
		return block.TypeGrass
	}

	if y < ground-10 {
		return block.TypeStone
	}

	if y < ground {
		return block.TypeDirt
	}

	return block.TypeNone
}
