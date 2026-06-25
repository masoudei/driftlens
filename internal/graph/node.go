package graph

import "time"

type NodeType string

const (
	NodeCommit          NodeType = "Commit"
	NodePipeline        NodeType = "Pipeline"
	NodeTerraformRun    NodeType = "TerraformRun"
	NodeArgoSync        NodeType = "ArgoSync"
	NodeKubernetesEvent NodeType = "KubernetesEvent"
	NodeDriftEvent      NodeType = "DriftEvent"
)

type Node struct {
	ID         string
	Type       NodeType
	Properties map[string]string
	Timestamp  time.Time
}

func NewNode(id string, typ NodeType, props map[string]string) *Node {
	return &Node{
		ID:         id,
		Type:       typ,
		Properties: props,
		Timestamp:  time.Now().UTC(),
	}
}
