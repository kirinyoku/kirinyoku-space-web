package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// handleGetPosts handles GET requests for retrieving paginated posts.
// It extracts pagination parameters and returns a list of posts.
func (s *Server) handleGetPosts(ctx *gin.Context) {
	page, limit := getPaginationParams(ctx)

	posts, err := s.db.GetPosts(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// handleSearchPosts handles GET requests for searching posts by query string.
// If no search query is provided, it falls back to handleGetPosts.
func (s *Server) handleSearchPosts(ctx *gin.Context) {
	query := ctx.Query("search")
	if query != "" {
		page, limit := getPaginationParams(ctx)

		posts, err := s.db.SearchPosts(query, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, posts)

		return
	}

	// If no search query, fall back to GetPosts (handled by the same route)
	s.handleGetPosts(ctx)
}

// handleGetPostsByTag handles GET requests for retrieving posts by tag.
// If no tag is provided, it falls back to handleGetPosts.
func (s *Server) handleGetPostsByTag(ctx *gin.Context) {
	tag := ctx.Query("tag")
	if tag != "" {
		page, limit := getPaginationParams(ctx)

		posts, err := s.db.GetPostsByTag(tag, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, posts)

		return
	}

	// If no tag, fall back to GetPosts (handled by the same route)
	s.handleGetPosts(ctx)
}

// handleGetTags handles GET requests for retrieving all unique tags.
func (s *Server) handleGetTags(c *gin.Context) {
	tags, err := s.db.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// getPaginationParams extracts and validates pagination parameters from the request.
// Returns page number and limit with defaults of page=1 and limit=10 if not specified.
func getPaginationParams(ctx *gin.Context) (int, int) {
	page := 1
	if p := ctx.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10

	if l := ctx.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	return page, limit
}
