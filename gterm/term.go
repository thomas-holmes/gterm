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
	FontSize     int
	heightPixel  int
	widthPixel   int
	font         string
	sdlWindow    *sdl.Window
	sdlRenderer  *sdl.Renderer
	eventManager *eventManager
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, tileSize int, font string) *Window {
	eventManager := newEventManager()
	window := &Window{
		Columns:      columns,
		Rows:         rows,
		FontSize:     tileSize,
		font:         font,
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

func computeTileSize(font string, fontSize int) (width int, height int, err error) {
	fontFile, err := ttf.OpenFont(font, fontSize)
	if err != nil {
		return 0, 0, err
	}
	atGlyph, err := fontFile.RenderUTF8_Blended("@", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0, 0, err
	}
	return int(atGlyph.W), int(atGlyph.H), nil
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
	tileWidth, tileHeight, err := computeTileSize(window.font, window.FontSize)
	window.heightPixel = tileHeight * window.Rows
	window.widthPixel = tileWidth * window.Columns

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
