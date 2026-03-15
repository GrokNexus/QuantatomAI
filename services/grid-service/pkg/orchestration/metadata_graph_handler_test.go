package orchestration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetadataGraphHandler_ValidRequest(t *testing.T) {
	handler := RequireTenantHeader(MetadataGraphHandler(&MockMetadataGraphStore{}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metadata/graph?appId=app-1&dimension=region&rootMember=Global", nil)
	req.Header.Set("X-Tenant-ID", "tenant-ultra")
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp MetadataGraphResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.TenantID != "tenant-ultra" {
		t.Fatalf("expected tenant-ultra, got %q", resp.TenantID)
	}

	if len(resp.Nodes) == 0 {
		t.Fatalf("expected non-empty nodes")
	}

	if len(resp.Edges) == 0 {
		t.Fatalf("expected non-empty edges")
	}
}

func TestMetadataGraphHandler_RejectsMissingQueryParameters(t *testing.T) {
	handler := RequireTenantHeader(MetadataGraphHandler(&MockMetadataGraphStore{}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metadata/graph?appId=app-1&dimension=region", nil)
	req.Header.Set("X-Tenant-ID", "tenant-ultra")
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
