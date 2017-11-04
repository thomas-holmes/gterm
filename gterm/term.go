package gterm

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// EventManager is
type eventManager struct {
	id            int
	inputHandlers map[int]InputHandler
}

func newEventManager() eventManager {
	eventManager := eventManager{
		id:            0,
		inputHandlers: make(map[int]InputHandler),
	}

	return eventManager
}

// RegisterInputHandler registers an InputHandler
func (eventManager *eventManager) RegisterInputHandler(handler InputHandler) int {
	eventManager.id++
	eventManager.inputHandlers[eventManager.id] = handler
	return eventManager.id
}

// UnregisterInputHandler removes an input handler, returning true if an InputHandler
// was registered with the provided id, false otherwise
func (eventManager *eventManager) UnregisterInputHandler(id int) bool {
	_, ok := eventManager.inputHandlers[id]
	if ok {
		delete(eventManager.inputHandlers, id)
		return true
	}

	return false
}

// ProcessInputEvent provides the input event to each handler and returns the number of
// handlers executed
func (eventManager *eventManager) ProcessInputEvent(event sdl.Event) int {
	count := 0
	for _, handler := range eventManager.inputHandlers {
		handler(event)
		count++
	}
	return count
}

// InputHandler a function that processes an sdl.Event
type InputHandler func(sdl.Event)

// Window represents the base window object
type Window struct {
	Columns      int
	Rows         int
	TileSize     int
	heightPixel  int
	widthPixel   int
	sdlWindow    *sdl.Window
	sdlRenderer  *sdl.Renderer
	eventManager *eventManager
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, tileSize int) *Window {
	eventManager := newEventManager()
	window := &Window{
		Columns:      columns,
		Rows:         rows,
		TileSize:     tileSize,
		heightPixel:  rows * tileSize,
		widthPixel:   columns * tileSize,
		eventManager: &eventManager,
	}

	return window
}

// RegisterInputHandler save an input handler to be processed during the main event loop
// Returns an int identifier for the input handler so it can be removed later.
func (window *Window) RegisterInputHandler(handler InputHandler) int {
	return window.eventManager.RegisterInputHandler(handler)
}

// UnregisterInputHandler Unregister a handler by its id. Returns true if it was found
// in the map of handlers and removed otherwise returns false
func (window *Window) UnregisterInputHandler(handlerID int) bool {
	return window.eventManager.UnregisterInputHandler(handlerID)
}

// Run is a blocking call that starts the SDL rendering loop
func (window *Window) Run() {
	window.startRenderLoop()
}
func (window *Window) startRenderLoop() {
	err := window.sdlRenderer.SetDrawColor(200, 200, 225, 255)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}
	for {
		window.sdlRenderer.Clear()

		if event := sdl.PollEvent(); event != nil {
			window.eventManager.ProcessInputEvent(event)
		}
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

	return nil
}
