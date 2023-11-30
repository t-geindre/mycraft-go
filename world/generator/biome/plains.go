package biome

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/block"
	"mycraft/world/generator/noise"
	"mycraft/world/generator/noise/normalized"
)

type Plains struct {
	rangeFrom float32
	rangeTo   float32
	ground    float32
	treeNoise noise.Noise
}

func NewPlains(rangeFrom, rangeTo float32, seed int64) *Plains {
	p := new(Plains)
	p.rangeFrom = rangeFrom
	p.rangeTo = rangeTo

	p.treeNoise = normalized.NewSimplexNoise(seed)
	p.treeNoise = noise.NewScale(p.treeNoise, 1)
	p.treeNoise = noise.NewAmplify(p.treeNoise, 100)

	return p
}

func (p *Plains) SetGround(level float32) {
	p.ground = level
}

func (p *Plains) Match(level float32) bool {
	return level >= p.rangeFrom && level <= p.rangeTo
}

func (p *Plains) GetBlockAt(x, y, z float32) uint16 {
	treeBlock := p.getTreeBlockAt(x, y, z)
	if treeBlock != block.TypeNone {
		return treeBlock
	}

	return p.getGroundBlockAt(x, y, z)
}

func (p *Plains) getGroundBlockAt(x, y, z float32) uint16 {
	if y == p.ground {
		return block.TypeGrass
	}

	if y < p.ground {
		return block.TypeDirt
	}

	return block.TypeNone
}

func (p *Plains) getTreeBlockAt(x, y, z float32) uint16 {
	treeSpawn := math32.Round(p.treeNoise.Eval2(x, z))
	trunk := float32(0)
	switch treeSpawn {
	case 100:
		trunk = 6
	}
	if y > p.ground && y < p.ground+trunk && trunk > 0 {
		return block.TypeSpruceLog
	}

	return block.TypeNone
}

/*
// generateTree génère un arbre à la position (x, y, z) en utilisant une fonction de bruit.
func generateTree(x, y, z int) {
    // Hauteur du tronc déterminée par une fonction de bruit
    trunkHeight := int(simplexNoise(x, y, z) * MaxTrunkHeight)

    // Placer le tronc
    for i := 0; i < trunkHeight; i++ {
        setBlock(x, y+i, z, BlockTypeLog)
    }

    // Générer le feuillage autour du sommet du tronc
    generateLeaves(x, y+trunkHeight, z, LeafRadius)
}

// generateLeaves génère le feuillage autour de la position (x, y, z) avec un certain rayon.
func generateLeaves(x, y, z, radius int) {
    for offsetX := -radius; offsetX <= radius; offsetX++ {
        for offsetY := -radius; offsetY <= radius; offsetY++ {
            for offsetZ := -radius; offsetZ <= radius; offsetZ++ {
                // Utiliser une fonction de bruit pour déterminer si un bloc de feuillage doit être placé à cette position
                if simplexNoise(x+offsetX, y+offsetY, z+offsetZ) > LeafThreshold {
                    setBlock(x+offsetX, y+offsetY, z+offsetZ, BlockTypeLeaves)
                }
            }
        }
    }
}
*/
