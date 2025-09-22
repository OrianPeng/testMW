# AI Assistant Integration

## 概述

AI助手功能已成功集成到前端界面中，使用iframe方式嵌入外部聊天机器人。

## 功能特性

### 🎯 主要功能
- **嵌入式聊天机器人**: 使用iframe集成外部AI聊天机器人
- **快速操作按钮**: 提供常用的业务操作快捷入口
- **响应式设计**: 支持桌面和移动设备
- **实时刷新**: 支持聊天界面刷新功能

### 🚀 快速操作
- **Create PR**: 创建采购请求
- **Search PR Status**: 查询采购请求状态
- **Query PO Status**: 查询采购订单状态
- **Process Bottleneck Search**: 查看流程瓶颈

## 技术实现

### 文件结构
```
frontend/src/views/Chat.vue          # 主要的AI助手页面
frontend/test-iframe.html            # iframe测试页面
frontend/AI_ASSISTANT_INTEGRATION.md # 本文档
```

### 核心组件
- **Chat.vue**: 主要的AI助手界面组件
- **CreatePRModal.vue**: 创建采购请求的模态框
- **iframe集成**: 嵌入外部聊天机器人

### 技术特点
- **Vue 3 Composition API**: 使用最新的Vue 3语法
- **TypeScript支持**: 完整的类型安全
- **Ant Design Vue**: 现代化的UI组件库
- **响应式布局**: 适配不同屏幕尺寸

## 使用方法

### 1. 访问AI助手
- 启动前端应用后，访问 `/chat` 路由
- 或者直接访问首页，会自动重定向到聊天页面

### 2. 使用快速操作
- 点击顶部的快速操作按钮
- 系统会自动向聊天机器人发送相应的提示信息

### 3. 创建采购请求
- 点击"Create PR"按钮
- 填写采购请求信息
- 创建成功后会在聊天中显示确认信息

## 配置说明

### iframe配置
```html
<iframe
  src="http://localhost/chatbot/HHOn2vBJZJHHTSug"
  style="width: 100%; height: 100%; min-height: 700px"
  frameborder="0"
  allow="microphone"
/>
```

### 快速操作配置
```javascript
const quickActions = [
  {
    key: 'create-pr',
    label: 'Create PR',
    icon: h(FileTextOutlined),
    type: 'primary',
    prompt: 'Please help me create a purchase request...'
  },
  // ... 其他操作
]
```

## 测试

### 1. 本地测试
```bash
# 启动前端开发服务器
cd frontend
npm run dev

# 访问测试页面
open http://localhost:3000/chat
```

### 2. iframe测试
```bash
# 直接打开测试页面
open frontend/test-iframe.html
```

## 故障排除

### 常见问题

1. **iframe无法加载**
   - 检查聊天机器人URL是否正确
   - 确认网络连接正常
   - 查看浏览器控制台错误信息

2. **快速操作无响应**
   - 检查iframe是否完全加载
   - 确认postMessage通信是否正常
   - 查看浏览器控制台错误信息

3. **样式显示异常**
   - 检查CSS样式是否正确应用
   - 确认响应式断点设置
   - 验证浏览器兼容性

### 调试方法

1. **浏览器开发者工具**
   - 打开F12开发者工具
   - 查看Console标签页的错误信息
   - 检查Network标签页的网络请求

2. **Vue DevTools**
   - 安装Vue DevTools浏览器扩展
   - 查看组件状态和props
   - 监控事件触发情况

## 未来改进

### 计划功能
- [ ] 支持更多快速操作
- [ ] 添加聊天历史记录
- [ ] 实现消息通知功能
- [ ] 支持文件上传到聊天机器人
- [ ] 添加语音输入支持

### 技术优化
- [ ] 优化iframe加载性能
- [ ] 添加错误重试机制
- [ ] 实现消息队列管理
- [ ] 添加用户偏好设置

## 联系支持

如有问题或建议，请联系开发团队或查看项目文档。

