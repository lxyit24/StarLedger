package contract

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"starledger/internal/middleware"
	"starledger/internal/model"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	svc *service.ContractService
}

func NewHandler(svc *service.ContractService) *Handler {
	return &Handler{svc: svc}
}

// List returns contracts with pagination.
func (h *Handler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := h.svc.List(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询合同失败")
		return
	}

	pkg.Success(c, model.PageResp{
		Total: total,
		Items: items,
	})
}

// Get returns a contract by ID.
func (h *Handler) Get(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的合同ID")
		return
	}

	item, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusNotFound, "合同不存在")
		return
	}

	pkg.Success(c, item)
}

// Create creates a new contract.
func (h *Handler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req model.CreateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	item, err := h.svc.Create(c.Request.Context(), tenantID, req.Title, req.PartyA, req.PartyB, req.Amount, startDate, endDate, req.FileURL, req.Remark)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "创建合同失败")
		return
	}

	pkg.Success(c, item)
}

// Update updates a contract.
func (h *Handler) Update(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的合同ID")
		return
	}

	var req model.UpdateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	item, err := h.svc.Update(c.Request.Context(), id, tenantID, req.Title, req.PartyA, req.PartyB, req.Amount, startDate, endDate, req.Status, req.FileURL, req.Remark)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "更新合同失败")
		return
	}

	pkg.Success(c, item)
}

// Delete deletes a contract.
func (h *Handler) Delete(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的合同ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "删除合同失败")
		return
	}

	pkg.SuccessWithMessage(c, "合同已删除")
}
