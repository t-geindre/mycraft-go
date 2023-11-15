package world

import (
	"github.com/g3n/engine/math32"
	"mycraft/block"
)

const ChunkletSize = 16

type Chunklet struct {
	Position *math32.Vector3
	Blocks   [ChunkletSize][ChunkletSize][ChunkletSize]*block.Block
	Size     int
}

func NewChunklet(position *math32.Vector3) *Chunklet {
	c := new(Chunklet)
	c.Position = position
	c.Size = ChunkletSize

	return c
}
