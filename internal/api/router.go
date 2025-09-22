package api

import (
	"fmt"
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter 设置路由
func SetupRouter(queueManager interfaces.QueueManager, rpaClient interfaces.RPAClient, logger *logrus.Logger) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())

	// 创建处理器
	handler := NewHandler(queueManager, rpaClient, logger)

	// API 路由组
	v1 := router.Group("/api/v1")
	{
		// 采购请求队列管理
		v1.POST("/purchase-requests/queue", handler.SubmitPurchaseRequest)
		v1.POST("/purchase-requests/queue/clear-all", handler.ClearQueue)
		v1.GET("/purchase-requests/queue", handler.ListPurchaseRequests)
		v1.GET("/purchase-requests/queue/:id", handler.GetPurchaseRequestStatus)

		// 系统状态
		v1.GET("/status", handler.GetQueueStats)
		v1.GET("/health", handler.HealthCheck)

	}

	return router
}

// SetupRouterWithPurchaseRequests 设置包含采购请求的路由
func SetupRouterWithPurchaseRequests(
	queueManager interfaces.QueueManager,
	rpaClient interfaces.RPAClient,
	purchaseRepo *repository.PurchaseRequestRepository,
	logger *logrus.Logger,
) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())

	// 创建处理器
	handler := NewHandler(queueManager, rpaClient, logger)
	purchaseHandler := NewPurchaseRequestHandler(purchaseRepo, logger)

	// API 路由组
	v1 := router.Group("/api/v1")
	{
		// 采购请求队列管理
		v1.POST("/purchase-requests/queue", handler.SubmitPurchaseRequest)
		v1.POST("/purchase-requests/queue/clear-all", handler.ClearQueue)
		v1.GET("/purchase-requests/queue", handler.ListPurchaseRequests)
		v1.GET("/purchase-requests/queue/:id", handler.GetPurchaseRequestStatus)

		// 系统状态
		v1.GET("/status", handler.GetQueueStats)
		v1.GET("/health", handler.HealthCheck)

		// 采购请求数据库管理
		purchase := v1.Group("/purchase-requests")
		{
			purchase.POST("", purchaseHandler.CreatePurchaseRequest)
			purchase.GET("/:request_id", purchaseHandler.GetPurchaseRequest)
			purchase.GET("", purchaseHandler.ListPurchaseRequests)
			purchase.PUT("/:request_id", purchaseHandler.UpdatePurchaseRequest)
			purchase.DELETE("/:request_id", purchaseHandler.DeletePurchaseRequest)
		}
	}

	return router
}

