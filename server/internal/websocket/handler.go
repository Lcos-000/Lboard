package websocket

import (
	"context"
	"errors"

	"whiteboard/server/internal/repository"
	"whiteboard/server/internal/service"

	"go.uber.org/zap"
)

// DefaultMessageHandler 默认消息处理结构体
type DefaultMessageHandler struct {
	roomService *service.RoomService
	logger      *zap.Logger
}

// NewDefaultMessageHandler 创建默认消息处理结构体
func NewDefaultMessageHandler(
	roomService *service.RoomService,
	logger *zap.Logger,
) *DefaultMessageHandler {
	return &DefaultMessageHandler{
		roomService: roomService,
		logger:      logger,
	}
}

// HandleMessage 处理消息
func (h *DefaultMessageHandler) HandleMessage(ctx context.Context, conn *Connection, msg Message) {
	switch msg.Type {
	// 处理Ping消息
	case MessageTypePing:
		h.handlePing(conn, msg)

	// 处理JoinRoom消息
	case MessageTypeJoinRoom:
		h.handleJoinRoom(ctx, conn, msg)

	// 未知消息类型
	default:
		h.sendOrClose(conn, Message{
			Type:      MessageTypeError,
			RequestID: msg.RequestID,
			RoomID:    msg.RoomID,
			// 转换为 JSON 字符串，不检查错误，因为marshal几乎不会出错
			Payload: MustJSON(ErrorPayload{
				Code:    "unknown_message_type",
				Message: "unknown websocket message type",
			}),
		})
	}
}

// handlePing 处理Ping消息
func (h *DefaultMessageHandler) handlePing(conn *Connection, msg Message) {
	h.sendOrClose(conn, Message{
		Type:      MessageTypePong,
		RequestID: msg.RequestID,
		RoomID:    msg.RoomID,
		Payload:   msg.Payload,
	})
}

// handleJoinRoom 处理JoinRoom消息
func (h *DefaultMessageHandler) handleJoinRoom(ctx context.Context, conn *Connection, msg Message) {
	if msg.RoomID == "" {
		// 调用sendOrClose，假如发送失败，关闭连接
		h.sendOrClose(conn, Message{
			Type:      MessageTypeError,
			RequestID: msg.RequestID,
			Payload: MustJSON(ErrorPayload{
				Code:    "missing_room_id",
				Message: "roomId is required",
			}),
		})
		return
	}
	// 调用roomService.GetRoom，获取房间信息
	_, err := h.roomService.GetRoom(ctx, conn.UserID(), msg.RoomID)
	// 如果获取房间信息失败，关闭连接
	if err != nil {
		code := "join_room_failed"
		message := err.Error()

		// 如果是service.ErrForbidden，说明用户不是房间的成员
		if errors.Is(err, service.ErrForbidden) {
			code = "forbidden"
			message = "you are not a member of this room"
		}

		// 如果是repository.ErrNotFound，说明房间不存在
		if errors.Is(err, repository.ErrNotFound) {
			code = "room_not_found"
			message = "room not found"
		}

		// 调用sendOrClose，假如发送失败，关闭连接
		h.sendOrClose(conn, Message{
			Type:      MessageTypeError,
			RequestID: msg.RequestID,
			RoomID:    msg.RoomID,
			Payload: MustJSON(ErrorPayload{
				Code:    code,
				Message: message,
			}),
		})
		return
	}
	// 加入房间成功，记录日志
	h.logger.Info("websocket join room accepted",
		zap.String("connID", conn.ID()),
		zap.String("userID", conn.UserID()),
		zap.String("roomID", msg.RoomID),
	)
	// 加入房间成功，发送JoinAck消息
	h.sendOrClose(conn, Message{
		Type:      MessageTypeJoinAck,
		RequestID: msg.RequestID,
		RoomID:    msg.RoomID,
		Payload: MustJSON(JoinRoomAckPayload{
			RoomID: msg.RoomID,
			Status: "joined",
		}),
	})
}

// sendOrClose 发送消息或关闭连接
func (h *DefaultMessageHandler) sendOrClose(conn *Connection, msg Message) {
	// 假如发送队列满，关闭连接
	if err := conn.Send(msg); err != nil {
		// 记录日志
		h.logger.Warn("websocket send queue full, closing connection",
			zap.String("connID", conn.ID()),
			zap.String("userID", conn.UserID()),
			zap.Error(err),
		)
		// 满了，关闭连接
		conn.Close()
	}
}
