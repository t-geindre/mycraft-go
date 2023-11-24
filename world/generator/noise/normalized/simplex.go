package normalized

import "github.com/ojrac/opensimplex-go"

type Simplex struct {
	opensimplex.Noise32
}

func NewSimplexNoise(seed int64) *Simplex {
	s := new(Simplex)
	s.Noise32 = opensimplex.NewNormalized32(seed)

	return s
}
