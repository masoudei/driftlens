package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/masoudei/driftlens/internal/api"
	"github.com/masoudei/driftlens/internal/collector/k8s"
	"github.com/masoudei/driftlens/internal/collector/mock"
	"github.com/masoudei/driftlens/internal/correlator"
	"github.com/masoudei/driftlens/internal/drift"
	"github.com/masoudei/driftlens/internal/eventbus"
	"github.com/masoudei/driftlens/internal/graph"
	"github.com/masoudei/driftlens/internal/store/postgres"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	bus := eventbus.NewInMemory()

	store := openStore()
	defer store.Close()

	if err := runMigrations(store); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	corr := correlator.New(store)
	corr.SubscribeTo(bus)

	det := drift.New(store)
	det.SubscribeTo(bus)

	if os.Getenv("DRIFTLENS_DEV") == "true" {
		runDev(bus)
	} else {
		runK8s(bus)
	}

	srv := api.New(store)
	go func() {
		addr := os.Getenv("DRIFTLENS_ADDR")
		if addr == "" {
			addr = ":8080"
		}
		log.Printf("API listening on %s", addr)
		if err := srv.Run(addr); err != nil {
			log.Fatalf("api server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Print("DriftLens started")

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				nc, _ := store.NodeCount()
				rc, _ := store.RelationshipCount()
				log.Printf("Graph: %d nodes, %d relationships", nc, rc)
			}
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down...")
	bus.Close()
}

func openStore() graph.Store {
	dsn := os.Getenv("DRIFTLENS_DATABASE_URL")
	if dsn == "" {
		log.Print("DRIFTLENS_DATABASE_URL not set, using in-memory store")
		return graph.New()
	}

	pg, err := postgres.Open(dsn)
	if err != nil {
		log.Fatalf("postgres connection: %v", err)
	}
	return pg
}

func runMigrations(store graph.Store) error {
	type migrator interface {
		Migrate(ctx context.Context) error
	}
	m, ok := store.(migrator)
	if !ok {
		return nil
	}
	return m.Migrate(context.Background())
}

func runDev(bus eventbus.Bus) {
	log.Print("DEV mode — using mock collector")
	m := mock.New(bus)
	if err := m.Start(context.Background()); err != nil {
		log.Fatalf("mock collector failed: %v", err)
	}
}

func runK8s(bus eventbus.Bus) {
	log.Print("K8s mode — connecting to cluster")
	clientset, err := newClientset()
	if err != nil {
		log.Fatalf("failed to create k8s client: %v", err)
	}

	namespace := os.Getenv("DRIFTLENS_NAMESPACE")
	if namespace == "" {
		log.Print("DRIFTLENS_NAMESPACE not set, watching all namespaces")
	}

	collector := k8s.New(clientset, bus, namespace)
	if err := collector.Start(context.Background()); err != nil {
		log.Fatalf("k8s collector failed: %v", err)
	}
	log.Printf("Watching namespace %s", namespace)
}

func newClientset() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(config)
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath())
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func kubeconfigPath() string {
	if k := os.Getenv("KUBECONFIG"); k != "" {
		return k
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
