package camera

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
	"time"
)

const rad90 = 100 * math32.Pi / 180

type DirectionStates struct {
	Forward  bool
	Backward bool
	Left     bool
	Right    bool
}

type WASMControl struct {
	Camera              *camera.Camera
	VelocityForward     float32
	VelocityForwardCap  float32
	VelocitySideward    float32
	VelocitySidewardCap float32
	Acceleration        float32
	Sensitivity         float32

	directions          DirectionStates
	capturedMouseWindow *window.GlfwWindow
	cursorLastX         float32
	cursorLastY         float32
	yRotation           float32
	xRotation           float32
}

func NewWASMControl(cam *camera.Camera) *WASMControl {
	wc := new(WASMControl)
	wc.directions = DirectionStates{Forward: false, Backward: false, Left: false, Right: false}
	wc.Acceleration = 1
	wc.VelocityForwardCap = 10
	wc.VelocitySidewardCap = 10
	wc.Camera = cam
	wc.Sensitivity = 1

	// Subscribe to events
	gui.Manager().SubscribeID(window.OnKeyDown, &wc, wc.onKey)
	gui.Manager().SubscribeID(window.OnKeyUp, &wc, wc.onKey)
	gui.Manager().SubscribeID(window.OnCursor, &wc, wc.onCursor)

	return wc
}

func (wc *WASMControl) CaptureMouse(window *window.GlfwWindow) {
	wc.capturedMouseWindow = window
	wc.capturedMouseWindow.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func (wc *WASMControl) ReleaseMouse() {
	wc.capturedMouseWindow.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	wc.capturedMouseWindow = nil
}

func (wc *WASMControl) Update(deltaTime time.Duration) {
	deltaFact := float32(deltaTime.Milliseconds()) / 1000

	wc.updateVelocity()
	wc.updateRotation()
	wc.updatePosition(deltaFact)
	wc.doCaptureMouse()
}

func (wc *WASMControl) onCursor(_ string, ev interface{}) {
	mev := ev.(*window.CursorEvent)

	wc.yRotation += math32.DegToRad(wc.cursorLastX-mev.Xpos) * wc.Sensitivity / 10
	wc.xRotation += math32.DegToRad(wc.cursorLastY-mev.Ypos) * wc.Sensitivity / 10

	wc.xRotation = math32.Max(-rad90, math32.Min(rad90, wc.xRotation))

	wc.cursorLastX = mev.Xpos
	wc.cursorLastY = mev.Ypos
}

func (wc *WASMControl) onKey(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	state := false
	if evname == "w.OnKeyDown" {
		state = true
	}

	switch kev.Key {
	case window.KeyUp, window.KeyW:
		wc.directions.Forward = state
	case window.KeyDown, window.KeyS:
		wc.directions.Backward = state
	case window.KeyLeft, window.KeyA:
		wc.directions.Left = state
	case window.KeyRight, window.KeyD:
		wc.directions.Right = state
	}
}

func (wc *WASMControl) updateVelocity() {
	// Accelerate on active directions
	if wc.directions.Forward {
		wc.VelocityForward -= wc.Acceleration
	}

	if wc.directions.Backward {
		wc.VelocityForward += wc.Acceleration
	}

	if wc.directions.Right {
		wc.VelocitySideward += wc.Acceleration
	}

	if wc.directions.Left {
		wc.VelocitySideward -= wc.Acceleration
	}

	// Return to 0 velocity on inactive directions
	if !wc.directions.Left && !wc.directions.Right && wc.VelocitySideward != 0 {
		if wc.VelocitySideward > 0 {
			wc.VelocitySideward = math32.Max(0, wc.VelocitySideward-wc.Acceleration)
		} else {
			wc.VelocitySideward = math32.Min(0, wc.VelocitySideward+wc.Acceleration)
		}
	}

	if !wc.directions.Forward && !wc.directions.Backward && wc.VelocityForward != 0 {
		if wc.VelocityForward > 0 {
			wc.VelocityForward = math32.Max(0, wc.VelocityForward-wc.Acceleration)
		} else {
			wc.VelocityForward = math32.Min(0, wc.VelocityForward+wc.Acceleration)
		}
	}

	// Enforce velocity cap
	if wc.VelocityForward > wc.VelocityForwardCap {
		wc.VelocityForward = wc.VelocityForwardCap
	}

	if wc.VelocityForward < -wc.VelocityForwardCap {
		wc.VelocityForward = -wc.VelocityForwardCap
	}

	if wc.VelocitySideward > wc.VelocitySidewardCap {
		wc.VelocitySideward = wc.VelocitySidewardCap
	}

	if wc.VelocitySideward < -wc.VelocitySidewardCap {
		wc.VelocitySideward = -wc.VelocitySidewardCap
	}
}

func (wc *WASMControl) updateRotation() {
	wc.Camera.SetRotation(wc.xRotation, wc.yRotation, 0)
}

func (wc *WASMControl) updatePosition(deltaFact float32) {
	if wc.VelocityForward != 0 {
		phi, theta := math32.DegToRad(360)-wc.xRotation, wc.yRotation
		camPos := wc.Camera.Position()

		moveVec := math32.Vector3{
			X: camPos.X + (math32.Sin(theta) * wc.VelocityForward * deltaFact),
			Y: camPos.Y + (math32.Sin(phi) * wc.VelocityForward * deltaFact),
			Z: camPos.Z + (math32.Cos(phi) * math32.Cos(theta) * wc.VelocityForward * deltaFact),
		}

		wc.Camera.SetPositionVec(&moveVec)
	}

	if wc.VelocitySideward != 0 {
		phi, theta := math32.DegToRad(360)-wc.xRotation, wc.yRotation+math32.DegToRad(90)
		camPos := wc.Camera.Position()

		moveVec := math32.Vector3{
			X: camPos.X + (math32.Sin(theta) * wc.VelocitySideward * deltaFact),
			Y: camPos.Y,
			Z: camPos.Z + (math32.Cos(phi) * math32.Cos(theta) * wc.VelocitySideward * deltaFact),
		}

		wc.Camera.SetPositionVec(&moveVec)
	}

}

func (wc *WASMControl) doCaptureMouse() {
	if wc.capturedMouseWindow != nil {
		width, height := wc.capturedMouseWindow.GetSize()
		wc.cursorLastX, wc.cursorLastY = float32(width)/2, float32(height)/2
		wc.capturedMouseWindow.SetCursorPos(float64(wc.cursorLastX), float64(wc.cursorLastY))
	}
}
