package noise

import "github.com/g3n/engine/math32"

type Clamp struct {
	noise    Noise
	min, max float32
}

func NewClamp(noise Noise, min, max float32) *Clamp {
	c := new(Clamp)
	c.noise = noise
	c.min = min
	c.max = max
	return c
}

func (c *Clamp) Eval2(x, y float32) float32 {
	return math32.Clamp(c.noise.Eval2(x, y), c.min, c.max)
}
