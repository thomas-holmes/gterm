package main

import (
	"github.com/thomas-holmes/gterm"
)

type PopMenu struct {
	done bool

	X int
	Y int
	W int
	H int
}

func (pop PopMenu) Done() bool {
	return pop.done
}

type Menu interface {
	Update(InputEvent) bool
	Render(window *gterm.Window)
	Done() bool
}
