package main

func (level Level) getTile(x int, y int) *Tile {
	return &level.tiles[y*level.Columns+x]
}

type Level struct {
	Columns   int
	Rows      int
	VisionMap VisionMap
	ScentMap  ScentMap
	tiles     []Tile
}

func loadFromString(levelString string) Level {
	level := Level{}

	tiles := make([]Tile, 0, len(levelString))

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
			t.TileKind = UpStair
			t.TileGlyph = UpStairGlyph
		case DownStairGlyph:
			t.TileKind = DownStair
			t.TileGlyph = DownStairGlyph
		}

		tiles = append(tiles, t)
		c++
	}
	level.Columns = c
	level.Rows = r + 1
	level.tiles = tiles

	level.VisionMap = *NewVisionMap(level.Columns, level.Rows)
	level.ScentMap = NewScentMap(level.Columns, level.Rows)

	return level
}

// TODO: Start generating levels soon instead of using a hard coded grid :)
var LevelMask1 = "" + // Make gofmt happy
	"########################################\n" +
	"#...................#....#...#.........#\n" +
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
	"#........................#....>........#\n" +
	"########################################"

var LevelMask2 = "" + // Make gofmt happy
	"##############################################\n" +
	"#............................................#\n" +
	"#.............................<..............#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#............................................#\n" +
	"#.........<..................................#\n" +
	"#............................................#\n" +
	"##############################################"
