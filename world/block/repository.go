package block

import (
	"fmt"
)

type Repository struct {
	blocks map[uint16]func() *Block
}

var RepositoryInstance *Repository

func (r *Repository) Get(id uint16) *Block {
	if block, ok := r.blocks[id]; ok {
		return block()
	}
	panic(fmt.Errorf(`unknown block "%s"`, id))
}

func newRepository() *Repository {
	r := new(Repository)
	r.blocks = make(map[uint16]func() *Block)

	for id, def := range blockReference() {
		r.blocks[id] = func(def blockDef) func() *Block {
			return func() *Block {
				b := Block{}
				b.Materials = BlockMaterials{
					Top:    def.Materials.Top,
					Bottom: def.Materials.Bottom,
					North:  def.Materials.North,
					South:  def.Materials.South,
					East:   def.Materials.East,
					West:   def.Materials.West,
				}

				return &b
			}
		}(def)
	}

	return r
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = newRepository()
	}

	return RepositoryInstance
}
