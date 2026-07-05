package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/internal/pkg"
)

// TenantIsolation extracts tenant_id from JWT claims and ensures it's set in context.
// This middleware should be used after JWTAuth.
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		if !exists || tenantID.(int) == 0 {
			pkg.Fail(c, http.StatusForbidden, "租户信息缺失")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetTenantID is a helper to get tenant_id from gin context.
func GetTenantID(c *gin.Context) int {
	if v, exists := c.Get("tenant_id"); exists {
		return v.(int)
	}
	return 0
}

// GetUserID is a helper to get user_id from gin context.
func GetUserID(c *gin.Context) int {
	if v, exists := c.Get("user_id"); exists {
		return v.(int)
	}
	return 0
}
