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

const ChunkletSize = 16

type Chunklet struct {
	*graphic.Mesh
	quads            []*geometry.Quad
	quadsOptimized   []*geometry.Quad
	quadsByMatOrient map[material.IMaterial]map[uint8][]*geometry.Quad
	centerChunk      *world.Chunk
	eastChunk        *world.Chunk
	westChunk        *world.Chunk
	northChunk       *world.Chunk
	southChunk       *world.Chunk
	index            float32
}

func NewChunklet(chunk, east, west, north, south *world.Chunk, index float32) *Chunklet {
	c := new(Chunklet)
	c.centerChunk = chunk
	c.eastChunk = east
	c.westChunk = west
	c.northChunk = north
	c.southChunk = south
	c.index = index
	c.computeQuads()

	if len(c.quads) == 0 {
		return nil
	}

	c.Mesh = c.GetMesh()
	c.Mesh.SetPosition(chunk.Position().X, index, chunk.Position().Y)

	return c
}

func (c *Chunklet) computeQuads() {
	c.quads = make([]*geometry.Quad, 0)
	c.quadsOptimized = make([]*geometry.Quad, 0)
	c.quadsByMatOrient = make(map[material.IMaterial]map[uint8][]*geometry.Quad)

	for x := float32(0); x < c.centerChunk.Size().X; x++ {
		for y := c.index; y < c.index+ChunkletSize; y++ {
			for z := float32(0); z < c.centerChunk.Size().Z; z++ {
				iX, iY, iZ := int(x), int(y), int(z)

				if c.centerChunk.GetBlockAt(iX, iY, iZ) == nil {
					continue
				}

				// todo check neighbors

				if y == c.centerChunk.Size().Y-1 || c.centerChunk.GetBlockAt(iX, iY+1, iZ) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x, Y: y + 1, Z: z},
						geometry.QuadFaceUp,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.Top,
					)
				}

				if y == 0 || c.centerChunk.GetBlockAt(iX, iY-1, iZ) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x, Y: y, Z: z},
						geometry.QuadFaceDown,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.Bottom,
					)
				}

				if z == 0 || c.centerChunk.GetBlockAt(iX, iY, iZ-1) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x, Y: y, Z: z},
						geometry.QuadFaceSouth,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.South,
					)
				}

				if z == c.centerChunk.Size().Z-1 || c.centerChunk.GetBlockAt(iX, iY, iZ+1) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x, Y: y, Z: z + 1},
						geometry.QuadFaceNorth,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.North,
					)
				}

				if x == 0 || c.centerChunk.GetBlockAt(iX-1, iY, iZ) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x, Y: y, Z: z},
						geometry.QuadFaceEast,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.East,
					)
				}

				if x == c.centerChunk.Size().X-1 || c.centerChunk.GetBlockAt(iX+1, iY, iZ) == nil {
					c.AddNewQuad(
						math32.Vector3{X: x + 1, Y: y, Z: z},
						geometry.QuadFaceWest,
						c.centerChunk.GetBlockAt(iX, iY, iZ).Materials.West,
					)
				}
			}
		}
	}

	//c.MergeQuads()
}

func (c *Chunklet) AddNewQuad(pos math32.Vector3, orientation uint8, material material.IMaterial) {
	c.quads = append(c.quads, geometry.NewQuad(pos, orientation, material))

	if _, ok := c.quadsByMatOrient[material]; !ok {
		c.quadsByMatOrient[material] = make(map[uint8][]*geometry.Quad)
	}

	if _, ok := c.quadsByMatOrient[material][orientation]; !ok {
		c.quadsByMatOrient[material][orientation] = make([]*geometry.Quad, 0)
	}

	c.quadsByMatOrient[material][orientation] = append(c.quadsByMatOrient[material][orientation], c.quads[len(c.quads)-1])
}

func (c *Chunklet) GetMesh() *graphic.Mesh {
	geo := geometry2.NewGeometry()

	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	materialMap := make(map[material.IMaterial]int)
	quadsByMaterial := make(map[material.IMaterial][]*geometry.Quad)
	for _, quad := range c.quads {
		if _, ok := quadsByMaterial[quad.Material()]; !ok {
			quadsByMaterial[quad.Material()] = make([]*geometry.Quad, 0)
		}
		quadsByMaterial[quad.Material()] = append(quadsByMaterial[quad.Material()], quad)
	}

	offset := uint32(0)
	matGroupId := 0
	for mat, quads := range quadsByMaterial {
		materialMap[mat] = matGroupId
		indicesLen := 0
		indicesStart := indices.Len()
		for _, quad := range quads {
			quadIndices := quad.Indices(offset)
			indicesLen += len(quadIndices)
			positions.Append(quad.Vertices()...)
			normals.Append(quad.Normals()...)
			uvs.Append(quad.Uvs()...)
			indices.Append(quadIndices...)
			offset = uint32(positions.Len() / 3)
		}
		geo.AddGroup(indicesStart, indicesLen, matGroupId)
		matGroupId++
	}

	geo.SetIndices(indices)
	geo.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	geo.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	geo.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	mesh := graphic.NewMesh(geo, nil)
	for mat, groupId := range materialMap {
		mesh.AddGroupMaterial(mat, groupId)
	}

	return mesh
}

/*
func (cm *ChunkMesher) mergeQuads(quads []*geometry.Quad) []*geometry.Quad {
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
*/
