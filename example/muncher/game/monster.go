package game

import (
	"log"
	"math/rand"
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

func (monster *Monster) Pursue(turn int64, scent ScentMap) {
	candidates := scent.track(turn, monster.xPos, monster.yPos)

	log.Printf("Monster %#v found tracking candidates: %v", *monster, candidates)

	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]

		oldX := monster.xPos
		oldY := monster.yPos
		monster.xPos = choice.XPos
		monster.yPos = choice.YPos

		monster.Broadcast(MoveEntity, MoveEntityMessage{
			ID:   monster.id,
			OldX: oldX,
			OldY: oldY,
			NewX: monster.xPos,
			NewY: monster.yPos,
		})
	}
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
