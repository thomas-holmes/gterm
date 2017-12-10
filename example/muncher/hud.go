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

	nextFreeRow int
}

func NewHud(player *Player, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player:      player,
		World:       world,
		XPos:        xPos,
		YPos:        yPos,
		nextFreeRow: 0,
	}

	world.MessageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) GetNextRow() int {
	hud.nextFreeRow++
	return hud.nextFreeRow + hud.YPos - 1
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
	}
}
func (hud *HUD) renderPlayerName(world *World) {
	world.Window.PutString(hud.XPos, hud.GetNextRow(), world.Player.Name, Yellow)
}

func (hud *HUD) renderPlayerPosition(world *World) {
	position := fmt.Sprintf("(%v, %v)", hud.Player.X, hud.Player.Y)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), position, Yellow)
}

func (hud *HUD) renderPlayerHealth(world *World) {
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

	world.Window.PutString(hud.XPos, hud.GetNextRow(), hp, hpColor)
}

func (hud *HUD) renderPlayerLevel(world *World) {
	level := fmt.Sprintf("Level: %v", hud.Player.Level)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), level, Yellow)
}

func (hud *HUD) renderTurnCount(world *World) {
	turnCount := fmt.Sprintf("Turn: %v", world.turnCount)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), turnCount, Yellow)
}

func (hud *HUD) renderItemDisplay(world *World) {
	hud.nextFreeRow += 3
	offsetY := hud.GetNextRow()
	offsetX := hud.XPos

	items := make([]*Item, 0)
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			if world.CurrentLevel.VisionMap.VisibilityAt(x, y) == Visible {
				tile := world.GetTile(x, y)
				if tile.Item != nil {
					items = append(items, tile.Item)
				}
			}
		}
	}
	for _, item := range items {
		world.Window.PutRune(hud.XPos, offsetY, item.Symbol, item.Color, gterm.NoColor)
		name := item.Name
		offsetX = hud.XPos
		for {
			offsetX = hud.XPos + 2
			if len(name) == 0 {
				break
			}
			maxLength := (world.Window.Columns - offsetX)
			cut := min(len(name), maxLength)
			printable := name[:cut]
			name = name[cut:]
			world.Window.PutString(offsetX, offsetY, printable, Yellow)
			offsetY = hud.GetNextRow()
			offsetX = hud.XPos + 1
		}
	}
	offsetX = hud.XPos
}

func (hud *HUD) Render(world *World) {
	hud.nextFreeRow = 0
	defer timeMe(time.Now(), "HUD.Render")
	hud.renderPlayerName(world)
	hud.renderPlayerPosition(world)
	hud.renderPlayerHealth(world)
	hud.renderPlayerLevel(world)
	hud.renderTurnCount(world)
	hud.renderItemDisplay(world)
}
