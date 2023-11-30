package biome

import "mycraft/world/block"

type Water struct {
	rangeFrom float32
	rangeTo   float32
	ground    float32
}

func NewWater(rangeFrom, rangeTo float32) *Water {
	w := new(Water)
	w.rangeFrom = rangeFrom
	w.rangeTo = rangeTo

	return w
}

func (w *Water) SetGround(ground float32) {
	w.ground = ground
}

func (w *Water) Match(ground float32) bool {
	return ground >= w.rangeFrom && ground <= w.rangeTo
}

func (w *Water) GetBlockAt(x, y, z float32) uint16 {
	if y > w.rangeTo {
		return block.TypeNone
	}

	if y > w.ground {
		return block.TypeWater
	}

	if y == w.ground {
		if w.rangeTo-w.ground > 10 {
			return block.TypeGravel
		}
		if w.rangeTo+1 == w.ground {
			return block.TypeWater
		}
	}

	return block.TypeSand
}
