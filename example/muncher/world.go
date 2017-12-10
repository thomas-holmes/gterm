package main

import (
	"log"
	"strconv"
	"time"

	"github.com/MichaelTJones/pcg"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Position struct {
	XPos int
	YPos int
}

type World struct {
	Window     *gterm.Window
	MessageBus MessageBus

	turnCount uint64

	rng *pcg.PCG64

	Player *Player

	CurrentLevel *Level

	Levels       []Level
	LevelChanged bool

	CameraCentered bool
	CameraWidth    int
	CameraHeight   int
	CameraOffsetX  int
	CameraOffsetY  int
	CameraX        int
	CameraY        int

	nextID int

	pop *PopUp

	Suspended        bool
	showScentOverlay bool

	InputBuffer []sdl.Event

	MenuStack []Menu

	GameOver bool
	QuitGame bool

	needInput bool
}

func (world *World) GetNextID() int {
	world.nextID++
	return world.nextID
}

func (world *World) PopInput() (sdl.Event, bool) {
	if len(world.InputBuffer) > 0 {
		input := world.InputBuffer[0]
		world.InputBuffer = world.InputBuffer[1:]
		return input, true
	}

	return nil, false
}

func (world *World) AddInput(event sdl.Event) {
	world.InputBuffer = append(world.InputBuffer, event)
}

func (world *World) SetCurrentLevel(index int) {
	world.CurrentLevel = &world.Levels[index]
}

func (world *World) AddLevelFromString(levelString string) {
	world.Levels = append(world.Levels, LoadFromString(levelString))

	levels := len(world.Levels)
	if levels > 1 {
		connectTwoLevels(&world.Levels[levels-2], &world.Levels[levels-1])
	}
}

func (world *World) AddLevelFromCandidate(level *CandidateLevel) {
	world.Levels = append(world.Levels, LoadCandidateLevel(level))

	levels := len(world.Levels)
	if levels > 1 {
		connectTwoLevels(&world.Levels[levels-2], &world.Levels[levels-1])
	}

}

func (world *World) GetTile(column int, row int) *Tile {
	tile := world.CurrentLevel.getTile(column, row)
	return tile
}

func (world *World) GetCreatureAtTile(xPos int, yPos int) (*Creature, bool) {
	if world.Player.X == xPos && world.Player.Y == yPos {
		return &world.Player.Creature, true
	} else if monster, ok := world.GetMonsterAtTile(xPos, yPos); ok {
		return &monster.Creature, ok
	}
	return nil, false
}

func (world *World) GetMonsterAtTile(x int, y int) (*Monster, bool) {
	if monster, ok := world.GetTile(x, y).Creature.(*Monster); ok {
		return monster, ok
	}
	return nil, false
}

func (world World) IsTileOccupied(x int, y int) bool {
	return world.GetTile(x, y).Creature != nil
}

func (world *World) IsTileMonster(x int, y int) bool {
	_, ok := world.GetMonsterAtTile(x, y)
	return ok
}

func (world *World) CanStandOnTile(column int, row int) bool {
	return !world.GetTile(column, row).IsWall() && !world.IsTileOccupied(column, row)
}

func (world *World) Suspend() {
	log.Println("Suspending world")
	world.Suspended = true
}

func (world *World) Resume() {
	log.Println("Resuming world")
	world.Suspended = false
}

func (world *World) HandleInput(event sdl.Event) {
	// TODO: Do better here, we should check keyboard/mouse/modifier/etc... state
	if world.Suspended {
		return
	}
	if event != nil {
		for _, entity := range world.CurrentLevel.Entities {
			if inputtable, ok := entity.(Inputtable); ok {
				inputtable.HandleInput(event, world)
			}
		}
	}
}

// TODO: The way this works kind of sucks. Don't have great ideas for how to make it better
// but need to come up with something. I made the "PopUp" specifically for showing the game over
// overlay and have a bunch of code that ~sort of~ suspends the game. Need to probably come up
// with a better system so they can take input, spawn other popups, etc...
func (world *World) ShowPopUp(pop PopUp) {
	pop.SetMessageBus(&world.MessageBus)
	world.pop = &pop
	world.pop.Show()
}

func (world *World) ClosePopUp() {
	if world.pop == nil {
		return
	}

	world.pop.Hide()
	world.pop.RemoveMessageBus()
	world.pop = nil
}

func (world *World) AddPlayer(player *Player) {
	world.Player = player
	if !world.CanStandOnTile(player.X, player.Y) {
		for _, t := range world.CurrentLevel.tiles {
			if !t.IsWall() && !(t.Creature != nil) {
				player.X = t.X
				player.Y = t.Y
				log.Printf("Player position adjusted to (%v,%v)", player.X, player.Y)
				break
			}
		}
	}

	// Center the camera on the player
	if world.CameraCentered {
		world.CameraX = player.X
		world.CameraY = player.Y
	}

	world.CurrentLevel.VisionMap.UpdateVision(6, world)
	world.CurrentLevel.ScentMap.UpdateScents(world)

	world.AddEntity(player)
}

// TODO: This is kinda janky. Figure out something better for this. Probably don't need
// the Renderable interface any more
func (world *World) AddRenderable(entity Entity, x int, y int) {
	world.GetTile(x, y).Creature = entity
}

func (world *World) AddEntity(e Entity) {
	e.SetIdentity(world.GetNextID())
	log.Printf("Adding entity %+v", e)

	if n, ok := e.(Notifier); ok {
		n.SetMessageBus(&world.MessageBus)
	}

	if l, ok := e.(Listener); ok {
		world.MessageBus.Subscribe(l)
	}

	// TODO: Clean this up too
	switch actual := e.(type) {
	case Renderable:
		world.AddRenderable(e, actual.XPos(), actual.YPos())
	}

	world.CurrentLevel.Entities = append(world.CurrentLevel.Entities, e)
}

// Track down these "out of bounds" errors and try to squash them. Either this function
// needs to know enough to not write to Window or the caller should not try to write outside
// of its viewport
func (world *World) RenderRuneAt(x int, y int, out rune, fColor sdl.Color, bColor sdl.Color) {
	err := world.Window.PutRune(x-world.CameraX+world.CameraOffsetX, y-world.CameraY+world.CameraOffsetY, out, fColor, bColor)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}

// Same above
func (world *World) RenderStringAt(x int, y int, out string, color sdl.Color) {
	err := world.Window.PutString(x-world.CameraX+world.CameraOffsetX, y-world.CameraY+world.CameraOffsetY, out, color)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}

func (world *World) tidyMenus() bool {
	for i := len(world.MenuStack) - 1; i >= 0; i-- {
		log.Printf("%v", world.MenuStack)
		if world.MenuStack[i].Done() {
			world.MenuStack = world.MenuStack[:i]
		}
	}

	return len(world.MenuStack) > 0
}

func (world *World) Update() bool {
	log.Printf("Updating turn [%v]", world.turnCount)

	if world.tidyMenus() {
		currentMenu := world.MenuStack[len(world.MenuStack)-1]
		if input, ok := world.PopInput(); ok {
			world.needInput = currentMenu.Update(input)
			world.tidyMenus()
		} else {
			world.needInput = true
		}
		return world.needInput
	}

	if world.CurrentLevel.NextEntity == 0 && world.CurrentLevel.NextEnergy == 0 {
		world.turnCount++
	}
	turn := world.turnCount

	world.needInput = false
	for i := world.CurrentLevel.NextEntity; i < len(world.CurrentLevel.Entities); i++ {
		e := world.CurrentLevel.Entities[i]

		log.Printf("EntityLen: %v, nextEntity: %v, nextEnergy: %v", len(world.CurrentLevel.Entities), world.CurrentLevel.NextEntity, world.CurrentLevel.NextEnergy)
		energized, isEnergized := e.(Energized)
		if isEnergized && world.CurrentLevel.NextEnergy == i {
			energized.AddEnergy(100)
			world.CurrentLevel.NextEnergy = i + 1
		}

		// TODO: Need to unravel this somehow, this is gross.
		// Maybe I put proxies for some of these methods, CanAct, NeedsInput etc... on Entity
		// but not on the interface and they can do cast and ask or return the proper default
		if a, ok := e.(Actor); ok {
			if a.CanAct() {
				if inputtable, ok := e.(Inputtable); ok {
					if inputtable.NeedsInput() {
						if input, ok := world.PopInput(); ok {
							if a.Update(turn, input, world) {
								world.CurrentLevel.NextEntity = i + 1
							} else {
								world.needInput = true
								break
							}
						} else {
							world.needInput = true
							break
						}
					} else {
						a.Update(turn, nil, world)
						world.CurrentLevel.NextEntity = i + 1
					}
				} else {
					a.Update(turn, nil, world)
					world.CurrentLevel.NextEntity = i + 1
				}
			}
		}

		if _, ok := e.(*Player); ok {
			world.CurrentLevel.VisionMap.UpdateVision(6, world)
			world.CurrentLevel.ScentMap.UpdateScents(world)
		}
		if world.LevelChanged {
			world.LevelChanged = false
			// Reset these or the player can end up in a spot where they have no energy but need input
			world.CurrentLevel.NextEnergy = 0
			world.CurrentLevel.NextEntity = 0
			log.Printf("Restarting loop after changing level")
			return true
		}
	}
	if world.CurrentLevel.NextEntity >= len(world.CurrentLevel.Entities) {
		log.Printf("Looping around")
		world.CurrentLevel.NextEntity = 0
		world.CurrentLevel.NextEnergy = 0
	}
	return world.needInput
}

func (world *World) UpdateCamera() {
	if world.CameraCentered {
		world.CameraOffsetX = 30
		world.CameraOffsetY = 15
		world.CameraX = world.Player.X
		world.CameraY = world.Player.Y
	} else {
		world.CameraOffsetX = 0
		world.CameraOffsetY = 0
		world.CameraX = 0
		world.CameraY = 0
	}
}

// Render redrwas everything!
func (world *World) Render() {
	world.UpdateCamera()
	defer timeMe(time.Now(), "World.Render.TileLoop")
	var minX, minY, maxX, maxY int
	if world.CameraCentered {
		minY, maxY = max(0, world.CameraY-(world.CameraHeight/2)), min(world.CurrentLevel.Rows, world.CameraY+(world.CameraWidth/2))
		minX, maxX = max(0, world.CameraX-(world.CameraWidth/2)), min(world.CurrentLevel.Columns, world.CameraX+(world.CameraHeight/2))
	} else {
		minY, maxY = 0, world.CurrentLevel.Rows
		minX, maxX = 0, world.CurrentLevel.Columns
	}
	log.Printf(" min x/y (%v,%v)  max x/y(%v,%v)", minX, minY, maxX, maxY)
	for row := minY; row < maxY; row++ {
		for col := minX; col < maxX; col++ {
			tile := world.GetTile(col, row)

			visibility := world.CurrentLevel.VisionMap.VisibilityAt(col, row)
			tile.Render(world, visibility)
		}
	}

	if world.pop != nil && world.pop.Shown {
		world.pop.Render(world.Window)
	}

	// Render bottom to top
	for _, m := range world.MenuStack {
		m.Render()
	}

	if world.showScentOverlay {
		world.OverlayScentMap()
	}
}

func (world *World) OverlayVisionMap() {
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			world.RenderRuneAt(x, y, []rune(strconv.Itoa(int(world.CurrentLevel.VisionMap.Map[y*world.CurrentLevel.Columns+x])))[0], Blue, gterm.NoColor)
		}
	}
}