// SetupFullRouter 设置完整的API路由（包含所有功能）
func SetupFullRouter(
	queueManager interfaces.QueueManager,
	rpaClient interfaces.RPAClient,
	uipathClient interfaces.UiPathClient,
	purchaseRepo *repository.PurchaseRequestRepository,
	poRepo *repository.PurchaseOrderRepository,
	supplierRepo *repository.SupplierRepository,
	logger *logrus.Logger,
) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())

	// 创建处理器
	handler := NewHandler(queueManager, rpaClient, logger)
	purchaseHandler := NewPurchaseRequestHandler(purchaseRepo, logger)
	poHandler := NewPurchaseOrderHandler(poRepo, logger)
	supplierHandler := NewSupplierHandler(supplierRepo, logger)
	dashboardHandler := NewDashboardHandler(purchaseRepo, poRepo, logger)
	processHandler := NewProcessMonitorHandler(logger)
	aiHandler := NewAIAssistantHandler(logger)
	uipathHandler := NewUiPathHandlers(uipathClient, logger)

	// API 路由组
	api := router.Group("/api/v1")
	{
		// 采购请求管理
		purchase := api.Group("/purchase-requests")
		{
			purchase.POST("", purchaseHandler.CreatePurchaseRequest)
			purchase.GET("/:request_id", purchaseHandler.GetPurchaseRequest)
			purchase.GET("", purchaseHandler.ListPurchaseRequests)
			purchase.PUT("/:request_id", purchaseHandler.UpdatePurchaseRequest)
			purchase.DELETE("/:request_id", purchaseHandler.DeletePurchaseRequest)
			purchase.GET("/statistics", purchaseHandler.GetPRStatistics)
		}

		// 采购订单管理
		po := api.Group("/purchase-orders")
		{
			po.POST("", poHandler.CreatePurchaseOrder)
			po.GET("/:po_number", poHandler.GetPurchaseOrder)
			po.GET("", poHandler.ListPurchaseOrders)
			po.PUT("/:po_number", poHandler.UpdatePurchaseOrder)
			po.DELETE("/:po_number", poHandler.DeletePurchaseOrder)
			po.GET("/statistics", poHandler.GetPOStatistics)
		}

		// 供应商管理
		suppliers := api.Group("/suppliers")
		{
			suppliers.POST("", supplierHandler.CreateSupplier)
			suppliers.GET("/:id", supplierHandler.GetSupplier)
			suppliers.GET("", supplierHandler.ListSuppliers)
			suppliers.PUT("/:id", supplierHandler.UpdateSupplier)
			suppliers.DELETE("/:id", supplierHandler.DeleteSupplier)
			suppliers.GET("/statistics", supplierHandler.GetSupplierStatistics)
		}

		// 流程监控
		process := api.Group("/process-monitor")
		{
			process.GET("/overview", processHandler.GetProcessOverview)
			process.GET("/blocked", processHandler.GetBlockedProcesses)
			process.GET("/steps", processHandler.GetProcessSteps)
		}

		// 仪表板
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/metrics", dashboardHandler.GetDashboardMetrics)
			dashboard.GET("/charts/pr-trend", dashboardHandler.GetPRTrendChart)
			dashboard.GET("/charts/po-status", dashboardHandler.GetPOStatusDistribution)
			dashboard.GET("/charts/supplier-ranking", dashboardHandler.GetSupplierRanking)
			dashboard.GET("/charts/department-stats", dashboardHandler.GetDepartmentStats)
		}

		// AI助手
		ai := api.Group("/ai-assistant")
		{
			ai.POST("/message", aiHandler.SendMessage)
		}

		// UiPath RPA集成
		uipath := api.Group("/uipath")
		{
			uipath.POST("/queue/add", uipathHandler.AddQueueItem)
			uipath.GET("/status", uipathHandler.CheckUiPathStatus)
			uipath.GET("/available", uipathHandler.IsUiPathAvailable)
			uipath.GET("/test", uipathHandler.TestUiPathConnection)
			uipath.POST("/queue/test", uipathHandler.AddTestQueueItem)
		}

		// 文件上传
		api.POST("/upload", func(c *gin.Context) {
			// 简单的文件上传处理
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "VALIDATION_ERROR",
						"message": "No file uploaded",
					},
				})
				return
			}

			// 生成文件ID
			fileID := "file_" + fmt.Sprintf("%d", time.Now().Unix())

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"file_id":  fileID,
					"filename": file.Filename,
					"size":     file.Size,
					"url":      "/files/" + fileID,
				},
			})
		})

		// 队列管理（保持原有功能）
		queue := api.Group("/v1")
		{
			queue.POST("/purchase-requests/queue", handler.SubmitPurchaseRequest)
			queue.POST("/purchase-requests/queue/clear-all", handler.ClearQueue)
			queue.GET("/purchase-requests/queue", handler.ListPurchaseRequests)
			queue.GET("/purchase-requests/queue/:id", handler.GetPurchaseRequestStatus)
			queue.GET("/status", handler.GetQueueStats)
			queue.GET("/health", handler.HealthCheck)
		}
	}

	return router
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 记录请求
		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Info("API request received")

		c.Next()

		// 记录响应
		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		}).Info("API request completed")
	})
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
