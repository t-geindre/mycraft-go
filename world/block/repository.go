package block

import (
	"fmt"
)

type Repository struct {
	blocks map[uint8]*Block
}

var RepositoryInstance *Repository

func (r *Repository) Get(id uint8) *Block {
	if block, ok := r.blocks[id]; ok {
		return block
	}
	panic(fmt.Errorf(`unknown block "%d"`, id))
}

func newRepository() *Repository {
	r := new(Repository)
	r.blocks = getReference()
	for id, block := range r.blocks {
		block.id = id
	}
	return r
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = newRepository()
	}

	return RepositoryInstance
}
