package processor

import (
	"context"
	"fmt"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// QueueProcessor 队列处理器
type QueueProcessor struct {
	queueManager        interfaces.QueueManager
	rpaClient           interfaces.RPAClient
	notificationService interfaces.NotificationService
	logger              *logrus.Logger

	// 配置参数
	checkInterval time.Duration // 检查间隔
	maxRetries    int           // 最大重试次数
	retryInterval time.Duration // 重试间隔

	// 控制
	stopCh chan struct{}
	doneCh chan struct{}
}

// NewQueueProcessor 创建新的队列处理器
func NewQueueProcessor(
	queueManager interfaces.QueueManager,
	rpaClient interfaces.RPAClient,
	notificationService interfaces.NotificationService,
	logger *logrus.Logger,
) *QueueProcessor {
	return &QueueProcessor{
		queueManager:        queueManager,
		rpaClient:           rpaClient,
		notificationService: notificationService,
		logger:              logger,
		checkInterval:       5 * time.Second,
		maxRetries:          3,
		retryInterval:       10 * time.Second,
		stopCh:              make(chan struct{}),
		doneCh:              make(chan struct{}),
	}
}

// Start 启动队列处理器
func (p *QueueProcessor) Start(ctx context.Context) {
	p.logger.Info("Starting purchase request queue processor")

	go p.processLoop(ctx)
}

// Stop 停止队列处理器
func (p *QueueProcessor) Stop() {
	p.logger.Info("Stopping purchase request queue processor")
	close(p.stopCh)
	<-p.doneCh
}

// processLoop 处理循环
func (p *QueueProcessor) processLoop(ctx context.Context) {
	defer close(p.doneCh)

	ticker := time.NewTicker(p.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Context cancelled, stopping processor")
			return
		case <-p.stopCh:
			p.logger.Info("Stop signal received, stopping processor")
			return
		case <-ticker.C:
			p.processPurchaseRequests(ctx)
		}
	}
}

// processPurchaseRequests 处理队列中的采购请求
func (p *QueueProcessor) processPurchaseRequests(ctx context.Context) {
	// 检查 RPA 系统是否可用
	available, err := p.rpaClient.IsAvailable(ctx)
	if err != nil {
		p.logger.WithError(err).Warn("Failed to check RPA availability")
		return
	}

	if !available {
		p.logger.Debug("RPA system is not available, skipping request processing")
		return
	}

	// 获取待处理请求数量
	pendingCount, err := p.queueManager.GetPendingCount(ctx)
	if err != nil {
		p.logger.WithError(err).Error("Failed to get pending count")
		return
	}

	if pendingCount == 0 {
		p.logger.Debug("No pending purchase requests to process")
		return
	}

	p.logger.WithField("pending_count", pendingCount).Info("Processing pending purchase requests")

	// 逐个处理请求
	for {
		req, err := p.queueManager.DequeuePurchaseRequest(ctx)
		if err != nil {
			// 没有更多请求或出现错误
			break
		}

		p.processPurchaseRequest(ctx, req)

		// 检查 RPA 是否仍然可用
		if available, _ := p.rpaClient.IsAvailable(ctx); !available {
			p.logger.Info("RPA system became unavailable during processing")
			break
		}
	}
}

// processPurchaseRequest 处理单个采购请求
func (p *QueueProcessor) processPurchaseRequest(ctx context.Context, req *models.QueuedPurchaseRequest) {
	logger := p.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"requester":  req.Requester,
		"retry":      req.RetryCount,
	})

	logger.Info("Processing purchase request")

	// 构建 RPA 请求
	rpaReq := &models.RPARequest{
		ID:        req.RequestID,
		RequestID: req.RequestID,
		Data: map[string]interface{}{
			"requester":             req.Requester,
			"doc_type":              req.DocType,
			"plant":                 req.Plant,
			"material":              req.Material,
			"quantity":              req.Quantity,
			"unit_price":            req.UnitPrice,
			"total_amount":          req.TotalAmount,
			"delivery_date":         req.DeliveryDate,
			"vendor_code":           req.VendorCode,
			"short_text":            req.ShortText,
			"material_group":        req.MaterialGroup,
			"unit_type":             req.UnitType,
			"purchase_organization": req.PurchaseOrganization,
			"currency":              req.Currency,
			"priority":              req.Priority,
			"status":                req.Status,
			"urgency":               req.Urgency,
			"comments":              req.Comments,
		},
		Metadata: map[string]string{
			"retry_count":    strconv.Itoa(req.RetryCount),
			"created_at":     req.CreatedAt.Format(time.RFC3339),
			"execution_mode": "queued",
		},
	}

	// 发送到 RPA 系统
	result, err := p.rpaClient.SendPurchaseRequest(ctx, rpaReq)
	if err != nil {
		logger.WithError(err).Error("Failed to send purchase request to RPA")
		p.handlePurchaseRequestError(ctx, req, err)
		return
	}

	// 处理结果
	if result.Success {
		logger.Info("Purchase request processed successfully")

		// 更新状态为已完成
		if err := p.queueManager.UpdatePurchaseRequestStatus(ctx, req.RequestID, models.PurchaseStatusCompleted, ""); err != nil {
			logger.WithError(err).Error("Failed to update request status to completed")
		}

		// 删除队列中的请求
		if err := p.queueManager.DeletePurchaseRequest(ctx, req.RequestID); err != nil {
			logger.WithError(err).Error("Failed to delete request from queue")
		} else {
			logger.Info("Purchase request deleted from queue after successful processing")
		}

		// 发送完成通知
		if err := p.notificationService.NotifyPurchaseRequestCompletion(ctx, req, result); err != nil {
			logger.WithError(err).Warn("Failed to send completion notification")
		}
	} else {
		logger.WithField("error", result.Error).Warn("RPA processing failed")
		p.handlePurchaseRequestError(ctx, req, fmt.Errorf("RPA processing failed: %s", result.Error))
	}
}

// handlePurchaseRequestError 处理采购请求错误
func (p *QueueProcessor) handlePurchaseRequestError(ctx context.Context, req *models.QueuedPurchaseRequest, err error) {
	logger := p.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"retry":      req.RetryCount,
		"max_retry":  p.maxRetries,
	})

	if req.RetryCount < p.maxRetries {
		// 可以重试
		logger.Info("Purchase request will be retried")

		// 增加重试次数
		req.RetryCount++
		req.Status = models.PurchaseStatusProcessing

		// 重新加入队列
		if queueErr := p.queueManager.EnqueuePurchaseRequest(ctx, req.PurchaseRequest); queueErr != nil {
			logger.WithError(queueErr).Error("Failed to re-enqueue purchase request for retry")
		}
	} else {
		// 超过最大重试次数，标记为失败
		logger.Error("Purchase request failed after maximum retries")

		// 更新状态为失败
		if err := p.queueManager.UpdatePurchaseRequestStatus(ctx, req.RequestID, models.PurchaseStatusFailed, err.Error()); err != nil {
			logger.WithError(err).Error("Failed to update request status to failed")
		}

		// 删除队列中的请求
		if err := p.queueManager.DeletePurchaseRequest(ctx, req.RequestID); err != nil {
			logger.WithError(err).Error("Failed to delete failed request from queue")
		}

		// 由于已移除callback字段，不再发送失败通知
		logger.Info("Callback functionality removed, skipping failure notification")
	}
}
