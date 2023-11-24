package generator

import (
	"mycraft/world"
	"mycraft/world/block"
	"mycraft/world/generator/noise"
)

type Simple struct {
	noise      *noise.Cached
	waterLevel float32
}

func NewSimpleGenerator(baseNoise noise.Noise, waterLevel float32) *Simple {
	s := new(Simple)
	s.waterLevel = waterLevel

	// Add noise caching
	switch baseNoise.(type) {
	case *noise.Cached:
		s.noise = baseNoise.(*noise.Cached)
	default:
		s.noise = noise.NewCachedNoise(baseNoise, world.ChunkDepth)
	}

	return s
}

func (s *Simple) GetBlockAt(x, y, z float32) uint16 {
	groundLevel := s.noise.Eval2(x, z)
	if y == groundLevel {
		return block.BlockGrass
	}

	if y < groundLevel {
		return block.BlockDirt
	}

	if y > groundLevel && y < s.waterLevel {
		return block.BlockWater
	}

	return block.BlockNone
}

func (s *Simple) Reset() {
	s.noise.Clear()
}
