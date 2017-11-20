package game

import (
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

func (tile *Tile) Render(world *World) {
	pos := Position{XPos: tile.XPos, YPos: tile.YPos}
	items := world.renderItems[pos]
	if len(items) > 0 {
		for _, item := range items {
			item.Render(world)
		}
	} else {
		tile.RenderBackground(world) // bad API, refactor
	}
}

func cellDistance(xPos1 int, yPos1 int, xPos2 int, yPos2 int) (x int, y int) {
	return abs(xPos2 - xPos1), abs(yPos2 - yPos1)
}

func (tile Tile) Visible(xPos int, yPos int, world World) bool {
	x, y := cellDistance(xPos, yPos, world.Player.XPos(), world.Player.YPos())

	if x > 4 || y > 4 {
		return false
	}

	return true
}

func (tile Tile) RenderBackground(world *World) {
	world.RenderAt(tile.XPos, tile.YPos, tile.BackgroundGlyph, tile.BackgroundColor)
}
