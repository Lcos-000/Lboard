package handlers

import (
	"errors"
	"net/http"
	"whiteboard/server/internal/http/middleware"
	"whiteboard/server/internal/repository"
	"whiteboard/server/internal/service"

	"github.com/gin-gonic/gin"
)

// 依赖注入
type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

type createRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type inviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// CreateRoom 创建房间
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req createRoomRequest
	// 解析请求体参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 创建房间
	room, err := h.roomService.CreateRoom(c.Request.Context(), userID, service.CreateRoomInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidArgument) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"room": room,
	})
}

// ListRooms 列出该用户的房间列表
func (h *RoomHandler) ListRooms(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	rooms, err := h.roomService.ListRooms(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

// GetRoom 获取房间详情
func (h *RoomHandler) GetRoom(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	roomID := c.Param("roomId")
	room, err := h.roomService.GetRoom(c.Request.Context(), userID, roomID)
	// 如果用户不存在或不是房间成员，返回403 Forbidden
	// 如果房间不存在，返回404 Not Found
	// 其他错误，返回500 Internal Server Error
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrForbidden) {
			status = http.StatusForbidden
		}
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

// InviteMember 邀请成员加入房间
func (h *RoomHandler) InviteMember(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	roomID := c.Param("roomId")
	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// 邀请成员
	err := h.roomService.InviteMember(c.Request.Context(), userID, roomID, service.InviteMemberInput{
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrForbidden) {
			status = http.StatusForbidden
		}
		if errors.Is(err, service.ErrInvalidArgument) {
			status = http.StatusBadRequest
		}
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
