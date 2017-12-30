package main

import "log"
import "github.com/thomas-holmes/gterm"
import "github.com/veandco/go-sdl2/sdl"

type Animation interface {
	Done() bool
	Start(time uint32)
	Update(delta uint32)
	Render(world *World)
}

// Maybe make a LinearAnimation instead?
type LinearSpellAnimation struct {
	x int
	y int

	startTime       uint32
	accumulatedTime uint32
	path            []Position
	step            int

	Delay uint32
	Speed uint32

	Color sdl.Color
	Glyph rune
}

func (a *LinearSpellAnimation) Start(time uint32) {
	a.startTime = time
	a.accumulatedTime = time
}

func NewLinearSpellAnimation(startX, startY, endX, endY int, speed uint32, delay uint32, glyph rune, color sdl.Color) LinearSpellAnimation {
	return LinearSpellAnimation{
		Speed: speed,
		Color: color,
		Glyph: glyph,
		Delay: delay,
		path:  PlotLine(startX, startY, endX, endY)[1:],
	}
}

func (a *LinearSpellAnimation) Done() bool {
	// Doing this math with int64 due to uint32 underflows making this misbehave
	return (int64(a.accumulatedTime) - int64((a.startTime + a.Delay))) >= int64((a.Speed * uint32(len(a.path))))
}

func (a *LinearSpellAnimation) Update(delta uint32) {
	if a.startTime == 0 {
		log.Panicf("You must first call start on the animation before running it")
	}

	a.accumulatedTime += delta
	elapsed := a.accumulatedTime - (a.startTime + a.Delay) // beware underflow
	a.step = int((elapsed) / a.Speed)
}

func (a *LinearSpellAnimation) Ready() bool {
	return a.accumulatedTime > a.startTime+a.Delay
}

// Interp?
func (a *LinearSpellAnimation) Render(world *World) {
	if a.Done() || !a.Ready() {
		return
	}

	pos := a.path[a.step]
	world.RenderRuneAt(pos.X, pos.Y, a.Glyph, a.Color, gterm.NoColor)
}
