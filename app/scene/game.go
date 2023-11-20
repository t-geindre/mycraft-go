package scene

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"mycraft/app"
	"mycraft/camera"
	"mycraft/mesh"
	"mycraft/world"
	"mycraft/world/generator"
	"time"
)

const renderingDistance = 4 // chunks

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
	g.world = world.NewWorld(renderingDistance, generator.NewInfiniteRandomGenerator())
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

	g.world.UpdateFromVec3(g.app.Cam.Position())
	select {
	case chunklets := <-g.world.GetChunkletToAdd():
		for _, chunklet := range chunklets {
			chunklet.Mesh = mesh.NewChunkletMesh(chunklet)
			g.container.Add(chunklet.Mesh)
		}
	case chunklets := <-g.world.GetChunkletToRemove():
		for _, chunklet := range chunklets {
			if ok := g.container.Remove(chunklet.Mesh); !ok {
				fmt.Println("Failed to remove chunklet")
			}
		}
	default:
	}
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)
}
