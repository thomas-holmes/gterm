package main

import (
	"github.com/MichaelTJones/pcg"
)

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

func loadFromString(levelString string) Level {
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

// TODO: Start generating levels soon instead of using a hard coded grid :)
var LevelMask1 = "" + // Make gofmt happy
	"########################################\n" +
	"#..>................#....#...#.........#\n" +
	"#.............#####.#......#.#.........#\n" +
	"#...#.........#...#.#....###.#.........#\n" +
	"#.............#...#.#....###.#.........#\n" +
	"#.............#...#.#....###.#.........#\n" +
	"#.............#...#.#....#.............#\n" +
	"#.....>.......##.##.#....#.............#\n" +
	"#..............#.#..####.#.............#\n" +
	"#.....###......#.#.....#.#.............#\n" +
	"#.....###......#.#.....#.#....##.......#\n" +
	"#.....###......#.#.....#.#....##.......#\n" +
	"#.....###......#.#.....#.#.............#\n" +
	"#..............#.#.....#.#.............#\n" +
	"#..............#.#.....#.#..........####\n" +
	"#........#######.#######.#..........#..#\n" +
	"#........................#.............#\n" +
	"########################################"

var LevelMask2 = "" + // Make gofmt happy
	"##############################################\n" +
	"#.......######...............................#\n" +
	"#.......#.<..#...............................#\n" +
	"#.......##...#...............................#\n" +
	"#........#.###...............................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#.........<..................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"##############################################"

type Room struct {
	ID          int
	X           int
	Y           int
	W           int
	H           int
	connectedTo []int
}
type CandidateLevel struct {
	rng *pcg.PCG64

	W int
	H int

	nextRoomID int

	rooms map[int]*Room
	tiles []TileKind
}

const (
	MinRoomWidth      = 3
	MinRoomHeight     = 3
	MaxRoomWidth      = 20
	MaxRoomHeight     = 20
	MaxRoomIterations = 200
)

func (level *CandidateLevel) genNextRoomID() int {
	id := level.nextRoomID
	level.nextRoomID++
	return id
}

// TODO: Exclude the edges of the whole level
func (level *CandidateLevel) tryAddRandomRoom() {
	widthBound := uint64(MaxRoomWidth - MinRoomWidth)
	heightBound := uint64(MaxRoomHeight - MinRoomHeight)

	randomWidth := int(level.rng.Bounded(widthBound)) + MinRoomWidth
	randomHeight := int(level.rng.Bounded(heightBound)) + MinRoomHeight

	topLeftX := int(level.rng.Bounded(uint64(level.W-randomWidth-1))) + 1
	topLeftY := int(level.rng.Bounded(uint64(level.H-randomHeight-1))) + 1

	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			if level.tiles[y*level.W+x] != Wall {
				// Just quit if we run into a non-wall feature
				return
			}
		}
	}

	// We can place our room, so make it all floors.
	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			level.tiles[y*level.W+x] = Floor
		}
	}

	room := &Room{
		ID: level.genNextRoomID(),
		X:  topLeftX,
		Y:  topLeftY,
		W:  randomWidth,
		H:  randomHeight,
	}

	level.rooms[room.ID] = room

}

func (level *CandidateLevel) addRooms() {
	for i := 0; i < MaxRoomIterations; i++ {
		level.tryAddRandomRoom()
	}
}

func (level *CandidateLevel) connectRooms() {

}

func (level *CandidateLevel) encodeAsString() string {
	levelStr := ""
	for y := 0; y < level.H; y++ {
		if y != 0 {
			levelStr += "\n"
		}
		for x := 0; x < level.W; x++ {
			switch level.tiles[y*level.W+x] {
			case Wall:
				levelStr += string(WallGlyph)
			case Floor:
				levelStr += string(FloorGlyph)
			case DownStair:
				levelStr += string(DownStairGlyph)
			case UpStair:
				levelStr += string(UpStairGlyph)
			}
		}
	}

	return levelStr
}

func GenLevel(rng *pcg.PCG64, maxX int, maxY int) string {
	subX := rng.Bounded(uint64(maxX / 4))
	subY := rng.Bounded(uint64(maxY / 4))

	W := maxX - int(subX)
	H := maxY - int(subY)

	level := &CandidateLevel{
		rng: rng,
		W:   W,
		H:   H,

		rooms: make(map[int]*Room),
		tiles: make([]TileKind, W*H, W*H),
	}

	level.addRooms()

	level.connectRooms()

	return level.encodeAsString()
}
