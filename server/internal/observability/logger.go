package observability

import (
	"strings"

	"whiteboard/server/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 创建日志记录器
func NewLogger(cfg config.LogConfig) (*zap.Logger, error) {
	var zapCfg zap.Config
	// 设置日志格式
	if strings.EqualFold(cfg.Format, "json") {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}
	// 设置日志级别
	level := zapcore.DebugLevel
	//解析日志级别配置
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	// 设置日志级别
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.Encoding = cfg.Format
	// 设置日志输出方式
	if zapCfg.Encoding == "" {
		zapCfg.Encoding = "console"
	}
	// 构建日志记录器
	return zapCfg.Build()
}
