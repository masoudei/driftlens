package k8s

import (
	"context"
	"fmt"
	"sync"

	"github.com/masoudei/driftlens/internal/eventbus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type resourceWatch struct {
	resource  string
	watchFunc func(metav1.ListOptions) (watch.Interface, error)
}

type Collector struct {
	clientset kubernetes.Interface
	bus       eventbus.Bus
	namespace string
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func New(clientset kubernetes.Interface, bus eventbus.Bus, namespace string) *Collector {
	return &Collector{
		clientset: clientset,
		bus:       bus,
		namespace: namespace,
	}
}

func (c *Collector) Start(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)

	resources := []resourceWatch{
		{resource: "Deployment", watchFunc: c.watchDeployments},
		{resource: "StatefulSet", watchFunc: c.watchStatefulSets},
		{resource: "DaemonSet", watchFunc: c.watchDaemonSets},
		{resource: "ConfigMap", watchFunc: c.watchConfigMaps},
		{resource: "Secret", watchFunc: c.watchSecrets},
	}

	for _, r := range resources {
		if err := c.startWatch(ctx, r); err != nil {
			return fmt.Errorf("k8s collector: %s: %w", r.resource, err)
		}
	}

	return nil
}

func (c *Collector) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
}

func (c *Collector) startWatch(ctx context.Context, r resourceWatch) error {
	wi, err := r.watchFunc(metav1.ListOptions{})
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer wi.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-wi.ResultChan():
				if !ok {
					return
				}
				c.handleEvent(r.resource, ev)
			}
		}
	}()

	return nil
}

func (c *Collector) handleEvent(resource string, ev watch.Event) {
	var eventType eventbus.EventType
	switch ev.Type {
	case watch.Added:
		eventType = eventbus.EventResourceCreated
	case watch.Modified:
		eventType = eventbus.EventResourceUpdated
	case watch.Deleted:
		eventType = eventbus.EventResourceDeleted
	default:
		return
	}

	name := nameFromObject(ev.Object)
	ns := namespaceFromObject(ev.Object)

	c.bus.Publish(eventbus.Event{
		Type:   eventType,
		Source: "k8s-collector",
		Data: map[string]any{
			"resource":  resource,
			"name":      name,
			"namespace": ns,
		},
	})
}

func nameFromObject(obj runtime.Object) string {
	if o, ok := obj.(metav1.Object); ok {
		return o.GetName()
	}
	return "unknown"
}

func namespaceFromObject(obj runtime.Object) string {
	if o, ok := obj.(metav1.Object); ok {
		return o.GetNamespace()
	}
	return "default"
}

func (c *Collector) watchDeployments(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.AppsV1().Deployments(c.namespace).Watch(context.TODO(), opts)
}

func (c *Collector) watchStatefulSets(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.AppsV1().StatefulSets(c.namespace).Watch(context.TODO(), opts)
}

func (c *Collector) watchDaemonSets(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.AppsV1().DaemonSets(c.namespace).Watch(context.TODO(), opts)
}

func (c *Collector) watchConfigMaps(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.CoreV1().ConfigMaps(c.namespace).Watch(context.TODO(), opts)
}

func (c *Collector) watchSecrets(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.CoreV1().Secrets(c.namespace).Watch(context.TODO(), opts)
}
