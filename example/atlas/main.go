package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window := gterm.NewWindow(80, 40, "fonts/cp437_12x12.png", 12, 12, false)

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to init window", err)
	}

	w := window.SdlWindow
	_ = w
	r := window.SdlRenderer
	_ = r

	w.SetResizable(true)

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

	window.PutRune(40, 10, 'X', white, black)
	window.PutRune(41, 10, rune(16), white, black)
	window.PutRune(42, 10, ' ', white, black)
	window.PutRune(43, 10, rune('▓'), white, black)

	window.PutRune(40, 11, 'X', black, white)
	window.PutRune(41, 11, rune(16), black, white)
	window.PutRune(42, 11, ' ', black, white)
	window.PutRune(43, 11, '▓', black, white)

	window.PutRune(40, 12, 'X', green, blue)
	window.PutRune(41, 12, rune(16), green, blue)
	window.PutRune(42, 12, ' ', green, blue)
	window.PutRune(43, 12, '▓', green, blue)

	window.PutRune(40, 13, 'X', red, white)
	window.PutRune(41, 13, rune(16), red, white)
	window.PutRune(42, 13, ' ', red, white)
	window.PutRune(43, 13, '▓', red, white)

	window.PutRune(40, 14, 'X', green, white)
	window.PutRune(41, 14, rune(16), green, white)
	window.PutRune(42, 14, ' ', green, white)
	window.PutRune(43, 14, '▓', green, white)

	window.PutRune(50, 20, rune('█'), green, white)
	window.PutRune(50, 21, rune('█'), green, white)
	window.PutRune(50, 22, rune('█'), green, white)
	window.PutRune(50, 23, rune('█'), green, white)
	window.PutRune(50, 24, rune('█'), green, white)
	window.PutRune(51, 20, rune('█'), green, white)
	window.PutRune(52, 20, rune('█'), green, white)
	window.PutRune(53, 20, rune('█'), green, white)
	window.PutRune(54, 20, rune('█'), green, white)

	window.PutRune(30, 30, rune(0x2500), green, white)
	window.PutRune(31, 30, rune(0x2500), green, white)
	window.PutRune(32, 30, rune(0x2500), green, white)
	window.PutRune(33, 30, rune(0x2500), green, white)
	window.PutRune(34, 30, rune(0x2500), green, white)
	window.PutRune(35, 30, rune(0x2500), green, white)
	window.PutRune(36, 30, rune(0x2500), green, white)

	window.PutRune(30, 31, rune(0x2500), green, black)
	window.PutRune(31, 31, rune(0x2500), green, black)
	window.PutRune(32, 31, rune(0x2500), green, black)
	window.PutRune(33, 31, rune(0x2500), green, black)
	window.PutRune(34, 31, rune(0x2500), green, black)
	window.PutRune(35, 31, rune(0x2500), green, black)
	window.PutRune(36, 31, rune(0x2500), green, black)

	frames := 0
	last := sdl.GetTicks()
	for {
		frames++
		now := sdl.GetTicks()
		if now-last > 1000 {
			log.Println("FPS:", frames)
			last = now
			frames = 0
		}
		for {
			e := sdl.PollEvent()
			switch v := e.(type) {
			case *sdl.WindowEvent:
				if v.Event == sdl.WINDOWEVENT_RESIZED {
					//					window.UpdateSize()
				}
			case *sdl.KeyDownEvent:
				switch v.Keysym.Sym {
				case sdl.K_1:
					window.ChangeFont("fonts/cp437_8x8.png", 8, 8)
				case sdl.K_2:
					window.ChangeFont("fonts/cp437_12x12.png", 12, 12)
				case sdl.K_3:
					window.ChangeFont("fonts/cp437_16x16.png", 16, 16)
				}
			}
			if e == nil {
				break
			}
		}
		window.Refresh()
	}
}
