package main

import (
	"go-pong/gui"
	"log"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	err := gui.Init(app)
	if err != nil {
		log.Fatal("Failed to init app")
	}
}
