package game

type Entity interface {
	ID() int
	SetID(int)
}