package main

import (
	"fmt"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type InspectionPop struct {
	World *World

	PopMenu

	InspectX int
	InspectY int
}

func (pop *InspectionPop) Update(input InputEvent) bool {
	newX, newY := pop.InspectX, pop.InspectY
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			pop.done = true
			return true
		case sdl.K_h:
			newX = pop.InspectX - 1
		case sdl.K_j:
			newY = pop.InspectY + 1
		case sdl.K_k:
			newY = pop.InspectY - 1
		case sdl.K_l:
			newX = pop.InspectX + 1
		case sdl.K_b:
			newX, newY = pop.InspectX-1, pop.InspectY+1
		case sdl.K_n:
			newX, newY = pop.InspectX+1, pop.InspectY+1
		case sdl.K_y:
			newX, newY = pop.InspectX-1, pop.InspectY-1
		case sdl.K_u:
			newX, newY = pop.InspectX+1, pop.InspectY-1
		}
	}

	if (newX != pop.InspectX || newY != pop.InspectY) &&
		(newX > 0 && newX < pop.World.CurrentLevel.Columns) &&
		(newY > 0 && newY < pop.World.CurrentLevel.Rows) {
		// Guard against level boundaries
		pop.InspectX = newX
		pop.InspectY = newY

		return true
	}

	return false
}

func (pop *InspectionPop) RenderTileDescription(tile *Tile) {
	if pop.World.CurrentLevel.VisionMap.VisibilityAt(tile.X, tile.Y) == Unseen {
		return
	}
	yOffset := 0
	if c := tile.Creature; c != nil {
		xOffset := 0
		pop.World.Window.PutRune(pop.X+xOffset, pop.Y+yOffset, c.RenderGlyph, c.RenderColor, gterm.NoColor)
		xOffset += 2

		creatureLine1 := fmt.Sprintf("%v (%v/%v)", c.Name, c.HP.Current, c.HP.Max)
		pop.World.Window.PutString(pop.X+xOffset, pop.Y+yOffset, creatureLine1, Yellow)

		yOffset++
	}
	if i := tile.Item; i != nil {
		xOffset := 0
		pop.World.Window.PutRune(pop.X+xOffset, pop.Y+yOffset, i.Symbol, i.Color, gterm.NoColor)

		itemLine1 := fmt.Sprintf("- %v (%v)", i.Name, i.Power)
		yOffset += putWrappedText(pop.World.Window, itemLine1, pop.X, pop.Y+yOffset, 2, 4, pop.W-xOffset, Yellow)
	}
	{
		terrainLine1 := ""
		switch tile.TileKind {
		case Floor:
			terrainLine1 = "Stone floor"
		case Wall:
			terrainLine1 = "A solid rock wall"
		case UpStair:
			terrainLine1 = "Stairs leading up"
		case DownStair:
			terrainLine1 = "Stairs leading down"
		}

		if len(terrainLine1) > 0 {
			pop.World.Window.PutString(pop.X, pop.Y+yOffset, terrainLine1, Yellow)
			yOffset++
		}
	}
	{
		scentStrength := pop.World.CurrentLevel.ScentMap.getScent(pop.InspectX, pop.InspectY)
		scentLine1 := fmt.Sprintf("Scent Strength: (%v)", scentStrength)
		pop.World.Window.PutString(pop.X, pop.Y+yOffset, scentLine1, Yellow)
	}
}

// Maybe should interact with world/tiles than window directly
func (pop *InspectionPop) RenderCursor(window *gterm.Window) {
	white := White
	white.A = 50
	yellow := Yellow
	yellow.A = 200
	positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, pop.InspectX, pop.InspectY)
	for _, pos := range positions {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, white)
	}
	pop.World.RenderRuneAt(pop.InspectX, pop.InspectY, ' ', gterm.NoColor, yellow)
}

func (pop *InspectionPop) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)

	pop.RenderCursor(window)

	tile := pop.World.CurrentLevel.GetTile(pop.InspectX, pop.InspectY)

	pop.RenderTileDescription(tile)
}
