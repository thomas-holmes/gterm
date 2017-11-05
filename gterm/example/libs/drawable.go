package libs

import "github.com/veandco/go-sdl2/sdl"

type Drawable interface {
	X() int
	Y() int
	Name() string
	Glyph() string
	ForegroundColor() sdl.Color
	BackgroundColor() sdl.Color
}

// A panel might look like:
//   ┌──────────┐
//   │          │
//   │          │
//   │          │
//   │          │
//   └──────────┘

// Box drawing characters
const BoxVertical = "│"
const BoxHorizontal = "─"
const BoxTopLeft = "┌"
const BoxTopRight = "┐"
const BoxBottomLeft = "└"
const BoxBottomRight = "┘"
