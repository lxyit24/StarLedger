package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/ent"
	"starledger/ent/tenantmodule"
	"starledger/internal/pkg"
)

// ModuleAccess checks if the current tenant has enabled the given module.
// Core modules (billing) are always accessible. This middleware should be used after JWTAuth + TenantIsolation.
func ModuleAccess(client *ent.Client, moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == 0 {
			c.Next()
			return
		}

		// Check if the module is enabled for this tenant
		enabled, err := client.TenantModule.Query().
			Where(
				tenantmodule.TenantID(tenantID),
				tenantmodule.ModuleName(moduleName),
				tenantmodule.Enabled(true),
			).
			Exist(c.Request.Context())
		if err != nil {
			pkg.Fail(c, http.StatusInternalServerError, "模块权限检查失败")
			c.Abort()
			return
		}

		if !enabled {
			pkg.Fail(c, http.StatusForbidden, "该模块未启用，请在市场中启用")
			c.Abort()
			return
		}

		c.Next()
	}
}
