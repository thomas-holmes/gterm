package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const tileSize = 40
const windowWidth = 640
const windowHeight = 480

func getResource(root string, asset string) string {
	return path.Join("assets", root, asset)
}

func loadTexture(file string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	texture, err := img.LoadTexture(renderer, file)
	if err != nil {
		return nil, err
	}

	return texture, nil
}

func renderTexture(texture *sdl.Texture, renderer *sdl.Renderer, x int, y int) error {
	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}
	// This is silly, casting between int and int32, oh well
	return renderTextureScaled(texture, renderer, x, y, int(width), int(height))
}

func renderTextureScaled(texture *sdl.Texture, renderer *sdl.Renderer, x int, y int, w int, h int) error {
	dest := sdl.Rect{H: int32(h), W: int32(w), X: int32(x), Y: int32(y)}
	err := renderer.Copy(texture, nil, &dest)
	if err != nil {
		return err
	}
	return nil
}

func setSdlLogger() {
	/*
		LOG_PRIORITY_VERBOSE
		LOG_PRIORITY_DEBUG
		LOG_PRIORITY_INFO
		LOG_PRIORITY_WARN
		LOG_PRIORITY_ERROR
		LOG_PRIORITY_CRITICAL
	*/
	sdl.LogSetOutputFunction(func(data interface{}, cat int, pri sdl.LogPriority, message string) {
		priArray := [6]string{"VERBOSE", "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"}
		log.Println("[SDL]", fmt.Sprintf("[%v]", priArray[pri-1]), message)
	}, nil)
}

func tileBackground(background *sdl.Texture, renderer *sdl.Renderer, width int, height int, tileSize int) error {
	xTiles := width / tileSize
	yTiles := height / tileSize

	for tile := 0; tile < xTiles*yTiles; tile++ {
		x := tile % xTiles
		y := tile / xTiles

		err := renderTextureScaled(background, renderer, x*tileSize, y*tileSize, tileSize, tileSize)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	background, err := loadTexture(getResource("img", "background.png"), renderer)
	if err != nil {
		panic(err)
	}
	defer background.Destroy()

	tileBackground(background, renderer, windowWidth, windowHeight, tileSize)

	foreground, err := loadTexture(getResource("img", "image.png"), renderer)
	if err != nil {
		panic(err)
	}
	defer foreground.Destroy()

	renderTexture(foreground, renderer, 200, 200)

	quit := false
	frameCount := 0
	go startFpsCounter(&frameCount)
	for quit == false {
		event := sdl.PollEvent()
		if event != nil {
			log.Println(fmt.Sprintf("got event %#v", event))
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				quit = true
			case *sdl.MouseButtonEvent:
				quit = true
			}
		}
		frameCount++

		renderer.Present()
		window.UpdateSurface()
	}
}

func startFpsCounter(i *int) {
	timer := time.Tick(1 * time.Second)
	go func() {
		for {
			<-timer
			log.Println(*i)
			*i = 0
		}
	}()
}

func init() {
	setSdlLogger()
}
