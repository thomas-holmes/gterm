package main

import (
	"testing"
)

func TestLoadFromString(t *testing.T) {
	smallLevel := "" +
		"#####\n" +
		"#...#\n" +
		"#...#\n" +
		"#####"

	level := LoadFromString(smallLevel)

	if level.Columns != 5 {
		t.Errorf("Expected columns of 5, got %v", level.Columns)
	}
	if level.Rows != 4 {
		t.Errorf("Expected rows of 4, got %v", level.Rows)
	}
	if len(level.tiles) != 20 {
		t.Errorf("Expected tile slice to be 20, got %v", len(level.tiles))
	}
}

func xyIndex(x int, y int, cols int) int {
	return y*cols + x
}

func TestLoadFromStringTilePlacement(t *testing.T) {
	checkTile := func(level Level, x int, y int, expected TileKind) {
		tile := level.tiles[xyIndex(x, y, level.Columns)]
		if tile.TileKind != expected {
			t.Errorf("Expected to find TileKind(%v) at (%v,%v) but instead found %+v", expected, x, y, tile)
		}
	}

	levelStr := "" +
		"######\n" +
		"#....#\n" +
		"#....#\n" +
		"#..>.#\n" +
		"#..#.#\n" +
		"#<...#\n" +
		"######"

	level := LoadFromString(levelStr)

	if level.Columns != 6 {
		t.Errorf("Expected columns of 6, got %v", level.Columns)
	}

	if level.Rows != 7 {
		t.Errorf("Expected rows of 7, got %v", level.Rows)
	}

	if len(level.tiles) != 42 {
		t.Errorf("Expected tile slice to be 42, got %v", len(level.tiles))
	}

	checkTile(level, 0, 0, Wall)
	checkTile(level, 3, 4, Wall)
	checkTile(level, 3, 3, DownStair)
	checkTile(level, 1, 5, UpStair)
	checkTile(level, 4, 5, Floor)
}
