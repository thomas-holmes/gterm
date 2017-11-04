package main

import "github.com/thomas-holmes/sneaker/gterm"
import "time"

func main() {
	window := gterm.NewWindow(40, 40, 16)

	window.Init()

	time.Sleep(2 * time.Second)
}
