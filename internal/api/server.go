package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/masoudei/driftlens/internal/graph"
)

type Server struct {
	store graph.Store
	engine *gin.Engine
}

func New(store graph.Store) *Server {
	s := &Server{store: store}
	s.engine = gin.Default()
	s.engine.Use(corsMiddleware())
	s.registerRoutes()
	return s
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", s.health)
	s.engine.GET("/drifts", s.listDrifts)
	s.engine.GET("/drifts/:id", s.getDrift)
	s.engine.GET("/timeline/:resource", s.timeline)
	s.engine.GET("/risk/:resource", s.risk)
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func (s *Server) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (s *Server) listDrifts(c *gin.Context) {
	nodes, err := s.store.ListNodes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var drifts []gin.H
	for _, n := range nodes {
		if n.Type == graph.NodeDriftEvent {
			rels, _ := s.store.RelationshipsTo(n.ID)
			var sources []string
			for _, r := range rels {
				sources = append(sources, r.SourceID)
			}
			drifts = append(drifts, gin.H{
				"id":         n.ID,
				"timestamp":  n.Timestamp,
				"properties": n.Properties,
				"sources":    sources,
			})
		}
	}
	if drifts == nil {
		drifts = []gin.H{}
	}
	c.JSON(200, gin.H{"drifts": drifts})
}

func (s *Server) getDrift(c *gin.Context) {
	id := c.Param("id")

	node, ok, err := s.store.GetNode(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(404, gin.H{"error": "drift not found"})
		return
	}

	rels, _ := s.store.RelationshipsTo(id)
	var sources []string
	for _, r := range rels {
		sources = append(sources, r.SourceID)
	}

	chain, _ := s.store.Traverse(id, 10)

	c.JSON(200, gin.H{
		"drift": gin.H{
			"id":         node.ID,
			"type":       node.Type,
			"timestamp":  node.Timestamp,
			"properties": node.Properties,
			"sources":    sources,
			"chain":      chain,
		},
	})
}

func (s *Server) timeline(c *gin.Context) {
	resource := c.Param("resource")
	maxDepth, _ := strconv.Atoi(c.DefaultQuery("depth", "10"))

	allNodes, err := s.store.ListNodes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, n := range allNodes {
		if n.Properties["resource"] == resource || n.Properties["name"] == resource || string(n.Type) == resource {
			chain, _ := s.store.Traverse(n.ID, maxDepth)
			result = append(result, gin.H{
				"id":        n.ID,
				"type":      n.Type,
				"timestamp": n.Timestamp,
				"props":     n.Properties,
				"chain":     chain,
			})
		}
	}

	if result == nil {
		result = []gin.H{}
	}
	c.JSON(200, gin.H{"timeline": result})
}

func (s *Server) risk(c *gin.Context) {
	resource := c.Param("resource")

	allNodes, err := s.store.ListNodes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var driftCount int
	for _, n := range allNodes {
		if n.Type == graph.NodeDriftEvent && n.Properties["resource"] == resource {
			driftCount++
		}
	}

	severity := "low"
	if driftCount > 5 {
		severity = "high"
	} else if driftCount > 2 {
		severity = "medium"
	}

	c.JSON(200, gin.H{
		"resource":       resource,
		"drift_count":    driftCount,
		"risk_severity":  severity,
	})
}
