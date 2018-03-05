package gterm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var White = sdl.Color{R: 225, G: 225, B: 225, A: 255}

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
	drawInterval    uint32
	vsync           bool

	texturePages map[int]fontPage
}

type fontPage struct {
	lookup  map[rune]sdl.Rect
	texture *sdl.Texture
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
		texturePages: make(map[int]fontPage),
	}

	return window
}

func (window *Window) AddTexturePage(start int) error {
	pageBase := start & 0xFF00
	pageMax := pageBase + 0x00FF

	log.Printf("Generating texture page from %v to %v", pageBase, pageMax)

	page := fontPage{lookup: make(map[rune]sdl.Rect)}

	var surfaces []*sdl.Surface

	textureWidth := int32(0)
	height := int32(0)
	for b := pageBase; b <= pageMax; b++ {
		surface, err := window.font.RenderUTF8_Blended(string(b), White)
		if err != nil {
			if strings.Contains(err.Error(), "zero width") {
				continue
			}
			return err
		}
		defer surface.Free()

		log.Println(string(b))
		rect := sdl.Rect{X: textureWidth, Y: 0, W: surface.W, H: surface.H}
		page.lookup[rune(b)] = rect

		surfaces = append(surfaces, surface)
		textureWidth += surface.W
		if surface.H > height {
			height = surface.H
		}
	}

	texture, err := window.SdlRenderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int(textureWidth), int(height))
	if err != nil {
		return err
	}

	if err := texture.SetBlendMode(sdl.BLENDMODE_ADD); err != nil {
		return err
	}

	region, lockPitch, err := texture.Lock(nil)
	if err != nil {
		return err
	}

	accW := int32(0)
	bytesPerPixel := int32(lockPitch) / textureWidth
	for _, surface := range surfaces {
		count := int32(0)
		offset := int32(0)
		charOffset := accW * bytesPerPixel
		for _, b := range surface.Pixels() {
			if count == surface.Pitch {
				offset += int32(lockPitch)
				count = 0
			}
			region[count+offset+charOffset] = b
			count++
		}
		accW += surface.W
	}
	texture.Unlock()

	page.texture = texture
	window.texturePages[pageBase] = page

	return nil
}

func computeCellSize(font *ttf.Font) (width int, height int, err error) {
	w, h, err := font.SizeUTF8("@")
	if err != nil {
		log.Printf("Computed cell size of w: %v, h: %v", w, h)
		return 0, 0, err
	}
	return w, h, nil
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

	if err := window.AddTexturePage(0); err != nil {
		return err
	}
	if err := window.AddTexturePage(0x2500); err != nil {
		return err
	}

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

func (window *Window) renderCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	cell := window.cells[index]
	renderItems := cell.renderItems

	if cell.bgColor.A != 0 {
		window.drawBackground(col, row, cell.bgColor)
	}

	for index := range renderItems {
		renderItem := &renderItems[index]

		idx := int(renderItem.Glyph) & 0xFF00
		page := window.texturePages[idx]

		destinationRect := sdl.Rect{
			X: int32(col * window.tileWidthPixel),
			Y: int32(row * window.tileHeightPixel),
			W: int32(window.tileWidthPixel),
			H: int32(window.tileHeightPixel),
		}

		sourceRect, ok := page.lookup[renderItem.Glyph]
		if !ok {
			log.Printf("Couldn't lookup glyph %v", renderItem.Glyph)
		}

		atlas := page.texture
		atlas.SetColorMod(renderItem.FColor.R, renderItem.FColor.G, renderItem.FColor.B)
		err = window.SdlRenderer.Copy(atlas, &sourceRect, &destinationRect)
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

func (window *Window) PutString(col int, row int, content string, fColor sdl.Color) error {
	for step, rune := range content {
		if err := window.PutRune(col+step, row, rune, fColor, NoColor); err != nil {
			return err
		}
	}

	return nil
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
	font          *ttf.Font
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

		r, g, b := uint8(fps.color.R), uint8(fps.color.G), uint8(fps.color.B)

		page := window.texturePages[0]
		atlas, lookup := page.texture, page.lookup
		for i, rn := range fpsString {
			destination.X = int32(xStart + (window.tileWidthPixel * i))
			source := lookup[rn]

			atlas.SetColorMod(r, g, b)
			if err := window.SdlRenderer.Copy(atlas, &source, &destination); err != nil {
				log.Println("Couldn't copy FPS to renderer", err)
			}
		}
	}
}

func (window *Window) renderNewDebugFontTexture(pageBase int, row int) {
	page, ok := window.texturePages[pageBase]
	if !ok {
		log.Panicln("no texture?")
	}

	tex := page.texture
	_, _, width, height, err := tex.Query()

	if err != nil {
		log.Panicln("Couldn't render new debug font texture", err)
	}

	destRect := sdl.Rect{X: int32(0), Y: int32(row) * height, W: width, H: height}
	window.SdlRenderer.Copy(tex, nil, &destRect)
}

var lastDraw = sdl.GetTicks()

func min(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func (window *Window) drawBackground(x int, y int, color sdl.Color) {
	r, g, b, a := uint8(color.R), uint8(color.G), uint8(color.B), uint8(color.A)
	window.SdlRenderer.SetDrawColor(r, g, b, a)
	if err := window.SdlRenderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		log.Println("failed to sent blendmode", err)
	}
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

	if err := window.renderCells(); err != nil {
		log.Println("Failed to render cells", err)
	}

	window.fps.MaybeRender(window) // Ok this is dumb

	/*
		window.renderNewDebugFontTexture(0, 0)
		window.renderNewDebugFontTexture(0x2500, 1)
	*/

	window.SdlRenderer.Present()
}
