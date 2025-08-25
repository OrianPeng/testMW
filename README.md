# RPA Middleware

一个高性能的RPA（机器人流程自动化）中间件系统，用于管理AI Agent与RPA系统之间的任务队列和通信。

## 📋 项目概述

RPA Middleware 是一个用 Go 语言开发的中间件服务，主要功能包括：

- **任务队列管理**：支持内存队列和Redis队列两种模式
- **智能任务调度**：根据RPA系统状态自动选择直接执行或排队等待
- **高可用性**：支持任务重试、错误处理和状态监控
- **RESTful API**：提供完整的HTTP API接口
- **实时通知**：支持任务完成后的回调通知
- **健康监控**：内置健康检查和状态监控

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
```

### 核心组件

- **API Layer** (`internal/api/`): HTTP API接口和路由管理
- **Queue Manager** (`internal/queue/`): 队列管理，支持内存和Redis
- **Processor** (`internal/processor/`): 队列处理器，负责任务调度和执行
- **RPA Client** (`internal/rpa/`): RPA系统客户端
- **Notification** (`internal/notification/`): 通知服务
- **Config** (`internal/config/`): 配置管理

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Redis (可选，用于持久化队列)
- Docker (可选)

### 安装

1. **克隆项目**
```bash
git clone <repository-url>
cd rpa-middleware
```

2. **安装依赖**
```bash
go mod download
```

3. **构建项目**
```bash
make build
```

### 配置

项目使用 `config.yaml` 进行配置：

```yaml
server:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s

rpa:
  base_url: "http://localhost:8084"
  timeout: 60s
  check_interval: 60s

queue:
  type: "redis"  # memory 或 redis
  check_interval: 5s
  max_retries: 3
  retry_interval: 10s

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
```

### 运行

#### 使用内存队列（开发模式）
```bash
make run
```

#### 使用Redis队列（生产模式）
```bash
# 启动Redis
docker run -d -p 6379:6379 redis:7-alpine

# 运行中间件
make run-redis
```

#### 使用Docker Compose
```bash
docker-compose up -d
```

## 📚 API 文档

### 基础URL
```
http://localhost:8082/api/v1
```

### 主要接口

#### 1. 提交任务请求
```http
POST /requests
Content-Type: application/json

{
  "id": "req-123",
  "agent_id": "agent-001",
  "data": {
    "action": "process_document",
    "document_id": "doc-123"
  },
  "priority": 5,
  "callback": "http://agent/callback"
}
```

**响应示例：**
```json
{
  "request_id": "req-123",
  "status": "pending",
  "message": "Request queued successfully"
}
```

#### 2. 查询任务状态
```http
GET /requests/{request_id}
```

**响应示例：**
```json
{
  "request_id": "req-123",
  "status": "completed",
  "message": "Request processed successfully",
  "data": {
    "result": "document processed"
  }
}
```

#### 3. 列出所有任务
```http
GET /requests?limit=10&offset=0
```

#### 4. 系统状态
```http
GET /status
```

#### 5. 健康检查
```http
GET /health
```

## 🔧 开发指南

### 项目结构
```
rpa-middleware/
├── cmd/                    # 可执行文件
│   ├── server/            # 主服务器
│   ├── rpa_mock/          # RPA模拟服务
│   └── redis_test/        # Redis测试工具
├── internal/              # 内部包
│   ├── api/              # API层
│   ├── config/           # 配置管理
│   ├── interfaces/       # 接口定义
│   ├── models/           # 数据模型
│   ├── notification/     # 通知服务
│   ├── processor/        # 队列处理器
│   ├── queue/            # 队列管理
│   └── rpa/              # RPA客户端
├── examples/             # 示例代码
├── config.yaml          # 配置文件
├── docker-compose.yml   # Docker编排
└── Makefile            # 构建脚本
```

### 常用命令

```bash
# 构建
make build

# 运行
make run

# 测试
make test

# 代码格式化
make format

# 清理
make clean

# 开发模式（需要安装air）
make dev
```

### 添加新的队列类型

1. 实现 `interfaces.QueueManager` 接口
2. 在 `internal/queue/factory.go` 中添加新的队列类型
3. 更新配置结构以支持新队列类型

### 添加新的通知方式

1. 实现 `interfaces.NotificationService` 接口
2. 在处理器中集成新的通知服务

## 🧪 测试

### 运行测试
```bash
# 单元测试
make test

# 测试覆盖率
make test-coverage

# 快速功能测试
make test-quick

# 全面测试
make test-full
```

### 使用RPA模拟服务

项目包含一个RPA模拟服务，用于测试：

```bash
# 启动RPA模拟服务
go run cmd/rpa_mock/main.go

# 在另一个终端启动中间件
make run
```

### 测试示例

参考 `examples/client_example.go` 了解如何使用API：

```bash
go run examples/client_example.go
```

## 📊 监控和日志

### 日志配置

系统使用结构化JSON日志，支持以下级别：
- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 监控指标

- 队列长度
- 处理成功率
- 平均处理时间
- RPA系统状态
- 错误率

## 🔒 安全考虑

- 所有API请求都记录在日志中
- 支持CORS配置
- 建议在生产环境中添加认证和授权
- 敏感配置应通过环境变量注入

## 🚀 部署

### Docker部署
```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run

# 使用Redis的完整部署
docker-compose up -d
```

### 生产环境建议

1. 使用Redis作为队列存储
2. 配置适当的日志级别
3. 设置监控和告警
4. 使用负载均衡器
5. 配置健康检查
6. 设置资源限制

## 🤝 贡献

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

[许可证信息]

## 📞 支持

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件
- 项目讨论区

---

**注意**: 这是一个生产就绪的RPA中间件系统，支持高并发、高可用性和可扩展性。建议在生产环境中使用Redis队列以确保数据持久化。