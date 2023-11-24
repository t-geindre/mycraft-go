package infinite

import (
	"github.com/g3n/engine/math32"
	"github.com/ojrac/opensimplex-go"
	"mycraft/world"
	"mycraft/world/block"
)

type Noise struct {
	seed  int64
	noise opensimplex.Noise32
}

func NewNoiseGenerator(seed int64) *Noise {
	n := new(Noise)
	n.seed = seed

	// Normalized noise will return 0-1 values
	n.noise = opensimplex.NewNormalized32(n.seed)
	return n
}

func (n *Noise) Populate(chunk *world.Chunk) {
	blockRepository := block.GetRepository()
	for x := 0; x < world.ChunkWith; x++ {
		for z := 0; z < world.ChunkDepth; z++ {
			cX, cZ := chunk.Position().X+float32(x), chunk.Position().Y+float32(z)
			groundLevel := n.GetGroundLevelAt(cX, cZ)
			for y := 0; y < world.ChunkHeight; y++ {
				cY := float32(y)

				if cY == groundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.BlockGrass))
					continue
				}
				if cY < groundLevel {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.BlockDirt))
					continue
				}

				if cY > groundLevel && cY < 90 {
					chunk.SetBlockAt(x, y, z, blockRepository.Get(block.BlockWater))
				}
			}
		}
	}
}

func (n *Noise) GetGroundLevelAt(x, z float32) float32 {
	noise := float32(50)
	pvNoise := n.noise.Eval2(x/200, z/200)

	noise += n.noise.Eval2(x/50, z/50) * 100 * pvNoise
	noise += n.noise.Eval2(x/5000, z/5000) * 100 * (1 - pvNoise)

	noise += n.noise.Eval2(x/5, z/5) * 2

	return math32.Floor(noise)
}
