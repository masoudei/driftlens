package drift

import (
	"fmt"
	"sync"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
)

type Detector struct {
	store     graph.Store
	baselines map[string]map[string]string
	mu        sync.RWMutex
}

func New(store graph.Store) *Detector {
	return &Detector{
		store:     store,
		baselines: make(map[string]map[string]string),
	}
}

func (d *Detector) SubscribeTo(bus eventbus.Bus) {
	bus.Subscribe(eventbus.EventResourceCreated, d.handleCreated)
	bus.Subscribe(eventbus.EventResourceUpdated, d.handleUpdated)
	bus.Subscribe(eventbus.EventResourceDeleted, d.handleDeleted)
}

func (d *Detector) handleCreated(ev eventbus.Event) {
	resource, _ := ev.Data["resource"].(string)
	name, _ := ev.Data["name"].(string)
	ns, _ := ev.Data["namespace"].(string)
	nodeID := fmt.Sprintf("%s/%s/%s", ns, resource, name)

	d.mu.Lock()
	d.baselines[nodeID] = cloneProps(ev.Data)
	d.mu.Unlock()
}

func (d *Detector) handleUpdated(ev eventbus.Event) {
	resource, _ := ev.Data["resource"].(string)
	name, _ := ev.Data["name"].(string)
	ns, _ := ev.Data["namespace"].(string)
	nodeID := fmt.Sprintf("%s/%s/%s", ns, resource, name)

	d.mu.RLock()
	baseline, exists := d.baselines[nodeID]
	d.mu.RUnlock()

	if !exists {
		d.mu.Lock()
		d.baselines[nodeID] = cloneProps(ev.Data)
		d.mu.Unlock()
		return
	}

	current := cloneProps(ev.Data)
	diffs := compareProps(baseline, current)
	if len(diffs) == 0 {
		return
	}

	driftID := fmt.Sprintf("drift-%d", time.Now().UnixNano())
	driftNode := graph.NewNode(driftID, graph.NodeDriftEvent, map[string]string{
		"resource":  resource,
		"name":      name,
		"namespace": ns,
		"drift":     fmt.Sprintf("%v", diffs),
		"severity":  "medium",
		"detected":  time.Now().UTC().Format(time.RFC3339),
	})
	d.store.AddNode(driftNode)

	rel := graph.NewRelationship(
		fmt.Sprintf("rel-%d", time.Now().UnixNano()),
		nodeID, driftID, graph.Drifted, map[string]string{
			"reason": fmt.Sprintf("changed: %v", diffs),
		},
	)
	d.store.AddRelationship(rel)

	d.mu.Lock()
	d.baselines[nodeID] = current
	d.mu.Unlock()
}

func (d *Detector) handleDeleted(ev eventbus.Event) {
	resource, _ := ev.Data["resource"].(string)
	name, _ := ev.Data["name"].(string)
	ns, _ := ev.Data["namespace"].(string)
	nodeID := fmt.Sprintf("%s/%s/%s", ns, resource, name)

	d.mu.Lock()
	delete(d.baselines, nodeID)
	d.mu.Unlock()
}

func cloneProps(data map[string]any) map[string]string {
	out := make(map[string]string, len(data))
	for k, v := range data {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

func compareProps(baseline, current map[string]string) []string {
	var diffs []string
	for k, v := range current {
		if old, ok := baseline[k]; !ok || old != v {
			diffs = append(diffs, fmt.Sprintf("%s: %s -> %s", k, old, v))
		}
	}
	return diffs
}
