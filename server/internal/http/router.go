package httpapi

import (
	"time"

	"whiteboard/server/internal/auth"
	"whiteboard/server/internal/http/handlers"
	"whiteboard/server/internal/http/middleware"
	ws "whiteboard/server/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type RouterDeps struct {
	Logger        *zap.Logger
	JWTManager    *auth.JWTManager
	HealthHandler *handlers.HealthHandler
	AuthHandler   *handlers.AuthHandler
	RoomHandler   *handlers.RoomHandler
	WSGateway     *ws.Gateway
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(accessLogMiddleware(deps.Logger))

	r.GET("/healthz", deps.HealthHandler.Healthz)
	r.GET("/readyz", deps.HealthHandler.Readyz)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ws", deps.WSGateway.HandleWebSocket)

	api := r.Group("/api/v1")

	api.POST("/auth/register", deps.AuthHandler.Register)
	api.POST("/auth/login", deps.AuthHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(deps.JWTManager))

	protected.GET("/auth/me", deps.AuthHandler.Me)

	protected.POST("/rooms", deps.RoomHandler.CreateRoom)
	protected.GET("/rooms", deps.RoomHandler.ListRooms)
	protected.GET("/rooms/:roomId", deps.RoomHandler.GetRoom)
	protected.POST("/rooms/:roomId/members", deps.RoomHandler.InviteMember)

	return r
}

func accessLogMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("http_request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
