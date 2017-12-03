package main

type Energy struct {
	currentEnergy int
	maxEnergy     int
}

type Energized interface {
	CurrentEnergy() int
	AddEnergy(int)
}

func (energy Energy) CurrentEnergy() int {
	return energy.currentEnergy
}

func (energy *Energy) AddEnergy(e int) {
	newEnergy := min(energy.maxEnergy, energy.currentEnergy+e)
	energy.currentEnergy = newEnergy
}
