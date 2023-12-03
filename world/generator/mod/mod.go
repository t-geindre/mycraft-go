package mod

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/chunk"
)

type Mod interface {
	Apply(chunk *chunk.Chunk)
	GetBounds() *math32.Vector4
	IsEmpty() bool
}
