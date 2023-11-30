package biome

type Biome interface {
	Match(level float32) bool
	SetGround(level float32)
	GetBlockAt(x, y, z float32) uint16
}
