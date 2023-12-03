package biome

import (
	"mycraft/world/block"
	"mycraft/world/chunk"
	"mycraft/world/generator/mod"
)

type Water struct {
	waterLvl float32
}

func NewWater(waterLvl float32) *Water {
	w := new(Water)
	w.waterLvl = waterLvl

	return w
}

func (w *Water) FillGround(chunk *chunk.Chunk, ground, x, z float32) {
	for y := w.waterLvl; y > 0; y-- {
		chunk.SetBlockAtF(x, y, z, w.getBlockAt(ground, x, y, z))
	}
}

func (w *Water) getBlockAt(ground, x, y, z float32) uint8 {
	if y > w.waterLvl {
		return block.TypeNone
	}

	if y > ground {
		return block.TypeWater
	}

	if y == ground {
		if w.waterLvl-ground > 10 {
			return block.TypeGravel
		}
		if w.waterLvl+1 == ground {
			return block.TypeWater
		}
	}

	return block.TypeSand
}

func (w *Water) GetMod(ground, x, z float32) mod.Mod {
	return nil
}
