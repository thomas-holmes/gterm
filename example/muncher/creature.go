package main

import (
	"fmt"
	"log" // Replace w/ PCG deterministic random
	"strconv"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Team int

const (
	NeutralTeam Team = iota
	PlayerTeam
	MonsterTeam
)

type Health struct {
	Current int
	Max     int
}

type Magic struct {
	Current int
	Max     int
}

type Creature struct {
	Identifiable

	IsPlayer bool

	Experience int

	RenderGlyph rune
	RenderColor sdl.Color

	Depth int

	Team

	State MonsterBehavior

	X int
	Y int

	Energy

	Inventory

	Equipment

	HP Health
	MP Magic

	Spells []Spell

	Level int

	Name string

	Messaging
}

func (c Creature) CanAct() bool {
	return c.currentEnergy >= 100
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c *Creature) Damage(damage int) {
	log.Printf("%v is Taking damage of %v", *c, damage)
	c.HP.Current = max(0, c.HP.Current-damage)

	if c.IsPlayer {
		c.Broadcast(PlayerUpdate, nil)
		if c.HP.Current == 0 {
			c.Broadcast(PlayerDead, nil)
		}
	}
}

func (c *Creature) TryMove(newX int, newY int, world *World) (MoveResult, interface{}) {
	if world.CurrentLevel.CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.CurrentLevel.GetCreatureAtTile(newX, newY); ok {
		log.Printf("Got a creature in TryMove, %+v", *defender)
		if c.Team != defender.Team {
			a, aOk := world.GetEntity(c.ID)
			d, dOk := world.GetEntity(defender.ID)
			if aOk && dOk {
				return MoveIsEnemy, MoveEnemy{Attacker: a, Defender: d}
			}
		}
	}

	return MoveIsInvalid, nil
}

func NewCreature(level int, maxHP int) Creature {
	return Creature{
		Level: level,
		Team:  NeutralTeam,
		Energy: Energy{
			currentEnergy: 100,
			maxEnergy:     100,
		},
		HP:        Health{Current: 5, Max: 5},
		Equipment: NewEquipment(),
	}
}

func NewPlayer() Creature {
	player := NewCreature(1, 5)
	player.Team = PlayerTeam
	player.RenderGlyph = '@'
	player.RenderColor = Red
	player.IsPlayer = true
	player.Spells = DefaultSpells

	return player
}

func NewMonster(xPos int, yPos int, level int, hp int) Creature {
	monster := NewCreature(level, hp)

	monster.X = xPos
	monster.Y = yPos
	monster.Team = MonsterTeam
	monster.RenderColor = Green
	monster.RenderGlyph = []rune(strconv.Itoa(monster.Level))[0]

	return monster
}

func (player *Creature) LevelUp() {
	player.Experience = max(0, player.Experience-player.Level)
	player.Level++
	player.HP.Max = player.HP.Max + max(1, int(float64(player.HP.Max)*0.1))
	player.HP.Current = player.HP.Max
	player.Broadcast(GameLogAppend, GameLogAppendMessage{[]string{fmt.Sprintf("You are now level %v", player.Level)}})
}

func (player *Creature) GainExp(exp int) {
	player.Experience += exp
	if player.Experience >= player.Level {
		player.LevelUp()
		player.Broadcast(PlayerUpdate, nil)
	}
}

func (player Creature) HealthPercentage() float32 {
	current := float32(player.HP.Current)
	max := float32(player.HP.Max)
	return current / max
}

func (player *Creature) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp

	player.Broadcast(PlayerUpdate, nil)
}

func (player *Creature) PickupItem(world *World) bool {
	tile := world.CurrentLevel.GetTile(player.X, player.Y)
	if tile.Item == nil {
		return false
	}

	player.Items = append(player.Items, tile.Item)
	tile.Item = nil
	return true
}

// Update returns true if an action that would constitute advancing the turn took place
func (creature *Creature) Update(turn uint64, input InputEvent, world *World) bool {
	success := false
	if creature.IsPlayer {
		success = creature.HandleInput(input, world)
	} else {
		success = creature.Pursue(turn, world)
	}

	if success {
		creature.currentEnergy -= 100
		return true
	}

	return false
}

func (creature *Creature) CastSpell(spell Spell, world *World) {
	menu := &SpellTargeting{X: 0, Y: 0, W: 0, H: 0, TargetX: creature.X, TargetY: creature.Y, World: world, Spell: spell}
	creature.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
}

