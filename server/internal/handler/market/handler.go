package market

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/internal/middleware"
	"starledger/internal/module"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	svc      *service.MarketService
	registry *module.Registry
}

func NewHandler(svc *service.MarketService, registry *module.Registry) *Handler {
	return &Handler{svc: svc, registry: registry}
}

// ListModules returns all available modules with enabled status for current tenant.
func (h *Handler) ListModules(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	infos := h.registry.ListModuleInfos()

	// Get tenant's module settings
	tenantMods, err := h.svc.GetTenantModules(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询模块状态失败")
		return
	}

	// Merge enabled status
	for i := range infos {
		if enabled, ok := tenantMods[infos[i].Name]; ok {
			infos[i].Enabled = enabled
		} else {
			// If no record exists, core modules are enabled by default
			infos[i].Enabled = infos[i].IsCore
		}
	}

	pkg.Success(c, infos)
}

// EnableModule enables a module for the current tenant.
func (h *Handler) EnableModule(c *gin.Context) {
	name := c.Param("name")
	tenantID := middleware.GetTenantID(c)

	// Check if module exists in registry
	mod, ok := h.registry.GetModule(name)
	if !ok {
		pkg.Fail(c, http.StatusNotFound, "模块不存在")
		return
	}
	_ = mod

	if err := h.svc.EnableModule(c.Request.Context(), tenantID, name); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "启用模块失败")
		return
	}

	pkg.SuccessWithMessage(c, "模块已启用")
}

// DisableModule disables a module for the current tenant.
func (h *Handler) DisableModule(c *gin.Context) {
	name := c.Param("name")
	tenantID := middleware.GetTenantID(c)

	// Check if module exists in registry
	mod, ok := h.registry.GetModule(name)
	if !ok {
		pkg.Fail(c, http.StatusNotFound, "模块不存在")
		return
	}

	// Core modules cannot be disabled
	if mod.IsCore() {
		pkg.Fail(c, http.StatusForbidden, "核心模块不可停用")
		return
	}

	if err := h.svc.DisableModule(c.Request.Context(), tenantID, name); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "停用模块失败")
		return
	}

	pkg.SuccessWithMessage(c, "模块已停用")
}
