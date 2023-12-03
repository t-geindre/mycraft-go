package generator

import (
	"mycraft/world"
	"mycraft/world/block"
	"mycraft/world/generator/noise"
)

type Simple struct {
	noise      noise.Noise
	waterLevel float32
}

func NewSimpleGenerator(baseNoise noise.Noise, waterLevel float32) *Simple {
	s := new(Simple)
	s.waterLevel = waterLevel
	s.noise = baseNoise

	return s
}

func (s *Simple) Populate(chunk *world.Chunk) {
	for x := float32(0); x < world.ChunkWidth; x++ {
		for z := float32(0); z < world.ChunkWidth; z++ {
			sampleX := x + chunk.Position().X
			sampleZ := z + chunk.Position().Y
			for y := float32(0); y < world.ChunkHeight; y++ {
				if y == 0 {
					chunk.SetBlockAtF(x, y, z, block.TypeBedrock)
					continue
				}
				chunk.SetBlockAtF(x, y, z, s.getBlockAt(sampleX, y, sampleZ))
			}
		}
	}
}

func (s *Simple) getBlockAt(x, y, z float32) uint8 {
	groundLevel := s.noise.Eval2(x, z)
	if y == groundLevel {
		return block.TypeGrass
	}

	if y < groundLevel {
		return block.TypeDirt
	}

	if y > groundLevel && y < s.waterLevel {
		return block.TypeWater
	}

	return block.TypeNone
}
