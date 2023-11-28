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

type Chunklet struct {
	*geometry.Geometry
	materialMap map[material.IMaterial]int
	repository  *block.Repository
}

func NewChunkletGeometry(chunk, east, west, north, south *world.Chunk, index float32) *Chunklet {
	c := new(Chunklet)
	c.repository = block.GetRepository()

	if chunk.AreLayersEmpty(int(index), int(index)+ChunkletSize) {
		return nil
	}

	quads := c.computeQuads(chunk, east, west, north, south, index)
	if len(quads) == 0 {
		return nil
	}

	c.computeGeometry(quads)

	return c
}

func (c *Chunklet) computeQuads(chunk, east, west, north, south *world.Chunk, index float32) []*Quad {
	quads := make([]*Quad, 0, 16*16*16)

	// Todo there must be a way to factorize this code
	// Todo translate x,y to center chunklet
	for x := float32(0); x < chunk.Size().X; x++ {
		for y := index; y < index+ChunkletSize; y++ {
			if chunk.IsLayerEmpty(int(y)) {
				continue
			}
			for z := float32(0); z < chunk.Size().Z; z++ {
				iX, iY, iZ := int(x), int(y), int(z)
				currentBlock := c.getBlockAt(chunk, iX, iY, iZ)

				if currentBlock == nil {
					continue
				}

				var topBlock *block.Block
				if y < chunk.Size().Y-1 {
					topBlock = c.getBlockAt(chunk, iX, iY+1, iZ)
				}
				if topBlock == nil || (topBlock.Transparent && topBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x, Y: y + 1 - index, Z: z},
						QuadFaceUp,
						currentBlock.Materials.Top,
					))
				}

				var bottomBlock *block.Block
				if y > 0 {
					bottomBlock = c.getBlockAt(chunk, iX, iY-1, iZ)
				}
				if bottomBlock == nil || (bottomBlock.Transparent && bottomBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x, Y: y - index, Z: z},
						QuadFaceDown,
						currentBlock.Materials.Bottom,
					))
				}

				var southBlock *block.Block = nil
				if z == 0 {
					southBlock = c.getBlockAt(south, iX, iY, int(south.Size().Z)-1)
				} else {
					southBlock = c.getBlockAt(chunk, iX, iY, iZ-1)
				}
				if southBlock == nil || (southBlock.Transparent && southBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x, Y: y - index, Z: z},
						QuadFaceSouth,
						currentBlock.Materials.South,
					))
				}

				var northBlock *block.Block = nil
				if z == chunk.Size().Z-1 {
					northBlock = c.getBlockAt(north, iX, iY, 0)
				} else {
					northBlock = c.getBlockAt(chunk, iX, iY, iZ+1)
				}
				if northBlock == nil || (northBlock.Transparent && northBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x, Y: y - index, Z: z + 1},
						QuadFaceNorth,
						currentBlock.Materials.North,
					))
				}

				var westBlock *block.Block = nil
				if x == 0 {
					westBlock = c.getBlockAt(west, int(west.Size().X)-1, iY, iZ)
				} else {
					westBlock = c.getBlockAt(chunk, iX-1, iY, iZ)
				}
				if westBlock == nil || (westBlock.Transparent && westBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x, Y: y - index, Z: z},
						QuadFaceEast,
						currentBlock.Materials.East,
					))
				}

				var eastBlock *block.Block = nil
				if x == chunk.Size().X-1 {
					eastBlock = c.getBlockAt(east, 0, iY, iZ)
				} else {
					eastBlock = c.getBlockAt(chunk, iX+1, iY, iZ)
				}
				if eastBlock == nil || (eastBlock.Transparent && eastBlock.Id != currentBlock.Id) {
					quads = append(quads, NewQuad(
						math32.Vector3{X: x + 1, Y: y - index, Z: z},
						QuadFaceWest,
						currentBlock.Materials.West,
					))
				}
			}
		}
	}

	return quads
}

func (c *Chunklet) computeGeometry(quads []*Quad) {
	c.Geometry = geometry.NewGeometry()

	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	quadsByMaterial := make(map[material.IMaterial][]*Quad)
	for _, quad := range quads {
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

func (c *Chunklet) getBlockAt(chunk *world.Chunk, x, y, z int) *block.Block {
	b := chunk.GetBlockAt(x, y, z)

	if b == block.BlockNone {
		return nil
	}

	return c.repository.Get(b)
}

func (c *Chunklet) MaterialMap() map[material.IMaterial]int {
	return c.materialMap
}
