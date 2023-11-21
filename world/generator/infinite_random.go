package generator

import (
	"math/rand"
	"mycraft/world"
	"mycraft/world/block"
)

type InfiniteRandom struct {
}

func (g *InfiniteRandom) Populate(chunk *world.Chunk) {
	blockRepository := block.GetRepository()
	for x := 0; x < world.ChunkWith; x++ {
		for z := 0; z < world.ChunkDepth; z++ {
			gLevel := rand.Intn(10) + 4
			for y := 0; y < world.ChunkHeight; y++ {
				if y == gLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.GrassBlock))
					continue
				}
				if y < gLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.DirtBlock))
				}
			}
		}
	}
}

func NewInfiniteRandomGenerator() *InfiniteRandom {
	return new(InfiniteRandom)
}
