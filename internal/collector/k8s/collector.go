package k8s

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/masoudei/driftlens/internal/eventbus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
		{resource: "Namespace", watchFunc: c.watchNamespaces},
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

	data := map[string]any{
		"resource":  resource,
		"name":      name,
		"namespace": ns,
	}

	for k, v := range extractProperties(resource, ev.Object) {
		data[k] = v
	}

	c.bus.Publish(eventbus.Event{
		Type:   eventType,
		Source: "k8s-collector",
		Data:   data,
	})
}

func (c *Collector) watchNamespaces(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientset.CoreV1().Namespaces().Watch(context.TODO(), opts)
}

func extractProperties(resource string, obj runtime.Object) map[string]any {
	props := map[string]any{}

	switch resource {
	case "Namespace":
		if ns, ok := obj.(*corev1.Namespace); ok {
			props["status_phase"] = string(ns.Status.Phase)
		}
	case "Deployment":
		if d, ok := obj.(*appsv1.Deployment); ok {
			if d.Spec.Replicas != nil {
				props["replicas"] = int(*d.Spec.Replicas)
			}
			props["images"] = containerImages(d.Spec.Template.Spec.Containers)
		}
	case "StatefulSet":
		if s, ok := obj.(*appsv1.StatefulSet); ok {
			if s.Spec.Replicas != nil {
				props["replicas"] = int(*s.Spec.Replicas)
			}
			props["images"] = containerImages(s.Spec.Template.Spec.Containers)
		}
	case "DaemonSet":
		if d, ok := obj.(*appsv1.DaemonSet); ok {
			props["images"] = containerImages(d.Spec.Template.Spec.Containers)
		}
	case "ConfigMap":
		if cm, ok := obj.(*corev1.ConfigMap); ok {
			props["data_keys"] = sortedKeys(cm.Data)
		}
	case "Secret":
		if s, ok := obj.(*corev1.Secret); ok {
			props["secret_keys"] = sortedStringKeys(s.Data)
		}
	}

	return props
}

func containerImages(containers []corev1.Container) string {
	images := make([]string, len(containers))
	for i, c := range containers {
		images[i] = c.Image
	}
	sort.Strings(images)
	return strings.Join(images, ",")
}

func sortedKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func sortedStringKeys(m map[string][]byte) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
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
