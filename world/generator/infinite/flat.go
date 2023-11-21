package infinite

import (
	"mycraft/world"
	"mycraft/world/block"
)

const GroundLevel = 4

type Flat struct {
}

func (g *Flat) Populate(chunk *world.Chunk) {
	blockRepository := block.GetRepository()
	for x := 0; x < world.ChunkWith; x++ {
		for z := 0; z < world.ChunkDepth; z++ {
			for y := 0; y < world.ChunkHeight; y++ {
				if y == GroundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.GrassBlock))
					continue
				}
				if y < GroundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.DirtBlock))
				}
			}
		}
	}
}

func NewFlatGenerator() *Flat {
	return new(Flat)
}
