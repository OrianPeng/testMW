package api

import (
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler API处理器
type Handler struct {
	queueManager interfaces.QueueManager
	rpaClient    interfaces.RPAClient
	logger       *logrus.Logger
}

// NewHandler 创建新的处理器
func NewHandler(queueManager interfaces.QueueManager, rpaClient interfaces.RPAClient, logger *logrus.Logger) *Handler {
	return &Handler{
		queueManager: queueManager,
		rpaClient:    rpaClient,
		logger:       logger,
	}
}

// SubmitPurchaseRequest 提交采购请求到队列
func (h *Handler) SubmitPurchaseRequest(c *gin.Context) {
	var req models.CreatePurchaseRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 计算总金额
	totalAmount := float64(req.Quantity) * req.UnitPrice

	// 设置默认货币
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	// 创建采购请求
	purchaseReq := &models.PurchaseRequest{
		RequestID:            req.RequestID,
		DocType:              req.DocType,
		Plant:                req.Plant,
		Quantity:             req.Quantity,
		UnitPrice:            req.UnitPrice,
		Material:             req.Material,
		DeliveryDate:         *req.DeliveryDate,
		VendorCode:           req.VendorCode,
		ShortText:            req.ShortText,
		MaterialGroup:        req.MaterialGroup,
		UnitType:             req.UnitType,
		Requester:            req.Requester,
		PurchaseOrganization: req.PurchaseOrganization,
		Currency:             currency,
		TotalAmount:          totalAmount,
		Priority:             req.Priority,
		Status:               models.PurchaseStatusPending,
		Urgency:              req.Urgency,
		Comments:             req.Comments,
		RetryCount:           0,
	}

	// 将请求加入队列
	if err := h.queueManager.EnqueuePurchaseRequest(c.Request.Context(), purchaseReq); err != nil {
		h.logger.WithError(err).Error("Failed to enqueue purchase request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit request to queue",
		})
		return
	}

	h.logger.WithField("request_id", req.RequestID).Info("Purchase request submitted to queue")

	c.JSON(http.StatusOK, gin.H{
		"request_id": req.RequestID,
		"status":     models.PurchaseStatusProcessing,
		"message":    "Purchase request queued successfully",
	})
}

// GetPurchaseRequestStatus 获取采购请求状态
func (h *Handler) GetPurchaseRequestStatus(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	// 从队列中获取请求状态
	queuedReq, err := h.queueManager.GetPurchaseRequestStatus(c.Request.Context(), requestID)
	if err != nil {
		h.logger.WithError(err).WithField("request_id", requestID).Error("Failed to get request status")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Request not found in queue",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id":     queuedReq.RequestID,
		"status":         queuedReq.Status,
		"priority":       queuedReq.Priority,
		"queue_position": queuedReq.QueuePosition,
		"enqueued_at":    queuedReq.EnqueuedAt,
		"retry_count":    queuedReq.RetryCount,
		"error_msg":      queuedReq.ErrorMsg,
		"processed_at":   queuedReq.ProcessedAt,
	})
}

// ListPurchaseRequests 列出队列中的采购请求
func (h *Handler) ListPurchaseRequests(c *gin.Context) {
	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	statusStr := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// 从队列中获取请求列表
	requests, err := h.queueManager.ListPurchaseRequests(c.Request.Context(), statusStr, limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list purchase requests")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list requests",
		})
		return
	}

	// 获取待处理数量
	pendingCount, err := h.queueManager.GetPendingCount(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pending count")
	}

	// 构建响应
	var responseRequests []gin.H
	for _, req := range requests {
		responseRequests = append(responseRequests, gin.H{
			"request_id":            req.RequestID,
			"doc_type":              req.DocType,
			"plant":                 req.Plant,
			"material":              req.Material,
			"quantity":              req.Quantity,
			"unit_price":            req.UnitPrice,
			"total_amount":          req.TotalAmount,
			"currency":              req.Currency,
			"vendor_code":           req.VendorCode,
			"material_group":        req.MaterialGroup,
			"requester":             req.Requester,
			"purchase_organization": req.PurchaseOrganization,
			"priority":              req.Priority,
			"status":                req.Status,
			"urgency":               req.Urgency,
			"queue_position":        req.QueuePosition,
			"enqueued_at":           req.EnqueuedAt,
			"retry_count":           req.RetryCount,
			"error_msg":             req.ErrorMsg,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"requests":      responseRequests,
		"pending_count": pendingCount,
		"limit":         limit,
		"offset":        offset,
	})
}

// GetQueueStats 获取队列统计信息
func (h *Handler) GetQueueStats(c *gin.Context) {
	// 获取待处理数量
	pendingCount, err := h.queueManager.GetPendingCount(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pending count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get queue stats",
		})
		return
	}

	// 检查RPA系统状态
	rpaAvailable, err := h.rpaClient.IsAvailable(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to check RPA availability")
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_stats": gin.H{
			"pending_count": pendingCount,
		},
		"rpa_status": gin.H{
			"available": rpaAvailable,
		},
	})
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	// 检查队列连接
	_, err := h.queueManager.GetPendingCount(c.Request.Context())
	queueHealthy := err == nil

	// 检查RPA连接
	rpaHealthy := false
	if available, err := h.rpaClient.IsAvailable(c.Request.Context()); err == nil {
		rpaHealthy = available
	}

	overallHealth := queueHealthy && rpaHealthy
	statusCode := http.StatusOK
	if !overallHealth {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status": "healthy",
		"components": gin.H{
			"queue": gin.H{
				"healthy": queueHealthy,
				"error":   err,
			},
			"rpa": gin.H{
				"healthy": rpaHealthy,
			},
		},
	})
}

// ClearQueue 清空队列
func (h *Handler) ClearQueue(c *gin.Context) {
	if err := h.queueManager.ClearAll(c.Request.Context()); err != nil {
		h.logger.WithError(err).Error("Failed to clear queue")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear queue",
		})
		return
	}

	h.logger.Info("Queue cleared")
	c.JSON(http.StatusOK, gin.H{
		"message": "Queue cleared successfully",
	})
}
