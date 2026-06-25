package graph

import "fmt"

// Graph is an in-memory implementation of Store.
type Graph struct {
	nodes map[string]*Node
	rels  map[string]*Relationship
}

func New() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		rels:  make(map[string]*Relationship),
	}
}

func (g *Graph) AddNode(n *Node) error {
	g.nodes[n.ID] = n
	return nil
}

func (g *Graph) GetNode(id string) (*Node, bool, error) {
	n, ok := g.nodes[id]
	return n, ok, nil
}

func (g *Graph) RemoveNode(id string) error {
	delete(g.nodes, id)
	for k, r := range g.rels {
		if r.SourceID == id || r.TargetID == id {
			delete(g.rels, k)
		}
	}
	return nil
}

func (g *Graph) AddRelationship(r *Relationship) error {
	g.rels[r.ID] = r
	return nil
}

func (g *Graph) GetRelationship(id string) (*Relationship, bool, error) {
	r, ok := g.rels[id]
	return r, ok, nil
}

func (g *Graph) RemoveRelationship(id string) error {
	delete(g.rels, id)
	return nil
}

func (g *Graph) ListNodes() ([]*Node, error) {
	out := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		out = append(out, n)
	}
	return out, nil
}

func (g *Graph) ListRelationships() ([]*Relationship, error) {
	out := make([]*Relationship, 0, len(g.rels))
	for _, r := range g.rels {
		out = append(out, r)
	}
	return out, nil
}

func (g *Graph) RelationshipsFrom(id string) ([]*Relationship, error) {
	var out []*Relationship
	for _, r := range g.rels {
		if r.SourceID == id {
			out = append(out, r)
		}
	}
	return out, nil
}

func (g *Graph) RelationshipsTo(id string) ([]*Relationship, error) {
	var out []*Relationship
	for _, r := range g.rels {
		if r.TargetID == id {
			out = append(out, r)
		}
	}
	return out, nil
}

func (g *Graph) Traverse(startID string, maxDepth int) ([]string, error) {
	visited := make(map[string]bool)
	var order []string
	var dfs func(id string, depth int)
	dfs = func(id string, depth int) {
		if depth > maxDepth || visited[id] {
			return
		}
		visited[id] = true
		order = append(order, id)
		for _, r := range g.rels {
			if r.SourceID == id {
				dfs(r.TargetID, depth+1)
			}
		}
	}
	dfs(startID, 0)
	return order, nil
}

func (g *Graph) NodeCount() (int, error) {
	return len(g.nodes), nil
}

func (g *Graph) RelationshipCount() (int, error) {
	return len(g.rels), nil
}

func (g *Graph) Close() error {
	return nil
}

func (g *Graph) String() string {
	return fmt.Sprintf("Graph{nodes=%d, relationships=%d}", len(g.nodes), len(g.rels))
}
