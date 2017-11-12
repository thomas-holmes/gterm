package game

import (
	"log"

	"github.com/thomas-holmes/sneaker/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Tile struct {
	Contents        []Renderable
	Dirty           bool
	BackgroundColor sdl.Color
	BackgroundGlyph string
	Wall            bool
}

func NewTile() Tile {
	return Tile{
		Dirty:           true,
		BackgroundColor: sdl.Color{R: 192, G: 192, B: 192, A: 255},
		BackgroundGlyph: ".",
	}
}

func (tile Tile) RenderBackground(col int, row int, window *gterm.Window) {
	err := window.AddToCell(col, row, tile.BackgroundGlyph, tile.BackgroundColor)
	if err != nil {
		log.Println("Failed to render background", err)
	}
}

func (tile *Tile) ClearContents() {
	tile.Contents = tile.Contents[:0]
	tile.Dirty = true
}

func (tile *Tile) AddRenderable(renderable Renderable) {
	tile.Contents = append(tile.Contents, renderable)
	tile.Dirty = true
}

type World struct {
	Window  *gterm.Window
	Columns int
	Rows    int
	Tiles   []Tile
	Dirty   bool
}

var LevelMask = []int{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
}

func (world *World) BuildLevelFromMask(mask []int) {
	for index := range mask {
		if mask[index] == 1 {
			tile := &world.Tiles[index]

			tile.BackgroundGlyph = "#"
			tile.Wall = true
			tile.BackgroundColor = sdl.Color{R: 225, G: 225, B: 225, A: 255}
			tile.Dirty = true
		}
	}
}
func (world *World) BuildLevel() {
	for row := 0; row < world.Rows; row++ {
		for col := 0; col < world.Columns; col++ {
			if row == 0 || (row == world.Rows-1) || (col == 0 || col == world.Columns-1) {
				world.Tiles[row*world.Columns+col].BackgroundGlyph = "#"
				world.Tiles[row*world.Columns+col].Wall = true
				world.Tiles[row*world.Columns+col].BackgroundColor = sdl.Color{R: 225, G: 225, B: 225, A: 255}
				world.Tiles[row*world.Columns+col].Dirty = true
			}
		}
	}
}
func (world World) TileIndex(column int, row int) int {
	return row*world.Columns + column
}

func (world *World) GetTile(column int, row int) *Tile {
	tile := &world.Tiles[world.TileIndex(column, row)]
	return tile
}

func (world *World) CanStandOnTile(column int, row int) bool {
	return !world.GetTile(column, row).Wall
}

func (world *World) ClearTile(column int, row int) {
	world.GetTile(column, row).ClearContents()
}

func (world *World) DirtyTile(column int, row int) {
	world.GetTile(column, row).Dirty = true
}

func (world *World) AddRenderableToTile(column int, row int, renderable Renderable) {
	world.GetTile(column, row).AddRenderable(renderable)
}

func (world *World) Render() {
	for row := 0; row < world.Rows; row++ {
		for col := 0; col < world.Columns; col++ {
			tile := &world.Tiles[world.TileIndex(col, row)]
			if tile.Dirty {
				world.Window.ClearCell(col, row)
				tile.RenderBackground(col, row, world.Window)
				for _, renderable := range tile.Contents {
					renderable.Render(world)
				}
				tile.Dirty = false
			}
		}
	}
	world.Dirty = false
}

func NewWorld(window *gterm.Window, columns int, rows int) World {
	tiles := make([]Tile, columns*rows, columns*rows)
	for index := range tiles {
		tiles[index] = NewTile()
	}

	world := World{
		Window:  window,
		Columns: columns,
		Rows:    rows,
		Dirty:   true,
		Tiles:   tiles,
	}

	return world
}
