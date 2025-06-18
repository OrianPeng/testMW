package rpa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPClient RPA HTTP 客户端
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewHTTPClient 创建新的 RPA HTTP 客户端
func NewHTTPClient(baseURL string, timeout time.Duration, logger *logrus.Logger) interfaces.RPAClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// SendRequest 发送请求到 RPA 系统
func (c *HTTPClient) SendRequest(ctx context.Context, req *models.RPARequest) (*models.RPAResponse, error) {
	c.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"agent_id":   req.AgentID,
	}).Info("Sending request to RPA system")

	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/execute", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", req.ID)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to RPA: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RPA system returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var rpaResp models.RPAResponse
	if err := json.Unmarshal(respBody, &rpaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"success":    rpaResp.Success,
	}).Info("Received response from RPA system")

	return &rpaResp, nil
}

// CheckStatus 检查 RPA 系统状态
func (c *HTTPClient) CheckStatus(ctx context.Context) (*models.RPAStatus, error) {
	c.logger.Debug("Checking RPA system status")

	// 创建状态检查请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/status", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// 如果连接失败，认为 RPA 不可用
		return &models.RPAStatus{
			IsAvailable: false,
			LastCheck:   time.Now(),
		}, nil
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.RPAStatus{
			IsAvailable: false,
			LastCheck:   time.Now(),
		}, nil
	}

	// 解析响应
	var status models.RPAStatus
	if err := json.Unmarshal(respBody, &status); err != nil {
		// 如果解析失败，但连接成功，认为部分可用
		return &models.RPAStatus{
			IsAvailable: resp.StatusCode == http.StatusOK,
			LastCheck:   time.Now(),
		}, nil
	}

	status.LastCheck = time.Now()
	
	c.logger.WithFields(logrus.Fields{
		"is_available": status.IsAvailable,
		"current_task": status.CurrentTask,
	}).Debug("RPA system status checked")

	return &status, nil
}

// IsAvailable 检查 RPA 系统是否可用
func (c *HTTPClient) IsAvailable(ctx context.Context) (bool, error) {
	status, err := c.CheckStatus(ctx)
	if err != nil {
		return false, err
	}
	return status.IsAvailable, nil
}