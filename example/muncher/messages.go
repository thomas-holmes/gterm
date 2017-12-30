package main

import "log"

type Message int

const (
	PlayerUpdate Message = iota
	ClearRegion
	MoveEntity
	AttackEntity
	SpellLaunch
	PlayerDead
	KillEntity
	PlayerFloorChange
	ShowMenu
	EquipItem
	GameLogAppend
	ShowFullGameLog
)

type ClearRegionMessage struct {
	X int
	Y int
	W int
	H int
}

type MoveEntityMessage struct {
	ID   int
	OldX int
	OldY int
	NewX int
	NewY int
}

type AttackEntityMesasge struct {
	Attacker Entity
	Defender Entity
}

type SpellLaunchMessage struct {
	Caster Entity
	X      int
	Y      int
	Spell  Spell
}

type KillEntityMessage struct {
	Attacker Entity
	Defender Entity
}

type Listener interface {
	Notify(message Message, data interface{})
}

type Messaging struct {
	messageBus *MessageBus
}

type PlayerFloorChangeMessage struct {
	Stair
}

type ShowMenuMessage struct {
	Menu Menu
}

type EquipItemMessage struct {
	Item
}

type GameLogAppendMessage struct {
	Messages []string
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
