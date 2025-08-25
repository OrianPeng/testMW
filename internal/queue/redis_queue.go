package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rpa-middleware/internal/interfaces"
	"rpa-middleware/internal/models"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisQueueManager Redis队列管理器
type RedisQueueManager struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewRedisQueueManager 创建新的Redis队列管理器
func NewRedisQueueManager(addr, password string, db int, logger *logrus.Logger) (interfaces.QueueManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis queue manager initialized successfully")

	return &RedisQueueManager{
		client: client,
		logger: logger,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisQueueManager) Close() error {
	return r.client.Close()
}

// EnqueuePurchaseRequest 将采购请求加入队列
func (r *RedisQueueManager) EnqueuePurchaseRequest(ctx context.Context, req *models.PurchaseRequest) error {
	// 检查请求是否已存在
	exists, err := r.client.Exists(ctx, r.getRequestKey(req.RequestID)).Result()
	if err != nil {
		return fmt.Errorf("failed to check request existence: %w", err)
	}
	if exists == 1 {
		return errors.New("request already exists in queue")
	}

	// 创建队列请求
	queuedReq := &models.QueuedPurchaseRequest{
		PurchaseRequest: req,
		QueuePosition:   0, // 将在入队时计算
		EnqueuedAt:      time.Now(),
	}

	// 设置状态为处理中
	queuedReq.Status = models.PurchaseStatusProcessing

	// 序列化请求
	reqData, err := json.Marshal(queuedReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 使用管道批量操作
	pipe := r.client.Pipeline()

	// 存储请求数据
	pipe.Set(ctx, r.getRequestKey(req.RequestID), reqData, 0)

	// 添加到优先级队列（使用优先级作为分数）
	score := float64(req.Priority)
	pipe.ZAdd(ctx, r.getQueueKey(), &redis.Z{
		Score:  score,
		Member: req.RequestID,
	})

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to enqueue request: %w", err)
	}

	// 更新队列位置
	r.updateQueuePositions(ctx)

	r.logger.WithFields(logrus.Fields{
		"request_id": req.RequestID,
		"priority":   req.Priority,
	}).Info("Purchase request enqueued to Redis")

	return nil
}

// DequeuePurchaseRequest 从队列中取出采购请求（按优先级）
func (r *RedisQueueManager) DequeuePurchaseRequest(ctx context.Context) (*models.QueuedPurchaseRequest, error) {
	// 获取最高优先级的请求
	result, err := r.client.ZRevRangeWithScores(ctx, r.getQueueKey(), 0, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get request from queue: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("queue is empty")
	}

	requestID := result[0].Member.(string)

	// 获取请求数据
	reqData, err := r.client.Get(ctx, r.getRequestKey(requestID)).Result()
	if err != nil {
		if err == redis.Nil {
			// 请求数据不存在，从队列中移除
			r.client.ZRem(ctx, r.getQueueKey(), requestID)
			return nil, errors.New("request data not found")
		}
		return nil, fmt.Errorf("failed to get request data: %w", err)
	}

	// 反序列化请求
	var queuedReq models.QueuedPurchaseRequest
	if err := json.Unmarshal([]byte(reqData), &queuedReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 检查状态
	if queuedReq.Status != models.PurchaseStatusProcessing {
		// 从队列中移除非处理状态的请求
		r.client.ZRem(ctx, r.getQueueKey(), requestID)
		return nil, errors.New("request is not in processing status")
	}

	// 使用管道批量操作
	pipe := r.client.Pipeline()

	// 从队列中移除
	pipe.ZRem(ctx, r.getQueueKey(), requestID)

	// 删除请求数据
	pipe.Del(ctx, r.getRequestKey(requestID))

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue request: %w", err)
	}

	// 更新队列位置
	r.updateQueuePositions(ctx)

	r.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"priority":   queuedReq.Priority,
	}).Info("Purchase request dequeued from Redis")

	return &queuedReq, nil
}

// UpdatePurchaseRequestStatus 更新采购请求状态
func (r *RedisQueueManager) UpdatePurchaseRequestStatus(ctx context.Context, requestID string, status models.PurchaseRequestStatus, errorMsg string) error {
	// 获取请求数据
	reqData, err := r.client.Get(ctx, r.getRequestKey(requestID)).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("request not found in queue")
		}
		return fmt.Errorf("failed to get request data: %w", err)
	}

	// 反序列化请求
	var queuedReq models.QueuedPurchaseRequest
	if err := json.Unmarshal([]byte(reqData), &queuedReq); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 更新状态
	queuedReq.Status = status
	queuedReq.ErrorMsg = errorMsg
	queuedReq.UpdatedAt = time.Now()

	if status == models.PurchaseStatusCompleted || status == models.PurchaseStatusFailed {
		queuedReq.ProcessedAt = &time.Time{}
		*queuedReq.ProcessedAt = time.Now()
	}

	// 序列化更新后的请求
	updatedData, err := json.Marshal(queuedReq)
	if err != nil {
		return fmt.Errorf("failed to marshal updated request: %w", err)
	}

	// 更新存储
	if err := r.client.Set(ctx, r.getRequestKey(requestID), updatedData, 0).Err(); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	// 如果状态不是处理中，从队列中移除
	if status != models.PurchaseStatusProcessing {
		r.client.ZRem(ctx, r.getQueueKey(), requestID)
		r.updateQueuePositions(ctx)
	}

	r.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"status":     status,
		"error_msg":  errorMsg,
	}).Info("Purchase request status updated in Redis")

	return nil
}

