package generator

import (
	"github.com/g3n/engine/math32"
	"math/rand"
	"mycraft/world"
	"mycraft/world/block"
)

type InfiniteRandom struct {
}

func (g *InfiniteRandom) Populate(chunk *world.Chunk, addChunkletChan chan []*world.Chunklet) {
	chunkPos := chunk.Position()
	chunklets := make([]*world.Chunklet, 0)
	for cX := 0; cX < world.ChunkSize*world.ChunkletSize; cX += world.ChunkletSize {
		for cY := 0; cY < world.ChunkSize*world.ChunkletSize; cY += world.ChunkletSize {
			position := math32.Vector3{
				X: chunkPos.X + float32(cX),
				Y: 0,
				Z: chunkPos.Y + float32(cY),
			}
			chunklet := world.NewChunklet(position)
			for x := 0; x < world.ChunkletSize; x++ {
				for y := 0; y < world.ChunkletSize; y++ {
					height := rand.Intn(5) + 1
					for i := height; i > 0; i-- {
						chunklet.SetBlockAt(x, i, y, block.GetRepository().Get(block.GrassBlock))
					}
				}
			}
			chunk.AddChunklet(chunklet)
			chunklets = append(chunklets, chunklet)
			if len(chunklets) == 50 {
				addChunkletChan <- chunklets
				chunklets = make([]*world.Chunklet, 0)
			}
		}
	}
	if len(chunklets) != 0 {
		addChunkletChan <- chunklets
	}
}

func NewInfiniteRandomGenerator() *InfiniteRandom {
	return new(InfiniteRandom)
}
