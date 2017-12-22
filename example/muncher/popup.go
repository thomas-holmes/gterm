package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type PopUp struct {
	X                int
	Y                int
	W                int
	H                int
	Content          []string
	ContentColor     sdl.Color
	ContentRelativeX int
	ContentRelativeY int

	Shown bool

	done bool

	Messaging
}

func NewPopUp(xPos int, yPos int, width int, height int, color sdl.Color, contents ...string) PopUp {
	contentLen := len(contents)
	maxWidth := 0

	for _, content := range contents {
		thisLen := len(content)
		if thisLen > maxWidth {
			maxWidth = thisLen
		}
	}

	centeredXOffset := (width - maxWidth) / 2
	centeredYOffset := (height - contentLen) / 2

	pop := PopUp{
		X:                xPos,
		Y:                yPos,
		W:                width,
		H:                height,
		ContentRelativeX: centeredXOffset,
		ContentRelativeY: centeredYOffset,
		Content:          contents,
		ContentColor:     color,
	}

	return pop
}

func (pop PopUp) Done() bool {
	return pop.done
}

func (pop PopUp) ClearUnderlying() {
	pop.Broadcast(ClearRegion, ClearRegionMessage{
		X: pop.X,
		Y: pop.Y,
		W: pop.W,
		H: pop.H,
	})
}

func (pop *PopUp) Show() {
	pop.Shown = true
	pop.Broadcast(PopUpShown, nil)
	// Broadcast a message stopping other game systems
}

func (pop *PopUp) Hide() {
	pop.Shown = false
	pop.Broadcast(PopUpHidden, nil)
	// Broadcast message resuming other game systems
}

func (pop *PopUp) RenderBorder(window *gterm.Window) {
	leftX := pop.X
	rightX := pop.X + pop.W - 1

	topY := pop.Y
	bottomY := pop.Y + pop.H - 1

	for y := topY; y <= bottomY; y++ {
		for x := leftX; x <= rightX; x++ {
			if y == topY || y == bottomY || x == leftX || x == rightX {
				window.PutRune(x, y, '%', pop.ContentColor, gterm.NoColor)
			}
		}
	}
}

func (pop PopUp) RenderContents(window *gterm.Window) {
	xOffset := pop.X + pop.ContentRelativeX
	yOffset := pop.Y + pop.ContentRelativeY

	for line, content := range pop.Content {
		window.PutString(xOffset, yOffset+line, content, pop.ContentColor)
	}
}

func (pop PopUp) Render(window *gterm.Window) {
	pop.ClearUnderlying()
	pop.RenderBorder(window)
	pop.RenderContents(window)
}

func (pop *PopUp) Update(event sdl.Event) bool {
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
