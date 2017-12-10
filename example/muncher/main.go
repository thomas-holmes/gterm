package main

import (
	"log"
	"math/rand"
	"path"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"

	"net/http"
	_ "net/http/pprof"
)

var quit = false

func eventActionable(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.KeyDownEvent:
		return true
	case *sdl.QuitEvent:
		return true
	}
	return false
}

func handleInput(event sdl.Event, world *World) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			quit = true
		case sdl.K_5:
			spawnRandomMonster(world)
		case sdl.K_BACKSLASH:
			world.ToggleScentOverlay()
		case sdl.K_g:
			log.Printf("\n%v", GenLevel(world.rng, 80, 40, GenDownStairs|GenUpStairs))
		}
	case *sdl.QuitEvent:
		quit = true
	}
}

func spawnRandomMonster(world *World) {
	for tries := 0; tries < 100; tries++ {
		x := rand.Intn(world.CurrentLevel.Columns)
		y := rand.Intn(world.CurrentLevel.Rows)

		if world.CanStandOnTile(x, y) {
			level := rand.Intn(8) + 1
			monster := NewMonster(x, y, level, Green, level)
			world.AddEntity(&monster)
			return
		}
	}
}

func main() {
	// Disable FPS limit, generally, so I can monitor performance.
	window := gterm.NewWindow(100, 30, path.Join("assets", "font", "DejaVuSansMono.ttf"), 24, false)

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to Init() window", err)
	}

	window.SetTitle("Muncher")

	window.SetBackgroundColor(gterm.NoColor)

	window.ShouldRenderFps(true)

	world := NewWorld(window, true, 99)
	{
		// TODO: Roll this up into some kind of registering a system function on the world
		combat := CombatSystem{}

		combat.SetMessageBus(&world.MessageBus)
		world.MessageBus.Subscribe(combat)
	}

	player := NewPlayer()
	player.LevelUp()
	player.LevelUp()
	player.LevelUp()
	player.LevelUp()
	player.LevelUp()
	player.LevelUp()

	player.Name = "Euclid"
	level1 := GenLevel(world.rng, 80, 40, GenDownStairs)
	log.Printf("Level1\n%v", level1)
	world.AddLevelFromString(level1)
	level2 := GenLevel(world.rng, 80, 40, GenUpStairs)
	log.Printf("Level2\n%v", level2)
	world.AddLevelFromString(level2)
	world.SetCurrentLevel(0)

	// TODO: Fix Player needs to be added after level creation now
	world.AddPlayer(&player)

	// TODO: Fix Monsters need to be added after level creation now
	for i := 0; i < 10; i++ {
		spawnRandomMonster(world)
	}

	hud := NewHud(&player, world, 60, 0)

	for !quit {

		event := sdl.PollEvent()
		// TODO: Don't advance turn count if we don't actually do something worthwhile
		// Tried to fix this for a while but ended up with some weird issues. Will take
		// some more thinking and probably a further refactor around whether or not the
		// player actually performed an action. Need to decouple turn advancement from
		// input acquisition.
		if world.turnCount == 0 || eventActionable(event) {
			window.ClearWindow()

			handleInput(event, world)

			world.AddInput(event)

			updateLoops := 0
			for !world.Update() && !world.GameOver {
				updateLoops++
			}
			log.Printf("Ran %v update loops", updateLoops)

			world.Render()

			hud.Render(world)
		}

		window.Refresh()
	}
}

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}