// TODO: Something isn't quite right with these colors or how I am selecting them
// down below. The progress seems weird when I bring up the scent map display. The
// second to last band of color seems to not be what I expect.
var ScentColors = []sdl.Color{
	sdl.Color{R: 175, G: 50, B: 50, A: 1},
	sdl.Color{R: 225, G: 50, B: 25, A: 1},
	sdl.Color{R: 255, G: 0, B: 0, A: 1},
	sdl.Color{R: 100, G: 175, B: 50, A: 1},
	sdl.Color{R: 50, G: 255, B: 100, A: 1},
	sdl.Color{R: 0, G: 150, B: 175, A: 1},
	sdl.Color{R: 0, G: 50, B: 255, A: 1},
}

func (world *World) ToggleScentOverlay() {
	world.showScentOverlay = !world.showScentOverlay
}

func (world *World) OverlayScentMap() {
	turn := world.turnCount
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			scent := world.CurrentLevel.ScentMap.getScent(x, y)

			maxScent := float64(turn * 32)
			recent := float64((turn - 10) * 32)

			turnsAgo := min(int((maxScent-scent)/32), len(ScentColors)-1)
			distance := ((turn - uint64(turnsAgo)) * 32) - uint64(scent)

			bgColor := ScentColors[turnsAgo]

			bgColor.R /= 2
			bgColor.G /= 2
			bgColor.B /= 2

			if bgColor.R > bgColor.G && bgColor.R > bgColor.B {
				bgColor.R -= uint8(distance * 5)
			} else if bgColor.G > bgColor.B {
				bgColor.G -= uint8(distance * 5)
			} else {
				bgColor.B -= uint8(distance * 5)
			}
			if scent > 0 && scent > recent {
				world.RenderRuneAt(x, y, ' ', Purple, bgColor)
			}
		}
	}
}

