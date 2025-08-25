# RPA Middleware 项目介绍

## 📋 项目概述

RPA Middleware 是一个基于 Go 语言开发的高性能中间件系统，专门用于管理采购请求的自动化处理流程。该系统作为 AI Agent 与 RPA（机器人流程自动化）系统之间的桥梁，提供完整的任务队列管理、状态跟踪和回调通知功能。

## 🏗️ 系统架构

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│  AI Agent   │───▶│   Middleware │───▶│  RPA System │
│             │    │              │    │             │
└─────────────┘    └──────────────┘    └─────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │    Queue     │
                   │ (Memory/Redis)│
                   └──────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │   MySQL DB   │
                   │ (持久化存储)  │
                   └──────────────┘
```

## 🚀 核心功能

### 1 采购请求管理
- **创建采购请求**：支持完整的采购信息录入，包括物料、数量、价格、供应商等
- **状态跟踪**：实时跟踪采购请求的处理状态（待处理、处理中、已完成、失败等）
- **优先级管理**：支持14级优先级设置，确保重要请求优先处理
- **信息完整性检查**：自动检查采购请求信息的完整性，标记为"信息完整或信息不完整###2. 智能队列系统
- **双队列模式**：支持内存队列（开发模式）和Redis队列（生产模式）
- **优先级队列**：基于优先级的智能排序，高优先级请求优先处理
- **重试机制**：失败请求自动重试，最多3次，支持自定义重试间隔
- **队列监控**：实时监控队列状态，包括待处理数量、处理进度等

### 3RPA系统集成
- **RPA客户端**：HTTP客户端与RPA系统通信
- **状态检查**：定期检查RPA系统可用性
- **任务分发**：将采购请求转换为RPA可执行的格式并发送
- **响应处理**：处理RPA系统的执行结果和错误信息

### 4. 通知服务
- **回调通知**：任务完成后自动发送回调通知到指定URL
- **失败通知**：处理失败时发送详细错误信息
- **HTTP通知器**：基于HTTP的通知服务，支持超时和重试

### 5. 数据库持久化
- **MySQL存储**：所有采购请求数据持久化到MySQL数据库
- **事务支持**：确保数据一致性和完整性
- **索引优化**：针对查询场景优化的数据库索引

### 6 健康监控
- **系统状态检查**：监控队列和RPA系统状态
- **健康检查接口**：提供标准的健康检查端点
- **详细日志**：结构化JSON日志，支持不同级别

## 📊 MySQL 数据库表结构

### purchase_requests 表

```sql
CREATE TABLE purchase_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL UNIQUE,           -- 请求唯一标识
    doc_type VARCHAR(50) NOT NULL,                     -- 文档类型
    plant VARCHAR(50) NOT NULL,                        -- 工厂代码
    quantity INT NOT NULL,                             -- 数量
    unit_price DECIMAL(10,2) NOT NULL,                 -- 单价
    material VARCHAR(100) NOT NULL,                    -- 物料编码
    delivery_date DATETIME NOT NULL,                   -- 交付日期
    vendor_code VARCHAR(50) NOT NULL,                  -- 供应商代码
    short_text TEXT NOT NULL,                          -- 简短描述
    material_group VARCHAR(50) NOT NULL,               -- 物料组
    unit_type VARCHAR(20) NOT NULL,                    -- 单位类型
    requester VARCHAR(100) NOT NULL,                   -- 申请人
    purchase_organization VARCHAR(50) NOT NULL,        -- 采购组织
    currency VARCHAR(10) NOT NULL DEFAULT CNY,   -- 货币
    total_amount DECIMAL(10,2) NOT NULL,               -- 总金额
    priority INT DEFAULT 2,                            -- 优先级(1-4)
    status VARCHAR(20) DEFAULT 'pending,              -- 状态
    urgency VARCHAR(20 DEFAULT 'normal,              -- 紧急程度
    approver_id VARCHAR(100) NULL,                     -- 审批人ID
    approver_name VARCHAR(200) NULL,                   -- 审批人姓名
    approved_at DATETIME NULL,                         -- 审批时间
    rejection_reason TEXT NULL,                        -- 拒绝原因
    comments TEXT,                                     -- 备注
    attachments TEXT,                                  -- 附件信息
    metadata JSON,                                     -- 元数据
    callback VARCHAR(500),                             -- 回调URL
    retry_count INT DEFAULT 0,                         -- 重试次数
    error_msg TEXT,                                    -- 错误信息
    processed_at DATETIME NULL,                        -- 处理完成时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    
    -- 索引
    INDEX idx_request_id (request_id),
    INDEX idx_requester (requester),
    INDEX idx_plant (plant),
    INDEX idx_status (status),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at),
    INDEX idx_status_priority (status, priority),
    INDEX idx_vendor_code (vendor_code),
    INDEX idx_material_group (material_group)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4TE=utf8mb4_unicode_ci;
