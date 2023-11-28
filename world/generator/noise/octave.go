package noise

type Octave struct {
	noise       Noise
	octaves     int8
	lacunarity  float32
	persistence float32
}

func NewOctave(noise Noise, lacunarity, persistence float32, octaves int8) *Octave {
	o := new(Octave)
	o.noise = noise
	o.lacunarity = lacunarity
	o.persistence = persistence
	o.octaves = octaves

	return o
}

func (o *Octave) Eval2(x, y float32) float32 {
	val, normalizer := float32(0), float32(0)

	frequency := o.lacunarity
	amplitude := o.persistence

	for i := o.octaves; i > 0; i-- {
		val += o.noise.Eval2(x*frequency, y*frequency) * amplitude

		normalizer += amplitude
		frequency *= o.lacunarity
		amplitude *= o.persistence
	}

	return val / normalizer
}
