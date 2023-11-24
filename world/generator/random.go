package generator

import (
	"github.com/g3n/engine/math32"
	"math/rand"
	"mycraft/world"
	"mycraft/world/block"
)

type Random struct {
	elevation *world.Vector2FloatCache
}

func NewRandomGenerator() *Random {
	r := new(Random)
	r.elevation = world.NewElevationCache(func(x, z float32) float32 {
		return math32.Floor(rand.Float32()*10) + 50
	})
	return r
}

func (r *Random) GetBlockAt(x, y, z float32) uint16 {
	gLevel := r.elevation.GetValue(x, z)
	if y == gLevel {
		return block.BlockGrass
	}
	if y < gLevel-2 {
		return block.BlockStone
	}
	if y < gLevel {
		return block.BlockDirt
	}
	return block.BlockNone
}

func (r *Random) Reset() {
	r.elevation.Reset()
}
