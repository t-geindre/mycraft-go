package main

/*

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
	"time"
)

func main() {
	// App
	a := app.App()
	scene := core.NewNode()

	// Cam
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

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

	// Set background color
	a.Gls().ClearColor(0, 0, 0, 1.0)

	// Light
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1))

	// Flower
	sprite, _ := texture.NewTexture2DFromImage("assets/block/dandelion.png")
	sprite.SetMagFilter(gls.NEAREST)

	mat := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
	mat.AddTexture(sprite)
	mat.SetTransparent(true)

	plan := geometry.NewPlane(1, 1)
	for i := float32(0); i < 4; i++ {
		mesh := graphic.NewMesh(plan, mat)
		mesh.SetRotation(0, math32.DegToRad(90*i), 0)
		mesh.SetPosition(0, 0, 0)
		scene.Add(mesh)
	}

	// Run app
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		fmt.Println("render")
		renderer.Render(scene, cam)
	})
}
*/
