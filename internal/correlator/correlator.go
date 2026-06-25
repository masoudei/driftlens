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

func (c *Correlator) SubscribeTo(bus eventbus.Bus) {
	bus.Subscribe(eventbus.EventResourceCreated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceUpdated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceDeleted, c.HandleEvent)
}
