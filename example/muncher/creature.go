package main

import "log"

type Team int

const (
	PlayerTeam Team = iota
	MonsterTeam
)

type Health struct {
	Current int
	Max     int
}

type Creature struct {
	Identifiable

	Team Team

	X int
	Y int

	Energy

	Equipment

	HP Health

	Level int

	Name string
}

func (c *Creature) NeedsInput() bool {
	log.Printf("Called Creature NeedsInput %+v", c)
	return false
}

func (c Creature) CanAct() bool {
	return c.currentEnergy >= 100
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c *Creature) Damage(damage int) {
	log.Printf("%v is Taking damage of %v", *c, damage)
	c.HP.Current = max(0, c.HP.Current-damage)
}

// Name collisions, man.
func (c *Creature) Combatant() *Creature {
	return c
}

func (c *Creature) TryMove(newX int, newY int, world *World) (MoveResult, interface{}) {

	if world.CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.GetCreatureAtTile(newX, newY); ok {
		if c.Team != defender.Team {
			a, aOk := world.GetEntity(c.ID)
			d, dOk := world.GetEntity(defender.ID)
			if aOk && dOk {
				return MoveIsEnemy, MoveEnemy{Attacker: a, Defender: d}
			}
		}
	}

	return MoveIsInvalid, nil
}

func NewCreature(level int, maxEnergy int, maxHP int) Creature {
	return Creature{
		Level: level,
		Energy: Energy{
			currentEnergy: maxEnergy,
			maxEnergy:     maxEnergy,
		},
		HP: Health{Current: 5, Max: 5},
		Equipment: NewEquipment(),
	}
}

type MoveResult int

const (
	MoveIsInvalid MoveResult = iota
	MoveIsSuccess
	MoveIsEnemy
)

type MoveEnemy struct {
	Attacker Entity
	Defender Entity
}
