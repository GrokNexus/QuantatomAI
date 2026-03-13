package handlers

import (
	"net/http"
	"strconv"

	"quantatomai/grid-service/src/planner"

	"github.com/gin-gonic/gin"
)

// MetadataHandler exposes dimension and member discovery for UI builders.
type MetadataHandler struct {
	metadata planner.MetadataResolver
}

func NewMetadataHandler(m planner.MetadataResolver) *MetadataHandler {
	return &MetadataHandler{metadata: m}
}

// RegisterRoutes mounts metadata routes.
func (h *MetadataHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/metadata/models/:modelId")
	{
		api.GET("/dimensions", h.ListDimensions)
		api.GET("/dimensions/:dim/members", h.ListMembers)
	}
}

func (h *MetadataHandler) ListDimensions(c *gin.Context) {
	// modelId is present in path for future multi-model support; resolver is already scoped.
	dims, err := h.metadata.ListDimensions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dims)
}

func (h *MetadataHandler) ListMembers(c *gin.Context) {
	dim := c.Param("dim")
	branch := c.Query("branchId")
	parent := c.Query("parentCode")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	opts := planner.MemberListOptions{
		ParentCode: parent,
		Limit:      limit,
		Offset:     offset,
	}

	nodes, err := h.metadata.ListMembers(c.Request.Context(), dim, branch, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}
