package game

type Renderable interface {
	XPos() int
	YPos() int
	Render(world *World)
}
