package main

import (
	"math"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/thomas-holmes/gterm"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func max64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func maxu64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
func maxf32(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
func min64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func minu64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
func minf32(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -1 * a
}

func timeMe(start time.Time, name string) {
	// elapsed := time.Since(start)
	// log.Printf("%v %s", name, elapsed)

}

func squareDistance(x0 int, y0 int, x1 int, y1 int) int {
	xDelta := abs(x0 - x1)
	yDelta := abs(y0 - y1)

	if xDelta > yDelta {
		return xDelta
	}

	return yDelta
}

func euclideanDistance(x0 int, y0 int, x1 int, y1 int) float64 {
	x := x1 - x0
	y := y1 - y0
	return math.Sqrt(float64(x*x) + float64(y*y))
}

func putWrappedText(window *gterm.Window, content string, x int, y int, firstIndent int, afterIndent int, width int, color sdl.Color) int {
	offsetX := x + firstIndent
	offsetY := y

	for {
		if len(content) == 0 {
			break
		}
		maxLength := width - (offsetX - x)
		cut := min(len(content), maxLength)
		printable := content[:cut]
		lastSpace := strings.LastIndexAny(printable, " ")
		if printable != content && lastSpace > -1 {
			printable = printable[:lastSpace]
			content = strings.TrimSpace(content[lastSpace:])
		} else {
			content = strings.TrimSpace(content[cut:])
		}
		window.PutString(offsetX, offsetY, printable, color)
		offsetY++
		offsetX = x + afterIndent
	}

	return offsetY - y
}
