package main

import "github.com/thomas-holmes/gterm"
import "github.com/veandco/go-sdl2/sdl"

type EndGameMenu struct {
	world *World

	Content          []string
	ContentColor     sdl.Color
	ContentRelativeX int
	ContentRelativeY int

	PopMenu
}

func (pop *EndGameMenu) Update(input InputEvent) bool {
	return false
}

func (pop *EndGameMenu) RenderBorder(window *gterm.Window) {
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

func (pop *EndGameMenu) RenderContents(window *gterm.Window) {
	xOffset := pop.X + pop.ContentRelativeX
	yOffset := pop.Y + pop.ContentRelativeY

	for line, content := range pop.Content {
		window.PutString(xOffset, yOffset+line, content, pop.ContentColor)
	}
}

func (pop *EndGameMenu) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.RenderBorder(window)
	pop.RenderContents(window)
}

func NewEndGameMenu(x int, y int, w int, h int, color sdl.Color, contents ...string) EndGameMenu {
	contentLen := len(contents)
	maxWidth := 0

	for _, content := range contents {
		thisLen := len(content)
		if thisLen > maxWidth {
			maxWidth = thisLen
		}
	}

	centeredXOffset := (w - maxWidth) / 2
	centeredYOffset := (h - contentLen) / 2

	pop := EndGameMenu{
		PopMenu: PopMenu{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
		ContentRelativeX: centeredXOffset,
		ContentRelativeY: centeredYOffset,
		Content:          contents,
		ContentColor:     color,
	}

	return pop
}
