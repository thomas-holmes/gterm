package game

import (
	"time"
)

type ScentMap struct {
	columns int
	rows    int
	scent   []int64
}

func (scentMap ScentMap) getScent(xPos int, yPos int) int64 {
	return scentMap.scent[yPos*scentMap.columns+xPos]
}

func (scentMap ScentMap) dirty(xPos int, yPos int, turn int64, distance int) {
	scentMap.scent[yPos*scentMap.columns+xPos] = turn*32 - int64(distance)
}

func (scentMap ScentMap) track(turn int64, xPos int, yPos int) []Position {
	minX := max(0, xPos-1)
	maxX := min(scentMap.columns, xPos+2)
	minY := max(0, yPos-1)
	maxY := min(scentMap.columns, yPos+2)

	candidates := make([]Position, 0, 8)
	strongest := turn * 31
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
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

func (scentMap ScentMap) UpdateScents(turn int64, vision VisionMap, player Player) {
	defer timeMe(time.Now(), "ScentMap.UpdateScents")
	for y := 0; y < scentMap.rows; y++ {
		for x := 0; x < scentMap.columns; x++ {
			vision := vision.VisibilityAt(x, y)
			if vision == Visible {
				scentMap.dirty(x, y, turn, distance(player.xPos, player.yPos, x, y))
			}
		}
	}
}

func newScentMap(columns int, rows int) ScentMap {
	return ScentMap{
		columns: columns,
		rows:    rows,
		scent:   make([]int64, columns*rows, columns*rows),
	}
}
