package game

import "github.com/veandco/go-sdl2/sdl"

func (player *Player) UpdatePosition(xPos int, yPos int) {
	player.World.ClearTile(player.XPos, player.YPos)
	if xPos >= 0 && xPos < player.World.Columns &&
		yPos >= 0 && yPos < player.World.Rows &&
		player.World.CanStandOnTile(xPos, yPos) {
		player.XPos = xPos
		player.YPos = yPos
	}
	player.World.AddRenderableToTile(player.XPos, player.YPos, player)
	player.World.MessageBus.Broadcast(PlayerUpdate, nil)
}

func (player *Player) Render(world *World) {
	world.Window.AddToCell(player.XPos, player.YPos, player.RenderGlyph, player.RenderColor)
}

// Player pepresents the player
type Player struct {
	World       *World
	Name        string
	XPos        int
	YPos        int
	RenderGlyph string
	RenderColor sdl.Color
}

func NewPlayer(world *World, xPos int, yPos int) Player {
	player := Player{
		World:       world,
		RenderGlyph: "@",
		RenderColor: sdl.Color{R: 255, G: 0, B: 0, A: 0},
		XPos:        xPos,
		YPos:        yPos,
	}

	world.AddRenderableToTile(xPos, yPos, &player)

	return player
}

// HandleInput updates player position based on user input
func (player *Player) HandleInput(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_h:
			player.UpdatePosition(player.XPos-1, player.YPos)
		case sdl.K_j:
			player.UpdatePosition(player.XPos, player.YPos+1)
		case sdl.K_k:
			player.UpdatePosition(player.XPos, player.YPos-1)
		case sdl.K_l:
			player.UpdatePosition(player.XPos+1, player.YPos)
		case sdl.K_b:
			player.UpdatePosition(player.XPos-1, player.YPos+1)
		case sdl.K_n:
			player.UpdatePosition(player.XPos+1, player.YPos+1)
		case sdl.K_y:
			player.UpdatePosition(player.XPos-1, player.YPos-1)
		case sdl.K_u:
			player.UpdatePosition(player.XPos+1, player.YPos-1)
		}
	}
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.RenderColor = color
}
