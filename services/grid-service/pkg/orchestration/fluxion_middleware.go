package orchestration

import (
	"context"
	"encoding/json"
	"net/http"
)

// In a real scenario, this would be a DB connection pool executing:
// SELECT fluxion_ai_enabled FROM tenants WHERE id = $1
type TenantConfigStore interface {
	IsFluxionEnabled(ctx context.Context, tenantID string) (bool, error)
}

// Simulated Store for the sake of the Moat architecture implementation
type MockTenantStore struct{}

func (s *MockTenantStore) IsFluxionEnabled(ctx context.Context, tenantID string) (bool, error) {
	// For testing Phase 8.1, only tenant "tenant-ultra" has Fluxion enabled.
	if tenantID == "tenant-ultra" {
		return true, nil
	}
	return false, nil
}

// FluxionMiddleware enforces Phase 8.1 Enterprise AI Data Sovereignty
// by explicitly hard-blocking ANY traffic to AI inference endpoints
// unless the Tenant has explicitly opted-in via the 'fluxion_ai_enabled' DB flag.
func FluxionMiddleware(store TenantConfigStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Tenant ID from injected JWT context
		// (In previous phases, we set this in context or header)
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			http.Error(w, `{"error": "Unauthorized: Missing Tenant ID"}`, http.StatusUnauthorized)
			return
		}

		// 2. Query Postgres for the explicit kill-switch flag
		enabled, err := store.IsFluxionEnabled(r.Context(), tenantID)
		if err != nil {
			http.Error(w, `{"error": "Internal Server Error verifying AI Governance"}`, http.StatusInternalServerError)
			return
		}

		// 3. The Generative Kill-Switch Enforcement
		if !enabled {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error":   "Fluxion AI Engine Blocked",
				"message": "Enterprise AI features are strictly Opt-In. Please enable 'fluxion_ai_enabled' in your Tenant configuration to preserve data sovereignty.",
				"code":    "ERR_FLUXION_GOVERNANCE_BLOCKED",
			}); err != nil {
				http.Error(w, `{"error": "failed to encode governance response"}`, http.StatusInternalServerError)
			}
			return
		}

		// 4. Authorized: Pass through to the AI Endpoints (Forecast/RAG)
		next(w, r)
	}
}
