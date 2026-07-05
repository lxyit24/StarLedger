package user

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
	userSvc *service.UserService
}

func NewHandler(userSvc *service.UserService) *Handler {
	return &Handler{userSvc: userSvc}
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
	items, total, err := h.userSvc.List(c.Request.Context(), tenantID, req.Page, req.PageSize)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, model.PageResp{Total: total, Items: items})
}

func (h *Handler) Create(c *gin.Context) {
	var req model.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	u, err := h.userSvc.Create(c.Request.Context(), tenantID, req.Username, req.Password, req.RealName, req.Email, req.Phone, req.RoleIDs)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.Success(c, u)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	u, err := h.userSvc.Update(c.Request.Context(), id, tenantID, req.RealName, req.Email, req.Phone, req.Status, req.RoleIDs)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.Success(c, u)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	pkg.SuccessWithMessage(c, "用户已禁用")
}

// RoleHandler handles role management.
type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	items, err := h.svc.List(c.Request.Context(), tenantID)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, items)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req model.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	r, err := h.svc.Create(c.Request.Context(), tenantID, req.Name, req.Description, req.Permissions)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.Success(c, r)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	r, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Description, req.Permissions)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.Success(c, r)
}
