package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/masoudei/driftlens/internal/graph"
)

type Store struct {
	db *sql.DB
}

func Open(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS nodes (
			id          TEXT PRIMARY KEY,
			type        TEXT NOT NULL,
			properties  JSONB DEFAULT '{}',
			timestamp   TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS relationships (
			id          TEXT PRIMARY KEY,
			source_id   TEXT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
			target_id   TEXT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
			type        TEXT NOT NULL,
			properties  JSONB DEFAULT '{}',
			timestamp   TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_rels_source ON relationships(source_id);
		CREATE INDEX IF NOT EXISTS idx_rels_target ON relationships(target_id);
	`)
	return err
}

func (s *Store) AddNode(n *graph.Node) error {
	props, _ := json.Marshal(n.Properties)
	_, err := s.db.Exec(
		`INSERT INTO nodes (id, type, properties, timestamp)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE SET type=$2, properties=$3, timestamp=$4`,
		n.ID, string(n.Type), props, n.Timestamp,
	)
	return err
}

func (s *Store) GetNode(id string) (*graph.Node, bool, error) {
	var n graph.Node
	var props []byte
	var ts time.Time
	err := s.db.QueryRow(
		`SELECT id, type, properties, timestamp FROM nodes WHERE id=$1`, id,
	).Scan(&n.ID, (*string)(&n.Type), &props, &ts)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	n.Timestamp = ts
	json.Unmarshal(props, &n.Properties)
	return &n, true, nil
}

func (s *Store) RemoveNode(id string) error {
	_, err := s.db.Exec(`DELETE FROM nodes WHERE id=$1`, id)
	return err
}

func (s *Store) AddRelationship(r *graph.Relationship) error {
	props, _ := json.Marshal(r.Properties)
	_, err := s.db.Exec(
		`INSERT INTO relationships (id, source_id, target_id, type, properties, timestamp)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (id) DO UPDATE SET source_id=$2, target_id=$3, type=$4, properties=$5, timestamp=$6`,
		r.ID, r.SourceID, r.TargetID, string(r.Type), props, r.Timestamp,
	)
	return err
}

func (s *Store) GetRelationship(id string) (*graph.Relationship, bool, error) {
	var r graph.Relationship
	var props []byte
	var ts time.Time
	err := s.db.QueryRow(
		`SELECT id, source_id, target_id, type, properties, timestamp FROM relationships WHERE id=$1`, id,
	).Scan(&r.ID, &r.SourceID, &r.TargetID, (*string)(&r.Type), &props, &ts)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	r.Timestamp = ts
	json.Unmarshal(props, &r.Properties)
	return &r, true, nil
}

func (s *Store) RemoveRelationship(id string) error {
	_, err := s.db.Exec(`DELETE FROM relationships WHERE id=$1`, id)
	return err
}

func (s *Store) ListNodes() ([]*graph.Node, error) {
	rows, err := s.db.Query(`SELECT id, type, properties, timestamp FROM nodes ORDER BY timestamp`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNodes(rows)
}

func (s *Store) ListRelationships() ([]*graph.Relationship, error) {
	rows, err := s.db.Query(`SELECT id, source_id, target_id, type, properties, timestamp FROM relationships ORDER BY timestamp`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRelationships(rows)
}

func (s *Store) RelationshipsFrom(id string) ([]*graph.Relationship, error) {
	rows, err := s.db.Query(
		`SELECT id, source_id, target_id, type, properties, timestamp FROM relationships WHERE source_id=$1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRelationships(rows)
}

func (s *Store) RelationshipsTo(id string) ([]*graph.Relationship, error) {
	rows, err := s.db.Query(
		`SELECT id, source_id, target_id, type, properties, timestamp FROM relationships WHERE target_id=$1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRelationships(rows)
}

func (s *Store) Traverse(startID string, maxDepth int) ([]string, error) {
	query := `
		WITH RECURSIVE chain AS (
			SELECT id, 0 AS depth FROM nodes WHERE id = $1
			UNION
			SELECT r.target_id, c.depth + 1
			FROM chain c
			JOIN relationships r ON r.source_id = c.id
			WHERE c.depth < $2
		)
		SELECT DISTINCT id FROM chain ORDER BY depth`
	rows, err := s.db.Query(query, startID, maxDepth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *Store) NodeCount() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM nodes`).Scan(&count)
	return count, err
}

func (s *Store) RelationshipCount() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM relationships`).Scan(&count)
	return count, err
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) NodesByType(typ graph.NodeType) ([]*graph.Node, error) {
	rows, err := s.db.Query(`SELECT id, type, properties, timestamp FROM nodes WHERE type=$1`, string(typ))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNodes(rows)
}

func scanNodes(rows *sql.Rows) ([]*graph.Node, error) {
	var nodes []*graph.Node
	for rows.Next() {
		var n graph.Node
		var props []byte
		var ts time.Time
		if err := rows.Scan(&n.ID, (*string)(&n.Type), &props, &ts); err != nil {
			return nil, err
		}
		n.Timestamp = ts
		json.Unmarshal(props, &n.Properties)
		if n.Properties == nil {
			n.Properties = make(map[string]string)
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func scanRelationships(rows *sql.Rows) ([]*graph.Relationship, error) {
	var rels []*graph.Relationship
	for rows.Next() {
		var r graph.Relationship
		var props []byte
		var ts time.Time
		if err := rows.Scan(&r.ID, &r.SourceID, &r.TargetID, (*string)(&r.Type), &props, &ts); err != nil {
			return nil, err
		}
		r.Timestamp = ts
		json.Unmarshal(props, &r.Properties)
		rels = append(rels, &r)
	}
	return rels, nil
}

var _ graph.Store = (*Store)(nil)
