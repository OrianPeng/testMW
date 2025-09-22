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

// PurchaseRequestRepository 采购请求数据访问层
type PurchaseRequestRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewPurchaseRequestRepository 创建新的采购请求仓库
func NewPurchaseRequestRepository(db *sql.DB, logger *logrus.Logger) *PurchaseRequestRepository {
	return &PurchaseRequestRepository{
		db:     db,
		logger: logger,
	}
}

// Create 创建采购请求
func (r *PurchaseRequestRepository) Create(ctx context.Context, req *models.CreatePurchaseRequestRequest) (*models.PurchaseRequest, error) {
	// 计算总金额
	totalAmount := float64(req.Quantity) * req.UnitPrice

	// 设置默认货币
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	// 检查字段完整性并设置状态
	status := models.CheckCreateRequestCompleteness(req)

	query := `
		INSERT INTO purchase_requests (
			request_id, doc_type, plant, quantity, unit_price, material,
			delivery_date, vendor_code, short_text, material_group, unit_type,
			requester, purchase_organization, currency, total_amount,
			priority, status, urgency, comments, retry_count, error_msg, processed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		req.RequestID,
		req.DocType,
		req.Plant,
		req.Quantity,
		req.UnitPrice,
		req.Material,
		req.DeliveryDate,
		req.VendorCode,
		req.ShortText,
		req.MaterialGroup,
		req.UnitType,
		req.Requester,
		req.PurchaseOrganization,
		currency,
		totalAmount,
		req.Priority,
		status,
		req.Urgency,
		req.Comments,
		0,   // retry_count
		"",  // error_msg
		nil, // processed_at
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase request: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// 查询创建的记录
	return r.GetByID(ctx, id)
}

// GetByID 根据ID获取采购请求
func (r *PurchaseRequestRepository) GetByID(ctx context.Context, id int64) (*models.PurchaseRequest, error) {
	query := `
		SELECT id, request_id, doc_type, plant, quantity, unit_price, material,
			   delivery_date, vendor_code, short_text, material_group, unit_type,
			   requester, purchase_organization, currency, total_amount,
			   priority, status, urgency, approver_id, approver_name, approved_at,
			   rejection_reason, comments, retry_count, error_msg, processed_at, created_at, updated_at
		FROM purchase_requests
		WHERE id = ?
	`

	var req models.PurchaseRequest
	var errorMsg, comments sql.NullString
	var approverID, approverName, rejectionReason sql.NullString
	var approvedAt, processedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.RequestID, &req.DocType, &req.Plant, &req.Quantity, &req.UnitPrice, &req.Material,
		&req.DeliveryDate, &req.VendorCode, &req.ShortText, &req.MaterialGroup, &req.UnitType,
		&req.Requester, &req.PurchaseOrganization, &req.Currency, &req.TotalAmount,
		&req.Priority, &req.Status, &req.Urgency, &approverID, &approverName, &approvedAt,
		&rejectionReason, &comments, &req.RetryCount, &errorMsg, &processedAt, &req.CreatedAt, &req.UpdatedAt,
	)

	// 处理NULL值
	if approverID.Valid {
		req.ApproverID = &approverID.String
	}
	if approverName.Valid {
		req.ApproverName = &approverName.String
	}
	if approvedAt.Valid {
		req.ApprovedAt = &approvedAt.Time
	}
	if rejectionReason.Valid {
		req.RejectionReason = &rejectionReason.String
	}
	if comments.Valid {
		req.Comments = comments.String
	}
	if errorMsg.Valid {
		req.ErrorMsg = errorMsg.String
	}
	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase request not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get purchase request: %w", err)
	}

	return &req, nil
}

// GetByRequestID 根据请求ID获取采购请求
func (r *PurchaseRequestRepository) GetByRequestID(ctx context.Context, requestID string) (*models.PurchaseRequest, error) {
	query := `
		SELECT id, request_id, doc_type, plant, quantity, unit_price, material,
			   delivery_date, vendor_code, short_text, material_group, unit_type,
			   requester, purchase_organization, currency, total_amount,
			   priority, status, urgency, approver_id, approver_name, approved_at,
			   rejection_reason, comments, retry_count, error_msg, processed_at, created_at, updated_at
		FROM purchase_requests
		WHERE request_id = ?
	`

	var req models.PurchaseRequest
	var errorMsg, comments sql.NullString
	var approverID, approverName, rejectionReason sql.NullString
	var approvedAt, processedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.ID, &req.RequestID, &req.DocType, &req.Plant, &req.Quantity, &req.UnitPrice, &req.Material,
		&req.DeliveryDate, &req.VendorCode, &req.ShortText, &req.MaterialGroup, &req.UnitType,
		&req.Requester, &req.PurchaseOrganization, &req.Currency, &req.TotalAmount,
		&req.Priority, &req.Status, &req.Urgency, &approverID, &approverName, &approvedAt,
		&rejectionReason, &comments, &req.RetryCount, &errorMsg, &processedAt, &req.CreatedAt, &req.UpdatedAt,
	)

	// 处理NULL值
	if approverID.Valid {
		req.ApproverID = &approverID.String
	}
	if approverName.Valid {
		req.ApproverName = &approverName.String
	}
	if approvedAt.Valid {
		req.ApprovedAt = &approvedAt.Time
	}
	if rejectionReason.Valid {
		req.RejectionReason = &rejectionReason.String
	}
	if comments.Valid {
		req.Comments = comments.String
	}
	if errorMsg.Valid {
		req.ErrorMsg = errorMsg.String
	}
	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase request not found: %s", requestID)
		}
		return nil, fmt.Errorf("failed to get purchase request: %w", err)
	}

	return &req, nil
}

// List 查询采购请求列表
func (r *PurchaseRequestRepository) List(ctx context.Context, query *models.PurchaseRequestQuery) (*models.PurchaseRequestListResponse, error) {
	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.Status != nil {
		whereClause += " AND status = ?"
		args = append(args, *query.Status)
	}

	if query.Requester != nil {
		whereClause += " AND requester = ?"
		args = append(args, *query.Requester)
	}

	if query.Plant != nil {
		whereClause += " AND plant = ?"
		args = append(args, *query.Plant)
	}

	if query.VendorCode != nil {
		whereClause += " AND vendor_code = ?"
		args = append(args, *query.VendorCode)
	}

	if query.MaterialGroup != nil {
		whereClause += " AND material_group = ?"
		args = append(args, *query.MaterialGroup)
	}

	if query.Priority != nil {
		whereClause += " AND priority = ?"
		args = append(args, *query.Priority)
	}

	if query.Urgency != nil {
		whereClause += " AND urgency = ?"
		args = append(args, *query.Urgency)
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
			"id": true, "request_id": true, "requester": true, "plant": true,
			"material": true, "quantity": true, "unit_price": true, "total_amount": true,
			"priority": true, "status": true, "urgency": true, "delivery_date": true,
			"vendor_code": true, "material_group": true, "created_at": true, "updated_at": true,
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM purchase_requests %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count purchase requests: %w", err)
	}

	// 查询数据
	dataQuery := fmt.Sprintf(`
		SELECT id, request_id, doc_type, plant, quantity, unit_price, material,
			   delivery_date, vendor_code, short_text, material_group, unit_type,
			   requester, purchase_organization, currency, total_amount,
			   priority, status, urgency, approver_id, approver_name, approved_at,
			   rejection_reason, comments, retry_count, error_msg, processed_at, created_at, updated_at
		FROM purchase_requests
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, sortClause)

	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchase requests: %w", err)
	}
	defer rows.Close()

	var requests []*models.PurchaseRequest
	for rows.Next() {
		var req models.PurchaseRequest
		var errorMsg, comments sql.NullString
		var approverID, approverName, rejectionReason sql.NullString
		var approvedAt, processedAt sql.NullTime

		err := rows.Scan(
			&req.ID, &req.RequestID, &req.DocType, &req.Plant, &req.Quantity, &req.UnitPrice, &req.Material,
			&req.DeliveryDate, &req.VendorCode, &req.ShortText, &req.MaterialGroup, &req.UnitType,
			&req.Requester, &req.PurchaseOrganization, &req.Currency, &req.TotalAmount,
			&req.Priority, &req.Status, &req.Urgency, &approverID, &approverName, &approvedAt,
			&rejectionReason, &comments, &req.RetryCount, &errorMsg, &processedAt, &req.CreatedAt, &req.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase request: %w", err)
		}

		// 处理NULL值
		if approverID.Valid {
			req.ApproverID = &approverID.String
		}
		if approverName.Valid {
			req.ApproverName = &approverName.String
		}
		if approvedAt.Valid {
			req.ApprovedAt = &approvedAt.Time
		}
		if rejectionReason.Valid {
			req.RejectionReason = &rejectionReason.String
		}
		if comments.Valid {
			req.Comments = comments.String
		}
		if errorMsg.Valid {
			req.ErrorMsg = errorMsg.String
		}
		if processedAt.Valid {
			req.ProcessedAt = &processedAt.Time
		}

		requests = append(requests, &req)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate purchase requests: %w", err)
	}

	// 转换为响应格式
	var responses []*models.PurchaseRequestResponse
	for _, req := range requests {
		response := &models.PurchaseRequestResponse{
			ID:                   req.ID,
			RequestID:            req.RequestID,
			DocType:              req.DocType,
			Plant:                req.Plant,
			Quantity:             req.Quantity,
			UnitPrice:            req.UnitPrice,
			Material:             req.Material,
			DeliveryDate:         req.DeliveryDate,
			VendorCode:           req.VendorCode,
			ShortText:            req.ShortText,
			MaterialGroup:        req.MaterialGroup,
			UnitType:             req.UnitType,
			Requester:            req.Requester,
			PurchaseOrganization: req.PurchaseOrganization,
			Currency:             req.Currency,
			TotalAmount:          req.TotalAmount,
			Priority:             req.Priority,
			Status:               req.Status,
			Urgency:              req.Urgency,
			ApproverID:           req.ApproverID,
			ApproverName:         req.ApproverName,
			ApprovedAt:           req.ApprovedAt,
			RejectionReason:      req.RejectionReason,
			Comments:             req.Comments,
			RetryCount:           req.RetryCount,
			ErrorMsg:             req.ErrorMsg,
			ProcessedAt:          req.ProcessedAt,
			CreatedAt:            req.CreatedAt,
			UpdatedAt:            req.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return &models.PurchaseRequestListResponse{
		Requests: responses,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// Update 更新采购请求
func (r *PurchaseRequestRepository) Update(ctx context.Context, id int64, updateReq *models.UpdatePurchaseRequestRequest) (*models.PurchaseRequest, error) {
	// 打印接收到的更新请求，便于调试
	r.logger.WithFields(logrus.Fields{
		"id":        id,
		"updateReq": updateReq,
	}).Info("Repository Update called")

	// 构建更新字段
	updateFields := []string{}
	args := []interface{}{}

	// 基本信息字段
	if updateReq.DocType != nil {
		updateFields = append(updateFields, "doc_type = ?")
		args = append(args, *updateReq.DocType)
	}

	if updateReq.Plant != nil {
		updateFields = append(updateFields, "plant = ?")
		args = append(args, *updateReq.Plant)
	}

	if updateReq.Quantity != nil {
		updateFields = append(updateFields, "quantity = ?")
		args = append(args, *updateReq.Quantity)
	}

	if updateReq.UnitPrice != nil {
		updateFields = append(updateFields, "unit_price = ?")
		args = append(args, *updateReq.UnitPrice)
	}

	if updateReq.Material != nil {
		updateFields = append(updateFields, "material = ?")
		args = append(args, *updateReq.Material)
	}

	if updateReq.DeliveryDate != nil {
		updateFields = append(updateFields, "delivery_date = ?")
		args = append(args, *updateReq.DeliveryDate)
	}

	if updateReq.VendorCode != nil {
		updateFields = append(updateFields, "vendor_code = ?")
		args = append(args, *updateReq.VendorCode)
	}

	if updateReq.ShortText != nil {
		updateFields = append(updateFields, "short_text = ?")
		args = append(args, *updateReq.ShortText)
	}

	if updateReq.MaterialGroup != nil {
		updateFields = append(updateFields, "material_group = ?")
		args = append(args, *updateReq.MaterialGroup)
	}

	if updateReq.UnitType != nil {
		updateFields = append(updateFields, "unit_type = ?")
		args = append(args, *updateReq.UnitType)
	}

	if updateReq.Requester != nil {
		updateFields = append(updateFields, "requester = ?")
		args = append(args, *updateReq.Requester)
	}

	if updateReq.PurchaseOrganization != nil {
		updateFields = append(updateFields, "purchase_organization = ?")
		args = append(args, *updateReq.PurchaseOrganization)
	}

	if updateReq.Currency != nil {
		updateFields = append(updateFields, "currency = ?")
		args = append(args, *updateReq.Currency)
	}

	if updateReq.Priority != nil {
		updateFields = append(updateFields, "priority = ?")
		args = append(args, *updateReq.Priority)
	}

	if updateReq.Urgency != nil {
		updateFields = append(updateFields, "urgency = ?")
		args = append(args, *updateReq.Urgency)
	}

	if updateReq.Comments != nil {
		updateFields = append(updateFields, "comments = ?")
		args = append(args, *updateReq.Comments)
	}

	// 状态和审批相关字段
	if updateReq.Status != nil {
		updateFields = append(updateFields, "status = ?")
		args = append(args, string(*updateReq.Status))
	}

	if updateReq.ApproverID != nil {
		updateFields = append(updateFields, "approver_id = ?")
		args = append(args, *updateReq.ApproverID)
	}

	if updateReq.ApproverName != nil {
		updateFields = append(updateFields, "approver_name = ?")
		args = append(args, *updateReq.ApproverName)
	}

	if updateReq.RejectionReason != nil {
		updateFields = append(updateFields, "rejection_reason = ?")
		args = append(args, *updateReq.RejectionReason)
	}

	// 处理相关字段 - 确保正确处理
	if updateReq.RetryCount != nil {
		updateFields = append(updateFields, "retry_count = ?")
		args = append(args, *updateReq.RetryCount)
	}

	if updateReq.ErrorMsg != nil {
		updateFields = append(updateFields, "error_msg = ?")
		args = append(args, *updateReq.ErrorMsg)
	}

	if updateReq.ProcessedAt != nil {
		updateFields = append(updateFields, "processed_at = ?")
		args = append(args, *updateReq.ProcessedAt)
	}

	// 如果状态是已批准，设置批准时间
	if updateReq.Status != nil && *updateReq.Status == models.PurchaseStatusApproved {
		updateFields = append(updateFields, "approved_at = ?")
		args = append(args, time.Now())
	}

	// 获取当前记录（只获取一次，避免重复查询）
	var currentReq *models.PurchaseRequest
	var err error
	if updateReq.Quantity != nil || updateReq.UnitPrice != nil || updateReq.Status == nil {
		currentReq, err = r.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get current purchase request: %w", err)
		}
	}

	// 如果更新了数量或单价，重新计算总金额
	if updateReq.Quantity != nil || updateReq.UnitPrice != nil {
		quantity := currentReq.Quantity
		unitPrice := currentReq.UnitPrice

		if updateReq.Quantity != nil {
			quantity = *updateReq.Quantity
		}
		if updateReq.UnitPrice != nil {
			unitPrice = *updateReq.UnitPrice
		}

		totalAmount := float64(quantity) * unitPrice
		updateFields = append(updateFields, "total_amount = ?")
		args = append(args, totalAmount)
	}

	// 检查是否需要更新状态（如果没有明确设置状态）
	if updateReq.Status == nil {

		// 创建临时请求对象用于检查完整性
		tempReq := &models.PurchaseRequest{
			RequestID:            currentReq.RequestID,
			DocType:              currentReq.DocType,
			Plant:                currentReq.Plant,
			Quantity:             currentReq.Quantity,
			UnitPrice:            currentReq.UnitPrice,
			Material:             currentReq.Material,
			DeliveryDate:         currentReq.DeliveryDate,
			VendorCode:           currentReq.VendorCode,
			ShortText:            currentReq.ShortText,
			MaterialGroup:        currentReq.MaterialGroup,
			UnitType:             currentReq.UnitType,
			Requester:            currentReq.Requester,
			PurchaseOrganization: currentReq.PurchaseOrganization,
			Currency:             currentReq.Currency,
			TotalAmount:          currentReq.TotalAmount,
		}

		// 应用更新字段
		if updateReq.DocType != nil {
			tempReq.DocType = *updateReq.DocType
		}
		if updateReq.Plant != nil {
			tempReq.Plant = *updateReq.Plant
		}
		if updateReq.Quantity != nil {
			tempReq.Quantity = *updateReq.Quantity
		}
		if updateReq.UnitPrice != nil {
			tempReq.UnitPrice = *updateReq.UnitPrice
		}
		if updateReq.Material != nil {
			tempReq.Material = *updateReq.Material
		}
		if updateReq.DeliveryDate != nil {
			tempReq.DeliveryDate = *updateReq.DeliveryDate
		}
		if updateReq.VendorCode != nil {
			tempReq.VendorCode = *updateReq.VendorCode
		}
		if updateReq.ShortText != nil {
			tempReq.ShortText = *updateReq.ShortText
		}
		if updateReq.MaterialGroup != nil {
			tempReq.MaterialGroup = *updateReq.MaterialGroup
		}
		if updateReq.UnitType != nil {
			tempReq.UnitType = *updateReq.UnitType
		}
		if updateReq.Requester != nil {
			tempReq.Requester = *updateReq.Requester
		}
		if updateReq.PurchaseOrganization != nil {
			tempReq.PurchaseOrganization = *updateReq.PurchaseOrganization
		}
		if updateReq.Currency != nil {
			tempReq.Currency = *updateReq.Currency
		}
		if updateReq.Quantity != nil || updateReq.UnitPrice != nil {
			quantity := currentReq.Quantity
			unitPrice := currentReq.UnitPrice
			if updateReq.Quantity != nil {
				quantity = *updateReq.Quantity
			}
			if updateReq.UnitPrice != nil {
				unitPrice = *updateReq.UnitPrice
			}
			tempReq.TotalAmount = float64(quantity) * unitPrice
		}

		// 检查字段完整性并设置状态
		newStatus := models.CheckFieldsCompleteness(tempReq)
		updateFields = append(updateFields, "status = ?")
		args = append(args, string(newStatus))
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updateFields = append(updateFields, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := fmt.Sprintf("UPDATE purchase_requests SET %s WHERE id = ?", strings.Join(updateFields, ", "))

	// 打印 SQL 语句和参数，便于调试
	fmt.Printf("=== SQL DEBUG ===\n")
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Args: %+v\n", args)
	fmt.Printf("================\n")

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).WithField("query", query).Error("Failed to execute update query")
		return nil, fmt.Errorf("failed to update purchase request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("purchase request not found: %d", id)
	}

	// 查询更新后的记录
	return r.GetByID(ctx, id)
}

// UpdateByRequestID 通过request_id更新采购请求
func (r *PurchaseRequestRepository) UpdateByRequestID(ctx context.Context, requestID string, updateReq *models.UpdatePurchaseRequestRequest) (*models.PurchaseRequest, error) {
	// 首先通过request_id获取记录
	purchaseReq, err := r.GetByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("purchase request not found with request_id: %s", requestID)
	}

	// 使用现有的Update方法进行更新
	return r.Update(ctx, purchaseReq.ID, updateReq)
}

// Delete 删除采购请求
func (r *PurchaseRequestRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM purchase_requests WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase request not found: %d", id)
	}

	return nil
}

// DeleteByRequestID 通过request_id删除采购请求
func (r *PurchaseRequestRepository) DeleteByRequestID(ctx context.Context, requestID string) error {
	// 首先通过request_id获取记录
	purchaseReq, err := r.GetByRequestID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("purchase request not found with request_id: %s", requestID)
	}

	// 使用现有的Delete方法进行删除
	return r.Delete(ctx, purchaseReq.ID)
}
