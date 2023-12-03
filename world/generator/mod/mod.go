package mod

import (
	"github.com/g3n/engine/math32"
	"mycraft/world"
)

type Mod interface {
	Apply(chunk *world.Chunk)
	GetBounds() math32.Vector4
}
