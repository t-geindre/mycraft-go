package biome

import (
	"mycraft/world/block"
)

type Desert struct {
}

func NewDesert() *Desert {
	d := new(Desert)

	return d
}

func (d *Desert) GetBlockAt(x, y, z, strength float32) uint16 {
	if y < 10 {
		return block.TypeSand
	}
	return block.TypeNone
}

func (d *Desert) Reset() {
}
