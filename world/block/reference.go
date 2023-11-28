package block

import (
	"github.com/g3n/engine/material"
	matRef "mycraft/mesh/material"
)

const (
	BlockNone = iota
	BlockSand
	BlockGrass
	BlockGrassSnow
	BlockDirt
	BlockStone
	BlockWater
	BlockAndesite
	BlockBarrel
	BlockBedrock
	BlockDiamondBlock
	BlockDiamondOre
	BlockIronOre
	BlockNetherrack
	BlockOakPlanks
	BlockPackedIce
	BlockSmithingTable
	BlockSmoothBasalt
	BlockSnow
	BlockSoulSand
	BlockSpruceLog
	BlockStoneBricks
	BlockTnt
)

type BlockMaterials struct {
	Top    material.IMaterial
	Bottom material.IMaterial
	North  material.IMaterial
	South  material.IMaterial
	East   material.IMaterial
	West   material.IMaterial
}

type Block struct {
	Id          uint16
	Transparent bool
	Materials   BlockMaterials
}

func blockReference() map[uint16]*Block {
	ref := make(map[uint16]*Block)
	matRepo := matRef.GetRepository()

	ref[BlockGrass] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BlockGrassTop),
			Bottom: matRepo.Get(matRef.BlockDirt),
			North:  matRepo.Get(matRef.BlockGrassSide),
			South:  matRepo.Get(matRef.BlockGrassSide),
			East:   matRepo.Get(matRef.BlockGrassSide),
			West:   matRepo.Get(matRef.BlockGrassSide),
		},
	}

	ref[BlockGrassSnow] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.Snow),
			Bottom: matRepo.Get(matRef.BlockDirt),
			North:  matRepo.Get(matRef.BlockGrassSideSnow),
			South:  matRepo.Get(matRef.BlockGrassSideSnow),
			East:   matRepo.Get(matRef.BlockGrassSideSnow),
			West:   matRepo.Get(matRef.BlockGrassSideSnow),
		},
	}

	ref[BlockDirt] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BlockDirt),
			Bottom: matRepo.Get(matRef.BlockDirt),
			North:  matRepo.Get(matRef.BlockDirt),
			South:  matRepo.Get(matRef.BlockDirt),
			East:   matRepo.Get(matRef.BlockDirt),
			West:   matRepo.Get(matRef.BlockDirt),
		},
	}

	ref[BlockStone] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BlockStone),
			Bottom: matRepo.Get(matRef.BlockStone),
			North:  matRepo.Get(matRef.BlockStone),
			South:  matRepo.Get(matRef.BlockStone),
			East:   matRepo.Get(matRef.BlockStone),
			West:   matRepo.Get(matRef.BlockStone),
		},
	}

	ref[BlockWater] = &Block{
		Transparent: true,
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BlockWater),
			Bottom: matRepo.Get(matRef.BlockWater),
			North:  matRepo.Get(matRef.BlockWater),
			South:  matRepo.Get(matRef.BlockWater),
			East:   matRepo.Get(matRef.BlockWater),
			West:   matRepo.Get(matRef.BlockWater),
		},
	}

	ref[BlockSand] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BlockSand),
			Bottom: matRepo.Get(matRef.BlockSand),
			North:  matRepo.Get(matRef.BlockSand),
			South:  matRepo.Get(matRef.BlockSand),
			East:   matRepo.Get(matRef.BlockSand),
			West:   matRepo.Get(matRef.BlockSand),
		},
	}

	ref[BlockAndesite] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.Andesite),
			Bottom: matRepo.Get(matRef.Andesite),
			East:   matRepo.Get(matRef.Andesite),
			West:   matRepo.Get(matRef.Andesite),
			North:  matRepo.Get(matRef.Andesite),
			South:  matRepo.Get(matRef.Andesite),
		},
	}
	ref[BlockBarrel] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.BarrelTop),
			Bottom: matRepo.Get(matRef.BarrelBottom),
			East:   matRepo.Get(matRef.BarrelSide),
			West:   matRepo.Get(matRef.BarrelSide),
			North:  matRepo.Get(matRef.BarrelSide),
			South:  matRepo.Get(matRef.BarrelSide),
		},
	}
	ref[BlockBedrock] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.Bedrock),
			Bottom: matRepo.Get(matRef.Bedrock),
			East:   matRepo.Get(matRef.Bedrock),
			West:   matRepo.Get(matRef.Bedrock),
			North:  matRepo.Get(matRef.Bedrock),
			South:  matRepo.Get(matRef.Bedrock),
		},
	}
	ref[BlockDiamondBlock] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.DiamondBlock),
			Bottom: matRepo.Get(matRef.DiamondBlock),
			East:   matRepo.Get(matRef.DiamondBlock),
			West:   matRepo.Get(matRef.DiamondBlock),
			North:  matRepo.Get(matRef.DiamondBlock),
			South:  matRepo.Get(matRef.DiamondBlock),
		},
	}
	ref[BlockDiamondOre] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.DiamondOre),
			Bottom: matRepo.Get(matRef.DiamondOre),
			East:   matRepo.Get(matRef.DiamondOre),
			West:   matRepo.Get(matRef.DiamondOre),
			North:  matRepo.Get(matRef.DiamondOre),
			South:  matRepo.Get(matRef.DiamondOre),
		},
	}
	ref[BlockIronOre] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.IronOre),
			Bottom: matRepo.Get(matRef.IronOre),
			East:   matRepo.Get(matRef.IronOre),
			West:   matRepo.Get(matRef.IronOre),
			North:  matRepo.Get(matRef.IronOre),
			South:  matRepo.Get(matRef.IronOre),
		},
	}
	ref[BlockNetherrack] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.Netherrack),
			Bottom: matRepo.Get(matRef.Netherrack),
			East:   matRepo.Get(matRef.Netherrack),
			West:   matRepo.Get(matRef.Netherrack),
			North:  matRepo.Get(matRef.Netherrack),
			South:  matRepo.Get(matRef.Netherrack),
		},
	}
	ref[BlockOakPlanks] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.OakPlanks),
			Bottom: matRepo.Get(matRef.OakPlanks),
			East:   matRepo.Get(matRef.OakPlanks),
			West:   matRepo.Get(matRef.OakPlanks),
			North:  matRepo.Get(matRef.OakPlanks),
			South:  matRepo.Get(matRef.OakPlanks),
		},
	}
	ref[BlockPackedIce] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.PackedIce),
			Bottom: matRepo.Get(matRef.PackedIce),
			East:   matRepo.Get(matRef.PackedIce),
			West:   matRepo.Get(matRef.PackedIce),
			North:  matRepo.Get(matRef.PackedIce),
			South:  matRepo.Get(matRef.PackedIce),
		},
	}
	ref[BlockSmithingTable] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.SmithingTableTop),
			Bottom: matRepo.Get(matRef.SmithingTableBottom),
			East:   matRepo.Get(matRef.SmithingTableSide),
			West:   matRepo.Get(matRef.SmithingTableSide),
			North:  matRepo.Get(matRef.SmithingTableFront),
			South:  matRepo.Get(matRef.SmithingTableFront),
		},
	}
	ref[BlockSmoothBasalt] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.SmoothBasalt),
			Bottom: matRepo.Get(matRef.SmoothBasalt),
			East:   matRepo.Get(matRef.SmoothBasalt),
			West:   matRepo.Get(matRef.SmoothBasalt),
			North:  matRepo.Get(matRef.SmoothBasalt),
			South:  matRepo.Get(matRef.SmoothBasalt),
		},
	}
	ref[BlockSnow] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.Snow),
			Bottom: matRepo.Get(matRef.Snow),
			East:   matRepo.Get(matRef.Snow),
			West:   matRepo.Get(matRef.Snow),
			North:  matRepo.Get(matRef.Snow),
			South:  matRepo.Get(matRef.Snow),
		},
	}
	ref[BlockSoulSand] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.SoulSand),
			Bottom: matRepo.Get(matRef.SoulSand),
			East:   matRepo.Get(matRef.SoulSand),
			West:   matRepo.Get(matRef.SoulSand),
			North:  matRepo.Get(matRef.SoulSand),
			South:  matRepo.Get(matRef.SoulSand),
		},
	}
	ref[BlockSpruceLog] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.SpruceLogTop),
			Bottom: matRepo.Get(matRef.SpruceLogTop),
			East:   matRepo.Get(matRef.SpruceLogSide),
			West:   matRepo.Get(matRef.SpruceLogSide),
			North:  matRepo.Get(matRef.SpruceLogSide),
			South:  matRepo.Get(matRef.SpruceLogSide),
		},
	}
	ref[BlockStoneBricks] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.StoneBricks),
			Bottom: matRepo.Get(matRef.StoneBricks),
			East:   matRepo.Get(matRef.StoneBricks),
			West:   matRepo.Get(matRef.StoneBricks),
			North:  matRepo.Get(matRef.StoneBricks),
			South:  matRepo.Get(matRef.StoneBricks),
		},
	}
	ref[BlockTnt] = &Block{
		Materials: BlockMaterials{
			Top:    matRepo.Get(matRef.TntTop),
			Bottom: matRepo.Get(matRef.TntBottom),
			East:   matRepo.Get(matRef.TntSide),
			West:   matRepo.Get(matRef.TntSide),
			North:  matRepo.Get(matRef.TntSide),
			South:  matRepo.Get(matRef.TntSide),
		},
	}

	return ref
}
