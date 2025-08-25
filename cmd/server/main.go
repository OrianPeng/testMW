package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rpa-middleware/internal/api"
	"rpa-middleware/internal/config"
	"rpa-middleware/internal/database"
	"rpa-middleware/internal/notification"
	"rpa-middleware/internal/processor"
	"rpa-middleware/internal/queue"
	"rpa-middleware/internal/repository"
	"rpa-middleware/internal/rpa"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting RPA Middleware Server")

	// 加载配置
	cfg := loadConfig(logger)

	// 可以从环境变量覆盖配置
	if rpaURL := os.Getenv("RPA_BASE_URL"); rpaURL != "" {
		cfg.RPA.BaseURL = rpaURL
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}
	if queueType := os.Getenv("QUEUE_TYPE"); queueType != "" {
		cfg.Queue.Type = queueType
	}

	// 初始化MySQL数据库
	mysqlManager, err := database.NewMySQLManager(&cfg.MySQL, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize MySQL database")
	}
	defer mysqlManager.Close()

	// 初始化数据库表
	if err := mysqlManager.InitTables(); err != nil {
		logger.WithError(err).Fatal("Failed to initialize database tables")
	}

	// 创建采购请求仓库
	purchaseRepo := repository.NewPurchaseRequestRepository(mysqlManager.GetDB(), logger)

	// 创建队列管理器
	queueManager, err := queue.NewQueueManager(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create queue manager")
	}

	// 创建组件
	rpaClient := rpa.NewHTTPClient(cfg.RPA.BaseURL, cfg.RPA.Timeout, logger)
	notificationService := notification.NewHTTPNotifier(30*time.Second, logger)

	// 创建队列处理器
	processor := processor.NewQueueProcessor(queueManager, rpaClient, notificationService, logger)

	// 设置路由（包含采购请求功能和队列管理）
	router := api.SetupRouterWithPurchaseRequests(queueManager, rpaClient, purchaseRepo, logger)

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

	// 关闭队列管理器（如果是 Redis）
	if redisQueue, ok := queueManager.(*queue.RedisQueueManager); ok {
		if err := redisQueue.Close(); err != nil {
			logger.WithError(err).Error("Failed to close Redis connection")
		}
	}

	// 关闭 HTTP 服务器
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

// loadConfig 加载配置文件
func loadConfig(logger *logrus.Logger) *config.Config {
	// 首先尝试从配置文件加载
	if data, err := os.ReadFile("config.yaml"); err == nil {
		var cfg config.Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			logger.WithError(err).Warn("Failed to parse config.yaml, using default config")
			return config.DefaultConfig()
		}
		logger.Info("Loaded configuration from config.yaml")
		return &cfg
	}

	// 如果配置文件不存在，使用默认配置
	logger.Info("Using default configuration")
	return config.DefaultConfig()
}
