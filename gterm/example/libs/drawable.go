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
