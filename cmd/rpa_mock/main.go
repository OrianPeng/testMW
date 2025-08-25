package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// MockRPARequest 模拟RPA请求
type MockRPARequest struct {
	ID       string                 `json:"id"`
	AgentID  string                 `json:"agent_id"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// MockRPAResponse 模拟RPA响应
type MockRPAResponse struct {
	ID      string                 `json:"id"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// MockRPAStatus 模拟RPA状态
type MockRPAStatus struct {
	IsAvailable bool      `json:"is_available"`
	LastCheck   time.Time `json:"last_check"`
	CurrentTask string    `json:"current_task,omitempty"`
	Version     string    `json:"version"`
	Uptime      string    `json:"uptime"`
}

// MockRPAServer 模拟RPA服务器
type MockRPAServer struct {
	logger       *logrus.Logger
	startTime    time.Time
	currentTask  string
	taskMutex    sync.RWMutex
	requestCount int
	countMutex   sync.Mutex
}

// NewMockRPAServer 创建新的模拟RPA服务器
func NewMockRPAServer() *MockRPAServer {
	return &MockRPAServer{
		logger:    logrus.New(),
		startTime: time.Now(),
	}
}

// executeHandler 处理RPA执行请求
func (s *MockRPAServer) executeHandler(w http.ResponseWriter, r *http.Request) {
	s.countMutex.Lock()
	s.requestCount++
	count := s.requestCount
	s.countMutex.Unlock()

	// 解析请求
	var req MockRPARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.WithError(err).Error("Failed to decode request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"request_id": req.ID,
		"agent_id":   req.AgentID,
		"count":      count,
	}).Info("Received RPA execution request")

	// 设置当前任务
	s.taskMutex.Lock()
	s.currentTask = fmt.Sprintf("Processing request %s from agent %s", req.ID, req.AgentID)
	s.taskMutex.Unlock()

	// 模拟处理时间
	processingTime := s.simulateProcessingTime(req)
	time.Sleep(processingTime)

	// 生成响应
	response := s.generateResponse(req, count)

	// 清除当前任务
	s.taskMutex.Lock()
	s.currentTask = ""
	s.taskMutex.Unlock()

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	s.logger.WithFields(logrus.Fields{
		"request_id":      req.ID,
		"success":         response.Success,
		"processing_time": processingTime,
	}).Info("RPA execution completed")
}

