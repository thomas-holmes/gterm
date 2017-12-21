package main

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
	// This is all garbage. creature is a mess and I hate it
	creature := NewCreature(level, 100, hp)

	creature.X = xPos
	creature.Y = yPos
	creature.Team = MonsterTeam
	monster := Monster{
		Color:    color,
		Creature: creature,
	}

	return monster
}

// TODO: If a monster is blocking the ideal path our monster should go around
func (monster *Monster) Pursue(turn uint64, world *World) bool {
	if world.CurrentLevel.VisionMap.VisibilityAt(monster.X, monster.Y) == Visible {
		monster.State = Pursuing
	}

	if monster.State != Pursuing {
		return true
	}

	scent := world.CurrentLevel.ScentMap

	// TODO: Maybe short circuit tracking here and just attack the player instead
	// if in ranger?
	candidates := scent.track(turn, monster.X, monster.Y)

	// TODO: Sometimes the monster takes a suboptimal path
	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]
		if len(candidates) > 1 {
			// TODO: Not actually sure if this is invalid but for now I want to know if it happens.
			log.Printf("More than one candidate, %+v", candidates)
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

func (monster *Monster) Update(turn uint64, _ sdl.Event, world *World) bool {
	if monster.Pursue(turn, world) {
		monster.currentEnergy -= 100
		return true
	}
	return false
}

func (monster *Monster) NeedsInput() bool {
	return false
}

func (monster *Monster) Render(world *World) {
	glyph := []rune(strconv.Itoa(monster.Level))[0]
	world.RenderRuneAt(monster.X, monster.Y, glyph, monster.Color, gterm.NoColor)
}
