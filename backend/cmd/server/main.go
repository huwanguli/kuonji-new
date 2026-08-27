package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/internal/router"
	"github.com/kuonji/blog/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	logger.Info("Starting server...",
		zap.Int("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

	// 初始化数据库连接
	if err := repository.InitDatabase(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}

	// 自动迁移数据库表结构
	// 已知问题：GORM v1.30 的 AutoMigrate 与 MySQL 的 uniqueIndex 存在兼容性 bug，
	// 会错误地尝试 DROP CONSTRAINT（MySQL 将其解释为 DROP FOREIGN KEY）导致 Error 1091。
	// 解决方案：首次部署用 go run ./cmd/seed 初始化表结构，之后关闭 auto_migrate。
	if cfg.Database.AutoMigrate {
		if err := repository.AutoMigrate(); err != nil {
			logger.Warn("Auto migrate encountered issues (non-fatal, tables may already exist)", zap.Error(err))
		} else {
			logger.Info("Database migrated")
		}
	} else {
		logger.Info("Auto migrate disabled (database.auto_migrate = false)")
	}

	// 初始化 Redis 连接（缓存非核心依赖，失败时降级为警告）
	if err := repository.InitRedis(&cfg.Redis); err != nil {
		logger.Warn("Failed to init redis, cache disabled", zap.Error(err))
	} else {
		logger.Info("Redis connected")
	}

	// 初始化路由
	r := router.Setup(cfg.Server.Mode)

	// 创建 HTTP 服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 在 goroutine 中启动服务
	go func() {
		logger.Info("Server is running", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
