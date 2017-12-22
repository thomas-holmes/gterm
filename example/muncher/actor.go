package main

import "github.com/veandco/go-sdl2/sdl"

type Actor interface {
	CanAct() bool
	Update(turn uint64, event sdl.Event, world *World) bool
	NeedsInput() bool
}
