package game

type Team int

const (
	PlayerTeam Team = iota
	MonsterTeam
)

type Creature struct {
	Identifiable

	Team Team

	X int
	Y int

	HP Health

	Level int

	Name string
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

func (c *Creature) TryMove(newX int, newY int, world World) (MoveResult, interface{}) {

	if world.CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.GetCreatureAtTile(newX, newY); ok {
		if c.Team != defender.Team {
			return MoveIsEnemy, MoveEnemy{Attacker: c, Defender: defender}
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
	Attacker *Creature
	Defender *Creature
}
