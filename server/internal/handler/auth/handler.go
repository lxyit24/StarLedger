package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"starledger/internal/middleware"
	"starledger/internal/model"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type Handler struct {
	authSvc *service.AuthService
	userSvc *service.UserService
}

func NewHandler(authSvc *service.AuthService, userSvc *service.UserService) *Handler {
	return &Handler{authSvc: authSvc, userSvc: userSvc}
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	u, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		pkg.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := middleware.GenerateToken(u.ID, u.TenantID, u.Username)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	pkg.Success(c, model.LoginResp{
		Token:    token,
		UserID:   u.ID,
		TenantID: u.TenantID,
		Username: u.Username,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	u, err := h.authSvc.Register(c.Request.Context(),
		req.TenantName, req.Contact, req.Phone, req.Email,
		req.Username, req.Password, req.RealName)
	if err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := middleware.GenerateToken(u.ID, u.TenantID, u.Username)
	if err != nil {
		pkg.Fail(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	pkg.Success(c, model.LoginResp{
		Token:    token,
		UserID:   u.ID,
		TenantID: u.TenantID,
		Username: u.Username,
	})
}

func (h *Handler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	u, err := h.userSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		pkg.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}
	pkg.Success(c, gin.H{
		"id":        u.ID,
		"username":  u.Username,
		"real_name": u.RealName,
		"email":     u.Email,
		"phone":     u.Phone,
		"status":    u.Status,
		"tenant_id": u.TenantID,
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req model.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		pkg.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.SuccessWithMessage(c, "密码修改成功")
}
