package block

import (
	"github.com/g3n/engine/core"
)

type Block struct {
	Id   string
	Type string
	*core.Node
}

func (b *Block) Clone() *Block {
	return &Block{
		Id:   b.Id,
		Type: b.Type,
		Node: b.Node.Clone().(*core.Node),
	}
}
