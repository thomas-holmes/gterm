package game

type Creature struct {
	Identifiable

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
	if world.IsTileMonster(newX, newY) {
		monster := world.GetMonsterAtTile(newX, newY)
		return MoveIsEnemy, MoveEnemy{Attacker: c, Defender: &monster.Creature}
	} else if world.CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
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
