package game

import (
	"log"
)

type Visibility int

const (
	Unseen Visibility = iota
	Visible
	Seen // Unused
)

type VisionMap struct {
	Columns int
	Rows    int
	Map     []Visibility
}

func (vision VisionMap) VisibilityAt(x int, y int) Visibility {
	return vision.Map[y*vision.Columns+x]
}

func (vision *VisionMap) UpdateVision(viewDistance int, player *Player, world *World) {
	playerX := player.XPos()
	playerY := player.YPos()

	// Go one bast the min/max so we update cells we are moving away from
	minX := max(playerX-viewDistance-2, 0)
	maxX := min(playerX+viewDistance+2, vision.Columns)

	minY := max(playerY-viewDistance-2, 0)
	maxY := min(playerY+viewDistance+2, vision.Rows)

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// TODO: This is gross, but whatever.
			if abs(x-playerX) > viewDistance || abs(y-playerY) > viewDistance {
				vision.Map[y*vision.Columns+x] = Unseen
			} else {
				vision.Map[y*vision.Columns+x] = CheckVision(playerX, playerY, x, y, world)
			}
		}
	}
}

func CheckVision(playerX int, playerY int, candidateX int, candidateY int, world *World) Visibility {
	cells := PlotLine(playerX, playerY, candidateX, candidateY)

	foundWall := false

	for _, cell := range cells {
		if foundWall {
			return Unseen
		}
		if world.GetTile(cell.XPos, cell.YPos).Wall {
			foundWall = true
		}
	}
	return Visible
}

func NewVisionMap(columns int, rows int) VisionMap {
	return VisionMap{
		Columns: columns,
		Rows:    rows,
		Map:     make([]Visibility, columns*rows),
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
		coordinates = append(coordinates, Position{XPos: correctedX, YPos: correctedY})

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
			} else {
				return 1
			}
		} else {
			if (y0 - y1) < (x1 - x0) {
				return 7
			} else {
				return 6
			}
		}
	} else {
		if y1 > y0 {
			if (y1 - y0) < (x0 - x1) {
				return 3
			} else {
				return 2
			}
		} else {
			if (y0 - y1) < (x0 - x1) {
				return 4
			} else {
				return 5
			}
		}
	}
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
