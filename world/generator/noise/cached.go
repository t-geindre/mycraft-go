package noise

type Cached struct {
	cache    map[float32]map[float32]*float32
	noise    Noise
	capacity int
}

func NewCachedNoise(noise Noise, capacity int) *Cached {
	e := new(Cached)

	e.noise = noise
	e.capacity = capacity

	e.Clear()

	return e
}

func (e *Cached) Eval2(x, z float32) float32 {
	if e.cache[x][z] == nil {
		if e.cache[x] == nil {
			e.cache[x] = make(map[float32]*float32, e.capacity)
		}
		value := e.noise.Eval2(x, z)
		e.cache[x][z] = &value
	}

	return *e.cache[x][z]
}

func (e *Cached) Eval3(x, y, z float32) float32 {
	panic("noise cache not implemented for 3D noise")
}

func (e *Cached) Eval4(x, y, z, w float32) float32 {
	panic("noise cache not implemented for 4D noise")
}

func (e *Cached) Clear() {
	e.cache = make(map[float32]map[float32]*float32, e.capacity)
}
