@echo off
echo === 时间格式测试脚本 ===

set BASE_URL=http://localhost:8082/api/v1
set TIMESTAMP=%time:~0,2%%time:~3,2%%time:~6,2%

echo 创建测试采购请求...
curl -X POST %BASE_URL%/purchase-requests ^
  -H "Content-Type: application/json" ^
  -d "{\"request_id\":\"PR-TIME-TEST-%TIMESTAMP%\",\"doc_type\":\"PR\",\"plant\":\"1000\",\"quantity\":1,\"unit_price\":100.00,\"material\":\"TIME-TEST\",\"delivery_date\":\"2026-05-19\",\"vendor_code\":\"VENDOR-TEST\",\"short_text\":\"Time format test\",\"material_group\":\"TEST\",\"unit_type\":\"EA\",\"requester\":\"Test User\",\"purchase_organization\":\"TEST-ORG\",\"currency\":\"CNY\",\"priority\":1,\"urgency\":\"normal\",\"comments\":\"Time format test\"}"

echo.
echo 测试时间格式: 2026-05-19
curl -X PUT %BASE_URL%/purchase-requests/PR-TIME-TEST-%TIMESTAMP% ^
  -H "Content-Type: application/json" ^
  -d "{\"delivery_date\":\"2026-05-19\",\"comments\":\"Testing date format\"}"

echo.
echo 测试时间格式: 2026-05-19T15:30:00
curl -X PUT %BASE_URL%/purchase-requests/PR-TIME-TEST-%TIMESTAMP% ^
  -H "Content-Type: application/json" ^
  -d "{\"delivery_date\":\"2026-05-19T15:30:00\",\"comments\":\"Testing datetime format\"}"

echo.
echo 测试时间格式: 2026/05/19
curl -X PUT %BASE_URL%/purchase-requests/PR-TIME-TEST-%TIMESTAMP% ^
  -H "Content-Type: application/json" ^
  -d "{\"delivery_date\":\"2026/05/19\",\"comments\":\"Testing slash format\"}"

echo.
echo 测试时间格式: 2026-05-19T15:30:00Z
curl -X PUT %BASE_URL%/purchase-requests/PR-TIME-TEST-%TIMESTAMP% ^
  -H "Content-Type: application/json" ^
  -d "{\"delivery_date\":\"2026-05-19T15:30:00Z\",\"comments\":\"Testing UTC format\"}"

echo.
echo 测试空字符串
curl -X PUT %BASE_URL%/purchase-requests/PR-TIME-TEST-%TIMESTAMP% ^
  -H "Content-Type: application/json" ^
  -d "{\"delivery_date\":\"\",\"comments\":\"Testing empty string\"}"

echo.
echo === 时间格式测试完成 === 