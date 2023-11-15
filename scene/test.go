package scene

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"mycraft/app"
	"mycraft/block"
	mesh2 "mycraft/mesh"
	"mycraft/world"
	"time"
)

type Test struct {
}

func (t *Test) Setup(container *core.Node, app *app.App) {

	chunk := world.NewChunklet(math32.NewVector3(0, 0, 0))

	grassBlock := block.GetRepository().Get("green_grass")
	dirtBlock := block.GetRepository().Get("dirt")

	for x := 0; x < chunk.Size; x++ {
		for y := 0; y < chunk.Size; y++ {
			for z := 0; z < chunk.Size; z++ {
				if y == 2 || (y == 3 && x%2 == 0 && z%3 == 0) {
					chunk.Blocks[x][y][z] = grassBlock
					chunk.Blocks[x][y-1][z] = dirtBlock
					continue
				}
			}
		}
	}

	container.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.2))

	mesher := mesh2.NewChunkletMesher(chunk)
	mesher.ComputeQuads()
	mesh := mesher.GetMesh()
	mesh.SetPositionVec(math32.NewVector3(0, 0, -5))
	container.Add(mesh)

	orb := camera.NewOrbitControl(app.Cam)
	orb.SetTarget(*math32.NewVector3(float32(chunk.Size/2), 2, float32(chunk.Size/2)))
}

func (t *Test) Update(deltaTime time.Duration) {
}

func (t *Test) Dispose() {
}

func NewTestScene() *Test {
	return new(Test)
}
