package noise

type Scale struct {
	noise Noise
	scale float32
}

func NewScale(noise Noise, scale float32) *Scale {
	s := new(Scale)
	s.noise = noise
	s.scale = scale

	return s
}

func (s *Scale) Eval2(x, y float32) float32 {
	return s.noise.Eval2(x/s.scale, y/s.scale)
}
