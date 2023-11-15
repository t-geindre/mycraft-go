package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/material"
)

const (
	KindEmpty = iota
	KindOpaque
	KindTransparent
)

type BlockMaterials struct {
	Up    material.IMaterial
	Down  material.IMaterial
	North material.IMaterial
	South material.IMaterial
	West  material.IMaterial
	East  material.IMaterial
}

type Block struct {
	Id        string
	Type      string
	Kind      uint8
	Materials BlockMaterials
	*core.Node
}
