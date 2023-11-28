package world

import (
	"github.com/g3n/engine/math32"
)

const ChunkWidth = 16
const ChunkDepth = 16
const ChunkHeight = 256

type Chunk struct {
	blocks       [ChunkWidth][ChunkHeight][ChunkDepth]uint16
	position     *math32.Vector2
	size         *math32.Vector3
	isEmpty      bool
	filledLayers map[int]bool
}

func NewChunk(pos math32.Vector2) *Chunk {
	c := new(Chunk)
	c.position = &pos
	c.size = &math32.Vector3{
		X: ChunkWidth,
		Y: ChunkHeight,
		Z: ChunkDepth,
	}
	c.filledLayers = make(map[int]bool, ChunkHeight)

	return c
}

func (c *Chunk) Position() *math32.Vector2 {
	return c.position
}

func (c *Chunk) SetBlockAt(x, y, z int, b uint16) {
	c.filledLayers[y] = true
	c.blocks[x][y][z] = b
}

func (c *Chunk) GetBlockAt(x, y, z int) uint16 {
	return c.blocks[x][y][z]
}

func (c *Chunk) Size() *math32.Vector3 {
	return c.size
}

func (c *Chunk) AreLayersEmpty(from, to int) bool {
	for i := from; i < to; i++ {
		if _, ok := c.filledLayers[i]; ok {
			return false
		}
	}

	return true
}

func (c *Chunk) IsLayerEmpty(l int) bool {
	_, ok := c.filledLayers[l]
	return !ok
}
