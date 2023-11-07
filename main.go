package main

import (
	"sort"
	"strings"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
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

	glWindow := a.IWindow.(*window.GlfwWindow)
	glWindow.SetTitle("MyCraft - version -2.1.125.5-rev4-alpa0")

	//	a.Gls().Enable(gls.MULTISAMPLE)

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Load blocks
	blocks := LoadBlocks()

	// Blocks controls
	var currentBlock string
	setCurrentBlock := func(newBlock string) {
		scene.Remove(blocks[currentBlock])
		currentBlock = newBlock
		scene.Add(blocks[currentBlock])
	}

	blockNames := make([]string, 0, len(blocks))
	for name := range blocks {
		blockNames = append(blockNames, name)
	}
	sort.Strings(blockNames)
	currentBlock = blockNames[0]

	var btnSpacing float32 = 0.0
	for _, blockName := range blockNames {
		blockLabel := strings.Replace(blockName, "_", " ", -1)
		blockLabel = strings.ToUpper(blockLabel)

		btn := gui.NewButton(blockLabel)
		btn.SetPosition(0, btnSpacing)
		btn.SetSize(150, 30)
		btnSpacing += 30

		btn.Subscribe(gui.OnClick, func(blockName string) func(name string, ev interface{}) {
			return func(name string, ev interface{}) {
				setCurrentBlock(blockName)
			}
		}(blockName))

		scene.Add(btn)
	}

	setCurrentBlock(currentBlock)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, -1, 2)
	scene.Add(pointLight)

	// Set background color to gray
	a.Gls().ClearColor(.5, .5, .8, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {

		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
