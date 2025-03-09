package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// handleGetPosts handles HTTP GET requests for retrieving filtered posts
func (s *Server) handleGetPosts(ctx *gin.Context) {
	page, limit := getPaginationParams(ctx)

	query := ctx.Query("search")
	tag := ctx.Query("tag")
	postType := ctx.Query("type")
	language := ctx.Query("language")

	response, err := s.db.GetPostsWithFilters(query, tag, postType, language, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"posts":       response.Posts,
		"total_count": response.TotalCount,
	})
}

// handleGetTags handles HTTP GET requests for retrieving all unique tags
func (s *Server) handleGetTags(ctx *gin.Context) {
	tags, err := s.db.GetTags()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tags)
}

// handleGetLanguages handles HTTP GET requests for retrieving all unique languages
func (s *Server) handleGetLanguages(ctx *gin.Context) {
	languages, err := s.db.GetLanguages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, languages)
}

// getPaginationParams extracts and validates pagination parameters from request
// Returns page number and limit with default values if not provided or invalid
func getPaginationParams(c *gin.Context) (int, int) {
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	return page, limit
}
