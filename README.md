# gterm

gterm is intended to be an abstraction for rendering a cell/tile addressable grid of output. Ideally it is for something like a roguelike game.

To compile/run you will need sdl2, sdl_ttf, and sdl_image, and gcc/mingw

The sdl2 bindings are vendored but you will still need to ensure that the shared libraries are installed and available on your system. Go to [veandco/go-sdl2](https://github.com/veandco/go-sdl2) and follow the SDL installation instructions for your platform.

To run the example app go to the example/ directory and run `go run muncher/main.go`

## Disclaimers

This is currently completely untested and mostly wild, rapid, speculative work.

I have no idea what I actually need to build a roguelike so I am iterating aggressively.
