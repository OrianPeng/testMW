package rpa

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

// HTTPClient RPA HTTP客户端
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(baseURL string, timeout time.Duration, logger *logrus.Logger) interfaces.RPAClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// SendPurchaseRequest 发送采购请求到RPA系统
func (c *HTTPClient) SendPurchaseRequest(ctx context.Context, req *models.RPARequest) (*models.RPAResponse, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/api/v1/execute", c.baseURL)

	// 构建RPA模拟器期望的请求格式
	rpaReq := map[string]interface{}{
		"id":       req.ID,
		"agent_id": "middleware-agent",
		"data": map[string]interface{}{
			"action": "process_purchase_request",
			"params": req.Data,
		},
		"metadata": req.Metadata,
	}

	// 序列化请求
	reqBody, err := json.Marshal(rpaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return &models.RPAResponse{
			ID:        req.ID,
			RequestID: req.RequestID,
			Success:   false,
			Error:     fmt.Sprintf("HTTP error: %d", resp.StatusCode),
		}, nil
	}

	// 解析响应
	var rpaResp models.RPAResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 设置请求ID
	rpaResp.ID = req.ID
	rpaResp.RequestID = req.RequestID

	c.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"success":    rpaResp.Success,
		"error":      rpaResp.Error,
	}).Info("RPA purchase request sent")

	return &rpaResp, nil
}

// CheckStatus 检查RPA系统状态
func (c *HTTPClient) CheckStatus(ctx context.Context) (*models.RPAStatus, error) {
	url := fmt.Sprintf("%s/api/v1/status", c.baseURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to check RPA status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.RPAStatus{
			IsAvailable: false,
			LastCheck:   time.Now(),
		}, nil
	}

	var status models.RPAStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	status.LastCheck = time.Now()

	return &status, nil
}

// IsAvailable 检查RPA系统是否可用
func (c *HTTPClient) IsAvailable(ctx context.Context) (bool, error) {
	status, err := c.CheckStatus(ctx)
	if err != nil {
		return false, err
	}

	return status.IsAvailable, nil
}
