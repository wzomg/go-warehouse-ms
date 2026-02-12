package api

import (
	"net/http"

	"go-warehouse-ms/internal/model"
	"go-warehouse-ms/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type LoginRequest struct {
	UserID  string `json:"userId"`
	UserPwd string `json:"userPwd"`
}

type RegisterRequest struct {
	UserID  string `json:"userId"`
	UserPwd string `json:"userPwd"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" || req.UserPwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	user, err := h.auth.Login(req.UserID, req.UserPwd)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		case service.ErrPasswordWrong:
			c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "登录失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "user": user})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" || req.UserPwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名和密码必填"})
		return
	}
	user := &model.User{
		UserID:  req.UserID,
		UserPwd: req.UserPwd,
	}
	if err := h.auth.Register(user); err != nil {
		switch err {
		case service.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"message": "用户已存在"})
		case service.ErrInvalidPayload:
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
