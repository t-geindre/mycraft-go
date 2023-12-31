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
	BlockSand
	Bedrock
	Dandelion
	Gravel
	OakLogSide
	OakSpruceLogTop
	OakLeaves
)

type materialDef struct {
	TextureFile string
	Setup       func(m *material.Standard)
}

func materialReference() map[uint16]materialDef {
	ref := make(map[uint16]materialDef)

	ref[BlockGrassTop] = materialDef{
		TextureFile: "assets/block/grass_block_top.png",
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.49, G: 0.74, B: 0.42})
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

	ref[BlockSand] = materialDef{
		TextureFile: "assets/block/sand.png",
	}

	ref[BlockWater] = materialDef{
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.2, G: 0.4, B: 0.8})
			m.SetOpacity(0.8)
			m.SetTransparent(true)
		},
	}

	ref[Bedrock] = materialDef{
		TextureFile: "assets/block/bedrock.png",
	}

	ref[Dandelion] = materialDef{
		TextureFile: "assets/block/dandelion.png",
	}

	ref[Gravel] = materialDef{
		TextureFile: "assets/block/gravel.png",
	}

	ref[OakLogSide] = materialDef{
		TextureFile: "assets/block/oak_log.png",
	}

	ref[OakSpruceLogTop] = materialDef{
		TextureFile: "assets/block/oak_log_top.png",
	}

	ref[OakLeaves] = materialDef{
		TextureFile: "assets/block/oak_leaves.png",
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.66, G: 0.99, B: 0.59})
			m.SetTransparent(true)
		},
	}

	return ref
}
