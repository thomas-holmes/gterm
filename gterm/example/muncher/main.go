package main

import (
	"log"
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

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to Init() window", err)
	}

	window.SetTitle("Muncher")

	window.SetBackgroundColor(sdl.Color{R: 0, G: 0, B: 0, A: 0})

	window.ShouldRenderFps(true)

	world := game.NewWorld(window, 40, 18)

	player := game.NewPlayer(world.GetNextID(), 5, 5)
	player.Name = "Euclid"

	monster := game.NewMonster(world.GetNextID(), 10, 10, "1", game.Green, 1)

	world.BuildLevelFromMask(game.LevelMask)

	hud := game.NewHud(&player, &world, 60, 0)

	world.AddEntity(&player)
	world.AddEntity(&monster)

	for !quit {
		if event := sdl.PollEvent(); event != nil {
			handleInput(event)
			player.HandleInput(event, &world)
		}
		world.Render()

		hud.Render(&world)

		window.Render()
	}
}

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}
