package noise

type Noise interface {
	Eval2(x, y float32) float32
	Eval3(x, y, z float32) float32
	Eval4(x, y, z, w float32) float32
}
