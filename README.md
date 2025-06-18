# RPA 中间件

这是一个用于处理 AI Agent 和 RPA 系统之间请求调度的中间件服务。当 RPA 系统繁忙时，中间件会将请求缓存在队列中，等待 RPA 系统可用时再进行处理。

## 功能特性

- 🔄 **请求队列管理**：支持基于优先级的请求队列
- 📊 **RPA 状态监控**：实时监控 RPA 系统可用性
- 🔁 **自动重试机制**：失败请求自动重试
- 📞 **回调通知**：支持异步回调通知
- 🔍 **请求追踪**：完整的请求生命周期管理
- 🎯 **RESTful API**：清晰的 API 接口设计

## 架构设计

```
AI Agent → 中间件 → RPA 系统
          ↓
        请求队列
```

### 核心组件

1. **Queue Manager**：队列管理器，负责请求的入队、出队和状态管理
2. **RPA Client**：RPA 系统客户端，负责与 RPA 系统通信
3. **Queue Processor**：队列处理器，负责监控和处理队列中的请求
4. **Notification Service**：通知服务，负责发送回调通知
5. **API Handler**：HTTP API 处理器，提供 RESTful 接口

## API 接口

### 提交请求
```bash
POST /api/v1/requests
Content-Type: application/json

{
  "id": "request-123",
  "agent_id": "agent-001",
  "data": {
    "action": "process_document",
    "params": {...}
  },
  "priority": 5,
  "callback": "http://ai-agent/callback"
}
```

### 查询请求状态
```bash
GET /api/v1/requests/{request_id}
```

### 列出请求
```bash
GET /api/v1/requests?status=pending&limit=20&offset=0
```

### 获取系统状态
```bash
GET /api/v1/status
```

### 健康检查
```bash
GET /api/v1/health
```

## 数据模型

### 请求状态
- `pending`: 等待处理
- `processing`: 正在处理
- `completed`: 已完成
- `failed`: 处理失败

### 请求优先级
数值越大优先级越高，默认为 0。

## 运行项目

### 本地运行

1. 安装依赖
```bash
go mod tidy
```

2. 运行服务
```bash
go run cmd/server/main.go
```

3. 或者编译后运行
```bash
go build -o rpa-middleware cmd/server/main.go
./rpa-middleware
```

### 环境变量配置

- `RPA_BASE_URL`: RPA 系统基础 URL (默认: http://localhost:8081)
- `SERVER_PORT`: 服务端口 (默认: 8080)

### Docker 运行

可以创建 Dockerfile 进行容器化部署。

## 配置说明

配置文件 `config.yaml` 支持以下参数：

```yaml
server:
  port: 8080                # HTTP 服务端口
  read_timeout: 30s         # 读超时
  write_timeout: 30s        # 写超时

rpa:
  base_url: "http://localhost:8081"  # RPA 系统 URL
  timeout: 60s              # 请求超时
  check_interval: 5s        # 状态检查间隔

queue:
  check_interval: 5s        # 队列检查间隔
  max_retries: 3            # 最大重试次数
  retry_interval: 10s       # 重试间隔

logger:
  level: "info"             # 日志级别
  format: "json"            # 日志格式
```

## 与外部系统集成

### AI Agent 集成

AI Agent 需要调用中间件的 `/api/v1/requests` 接口提交请求，并可选择提供回调 URL 接收处理结果。

### RPA 系统集成

RPA 系统需要提供以下接口：

1. **执行接口**: `POST /api/v1/execute`
   ```json
   {
     "id": "request-123",
     "agent_id": "agent-001",
     "data": {...},
     "metadata": {...}
   }
   ```

2. **状态接口**: `GET /api/v1/status`
   ```json
   {
     "is_available": true,
     "current_task": "task-456"
   }
   ```

## 监控和日志

- 使用结构化日志记录所有关键操作
- 支持请求追踪和性能监控
- 提供健康检查和状态查询接口

## 扩展性

- 队列管理器可扩展为 Redis 或数据库实现
- 支持多实例部署和负载均衡
- 可集成消息队列系统（如 RabbitMQ、Kafka）