package models

import "time"

// SupplierType 供应商类型
type SupplierType string

const (
	SupplierTypeManufacturer SupplierType = "manufacturer" // 制造商
	SupplierTypeDistributor  SupplierType = "distributor"  // 分销商
	SupplierTypeService      SupplierType = "service"      // 服务商
)

// SupplierStatus 供应商状态
type SupplierStatus string

const (
	SupplierStatusActive   SupplierStatus = "active"   // 活跃
	SupplierStatusInactive SupplierStatus = "inactive" // 非活跃
	SupplierStatusPending  SupplierStatus = "pending"  // 待审核
)

// Supplier 供应商
type Supplier struct {
	ID            string         `json:"id" db:"id"`
	Code          string         `json:"code" db:"code" binding:"required"`
	Name          string         `json:"name" db:"name" binding:"required"`
	Type          SupplierType   `json:"type" db:"type"`
	Status        SupplierStatus `json:"status" db:"status"`
	ContactPerson string         `json:"contact_person" db:"contact_person"`
	Phone         string         `json:"phone" db:"phone"`
	Email         string         `json:"email" db:"email"`
	Address       string         `json:"address" db:"address"`
	Categories    []string       `json:"categories" db:"categories"`
	Notes         string         `json:"notes" db:"notes"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
}

// CreateSupplierRequest 创建供应商的请求
type CreateSupplierRequest struct {
	Code          string         `json:"code" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	Type          SupplierType   `json:"type"`
	Status        SupplierStatus `json:"status"`
	ContactPerson string         `json:"contact_person"`
	Phone         string         `json:"phone"`
	Email         string         `json:"email"`
	Address       string         `json:"address"`
	Categories    []string       `json:"categories"`
	Notes         string         `json:"notes"`
}

// UpdateSupplierRequest 更新供应商的请求
type UpdateSupplierRequest struct {
	Code          *string         `json:"code"`
	Name          *string         `json:"name"`
	Type          *SupplierType   `json:"type"`
	Status        *SupplierStatus `json:"status"`
	ContactPerson *string         `json:"contact_person"`
	Phone         *string         `json:"phone"`
	Email         *string         `json:"email"`
	Address       *string         `json:"address"`
	Categories    *[]string       `json:"categories"`
	Notes         *string         `json:"notes"`
}

// SupplierResponse 供应商响应
type SupplierResponse struct {
	ID            string         `json:"id"`
	Code          string         `json:"code"`
	Name          string         `json:"name"`
	Type          SupplierType   `json:"type"`
	Status        SupplierStatus `json:"status"`
	ContactPerson string         `json:"contact_person"`
	Phone         string         `json:"phone"`
	Email         string         `json:"email"`
	Address       string         `json:"address"`
	Categories    []string       `json:"categories"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
}

// SupplierListResponse 供应商列表响应
type SupplierListResponse struct {
	Items      []*SupplierResponse `json:"items"`
	Pagination PaginationResponse  `json:"pagination"`
}

// SupplierQuery 供应商查询条件
type SupplierQuery struct {
	Code      *string         `json:"code"`
	Name      *string         `json:"name"`
	Type      *SupplierType   `json:"type"`
	Status    *SupplierStatus `json:"status"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	SortBy    string          `json:"sort_by"`
	SortOrder string          `json:"sort_order"`
}
