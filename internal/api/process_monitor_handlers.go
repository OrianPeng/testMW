package api

import (
	"net/http"
	"rpa-middleware/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ProcessMonitorHandler 流程监控API处理器
type ProcessMonitorHandler struct {
	logger *logrus.Logger
}

// NewProcessMonitorHandler 创建新的流程监控处理器
func NewProcessMonitorHandler(logger *logrus.Logger) *ProcessMonitorHandler {
	return &ProcessMonitorHandler{
		logger: logger,
	}
}

// GetProcessOverview 获取流程概览
func (h *ProcessMonitorHandler) GetProcessOverview(c *gin.Context) {
	// 模拟流程概览数据
	overview := models.ProcessOverview{
		TodayPending:    15,
		Overdue:         3,
		AvgApprovalTime: 24,
		CompletionRate:  85,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overview,
	})
}

// GetBlockedProcesses 获取被阻塞的流程
func (h *ProcessMonitorHandler) GetBlockedProcesses(c *gin.Context) {
	// 模拟被阻塞的流程数据
	blockedProcesses := []models.BlockedProcess{
		{
			ID:              "1",
			RequestID:       "PR-20231201-001",
			CurrentStep:     "Department Approval",
			Responsible:     "Mike Johnson",
			Department:      "Procurement Department",
			BlockedDuration: "2 days",
			Priority:        4,
			Reason:          "Waiting for department manager approval",
		},
		{
			ID:              "2",
			RequestID:       "PR-20231201-002",
			CurrentStep:     "Budget Review",
			Responsible:     "Sarah Wilson",
			Department:      "Finance Department",
			BlockedDuration: "1 day",
			Priority:        3,
			Reason:          "Budget allocation pending",
		},
		{
			ID:              "3",
			RequestID:       "PR-20231201-003",
			CurrentStep:     "Technical Review",
			Responsible:     "David Chen",
			Department:      "IT Department",
			BlockedDuration: "3 days",
			Priority:        2,
			Reason:          "Technical specifications under review",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blockedProcesses,
	})
}

// GetProcessSteps 获取流程步骤
func (h *ProcessMonitorHandler) GetProcessSteps(c *gin.Context) {
	processType := c.Query("process_type")
	if processType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Process type is required",
			},
		})
		return
	}

	var steps []models.ProcessStep

	if processType == "pr" {
		// 采购请求流程步骤
		steps = []models.ProcessStep{
			{
				Title:       "Submit Application",
				Status:      models.ProcessStepCompleted,
				Responsible: "Applicant",
				Time:        "2 hours",
				Count:       0,
			},
			{
				Title:       "Department Approval",
				Status:      models.ProcessStepCurrent,
				Responsible: "Department Manager",
				Time:        "4 hours",
				Count:       5,
			},
			{
				Title:       "Budget Review",
				Status:      models.ProcessStepPending,
				Responsible: "Finance Manager",
				Time:        "2 hours",
				Count:       0,
			},
			{
				Title:       "Final Approval",
				Status:      models.ProcessStepPending,
				Responsible: "Procurement Manager",
				Time:        "1 hour",
				Count:       0,
			},
		}
	} else if processType == "po" {
		// 采购订单流程步骤
		steps = []models.ProcessStep{
			{
				Title:       "Create Order",
				Status:      models.ProcessStepCompleted,
				Responsible: "Procurement Staff",
				Time:        "1 hour",
				Count:       0,
			},
			{
				Title:       "Send to Supplier",
				Status:      models.ProcessStepCurrent,
				Responsible: "Procurement Staff",
				Time:        "30 minutes",
				Count:       3,
			},
			{
				Title:       "Supplier Confirmation",
				Status:      models.ProcessStepPending,
				Responsible: "Supplier",
				Time:        "24 hours",
				Count:       0,
			},
			{
				Title:       "Delivery",
				Status:      models.ProcessStepPending,
				Responsible: "Supplier",
				Time:        "3-5 days",
				Count:       0,
			},
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid process type. Must be 'pr' or 'po'",
			},
		})
		return
	}

	response := models.ProcessStepsResponse{
		Steps: steps,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.Steps,
	})
}
