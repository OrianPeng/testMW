package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rpa-middleware/internal/api"
	"rpa-middleware/internal/config"
	"rpa-middleware/internal/notification"
	"rpa-middleware/internal/processor"
	"rpa-middleware/internal/queue"
	"rpa-middleware/internal/rpa"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting RPA Middleware Server")

	// 加载配置
	cfg := config.DefaultConfig()
	
	// 可以从环境变量或配置文件加载配置
	if rpaURL := os.Getenv("RPA_BASE_URL"); rpaURL != "" {
		cfg.RPA.BaseURL = rpaURL
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// 这里可以解析端口
	}

	// 创建组件
	queueManager := queue.NewMemoryQueueManager()
	rpaClient := rpa.NewHTTPClient(cfg.RPA.BaseURL, cfg.RPA.Timeout, logger)
	notificationService := notification.NewHTTPNotifier(30*time.Second, logger)

	// 创建队列处理器
	processor := processor.NewQueueProcessor(queueManager, rpaClient, notificationService, logger)

	// 设置路由
	router := api.SetupRouter(queueManager, rpaClient, logger)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动队列处理器
	processor.Start(ctx)

	// 启动 HTTP 服务器
	go func() {
		logger.WithField("port", cfg.Server.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 关闭队列处理器
	processor.Stop()

	// 关闭 HTTP 服务器
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}