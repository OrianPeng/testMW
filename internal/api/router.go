package api

import (
	"rpa-middleware/internal/interfaces"

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
		// 请求管理
		v1.POST("/requests", handler.SubmitRequest)
		v1.GET("/requests/:id", handler.GetRequestStatus)
		v1.GET("/requests", handler.ListRequests)
		
		// 系统状态
		v1.GET("/status", handler.GetQueueStats)
		v1.GET("/health", handler.HealthCheck)
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