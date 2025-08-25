@echo off
echo Testing CreatePurchaseRequest API with empty delivery_date...
echo.

curl -X POST http://localhost:8082/api/v1/purchase-requests ^
  -H "Content-Type: application/json" ^
  -d "{\"request_id\": \"PR-DB-1703123456-001\", \"doc_type\": \"PR\", \"plant\": \"\", \"quantity\": 6, \"unit_price\": 7000.00, \"material\": \"\", \"delivery_date\": \"\", \"vendor_code\": \"APPLE-001\", \"short_text\": \"\", \"material_group\": \"\", \"unit_type\": \"\", \"requester\": \"\", \"purchase_organization\": \"\", \"currency\": \"\", \"priority\": 2, \"urgency\": \"normal\", \"comments\": \"\"}"

echo.
echo Test completed.
pause 