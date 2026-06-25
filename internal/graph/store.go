package graph

// Store is the persistence interface for the causality graph.
// Both the in-memory Graph and PostgreSQL store implement it.
type Store interface {
	AddNode(n *Node) error
	GetNode(id string) (*Node, bool, error)
	RemoveNode(id string) error
	AddRelationship(r *Relationship) error
	GetRelationship(id string) (*Relationship, bool, error)
	RemoveRelationship(id string) error
	ListNodes() ([]*Node, error)
	ListRelationships() ([]*Relationship, error)
	RelationshipsFrom(id string) ([]*Relationship, error)
	RelationshipsTo(id string) ([]*Relationship, error)
	Traverse(startID string, maxDepth int) ([]string, error)
	NodeCount() (int, error)
	RelationshipCount() (int, error)
	Close() error
}
