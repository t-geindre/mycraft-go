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

	ground += bg.mountainNoise.Eval2(x, z)
	//ground += math32.Floor(bg.noise.Eval2(x/100, z/100) * 4)

	if y == ground {
		if y > 90 {
			return block.BlockGrassSnow
		}
		return block.BlockGrass
	}

	if y < ground-2 {
		return block.BlockStone
	}

	if y < ground {
		return block.BlockDirt
	}

	if y > ground && y < 50 {
		return block.BlockWater
	}

	return block.BlockNone
}

func (bg BiomeGenerator) Reset() {
	bg.noise.Clear()
}
