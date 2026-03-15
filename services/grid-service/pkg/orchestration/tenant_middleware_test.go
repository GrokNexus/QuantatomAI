package orchestration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireTenantHeader_MissingHeader(t *testing.T) {
	handler := RequireTenantHeader(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metadata/graph", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if body["code"] != "ERR_MISSING_TENANT" {
		t.Fatalf("expected ERR_MISSING_TENANT, got %q", body["code"])
	}
}

func TestRequireTenantHeader_PassesWithHeader(t *testing.T) {
	handler := RequireTenantHeader(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metadata/graph", nil)
	req.Header.Set("X-Tenant-ID", "tenant-ultra")
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
