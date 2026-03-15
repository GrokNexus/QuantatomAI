package orchestration

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

// MetadataGraphNode describes a single member node in the metadata graph.
type MetadataGraphNode struct {
	ID        string `json:"id"`
	Dimension string `json:"dimension"`
	Name      string `json:"name"`
	Path      string `json:"path"`
}

// MetadataGraphEdge describes a directional relationship between two graph nodes.
type MetadataGraphEdge struct {
	FromID string `json:"fromId"`
	ToID   string `json:"toId"`
	Type   string `json:"type"`
}

// MetadataGraphResponse is a tenant-scoped lineage/impact view used by Phase 5.
type MetadataGraphResponse struct {
	TenantID    string              `json:"tenantId"`
	AppID       string              `json:"appId"`
	Dimension   string              `json:"dimension"`
	RootMember  string              `json:"rootMember"`
	Nodes       []MetadataGraphNode `json:"nodes"`
	Edges       []MetadataGraphEdge `json:"edges"`
	Ancestors   []string            `json:"ancestors"`
	Descendants []string            `json:"descendants"`
}

// MetadataGraphStore abstracts the source of metadata members.
type MetadataGraphStore interface {
	GetMemberPaths(tenantID, appID, dimension string) ([]string, error)
}

// MockMetadataGraphStore provides deterministic graph samples until DB-backed retrieval is wired.
type MockMetadataGraphStore struct{}

func (s *MockMetadataGraphStore) GetMemberPaths(tenantID, appID, dimension string) ([]string, error) {
	_ = tenantID
	_ = appID
	if strings.EqualFold(dimension, "region") {
		return []string{
			"Global",
			"Global.NorthAmerica",
			"Global.NorthAmerica.USA",
			"Global.NorthAmerica.Canada",
			"Global.EMEA",
			"Global.EMEA.UK",
			"Global.EMEA.Germany",
		}, nil
	}

	return []string{
		"All",
		"All.Topline",
		"All.Topline.Revenue",
		"All.Topline.COGS",
		"All.Bottomline",
		"All.Bottomline.NetIncome",
	}, nil
}

// MetadataGraphHandler returns a tenant-safe metadata graph for the requested root member.
func MetadataGraphHandler(store MetadataGraphStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get(tenantHeader)
		appID := strings.TrimSpace(r.URL.Query().Get("appId"))
		dimension := strings.TrimSpace(r.URL.Query().Get("dimension"))
		rootMember := strings.TrimSpace(r.URL.Query().Get("rootMember"))

		if appID == "" || dimension == "" || rootMember == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "Bad Request",
				"message": "appId, dimension, and rootMember are required query parameters",
				"code":    "ERR_INVALID_METADATA_GRAPH_QUERY",
			})
			return
		}

		paths, err := store.GetMemberPaths(tenantID, appID, dimension)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "Internal Server Error",
				"message": "Failed to load metadata graph",
				"code":    "ERR_METADATA_GRAPH_STORE",
			})
			return
		}

		response := buildMetadataGraphResponse(tenantID, appID, dimension, rootMember, paths)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}

func buildMetadataGraphResponse(tenantID, appID, dimension, rootMember string, paths []string) MetadataGraphResponse {
	nodes := make([]MetadataGraphNode, 0, len(paths))
	edges := make([]MetadataGraphEdge, 0, len(paths))
	ancestors := []string{}
	descendants := []string{}

	rootLower := strings.ToLower(rootMember)
	seenAncestors := map[string]struct{}{}

	for _, fullPath := range paths {
		parts := strings.Split(fullPath, ".")
		name := parts[len(parts)-1]
		id := strings.ToLower(strings.ReplaceAll(fullPath, ".", "_"))

		nodes = append(nodes, MetadataGraphNode{
			ID:        id,
			Dimension: dimension,
			Name:      name,
			Path:      fullPath,
		})

		if len(parts) > 1 {
			parentPath := strings.Join(parts[:len(parts)-1], ".")
			edges = append(edges, MetadataGraphEdge{
				FromID: strings.ToLower(strings.ReplaceAll(parentPath, ".", "_")),
				ToID:   id,
				Type:   "parent-child",
			})
		}

		if strings.EqualFold(name, rootMember) || strings.EqualFold(fullPath, rootMember) {
			for i := 0; i < len(parts)-1; i++ {
				ancestor := parts[i]
				if _, ok := seenAncestors[strings.ToLower(ancestor)]; !ok {
					ancestors = append(ancestors, ancestor)
					seenAncestors[strings.ToLower(ancestor)] = struct{}{}
				}
			}
		}

		if strings.HasPrefix(strings.ToLower(fullPath), rootLower+".") {
			descendants = append(descendants, name)
		}
	}

	sort.Strings(ancestors)
	sort.Strings(descendants)

	return MetadataGraphResponse{
		TenantID:    tenantID,
		AppID:       appID,
		Dimension:   dimension,
		RootMember:  rootMember,
		Nodes:       nodes,
		Edges:       edges,
		Ancestors:   ancestors,
		Descendants: descendants,
	}
}