func (world *World) RemoveEntity(entity Entity) {
	log.Printf("Removing entity")
	foundIndex := -1
	var foundEntity Entity
	for i, e := range world.CurrentLevel.Entities {
		if e.Identity() == entity.Identity() {
			foundIndex = i
			foundEntity = e
			break
		}
	}

	if foundIndex > -1 {
		world.CurrentLevel.Entities = append(world.CurrentLevel.Entities[:foundIndex], world.CurrentLevel.Entities[foundIndex+1:]...)
	}

	if renderable, ok := foundEntity.(Renderable); ok {
		world.GetTile(renderable.XPos(), renderable.YPos()).Creature = nil
	}
	log.Printf("Finisehd remove")
}

func (world *World) MoveEntity(message MoveEntityMessage) {
	oldTile := world.GetTile(message.OldX, message.OldY)
	newTile := world.GetTile(message.NewX, message.NewY)
	newTile.Creature = oldTile.Creature
	oldTile.Creature = nil
}

// TODO: This is has to perform a linear search which is less than ideal
// but I wanted ordered traversal, which you don't get with maps in go.
// Keep an eye on the performance of this.
func (world *World) GetEntity(id int) (Entity, bool) {
	for _, e := range world.CurrentLevel.Entities {
		if e.Identity() == id {
			return e, true
		}
	}
	return nil, false
}

