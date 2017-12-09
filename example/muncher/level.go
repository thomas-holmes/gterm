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

func LoadFromString(levelString string) Level {
	level := Level{}

	tiles := make([]Tile, 0, len(levelString))
	var stairs []Stair

	r, c := 0, 0
	t := NewTile(c, r)
	for _, s := range levelString {
		if s == '\n' {
			r++
			c = 0
			continue
		}
		t.X = c
		t.Y = r
		switch s {
		case WallGlyph:
			t.TileKind = Wall
			t.TileGlyph = WallGlyph
		case FloorGlyph:
			t.TileKind = Floor
			t.TileGlyph = FloorGlyph
		case UpStairGlyph:
			stair := Stair{
				X:    c,
				Y:    r,
				Down: false,
			}
			t.TileKind = UpStair
			t.TileGlyph = UpStairGlyph
			stairs = append(stairs, stair)
		case DownStairGlyph:
			stair := Stair{
				X:    c,
				Y:    r,
				Down: true,
			}
			t.TileKind = DownStair
			t.TileGlyph = DownStairGlyph
			stairs = append(stairs, stair)
		}

		tiles = append(tiles, t)
		c++
	}
	level.Columns = c
	level.Rows = r + 1
	level.tiles = tiles
	level.stairs = stairs

	level.VisionMap = NewVisionMap(level.Columns, level.Rows)
	level.ScentMap = NewScentMap(level.Columns, level.Rows)

	return level
}
