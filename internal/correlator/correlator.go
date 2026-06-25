package correlator

import (
	"fmt"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
)

type Correlator struct {
	store graph.Store
}

func New(store graph.Store) *Correlator {
	return &Correlator{store: store}
}

func (c *Correlator) HandleEvent(ev eventbus.Event) {
	resource, _ := ev.Data["resource"].(string)
	name, _ := ev.Data["name"].(string)
	ns, _ := ev.Data["namespace"].(string)

	nodeID := fmt.Sprintf("%s/%s/%s", ns, resource, name)
	props := map[string]string{
		"resource":  resource,
		"name":      name,
		"namespace": ns,
		"source":    ev.Source,
	}

	switch ev.Type {
	case eventbus.EventResourceCreated:
		props["action"] = "created"
	case eventbus.EventResourceUpdated:
		props["action"] = "updated"
	case eventbus.EventResourceDeleted:
		props["action"] = "deleted"
	}

	existing, found, _ := c.store.GetNode(nodeID)
	if found {
		existing.Timestamp = time.Now().UTC()
		existing.Properties = props
		c.store.AddNode(existing)
		return
	}

	node := graph.NewNode(nodeID, graph.NodeKubernetesEvent, props)
	c.store.AddNode(node)
}

func (c *Correlator) HandleArgoEvent(ev eventbus.Event) {
	app, _ := ev.Data["app"].(string)
	nodeID := fmt.Sprintf("argocd/%s", app)

	props := map[string]string{
		"app":         app,
		"resource":    "ArgoSync",
		"name":        app,
		"source":      ev.Source,
		"sync_status": fmt.Sprintf("%v", ev.Data["curr_sync"]),
		"health":      fmt.Sprintf("%v", ev.Data["curr_health"]),
	}

	switch ev.Type {
	case eventbus.EventArgoSyncStarted:
		props["action"] = "sync-started"
	case eventbus.EventArgoSyncSucceeded:
		props["action"] = "sync-succeeded"
	case eventbus.EventArgoSyncFailed:
		props["action"] = "sync-failed"
	case eventbus.EventArgoHealthChanged:
		props["action"] = "health-changed"
	}

	existing, found, _ := c.store.GetNode(nodeID)
	if found {
		existing.Timestamp = time.Now().UTC()
		existing.Properties = props
		c.store.AddNode(existing)
		return
	}

	node := graph.NewNode(nodeID, graph.NodeArgoSync, props)
	c.store.AddNode(node)

	allNodes, _ := c.store.ListNodes()
	for _, n := range allNodes {
		if n.Type == graph.NodeKubernetesEvent && n.Properties["namespace"] != "" {
			relID := fmt.Sprintf("rel-argo-%s-%s-%d", app, n.ID, time.Now().UnixNano())
			rel := graph.NewRelationship(relID, nodeID, n.ID, graph.Caused, map[string]string{
				"reason": fmt.Sprintf("ArgoCD sync of %s caused %s/%s change", app, n.Properties["resource"], n.Properties["name"]),
			})
			c.store.AddRelationship(rel)
		}
	}
}

func (c *Correlator) SubscribeTo(bus eventbus.Bus) {
	bus.Subscribe(eventbus.EventResourceCreated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceUpdated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceDeleted, c.HandleEvent)
	bus.Subscribe(eventbus.EventArgoSyncStarted, c.HandleArgoEvent)
	bus.Subscribe(eventbus.EventArgoSyncSucceeded, c.HandleArgoEvent)
	bus.Subscribe(eventbus.EventArgoSyncFailed, c.HandleArgoEvent)
	bus.Subscribe(eventbus.EventArgoHealthChanged, c.HandleArgoEvent)
}
