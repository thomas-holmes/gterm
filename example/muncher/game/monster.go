package game

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type Monster struct {
	id    int
	xPos  int
	yPos  int
	HP    Health
	Level int
	Glyph rune
	Color sdl.Color
	Messaging
}

func NewMonster(xPos int, yPos int, level int, color sdl.Color, hp int) Monster {
	monster := Monster{
		xPos:  xPos,
		yPos:  yPos,
		Color: color,
		HP: Health{
			Current: hp,
			Max:     hp,
		},
		Level: level,
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

func (monster *Monster) SetID(id int) {
	monster.id = id
}

func (monster *Monster) Kill() {
	monster.Broadcast(KillMonster, KillMonsterMessage{ID: monster.ID()})
}

func (monster *Monster) UpdatePosition(xPos int, yPos int) {
	monster.xPos = xPos
	monster.yPos = yPos
}

func (monster *Monster) Render(world *World) {
	glyph := []rune(strconv.Itoa(monster.Level))[0]
	world.RenderRuneAt(monster.xPos, monster.yPos, glyph, monster.Color)
}
