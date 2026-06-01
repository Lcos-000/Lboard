package httpapi

import (
	"time"

	"whiteboard/server/internal/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// RouterDeps 路由依赖项，结构体注入
// 注意这里没有NewRouterDeps函数所以在使用NewRouter的时候要手动组装结构体
type RouterDeps struct {
	Logger        *zap.Logger
	HealthHandler *handlers.HealthHandler
}

// NewRouter 创建HTTP路由
func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(accessLogMiddleware(deps.Logger))

	r.GET("/healthz", deps.HealthHandler.Healthz)
	r.GET("/readyz", deps.HealthHandler.Readyz)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

// accessLogMiddleware 访问日志中间件
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
