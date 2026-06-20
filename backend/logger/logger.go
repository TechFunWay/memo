package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugarLogger *zap.SugaredLogger
)

// InitLogger 初始化日志
func InitLogger(dataDir string) error {
	logDir := filepath.Join(dataDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{
		filepath.Join(logDir, "memo.log"),
		"stdout",
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	sugarLogger = logger.Sugar()
	return nil
}

// GetSugarLogger 获取日志记录器
func GetSugarLogger() *zap.SugaredLogger {
	return sugarLogger
}

// Sync 同步日志
func Sync() {
	if sugarLogger != nil {
		sugarLogger.Sync()
	}
}
