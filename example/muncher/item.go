package main

import (
	"fmt"
	"strconv"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

// This is used for empty hands, maybe?
var NoItem = Item{Power: 1}

var SampleItems = []Item{
	Item{Symbol: ')', Color: Red, Name: "Dagger", Description: "A plain dagger.", Power: 2},
	Item{Symbol: ')', Color: Green, Name: "Rapier", Description: "A slender weapon, perfect for thursting.", Power: 4},
	Item{Symbol: ')', Color: Yellow, Name: "Warhammer", Description: "A greusome warhammer. Few problems can't be solved by crushing", Power: 6},
	Item{Symbol: ')', Color: Blue, Name: "Shortsword named Sting", Description: "This shortsword glows blue when orcs are nearby. It is razor sharp.", Power: 8},
	Item{Symbol: ')', Color: Purple, Name: "Thunderfury, Blessed Blade of the Windseeker", Description: "A two-pronged blade that courses with electrical power. You feel unstoppable while holding this weapon.", Power: 10},
}

type Item struct {
	Name        string
	Description string
	Symbol      rune
	Color       sdl.Color

	Power int
}

type ItemDetails struct {
	*Item

	done bool

	X int
	Y int
	W int
	H int
}

func (pop *ItemDetails) Done() bool {
	return pop.done
}

func (pop *ItemDetails) Update(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			pop.done = true
			return true
		}
	}

	return false
}

func (pop *ItemDetails) renderShortDescription(row int, window *gterm.Window) int {
	offsetX := pop.X + 1
	offsetY := row
	window.PutRune(offsetX, offsetY, pop.Item.Symbol, pop.Item.Color, gterm.NoColor)
	nameStr := fmt.Sprintf(" - %v", pop.Item.Name)
	offsetX += 2

	window.PutString(offsetX, offsetY, nameStr, White)
	return offsetY + 1
}

func (pop *ItemDetails) renderLongDescription(row int, window *gterm.Window) int {
	offsetX := pop.X + 4
	offsetY := row + 1

	description := pop.Item.Description
	for {
		if len(description) == 0 {
			break
		}
		maxLength := pop.W - offsetX + pop.X - 1
		cut := min(len(description), maxLength)
		printable := description[:cut]
		description = description[cut:]
		window.PutString(offsetX, offsetY, printable, White)
		offsetY++
		offsetX = pop.X
	}

	return offsetY + 1
}

func (pop *ItemDetails) renderPower(row int, window *gterm.Window) int {
	offsetX := pop.X + 1
	offsetY := row + 1

	powerString := "Power: "
	window.PutString(offsetX, offsetY, powerString, White)

	offsetX += len(powerString)
	window.PutString(offsetX, offsetY, strconv.Itoa(pop.Item.Power), pop.Item.Color)

	return offsetY + 1
}

func (pop *ItemDetails) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	row := pop.Y + 1
	row = pop.renderShortDescription(row, window)
	row = pop.renderLongDescription(row, window)
	_ = pop.renderPower(row, window)
}
