package infinite

import (
	"math/rand"
	"mycraft/world"
	"mycraft/world/block"
)

type Random struct {
}

func (g *Random) Populate(chunk *world.Chunk) {
	blockRepository := block.GetRepository()
	for x := 0; x < world.ChunkWith; x++ {
		for z := 0; z < world.ChunkDepth; z++ {
			gLevel := rand.Intn(31) + 1
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

func NewRandomGenerator() *Random {
	return new(Random)
}
