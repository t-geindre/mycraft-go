package main

import (
	"fmt"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/util"
	"mycraft/block"
	"mycraft/camera"
	"mycraft/material"
	"mycraft/world"
	"sort"
	"strings"
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

	demoWorld := world.NewDemoWorld(scene, 20, &blocksRepository)

	// Current block switcher
	var currentBlock *graphic.Mesh
	setCurrentBlock := func(id string) {
		scene.Remove(currentBlock)
		currentBlock = blocksRepository.Get(id).CreateMesh()
		scene.Add(currentBlock)
		currentBlock.SetPosition(0, 0, 0)
	}

	// Create sorted blocks list
	blockIds := make([]string, 0, len(blocksRepository.Blocks))
	for name := range blocksRepository.Blocks {
		blockIds = append(blockIds, name)
	}
	sort.Strings(blockIds)

	// First sorted block set as default
	setCurrentBlock(blockIds[0])

	// GUI
	_, height := a.GetSize()
	btnSpacing := float32(height / len(blockIds))

	for idx, blockName := range blockIds {
		blockLabel := strings.Replace(blockName, "_", " ", -1)
		blockLabel = strings.ToUpper(blockLabel)

		btn := gui.NewButton(blockLabel)
		btn.SetPosition(0, btnSpacing*float32(idx))
		btn.SetSize(150, btnSpacing)

		btn.Subscribe(
			gui.OnClick,
			func(blockName string) func(name string, ev interface{}) {
				return func(_ string, _ interface{}) {
					setCurrentBlock(blockName)
					// remove button focus to get window events
					gui.Manager().SetKeyFocus(nil)
				}
			}(blockName),
		)

		scene.Add(btn)
	}

	controlEnabled := false
	gui.Manager().Subscribe(window.OnKeyDown, func(_ string, ev any) {
		kev := ev.(*window.KeyEvent)
		if kev.Key != window.KeyEscape {
			return
		}

		if !controlEnabled {
			WASMControl.ReleaseMouse()
			controlEnabled = true
			return
		}
		WASMControl.CaptureMouse(glWindow)
		controlEnabled = false

	})

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

	// === test
	geo := geometry.NewPlane(1, 1)

	for j := float32(0); j < 10; j++ {
		for k := float32(0); k < 10; k++ {
			for i := float32(0); i < 4; i++ {
				mesh := graphic.NewMesh(geo, materialRepository.Get("oxeye_daisy"))
				mesh.SetPosition(2+j, -1, k)
				mesh.SetRotation(0, math32.DegToRad(90*i), 0)
				scene.Add(mesh)
			}
		}
	}

	// Set background color to some blue
	a.Gls().ClearColor(.5, .5, .8, 1.0)

	// Run the application
	rotation := float32(0)
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		framerater.Start()

		fps, _, ok := framerater.FPS(deltaTime)
		if ok {
			label.SetText(fmt.Sprintf("FPS %d", int(fps)))
		}

		WASMControl.Update(deltaTime)
		demoWorld.Update(cam.Position())

		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)

		rotation += float32(deltaTime.Milliseconds()) / 1000
		currentBlock.SetRotation(rotation, rotation/2, 0)

		framerater.Wait()
	})
}
