# 前端API文档

## 概述

本文档描述了RPA中间件系统的前端API接口，包括采购请求、采购订单、供应商管理、流程监控、仪表板和AI助手等功能模块。

## 基础配置

### API基础URL
```
/api/v1
```

### 请求配置
- **超时时间**: 10秒
- **Content-Type**: application/json
- **认证方式**: Bearer Token (存储在localStorage)

### 请求拦截器
- 自动添加Authorization头
- 处理401未授权错误，自动跳转到登录页

## API接口

### 1. 采购请求管理 (Purchase Requests)

#### 1.1 获取采购请求列表
```typescript
GET /purchase-requests
```

**查询参数:**
- `status` (string): 状态筛选
- `requester` (string): 申请人筛选
- `plant` (string): 工厂筛选
- `vendor_code` (string): 供应商代码筛选
- `material_group` (string): 物料组筛选
- `priority` (number): 优先级筛选
- `urgency` (string): 紧急程度筛选
- `start_date` (string): 开始日期
- `end_date` (string): 结束日期
- `page` (number): 页码
- `limit` (number): 每页数量
- `sort_by` (string): 排序字段
- `sort_order` (string): 排序方向 ('asc' | 'desc')

**响应示例:**
```json
{
  "data": [
    {
      "request_id": "PR-12345",
      "requester": "张三",
      "material": "MAT001",
      "quantity": 100,
      "unit_price": 10.50,
      "total_amount": 1050.00,
      "status": "pending",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

#### 1.2 获取采购请求详情
```typescript
GET /purchase-requests/{requestId}
```

#### 1.3 创建采购请求
```typescript
POST /purchase-requests
```

**请求体:**
```json
{
  "request_id": "PR-12345",
  "requester": "张三",
  "material": "MAT001",
  "quantity": 100,
  "unit_price": 10.50,
  "delivery_date": "2024-02-01",
  "vendor_code": "V001",
  "short_text": "采购说明",
  "plant": "P001",
  "doc_type": "PR",
  "currency": "CNY"
}
```

#### 1.4 更新采购请求
```typescript
PUT /purchase-requests/{requestId}
```

#### 1.5 删除采购请求
```typescript
DELETE /purchase-requests/{requestId}
```

#### 1.6 获取采购请求统计
```typescript
GET /purchase-requests/statistics
```

### 2. 采购订单管理 (Purchase Orders)

#### 2.1 获取采购订单列表
```typescript
GET /purchase-orders
```

**查询参数:**

- `status` (string): 状态筛选
- `supplier_name` (string): 供应商名称筛选
- `supplier_code` (string): 供应商代码筛选
- `page` (number): 页码
- `limit` (number): 每页数量
- `sort_by` (string): 排序字段
- `sort_order` (string): 排序方向

#### 2.2 获取采购订单详情
```typescript
GET /purchase-orders/{poNumber}
```

#### 2.3 创建采购订单
```typescript
POST /purchase-orders
```

#### 2.4 更新采购订单
```typescript
PUT /purchase-orders/{poNumber}
```

#### 2.5 删除采购订单
```typescript
DELETE /purchase-orders/{poNumber}
```

#### 2.6 获取采购订单统计
```typescript
GET /purchase-orders/statistics
```

### 3. 供应商管理 (Suppliers)

#### 3.1 获取供应商列表
```typescript
GET /suppliers
```

**查询参数:**
- `status` (string): 状态筛选
- `type` (string): 类型筛选
- `page` (number): 页码
- `limit` (number): 每页数量
- `sort_by` (string): 排序字段
- `sort_order` (string): 排序方向

#### 3.2 获取供应商详情
```typescript
GET /suppliers/{id}
```

#### 3.3 创建供应商
```typescript
POST /suppliers
```

#### 3.4 更新供应商
```typescript
PUT /suppliers/{id}
```

#### 3.5 删除供应商
```typescript
DELETE /suppliers/{id}
```

#### 3.6 获取供应商统计
```typescript
GET /suppliers/statistics
```

### 4. 流程监控 (Process Monitor)

#### 4.1 获取流程概览
```typescript
GET /process-monitor/overview
```

#### 4.2 获取被阻塞的流程
```typescript
GET /process-monitor/blocked
```

#### 4.3 获取流程步骤
```typescript
GET /process-monitor/steps
```

**查询参数:**
- `process_type` (string): 流程类型筛选

### 5. 仪表板 (Dashboard)

#### 5.1 获取仪表板指标
```typescript
GET /dashboard/metrics
```

#### 5.2 获取PR趋势图表
```typescript
GET /dashboard/charts/pr-trend
```

#### 5.3 获取PO状态分布
```typescript
GET /dashboard/charts/po-status
```

#### 5.4 获取供应商排名
```typescript
GET /dashboard/charts/supplier-ranking
```

#### 5.5 获取部门统计
```typescript
GET /dashboard/charts/department-stats
```

### 6. AI助手 (AI Assistant)

#### 6.1 发送AI消息
```typescript
POST /ai-assistant/message
```

**请求体:**
```json
{
  "message": "请帮我查询采购请求状态",
  "conversation_id": "conv_12345"
}
```

**响应示例:**
```json
{
  "response": "根据您的查询，当前有5个待处理的采购请求...",
  "conversation_id": "conv_12345",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 7. 系统状态 (System Status)

#### 7.1 获取队列统计
```typescript
GET /status
```

#### 7.2 清空队列
```typescript
POST /purchase-requests/queue/clear-all
```

#### 7.3 健康检查
```typescript
GET /health
```

### 8. 文件上传

#### 8.1 上传文件
```typescript
POST /upload
```

**请求格式:** multipart/form-data

**响应示例:**
```json
{
  "success": true,
  "data": {
    "file_id": "file_1704067200",
    "filename": "document.pdf",
    "size": 1024,
    "url": "/files/file_1704067200"
  }
}
```

## 错误处理

### 标准错误响应格式
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述",
    "details": "详细错误信息"
  }
}
```

### 常见错误码
- `VALIDATION_ERROR`: 请求参数验证失败
- `NOT_FOUND`: 资源不存在
- `UNAUTHORIZED`: 未授权访问
- `FORBIDDEN`: 禁止访问
- `INTERNAL_ERROR`: 服务器内部错误

## 使用示例

### JavaScript/TypeScript 使用示例

```typescript
import { purchaseApi } from '@/api/purchase'

