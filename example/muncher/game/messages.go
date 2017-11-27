package game

import "log"

type Message int

const (
	PlayerUpdate Message = iota
	ClearRegion
	MoveEntity
	PlayerMove
	CreatureAttack
	PlayerDead
	KillEntity
	PopUpShown
	PopUpHidden
)

type ClearRegionMessage struct {
	XPos   int
	YPos   int
	Width  int
	Height int
}

type MoveEntityMessage struct {
	ID   int
	OldX int
	OldY int
	NewX int
	NewY int
}

type PlayerMoveMessage struct {
	ID   int
	OldX int
	OldY int
	NewX int
	NewY int
}

type CreatureAttackMessage struct {
	Attacker *Creature
	Defender *Creature
}

type KillEntityMessage struct {
	Attacker *Creature
	Defender *Creature
}

type Listener interface {
	Notify(message Message, data interface{})
}

type Messaging struct {
	messageBus *MessageBus
}

type Notifier interface {
	SetMessageBus(messageBus *MessageBus)
	RemoveMessageBus()
	Broadcast(message Message, data interface{})
}

func (messaging *Messaging) SetMessageBus(messageBus *MessageBus) {
	messaging.messageBus = messageBus
}

func (messaging *Messaging) Broadcast(message Message, data interface{}) {
	if messaging.messageBus != nil {
		messaging.messageBus.Broadcast(message, data)
	} else {
		log.Printf("Debug, no message bus for message [%+v] data [%+v]", message, data)
	}
}

func (messaging *Messaging) RemoveMessageBus() {
	messaging.messageBus = nil
}

type MessageBus struct {
	Listeners []Listener
}

func (messageBus *MessageBus) Subscribe(listener Listener) {
	messageBus.Listeners = append(messageBus.Listeners, listener)
}

// Broadcast notifie all listeners. This is synchronous.
func (messageBus MessageBus) Broadcast(message Message, data interface{}) {
	for _, listener := range messageBus.Listeners {
		listener.Notify(message, data)
	}
}
