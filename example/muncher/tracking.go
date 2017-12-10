package main

import (
	"log"
	"time"
)

type ScentMap struct {
	columns int
	rows    int
	scent   []float64
}

func (scentMap ScentMap) getScent(xPos int, yPos int) float64 {
	return scentMap.scent[yPos*scentMap.columns+xPos]
}

func (scentMap ScentMap) dirty(xPos int, yPos int, turn uint64, distance float64) {
	scentMap.scent[yPos*scentMap.columns+xPos] = float64(turn*32) - distance
}

func (scentMap ScentMap) track(turn uint64, xPos int, yPos int) []Position {
	minX := max(0, xPos-1)
	maxX := min(scentMap.columns, xPos+2)
	minY := max(0, yPos-1)
	maxY := min(scentMap.columns, yPos+2)

	candidates := make([]Position, 0, 8)
	log.Printf("Turn %v", turn)
	recent := float64((turn - (minu64(turn, 50))) * 32.0)
	strongest := recent
	log.Printf("Strongest is %v", strongest)
	log.Printf("Scanning (%v,%x) to (%v,%v)", minX, minY, maxX, maxY)
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			strength := scentMap.getScent(x, y)
			log.Printf("Found Strength at (%v,%v) as %v", x, y, strength)

			if strength > strongest {
				candidates = candidates[:0]
				strongest = strength
			}

			if strength == strongest {
				log.Printf("Found a candidate at (%v,%v)", x, y)
				candidates = append(candidates,
					Position{XPos: x, YPos: y},
				)
			}
		}
	}
	return candidates
}

func (scentMap ScentMap) UpdateScents(world *World) {
	defer timeMe(time.Now(), "ScentMap.UpdateScents")

	vision := world.CurrentLevel.VisionMap
	player := world.Player
	for y := 0; y < scentMap.rows; y++ {
		for x := 0; x < scentMap.columns; x++ {
			vision := vision.VisibilityAt(x, y)
			if vision == Visible && !world.GetTile(x, y).IsWall() {
				scentMap.dirty(x, y, world.turnCount, distance(player.X, player.Y, x, y))
			}
		}
	}
}

func NewScentMap(columns int, rows int) ScentMap {
	return ScentMap{
		columns: columns,
		rows:    rows,
		scent:   make([]float64, columns*rows, columns*rows),
	}
}
