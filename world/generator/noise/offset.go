package noise

type Offset struct {
	noise Noise
	x, y  float32
}

func NewOffset(noise Noise, x, y float32) *Offset {
	o := new(Offset)
	o.noise = noise
	o.x = x
	o.y = y

	return o
}

func (o *Offset) Eval2(x, y float32) float32 {
	return o.noise.Eval2(x+o.x, y+o.y)
}
