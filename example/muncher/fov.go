package main

import (
	"log"
	"time"
)

type Visibility int

const (
	Unseen Visibility = iota
	Visible
	Seen
)

type VisionMap struct {
	Columns int
	Rows    int
	Current int64
	Map     []int64
}

func (vision VisionMap) VisibilityAt(x int, y int) Visibility {
	switch vision.lastSeenAt(x, y) {
	case vision.Current:
		return Visible
	case 0:
		return Unseen
	default:
		return Seen
	}
}

func (vision VisionMap) lastSeenAt(x int, y int) int64 {
	return vision.Map[y*vision.Columns+x]
}

func (vision *VisionMap) UpdateVision(viewDistance int, world *World) {
	defer timeMe(time.Now(), "VisionMap.UpdateVision")
	playerX := world.Player.X
	playerY := world.Player.Y

	// Go beyond the min/max so we update cells we are moving away from
	minX := max(playerX-viewDistance, 0)
	maxX := min(playerX+viewDistance, vision.Columns)

	minY := max(playerY-viewDistance, 0)
	maxY := min(playerY+viewDistance, vision.Rows)

	vision.Current++
	current := vision.Current

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			previousVision := vision.lastSeenAt(x, y)

			if previousVision == current {
				continue
			}

			newVision := previousVision

			if vision.CheckVision(playerX, playerY, x, y, world) {
				newVision = current
			}

			vision.Map[y*vision.Columns+x] = newVision
		}
	}
}

func (vision *VisionMap) CheckVision(playerX int, playerY int, candidateX int, candidateY int, world *World) bool {
	cells := PlotLine(playerX, playerY, candidateX, candidateY)

	foundWall := false

	for _, cell := range cells {
		if foundWall {
			return false
		}
		// Either a wall or on the way to a wall, so we can see it.
		vision.Map[cell.Y*vision.Columns+cell.X] = vision.Current

		tile := world.CurrentLevel.GetTile(cell.X, cell.Y)

		if tile.IsWall() {
			foundWall = true
		}
	}
	return true
}

func NewVisionMap(columns int, rows int) *VisionMap {
	return &VisionMap{
		Columns: columns,
		Rows:    rows,
		Map:     make([]int64, columns*rows),
	}
}

func PlotLine(x0 int, y0 int, x1 int, y1 int) []Position {
	octant := computeOctant(x0, y0, x1, y1)
	x0, y0 = toOctantZero(octant, x0, y0)
	x1, y1 = toOctantZero(octant, x1, y1)

	dx := x1 - x0
	dy := y1 - y0
	d := 2*dy - dx
	y := y0

	coordinates := make([]Position, 0)

	for x := x0; x <= x1; x++ {
		correctedX, correctedY := fromOctantZero(octant, x, y)
		coordinates = append(coordinates, Position{X: correctedX, Y: correctedY})

		if d > 0 {
			y++
			d -= (2 * dx)
		}
		d += (2 * dy)
	}

	return coordinates
}

func computeOctant(x0 int, y0 int, x1 int, y1 int) int {
	if x1 > x0 {
		if y1 > y0 {
			if (y1 - y0) < (x1 - x0) {
				return 0
			}
			return 1
		}
		if (y0 - y1) < (x1 - x0) {
			return 7
		}
		return 6
	}
	if y1 > y0 {
		if (y1 - y0) < (x0 - x1) {
			return 3
		}
		return 2
	}
	if (y0 - y1) < (x0 - x1) {
		return 4
	}
	return 5
}

func toOctantZero(octant int, x int, y int) (int, int) {
	switch octant {
	case 0:
		return x, y
	case 1:
		return y, x
	case 2:
		return y, -x
	case 3:
		return -x, y
	case 4:
		return -x, -y
	case 5:
		return -y, -x
	case 6:
		return -y, x
	case 7:
		return x, -y
	}
	log.Fatalf("Received invalid octant, %v for (%v,%v)", octant, x, y)
	return x, y // Unreachable
}

func fromOctantZero(octant int, x int, y int) (int, int) {
	switch octant {
	case 0:
		return x, y
	case 1:
		return y, x
	case 2:
		return -y, x
	case 3:
		return -x, y
	case 4:
		return -x, -y
	case 5:
		return -y, -x
	case 6:
		return y, -x
	case 7:
		return x, -y
	}

	log.Fatalf("Received invalid octant, %v for (%v,%v)", octant, x, y)
	return x, y // Unreachable
}
