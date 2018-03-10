package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window := gterm.NewWindow(80, 80, "fonts/cp437_12x12.png", 12, 12, true)

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to init window", err)
	}

	/*
		row := 0
		col := 0
		for i := 32; i <= 255; i++ {
			if (i-0x2500)%50 == 0 {
				row++
				col = 0
			}

			window.PutRune(col, row, rune(i), gterm.White, gterm.NoColor)
			col++
		}
	*/

	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	black := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	green := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	blue := sdl.Color{R: 0, G: 0, B: 255, A: 255}
	red := sdl.Color{R: 255, G: 0, B: 0, A: 255}

	window.PutRune(40, 40, 'X', white, black)
	window.PutRune(41, 40, rune(16), white, black)
	window.PutRune(42, 40, ' ', white, black)
	window.PutRune(43, 40, rune('▓'), white, black)

	window.PutRune(40, 41, 'X', black, white)
	window.PutRune(41, 41, rune(16), black, white)
	window.PutRune(42, 41, ' ', black, white)
	window.PutRune(43, 41, '▓', black, white)

	window.PutRune(40, 42, 'X', green, blue)
	window.PutRune(41, 42, rune(16), green, blue)
	window.PutRune(42, 42, ' ', green, blue)
	window.PutRune(43, 42, '▓', green, blue)

	window.PutRune(40, 43, 'X', red, white)
	window.PutRune(41, 43, rune(16), red, white)
	window.PutRune(42, 43, ' ', red, white)
	window.PutRune(43, 43, '▓', red, white)

	window.PutRune(40, 44, 'X', green, white)
	window.PutRune(41, 44, rune(16), green, white)
	window.PutRune(42, 44, ' ', green, white)
	window.PutRune(43, 44, '▓', green, white)

	window.PutRune(50, 50, rune('█'), green, white)
	window.PutRune(50, 51, rune('█'), green, white)
	window.PutRune(50, 52, rune('█'), green, white)
	window.PutRune(50, 53, rune('█'), green, white)
	window.PutRune(50, 54, rune('█'), green, white)
	window.PutRune(51, 50, rune('█'), green, white)
	window.PutRune(52, 50, rune('█'), green, white)
	window.PutRune(53, 50, rune('█'), green, white)
	window.PutRune(54, 50, rune('█'), green, white)

	window.PutRune(30, 60, rune(0x2500), green, white)
	window.PutRune(31, 60, rune(0x2500), green, white)
	window.PutRune(32, 60, rune(0x2500), green, white)
	window.PutRune(33, 60, rune(0x2500), green, white)
	window.PutRune(34, 60, rune(0x2500), green, white)
	window.PutRune(35, 60, rune(0x2500), green, white)
	window.PutRune(36, 60, rune(0x2500), green, white)

	window.PutRune(30, 61, rune(0x2500), green, black)
	window.PutRune(31, 61, rune(0x2500), green, black)
	window.PutRune(32, 61, rune(0x2500), green, black)
	window.PutRune(33, 61, rune(0x2500), green, black)
	window.PutRune(34, 61, rune(0x2500), green, black)
	window.PutRune(35, 61, rune(0x2500), green, black)
	window.PutRune(36, 61, rune(0x2500), green, black)

	for {
		window.Refresh()
	}
}
