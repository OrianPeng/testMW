# 后端API文档

## 概述

本文档描述了RPA中间件系统的后端API接口，包括采购请求管理、采购订单管理、供应商管理、流程监控、仪表板、AI助手、UiPath集成和队列管理等功能模块。

## 基础信息

### 服务器配置
- **端口**: 8082 (默认)
- **基础路径**: `/api/v1`
- **协议**: HTTP/HTTPS
- **内容类型**: application/json

### 中间件
- **CORS**: 支持跨域请求
- **日志记录**: 自动记录API请求和响应
- **错误恢复**: 自动恢复panic错误

## API接口

### 1. 采购请求管理 (Purchase Requests)

#### 1.1 创建采购请求
```http
POST /api/v1/purchase-requests
```

**请求体:**
```json
{
  "request_id": "PR-12345",
  "requester": "张三",
  "material": "MAT001",
  "quantity": 100,
  "unit_price": 10.50,
  "delivery_date": "2024-02-01T00:00:00Z",
  "vendor_code": "V001",
  "short_text": "采购说明",
  "plant": "P001",
  "doc_type": "PR",
  "currency": "CNY",
  "priority": 1,
  "urgency": "normal",
  "material_group": "MG001",
  "department": "IT",
  "cost_center": "CC001",
  "project_code": "PROJ001",
  "approver": "李四",
  "notes": "备注信息"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "request_id": "PR-12345",
    "requester": "张三",
    "material": "MAT001",
    "quantity": 100,
    "unit_price": 10.50,
    "total_amount": 1050.00,
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 1.2 获取采购请求详情
```http
GET /api/v1/purchase-requests/{request_id}
```

#### 1.3 获取采购请求列表
```http
GET /api/v1/purchase-requests
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
- `sort_order` (string): 排序方向

#### 1.4 更新采购请求
```http
PUT /api/v1/purchase-requests/{request_id}
```

#### 1.5 删除采购请求
```http
DELETE /api/v1/purchase-requests/{request_id}
```

#### 1.6 获取采购请求统计
```http
GET /api/v1/purchase-requests/statistics
```

### 2. 采购订单管理 (Purchase Orders)

#### 2.1 创建采购订单
```http
POST /api/v1/purchase-orders
```

**请求体:**
```json
{
  "po_number": "PO-12345",
  "supplier_id": "SUP001",
  "supplier_name": "供应商A",
  "supplier_code": "V001",
  "total_amount": 1050.00,
  "currency": "CNY",
  "status": "pending",
  "order_date": "2024-01-01T10:00:00Z",
  "expected_delivery": "2024-02-01T00:00:00Z",
  "items": [
    {
      "material": "MAT001",
      "quantity": 100,
      "unit_price": 10.50,
      "total_price": 1050.00
    }
  ]
}
```

#### 2.2 获取采购订单详情
```http
GET /api/v1/purchase-orders/{po_number}
```

#### 2.3 获取采购订单列表
```http
GET /api/v1/purchase-orders
```

#### 2.4 更新采购订单
```http
PUT /api/v1/purchase-orders/{po_number}
```

#### 2.5 删除采购订单
```http
DELETE /api/v1/purchase-orders/{po_number}
```

#### 2.6 获取采购订单统计
```http
GET /api/v1/purchase-orders/statistics
```

### 3. 供应商管理 (Suppliers)

#### 3.1 创建供应商
```http
POST /api/v1/suppliers
```

**请求体:**
```json
{
  "supplier_code": "V001",
  "supplier_name": "供应商A",
  "contact_person": "联系人",
  "email": "contact@supplier.com",
  "phone": "1234567890",
  "address": "供应商地址",
  "status": "active",
  "supplier_type": "vendor",
  "payment_terms": "30天",
  "currency": "CNY"
}
```

#### 3.2 获取供应商详情
```http
GET /api/v1/suppliers/{id}
```

#### 3.3 获取供应商列表
```http
GET /api/v1/suppliers
```

#### 3.4 更新供应商
```http
PUT /api/v1/suppliers/{id}
```

#### 3.5 删除供应商
```http
DELETE /api/v1/suppliers/{id}
```

#### 3.6 获取供应商统计
```http
GET /api/v1/suppliers/statistics
```

### 4. 流程监控 (Process Monitor)

#### 4.1 获取流程概览
```http
GET /api/v1/process-monitor/overview
```

**响应:**
```json
{
  "total_processes": 100,
  "active_processes": 80,
  "blocked_processes": 5,
  "completed_processes": 15,
  "process_types": {
    "purchase_request": 60,
    "purchase_order": 30,
    "supplier_approval": 10
  }
}
```

#### 4.2 获取被阻塞的流程
```http
GET /api/v1/process-monitor/blocked
```

#### 4.3 获取流程步骤
```http
GET /api/v1/process-monitor/steps
```

**查询参数:**
- `process_type` (string): 流程类型筛选

### 5. 仪表板 (Dashboard)

#### 5.1 获取仪表板指标
```http
GET /api/v1/dashboard/metrics
```

**响应:**
```json
{
  "total_pr": 1000,
  "pending_pr": 50,
  "approved_pr": 800,
  "rejected_pr": 150,
  "total_po": 500,
  "pending_po": 20,
  "completed_po": 450,
  "total_suppliers": 100,
  "active_suppliers": 95
}
```

#### 5.2 获取PR趋势图表
```http
GET /api/v1/dashboard/charts/pr-trend
```

#### 5.3 获取PO状态分布
```http
GET /api/v1/dashboard/charts/po-status
```

#### 5.4 获取供应商排名
```http
GET /api/v1/dashboard/charts/supplier-ranking
```

