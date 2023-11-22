package mesh

import (
	"github.com/g3n/engine/graphic"
	"mycraft/mesh/geometry"
	"mycraft/world"
)

const ChunkletSize = 16

type Chunklet struct {
	*graphic.Mesh
}

func NewChunklet(chunk, east, west, north, south *world.Chunk, index float32) *Chunklet {
	geo := geometry.NewChunkletGeometry(chunk, east, west, north, south, index)
	if geo == nil {
		return nil
	}

	c := new(Chunklet)

	mesh := graphic.NewMesh(geo, nil)

	for mat, groupId := range geo.MaterialMap() {
		mesh.AddGroupMaterial(mat, groupId)
	}

	c.Mesh = mesh
	c.Mesh.SetPosition(chunk.Position().X, index, chunk.Position().Y)

	return c
}
