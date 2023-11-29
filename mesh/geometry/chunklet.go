package geometry

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"mycraft/world"
	"mycraft/world/block"
)

const ChunkletSize = 16
const BlockSize = 1

type Chunklet struct {
	*geometry.Geometry
	materialMap map[material.IMaterial]int
	repository  *block.Repository
	chunk       *world.Chunk
	chunkEast   *world.Chunk
	chunkWest   *world.Chunk
	chunkNorth  *world.Chunk
	chunkSouth  *world.Chunk
	yIndex      float32
	quads       []*Quad
}

func NewChunkletGeometry(chunk, east, west, north, south *world.Chunk, index float32) (*geometry.Geometry, map[material.IMaterial]int) {
	if chunk.AreLayersEmpty(int(index), int(index)+ChunkletSize) {
		return nil, nil
	}

	c := new(Chunklet)
	c.repository = block.GetRepository()
	c.chunk = chunk
	c.chunkEast = east
	c.chunkWest = west
	c.chunkNorth = north
	c.chunkSouth = south
	c.yIndex = index
	c.quads = make([]*Quad, 0)

	c.computeQuads()

	if len(c.quads) == 0 {
		return nil, nil
	}

	c.computeGeometry()

	return c.Geometry, c.materialMap
}

func (c *Chunklet) computeQuads() {
	// Todo translate x,y to center chunklet
	for y := c.yIndex; y < c.yIndex+ChunkletSize; y++ {
		if c.chunk.IsLayerEmpty(int(y)) {
			continue
		}
		for x := float32(0); x < c.chunk.Size().X; x++ {
			for z := float32(0); z < c.chunk.Size().Z; z++ {
				iX, iY, iZ := int(x), int(y), int(z)
				b := c.getBlock(iX, iY, iZ)

				if b == nil {
					continue
				}

				faces := map[uint8]*block.Block{
					QuadFaceUp:    c.getBlock(iX, iY+1, iZ),
					QuadFaceDown:  c.getBlock(iX, iY-1, iZ),
					QuadFaceSouth: c.getBlock(iX, iY, iZ-1),
					QuadFaceNorth: c.getBlock(iX, iY, iZ+1),
					QuadFaceEast:  c.getBlock(iX-1, iY, iZ),
					QuadFaceWest:  c.getBlock(iX+1, iY, iZ),
				}

				for face, adjB := range faces {
					if adjB == nil || (adjB.IsTransparent() && !adjB.IsSame(b)) {
						var pos math32.Vector3
						var mat material.IMaterial

						switch face {
						case QuadFaceUp:
							pos = math32.Vector3{X: x, Y: y + BlockSize - c.yIndex, Z: z}
							mat = b.GetMaterial(block.MaterialTop)
						case QuadFaceDown:
							pos = math32.Vector3{X: x, Y: y - c.yIndex, Z: z}
							mat = b.GetMaterial(block.MaterialBottom)
						case QuadFaceSouth:
							pos = math32.Vector3{X: x, Y: y - c.yIndex, Z: z}
							mat = b.GetMaterial(block.MaterialSouth)
						case QuadFaceNorth:
							pos = math32.Vector3{X: x, Y: y - c.yIndex, Z: z + BlockSize}
							mat = b.GetMaterial(block.MaterialNorth)
						case QuadFaceEast:
							pos = math32.Vector3{X: x, Y: y - c.yIndex, Z: z}
							mat = b.GetMaterial(block.MaterialEast)
						case QuadFaceWest:
							pos = math32.Vector3{X: x + BlockSize, Y: y - c.yIndex, Z: z}
							mat = b.GetMaterial(block.MaterialWest)
						default:
							panic("unknown quad face")
						}

						c.quads = append(c.quads, NewQuad(pos, face, mat))
					}
				}
			}
		}
	}
}

func (c *Chunklet) computeGeometry() {
	c.Geometry = geometry.NewGeometry()

	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	quadsByMaterial := make(map[material.IMaterial][]*Quad)
	for _, quad := range c.quads {
		if _, ok := quadsByMaterial[quad.Material()]; !ok {
			quadsByMaterial[quad.Material()] = make([]*Quad, 0)
		}
		quadsByMaterial[quad.Material()] = append(quadsByMaterial[quad.Material()], quad)
	}

	c.materialMap = make(map[material.IMaterial]int)
	offset := uint32(0)
	matGroupId := 0
	for mat, quads := range quadsByMaterial {
		c.materialMap[mat] = matGroupId
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
		c.Geometry.AddGroup(indicesStart, indicesLen, matGroupId)
		matGroupId++
	}

	c.Geometry.SetIndices(indices)
	c.Geometry.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	c.Geometry.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	c.Geometry.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
}

func (c *Chunklet) getBlock(x, y, z int) *block.Block {
	if y < 0 || y >= int(c.chunk.Size().Y) {
		return nil
	}
	if x < 0 {
		return c.getBlockAt(c.chunkWest, int(c.chunkWest.Size().X)+x, y, z)
	}
	if x >= int(c.chunk.Size().X) {
		return c.getBlockAt(c.chunkEast, x-int(c.chunk.Size().X), y, z)
	}
	if z < 0 {
		return c.getBlockAt(c.chunkSouth, x, y, int(c.chunkSouth.Size().Z)+z)
	}
	if z >= int(c.chunk.Size().Z) {
		return c.getBlockAt(c.chunkNorth, x, y, z-int(c.chunk.Size().Z))
	}

	return c.getBlockAt(c.chunk, x, y, z)
}

func (c *Chunklet) getBlockAt(chunk *world.Chunk, x, y, z int) *block.Block {
	b := chunk.GetBlockAt(x, y, z)
	if b == block.TypeNone {
		return nil
	}
	return c.repository.Get(b)
}
