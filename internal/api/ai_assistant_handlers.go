package api

import (
	"net/http"
	"rpa-middleware/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AIAssistantHandler AI助手API处理器
type AIAssistantHandler struct {
	logger *logrus.Logger
}

// NewAIAssistantHandler 创建新的AI助手处理器
func NewAIAssistantHandler(logger *logrus.Logger) *AIAssistantHandler {
	return &AIAssistantHandler{
		logger: logger,
	}
}

// SendMessage 发送消息给AI助手
func (h *AIAssistantHandler) SendMessage(c *gin.Context) {
	var req models.AIMessageRequest
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

	// 生成对话ID（如果没有提供）
	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = "conv_" + uuid.New().String()
	}

	// 记录日志
	h.logger.WithFields(logrus.Fields{
		"conversation_id": conversationID,
		"message":         req.Message,
	}).Info("Processing AI assistant message")

	// 模拟AI响应
	response := h.generateAIResponse(req.Message, conversationID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// generateAIResponse 生成AI响应（模拟）
func (h *AIAssistantHandler) generateAIResponse(message, conversationID string) models.AIMessageResponse {
	// 简单的关键词匹配来生成不同类型的响应
	messageLower := strings.ToLower(message)

	var response interface{}
	var responseType models.AIResponseType
	var suggestions []string

	if strings.Contains(messageLower, "create") && strings.Contains(messageLower, "pr") {
		// 创建PR的响应
		response = models.AITextResponse{
			Type:    models.AIResponseTypeText,
			Content: "I'll help you create a purchase request. Please provide the following information: material description, quantity, unit price, supplier, and delivery date.",
		}
		responseType = models.AIResponseTypeText
		suggestions = []string{"Create PR", "View PR Template", "Check Supplier List"}

	} else if strings.Contains(messageLower, "status") || strings.Contains(messageLower, "check") {
		// 查询状态的响应
		response = models.AITableResponse{
			Type: models.AIResponseTypeTable,
			TableColumns: []models.AITableColumn{
				{Title: "PR Number", DataIndex: "pr_number", Key: "pr_number"},
				{Title: "Status", DataIndex: "status", Key: "status"},
				{Title: "Requester", DataIndex: "requester", Key: "requester"},
			},
			TableData: []models.AITableData{
				{
					Key: "1",
					Data: map[string]interface{}{
						"pr_number": "PR-20231201-001",
						"status":    "Pending",
						"requester": "John Smith",
					},
				},
				{
					Key: "2",
					Data: map[string]interface{}{
						"pr_number": "PR-20231201-002",
						"status":    "Approved",
						"requester": "Alice Johnson",
					},
				},
			},
		}
		responseType = models.AIResponseTypeTable
		suggestions = []string{"View Details", "Update Status", "Export Data"}

	} else if strings.Contains(messageLower, "process") || strings.Contains(messageLower, "workflow") {
		// 流程相关的响应
		response = models.AIProcessResponse{
			Type: models.AIResponseTypeProcess,
			ProcessSteps: []models.AIProcessStep{
				{
					Title:       "Submit Application",
					Status:      "completed",
					Responsible: "Applicant",
				},
				{
					Title:       "Department Approval",
					Status:      "current",
					Responsible: "Department Manager",
				},
				{
					Title:       "Budget Review",
					Status:      "pending",
					Responsible: "Finance Manager",
				},
			},
		}
		responseType = models.AIResponseTypeProcess
		suggestions = []string{"View Process Details", "Check Blocked Items", "Process Statistics"}

	} else if strings.Contains(messageLower, "supplier") || strings.Contains(messageLower, "vendor") {
		// 供应商相关的响应
		response = models.AICardResponse{
			Type:      models.AIResponseTypeCard,
			CardTitle: "Supplier Information",
			CardData: map[string]string{
				"Supplier": "ABC Manufacturing Co., Ltd.",
				"Code":     "SUP001",
				"Status":   "Active",
				"Contact":  "John Smith",
			},
		}
		responseType = models.AIResponseTypeCard
		suggestions = []string{"View Supplier Details", "Check Orders", "Update Information"}

	} else {
		// 默认文本响应
		response = models.AITextResponse{
			Type:    models.AIResponseTypeText,
			Content: "I'm here to help you with procurement management. You can ask me about creating purchase requests, checking status, managing suppliers, or viewing process workflows. How can I assist you today?",
		}
		responseType = models.AIResponseTypeText
		suggestions = []string{"Create PR", "Check Status", "View Suppliers", "Process Monitor"}
	}

	return models.AIMessageResponse{
		Response:       response,
		Type:           responseType,
		ConversationID: conversationID,
		Suggestions:    suggestions,
	}
}
