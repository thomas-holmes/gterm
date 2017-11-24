package game

import (
	"time"
)

type ScentMap struct {
	columns int
	rows    int
	scent   []int
}

func (scentMap ScentMap) getScent(xPos int, yPos int) int {
	return scentMap.scent[yPos*scentMap.columns+xPos]
}

const (
	MinScent int = 0
	MaxScent int = 8
)

func (scentMap ScentMap) dirty(xPos int, yPos int) {
	scentMap.scent[yPos*scentMap.columns+xPos] = MaxScent
}

func (scentMap ScentMap) freshen(xPos int, yPos int) {
	current := scentMap.getScent(xPos, yPos)
	newScent := max(MinScent, current-1)
	scentMap.scent[yPos*scentMap.columns+xPos] = newScent
}

func (scentMap ScentMap) track(xPos int, yPos int) []Position {
	minX := max(0, xPos-1)
	maxX := min(scentMap.columns, xPos+2)
	minY := max(0, yPos-1)
	maxY := min(scentMap.columns, yPos+2)

	candidates := make([]Position, 0, 8)
	strongest := MinScent + 1 // Skip 0s
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; y++ {
			strength := scentMap.getScent(x, y)

			if strength > strongest {
				candidates = candidates[:0]
				strongest = strength
			}

			if strength == strongest {
				candidates = append(candidates,
					Position{XPos: x, YPos: y},
				)
			}
		}
	}
	return candidates
}

func (scentMap ScentMap) UpdateScents(vision VisionMap) {
	defer timeMe(time.Now(), "ScentMap.UpdateScents")
	for y := 0; y < scentMap.rows; y++ {
		for x := 0; x < scentMap.columns; x++ {
			vision := vision.VisibilityAt(x, y)
			if vision == Visible {
				scentMap.dirty(x, y)
			} else {
				scentMap.freshen(x, y)
			}
		}
	}
}

func newScentMap(columns int, rows int) ScentMap {
	return ScentMap{
		columns: columns,
		rows:    rows,
		scent:   make([]int, columns*rows, columns*rows),
	}
}
