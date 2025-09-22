# 部署启动文档

## 概述

本文档描述了RPA中间件系统的部署和启动流程，包括环境配置、依赖安装、数据库设置、服务启动等步骤。

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端 (Vue3)   │    │   后端 (Go)     │    │   数据库 (MySQL)│
│   Port: 3000    │◄──►│   Port: 8082    │◄──►│   Port: 3306    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  队列 (Redis)    │
                       │   Port: 6379    │
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  RPA (UiPath)   │
                       │   Port: 8084    │
                       └─────────────────┘
```

## 环境要求

### 系统要求
- **操作系统**: Windows 10/11, Linux (Ubuntu 18.04+), macOS 10.15+
- **内存**: 最少4GB，推荐8GB+
- **磁盘空间**: 最少10GB可用空间
- **网络**: 稳定的网络连接

### 软件依赖

#### 后端依赖
- **Go**: 1.19+ (推荐1.21+)
- **MySQL**: 8.0+ (或5.7+)
- **Redis**: 6.0+ (可选，用于队列管理)

#### 前端依赖
- **Node.js**: 16.0+ (推荐18.0+)
- **npm**: 8.0+ 或 **yarn**: 1.22+

#### 其他工具
- **Git**: 2.0+
- **Make**: 3.81+ (Windows用户需要安装)

## 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd testMW
```

### 2. 配置环境
```bash
# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置文件
nano config.yaml
```

### 3. 启动服务
```bash
# 使用Makefile启动所有服务
make start-all

# 或者分别启动
make start-backend
make start-frontend
```

## 详细部署步骤

### 第一步：环境准备

#### 1.1 安装Go
```bash
# 下载Go (以1.21.5为例)
wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# 设置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 1.2 安装Node.js
```bash
# 使用NodeSource仓库安装Node.js 18
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证安装
node --version
npm --version
```

#### 1.3 安装MySQL
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server

# 启动MySQL服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation
```

#### 1.4 安装Redis (可选)
```bash
# Ubuntu/Debian
sudo apt install redis-server

# 启动Redis服务
sudo systemctl start redis-server
sudo systemctl enable redis-server

# 验证安装
redis-cli ping
```

### 第二步：数据库配置

#### 2.1 创建数据库
```sql
-- 登录MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE middleware CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户 (可选)
CREATE USER 'middleware_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON middleware.* TO 'middleware_user'@'localhost';
FLUSH PRIVILEGES;
```

#### 2.2 初始化数据库表
```bash
# 运行数据库初始化脚本
mysql -u root -p middleware < init_database.sql

# 创建额外表 (如果需要)
mysql -u root -p middleware < create_additional_tables.sql

# 插入测试数据 (可选)
mysql -u root -p middleware < insert_mock_data.sql
```

### 第三步：后端部署

#### 3.1 安装Go依赖
```bash
# 进入项目目录
cd testMW

# 下载依赖
go mod download
go mod tidy
```

#### 3.2 配置后端
```bash
# 编辑配置文件
nano config.yaml
```

**config.yaml 示例配置:**
```yaml
server:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s

rpa:
  base_url: "http://localhost:8084"
  timeout: 60s
  check_interval: 60s

uipath:
  orch_base_url: "https://your-uipath-orchestrator.com"
  tenancy_name: "Default"
  username: "your_username"
  password: "your_password"
  folder_id: 31
  queue_name: "CreationPR"
  verify_ssl: false
  timeout: 60s

queue:
  type: "memory"  # 或 "redis"
  check_interval: 5s
  max_retries: 3
  retry_interval: 10s

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5

mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_mysql_password"
  database: "middleware"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 300s

logger:
  level: "info"
  format: "json"
```

#### 3.3 构建后端
```bash
# 构建所有后端服务
make build-backend

# 或者单独构建
go build -o build/server.exe cmd/server/main.go
go build -o build/full_server.exe cmd/full_server/main.go
go build -o build/quick_test.exe cmd/quick_test/main.go
```

#### 3.4 启动后端服务
```bash
# 启动完整服务器 (推荐)
./build/full_server.exe

# 或者启动基础服务器
./build/server.exe

# 后台运行
nohup ./build/full_server.exe > server.log 2>&1 &
```

### 第四步：前端部署

#### 4.1 安装前端依赖
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 或者使用yarn
yarn install
```

#### 4.2 配置前端
```bash
# 编辑环境配置
nano .env.local
```

**.env.local 示例配置:**
```env
VITE_API_BASE_URL=http://localhost:8082
VITE_APP_TITLE=RPA中间件系统
VITE_APP_VERSION=1.0.0
```

#### 4.3 构建前端
```bash
# 开发环境构建
npm run build:dev

# 生产环境构建
npm run build:prod

# 或者使用yarn
yarn build:dev
yarn build:prod
```

#### 4.4 启动前端服务
```bash
# 开发模式
npm run dev

# 生产模式
npm run preview

