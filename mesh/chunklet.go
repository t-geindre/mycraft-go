package mesh

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"mycraft/mesh/geometry"
	"mycraft/world/chunk"
)

const ChunkletSize = 16

type Chunklet struct {
	*graphic.Mesh
}

func NewChunklet(chunk, east, west, north, south *chunk.Chunk, pos math32.Vector3) *Chunklet {
	geo, matMap := geometry.NewChunkletGeometry(chunk, east, west, north, south, pos.Y)
	if geo == nil {
		return nil
	}

	c := new(Chunklet)
	c.Mesh = graphic.NewMesh(geo, nil)
	c.Mesh.SetPositionVec(&pos)

	for mat, groupId := range matMap {
		c.Mesh.AddGroupMaterial(mat, groupId)
	}

	return c
}
