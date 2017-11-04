package gterm

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// InputHandler a function that processes an sdl.Event
type InputHandler func(sdl.Event)

type inputHandlerRequest interface {
	responseChannel() chan<- int
}

type inputHandlerAddRequest struct {
	handler  InputHandler
	response chan<- int
}

func (request inputHandlerAddRequest) responseChannel() chan<- int {
	return request.response
}

type inputHandlerRemoveRequest struct {
	id       int
	response chan<- int
}

func (request inputHandlerRemoveRequest) responseChannel() chan<- int {
	return request.response
}

type inputHandlerProcessHandlers struct {
	event    sdl.Event
	response chan<- int
}

func (request inputHandlerProcessHandlers) responseChannel() chan<- int {
	return request.response
}

// Window represents the base window object
type Window struct {
	Columns          int
	Rows             int
	TileSize         int
	heightPixel      int
	widthPixel       int
	sdlWindow        *sdl.Window
	sdlRenderer      *sdl.Renderer
	inputHandlers    map[int]InputHandler
	inputHandlerChan chan inputHandlerRequest
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, tileSize int) *Window {
	window := &Window{
		Columns:          columns,
		Rows:             rows,
		TileSize:         tileSize,
		heightPixel:      rows * tileSize,
		widthPixel:       columns * tileSize,
		inputHandlers:    make(map[int]InputHandler),
		inputHandlerChan: make(chan inputHandlerRequest),
	}

	return window
}

// RegisterInputHandler save an input handler to be processed during the main event loop
// Returns an int identifier for the input handler so it can be removed later.
func (window *Window) RegisterInputHandler(handler InputHandler) int {
	responseChannel := make(chan int)
	request := inputHandlerAddRequest{handler, responseChannel}
	window.inputHandlerChan <- request
	return <-responseChannel
}

// UnregisterInputHandler Unregister a handler by its id. Returns true if it was found
// in the map of handlers and removed otherwise returns false
func (window *Window) UnregisterInputHandler(handlerID int) bool {
	responseChannel := make(chan int)
	request := inputHandlerRemoveRequest{handlerID, responseChannel}
	window.inputHandlerChan <- request
	response := <-responseChannel
	return response == handlerID
}

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
			response := make(chan int)
			request := inputHandlerProcessHandlers{event, response}
			window.inputHandlerChan <- request
			processed := <-response
			log.Printf("Processed %v handlers", processed)
		}
		// window.sdlWindow.UpdateSurface()

		window.sdlRenderer.Present()
	}
}

func (window *Window) startEventManager() {
	id := 1
	for request := range window.inputHandlerChan {
		switch r := request.(type) {
		case inputHandlerAddRequest:
			window.inputHandlers[id] = r.handler
			r.response <- id
			id++
		case inputHandlerRemoveRequest:
			_, ok := window.inputHandlers[r.id]
			if ok {
				delete(window.inputHandlers, r.id)
				r.response <- r.id
			} else {
				r.response <- 0
			}
		case inputHandlerProcessHandlers:
			count := 0
			for _, handler := range window.inputHandlers {
				handler(r.event)
				count++
			}
			r.response <- count
		}
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

	go window.startEventManager()
	// go window.startRenderLoop()

	return nil
}
