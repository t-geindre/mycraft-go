package app

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/window"
	"log"
	"time"
)

type Scene interface {
	Setup(container *core.Node, app *App)
	Update(deltaTime time.Duration)
	Dispose()
}

type App struct {
	Engine      *app.Application
	Cam         *camera.Camera
	GlsWindow   *window.GlfwWindow
	RootNode    *core.Node
	Scenes      map[Scene]Scene
	ScenesNodes map[Scene]*core.Node
	Framerater  *util.FrameRater
}

func NewApp(appName string, targetFps uint) *App {
	a := new(App)

	a.Engine = app.App()
	a.RootNode = core.NewNode()
	gui.Manager().Set(a.RootNode)

	a.GlsWindow = a.Engine.IWindow.(*window.GlfwWindow)
	a.GlsWindow.SetTitle(appName)
	a.GlsWindow.SetSize(800, 600)

	a.Cam = camera.New(1)
	a.Cam.SetPosition(0, 0, 0)
	a.RootNode.Add(a.Cam)

	a.Scenes = make(map[Scene]Scene)
	a.ScenesNodes = make(map[Scene]*core.Node)

	a.Engine.Subscribe(window.OnKeyDown, a.OnKeyDown)
	a.Engine.Subscribe(window.OnWindowSize, a.OnResize)

	a.Framerater = util.NewFrameRater(targetFps)

	a.OnResize("", nil)

	return a
}

func (a *App) ToggleFullscreen() {
	a.GlsWindow.SetFullscreen(!a.GlsWindow.Fullscreen())
}

func (a *App) OnKeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)
	if kev.Key == window.KeyEnter && kev.Mods == window.ModAlt {
		a.ToggleFullscreen()
		return
	}
}

func (a *App) OnResize(_ string, _ interface{}) {
	// Get framebuffer size and update viewport accordingly
	wWidth, wHeight := a.Engine.GetSize()
	a.Engine.Gls().Viewport(0, 0, int32(wWidth), int32(wHeight))
	// Update the camera's aspect ratio
	a.Cam.SetAspect(float32(wWidth) / float32(wHeight))
}

func (a *App) AddScene(scene Scene) {
	node := core.NewNode()
	scene.Setup(node, a)
	a.RootNode.Add(node)
	a.Scenes[scene] = scene
	a.ScenesNodes[scene] = node
}

func (a *App) RemoveScene(scene Scene) {
	scene.Dispose()
	a.RootNode.Remove(a.ScenesNodes[scene])

	delete(a.Scenes, scene)
	delete(a.ScenesNodes, scene)
}

func (a *App) Run() {
	a.Engine.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Framerater.Start()
		a.Engine.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		for _, scene := range a.Scenes {
			scene.Update(deltaTime)
		}

		if err := renderer.Render(a.RootNode, a.Cam); err != nil {
			log.Fatal(err)
		}
		a.Framerater.Wait()
	})
}
