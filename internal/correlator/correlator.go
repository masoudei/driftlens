// Package correlator listens to the event bus and builds the causality graph.
//
// It subscribes to collector events and creates corresponding graph nodes
// and relationships, forming the chain:
//
//	KubernetesEvent ──CAUSED──► DriftEvent
package correlator

import (
	"fmt"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
)

type Correlator struct {
	graph *graph.Graph
}

func New(g *graph.Graph) *Correlator {
	return &Correlator{graph: g}
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

	existing, found := c.graph.GetNode(nodeID)
	if found {
		existing.Timestamp = time.Now().UTC()
		existing.Properties = props
		return
	}

	var nodeType graph.NodeType
	switch resource {
	case "Deployment", "StatefulSet", "DaemonSet":
		nodeType = graph.NodeKubernetesEvent
	case "ConfigMap", "Secret":
		nodeType = graph.NodeKubernetesEvent
	default:
		nodeType = graph.NodeKubernetesEvent
	}

	node := graph.NewNode(nodeID, nodeType, props)
	c.graph.AddNode(node)
}

func (c *Correlator) SubscribeTo(bus eventbus.Bus) {
	bus.Subscribe(eventbus.EventResourceCreated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceUpdated, c.HandleEvent)
	bus.Subscribe(eventbus.EventResourceDeleted, c.HandleEvent)
}
