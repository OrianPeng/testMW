package queue

import (
	"context"
	"fmt"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"sort"
	"sync"
	"time"
)

// MemoryQueueManager 基于内存的队列管理器
type MemoryQueueManager struct {
	mu       sync.RWMutex
	requests map[string]*models.QueuedRequest
	queue    []*models.QueuedRequest
}

// NewMemoryQueueManager 创建新的内存队列管理器
func NewMemoryQueueManager() interfaces.QueueManager {
	return &MemoryQueueManager{
		requests: make(map[string]*models.QueuedRequest),
		queue:    make([]*models.QueuedRequest, 0),
	}
}

// EnqueueRequest 将请求加入队列
func (m *MemoryQueueManager) EnqueueRequest(ctx context.Context, req *models.AgentRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	queuedReq := &models.QueuedRequest{
		AgentRequest: req,
		Status:       models.StatusPending,
		RetryCount:   0,
		UpdatedAt:    time.Now(),
	}

	m.requests[req.ID] = queuedReq
	m.queue = append(m.queue, queuedReq)
	
	// 按优先级排序（高优先级在前）
	sort.Slice(m.queue, func(i, j int) bool {
		if m.queue[i].Priority == m.queue[j].Priority {
			return m.queue[i].CreatedAt.Before(m.queue[j].CreatedAt)
		}
		return m.queue[i].Priority > m.queue[j].Priority
	})

	return nil
}

// DequeueRequest 从队列中取出请求
func (m *MemoryQueueManager) DequeueRequest(ctx context.Context) (*models.QueuedRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 找到第一个待处理的请求
	for i, req := range m.queue {
		if req.Status == models.StatusPending {
			req.Status = models.StatusProcessing
			req.UpdatedAt = time.Now()
			
			// 从队列中移除
			m.queue = append(m.queue[:i], m.queue[i+1:]...)
			return req, nil
		}
	}

	return nil, fmt.Errorf("no pending requests in queue")
}

// UpdateRequestStatus 更新请求状态
func (m *MemoryQueueManager) UpdateRequestStatus(ctx context.Context, requestID string, status models.RequestStatus, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	req, exists := m.requests[requestID]
	if !exists {
		return fmt.Errorf("request not found: %s", requestID)
	}

	req.Status = status
	req.UpdatedAt = time.Now()
	req.ErrorMsg = errorMsg

	if status == models.StatusFailed {
		req.RetryCount++
	}

	return nil
}

// GetRequestStatus 获取请求状态
func (m *MemoryQueueManager) GetRequestStatus(ctx context.Context, requestID string) (*models.QueuedRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	req, exists := m.requests[requestID]
	if !exists {
		return nil, fmt.Errorf("request not found: %s", requestID)
	}

	// 返回副本以避免并发修改
	result := *req
	return &result, nil
}

// GetPendingCount 获取待处理请求数量
func (m *MemoryQueueManager) GetPendingCount(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, req := range m.queue {
		if req.Status == models.StatusPending {
			count++
		}
	}

	return count, nil
}

// ListRequests 列出请求
func (m *MemoryQueueManager) ListRequests(ctx context.Context, status models.RequestStatus, limit, offset int) ([]*models.QueuedRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []*models.QueuedRequest
	for _, req := range m.requests {
		if status == "" || req.Status == status {
			filtered = append(filtered, req)
		}
	}

	// 按创建时间排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	// 分页
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := make([]*models.QueuedRequest, end-start)
	for i := start; i < end; i++ {
		// 返回副本
		req := *filtered[i]
		result[i-start] = &req
	}

	return result, nil
}