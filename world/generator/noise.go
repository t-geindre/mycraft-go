package generator

import (
	"github.com/g3n/engine/math32"
	"github.com/ojrac/opensimplex-go"
	"mycraft/world"
	"mycraft/world/block"
)

type Noise struct {
	seed      int64
	noise     opensimplex.Noise32
	blocks    block.Repository
	elevation *world.Vector2FloatCache
}

func NewNoiseGenerator(seed int64) *Noise {
	n := new(Noise)
	n.seed = seed
	n.elevation = world.NewElevationCache(func(x, z float32) float32 {
		return n.GetGroundLevelAt(x, z)
	})

	// Normalized noise will return 0-1 values
	n.noise = opensimplex.NewNormalized32(n.seed)
	return n
}

func (n *Noise) GetBlockAt(x, y, z float32) uint16 {
	groundLevel := n.elevation.GetValue(x, z)
	if y == groundLevel {
		return block.BlockGrass
	}

	if y < groundLevel {
		return block.BlockDirt
	}

	if y > groundLevel && y < 90 {
		return block.BlockWater
	}

	return block.BlockNone
}

func (n *Noise) GetGroundLevelAt(x, z float32) float32 {
	noise := float32(50)
	pvNoise := n.noise.Eval2(x/200, z/200)

	noise += n.noise.Eval2(x/50, z/50) * 100 * pvNoise
	noise += n.noise.Eval2(x/5000, z/5000) * 100 * (1 - pvNoise)

	noise += n.noise.Eval2(x/5, z/5) * 2

	return math32.Floor(noise)
}

func (n *Noise) Reset() {
	n.elevation.Reset()
}
