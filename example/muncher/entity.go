package main

type Identifiable struct {
	ID int
}

func (e *Identifiable) SetIdentity(id int) {
	e.ID = id
}

func (e Identifiable) Identity() int {
	return e.ID
}

// TODO: Too much stuff on this. Need to figure out how to do entities properly
type Entity interface {
	Identity() int
	SetIdentity(int)
}
