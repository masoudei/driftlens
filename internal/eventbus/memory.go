package eventbus

import "sync"

type InMemory struct {
	mu          sync.RWMutex
	subscribers map[EventType][]Handler
}

func NewInMemory() *InMemory {
	return &InMemory{
		subscribers: make(map[EventType][]Handler),
	}
}

func (b *InMemory) Publish(event Event) {
	b.mu.RLock()
	handlers := b.subscribers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

func (b *InMemory) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
	b.mu.Unlock()
}

func (b *InMemory) Close() {
	b.mu.Lock()
	b.subscribers = nil
	b.mu.Unlock()
}
