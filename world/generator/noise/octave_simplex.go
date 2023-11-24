package noise

import (
	"github.com/g3n/engine/math32"
	"mycraft/world/generator/noise/normalized"
)

type OctaveSimplex struct {
	noise   normalized.Simplex
	octaves []parametricSimplexOctave
	ground  float32
}

type parametricSimplexOctave struct {
	period    float32
	amplitude float32
}

func NewOctaveSimplexNoise(seed int64) *OctaveSimplex {
	os := new(OctaveSimplex)
	os.noise = *normalized.NewSimplexNoise(seed)

	// Default configuration
	os.SetGroundLevel(50)
	os.ClearOctaves()
	os.AddOctave(5, 2)
	os.AddOctave(50, 50)
	os.AddOctave(5000, -100)

	return os
}

func (os *OctaveSimplex) AddOctave(period, amplitude float32) {
	os.octaves = append(os.octaves, parametricSimplexOctave{period, amplitude})
}

func (os *OctaveSimplex) ClearOctaves() {
	os.octaves = make([]parametricSimplexOctave, 0)
}

func (os *OctaveSimplex) SetGroundLevel(ground float32) {
	os.ground = ground
}

func (os *OctaveSimplex) Eval2(x, y float32) float32 {
	value := os.ground

	for _, octave := range os.octaves {
		value += os.noise.Eval2(x/octave.period, y/octave.period) * octave.amplitude
	}

	return math32.Floor(value)
}

func (os *OctaveSimplex) Eval3(x, y, z float32) float32 {
	panic("Parametric simplex noise not implemented for 3D noise")
}

func (os *OctaveSimplex) Eval4(x, y, z, w float32) float32 {
	panic("Parametric simplex noise not implemented for 4D noise")
}
