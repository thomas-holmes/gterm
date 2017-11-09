package game

import (
	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

/*
Right now Player doesn't have a separate idea of screen and world position.
It is also just a dumb free floating struct that I'll have to track in the
main loop. I'll need to setup a system to register actors in the "scene" in
addition to tracking both the player's world position and its render position.

Once I separate the viewport from the world the player position likely won't
move very often as roguelikes tend to keep the player centered.
*/

// Player pepresents the player
type Player struct {
	Column int
	Row    int
	Glyph  string
	Color  sdl.Color
}

// MoveTo updates the player to a new position
func (player *Player) MoveTo(col int, row int) {
	player.Column = col
	player.Row = row
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.Color = color
}

// Render places the player on the gterm Window
func (player Player) Render(window *gterm.Window) {
	window.AddToCell(player.Column, player.Row, player.Glyph, player.Color)
}
