package game

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/thomas-holmes/gterm"
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

func (monster *Monster) Pursue(turn int64, world World) {
	scent := world.ScentMap
	candidates := scent.track(turn, monster.xPos, monster.yPos)

	log.Printf("Monster %#v found tracking candidates: %v", *monster, candidates)

	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]

		monster.UpdatePosition(choice.XPos, choice.YPos, world)
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

func (monster *Monster) UpdatePosition(xPos int, yPos int, world World) {
	if xPos >= 0 && xPos < world.Columns &&
		yPos >= 0 && yPos < world.Rows {
		if world.IsTileMonster(xPos, yPos) {
			// Nothing
		} else if world.CanStandOnTile(xPos, yPos) {
			oldX := monster.xPos
			oldY := monster.yPos
			monster.xPos = xPos
			monster.yPos = yPos

			monster.Broadcast(MoveEntity, MoveEntityMessage{
				ID:   monster.id,
				OldX: oldX,
				OldY: oldY,
				NewX: monster.xPos,
				NewY: monster.yPos,
			})
		}
	}
}

func (monster *Monster) Render(world *World) {
	glyph := []rune(strconv.Itoa(monster.Level))[0]
	world.RenderRuneAt(monster.xPos, monster.yPos, glyph, monster.Color, gterm.NoColor)
}
