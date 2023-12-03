package biome

type biomeMatch struct {
	biome              Biome
	rangeFrom, rangeTo float32
}

type Selector struct {
	biomes             []*biomeMatch
	rangeFrom, rangeTo float32
}

func NewSelector() *Selector {
	b := new(Selector)
	b.biomes = make([]*biomeMatch, 0)

	return b
}

func (s *Selector) Add(biome Biome, rangeFrom, rangeTo float32) {
	s.biomes = append(s.biomes, &biomeMatch{biome, rangeFrom, rangeTo})
}

func (s *Selector) Match(val float32) Biome {
	for _, b := range s.biomes {
		if val >= b.rangeFrom && val <= b.rangeTo {
			return b.biome
		}
	}

	return nil
}
