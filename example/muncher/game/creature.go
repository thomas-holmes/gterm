package game

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

	CurrentEnergy int
	MaxEnergy     int

	HP Health

	Level int

	Name string
}

func (c Creature) NeedsInput() bool {
	log.Printf("Called Creature NeedsInput %+v", c)
	return false
}

func (c *Creature) AddEnergy(energy int) {
	c.CurrentEnergy = min(c.MaxEnergy, c.CurrentEnergy+energy)
}
func (c Creature) Energy() int {
	return c.CurrentEnergy
}
func (c Creature) CanAct() bool {
	return c.CurrentEnergy >= 100
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c *Creature) Damage(damage int) {
	c.HP.Current = max(0, c.HP.Current-damage)
}

func (c *Creature) Fighter() *Creature {
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
