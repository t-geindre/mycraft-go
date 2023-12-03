package biome

import (
	"mycraft/world/chunk"
	"mycraft/world/generator/mod"
)

type Biome interface {
	FillGround(chunk *chunk.Chunk, ground, x, z float32)
	GetMod(ground, x, z float32) mod.Mod
}
