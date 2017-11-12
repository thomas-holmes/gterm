package main

import (
	"log"
	"path"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/thomas-holmes/sneaker/gterm/example/muncher/game"
	"github.com/veandco/go-sdl2/sdl"
)

var quit = false

func handleInput(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			quit = true
		}
	case *sdl.QuitEvent:
		quit = true
	}
}

var red = sdl.Color{R: 255, G: 0, B: 0, A: 255}

func render(window *gterm.Window, renderable game.Renderable) {
	if renderable.ShouldRender() {
		log.Println("Actually rendering", renderable)
		window.AddToCell(renderable.RenderCol(), renderable.RenderRow(), renderable.RenderGlyph(), renderable.RenderColor())
		renderable.Rendered()
	}
}

func main() {
	window := gterm.NewWindow(80, 24, path.Join("assets", "font", "FiraMono-Regular.ttf"), 16)

	window.Init()

	window.SetTitle("Muncher")

	window.SetBackgroundColor(sdl.Color{R: 0, G: 0, B: 0, A: 0})

	window.ShouldRenderFps(true)

	player := game.NewPlayer(window, 0, 0)

	inputtables := []game.Inputtable{&player}
	renderables := []game.Renderable{&player}

	for !quit {
		if event := sdl.PollEvent(); event != nil {
			handleInput(event)
			for _, inputtable := range inputtables {
				inputtable.HandleInput(event)
			}
		}

		for _, renderable := range renderables {
			render(window, renderable)
		}

		window.Render()
	}
}
