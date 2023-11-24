package mesh

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"mycraft/world"
)

type WorldMesher struct {
	chunks         map[math32.Vector2]*world.Chunk
	meshes         map[math32.Vector3]*Chunklet
	meshesQueue    *ChunkletQueue
	container      *core.Node
	renderDistance float32 // chunks horizontal, chunklet vertical
	lastPos        *math32.Vector3
}

const meshQueryPackSize = 20

const (
	chunkCenter = iota
	chunkEast
	chunkWest
	chunkNorth
	chunkSouth
)

func NewWorldMesher(rd float32) *WorldMesher {
	// todo add a meshing go routine to avoid rendering lag
	wm := new(WorldMesher)
	wm.chunks = make(map[math32.Vector2]*world.Chunk)
	wm.meshes = make(map[math32.Vector3]*Chunklet)
	wm.container = core.NewNode()
	wm.renderDistance = rd

	wm.meshesQueue = NewChunkletQueue(20, 1, func(chunklet *Chunklet) {
		if chunklet == nil {
			return
		}
		wm.meshes[chunklet.Position()] = chunklet
		wm.container.Add(chunklet)
	})

	return wm
}

func (wm *WorldMesher) Update(pos math32.Vector3) {
	wm.meshesQueue.Pull()

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
		requiredChunks := [...]math32.Vector2{
			chunkCenter: {X: meshPos.X, Y: meshPos.Z},
			chunkEast:   {X: meshPos.X + world.ChunkWith, Y: meshPos.Z},
			chunkWest:   {X: meshPos.X - world.ChunkWith, Y: meshPos.Z},
			chunkNorth:  {X: meshPos.X, Y: meshPos.Z + world.ChunkDepth},
			chunkSouth:  {X: meshPos.X, Y: meshPos.Z - world.ChunkDepth},
		}

		for _, chunkPos := range requiredChunks {
			if _, ok := wm.chunks[chunkPos]; !ok {
				continue MeshLoop
			}
		}

		// avoid further generation query
		wm.meshes[meshPos] = nil

		wm.meshesQueue.Push(ChunkletQuery{
			Pos:    meshPos,
			Center: wm.chunks[requiredChunks[chunkCenter]],
			East:   wm.chunks[requiredChunks[chunkEast]],
			West:   wm.chunks[requiredChunks[chunkWest]],
			North:  wm.chunks[requiredChunks[chunkNorth]],
			South:  wm.chunks[requiredChunks[chunkSouth]],
		})
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
	missing := make([]math32.Vector3, 0, 100)
	for x := pos.X - wm.renderDistance; x <= pos.X+wm.renderDistance; x += world.ChunkWith {
		for z := pos.Z - wm.renderDistance; z <= pos.Z+wm.renderDistance; z += world.ChunkDepth {
			for y := math32.Max(0, pos.Y-wm.renderDistance); y < math32.Min(pos.Y+wm.renderDistance, world.ChunkHeight); y += ChunkletSize {
				newPos := math32.Vector3{X: x, Y: y, Z: z}
				if _, ok := wm.meshes[newPos]; ok {
					continue
				}
				missing = append(missing, newPos)
			}
		}
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
