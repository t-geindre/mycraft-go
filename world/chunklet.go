package world

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"mycraft/world/block"
)

const ChunkletSize = 16

const OnDispose = "onDispose"

type Chunklet struct {
	Position *math32.Vector3
	Blocks   [ChunkletSize][ChunkletSize][ChunkletSize]*block.Block
	Size     int
	Empty    bool
	core.Dispatcher
}

func NewChunklet(position math32.Vector3) *Chunklet {
	c := new(Chunklet)
	c.Position = &position
	c.Size = ChunkletSize
	c.Empty = true
	c.Dispatcher.Initialize()

	return c
}

func (c *Chunklet) SetBlockAt(x, y, z int, b *block.Block) {
	c.Blocks[x][y][z] = b
	c.Empty = false
}

func (c *Chunklet) Dispose() {
	c.Dispatcher.Dispatch(OnDispose, c)
}
