# 订单中心前端

基于 Vue 3 + TypeScript + Vite + Ant Design Vue 构建的现代化订单中心前端应用。

## 功能特性

### 1. 智能助手
- 集成 Dify 聊天机器人
- 支持自然语言交互
- 预定义快捷操作
- 支持文件上传
- 多种消息类型渲染（文本、表格、卡片、流程图）

### 2. PR查询管理
- 采购请求（PR）的创建、查询、编辑
- 高级筛选和搜索
- 状态管理和审批流程
- 详情查看和批量操作

### 3. PO查询管理
- 采购订单（PO）的查询和管理
- 供应商发送功能
- 状态跟踪和确认
- 导出和统计功能

### 4. 流程监控
- 实时流程状态监控
- 阻塞流程识别和处理
- 审批时长统计
- 流程步骤可视化

### 5. 供应商管理
- 供应商信息管理
- 联系人信息维护
- 采购品类分类
- 状态管理和权限控制

### 6. 数据看板
- 关键指标展示
- 趋势图表分析
- 实时数据更新
- 多维度统计

## 技术栈

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI组件库**: Ant Design Vue 4.x
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **图表**: ECharts
- **HTTP客户端**: Axios
- **时间处理**: Day.js
- **样式**: CSS3 + Ant Design 主题

## 项目结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API接口
│   ├── components/        # 公共组件
│   ├── layouts/          # 布局组件
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理
│   ├── styles/           # 全局样式
│   ├── types/            # TypeScript类型定义
│   ├── views/            # 页面组件
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── index.html            # HTML模板
├── package.json          # 依赖配置
├── tsconfig.json         # TypeScript配置
├── vite.config.ts        # Vite配置
└── README.md            # 项目说明
```

## 开发指南

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0 或 yarn >= 1.22.0

### 安装依赖

```bash
cd frontend
npm install
```

### 开发环境启动

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 代码检查

```bash
npm run lint
```

## 配置说明

### 环境变量

创建 `.env.local` 文件：

```env
VITE_API_BASE_URL=http://localhost:8082
VITE_DIFY_API_URL=your_dify_api_url
VITE_DIFY_API_KEY=your_dify_api_key
```

### 代理配置

开发环境下，API请求会自动代理到后端服务（localhost:8082）。

### 主题定制

可以通过修改 `src/styles/index.css` 文件来定制主题样式。

## 功能模块

### 权限控制

系统支持基于角色的权限控制：

- **管理员**: 所有权限
- **采购员**: PR/PO创建、查询、供应商管理
- **审批员**: PR审批、流程监控
- **普通用户**: 基础查询功能

### 响应式设计

- 支持桌面端和移动端
- 自适应布局
- 触摸友好的交互

### 国际化支持

- 支持中英文切换
- 可扩展多语言支持

## 部署说明

### Docker部署

```bash
# 构建镜像
docker build -t order-center-frontend .

# 运行容器
docker run -p 3000:80 order-center-frontend
```

### Nginx配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 开发规范

### 代码规范

- 使用 TypeScript 进行类型检查
- 遵循 Vue 3 Composition API 规范
- 使用 ESLint 进行代码检查
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case

### 提交规范

使用 Conventional Commits 规范：

```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建过程或辅助工具的变动
```

## 常见问题

### 1. 开发环境启动失败

检查 Node.js 版本是否符合要求，清除 node_modules 重新安装：

```bash
rm -rf node_modules package-lock.json
npm install
```

### 2. API请求失败

检查后端服务是否启动，确认代理配置是否正确。

### 3. 样式问题

确认 Ant Design Vue 版本兼容性，检查全局样式是否正确引入。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
