package biome

import "mycraft/world/block"

type Beach struct {
	rangeFrom float32
	rangeTo   float32
	ground    float32
}

func NewBeach(rangeFrom, rangeTo float32) *Beach {
	b := new(Beach)
	b.rangeFrom = rangeFrom
	b.rangeTo = rangeTo

	return b
}

func (b *Beach) SetGround(level float32) {
	b.ground = level
}

func (b *Beach) Match(level float32) bool {
	return level >= b.rangeFrom && level <= b.rangeTo
}

func (b *Beach) GetBlockAt(x, y, z float32) uint16 {
	if y == b.ground {
		return block.TypeSand
	}

	if y < b.ground {
		return block.TypeStone
	}

	return block.TypeNone
}
