package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// runTest 运行单个测试
func runTest(suite *TestSuite, testName, method, endpoint string, data interface{}, logger *logrus.Logger) TestResult {
	startTime := time.Now()

	var req *http.Request
	var err error

	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, suite.ServerURL+endpoint, bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(method, suite.ServerURL+endpoint, nil)
	}

	result := TestResult{
		TestName:     testName,
		Endpoint:     endpoint,
		Method:       method,
		ResponseTime: time.Since(startTime),
		Timestamp:    time.Now(),
	}

	if err != nil {
		result.Status = "FAIL"
		result.ErrorMessage = err.Error()
		suite.Results = append(suite.Results, result)
		logger.WithFields(logrus.Fields{
			"test":   testName,
			"status": "FAIL",
			"error":  err.Error(),
		}).Error("Test failed")
		return result
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	duration := time.Since(startTime)
	result.ResponseTime = duration

	if err != nil {
		result.Status = "FAIL"
		result.ErrorMessage = err.Error()
	} else {
		defer resp.Body.Close()
		result.StatusCode = resp.StatusCode

		body, _ := io.ReadAll(resp.Body)
		result.Response = string(body)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Status = "PASS"
		} else {
			result.Status = "FAIL"
			result.ErrorMessage = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
		}
	}

	suite.Results = append(suite.Results, result)

	logger.WithFields(logrus.Fields{
		"test":     testName,
		"status":   result.Status,
		"duration": result.ResponseTime,
		"code":     result.StatusCode,
	}).Info("Test completed")

	return result
}

// generateReport 生成测试报告
func generateReport(suite *TestSuite, logger *logrus.Logger) {
	logger.Info("=== 生成测试报告 ===")

	// 生成JSON报告
	jsonReport, _ := json.MarshalIndent(suite, "", "  ")

	// 保存JSON报告
	jsonFile := fmt.Sprintf("test_report_%s.json", time.Now().Format("20060102_150405"))
	err := os.WriteFile(jsonFile, jsonReport, 0644)
	if err != nil {
		logger.WithError(err).Error("Failed to save JSON report")
	} else {
		logger.WithField("file", jsonFile).Info("JSON report saved")
	}

	// 生成HTML报告
	htmlReport := generateHTMLReport(suite)
	htmlFile := fmt.Sprintf("test_report_%s.html", time.Now().Format("20060102_150405"))
	err = os.WriteFile(htmlFile, []byte(htmlReport), 0644)
	if err != nil {
		logger.WithError(err).Error("Failed to save HTML report")
	} else {
		logger.WithField("file", htmlFile).Info("HTML report saved")
	}

	// 生成控制台报告
	printConsoleReport(suite, logger)
}

// generateHTMLReport 生成HTML报告
func generateHTMLReport(suite *TestSuite) string {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>API测试报告</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .test-result { margin: 10px 0; padding: 10px; border-radius: 3px; }
        .pass { background-color: #d4edda; border-left: 4px solid #28a745; }
        .fail { background-color: #f8d7da; border-left: 4px solid #dc3545; }
        .skip { background-color: #fff3cd; border-left: 4px solid #ffc107; }
        .stats { display: flex; gap: 20px; }
        .stat-box { background-color: #e9ecef; padding: 15px; border-radius: 5px; text-align: center; }
        table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .code { background-color: #f8f9fa; padding: 2px 4px; border-radius: 3px; font-family: monospace; }
    </style>
</head>
<body>
    <div class="header">
        <h1>API测试报告</h1>
        <p><strong>服务器地址:</strong> %s</p>
        <p><strong>测试时间:</strong> %s - %s</p>
        <p><strong>总耗时:</strong> %v</p>
    </div>

    <div class="summary">
        <h2>测试摘要</h2>
        <div class="stats">
            <div class="stat-box">
                <h3>%d</h3>
                <p>总测试数</p>
            </div>
            <div class="stat-box" style="background-color: #d4edda;">
                <h3>%d</h3>
                <p>通过</p>
            </div>
            <div class="stat-box" style="background-color: #f8d7da;">
                <h3>%d</h3>
                <p>失败</p>
            </div>
            <div class="stat-box" style="background-color: #fff3cd;">
                <h3>%d</h3>
                <p>跳过</p>
            </div>
        </div>
    </div>

    <div>
        <h2>详细测试结果</h2>
        <table>
            <tr>
                <th>测试名称</th>
                <th>方法</th>
                <th>端点</th>
                <th>状态</th>
                <th>状态码</th>
                <th>响应时间</th>
                <th>错误信息</th>
            </tr>`,
		suite.ServerURL,
		suite.StartTime.Format("2006-01-02 15:04:05"),
		suite.EndTime.Format("2006-01-02 15:04:05"),
		suite.Duration,
		suite.TotalTests,
		suite.PassedTests,
		suite.FailedTests,
		suite.SkippedTests,
	)

	for _, result := range suite.Results {
		statusClass := strings.ToLower(result.Status)
		html += fmt.Sprintf(`
            <tr class="%s">
                <td>%s</td>
                <td><span class="code">%s</span></td>
                <td><span class="code">%s</span></td>
                <td><strong>%s</strong></td>
                <td>%d</td>
                <td>%v</td>
                <td>%s</td>
            </tr>`,
			statusClass,
			result.TestName,
			result.Method,
			result.Endpoint,
			result.Status,
			result.StatusCode,
			result.ResponseTime,
			result.ErrorMessage,
		)
	}

	html += `
        </table>
    </div>
</body>
</html>`

	return html
}

// printConsoleReport 打印控制台报告
func printConsoleReport(suite *TestSuite, logger *logrus.Logger) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("                           API测试报告")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("服务器地址: %s\n", suite.ServerURL)
	fmt.Printf("测试时间: %s - %s\n",
		suite.StartTime.Format("2006-01-02 15:04:05"),
		suite.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("总耗时: %v\n", suite.Duration)

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("测试摘要:")
	fmt.Printf("  总测试数: %d\n", suite.TotalTests)
	fmt.Printf("  通过: %d (%.1f%%)\n", suite.PassedTests, float64(suite.PassedTests)/float64(suite.TotalTests)*100)
	fmt.Printf("  失败: %d (%.1f%%)\n", suite.FailedTests, float64(suite.FailedTests)/float64(suite.TotalTests)*100)
	fmt.Printf("  跳过: %d (%.1f%%)\n", suite.SkippedTests, float64(suite.SkippedTests)/float64(suite.TotalTests)*100)

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("详细结果:")
	fmt.Printf("%-30s %-6s %-40s %-6s %-8s %-12s %s\n",
		"测试名称", "方法", "端点", "状态", "状态码", "响应时间", "错误信息")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range suite.Results {
		statusColor := ""
		switch result.Status {
		case "PASS":
			statusColor = "✓"
		case "FAIL":
			statusColor = "✗"
		case "SKIP":
			statusColor = "-"
		}

		endpoint := result.Endpoint
		if len(endpoint) > 40 {
			endpoint = endpoint[:37] + "..."
		}

		fmt.Printf("%-30s %-6s %-40s %-6s %-8d %-12v %s\n",
			result.TestName,
			result.Method,
			endpoint,
			statusColor,
			result.StatusCode,
			result.ResponseTime,
			result.ErrorMessage,
		)
	}

	fmt.Println(strings.Repeat("=", 80))

	if suite.FailedTests > 0 {
		fmt.Printf("\n❌ 发现 %d 个失败的测试，请检查上述错误信息\n", suite.FailedTests)
	} else {
		fmt.Println("\n✅ 所有测试都通过了！")
	}
}
