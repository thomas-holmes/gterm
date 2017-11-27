package game

import (
	"log"
	"math/rand"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

func getRandomColor() sdl.Color {
	return sdl.Color{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func (player *Player) Render(world *World) {
	playerBg := gterm.NoColor // playerBg := getRandomColor()
	world.RenderRuneAt(player.X, player.Y, player.RenderGlyph, player.RenderColor, playerBg)
}

type Health struct {
	Current int
	Max     int
}

// Player pepresents the player
type Player struct {
	Experience  int
	RenderGlyph rune
	RenderColor sdl.Color

	Creature

	Messaging
}

func (player *Player) LevelUp() {
	player.Experience -= player.Level
	player.Level++
	player.HP.Max = int(float32(player.HP.Max) * 1.5)
	player.HP.Current = player.HP.Max
}

func (player *Player) GainExp(exp int) {
	player.Experience += exp
	log.Println("Got some exp", exp)
	if player.Experience >= player.Level {
		player.LevelUp()
		player.Broadcast(PlayerUpdate, nil)
	}
}

func NewPlayer(xPos int, yPos int) Player {
	player := Player{
		RenderGlyph: '@',
		RenderColor: sdl.Color{R: 255, G: 0, B: 0, A: 0},
		Creature: Creature{
			Level: 1,
			HP:    Health{Current: 5, Max: 5},
			X:     xPos,
			Y:     yPos,
		},
	}

	log.Printf("Made a player, %#v", player)
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

	player.Broadcast(PlayerUpdate, nil)
	if newHp == 0 {
		player.Broadcast(PlayerDead, nil)
	}

}

func (player *Player) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp

	player.Broadcast(PlayerUpdate, nil)
}

// HandleInput updates player position based on user input
func (player *Player) HandleInput(event sdl.Event, world *World) {
	newX := player.X
	newY := player.Y

	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_h:
			newX = player.X - 1
		case sdl.K_j:
			newY = player.Y + 1
		case sdl.K_k:
			newY = player.Y - 1
		case sdl.K_l:
			newX = player.X + 1
		case sdl.K_b:
			newX, newY = player.X-1, player.Y+1
		case sdl.K_n:
			newX, newY = player.X+1, player.Y+1
		case sdl.K_y:
			newX, newY = player.X-1, player.Y-1
		case sdl.K_u:
			newX, newY = player.X+1, player.Y-1
		case sdl.K_1:
			player.Damage(1)
		case sdl.K_2:
			player.Heal(1)
		}

		if newX != player.X || newY != player.Y {
			result, data := player.TryMove(newX, newY, *world)
			switch result {
			case MoveIsInvalid:
			case MoveIsSuccess:
				oldX := player.X
				oldY := player.Y
				player.X = newX
				player.Y = newY
				player.Broadcast(PlayerMove, PlayerMoveMessage{ID: player.ID, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					player.Broadcast(CreatureAttack, CreatureAttackMessage{
						Attacker: data.Attacker,
						Defender: data.Defender,
					})
				}
			}
		}
	}
}

func (player *Player) Notify(message Message, data interface{}) {
	switch message {
	case KillEntity:
		if d, ok := data.(KillEntityMessage); ok {
			if d.Defender.ID == player.ID {
				player.Broadcast(PlayerDead, nil)
				return
			}
			if d.Attacker.ID != player.ID {
				return
			}
			if player.Level > d.Defender.Level {
				player.GainExp((d.Defender.Level + 1) / 4)
			} else {
				player.GainExp((d.Defender.Level + 1) / 2)
			}
		}
	}
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.RenderColor = color
}
