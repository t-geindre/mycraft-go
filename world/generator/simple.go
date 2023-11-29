package generator

import (
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

func (s *Simple) GetBlockAt(x, y, z float32) uint16 {
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
