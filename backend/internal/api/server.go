// Package api provides:
// - HTTP server functionality for the application
// - HTTP handlers for the API server
package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/db"
)

// Server represents the HTTP server and its dependencies.
type Server struct {
	db     *db.DB      // Database connection
	router *gin.Engine // HTTP router instance
}

// NewServer creates and initializes a new Server instance.
// It takes a database connection as parameter and sets up the routes.
func NewServer(db *db.DB) *Server {
	router := gin.Default()

	router.Use(corsMiddleware())

	s := &Server{
		db:     db,
		router: router,
	}

	s.setupRoutes()

	return s
}

// corsMiddleware adds CORS headers to allow requests from the frontend
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from the frontend origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Allow specific headers (if needed)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// setupRoutes configures all the routes for the HTTP server.
// It sets up endpoints for retrieving posts (with search and tag filtering)
// and getting all available tags.
func (s *Server) setupRoutes() {
	s.router.GET("/posts", s.handleGetPosts)
	s.router.GET("/tags", s.handleGetTags)
	s.router.GET("/languages", s.handleGetLanguages)

}

// Start begins listening for HTTP requests on the specified address.
// Returns an error if the server fails to start.
func (s *Server) Start(addr string) error {
	log.Printf("Starting API server on %s", addr)
	return s.router.Run(addr)
}
