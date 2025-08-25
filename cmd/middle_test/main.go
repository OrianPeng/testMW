package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// 测试采购请求API
func main() {
	baseURL := "http://localhost:8082/api/v1"
	timestamp := time.Now().Unix()

	fmt.Println("=== 重构后采购请求API测试 ===")

	// 测试0: 清理队列
	fmt.Println("\n=== 测试0: 清理队列 ===")
	clearResp, err := makeTestRequest("POST", baseURL+"/purchase-requests/queue/clear-all", nil)
	if err != nil {
		fmt.Printf("清理队列失败: %v\n", err)
	} else {
		fmt.Printf("清理队列响应: %s\n", clearResp)
	}

	// 等待一秒确保清理完成
	time.Sleep(1 * time.Second)

	// 验证队列是否已清空
	fmt.Println("\n=== 验证队列是否已清空 ===")
	verifyResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/queue", nil)
	if err != nil {
		fmt.Printf("验证队列失败: %v\n", err)
	} else {
		fmt.Printf("验证队列响应: %s\n", verifyResp)
	}

	// 测试1: 提交采购请求到队列（完整信息）
	fmt.Println("\n=== 测试1: 提交采购请求到队列（完整信息） ===")
	queueRequest := map[string]interface{}{
		"request_id":            "PR-" + strconv.FormatInt(timestamp, 10) + "-001",
		"doc_type":              "PR",
		"plant":                 "1000",
		"quantity":              2,
		"unit_price":            5999.00,
		"material":              "LAPTOP-001",
		"delivery_date":         time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
		"vendor_code":           "VENDOR-001",
		"short_text":            "高性能开发用笔记本电脑",
		"material_group":        "IT-EQUIPMENT",
		"unit_type":             "PC",
		"requester":             "张三",
		"purchase_organization": "PURCHASE-ORG-001",
		"currency":              "CNY",
		"priority":              3,
		"urgency":               "normal",
		"comments":              "请尽快处理",
	}

	queueResp, err := makeTestRequest("POST", baseURL+"/purchase-requests/queue", queueRequest)
	if err != nil {
		fmt.Printf("提交采购请求到队列失败: %v\n", err)
	} else {
		fmt.Printf("提交到队列响应: %s\n", queueResp)
	}

	// 测试2: 提交另一个采购请求到队列（完整信息）
	fmt.Println("\n=== 测试2: 提交另一个采购请求到队列（完整信息） ===")
	queueRequest2 := map[string]interface{}{
		"request_id":            "PR-" + strconv.FormatInt(timestamp, 10) + "-002",
		"doc_type":              "PR",
		"plant":                 "2000",
		"quantity":              1,
		"unit_price":            2999.00,
		"material":              "PROJECTOR-001",
		"delivery_date":         time.Now().AddDate(0, 0, 3).Format(time.RFC3339),
		"vendor_code":           "VENDOR-002",
		"short_text":            "会议室用投影仪",
		"material_group":        "OFFICE-EQUIPMENT",
		"unit_type":             "EA",
		"requester":             "李四",
		"purchase_organization": "PURCHASE-ORG-002",
		"currency":              "CNY",
		"priority":              4, // 更高优先级
		"urgency":               "urgent",
		"comments":              "紧急需求",
	}

	queueResp2, err := makeTestRequest("POST", baseURL+"/purchase-requests/queue", queueRequest2)
	if err != nil {
		fmt.Printf("提交第二个采购请求到队列失败: %v\n", err)
	} else {
		fmt.Printf("提交第二个到队列响应: %s\n", queueResp2)
	}

	// 测试3: 查询队列状态
	fmt.Println("\n=== 测试3: 查询队列状态 ===")
	statusResp, err := makeTestRequest("GET", baseURL+"/status", nil)
	if err != nil {
		fmt.Printf("查询队列状态失败: %v\n", err)
	} else {
		fmt.Printf("队列状态响应: %s\n", statusResp)
	}

	// 测试4: 查询队列中的采购请求列表
	fmt.Println("\n=== 测试4: 查询队列中的采购请求列表 ===")
	queueListResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/queue", nil)
	if err != nil {
		fmt.Printf("查询队列列表失败: %v\n", err)
	} else {
		fmt.Printf("队列列表响应: %s\n", queueListResp)
	}

	// 测试5: 查询特定采购请求的队列状态
	fmt.Println("\n=== 测试5: 查询特定采购请求的队列状态 ===")
	requestID := "PR-" + strconv.FormatInt(timestamp, 10) + "-001"
	queueStatusResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/queue/"+requestID, nil)
	if err != nil {
		fmt.Printf("查询队列状态失败: %v\n", err)
	} else {
		fmt.Printf("队列状态响应: %s\n", queueStatusResp)
	}

	// 测试6: 创建采购请求到数据库（完整信息）
	fmt.Println("\n=== 测试6: 创建采购请求到数据库（完整信息） ===")
	dbRequest := map[string]interface{}{
		"request_id":            "PR-DB-" + strconv.FormatInt(timestamp, 10) + "-001",
		"doc_type":              "PR",
		"plant":                 "3000",
		"quantity":              1,
		"unit_price":            1500.00,
		"material":              "PRINTER-001",
		"delivery_date":         time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
		"vendor_code":           "VENDOR-003",
		"short_text":            "办公室用打印机",
		"material_group":        "OFFICE-EQUIPMENT",
		"unit_type":             "EA",
		"requester":             "王五",
		"purchase_organization": "PURCHASE-ORG-003",
		"currency":              "CNY",
		"priority":              2,
		"urgency":               "normal",
		"comments":              "常规采购",
	}
	dbResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", dbRequest)
	if err != nil {
		fmt.Printf("创建采购请求到数据库失败: %v\n", err)
	} else {
		fmt.Printf("创建到数据库响应: %s\n", dbResp)
	}

	// 自动提取最新创建的采购请求request_id
	var createdRequestID string
	if dbResp != "" {
		var dbRespObj map[string]interface{}
		if err := json.Unmarshal([]byte(dbResp), &dbRespObj); err == nil {
			if data, ok := dbRespObj["data"].(map[string]interface{}); ok {
				if requestIDVal, ok := data["request_id"].(string); ok {
					createdRequestID = requestIDVal
				}
			}
		}
	}
	if createdRequestID == "" {
		fmt.Println("未能获取到新建采购请求request_id，update测试将跳过！")
		return
	}

	// ========== Update API 详细测试 ========== //
	fmt.Printf("\n=== Update API测试将使用request_id: %s ===\n", createdRequestID)

	fmt.Println("\n=== 测试11: 使用 request_id 更新 ===")
	updateReq1 := map[string]interface{}{
		"material": "MATERIAL-UPDATED",
	}
	updateResp1, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq1)
	if err != nil {
		fmt.Printf("request_id更新失败: %v\n", err)
	} else {
		fmt.Printf("request_id更新响应: %s\n", updateResp1)
	}

	fmt.Println("\n=== 测试12: 测试状态检查 - 完整信息（应设置为enough） ===")
	updateReq2 := map[string]interface{}{
		"material":              "COMPLETE-MATERIAL",
		"quantity":              10,
		"unit_price":            100.00,
		"currency":              "CNY",
		"short_text":            "Complete information test",
		"material_group":        "TEST-GROUP",
		"unit_type":             "EA",
		"requester":             "Test User",
		"purchase_organization": "TEST-ORG",
	}
	updateResp2, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq2)
	if err != nil {
		fmt.Printf("完整信息更新失败: %v\n", err)
	} else {
		fmt.Printf("完整信息更新响应: %s\n", updateResp2)
	}

	fmt.Println("\n=== 测试13: 更新处理信息 ===")
	updateReq3 := map[string]interface{}{
		"retry_count":  2,
		"error_msg":    "Test error message",
		"processed_at": "2026-05-19T15:30:00",
	}
	updateResp3, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq3)
	if err != nil {
		fmt.Printf("处理信息更新失败: %v\n", err)
	} else {
		fmt.Printf("处理信息更新响应: %s\n", updateResp3)
	}

	fmt.Println("\n=== 测试14: 批量字段更新 ===")
	updateReq4 := map[string]interface{}{
		"material":       "BATCH-UPDATED-MATERIAL",
		"quantity":       150,
		"unit_price":     25.50,
		"priority":       3,
		"urgency":        "urgent",
		"comments":       "Batch update test",
		"vendor_code":    "VENDOR-BATCH",
		"material_group": "GROUP-BATCH",
	}
	updateResp4, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq4)
	if err != nil {
		fmt.Printf("批量更新失败: %v\n", err)
	} else {
		fmt.Printf("批量更新响应: %s\n", updateResp4)
	}

	fmt.Println("\n=== 测试15: 测试状态检查 - 不完整信息（应设置为lack） ===")
	updateReq5 := map[string]interface{}{
		"material":   "INCOMPLETE-MATERIAL",
		"quantity":   5,
		"unit_price": 50.00,
		// 故意不设置currency，使其不完整
		"short_text":            "Incomplete information test",
		"material_group":        "TEST-GROUP",
		"unit_type":             "EA",
		"requester":             "Test User",
		"purchase_organization": "TEST-ORG",
	}
	updateResp5, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq5)
	if err != nil {
		fmt.Printf("不完整信息更新失败: %v\n", err)
	} else {
		fmt.Printf("不完整信息更新响应: %s\n", updateResp5)
	}

	fmt.Println("\n=== 测试16: 更新不存在的request_id（应失败） ===")
	updateReq6 := map[string]interface{}{
		"material": "SHOULD-NOT-EXIST",
	}
	updateResp6, err := makeTestRequest("PUT", baseURL+"/purchase-requests/NON-EXISTENT-REQUEST-ID", updateReq6)
	if err != nil {
		fmt.Printf("不存在request_id更新失败: %v\n", err)
	} else {
		fmt.Printf("不存在request_id更新响应: %s\n", updateResp6)
	}

	// 测试17: 创建不完整的采购请求（应设置为lack）
	fmt.Println("\n=== 测试17: 创建不完整的采购请求（应设置为lack） ===")
	incompleteRequestID := "PR-INCOMPLETE-" + strconv.FormatInt(timestamp, 10)
	incompleteRequest := map[string]interface{}{
		"request_id":    incompleteRequestID,
		"doc_type":      "PR",
		"plant":         "4000",
		"quantity":      1,
		"unit_price":    500.00,
		"material":      "INCOMPLETE-ITEM",
		"delivery_date": time.Now().AddDate(0, 0, 5).Format(time.RFC3339),
		"vendor_code":   "VENDOR-004",
		"short_text":    "不完整的采购请求",
		// 故意不设置material_group, unit_type, requester, purchase_organization, currency
		"priority": 1,
		"urgency":  "normal",
		"comments": "测试不完整信息",
	}
	incompleteResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", incompleteRequest)
	if err != nil {
		fmt.Printf("创建不完整采购请求失败: %v\n", err)
	} else {
		fmt.Printf("创建不完整采购请求响应: %s\n", incompleteResp)
	}

	// 测试18: 更新采购请求为已批准状态
	fmt.Println("\n=== 测试18: 更新采购请求为已批准状态 ===")
	approveRequest := map[string]interface{}{
		"status":        "approved",
		"approver_id":   "approver_002",
		"approver_name": "Senior Approver",
		"comments":      "Approved after review",
	}

	approveResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, approveRequest)
	if err != nil {
		fmt.Printf("更新为已批准状态失败: %v\n", err)
	} else {
		fmt.Printf("更新为已批准状态响应: %s\n", approveResp)
	}

	// 测试19: 错误测试 - 更新不存在的采购请求
	fmt.Println("\n=== 测试19: 错误测试 - 更新不存在的采购请求 ===")
	updateNonExistentRequest := map[string]interface{}{
		"comments": "This should fail",
	}

	updateNonExistentResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/NON-EXISTENT-REQUEST-ID", updateNonExistentRequest)
	if err != nil {
		fmt.Printf("更新不存在的采购请求失败（预期）: %v\n", err)
	} else {
		fmt.Printf("更新不存在的采购请求响应: %s\n", updateNonExistentResp)
	}

	// 测试20: 错误测试 - 无效的JSON数据
	fmt.Println("\n=== 测试20: 错误测试 - 无效的JSON数据 ===")
	invalidJSONRequest := map[string]interface{}{
		"quantity": "invalid_number", // 应该是数字
	}

	invalidJSONResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, invalidJSONRequest)
	if err != nil {
		fmt.Printf("无效JSON数据测试失败（预期）: %v\n", err)
	} else {
		fmt.Printf("无效JSON数据测试响应: %s\n", invalidJSONResp)
	}

	// 测试21: 验证更新后的采购请求
	fmt.Println("\n=== 测试21: 验证更新后的采购请求 ===")
	verifyUpdatedResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/"+createdRequestID, nil)
	if err != nil {
		fmt.Printf("验证更新后的采购请求失败: %v\n", err)
	} else {
		fmt.Printf("验证更新后的采购请求响应: %s\n", verifyUpdatedResp)
	}

	// 用户自定义测试：部分字段为空的采购请求
	fmt.Println("\n=== 用户自定义测试: 部分字段为空的采购请求 ===")
	userRequest := map[string]interface{}{
		"request_id":            "PR-DB-1703123456-001",
		"doc_type":              "PR",
		"plant":                 "",
		"quantity":              6,
		"unit_price":            7000.00,
		"material":              "",
		"delivery_date":         "",
		"vendor_code":           "APPLE-001",
		"short_text":            "",
		"material_group":        "",
		"unit_type":             "",
		"requester":             "",
		"purchase_organization": "",
		"currency":              "",
		"priority":              2,
		"urgency":               "normal",
		"comments":              "",
	}
	userResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", userRequest)
	if err != nil {
		fmt.Printf("用户自定义采购请求测试失败: %v\n", err)
	} else {
		fmt.Printf("用户自定义采购请求响应: %s\n", userResp)
	}

	// 测试22: 测试新的时间格式支持
	fmt.Println("\n=== 测试22: 测试新的时间格式支持 ===")

	// 测试不同的时间格式
	timeFormats := []string{
		"2026-05-19",                // 日期格式
		"2026-05-19T15:30:00",       // 日期时间格式
		"2026-05-19 15:30:00",       // 空格分隔格式
		"2026/05/19",                // 斜杠分隔格式
		"2026-05-19T15:30:00Z",      // UTC格式
		"2026-05-19T15:30:00+08:00", // 带时区格式
	}

	for i, timeStr := range timeFormats {
		fmt.Printf("\n--- 测试时间格式 %d: %s ---\n", i+1, timeStr)
		updateReq := map[string]interface{}{
			"delivery_date": timeStr,
			"comments":      fmt.Sprintf("Testing time format: %s", timeStr),
		}
		updateResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, updateReq)
		if err != nil {
			fmt.Printf("时间格式 '%s' 测试失败: %v\n", timeStr, err)
		} else {
			fmt.Printf("时间格式 '%s' 测试成功: %s\n", timeStr, updateResp)
		}
	}

	// 测试23: 测试空字符串和null值
	fmt.Println("\n=== 测试23: 测试空字符串和null值 ===")
	emptyUpdateReq := map[string]interface{}{
		"delivery_date": "",
		"processed_at":  nil,
		"comments":      "Testing empty and null values",
	}
	emptyUpdateResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+createdRequestID, emptyUpdateReq)
	if err != nil {
		fmt.Printf("空值测试失败: %v\n", err)
	} else {
		fmt.Printf("空值测试成功: %s\n", emptyUpdateResp)
	}

	// 测试24: 专门测试状态从lack到enough的更新
	fmt.Println("\n=== 测试24: 专门测试状态从lack到enough的更新 ===")

	// 首先创建一个不完整的采购请求（应该状态为lack）
	statusTestRequestID := "PR-STATUS-TEST-" + strconv.FormatInt(timestamp, 10)
	statusTestRequest := map[string]interface{}{
		"request_id":            statusTestRequestID,
		"doc_type":              "PR",
		"plant":                 "5000",
		"quantity":              5,
		"unit_price":            200.00,
		"material":              "STATUS-TEST-MATERIAL",
		"delivery_date":         time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
		"vendor_code":           "VENDOR-STATUS",
		"short_text":            "状态测试物料",
		"material_group":        "GROUP-STATUS",
		"unit_type":             "EA",
		"requester":             "STATUS-TEST-USER",
		"purchase_organization": "ORG-STATUS",
		// 故意不设置currency，使其状态为lack
		"priority": 2,
		"urgency":  "normal",
		"comments": "测试状态更新逻辑",
	}

	statusTestResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", statusTestRequest)
	if err != nil {
		fmt.Printf("创建状态测试采购请求失败: %v\n", err)
	} else {
		fmt.Printf("创建状态测试采购请求响应: %s\n", statusTestResp)
	}

	// 等待一秒确保数据写入
	time.Sleep(1 * time.Second)

	// 查询创建的不完整请求，确认状态为lack
	fmt.Println("\n--- 查询不完整请求状态 ---")
	queryIncompleteResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/"+statusTestRequestID, nil)
	if err != nil {
		fmt.Printf("查询不完整请求失败: %v\n", err)
	} else {
		fmt.Printf("不完整请求状态: %s\n", queryIncompleteResp)
	}

	// 更新请求，添加缺失的currency字段
	fmt.Println("\n--- 更新请求，添加currency字段 ---")
	statusUpdateReq := map[string]interface{}{
		"currency": "CNY",
		"comments": "添加了货币字段，状态应该变为enough",
	}

	statusUpdateResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+statusTestRequestID, statusUpdateReq)
	if err != nil {
		fmt.Printf("状态更新失败: %v\n", err)
	} else {
		fmt.Printf("状态更新响应: %s\n", statusUpdateResp)
	}

	// 等待一秒确保数据更新
	time.Sleep(1 * time.Second)

	// 再次查询，确认状态已更新为enough
	fmt.Println("\n--- 查询更新后的状态 ---")
	queryUpdatedResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/"+statusTestRequestID, nil)
	if err != nil {
		fmt.Printf("查询更新后状态失败: %v\n", err)
	} else {
		fmt.Printf("更新后状态: %s\n", queryUpdatedResp)
	}

	// 测试25: 测试多个字段同时更新时的状态检查
	fmt.Println("\n=== 测试25: 测试多个字段同时更新时的状态检查 ===")

	// 创建另一个不完整的请求
	multiFieldTestID := "PR-MULTI-TEST-" + strconv.FormatInt(timestamp, 10)
	multiFieldRequest := map[string]interface{}{
		"request_id":    multiFieldTestID,
		"doc_type":      "PR",
		"plant":         "6000",
		"quantity":      3,
		"unit_price":    150.00,
		"material":      "MULTI-TEST-MATERIAL",
		"delivery_date": time.Now().AddDate(0, 0, 15).Format(time.RFC3339),
		"vendor_code":   "VENDOR-MULTI",
		"short_text":    "多字段测试",
		// 故意不设置多个字段
		"priority": 1,
		"urgency":  "normal",
		"comments": "测试多字段更新",
	}

	multiFieldResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", multiFieldRequest)
	if err != nil {
		fmt.Printf("创建多字段测试请求失败: %v\n", err)
	} else {
		fmt.Printf("创建多字段测试请求响应: %s\n", multiFieldResp)
	}

	time.Sleep(1 * time.Second)

	// 同时更新多个缺失字段
	fmt.Println("\n--- 同时更新多个缺失字段 ---")
	multiFieldUpdateReq := map[string]interface{}{
		"material_group":        "GROUP-MULTI",
		"unit_type":             "EA",
		"requester":             "MULTI-TEST-USER",
		"purchase_organization": "ORG-MULTI",
		"currency":              "USD",
		"comments":              "同时更新了多个字段，状态应该变为enough",
	}

	multiFieldUpdateResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+multiFieldTestID, multiFieldUpdateReq)
	if err != nil {
		fmt.Printf("多字段更新失败: %v\n", err)
	} else {
		fmt.Printf("多字段更新响应: %s\n", multiFieldUpdateResp)
	}

	time.Sleep(1 * time.Second)

	// 查询多字段更新后的状态
	fmt.Println("\n--- 查询多字段更新后的状态 ---")
	queryMultiFieldResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/"+multiFieldTestID, nil)
	if err != nil {
		fmt.Printf("查询多字段更新后状态失败: %v\n", err)
	} else {
		fmt.Printf("多字段更新后状态: %s\n", queryMultiFieldResp)
	}

	// 测试26: 测试数量/单价更新时的总金额重新计算和状态检查
	fmt.Println("\n=== 测试26: 测试数量/单价更新时的总金额重新计算和状态检查 ===")

	// 创建一个完整的请求
	calcTestID := "PR-CALC-TEST-" + strconv.FormatInt(timestamp, 10)
	calcRequest := map[string]interface{}{
		"request_id":            calcTestID,
		"doc_type":              "PR",
		"plant":                 "7000",
		"quantity":              2,
		"unit_price":            100.00,
		"material":              "CALC-TEST-MATERIAL",
		"delivery_date":         time.Now().AddDate(0, 0, 20).Format(time.RFC3339),
		"vendor_code":           "VENDOR-CALC",
		"short_text":            "计算测试物料",
		"material_group":        "GROUP-CALC",
		"unit_type":             "EA",
		"requester":             "CALC-TEST-USER",
		"purchase_organization": "ORG-CALC",
		"currency":              "EUR",
		"priority":              3,
		"urgency":               "urgent",
		"comments":              "测试总金额计算和状态检查",
	}

	calcResp, err := makeTestRequest("POST", baseURL+"/purchase-requests", calcRequest)
	if err != nil {
		fmt.Printf("创建计算测试请求失败: %v\n", err)
	} else {
		fmt.Printf("创建计算测试请求响应: %s\n", calcResp)
	}

	time.Sleep(1 * time.Second)

	// 更新数量和单价
	fmt.Println("\n--- 更新数量和单价 ---")
	calcUpdateReq := map[string]interface{}{
		"quantity":   5,
		"unit_price": 150.00,
		"comments":   "更新了数量和单价，总金额应该重新计算为750.00",
	}

	calcUpdateResp, err := makeTestRequest("PUT", baseURL+"/purchase-requests/"+calcTestID, calcUpdateReq)
	if err != nil {
		fmt.Printf("计算更新失败: %v\n", err)
	} else {
		fmt.Printf("计算更新响应: %s\n", calcUpdateResp)
	}

	time.Sleep(1 * time.Second)

	// 查询计算更新后的结果
	fmt.Println("\n--- 查询计算更新后的结果 ---")
	queryCalcResp, err := makeTestRequest("GET", baseURL+"/purchase-requests/"+calcTestID, nil)
	if err != nil {
		fmt.Printf("查询计算更新后结果失败: %v\n", err)
	} else {
		fmt.Printf("计算更新后结果: %s\n", queryCalcResp)
	}

	fmt.Println("\n=== 重构后测试完成（包含更新接口测试、时间格式优化和状态更新测试） ===")
}

// makeTestRequest 发送HTTP请求
func makeTestRequest(method, url string, body interface{}) (string, error) {
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
