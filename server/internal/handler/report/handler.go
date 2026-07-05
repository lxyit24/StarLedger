package report

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/internal/middleware"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	svc *service.ReportService
}

func NewHandler(svc *service.ReportService) *Handler {
	return &Handler{svc: svc}
}

// Overview returns dashboard overview statistics.
func (h *Handler) Overview(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.Overview(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}

// MonthlyTrend returns monthly bill trend data.
func (h *Handler) MonthlyTrend(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.MonthlyTrend(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}

// BillTypeDistribution returns bill amount by type.
func (h *Handler) BillTypeDistribution(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.BillTypeDistribution(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}

// BillStatusSummary returns bill summary by status.
func (h *Handler) BillStatusSummary(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.BillStatusSummary(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}

// ServerCost returns server cost analysis by provider.
func (h *Handler) ServerCost(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.ServerCostAnalysis(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}

// TaskStats returns task statistics.
func (h *Handler) TaskStats(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	data, err := h.svc.TaskStatistics(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}
