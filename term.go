package gterm

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/text/encoding/charmap"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var White = sdl.Color{R: 225, G: 225, B: 225, A: 255}

var CP437 = charmap.CodePage437

// Window represents the base window object
type Window struct {
	Columns         int
	Rows            int
	FontSize        int
	FontHPixel      int
	FontWPixel      int
	HeightPixel     int
	WidthPixel      int
	fontPath        string
	fontSheet       *sdl.Texture
	spritesPerRow   int
	SdlWindow       *sdl.Window
	SdlRenderer     *sdl.Renderer
	backgroundColor sdl.Color
	cells           []cell
	fps             fpsCounter
	vsync           bool
}

type cell struct {
	bgColor     sdl.Color
	renderItems []renderItem
}

type renderItem struct {
	FColor sdl.Color
	Glyph  rune
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, fontPath string, fontX int, fontY int, vsync bool) *Window {
	numCells := columns * rows
	cells := make([]cell, numCells, numCells)

	window := &Window{
		Columns:     columns,
		Rows:        rows,
		fontPath:    fontPath,
		cells:       cells,
		vsync:       vsync,
		FontHPixel:  fontX,
		FontWPixel:  fontY,
		WidthPixel:  columns * fontX,
		HeightPixel: rows * fontY,
	}

	return window
}

func (window *Window) SetTitle(title string) {
	window.SdlWindow.SetTitle(title)
}

// Init initialized the window for drawing
func (window *Window) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	if flags := img.Init(img.INIT_PNG); flags&img.INIT_PNG == 0 {
		return errors.New("Failed to initialize sdl2_img for PNG")
	}

	sdlWindow, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, window.WidthPixel, window.HeightPixel, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	var flags uint32 = sdl.RENDERER_ACCELERATED
	if window.vsync {
		flags = sdl.RENDERER_PRESENTVSYNC
	}
	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, flags)
	if err != nil {
		return err
	}
	if err := sdlRenderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return err
	}

	rwops := sdl.RWFromFile(window.fontPath, "rb")
	if rwops == nil {
		return fmt.Errorf("Failed to load image from %s", window.fontPath)
	}

	surface, err := img.LoadPNG_RW(rwops)
	if err != nil {
		return err
	}
	defer surface.Free()
	if err := surface.SetColorKey(sdl.ENABLE, 0); err != nil {
		return err
	}

	window.spritesPerRow = int(surface.W) / window.FontWPixel

	texture, err := sdlRenderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	if err := texture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return err
	}
	window.fontSheet = texture

	err = sdlRenderer.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}

	window.SdlWindow = sdlWindow
	window.SdlRenderer = sdlRenderer

	window.fps = newFpsCounter()

	return nil
}

func (window *Window) SetBackgroundColor(color sdl.Color) {
	window.backgroundColor = color
}

func (window *Window) cellIndex(col int, row int) (int, error) {
	if col >= window.Columns || col < 0 || row >= window.Rows || row < 0 {
		return 0, fmt.Errorf("Requested invalid position (%v,%v) on board of dimensions %vx%v", col, row, window.Columns, window.Rows)
	}
	return col + window.Columns*row, nil
}

func (window *Window) renderCell(cellCol int, cellRow int) error {
	idx, err := window.cellIndex(cellCol, cellRow)
	if err != nil {
		return err
	}

	destX := cellCol * window.FontWPixel
	destY := cellRow * window.FontHPixel
	destRect := sdl.Rect{X: int32(destX), Y: int32(destY), W: int32(window.FontWPixel), H: int32(window.FontHPixel)}

	cell := window.cells[idx]
	for _, item := range cell.renderItems {
		runeByte, ok := CP437.EncodeRune(item.Glyph)
		if !ok {
			log.Println("Could not encode rune", item.Glyph)
		}

		row := int(runeByte) / window.spritesPerRow
		col := int(runeByte) % window.spritesPerRow
		sX := col * window.FontWPixel
		sY := row * window.FontHPixel

		sourceRect := sdl.Rect{X: int32(sX), Y: int32(sY), W: int32(window.FontWPixel), H: int32(window.FontHPixel)}

		{
			color := cell.bgColor
			r, g, b, a := uint8(color.R), uint8(color.G), uint8(color.B), uint8(color.A)
			window.SdlRenderer.SetDrawColor(r, g, b, a)
			window.SdlRenderer.FillRect(&destRect)
		}

		color := item.FColor
		r, g, b := uint8(color.R), uint8(color.G), uint8(color.B)
		window.fontSheet.SetColorMod(r, g, b)
		if err := window.SdlRenderer.Copy(window.fontSheet, &sourceRect, &destRect); err != nil {
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

// NoColor is used to represent no background color
var NoColor = sdl.Color{R: 0, G: 0, B: 0, A: 0}

func (window *Window) PutRune(col int, row int, glyph rune, fColor sdl.Color, bColor sdl.Color) error {
	renderItem := renderItem{Glyph: glyph, FColor: fColor}
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}
	window.cells[index].renderItems = append(window.cells[index].renderItems, renderItem)
	window.cells[index].bgColor = bColor

	return nil
}

func (window *Window) PutStringBg(col int, row int, content string, fColor sdl.Color, bColor sdl.Color) error {
	for step, rune := range content {
		if err := window.PutRune(col+step, row, rune, fColor, bColor); err != nil {
			return err
		}
	}

	return nil
}

func (window *Window) PutString(col int, row int, content string, fColor sdl.Color) error {
	return window.PutStringBg(col, row, content, fColor, NoColor)
}

func (window *Window) ClearRegion(col int, row int, width int, height int) error {
	for y := row; y < row+height; y++ {
		for x := col; x < col+width; x++ {
			if err := window.ClearCell(x, y); err != nil {
				return err
			}
		}
	}
	return nil
}

func (window *Window) ClearCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	window.cells[index] = cell{}

	return nil
}

func (window *Window) ClearWindow() {
	window.cells = make([]cell, window.Columns*window.Rows, window.Columns*window.Rows)
}

func (window *Window) ShouldRenderFps(shouldRender bool) {
	window.fps.shouldRender(shouldRender)
}

// Render updates the display based on new information since last Render
type fpsCounter struct {
	renderFps     bool
	framesElapsed int
	currentFps    int
	lastTicks     uint32
	color         sdl.Color
}

func newFpsCounter() fpsCounter {
	return fpsCounter{
		lastTicks: sdl.GetTicks(),
		color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
	}
}

func (fps *fpsCounter) shouldRender(shouldRender bool) {
	fps.renderFps = shouldRender
}

func min(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func (window *Window) DebugDrawSpriteSheet() error {
	_, _, w, h, err := window.fontSheet.Query()
	if err != nil {
		return err
	}

	return window.SdlRenderer.Copy(window.fontSheet, nil, &sdl.Rect{X: 0, Y: 0, W: w, H: h})
}

func (window *Window) Refresh() {
	err := window.SdlRenderer.SetDrawColor(window.backgroundColor.R, window.backgroundColor.G, window.backgroundColor.B, window.backgroundColor.A)
	if err != nil {
		log.Fatal(err)
	}
	window.SdlRenderer.Clear()

	if err := window.renderCells(); err != nil {
		log.Println("Failed to render cells", err)
	}

	window.SdlRenderer.Present()
}
