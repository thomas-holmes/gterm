package game

import (
	"log"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Tile struct {
	Dirty           bool
	BackgroundColor sdl.Color
	BackgroundGlyph string
	Wall            bool
}

func NewTile() Tile {
	return Tile{
		Dirty:           true,
		BackgroundColor: sdl.Color{R: 192, G: 192, B: 192, A: 255},
		BackgroundGlyph: ".",
	}
}

func (tile Tile) RenderBackground(col int, row int, window *gterm.Window) {
	err := window.AddToCell(col, row, tile.BackgroundGlyph, tile.BackgroundColor)
	if err != nil {
		log.Println("Failed to render background", err)
	}
}
