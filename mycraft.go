package main

import (
	"mycraft/app"
	"mycraft/app/scene"
)

func main() {
	a := app.NewApp("Mycraft", 60)
	a.AddScene(scene.NewGameScene())
	a.AddScene(scene.NewDebugScene())
	a.Run()
}
