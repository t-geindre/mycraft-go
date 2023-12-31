package generator

import (
	"mycraft/world/block"
	"mycraft/world/chunk"
	"mycraft/world/generator/biome"
	"mycraft/world/generator/mod"
	"mycraft/world/generator/noise"
	"mycraft/world/generator/noise/normalized"
)

type BiomeGenerator struct {
	elevationNoise noise.Noise
	biomeNoise     noise.Noise
	biomes         *biome.Selector
}

func NewBiomeGenerator(seed int64) *BiomeGenerator {
	bg := new(BiomeGenerator)

	bg.elevationNoise = normalized.NewSimplexNoise(seed)
	bg.elevationNoise = noise.NewOctave(bg.elevationNoise, 2.1, 0.7, 5)
	bg.elevationNoise = noise.NewScale(bg.elevationNoise, 1600)
	bg.elevationNoise = noise.NewAmplify(bg.elevationNoise, 70)
	bg.elevationNoise = noise.NewFloor(bg.elevationNoise)
	bg.elevationNoise = noise.NewClamp(bg.elevationNoise, 1, chunk.Height-1)

	bg.biomes = biome.NewSelector()
	bg.biomes.Add(biome.NewWater(25), 1, 25)
	bg.biomes.Add(biome.NewBeach(), 26, 26)
	bg.biomes.Add(biome.NewPlains(), 26, 255)

	return bg
}

func (bg BiomeGenerator) Populate(chk *chunk.Chunk) []*mod.Mod {
	for x := float32(0); x < chunk.Width; x++ {
		for z := float32(0); z < chunk.Width; z++ {
			sampleX := x + chk.Position().X
			sampleZ := z + chk.Position().Y

			ground := bg.elevationNoise.Eval2(sampleX, sampleZ)

			localBiome := bg.biomes.Match(ground)
			if localBiome == nil {
				continue
			}

			localBiome.FillGround(chk, ground, x, z)
			localMod := localBiome.GetMod(ground, x, z)
			if localMod != nil {
				localMod.Apply(chk)
			}
			chk.SetBlockAtF(x, 0, z, block.TypeBedrock)
		}
	}

	return nil
}
