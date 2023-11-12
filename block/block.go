package block

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

type Block struct {
	Id         string
	Type       string
	CreateMesh func() *graphic.Mesh
	Meshes     []*graphic.Mesh
}

func (b *Block) Clone() *Block {
	meshes := make([]*graphic.Mesh, len(b.Meshes))
	for i, mesh := range b.Meshes {
		meshes[i] = mesh.Clone().(*graphic.Mesh)
	}

	return &Block{
		Id:         b.Id,
		Type:       b.Type,
		CreateMesh: b.CreateMesh,
		Meshes:     meshes,
	}
}

func (b *Block) SetPosition(pos math32.Vector3) {
	for _, mesh := range b.Meshes {
		mesh.SetPositionVec(&pos)
	}
}
