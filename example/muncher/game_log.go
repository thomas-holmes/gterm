package main

import (
	"log"

	"github.com/thomas-holmes/gterm"
)

// This is going to result in a lot of heap allocations
type GameLog struct {
	world *World

	X int
	Y int
	W int
	H int

	// Consider switching to fixed size circular buffer
	Messages []string

	Messaging
}

func NewGameLog(x int, y int, w int, h int, world *World, messageBus *MessageBus) *GameLog {
	gameLog := GameLog{
		world: world,
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
		Messaging: Messaging{
			messageBus: messageBus,
		},
	}

	gameLog.messageBus.Subscribe(&gameLog)

	return &gameLog
}

func (gameLog *GameLog) Render(window *gterm.Window) {
	for i := 0; i < gameLog.H && i < len(gameLog.Messages); i++ {
		message := gameLog.Messages[i]
		cut := min(len(message), gameLog.W)
		err := window.PutString(gameLog.X, gameLog.Y+gameLog.H-i, gameLog.Messages[i][:cut], White)
		if err != nil {
			log.Println("Failed to render log?", err)
		}
	}
}

// TODO: This will probably suck perf/allocation wise? Might be constantly
// reallocating I push things back.
func (gameLog *GameLog) appendMessages(messages []string) {
	gameLog.Messages = append(messages, gameLog.Messages...)
}

func (gameLog *GameLog) Notify(message Message, data interface{}) {
	switch message {
	case GameLogAppend:
		if d, ok := data.(GameLogAppendMessage); ok {
			gameLog.appendMessages(d.Messages)
		}
	case ShowFullGameLog:
		menu := &FullGameLog{PopMenu: PopMenu{X: 5, Y: 0, W: 80, H: gameLog.world.Window.Rows - 2}, GameLog: gameLog}
		gameLog.Broadcast(ShowMenu, ShowMenuMessage{menu})
	}
}
