package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"whiteboard/server/internal/auth"
	"whiteboard/server/internal/config"
	httpapi "whiteboard/server/internal/http"
	"whiteboard/server/internal/http/handlers"
	"whiteboard/server/internal/repository"
	"whiteboard/server/internal/service"
	ws "whiteboard/server/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// App 应用入口,负责初始化和运行应用，这里将main和app分层
type App struct {
	cfg    *config.Config
	logger *zap.Logger

	pg    *pgxpool.Pool
	redis *redis.Client
	minio *minio.Client

	httpServer *http.Server
}

// New 创建应用实例
func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*App, error) {
	gin.SetMode(cfg.Server.GinMode)
	// 初始化数据库连接池
	pgPool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	// 测试数据库连接
	if err := pgPool.Ping(ctx); err != nil {
		pgPool.Close()
		return nil, err
	}
	// 初始化Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		pgPool.Close()
		_ = redisClient.Close()
		return nil, err
	}

	if err := repository.EnsureSchema(ctx, pgPool); err != nil {
		pgPool.Close()
		return nil, err
	}
	// 初始化MinIO连接
	minioClient, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		pgPool.Close()
		_ = redisClient.Close()
		return nil, err
	}
	// 测试MinIO连接
	if err := ensureBucket(ctx, minioClient, cfg.MinIO.Bucket); err != nil {
		pgPool.Close()
		_ = redisClient.Close()
		return nil, err
	}
	// 初始化健康检查处理
	healthHandler := handlers.NewHealthHandler(
		pgPool,
		redisClient,
		minioClient,
		cfg.MinIO.Bucket,
	)
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenTTL)

	userRepo := repository.NewUserRepository(pgPool)
	roomRepo := repository.NewRoomRepository(pgPool)

	authService := service.NewAuthService(userRepo, jwtManager)
	roomService := service.NewRoomService(roomRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	roomHandler := handlers.NewRoomHandler(roomService)

	// 初始化WebSocket连接网关
	wsHub := ws.NewHub(10000)
	wsHandler := ws.NewDefaultMessageHandler(roomService, logger)
	wsGateway := ws.NewGateway(
		wsHub,
		jwtManager,
		roomService,
		wsHandler,
		logger,
	)

	// 初始化HTTP路由
	router := httpapi.NewRouter(httpapi.RouterDeps{
		Logger:        logger,
		JWTManager:    jwtManager,
		HealthHandler: healthHandler,
		AuthHandler:   authHandler,
		RoomHandler:   roomHandler,
		WSGateway:     wsGateway,
	})

	// 初始化HTTP服务器
	httpServer := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		cfg:        cfg,
		logger:     logger,
		pg:         pgPool,
		redis:      redisClient,
		minio:      minioClient,
		httpServer: httpServer,
	}, nil
}

// Run 启动应用
func (a *App) Run() error {
	a.logger.Info("starting whiteboard server",
		zap.String("addr", a.cfg.Server.Addr),
	)

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shutdown 关闭应用
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down whiteboard server")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	// 关闭数据库连接池
	if a.pg != nil {
		a.pg.Close()
	}
	// 关闭Redis连接
	if a.redis != nil {
		_ = a.redis.Close()
	}
	// 注意minio不需要关闭，因为minio-go是基于http的客户端，
	// 它的连接是在需要时才创建的，而不是在初始化时创建的。
	// 因此，不需要在应用关闭时手动关闭minio连接。
	return nil
}

// ensureBucket 确保MinIO桶存在
func ensureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	// 检查桶是否存在
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	// 如果桶存在，直接返回
	if exists {
		return nil
	}
	//如果桶不存在，创建桶
	return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}
