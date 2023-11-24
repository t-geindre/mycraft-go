package normalized

import "github.com/g3n/engine/math32"

type Sin struct {
}

func NewSinNoise() *Sin {
	return new(Sin)
}

func (s Sin) Eval2(x, y float32) float32 {
	return (math32.Sin(x) + math32.Sin(y)) / 2
}

func (s Sin) Eval3(x, y, z float32) float32 {
	return (math32.Sin(x) + math32.Sin(y) + math32.Sin(z)) / 3
}

func (s Sin) Eval4(x, y, z, w float32) float32 {
	return (math32.Sin(x) + math32.Sin(y) + math32.Sin(z) + math32.Sin(w)) / 4
}
