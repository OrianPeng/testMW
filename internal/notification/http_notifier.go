package notification

import (
	"context"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPNotifier 实现HTTP通知服务
type HTTPNotifier struct {
	timeout time.Duration
	logger  *logrus.Logger
}

// NewHTTPNotifier 创建新的HTTP通知器
func NewHTTPNotifier(timeout time.Duration, logger *logrus.Logger) interfaces.NotificationService {
	return &HTTPNotifier{
		timeout: timeout,
		logger:  logger,
	}
}

// SendPurchaseRequestCallback 发送采购请求回调通知
func (h *HTTPNotifier) SendPurchaseRequestCallback(ctx context.Context, callbackURL string, response *models.PurchaseRequestCallbackResponse) error {
	h.logger.WithFields(logrus.Fields{
		"callback_url": callbackURL,
		"request_id":   response.RequestID,
	}).Info("Sending purchase request callback")
	// TODO: 实现实际的HTTP回调逻辑
	return nil
}

// NotifyPurchaseRequestCompletion 通知采购请求完成
func (h *HTTPNotifier) NotifyPurchaseRequestCompletion(ctx context.Context, req *models.QueuedPurchaseRequest, result *models.RPAResponse) error {
	h.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"success":    result.Success,
	}).Info("Notifying purchase request completion")
	// TODO: 实现实际的通知逻辑
	return nil
}
