package main

import (
	"github.com/thomas-holmes/gterm"
)

type Menu interface {
	Update(InputEvent) bool
	Render(window *gterm.Window)
	Done() bool
}
