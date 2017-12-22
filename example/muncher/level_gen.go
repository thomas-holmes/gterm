package main

import (
	"time"

	"github.com/MichaelTJones/pcg"
)

type Room struct {
	ID          int
	X           int
	Y           int
	W           int
	H           int
	connectedTo []int
}

type CandidateTile struct {
	TileKind TileKind
	Item     *Item
}

type CandidateLevel struct {
	rng *pcg.PCG64

	W int
	H int

	nextRoomID int

	flags LevelGenFlag

	rooms map[int]*Room
	tiles []CandidateTile
}

const (
	MinRoomWidth      = 3
	MinRoomHeight     = 3
	MaxRoomWidth      = 20
	MaxRoomHeight     = 20
	MaxRoomIterations = 200
	MaxItemPlacement  = 50 // Let's go overboard at first.
)

func (level *CandidateLevel) genNextRoomID() int {
	id := level.nextRoomID
	level.nextRoomID++
	return id
}

func (level *CandidateLevel) tryAddRandomRoom() {
	widthBound := uint64(MaxRoomWidth - MinRoomWidth)
	heightBound := uint64(MaxRoomHeight - MinRoomHeight)

	randomWidth := int(level.rng.Bounded(widthBound)) + MinRoomWidth
	randomHeight := int(level.rng.Bounded(heightBound)) + MinRoomHeight

	topLeftX := int(level.rng.Bounded(uint64(level.W-randomWidth-1))) + 1
	topLeftY := int(level.rng.Bounded(uint64(level.H-randomHeight-1))) + 1

	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			if level.tiles[y*level.W+x].TileKind != Wall {
				// Just quit if we run into a non-wall feature
				return
			}
		}
	}

	// We can place our room, so make it all floors.
	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			level.tiles[y*level.W+x].TileKind = Floor
		}
	}

	room := &Room{
		ID: level.genNextRoomID(),
		X:  topLeftX,
		Y:  topLeftY,
		W:  randomWidth,
		H:  randomHeight,
	}

	level.rooms[room.ID] = room

	if room.ID > 0 {
		level.connectRooms(room.ID, room.ID-1)
	}
}

func (level *CandidateLevel) addRooms() {
	for i := 0; i < MaxRoomIterations; i++ {
		level.tryAddRandomRoom()
	}
}

func (room *Room) chooseTopWall(rng *pcg.PCG64) (int, int) {
	x := int(rng.Bounded(uint64(room.W))) + room.X
	y := room.Y
	return x, y
}
func (room *Room) chooseRightWall(rng *pcg.PCG64) (int, int) {
	x := room.X + room.W - 1
	y := int(rng.Bounded(uint64(room.H))) + room.Y
	return x, y
}
func (room *Room) chooseBottomWall(rng *pcg.PCG64) (int, int) {
	x := int(rng.Bounded(uint64(room.W))) + room.X
	y := room.Y + room.H - 1
	return x, y
}
func (room *Room) chooseLeftWall(rng *pcg.PCG64) (int, int) {
	x := room.X
	y := int(rng.Bounded(uint64(room.H))) + room.Y
	return x, y
}

