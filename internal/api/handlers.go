package api

import (
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Handler API 处理器
type Handler struct {
	queueManager interfaces.QueueManager
	rpaClient    interfaces.RPAClient
	logger       *logrus.Logger
}

// NewHandler 创建新的 API 处理器
func NewHandler(queueManager interfaces.QueueManager, rpaClient interfaces.RPAClient, logger *logrus.Logger) *Handler {
	return &Handler{
		queueManager: queueManager,
		rpaClient:    rpaClient,
		logger:       logger,
	}
}

// SubmitRequest 提交请求到队列
func (h *Handler) SubmitRequest(c *gin.Context) {
	var req models.AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 生成请求 ID（如果没有提供）
	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	// 设置创建时间
	req.CreatedAt = time.Now()

	// 记录日志
	h.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"agent_id":   req.AgentID,
		"priority":   req.Priority,
	}).Info("Received request from AI agent")

	// 加入队列
	if err := h.queueManager.EnqueueRequest(c.Request.Context(), &req); err != nil {
		h.logger.WithError(err).Error("Failed to enqueue request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to enqueue request",
		})
		return
	}

	// 返回响应
	response := models.AgentResponse{
		RequestID: req.ID,
		Status:    models.StatusPending,
		Message:   "Request queued successfully",
	}

	c.JSON(http.StatusAccepted, response)
}

// GetRequestStatus 获取请求状态
func (h *Handler) GetRequestStatus(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request ID is required",
		})
		return
	}

	req, err := h.queueManager.GetRequestStatus(c.Request.Context(), requestID)
	if err != nil {
		h.logger.WithError(err).WithField("request_id", requestID).Warn("Request not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Request not found",
		})
		return
	}

	response := models.AgentResponse{
		RequestID: req.ID,
		Status:    req.Status,
		Message:   h.getStatusMessage(req.Status),
	}

	if req.Status == models.StatusFailed && req.ErrorMsg != "" {
		response.Message = req.ErrorMsg
	}

	c.JSON(http.StatusOK, response)
}

// ListRequests 列出请求
func (h *Handler) ListRequests(c *gin.Context) {
	// 解析查询参数
	status := models.RequestStatus(c.Query("status"))
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// 获取请求列表
	requests, err := h.queueManager.ListRequests(c.Request.Context(), status, limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list requests")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list requests",
		})
		return
	}

	// 转换为响应格式
	var responses []models.AgentResponse
	for _, req := range requests {
		response := models.AgentResponse{
			RequestID: req.ID,
			Status:    req.Status,
			Message:   h.getStatusMessage(req.Status),
			Data: map[string]interface{}{
				"agent_id":    req.AgentID,
				"priority":    req.Priority,
				"retry_count": req.RetryCount,
				"created_at":  req.CreatedAt,
				"updated_at":  req.UpdatedAt,
			},
		}

		if req.Status == models.StatusFailed && req.ErrorMsg != "" {
			response.Message = req.ErrorMsg
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": responses,
		"metadata": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(responses),
		},
	})
}

// GetQueueStats 获取队列统计信息
func (h *Handler) GetQueueStats(c *gin.Context) {
	pendingCount, err := h.queueManager.GetPendingCount(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pending count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get queue statistics",
		})
		return
	}

	// 获取 RPA 状态
	rpaStatus, err := h.rpaClient.CheckStatus(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Warn("Failed to check RPA status")
		rpaStatus = &models.RPAStatus{
			IsAvailable: false,
			LastCheck:   time.Now(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"queue": gin.H{
			"pending_requests": pendingCount,
		},
		"rpa_status": rpaStatus,
		"timestamp":  time.Now(),
	})
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

// getStatusMessage 获取状态消息
func (h *Handler) getStatusMessage(status models.RequestStatus) string {
	switch status {
	case models.StatusPending:
		return "Request is pending processing"
	case models.StatusProcessing:
		return "Request is being processed"
	case models.StatusCompleted:
		return "Request completed successfully"
	case models.StatusFailed:
		return "Request processing failed"
	default:
		return "Unknown status"
	}
}