package infinite

type Noise struct {
	seed int64
}

func NewNoiseGenerator(seed int64) *Noise {
	i := new(Noise)
	i.seed = seed
	return i
}
