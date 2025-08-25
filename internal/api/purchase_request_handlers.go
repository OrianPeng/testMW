package api

import (
	"fmt"
	"net/http"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PurchaseRequestHandler 采购请求API处理器
type PurchaseRequestHandler struct {
	repo   *repository.PurchaseRequestRepository
	logger *logrus.Logger
}

// NewPurchaseRequestHandler 创建新的采购请求处理器
func NewPurchaseRequestHandler(repo *repository.PurchaseRequestRepository, logger *logrus.Logger) *PurchaseRequestHandler {
	return &PurchaseRequestHandler{
		repo:   repo,
		logger: logger,
	}
}

// CreatePurchaseRequest 创建采购请求
func (h *PurchaseRequestHandler) CreatePurchaseRequest(c *gin.Context) {
	var req models.CreatePurchaseRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// 验证并设置默认值
	if err := req.ValidateAndSetDefaults(); err != nil {
		h.logger.WithError(err).Error("Failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// 生成请求ID（如果没有提供）
	if req.RequestID == "" {
		req.RequestID = "PR-" + uuid.New().String()
	}

	// 记录日志
	h.logger.WithFields(logrus.Fields{
		"request_id":   req.RequestID,
		"requester":    req.Requester,
		"material":     req.Material,
		"quantity":     req.Quantity,
		"total_amount": float64(req.Quantity) * req.UnitPrice,
	}).Info("Creating purchase request")

	// 创建采购请求
	purchaseReq, err := h.repo.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create purchase request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create purchase request",
		})
		return
	}

	// 转换为响应格式
	response := h.convertToResponse(purchaseReq)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Purchase request created successfully",
		"data":    response,
	})
}

// GetPurchaseRequest 获取采购请求详情
func (h *PurchaseRequestHandler) GetPurchaseRequest(c *gin.Context) {
	requestID := c.Param("request_id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Purchase request ID is required",
		})
		return
	}

	// 使用request_id查询
	purchaseReq, err := h.repo.GetByRequestID(c.Request.Context(), requestID)
	if err != nil {
		h.logger.WithError(err).WithField("request_id", requestID).Warn("Purchase request not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Purchase request not found",
		})
		return
	}

	response := h.convertToResponse(purchaseReq)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ListPurchaseRequests 查询采购请求列表
func (h *PurchaseRequestHandler) ListPurchaseRequests(c *gin.Context) {
	// 解析查询参数
	query := &models.PurchaseRequestQuery{}

	// 分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			query.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			query.PageSize = pageSize
		}
	}

	// 过滤参数
	if status := c.Query("status"); status != "" {
		query.Status = (*models.PurchaseRequestStatus)(&status)
	}
	if requester := c.Query("requester"); requester != "" {
		query.Requester = &requester
	}
	if plant := c.Query("plant"); plant != "" {
		query.Plant = &plant
	}
	if vendorCode := c.Query("vendor_code"); vendorCode != "" {
		query.VendorCode = &vendorCode
	}
	if materialGroup := c.Query("material_group"); materialGroup != "" {
		query.MaterialGroup = &materialGroup
	}
	if priorityStr := c.Query("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			query.Priority = &priority
		}
	}
	if urgency := c.Query("urgency"); urgency != "" {
		query.Urgency = &urgency
	}

	// 日期范围
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			query.StartDate = &startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			query.EndDate = &endDate
		}
	}

	// 排序参数
	if sortBy := c.Query("sort_by"); sortBy != "" {
		query.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		query.SortOrder = sortOrder
	}

	// 查询数据
	result, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list purchase requests")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list purchase requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// UpdatePurchaseRequest 更新采购请求
func (h *PurchaseRequestHandler) UpdatePurchaseRequest(c *gin.Context) {
	requestID := c.Param("request_id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Purchase request ID is required",
		})
		return
	}

	var updateReq models.UpdatePurchaseRequestRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 打印 updateReq 结构体内容，便于调试
	h.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"updateReq":  updateReq,
	}).Info("Update request received")

	// 直接打印到控制台，确保能看到调试信息
	fmt.Printf("=== DEBUG: Update request for request_id %s ===\n", requestID)
	fmt.Printf("RetryCount: %+v\n", updateReq.RetryCount)
	fmt.Printf("ErrorMsg: %+v\n", updateReq.ErrorMsg)
	fmt.Printf("ProcessedAt: %+v\n", updateReq.ProcessedAt)
	fmt.Printf("=====================================\n")

	// 记录更新日志
	logFields := logrus.Fields{
		"request_id": requestID,
	}

	// 记录更新的字段
	if updateReq.Status != nil {
		logFields["status"] = *updateReq.Status
	}
	if updateReq.Quantity != nil {
		logFields["quantity"] = *updateReq.Quantity
	}
	if updateReq.UnitPrice != nil {
		logFields["unit_price"] = *updateReq.UnitPrice
	}
	if updateReq.Material != nil {
		logFields["material"] = *updateReq.Material
	}
	if updateReq.VendorCode != nil {
		logFields["vendor_code"] = *updateReq.VendorCode
	}

	h.logger.WithFields(logFields).Info("Updating purchase request")

	// 更新采购请求
	purchaseReq, err := h.repo.UpdateByRequestID(c.Request.Context(), requestID, &updateReq)
	if err != nil {
		h.logger.WithError(err).WithField("request_id", requestID).Error("Failed to update purchase request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update purchase request",
		})
		return
	}

	response := h.convertToResponse(purchaseReq)
	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase request updated successfully",
		"data":    response,
	})
}

// DeletePurchaseRequest 删除采购请求
func (h *PurchaseRequestHandler) DeletePurchaseRequest(c *gin.Context) {
	requestID := c.Param("request_id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Purchase request ID is required",
		})
		return
	}

	// 删除采购请求
	if err := h.repo.DeleteByRequestID(c.Request.Context(), requestID); err != nil {
		h.logger.WithError(err).Error("Failed to delete purchase request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete purchase request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase request deleted successfully",
	})
}

// convertToResponse 转换为响应格式
func (h *PurchaseRequestHandler) convertToResponse(req *models.PurchaseRequest) *models.PurchaseRequestResponse {
	return &models.PurchaseRequestResponse{
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
}
