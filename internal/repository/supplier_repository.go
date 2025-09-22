package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"rpa-middleware/internal/models"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// SupplierRepository 供应商数据访问层
type SupplierRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewSupplierRepository 创建新的供应商仓库
func NewSupplierRepository(db *sql.DB, logger *logrus.Logger) *SupplierRepository {
	return &SupplierRepository{
		db:     db,
		logger: logger,
	}
}

// Create 创建供应商
func (r *SupplierRepository) Create(ctx context.Context, req *models.CreateSupplierRequest) (*models.Supplier, error) {
	// 设置默认值
	supplierType := req.Type
	if supplierType == "" {
		supplierType = models.SupplierTypeManufacturer
	}

	status := req.Status
	if status == "" {
		status = models.SupplierStatusActive
	}

	// 序列化分类
	categoriesJSON, err := json.Marshal(req.Categories)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal categories: %w", err)
	}

	// 生成ID
	id := fmt.Sprintf("SUP_%d", time.Now().Unix())

	query := `
		INSERT INTO suppliers (
			id, code, name, type, status, contact_person, phone, email,
			address, categories, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		id,
		req.Code,
		req.Name,
		supplierType,
		status,
		req.ContactPerson,
		req.Phone,
		req.Email,
		req.Address,
		categoriesJSON,
		req.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	// 查询创建的记录
	return r.GetByID(ctx, id)
}

// GetByID 根据ID获取供应商
func (r *SupplierRepository) GetByID(ctx context.Context, id string) (*models.Supplier, error) {
	query := `
		SELECT id, code, name, type, status, contact_person, phone, email,
			   address, categories, notes, created_at
		FROM suppliers
		WHERE id = ?
	`

	var supplier models.Supplier
	var categoriesJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&supplier.ID, &supplier.Code, &supplier.Name, &supplier.Type, &supplier.Status,
		&supplier.ContactPerson, &supplier.Phone, &supplier.Email, &supplier.Address,
		&categoriesJSON, &supplier.Notes, &supplier.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	// 反序列化分类
	if categoriesJSON != "" {
		if err := json.Unmarshal([]byte(categoriesJSON), &supplier.Categories); err != nil {
			r.logger.WithError(err).Warn("Failed to unmarshal categories")
			supplier.Categories = []string{}
		}
	}

	return &supplier, nil
}

// GetByCode 根据代码获取供应商
func (r *SupplierRepository) GetByCode(ctx context.Context, code string) (*models.Supplier, error) {
	query := `
		SELECT id, code, name, type, status, contact_person, phone, email,
			   address, categories, notes, created_at
		FROM suppliers
		WHERE code = ?
	`

	var supplier models.Supplier
	var categoriesJSON string

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&supplier.ID, &supplier.Code, &supplier.Name, &supplier.Type, &supplier.Status,
		&supplier.ContactPerson, &supplier.Phone, &supplier.Email, &supplier.Address,
		&categoriesJSON, &supplier.Notes, &supplier.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	// 反序列化分类
	if categoriesJSON != "" {
		if err := json.Unmarshal([]byte(categoriesJSON), &supplier.Categories); err != nil {
			r.logger.WithError(err).Warn("Failed to unmarshal categories")
			supplier.Categories = []string{}
		}
	}

	return &supplier, nil
}

// List 查询供应商列表
func (r *SupplierRepository) List(ctx context.Context, query *models.SupplierQuery) (*models.SupplierListResponse, error) {
	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.Code != nil {
		whereClause += " AND code LIKE ?"
		args = append(args, "%"+*query.Code+"%")
	}

	if query.Name != nil {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+*query.Name+"%")
	}

	if query.Type != nil {
		whereClause += " AND type = ?"
		args = append(args, *query.Type)
	}

	if query.Status != nil {
		whereClause += " AND status = ?"
		args = append(args, *query.Status)
	}

	// 构建排序
	sortClause := "ORDER BY created_at DESC"
	if query.SortBy != "" {
		validSortFields := map[string]bool{
			"id": true, "code": true, "name": true, "type": true, "status": true, "created_at": true,
		}

		if validSortFields[query.SortBy] {
			sortOrder := "DESC"
			if strings.ToUpper(query.SortOrder) == "ASC" {
				sortOrder = "ASC"
			}
			sortClause = fmt.Sprintf("ORDER BY %s %s", query.SortBy, sortOrder)
		}
	}

	// 设置分页默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	offset := (query.Page - 1) * query.PageSize

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM suppliers %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count suppliers: %w", err)
	}

	// 查询数据
	dataQuery := fmt.Sprintf(`
		SELECT id, code, name, type, status, contact_person, phone, email,
			   address, categories, notes, created_at
		FROM suppliers
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, sortClause)

	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		var supplier models.Supplier
		var categoriesJSON string

		err := rows.Scan(
			&supplier.ID, &supplier.Code, &supplier.Name, &supplier.Type, &supplier.Status,
			&supplier.ContactPerson, &supplier.Phone, &supplier.Email, &supplier.Address,
			&categoriesJSON, &supplier.Notes, &supplier.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier: %w", err)
		}

		// 反序列化分类
		if categoriesJSON != "" {
			if err := json.Unmarshal([]byte(categoriesJSON), &supplier.Categories); err != nil {
				r.logger.WithError(err).Warn("Failed to unmarshal categories")
				supplier.Categories = []string{}
			}
		}

		suppliers = append(suppliers, &supplier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate suppliers: %w", err)
	}

	// 转换为响应格式
	var responses []*models.SupplierResponse
	for _, supplier := range suppliers {
		response := &models.SupplierResponse{
			ID:            supplier.ID,
			Code:          supplier.Code,
			Name:          supplier.Name,
			Type:          supplier.Type,
			Status:        supplier.Status,
			ContactPerson: supplier.ContactPerson,
			Phone:         supplier.Phone,
			Email:         supplier.Email,
			Address:       supplier.Address,
			Categories:    supplier.Categories,
			Notes:         supplier.Notes,
			CreatedAt:     supplier.CreatedAt,
		}
		responses = append(responses, response)
	}

	// 计算总页数
	pages := int(total) / query.PageSize
	if int(total)%query.PageSize > 0 {
		pages++
	}

	return &models.SupplierListResponse{
		Items: responses,
		Pagination: models.PaginationResponse{
			Current:  query.Page,
			PageSize: query.PageSize,
			Total:    int(total),
			Pages:    pages,
		},
	}, nil
}

