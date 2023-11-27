package generator

import (
	"mycraft/world"
	"mycraft/world/block"
	"mycraft/world/generator/biome"
	"mycraft/world/generator/noise"
	"mycraft/world/generator/noise/normalized"
)

type BiomeGenerator struct {
	mountainNoise noise.Noise
	noise         *noise.Cached
	biomes        []biome.Biome
}

func NewBiomeGenerator(seed int64) *BiomeGenerator {
	bg := new(BiomeGenerator)
	bg.noise = noise.GetCachedNoise(normalized.NewSimplexNoise(seed), world.ChunkDepth)
	bg.mountainNoise = noise.GetCachedNoise(noise.NewMountainNoise(seed), world.ChunkDepth)

	return bg
}

func (bg BiomeGenerator) GetBlockAt(x, y, z float32) uint16 {
	ground := float32(50)
	ground += bg.noise.Eval2(x/100, z/100)*4 - 2

	ground += bg.mountainNoise.Eval2(x, z)

	if y < ground {
		if ground > 80 {
			return block.BlockStone
		}
		return block.BlockGrass
	}

	return block.BlockNone
}

func (bg BiomeGenerator) Reset() {
	bg.noise.Clear()
}
