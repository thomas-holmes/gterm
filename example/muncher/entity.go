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

type Entity interface {
	Identity() int
	SetIdentity(int)
}

func foo(entity Entity) {

}
