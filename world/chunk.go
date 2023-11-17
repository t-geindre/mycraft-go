package world

import "github.com/g3n/engine/math32"

const ChunkSize = 4 // x Chunklet

type Chunk struct {
	chunklets []*Chunklet
	position  *math32.Vector2
	size      uint16
}

func NewChunk(pos math32.Vector2) *Chunk {
	c := new(Chunk)
	c.position = &pos
	c.size = ChunkSize

	return c
}

func (c *Chunk) Position() *math32.Vector2 {
	return c.position
}

func (c *Chunk) AddChunklet(chunklet *Chunklet) {
	c.chunklets = append(c.chunklets, chunklet)
}