// Update 更新供应商
func (r *SupplierRepository) Update(ctx context.Context, id string, updateReq *models.UpdateSupplierRequest) (*models.Supplier, error) {
	// 构建更新字段
	updateFields := []string{}
	args := []interface{}{}

	if updateReq.Code != nil {
		updateFields = append(updateFields, "code = ?")
		args = append(args, *updateReq.Code)
	}

	if updateReq.Name != nil {
		updateFields = append(updateFields, "name = ?")
		args = append(args, *updateReq.Name)
	}

	if updateReq.Type != nil {
		updateFields = append(updateFields, "type = ?")
		args = append(args, *updateReq.Type)
	}

	if updateReq.Status != nil {
		updateFields = append(updateFields, "status = ?")
		args = append(args, *updateReq.Status)
	}

	if updateReq.ContactPerson != nil {
		updateFields = append(updateFields, "contact_person = ?")
		args = append(args, *updateReq.ContactPerson)
	}

	if updateReq.Phone != nil {
		updateFields = append(updateFields, "phone = ?")
		args = append(args, *updateReq.Phone)
	}

	if updateReq.Email != nil {
		updateFields = append(updateFields, "email = ?")
		args = append(args, *updateReq.Email)
	}

	if updateReq.Address != nil {
		updateFields = append(updateFields, "address = ?")
		args = append(args, *updateReq.Address)
	}

	if updateReq.Categories != nil {
		categoriesJSON, err := json.Marshal(*updateReq.Categories)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal categories: %w", err)
		}
		updateFields = append(updateFields, "categories = ?")
		args = append(args, categoriesJSON)
	}

	if updateReq.Notes != nil {
		updateFields = append(updateFields, "notes = ?")
		args = append(args, *updateReq.Notes)
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updateFields = append(updateFields, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := fmt.Sprintf("UPDATE suppliers SET %s WHERE id = ?", strings.Join(updateFields, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("supplier not found: %s", id)
	}

	// 查询更新后的记录
	return r.GetByID(ctx, id)
}

// Delete 删除供应商
func (r *SupplierRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM suppliers WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier not found: %s", id)
	}

	return nil
}

// GetStatistics 获取供应商统计
func (r *SupplierRepository) GetStatistics(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END) as inactive,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH) THEN 1 ELSE 0 END) as added_this_month
		FROM suppliers
	`

	var total, active, inactive, pending, addedThisMonth int
	err := r.db.QueryRowContext(ctx, query).Scan(&total, &active, &inactive, &pending, &addedThisMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier statistics: %w", err)
	}

	return map[string]int{
		"total":            total,
		"active":           active,
		"inactive":         inactive,
		"pending":          pending,
		"added_this_month": addedThisMonth,
	}, nil
}
