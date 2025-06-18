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

// HTTPNotifier HTTP 通知服务
type HTTPNotifier struct {
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewHTTPNotifier 创建新的 HTTP 通知服务
func NewHTTPNotifier(timeout time.Duration, logger *logrus.Logger) interfaces.NotificationService {
	return &HTTPNotifier{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// SendCallback 发送回调通知
func (n *HTTPNotifier) SendCallback(ctx context.Context, callbackURL string, response *models.AgentResponse) error {
	if callbackURL == "" {
		return nil // 没有回调 URL，跳过通知
	}

	n.logger.WithFields(logrus.Fields{
		"request_id":   response.RequestID,
		"callback_url": callbackURL,
		"status":       response.Status,
	}).Info("Sending callback notification")

	// 序列化响应
	reqBody, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal callback response: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", callbackURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create callback request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Callback-Type", "rpa-middleware")
	httpReq.Header.Set("X-Request-ID", response.RequestID)

	// 发送请求
	resp, err := n.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("callback failed with status %d", resp.StatusCode)
	}

	n.logger.WithFields(logrus.Fields{
		"request_id":   response.RequestID,
		"callback_url": callbackURL,
	}).Info("Callback notification sent successfully")

	return nil
}

// NotifyCompletion 通知请求完成
func (n *HTTPNotifier) NotifyCompletion(ctx context.Context, req *models.QueuedRequest, result *models.RPAResponse) error {
	// 构建响应
	response := &models.AgentResponse{
		RequestID: req.ID,
		Status:    req.Status,
		Message:   "Request processed successfully",
	}

	if result != nil {
		if result.Success {
			response.Data = result.Data
		} else {
			response.Status = models.StatusFailed
			response.Message = result.Error
		}
	}

	// 发送回调通知
	return n.SendCallback(ctx, req.Callback, response)
}