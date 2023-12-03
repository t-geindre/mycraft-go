package chunk

import (
	"github.com/g3n/engine/math32"
)

const Width = 16
const Depth = 16
const Height = 256

type Chunk struct {
	blocks       [Width][Height][Depth]uint8
	position     *math32.Vector2
	size         *math32.Vector3
	isEmpty      bool
	filledLayers map[int]bool
}

func NewChunk(pos math32.Vector2) *Chunk {
	c := new(Chunk)
	c.position = &pos
	c.size = &math32.Vector3{
		X: Width,
		Y: Height,
		Z: Depth,
	}
	c.filledLayers = make(map[int]bool, Height)

	return c
}

func (c *Chunk) Position() *math32.Vector2 {
	return c.position
}

func (c *Chunk) SetBlockAtF(x, y, z float32, b uint8) {
	c.SetBlockAt(int(x), int(y), int(z), b)
}

func (c *Chunk) SetBlockAt(x, y, z int, b uint8) {
	c.filledLayers[y] = true
	c.blocks[x][y][z] = b
}

func (c *Chunk) GetBlockAt(x, y, z int) uint8 {
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

func (c *Chunk) InBounds(x, y, z float32) bool {
	return x >= c.position.X && x < c.position.X+Width && y >= 0 && y < Height && z >= c.position.Y && z < c.position.Y+Depth
}
