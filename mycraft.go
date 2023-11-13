package main

import (
	"mycraft/app"
	"mycraft/scene"
)

func main() {
	a := app.NewApp("Mycraft", 60)
	a.AddScene(scene.NewGameScene())
	a.AddScene(scene.NewDebugScene(true))
	a.Run()
}
