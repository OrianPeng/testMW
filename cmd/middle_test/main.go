package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// TestResult 测试结果
type TestResult struct {
	TestName     string        `json:"test_name"`
	Endpoint     string        `json:"endpoint"`
	Method       string        `json:"method"`
	Status       string        `json:"status"` // PASS, FAIL, SKIP
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Response     string        `json:"response,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// TestSuite 测试套件
type TestSuite struct {
	ServerURL    string        `json:"server_url"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	TotalTests   int           `json:"total_tests"`
	PassedTests  int           `json:"passed_tests"`
	FailedTests  int           `json:"failed_tests"`
	SkippedTests int           `json:"skipped_tests"`
	Results      []TestResult  `json:"results"`
	Duration     time.Duration `json:"duration"`
}

// TestData 测试数据
type TestData struct {
	PurchaseRequestID string
	PurchaseOrderID   string
	SupplierID        string
	FileID            string
}

var testData TestData

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 获取服务器地址
	serverURL := "http://localhost:8082"
	if len(os.Args) > 1 {
		serverURL = os.Args[1]
	}

	logger.WithField("server_url", serverURL).Info("Starting comprehensive API tests")

	// 创建测试套件
	suite := &TestSuite{
		ServerURL: serverURL,
		StartTime: time.Now(),
		Results:   make([]TestResult, 0),
	}

	// 运行所有测试
	runAllTests(suite, logger)

	// 生成测试报告
	generateReport(suite, logger)
}
