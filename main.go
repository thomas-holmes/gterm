package main

import (
	"fmt"
	"log"
	"path"

	"github.com/veandco/go-sdl2/sdl"
)

func getResource(root string, asset string) string {
	return path.Join("assets", root, asset)
}

func loadTexture(file string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	bmp, err := sdl.LoadBMP(file)
	if err != nil {
		return nil, err
	}
	defer bmp.Free()

	tex, err := renderer.CreateTextureFromSurface(bmp)
	if err != nil {
		return nil, err
	}
	return tex, nil
}

func renderTexture(texture *sdl.Texture, renderer *sdl.Renderer, x int, y int) error {
	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}

	dest := sdl.Rect{H: height, W: width, X: int32(x), Y: int32(y)}
	err = renderer.Copy(texture, nil, &dest)
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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.LogWarn(0, "omg a warning")

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

	background, err := loadTexture(getResource("img", "background.bmp"), renderer)
	if err != nil {
		panic(err)
	}
	defer background.Destroy()

	foreground, err := loadTexture(getResource("img", "image.bmp"), renderer)
	if err != nil {
		panic(err)
	}
	defer foreground.Destroy()

	if err = renderTexture(background, renderer, 0, 0); err != nil {
		panic(err)
	}
	if err = renderTexture(background, renderer, 320, 0); err != nil {
		panic(err)
	}
	if err = renderTexture(background, renderer, 0, 240); err != nil {
		panic(err)
	}
	if err = renderTexture(background, renderer, 320, 240); err != nil {
		panic(err)
	}

	if err = renderTexture(foreground, renderer, 50, 50); err != nil {
		panic(err)
	}

	renderer.Present()
	window.UpdateSurface()

	sdl.Delay(2500)
	renderer.Destroy()
}

func init() {
	setSdlLogger()
}