```

### 状态说明
- `pending`: 待处理
- `processing`: 处理中（队列中）
- `completed`: 已完成
- `failed`: 处理失败
- `approved`: 已批准
- `rejected`: 已拒绝
- `cancelled`: 已取消
- `enough`: 信息完整
- `lack`: 信息不完整

## 📚 API 接口文档

### 基础信息
- **基础URL**: `http://localhost:882pi/v1`
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`

### 1. 采购请求队列管理

####1.1 提交采购请求到队列
```http
POST /purchase-requests/queue
```

**请求体**:
```json[object Object]
  request_id: PR-202401,
  doc_type": PR,
  plant": "1000,
 quantity":10,
  unit_price": 100.50 material":MAT-1,
  delivery_date: 2024123159Z",
 vendor_code:VENDOR-01,
 short_text": "办公用品采购",
 material_group":OFFICE",
unit_type": EA",
requester": 张三chase_organization":PURCHASE-ORG01,
currency": CNY,
  priority:3urgency": "normal,
  comments:请尽快处理"
}
```

**响应**:
```json[object Object]
  request_id: PR-202401atus:processing",
  "message":Purchase request queued successfully}
```

####12队列中的采购请求状态
```http
GET /purchase-requests/queue/{request_id}
```

**响应**:
```json[object Object]
  request_id: PR-202401atus:processing,
  priority": 3,
  queue_position":2,
enqueued_at: 202401110:00:00Z",
  retry_count:0 error_msg": ",processed_at": null
}
```

#### 1.3列出队列中的采购请求
```http
GET /purchase-requests/queue?limit=10&status=processing
```

**查询参数**:
- `limit`: 每页数量（默认10）
- `offset`: 偏移量（默认0）
- `status`: 状态过滤（可选）

**响应**:
```json
[object Object]requests: [[object Object]
      request_id: R-2024,
      doc_type": PR",
     plant: ,
     material": MAT-001,
    quantity: 10,
      unit_price: 100.5   total_amount": 1005.00
    currency": "CNY",
     vendor_code": "VENDOR-001,  material_group:OFFICE",
    requester": "张三,
      "purchase_organization":PURCHASE-ORG-001,
   priority": 3status": "processing",
      urgency:normal",
      queue_position:1    enqueued_at: 20241T1000      retry_count: 0,
    error_msg: }
  ],
  pending_count": 5,
limit: 10,
  "offset: 0
}
```

####1.4 清空队列
```http
POST /purchase-requests/queue/clear-all
```

**响应**:
```json
{
 message": "Queue cleared successfully}
