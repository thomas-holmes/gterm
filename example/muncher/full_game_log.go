package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type FullGameLog struct {
	*GameLog

	done bool

	X int
	Y int
	W int
	H int

	ScrollPosition int
}

func (pop *FullGameLog) Update(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_ESCAPE:
			pop.done = true
			return true
		}
	}
	return false
}

func (pop *FullGameLog) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Println("Got an error clearing FullGameLog region", err)
	}

	messagesToRender := min(len(pop.Messages), pop.H)

	yOffset := 0
	for i := messagesToRender - 1; i >= 0; i-- {
		message := pop.Messages[i]
		window.PutString(pop.X, pop.Y+yOffset, message, White)
		yOffset++
	}
}

func (pop FullGameLog) Done() bool {
	return pop.done
}
