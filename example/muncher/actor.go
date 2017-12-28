package main

type Actor interface {
	CanAct() bool
	Update(turn uint64, input InputEvent, world *World) bool
	NeedsInput() bool
}
