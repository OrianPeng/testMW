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
	baseURL := "http://localhost:8080/api/v1/purchase-requests"

	// 创建测试采购请求
	createRequest := map[string]interface{}{
		"request_id":            "PR-TEST-" + fmt.Sprintf("%d", time.Now().Unix()),
		"doc_type":              "PR",
		"plant":                 "1000",
		"quantity":              10,
		"unit_price":            25.50,
		"material":              "TEST-MATERIAL",
		"delivery_date":         time.Now().AddDate(0, 0, 5).Format(time.RFC3339),
		"vendor_code":           "VENDOR-TEST",
		"short_text":            "Test purchase request",
		"material_group":        "TEST-GROUP",
		"unit_type":             "PCS",
		"requester":             "Test User",
		"purchase_organization": "TEST-ORG",
		"currency":              "CNY",
		"priority":              2,
		"urgency":               "normal",
		"comments":              "Test request for request_id update",
	}

	fmt.Println("=== 创建测试采购请求 ===")
	createResp, err := makeRequest("POST", baseURL, createRequest)
	if err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}
	fmt.Printf("创建响应: %s\n", createResp)

	// 提取request_id
	var createRespObj map[string]interface{}
	if err := json.Unmarshal([]byte(createResp), &createRespObj); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}

	var requestID string
	if data, ok := createRespObj["data"].(map[string]interface{}); ok {
		if id, ok := data["request_id"].(string); ok {
			requestID = id
		}
	}

	if requestID == "" {
		fmt.Println("未能获取request_id，测试终止")
		return
	}

	fmt.Printf("\n=== 使用request_id更新: %s ===\n", requestID)

	// 测试更新
	updateRequest := map[string]interface{}{
		"quantity":   20,
		"unit_price": 30.00,
		"comments":   "Updated via request_id",
		"priority":   3,
		"urgency":    "urgent",
	}

	updateResp, err := makeRequest("PUT", baseURL+"/"+requestID, updateRequest)
	if err != nil {
		fmt.Printf("更新失败: %v\n", err)
		return
	}
	fmt.Printf("更新响应: %s\n", updateResp)

	// 验证更新结果
	fmt.Printf("\n=== 验证更新结果 ===\n")
	verifyResp, err := makeRequest("GET", baseURL+"/"+requestID, nil)
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
		return
	}
	fmt.Printf("验证响应: %s\n", verifyResp)

	fmt.Println("\n=== 测试完成 ===")
}

func makeRequest(method, url string, body interface{}) (string, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(respBody), nil
}
