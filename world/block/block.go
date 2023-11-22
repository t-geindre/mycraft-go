package block

import "github.com/g3n/engine/material"

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
