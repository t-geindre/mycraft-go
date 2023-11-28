package normalized

type Flat struct {
	val float32
}

func NewFlatNoise(val float32) *Flat {
	if val > 1 || val < 0 {
		panic("Flat noise value must be between 0 and 1 (normalized)")
	}

	flat := new(Flat)
	flat.val = val

	return flat
}

func (f Flat) Eval2(x, y float32) float32 {
	return f.val
}
