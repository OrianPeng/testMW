package interfaces

import (
	"context"
	"rpa-middleware/internal/models"
)

// QueueManager 队列管理器接口
type QueueManager interface {
	// EnqueuePurchaseRequest 将采购请求加入队列
	EnqueuePurchaseRequest(ctx context.Context, req *models.PurchaseRequest) error

	// DequeuePurchaseRequest 从队列中取出采购请求（按优先级）
	DequeuePurchaseRequest(ctx context.Context) (*models.QueuedPurchaseRequest, error)

	// UpdatePurchaseRequestStatus 更新采购请求状态
	UpdatePurchaseRequestStatus(ctx context.Context, requestID string, status models.PurchaseRequestStatus, errorMsg string) error

	// GetPurchaseRequestStatus 获取采购请求状态
	GetPurchaseRequestStatus(ctx context.Context, requestID string) (*models.QueuedPurchaseRequest, error)

	// GetPendingCount 获取待处理请求数量
	GetPendingCount(ctx context.Context) (int, error)

	// ListPurchaseRequests 列出采购请求（支持分页和过滤）
	ListPurchaseRequests(ctx context.Context, status string, limit, offset int) ([]*models.QueuedPurchaseRequest, error)

	// DeletePurchaseRequest 删除队列中的采购请求
	DeletePurchaseRequest(ctx context.Context, requestID string) error

	// ClearAll 清空所有采购请求
	ClearAll(ctx context.Context) error
}

// RPAClient RPA 系统客户端接口
type RPAClient interface {
	// SendPurchaseRequest 发送采购请求到 RPA 系统
	SendPurchaseRequest(ctx context.Context, req *models.RPARequest) (*models.RPAResponse, error)

	// CheckStatus 检查 RPA 系统状态
	CheckStatus(ctx context.Context) (*models.RPAStatus, error)

	// IsAvailable 检查 RPA 系统是否可用
	IsAvailable(ctx context.Context) (bool, error)
}

// NotificationService 通知服务接口
type NotificationService interface {
	// SendPurchaseRequestCallback 发送采购请求回调通知
	SendPurchaseRequestCallback(ctx context.Context, callbackURL string, response *models.PurchaseRequestCallbackResponse) error

	// NotifyPurchaseRequestCompletion 通知采购请求完成
	NotifyPurchaseRequestCompletion(ctx context.Context, req *models.QueuedPurchaseRequest, result *models.RPAResponse) error
}
