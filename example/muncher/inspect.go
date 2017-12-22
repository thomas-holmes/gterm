package main

import (
	"fmt"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type InspectionPop struct {
	World *World

	done bool

	X int
	Y int
	W int
	H int

	InspectX int
	InspectY int
}

func (pop *InspectionPop) Done() bool {
	return pop.done
}

func (pop *InspectionPop) Update(event sdl.Event) bool {
	newX, newY := pop.InspectX, pop.InspectY
	switch e := event.(type) {
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

	if newX != pop.InspectX || newY != pop.InspectY {
		// Guard against level boundaries
		pop.InspectX = newX
		pop.InspectY = newY

		return true
	}

	return false
}

func (pop *InspectionPop) RenderTileDescription(tile *Tile) {
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
		xOffset += 2
		itemLine1 := fmt.Sprintf("- %v (%v)", i.Name, i.Power)
		pop.World.Window.PutString(pop.X+xOffset, pop.Y+yOffset, itemLine1, Yellow)
		yOffset++
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
	pop.World.RenderRuneAt(pop.InspectX, pop.InspectY, 'X', Yellow, White)
}

func (pop *InspectionPop) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)

	pop.RenderCursor(window)

	tile := pop.World.GetTile(pop.InspectX, pop.InspectY)

	pop.RenderTileDescription(tile)
}
