package material

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

const (
	BlockGrassTop = iota
	BlockGrassSide
	BlockDirt
	BlockStone
	BlockWater
)

type materialDef struct {
	TextureFile string
	Transparent bool
	Setup       func(m *material.Standard)
}

func materialReference() map[uint16]materialDef {
	ref := make(map[uint16]materialDef)

	ref[BlockGrassTop] = materialDef{
		TextureFile: "assets/block/grass_block_top.png",
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.66, G: 0.99, B: 0.59})
		},
	}

	ref[BlockGrassSide] = materialDef{
		TextureFile: "assets/block/grass_block_side.png",
	}

	ref[BlockDirt] = materialDef{
		TextureFile: "assets/block/dirt.png",
	}

	ref[BlockStone] = materialDef{
		TextureFile: "assets/block/stone.png",
	}

	ref[BlockWater] = materialDef{
		Transparent: true,
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.2, G: 0.4, B: 0.8})
			m.SetOpacity(0.8)
			m.SetTransparent(true)
		},
	}

	return ref
}
