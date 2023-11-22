package scene

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"mycraft/app"
	"mycraft/camera"
	"mycraft/mesh"
	"mycraft/world"
	"mycraft/world/generator/infinite"
	"time"
)

const renderingDistance = 10 // chunks

type Game struct {
	container      *core.Node
	app            *app.App
	camControl     *camera.WASMControl
	cursorCaptured bool
	world          *world.World
	worldMesher    *mesh.WorldMesher
}

func NewGameScene() *Game {
	return new(Game)
}

func (g *Game) Setup(container *core.Node, app *app.App) {
	g.container = container
	g.app = app

	// Set WASM camera control
	g.app.Cam.SetPosition(0, 100, 0)
	g.camControl = camera.NewWASMControl(g.app.Cam)
	g.camControl.CaptureMouse(g.app.GlsWindow)

	// Toggle mouse capture on escape key
	g.cursorCaptured = false
	g.app.Engine.SubscribeID(window.OnKeyDown, &g, g.onKeyDown)

	// Create and add lights
	g.container.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1))

	// Set background color to some blue todo : add skybox
	g.app.Engine.Gls().ClearColor(.5, .5, .8, 1.0)

	// Create world
	// Rendering distance is increased by 1 to avoid chunks not being rendered
	g.world = world.NewWorld(renderingDistance+1, infinite.NewNoiseGenerator(1)) //infinite.NewSinGenerator(10, 20))

	// Create world mesher
	g.worldMesher = mesh.NewWorldMesher(renderingDistance * mesh.ChunkletSize)
	g.container.Add(g.worldMesher.Container())
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
	g.world.UpdateFromVec3(g.app.Cam.Position())

	select {
	case chunks := <-g.world.AddChunkChannel():
		for _, chunk := range chunks {
			g.worldMesher.AddChunk(chunk)
		}
	case chunks := <-g.world.RemoveChunkChannel():
		for _, chunk := range chunks {
			g.worldMesher.RemoveChunk(chunk)
		}
	default:
	}

	g.camControl.Update(deltaTime)
	g.worldMesher.Update(g.app.Cam.Position())
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)
}
