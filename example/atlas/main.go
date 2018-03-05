package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
)

func main() {
	window := gterm.NewWindow(300, 20, "DejaVuSansMono.ttf", 16, true)

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

	for {
		window.Refresh()
	}
}
