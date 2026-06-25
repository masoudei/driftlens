package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/masoudei/driftlens/internal/eventbus"
)

type syncInfo struct {
	Status   string `json:"status"`
	Revision string `json:"revision,omitempty"`
}

type healthInfo struct {
	Status string `json:"status"`
}

type appStatus struct {
	Sync   syncInfo   `json:"sync"`
	Health healthInfo `json:"health"`
}

type trackedApp struct {
	Name string
	Sync string
	Health string
}

type Collector struct {
	serverURL string
	token     string
	bus       eventbus.Bus
	interval  time.Duration
	prevState map[string]trackedApp
	mu        sync.RWMutex
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func New(serverURL, token string, bus eventbus.Bus, interval time.Duration) *Collector {
	return &Collector{
		serverURL: serverURL,
		token:     token,
		bus:       bus,
		interval:  interval,
		prevState: make(map[string]trackedApp),
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.pollOnce(ctx)
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.pollOnce(ctx)
			}
		}
	}()

	return nil
}

func (c *Collector) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
}

func (c *Collector) pollOnce(ctx context.Context) {
	apps, err := c.fetchApplications(ctx)
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, app := range apps {
		name := app.Metadata.Name
		curr := trackedApp{Name: name, Sync: app.Status.Sync.Status, Health: app.Status.Health.Status}
		prev, exists := c.prevState[name]
		c.prevState[name] = curr
		if !exists {
			continue
		}
		c.detectChanges(prev, curr)
	}
}

type argocdAppMeta struct {
	Name string `json:"name"`
}

type argocdAppItem struct {
	Metadata argocdAppMeta `json:"metadata"`
	Status   appStatus     `json:"status"`
}

type argocdResponse struct {
	Items []argocdAppItem `json:"items"`
}

func (c *Collector) fetchApplications(ctx context.Context) ([]argocdAppItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.serverURL+"/api/v1/applications", nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("argocd api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("argocd api: %s", resp.Status)
	}

	var result argocdResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("argocd decode: %w", err)
	}
	return result.Items, nil
}

func (c *Collector) detectChanges(prev, curr trackedApp) {
	if prev.Sync != curr.Sync {
		var evType eventbus.EventType
		switch curr.Sync {
		case "Synced":
			evType = eventbus.EventArgoSyncSucceeded
		case "OutOfSync":
			evType = eventbus.EventArgoSyncFailed
		default:
			evType = eventbus.EventArgoSyncStarted
		}
		c.bus.Publish(eventbus.Event{
			Type:   evType,
			Source: "argocd-collector",
			Data: map[string]any{
				"app":        curr.Name,
				"prev_sync":  prev.Sync,
				"curr_sync":  curr.Sync,
				"prev_health": prev.Health,
				"curr_health": curr.Health,
			},
		})
	}
	if prev.Health != curr.Health {
		c.bus.Publish(eventbus.Event{
			Type:   eventbus.EventArgoHealthChanged,
			Source: "argocd-collector",
			Data: map[string]any{
				"app":         curr.Name,
				"prev_health": prev.Health,
				"curr_health": curr.Health,
			},
		})
	}
}
