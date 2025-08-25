package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPNotifier HTTP通知服务
type HTTPNotifier struct {
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewHTTPNotifier 创建新的HTTP通知服务
func NewHTTPNotifier(timeout time.Duration, logger *logrus.Logger) interfaces.NotificationService {
	return &HTTPNotifier{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// SendPurchaseRequestCallback 发送采购请求回调通知
func (n *HTTPNotifier) SendPurchaseRequestCallback(ctx context.Context, callbackURL string, response *models.PurchaseRequestCallbackResponse) error {
	// 序列化响应
	responseBody, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal callback response: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, bytes.NewBuffer(responseBody))
	if err != nil {
		return fmt.Errorf("failed to create callback request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Callback-Type", "purchase-request")
	req.Header.Set("X-Request-ID", response.RequestID)

	// 发送请求
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("callback failed with status %d", resp.StatusCode)
	}

	n.logger.WithFields(logrus.Fields{
		"request_id":   response.RequestID,
		"status":       response.Status,
		"callback_url": callbackURL,
		"http_status":  resp.StatusCode,
	}).Info("Purchase request callback sent successfully")

	return nil
}

// NotifyPurchaseRequestCompletion 通知采购请求完成
func (n *HTTPNotifier) NotifyPurchaseRequestCompletion(ctx context.Context, req *models.QueuedPurchaseRequest, result *models.RPAResponse) error {
	// 由于已移除callback字段，直接返回
	n.logger.WithField("request_id", req.RequestID).Debug("Callback functionality removed, skipping notification")
	return nil
}
