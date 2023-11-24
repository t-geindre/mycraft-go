package generator

import "github.com/ojrac/opensimplex-go"

type BiomeGenerator struct {
	seed       int64
	biomeNoise opensimplex.Noise32
}

func NewBiomeGenerator(seed int64) *BiomeGenerator {
	b := new(BiomeGenerator)
	b.seed = seed
	b.biomeNoise = opensimplex.NewNormalized32(seed)
	// biome strength calculated accordig to the more centeral value in given ranges (0-1)
	return b
}

func (b BiomeGenerator) GetBlockAt(x, y, z float32) uint16 {
	panic("implement me")
}

func (b BiomeGenerator) Reset() {
	panic("implement me")
}
