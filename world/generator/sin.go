package generator

import (
	"github.com/g3n/engine/math32"
	"mycraft/world"
	"mycraft/world/block"
)

type Sin struct {
	elevation *world.Vector2FloatCache
}

func NewSinGenerator(amplitude, period float32) *Sin {
	i := new(Sin)
	i.elevation = world.NewElevationCache(func(x, z float32) float32 {
		return math32.Ceil(amplitude * (math32.Sin(x/period) + math32.Sin(z/period) + 2))
	})
	return i
}

func (i *Sin) GetBlockAt(x, y, z float32) uint16 {
	groundLevel := i.elevation.GetValue(x, z)
	if y < groundLevel {
		return block.BlockDirt
	}
	if y == groundLevel {
		return block.BlockGrass
	}

	return block.BlockNone
}

func (i *Sin) Reset() {
	i.elevation.Reset()
}
