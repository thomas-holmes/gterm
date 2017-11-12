package main

import (
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/thomas-holmes/sneaker/gterm/example/libs"
	"github.com/veandco/go-sdl2/sdl"

	_ "net/http/pprof"
)

func logOnKeyPress(event sdl.Event) {
	switch e := event.(type) {
	case sdl.KeyDownEvent:
		log.Println(e)
	}
}

func getFpsHandler() libs.RenderHandler {
	frames := 0
	last := time.Now().Second()
	handler := func(window *gterm.Window) {
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

func panelAdjusterInputHandler(panel *libs.Panel) libs.InputHandler {
	return func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.KeyDownEvent:
			switch e.Keysym.Sym {
			case sdl.K_UP:
				panel.Update(panel.XPos, panel.YPos-1, panel.Width, panel.Height)
			case sdl.K_DOWN:
				panel.Update(panel.XPos, panel.YPos+1, panel.Width, panel.Height)
			case sdl.K_LEFT:
				panel.Update(panel.XPos-1, panel.YPos, panel.Width, panel.Height)
			case sdl.K_RIGHT:
				panel.Update(panel.XPos+1, panel.YPos, panel.Width, panel.Height)
			case sdl.K_EQUALS:
				panel.Update(panel.XPos, panel.YPos, panel.Width+1, panel.Height)
			case sdl.K_MINUS:
				panel.Update(panel.XPos, panel.YPos, panel.Width-1, panel.Height)
			case sdl.K_LEFTBRACKET:
				panel.Update(panel.XPos, panel.YPos, panel.Width, panel.Height-1)
			case sdl.K_RIGHTBRACKET:
				panel.Update(panel.XPos, panel.YPos, panel.Width, panel.Height+1)
			}
		}
	}
}

type Player struct {
	XPos   int
	YPos   int
	FColor sdl.Color
	Glyph  string
}

func (player *Player) handleInput(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_h:
			player.XPos--
		case sdl.K_j:
			player.YPos++
		case sdl.K_k:
			player.YPos--
		case sdl.K_l:
			player.XPos++
		case sdl.K_b:
			player.XPos--
			player.YPos++
		case sdl.K_n:
			player.XPos++
			player.YPos++
		case sdl.K_y:
			player.XPos--
			player.YPos--
		case sdl.K_u:
			player.XPos++
			player.YPos--
		}
	}
}

func (player *Player) handleRender(window *gterm.Window) {
	window.AddToCell(player.XPos, player.YPos, player.Glyph, player.FColor)
}

func renderEverywhere(window *gterm.Window) {
	for row := 0; row < window.Rows; row++ {
		err := window.AddToCell(0, row, strings.Repeat(".", window.Columns), sdl.Color{R: 115, G: 115, B: 115, A: 255})
		if err != nil {
			log.Fatalln("Failed to render everywhere", err)
		}
	}
}

func quitHandler(event sdl.Event) {
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

var quit = false

func main() {
	go http.ListenAndServe("localhost:6060", nil)

	window := gterm.NewWindow(80, 24, path.Join("assets", "font", "FiraMono-Regular.ttf"), 16, 0)

	if err := window.Init(); err != nil {
		log.Fatal(err)
	}

	window.SetTitle("gterm Example App")

	panelManager := libs.NewPanelManager(window)
	eventManager := libs.NewEventManager(window)

	player := Player{10, 10, sdl.Color{R: 255, G: 25, B: 55, A: 255}, "@"}

	panel := panelManager.NewPanel(10, 2, 20, 20, 1)
	eventManager.RegisterInputHandler(panelAdjusterInputHandler(panel))
	eventManager.RegisterInputHandler(player.handleInput)
	eventManager.RegisterInputHandler(quitHandler)
	eventManager.RegisterRenderHandler(getFpsHandler())
	eventManager.RegisterRenderHandler(panelManager.HandleRender)
	eventManager.RegisterRenderHandler(player.handleRender)

	window.ShouldRenderFps(true)

	for !quit {
		window.ClearWindow()

		if event := sdl.PollEvent(); event != nil {
			eventManager.RunInputHandlers(event)
		}
		renderEverywhere(window)

		eventManager.RunRenderHandlers()

		window.Render()
	}
}
