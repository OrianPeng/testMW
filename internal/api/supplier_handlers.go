package api

import (
	"net/http"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SupplierHandler 供应商API处理器
type SupplierHandler struct {
	repo   *repository.SupplierRepository
	logger *logrus.Logger
}

// NewSupplierHandler 创建新的供应商处理器
func NewSupplierHandler(repo *repository.SupplierRepository, logger *logrus.Logger) *SupplierHandler {
	return &SupplierHandler{
		repo:   repo,
		logger: logger,
	}
}

// CreateSupplier 创建供应商
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req models.CreateSupplierRequest
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
		"code": req.Code,
		"name": req.Name,
		"type": req.Type,
	}).Info("Creating supplier")

	// 创建供应商
	supplier, err := h.repo.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create supplier")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create supplier",
			},
		})
		return
	}

	// 转换为响应格式
	response := h.convertToResponse(supplier)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Supplier created successfully",
	})
}

// GetSupplier 获取供应商详情
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Supplier ID is required",
			},
		})
		return
	}

	// 使用ID查询
	supplier, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Warn("Supplier not found")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Supplier not found",
			},
		})
		return
	}

	response := h.convertToResponse(supplier)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListSuppliers 查询供应商列表
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	// 解析查询参数
	query := &models.SupplierQuery{}

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
	if code := c.Query("code"); code != "" {
		query.Code = &code
	}
	if name := c.Query("name"); name != "" {
		query.Name = &name
	}
	if supplierType := c.Query("type"); supplierType != "" {
		query.Type = (*models.SupplierType)(&supplierType)
	}
	if status := c.Query("status"); status != "" {
		query.Status = (*models.SupplierStatus)(&status)
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
		h.logger.WithError(err).Error("Failed to list suppliers")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to list suppliers",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateSupplier 更新供应商
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Supplier ID is required",
			},
		})
		return
	}

	var updateReq models.UpdateSupplierRequest
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

	// 更新供应商
	supplier, err := h.repo.Update(c.Request.Context(), id, &updateReq)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update supplier")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update supplier",
			},
		})
		return
	}

	response := h.convertToResponse(supplier)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Supplier updated successfully",
	})
}

// DeleteSupplier 删除供应商
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Supplier ID is required",
			},
		})
		return
	}

	// 删除供应商
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("Failed to delete supplier")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete supplier",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Supplier deleted successfully",
	})
}

// GetSupplierStatistics 获取供应商统计
func (h *SupplierHandler) GetSupplierStatistics(c *gin.Context) {
	stats, err := h.repo.GetStatistics(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get supplier statistics")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get supplier statistics",
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
func (h *SupplierHandler) convertToResponse(supplier *models.Supplier) *models.SupplierResponse {
	return &models.SupplierResponse{
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
}
