package mock

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
)

var resources = []string{"Deployment", "StatefulSet", "DaemonSet", "ConfigMap", "Secret"}
var names = []string{"payment-api", "checkout-api", "users-svc", "orders-svc", "inventory-svc"}

type Collector struct {
	bus       eventbus.Bus
	cancel    context.CancelFunc
	events    int
	replicas  map[string]int
	replicasMu sync.Mutex
}

func New(bus eventbus.Bus) *Collector {
	return &Collector{
		bus:      bus,
		replicas: make(map[string]int),
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)

	for i, r := range resources {
		c.emit(ctx, r, names[i], eventbus.EventResourceCreated, 3)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r := resources[rand.Intn(len(resources))]
				n := names[rand.Intn(len(names))]

				if rand.Intn(3) == 0 {
					c.emit(ctx, r, n, eventbus.EventResourceDeleted, 0)
				} else {
					c.emit(ctx, r, n, eventbus.EventResourceUpdated, rand.Intn(10)+1)
				}
			}
		}
	}()

	return nil
}

func (c *Collector) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Collector) emit(ctx context.Context, resource, name string, eventType eventbus.EventType, replicas int) {
	id := fmt.Sprintf("%s/%s/%s", "default", resource, name)

	c.replicasMu.Lock()
	prevReplicas := c.replicas[id]
	if eventType == eventbus.EventResourceCreated || eventType == eventbus.EventResourceUpdated {
		c.replicas[id] = replicas
	} else if eventType == eventbus.EventResourceDeleted {
		delete(c.replicas, id)
	}
	c.replicasMu.Unlock()

	data := map[string]any{
		"resource":  resource,
		"name":      name,
		"namespace": "default",
		"id":        id,
	}

	if resource == "Deployment" || resource == "StatefulSet" || resource == "DaemonSet" {
		data["replicas"] = fmt.Sprintf("%d", replicas)
		images := []string{"nginx:1.25", "nginx:1.26", "nginx:1.27", "app:v2", "app:v3"}
		data["image"] = images[replicas%len(images)]
	}

	if resource == "ConfigMap" {
		data["data_keys"] = fmt.Sprintf("config-%d", replicas)
	}

	if eventType == eventbus.EventResourceUpdated {
		data["prev_replicas"] = fmt.Sprintf("%d", prevReplicas)
	}

	c.bus.Publish(eventbus.Event{
		Type:      eventType,
		Source:    "mock-collector",
		Timestamp: time.Now().UTC(),
		Data:      data,
	})
	c.events++
}
