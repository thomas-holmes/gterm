package libs

import (
	"log"

	"github.com/thomas-holmes/sneaker/gterm"

	"github.com/veandco/go-sdl2/sdl"
)

type PanelManager struct {
	heightPixels int
	widthPixels  int
	window       *gterm.Window
	panels       []*Panel
}

func NewPanelManager(window *gterm.Window) PanelManager {
	return PanelManager{window: window, panels: make([]*Panel, 0, 0)}
}

// func (panel *Panel) draw(renderer *sdl.Renderer, font *ttf.Font, fontSize int, heightPixels int, widthPixels int) error {
func (panelManager *PanelManager) HandleRender(window *gterm.Window) {
	for _, panel := range panelManager.panels {
		if err := panel.draw(panelManager.window); err != nil {
			log.Fatalln(err)
		}
	}
}

// Panel represents an application "drawable" region of the virtual terminal
type Panel struct {
	XPos   int
	YPos   int
	Width  int
	Height int
	Z      int
	dirty  bool
}

func (panelManager *PanelManager) NewPanel(xPos int, yPos int, width int, height int, z int) *Panel {
	panel := Panel{
		XPos:   xPos,
		YPos:   yPos,
		Width:  width,
		Height: height,
		Z:      z,
		dirty:  true,
	}

	panelManager.panels = append(panelManager.panels, &panel)

	return &panel
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Update adjust the dimensions of a panel, super unsafe, but fun
func (panel *Panel) Update(xPos int, yPos int, width int, height int) {
	panel.XPos = xPos
	panel.YPos = yPos
	panel.Width = max(width, 2)
	panel.Height = max(height, 2)

	panel.dirty = true
}

// A panel might look like:
//   ┌──────────┐
//   │          │
//   │          │
//   │          │
//   │          │
//   └──────────┘

// Box drawing characters
const BoxVertical = "│"
const BoxHorizontal = "─"
const BoxTopLeft = "┌"
const BoxTopRight = "┐"
const BoxBottomLeft = "└"
const BoxBottomRight = "┘"

func (panel *Panel) drawTopRow(window *gterm.Window) error {
	fColor := sdl.Color{R: 20, G: 200, B: 50, A: 255}

	leftCol := panel.XPos
	rightCol := panel.XPos + panel.Width - 1
	topRow := panel.YPos

	if err := window.AddToCell(leftCol, topRow, BoxTopLeft, fColor); err != nil {
		return err
	}
	for col := leftCol + 1; col < rightCol; col++ {
		if err := window.AddToCell(col, topRow, BoxHorizontal, fColor); err != nil {
			return err
		}
	}
	err := window.AddToCell(rightCol, topRow, BoxTopRight, fColor)

	return err
}

func (panel *Panel) drawBody(window *gterm.Window) error {
	fColor := sdl.Color{R: 20, G: 200, B: 50, A: 255}
	leftCol := panel.XPos
	rightCol := panel.XPos + panel.Width - 1
	topRow := panel.YPos
	bottomRow := panel.YPos + panel.Height - 1

	for row := topRow + 1; row < bottomRow; row++ {
		if err := window.AddToCell(leftCol, row, BoxVertical, fColor); err != nil {
			return err
		}
		if err := window.AddToCell(rightCol, row, BoxVertical, fColor); err != nil {
			return err
		}
	}

	return nil
}
func (panel *Panel) drawBottomRow(window *gterm.Window) error {
	fColor := sdl.Color{R: 20, G: 200, B: 50, A: 255}

	leftCol := panel.XPos
	rightCol := panel.XPos + panel.Width - 1
	bottomRow := panel.YPos + panel.Height - 1

	window.AddToCell(leftCol, bottomRow, BoxBottomLeft, fColor)
	for col := leftCol + 1; col < rightCol; col++ {
		window.AddToCell(col, bottomRow, BoxHorizontal, fColor)
	}
	window.AddToCell(rightCol, bottomRow, BoxBottomRight, fColor)

	return nil
}

func (panel *Panel) draw(window *gterm.Window) error {
	if err := panel.drawTopRow(window); err != nil {
		return err
	}
	if err := panel.drawBody(window); err != nil {
		return err
	}
	err := panel.drawBottomRow(window)

	panel.dirty = false
	return err
}
