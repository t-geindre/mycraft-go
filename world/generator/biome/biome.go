package biome

import (
	"mycraft/world/chunk"
)

type Biome interface {
	FillGround(chunk *chunk.Chunk, ground, x, z float32)
}
