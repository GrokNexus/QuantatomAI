package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const TenantHeaderName = "X-Tenant-ID"

// RequireTenantHeader ensures requests include tenant context.
func RequireTenantHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantHeaderName)
		if tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing required X-Tenant-ID header",
				"code":    "ERR_MISSING_TENANT",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
