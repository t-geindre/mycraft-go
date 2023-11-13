package scene

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"mycraft/app"
	"mycraft/block"
	"mycraft/camera"
	"mycraft/world"
	"time"
)

const renderingDistance = 20
const chanPackSize = 100

type Game struct {
	container *core.Node
	app       *app.App

	addMeshChannel    chan []*block.Block
	removeMeshChannel chan []*block.Block
	positionChannel   chan math32.Vector3
	camControl        *camera.WASMControl
	cursorCaptured    bool
}

func NewGameScene() *Game {
	return new(Game)
}

func (g *Game) Setup(container *core.Node, app *app.App) {
	g.container = container
	g.app = app

	// Set WASM camera control
	g.camControl = camera.NewWASMControl(g.app.Cam)
	g.camControl.CaptureMouse(g.app.GlsWindow)

	// Toggle mouse capture on escape key
	g.cursorCaptured = false
	g.app.Engine.SubscribeID(window.OnKeyDown, &g, g.onKeyDown)

	// Create and add lights
	g.container.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.1))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 2, 2)
	g.container.Add(pointLight)

	// Set background color to some blue
	g.app.Engine.Gls().ClearColor(.5, .5, .8, 1.0)

	// World channel setup
	g.addMeshChannel = make(chan []*block.Block, 200)
	g.removeMeshChannel = make(chan []*block.Block, 200)
	g.positionChannel = make(chan math32.Vector3, 1)

	// Start world routine
	demoWorld := world.NewDemoWorld(renderingDistance, g.addMeshChannel, g.removeMeshChannel, g.positionChannel, chanPackSize)
	go demoWorld.Run()
	g.positionChannel <- g.app.Cam.Position()
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
	// Cam control update
	g.camControl.Update(deltaTime)

	// World routine communication
	select {
	case blocks := <-g.addMeshChannel:
		for _, bl := range blocks {
			for _, mesh := range bl.Meshes {
				g.container.Add(mesh)
			}
		}
	case blocks := <-g.removeMeshChannel:
		for _, bl := range blocks {
			for _, mesh := range bl.Meshes {
				if !g.container.Remove(mesh) {
					fmt.Println("MESH NOT FOUND")
				}
			}
		}
	default:
		if len(g.positionChannel) == 0 {
			// If the world routine is still handling the previous position
			// there is no need to send a new one yet
			g.positionChannel <- g.app.Cam.Position()
		}
	}
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)

	close(g.addMeshChannel)
	close(g.removeMeshChannel)
	close(g.positionChannel)
}
