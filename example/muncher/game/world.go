package game

import (
	"log"
	"strconv"
	"time"

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
	VisionMap  *VisionMap

	ScentMap ScentMap

	Player *Player

	Columns int
	Rows    int
	Tiles   []Tile

	CameraWidth  int
	CameraHeight int
	CameraX      int
	CameraY      int

	nextID int

	pop *PopUp

	Suspended        bool
	showScentOverlay bool

	InputBuffer []sdl.Event

	GameOver bool

	renderItems map[Position][]Renderable
	entities    []Entity
	nextEntity  int
	nextEnergy  int
	needInput   bool
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

func (world *World) BuildLevelFromMask(mask []int) {
	for index := range mask {
		if mask[index] == 1 {
			tile := &world.Tiles[index]

			tile.BackgroundGlyph = '#'
			tile.Wall = true
			tile.BackgroundColor = sdl.Color{R: 225, G: 225, B: 225, A: 255}
		}
	}
}

func (world *World) BuildLevel() {
	for row := 0; row < world.Rows; row++ {
		for col := 0; col < world.Columns; col++ {
			if row == 0 || (row == world.Rows-1) || (col == 0 || col == world.Columns-1) {
				world.Tiles[row*world.Columns+col].BackgroundGlyph = '#'
				world.Tiles[row*world.Columns+col].Wall = true
			}
		}
	}
}
func (world World) TileIndex(column int, row int) int {
	return row*world.Columns + column
}

func (world *World) GetTile(column int, row int) *Tile {
	tile := &world.Tiles[world.TileIndex(column, row)]
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

func (world *World) GetMonsterAtTile(column int, row int) (*Monster, bool) {
	pos := Position{XPos: column, YPos: row}
	renderItems := world.renderItems[pos]
	for _, item := range renderItems {
		if monster, ok := item.(*Monster); ok {
			return monster, ok
		}
	}
	return nil, false
}

func (world World) IsTileOccupied(column int, row int) bool {
	pos := Position{XPos: column, YPos: row}
	renderItems := world.renderItems[pos]
	return len(renderItems) > 0
}

func (world *World) IsTileMonster(column int, row int) bool {
	pos := Position{XPos: column, YPos: row}
	renderItems := world.renderItems[pos]
	isMonster := false
	for _, item := range renderItems {
		if _, ok := item.(*Monster); ok {
			isMonster = true
		}
	}
	return isMonster
}

func (world *World) CanStandOnTile(column int, row int) bool {
	return !world.GetTile(column, row).Wall && !world.IsTileOccupied(column, row)
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
		for _, entity := range world.entities {
			if inputtable, ok := entity.(Inputtable); ok {
				inputtable.HandleInput(event, world)
			}
		}
	}
}

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

func (world *World) AddRenderable(renderable Renderable) {
	pos := Position{XPos: renderable.XPos(), YPos: renderable.YPos()}
	slice := world.renderItems[pos]
	world.renderItems[pos] = append(slice, renderable)
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

	switch actual := e.(type) {
	case Renderable:
		world.AddRenderable(actual)
	}

	world.entities = append(world.entities, e)
}

func (world *World) RenderRuneAt(x int, y int, out rune, fColor sdl.Color, bColor sdl.Color) {
	err := world.Window.PutRune(x+world.CameraX, y+world.CameraY, out, fColor, bColor)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}

func (world *World) RenderStringAt(x int, y int, out string, color sdl.Color) {
	err := world.Window.PutString(x+world.CameraX, y+world.CameraY, out, color)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}

func (world *World) Update(turn int64) bool {
	world.VisionMap.UpdateVision(6, world.Player, world)
	world.ScentMap.UpdateScents(turn, *world)

	world.needInput = false
	for i := world.nextEntity; i < len(world.entities); i++ {
		e := world.entities[i]

		energized, isEnergized := e.(Energized)

		if isEnergized && world.nextEnergy == i {
			energized.AddEnergy(100)
			world.nextEnergy = i + 1
		}

		if e.CanAct() {
			if e.NeedsInput() {
				log.Printf("Found one that needs input %+v", e)
				if input, ok := world.PopInput(); ok {
					e.Update(turn, input, world)
					world.nextEntity = i + 1
				} else {
					world.needInput = true
					break
				}
			} else {
				e.Update(turn, nil, world)
				world.nextEntity = i + 1
			}
		}
	}
	if world.nextEntity == len(world.entities) {
		log.Printf("Looping around")
		world.nextEntity = 0
		world.nextEnergy = 0
	}
	return world.needInput
}

// Render redrwas everything!
func (world *World) Render(turnCount int64) {
	defer timeMe(time.Now(), "World.Render.TileLoop")
	for row := max(0, 0-world.CameraY); row < min(world.Rows, world.Rows-world.CameraY); row++ {
		for col := max(0, 0-world.CameraX); col < min(world.Columns, world.Columns-world.CameraX); col++ {
			tile := world.GetTile(col, row)

			visibility := world.VisionMap.VisibilityAt(col, row)
			tile.Render(world, visibility)
		}
	}

	if world.pop != nil && world.pop.Shown {
		world.pop.Render(world.Window)
	}

	if world.showScentOverlay {
		world.OverlayScentMap(turnCount)
	}
}

func (world *World) OverlayVisionMap() {
	blue := sdl.Color{R: 0, G: 0, B: 200, A: 255}
	for y := 0; y < world.Rows; y++ {
		for x := 0; x < world.Columns; x++ {
			world.RenderRuneAt(x, y, []rune(strconv.Itoa(int(world.VisionMap.Map[y*world.Columns+x])))[0], blue, gterm.NoColor)
		}
	}
}

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

func (world *World) OverlayScentMap(turn int64) {
	purple := sdl.Color{R: 200, G: 0, B: 200, A: 255}
	for y := 0; y < world.Rows; y++ {
		for x := 0; x < world.Columns; x++ {
			scent := world.ScentMap.getScent(x, y)

			maxScent := float32(turn * 32)
			recent := float32((turn - 10) * 32)

			turnsAgo := min(int((maxScent-scent)/32), len(ScentColors)-1)
			distance := ((turn - int64(turnsAgo)) * 32) - int64(scent)

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
				world.RenderRuneAt(x, y, ' ', purple, bgColor)
			}
		}
	}
}

