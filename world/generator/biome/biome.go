package biome

type Biome interface {
	// GetBlockAt
	// strength 0-1 biome strength on terrain
	GetBlockAt(x, y, z, strength float32) uint16
	Reset()
}
