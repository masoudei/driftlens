package graph_test

import (
	"testing"

	"github.com/masoudei/driftlens/internal/graph"
)

func TestNewGraph(t *testing.T) {
	g := graph.New()
	count, _ := g.NodeCount()
	if count != 0 {
		t.Fatalf("expected 0 nodes, got %d", count)
	}
	rc, _ := g.RelationshipCount()
	if rc != 0 {
		t.Fatalf("expected 0 rels, got %d", rc)
	}
}

func TestAddAndGetNode(t *testing.T) {
	g := graph.New()
	n := graph.NewNode("n1", graph.NodeCommit, map[string]string{"sha": "abc"})
	g.AddNode(n)

	got, ok, _ := g.GetNode("n1")
	if !ok {
		t.Fatal("expected to find n1")
	}
	if got.Type != graph.NodeCommit {
		t.Fatalf("expected Commit, got %s", got.Type)
	}
}

func TestRemoveNodeCascadesRelationships(t *testing.T) {
	g := graph.New()
	g.AddNode(graph.NewNode("a", graph.NodeCommit, nil))
	g.AddNode(graph.NewNode("b", graph.NodePipeline, nil))
	g.AddRelationship(graph.NewRelationship("r1", "a", "b", graph.Caused, nil))

	g.RemoveNode("a")

	if _, ok, _ := g.GetNode("a"); ok {
		t.Fatal("expected node a to be removed")
	}
	rc, _ := g.RelationshipCount()
	if rc != 0 {
		t.Fatal("expected relationships to cascade on node removal")
	}
}

func TestTraverse(t *testing.T) {
	g := graph.New()
	g.AddNode(graph.NewNode("a", graph.NodeCommit, nil))
	g.AddNode(graph.NewNode("b", graph.NodePipeline, nil))
	g.AddNode(graph.NewNode("c", graph.NodeTerraformRun, nil))
	g.AddNode(graph.NewNode("d", graph.NodeArgoSync, nil))

	g.AddRelationship(graph.NewRelationship("r1", "a", "b", graph.Caused, nil))
	g.AddRelationship(graph.NewRelationship("r2", "b", "c", graph.Triggered, nil))
	g.AddRelationship(graph.NewRelationship("r3", "c", "d", graph.Deployed, nil))

	path, _ := g.Traverse("a", 10)
	if len(path) != 4 {
		t.Fatalf("expected 4 nodes in path, got %d: %v", len(path), path)
	}
	if path[0] != "a" || path[3] != "d" {
		t.Fatalf("unexpected traversal order: %v", path)
	}
}

func TestRelationshipsFromTo(t *testing.T) {
	g := graph.New()
	g.AddNode(graph.NewNode("a", graph.NodeCommit, nil))
	g.AddNode(graph.NewNode("b", graph.NodePipeline, nil))
	g.AddNode(graph.NewNode("c", graph.NodeDriftEvent, nil))

	g.AddRelationship(graph.NewRelationship("r1", "a", "b", graph.Caused, nil))
	g.AddRelationship(graph.NewRelationship("r2", "a", "c", graph.Caused, nil))

	from, _ := g.RelationshipsFrom("a")
	if len(from) != 2 {
		t.Fatalf("expected 2 rels from a, got %d", len(from))
	}

	to, _ := g.RelationshipsTo("c")
	if len(to) != 1 {
		t.Fatalf("expected 1 rel to c, got %d", len(to))
	}
}
