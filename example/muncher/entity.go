package main

import "github.com/veandco/go-sdl2/sdl"

type Identifiable struct {
	ID int
}

func (e *Identifiable) SetIdentity(id int) {
	e.ID = id
}

func (e Identifiable) Identity() int {
	return e.ID
}

// TODO: Too much stuff on this. Need to figure out how to do entities properly
type Entity interface {
	Identity() int
	SetIdentity(int)
	NeedsInput() bool
	CanAct() bool
	Update(turn int64, event sdl.Event, world *World) bool
	Fighter() *Creature
}
