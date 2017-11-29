package game

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type MonsterBehavior int

const (
	Idle MonsterBehavior = iota
	Pursuing
)

type Monster struct {
	Glyph rune
	Color sdl.Color

	State MonsterBehavior

	Creature

	Messaging
}

func NewMonster(xPos int, yPos int, level int, color sdl.Color, hp int) Monster {
	monster := Monster{
		Color: color,
		Creature: Creature{
			MaxEnergy: 100,
			Team:      MonsterTeam,
			Level:     level,
			X:         xPos,
			Y:         yPos,
			HP: Health{
				Current: hp,
				Max:     hp,
			},
		},
	}

	return monster
}

func (monster *Monster) Pursue(turn int64, world *World) bool {
	if world.VisionMap.VisibilityAt(monster.X, monster.Y) == Visible {
		monster.State = Pursuing
	}

	if monster.State != Pursuing {
		return true
	}

	scent := world.ScentMap
	candidates := scent.track(turn, monster.X, monster.Y)

	log.Printf("Monster %#v found tracking candidates: %v", *monster, candidates)

	// TODO: Sometimes the monster takes a suboptimal path
	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]
		if len(candidates) > 1 {
			log.Panicf("More than one candidate, %+v", candidates)
		}

		result, data := monster.TryMove(choice.XPos, choice.YPos, world)
		log.Printf("Tried to move %#v, got result: %v, data %#v", monster, result, data)
		switch result {
		case MoveIsInvalid:
			log.Panicf("Monsters aren't allowed to yield their turn")
			return false
		case MoveIsSuccess:
			oldX := monster.X
			oldY := monster.Y
			monster.X = choice.XPos
			monster.Y = choice.YPos
			monster.Broadcast(MoveEntity, MoveEntityMessage{
				ID:   monster.ID,
				OldX: oldX,
				OldY: oldY,
				NewX: choice.XPos,
				NewY: choice.YPos,
			})
		case MoveIsEnemy:
			if data, ok := data.(MoveEnemy); ok {
				monster.Broadcast(AttackEntity, AttackEntityMesasge{
					Attacker: data.Attacker,
					Defender: data.Defender,
				})
			}
		}
		return true
	}
	return false
}

func (monster *Monster) Update(turn int64, _ sdl.Event, world *World) bool {
	log.Printf("Updating Monster %+v", *monster)
	if monster.Pursue(turn, world) {
		monster.CurrentEnergy -= 100
		return true
	}
	return false
}

func (monster Monster) NeedsInput() bool {
	log.Printf("Called monster NeedsInput %+v", monster)
	return false
}

func (monster *Monster) Render(world *World) {
	glyph := []rune(strconv.Itoa(monster.Level))[0]
	world.RenderRuneAt(monster.X, monster.Y, glyph, monster.Color, gterm.NoColor)
}
