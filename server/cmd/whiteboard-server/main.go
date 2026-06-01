package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"whiteboard/server/internal/app"
	"whiteboard/server/internal/config"
	"whiteboard/server/internal/observability"

	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	// 这里的flag返回的是一个指针，指向命令行参数的值，之所以要返回指针，是因为go的作者希望能够区分传入""和未传入参数的区别
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()
	// 加载配置文件
	// 解引用configPath指针
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	// 创建日志记录器
	logger, err := observability.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("create logger failed: %v", err)
	}
	defer func() {
		// 确保日志记录器在程序退出时同步写入
		_ = logger.Sync()
	}()
	// 注册信号处理函数
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 创建应用实例
	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("create app failed", zap.Error(err))
	}
	// 创建错误通道
	errCh := make(chan error, 1)
	// 启动应用
	go func() {
		errCh <- application.Run()
	}()
	// 等待信号或错误
	select {
	case <-ctx.Done():
		logger.Info("received shutdown signal")
	case err := <-errCh:
		if err != nil {
			logger.Fatal("server stopped with error", zap.Error(err))
		}
	}
	// 设置关闭超时时间，防止http.server关闭时阻塞
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 关闭应用
	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", zap.Error(err))
	}
}
