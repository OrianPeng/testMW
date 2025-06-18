package models

import (
	"time"
)

// RequestStatus 表示请求状态
type RequestStatus string

const (
	StatusPending    RequestStatus = "pending"    // 等待处理
	StatusProcessing RequestStatus = "processing" // 正在处理
	StatusCompleted  RequestStatus = "completed"  // 已完成
	StatusFailed     RequestStatus = "failed"     // 处理失败
)

// AgentRequest 来自 AI Agent 的请求
type AgentRequest struct {
	ID        string                 `json:"id" binding:"required"`
	AgentID   string                 `json:"agent_id" binding:"required"`
	Data      map[string]interface{} `json:"data" binding:"required"`
	Priority  int                    `json:"priority,omitempty"` // 优先级，数值越大优先级越高
	Callback  string                 `json:"callback,omitempty"` // 回调 URL
	CreatedAt time.Time              `json:"created_at"`
}

// RPARequest 发送给 RPA 系统的请求
type RPARequest struct {
	ID       string                 `json:"id"`
	AgentID  string                 `json:"agent_id"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// QueuedRequest 队列中的请求
type QueuedRequest struct {
	*AgentRequest
	Status     RequestStatus `json:"status"`
	RetryCount int           `json:"retry_count"`
	UpdatedAt  time.Time     `json:"updated_at"`
	ErrorMsg   string        `json:"error_msg,omitempty"`
}

// RPAResponse RPA 系统的响应
type RPAResponse struct {
	ID      string                 `json:"id"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// RPAStatus RPA 系统状态
type RPAStatus struct {
	IsAvailable bool      `json:"is_available"`
	LastCheck   time.Time `json:"last_check"`
	CurrentTask string    `json:"current_task,omitempty"`
}

// AgentResponse 返回给 AI Agent 的响应
type AgentResponse struct {
	RequestID string        `json:"request_id"`
	Status    RequestStatus `json:"status"`
	Message   string        `json:"message"`
	Data      interface{}   `json:"data,omitempty"`
}