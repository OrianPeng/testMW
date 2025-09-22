package main

import (
	"fmt"
	"os"
	"rpa-middleware/internal/api"
	"rpa-middleware/internal/config"
	"rpa-middleware/internal/database"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/queue"
	"rpa-middleware/internal/repository"
	"rpa-middleware/internal/rpa"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.Info("Starting full API server...")

	// 加载配置
	cfg := loadConfig(logger)

	// 初始化数据库
	db, err := database.NewMySQLManager(&cfg.MySQL, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()

	// 初始化数据库表
	if err := db.InitTables(); err != nil {
		logger.WithError(err).Fatal("Failed to initialize database tables")
	}

	// 初始化队列管理器
	queueManager, err := queue.NewQueueManager(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize queue manager")
	}

	// 初始化RPA客户端
	rpaClient := rpa.NewHTTPClient(cfg.RPA.BaseURL, cfg.RPA.Timeout, logger)

	// 初始化UiPath客户端
	uipathConfig := &models.UiPathConfig{
		OrchBaseURL: cfg.UiPath.OrchBaseURL,
		TenancyName: cfg.UiPath.TenancyName,
		Username:    cfg.UiPath.Username,
		Password:    cfg.UiPath.Password,
		FolderID:    cfg.UiPath.FolderID,
		QueueName:   cfg.UiPath.QueueName,
		VerifySSL:   cfg.UiPath.VerifySSL,
		Timeout:     cfg.UiPath.Timeout,
	}
	uipathClient := rpa.NewUiPathHTTPClient(uipathConfig, logger)

	// 初始化仓库
	purchaseRepo := repository.NewPurchaseRequestRepository(db.GetDB(), logger)
	poRepo := repository.NewPurchaseOrderRepository(db.GetDB(), logger)
	supplierRepo := repository.NewSupplierRepository(db.GetDB(), logger)

	// 设置路由
	router := api.SetupFullRouter(queueManager, rpaClient, uipathClient, purchaseRepo, poRepo, supplierRepo, logger)

	// 启动服务器
	logger.WithField("port", cfg.Server.Port).Info("Starting server...")
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}

// loadConfig 加载配置
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
