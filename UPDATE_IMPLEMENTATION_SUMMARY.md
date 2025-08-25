# 采购请求更新接口实现总结

## 概述

成功实现了一个全面的采购请求更新接口，支持通过JSON格式的数据更新MySQL数据库中的PurchaseRequest记录。接口使用request_id作为标识符，更加直观和易用。

## 实现的功能

### 1. 扩展的更新请求模型

在 `internal/models/purchase_request.go` 中扩展了 `UpdatePurchaseRequestRequest` 结构体，支持更新以下字段：

**基本信息字段：**
- `doc_type` - 文档类型
- `plant` - 工厂
- `quantity` - 数量
- `unit_price` - 单价
- `material` - 物料
- `delivery_date` - 交付日期
- `vendor_code` - 供应商代码
- `short_text` - 简短描述
- `material_group` - 物料组
- `unit_type` - 单位类型
- `requester` - 申请人
- `purchase_organization` - 采购组织
- `currency` - 货币
- `priority` - 优先级
- `urgency` - 紧急程度
- `comments` - 评论
- `attachments` - 附件
- `metadata` - 元数据
- `callback` - 回调URL

**状态和审批字段：**
- `status` - 状态
- `approver_id` - 审批人ID
- `approver_name` - 审批人姓名
- `rejection_reason` - 拒绝原因

**处理相关字段：**
- `retry_count` - 重试次数
- `error_msg` - 错误信息
- `processed_at` - 处理时间

### 2. 增强的Repository层

在 `internal/repository/purchase_request_repository.go` 中更新了 `Update` 方法并新增了 `UpdateByRequestID` 方法：

- **部分更新支持**：只更新JSON中提供的字段
- **自动计算**：当更新数量或单价时，自动重新计算总金额
- **状态处理**：当状态更新为"approved"时，自动设置批准时间
- **元数据序列化**：支持JSON格式的元数据更新
- **错误处理**：完善的错误处理和日志记录
- **request_id支持**：通过request_id进行更新操作

### 3. 改进的Handler层

在 `internal/api/purchase_request_handlers.go` 中增强了 `UpdatePurchaseRequest` 方法：

- **request_id支持**：使用request_id进行更新操作
- **详细日志记录**：记录更新操作的关键信息
- **错误处理**：完善的错误处理和响应

### 4. 路由配置

在 `internal/api/router.go` 中更新了路由配置：
```go
purchase.PUT("/:request_id", purchaseHandler.UpdatePurchaseRequest)
purchase.GET("/:request_id", purchaseHandler.GetPurchaseRequest)
purchase.DELETE("/:request_id", purchaseHandler.DeletePurchaseRequest)
```

## 接口特性

### 1. 部分更新
- 只更新JSON中提供的字段
- 未提供的字段保持不变
- 支持任意字段组合的更新

### 2. 自动计算
- 当更新数量或单价时，自动重新计算总金额
- 当状态更新为"approved"时，自动设置批准时间

### 3. request_id标识
- 使用request_id进行更新操作
- 更加直观和易用
- 支持字符串格式的标识符

### 4. 数据验证
- JSON格式验证
- 数据类型验证
- 字段存在性验证

## 使用示例

### 基本更新
```bash
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "unit_price": 30.50,
    "comments": "Updated via API"
  }'
```

### 状态更新
```bash
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "approver_id": "approver_001",
    "approver_name": "John Doe"
  }'
```

### 元数据更新
```bash
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "source": "api",
      "updated_by": "user123",
      "update_reason": "testing"
    }
  }'
```

## 响应格式

### 成功响应 (200 OK)
```json
{
  "message": "Purchase request updated successfully",
  "data": {
    "id": 1,
    "request_id": "PR-123456",
    "quantity": 100,
    "unit_price": 30.50,
    "total_amount": 3050.00,
    "status": "pending",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### 错误响应
```json
{
  "error": "Purchase request not found"
}
```

## 测试文件

创建了以下测试文件：

1. **cmd/middle_test/main.go** - 完整的API测试，包括更新接口测试
2. **test_update.bat** - Windows批处理脚本，用于快速测试更新接口
3. **test_requests.json** - 测试用的请求数据
4. **test_update_request.json** - 测试用的更新请求数据

## 新增的Repository方法

### UpdateByRequestID
```go
func (r *PurchaseRequestRepository) UpdateByRequestID(ctx context.Context, requestID string, updateReq *models.UpdatePurchaseRequestRequest) (*models.PurchaseRequest, error)
```

### DeleteByRequestID
```go
func (r *PurchaseRequestRepository) DeleteByRequestID(ctx context.Context, requestID string) error
```

## 主要变更

1. **路由参数**：从`:id`改为`:request_id`
2. **处理器方法**：使用request_id参数而不是数字ID
3. **Repository方法**：新增通过request_id操作的方法
4. **测试用例**：更新所有测试用例使用request_id
5. **文档更新**：更新API文档和实现总结

## 优势

1. **更直观**：使用request_id比数字ID更容易理解和记忆
2. **更安全**：避免数字ID可能带来的安全问题
3. **更灵活**：request_id可以包含业务含义
4. **更一致**：与其他业务系统保持一致 