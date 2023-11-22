package infinite

import (
	"github.com/g3n/engine/math32"
	"mycraft/world"
	"mycraft/world/block"
)

type Sin struct {
	amplitude float32
	period    float32
}

func NewSinGenerator(amplitude, period float32) *Sin {
	i := new(Sin)
	i.amplitude = amplitude
	i.period = period
	return i
}

func (i *Sin) Populate(chunk *world.Chunk) {
	blockRepository := block.GetRepository()
	for x := 0; x < world.ChunkWith; x++ {
		for z := 0; z < world.ChunkDepth; z++ {
			cX, cZ := chunk.Position().X+float32(x), chunk.Position().Y+float32(z)
			groundLevel := math32.Ceil(i.amplitude * (math32.Sin(cX/i.period) + math32.Sin(cZ/i.period) + 2))
			for y := 0; y < world.ChunkHeight; y++ {
				cY := float32(y)
				if cY < groundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.BlockDirt))
					continue
				}
				if cY == groundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.BlockGrass))
					continue
				}
			}
		}
	}
}
