package world

import (
	"github.com/g3n/engine/math32"
	"testing"
)

func TestChunk_IsLayerEmpty(t *testing.T) {
	allEmpty := make(map[int]bool, ChunkHeight)
	for i := 0; i < ChunkHeight; i++ {
		allEmpty[i] = true
	}

	tests := []struct {
		name  string
		calls []int
		want  map[int]bool
	}{
		{"All layers empty", []int{}, map[int]bool{}},
		{"Some layers not empty", []int{8, 9, 150}, map[int]bool{8: false, 9: false, 150: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChunk(math32.Vector2{})
			for _, call := range tt.calls {
				c.SetBlockAt(0, call, 0, 1)
			}
			for layer, want := range allEmpty {
				if _, ok := tt.want[layer]; ok {
					want = tt.want[layer]
				}
				if got := c.IsLayerEmpty(layer); got != want {
					t.Errorf("IsLayerEmpty() = %v, want %v", got, want)
				}
			}
		})
	}
}
