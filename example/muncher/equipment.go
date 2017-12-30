package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Equipment struct {
	Weapon Item
}

func NewEquipment() Equipment {
	return Equipment{
		Weapon: NoItem,
	}
}

type EquipmentPop struct {
	Player *Creature

	PopMenu

	Messaging
}

func (pop *EquipmentPop) equipItem(index int) {
	// Should probably filter on equippable and change the list, whatever
	// Consider doing this w/o message broadcast. We do have a ref to the player, after all
	if index < len(pop.Player.Inventory.Items) {
		item := pop.Player.Inventory.Items[index]
		log.Printf("Equipping item %+v", item)
		pop.Broadcast(EquipItem, EquipItemMessage{*item})
		pop.done = true
	}
}

func (pop *EquipmentPop) Update(input InputEvent) bool {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_ESCAPE:
			pop.done = true
		case k >= sdl.K_a && k <= sdl.K_z:
			pop.equipItem(int(k - sdl.K_a))
		}

	}

	return false
}

func (pop *EquipmentPop) Render(window *gterm.Window) {
	// TODO: Don't do this
	inventoryPop := InventoryPop{
		Inventory: pop.Player.Inventory,
		PopMenu: PopMenu{
			X: pop.X,
			Y: pop.Y,
			W: pop.W,
			H: pop.H,
		},
	}

	inventoryPop.Render(window)
}
