package task

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
	svc *service.TaskService
}

func NewHandler(svc *service.TaskService) *Handler {
	return &Handler{svc: svc}
}

// List returns tasks with pagination.
func (h *Handler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	assigneeID, _ := strconv.Atoi(c.DefaultQuery("assignee_id", "0"))

	items, total, err := h.svc.List(c.Request.Context(), tenantID, page, pageSize, assigneeID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询任务失败")
		return
	}

	pkg.Success(c, model.PageResp{
		Total: total,
		Items: items,
	})
}

// Get returns a task by ID.
func (h *Handler) Get(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	item, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusNotFound, "任务不存在")
		return
	}

	pkg.Success(c, item)
}

// Create creates a new task.
func (h *Handler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	var req model.CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		dueDate = &req.DueDate.Time
	}

	item, err := h.svc.Create(c.Request.Context(), tenantID, userID, req.Title, req.Description, req.AssigneeID, req.Priority, dueDate)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "创建任务失败")
		return
	}

	pkg.Success(c, item)
}

// Update updates a task.
func (h *Handler) Update(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	var req model.UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		dueDate = &req.DueDate.Time
	}

	item, err := h.svc.Update(c.Request.Context(), id, tenantID, req.Title, req.Description, req.AssigneeID, req.Status, req.Priority, dueDate)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "更新任务失败")
		return
	}

	pkg.Success(c, item)
}

// Delete deletes a task.
func (h *Handler) Delete(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "删除任务失败")
		return
	}

	pkg.SuccessWithMessage(c, "任务已删除")
}

// Assign assigns a task to a user.
func (h *Handler) Assign(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	var req struct {
		AssigneeID int `json:"assignee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	item, err := h.svc.Assign(c.Request.Context(), id, tenantID, req.AssigneeID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "分配任务失败")
		return
	}

	pkg.Success(c, item)
}

// MyTasks returns tasks assigned to the current user.
func (h *Handler) MyTasks(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	items, err := h.svc.MyTasks(c.Request.Context(), tenantID, userID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询我的任务失败")
		return
	}

	pkg.Success(c, items)
}
