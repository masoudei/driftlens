package eventbus_test

import (
	"sync"
	"testing"

	"github.com/masoudei/driftlens/internal/eventbus"
)

func TestPublishSubscribe(t *testing.T) {
	bus := eventbus.NewInMemory()
	var mu sync.Mutex
	received := 0

	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		mu.Lock()
		received++
		mu.Unlock()
	})

	bus.Publish(eventbus.Event{
		Type:   eventbus.EventResourceCreated,
		Source: "test",
		Data:   map[string]any{"resource": "Deployment"},
	})

	bus.Publish(eventbus.Event{
		Type:   eventbus.EventResourceCreated,
		Source: "test",
		Data:   map[string]any{"resource": "ConfigMap"},
	})

	if received != 2 {
		t.Fatalf("expected 2 events, got %d", received)
	}
}

func TestSubscribeFilteredByType(t *testing.T) {
	bus := eventbus.NewInMemory()
	var mu sync.Mutex
	created := 0
	deleted := 0

	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		mu.Lock()
		created++
		mu.Unlock()
	})
	bus.Subscribe(eventbus.EventResourceDeleted, func(ev eventbus.Event) {
		mu.Lock()
		deleted++
		mu.Unlock()
	})

	bus.Publish(eventbus.Event{Type: eventbus.EventResourceCreated})
	bus.Publish(eventbus.Event{Type: eventbus.EventResourceDeleted})
	bus.Publish(eventbus.Event{Type: eventbus.EventResourceUpdated})

	if created != 1 {
		t.Fatalf("expected 1 created, got %d", created)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted, got %d", deleted)
	}
}

func TestCloseStopsDelivery(t *testing.T) {
	bus := eventbus.NewInMemory()
	received := 0

	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		received++
	})

	bus.Close()
	bus.Publish(eventbus.Event{Type: eventbus.EventResourceCreated})

	if received != 0 {
		t.Fatal("expected no events after close")
	}
}
