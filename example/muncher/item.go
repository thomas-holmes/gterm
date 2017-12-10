package main

import "github.com/veandco/go-sdl2/sdl"

var SampleItems = []Item{
	Item{Symbol: ')', Color: Red, Name: "Dagger", Power: 2},
	Item{Symbol: ')', Color: Green, Name: "Rapier", Power: 4},
	Item{Symbol: ')', Color: Yellow, Name: "Warhammer", Power: 6},
	Item{Symbol: ')', Color: Blue, Name: "Shortsword named Sting", Power: 8},
	Item{Symbol: ')', Color: Purple, Name: "Thunderfury, Blessed Blade of the Windseeker", Power: 10},
}

type Item struct {
	Name   string
	Symbol rune
	Color  sdl.Color

	Power int
}
