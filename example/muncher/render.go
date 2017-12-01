package main

type Renderable interface {
	Identity() int
	XPos() int
	YPos() int
	Render(world *World)
}
