package game

type Message int

const (
	PlayerUpdate Message = iota
	TileInvalidated
)

type TileInvalidatedMessage struct {
	XPos int
	YPos int
}

type Listener interface {
	Notify(message Message, data interface{})
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
