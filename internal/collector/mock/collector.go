package mock

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
)

var resources = []string{"Deployment", "StatefulSet", "DaemonSet", "ConfigMap", "Secret"}
var names = []string{"payment-api", "checkout-api", "users-svc", "orders-svc", "inventory-svc"}

type Collector struct {
	bus    eventbus.Bus
	cancel context.CancelFunc
	events int
}

func New(bus eventbus.Bus) *Collector {
	return &Collector{bus: bus}
}

func (c *Collector) Start(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)

	for i, r := range resources {
		c.emit(ctx, r, names[i], eventbus.EventResourceCreated)
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
				types := []eventbus.EventType{
					eventbus.EventResourceUpdated,
					eventbus.EventResourceUpdated,
					eventbus.EventResourceDeleted,
				}
				c.emit(ctx, r, n, types[rand.Intn(len(types))])
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

func (c *Collector) emit(ctx context.Context, resource, name string, eventType eventbus.EventType) {
	id := fmt.Sprintf("%s/%s/%s", "default", resource, name)
	c.bus.Publish(eventbus.Event{
		Type:      eventType,
		Source:    "mock-collector",
		Timestamp: time.Now().UTC(),
		Data: map[string]any{
			"resource":  resource,
			"name":      name,
			"namespace": "default",
			"id":        id,
		},
	})
	c.events++
}
