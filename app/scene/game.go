package scene

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"mycraft/app"
	"mycraft/camera"
	mesh2 "mycraft/mesh"
	"mycraft/world"
	"mycraft/world/generator"
	"time"
)

const renderingDistance = 5 // chunks

type Game struct {
	container      *core.Node
	app            *app.App
	camControl     *camera.WASMControl
	cursorCaptured bool
	world          *world.World
}

func NewGameScene() *Game {
	return new(Game)
}

func (g *Game) Setup(container *core.Node, app *app.App) {
	g.container = container
	g.app = app

	// Set WASM camera control
	g.app.Cam.SetPosition(0, 10, 0)
	g.camControl = camera.NewWASMControl(g.app.Cam)
	g.camControl.CaptureMouse(g.app.GlsWindow)

	// Toggle mouse capture on escape key
	g.cursorCaptured = false
	g.app.Engine.SubscribeID(window.OnKeyDown, &g, g.onKeyDown)

	// Create and add lights
	g.container.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.1))

	// Set background color to some blue todo : add skybox
	g.app.Engine.Gls().ClearColor(.5, .5, .8, 1.0)

	// Create world
	g.world = world.NewWorld(renderingDistance, generator.NewInfiniteFlatGenerator())
}

func (g *Game) onKeyDown(_ string, ev interface{}) {
	kev := ev.(*window.KeyEvent)
	if kev.Key == window.KeyEscape {
		if g.cursorCaptured {
			g.cursorCaptured = false
			g.camControl.ReleaseMouse()
		} else {
			g.cursorCaptured = true
			g.camControl.CaptureMouse(g.app.GlsWindow)
		}
	}
}

func (g *Game) Update(deltaTime time.Duration) {
	g.camControl.Update(deltaTime)

	chunklets := g.world.UpdateFromVec3(g.app.Cam.Position())
	if chunklets != nil && len(chunklets) > 0 {
		for _, chunklet := range chunklets {
			if !chunklet.Empty {
				mesher := mesh2.NewChunkletMesher(chunklet)
				mesher.ComputeQuads()
				mesh := mesher.GetMesh()
				g.container.Add(mesh)

				chunklet.Subscribe(world.OnDispose, func(mesh *graphic.Mesh) func(_ string, ev interface{}) {
					return func(_ string, ev interface{}) {
						g.container.Remove(mesh)
						mesh.GetGeometry().Dispose()
					}
				}(mesh))
			}
		}
	}
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)
}
