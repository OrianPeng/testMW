package api

import (
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UiPathHandlers UiPath处理器
type UiPathHandlers struct {
	uipathClient interfaces.UiPathClient
	logger       *logrus.Logger
}

// NewUiPathHandlers 创建新的UiPath处理器
func NewUiPathHandlers(uipathClient interfaces.UiPathClient, logger *logrus.Logger) *UiPathHandlers {
	return &UiPathHandlers{
		uipathClient: uipathClient,
		logger:       logger,
	}
}

// AddQueueItem 添加项目到UiPath队列
// @Summary 添加项目到UiPath队列
// @Description 将指定的数据添加到UiPath Orchestrator队列中
// @Tags UiPath
// @Accept json
// @Produce json
// @Param request body models.UiPathAddQueueItemRequest true "队列项目请求"
// @Success 200 {object} models.UiPathAddQueueItemResponse "成功添加队列项目"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/uipath/queue/add [post]
func (h *UiPathHandlers) AddQueueItem(c *gin.Context) {
	var req models.UiPathAddQueueItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// 验证必需字段
	if req.QueueName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": "queue_name is required",
		})
		return
	}

	if req.SpecificContent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": "specific_content is required",
		})
		return
	}

	// 设置默认值
	if req.Priority == "" {
		req.Priority = "Normal"
	}

	h.logger.WithFields(logrus.Fields{
		"queue_name": req.QueueName,
		"priority":   req.Priority,
		"reference":  req.Reference,
	}).Info("Processing UiPath queue item request")

	// 调用UiPath客户端
	response, err := h.uipathClient.AddQueueItem(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add item to UiPath queue")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add queue item",
			"message": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id":    response.ID,
		"queue_name": response.QueueName,
		"status":     response.Status,
	}).Info("Successfully added item to UiPath queue")

	c.JSON(http.StatusOK, response)
}

// CheckUiPathStatus 检查UiPath系统状态
// @Summary 检查UiPath系统状态
// @Description 检查UiPath Orchestrator系统的可用性状态
// @Tags UiPath
// @Produce json
// @Success 200 {object} models.UiPathStatusResponse "系统状态信息"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/uipath/status [get]
func (h *UiPathHandlers) CheckUiPathStatus(c *gin.Context) {
	h.logger.Info("Checking UiPath system status")

	status, err := h.uipathClient.CheckUiPathStatus(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to check UiPath status")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check UiPath status",
			"message": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"is_available": status.IsAvailable,
		"last_check":   status.LastCheck,
		"message":      status.Message,
	}).Info("UiPath status checked")

	c.JSON(http.StatusOK, status)
}

// IsUiPathAvailable 检查UiPath系统是否可用
// @Summary 检查UiPath系统是否可用
// @Description 简单检查UiPath系统是否可用，返回布尔值
// @Tags UiPath
// @Produce json
// @Success 200 {object} map[string]bool "可用性状态"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/uipath/available [get]
func (h *UiPathHandlers) IsUiPathAvailable(c *gin.Context) {
	h.logger.Info("Checking if UiPath is available")

	available, err := h.uipathClient.IsUiPathAvailable(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to check UiPath availability")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check UiPath availability",
			"message": err.Error(),
		})
		return
	}

	h.logger.WithField("available", available).Info("UiPath availability checked")

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"timestamp": time.Now(),
	})
}

// TestUiPathConnection 测试UiPath连接
// @Summary 测试UiPath连接
// @Description 测试与UiPath Orchestrator的连接，包括认证和基本API访问
// @Tags UiPath
// @Produce json
// @Success 200 {object} map[string]interface{} "连接测试结果"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/uipath/test [get]
func (h *UiPathHandlers) TestUiPathConnection(c *gin.Context) {
	h.logger.Info("Testing UiPath connection")

	// 测试认证
	token, err := h.uipathClient.Authenticate(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("UiPath authentication failed")
		c.JSON(http.StatusOK, gin.H{
			"success":   false,
			"message":   "Authentication failed",
			"error":     err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	// 测试状态检查
	status, err := h.uipathClient.CheckUiPathStatus(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("UiPath status check failed")
		c.JSON(http.StatusOK, gin.H{
			"success":   false,
			"message":   "Status check failed",
			"error":     err.Error(),
			"token":     token[:min(len(token), 20)] + "...", // 只显示前20个字符
			"timestamp": time.Now(),
		})
		return
	}

	h.logger.Info("UiPath connection test successful")

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Connection test successful",
		"token_length": len(token),
		"is_available": status.IsAvailable,
		"last_check":   status.LastCheck,
		"timestamp":    time.Now(),
	})
}

// AddTestQueueItem 添加测试队列项目
// @Summary 添加测试队列项目
// @Description 使用默认测试数据添加一个队列项目到UiPath
// @Tags UiPath
// @Produce json
// @Param queue_name query string false "队列名称" default("CreationPR")
// @Param priority query string false "优先级" Enums(Normal, High, Critical) default("Normal")
// @Success 200 {object} models.UiPathAddQueueItemResponse "成功添加测试队列项目"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/uipath/queue/test [post]
func (h *UiPathHandlers) AddTestQueueItem(c *gin.Context) {
	// 获取查询参数
	queueName := c.DefaultQuery("queue_name", "CreationPR")
	priority := c.DefaultQuery("priority", "Normal")

	// 创建测试数据
	testData := map[string]interface{}{
		"Documento": "TEST-" + strconv.FormatInt(time.Now().Unix(), 10),
		"ClienteId": 123,
		"Total":     99.99,
		"Moeda":     "EUR",
		"TestMode":  true,
		"Timestamp": time.Now().Format(time.RFC3339),
	}

	req := &models.UiPathAddQueueItemRequest{
		QueueName:       queueName,
		Priority:        priority,
		SpecificContent: testData,
		Reference:       "Test-" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	h.logger.WithFields(logrus.Fields{
		"queue_name": req.QueueName,
		"priority":   req.Priority,
		"reference":  req.Reference,
	}).Info("Adding test item to UiPath queue")

	// 调用UiPath客户端
	response, err := h.uipathClient.AddQueueItem(c.Request.Context(), req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add test item to UiPath queue")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add test queue item",
			"message": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id":    response.ID,
		"queue_name": response.QueueName,
		"status":     response.Status,
	}).Info("Successfully added test item to UiPath queue")

	c.JSON(http.StatusOK, response)
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
