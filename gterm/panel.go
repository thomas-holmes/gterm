package gterm

// Panel represents an application "drawable" region of the virtual terminal
type Panel struct {
	XPos     int
	YPos     int
	Width    int
	Height   int
	Z        int
	contents [][][]interface{}
	Fill     Drawable
}

func newPanel(xPos int, yPos int, width int, height int, z int, fill Drawable) Panel {
	panelContents := make([][][]interface{}, width, width)
	for i := range panelContents {
		panelContents[i] = make([][]interface{}, height, height)
	}
	return Panel{
		XPos:     xPos,
		YPos:     yPos,
		Width:    width,
		Height:   height,
		Z:        z,
		contents: panelContents,
	}
}

func (panel *Panel) getContentsAt(xPos int, yPos int) []interface{} {
	return panel.contents[xPos][yPos]
}
