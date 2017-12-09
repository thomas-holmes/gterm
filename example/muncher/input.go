package main

import "github.com/veandco/go-sdl2/sdl"

type Inputtable interface {
	HandleInput(event sdl.Event, world *World) bool
	NeedsInput() bool
}
