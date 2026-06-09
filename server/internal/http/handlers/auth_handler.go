package handlers

import (
	"errors"
	"net/http"

	"whiteboard/server/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理程序
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register 注册用户
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	// 绑定JSON体到req结构体
	// 如果JSON体无效，返回400错误
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// 调用服务层注册用户
	result, err := h.authService.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidArgument) {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	// 返回注册结果
	c.JSON(http.StatusCreated, result)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	// 绑定JSON体到req结构体
	// 如果JSON体无效，返回400错误
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// 调用服务层登录用户
	result, err := h.authService.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidLogin) {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	// 返回登录结果
	c.JSON(http.StatusOK, result)
}

// Me 获取当前用户信息，提供给前端调用
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// 调用服务层获取当前用户信息
	user, err := h.authService.Me(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
