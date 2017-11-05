package gterm

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type panelManager struct {
	heightPixels int
	widthPixels  int
	panels       []*Panel
}

func newPanelManager() panelManager {
	return panelManager{panels: make([]*Panel, 0, 0)}
}

// func (panel *Panel) draw(renderer *sdl.Renderer, font *ttf.Font, fontSize int, heightPixels int, widthPixels int) error {
func (panelManager *panelManager) panelManagerRenderHandler() RenderHandler {
	return func(renderer *sdl.Renderer) {
		for _, panel := range panelManager.panels {
			panel.draw(renderer)
		}
	}
}

// Panel represents an application "drawable" region of the virtual terminal
type Panel struct {
	XPos         int
	YPos         int
	Width        int
	Height       int
	Z            int
	Font         *ttf.Font
	heightPixels int
	widthPixels  int
	contents     [][][]interface{}
}

func (panelManager *panelManager) newPanel(xPos int, yPos int, width int, height int, z int, font *ttf.Font) *Panel {
	panelContents := make([][][]interface{}, width, width)
	for i := range panelContents {
		panelContents[i] = make([][]interface{}, height, height)
	}
	widthPixels, heightPixels, err := computeCellSizeFromFont(font)
	if err != nil {
		panic(err)
	}
	panel := Panel{
		XPos:         xPos,
		YPos:         yPos,
		Width:        width,
		Height:       height,
		Z:            z,
		Font:         font,
		widthPixels:  widthPixels,
		heightPixels: heightPixels,
		contents:     panelContents,
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
}

func (panel *Panel) getContentsAt(xPos int, yPos int) []interface{} {
	return panel.contents[xPos][yPos]
}

func (panel *Panel) renderRow(row int) string {
	return fmt.Sprintf("%v%v%v", BoxVertical, strings.Repeat(" ", panel.Width-2), BoxVertical)
}

func (panel *Panel) draw(renderer *sdl.Renderer) error {
	panelRows := make([]string, panel.Height, panel.Height)
	panelRows[0] = fmt.Sprintf("%v%v%v", BoxTopLeft, strings.Repeat(BoxHorizontal, panel.Width-2), BoxTopRight)
	for row := 1; row < panel.Height; row++ {
		panelRows[row] = panel.renderRow(row)
	}
	panelRows[panel.Height-1] = fmt.Sprintf("%v%v%v", BoxBottomLeft, strings.Repeat(BoxHorizontal, panel.Width-2), BoxBottomRight)

	boxColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	for rowIndex, row := range panelRows {
		rendered, err := panel.Font.RenderUTF8_Blended(row, boxColor)
		if err != nil {
			return err
		}
		defer rendered.Free()
		texture, err := renderer.CreateTextureFromSurface(rendered)
		if err != nil {
			return err
		}
		defer texture.Destroy()

		_, _, width, height, err := texture.Query()
		if err != nil {
			return err
		}
		dest := sdl.Rect{
			W: int32(width),
			H: int32(height),
			X: int32(panel.widthPixels * panel.XPos),
			Y: int32(panel.heightPixels * (panel.YPos + rowIndex)),
		}
		err = renderer.Copy(texture, nil, &dest)
		if err != nil {
			return err
		}
	}
	return nil
}
