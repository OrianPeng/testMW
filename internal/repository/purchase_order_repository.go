package repository

import (
	"context"
	"database/sql"
	"fmt"
	"rpa-middleware/internal/models"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// PurchaseOrderRepository 采购订单数据访问层
type PurchaseOrderRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewPurchaseOrderRepository 创建新的采购订单仓库
func NewPurchaseOrderRepository(db *sql.DB, logger *logrus.Logger) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{
		db:     db,
		logger: logger,
	}
}

// Create 创建采购订单
func (r *PurchaseOrderRepository) Create(ctx context.Context, req *models.CreatePurchaseOrderRequest) (*models.PurchaseOrder, error) {
	// 设置默认货币
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	// 设置默认状态
	status := req.Status
	if status == "" {
		status = models.POStatusDraft
	}

	query := `
		INSERT INTO purchase_orders (
			po_number, status, supplier_name, supplier_code, total_amount,
			currency, created_by, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		req.PONumber,
		status,
		req.SupplierName,
		req.SupplierCode,
		req.TotalAmount,
		currency,
		req.CreatedBy,
		req.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// 查询创建的记录
	return r.GetByID(ctx, id)
}

// GetByID 根据ID获取采购订单
func (r *PurchaseOrderRepository) GetByID(ctx context.Context, id int64) (*models.PurchaseOrder, error) {
	query := `
		SELECT id, po_number, status, supplier_name, supplier_code, total_amount,
			   currency, created_by, created_at, sent_at, confirmed_at, notes
		FROM purchase_orders
		WHERE id = ?
	`

	var po models.PurchaseOrder

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&po.ID, &po.PONumber, &po.Status, &po.SupplierName, &po.SupplierCode, &po.TotalAmount,
		&po.Currency, &po.CreatedBy, &po.CreatedAt, &po.SentAt, &po.ConfirmedAt, &po.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return &po, nil
}

// GetByPONumber 根据PO编号获取采购订单
func (r *PurchaseOrderRepository) GetByPONumber(ctx context.Context, poNumber string) (*models.PurchaseOrder, error) {
	query := `
		SELECT id, po_number, status, supplier_name, supplier_code, total_amount,
			   currency, created_by, created_at, sent_at, confirmed_at, notes
		FROM purchase_orders
		WHERE po_number = ?
	`

	var po models.PurchaseOrder

	err := r.db.QueryRowContext(ctx, query, poNumber).Scan(
		&po.ID, &po.PONumber, &po.Status, &po.SupplierName, &po.SupplierCode, &po.TotalAmount,
		&po.Currency, &po.CreatedBy, &po.CreatedAt, &po.SentAt, &po.ConfirmedAt, &po.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found: %s", poNumber)
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return &po, nil
}

// List 查询采购订单列表
func (r *PurchaseOrderRepository) List(ctx context.Context, query *models.PurchaseOrderQuery) (*models.PurchaseOrderListResponse, error) {
	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.PONumber != nil {
		whereClause += " AND po_number LIKE ?"
		args = append(args, "%"+*query.PONumber+"%")
	}

	if query.Status != nil {
		whereClause += " AND status = ?"
		args = append(args, *query.Status)
	}

	if query.Supplier != nil {
		whereClause += " AND (supplier_name LIKE ? OR supplier_code LIKE ?)"
		args = append(args, "%"+*query.Supplier+"%", "%"+*query.Supplier+"%")
	}

	if query.StartDate != nil {
		whereClause += " AND created_at >= ?"
		args = append(args, *query.StartDate)
	}

	if query.EndDate != nil {
		whereClause += " AND created_at <= ?"
		args = append(args, *query.EndDate)
	}

	// 构建排序
	sortClause := "ORDER BY created_at DESC"
	if query.SortBy != "" {
		validSortFields := map[string]bool{
			"id": true, "po_number": true, "supplier_name": true, "total_amount": true,
			"status": true, "created_at": true, "sent_at": true, "confirmed_at": true,
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM purchase_orders %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	// 查询数据
	dataQuery := fmt.Sprintf(`
		SELECT id, po_number, status, supplier_name, supplier_code, total_amount,
			   currency, created_by, created_at, sent_at, confirmed_at, notes
		FROM purchase_orders
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, sortClause)

	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchase orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.PurchaseOrder
	for rows.Next() {
		var po models.PurchaseOrder

		err := rows.Scan(
			&po.ID, &po.PONumber, &po.Status, &po.SupplierName, &po.SupplierCode, &po.TotalAmount,
			&po.Currency, &po.CreatedBy, &po.CreatedAt, &po.SentAt, &po.ConfirmedAt, &po.Notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order: %w", err)
		}

		orders = append(orders, &po)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate purchase orders: %w", err)
	}

	// 转换为响应格式
	var responses []*models.PurchaseOrderResponse
	for _, po := range orders {
		response := &models.PurchaseOrderResponse{
			ID:           po.ID,
			PONumber:     po.PONumber,
			Status:       po.Status,
			SupplierName: po.SupplierName,
			SupplierCode: po.SupplierCode,
			TotalAmount:  po.TotalAmount,
			Currency:     po.Currency,
			CreatedBy:    po.CreatedBy,
			CreatedAt:    po.CreatedAt,
			SentAt:       po.SentAt,
			ConfirmedAt:  po.ConfirmedAt,
			Notes:        po.Notes,
		}
		responses = append(responses, response)
	}

	// 计算总页数
	pages := int(total) / query.PageSize
	if int(total)%query.PageSize > 0 {
		pages++
	}

	return &models.PurchaseOrderListResponse{
		Items: responses,
		Pagination: models.PaginationResponse{
			Current:  query.Page,
			PageSize: query.PageSize,
			Total:    int(total),
			Pages:    pages,
		},
	}, nil
}

