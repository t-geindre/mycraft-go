package world

type Generator interface {
	Populate(chunk *Chunk, addChunkletChan chan []*Chunklet)
}
