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
		return block.BlockSand
	}
	return block.BlockNone
}

func (d *Desert) Reset() {
}