func (level *CandidateLevel) connectRooms(roomId1 int, roomId2 int) {
	room1 := level.rooms[roomId1]
	room2 := level.rooms[roomId2]

	const (
		Top int = iota
		Right
		Bottom
		Left
	)

	yDelta := room1.Y - room2.Y
	xDelta := room1.X - room2.X

	var r1x, r1y, r2x, r2y int

	if yDelta > 0 { // Below
		if xDelta > 0 { // Right
			if xDelta > yDelta { // further right than below
				r1x, r1y = room1.chooseLeftWall(level.rng)
				r2x, r2y = room2.chooseRightWall(level.rng)
			} else { // further below than right
				r1x, r1y = room1.chooseTopWall(level.rng)
				r2x, r2y = room2.chooseBottomWall(level.rng)
			}
		} else { // Left
			if -xDelta > yDelta { // further left than below
				r1x, r1y = room1.chooseRightWall(level.rng)
				r2x, r2y = room2.chooseLeftWall(level.rng)
			} else { // further below than left
				r1x, r1y = room1.chooseTopWall(level.rng)
				r2x, r2y = room2.chooseBottomWall(level.rng)
			}
		}
	} else { // Above
		if xDelta > 0 { // Right
			if xDelta > -yDelta { // further right than above
				r1x, r1y = room1.chooseLeftWall(level.rng)
				r2x, r2y = room2.chooseRightWall(level.rng)
			} else { // further above than right
				r1x, r1y = room1.chooseBottomWall(level.rng)
				r2x, r2y = room2.chooseTopWall(level.rng)
			}
		} else { // Left
			if -xDelta > -yDelta { // further left than above
				r1x, r1y = room1.chooseRightWall(level.rng)
				r2x, r2y = room2.chooseLeftWall(level.rng)
			} else { // further above than left
				r1x, r1y = room1.chooseBottomWall(level.rng)
				r2x, r2y = room2.chooseTopWall(level.rng)
			}
		}
	}

	if r1x > r2x {
		r1x, r2x = r2x, r1x
		r1y, r2y = r2y, r1y
	}
	for ; r1x <= r2x; r1x++ {
		level.tiles[r1y*level.W+r1x].TileKind = Floor
	}
	r1x--
	if r1y > r2y {
		r1x, r2x = r2x, r1x
		r1y, r2y = r2y, r1y
	}
	for ; r1y <= r2y; r1y++ {
		level.tiles[r1y*level.W+r1x].TileKind = Floor
	}
}

func (level *CandidateLevel) addStairs() {
	levelSize := uint64(len(level.tiles))
	if level.flags&GenUpStairs != 0 {
		for i := 0; i < 3; {
			for {
				candidate := level.rng.Bounded(levelSize)
				if level.tiles[candidate].TileKind == Floor {
					level.tiles[candidate].TileKind = UpStair
					i++
					break
				}
			}
		}
	}

	if level.flags&GenDownStairs != 0 {
		for i := 0; i < 3; {
			for {
				candidate := level.rng.Bounded(levelSize)
				if level.tiles[candidate].TileKind == Floor {
					level.tiles[candidate].TileKind = DownStair
					i++
					break
				}
			}
		}
	}
}

func (level *CandidateLevel) addItems() {
	for i := 0; i < MaxItemPlacement; i++ {
		itemIndex := level.rng.Bounded(uint64(len(SampleItems)))
		randomItem := SampleItems[itemIndex]

		tileIndex := level.rng.Bounded(uint64(len(level.tiles)))
		if level.tiles[tileIndex].TileKind == Floor {
			level.tiles[tileIndex].Item = &randomItem
		}
	}
}

func (level *CandidateLevel) encodeAsString() string {
	levelStr := ""
	for y := 0; y < level.H; y++ {
		if y != 0 {
			levelStr += "\n"
		}
		for x := 0; x < level.W; x++ {
			switch level.tiles[y*level.W+x].TileKind {
			case Wall:
				levelStr += string(WallGlyph)
			case Floor:
				item := level.tiles[y*level.W+x].Item
				if item != nil {
					levelStr += string(item.Symbol)
				} else {
					levelStr += string(FloorGlyph)
				}
			case DownStair:
				levelStr += string(DownStairGlyph)
			case UpStair:
				levelStr += string(UpStairGlyph)
			}
		}
	}

	return levelStr
}

type LevelGenFlag int

const (
	GenUpStairs = 1 << iota
	GenDownStairs
)

func GenLevel(rng *pcg.PCG64, maxX int, maxY int, flags LevelGenFlag) *CandidateLevel {
	defer timeMe(time.Now(), "GenLevel")
	subX := rng.Bounded(uint64(maxX / 4))
	subY := rng.Bounded(uint64(maxY / 4))

	W := maxX - int(subX)
	H := maxY - int(subY)

	level := &CandidateLevel{
		rng: rng,
		W:   W,
		H:   H,

		flags: flags,

		rooms: make(map[int]*Room),
		tiles: make([]CandidateTile, W*H, W*H),
	}

	level.addRooms()

	level.addStairs()

	level.addItems()

	return level
}
