package billing

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
	svc *service.BillService
}

func NewHandler(svc *service.BillService) *Handler {
	return &Handler{svc: svc}
}

type CreateBillReq struct {
	BillType   string     `json:"bill_type" binding:"required"`
	ResourceID *int       `json:"resource_id"`
	Amount     float64    `json:"amount" binding:"required"`
	BillDate   model.Date `json:"bill_date"`
	DueDate    model.Date `json:"due_date" binding:"required"`
	Remark     string     `json:"remark"`
}

type PayBillReq struct {
	PaymentMethod string `json:"payment_method"`
	InvoiceNo     string `json:"invoice_no"`
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
	billType := c.Query("bill_type")
	status := c.Query("status")

	items, total, err := h.svc.List(c.Request.Context(), tenantID, req.Page, req.PageSize, billType, status)
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
		pkg.Fail(c, http.StatusNotFound, "账单不存在")
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Create(c.Request.Context(), tenantID, req.BillType, req.ResourceID, req.Amount, req.BillDate.Time, req.DueDate.Time, req.Remark)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "创建失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Pay(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req PayBillReq
	_ = c.ShouldBindJSON(&req)

	tenantID := middleware.GetTenantID(c)
	item, err := h.svc.Pay(c.Request.Context(), id, tenantID, req.PaymentMethod, req.InvoiceNo)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "支付失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
}

func (h *Handler) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tenantID := middleware.GetTenantID(c)

	item, err := h.svc.Cancel(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "取消失败: "+err.Error())
		return
	}
	pkg.Success(c, item)
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

func (h *Handler) Overdue(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	items, err := h.svc.Overdue(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, items)
}

type BatchIDsReq struct {
	IDs []int `json:"ids" binding:"required"`
}

type BatchPayReq struct {
	IDs           []int  `json:"ids" binding:"required"`
	PaymentMethod string `json:"payment_method"`
}

func (h *Handler) BatchPay(c *gin.Context) {
	var req BatchPayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	tenantID := middleware.GetTenantID(c)
	count, err := h.svc.BatchPay(c.Request.Context(), tenantID, req.IDs, req.PaymentMethod)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "批量支付失败: "+err.Error())
		return
	}
	pkg.Success(c, gin.H{"paid_count": count})
}

func (h *Handler) BatchDelete(c *gin.Context) {
	var req BatchIDsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	tenantID := middleware.GetTenantID(c)
	count, err := h.svc.BatchDelete(c.Request.Context(), tenantID, req.IDs)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "批量删除失败: "+err.Error())
		return
	}
	pkg.Success(c, gin.H{"deleted_count": count})
}
