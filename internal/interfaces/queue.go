package interfaces

import (
	"context"
	"rpa-middleware/internal/models"
)

// QueueManager 队列管理器接口
type QueueManager interface {
	// EnqueueRequest 将请求加入队列
	EnqueueRequest(ctx context.Context, req *models.AgentRequest) error
	
	// DequeueRequest 从队列中取出请求（按优先级）
	DequeueRequest(ctx context.Context) (*models.QueuedRequest, error)
	
	// UpdateRequestStatus 更新请求状态
	UpdateRequestStatus(ctx context.Context, requestID string, status models.RequestStatus, errorMsg string) error
	
	// GetRequestStatus 获取请求状态
	GetRequestStatus(ctx context.Context, requestID string) (*models.QueuedRequest, error)
	
	// GetPendingCount 获取待处理请求数量
	GetPendingCount(ctx context.Context) (int, error)
	
	// ListRequests 列出请求（支持分页和过滤）
	ListRequests(ctx context.Context, status models.RequestStatus, limit, offset int) ([]*models.QueuedRequest, error)
}

// RPAClient RPA 系统客户端接口
type RPAClient interface {
	// SendRequest 发送请求到 RPA 系统
	SendRequest(ctx context.Context, req *models.RPARequest) (*models.RPAResponse, error)
	
	// CheckStatus 检查 RPA 系统状态
	CheckStatus(ctx context.Context) (*models.RPAStatus, error)
	
	// IsAvailable 检查 RPA 系统是否可用
	IsAvailable(ctx context.Context) (bool, error)
}

// NotificationService 通知服务接口
type NotificationService interface {
	// SendCallback 发送回调通知
	SendCallback(ctx context.Context, callbackURL string, response *models.AgentResponse) error
	
	// NotifyCompletion 通知请求完成
	NotifyCompletion(ctx context.Context, req *models.QueuedRequest, result *models.RPAResponse) error
}