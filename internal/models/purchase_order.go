package models

import "time"

// PurchaseOrderStatus 采购订单状态
type PurchaseOrderStatus string

const (
	POStatusDraft     PurchaseOrderStatus = "draft"     // 草稿
	POStatusSent      PurchaseOrderStatus = "sent"      // 已发送
	POStatusConfirmed PurchaseOrderStatus = "confirmed" // 已确认
	POStatusDelivered PurchaseOrderStatus = "delivered" // 已交付
	POStatusCancelled PurchaseOrderStatus = "cancelled" // 已取消
)

// PurchaseOrder 采购订单
type PurchaseOrder struct {
	ID           int64               `json:"id" db:"id"`
	PONumber     string              `json:"po_number" db:"po_number" binding:"required"`
	Status       PurchaseOrderStatus `json:"status" db:"status"`
	SupplierName string              `json:"supplier_name" db:"supplier_name" binding:"required"`
	SupplierCode string              `json:"supplier_code" db:"supplier_code" binding:"required"`
	TotalAmount  float64             `json:"total_amount" db:"total_amount"`
	Currency     string              `json:"currency" db:"currency"`
	CreatedBy    string              `json:"created_by" db:"created_by"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	SentAt       *time.Time          `json:"sent_at" db:"sent_at"`
	ConfirmedAt  *time.Time          `json:"confirmed_at" db:"confirmed_at"`
	Notes        string              `json:"notes" db:"notes"`
}

// CreatePurchaseOrderRequest 创建采购订单的请求
type CreatePurchaseOrderRequest struct {
	PONumber     string              `json:"po_number" binding:"required"`
	Status       PurchaseOrderStatus `json:"status"`
	SupplierName string              `json:"supplier_name" binding:"required"`
	SupplierCode string              `json:"supplier_code" binding:"required"`
	TotalAmount  float64             `json:"total_amount" binding:"required,min=0"`
	Currency     string              `json:"currency"`
	CreatedBy    string              `json:"created_by" binding:"required"`
	Notes        string              `json:"notes"`
}

// UpdatePurchaseOrderRequest 更新采购订单的请求
type UpdatePurchaseOrderRequest struct {
	Status       *PurchaseOrderStatus `json:"status"`
	SupplierName *string              `json:"supplier_name"`
	SupplierCode *string              `json:"supplier_code"`
	TotalAmount  *float64             `json:"total_amount"`
	Currency     *string              `json:"currency"`
	Notes        *string              `json:"notes"`
	SentAt       *time.Time           `json:"sent_at"`
	ConfirmedAt  *time.Time           `json:"confirmed_at"`
}

// PurchaseOrderResponse 采购订单响应
type PurchaseOrderResponse struct {
	ID           int64               `json:"id"`
	PONumber     string              `json:"po_number"`
	Status       PurchaseOrderStatus `json:"status"`
	SupplierName string              `json:"supplier_name"`
	SupplierCode string              `json:"supplier_code"`
	TotalAmount  float64             `json:"total_amount"`
	Currency     string              `json:"currency"`
	CreatedBy    string              `json:"created_by"`
	CreatedAt    time.Time           `json:"created_at"`
	SentAt       *time.Time          `json:"sent_at"`
	ConfirmedAt  *time.Time          `json:"confirmed_at"`
	Notes        string              `json:"notes"`
}

// PurchaseOrderListResponse 采购订单列表响应
type PurchaseOrderListResponse struct {
	Items      []*PurchaseOrderResponse `json:"items"`
	Pagination PaginationResponse       `json:"pagination"`
}

// PurchaseOrderQuery 采购订单查询条件
type PurchaseOrderQuery struct {
	PONumber  *string              `json:"po_number"`
	Status    *PurchaseOrderStatus `json:"status"`
	Supplier  *string              `json:"supplier"`
	StartDate *time.Time           `json:"start_date"`
	EndDate   *time.Time           `json:"end_date"`
	Page      int                  `json:"page"`
	PageSize  int                  `json:"page_size"`
	SortBy    string               `json:"sort_by"`
	SortOrder string               `json:"sort_order"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Current  int `json:"current"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}
