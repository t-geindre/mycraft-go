package block

import (
	"fmt"
)

type Repository struct {
	blocks map[uint16]*Block
}

var RepositoryInstance *Repository

func (r *Repository) Get(id uint16) *Block {
	if block, ok := r.blocks[id]; ok {
		return block
	}
	panic(fmt.Errorf(`unknown block "%s"`, id))
}

func newRepository() *Repository {
	r := new(Repository)
	r.blocks = blockReference()
	for id, block := range r.blocks {
		block.Id = id
	}
	return r
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = newRepository()
	}

	return RepositoryInstance
}
