package main

import (
	"log"

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
	window := gterm.NewWindow(40, 40, 16)

	window.Init()
	window.RegisterInputHandler(logOnKeyPress)

	window.Run()
}
