package generator

import (
	"mycraft/world/chunk"
	"mycraft/world/generator/mod"
)

type Generator interface {
	Populate(chunk *chunk.Chunk) []*mod.Mod
}
