package tenant

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"starledger/internal/model"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	svc *service.TenantService
}

func NewHandler(svc *service.TenantService) *Handler {
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

	items, total, err := h.svc.List(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	pkg.Success(c, model.PageResp{Total: total, Items: items})
}

func (h *Handler) Create(c *gin.Context) {
	var req model.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	t, err := h.svc.Create(c.Request.Context(), req.Name, req.Contact, req.Phone, req.Email)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	pkg.Success(c, t)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	t, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Contact, req.Phone, req.Email, req.Status)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	pkg.Success(c, t)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	pkg.SuccessWithMessage(c, "租户已禁用")
}
