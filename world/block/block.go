package block

import "github.com/g3n/engine/material"

const (
	TypeNone = iota
	TypeSand
	TypeGrass
	TypeDirt
	TypeStone
	TypeWater
	TypeBedrock
	TypeDandelion
)

const (
	KindCube = iota
	KindPlant
)

const (
	MaterialTop = iota
	MaterialBottom
	MaterialNorth
	MaterialSouth
	MaterialEast
	MaterialWest
)

type Block struct {
	id          uint16
	kind        uint16
	transparent bool
	materials   map[uint16]material.IMaterial
}

func NewBlock(id uint16, kind uint16) *Block {
	return &Block{
		id:          id,
		kind:        kind,
		transparent: false,
		materials:   make(map[uint16]material.IMaterial),
	}
}

func (b *Block) IsSame(block *Block) bool {
	return b.id == block.id
}

func (b *Block) GetId() uint16 {
	return b.id
}

func (b *Block) GetKind() uint16 {
	return b.kind
}

func (b *Block) IsTransparent() bool {
	return b.transparent
}

func (b *Block) SetTransparent(t bool) {
	b.transparent = t
}

func (b *Block) SetMaterial(m material.IMaterial, ids ...uint16) {
	for _, id := range ids {
		b.materials[id] = m
	}
}

func (b *Block) GetMaterial(id uint16) material.IMaterial {
	if mat, ok := b.materials[id]; ok {
		return mat
	}

	panic("Material not found")
}
