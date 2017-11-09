package gterm

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Window represents the base window object
type Window struct {
	Columns         int
	Rows            int
	FontSize        int
	tileHeightPixel int
	tileWidthPixel  int
	heightPixel     int
	widthPixel      int
	fontPath        string
	font            *ttf.Font
	SdlWindow       *sdl.Window
	SdlRenderer     *sdl.Renderer
	cells           [][]renderItem
}

type renderItem struct {
	FColor sdl.Color
	BColor sdl.Color
	Glyph  string
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, fontPath string, fontSize int) *Window {

	numCells := columns * rows
	cells := make([][]renderItem, numCells, numCells)

	window := &Window{
		Columns:  columns,
		Rows:     rows,
		FontSize: fontSize,
		fontPath: fontPath,
		cells:    cells,
	}

	return window
}

func computeCellSize(font *ttf.Font) (width int, height int, err error) {
	atGlyph, err := font.RenderUTF8_Blended("@", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0, 0, err
	}
	return int(atGlyph.W), int(atGlyph.H), nil
}

func (window *Window) SetTitle(title string) {
	window.SdlWindow.SetTitle(title)
}

// Init initialized the window for drawing
func (window *Window) Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING) // not sure where to do this
	if err != nil {
		return err
	}
	err = ttf.Init()
	if err != nil {
		return nil
	}
	openedFont, err := ttf.OpenFont(window.fontPath, window.FontSize)
	if err != nil {
		return err
	}

	window.font = openedFont
	tileWidth, tileHeight, err := computeCellSize(window.font)
	if err != nil {
		return err
	}

	window.tileWidthPixel = tileWidth
	window.tileHeightPixel = tileHeight
	window.heightPixel = tileHeight * window.Rows
	window.widthPixel = tileWidth * window.Columns

	log.Printf("Creating window w:%v, h:%v", window.widthPixel, window.heightPixel)
	sdlWindow, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, window.widthPixel, window.heightPixel, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	err = sdlRenderer.SetDrawColor(10, 10, 25, 255)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}

	window.SdlWindow = sdlWindow
	window.SdlRenderer = sdlRenderer

	return nil
}

func (window *Window) cellIndex(col int, row int) (int, error) {
	if col >= window.Columns || col < 0 || row >= window.Rows || row < 0 {
		return 0, fmt.Errorf("Requested invalid position (%v,%v) on board of dimensions %vx%v", col, row, window.Columns, window.Rows)
	}
	return col + window.Columns*row, nil
}

func (window *Window) renderCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	renderItems := window.cells[index]

	destinationRect := sdl.Rect{
		X: int32(col * window.tileWidthPixel),
		Y: int32(row * window.tileHeightPixel),
		W: int32(window.tileWidthPixel),
		H: int32(window.tileHeightPixel),
	}

	for _, renderItem := range renderItems {
		// surface, err := window.font.RenderUTF8_Blended(renderItem.Glyph, renderItem.FColor)
		surface, err := window.font.RenderUTF8_Solid(renderItem.Glyph, renderItem.FColor)
		if err != nil {
			return err
		}
		defer surface.Free()

		texture, err := window.SdlRenderer.CreateTextureFromSurface(surface)
		if err != nil {
			return err
		}
		defer texture.Destroy()

		_, _, width, height, err := texture.Query()
		if err != nil {
			return err
		}

		destinationRect.W = width
		destinationRect.H = height

		// log.Println(destinationRect)
		err = window.SdlRenderer.Copy(texture, nil, &destinationRect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (window *Window) renderCells() error {
	for col := 0; col < window.Columns; col++ {
		for row := 0; row < window.Rows; row++ {
			if err := window.renderCell(col, row); err != nil {
				return err
			}
		}
	}
	return nil
}

func (window *Window) AddToCell(col int, row int, glyph string, fColor sdl.Color, bColor sdl.Color) error {
	renderItem := renderItem{Glyph: glyph, FColor: fColor, BColor: bColor}
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}
	window.cells[index] = append(window.cells[index], renderItem)

	return nil
}

func (window *Window) EraseCell(col, int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	window.cells[index] = make([]renderItem, 0, 0)

	return nil
}

func (window *Window) ClearWindow() {
	window.cells = make([][]renderItem, window.Columns*window.Rows, window.Columns*window.Rows)
}

// Render updates the display based on new information since last Render
func (window *Window) Render() {
	window.SdlRenderer.Clear()

	window.renderCells()

	window.SdlRenderer.Present()
}
