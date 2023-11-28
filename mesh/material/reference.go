package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

const (
	BlockGrassTop = iota
	BlockGrassSide
	BlockGrassSideSnow
	BlockDirt
	BlockStone
	BlockWater
	BlockSand
	Andesite
	BarrelSide
	BarrelBottom
	BarrelTop
	Bedrock
	Dandelion
	DiamondBlock
	DiamondOre
	IronOre
	Netherrack
	OakPlanks
	OrangeTulip
	OxeyeDaisy
	PackedIce
	SmithingTableBottom
	SmithingTableFront
	SmithingTableSide
	SmithingTableTop
	SmoothBasalt
	Snow
	SoulSand
	SpruceLogSide
	SpruceLogTop
	StoneBricks
	TntBottom
	TntSide
	TntTop
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

	ref[BlockGrassSideSnow] = materialDef{
		TextureFile: "assets/block/dirt.png",
		Setup: func(m *material.Standard) {
			t, _ := texture.NewTexture2DFromImage("assets/block/grass_block_side_overlay.png")
			t.SetMagFilter(gls.NEAREST)
			t.SetWrapT(gls.REPEAT)
			t.SetWrapS(gls.REPEAT)
			m.AddTexture(t)
		},
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
		Transparent: true,
		Setup: func(m *material.Standard) {
			m.SetColor(&math32.Color{R: 0.2, G: 0.4, B: 0.8})
			m.SetOpacity(0.8)
			m.SetTransparent(true)
		},
	}

	ref[Andesite] = materialDef{
		TextureFile: "assets/block/andesite.png",
	}

	ref[BarrelSide] = materialDef{
		TextureFile: "assets/block/barrel_side.png",
	}

	ref[BarrelBottom] = materialDef{
		TextureFile: "assets/block/barrel_bottom.png",
	}

	ref[BarrelTop] = materialDef{
		TextureFile: "assets/block/barrel_top.png",
	}

	ref[Bedrock] = materialDef{
		TextureFile: "assets/block/bedrock.png",
	}

	ref[Dandelion] = materialDef{
		TextureFile: "assets/block/dandelion.png",
	}

	ref[DiamondBlock] = materialDef{
		TextureFile: "assets/block/diamond_block.png",
	}

	ref[DiamondOre] = materialDef{
		TextureFile: "assets/block/diamond_ore.png",
	}

	ref[IronOre] = materialDef{
		TextureFile: "assets/block/iron_ore.png",
	}

	ref[Netherrack] = materialDef{
		TextureFile: "assets/block/netherrack.png",
	}

	ref[OakPlanks] = materialDef{
		TextureFile: "assets/block/oak_planks.png",
	}

	ref[OrangeTulip] = materialDef{
		TextureFile: "assets/block/orange_tulip.png",
	}

	ref[OxeyeDaisy] = materialDef{
		TextureFile: "assets/block/oxeye_daisy.png",
	}

	ref[PackedIce] = materialDef{
		TextureFile: "assets/block/packed_ice.png",
	}

	ref[SmithingTableBottom] = materialDef{
		TextureFile: "assets/block/smithing_table_bottom.png",
	}

	ref[SmithingTableFront] = materialDef{
		TextureFile: "assets/block/smithing_table_front.png",
	}

	ref[SmithingTableSide] = materialDef{
		TextureFile: "assets/block/smithing_table_side.png",
	}

	ref[SmithingTableTop] = materialDef{
		TextureFile: "assets/block/smithing_table_top.png",
	}

	ref[SmoothBasalt] = materialDef{
		TextureFile: "assets/block/smooth_basalt.png",
	}

	ref[Snow] = materialDef{
		TextureFile: "assets/block/snow.png",
	}

	ref[SoulSand] = materialDef{
		TextureFile: "assets/block/soul_sand.png",
	}

	ref[SpruceLogSide] = materialDef{
		TextureFile: "assets/block/spruce_log.png",
	}

	ref[SpruceLogTop] = materialDef{
		TextureFile: "assets/block/spruce_log_top.png",
	}

	ref[StoneBricks] = materialDef{
		TextureFile: "assets/block/stone_bricks.png",
	}

	ref[TntBottom] = materialDef{
		TextureFile: "assets/block/tnt_bottom.png",
	}

	ref[TntSide] = materialDef{
		TextureFile: "assets/block/tnt_side.png",
	}

	ref[TntTop] = materialDef{
		TextureFile: "assets/block/tnt_top.png",
	}

	return ref
}
