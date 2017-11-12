package game

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type Monster struct {
	id    int
	xPos  int
	yPos  int
	HP    Health
	Glyph string
	Color sdl.Color
}

func NewMonster(xPos int, yPos int, glyph string, color sdl.Color, hp int) Monster {
	monster := Monster{
		xPos:  xPos,
		yPos:  yPos,
		Glyph: glyph,
		Color: color,
		HP: Health{
			Current: hp,
			Max:     hp,
		},
	}

	return monster
}

func (monster Monster) XPos() int {
	return monster.xPos
}

func (monster Monster) YPos() int {
	return monster.yPos
}

func (monster Monster) ID() int {
	return monster.id
}

func (monster *Monster) UpdatePosition(xPos int, yPos int) {
	monster.xPos = xPos
	monster.yPos = yPos
}

func (monster *Monster) Render(world *World) {
	if err := world.Window.AddToCell(monster.xPos, monster.yPos, monster.Glyph, monster.Color); err != nil {
		log.Println("Failed to render monster", monster)
	}
}