```

### 2. 采购请求数据库管理

#### 20.1建采购请求（数据库）
```http
POST /purchase-requests
```

**请求体**: 同队列提交接口

**响应**:
```json
{
  "message":Purchase request created successfully,
 data:[object Object]id:1
    request_id:PR-2241,
    doc_type": PR",
   plant": "100,
  quantity": 10
    unit_price": 1000.50
   material": "MAT-01    delivery_date: 20241231T23:59:59Z",
   vendor_code":VENDOR-001,
   short_text: 采购",
   material_group: OFFICE",
  unit_type: A",
  requester": "张三",
    "purchase_organization":PURCHASE-ORG-01,
  currency":CNY
  total_amount": 100500,
    priority": 3status": "pending",
    urgency: al,
    comments":请尽快处理",
    retry_count:0
    created_at: 202411T10:00:0,
    updated_at: 2024-11T10:00Z
  }
}
```

#### 2.2 获取采购请求详情
```http
GET /purchase-requests/{request_id}
```

**响应**:
```json[object Object]
 data:[object Object]id:1
    request_id:PR-2241,
    doc_type": PR",
   plant": "100,
  quantity": 10
    unit_price": 1000.50
   material": "MAT-01    delivery_date: 20241231T23:59:59Z",
   vendor_code":VENDOR-001,
   short_text: 采购",
   material_group: OFFICE",
  unit_type: A",
  requester": "张三",
    "purchase_organization":PURCHASE-ORG-01,
  currency":CNY
  total_amount": 100500,
    priority": 3status": "pending",
    urgency: al,
    comments":请尽快处理",
    retry_count:0
    created_at: 202411T10:00:0,
    updated_at: 2024-11T10:00Z
  }
}
```

#### 2.3 查询采购请求列表
```http
GET /purchase-requests?page=1&page_size=10tatus=pending&requester=张三&sort_by=created_at&sort_order=desc
```

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）
- `status`: 状态过滤
- `requester`: 申请人过滤
- `plant`: 工厂过滤
- `vendor_code`: 供应商过滤
- `material_group`: 物料组过滤
- `priority`: 优先级过滤
- `urgency`: 紧急程度过滤
- `start_date`: 开始日期（YYYY-MM-DD）
- `end_date`: 结束日期（YYYY-MM-DD）
- `sort_by`: 排序字段
- `sort_order`: 排序方向（asc/desc）

**响应**:
```json[object Object]data:[object Object]requests":
     [object Object]
        id,
        request_id: 20241
        doc_type: R",
       plant:1000
      quantity":10       unit_price:100.5
       material": MAT-1
        delivery_date: 202412359,     vendor_code": "VENDOR-001      short_text":办公用品采购",
       material_group":OFFICE",
      unit_type: EA",
      requester: "张三",
purchase_organization":PURCHASE-ORG-001",
      currency": "CNY",
      total_amount": 1005
     priority": 3
        status": pending",
        urgency":normal",
        comments": 请尽快处理",
        retry_count": 0,
        created_at: 2024100,
        updated_at: 2024-1100Z"
      }
    ],
   total: 1  page":1,
   page_size": 10
  }
}
```

#### 20.4更新采购请求
```http
PUT /purchase-requests/{request_id}
```

**请求体**:
```json
{
  "status: pproved",
approver_id":APPROVER-1,
  approver_name":李四,
  comments: 批准，请按计划执行"
}
```

**响应**:
```json
{
  "message":Purchase request updated successfully,
 data:[object Object]id:1
    request_id:PR-202401,
    status:approved",
  approver_id": "APPROVER-01    approver_name": "李四",
  approved_at: 202411T11:00:00,
    comments":已批准，请按计划执行,
    updated_at: 2024-11T11:00Z
  }
}
```

#### 20.5购请求
```http
DELETE /purchase-requests/{request_id}
```

**响应**:
```json
{
  "message":Purchase request deleted successfully}
```

### 3. 系统状态和监控

#### 30.1获取系统状态
```http
GET /status
```

**响应**:
```json
[object Object]queue_stats": {
    pending_count": 5
  },
rpa_status": [object Object] available": true
  }
}
```

####32健康检查
```http
GET /health
```

**响应**:
```json
{
  status": "healthy,components":[object Object]  queue: [object Object]  healthy": true,
     error: null
    },
  rpa: [object Object] healthy: true
    }
  }
}
```

## ⚙️ 配置说明

### config.yaml 配置文件
```yaml
server:
  port: 8082                    # 服务器端口
  read_timeout: 30s             # 读取超时
  write_timeout:30s            # 写入超时

