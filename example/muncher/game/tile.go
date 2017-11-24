package game

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Tile struct {
	XPos            int
	YPos            int
	BackgroundColor sdl.Color
	BackgroundGlyph rune
	Wall            bool
	WasVisible      bool
}

func NewTile(x int, y int) Tile {
	return Tile{
		XPos:            x,
		YPos:            y,
		BackgroundColor: sdl.Color{R: 192, G: 192, B: 192, A: 255},
		BackgroundGlyph: '.',
	}
}

func (tile *Tile) Render(world *World, visibility Visibility) {
	if visibility == Unseen {
		return
	}

	pos := Position{XPos: tile.XPos, YPos: tile.YPos}
	items := world.renderItems[pos]
	if len(items) > 0 && visibility == Visible {
		for _, item := range items {
			item.Render(world)
		}
	} else {
		tile.RenderBackground(world, visibility) // bad API, refactor
	}
}

func (tile Tile) RenderBackground(world *World, visibility Visibility) {
	color := tile.BackgroundColor

	if visibility == Seen {
		color.R /= 2
		color.G /= 2
		color.B /= 2
	}

	world.RenderRuneAt(tile.XPos, tile.YPos, tile.BackgroundGlyph, color, gterm.NoColor)
}
