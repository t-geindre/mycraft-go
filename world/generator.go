package world

type Generator interface {
	GetBlockAt(x, y, z float32) uint16
	Reset()
}
