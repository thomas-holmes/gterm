package main

import (
	"log"
	"math/rand"
	"path"

	"github.com/thomas-holmes/gterm"
	"github.com/thomas-holmes/gterm/example/muncher/game"
	"github.com/veandco/go-sdl2/sdl"

	"net/http"
	_ "net/http/pprof"
)

var quit = false

func handleInput(event sdl.Event, world *game.World) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			quit = true
		case sdl.K_5:
			spawnRandomMonster(world)
		}
	case *sdl.QuitEvent:
		quit = true
	}
}

var red = sdl.Color{R: 255, G: 0, B: 0, A: 255}

func spawnRandomMonster(world *game.World) {
	for tries := 0; tries < 100; tries++ {
		x := rand.Intn(world.Columns)
		y := rand.Intn(world.Rows)

		if world.CanStandOnTile(x, y) {
			level := rand.Intn(8) + 1
			monster := game.NewMonster(x, y, level, game.Green, level)
			world.AddEntity(&monster)
			return
		}

		log.Println("Faield to find valid position for monster after 100 attempts")
	}
}

func addMonsters(world *game.World) {
	// For some reassigning a single var keeps giving same memory address.
	// I guess go is copying the NewMonster into the same stack allocated struct.
	// I think this makes sense but I'm annoyed that I didn't realize it would happen.
	m1 := game.NewMonster(10, 6, 1, game.Green, 1)
	world.AddEntity(&m1)
	m2 := game.NewMonster(10, 7, 1, game.Green, 1)
	world.AddEntity(&m2)
	m3 := game.NewMonster(10, 8, 2, game.Green, 1)
	world.AddEntity(&m3)
	m4 := game.NewMonster(10, 9, 3, game.Green, 1)
	world.AddEntity(&m4)
	m5 := game.NewMonster(10, 10, 3, game.Green, 1)
	world.AddEntity(&m5)
	m6 := game.NewMonster(10, 11, 4, game.Green, 1)
	world.AddEntity(&m6)
	m7 := game.NewMonster(10, 12, 5, game.Green, 1)
	world.AddEntity(&m7)

}

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
	{
		// TODO: Roll this up into some kind of registering a system function on the world
		combat := game.CombatSystem{}

		combat.SetMessageBus(&world.MessageBus)
		world.MessageBus.Subscribe(combat)
	}

	player := game.NewPlayer(5, 5)
	player.Name = "Euclid"

	world.AddEntity(&player)

	addMonsters(world)

	world.BuildLevelFromMask(LevelMask)

	hud := game.NewHud(&player, world, 60, 0)

	for !quit {
		event := sdl.PollEvent()

		handleInput(event, world)
		world.HandleInput(event)

		world.Render()

		hud.Render(world)

		window.Render()
	}
}

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}

var LevelMask = []int{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
}
