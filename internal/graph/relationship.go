package graph

import "time"

type RelType string

const (
	Caused    RelType = "CAUSED"
	Triggered RelType = "TRIGGERED"
	Modified  RelType = "MODIFIED"
	Deployed  RelType = "DEPLOYED"
	Drifted   RelType = "DRIFTED"
	Impacts   RelType = "IMPACTS"
)

type Relationship struct {
	ID         string
	SourceID   string
	TargetID   string
	Type       RelType
	Properties map[string]string
	Timestamp  time.Time
}

func NewRelationship(id, sourceID, targetID string, typ RelType, props map[string]string) *Relationship {
	return &Relationship{
		ID:         id,
		SourceID:   sourceID,
		TargetID:   targetID,
		Type:       typ,
		Properties: props,
		Timestamp:  time.Now().UTC(),
	}
}