// HandleInput updates player position based on user input
func (player *Creature) HandleInput(input InputEvent, world *World) bool {
	newX := player.X
	newY := player.Y

	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_COMMA:
			if input.Keymod&sdl.KMOD_SHIFT > 0 {
				tile := world.CurrentLevel.GetTile(player.X, player.Y)
				if tile.TileKind == UpStair {
					if stair, ok := world.CurrentLevel.getStair(player.X, player.Y); ok {
						player.Broadcast(PlayerFloorChange, PlayerFloorChangeMessage{
							Stair: stair,
						})
					} else {
						return false
					}
				}
			}
			return false
		case sdl.K_PERIOD:
			if input.Keymod&sdl.KMOD_SHIFT > 0 {
				tile := world.CurrentLevel.GetTile(player.X, player.Y)
				if tile.TileKind == DownStair {
					if stair, ok := world.CurrentLevel.getStair(player.X, player.Y); ok {
						player.Broadcast(PlayerFloorChange, PlayerFloorChangeMessage{
							Stair: stair,
						})
					} else {
						return false
					}
				}
			}
			// Period returns true because it means "wait"
			return true
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
			return false
		case sdl.K_2:
			player.Heal(1)
			return false
		case sdl.K_g:
			return player.PickupItem(world)
		case sdl.K_i:
			menu := &InventoryPop{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4, Inventory: player.Inventory}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_e:
			menu := &EquipmentPop{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4, Player: player}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_x:
			menu := &InspectionPop{X: 60, Y: 20, W: 30, H: 5, World: world, InspectX: player.X, InspectY: player.Y}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_z:
			menu := &SpellPop{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4, World: world}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_m:
			player.Broadcast(ShowFullGameLog, nil)
			return false
		case sdl.K_ESCAPE:
			world.GameOver = true
			world.QuitGame = true
			return true
		default:
			return false
		}

		if newX != player.X || newY != player.Y {
			result, data := player.TryMove(newX, newY, world)
			switch result {
			case MoveIsInvalid:
				return false
			case MoveIsSuccess:
				oldX := player.X
				oldY := player.Y
				player.X = newX
				player.Y = newY
				player.Broadcast(MoveEntity, MoveEntityMessage{ID: player.ID, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					player.Broadcast(AttackEntity, AttackEntityMesasge{
						Attacker: data.Attacker,
						Defender: data.Defender,
					})
				}
			}
		}
		return true
	}
	return false
}

func (creature *Creature) Notify(message Message, data interface{}) {
	if !creature.IsPlayer {
		return
	}
	switch message {
	case KillEntity:
		if d, ok := data.(KillEntityMessage); ok {
			attacker, ok := d.Attacker.(*Creature)
			if !ok {
				return
			}
			defender, ok := d.Defender.(*Creature)
			if !ok {
				return
			}

			if defender.ID == creature.ID {
				creature.Broadcast(PlayerDead, nil)
				return
			}
			if attacker.ID != creature.ID {
				return
			}
			attacker.GainExp(defender.Level)
			if creature.Level > defender.Level {
				//creature.GainExp((defender.Level + 1) / 4)
			} else {
				//creature.GainExp((defender.Level + 1) / 2)
			}
		}
	case EquipItem:
		if d, ok := data.(EquipItemMessage); ok {
			creature.Equipment.Weapon = d.Item // This is super low effort, but should work?
		}
	}
}

func (c *Creature) NeedsInput() bool {
	return c.IsPlayer
}

// SetColor updates the render color of the player
func (player *Creature) SetColor(color sdl.Color) {
	player.RenderColor = color
}

func (creature *Creature) Render(world *World) {
	world.RenderRuneAt(creature.X, creature.Y, creature.RenderGlyph, creature.RenderColor, gterm.NoColor)
}

func (monster *Creature) Pursue(turn uint64, world *World) bool {
	if world.CurrentLevel.VisionMap.VisibilityAt(monster.X, monster.Y) == Visible {
		monster.State = Pursuing
	}

	if monster.State != Pursuing {
		return true
	}

	scent := world.CurrentLevel.ScentMap

	// TODO: Maybe short circuit tracking here and just attack the player instead
	// if in range?
	candidates := scent.track(turn, monster.X, monster.Y)

	if len(candidates) > 0 {
		for _, choice := range candidates {
			result, data := monster.TryMove(choice.X, choice.Y, world)
			switch result {
			case MoveIsInvalid:
				continue
			case MoveIsSuccess:
				oldX := monster.X
				oldY := monster.Y
				monster.X = choice.X
				monster.Y = choice.Y
				monster.Broadcast(MoveEntity, MoveEntityMessage{
					ID:   monster.ID,
					OldX: oldX,
					OldY: oldY,
					NewX: choice.X,
					NewY: choice.Y,
				})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					monster.Broadcast(AttackEntity, AttackEntityMesasge{
						Attacker: data.Attacker,
						Defender: data.Defender,
					})
				}
			}
			return true
		}
	} else {
		return true
	}
	return false
}

type MoveResult int

const (
	MoveIsInvalid MoveResult = iota
	MoveIsSuccess
	MoveIsEnemy
)

type MoveEnemy struct {
	Attacker Entity
	Defender Entity
}

type MonsterBehavior int

const (
	Idle MonsterBehavior = iota
	Pursuing
)
