package main

import (
	"log"
	"path"
	"time"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/thomas-holmes/sneaker/gterm/example/libs"
	"github.com/veandco/go-sdl2/sdl"
)

func logOnKeyPress(event sdl.Event) {
	// log.Printf("%#v", event)
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
			case sdl.K_SPACE:
				panel.Update(panel.XPos, panel.YPos+1, panel.Width, panel.Height)
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

func main() {
	window := gterm.NewWindow(80, 24, path.Join("assets", "font", "FiraMono-Regular.ttf"), 16)

	if err := window.Init(); err != nil {
		log.Fatal(err)
	}

	panelManager := libs.NewPanelManager(window)
	eventManager := libs.NewEventManager(window)

	panel := panelManager.NewPanel(20, 2, 40, 10, 1)
	eventManager.RegisterInputHandler(panelAdjusterInputHandler(panel))
	eventManager.RegisterRenderHandler(getFpsHandler())

	color := sdl.Color{R: 255, G: 0, B: 0, A: 0}
	window.AddToCell(0, 0, "A", color)
	window.AddToCell(1, 0, "B", color)
	window.AddToCell(0, 1, "C", color)

	quit := false
	for !quit {
		window.AddToCell(24, 12, "@", color)

		if event := sdl.PollEvent(); event != nil {
			eventManager.RunInputHandlers(event)
		}
		eventManager.RunRenderHandlers()

		panelManager.RenderPanels()

		window.Render()
	}
}
