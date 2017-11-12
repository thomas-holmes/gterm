package main

import (
	"path"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/thomas-holmes/sneaker/gterm/example/muncher/game"
	"github.com/veandco/go-sdl2/sdl"

	"net/http"
	_ "net/http/pprof"
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

func main() {
	// Disable FPS limit, generally, so I can monitor performance.
	window := gterm.NewWindow(80, 24, path.Join("assets", "font", "FiraMono-Regular.ttf"), 16, 0)

	window.Init()

	window.SetTitle("Muncher")

	window.SetBackgroundColor(sdl.Color{R: 0, G: 0, B: 0, A: 0})

	window.ShouldRenderFps(true)

	world := game.NewWorld(window, 40, 18)

	player := game.NewPlayer(&world, 5, 5)
	player.Name = "Euclid"

	world.BuildLevelFromMask(game.LevelMask)

	hud := game.NewHud(&player, &world, 60, 0)

	for !quit {
		if event := sdl.PollEvent(); event != nil {
			handleInput(event)
			player.HandleInput(event)
		}
		world.Render()

		hud.Render(&world)

		window.Render()
	}
}

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}
