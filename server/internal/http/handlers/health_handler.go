package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	startedAt time.Time

	pg     *pgxpool.Pool
	redis  *redis.Client
	minio  *minio.Client
	bucket string
}

func NewHealthHandler(
	pg *pgxpool.Pool,
	redisClient *redis.Client,
	minioClient *minio.Client,
	bucket string,
) *HealthHandler {
	return &HealthHandler{
		startedAt: time.Now(),
		pg:        pg,
		redis:     redisClient,
		minio:     minioClient,
		bucket:    bucket,
	}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "whiteboard-server",
		"startedAt": h.startedAt.Format(time.RFC3339),
	})
}

// Readyz 检查服务是否准备就绪
func (h *HealthHandler) Readyz(c *gin.Context) {
	// 设置超时时间,用来检查redis或者minio或者postgres连接是否超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]string)
	ready := true

	// 检查PostgreSQL连接
	if err := h.pg.Ping(ctx); err != nil {
		checks["postgres"] = err.Error()
		ready = false
	} else {
		checks["postgres"] = "ok"
	}

	// 检查Redis连接
	if err := h.redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = err.Error()
		ready = false
	} else {
		checks["redis"] = "ok"
	}

	// 检查MinIO连接(这里minio的检查比较特殊，检查的是bucket是否存在)
	if exists, err := h.minio.BucketExists(ctx, h.bucket); err != nil {
		checks["minio"] = err.Error()
		ready = false
	} else if !exists {
		checks["minio"] = "bucket not found"
		ready = false
	} else {
		checks["minio"] = "ok"
	}

	statusCode := http.StatusOK
	statusText := "ready"

	if !ready {
		statusCode = http.StatusServiceUnavailable
		statusText = "not_ready"
	}

	c.JSON(statusCode, gin.H{
		"status": statusText,
		// 嵌套检查结果
		"checks": checks,
	})
}
