package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Inventory struct {
	Items []*Item
}

type InventoryPop struct {
	Inventory

	done bool

	X int
	Y int
	W int
	H int

	Messaging
}

func (pop *InventoryPop) Done() bool {
	return pop.done
}

func (pop *InventoryPop) tryShowItem(index int) {
	if index < len(pop.Items) {
		menu := ItemDetails{X: 2, Y: 2, W: 50, H: 26, Item: pop.Items[index]}
		log.Printf("Trying to broadcast %+v", menu)
		pop.Broadcast(ShowMenu, ShowMenuMessage{Menu: &menu})
	}
}

func (pop *InventoryPop) Update(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_ESCAPE:
			pop.done = true
			return true
		case k >= sdl.K_a && k <= sdl.K_z:
			pop.tryShowItem(int(k - sdl.K_a))
		}

	}

	return false
}

func (pop *InventoryPop) renderItem(index int, row int, window *gterm.Window) int {
	offsetY := row
	offsetX := pop.X + 1

	item := pop.Items[index]

	selectionStr := fmt.Sprintf("%v - ", string('a'+index))

	window.PutString(offsetX, offsetY, selectionStr, White)

	name := item.Name

	offsetY += putWrappedText(window, name, offsetX, offsetY, len(selectionStr), 2, pop.W-offsetX+pop.X-1, White)
	return offsetY
}

func (pop *InventoryPop) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Printf("(%v,%v) (%v,%v)", pop.X, pop.Y, pop.W, pop.H)
		log.Println("Failed to render inventory", err)
	}

	nextRow := pop.Y + 1
	for i := 0; i < len(pop.Items); i++ {
		nextRow = pop.renderItem(i, nextRow, window)
	}

	window.PutString(pop.X, pop.Y, strings.Repeat("%", pop.W), White)
	for y := pop.Y + 1; y < pop.Y+pop.H-1; y++ {
		window.PutRune(pop.X, y, '%', White, gterm.NoColor)
		window.PutRune(pop.X+pop.W-1, y, '%', White, gterm.NoColor)
	}
	window.PutString(pop.X, pop.Y+pop.H-1, strings.Repeat("%", pop.W), White)
}
