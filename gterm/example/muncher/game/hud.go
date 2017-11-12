package game

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type HUD struct {
	World  *World
	Player *Player
	XPos   int
	YPos   int
	Width  int
	Height int
	Dirty  bool
}

func NewHud(player *Player, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player: player,
		World:  world,
		XPos:   xPos,
		YPos:   yPos,
		Dirty:  true,
	}

	world.MessageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
		hud.Dirty = true
	}
}

var YELLOW = sdl.Color{R: 255, G: 255, B: 0, A: 255}

func (hud *HUD) Render(world *World) {
	if hud.Dirty {
		log.Println("Hud is dirty")
		world.Window.ClearCell(hud.XPos, hud.YPos)
		world.Window.ClearCell(hud.XPos, hud.YPos+1)

		world.Window.AddToCell(hud.XPos, hud.YPos, hud.Player.Name, YELLOW)

		position := fmt.Sprintf("(%v, %v)", hud.Player.XPos, hud.Player.YPos)
		world.Window.AddToCell(hud.XPos, hud.YPos+1, position, YELLOW)

		hud.Dirty = false
	}
}
