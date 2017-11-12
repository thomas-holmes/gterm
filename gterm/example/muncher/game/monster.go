package game

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type Monster struct {
	XPos  int
	YPos  int
	HP    Health
	Glyph string
	Color sdl.Color
	Dirty bool
}

func (monster *Monster) UpdatePosition(xPos int, yPos int) {
	monster.XPos = xPos
	monster.YPos = yPos
	monster.Dirty = true
}

func (monster *Monster) Render(world *World) {
	if monster.Dirty {
		if err := world.Window.AddToCell(monster.XPos, monster.YPos, monster.Glyph, monster.Color); err != nil {
			log.Println("Failed to render monster", monster)
		}
		monster.Dirty = false
	}
}