// 获取采购请求列表
const getPurchaseRequests = async () => {
  try {
    const response = await purchaseApi.getPurchaseRequests({
      status: 'pending',
      page: 1,
      limit: 20
    })
    console.log(response.data)
  } catch (error) {
    console.error('获取采购请求失败:', error)
  }
}

// 创建采购请求
const createPurchaseRequest = async () => {
  try {
    const response = await purchaseApi.createPurchaseRequest({
      request_id: 'PR-12345',
      requester: '张三',
      material: 'MAT001',
      quantity: 100,
      unit_price: 10.50,
      delivery_date: '2024-02-01',
      vendor_code: 'V001',
      short_text: '采购说明',
      plant: 'P001',
      doc_type: 'PR',
      currency: 'CNY'
    })
    console.log('创建成功:', response)
  } catch (error) {
    console.error('创建失败:', error)
  }
}
```

## 注意事项

1. 所有API请求都需要包含正确的Content-Type头
2. 认证token会自动添加到请求头中
3. 请求超时时间为10秒
4. 分页参数page从1开始
5. 日期格式使用ISO 8601标准
6. 货币代码使用ISO 4217标准
7. 所有响应都遵循统一的JSON格式

## 更新日志

- **v1.0.0** (2024-01-01): 初始版本，包含基础API接口
- **v1.1.0** (2024-01-15): 添加AI助手和流程监控功能
- **v1.2.0** (2024-02-01): 添加仪表板和统计功能
