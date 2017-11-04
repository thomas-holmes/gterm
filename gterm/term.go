package gterm

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Window represents the base window object
type Window struct {
	Columns     int
	Rows        int
	TileSize    int
	heightPixel int
	widthPixel  int
	sdlWindow   *sdl.Window
	sdlRenderer *sdl.Renderer
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, tileSize int) *Window {
	return &Window{Columns: columns, Rows: rows, TileSize: tileSize, heightPixel: rows * tileSize, widthPixel: columns * tileSize}
}

func (window *Window) startRenderLoop() {
	err := window.sdlRenderer.SetDrawColor(200, 200, 225, 255)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}
	for {
		window.sdlRenderer.Clear()

		// window.sdlWindow.UpdateSurface()

		window.sdlRenderer.Present()
	}
}

// Init initialized the window for drawing
func (window *Window) Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}
	err = ttf.Init()
	if err != nil {
		return nil
	}

	log.Printf("Creating window w:%v, h:%v", window.widthPixel, window.heightPixel)
	sdlWindow, sdlRenderer, err := sdl.CreateWindowAndRenderer(window.widthPixel,
		window.heightPixel, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	if err != nil {
		return err
	}

	window.sdlWindow = sdlWindow
	window.sdlRenderer = sdlRenderer

	go window.startRenderLoop()

	return nil
}
