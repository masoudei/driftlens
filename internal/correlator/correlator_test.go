package correlator_test

import (
	"testing"

	"github.com/masoudei/driftlens/internal/correlator"
	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
)

func TestCorrelatorCreatesNode(t *testing.T) {
	g := graph.New()
	c := correlator.New(g)

	c.HandleEvent(eventbus.Event{
		Type:   eventbus.EventResourceCreated,
		Source: "test",
		Data: map[string]any{
			"resource":  "Deployment",
			"name":      "payment-api",
			"namespace": "default",
		},
	})

	count, _ := g.NodeCount()
	if count != 1 {
		t.Fatalf("expected 1 node, got %d", count)
	}

	node, ok, _ := g.GetNode("default/Deployment/payment-api")
	if !ok {
		t.Fatal("expected node to exist")
	}
	if node.Properties["action"] != "created" {
		t.Fatalf("expected action=created, got %s", node.Properties["action"])
	}
}

func TestCorrelatorUpdatesExistingNode(t *testing.T) {
	g := graph.New()
	c := correlator.New(g)

	c.HandleEvent(eventbus.Event{
		Type:   eventbus.EventResourceCreated,
		Source: "test",
		Data: map[string]any{
			"resource":  "Deployment",
			"name":      "payment-api",
			"namespace": "default",
		},
	})

	c.HandleEvent(eventbus.Event{
		Type:   eventbus.EventResourceUpdated,
		Source: "test",
		Data: map[string]any{
			"resource":  "Deployment",
			"name":      "payment-api",
			"namespace": "default",
		},
	})

	count, _ := g.NodeCount()
	if count != 1 {
		t.Fatalf("expected 1 node (updated), got %d", count)
	}

	node, _, _ := g.GetNode("default/Deployment/payment-api")
	if node.Properties["action"] != "updated" {
		t.Fatalf("expected action=updated, got %s", node.Properties["action"])
	}
}

func TestCorrelatorSubscribesToBus(t *testing.T) {
	g := graph.New()
	bus := eventbus.NewInMemory()
	c := correlator.New(g)
	c.SubscribeTo(bus)

	bus.Publish(eventbus.Event{
		Type:   eventbus.EventResourceCreated,
		Source: "test",
		Data: map[string]any{
			"resource":  "ConfigMap",
			"name":      "app-config",
			"namespace": "prod",
		},
	})

	count, _ := g.NodeCount()
	if count != 1 {
		t.Fatalf("expected 1 node via subscription, got %d", count)
	}
}
