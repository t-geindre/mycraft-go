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
	"mycraft/world/generator"
	"time"
)

const renderingDistance = 30 // chunks

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
	g.app.Cam.SetPosition(0, 60, 0)
	g.camControl = camera.NewWASMControl(g.app.Cam)
	g.camControl.CaptureMouse(g.app.GlsWindow)

	// Toggle mouse capture on escape key
	g.cursorCaptured = false
	g.app.Engine.SubscribeID(window.OnKeyDown, &g, g.onKeyDown)

	// Create and add lights
	// todo check how we can maybe set how materials react to light SetUseLights()
	g.container.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, .8))
	dl := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, .8)
	dl.SetDirectionVec(&math32.Vector3{X: 0, Y: -1, Z: 0})
	dl.SetPositionVec(&math32.Vector3{X: 0, Y: 300, Z: 0})
	g.container.Add(dl)

	// Add skybox
	g.container.Add(mesh.NewSkybox())

	// Create world
	g.world = world.NewWorld(renderingDistance, generator.NewBiomeGenerator(0))

	// Create world mesher
	g.worldMesher = mesh.NewWorldMesher(renderingDistance*mesh.ChunkletSize, g.world)
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
	g.camControl.Update(deltaTime)
	g.world.UpdateFromVec3(g.app.Cam.Position())
	g.worldMesher.Update(g.app.Cam.Position())
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)
	g.world.Dispose()
	g.worldMesher.Dispose()
}
