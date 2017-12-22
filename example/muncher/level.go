package main

func (level Level) getStair(x int, y int) (Stair, bool) {
	for _, s := range level.stairs {
		if s.X == x && s.Y == y {
			return s, true
		}
	}
	return Stair{}, false
}

func (level Level) getTile(x int, y int) *Tile {
	return &level.tiles[y*level.Columns+x]
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
	VisionMap VisionMap
	ScentMap  ScentMap
	tiles     []Tile
	stairs    []Stair

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

	level.VisionMap = NewVisionMap(level.Columns, level.Rows)
	level.ScentMap = NewScentMap(level.Columns, level.Rows)

	return level
}
