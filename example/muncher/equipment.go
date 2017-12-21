package main

type Equipment struct {
	Weapon Item
}

func NewEquipment() Equipment {
	return Equipment{
		Weapon: NoItem,
	}
}
