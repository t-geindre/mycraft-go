package scene

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	engineMaterial "github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"mycraft/app"
	"mycraft/block"
	"time"
)

const (
	faceUp = iota
	faceDown
	faceNorth
	faceSouth
	faceWest
	faceEast
)

const chunkSize = 16

type Test struct {
	materialMap map[engineMaterial.IMaterial][]int
	geo         *geometry.Geometry
}

func (t *Test) Setup(container *core.Node, app *app.App) {
	var chunk struct {
		blocks [chunkSize][chunkSize][chunkSize]*block.Block
	}

	t.materialMap = make(map[engineMaterial.IMaterial][]int)

	grassBlock := block.GetRepository().Get("green_grass")
	tntBlock := block.GetRepository().Get("tnt")
	dirtBlock := block.GetRepository().Get("dirt")

	for x := 0; x < chunkSize; x++ {
		for y := 0; y < chunkSize; y++ {
			for z := 0; z < chunkSize; z++ {
				if y == 2 || (y == 3 && x%2 == 0 && z%3 == 0) {
					chunk.blocks[x][y][z] = grassBlock
					chunk.blocks[x][y-1][z] = dirtBlock
					continue
				}
				chunk.blocks[x][y][z] = nil
			}
		}
	}

	chunk.blocks[8][8][8] = tntBlock

	container.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.2))

	orb := camera.NewOrbitControl(app.Cam)
	orb.SetTarget(*math32.NewVector3(chunkSize/2, 2, chunkSize/2))

	t.geo = geometry.NewGeometry()

	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	offset := uint32(0)
	for x := 0; x < chunkSize; x++ {
		for y := 0; y < chunkSize; y++ {
			for z := 0; z < chunkSize; z++ {
				if chunk.blocks[x][y][z] == nil {
					continue
				}

				cX, cY, cZ := float32(x), float32(y), float32(z)

				if y == chunkSize-1 || chunk.blocks[x][y+1][z] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceUp, chunk.blocks[x][y][z])
				}

				if y == 0 || chunk.blocks[x][y-1][z] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceDown, chunk.blocks[x][y][z])
				}

				if z == chunkSize-1 || chunk.blocks[x][y][z+1] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceNorth, chunk.blocks[x][y][z])
				}

				if z == 0 || chunk.blocks[x][y][z-1] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceSouth, chunk.blocks[x][y][z])
				}

				if x == chunkSize-1 || chunk.blocks[x+1][y][z] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceWest, chunk.blocks[x][y][z])
				}

				if x == 0 || chunk.blocks[x-1][y][z] == nil {
					t.addPlan(&positions, &normals, &uvs, &indices, cX, cY, cZ, &offset, faceEast, chunk.blocks[x][y][z])
				}
			}
		}
	}

	t.geo.SetIndices(indices)
	t.geo.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	t.geo.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	t.geo.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	mesh := graphic.NewMesh(t.geo, nil) //material.GetRepository().Get("tnt_side"))

	for mat, groupIds := range t.materialMap {
		for _, groupId := range groupIds {
			mesh.AddGroupMaterial(mat, groupId)
		}
	}

	mesh.SetPosition(0, 0, -5)

	container.Add(mesh)
	// vertex = sommet, plural vertices
	// Indices are numeric labels for the vertices in a given 3D scene. Indices allow us to tell WebGL how to connect vertices in order to produce a surface.
	// A normal in general is a unit vector whose direction is perpendicular to a surface at a specific point. Therefore it tells you in which direction a surface is facing. The main use case for normals are lighting calculations
	// uvs texture coordinates

}

func (t *Test) Update(deltaTime time.Duration) {
}

func (t *Test) Dispose() {
}

func NewTestScene() *Test {
	return new(Test)
}

func (t *Test) addPlan(
	positions, normals, uvs *math32.ArrayF32,
	indices *math32.ArrayU32,
	x, y, z float32,
	offset *uint32,
	orientation uint8,
	block *block.Block,
) {
	switch orientation {
	case faceUp:
		positions.Append(
			x-1, y, z,
			x, y, z,
			x, y, z-1,
			x-1, y, z-1,
		)
	case faceDown:
		positions.Append(
			x-1, y-1, z,
			x, y-1, z,
			x, y-1, z-1,
			x-1, y-1, z-1,
		)
	case faceNorth:
		positions.Append(
			x-1, y-1, z,
			x, y-1, z,
			x, y, z,
			x-1, y, z,
		)
	case faceSouth:
		positions.Append(
			x-1, y-1, z-1,
			x, y-1, z-1,
			x, y, z-1,
			x-1, y, z-1,
		)
	case faceWest:
		positions.Append(
			x, y-1, z,
			x, y-1, z-1,
			x, y, z-1,
			x, y, z,
		)
	case faceEast:
		positions.Append(
			x-1, y-1, z,
			x-1, y-1, z-1,
			x-1, y, z-1,
			x-1, y, z,
		)
	}

	normals.Append(
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	)

	uvs.Append(
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	)

	switch orientation {
	case faceUp, faceNorth, faceWest:
		indices.Append(
			*offset, *offset+1, *offset+2,
			*offset, *offset+2, *offset+3,
		)
	case faceDown, faceSouth, faceEast:
		indices.Append(
			*offset, *offset+2, *offset+1,
			*offset, *offset+3, *offset+2,
		)
	}

	var mat engineMaterial.IMaterial
	switch orientation {
	case faceUp:
		mat = block.Materials.Up
	case faceDown:
		mat = block.Materials.Down
	case faceNorth:
		mat = block.Materials.North
	case faceSouth:
		mat = block.Materials.South
	case faceWest:
		mat = block.Materials.West
	case faceEast:
		mat = block.Materials.East
	}

	if _, ok := t.materialMap[mat]; !ok {
		t.materialMap[mat] = make([]int, 0)
	}

	groupId := int(*offset) / 4
	t.materialMap[mat] = append(t.materialMap[mat], groupId)
	t.geo.AddGroup(int(*offset)/4*6, 6, groupId)

	*offset += 4
}
