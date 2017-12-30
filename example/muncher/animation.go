package main

import "log"
import "github.com/thomas-holmes/gterm"

type Animation interface {
	Done() bool
	Start(time uint32)
	Update(delta uint32)
	Render(world *World)
}

// Maybe make a LinearAnimation instead?
type FireBoltAnimation struct {
	x int
	y int

	startTime       uint32
	accumulatedTime uint32
	path            []Position
	step            int

	Speed uint32
}

func (a *FireBoltAnimation) Start(time uint32) {
	a.startTime = time
	a.accumulatedTime = time
}

func NewFireBoltAnimation(startX, startY, endX, endY int, speed uint32) FireBoltAnimation {
	return FireBoltAnimation{
		Speed: speed,
		path:  PlotLine(startX, startY, endX, endY)[1:],
	}
}

func (a *FireBoltAnimation) Done() bool {
	return (a.accumulatedTime - a.startTime) >= (a.Speed * uint32(len(a.path)))
}

func (a *FireBoltAnimation) Update(delta uint32) {
	if a.startTime == 0 {
		log.Panicf("You must first call start on the animation before running it")
	}

	a.accumulatedTime += delta
	elapsed := a.accumulatedTime - a.startTime
	a.step = int(elapsed / a.Speed)
}

// Interp?
func (a *FireBoltAnimation) Render(world *World) {
	if a.Done() {
		return
	}

	pos := a.path[a.step]
	world.RenderRuneAt(pos.X, pos.Y, '*', Red, gterm.NoColor)
}
