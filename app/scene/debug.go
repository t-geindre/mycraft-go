package scene

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	engineMaterial "github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"mycraft/app"
	"time"
)

type Debug struct {
	Stats []*DebugStat

	container       *core.Node
	toggleContainer *core.Node
	isActive        bool
	app             *app.App
	panel           *gui.Panel
	panelPadding    float32
	delays          map[*DebugStat]time.Duration
	wireFrame       bool
}

type DebugStat struct {
	label  string
	value  *gui.Label
	update func(deltaTime time.Duration, label *gui.Label)
	delay  time.Duration
}

func NewDebugScene() *Debug {
	d := new(Debug)
	d.isActive = true
	d.panelPadding = 5
	d.Stats = []*DebugStat{
		&DebugStat{label: "FPS", update: d.updateFps},
		&DebugStat{label: "Scenes", update: d.updateScenes, delay: 500 * time.Millisecond},
		&DebugStat{label: "Meshes", update: d.updateMeshes, delay: 1000 * time.Millisecond},
		&DebugStat{label: "Polygons", update: d.updatePolygons, delay: 1000 * time.Millisecond},
		&DebugStat{label: "Cam", update: d.updateCamPosition, delay: 300 * time.Millisecond},
	}

	return d
}

func (d *Debug) Setup(container *core.Node, app *app.App) {
	d.container = container
	d.app = app
	d.toggleContainer = core.NewNode()

	d.app.Engine.SubscribeID(window.OnKeyDown, &d, d.onKeyDown)
	d.app.Engine.SubscribeID(window.OnWindowSize, &d, d.onResize)

	d.delays = make(map[*DebugStat]time.Duration)
	for _, stat := range d.Stats {
		// Force stat update at first frame
		d.delays[stat] = stat.delay + time.Millisecond
	}

	d.setupGui()
	d.onResize("", nil)

	if d.isActive {
		d.isActive = false
		d.Toggle()
	}
}

func (d *Debug) Toggle() {
	if d.isActive {
		d.container.Remove(d.toggleContainer)
	} else {
		d.container.Add(d.toggleContainer)
	}
	d.isActive = !d.isActive
}

func (d *Debug) Update(deltaTime time.Duration) {
	if !d.isActive {
		return
	}

	for _, stat := range d.Stats {
		if stat.update != nil {
			if stat.delay > 0 {
				d.delays[stat] += deltaTime
				if d.delays[stat] < stat.delay {
					continue
				}
				d.delays[stat] = 0
			}
			stat.update(deltaTime, stat.value)
		}
	}
	d.alignStatsValues()
}

func (d *Debug) Dispose() {
	d.app.Engine.UnsubscribeID(window.OnKeyDown, &d)
	d.app.Engine.UnsubscribeID(window.OnWindowSize, &d)
}

func (d *Debug) onKeyDown(name string, ev interface{}) {
	kev := ev.(*window.KeyEvent)
	if kev.Key == window.KeyF3 {
		d.Toggle()
	}
	if kev.Key == window.KeyF2 {
		d.ToggleWireframe()
	}
}

func (d *Debug) onResize(name string, ev interface{}) {
	if d.panel != nil {
		wWidth, _ := d.app.Engine.GetSize()
		d.panel.SetPosition(float32(wWidth)-d.panel.Width(), 0)
	}
}

func (d *Debug) setupGui() {
	offset := float32(0)

	d.panel = gui.NewPanel(1, 1)
	d.panel.SetColor(&math32.Color{})
	d.panel.SetPaddings(d.panelPadding, d.panelPadding*2, d.panelPadding, d.panelPadding*2)
	d.panel.SetWidth(200)
	d.panel.SetHeight(200)
	d.toggleContainer.Add(d.panel)

	for _, stat := range d.Stats {
		label := gui.NewLabel(stat.label)
		label.SetColor(math32.NewColor("white"))
		label.SetPositionY(offset)
		d.panel.Add(label)

		if stat.update != nil {
			stat.value = gui.NewLabel("###")
			stat.value.SetColor(&math32.Color{R: .1, G: 1, B: .1})
			stat.value.SetPositionY(offset)
			d.panel.Add(stat.value)
		}

		offset += label.Height() + d.panelPadding*2
	}

	d.panel.SetSize(d.panel.Width(), offset+d.panelPadding)
	d.alignStatsValues()
}

func (d *Debug) alignStatsValues() {
	for _, stat := range d.Stats {
		stat.value.SetPositionX(d.panel.Width() - stat.value.Width() - d.panelPadding*4)
	}
}

func (d *Debug) updateFps(_ time.Duration, label *gui.Label) {
	fps, potFps, ok := d.app.Framerater.FPS(300 * time.Millisecond)
	if ok {
		label.SetText(fmt.Sprintf("%d/%d", int(fps), int(potFps)))
	}
}

func (d *Debug) updateScenes(deltaTime time.Duration, label *gui.Label) {
	label.SetText(fmt.Sprintf("%d", len(d.app.Scenes)))
}

func (d *Debug) updateMeshes(deltaTime time.Duration, label *gui.Label) {
	label.SetText(fmt.Sprintf("%d", d.countMeshes(d.app.RootNode)))
}

func (d *Debug) updatePolygons(deltaTime time.Duration, label *gui.Label) {
	label.SetText(fmt.Sprintf("%d", d.countPolygons(d.app.RootNode)))
}

func (d *Debug) countMeshes(node *core.Node) int {
	meshes := 0
	for _, child := range node.Children() {
		switch child.(type) {
		case *graphic.Mesh:
			meshes++
		case *core.Node:
			meshes += d.countMeshes(child.(*core.Node))
		}
	}
	return meshes
}

func (d *Debug) countPolygons(node *core.Node) int {
	meshes := 0
	for _, child := range node.Children() {
		switch child.(type) {
		case *graphic.Mesh:
			indices := child.(*graphic.Mesh).GetGeometry().Indices()
			meshes += indices.Len() / 3
		case *core.Node:
			meshes += d.countPolygons(child.(*core.Node))
		}
	}
	return meshes
}

func (d *Debug) updateCamPosition(deltaTime time.Duration, label *gui.Label) {
	pos := d.app.Cam.Position()
	label.SetText(fmt.Sprintf("%.2f, %.2f, %.2f", pos.X, pos.Y, pos.Z))
}

func (d *Debug) ToggleWireframe() {
	d.wireFrame = !d.wireFrame
	d.setWireFrame(d.app.RootNode, d.wireFrame)
}

func (d *Debug) setWireFrame(node *core.Node, wireframe bool) {
	for _, child := range node.Children() {
		switch child.(type) {
		case *graphic.Mesh:
			materials := child.(*graphic.Mesh).Materials()
			for _, mat := range materials {
				switch mat.IMaterial().(type) {
				case *engineMaterial.Standard:
					mat.IMaterial().(*engineMaterial.Standard).SetWireframe(wireframe)
				}
			}
		case *core.Node:
			d.setWireFrame(child.(*core.Node), wireframe)
		}
	}
}
