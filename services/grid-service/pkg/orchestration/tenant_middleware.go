package orchestration

import (
	"encoding/json"
	"net/http"
)

const tenantHeader = "X-Tenant-ID"

// RequireTenantHeader enforces tenant propagation for API handlers that must remain tenant-safe.
func RequireTenantHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get(tenantHeader)
		if tenantID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "Unauthorized",
				"message": "Missing required X-Tenant-ID header",
				"code":    "ERR_MISSING_TENANT",
			})
			return
		}

		next(w, r)
	}
}
