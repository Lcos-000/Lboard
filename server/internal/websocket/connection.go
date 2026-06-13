package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
	"whiteboard/server/pkg/uuid"

	gws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 定义连接相关常量
const (
	writeWait            = 10 * time.Second
	pongWait             = 60 * time.Second
	pingPeriod           = 25 * time.Second
	maxMessageSize       = 1 << 20 // 1MB
	defaultSendQueueSize = 64
)

// 定义send队列满误
var ErrSendQueueFull = errors.New("websocket send queue full")

// 定义能处理消息的接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, conn *Connection, msg Message)
}

// Connection 连接结构体
type Connection struct {
	id     string
	userID string
	email  string
	// 连接对象
	ws *gws.Conn
	// 注册中心
	hub *Hub
	// 处理消息的接口
	handler MessageHandler
	logger  *zap.Logger
	// 发送队列，有缓冲
	send chan Message
	// 关闭连接的once，确保只关闭一次
	closeOnce sync.Once
	// 广播关闭信号的channel
	done chan struct{}
}

// 初始化连接
func NewConnection(
	userID string,
	email string,
	ws *gws.Conn,
	hub *Hub,
	handler MessageHandler,
	logger *zap.Logger,
) *Connection {
	return &Connection{
		id:      uuid.NewString(),
		userID:  userID,
		email:   email,
		ws:      ws,
		hub:     hub,
		handler: handler,
		logger:  logger,
		send:    make(chan Message, defaultSendQueueSize),
		done:    make(chan struct{}),
	}
}
func (c *Connection) ID() string {
	return c.id
}
func (c *Connection) UserID() string {
	return c.userID
}
func (c *Connection) Email() string {
	return c.email
}

// Start 启动连接
func (c *Connection) Start(ctx context.Context) {
	// 启动写入循环
	go c.writeLoop()
	// 启动读取循环
	go c.readLoop(ctx)
	// 发送hello消息
	_ = c.Send(Message{
		Type: MessageTypeHello,
		Payload: MustJSON(HelloPayload{
			UserID: c.userID,
		}),
	})
}

// Send 发送消息
func (c *Connection) Send(msg Message) error {
	// 检查连接是否已关闭
	select {
	case <-c.done:
		return nil
	default:
	}

	// 尝试将消息发送入队列
	select {
	case c.send <- msg:
		// 消息发送成功
		return nil
	// 发送队列已满
	default:
		return ErrSendQueueFull
	}
}

// Close 关闭连接，关闭各种资源
func (c *Connection) Close() {
	// 关闭连接，确保只关闭一次
	c.closeOnce.Do(func() {
		close(c.done)
		// 从注册中心注销连接
		c.hub.Unregister(c)
		// 关闭websocket连接
		_ = c.ws.Close()
		// 记录关闭日志
		c.logger.Info("websocket connection closed",
			zap.String("connID", c.id),
			zap.String("userID", c.userID),
		)
	})
}

// 以下均为内部私有方法，不暴露给外部调用，只有同包内的方法可以调用它们

// readLoop 读取循环
func (c *Connection) readLoop(ctx context.Context) {
	defer c.Close()
	// 设置读取限制，防止内存泄漏
	c.ws.SetReadLimit(maxMessageSize)
	// 设置读取超时时间，防止读取超时
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
	// 设置pong处理函数，刷新读取超时时间
	c.ws.SetPongHandler(func(string) error {
		return c.ws.SetReadDeadline(time.Now().Add(pongWait))
	})
	// 循环读取消息
	for {
		// 检查连接是否已关闭
		select {
		case <-c.done:
			return
		// 检查上下文是否已取消
		case <-ctx.Done():
			return
		default:
		}
		// 读取消息
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			c.logger.Info("websocket read stopped",
				zap.String("connID", c.id),
				zap.String("userID", c.userID),
				zap.Error(err),
			)
			return
		}
		var msg Message
		// 解析消息，反序列化为Message结构体
		if err := json.Unmarshal(data, &msg); err != nil {
			// 解析失败，发送错误消息
			_ = c.Send(Message{
				Type: MessageTypeError,
				Payload: MustJSON(ErrorPayload{
					Code:    "invalid_json",
					Message: "invalid websocket message json",
				}),
			})
			continue
		}
		// 处理消息
		c.handler.HandleMessage(ctx, c, msg)
	}
}

// writeLoop 写入循环
func (c *Connection) writeLoop() {
	// 做一个ticker，用于发送ping消息，保持连接活跃
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case <-c.done:
			return
		// 假如有消息要发送
		case msg := <-c.send:
			// 发送消息
			if err := c.writeJSON(msg); err != nil {
				c.logger.Info("websocket write failed",
					zap.String("connID", c.id),
					zap.String("userID", c.userID),
					zap.Error(err),
				)
				return
			}
		// 假如有ping消息要发送
		case <-ticker.C:
			// 发送ping消息
			if err := c.writePing(); err != nil {
				c.logger.Info("websocket ping failed",
					zap.String("connID", c.id),
					zap.String("userID", c.userID),
					zap.Error(err),
				)
				return
			}
		}
	}
}

// writeJSON 发送JSON消息
func (c *Connection) writeJSON(msg Message) error {
	// 设置写入超时时间，防止写入超时，可能因为链接关闭导致设置失败
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return c.ws.WriteJSON(msg)
}

// writePing 发送ping消息
func (c *Connection) writePing() error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return c.ws.WriteMessage(gws.PingMessage, nil)
}
