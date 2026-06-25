// Package eventbus provides a publish-subscribe event bus for decoupling
// collectors from consumers (correlator, graph engine, etc.).
//
// The bus supports typed events so subscribers can filter by event type.
// Currently ships with an in-memory implementation; NATS will be the
// production backend.
package eventbus

import "time"

type EventType string

const (
	EventResourceCreated EventType = "resource.created"
	EventResourceUpdated EventType = "resource.updated"
	EventResourceDeleted EventType = "resource.deleted"
)

type Event struct {
	ID        string
	Type      EventType
	Source    string
	Data      map[string]any
	Timestamp time.Time
}

type Handler func(Event)

type Bus interface {
	Publish(event Event)
	Subscribe(eventType EventType, handler Handler)
	Close()
}
