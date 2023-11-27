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
	pv := m.noise.Eval2(x/500, y/500) * 10
	from := float32(1)
	to := float32(4)
	if pv >= from && pv <= to {
		fact := 1 - math32.Abs(pv-(from+to)/2)/((to-from)/2)
		return m.noise.Eval2(x/50, y/50)*50*fact + m.noise.Eval2(x/10, y/10)*6*fact
	}

	return 0
}
