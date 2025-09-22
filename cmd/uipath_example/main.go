package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	serverURL := "http://localhost:8082"

	fmt.Println("=== UiPath RPA 集成示例 ===")

	// 1. 检查UiPath状态
	fmt.Println("\n1. 检查UiPath状态...")
	checkStatus(serverURL)

	// 2. 测试连接
	fmt.Println("\n2. 测试UiPath连接...")
	testConnection(serverURL)

	// 3. 添加测试队列项目
	fmt.Println("\n3. 添加测试队列项目...")
	addTestItem(serverURL)

	// 4. 添加自定义队列项目
	fmt.Println("\n4. 添加自定义队列项目...")
	addCustomItem(serverURL)

	fmt.Println("\n=== 示例完成 ===")
}

// checkStatus 检查UiPath状态
func checkStatus(serverURL string) {
	resp, err := http.Get(serverURL + "/api/uipath/status")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(body))
}

// testConnection 测试连接
func testConnection(serverURL string) {
	resp, err := http.Get(serverURL + "/api/uipath/test")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(body))
}

// addTestItem 添加测试队列项目
func addTestItem(serverURL string) {
	resp, err := http.Post(serverURL+"/api/uipath/queue/test", "application/json", nil)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(body))
}

// addCustomItem 添加自定义队列项目
func addCustomItem(serverURL string) {
	// 创建测试数据
	testData := map[string]interface{}{
		"Documento": "EXAMPLE-" + fmt.Sprintf("%d", time.Now().Unix()),
		"ClienteId": 789,
		"Total":     299.99,
		"Moeda":     "CNY",
		"TestMode":  true,
		"Timestamp": time.Now().Format(time.RFC3339),
	}

	req := map[string]interface{}{
		"queue_name":       "CreationPR",
		"priority":         "High",
		"specific_content": testData,
		"reference":        "Example-" + fmt.Sprintf("%d", time.Now().Unix()),
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	resp, err := http.Post(serverURL+"/api/uipath/queue/add", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(body))
}
