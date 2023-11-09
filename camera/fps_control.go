package camera

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"time"
)

type DirectionStates struct {
	Forward  bool
	Backward bool
	Left     bool
	Right    bool
}

type FPSControl struct {
	Camera              *camera.Camera
	VelocityForward     float32
	VelocityForwardCap  float32
	VelocitySideward    float32
	VelocitySidewardCap float32
	Acceleration        float32
	Directions          DirectionStates
}

func NewFPSControl(cam *camera.Camera) *FPSControl {
	fc := new(FPSControl)
	fc.Directions = DirectionStates{Forward: false, Backward: false, Left: false, Right: false}
	fc.Acceleration = 1
	fc.VelocityForwardCap = 20
	fc.VelocitySidewardCap = 20
	fc.Camera = cam

	// Subscribe to events
	gui.Manager().SubscribeID(window.OnKeyDown, &fc, fc.onKey)
	gui.Manager().SubscribeID(window.OnKeyUp, &fc, fc.onKey)

	// side moves targets are +/-90° on y axis (I guess) according to forward direction
	fc.Camera.RotateX(360 * math32.Pi / 180) // 360 to radian conversion up/down
	fc.Camera.RotateY(360 * math32.Pi / 180) // 360 to radian conversion left/right

	return fc
}

func (fc *FPSControl) onKey(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	state := false
	if evname == "w.OnKeyDown" {
		state = true
	}

	switch kev.Key {
	case window.KeyUp:
		fc.Directions.Forward = state
	case window.KeyDown:
		fc.Directions.Backward = state
	case window.KeyLeft:
		fc.Directions.Left = state
	case window.KeyRight:
		fc.Directions.Right = state
	}
}

func (fc *FPSControl) velocityUpdate() {
	// Accelerate on active directions
	if fc.Directions.Forward {
		fc.VelocityForward -= fc.Acceleration
	}

	if fc.Directions.Backward {
		fc.VelocityForward += fc.Acceleration
	}

	if fc.Directions.Right {
		fc.VelocitySideward -= fc.Acceleration
	}

	if fc.Directions.Left {
		fc.VelocitySideward += fc.Acceleration
	}

	// Return to 0 velocity on inactive directions
	if !fc.Directions.Left && !fc.Directions.Right && fc.VelocitySideward != 0 {
		if fc.VelocitySideward > 0 {
			fc.VelocitySideward = math32.Max(0, fc.VelocitySideward-fc.Acceleration)
		} else {
			fc.VelocitySideward = math32.Min(0, fc.VelocitySideward+fc.Acceleration)
		}
	}

	if !fc.Directions.Forward && !fc.Directions.Backward && fc.VelocityForward != 0 {
		if fc.VelocityForward > 0 {
			fc.VelocityForward = math32.Max(0, fc.VelocityForward-fc.Acceleration)
		} else {
			fc.VelocityForward = math32.Min(0, fc.VelocityForward+fc.Acceleration)
		}
	}

	// Enforce velocity cap
	if fc.VelocityForward > fc.VelocityForwardCap {
		fc.VelocityForward = fc.VelocityForwardCap
	}

	if fc.VelocityForward < -fc.VelocityForwardCap {
		fc.VelocityForward = -fc.VelocityForwardCap
	}

	if fc.VelocitySideward > fc.VelocitySidewardCap {
		fc.VelocitySideward = fc.VelocitySidewardCap
	}

	if fc.VelocitySideward < -fc.VelocitySidewardCap {
		fc.VelocitySideward = -fc.VelocitySidewardCap
	}
}

func (fc *FPSControl) Update(deltaTime time.Duration) {
	fc.velocityUpdate()

	deltaFact := float32(deltaTime.Milliseconds()) / 1000

	position := fc.Camera.Position()
	position.X += fc.VelocitySideward * deltaFact
	position.Z += fc.VelocityForward * deltaFact

	fc.Camera.SetPositionVec(&position)
}
