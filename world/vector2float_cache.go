package world

type Vector2FloatCache struct {
	elevation map[float32]map[float32]*float32
	handler   func(x, z float32) float32
}

func NewElevationCache(handler func(x, z float32) float32) *Vector2FloatCache {
	e := new(Vector2FloatCache)
	e.Reset()
	e.handler = handler
	return e
}

func (e *Vector2FloatCache) Reset() {
	e.elevation = make(map[float32]map[float32]*float32, ChunkWidth)
}

func (e *Vector2FloatCache) GetValue(x, z float32) float32 {
	if e.elevation[x][z] == nil {
		if e.elevation[x] == nil {
			e.elevation[x] = make(map[float32]*float32, ChunkDepth)
		}
		el := e.handler(x, z)
		e.elevation[x][z] = &el
	}

	return *e.elevation[x][z]
}
