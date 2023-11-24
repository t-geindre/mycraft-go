package generator

import (
	"mycraft/world/block"
)

const GroundLevel = 4

type Flat struct {
}

func NewFlatGenerator() *Flat {
	return new(Flat)
}

func (g *Flat) GetBlockAt(x, y, z float32) uint16 {
	if y == GroundLevel {
		return block.BlockGrass
	}
	if y < GroundLevel {
		return block.BlockDirt
	}
	return block.BlockNone
}

func (g *Flat) Reset() {
}
