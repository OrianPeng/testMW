# 采购请求更新API接口文档

## 概述

新增了一个全面的采购请求更新接口，支持更新采购请求的所有字段。该接口接收JSON格式的数据，并更新MySQL数据库中的对应记录。

## 接口信息

- **URL**: `PUT /api/v1/purchase-requests/{request_id}`
- **Content-Type**: `application/json`
- **参数**: 
  - `{request_id}`: 采购请求的request_id（字符串）

## 请求格式

### 基本信息字段

```json
{
  "doc_type": "PR",
  "plant": "1000",
  "quantity": 50,
  "unit_price": 25.50,
  "material": "MATERIAL-001",
  "delivery_date": "2024-12-31T00:00:00Z",
  "vendor_code": "VENDOR-002",
  "short_text": "Updated purchase request",
  "material_group": "MG001",
  "unit_type": "PCS",
  "requester": "user123",
  "purchase_organization": "PO001",
  "currency": "CNY",
  "priority": 3,
  "urgency": "urgent",
  "comments": "Updated comments",
  "attachments": "[]",
  "metadata": {
    "source": "api",
    "updated_by": "user123"
  },
  "callback": "http://localhost:8080/callback"
}
```

### 状态和审批字段

```json
{
  "status": "approved",
  "approver_id": "approver_001",
  "approver_name": "John Doe",
  "rejection_reason": null
}
```

### 处理相关字段

```json
{
  "retry_count": 0,
  "error_msg": "",
  "processed_at": "2024-01-01T10:00:00Z"
}
```

## 特性

1. **部分更新**: 只更新JSON中提供的字段，未提供的字段保持不变
2. **自动计算**: 当更新数量或单价时，自动重新计算总金额
3. **状态处理**: 当状态更新为"approved"时，自动设置批准时间
4. **元数据支持**: 支持更新JSON格式的元数据
5. **request_id标识**: 使用request_id进行更新，更加直观和易用

## 响应格式

### 成功响应 (200 OK)

```json
{
  "message": "Purchase request updated successfully",
  "data": {
    "id": 1,
    "request_id": "PR-123456",
    "doc_type": "PR",
    "plant": "1000",
    "quantity": 50,
    "unit_price": 25.50,
    "material": "MATERIAL-001",
    "delivery_date": "2024-12-31T00:00:00Z",
    "vendor_code": "VENDOR-002",
    "short_text": "Updated purchase request",
    "material_group": "MG001",
    "unit_type": "PCS",
    "requester": "user123",
    "purchase_organization": "PO001",
    "currency": "CNY",
    "total_amount": 1275.00,
    "priority": 3,
    "status": "pending",
    "urgency": "urgent",
    "approver_id": null,
    "approver_name": null,
    "approved_at": null,
    "rejection_reason": null,
    "comments": "Updated comments",
    "attachments": "[]",
    "metadata": {
      "source": "api",
      "updated_by": "user123"
    },
    "callback": "http://localhost:8080/callback",
    "retry_count": 0,
    "error_msg": "",
    "processed_at": null,
    "created_at": "2024-01-01T09:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### 错误响应

#### 400 Bad Request
```json
{
  "error": "Invalid request format",
  "details": "具体错误信息"
}
```

#### 404 Not Found
```json
{
  "error": "Purchase request not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to update purchase request"
}
```

## 使用示例

### 使用curl更新采购请求

```bash
# 更新基本信息
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "unit_price": 30.50,
    "comments": "Updated quantity and price"
  }'

# 更新状态
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "approver_id": "approver_001",
    "approver_name": "John Doe"
  }'

# 使用不同的request_id更新
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-789012 \
  -H "Content-Type: application/json" \
  -d '{
    "priority": 4,
    "urgency": "critical"
  }'
```

### 使用JavaScript更新采购请求

```javascript
const updatePurchaseRequest = async (requestId, updateData) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/purchase-requests/${requestId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updateData)
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    console.log('Updated successfully:', result.data);
    return result.data;
  } catch (error) {
    console.error('Error updating purchase request:', error);
    throw error;
  }
};

// 使用示例
updatePurchaseRequest('PR-123456', {
  quantity: 150,
  unit_price: 25.00,
  comments: 'Updated via JavaScript'
});
```

## 注意事项

1. **字段验证**: 所有字段都是可选的，只更新提供的字段
2. **数据类型**: 确保提供的数据类型正确（如数字、字符串、日期等）
3. **日期格式**: 日期字段使用ISO 8601格式（如：`2024-12-31T00:00:00Z`）
4. **元数据**: metadata字段必须是有效的JSON对象
5. **并发安全**: 接口支持并发更新，但建议在高并发场景下使用适当的锁机制
6. **request_id格式**: request_id通常是字符串格式，如"PR-123456"

## 日志记录

接口会记录以下信息：
- 更新操作的request_id
- 更新的主要字段（状态、数量、单价、物料、供应商等）
- 错误信息（如果有）

日志格式示例：
```
INFO Updating purchase request request_id=PR-123456 status=pending quantity=100 unit_price=30.5
``` 