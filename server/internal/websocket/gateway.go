package websocket

import (
	"context"
	"errors"
	"net/http"
	"whiteboard/server/internal/auth"
	"whiteboard/server/internal/service"

	"github.com/gin-gonic/gin"
	gws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Gateway 连接网关
type Gateway struct {
	hub        *Hub
	jwtManager *auth.JWTManager
	roomSvc    *service.RoomService
	handler    MessageHandler
	logger     *zap.Logger
	upgrader   gws.Upgrader
}

// NewGateway 初始化连接网关
func NewGateway(
	hub *Hub,
	jwtManager *auth.JWTManager,
	roomSvc *service.RoomService,
	handler MessageHandler,
	logger *zap.Logger,
) *Gateway {
	return &Gateway{
		hub:        hub,
		jwtManager: jwtManager,
		roomSvc:    roomSvc,
		handler:    handler,
		logger:     logger,
		// 初始化升级器，用于升级HTTP连接为WebSocket连接
		upgrader: gws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 校验Origin，防止跨域攻击
			CheckOrigin: func(r *http.Request) bool {
				// Phase 2 先放开，方便本地开发。
				// 生产环境应校验 Origin。
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket连接请求
func (g *Gateway) HandleWebSocket(c *gin.Context) {
	// 从URL参数中获取token
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	// 校验token
	claims, err := g.jwtManager.Parse(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	// 升级HTTP连接为WebSocket连接
	wsConn, err := g.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		g.logger.Warn("websocket upgrade failed", zap.Error(err))
		return
	}
	// 创建连接
	conn := NewConnection(
		claims.UserID,
		claims.Email,
		wsConn,
		g.hub,
		g.handler,
		g.logger,
	)
	// 注册连接到注册中心
	if err := g.hub.Register(conn); err != nil {
		_ = wsConn.Close()
		if errors.Is(err, ErrTooManyConnections) {
			g.logger.Warn("too many websocket connections")
			return
		}
		g.logger.Error("register websocket connection failed", zap.Error(err))
		return
	}
	// 连接成功，启动连接
	g.logger.Info("websocket connected",
		zap.String("connID", conn.ID()),
		zap.String("userID", conn.UserID()),
		zap.Int("activeConnections", g.hub.ActiveConnections()),
	)
	// 开始读取和写入消息
	conn.Start(context.Background())
}