func (world *World) RemoveEntity(entity Entity) {
	// delete(world.entities, entity.Identity())
	foundIndex := -1
	var foundEntity Entity
	for i, e := range world.entities {
		if e.Identity() == entity.Identity() {
			foundIndex = i
			foundEntity = e
			break
		}
	}

	if foundIndex > -1 {
		world.entities = append(world.entities[:foundIndex], world.entities[foundIndex+1:]...)
	}

	if renderable, ok := foundEntity.(Renderable); ok {
		pos := Position{XPos: renderable.XPos(), YPos: renderable.YPos()}

		slice := world.renderItems[pos]

		foundIndex := -1
		for index, candidate := range slice {
			if candidate.Identity() == renderable.Identity() {
				foundIndex = index
				break
			}
		}

		if foundIndex == -1 {
			return
		}

		world.renderItems[pos] = append(slice[:foundIndex], slice[foundIndex+1:]...)
	}
}

func (world *World) BumpCameraX(amount int) {
	world.CameraX += amount
}

func (world *World) BumpCameraY(amount int) {
	world.CameraY += amount
}

// TODO: Dedup this with move entity, maybe?
func (world *World) MovePlayer(message PlayerMoveMessage) {
	oldPos := Position{XPos: message.OldX, YPos: message.OldY}
	slice := world.renderItems[oldPos]
	foundIndex := -1
	var foundItem Renderable
	for index, item := range slice {
		if item.Identity() == message.ID {
			foundIndex = index
			foundItem = item
			break
		}
	}

	if foundIndex != -1 {
		newSlice := append(slice[:foundIndex], slice[foundIndex+1:]...)
		world.renderItems[oldPos] = newSlice
	}

	newPos := Position{XPos: message.NewX, YPos: message.NewY}
	world.BumpCameraX(-(message.NewX - oldPos.XPos))
	world.BumpCameraY(-(message.NewY - oldPos.YPos))
	newSlice := world.renderItems[newPos]
	newSlice = append(newSlice, foundItem)
	world.renderItems[newPos] = newSlice
}

// TODO: Dedup this with move player
func (world *World) MoveEntity(message MoveEntityMessage) {
	log.Printf("Moving an entity, %#v", message)
	oldPos := Position{XPos: message.OldX, YPos: message.OldY}
	slice := world.renderItems[oldPos]
	foundIndex := -1
	var foundItem Renderable
	for index, item := range slice {
		if item.Identity() == message.ID {
			foundIndex = index
			foundItem = item
			break
		}
	}

	if foundIndex != -1 {
		newSlice := append(slice[:foundIndex], slice[foundIndex+1:]...)
		world.renderItems[oldPos] = newSlice
	}

	newPos := Position{XPos: message.NewX, YPos: message.NewY}
	newSlice := world.renderItems[newPos]
	newSlice = append(newSlice, foundItem)
	world.renderItems[newPos] = newSlice
}

func (world *World) GetEntity(id int) (Entity, bool) {
	for _, e := range world.entities {
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
			// log.Printf("Got invalidation %+v", d)
			world.Window.ClearRegion(d.XPos, d.YPos, d.Width, d.Height)
		}
	case PlayerMove:
		if d, ok := data.(PlayerMoveMessage); ok {
			world.MovePlayer(d)
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
	}
}

func NewWorld(window *gterm.Window, columns int, rows int, cameraWidth int, cameraHeight int) *World {
	tiles := make([]Tile, columns*rows, columns*rows)
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			tiles[row*columns+col] = NewTile(col, row)
		}
	}
	vision := NewVisionMap(columns, rows)

	world := World{
		Window:    window,
		VisionMap: &vision,
		ScentMap:  newScentMap(columns, rows),
		Columns:   columns,
		Rows:      rows,
		Tiles:     tiles,

		// TODO: Actually do something useful with the camera settings
		CameraX:      columns / 2,
		CameraY:      rows / 2,
		CameraWidth:  cameraWidth,
		CameraHeight: cameraHeight,

		renderItems: make(map[Position][]Renderable),
	}

	world.MessageBus.Subscribe(&world)

	return &world
}