// statusHandler 处理状态检查请求
func (s *MockRPAServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	s.taskMutex.RLock()
	currentTask := s.currentTask
	s.taskMutex.RUnlock()

	status := MockRPAStatus{
		IsAvailable: true,
		LastCheck:   time.Now(),
		CurrentTask: currentTask,
		Version:     "1.0.0-mock",
		Uptime:      time.Since(s.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// healthHandler 健康检查
func (s *MockRPAServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(s.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// simulateProcessingTime 模拟处理时间
func (s *MockRPAServer) simulateProcessingTime(req MockRPARequest) time.Duration {
	// 根据请求类型和参数模拟不同的处理时间
	action, ok := req.Data["action"].(string)
	if !ok {
		return 2 * time.Second // 默认处理时间
	}

	switch action {
	case "process_document":
		return 3 * time.Second
	case "urgent_task":
		return 1 * time.Second
	case "background_task":
		return 5 * time.Second
	case "data_analysis":
		return 8 * time.Second
	case "simple_task":
		return 1 * time.Second
	case "silent_task":
		return 2 * time.Second
	case "critical_task":
		return 500 * time.Millisecond
	case "long_id_test":
		return 1 * time.Second
	case "special_chars_test":
		return 2 * time.Second
	case "process_purchase_request":
		return 2 * time.Second
	default:
		return 2 * time.Second
	}
}

// generateResponse 生成响应
func (s *MockRPAServer) generateResponse(req MockRPARequest, count int) MockRPAResponse {
	action, ok := req.Data["action"].(string)
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing action in request data",
		}
	}

	// 模拟不同的响应场景
	switch action {
	case "process_document":
		return s.generateDocumentProcessingResponse(req, count)
	case "urgent_task":
		return s.generateUrgentTaskResponse(req, count)
	case "background_task":
		return s.generateBackgroundTaskResponse(req, count)
	case "data_analysis":
		return s.generateDataAnalysisResponse(req, count)
	case "simple_task":
		return s.generateSimpleTaskResponse(req, count)
	case "silent_task":
		return s.generateSilentTaskResponse(req, count)
	case "critical_task":
		return s.generateCriticalTaskResponse(req, count)
	case "long_id_test":
		return s.generateLongIDTestResponse(req, count)
	case "special_chars_test":
		return s.generateSpecialCharsTestResponse(req, count)
	case "process_purchase_request":
		return s.generatePurchaseRequestResponse(req, count)
	default:
		return MockRPAResponse{
			ID:      req.ID,
			Success: true,
			Data: map[string]interface{}{
				"message":       "Unknown action processed successfully",
				"action":        action,
				"processed_at":  time.Now(),
				"request_count": count,
			},
		}
	}
}

// generateDocumentProcessingResponse 生成文档处理响应
func (s *MockRPAServer) generateDocumentProcessingResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in document processing request",
		}
	}

	filePath, _ := params["file_path"].(string)
	outputFormat, _ := params["output_format"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "process_document",
			"file_path":     filePath,
			"output_format": outputFormat,
			"result": map[string]interface{}{
				"pages_processed": 15,
				"text_extracted":  "Sample extracted text content...",
				"processing_time": "3.2s",
				"file_size":       "2.5MB",
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateUrgentTaskResponse 生成紧急任务响应
func (s *MockRPAServer) generateUrgentTaskResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in urgent task request",
		}
	}

	taskType, _ := params["task_type"].(string)
	timeout, _ := params["timeout"].(float64)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":    "urgent_task",
			"task_type": taskType,
			"timeout":   timeout,
			"result": map[string]interface{}{
				"status":        "completed_urgently",
				"response_time": "0.8s",
				"priority":      "high",
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateBackgroundTaskResponse 生成后台任务响应
func (s *MockRPAServer) generateBackgroundTaskResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in background task request",
		}
	}

	taskType, _ := params["task_type"].(string)
	schedule, _ := params["schedule"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":    "background_task",
			"task_type": taskType,
			"schedule":  schedule,
			"result": map[string]interface{}{
				"status":         "scheduled",
				"next_run":       "2024-01-02T02:00:00Z",
				"estimated_time": "5.2s",
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateDataAnalysisResponse 生成数据分析响应
func (s *MockRPAServer) generateDataAnalysisResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in data analysis request",
		}
	}

	analysisType, _ := params["analysis_type"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "data_analysis",
			"analysis_type": analysisType,
			"result": map[string]interface{}{
				"datasets_processed": 1,
				"records_analyzed":   15000,
				"trends_found":       5,
				"charts_generated":   []string{"line", "bar", "pie"},
				"summary": map[string]interface{}{
					"total_sales":     "$1,250,000",
					"growth_rate":     "12.5%",
					"top_category":    "electronics",
					"recommendations": []string{"Increase marketing", "Optimize inventory"},
				},
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateSimpleTaskResponse 生成简单任务响应
func (s *MockRPAServer) generateSimpleTaskResponse(req MockRPARequest, count int) MockRPAResponse {
	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "simple_task",
			"result":        "Task completed successfully",
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateSilentTaskResponse 生成静默任务响应
func (s *MockRPAServer) generateSilentTaskResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in silent task request",
		}
	}

	mode, _ := params["mode"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "silent_task",
			"mode":          mode,
			"result":        "Silent operation completed",
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateCriticalTaskResponse 生成关键任务响应
func (s *MockRPAServer) generateCriticalTaskResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in critical task request",
		}
	}

	criticalLevel, _ := params["critical_level"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":         "critical_task",
			"critical_level": criticalLevel,
			"result": map[string]interface{}{
				"status":        "critical_operation_completed",
				"response_time": "0.5s",
				"priority":      "maximum",
				"verified":      true,
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateLongIDTestResponse 生成长ID测试响应
func (s *MockRPAServer) generateLongIDTestResponse(req MockRPARequest, count int) MockRPAResponse {
	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "long_id_test",
			"id_length":     len(req.ID),
			"result":        "Long ID processed successfully",
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generateSpecialCharsTestResponse 生成特殊字符测试响应
func (s *MockRPAServer) generateSpecialCharsTestResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in special chars test request",
		}
	}

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":        "special_chars_test",
			"params":        params,
			"result":        "Special characters processed successfully",
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// generatePurchaseRequestResponse 生成采购请求处理响应
func (s *MockRPAServer) generatePurchaseRequestResponse(req MockRPARequest, count int) MockRPAResponse {
	params, ok := req.Data["params"].(map[string]interface{})
	if !ok {
		return MockRPAResponse{
			ID:      req.ID,
			Success: false,
			Error:   "Missing params in purchase request",
		}
	}

	// 提取采购请求的关键信息
	requestID, _ := params["request_id"].(string)
	material, _ := params["material"].(string)
	quantity, _ := params["quantity"].(float64)
	unitPrice, _ := params["unit_price"].(float64)
	requester, _ := params["requester"].(string)

	return MockRPAResponse{
		ID:      req.ID,
		Success: true,
		Data: map[string]interface{}{
			"action":     "process_purchase_request",
			"request_id": requestID,
			"result": map[string]interface{}{
				"status":           "processed",
				"material":         material,
				"quantity":         quantity,
				"unit_price":       unitPrice,
				"total_amount":     quantity * unitPrice,
				"requester":        requester,
				"processing_time":  "2.1s",
				"approval_status":  "pending",
				"vendor_contacted": true,
			},
			"processed_at":  time.Now(),
			"request_count": count,
		},
	}
}

// setupRoutes 设置路由
func (s *MockRPAServer) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	// API 路由
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/execute", s.executeHandler).Methods("POST")
	api.HandleFunc("/status", s.statusHandler).Methods("GET")

	// 健康检查
	r.HandleFunc("/health", s.healthHandler).Methods("GET")

	// 根路径
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service":   "Mock RPA System",
			"version":   "1.0.0",
			"status":    "running",
			"uptime":    time.Since(s.startTime).String(),
			"endpoints": []string{"/api/v1/execute", "/api/v1/status", "/health"},
		})
	})

	return r
}

func main() {
	// 配置日志
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	// 创建模拟RPA服务器
	server := NewMockRPAServer()
	server.logger = logger

	// 设置路由
	router := server.setupRoutes()

	// 启动服务器
	port := 8084 // 使用不同的端口避免冲突
	addr := fmt.Sprintf(":%d", port)

	logger.WithFields(logrus.Fields{
		"port":    port,
		"address": "http://localhost" + addr,
	}).Info("Starting Mock RPA Server")

	logger.Info("Mock RPA Server endpoints:")
	logger.Info("  - POST /api/v1/execute  - Execute RPA task")
	logger.Info("  - GET  /api/v1/status   - Check RPA status")
	logger.Info("  - GET  /health          - Health check")
	logger.Info("  - GET  /                - Service info")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start Mock RPA Server:", err)
	}
}
