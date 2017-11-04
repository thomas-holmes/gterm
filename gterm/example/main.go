package main

import (
	"log"
	"path"
	"time"

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

func getFpsHandler() gterm.RenderHandler {
	frames := 0
	last := time.Now().Second()
	handler := func(renderer *sdl.Renderer) {
		frames++
		now := time.Now().Second()
		if now != last {
			last = now
			log.Printf("Frame rate is %v", frames)
			frames = 0
		}
	}
	return handler
}

func main() {
	window := gterm.NewWindow(80, 24, 16, path.Join("assets", "font", "FiraMono-Regular.ttf"))

	window.Init()
	window.RegisterInputHandler(logOnKeyPress)
	window.RegisterRenderHandler(getFpsHandler())

	window.Run()
}
