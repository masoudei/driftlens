package k8s_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/masoudei/driftlens/internal/collector/k8s"
	"github.com/masoudei/driftlens/internal/eventbus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func newDeployment(name, namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func newConfigMap(name, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func TestCollectorEmitsCreatedForDeployment(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	bus := eventbus.NewInMemory()
	collector := k8s.New(clientset, bus, "default")

	fakeWatcher := watch.NewFake()
	clientset.PrependWatchReactor("deployments", k8stesting.DefaultWatchReactor(fakeWatcher, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := collector.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	done := make(chan struct{})
	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		close(done)
	})

	fakeWatcher.Add(newDeployment("payment-api", "default"))

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for created event from deployment watch")
	}
}

func TestCollectorEmitsCreatedForConfigMap(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	bus := eventbus.NewInMemory()
	collector := k8s.New(clientset, bus, "default")

	fakeWatcher := watch.NewFake()
	clientset.PrependWatchReactor("configmaps", k8stesting.DefaultWatchReactor(fakeWatcher, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := collector.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	done := make(chan struct{})
	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		close(done)
	})

	fakeWatcher.Add(newConfigMap("app-config", "default"))

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for created event from configmap watch")
	}
}

func TestCollectorStop(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	bus := eventbus.NewInMemory()
	collector := k8s.New(clientset, bus, "default")

	ctx, cancel := context.WithCancel(context.Background())
	if err := collector.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	cancel()
	collector.Stop()
}

func TestCollectorIgnoresUnknownWatchEvents(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	bus := eventbus.NewInMemory()
	collector := k8s.New(clientset, bus, "default")

	fakeWatcher := watch.NewFake()
	clientset.PrependWatchReactor("deployments", k8stesting.DefaultWatchReactor(fakeWatcher, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := collector.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	var mu sync.Mutex
	received := 0
	bus.Subscribe(eventbus.EventResourceCreated, func(ev eventbus.Event) {
		mu.Lock()
		received++
		mu.Unlock()
	})

	fakeWatcher.Action(watch.Error, newDeployment("x", "default"))
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if received != 0 {
		t.Fatalf("expected 0 events for Error watch type, got %d", received)
	}
	mu.Unlock()
}
