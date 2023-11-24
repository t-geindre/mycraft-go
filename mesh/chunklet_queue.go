package mesh

import (
	"github.com/g3n/engine/math32"
	"mycraft/world"
)

type ChunkletQueue struct {
	chanIn          chan []ChunkletQuery
	chanOut         chan []*Chunklet
	awaitingQueries []ChunkletQuery
	packSize        int
	routines        int
	handler         func(*Chunklet)
}

type ChunkletQuery struct {
	Pos    math32.Vector3
	Center *world.Chunk
	East   *world.Chunk
	West   *world.Chunk
	North  *world.Chunk
	South  *world.Chunk
}

func NewChunkletQueue(packSize int, routines int, handler func(*Chunklet)) *ChunkletQueue {
	cq := new(ChunkletQueue)

	cq.packSize = packSize
	cq.routines = routines
	cq.handler = handler

	cq.chanIn = make(chan []ChunkletQuery, 1000)
	cq.chanOut = make(chan []*Chunklet, 1000)

	cq.awaitingQueries = make([]ChunkletQuery, 0, cq.packSize)

	go cq.Run()

	return cq
}

func (cq *ChunkletQueue) Push(q ChunkletQuery) {
	cq.Pull()
	cq.awaitingQueries = append(cq.awaitingQueries, q)

	if len(cq.awaitingQueries) >= cq.packSize {
		cq.chanIn <- cq.awaitingQueries
		cq.awaitingQueries = make([]ChunkletQuery, 0, cq.packSize)
		cq.Pull()
	}
}

func (cq *ChunkletQueue) Run() {
	for {
		queries := <-cq.chanIn
		meshes := make([]*Chunklet, len(queries))
		for _, q := range queries {
			meshes = append(meshes, NewChunklet(q.Center, q.East, q.West, q.North, q.South, q.Pos))
		}
		cq.chanOut <- meshes
	}
}

func (cq *ChunkletQueue) Pull() {
	for {
		select {
		case meshes := <-cq.chanOut:
			for _, mesh := range meshes {
				cq.handler(mesh)
			}
		default:
			return
		}
	}
}
