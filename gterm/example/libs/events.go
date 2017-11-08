package libs

import (
	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

// EventManager
type EventManager struct {
	id             int
	inputHandlers  map[int]InputHandler
	renderHandlers []RenderHandler
	window         *gterm.Window
}

func NewEventManager(window *gterm.Window) EventManager {
	eventManager := EventManager{
		id:             0,
		inputHandlers:  make(map[int]InputHandler),
		renderHandlers: make([]RenderHandler, 1),
		window:         window,
	}

	return eventManager
}

// RenderHandler supplies an SDL renderer for you draw on
type RenderHandler func(window *gterm.Window)

// RegisterRenderHandler register a rendering handler to draw to the screen
func (eventManager *EventManager) RegisterRenderHandler(handler RenderHandler) {
	eventManager.renderHandlers = append(eventManager.renderHandlers, handler)
}

func (eventManager *EventManager) RunRenderHandlers() {
	for _, handler := range eventManager.renderHandlers {
		if handler != nil {
			handler(eventManager.window)
		}
	}
}

// RegisterInputHandler registers an InputHandler
func (eventManager *EventManager) RegisterInputHandler(handler InputHandler) int {
	eventManager.id++
	eventManager.inputHandlers[eventManager.id] = handler
	return eventManager.id
}

// UnregisterInputHandler removes an input handler, returning true if an InputHandler
// was registered with the provided id, false otherwise
func (eventManager *EventManager) UnregisterInputHandler(id int) bool {
	_, ok := eventManager.inputHandlers[id]
	if ok {
		delete(eventManager.inputHandlers, id)
		return true
	}

	return false
}

// ProcessInputEvent provides the input event to each handler and returns the number of
// handlers executed
func (eventManager *EventManager) RunInputHandlers(event sdl.Event) int {
	count := 0
	for _, handler := range eventManager.inputHandlers {
		handler(event)
		count++
	}
	return count
}

// InputHandler a function that processes an sdl.Event
type InputHandler func(sdl.Event)
