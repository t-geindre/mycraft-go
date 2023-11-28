package mesh

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"mycraft/world"
	"sort"
	"time"
)

type WorldMesher struct {
	chunks            map[math32.Vector2]*world.Chunk
	meshes            map[math32.Vector3]*Chunklet
	world             *world.World
	renderDistance    float32 // chunks horizontal, chunklet vertical
	container         *core.Node
	lastPos           *math32.Vector3
	positionChannel   chan math32.Vector3
	addMeshChannel    chan []*Chunklet
	removeMeshChannel chan []*Chunklet
}

const (
	chunkCenter = iota
	chunkEast
	chunkWest
	chunkNorth
	chunkSouth
)

func NewWorldMesher(rd float32, w *world.World) *WorldMesher {
	wm := new(WorldMesher)

	wm.chunks = make(map[math32.Vector2]*world.Chunk)
	wm.meshes = make(map[math32.Vector3]*Chunklet)

	wm.positionChannel = make(chan math32.Vector3, 1)
	wm.addMeshChannel = make(chan []*Chunklet, 20)
	wm.removeMeshChannel = make(chan []*Chunklet, 20)

	wm.container = core.NewNode()
	wm.renderDistance = rd
	wm.world = w

	go wm.Run()

	return wm
}

func (wm *WorldMesher) Run() {
	for {
		select {
		case chunks := <-wm.world.AddChunkChannel():
			for _, chunk := range chunks {
				wm.AddChunk(chunk)
			}
			if wm.lastPos != nil {
				wm.doUpdate(*wm.lastPos)
			}
		case chunks := <-wm.world.RemoveChunkChannel():
			for _, chunk := range chunks {
				wm.RemoveChunk(chunk)
			}
		case pos, ok := <-wm.positionChannel:
			if !ok {
				return
			}
			if wm.lastPos != nil && pos.Equals(wm.lastPos) {
				break
			}
			wm.lastPos = &pos
			wm.doUpdate(pos)
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (wm *WorldMesher) Update(pos math32.Vector3) {
	pos = wm.getWorldPosition(pos)
	for len(wm.positionChannel) > 0 {
		// Still processing previous position
		// refresh the awaiting one
		<-wm.positionChannel
	}
	wm.positionChannel <- pos

	wm.UpdateMeshes()
}

func (wm *WorldMesher) UpdateMeshes() {
	for {
		select {
		case meshes := <-wm.addMeshChannel:
			for _, mesh := range meshes {
				wm.container.Add(mesh)
			}
		case meshes := <-wm.removeMeshChannel:
			for _, mesh := range meshes {
				mesh.GetGeometry().Dispose()
				wm.container.Remove(mesh)
			}
		default:
			return
		}
	}
}

func (wm *WorldMesher) doUpdate(pos math32.Vector3) {
	wm.addMissingMeshes(pos)
	wm.clearTooFarMeshes(pos)
}

func (wm *WorldMesher) addMissingMeshes(pos math32.Vector3) {
	toAdd := make([]*Chunklet, 0, 100)
MeshLoop:
	for _, meshPos := range wm.getMissingMeshesPos(pos) {
		requiredChunks := [...]math32.Vector2{
			chunkCenter: {X: meshPos.X, Y: meshPos.Z},
			chunkEast:   {X: meshPos.X + world.ChunkWidth, Y: meshPos.Z},
			chunkWest:   {X: meshPos.X - world.ChunkWidth, Y: meshPos.Z},
			chunkNorth:  {X: meshPos.X, Y: meshPos.Z + world.ChunkDepth},
			chunkSouth:  {X: meshPos.X, Y: meshPos.Z - world.ChunkDepth},
		}

		for _, chunkPos := range requiredChunks {
			if _, ok := wm.chunks[chunkPos]; !ok {
				continue MeshLoop
			}
		}

		wm.meshes[meshPos] = NewChunklet(
			wm.chunks[requiredChunks[chunkCenter]],
			wm.chunks[requiredChunks[chunkEast]],
			wm.chunks[requiredChunks[chunkWest]],
			wm.chunks[requiredChunks[chunkNorth]],
			wm.chunks[requiredChunks[chunkSouth]],
			meshPos,
		)

		if wm.meshes[meshPos] != nil {
			toAdd = append(toAdd, wm.meshes[meshPos])
		}
	}

	if len(toAdd) > 0 {
		wm.addMeshChannel <- toAdd
	}
}

func (wm *WorldMesher) Container() *core.Node {
	return wm.container
}

func (wm *WorldMesher) AddChunk(chunk *world.Chunk) {
	wm.chunks[*chunk.Position()] = chunk
}

func (wm *WorldMesher) RemoveChunk(chunk *world.Chunk) {
	delete(wm.chunks, *chunk.Position())
}

func (wm *WorldMesher) getMissingMeshesPos(pos math32.Vector3) []math32.Vector3 {
	missing := make([]math32.Vector3, 0, 100)
	for x := pos.X - wm.renderDistance; x <= pos.X+wm.renderDistance; x += world.ChunkWidth {
		for z := pos.Z - wm.renderDistance; z <= pos.Z+wm.renderDistance; z += world.ChunkDepth {
			for y := math32.Max(0, pos.Y-wm.renderDistance); y < math32.Min(pos.Y+wm.renderDistance, world.ChunkHeight); y += ChunkletSize {
				newPos := math32.Vector3{X: x, Y: y, Z: z}
				if newPos.DistanceToSquared(&pos) > wm.renderDistance*wm.renderDistance {
					continue
				}
				if _, ok := wm.meshes[newPos]; ok {
					continue
				}
				missing = append(missing, newPos)
			}
		}
	}

	sort.Slice(missing, func(i, j int) bool {
		return missing[i].DistanceTo(&pos) < missing[j].DistanceTo(&pos)
	})

	return missing
}

func (wm *WorldMesher) getWorldPosition(pos math32.Vector3) math32.Vector3 {
	return math32.Vector3{
		X: math32.Floor(pos.X/world.ChunkWidth) * world.ChunkWidth,
		Y: math32.Floor(pos.Y/ChunkletSize) * ChunkletSize,
		Z: math32.Floor(pos.Z/world.ChunkDepth) * world.ChunkDepth,
	}
}

func (wm *WorldMesher) clearTooFarMeshes(pos math32.Vector3) {
	toRemove := make([]*Chunklet, 0, 100)
	for meshPos, mesh := range wm.meshes {
		// todo clear farest first until meshcap is reached
		if math32.Abs(meshPos.X-pos.X) > wm.renderDistance ||
			math32.Abs(meshPos.Y-pos.Y) > wm.renderDistance ||
			math32.Abs(meshPos.Z-pos.Z) > wm.renderDistance {
			if mesh != nil {
				toRemove = append(toRemove, mesh)
			}
			delete(wm.meshes, meshPos)
		}
	}
	wm.removeMeshChannel <- toRemove
}

func (wm *WorldMesher) Dispose() {
	close(wm.positionChannel)
	close(wm.addMeshChannel)
	close(wm.removeMeshChannel)

	for _, mesh := range wm.meshes {
		if mesh != nil {
			mesh.GetGeometry().Dispose()
		}
	}
}
