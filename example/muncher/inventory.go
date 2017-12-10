package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Menu interface {
	Update(sdl.Event) bool
	Render()
	Done() bool
}

type Inventory struct {
	Items []*Item
}

type InventoryPop struct {
	Inventory

	window *gterm.Window

	done bool

	X int
	Y int
	W int
	H int
}

func (pop *InventoryPop) Done() bool {
	return pop.done
}

func (pop *InventoryPop) Update(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_a:
			pop.done = true
			return true
		}
	}

	return false
}

func (pop *InventoryPop) renderItem(index int, row int) int {
	offsetY := row
	offsetX := pop.X + 1

	item := pop.Items[index]

	selectionStr := fmt.Sprintf("%v - ", string('a'+index))

	pop.window.PutString(offsetX, offsetY, selectionStr, White)

	name := item.Name
	offsetX += len(selectionStr)
	for {
		if len(name) == 0 {
			break
		}
		maxLength := pop.W - offsetX + pop.X - 1
		cut := min(len(name), maxLength)
		printable := name[:cut]
		name = name[cut:]
		pop.window.PutString(offsetX, offsetY, printable, White)
		offsetY++
		offsetX = pop.X + len(selectionStr) - 2
	}
	return offsetY
}

func (pop *InventoryPop) Render() {
	if err := pop.window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Printf("(%v,%v) (%v,%v)", pop.X, pop.Y, pop.W, pop.H)
		log.Printf("Failed to render inventory", err)
	}

	nextRow := pop.Y + 1
	for i := 0; i < len(pop.Items); i++ {
		nextRow = pop.renderItem(i, nextRow)
	}

	pop.window.PutString(pop.X, pop.Y, strings.Repeat("%", pop.W), White)
	for y := pop.Y + 1; y < pop.Y+pop.H-1; y++ {
		pop.window.PutRune(pop.X, y, '%', White, gterm.NoColor)
		pop.window.PutRune(pop.X+pop.W-1, y, '%', White, gterm.NoColor)
	}
	pop.window.PutString(pop.X, pop.Y+pop.H-1, strings.Repeat("%", pop.W), White)
}
