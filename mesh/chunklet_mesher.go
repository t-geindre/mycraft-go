package mesh

import (
	geometry2 "github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"mycraft/mesh/geometry"
	"mycraft/world"
)

type ChunkletMesher struct {
	chunklet         *world.Chunklet
	quads            []*geometry.Quad
	quadsOptimized   []*geometry.Quad
	quadsByMatOrient map[material.IMaterial]map[uint8][]*geometry.Quad
}

func NewChunkletMesher(chunk *world.Chunklet) *ChunkletMesher {
	cm := new(ChunkletMesher)
	cm.chunklet = chunk

	return cm
}

func (cm *ChunkletMesher) ComputeQuads() {
	cm.quads = make([]*geometry.Quad, 0)
	cm.quadsOptimized = make([]*geometry.Quad, 0)
	cm.quadsByMatOrient = make(map[material.IMaterial]map[uint8][]*geometry.Quad)

	for x := 0; x < cm.chunklet.Size; x++ {
		for y := 0; y < cm.chunklet.Size; y++ {
			for z := 0; z < cm.chunklet.Size; z++ {
				if cm.chunklet.Blocks[x][y][z] == nil {
					continue
				}

				cX, cY, cZ := float32(x), float32(y), float32(z)

				// todo check neighbors

				if y == cm.chunklet.Size-1 || cm.chunklet.Blocks[x][y+1][z] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX, Y: cY + 1, Z: cZ},
						geometry.QuadFaceUp,
						cm.chunklet.Blocks[x][y][z].Materials.Top,
					)
				}

				if y == 0 || cm.chunklet.Blocks[x][y-1][z] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX, Y: cY, Z: cZ},
						geometry.QuadFaceDown,
						cm.chunklet.Blocks[x][y][z].Materials.Bottom,
					)
				}

				if z == 0 || cm.chunklet.Blocks[x][y][z-1] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX, Y: cY, Z: cZ},
						geometry.QuadFaceSouth,
						cm.chunklet.Blocks[x][y][z].Materials.South,
					)
				}

				if z == cm.chunklet.Size-1 || cm.chunklet.Blocks[x][y][z+1] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX, Y: cY, Z: cZ + 1},
						geometry.QuadFaceNorth,
						cm.chunklet.Blocks[x][y][z].Materials.North,
					)
				}

				if x == 0 || cm.chunklet.Blocks[x-1][y][z] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX, Y: cY, Z: cZ},
						geometry.QuadFaceEast,
						cm.chunklet.Blocks[x][y][z].Materials.East,
					)
				}

				if x == cm.chunklet.Size-1 || cm.chunklet.Blocks[x+1][y][z] == nil {
					cm.AddNewQuad(
						math32.Vector3{X: cX + 1, Y: cY, Z: cZ},
						geometry.QuadFaceWest,
						cm.chunklet.Blocks[x][y][z].Materials.West,
					)
				}
			}
		}
	}

	cm.MergeQuads()
}

func (cm *ChunkletMesher) AddNewQuad(pos math32.Vector3, orientation uint8, material material.IMaterial) {
	cm.quads = append(cm.quads, geometry.NewQuad(pos, orientation, material))

	if _, ok := cm.quadsByMatOrient[material]; !ok {
		cm.quadsByMatOrient[material] = make(map[uint8][]*geometry.Quad)
	}

	if _, ok := cm.quadsByMatOrient[material][orientation]; !ok {
		cm.quadsByMatOrient[material][orientation] = make([]*geometry.Quad, 0)
	}

	cm.quadsByMatOrient[material][orientation] = append(cm.quadsByMatOrient[material][orientation], cm.quads[len(cm.quads)-1])
}

func (cm *ChunkletMesher) MergeQuads() {
	for _, quadsByOrient := range cm.quadsByMatOrient {
		for _, quads := range quadsByOrient {
			cm.quadsOptimized = append(cm.quadsOptimized, cm.mergeQuads(quads)...)
		}
	}
}

func (cm *ChunkletMesher) mergeQuads(quads []*geometry.Quad) []*geometry.Quad {
	for {

		optimized := false

		for i := 0; i < len(quads); i++ {
			for j := i + 1; j < len(quads); j++ {
				iQuad := quads[i]
				jQuad := quads[j]
				iPos := iQuad.Position()
				jPos := jQuad.Position()
				iSize := iQuad.Size()
				jSize := jQuad.Size()

				switch iQuad.Orientation() {
				case geometry.QuadFaceUp, geometry.QuadFaceDown:
					if iPos.X+iSize.X == jPos.X && iPos.Z == jPos.Z && iPos.Y == jPos.Y && iSize.Y == jSize.Y {
						iQuad.SetSize(math32.Vector2{X: iSize.X + jSize.X, Y: iSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true

					}
					if iPos.Z+iSize.Y == jPos.Z && iPos.X == jPos.X && iPos.Y == jPos.Y && iSize.X == jSize.X {
						iQuad.SetSize(math32.Vector2{X: iSize.X, Y: iSize.Y + jSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true

					}
				case geometry.QuadFaceNorth, geometry.QuadFaceSouth:
					if iPos.X+iSize.X == jPos.X && iPos.Y == jPos.Y && iPos.Z == jPos.Z && iSize.Y == jSize.Y {
						iQuad.SetSize(math32.Vector2{X: iSize.X + jSize.X, Y: iSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true
					}
					if iPos.Y+iSize.Y == jPos.Y && iPos.X == jPos.X && iPos.Z == jPos.Z && iSize.X == jSize.X {
						iQuad.SetSize(math32.Vector2{X: iSize.X, Y: iSize.Y + jSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true
					}
				case geometry.QuadFaceWest, geometry.QuadFaceEast:
					if iPos.Z+iSize.X == jPos.Z && iPos.Y == jPos.Y && iPos.X == jPos.X && iSize.Y == jSize.Y {
						iQuad.SetSize(math32.Vector2{X: iSize.X + jSize.X, Y: iSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true
					}
					if iPos.Y+iSize.Y == jPos.Y && iPos.Z == jPos.Z && iPos.X == jPos.X && iSize.X == jSize.X {
						iQuad.SetSize(math32.Vector2{X: iSize.X, Y: iSize.Y + jSize.Y})
						quads = append(quads[:j], quads[j+1:]...)
						optimized = true
					}
				}
			}
		}

		if !optimized {
			break
		}
	}

	return quads
}

func (cm *ChunkletMesher) GetMesh() *graphic.Mesh {
	geo := geometry2.NewGeometry()

	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	materialMap := make(map[material.IMaterial][]int)

	offset := uint32(0)
	matGroupId := 0
	for _, quad := range cm.quadsOptimized {
		if quad == nil {
			continue
		}

		if _, ok := materialMap[quad.Material()]; !ok {
			materialMap[quad.Material()] = make([]int, 0)
		}

		materialMap[quad.Material()] = append(materialMap[quad.Material()], matGroupId)

		quadIndices := quad.Indices(offset)
		geo.AddGroup(indices.Len(), len(quadIndices), matGroupId)

		positions.Append(quad.Vertices()...)
		normals.Append(quad.Normals()...)
		uvs.Append(quad.Uvs()...)
		indices.Append(quadIndices...)

		offset = uint32(positions.Len() / 3)
		matGroupId++
	}

	geo.SetIndices(indices)
	geo.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	geo.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	geo.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	mesh := graphic.NewMesh(geo, nil)
	for mat, groupIds := range materialMap {
		for _, groupId := range groupIds {
			mesh.AddGroupMaterial(mat, groupId)
		}
	}

	mesh.SetPositionVec(cm.chunklet.Position)

	return mesh
}