#### 5.5 获取部门统计
```http
GET /api/v1/dashboard/charts/department-stats
```

### 6. AI助手 (AI Assistant)

#### 6.1 发送AI消息
```http
POST /api/v1/ai-assistant/message
```

**请求体:**
```json
{
  "message": "请帮我查询采购请求状态",
  "conversation_id": "conv_12345"
}
```

**响应:**
```json
{
  "response": "根据您的查询，当前有5个待处理的采购请求...",
  "conversation_id": "conv_12345",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 7. UiPath RPA集成

#### 7.1 添加队列项目
```http
POST /api/v1/uipath/queue/add
```

**请求体:**
```json
{
  "process_type": "purchase_request",
  "data": {
    "request_id": "PR-12345",
    "material": "MAT001",
    "quantity": 100
  },
  "priority": "normal"
}
```

#### 7.2 检查UiPath状态
```http
GET /api/v1/uipath/status
```

#### 7.3 检查UiPath可用性
```http
GET /api/v1/uipath/available
```

#### 7.4 测试UiPath连接
```http
GET /api/v1/uipath/test
```

#### 7.5 添加测试队列项目
```http
POST /api/v1/uipath/queue/test
```

### 8. 队列管理 (Queue Management)

#### 8.1 提交采购请求到队列
```http
POST /api/v1/v1/purchase-requests/queue
```

#### 8.2 清空队列
```http
POST /api/v1/v1/purchase-requests/queue/clear-all
```

#### 8.3 获取队列中的采购请求
```http
GET /api/v1/v1/purchase-requests/queue
```

#### 8.4 获取采购请求状态
```http
GET /api/v1/v1/purchase-requests/queue/{id}
```

#### 8.5 获取队列统计
```http
GET /api/v1/v1/status
```

#### 8.6 健康检查
```http
GET /api/v1/v1/health
```

**响应:**
```json
{
  "status": "healthy",
  "components": {
    "queue": {
      "healthy": true,
      "error": null
    },
    "rpa": {
      "healthy": true
    }
  }
}
```

### 9. 文件上传

#### 9.1 上传文件
```http
POST /api/v1/upload
```

**请求格式:** multipart/form-data

**响应:**
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

### HTTP状态码
- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未授权访问
- `403 Forbidden`: 禁止访问
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突
- `422 Unprocessable Entity`: 请求格式正确但语义错误
- `500 Internal Server Error`: 服务器内部错误
- `503 Service Unavailable`: 服务不可用

### 常见错误码
- `VALIDATION_ERROR`: 请求参数验证失败
- `NOT_FOUND`: 资源不存在
- `UNAUTHORIZED`: 未授权访问
- `FORBIDDEN`: 禁止访问
- `CONFLICT`: 资源冲突
- `INTERNAL_ERROR`: 服务器内部错误
- `SERVICE_UNAVAILABLE`: 服务不可用

## 数据模型

### 采购请求 (PurchaseRequest)
```json
{
  "request_id": "string",
  "requester": "string",
  "material": "string",
  "quantity": "number",
  "unit_price": "number",
  "total_amount": "number",
  "delivery_date": "datetime",
  "vendor_code": "string",
  "short_text": "string",
  "plant": "string",
  "doc_type": "string",
  "currency": "string",
  "status": "string",
  "priority": "number",
  "urgency": "string",
  "material_group": "string",
  "department": "string",
  "cost_center": "string",
  "project_code": "string",
  "approver": "string",
  "notes": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### 采购订单 (PurchaseOrder)
```json
{
  "po_number": "string",
  "supplier_id": "string",
  "supplier_name": "string",
  "supplier_code": "string",
  "total_amount": "number",
  "currency": "string",
  "status": "string",
  "order_date": "datetime",
  "expected_delivery": "datetime",
  "items": "array",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### 供应商 (Supplier)
```json
{
  "id": "string",
  "supplier_code": "string",
  "supplier_name": "string",
  "contact_person": "string",
  "email": "string",
  "phone": "string",
  "address": "string",
  "status": "string",
  "supplier_type": "string",
  "payment_terms": "string",
  "currency": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## 认证和授权

### 认证方式
- **Bearer Token**: 在Authorization头中传递
- **Token存储**: 建议存储在localStorage中

### 权限控制
- 所有API接口都需要有效的认证token
- 不同角色可能有不同的访问权限
- 敏感操作需要额外的权限验证

## 限流和配额

### 请求限制
- **频率限制**: 每分钟最多1000个请求
- **并发限制**: 最多100个并发请求
- **文件上传**: 单个文件最大10MB

### 配额管理
- 每个用户每天最多创建100个采购请求
- 每个用户每天最多创建50个采购订单
- 每个用户每天最多创建10个供应商

## 监控和日志

### 日志记录
- 所有API请求都会记录到日志中
- 包含请求方法、路径、IP地址、响应状态等信息
- 支持结构化日志格式

### 监控指标
- 请求响应时间
- 错误率统计
- 并发连接数
- 队列长度监控

## 版本控制

### API版本
- 当前版本: v1
- 版本路径: `/api/v1`
- 向后兼容性: 保持向后兼容至少6个月

### 版本升级
- 重大版本升级会提前30天通知
- 废弃的API会标记为deprecated
- 提供迁移指南和工具

## 更新日志

- **v1.0.0** (2024-01-01): 初始版本，包含基础API接口
- **v1.1.0** (2024-01-15): 添加AI助手和流程监控功能
- **v1.2.0** (2024-02-01): 添加仪表板和统计功能
- **v1.3.0** (2024-02-15): 添加UiPath集成和文件上传功能
