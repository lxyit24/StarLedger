package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/ent"
	"starledger/ent/tenantmodule"
	"starledger/internal/pkg"
)

// Module type to tenant type mapping (empty = all tenant types can access all modules)
var moduleTenantTypeMap = map[string][]string{}

// ModuleAccess checks if the current tenant has enabled the given module.
// Core modules (billing) are always accessible. This middleware should be used after JWTAuth + TenantIsolation.
func ModuleAccess(client *ent.Client, moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == 0 {
			c.Next()
			return
		}

		// Check if this module has tenant type restrictions
		if allowedTypes, ok := moduleTenantTypeMap[moduleName]; ok {
			t, err := client.Tenant.Get(c.Request.Context(), tenantID)
			if err != nil {
				pkg.Fail(c, http.StatusInternalServerError, "租户信息获取失败")
				c.Abort()
				return
			}

			allowed := false
			for _, typ := range allowedTypes {
				if string(t.Type) == typ {
					allowed = true
					break
				}
			}
			if !allowed {
				pkg.Fail(c, http.StatusForbidden, "该模块仅对特定租户类型可用")
				c.Abort()
				return
			}
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
