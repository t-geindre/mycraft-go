package noise

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/generator/noise/normalized"
)

type Mountain struct {
	noise Noise
}

func NewMountainNoise(seed int64) *Mountain {
	m := new(Mountain)
	m.noise = normalized.NewSimplexNoise(seed)

	return m
}

func (m Mountain) Eval2(x, y float32) float32 {
	val := float32(0)
	pv := m.noise.Eval2((x+100)/500, (y+100)/500) * 10
	from := float32(2)
	to := float32(4)
	if pv >= from && pv <= to {
		fact := 1 - math32.Abs(pv-(from+to)/2)/((to-from)/2)
		val += m.noise.Eval2(x/40, y/40) * 50 * fact
		val += m.noise.Eval2(x/10, y/10) * 5 * fact
	}

	return math32.Floor(val)
}
