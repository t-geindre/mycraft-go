package biome

type Biome interface {
	// Get biome selection strength 0-1
	Select(rainfall, temperature float32) float32
	GetBlockAt(x, y, z, strength float32) uint16
	Reset()
}
