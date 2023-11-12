package main

import (
	"fmt"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/util"
	"mycraft/block"
	"mycraft/camera"
	"mycraft/material"
	"mycraft/world"
	"time"

	"github.com/g3n/engine/app"
	g3ncamera "github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

func main() {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Window setup
	glWindow := a.IWindow.(*window.GlfwWindow)
	glWindow.SetTitle("MyCraft - version -2.1.125.5-rev4-alpa0")
	glWindow.SetSize(1600, 900)
	//glWindow.SetFullscreen(true)

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := g3ncamera.New(1)
	cam.SetPosition(0, 0, 0)
	scene.Add(cam)

	// Set up orbit control for the camera
	//g3ncamera.NewOrbitControl(cam)
	WASMControl := camera.NewWASMControl(cam)
	WASMControl.CaptureMouse(glWindow)

	// Load materials and blocks
	materialRepository := material.NewFromYamlFile("assets/materials.yaml")
	blocksRepository := block.NewFromYAMLFile("assets/blocks.yaml", materialRepository)

	// Framerate control/display
	framerater := util.NewFrameRater(60)
	label := gui.NewLabel("0")
	scene.Add(label)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		wWidth, wHeight := a.GetSize()
		a.Gls().Viewport(0, 0, int32(wWidth), int32(wHeight))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(wWidth) / float32(wHeight))
		// Update fps display position
		label.SetPosition(float32(wWidth)-60, 10)
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.1))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 2, 2)
	scene.Add(pointLight)

	// Set background color to some blue
	a.Gls().ClearColor(.5, .5, .8, 1.0)

	// World channels
	addMeshChannel := make(chan []*graphic.Mesh, 200)
	defer close(addMeshChannel)

	removeMeshChannel := make(chan []*graphic.Mesh, 200)
	defer close(removeMeshChannel)

	positionChannel := make(chan math32.Vector3, 1)
	defer close(positionChannel)

	demoWorld := world.NewDemoWorld(
		40,
		&blocksRepository,
		addMeshChannel,
		removeMeshChannel,
		positionChannel,
		200,
	)
	go demoWorld.Run()
	positionChannel <- cam.Position()

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		framerater.Start()

		fps, _, ok := framerater.FPS(deltaTime)
		if ok {
			label.SetText(fmt.Sprintf("FPS %d", int(fps)))
		}

		WASMControl.Update(deltaTime)

		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)

		select {
		case meshes := <-addMeshChannel:
			for _, mesh := range meshes {
				scene.Add(mesh)
			}
		case meshes := <-removeMeshChannel:
			for _, mesh := range meshes {
				scene.Remove(mesh)
			}
		default:
			if len(positionChannel) == 0 {
				// If the world go routine is still handling the last position, don't send a new one
				positionChannel <- cam.Position()
			}
		}

		framerater.Wait()
	})
}
