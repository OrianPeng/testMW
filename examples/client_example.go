package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// 示例请求结构
type ExampleRequest struct {
	ID       string                 `json:"id"`
	AgentID  string                 `json:"agent_id"`
	Data     map[string]interface{} `json:"data"`
	Priority int                    `json:"priority,omitempty"`
	Callback string                 `json:"callback,omitempty"`
}

// 示例响应结构
type ExampleResponse struct {
	RequestID string      `json:"request_id"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func main() {
	// 中间件服务地址
	middlewareURL := "http://localhost:8080"

	// 示例1: 提交请求
	fmt.Println("=== 示例1: 提交请求 ===")
	requestID, err := submitRequest(middlewareURL)
	if err != nil {
		log.Fatal("提交请求失败:", err)
	}
	fmt.Printf("请求已提交，ID: %s\n", requestID)

	// 示例2: 查询请求状态
	fmt.Println("\n=== 示例2: 查询请求状态 ===")
	status, err := getRequestStatus(middlewareURL, requestID)
	if err != nil {
		log.Fatal("查询状态失败:", err)
	}
	fmt.Printf("请求状态: %s - %s\n", status.Status, status.Message)

	// 示例3: 列出所有请求
	fmt.Println("\n=== 示例3: 列出所有请求 ===")
	requests, err := listRequests(middlewareURL)
	if err != nil {
		log.Fatal("列出请求失败:", err)
	}
	fmt.Printf("找到 %d 个请求\n", len(requests))
	for _, req := range requests {
		fmt.Printf("- %s: %s\n", req.RequestID, req.Status)
	}

	// 示例4: 获取系统状态
	fmt.Println("\n=== 示例4: 获取系统状态 ===")
	err = getSystemStatus(middlewareURL)
	if err != nil {
		log.Fatal("获取系统状态失败:", err)
	}

	// 示例5: 健康检查
	fmt.Println("\n=== 示例5: 健康检查 ===")
	err = healthCheck(middlewareURL)
	if err != nil {
		log.Fatal("健康检查失败:", err)
	}
}

// submitRequest 提交请求
func submitRequest(baseURL string) (string, error) {
	req := ExampleRequest{
		ID:      fmt.Sprintf("req-%d", time.Now().Unix()),
		AgentID: "example-agent",
		Data: map[string]interface{}{
			"action": "process_document",
			"document_id": "doc-123",
			"params": map[string]interface{}{
				"format": "pdf",
				"pages":  []int{1, 2, 3},
			},
		},
		Priority: 5,
		Callback: "http://example-agent/callback",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := http.Post(baseURL+"/api/v1/requests", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var response ExampleResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return response.RequestID, nil
}

// getRequestStatus 获取请求状态
func getRequestStatus(baseURL, requestID string) (*ExampleResponse, error) {
	resp, err := http.Get(baseURL + "/api/v1/requests/" + requestID)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var response ExampleResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &response, nil
}

// listRequests 列出请求
func listRequests(baseURL string) ([]ExampleResponse, error) {
	resp, err := http.Get(baseURL + "/api/v1/requests?limit=10")
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var response struct {
		Requests []ExampleResponse `json:"requests"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Requests, nil
}

// getSystemStatus 获取系统状态
func getSystemStatus(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/v1/status")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("系统状态: %s\n", string(respBody))
	return nil
}

// healthCheck 健康检查
func healthCheck(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/v1/health")
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("健康检查失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("健康状态: %s\n", string(respBody))
	return nil
}