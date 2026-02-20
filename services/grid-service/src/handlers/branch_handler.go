package handlers

import (
	"context"
	"net/http"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"
)

// Branch represents a workspace in the Delta-Branching architecture (Git-Flow).
// Layer 8.1 / Enterprise Wrap.
type Branch struct {
	ID           string    `json:"id"`
	AppID        string    `json:"appId"`
	Name         string    `json:"name"`
	BaseBranchID string    `json:"baseBranchId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type BranchHandler struct {
	db *sql.DB
}

func NewBranchHandler(db *sql.DB) *BranchHandler {
	return &BranchHandler{db: db}
}

// RegisterRoutes attaches the Git-Flow endpoints to the Gin router.
func (h *BranchHandler) RegisterRoutes(r *gin.Engine) {
	branches := r.Group("/api/v1/apps/:appId/branches")
	{
		branches.GET("/", h.ListBranches)
		branches.POST("/", h.CreateBranch)
		// MVP: Future routes for Merging and Committing
	}
}

// ListBranches returns all available sandboxes for a specific App.
func (h *BranchHandler) ListBranches(c *gin.Context) {
	appID := c.Param("appId")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, `
		SELECT id, name, base_branch_id, created_at 
		FROM branches 
		WHERE app_id = $1
		ORDER BY created_at DESC
	`, appID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list branches"})
		return
	}
	defer rows.Close()

	var results []Branch
	for rows.Next() {
		var b Branch
		var base sql.NullString
		if err := rows.Scan(&b.ID, &b.Name, &base, &b.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse branch row"})
			return
		}
		if base.Valid {
			b.BaseBranchID = base.String
		}
		b.AppID = appID
		results = append(results, b)
	}

	c.JSON(http.StatusOK, results)
}

// CreateBranch establishes a new isolated sandbox based off an existing branch (usually 'main').
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	appID := c.Param("appId")

	var req struct {
		Name         string `json:"name" binding:"required"`
		BaseBranchID string `json:"baseBranchId" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch creation payload"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var newBranchID string
	err := h.db.QueryRowContext(ctx, `
		INSERT INTO branches (app_id, name, base_branch_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, appID, req.Name, req.BaseBranchID).Scan(&newBranchID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create branch sandbox"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           newBranchID,
		"name":         req.Name,
		"baseBranchId": req.BaseBranchID,
		"status":       "Branch isolated successfully",
	})
}
