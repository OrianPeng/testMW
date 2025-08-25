package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PurchaseRequestStatus 采购请求状态
type PurchaseRequestStatus string

const (
	PurchaseStatusPending    PurchaseRequestStatus = "pending"    // 待审核
	PurchaseStatusApproved   PurchaseRequestStatus = "approved"   // 已批准
	PurchaseStatusRejected   PurchaseRequestStatus = "rejected"   // 已拒绝
	PurchaseStatusCompleted  PurchaseRequestStatus = "completed"  // 已完成
	PurchaseStatusCancelled  PurchaseRequestStatus = "cancelled"  // 已取消
	PurchaseStatusProcessing PurchaseRequestStatus = "processing" // 正在处理（队列中）
	PurchaseStatusFailed     PurchaseRequestStatus = "failed"     // 处理失败
	PurchaseStatusEnough     PurchaseRequestStatus = "enough"     // 信息完整
	PurchaseStatusLack       PurchaseRequestStatus = "lack"       // 信息不完整
)

// PurchaseRequest 采购请求
type PurchaseRequest struct {
	ID                   int64                 `json:"id" db:"id"`
	RequestID            string                `json:"request_id" db:"request_id" binding:"required"`
	DocType              string                `json:"doc_type" db:"doc_type" binding:"required"`
	Plant                string                `json:"plant" db:"plant" binding:"required"`
	Quantity             int                   `json:"quantity" db:"quantity" binding:"required,min=1"`
	UnitPrice            float64               `json:"unit_price" db:"unit_price" binding:"required,min=0"`
	Material             string                `json:"material" db:"material" binding:"required"`
	DeliveryDate         time.Time             `json:"delivery_date" db:"delivery_date" binding:"required"`
	VendorCode           string                `json:"vendor_code" db:"vendor_code" binding:"required"`
	ShortText            string                `json:"short_text" db:"short_text" binding:"required"`
	MaterialGroup        string                `json:"material_group" db:"material_group" binding:"required"`
	UnitType             string                `json:"unit_type" db:"unit_type" binding:"required"`
	Requester            string                `json:"requester" db:"requester" binding:"required"`
	PurchaseOrganization string                `json:"purchase_organization" db:"purchase_organization" binding:"required"`
	Currency             string                `json:"currency" db:"currency"`
	TotalAmount          float64               `json:"total_amount" db:"total_amount"`
	Priority             int                   `json:"priority" db:"priority"` // 优先级：1-低，2-中，3-高，4-紧急
	Status               PurchaseRequestStatus `json:"status" db:"status"`
	Urgency              string                `json:"urgency" db:"urgency"` // 紧急程度：normal, urgent, critical
	ApproverID           *string               `json:"approver_id" db:"approver_id"`
	ApproverName         *string               `json:"approver_name" db:"approver_name"`
	ApprovedAt           *time.Time            `json:"approved_at" db:"approved_at"`
	RejectionReason      *string               `json:"rejection_reason" db:"rejection_reason"`
	Comments             string                `json:"comments" db:"comments"`
	RetryCount           int                   `json:"retry_count" db:"retry_count"`             // 重试次数
	ErrorMsg             string                `json:"error_msg,omitempty" db:"error_msg"`       // 错误信息
	ProcessedAt          *time.Time            `json:"processed_at,omitempty" db:"processed_at"` // 处理完成时间

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePurchaseRequestRequest 创建采购请求的请求
type CreatePurchaseRequestRequest struct {
	RequestID            string     `json:"request_id"`
	DocType              string     `json:"doc_type"`
	Plant                string     `json:"plant"`
	Quantity             int        `json:"quantity"`
	UnitPrice            float64    `json:"unit_price"`
	Material             string     `json:"material"`
	DeliveryDate         *time.Time `json:"delivery_date"`
	VendorCode           string     `json:"vendor_code"`
	ShortText            string     `json:"short_text"`
	MaterialGroup        string     `json:"material_group"`
	UnitType             string     `json:"unit_type"`
	Requester            string     `json:"requester"`
	PurchaseOrganization string     `json:"purchase_organization"`
	Currency             string     `json:"currency"`
	Priority             int        `json:"priority"` // 优先级：1-低，2-中，3-高，4-紧急
	Urgency              string     `json:"urgency"`  // 紧急程度：normal, urgent, critical
	Comments             string     `json:"comments"`
}

// 自定义JSON反序列化，处理delivery_date为空字符串
func (r *CreatePurchaseRequestRequest) UnmarshalJSON(data []byte) error {
	type Alias CreatePurchaseRequestRequest
	tmp := &struct {
		DeliveryDate interface{} `json:"delivery_date"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	// 处理delivery_date
	switch v := tmp.DeliveryDate.(type) {
	case string:
		if v == "" {
			r.DeliveryDate = nil
			break
		}
		t, err := parseTimeFlexible(v)
		if err != nil {
			return fmt.Errorf("failed to parse delivery_date: %w", err)
		}
		r.DeliveryDate = &t
	case nil:
		r.DeliveryDate = nil
	}
	return nil
}

// ValidateAndSetDefaults 验证请求并设置默认值
func (req *CreatePurchaseRequestRequest) ValidateAndSetDefaults() error {
	// 设置默认值为不合理且一致的特殊值
	if req.DocType == "" {
		req.DocType = "INVALID"
	}
	if req.Plant == "" {
		req.Plant = "INVALID"
	}
	if req.Quantity == 0 {
		req.Quantity = -1
	}
	if req.UnitPrice == 0 {
		req.UnitPrice = -1
	}
	if req.Material == "" {
		req.Material = "INVALID"
	}
	if req.DeliveryDate == nil {
		// 设置为1970-01-01 00:00:00 UTC，表示无效
		invalidDate := time.Unix(0, 0).UTC()
		req.DeliveryDate = &invalidDate
	}
	if req.VendorCode == "" {
		req.VendorCode = "INVALID"
	}
	if req.ShortText == "" {
		req.ShortText = "INVALID"
	}
	if req.MaterialGroup == "" {
		req.MaterialGroup = "INVALID"
	}
	if req.UnitType == "" {
		req.UnitType = "INVALID"
	}
	if req.Requester == "" {
		req.Requester = "INVALID"
	}
	if req.PurchaseOrganization == "" {
		req.PurchaseOrganization = "INVALID"
	}
	if req.Currency == "" {
		req.Currency = "INVALID"
	}
	if req.Priority == 0 {
		req.Priority = -1
	}
	if req.Urgency == "" {
		req.Urgency = "INVALID"
	}
	return nil
}

// UpdatePurchaseRequestRequest 更新采购请求的请求
type UpdatePurchaseRequestRequest struct {
	// 基本信息
	DocType              *string    `json:"doc_type"`
	Plant                *string    `json:"plant"`
	Quantity             *int       `json:"quantity"`
	UnitPrice            *float64   `json:"unit_price"`
	Material             *string    `json:"material"`
	DeliveryDate         *time.Time `json:"delivery_date"`
	VendorCode           *string    `json:"vendor_code"`
	ShortText            *string    `json:"short_text"`
	MaterialGroup        *string    `json:"material_group"`
	UnitType             *string    `json:"unit_type"`
	Requester            *string    `json:"requester"`
	PurchaseOrganization *string    `json:"purchase_organization"`
	Currency             *string    `json:"currency"`
	Priority             *int       `json:"priority"`
	Urgency              *string    `json:"urgency"`
	Comments             *string    `json:"comments"`

	// 状态和审批相关
	Status          *PurchaseRequestStatus `json:"status"`
	ApproverID      *string                `json:"approver_id"`
	ApproverName    *string                `json:"approver_name"`
	RejectionReason *string                `json:"rejection_reason"`

	// 处理相关
	RetryCount  *int       `json:"retry_count"`
	ErrorMsg    *string    `json:"error_msg"`
	ProcessedAt *time.Time `json:"processed_at"`
}

// 自定义JSON反序列化，处理UpdatePurchaseRequestRequest中的时间字段
func (r *UpdatePurchaseRequestRequest) UnmarshalJSON(data []byte) error {
	type Alias UpdatePurchaseRequestRequest
	tmp := &struct {
		DeliveryDate interface{} `json:"delivery_date"`
		ProcessedAt  interface{} `json:"processed_at"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理delivery_date
	switch v := tmp.DeliveryDate.(type) {
	case string:
		if v == "" {
			r.DeliveryDate = nil
		} else {
			t, err := parseTimeFlexible(v)
			if err != nil {
				return fmt.Errorf("failed to parse delivery_date: %w", err)
			}
			r.DeliveryDate = &t
		}
	case nil:
		r.DeliveryDate = nil
	}

	// 处理processed_at
	switch v := tmp.ProcessedAt.(type) {
	case string:
		if v == "" {
			r.ProcessedAt = nil
		} else {
			t, err := parseTimeFlexible(v)
			if err != nil {
				return fmt.Errorf("failed to parse processed_at: %w", err)
			}
			r.ProcessedAt = &t
		}
	case nil:
		r.ProcessedAt = nil
	}

	return nil
}

// PurchaseRequestResponse 采购请求响应
type PurchaseRequestResponse struct {
	ID                   int64                 `json:"id"`
	RequestID            string                `json:"request_id"`
	DocType              string                `json:"doc_type"`
	Plant                string                `json:"plant"`
	Quantity             int                   `json:"quantity"`
	UnitPrice            float64               `json:"unit_price"`
	Material             string                `json:"material"`
	DeliveryDate         time.Time             `json:"delivery_date"`
	VendorCode           string                `json:"vendor_code"`
	ShortText            string                `json:"short_text"`
	MaterialGroup        string                `json:"material_group"`
	UnitType             string                `json:"unit_type"`
	Requester            string                `json:"requester"`
	PurchaseOrganization string                `json:"purchase_organization"`
	Currency             string                `json:"currency"`
	TotalAmount          float64               `json:"total_amount"`
	Priority             int                   `json:"priority"`
	Status               PurchaseRequestStatus `json:"status"`
	Urgency              string                `json:"urgency"`
	ApproverID           *string               `json:"approver_id"`
	ApproverName         *string               `json:"approver_name"`
	ApprovedAt           *time.Time            `json:"approved_at"`
	RejectionReason      *string               `json:"rejection_reason"`
	Comments             string                `json:"comments"`
	RetryCount           int                   `json:"retry_count"`
	ErrorMsg             string                `json:"error_msg,omitempty"`
	ProcessedAt          *time.Time            `json:"processed_at,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// PurchaseRequestListResponse 采购请求列表响应
type PurchaseRequestListResponse struct {
	Requests []*PurchaseRequestResponse `json:"requests"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
}

// PurchaseRequestQuery 采购请求查询条件
type PurchaseRequestQuery struct {
	Status        *PurchaseRequestStatus `json:"status"`
	Requester     *string                `json:"requester"`
	Plant         *string                `json:"plant"`
	VendorCode    *string                `json:"vendor_code"`
	MaterialGroup *string                `json:"material_group"`
	Priority      *int                   `json:"priority"`
	Urgency       *string                `json:"urgency"`
	StartDate     *time.Time             `json:"start_date"`
	EndDate       *time.Time             `json:"end_date"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	SortBy        string                 `json:"sort_by"`
	SortOrder     string                 `json:"sort_order"` // asc, desc
}

// QueuedPurchaseRequest 队列中的采购请求
type QueuedPurchaseRequest struct {
	*PurchaseRequest
	QueuePosition int       `json:"queue_position"` // 队列位置
	EnqueuedAt    time.Time `json:"enqueued_at"`    // 入队时间
}

// RPARequest 发送给 RPA 系统的请求
type RPARequest struct {
	ID        string                 `json:"id"`
	RequestID string                 `json:"request_id"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}

// RPAResponse RPA 系统的响应
type RPAResponse struct {
	ID        string                 `json:"id"`
	RequestID string                 `json:"request_id"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// RPAStatus RPA 系统状态
type RPAStatus struct {
	IsAvailable bool      `json:"is_available"`
	LastCheck   time.Time `json:"last_check"`
	CurrentTask string    `json:"current_task,omitempty"`
}

// PurchaseRequestCallbackResponse 采购请求回调响应
type PurchaseRequestCallbackResponse struct {
	RequestID string                `json:"request_id"`
	Status    PurchaseRequestStatus `json:"status"`
	Message   string                `json:"message"`
	Data      interface{}           `json:"data,omitempty"`
}

// CheckFieldsCompleteness 检查字段完整性并返回相应的状态
func CheckFieldsCompleteness(req *PurchaseRequest) PurchaseRequestStatus {
	// 检查必需字段是否都不为空且有效
	if req.RequestID != "" && req.RequestID != "INVALID" &&
		req.DocType != "" && req.DocType != "INVALID" &&
		req.Plant != "" && req.Plant != "INVALID" &&
		req.Quantity > 0 && req.Quantity != -1 &&
		req.UnitPrice > 0 && req.UnitPrice != -1 &&
		req.Material != "" && req.Material != "INVALID" &&
		!req.DeliveryDate.IsZero() && req.DeliveryDate.Unix() != 0 &&
		req.VendorCode != "" && req.VendorCode != "INVALID" &&
		req.ShortText != "" && req.ShortText != "INVALID" &&
		req.MaterialGroup != "" && req.MaterialGroup != "INVALID" &&
		req.UnitType != "" && req.UnitType != "INVALID" &&
		req.Requester != "" && req.Requester != "INVALID" &&
		req.PurchaseOrganization != "" && req.PurchaseOrganization != "INVALID" &&
		req.Currency != "" && req.Currency != "INVALID" &&
		req.TotalAmount > 0 && req.TotalAmount != -1 {
		return PurchaseStatusEnough
	}
	return PurchaseStatusLack
}

// CheckCreateRequestCompleteness 检查创建请求的字段完整性
func CheckCreateRequestCompleteness(req *CreatePurchaseRequestRequest) PurchaseRequestStatus {
	// 计算总金额
	totalAmount := float64(req.Quantity) * req.UnitPrice

	// 检查必需字段是否都不为空且有效
	if req.RequestID != "" && req.RequestID != "INVALID" &&
		req.DocType != "" && req.DocType != "INVALID" &&
		req.Plant != "" && req.Plant != "INVALID" &&
		req.Quantity > 0 && req.Quantity != -1 &&
		req.UnitPrice > 0 && req.UnitPrice != -1 &&
		req.Material != "" && req.Material != "INVALID" &&
		req.DeliveryDate != nil && req.DeliveryDate.Unix() != 0 &&
		req.VendorCode != "" && req.VendorCode != "INVALID" &&
		req.ShortText != "" && req.ShortText != "INVALID" &&
		req.MaterialGroup != "" && req.MaterialGroup != "INVALID" &&
		req.UnitType != "" && req.UnitType != "INVALID" &&
		req.Requester != "" && req.Requester != "INVALID" &&
		req.PurchaseOrganization != "" && req.PurchaseOrganization != "INVALID" &&
		req.Currency != "" && req.Currency != "INVALID" &&
		totalAmount > 0 && totalAmount != -1 {
		return PurchaseStatusEnough
	}
	return PurchaseStatusLack
}

// parseTimeFlexible 灵活解析时间字符串，支持多种格式
func parseTimeFlexible(timeStr string) (time.Time, error) {
	// 支持的日期时间格式列表
	formats := []string{
		time.RFC3339,                // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05.000Z",  // "2006-01-02T15:04:05.000Z"
		"2006-01-02T15:04:05",       // "2006-01-02T15:04:05"
		"2006-01-02 15:04:05",       // "2006-01-02 15:04:05"
		"2006-01-02",                // "2006-01-02"
		"2006/01/02",                // "2006/01/02"
		"2006-01-02T15:04:05Z",      // "2006-01-02T15:04:05Z"
		"2006-01-02T15:04:05-07:00", // "2006-01-02T15:04:05-07:00"
		"2006-01-02T15:04:05+07:00", // "2006-01-02T15:04:05+07:00"
		"01/02/2006",                // "01/02/2006" (MM/DD/YYYY)
		"02/01/2006",                // "02/01/2006" (DD/MM/YYYY)
		"2006-01-02 15:04",          // "2006-01-02 15:04"
		"2006-01-02T15:04",          // "2006-01-02T15:04"
	}

	// 去除首尾空格
	timeStr = strings.TrimSpace(timeStr)

	// 尝试每种格式
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string '%s' with any supported format", timeStr)
}
