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

	renderItems map[Position][]Renderable
	entities    map[int]Entity
}

func (world *World) GetNextID() int {
	world.nextID++
	return world.nextID
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

func (world *World) GetMonsterAtTile(column int, row int) *Monster {
	pos := Position{XPos: column, YPos: row}
	renderItems := world.renderItems[pos]
	for _, item := range renderItems {
		if monster, ok := item.(*Monster); ok {
			return monster
		}
	}
	return nil
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
	e.SetID(world.GetNextID())
	log.Printf("Adding entity %+v", e)

	if n, ok := e.(Notifier); ok {
		n.SetMessageBus(&world.MessageBus)
	}

	switch actual := e.(type) {
	case Renderable:
		world.AddRenderable(actual)
	}

	world.entities[e.ID()] = e
}

func (world *World) RenderRuneAt(x int, y int, out rune, fColor sdl.Color, bColor sdl.Color) {
	err := world.Window.PutRune(x+world.CameraX, y+world.CameraY, out, fColor, bColor)
	if err != nil {
		// log.Printf("Out of bounds %s", err)
		// log.Panicf("Could not add cell", err)
	}
}

func (world *World) RenderStringAt(x int, y int, out string, color sdl.Color) {
	err := world.Window.PutString(x+world.CameraX, y+world.CameraY, out, color)
	if err != nil {
		// log.Printf("Out of bounds %s", err)
		// log.Panicf("Could not add cell", err)
	}
}

func (world *World) Update(turn int64) {
	world.VisionMap.UpdateVision(6, world.Player, world)
	world.ScentMap.UpdateScents(turn, *world.VisionMap, *world.Player)

	for _, e := range world.entities {
		switch m := e.(type) {
		case *Monster:
			log.Printf("Got a monster, %v", m)
			m.Pursue(turn, world.ScentMap)
		}
	}
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
	sdl.Color{R: 175, G: 50, B: 50, A: 55},
	sdl.Color{R: 225, G: 25, B: 25, A: 55},
	sdl.Color{R: 255, G: 0, B: 0, A: 55},
	sdl.Color{R: 100, G: 175, B: 50, A: 55},
	sdl.Color{R: 50, G: 255, B: 100, A: 55},
	sdl.Color{R: 0, G: 150, B: 175, A: 55},
	sdl.Color{R: 0, G: 50, B: 255, A: 55},
}

func (world *World) ToggleScentOverlay() {
	world.showScentOverlay = !world.showScentOverlay
}

func (world *World) OverlayScentMap(turn int64) {
	purple := sdl.Color{R: 200, G: 0, B: 200, A: 255}
	for y := 0; y < world.Rows; y++ {
		for x := 0; x < world.Columns; x++ {
			scent := world.ScentMap.getScent(x, y)

			if x == 5 && y == 5 {
				maxScent := turn * 32
				log.Printf("max scent %v, found scent: %v", maxScent, scent)
			}

			bgColorIndex := min64((turn*32)-scent, 6)
			// log.Printf("scent %v, calculated %v", scent, bgColorIndex)

			bgColor := ScentColors[bgColorIndex]

			if scent > (turn-10)*32 {
				world.RenderRuneAt(x, y, ' ', purple, bgColor)
			}
		}
	}
}

func (world *World) RemoveEntity(entity Entity) {
	delete(world.entities, entity.ID())

	if renderable, ok := entity.(Renderable); ok {
		pos := Position{XPos: renderable.XPos(), YPos: renderable.YPos()}

		slice := world.renderItems[pos]

		foundIndex := -1
		for index, candidate := range slice {
			if candidate.ID() == renderable.ID() {
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
		if item.ID() == message.ID {
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
		if item.ID() == message.ID {
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

func (world *World) Notify(message Message, data interface{}) {
	switch message {
	case TileInvalidated:
		if d, ok := data.(TileInvalidatedMessage); ok {
			// log.Printf("Got invalidation %+v", d)
			world.Window.ClearCell(d.XPos, d.YPos)
		}
	case PlayerMove:
		if d, ok := data.(PlayerMoveMessage); ok {
			world.MovePlayer(d)
		}
	case MoveEntity:
		if d, ok := data.(MoveEntityMessage); ok {
			world.MoveEntity(d)
		}
	case KillMonster:
		if d, ok := data.(KillMonsterMessage); ok {
			monster := world.entities[d.ID]
			if m, ok := monster.(*Monster); ok {
				log.Println("remove an entity", m)
				world.RemoveEntity(m)
				world.MessageBus.Broadcast(TileInvalidated, TileInvalidatedMessage{XPos: m.XPos(), YPos: m.YPos()})
			}
		}
	case PopUpShown:
		log.Println("World, PopUp Shown")
		world.Suspend()
	case PopUpHidden:
		log.Println("World, PopUp Hidden")
		world.Resume()
	case PlayerDead:
		pop := NewPopUp(10, 5, 40, 6, Red, "YOU ARE VERY DEAD", "I AM SO SORRY :(")
		world.ShowPopUp(pop)
	}
}

func NewWorld(window *gterm.Window, columns int, rows int) *World {
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
		CameraX:      0,
		CameraY:      0,
		CameraWidth:  columns,
		CameraHeight: rows,

		renderItems: make(map[Position][]Renderable),
		entities:    make(map[int]Entity),
	}

	world.MessageBus.Subscribe(&world)

	return &world
}
