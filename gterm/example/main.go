package main

import (
	"log"
	"path"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

func logOnKeyPress(event sdl.Event) {
	log.Printf("%#v", event)
	switch e := event.(type) {
	case sdl.KeyDownEvent:
		log.Println(e)
	}
}

func main() {
	window := gterm.NewWindow(80, 24, 16, path.Join("assets", "font", "FiraMono-Regular.ttf"))

	window.Init()
	window.RegisterInputHandler(logOnKeyPress)

	window.Run()
}
