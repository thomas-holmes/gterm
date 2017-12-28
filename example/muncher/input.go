package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type InputEvent struct {
	sdl.Event
	sdl.Keymod
}
