# 重构总结：合并 Request 和 Purchase-Requests

## 重构目标

将原有的 `request` 和 `purchase-requests` 功能合并，以 `purchase-requests` 为准，增加队列交互功能，移除原有的通用 `request` 功能。

## 完成的重构内容

### 1. 数据模型更新

#### ✅ 更新 `PurchaseRequest` 模型
- 增加了队列相关字段：
  - `Callback` - 回调URL
  - `RetryCount` - 重试次数
  - `ErrorMsg` - 错误信息
  - `ProcessedAt` - 处理完成时间
- 增加了新的状态：
  - `PurchaseStatusProcessing` - 正在处理（队列中）
  - `PurchaseStatusFailed` - 处理失败
- 新增队列相关结构体：
  - `QueuedPurchaseRequest` - 队列中的采购请求
  - `RPARequest` - 发送给RPA系统的请求
  - `RPAResponse` - RPA系统的响应
  - `PurchaseRequestCallbackResponse` - 采购请求回调响应

#### ✅ 删除 `request.go`
- 移除了原有的通用请求模型
- 所有功能统一使用 `purchase-requests`

### 2. 数据库层更新

#### ✅ 更新数据库表结构
- 在 `purchase_requests` 表中增加字段：
  - `callback VARCHAR(500) NULL`
  - `retry_count INT DEFAULT 0`
  - `error_msg TEXT NULL`
  - `processed_at DATETIME NULL`
- 增加相关索引：
  - `idx_callback (callback)`
  - `idx_retry_count (retry_count)`

#### ✅ 更新数据访问层
- 更新 `PurchaseRequestRepository` 支持新字段
- 修改所有SQL查询语句包含新字段

### 3. 队列系统重构

#### ✅ 更新接口定义
- 修改 `QueueManager` 接口，改为处理采购请求
- 更新 `RPAClient` 接口，改为发送采购请求
- 更新 `NotificationService` 接口，改为处理采购请求回调

#### ✅ 更新队列实现
- **内存队列** (`MemoryQueueManager`)：
  - 改为处理 `QueuedPurchaseRequest`
  - 支持优先级排序
  - 增加队列位置管理
  - 支持重试机制

- **Redis队列** (`RedisQueueManager`)：
  - 改为处理 `QueuedPurchaseRequest`
  - 使用Redis有序集合实现优先级队列
  - 支持批量操作和位置更新

#### ✅ 更新队列处理器
- 修改 `QueueProcessor` 处理采购请求
- 增加采购请求特定的错误处理
- 支持采购请求的回调通知

### 4. RPA客户端更新

#### ✅ 更新RPA客户端
- 修改 `HTTPClient` 发送采购请求
- 更新请求URL为 `/api/process-purchase`
- 增加采购请求特定的数据映射

### 5. 通知服务更新

#### ✅ 更新通知服务
- 修改 `HTTPNotifier` 处理采购请求回调
- 增加采购请求特定的通知格式
- 支持失败重试通知

### 6. API层重构

#### ✅ 更新API处理器
- 移除原有的通用请求处理器
- 新增采购请求队列处理器：
  - `SubmitPurchaseRequest` - 提交采购请求到队列
  - `GetPurchaseRequestStatus` - 获取采购请求状态
  - `ListPurchaseRequests` - 列出队列中的采购请求
  - `GetQueueStats` - 获取队列统计信息
  - `HealthCheck` - 健康检查

#### ✅ 更新路由配置
- 移除原有的 `/requests` 路由
- 新增采购请求队列路由：
  - `POST /api/v1/purchase-requests/queue` - 提交到队列
  - `GET /api/v1/purchase-requests/queue/:id` - 获取状态
  - `GET /api/v1/purchase-requests/queue` - 列出队列
- 保留原有的采购请求数据库管理路由

### 7. 主程序更新

#### ✅ 更新主程序
- 确保使用新的路由配置
- 保持队列处理器和数据库的集成

## 新的API接口

### 采购请求队列管理
```
POST /api/v1/purchase-requests/queue          # 提交采购请求到队列
GET  /api/v1/purchase-requests/queue/:id      # 获取队列中的采购请求状态
GET  /api/v1/purchase-requests/queue          # 列出队列中的采购请求
```

### 采购请求数据库管理
```
POST   /api/v1/purchase-requests              # 创建采购请求（数据库）
GET    /api/v1/purchase-requests/:id          # 获取采购请求详情
GET    /api/v1/purchase-requests              # 查询采购请求列表
PUT    /api/v1/purchase-requests/:id          # 更新采购请求
DELETE /api/v1/purchase-requests/:id          # 删除采购请求
```

### 系统状态
```
GET /api/v1/status                            # 获取队列和RPA状态
GET /api/v1/health                            # 健康检查
```

## 工作流程

### 1. 采购请求提交流程
```
客户端 → POST /purchase-requests/queue → 队列管理器 → 队列存储
```

### 2. 采购请求处理流程
```
队列处理器 → 从队列取请求 → 发送给RPA系统 → 更新状态 → 通知回调
```

### 3. 状态管理
- `pending` → `processing` → `completed`/`failed`
- 支持重试机制（最多3次）
- 失败后发送回调通知

## 测试验证

创建了测试文件 `test_purchase_queue.go` 来验证：
- 采购请求创建
- 队列状态查询
- 采购请求状态查询
- 采购请求列表查询

## 总结

✅ **重构完成**：成功将 `request` 和 `purchase-requests` 合并
✅ **功能增强**：增加了完整的队列交互功能
✅ **代码清理**：移除了冗余的通用请求代码
✅ **接口统一**：所有功能统一使用采购请求模型
✅ **向后兼容**：保留了原有的采购请求数据库管理功能

现在系统专注于采购请求的处理，具有完整的队列管理、RPA集成、状态跟踪和回调通知功能。 