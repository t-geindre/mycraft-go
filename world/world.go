package world

import (
	"github.com/g3n/engine/math32"
)

type World struct {
	chunks          map[math32.Vector2]*Chunk
	rDist           float32 // chunks
	generator       Generator
	posChan         chan math32.Vector2
	lastPost        math32.Vector2
	initialized     bool
	addChunkChan    chan *Chunk
	removeChunkChan chan *Chunk
}

func NewWorld(rd float32, generator Generator) *World {
	w := new(World)
	w.chunks = make(map[math32.Vector2]*Chunk)
	w.rDist = rd
	w.generator = generator
	w.posChan = make(chan math32.Vector2, 1)
	w.addChunkChan = make(chan *Chunk, 10)
	w.removeChunkChan = make(chan *Chunk, 10)
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
		w.clearTooFarChunks(pos)
	}
}

func (w *World) Update(pos math32.Vector2) {
	if len(w.posChan) > 0 {
		<-w.posChan
	}
	w.posChan <- pos
}

func (w *World) UpdateFromVec3(pos math32.Vector3) {
	w.Update(w.GetWorldCoordinates(pos))
}

func (*World) GetWorldCoordinates(pos math32.Vector3) math32.Vector2 {
	return math32.Vector2{
		X: float32(int(pos.X) / ChunkWith),
		Y: float32(int(pos.Z) / ChunkDepth),
	}
}

func (w *World) AddChunkChannel() chan *Chunk {
	return w.addChunkChan
}

func (w *World) RemoveChunkChannel() chan *Chunk {
	return w.removeChunkChan
}

func (w *World) addMissingChunks(pos math32.Vector2) {
	for x := pos.X - w.rDist; x <= pos.X+w.rDist; x++ {
		for y := pos.Y - w.rDist; y <= pos.Y+w.rDist; y++ {
			chunkPos := math32.Vector2{X: x, Y: y}
			if _, ok := w.chunks[chunkPos]; !ok {
				chunkWorldPos := math32.Vector2{
					X: float32(int(x) * ChunkWith),
					Y: float32(int(y) * ChunkDepth),
				}
				w.chunks[chunkPos] = NewChunk(chunkWorldPos)
				w.generator.Populate(w.chunks[chunkPos])
				w.addChunkChan <- w.chunks[chunkPos]
			}
		}
	}
}

func (w *World) clearTooFarChunks(pos math32.Vector2) {
	for chunkPos, _ := range w.chunks {
		if math32.Abs(chunkPos.X-pos.X) > w.rDist || math32.Abs(pos.Y-chunkPos.Y) > w.rDist {
			w.removeChunkChan <- w.chunks[chunkPos]
			delete(w.chunks, chunkPos)
		}
	}
}
