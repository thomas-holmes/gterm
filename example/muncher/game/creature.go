package game

type Creature struct {
	Identifiable

	X int
	Y int

	HP Health

	Name string
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}
