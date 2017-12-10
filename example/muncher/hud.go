package main

import (
	"fmt"
	"time"

	"github.com/thomas-holmes/gterm"
)

type HUD struct {
	World  *World
	Player *Player
	XPos   int
	YPos   int
	Width  int
	Height int
}

func NewHud(player *Player, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player: player,
		World:  world,
		XPos:   xPos,
		YPos:   yPos,
	}

	world.MessageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
	}
}
func (hud HUD) renderPlayerName(world *World) {
	world.Window.PutString(hud.XPos, hud.YPos, world.Player.Name, Yellow)
}

func (hud HUD) renderPlayerPosition(world *World) {
	position := fmt.Sprintf("(%v, %v)", hud.Player.X, hud.Player.Y)
	world.Window.PutString(hud.XPos, hud.YPos+1, position, Yellow)
}

func (hud HUD) renderPlayerHealth(world *World) {
	hpColor := Red

	pct := hud.Player.HealthPercentage()
	switch {
	case pct >= 0.8:
		hpColor = Green
	case pct >= 0.6:
		hpColor = Yellow
	case pct >= 0.4:
		hpColor = Orange
	default:
		hpColor = Red
	}

	hp := fmt.Sprintf("%v/%v", hud.Player.HP.Current, hud.Player.HP.Max)
	if hud.Player.HP.Current == 0 {
		hp += " *DEAD*"
	}

	world.Window.PutString(hud.XPos, hud.YPos+2, hp, hpColor)
}

func (hud HUD) renderPlayerLevel(world *World) {
	level := fmt.Sprintf("Level: %v", hud.Player.Level)
	world.Window.PutString(hud.XPos, hud.YPos+3, level, Yellow)
}

func (hud HUD) renderTurnCount(world *World) {
	turnCount := fmt.Sprintf("Turn: %v", world.turnCount)
	world.Window.PutString(hud.XPos, hud.YPos+4, turnCount, Yellow)
}

func (hud HUD) renderStoodOnItem(world *World) {
	offsetY := hud.YPos + 8
	offsetX := hud.XPos
	tile := world.GetTile(world.Player.X, world.Player.Y)
	if tile.Item != nil {
		item := tile.Item
		world.Window.PutRune(hud.XPos, offsetY, item.Symbol, item.Color, gterm.NoColor)
		offsetX += 2
		name := item.Name
		maxLength := (world.Window.Columns - offsetX)
		for {
			if len(name) == 0 {
				break
			}
			cut := min(len(name), maxLength)
			printable := name[:cut]
			name = name[cut:]
			world.Window.PutString(offsetX, offsetY, printable, Yellow)
			offsetY++
			offsetX = hud.XPos + 1
		}
	}
}

func (hud *HUD) Render(world *World) {
	defer timeMe(time.Now(), "HUD.Render")
	hud.renderPlayerName(world)
	hud.renderPlayerPosition(world)
	hud.renderPlayerHealth(world)
	hud.renderPlayerLevel(world)
	hud.renderTurnCount(world)
	hud.renderStoodOnItem(world)
}
