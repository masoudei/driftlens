package graph

import "fmt"

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

func (g *Graph) AddNode(n *Node) {
	g.nodes[n.ID] = n
}

func (g *Graph) GetNode(id string) (*Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

func (g *Graph) RemoveNode(id string) {
	delete(g.nodes, id)
	for k, r := range g.rels {
		if r.SourceID == id || r.TargetID == id {
			delete(g.rels, k)
		}
	}
}

func (g *Graph) AddRelationship(r *Relationship) {
	g.rels[r.ID] = r
}

func (g *Graph) GetRelationship(id string) (*Relationship, bool) {
	r, ok := g.rels[id]
	return r, ok
}

func (g *Graph) RemoveRelationship(id string) {
	delete(g.rels, id)
}

func (g *Graph) ListNodes() []*Node {
	out := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		out = append(out, n)
	}
	return out
}

func (g *Graph) ListRelationships() []*Relationship {
	out := make([]*Relationship, 0, len(g.rels))
	for _, r := range g.rels {
		out = append(out, r)
	}
	return out
}

func (g *Graph) RelationshipsFrom(sourceID string) []*Relationship {
	var out []*Relationship
	for _, r := range g.rels {
		if r.SourceID == sourceID {
			out = append(out, r)
		}
	}
	return out
}

func (g *Graph) RelationshipsTo(targetID string) []*Relationship {
	var out []*Relationship
	for _, r := range g.rels {
		if r.TargetID == targetID {
			out = append(out, r)
		}
	}
	return out
}

func (g *Graph) Traverse(startID string, maxDepth int) []string {
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
	return order
}

func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

func (g *Graph) RelationshipCount() int {
	return len(g.rels)
}

func (g *Graph) String() string {
	return fmt.Sprintf("Graph{nodes=%d, relationships=%d}", len(g.nodes), len(g.rels))
}
