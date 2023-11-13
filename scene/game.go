package scene

import (
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

	addNodeChannel    chan []*block.Block
	removeNodeChannel chan []*block.Block
	positionChannel   chan math32.Vector3
	camControl        *camera.WASMControl
	cursorCaptured    bool
	blocksToSkip      map[*block.Block]bool
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
	g.addNodeChannel = make(chan []*block.Block, chanPackSize*2)
	g.removeNodeChannel = make(chan []*block.Block, chanPackSize*2)
	g.positionChannel = make(chan math32.Vector3, 1)

	// Block skip setup
	g.blocksToSkip = make(map[*block.Block]bool)

	// Start world routine
	demoWorld := world.NewDemoWorld(renderingDistance, g.addNodeChannel, g.removeNodeChannel, g.positionChannel, chanPackSize)
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
	g.camControl.Update(deltaTime)
	g.updateWorld()
}

func (g *Game) Dispose() {
	g.app.Engine.UnsubscribeID(window.OnKeyDown, &g)

	close(g.addNodeChannel)
	close(g.removeNodeChannel)
	close(g.positionChannel)
}

func (g *Game) updateWorld() {
	select {
	case nodes := <-g.addNodeChannel:
		for _, bl := range nodes {
			if g.blocksToSkip[bl] {
				delete(g.blocksToSkip, bl)
				continue
			}
			g.container.Add(bl.Node)
		}
	case nodes := <-g.removeNodeChannel:
		for _, bl := range nodes {
			if !g.container.Remove(bl.Node) {
				// Block not added yet, store it to skip it later
				g.blocksToSkip[bl] = true
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
