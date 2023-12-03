package noise

import "github.com/g3n/engine/math32"

type Floor struct {
	noise Noise
}

func NewFloor(noise Noise) *Floor {
	f := new(Floor)
	f.noise = noise
	return f
}

func (f *Floor) Eval2(x, y float32) float32 {
	return math32.Floor(f.noise.Eval2(x, y))
}
