package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MetadataGraphNode struct {
	ID        string `json:"id"`
	Dimension string `json:"dimension"`
	Name      string `json:"name"`
	Path      string `json:"path"`
}

type MetadataGraphEdge struct {
	FromID string `json:"fromId"`
	ToID   string `json:"toId"`
	Type   string `json:"type"`
}

type MetadataGraphResponse struct {
	TenantID    string              `json:"tenantId"`
	AppID       string              `json:"appId"`
	Dimension   string              `json:"dimension"`
	RootMember  string              `json:"rootMember"`
	BranchID    string              `json:"branchId,omitempty"`
	Nodes       []MetadataGraphNode `json:"nodes"`
	Edges       []MetadataGraphEdge `json:"edges"`
	Ancestors   []string            `json:"ancestors"`
	Descendants []string            `json:"descendants"`
}

type metadataPathRow struct {
	ID   string
	Name string
	Path string
}

type MetadataGraphHandler struct {
	db *sql.DB
}

func NewMetadataGraphHandler(db *sql.DB) *MetadataGraphHandler {
	return &MetadataGraphHandler{db: db}
}

func (h *MetadataGraphHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/v1/metadata")
	group.Use(RequireTenantHeader())
	group.GET("/graph", h.GetGraph)
}

func (h *MetadataGraphHandler) GetGraph(c *gin.Context) {
	tenantID := strings.TrimSpace(c.GetHeader(TenantHeaderName))
	appID := strings.TrimSpace(c.Query("appId"))
	dimension := strings.TrimSpace(c.Query("dimension"))
	rootMember := strings.TrimSpace(c.Query("rootMember"))
	branchID := strings.TrimSpace(c.Query("branchId"))

	if appID == "" || dimension == "" || rootMember == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "appId, dimension, and rootMember are required query parameters",
			"code":    "ERR_INVALID_METADATA_GRAPH_QUERY",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rootPath, err := h.resolveRootPath(ctx, tenantID, appID, dimension, rootMember, branchID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "rootMember not found for tenant/app/dimension"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve root member"})
		return
	}

	rows, err := h.queryGraphRows(ctx, tenantID, appID, dimension, rootPath, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query metadata graph"})
		return
	}

	resp := buildGraphResponse(tenantID, appID, dimension, rootMember, rootPath, branchID, rows)
	c.JSON(http.StatusOK, resp)
}

func (h *MetadataGraphHandler) resolveRootPath(ctx context.Context, tenantID, appID, dimension, rootMember, branchID string) (string, error) {
	const q = `
SELECT dm.path::text
FROM dimension_members dm
JOIN dimensions d ON d.id = dm.dimension_id
LEFT JOIN branches b ON b.id = $5::uuid
LEFT JOIN dimension_members child_override
  ON b.id IS NOT NULL
 AND child_override.branch_id = b.id
 AND child_override.dimension_id = dm.dimension_id
 AND child_override.name = dm.name
WHERE dm.tenant_id = $1::uuid
  AND dm.app_id = $2::uuid
  AND lower(d.name) = lower($3)
  AND (
    $5 = ''
    OR dm.branch_id = $5::uuid
    OR (
      dm.branch_id = b.base_branch_id
      AND child_override.id IS NULL
    )
  )
  AND COALESCE(dm.is_deleted, false) = false
  AND (
      lower(dm.name) = lower($4)
      OR lower(dm.path::text) = lower($4)
  )
ORDER BY nlevel(dm.path) ASC
LIMIT 1;
`

	var rootPath string
	err := h.db.QueryRowContext(ctx, q, tenantID, appID, dimension, rootMember, branchID).Scan(&rootPath)
	return rootPath, err
}

func (h *MetadataGraphHandler) queryGraphRows(ctx context.Context, tenantID, appID, dimension, rootPath, branchID string) ([]metadataPathRow, error) {
	const q = `
SELECT dm.id::text, dm.name, dm.path::text
FROM dimension_members dm
JOIN dimensions d ON d.id = dm.dimension_id
LEFT JOIN branches b ON b.id = $5::uuid
LEFT JOIN dimension_members child_override
  ON b.id IS NOT NULL
 AND child_override.branch_id = b.id
 AND child_override.dimension_id = dm.dimension_id
 AND child_override.name = dm.name
WHERE dm.tenant_id = $1::uuid
  AND dm.app_id = $2::uuid
  AND lower(d.name) = lower($3)
  AND (
    $5 = ''
    OR dm.branch_id = $5::uuid
    OR (
      dm.branch_id = b.base_branch_id
      AND child_override.id IS NULL
    )
  )
  AND COALESCE(dm.is_deleted, false) = false
  AND (
      dm.path <@ text2ltree($4)
      OR text2ltree($4) <@ dm.path
  )
ORDER BY nlevel(dm.path), dm.path;
`

	dbRows, err := h.db.QueryContext(ctx, q, tenantID, appID, dimension, rootPath, branchID)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := make([]metadataPathRow, 0, 32)
	for dbRows.Next() {
		var row metadataPathRow
		if err := dbRows.Scan(&row.ID, &row.Name, &row.Path); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	if err := dbRows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func buildGraphResponse(tenantID, appID, dimension, rootMember, rootPath, branchID string, rows []metadataPathRow) MetadataGraphResponse {
	nodes := make([]MetadataGraphNode, 0, len(rows))
	edges := make([]MetadataGraphEdge, 0, len(rows))
	ancestors := make([]string, 0, 8)
	descendants := make([]string, 0, len(rows))

	ancSet := make(map[string]struct{})
	nodeByPath := make(map[string]string, len(rows))

	for _, row := range rows {
		nodes = append(nodes, MetadataGraphNode{
			ID:        row.ID,
			Dimension: dimension,
			Name:      row.Name,
			Path:      row.Path,
		})
		nodeByPath[row.Path] = row.ID
	}

	for _, row := range rows {
		parts := strings.Split(row.Path, ".")
		if len(parts) > 1 {
			parentPath := strings.Join(parts[:len(parts)-1], ".")
			parentID, ok := nodeByPath[parentPath]
			if ok {
				edges = append(edges, MetadataGraphEdge{
					FromID: parentID,
					ToID:   row.ID,
					Type:   "parent-child",
				})
			}
		}

		if strings.HasPrefix(row.Path, rootPath+".") {
			descendants = append(descendants, row.Name)
		}
	}

	rootParts := strings.Split(rootPath, ".")
	for i := 0; i < len(rootParts)-1; i++ {
		anc := rootParts[i]
		if _, exists := ancSet[anc]; !exists {
			ancSet[anc] = struct{}{}
			ancestors = append(ancestors, anc)
		}
	}

	sort.Strings(ancestors)
	sort.Strings(descendants)

	return MetadataGraphResponse{
		TenantID:    tenantID,
		AppID:       appID,
		Dimension:   dimension,
		RootMember:  rootMember,
		BranchID:    branchID,
		Nodes:       nodes,
		Edges:       edges,
		Ancestors:   ancestors,
		Descendants: descendants,
	}
}

func (h *MetadataGraphHandler) String() string {
	return fmt.Sprintf("MetadataGraphHandler{%p}", h)
}
