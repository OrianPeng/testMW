package queue

import (
	"context"
	"errors"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// MemoryQueueManager 内存队列管理器
type MemoryQueueManager struct {
	requests map[string]*models.QueuedPurchaseRequest
	mutex    sync.RWMutex
	logger   *logrus.Logger
}

// NewMemoryQueueManager 创建新的内存队列管理器
func NewMemoryQueueManager() interfaces.QueueManager {
	return &MemoryQueueManager{
		requests: make(map[string]*models.QueuedPurchaseRequest),
		logger:   logrus.New(),
	}
}

// EnqueuePurchaseRequest 将采购请求加入队列
func (m *MemoryQueueManager) EnqueuePurchaseRequest(ctx context.Context, req *models.PurchaseRequest) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查请求是否已存在
	if _, exists := m.requests[req.RequestID]; exists {
		return errors.New("request already exists in queue")
	}

	// 创建队列请求
	queuedReq := &models.QueuedPurchaseRequest{
		PurchaseRequest: req,
		QueuePosition:   len(m.requests) + 1,
		EnqueuedAt:      time.Now(),
	}

	// 设置状态为处理中
	queuedReq.Status = models.PurchaseStatusProcessing

	m.requests[req.RequestID] = queuedReq

	m.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"priority":   req.Priority,
		"position":   queuedReq.QueuePosition,
	}).Info("Purchase request enqueued")

	return nil
}

// DequeuePurchaseRequest 从队列中取出采购请求（按优先级）
func (m *MemoryQueueManager) DequeuePurchaseRequest(ctx context.Context) (*models.QueuedPurchaseRequest, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.requests) == 0 {
		return nil, errors.New("queue is empty")
	}

	// 按优先级排序
	var requests []*models.QueuedPurchaseRequest
	for _, req := range m.requests {
		if req.Status == models.PurchaseStatusProcessing {
			requests = append(requests, req)
		}
	}

	if len(requests) == 0 {
		return nil, errors.New("no processing requests available")
	}

	// 按优先级排序（优先级高的先处理）
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].Priority != requests[j].Priority {
			return requests[i].Priority > requests[j].Priority
		}
		// 优先级相同时，按入队时间排序
		return requests[i].EnqueuedAt.Before(requests[j].EnqueuedAt)
	})

	// 取出第一个请求
	req := requests[0]
	delete(m.requests, req.RequestID)

	// 重新计算队列位置
	m.recalculatePositions()

	m.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"priority":   req.Priority,
	}).Info("Purchase request dequeued")

	return req, nil
}

// UpdatePurchaseRequestStatus 更新采购请求状态
func (m *MemoryQueueManager) UpdatePurchaseRequestStatus(ctx context.Context, requestID string, status models.PurchaseRequestStatus, errorMsg string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	req, exists := m.requests[requestID]
	if !exists {
		return errors.New("request not found in queue")
	}

	req.Status = status
	req.ErrorMsg = errorMsg
	req.UpdatedAt = time.Now()

	if status == models.PurchaseStatusCompleted || status == models.PurchaseStatusFailed {
		req.ProcessedAt = &time.Time{}
		*req.ProcessedAt = time.Now()
	}

	m.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"status":     status,
		"error_msg":  errorMsg,
	}).Info("Purchase request status updated")

	return nil
}

// GetPurchaseRequestStatus 获取采购请求状态
func (m *MemoryQueueManager) GetPurchaseRequestStatus(ctx context.Context, requestID string) (*models.QueuedPurchaseRequest, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	req, exists := m.requests[requestID]
	if !exists {
		return nil, errors.New("request not found in queue")
	}

	return req, nil
}

// GetPendingCount 获取待处理请求数量
func (m *MemoryQueueManager) GetPendingCount(ctx context.Context) (int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	count := 0
	for _, req := range m.requests {
		if req.Status == models.PurchaseStatusProcessing {
			count++
		}
	}

	return count, nil
}

// ListPurchaseRequests 列出采购请求（支持分页和过滤）
func (m *MemoryQueueManager) ListPurchaseRequests(ctx context.Context, status string, limit, offset int) ([]*models.QueuedPurchaseRequest, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var requests []*models.QueuedPurchaseRequest
	for _, req := range m.requests {
		// 如果状态为空字符串，返回所有请求；否则只返回匹配状态的请求
		if status == "" || req.Status == models.PurchaseRequestStatus(status) {
			requests = append(requests, req)
		}
	}

	// 按优先级和入队时间排序
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].Priority != requests[j].Priority {
			return requests[i].Priority > requests[j].Priority
		}
		return requests[i].EnqueuedAt.Before(requests[j].EnqueuedAt)
	})

	// 分页
	start := offset
	end := start + limit
	if start >= len(requests) {
		return []*models.QueuedPurchaseRequest{}, nil
	}
	if end > len(requests) {
		end = len(requests)
	}

	return requests[start:end], nil
}

// DeletePurchaseRequest 删除队列中的采购请求
func (m *MemoryQueueManager) DeletePurchaseRequest(ctx context.Context, requestID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.requests[requestID]; !exists {
		return errors.New("request not found in queue")
	}

	delete(m.requests, requestID)

	// 重新计算队列位置
	m.recalculatePositions()

	m.logger.WithField("request_id", requestID).Info("Purchase request deleted from queue")

	return nil
}

// ClearAll 清空所有采购请求
func (m *MemoryQueueManager) ClearAll(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.requests = make(map[string]*models.QueuedPurchaseRequest)
	return nil
}

// recalculatePositions 重新计算队列位置
func (m *MemoryQueueManager) recalculatePositions() {
	position := 1
	for _, req := range m.requests {
		req.QueuePosition = position
		position++
	}
}
