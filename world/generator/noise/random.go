package noise

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/generator/noise/normalized"
)

type Random struct {
	amplitude float32
	noise     *normalized.Random
	ground    float32
}

func NewRandomNoise(amplitude float32, ground float32) *Random {
	r := new(Random)
	r.amplitude = amplitude
	r.ground = ground
	r.noise = normalized.NewRandomNoise()
	return r
}

func (r Random) Eval2(x, y float32) float32 {
	return math32.Floor(r.noise.Eval2(x, y)*r.amplitude) + r.ground
}

func (r Random) Eval3(x, y, z float32) float32 {
	return math32.Floor(r.noise.Eval3(x, y, z)*r.amplitude) + r.ground
}

func (r Random) Eval4(x, y, z, w float32) float32 {
	return math32.Floor(r.noise.Eval4(x, y, z, w)*r.amplitude) + r.ground
}
