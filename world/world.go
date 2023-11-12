package world

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type World struct {
	Center            math32.Vector2
	RenderingDistance float32 // block
	ContainerNode     *core.Node
}

type IWord interface {
	Update(x, z float32)
}
