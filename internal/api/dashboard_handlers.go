package api

import (
	"net/http"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DashboardHandler 仪表板API处理器
type DashboardHandler struct {
	prRepo *repository.PurchaseRequestRepository
	poRepo *repository.PurchaseOrderRepository
	logger *logrus.Logger
}

// NewDashboardHandler 创建新的仪表板处理器
func NewDashboardHandler(prRepo *repository.PurchaseRequestRepository, poRepo *repository.PurchaseOrderRepository, logger *logrus.Logger) *DashboardHandler {
	return &DashboardHandler{
		prRepo: prRepo,
		poRepo: poRepo,
		logger: logger,
	}
}

// GetDashboardMetrics 获取仪表板指标
func (h *DashboardHandler) GetDashboardMetrics(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	// 计算时间范围
	var startDate time.Time
	switch period {
	case "7d":
		startDate = time.Now().AddDate(0, 0, -7)
	case "30d":
		startDate = time.Now().AddDate(0, 0, -30)
	case "90d":
		startDate = time.Now().AddDate(0, 0, -90)
	default:
		startDate = time.Now().AddDate(0, 0, -30)
	}

	// 获取今日PR数量
	todayPRQuery := &models.PurchaseRequestQuery{
		StartDate: &startDate,
		Page:      1,
		PageSize:  1000, // 获取所有数据用于统计
	}
	todayPRResult, err := h.prRepo.List(c.Request.Context(), todayPRQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get today PR count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get dashboard metrics",
			},
		})
		return
	}

	// 获取今日PO数量
	todayPOQuery := &models.PurchaseOrderQuery{
		StartDate: &startDate,
		Page:      1,
		PageSize:  1000,
	}
	todayPOResult, err := h.poRepo.List(c.Request.Context(), todayPOQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get today PO count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get dashboard metrics",
			},
		})
		return
	}

	// 获取待审批数量
	pendingStatus := models.PurchaseStatusPending
	pendingQuery := &models.PurchaseRequestQuery{
		Status:   &pendingStatus,
		Page:     1,
		PageSize: 1000,
	}
	pendingResult, err := h.prRepo.List(c.Request.Context(), pendingQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pending approval count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get dashboard metrics",
			},
		})
		return
	}

	// 计算平均审批时间（模拟数据）
	avgApprovalTime := 4.5

	metrics := models.DashboardMetrics{
		TodayPR:         len(todayPRResult.Requests),
		TodayPO:         len(todayPOResult.Items),
		PendingApproval: len(pendingResult.Requests),
		AvgApprovalTime: avgApprovalTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetPRTrendChart 获取PR趋势图表数据
func (h *DashboardHandler) GetPRTrendChart(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	// 计算时间范围和日期列表
	var days int
	switch period {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	default:
		days = 30
	}

	// 生成日期列表和对应的值（模拟数据）
	var dates []string
	var values []int

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, date)
		// 模拟数据：随机生成5-15之间的值
		value := 5 + (i % 10)
		values = append(values, value)
	}

	trendData := models.PRTrendData{
		Dates:  dates,
		Values: values,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trendData,
	})
}

// GetPOStatusDistribution 获取PO状态分布
func (h *DashboardHandler) GetPOStatusDistribution(c *gin.Context) {
	// 获取PO统计
	stats, err := h.poRepo.GetStatistics(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get PO statistics")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get PO status distribution",
			},
		})
		return
	}

	// 转换为图表数据格式
	distribution := []models.POStatusDistribution{
		{Name: "Draft", Value: stats["draft"]},
		{Name: "Sent", Value: stats["sent"]},
		{Name: "Confirmed", Value: stats["confirmed"]},
		{Name: "Delivered", Value: stats["delivered"]},
		{Name: "Cancelled", Value: stats["cancelled"]},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    distribution,
	})
}

// GetSupplierRanking 获取供应商排名
func (h *DashboardHandler) GetSupplierRanking(c *gin.Context) {
	// 模拟供应商排名数据
	ranking := []models.SupplierRanking{
		{
			Supplier:    "ABC Manufacturing Co., Ltd.",
			OrderCount:  15,
			TotalAmount: 500000,
		},
		{
			Supplier:    "XYZ Electronics Inc.",
			OrderCount:  12,
			TotalAmount: 350000,
		},
		{
			Supplier:    "DEF Materials Corp.",
			OrderCount:  8,
			TotalAmount: 200000,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ranking,
	})
}

// GetDepartmentStats 获取部门统计
func (h *DashboardHandler) GetDepartmentStats(c *gin.Context) {
	// 模拟部门统计数据
	stats := []models.DepartmentStats{
		{
			Department: "Procurement Department",
			PRCount:    25,
		},
		{
			Department: "IT Department",
			PRCount:    15,
		},
		{
			Department: "HR Department",
			PRCount:    8,
		},
		{
			Department: "Finance Department",
			PRCount:    12,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
