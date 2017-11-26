package game

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Monster struct {
	Level int
	Glyph rune
	Color sdl.Color

	Creature

	Messaging
}

func NewMonster(xPos int, yPos int, level int, color sdl.Color, hp int) Monster {
	monster := Monster{
		Color: color,
		Creature: Creature{
			X: xPos,
			Y: yPos,
			HP: Health{
				Current: hp,
				Max:     hp,
			},
		},
		Level: level,
	}

	return monster
}

func (monster *Monster) Pursue(turn int64, world World) {
	scent := world.ScentMap
	candidates := scent.track(turn, monster.X, monster.Y)

	log.Printf("Monster %#v found tracking candidates: %v", *monster, candidates)

	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]

		monster.UpdatePosition(choice.XPos, choice.YPos, world)
	}
}

func (monster *Monster) Kill() {
	monster.Broadcast(KillMonster, KillMonsterMessage{ID: monster.ID})
}

func (monster *Monster) UpdatePosition(xPos int, yPos int, world World) {
	if xPos >= 0 && xPos < world.Columns &&
		yPos >= 0 && yPos < world.Rows {
		if world.IsTileMonster(xPos, yPos) {
			// Nothing
		} else if world.CanStandOnTile(xPos, yPos) {
			oldX := monster.X
			oldY := monster.Y
			monster.X = xPos
			monster.Y = yPos

			monster.Broadcast(MoveEntity, MoveEntityMessage{
				ID:   monster.ID,
				OldX: oldX,
				OldY: oldY,
				NewX: monster.X,
				NewY: monster.Y,
			})
		}
	}
}

func (monster *Monster) Render(world *World) {
	glyph := []rune(strconv.Itoa(monster.Level))[0]
	world.RenderRuneAt(monster.X, monster.Y, glyph, monster.Color, gterm.NoColor)
}
