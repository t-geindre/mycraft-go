package block

import (
	"github.com/g3n/engine/material"
	matRef "mycraft/mesh/material"
)

const (
	BlockGrass = iota
	BlockDirt
	BlockStone
	BlockWater
)

type blockDefMaterials struct {
	Top    material.IMaterial
	Bottom material.IMaterial
	North  material.IMaterial
	South  material.IMaterial
	East   material.IMaterial
	West   material.IMaterial
}

type blockDef struct {
	Transparent bool
	Materials   blockDefMaterials
}

func blockReference() map[uint16]blockDef {
	ref := make(map[uint16]blockDef)
	matRepo := matRef.GetRepository()

	ref[BlockGrass] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.BlockGrassTop),
			Bottom: matRepo.Get(matRef.BlockDirt),
			North:  matRepo.Get(matRef.BlockGrassSide),
			South:  matRepo.Get(matRef.BlockGrassSide),
			East:   matRepo.Get(matRef.BlockGrassSide),
			West:   matRepo.Get(matRef.BlockGrassSide),
		},
	}

	ref[BlockDirt] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.BlockDirt),
			Bottom: matRepo.Get(matRef.BlockDirt),
			North:  matRepo.Get(matRef.BlockDirt),
			South:  matRepo.Get(matRef.BlockDirt),
			East:   matRepo.Get(matRef.BlockDirt),
			West:   matRepo.Get(matRef.BlockDirt),
		},
	}

	ref[BlockStone] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.BlockStone),
			Bottom: matRepo.Get(matRef.BlockStone),
			North:  matRepo.Get(matRef.BlockStone),
			South:  matRepo.Get(matRef.BlockStone),
			East:   matRepo.Get(matRef.BlockStone),
			West:   matRepo.Get(matRef.BlockStone),
		},
	}

	ref[BlockWater] = blockDef{
		Transparent: true,
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.BlockWater),
			Bottom: matRepo.Get(matRef.BlockWater),
			North:  matRepo.Get(matRef.BlockWater),
			South:  matRepo.Get(matRef.BlockWater),
			East:   matRepo.Get(matRef.BlockWater),
			West:   matRepo.Get(matRef.BlockWater),
		},
	}

	return ref
}
