@echo off
REM 测试更新采购请求API接口的批处理脚本

set BASE_URL=http://localhost:8080/api/v1/purchase-requests

echo === 采购请求更新API测试 ===
echo.

REM 测试1: 更新基本信息
echo 测试1: 更新基本信息（数量、单价、评论）
curl -X PUT "%BASE_URL%/PR-123456" ^
  -H "Content-Type: application/json" ^
  -d "{\"quantity\": 100, \"unit_price\": 30.50, \"comments\": \"Updated quantity and price via API test\"}"
echo.
echo ---
echo.

REM 测试2: 更新状态和审批信息
echo 测试2: 更新状态和审批信息
curl -X PUT "%BASE_URL%/PR-123456" ^
  -H "Content-Type: application/json" ^
  -d "{\"status\": \"pending\", \"approver_id\": \"approver_001\", \"approver_name\": \"Test Approver\"}"
echo.
echo ---
echo.

REM 测试3: 更新优先级和紧急程度
echo 测试3: 更新优先级和紧急程度
curl -X PUT "%BASE_URL%/PR-123456" ^
  -H "Content-Type: application/json" ^
  -d "{\"priority\": 3, \"urgency\": \"urgent\"}"
echo.
echo ---
echo.

REM 测试4: 更新元数据
echo 测试4: 更新元数据
curl -X PUT "%BASE_URL%/PR-123456" ^
  -H "Content-Type: application/json" ^
  -d "{\"metadata\": {\"source\": \"api_test\", \"updated_by\": \"test_script\", \"test_version\": \"1.0\"}}"
echo.
echo ---
echo.

REM 测试5: 使用不同的request_id更新
echo 测试5: 使用不同的request_id更新
curl -X PUT "%BASE_URL%/PR-789012" ^
  -H "Content-Type: application/json" ^
  -d "{\"comments\": \"Updated via different request_id\"}"
echo.
echo ---
echo.

REM 测试6: 错误测试 - 不存在的request_id
echo 测试6: 错误测试 - 不存在的request_id
curl -X PUT "%BASE_URL%/NON-EXISTENT-REQUEST-ID" ^
  -H "Content-Type: application/json" ^
  -d "{\"comments\": \"This should fail\"}"
echo.
echo ---
echo.

echo === 测试完成 ===
pause 