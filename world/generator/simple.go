package generator

import (
	"mycraft/world/block"
	"mycraft/world/chunk"
	"mycraft/world/generator/mod"
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

func (s *Simple) Populate(chk *chunk.Chunk) []*mod.Mod {
	for x := float32(0); x < chunk.Width; x++ {
		for z := float32(0); z < chunk.Width; z++ {
			sampleX := x + chk.Position().X
			sampleZ := z + chk.Position().Y
			for y := float32(0); y < chunk.Height; y++ {
				if y == 0 {
					chk.SetBlockAtF(x, y, z, block.TypeBedrock)
					continue
				}
				chk.SetBlockAtF(x, y, z, s.getBlockAt(sampleX, y, sampleZ))
			}
		}
	}

	return nil
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
