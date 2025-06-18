package processor

import (
	"context"
	"fmt"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// QueueProcessor 队列处理器
type QueueProcessor struct {
	queueManager        interfaces.QueueManager
	rpaClient          interfaces.RPAClient
	notificationService interfaces.NotificationService
	logger             *logrus.Logger
	
	// 配置参数
	checkInterval   time.Duration // 检查间隔
	maxRetries      int           // 最大重试次数
	retryInterval   time.Duration // 重试间隔
	
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
		rpaClient:          rpaClient,
		notificationService: notificationService,
		logger:             logger,
		checkInterval:      5 * time.Second,
		maxRetries:         3,
		retryInterval:      10 * time.Second,
		stopCh:             make(chan struct{}),
		doneCh:             make(chan struct{}),
	}
}

// Start 启动队列处理器
func (p *QueueProcessor) Start(ctx context.Context) {
	p.logger.Info("Starting queue processor")
	
	go p.processLoop(ctx)
}

// Stop 停止队列处理器
func (p *QueueProcessor) Stop() {
	p.logger.Info("Stopping queue processor")
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
			p.processRequests(ctx)
		}
	}
}

// processRequests 处理队列中的请求
func (p *QueueProcessor) processRequests(ctx context.Context) {
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
		p.logger.Debug("No pending requests to process")
		return
	}

	p.logger.WithField("pending_count", pendingCount).Info("Processing pending requests")

	// 逐个处理请求
	for {
		req, err := p.queueManager.DequeueRequest(ctx)
		if err != nil {
			// 没有更多请求或出现错误
			break
		}

		p.processRequest(ctx, req)

		// 检查 RPA 是否仍然可用
		if available, _ := p.rpaClient.IsAvailable(ctx); !available {
			p.logger.Info("RPA system became unavailable during processing")
			break
		}
	}
}

// processRequest 处理单个请求
func (p *QueueProcessor) processRequest(ctx context.Context, req *models.QueuedRequest) {
	logger := p.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"agent_id":   req.AgentID,
		"retry":      req.RetryCount,
	})

	logger.Info("Processing request")

	// 构建 RPA 请求
	rpaReq := &models.RPARequest{
		ID:      req.ID,
		AgentID: req.AgentID,
		Data:    req.Data,
		Metadata: map[string]string{
			"retry_count": string(rune(req.RetryCount)),
			"created_at":  req.CreatedAt.Format(time.RFC3339),
		},
	}

	// 发送到 RPA 系统
	result, err := p.rpaClient.SendRequest(ctx, rpaReq)
	if err != nil {
		logger.WithError(err).Error("Failed to send request to RPA")
		p.handleRequestError(ctx, req, err)
		return
	}

	// 处理结果
	if result.Success {
		logger.Info("Request processed successfully")
		
		// 更新状态为完成
		if err := p.queueManager.UpdateRequestStatus(ctx, req.ID, models.StatusCompleted, ""); err != nil {
			logger.WithError(err).Error("Failed to update request status to completed")
		}

		// 发送完成通知
		if err := p.notificationService.NotifyCompletion(ctx, req, result); err != nil {
			logger.WithError(err).Warn("Failed to send completion notification")
		}
	} else {
		logger.WithField("error", result.Error).Warn("RPA processing failed")
		p.handleRequestError(ctx, req, fmt.Errorf("RPA processing failed: %s", result.Error))
	}
}

// handleRequestError 处理请求错误
func (p *QueueProcessor) handleRequestError(ctx context.Context, req *models.QueuedRequest, err error) {
	logger := p.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"retry":      req.RetryCount,
		"max_retry":  p.maxRetries,
	})

	if req.RetryCount < p.maxRetries {
		// 可以重试
		logger.Info("Request will be retried")
		
		// 重新加入队列（状态改为 pending）
		req.Status = models.StatusPending
		if queueErr := p.queueManager.EnqueueRequest(ctx, req.AgentRequest); queueErr != nil {
			logger.WithError(queueErr).Error("Failed to re-enqueue request for retry")
		}
	} else {
		// 超过最大重试次数，标记为失败
		logger.Error("Request failed after maximum retries")
		
		if updateErr := p.queueManager.UpdateRequestStatus(ctx, req.ID, models.StatusFailed, err.Error()); updateErr != nil {
			logger.WithError(updateErr).Error("Failed to update request status to failed")
		}

		// 发送失败通知
		if notifyErr := p.notificationService.NotifyCompletion(ctx, req, nil); notifyErr != nil {
			logger.WithError(notifyErr).Warn("Failed to send failure notification")
		}
	}
}