rpa:
  base_url: http://localhost:8084 RPA系统地址
  timeout: 60s                       # RPA请求超时
  check_interval: 60s                # RPA状态检查间隔

queue:
  type: "memory"                     # 队列类型: memory 或 redis
  check_interval: 5                 # 队列检查间隔
  max_retries: 3                     # 最大重试次数
  retry_interval: 10s                # 重试间隔

redis:
  addr: "localhost:6379           # Redis地址
  password: ""                       # Redis密码
  db: 0                             # Redis数据库
  pool_size: 10                     # 连接池大小
  min_idle_conns:5                 # 最小空闲连接

mysql:
  host: "localhost                 # MySQL主机
  port: 3306                        # MySQL端口
  username: "root                 # 用户名
  password: "Test123ls             # 密码
  database: "middleware            # 数据库名
  charset: "utf8mb4                # 字符集
  parse_time: true                  # 解析时间
  loc: "Local"                      # 时区
  max_open_conns:25                # 最大连接数
  max_idle_conns:5                 # 最大空闲连接
  conn_max_lifetime:30s           # 连接最大生命周期

logger:
  level: "info"                     # 日志级别
  format: "json"                    # 日志格式
```

## 🚀 部署和运行

### 环境要求
- Go1.21MySQL80+
- Redis 7选，用于生产环境)
- Docker (可选)

### 快速启动

#### 1. 使用内存队列（开发模式）
```bash
# 构建项目
make build

# 运行服务
make run
```

####2. 使用Redis队列（生产模式）
```bash
# 启动Redis
docker run -d -p6379379 redis:7-alpine

# 运行服务
make run-redis
```

#### 3. 使用Docker Compose
```bash
docker-compose up -d
```

### 数据库初始化
```bash
# 创建数据库和表
mysql -u root -p < create_database.sql

# 或者使用初始化脚本
mysql -u root -p < init_database.sql
```

## 🧪 测试

### 运行测试
```bash
# 单元测试
make test

# 快速功能测试
make test-quick

# 全面测试
make test-full
```

### 使用RPA模拟服务
```bash
# 启动RPA模拟服务
go run cmd/rpa_mock/main.go

# 在另一个终端启动中间件
make run
```

## 📊 监控和日志

### 日志级别
- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 监控指标
- 队列长度和状态
- RPA系统可用性
- 请求处理成功率
- 平均处理时间
- 错误率和重试次数

## 🔒 安全考虑

- 所有API请求都记录在结构化日志中
- 支持CORS配置
- 建议在生产环境中添加认证和授权
- 敏感配置应通过环境变量注入
- 数据库连接使用连接池管理

## 📈 性能特性

- **高并发处理**: 支持大量并发采购请求
- **智能队列**: 基于优先级的任务调度
- **连接池**: MySQL和Redis连接池优化
- **异步处理**: 队列处理器异步处理任务
- **状态缓存**: 队列状态实时更新
- **错误恢复**: 自动重试和错误处理机制

## 🔄 工作流程

### 1. 采购请求提交流程
```
客户端 → POST /purchase-requests/queue → 队列管理器 → 队列存储 → 数据库持久化
```

### 2采购请求处理流程
```
队列处理器 → 检查RPA可用性 → 从队列取请求 → 发送给RPA系统 → 更新状态 → 通知回调
```

### 3. 状态管理流程
```
pending → processing → completed/failed
         ↓
    retry (最多3)
```

## 🤝 扩展性

### 添加新的队列类型
1 `interfaces.QueueManager` 接口
2. 在 `internal/queue/factory.go` 中添加新的队列类型3. 更新配置结构以支持新队列类型

### 添加新的通知方式
1faces.NotificationService` 接口
2在处理器中集成新的通知服务

### 添加新的RPA系统1. 实现 `interfaces.RPAClient` 接口2新RPA客户端工厂方法

---

**注意**: 这是一个生产就绪的RPA中间件系统，专注于采购请求的自动化处理。系统支持高并发、高可用性和可扩展性，建议在生产环境中使用Redis队列以确保数据持久化。 