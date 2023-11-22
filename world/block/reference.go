package block

import (
	"github.com/g3n/engine/material"
	matRef "mycraft/mesh/material"
)

const (
	GrassBlock = iota
	DirtBlock
	StoneBlock
	WaterBlock
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

	ref[GrassBlock] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.GrassBlockTop),
			Bottom: matRepo.Get(matRef.DirtBlock),
			North:  matRepo.Get(matRef.GrassBlockSide),
			South:  matRepo.Get(matRef.GrassBlockSide),
			East:   matRepo.Get(matRef.GrassBlockSide),
			West:   matRepo.Get(matRef.GrassBlockSide),
		},
	}

	ref[DirtBlock] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.DirtBlock),
			Bottom: matRepo.Get(matRef.WaterBlock),
			North:  matRepo.Get(matRef.DirtBlock),
			South:  matRepo.Get(matRef.DirtBlock),
			East:   matRepo.Get(matRef.DirtBlock),
			West:   matRepo.Get(matRef.DirtBlock),
		},
	}

	ref[StoneBlock] = blockDef{
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.StoneBlock),
			Bottom: matRepo.Get(matRef.StoneBlock),
			North:  matRepo.Get(matRef.StoneBlock),
			South:  matRepo.Get(matRef.StoneBlock),
			East:   matRepo.Get(matRef.StoneBlock),
			West:   matRepo.Get(matRef.StoneBlock),
		},
	}

	ref[WaterBlock] = blockDef{
		Transparent: true,
		Materials: blockDefMaterials{
			Top:    matRepo.Get(matRef.WaterBlock),
			Bottom: matRepo.Get(matRef.WaterBlock),
			North:  matRepo.Get(matRef.WaterBlock),
			South:  matRepo.Get(matRef.WaterBlock),
			East:   matRepo.Get(matRef.WaterBlock),
			West:   matRepo.Get(matRef.WaterBlock),
		},
	}

	return ref
}
