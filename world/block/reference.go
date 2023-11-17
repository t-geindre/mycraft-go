package block

import (
	"github.com/g3n/engine/material"
	matRef "mycraft/mesh/material"
)

const (
	GrassBlock = iota
	DirtBlock
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
	Materials blockDefMaterials
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
			Bottom: matRepo.Get(matRef.DirtBlock),
			North:  matRepo.Get(matRef.DirtBlock),
			South:  matRepo.Get(matRef.DirtBlock),
			East:   matRepo.Get(matRef.DirtBlock),
			West:   matRepo.Get(matRef.DirtBlock),
		},
	}

	return ref
}
