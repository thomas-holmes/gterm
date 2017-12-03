package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type TileKind int

const (
	Wall TileKind = iota
	Floor
	UpStair
	DownStair
)

const (
	WallGlyph      = '#'
	FloorGlyph     = '.'
	UpStairGlyph   = '<'
	DownStairGlyph = '>'
)

func NewTile(x int, y int) Tile {
	return Tile{
		X:     x,
		Y:     y,
		Color: sdl.Color{R: 192, G: 192, B: 192, A: 255},
	}
}

type Tile struct {
	X int
	Y int

	Color sdl.Color

	TileGlyph rune
	TileKind

	WasVisible bool
}

func (tile Tile) IsWall() bool {
	return tile.TileKind == Wall
}

func (tile *Tile) Render(world *World, visibility Visibility) {
	if visibility == Unseen {
		return
	}

	pos := Position{XPos: tile.X, YPos: tile.Y}
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
	color := tile.Color

	if visibility == Seen {
		color.R /= 2
		color.G /= 2
		color.B /= 2
	}

	world.RenderRuneAt(tile.X, tile.Y, tile.TileGlyph, color, gterm.NoColor)
}
