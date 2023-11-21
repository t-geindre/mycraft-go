package world

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/block"
)

const ChunkWith = 16
const ChunkDepth = 16
const ChunkHeight = 256

type Chunk struct {
	blocks   [ChunkWith][ChunkHeight][ChunkDepth]*block.Block
	position *math32.Vector2
	size     *math32.Vector3
}

func NewChunk(pos math32.Vector2) *Chunk {
	c := new(Chunk)
	c.position = &pos
	c.size = &math32.Vector3{
		X: ChunkWith,
		Y: ChunkHeight,
		Z: ChunkDepth,
	}

	return c
}

func (c *Chunk) Position() *math32.Vector2 {
	return c.position
}

func (c *Chunk) SetBlockAt(x, y, z int, b *block.Block) {
	c.blocks[x][y][z] = b
}

func (c *Chunk) GetBlockAt(x, y, z int) *block.Block {
	return c.blocks[x][y][z]
}

func (c *Chunk) Size() *math32.Vector3 {
	return c.size
}
