package audit

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
	svc *service.AuditLogService
}

func NewHandler(svc *service.AuditLogService) *Handler {
	return &Handler{svc: svc}
}

// List returns audit logs with pagination.
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
	action := c.Query("action")
	resourceType := c.Query("resource_type")

	items, total, err := h.svc.List(c.Request.Context(), tenantID, req.Page, req.PageSize, action, resourceType)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, model.PageResp{Total: total, Items: items})
}

// LogRequest is the request body for manually creating an audit log.
type LogRequest struct {
	Action       string `json:"action" binding:"required"`
	ResourceType string `json:"resource_type"`
	ResourceID   int    `json:"resource_id"`
	Detail       string `json:"detail"`
}

// Create manually creates an audit log entry.
func (h *Handler) Create(c *gin.Context) {
	var req LogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	username, _ := c.Get("username")
	ua := username.(string)

	err := h.svc.Create(c.Request.Context(), tenantID, userID, ua, req.Action, req.ResourceType, req.ResourceID, req.Detail, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "记录失败: "+err.Error())
		return
	}
	pkg.SuccessWithMessage(c, "记录成功")
}

// Stats returns audit log statistics.
func (h *Handler) Stats(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	_ = tenantID
	_ = strconv.Itoa(0)
	pkg.Success(c, gin.H{"message": "audit stats placeholder"})
}
