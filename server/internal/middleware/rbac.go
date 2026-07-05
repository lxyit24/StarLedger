package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"starledger/internal/pkg"
)

// RequirePermission checks if the user has the required permission.
// Permissions are stored in the user's role and loaded into context by a prior middleware.
func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			pkg.Fail(c, http.StatusForbidden, "权限信息未加载")
			c.Abort()
			return
		}

		perms, ok := permissions.([]string)
		if !ok || !slices.Contains(perms, perm) {
			pkg.Fail(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole checks if the user has the required role name.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			pkg.Fail(c, http.StatusForbidden, "角色信息未加载")
			c.Abort()
			return
		}

		roleList, ok := roles.([]string)
		if !ok || !slices.Contains(roleList, role) {
			pkg.Fail(c, http.StatusForbidden, "角色权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
