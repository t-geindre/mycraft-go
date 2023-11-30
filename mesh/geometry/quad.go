package geometry

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

const (
	QuadFaceUp = iota
	QuadFaceDown
	QuadFaceNorth
	QuadFaceSouth
	QuadFaceWest
	QuadFaceEast
)

type Quad struct {
	position    math32.Vector3
	orientation uint8
	material    material.IMaterial
	size        math32.Vector2
}

func NewQuad(position math32.Vector3, orientation uint8, material material.IMaterial) *Quad {
	q := new(Quad)
	q.position = position
	q.orientation = orientation
	q.material = material
	q.size = math32.Vector2{X: 1, Y: 1}

	return q
}

func (q *Quad) Orientation() uint8 {
	return q.orientation
}

func (q *Quad) Material() material.IMaterial {
	return q.material
}

func (q *Quad) Vertices() []float32 {
	switch q.orientation {
	case QuadFaceUp, QuadFaceDown:
		return []float32{
			q.position.X + q.size.X, q.position.Y, q.position.Z,
			q.position.X, q.position.Y, q.position.Z,
			q.position.X, q.position.Y, q.position.Z + q.size.Y,
			q.position.X + q.size.X, q.position.Y, q.position.Z + q.size.Y,
		}
	case QuadFaceNorth, QuadFaceSouth:
		return []float32{
			q.position.X + q.size.X, q.position.Y + q.size.Y, q.position.Z,
			q.position.X, q.position.Y + q.size.Y, q.position.Z,
			q.position.X, q.position.Y, q.position.Z,
			q.position.X + q.size.X, q.position.Y, q.position.Z,
		}
	case QuadFaceWest, QuadFaceEast:
		return []float32{
			q.position.X, q.position.Y + q.size.Y, q.position.Z,
			q.position.X, q.position.Y + q.size.Y, q.position.Z + q.size.X,
			q.position.X, q.position.Y, q.position.Z + q.size.X,
			q.position.X, q.position.Y, q.position.Z,
		}
	}

	panic("Invalid quad orientation")
}

func (q *Quad) Normals() []float32 {
	switch q.orientation {
	case QuadFaceUp:
		return []float32{
			0, 1, 0,
			0, 1, 0,
			0, 1, 0,
			0, 1, 0,
		}
	case QuadFaceDown:
		return []float32{
			0, -1, 0,
			0, -1, 0,
			0, -1, 0,
			0, -1, 0,
		}
	case QuadFaceNorth:
		return []float32{
			0, 0, -1,
			0, 0, -1,
			0, 0, -1,
			0, 0, -1,
		}
	case QuadFaceSouth:
		return []float32{
			0, 0, 1,
			0, 0, 1,
			0, 0, 1,
			0, 0, 1,
		}
	case QuadFaceWest:
		return []float32{
			-1, 0, 0,
			-1, 0, 0,
			-1, 0, 0,
			-1, 0, 0,
		}
	case QuadFaceEast:
		return []float32{
			1, 0, 0,
			1, 0, 0,
			1, 0, 0,
			1, 0, 0,
		}
	}

	panic("Invalid quad orientation")
}

func (q *Quad) Indices(offset uint32) []uint32 {
	switch q.Orientation() {
	case QuadFaceUp, QuadFaceNorth, QuadFaceWest:
		return []uint32{
			offset, offset + 1, offset + 2,
			offset, offset + 2, offset + 3,
		}
	case QuadFaceDown, QuadFaceSouth, QuadFaceEast:
		return []uint32{
			offset, offset + 2, offset + 1,
			offset, offset + 3, offset + 2,
		}
	}

	panic("Invalid quad orientation")
}

func (q *Quad) Uvs() []float32 {
	return []float32{
		q.size.X, q.size.Y,
		0, q.size.Y,
		0, 0,
		q.size.X, 0,
	}
}
