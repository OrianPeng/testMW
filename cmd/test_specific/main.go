package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 测试采购请求列表API
	resp, err := http.Get("http://localhost:8082/api/purchase-requests")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	// 测试采购订单列表API
	resp2, err := http.Get("http://localhost:8082/api/purchase-orders")
	if err != nil {
		fmt.Printf("Error making PO request: %v\n", err)
		return
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("Error reading PO response: %v\n", err)
		return
	}

	fmt.Printf("PO Status Code: %d\n", resp2.StatusCode)
	fmt.Printf("PO Response Body: %s\n", string(body2))
}
