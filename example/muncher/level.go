package main

import (
	"log"
	"sort"
)

type DistanceCandidate struct {
	Distance float64
	Position
}

func (level Level) getStair(x int, y int) (Stair, bool) {
	for _, s := range level.stairs {
		if s.X == x && s.Y == y {
			return s, true
		}
	}
	return Stair{}, false
}

func (level Level) GetTile(x int, y int) *Tile {
	return &level.tiles[y*level.Columns+x]
}

func (level *Level) IsTileOccupied(x int, y int) bool {
	return level.GetTile(x, y).Creature != nil
}

func (level *Level) CanStandOnTile(column int, row int) bool {
	if level == nil {
		log.Panicf("Wtf is going on")
	}
	return !level.GetTile(column, row).IsWall() && !level.IsTileOccupied(column, row)
}

func (level *Level) GetCreatureAtTile(xPos int, yPos int) (*Creature, bool) {
	if creature := level.GetTile(xPos, yPos).Creature; creature != nil {
		return creature, true
	}
	return nil, false
}

// GetVisibleCreatures returns a slice of creatures sorted so that the first is the closest
// based on euclidean distance.
func (level *Level) GetVisibleCreatures(originX int, originY int) []*Creature {
	candidates := make([]DistanceCandidate, 0, 8)
	for y := 0; y < level.VisionMap.Rows; y++ {
		for x := 0; x < level.VisionMap.Columns; x++ {
			if level.VisionMap.VisibilityAt(x, y) == Visible {
				candidates = append(candidates, DistanceCandidate{Position: Position{X: x, Y: y}, Distance: euclideanDistance(originX, originY, x, y)})
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Distance < candidates[j].Distance })

	creatures := make([]*Creature, 0, len(candidates))
	for _, candidate := range candidates {
		if creature, ok := level.GetCreatureAtTile(candidate.X, candidate.Y); ok {
			creatures = append(creatures, creature)
		}
	}

	return creatures
}

type Stair struct {
	Down bool
	X    int
	Y    int

	Connected bool
	DestX     int
	DestY     int
	DestLevel *Level
}

type Level struct {
	Columns   int
	Rows      int
	VisionMap *VisionMap
	ScentMap  *ScentMap
	tiles     []Tile
	stairs    []Stair

	MonsterDensity int

	Depth int

	NextEntity int
	NextEnergy int
	Entities   []Entity
}

// connectTwoLevels connects multiple levels arbitrarily. If there is an uneven number
// of stair cases you will end up with a dead stair.
func connectTwoLevels(upper *Level, lower *Level) {
	for i, downStair := range upper.stairs {
		if !downStair.Down || downStair.Connected {
			continue
		}

		for j, upStair := range lower.stairs {
			if upStair.Down || upStair.Connected {
				continue
			}

			upper.stairs[i].DestLevel = lower
			upper.stairs[i].DestX = upStair.X
			upper.stairs[i].DestY = upStair.Y
			upper.stairs[i].Connected = true

			lower.stairs[j].DestLevel = upper
			lower.stairs[j].DestX = downStair.X
			lower.stairs[j].DestY = downStair.Y
			lower.stairs[j].Connected = true

			break
		}
	}
}

func LoadCandidateLevel(candidate *CandidateLevel) Level {
	level := Level{}

	tiles := make([]Tile, 0, len(candidate.tiles))

	var stairs []Stair

	for y := 0; y < candidate.H; y++ {
		for x := 0; x < candidate.W; x++ {
			tile, cTile := NewTile(x, y), candidate.tiles[y*candidate.W+x]
			tile.TileKind = cTile.TileKind
			tile.TileGlyph = TileKindToGlyph(cTile.TileKind)
			tile.Item = cTile.Item

			switch tile.TileKind {
			case UpStair:
				stair := Stair{
					X:    x,
					Y:    y,
					Down: false,
				}
				stairs = append(stairs, stair)
			case DownStair:
				stair := Stair{
					X:    x,
					Y:    y,
					Down: true,
				}
				stairs = append(stairs, stair)
			}
			tiles = append(tiles, tile)
		}
	}

	level.Columns = candidate.W
	level.Rows = candidate.H
	level.tiles = tiles
	level.stairs = stairs
	level.MonsterDensity = 50

	level.VisionMap = NewVisionMap(level.Columns, level.Rows)
	level.ScentMap = NewScentMap(level.Columns, level.Rows)

	return level
}
