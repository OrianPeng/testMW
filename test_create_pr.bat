@echo off
echo Testing Create Purchase Request API...
echo.

echo Sending request to create purchase request...
curl -X POST http://localhost:8080/api/v1/purchase-requests ^
  -H "Content-Type: application/json" ^
  -d @test_create_pr.json

echo.
echo Test completed.
pause 