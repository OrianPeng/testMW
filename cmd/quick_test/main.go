package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// QuickTest 快速测试
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

	logger.WithField("server_url", serverURL).Info("Starting quick API tests")

	// 测试计数器
	totalTests := 0
	passedTests := 0
	failedTests := 0

	// 测试列表
	tests := []struct {
		name     string
		method   string
		endpoint string
		data     interface{}
	}{
		{"Health Check", "GET", "/api/v1/health", nil},
		{"System Status", "GET", "/api/v1/status", nil},
		{"List Purchase Requests", "GET", "/api/purchase-requests", nil},
		{"List Purchase Orders", "GET", "/api/purchase-orders", nil},
		{"List Suppliers", "GET", "/api/suppliers", nil},
		{"Process Overview", "GET", "/api/process-monitor/overview", nil},
		{"Dashboard Metrics", "GET", "/api/dashboard/metrics", nil},
		{"UiPath Status", "GET", "/api/uipath/status", nil},
		{"UiPath Available", "GET", "/api/uipath/available", nil},
		{"Queue Stats", "GET", "/api/v1/status", nil},
	}

	fmt.Println("=== 快速API测试 ===")
	fmt.Printf("服务器地址: %s\n", serverURL)
	fmt.Println(strings.Repeat("-", 60))

	// 运行测试
	for _, test := range tests {
		totalTests++
		startTime := time.Now()

		var req *http.Request
		var err error

		if test.data != nil {
			jsonData, _ := json.Marshal(test.data)
			req, err = http.NewRequest(test.method, serverURL+test.endpoint, bytes.NewBuffer(jsonData))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
			}
		} else {
			req, err = http.NewRequest(test.method, serverURL+test.endpoint, nil)
		}

		if err != nil {
			fmt.Printf("❌ %-30s %-6s %-30s [ERROR: %v]\n",
				test.name, test.method, test.endpoint, err)
			failedTests++
			continue
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ %-30s %-6s %-30s [ERROR: %v] (%v)\n",
				test.name, test.method, test.endpoint, err, duration)
			failedTests++
		} else {
			defer resp.Body.Close()

			status := "✓"
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				passedTests++
			} else {
				status = "✗"
				failedTests++
			}

			fmt.Printf("%s %-30s %-6s %-30s [%d] (%v)\n",
				status, test.name, test.method, test.endpoint, resp.StatusCode, duration)
		}
	}

	// 输出结果
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("总测试数: %d, 通过: %d, 失败: %d\n", totalTests, passedTests, failedTests)

	if failedTests == 0 {
		fmt.Println("✅ 所有测试都通过了！")
	} else {
		fmt.Printf("❌ 发现 %d 个失败的测试\n", failedTests)
	}
}
