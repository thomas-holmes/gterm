package game

import (
	"log"

	"github.com/thomas-holmes/gterm"
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

func (tile *Tile) Render(col int, row int, world *World) {
	if !tile.Dirty {
		return
	}

	pos := Position{XPos: col, YPos: row}
	items := world.renderItems[pos]
	if len(items) > 0 {
		for _, item := range items {
			item.Render(world)
		}
	} else {
		tile.RenderBackground(col, row, world.Window) // bad API, refactor
	}
	tile.Dirty = false
}

func (tile Tile) RenderBackground(col int, row int, window *gterm.Window) {
	err := window.AddToCell(col, row, tile.BackgroundGlyph, tile.BackgroundColor)
	if err != nil {
		log.Println("Failed to render background", err)
	}
}
