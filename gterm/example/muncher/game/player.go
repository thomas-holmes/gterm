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

type Health struct {
	Current int
	Max     int
}

// Player pepresents the player
type Player struct {
	World       *World
	HP          Health
	Name        string
	XPos        int
	YPos        int
	RenderGlyph string
	RenderColor sdl.Color
}

func NewPlayer(world *World, xPos int, yPos int) Player {
	player := Player{
		World:       world,
		HP:          Health{Current: 5, Max: 5},
		RenderGlyph: "@",
		RenderColor: sdl.Color{R: 255, G: 0, B: 0, A: 0},
		XPos:        xPos,
		YPos:        yPos,
	}

	world.AddRenderableToTile(xPos, yPos, &player)

	return player
}
func (player Player) HealthPercentage() float32 {
	current := float32(player.HP.Current)
	max := float32(player.HP.Max)
	return current / max
}

func (player *Player) Damage(amount int) {
	amount = max(amount, 0)

	newHp := max(player.HP.Current-amount, 0)
	player.HP.Current = newHp

	player.World.MessageBus.Broadcast(PlayerUpdate, nil)
}

func (player *Player) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp

	player.World.MessageBus.Broadcast(PlayerUpdate, nil)
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
		case sdl.K_1:
			player.Damage(1)
		case sdl.K_2:
			player.Heal(1)
		}
	}
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.RenderColor = color
}
