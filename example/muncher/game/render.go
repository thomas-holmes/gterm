package game

type Renderable interface {
	ID() int
	XPos() int
	YPos() int
	Render(world *World)
}