func (world *World) Notify(message Message, data interface{}) {
	switch message {
	case ClearRegion:
		if d, ok := data.(ClearRegionMessage); ok {
			world.Window.ClearRegion(d.XPos, d.YPos, d.Width, d.Height)
		}
	case MoveEntity:
		if d, ok := data.(MoveEntityMessage); ok {
			world.MoveEntity(d)
		}
	case KillEntity:
		if d, ok := data.(KillEntityMessage); ok {
			world.RemoveEntity(d.Defender)
		}
	case PopUpShown:
		log.Println("World, PopUp Shown")
		world.Suspend()
	case PopUpHidden:
		log.Println("World, PopUp Hidden")
		world.Resume()
	case PlayerDead:
		pop := NewPopUp(10, 5, 40, 6, Red, "YOU ARE VERY DEAD", "I AM SO SORRY :(")
		world.GameOver = true
		world.ShowPopUp(pop)
	case PlayerFloorChange:
		log.Printf("Changing floors %+v", data)
		if d, ok := data.(PlayerFloorChangeMessage); ok {
			if !d.Connected {
				break
			}
			world.RemoveEntity(world.Player)
			world.Player.X = d.DestX
			world.Player.Y = d.DestY
			world.CurrentLevel = d.DestLevel
			world.LevelChanged = true
			world.AddEntity(world.Player)
		}
	case ShowMenu:
		if d, ok := data.(ShowMenuMessage); ok {
			world.MenuStack = append(world.MenuStack, d.Menu)
		}
	}
}

const (
	DefaultSeq uint64 = iota * 1000
)

func NewWorld(window *gterm.Window, centered bool, seed uint64) *World {

	world := World{
		Window:         window,
		CameraCentered: centered,
		CameraX:        0,
		CameraY:        0,
		// TODO: Width/Height should probably be some function of the window dimensions
		CameraWidth:  56,
		CameraHeight: 28,
		rng:          pcg.NewPCG64(),
	}

	world.rng.Seed(seed, DefaultSeq, seed*seed, DefaultSeq+1)

	world.MessageBus.Subscribe(&world)

	return &world
}
