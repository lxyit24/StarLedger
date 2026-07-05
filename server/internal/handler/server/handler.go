package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"starledger/internal/middleware"
	"starledger/internal/model"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	svc *service.ServerLeaseService
}

func NewHandler(svc *service.ServerLeaseService) *Handler {
	return &Handler{svc: svc}
}

type CreateServerReq struct {
	ServerName  string    `json:"server_name" binding:"required"`
	Provider    string    `json:"provider"`
	Config      string    `json:"config"`
	IpAddress   string    `json:"ip_address"`
	Location    string    `json:"location"`
	MonthlyCost float64   `json:"monthly_cost"`
	StartDate   model.Date `json:"start_date"`
	EndDate     model.Date `json:"end_date" binding:"required"`
	RenewCycle  string    `json:"renew_cycle"`
	ContractNo  string    `json:"contract_no"`
	Remark      string    `json:"remark"`
}

type UpdateServerReq struct {
	ServerName  string    `json:"server_name"`
	Provider    string    `json:"provider"`
	Config      string    `json:"config"`
	IpAddress   string    `json:"ip_address"`
	Location    string    `json:"location"`
	MonthlyCost float64   `json:"monthly_cost"`
	EndDate     model.Date `json:"end_date"`
	ContractNo  string    `json:"contract_no"`
	Remark      string    `json:"remark"`
}

type RenewReq struct {
	Months int `json:"months" binding:"required,min=1"`
}

func (h *Handler) List(c *gin.Context) {
	var req model.PageReq
	_ = c.ShouldBindQuery(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	tenantID := middleware.GetTenantID(c)
	provider := c.Query("provider")
	status := c.Query("status")
	keyword := c.Query("keyword")

	items, total, err := h.svc.List(c.Request.Context(), tenantID, req.Page, req.PageSize, provider, status, keyword)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, model.PageResp{Total: total, Items: items})
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tenantID := middleware.GetTenantID(c)

	item, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusNotFound, "服务器不存在")
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Create(c.Request.Context(), tenantID,
		req.ServerName, req.Provider, req.Config, req.IpAddress, req.Location,
		req.MonthlyCost, req.StartDate.Time, req.EndDate.Time, req.RenewCycle, req.ContractNo, req.Remark)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "创建失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Update(c.Request.Context(), id, tenantID,
		req.ServerName, req.Provider, req.Config, req.IpAddress, req.Location,
		req.MonthlyCost, req.EndDate.Time, req.ContractNo, req.Remark)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "更新失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tenantID := middleware.GetTenantID(c)

	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "删除失败: "+err.Error())
		return
	}
	pkg.SuccessWithMessage(c, "服务器已退租")
}

func (h *Handler) Renew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req RenewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Renew(c.Request.Context(), id, tenantID, req.Months)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "续租失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Expiring(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	items, err := h.svc.Expiring(c.Request.Context(), tenantID, 30)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, items)
}
