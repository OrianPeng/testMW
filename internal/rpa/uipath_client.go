package rpa

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// UiPathHTTPClient UiPath HTTP客户端
type UiPathHTTPClient struct {
	config     *models.UiPathConfig
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewUiPathHTTPClient 创建新的UiPath HTTP客户端
func NewUiPathHTTPClient(config *models.UiPathConfig, logger *logrus.Logger) interfaces.UiPathClient {
	// 创建HTTP传输配置
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL},
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   config.Timeout,
	}

	return &UiPathHTTPClient{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
	}
}

// Authenticate 获取UiPath认证令牌
func (c *UiPathHTTPClient) Authenticate(ctx context.Context) (string, error) {
	authURL := c.config.OrchBaseURL + "/api/account/authenticate"

	authData := models.UiPathAuthRequest{
		TenancyName:            c.config.TenancyName,
		UsernameOrEmailAddress: c.config.Username,
		Password:               c.config.Password,
	}

	jsonData, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth data: %w", err)
	}

	c.logger.Info("Authenticating with UiPath Orchestrator...")

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp models.UiPathAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	if authResp.Result == "" {
		return "", fmt.Errorf("token not returned by the authentication endpoint")
	}

	c.logger.Info("Successfully authenticated with UiPath Orchestrator")
	return authResp.Result, nil
}

// AddQueueItem 添加项目到UiPath队列
func (c *UiPathHTTPClient) AddQueueItem(ctx context.Context, req *models.UiPathAddQueueItemRequest) (*models.UiPathAddQueueItemResponse, error) {
	// 获取认证令牌
	token, err := c.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	// 构建队列URL
	queueURL, err := url.JoinPath(c.config.OrchBaseURL, "/odata/Queues/UiPathODataSvc.AddQueueItem")
	if err != nil {
		return nil, fmt.Errorf("failed to build queue URL: %w", err)
	}

	// 设置优先级默认值
	priority := req.Priority
	if priority == "" {
		priority = "Normal"
	}

	// 设置参考值
	reference := req.Reference
	if reference == "" {
		if doc, ok := req.SpecificContent["Documento"].(string); ok {
			reference = doc
		} else {
			reference = "RefAutomatica"
		}
	}

	// 构建队列项目数据
	itemData := models.UiPathItemData{
		Name:            req.QueueName,
		Priority:        priority,
		SpecificContent: req.SpecificContent,
		Reference:       reference,
	}

	queueRequest := models.UiPathQueueItemRequest{
		ItemData: itemData,
	}

	jsonData, err := json.Marshal(queueRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal queue request: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"queue_name": req.QueueName,
		"priority":   priority,
		"reference":  reference,
	}).Info("Adding item to UiPath queue...")

	httpReq, err := http.NewRequestWithContext(ctx, "POST", queueURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create queue request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-UIPATH-OrganizationUnitId", fmt.Sprintf("%d", c.config.FolderID))

	if c.config.TenancyName != "" {
		httpReq.Header.Set("X-UIPATH-TenantName", c.config.TenancyName)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("queue item request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add queue item with status %d: %s", resp.StatusCode, string(body))
	}

	var queueResp models.UiPathQueueItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&queueResp); err != nil {
		return nil, fmt.Errorf("failed to decode queue response: %w", err)
	}

	response := &models.UiPathAddQueueItemResponse{
		ID:        queueResp.ID,
		QueueName: req.QueueName,
		Priority:  priority,
		Reference: reference,
		Status:    "success",
		Message:   "Item successfully added to queue",
		CreatedAt: time.Now(),
	}

	c.logger.WithFields(logrus.Fields{
		"item_id":    queueResp.ID,
		"queue_name": req.QueueName,
		"status":     "success",
	}).Info("Successfully added item to UiPath queue")

	return response, nil
}

// CheckUiPathStatus 检查UiPath系统状态
func (c *UiPathHTTPClient) CheckUiPathStatus(ctx context.Context) (*models.UiPathStatusResponse, error) {
	// 尝试认证来检查系统状态
	token, err := c.Authenticate(ctx)
	if err != nil {
		return &models.UiPathStatusResponse{
			IsAvailable: false,
			LastCheck:   time.Now(),
			Message:     fmt.Sprintf("Authentication failed: %v", err),
		}, nil
	}

	// 尝试访问一个简单的API端点来验证系统可用性
	statusURL := c.config.OrchBaseURL + "/odata/Users"

	req, err := http.NewRequestWithContext(ctx, "GET", statusURL, nil)
	if err != nil {
		return &models.UiPathStatusResponse{
			IsAvailable: false,
			LastCheck:   time.Now(),
			Message:     fmt.Sprintf("Failed to create status request: %v", err),
		}, nil
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &models.UiPathStatusResponse{
			IsAvailable: false,
			LastCheck:   time.Now(),
			Message:     fmt.Sprintf("Status check request failed: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	isAvailable := resp.StatusCode == http.StatusOK
	message := "System is available"
	if !isAvailable {
		message = fmt.Sprintf("System returned status %d", resp.StatusCode)
	}

	return &models.UiPathStatusResponse{
		IsAvailable: isAvailable,
		LastCheck:   time.Now(),
		Message:     message,
	}, nil
}

// IsUiPathAvailable 检查UiPath系统是否可用
func (c *UiPathHTTPClient) IsUiPathAvailable(ctx context.Context) (bool, error) {
	status, err := c.CheckUiPathStatus(ctx)
	if err != nil {
		return false, err
	}

	return status.IsAvailable, nil
}
