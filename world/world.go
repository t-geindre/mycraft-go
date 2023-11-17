package world

import (
	"github.com/g3n/engine/math32"
)

type World struct {
	chunks             map[math32.Vector2]*Chunk
	rDist              float32 // chunks
	generator          Generator
	posChan            chan math32.Vector2
	addChunkletChan    chan []*Chunklet
	removeChunkletChan chan []*Chunklet
	lastPost           math32.Vector2
	initialized        bool
}

func NewWorld(rd float32, generator Generator) *World {
	w := new(World)
	w.chunks = make(map[math32.Vector2]*Chunk)
	w.rDist = rd
	w.generator = generator
	w.posChan = make(chan math32.Vector2, 1)
	w.addChunkletChan = make(chan []*Chunklet, ChunkSize)
	w.removeChunkletChan = make(chan []*Chunklet, ChunkSize)
	go w.Run()
	return w
}

func (w *World) Run() {
	for {
		pos := <-w.posChan
		if pos.Equals(&w.lastPost) && w.initialized {
			continue
		}

		w.lastPost = pos
		w.initialized = true

		w.addMissingChunks(pos)
	}
}

func (w *World) Update(pos math32.Vector2) []*Chunklet {
	if len(w.posChan) > 0 {
		<-w.posChan
	}
	w.posChan <- pos

	w.clearTooFarChunks(pos)

	return w.getChunkletToAdd()
}

func (w *World) UpdateFromVec3(pos math32.Vector3) []*Chunklet {
	return w.Update(w.GetWorldCoordinates(pos))
}

func (*World) GetWorldCoordinates(pos math32.Vector3) math32.Vector2 {
	return math32.Vector2{
		X: float32(int(pos.X) / ChunkSize / ChunkletSize),
		Y: float32(int(pos.Z) / ChunkSize / ChunkletSize),
	}
}

func (w *World) addMissingChunks(pos math32.Vector2) {
	for x := pos.X - w.rDist; x <= pos.X+w.rDist; x++ {
		for y := pos.Y - w.rDist; y <= pos.Y+w.rDist; y++ {
			chunkPos := math32.Vector2{X: x, Y: y}
			if _, ok := w.chunks[chunkPos]; !ok {
				chunkWorldPos := math32.Vector2{
					X: float32(int(x) * ChunkSize * ChunkletSize),
					Y: float32(int(y) * ChunkSize * ChunkletSize),
				}
				w.chunks[chunkPos] = NewChunk(chunkWorldPos)
				w.generator.Populate(w.chunks[chunkPos], w.addChunkletChan)
			}
		}
	}
}

func (w *World) clearTooFarChunks(pos math32.Vector2) {
	for chunkPos, _ := range w.chunks {
		if math32.Abs(chunkPos.X-pos.X) > w.rDist || math32.Abs(pos.Y-chunkPos.Y) > w.rDist {
			for _, chunklet := range w.chunks[chunkPos].chunklets {
				chunklet.Dispose()
			}
			delete(w.chunks, chunkPos)
		}
	}
}

func (w *World) getChunkletToAdd() []*Chunklet {
	select {
	case chunklet := <-w.addChunkletChan:
		return chunklet
	default:
		return nil
	}
}
