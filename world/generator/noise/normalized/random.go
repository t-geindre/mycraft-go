package normalized

import "math/rand"

type Random struct {
}

func NewRandomNoise() *Random {
	return new(Random)
}

func (r Random) Eval2(_, _ float32) float32 {
	return rand.Float32()
}

func (r Random) Eval3(_, _, _ float32) float32 {
	return rand.Float32()
}

func (r Random) Eval4(_, _, _, _ float32) float32 {
	return rand.Float32()
}
