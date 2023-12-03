package biome

import "mycraft/world"

type Biome interface {
	FillGround(chunk *world.Chunk, ground, x, z float32)
}
