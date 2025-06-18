.PHONY: build run test clean deps help

# 构建配置
BINARY_NAME=rpa-middleware
MAIN_PATH=cmd/server/main.go
BUILD_DIR=build

# Go 配置
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## 显示帮助信息
	@echo "可用命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	$(GOMOD) download
	$(GOMOD) tidy

build: deps ## 构建项目
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

run: ## 运行项目
	$(GOCMD) run $(MAIN_PATH)

test: ## 运行测试
	$(GOTEST) -v -race ./...

test-coverage: ## 运行测试并生成覆盖率报告
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

clean: ## 清理构建文件
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

format: ## 格式化代码
	$(GOCMD) fmt ./...

lint: ## 代码检查 (需要安装 golint)
	golint ./...

docker-build: ## 构建 Docker 镜像
	docker build -t $(BINARY_NAME):latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

install: build ## 安装到系统
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

dev: ## 开发模式运行（带自动重载）
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请安装 air: go install github.com/cosmtrek/air@latest"; \
		$(MAKE) run; \
	fi