package mod

import (
	"mycraft/world/block"
)

type OakTree struct {
	BaseMod
}

func NewOakTree(trunkHeight, leafRadius, x, y, z float32) *OakTree {
	o := new(OakTree)
	o.BaseMod = *NewBaseMod()

	tTarget := y + trunkHeight
	for ; y < tTarget; y++ {
		o.BaseMod.AddBlock(x, y, z, block.TypeOakLog)
	}

	lTarget := y - leafRadius
	for ; y > lTarget; y-- {
		o.BaseMod.AddBlock(x, y, z, block.TypeOakLeaves)
	}

	return o
}
