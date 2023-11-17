package material

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

const (
	GrassBlockTop = iota
	GrassBlockSide
	DirtBlock
)

type materialDef struct {
	TextureFile string
	Transparent bool
	Setup       func(m *material.Standard)
}

func materialReference() map[uint16]materialDef {
	ref := make(map[uint16]materialDef)

	ref[GrassBlockTop] = materialDef{
		TextureFile: "assets/block/grass_block_top.png",
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.66, G: 0.99, B: 0.59})
		},
	}

	ref[GrassBlockSide] = materialDef{
		TextureFile: "assets/block/grass_block_side.png",
	}

	ref[DirtBlock] = materialDef{
		TextureFile: "assets/block/dirt.png",
	}

	return ref
}
