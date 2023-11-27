package noise

type Noise interface {
	Eval2(x, y float32) float32
}
