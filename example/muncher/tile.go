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

func TileKindToGlyph(kind TileKind) rune {
	switch kind {
	case Wall:
		return WallGlyph
	case Floor:
		return FloorGlyph
	case UpStair:
		return UpStairGlyph
	case DownStair:
		return DownStairGlyph
	}

	return WallGlyph // Default to a wall for now, I guess.
}

func NewTile(x int, y int) Tile {
	return Tile{
		X:     x,
		Y:     y,
		Color: White,
	}
}

type Tile struct {
	X int
	Y int

	Color sdl.Color

	Creature *Creature
	Item     *Item

	TileGlyph rune
	TileKind
}

func (tile Tile) IsWall() bool {
	return tile.TileKind == Wall
}

func (tile *Tile) Render(world *World, visibility Visibility) {
	if visibility == Unseen {
		return
	}

	if tile.Creature != nil && visibility == Visible {
		tile.Creature.Render(world)
	} else {
		tile.RenderBackground(world, visibility) // bad API, refactor
	}
}

func (tile Tile) RenderBackground(world *World, visibility Visibility) {
	var glyph rune
	var color sdl.Color

	if tile.Item != nil {
		glyph = tile.Item.Symbol
		color = tile.Item.Color
	} else {
		glyph = tile.TileGlyph
		color = tile.Color
	}

	if visibility == Seen {
		color.R /= 2
		color.G /= 2
		color.B /= 2
	}

	world.RenderRuneAt(tile.X, tile.Y, glyph, color, gterm.NoColor)
}
