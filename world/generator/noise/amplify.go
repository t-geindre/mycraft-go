package noise

type Amplify struct {
	noise     Noise
	amplitude float32
}

func NewAmplify(noise Noise, amplitude float32) *Amplify {
	a := new(Amplify)
	a.noise = noise
	a.amplitude = amplitude

	return a
}

func (a *Amplify) Eval2(x, y float32) float32 {
	return a.noise.Eval2(x, y) * a.amplitude
}
