package events

import (
	"sync"
)

// EventHandler defines the signature for functions that can handle events.
type EventHandler func(eventData interface{})

// EventBus stores the mapping of event topics to their subscribers.
type EventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewEventBus creates a new EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe adds a new event handler for a given topic.
func (bus *EventBus) Subscribe(topic string, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if _, found := bus.subscribers[topic]; !found {
		bus.subscribers[topic] = []EventHandler{}
	}
	bus.subscribers[topic] = append(bus.subscribers[topic], handler)
}

// Publish sends an event to all subscribers of a given topic.
// This is done asynchronously to avoid blocking the publisher.
func (bus *EventBus) Publish(topic string, data interface{}) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if handlers, found := bus.subscribers[topic]; found {
		// Run handlers in separate goroutines to avoid blocking
		for _, handler := range handlers {
			go handler(data)
		}
	}
}
