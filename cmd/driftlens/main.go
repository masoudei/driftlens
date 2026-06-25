package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masoudei/driftlens/internal/collector/k8s"
	"github.com/masoudei/driftlens/internal/collector/mock"
	"github.com/masoudei/driftlens/internal/correlator"
	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	bus := eventbus.NewInMemory()
	g := graph.New()

	corr := correlator.New(g)
	corr.SubscribeTo(bus)

	if os.Getenv("DRIFTLENS_DEV") == "true" {
		runDev(bus)
	} else {
		runK8s(bus)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("DriftLens started — graph: %s", g)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Printf("Graph: %d nodes, %d relationships",
					g.NodeCount(), g.RelationshipCount())
			}
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	bus.Close()
}

func runDev(bus eventbus.Bus) {
	log.Println("DEV mode — using mock collector")
	m := mock.New(bus)
	if err := m.Start(context.Background()); err != nil {
		log.Fatalf("mock collector failed: %v", err)
	}
}

func runK8s(bus eventbus.Bus) {
	log.Println("K8s mode — connecting to cluster")
	clientset, err := newClientset()
	if err != nil {
		log.Fatalf("failed to create k8s client: %v", err)
	}

	namespace := os.Getenv("DRIFTLENS_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	collector := k8s.New(clientset, bus, namespace)
	if err := collector.Start(context.Background()); err != nil {
		log.Fatalf("k8s collector failed: %v", err)
	}
	log.Printf("Watching namespace %s", namespace)
}

func newClientset() (kubernetes.Interface, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return kubernetes.NewForConfig(config)
	}

	config, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(config)
	}

	config, err = clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
