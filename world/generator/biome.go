package generator

import (
	"github.com/g3n/engine/math32"
	"mycraft/world"
	"mycraft/world/block"
	"mycraft/world/generator/noise"
	"mycraft/world/generator/noise/normalized"
)

type BiomeGenerator struct {
	noise      noise.Noise
	waterLevel float32
}

func NewBiomeGenerator(seed int64) *BiomeGenerator {
	bg := new(BiomeGenerator)

	bg.noise = normalized.NewSimplexNoise(seed)
	bg.noise = noise.NewScale(bg.noise, 200)
	bg.noise = noise.NewOctave(bg.noise, 2, 0.5, 4)
	bg.noise = noise.NewAmplify(bg.noise, 50)

	bg.waterLevel = 25

	return bg
}

func (bg BiomeGenerator) Populate(chunk *world.Chunk) {
	for x := float32(0); x < world.ChunkWidth; x++ {
		for z := float32(0); z < world.ChunkWidth; z++ {

			ground := bg.noise.Eval2(x+chunk.Position().X, z+chunk.Position().Y)
			ground = math32.Floor(ground)
			ground = math32.Clamp(ground, 1, world.ChunkHeight-1)

			for y := float32(0); y < world.ChunkHeight; y++ {
				if y > ground && y > bg.waterLevel {
					break
				}
				b := bg.getBlock(ground, x, y, z)
				chunk.SetBlockAt(int(x), int(y), int(z), b)
			}
		}
	}
}
func (bg BiomeGenerator) getBlock(ground, x, y, z float32) uint16 {
	if y == ground {
		if y <= bg.waterLevel {
			return block.BlockSand
		}
		return block.BlockGrass
	}
	if y < ground-2 {
		return block.BlockStone
	}

	if y < ground {
		if y <= bg.waterLevel {
			return block.BlockSand
		}
		return block.BlockDirt
	}

	if y > ground && y < bg.waterLevel {
		return block.BlockWater
	}

	return block.BlockNone
}
