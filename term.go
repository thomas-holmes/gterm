package gterm

import (
	"fmt"
	"log"
	"path"
	"strconv"

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
	backgroundColor sdl.Color
	cells           []cell
	fps             fpsCounter
	fontTexture     *sdl.Texture
	drawInterval    uint32
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
func NewWindow(columns int, rows int, fontPath string, fontSize int, vsync bool) *Window {
	numCells := columns * rows
	cells := make([]cell, numCells, numCells)
	drawInterval := uint32(0)

	window := &Window{
		Columns:      columns,
		Rows:         rows,
		FontSize:     fontSize,
		fontPath:     fontPath,
		cells:        cells,
		vsync:        vsync,
		drawInterval: drawInterval,
	}

	return window
}

func (window *Window) createFontAtlas(font *ttf.Font) (*sdl.Texture, error) {
	str := ""
	for i := 32; i < 126; i++ {
		str += string(i)
	}

	surface, err := font.RenderUTF8_Blended(str, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return nil, err
	}

	tex, err := window.SdlRenderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	_, _, width, height, err := tex.Query()
	if err != nil {
		log.Panicln("Failed to query texture", err)
	}
	log.Printf("Computed texture width of %v and height of %v", width, height)
	log.Printf("Computed tile width of %v and height of %v", window.tileWidthPixel, window.tileHeightPixel)

	return tex, err
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

	var flags uint32 = sdl.RENDERER_ACCELERATED
	if window.vsync {
		flags = sdl.RENDERER_PRESENTVSYNC
	}
	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, flags)
	if err != nil {
		return err
	}

	err = sdlRenderer.SetDrawColor(10, 10, 25, 255)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}

	window.SdlWindow = sdlWindow
	window.SdlRenderer = sdlRenderer

	texture, err := window.createFontAtlas(openedFont)
	if err != nil {
		return err
	}
	window.fontTexture = texture

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

// TODO: This totally doesn't render words any more so the HUD and Popup are broken.
func (window *Window) renderCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	cell := window.cells[index]
	renderItems := cell.renderItems

	for index := range renderItems {
		renderItem := &renderItems[index]
		charOffset := int(renderItem.Glyph - ' ')

		sourceRect := sdl.Rect{
			X: int32(charOffset * window.tileWidthPixel),
			Y: int32(0),
			W: int32(window.tileWidthPixel),
			H: int32(window.tileHeightPixel),
		}

		if err != nil {
			return err
		}

		destinationRect := sdl.Rect{
			X: int32(col * window.tileWidthPixel),
			Y: int32(row * window.tileHeightPixel),
			W: int32(window.tileWidthPixel),
			H: int32(window.tileHeightPixel),
		}

		window.fontTexture.SetColorMod(renderItem.FColor.R, renderItem.FColor.G, renderItem.FColor.B)
		err = window.SdlRenderer.Copy(window.fontTexture, &sourceRect, &destinationRect)
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

func (window *Window) PutRune(col int, row int, glyph rune, fColor sdl.Color) error {
	renderItem := renderItem{Glyph: glyph, FColor: fColor}
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}
	window.cells[index].renderItems = append(window.cells[index].renderItems, renderItem)

	return nil
}

func (window *Window) PutString(col int, row int, content string, fColor sdl.Color) error {
	for step, rune := range content {
		if err := window.PutRune(col+step, row, rune, fColor); err != nil {
			return err
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
	font          *ttf.Font
	color         sdl.Color
}

func newFpsCounter() fpsCounter {
	font, err := ttf.OpenFont(path.Join("assets/font/FiraMono-Regular.ttf"), 16)
	if err != nil {
		log.Println("Failed to open FPS font", err)
	}
	return fpsCounter{
		lastTicks: sdl.GetTicks(),
		color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
		font:      font,
	}
}

func (fps *fpsCounter) shouldRender(shouldRender bool) {
	fps.renderFps = shouldRender
}

func (fps *fpsCounter) MaybeRender(window *Window) {
	fps.framesElapsed++
	now := sdl.GetTicks()

	if fps.lastTicks < (now - 1000) {
		fps.lastTicks = now
		fps.currentFps = fps.framesElapsed
		fps.framesElapsed = 0
	}

	// TODO: This is basically a big dup of the code above in renderCell
	if fps.renderFps {
		fpsString := strconv.Itoa(fps.currentFps)
		xStart := window.widthPixel - (len(fpsString) * window.tileWidthPixel)

		destination := sdl.Rect{
			X: 0,
			Y: 0,
			W: int32(window.tileWidthPixel),
			H: int32(window.tileHeightPixel),
		}
		source := sdl.Rect{
			X: 0,
			Y: 0,
			W: int32(window.tileWidthPixel),
			H: int32(window.tileHeightPixel),
		}

		r, g, b := uint8(fps.color.R), uint8(fps.color.G), uint8(fps.color.B)

		for i, rune := range fpsString {
			destination.X = int32(xStart + (window.tileWidthPixel * i))
			source.X = int32(rune-' ') * int32(window.tileWidthPixel)

			window.fontTexture.SetColorMod(r, g, b)
			if err := window.SdlRenderer.Copy(window.fontTexture, &source, &destination); err != nil {
				log.Println("Couldn't copy FPS to renderer", err)
			}
		}
	}
}

func (window *Window) renderDebugFontTexture() {
	_, _, width, height, err := window.fontTexture.Query()
	if err != nil {
		log.Panicln("Couldn't render debug font texture", err)
	}

	destRect := sdl.Rect{X: int32(0), Y: int32(window.heightPixel) - height, W: width, H: height}
	window.SdlRenderer.Copy(window.fontTexture, nil, &destRect)

}

var lastDraw = sdl.GetTicks()

func min(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func (window *Window) RenderBackgroundSquare(x int, y int) {
	window.SdlRenderer.SetDrawColor(255, 200, 200, 255)
	destinationRect := sdl.Rect{
		X: int32(x * window.tileWidthPixel),
		Y: int32(y * window.tileHeightPixel),
		W: int32(window.tileWidthPixel),
		H: int32(window.tileHeightPixel),
	}
	window.SdlRenderer.FillRect(&destinationRect)
}

func (window *Window) Refresh() {
	err := window.SdlRenderer.SetDrawColor(window.backgroundColor.R, window.backgroundColor.G, window.backgroundColor.B, window.backgroundColor.A)
	if err != nil {
		log.Fatal(err)
	}
	window.SdlRenderer.Clear()

	window.renderCells()

	window.fps.MaybeRender(window) // Ok this is dumb

	window.renderDebugFontTexture()

	window.SdlRenderer.Present()
}
