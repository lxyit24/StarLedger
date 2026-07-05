package invoice

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
	svc *service.InvoiceService
}

func NewHandler(svc *service.InvoiceService) *Handler {
	return &Handler{svc: svc}
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
	status := c.Query("status")
	invoiceType := c.Query("invoice_type")

	items, total, err := h.svc.List(c.Request.Context(), tenantID, req.Page, req.PageSize, status, invoiceType)
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
		pkg.Fail(c, http.StatusNotFound, "发票不存在")
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Create(c *gin.Context) {
	var req service.CreateInvoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if req.InvoiceType == "" {
		req.InvoiceType = "vat_normal"
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "创建失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req service.UpdateInvoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Update(c.Request.Context(), id, tenantID, req)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "更新失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Issue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tenantID := middleware.GetTenantID(c)

	item, err := h.svc.Issue(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "开具失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tenantID := middleware.GetTenantID(c)

	item, err := h.svc.Cancel(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "作废失败: "+err.Error())
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
	pkg.Success(c, nil)
}

func (h *Handler) Summary(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	data, err := h.svc.Summary(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, data)
}
