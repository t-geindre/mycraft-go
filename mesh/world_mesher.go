package mesh

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"mycraft/world"
)

type WorldMesher struct {
	chunks         map[math32.Vector2]*world.Chunk
	meshes         map[math32.Vector3]*Chunklet
	container      *core.Node
	renderDistance float32 // chunks horizontal, chunklet vertical
	lastPos        *math32.Vector3
}

func NewWorldMesher(rd float32) *WorldMesher {
	// todo add a meshing go routine to avoid rendering lag
	wm := new(WorldMesher)
	wm.chunks = make(map[math32.Vector2]*world.Chunk)
	wm.meshes = make(map[math32.Vector3]*Chunklet)
	wm.container = core.NewNode()
	wm.renderDistance = rd
	return wm
}

func (wm *WorldMesher) Update(pos math32.Vector3) {
	pos = wm.getWorldPosition(pos)
	if wm.lastPos != nil && pos.Equals(wm.lastPos) {
		return
	}
	wm.lastPos = &pos
	wm.doUpdate(*wm.lastPos)
}

func (wm *WorldMesher) doUpdate(pos math32.Vector3) {
	wm.addMissingMeshes(pos)
	wm.clearToFarMeshes(pos)
}

func (wm *WorldMesher) addMissingMeshes(pos math32.Vector3) {
MeshLoop:
	for _, meshPos := range wm.getMissingMeshesPos(pos) {
		requiredChunks := map[string]math32.Vector2{
			"center": {X: meshPos.X, Y: meshPos.Z},
			"east":   {X: meshPos.X + world.ChunkWith, Y: meshPos.Z},
			"west":   {X: meshPos.X - world.ChunkWith, Y: meshPos.Z},
			"north":  {X: meshPos.X, Y: meshPos.Z + world.ChunkDepth},
			"south":  {X: meshPos.X, Y: meshPos.Z - world.ChunkDepth},
		}

		for _, chunkPos := range requiredChunks {
			if _, ok := wm.chunks[chunkPos]; !ok {
				continue MeshLoop
			}
		}

		mesh := NewChunklet(
			wm.chunks[requiredChunks["center"]],
			wm.chunks[requiredChunks["east"]],
			wm.chunks[requiredChunks["west"]],
			wm.chunks[requiredChunks["north"]],
			wm.chunks[requiredChunks["south"]],
			meshPos.Y,
		)
		wm.meshes[meshPos] = mesh
		if mesh == nil {
			continue
		}
		wm.container.Add(mesh)
	}
}

func (wm *WorldMesher) Container() *core.Node {
	return wm.container
}

func (wm *WorldMesher) AddChunk(chunk *world.Chunk) {
	wm.chunks[*chunk.Position()] = chunk
	wm.doUpdate(*wm.lastPos)
}

func (wm *WorldMesher) RemoveChunk(chunk *world.Chunk) {
	delete(wm.chunks, *chunk.Position())
}

func (wm *WorldMesher) getMissingMeshesPos(pos math32.Vector3) []math32.Vector3 {
	expectedPos := map[math32.Vector3]bool{}

	for westX, eastX := pos.X, pos.X; eastX <= pos.X+wm.renderDistance; westX, eastX = westX-world.ChunkWith, eastX+world.ChunkWith {
		for bottomY, topY := pos.Y, pos.Y; topY <= pos.Y+wm.renderDistance; bottomY, topY = bottomY-ChunkletSize, topY+ChunkletSize {
			for northZ, southZ := pos.Z, pos.Z; southZ <= pos.Z+wm.renderDistance; northZ, southZ = northZ-world.ChunkDepth, southZ+world.ChunkDepth {
				for _, meshPos := range []math32.Vector3{
					{X: eastX, Y: bottomY, Z: northZ},
					{X: eastX, Y: bottomY, Z: southZ},
					{X: eastX, Y: topY, Z: northZ},
					{X: eastX, Y: topY, Z: southZ},
					{X: westX, Y: bottomY, Z: northZ},
					{X: westX, Y: bottomY, Z: southZ},
					{X: westX, Y: topY, Z: northZ},
					{X: westX, Y: topY, Z: southZ},
				} {
					expectedPos[meshPos] = true
				}
			}
		}
	}

	missing := make([]math32.Vector3, 0)
	for pos := range expectedPos {
		if _, ok := wm.meshes[pos]; ok || pos.Y >= world.ChunkHeight || pos.Y < 0 {
			delete(expectedPos, pos)
			continue
		}
		missing = append(missing, pos)
	}

	return missing
}

func (wm *WorldMesher) getWorldPosition(pos math32.Vector3) math32.Vector3 {
	return math32.Vector3{
		X: math32.Floor(pos.X/world.ChunkWith) * world.ChunkWith,
		Y: math32.Floor(pos.Y/ChunkletSize) * ChunkletSize,
		Z: math32.Floor(pos.Z/world.ChunkDepth) * world.ChunkDepth,
	}
}

func (wm *WorldMesher) clearToFarMeshes(pos math32.Vector3) {
	meshCap := math32.Pow(wm.renderDistance/ChunkletSize, 3)
	if len(wm.container.Children()) < int(meshCap) {
		return
	}

	for meshPos, mesh := range wm.meshes {
		// todo clear farest first until meshcap is reached
		if math32.Abs(meshPos.X-pos.X) > wm.renderDistance ||
			math32.Abs(meshPos.Y-pos.Y) > wm.renderDistance ||
			math32.Abs(meshPos.Z-pos.Z) > wm.renderDistance {
			if mesh != nil {
				wm.container.Remove(mesh)
				mesh.GetGeometry().Dispose()
			}
			delete(wm.meshes, meshPos)
		}
	}
}