// GetPurchaseRequestStatus 获取采购请求状态
func (r *RedisQueueManager) GetPurchaseRequestStatus(ctx context.Context, requestID string) (*models.QueuedPurchaseRequest, error) {
	reqData, err := r.client.Get(ctx, r.getRequestKey(requestID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("request not found in queue")
		}
		return nil, fmt.Errorf("failed to get request data: %w", err)
	}

	var queuedReq models.QueuedPurchaseRequest
	if err := json.Unmarshal([]byte(reqData), &queuedReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return &queuedReq, nil
}

// GetPendingCount 获取待处理请求数量
func (r *RedisQueueManager) GetPendingCount(ctx context.Context) (int, error) {
	count, err := r.client.ZCard(ctx, r.getQueueKey()).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue count: %w", err)
	}

	return int(count), nil
}

// ListPurchaseRequests 列出采购请求（支持分页和过滤）
func (r *RedisQueueManager) ListPurchaseRequests(ctx context.Context, status string, limit, offset int) ([]*models.QueuedPurchaseRequest, error) {
	// 获取队列中的所有请求ID
	requestIDs, err := r.client.ZRevRange(ctx, r.getQueueKey(), int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get request IDs: %w", err)
	}

	var requests []*models.QueuedPurchaseRequest

	// 批量获取请求数据
	if len(requestIDs) > 0 {
		keys := make([]string, len(requestIDs))
		for i, id := range requestIDs {
			keys[i] = r.getRequestKey(id)
		}

		reqDataList, err := r.client.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get request data: %w", err)
		}

		for _, reqData := range reqDataList {
			if reqData == nil {
				continue
			}

			var queuedReq models.QueuedPurchaseRequest
			if err := json.Unmarshal([]byte(reqData.(string)), &queuedReq); err != nil {
				r.logger.WithError(err).Warn("Failed to unmarshal request data")
				continue
			}

			// 过滤状态
			if status == "" || queuedReq.Status == models.PurchaseRequestStatus(status) {
				requests = append(requests, &queuedReq)
			}
		}
	}

	// 按优先级和入队时间排序
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].Priority != requests[j].Priority {
			return requests[i].Priority > requests[j].Priority
		}
		return requests[i].EnqueuedAt.Before(requests[j].EnqueuedAt)
	})

	return requests, nil
}

// DeletePurchaseRequest 删除队列中的采购请求
func (r *RedisQueueManager) DeletePurchaseRequest(ctx context.Context, requestID string) error {
	pipe := r.client.Pipeline()

	// 从队列中移除
	pipe.ZRem(ctx, r.getQueueKey(), requestID)

	// 删除请求数据
	pipe.Del(ctx, r.getRequestKey(requestID))

	// 执行管道操作
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	// 更新队列位置
	r.updateQueuePositions(ctx)

	r.logger.WithField("request_id", requestID).Info("Purchase request deleted from Redis queue")

	return nil
}

// ClearAll 清空所有采购请求
func (r *RedisQueueManager) ClearAll(ctx context.Context) error {
	// 获取所有请求ID
	requestIDs, err := r.client.ZRange(ctx, r.getQueueKey(), 0, -1).Result()
	if err != nil {
		return err
	}
	if len(requestIDs) > 0 {
		// 删除所有请求数据
		keys := make([]string, len(requestIDs))
		for i, id := range requestIDs {
			keys[i] = r.getRequestKey(id)
		}
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	// 清空队列索引
	if err := r.client.Del(ctx, r.getQueueKey()).Err(); err != nil {
		return err
	}
	return nil
}

// getRequestKey 获取请求数据的Redis键
func (r *RedisQueueManager) getRequestKey(requestID string) string {
	return fmt.Sprintf("purchase_request:%s", requestID)
}

// getQueueKey 获取队列的Redis键
func (r *RedisQueueManager) getQueueKey() string {
	return "purchase_request_queue"
}

// updateQueuePositions 更新队列位置
func (r *RedisQueueManager) updateQueuePositions(ctx context.Context) {
	// 获取队列中的所有请求（按优先级排序）
	requestIDs, err := r.client.ZRevRange(ctx, r.getQueueKey(), 0, -1).Result()
	if err != nil {
		r.logger.WithError(err).Error("Failed to get queue positions")
		return
	}

	// 批量更新位置
	pipe := r.client.Pipeline()
	for i, requestID := range requestIDs {
		// 获取请求数据
		reqData, err := r.client.Get(ctx, r.getRequestKey(requestID)).Result()
		if err != nil {
			continue
		}

		var queuedReq models.QueuedPurchaseRequest
		if err := json.Unmarshal([]byte(reqData), &queuedReq); err != nil {
			continue
		}

		// 更新位置
		queuedReq.QueuePosition = i + 1

		// 序列化并更新
		updatedData, err := json.Marshal(queuedReq)
		if err != nil {
			continue
		}

		pipe.Set(ctx, r.getRequestKey(requestID), updatedData, 0)
	}

	pipe.Exec(ctx)
}
