package api

import (
	"net/http"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PurchaseOrderHandler 采购订单API处理器
type PurchaseOrderHandler struct {
	repo   *repository.PurchaseOrderRepository
	logger *logrus.Logger
}

// NewPurchaseOrderHandler 创建新的采购订单处理器
func NewPurchaseOrderHandler(repo *repository.PurchaseOrderRepository, logger *logrus.Logger) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		repo:   repo,
		logger: logger,
	}
}

// CreatePurchaseOrder 创建采购订单
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *gin.Context) {
	var req models.CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}

	// 记录日志
	h.logger.WithFields(logrus.Fields{
		"po_number":     req.PONumber,
		"supplier_name": req.SupplierName,
		"total_amount":  req.TotalAmount,
	}).Info("Creating purchase order")

	// 创建采购订单
	po, err := h.repo.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create purchase order")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create purchase order",
			},
		})
		return
	}

	// 转换为响应格式
	response := h.convertToResponse(po)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Purchase order created successfully",
	})
}

// GetPurchaseOrder 获取采购订单详情
func (h *PurchaseOrderHandler) GetPurchaseOrder(c *gin.Context) {
	poNumber := c.Param("po_number")
	if poNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Purchase order number is required",
			},
		})
		return
	}

	// 使用po_number查询
	po, err := h.repo.GetByPONumber(c.Request.Context(), poNumber)
	if err != nil {
		h.logger.WithError(err).WithField("po_number", poNumber).Warn("Purchase order not found")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Purchase order not found",
			},
		})
		return
	}

	response := h.convertToResponse(po)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListPurchaseOrders 查询采购订单列表
func (h *PurchaseOrderHandler) ListPurchaseOrders(c *gin.Context) {
	// 解析查询参数
	query := &models.PurchaseOrderQuery{}

	// 分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			query.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.PageSize = limit
		}
	}

	// 过滤参数
	if poNumber := c.Query("po_number"); poNumber != "" {
		query.PONumber = &poNumber
	}
	if status := c.Query("status"); status != "" {
		query.Status = (*models.PurchaseOrderStatus)(&status)
	}
	if supplier := c.Query("supplier"); supplier != "" {
		query.Supplier = &supplier
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
		h.logger.WithError(err).Error("Failed to list purchase orders")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to list purchase orders",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdatePurchaseOrder 更新采购订单
func (h *PurchaseOrderHandler) UpdatePurchaseOrder(c *gin.Context) {
	poNumber := c.Param("po_number")
	if poNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Purchase order number is required",
			},
		})
		return
	}

	var updateReq models.UpdatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}

	// 先获取现有记录
	existingPO, err := h.repo.GetByPONumber(c.Request.Context(), poNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Purchase order not found",
			},
		})
		return
	}

	// 更新采购订单
	po, err := h.repo.Update(c.Request.Context(), existingPO.ID, &updateReq)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update purchase order")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update purchase order",
			},
		})
		return
	}

	response := h.convertToResponse(po)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Purchase order updated successfully",
	})
}

// DeletePurchaseOrder 删除采购订单
func (h *PurchaseOrderHandler) DeletePurchaseOrder(c *gin.Context) {
	poNumber := c.Param("po_number")
	if poNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Purchase order number is required",
			},
		})
		return
	}

	// 先获取现有记录
	existingPO, err := h.repo.GetByPONumber(c.Request.Context(), poNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Purchase order not found",
			},
		})
		return
	}

	// 删除采购订单
	if err := h.repo.Delete(c.Request.Context(), existingPO.ID); err != nil {
		h.logger.WithError(err).Error("Failed to delete purchase order")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete purchase order",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Purchase order deleted successfully",
	})
}

// GetPOStatistics 获取采购订单统计
func (h *PurchaseOrderHandler) GetPOStatistics(c *gin.Context) {
	stats, err := h.repo.GetStatistics(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get purchase order statistics")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get purchase order statistics",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// convertToResponse 转换为响应格式
func (h *PurchaseOrderHandler) convertToResponse(po *models.PurchaseOrder) *models.PurchaseOrderResponse {
	return &models.PurchaseOrderResponse{
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
}
