package app

import (
	"github.com/g3n/engine/core"
	"time"
)

type Scene interface {
	Setup(container *core.Node, app *App)
	Update(deltaTime time.Duration)
	Dispose()
}