// Update 更新采购订单
func (r *PurchaseOrderRepository) Update(ctx context.Context, id int64, updateReq *models.UpdatePurchaseOrderRequest) (*models.PurchaseOrder, error) {
	// 构建更新字段
	updateFields := []string{}
	args := []interface{}{}

	if updateReq.Status != nil {
		updateFields = append(updateFields, "status = ?")
		args = append(args, *updateReq.Status)
	}

	if updateReq.SupplierName != nil {
		updateFields = append(updateFields, "supplier_name = ?")
		args = append(args, *updateReq.SupplierName)
	}

	if updateReq.SupplierCode != nil {
		updateFields = append(updateFields, "supplier_code = ?")
		args = append(args, *updateReq.SupplierCode)
	}

	if updateReq.TotalAmount != nil {
		updateFields = append(updateFields, "total_amount = ?")
		args = append(args, *updateReq.TotalAmount)
	}

	if updateReq.Currency != nil {
		updateFields = append(updateFields, "currency = ?")
		args = append(args, *updateReq.Currency)
	}

	if updateReq.Notes != nil {
		updateFields = append(updateFields, "notes = ?")
		args = append(args, *updateReq.Notes)
	}

	if updateReq.SentAt != nil {
		updateFields = append(updateFields, "sent_at = ?")
		args = append(args, *updateReq.SentAt)
	}

	if updateReq.ConfirmedAt != nil {
		updateFields = append(updateFields, "confirmed_at = ?")
		args = append(args, *updateReq.ConfirmedAt)
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updateFields = append(updateFields, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := fmt.Sprintf("UPDATE purchase_orders SET %s WHERE id = ?", strings.Join(updateFields, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("purchase order not found: %d", id)
	}

	// 查询更新后的记录
	return r.GetByID(ctx, id)
}

// Delete 删除采购订单
func (r *PurchaseOrderRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM purchase_orders WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found: %d", id)
	}

	return nil
}

// GetStatistics 获取采购订单统计
func (r *PurchaseOrderRepository) GetStatistics(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT 
			SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END) as draft,
			SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END) as sent,
			SUM(CASE WHEN status = 'confirmed' THEN 1 ELSE 0 END) as confirmed,
			SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END) as delivered,
			SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) as cancelled,
			COUNT(*) as total
		FROM purchase_orders
	`

	var draft, sent, confirmed, delivered, cancelled, total int
	err := r.db.QueryRowContext(ctx, query).Scan(&draft, &sent, &confirmed, &delivered, &cancelled, &total)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order statistics: %w", err)
	}

	return map[string]int{
		"draft":     draft,
		"sent":      sent,
		"confirmed": confirmed,
		"delivered": delivered,
		"cancelled": cancelled,
		"total":     total,
	}, nil
}
