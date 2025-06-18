# RPA 中间件快速启动指南

## 🚀 快速开始

### 1. 环境准备

确保已安装 Go 1.21 或更高版本：
```bash
go version
```

### 2. 获取代码

```bash
# 克隆代码（如果从 Git 仓库）
git clone <repository-url>
cd rpa-middleware

# 或者直接在现有目录中初始化
go mod init rpa-middleware
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置环境

复制环境变量示例文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件，配置 RPA 系统地址：
```bash
RPA_BASE_URL=http://your-rpa-system:8081
```

### 5. 启动服务

```bash
# 方法1: 直接运行
go run cmd/server/main.go

# 方法2: 使用 Makefile
make run

# 方法3: 构建后运行
make build
./build/rpa-middleware
```

### 6. 验证服务

打开浏览器或使用 curl 访问：

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 获取系统状态
curl http://localhost:8080/api/v1/status
```

预期响应：
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

## 📝 API 使用示例

### 提交请求

```bash
curl -X POST http://localhost:8080/api/v1/requests \
  -H "Content-Type: application/json" \
  -d '{
    "id": "req-001",
    "agent_id": "ai-agent-001",
    "data": {
      "action": "process_document",
      "document_id": "doc-123"
    },
    "priority": 5,
    "callback": "http://ai-agent/callback"
  }'
```

### 查询请求状态

```bash
curl http://localhost:8080/api/v1/requests/req-001
```

### 列出所有请求

```bash
# 列出所有请求
curl http://localhost:8080/api/v1/requests

# 按状态过滤
curl "http://localhost:8080/api/v1/requests?status=pending&limit=10"
```

## 🔧 使用客户端示例

运行客户端示例程序：

```bash
go run examples/client_example.go
```

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t rpa-middleware:latest .
```

### 运行容器

```bash
docker run -d \
  --name rpa-middleware \
  -p 8080:8080 \
  -e RPA_BASE_URL=http://your-rpa-system:8081 \
  rpa-middleware:latest
```

### 使用 Docker Compose

创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  rpa-middleware:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RPA_BASE_URL=http://rpa-system:8081
    restart: unless-stopped
```

启动：
```bash
docker-compose up -d
```

## 🛠️ 开发模式

### 热重载开发

安装 air：
```bash
go install github.com/cosmtrek/air@latest
```

启动开发模式：
```bash
make dev
# 或者
air
```

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make format
```

## 📊 监控和日志

### 查看日志

如果使用 Docker：
```bash
docker logs -f rpa-middleware
```

### 监控指标

访问系统状态接口：
```bash
curl http://localhost:8080/api/v1/status
```

## 🔧 故障排除

### 常见问题

1. **服务启动失败**
   - 检查端口 8080 是否被占用
   - 验证 Go 版本是否正确

2. **RPA 连接失败**
   - 确认 RPA_BASE_URL 配置正确
   - 检查 RPA 系统是否运行
   - 验证网络连接

3. **请求处理异常**
   - 查看日志输出
   - 检查 RPA 系统状态
   - 验证请求数据格式

### 日志级别调整

修改配置文件或环境变量：
```yaml
logger:
  level: "debug"  # trace, debug, info, warn, error
```

## 🔗 相关链接

- [完整文档](README.md)
- [API 文档](docs/api.md)
- [配置文档](docs/configuration.md)
- [部署文档](docs/deployment.md)

## 💡 下一步

1. 根据实际需求调整配置
2. 集成您的 AI Agent 和 RPA 系统
3. 设置监控和告警
4. 考虑生产环境部署策略