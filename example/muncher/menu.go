package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Menu interface {
	Update(sdl.Event) bool
	Render(window *gterm.Window)
	Done() bool
}
