package noise

type Flat struct {
	val float32
}

func NewFlatNoise(val float32) *Flat {
	flat := new(Flat)
	flat.val = val

	return flat
}

func (f Flat) Eval2(x, y float32) float32 {
	return f.val
}

func (f Flat) Eval3(x, y, z float32) float32 {
	return f.val
}

func (f Flat) Eval4(x, y, z, w float32) float32 {
	return f.val
}
