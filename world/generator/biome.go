package generator

import (
	"mycraft/world"
	"mycraft/world/block"
	"mycraft/world/generator/biome"
	"mycraft/world/generator/noise"
	"mycraft/world/generator/noise/normalized"
)

type BiomeGenerator struct {
	noise  noise.Noise
	biomes []biome.Biome
}

func NewBiomeGenerator(seed int64) *BiomeGenerator {
	bg := new(BiomeGenerator)

	bg.noise = normalized.NewSimplexNoise(seed)
	bg.noise = noise.NewOctave(bg.noise, 2.1, 0.7, 5)
	bg.noise = noise.NewScale(bg.noise, 1600)
	bg.noise = noise.NewAmplify(bg.noise, 70)
	bg.noise = noise.NewFloor(bg.noise)
	bg.noise = noise.NewClamp(bg.noise, 1, world.ChunkHeight-1)

	bg.biomes = []biome.Biome{
		biome.NewWater(1, 25),
		biome.NewBeach(26, 26),
		biome.NewPlains(26, 255, seed),
	}

	return bg
}

func (bg BiomeGenerator) Populate(chunk *world.Chunk) {
	for x := float32(0); x < world.ChunkWidth; x++ {
		for z := float32(0); z < world.ChunkWidth; z++ {
			sampleX := x + chunk.Position().X
			sampleZ := z + chunk.Position().Y

			ground := bg.noise.Eval2(sampleX, sampleZ)

			var localBiome biome.Biome
			for _, b := range bg.biomes {
				if b.Match(ground) {
					localBiome = b
					break
				}
			}

			if localBiome == nil {
				continue
			}

			localBiome.SetGround(ground)

			for y := float32(0); y < world.ChunkHeight; y++ {
				if y == 0 {
					chunk.SetBlockAtF(x, y, z, block.TypeBedrock)
					continue
				}
				chunk.SetBlockAtF(x, y, z, localBiome.GetBlockAt(sampleX, y, sampleZ))
			}
		}
	}
}
