package game

import (
	"log"

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

type PositionComponent struct {
	yPos int
	xPos int
}

type Positionable interface {
	YPos() int
	XPos() int
	UpdatePosition(int, int)
}

func (player *Player) XPos() int {
	return player.xPos
}

func (player *Player) YPos() int {
	return player.yPos
}

func (player *Player) UpdatePosition(xPos int, yPos int) {
	if err := player.window.ClearCell(player.xPos, player.yPos); err != nil {
		log.Println("Failed to remove player", err)
	}
	player.xPos = xPos
	player.yPos = yPos
	player.dirty = true
}

type RenderComponent struct {
	window      *gterm.Window
	renderGlyph string
	renderColor sdl.Color
	dirty       bool
}

type Renderable interface {
	ShouldRender() bool
	Render(window *gterm.Window)
}

func (player *Player) ShouldRender() bool {
	return player.dirty
}

func (player *Player) Render(window *gterm.Window) {
	window.AddToCell(player.xPos, player.yPos, player.renderGlyph, player.renderColor)
	player.dirty = false
}

type Inputtable interface {
	HandleInput(event sdl.Event)
}

// HandleInput updates player position based on user input
func (player *Player) HandleInput(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_h:
			player.UpdatePosition(player.XPos()-1, player.YPos())
		case sdl.K_j:
			player.UpdatePosition(player.XPos(), player.YPos()+1)
		case sdl.K_k:
			player.UpdatePosition(player.XPos(), player.YPos()-1)
		case sdl.K_l:
			player.UpdatePosition(player.XPos()+1, player.YPos())
		case sdl.K_b:
			player.UpdatePosition(player.XPos()-1, player.YPos()+1)
		case sdl.K_n:
			player.UpdatePosition(player.XPos()+1, player.YPos()+1)
		case sdl.K_y:
			player.UpdatePosition(player.XPos()-1, player.YPos()-1)
		case sdl.K_u:
			player.UpdatePosition(player.XPos()+1, player.YPos()-1)
		}
	}
}

// Player pepresents the player
type Player struct {
	RenderComponent
	PositionComponent
}

func NewPlayer(window *gterm.Window, xPos int, yPos int) Player {
	player := Player{
		RenderComponent: RenderComponent{
			window:      window,
			renderGlyph: "@",
			renderColor: sdl.Color{R: 255, G: 0, B: 0, A: 0},
			dirty:       true,
		},
		PositionComponent: PositionComponent{
			xPos: xPos,
			yPos: yPos,
		},
	}

	return player
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.renderColor = color
}