# 或者使用yarn
yarn dev
yarn preview
```

### 第五步：验证部署

#### 5.1 检查后端服务
```bash
# 检查健康状态
curl http://localhost:8082/api/v1/health

# 检查队列状态
curl http://localhost:8082/api/v1/status

# 检查API文档
curl http://localhost:8082/api/v1/purchase-requests
```

#### 5.2 检查前端服务
```bash
# 访问前端页面
curl http://localhost:3000

# 检查API连接
curl http://localhost:3000/api/v1/health
```

#### 5.3 检查数据库连接
```bash
# 使用测试工具
go run cmd/check_database/main.go

# 或者直接查询
mysql -u root -p -e "SELECT COUNT(*) FROM purchase_requests;" middleware
```

## 生产环境部署

### Docker部署 (推荐)

#### 1. 创建Dockerfile
```dockerfile
# 后端Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/full_server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .
EXPOSE 8082
CMD ["./server"]
```

#### 2. 创建docker-compose.yml
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: middleware
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init_database.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: .
    ports:
      - "8082:8082"
    depends_on:
      - mysql
      - redis
    environment:
      - MYSQL_HOST=mysql
      - REDIS_ADDR=redis:6379

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  mysql_data:
```

#### 3. 启动服务
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 系统服务部署

#### 1. 创建systemd服务文件
```bash
# 创建后端服务文件
sudo nano /etc/systemd/system/rpa-middleware.service
```

**rpa-middleware.service:**
```ini
[Unit]
Description=RPA Middleware Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/rpa-middleware
ExecStart=/opt/rpa-middleware/build/full_server.exe
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

#### 2. 启动服务
```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start rpa-middleware

# 设置开机自启
sudo systemctl enable rpa-middleware

# 查看状态
sudo systemctl status rpa-middleware
```

### Nginx反向代理

#### 1. 安装Nginx
```bash
sudo apt install nginx
```

#### 2. 配置Nginx
```bash
sudo nano /etc/nginx/sites-available/rpa-middleware
```

**Nginx配置:**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/rpa-middleware/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://localhost:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### 3. 启用配置
```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/rpa-middleware /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启Nginx
sudo systemctl restart nginx
```

## 监控和维护

### 日志管理
```bash
# 查看应用日志
tail -f server.log

# 查看系统日志
journalctl -u rpa-middleware -f

# 日志轮转配置
sudo nano /etc/logrotate.d/rpa-middleware
```

### 性能监控
```bash
# 监控系统资源
htop

# 监控数据库
mysqladmin -u root -p status

# 监控Redis
redis-cli info
```

### 备份策略
```bash
# 数据库备份
mysqldump -u root -p middleware > backup_$(date +%Y%m%d).sql

# 应用备份
tar -czf app_backup_$(date +%Y%m%d).tar.gz /opt/rpa-middleware

# 自动备份脚本
crontab -e
# 添加: 0 2 * * * /path/to/backup_script.sh
```

## 故障排除

### 常见问题

#### 1. 端口冲突
```bash
# 检查端口占用
netstat -tlnp | grep :8082
lsof -i :8082

# 杀死占用进程
sudo kill -9 <PID>
```

#### 2. 数据库连接失败
```bash
# 检查MySQL服务
sudo systemctl status mysql

# 检查连接
mysql -u root -p -h localhost

# 检查防火墙
sudo ufw status
```

#### 3. 前端构建失败
```bash
# 清理缓存
npm cache clean --force
rm -rf node_modules package-lock.json
npm install

# 检查Node.js版本
node --version
```

#### 4. 内存不足
```bash
# 检查内存使用
free -h

# 增加交换空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 调试模式

#### 1. 启用调试日志
```yaml
# config.yaml
logger:
  level: "debug"
  format: "text"
```

#### 2. 使用调试工具
```bash
# 启动调试服务器
go run cmd/debug_api/main.go

# 运行测试
go test ./...

# 性能分析
go run cmd/quick_test/main.go
```

## 更新和维护

### 更新应用
```bash
# 拉取最新代码
git pull origin main

# 重新构建
make build-all

# 重启服务
sudo systemctl restart rpa-middleware
```

### 数据库迁移
```bash
# 运行迁移脚本
mysql -u root -p middleware < update_database.sql

# 清理测试数据
mysql -u root -p middleware < cleanup_test_data.sql
```

### 配置更新
```bash
# 更新配置后重启服务
sudo systemctl restart rpa-middleware

# 验证配置
curl http://localhost:8082/api/v1/health
```

## 安全建议

### 1. 网络安全
- 使用HTTPS证书
- 配置防火墙规则
- 限制数据库访问IP

### 2. 应用安全
- 定期更新依赖
- 使用强密码
- 启用访问日志

### 3. 数据安全
- 定期备份数据
- 加密敏感信息
- 监控异常访问

## 联系支持

如果遇到问题，请提供以下信息：
- 系统环境信息
- 错误日志
- 配置文件
- 复现步骤

**技术支持邮箱**: support@example.com
**文档更新**: 2024-01-01
