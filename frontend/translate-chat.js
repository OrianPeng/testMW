import fs from 'fs';

// 读取智能助手页面文件
let content = fs.readFileSync('src/views/Chat.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 页面标题
  '智能助手': 'AI Assistant',
  '智能采购助手，帮您处理采购相关事务': 'Intelligent Procurement Assistant to Help You Handle Procurement Matters',
  
  // 快捷操作
  '快捷操作': 'Quick Actions',
  '创建PR': 'Create PR',
  '查询PO状态': 'Query PO Status',
  '查看流程卡点': 'View Process Bottlenecks',
  '发送PO给供应商': 'Send PO to Supplier',
  
  // 聊天界面
  '请输入您的问题...': 'Please enter your question...',
  '发送': 'Send',
  '上传文件': 'Upload File',
  '新对话': 'New Conversation',
  
  // 消息类型
  '文本消息': 'Text Message',
  '表格数据': 'Table Data',
  '卡片信息': 'Card Information',
  '流程步骤': 'Process Steps',
  
  // 表格列标题
  'PR编号': 'PR Number',
  '状态': 'Status',
  '创建时间': 'Created Time',
  '负责人': 'Responsible',
  
  // 流程步骤
  '提交申请': 'Submit Application',
  '部门审批': 'Department Approval',
  '财务审批': 'Finance Approval',
  '采购执行': 'Procurement Execution',
  '申请人': 'Applicant',
  
  // 状态
  '已完成': 'Completed',
  '进行中': 'Current',
  '待处理': 'Pending',
  
  // 操作按钮
  '查看详情': 'View Details',
  '复制': 'Copy',
  '重新生成': 'Regenerate',
  '点赞': 'Like',
  '点踩': 'Dislike',
  
  // 文件上传
  '点击或拖拽文件到此处上传': 'Click or drag files here to upload',
  '支持Excel、PDF等格式': 'Supports Excel, PDF and other formats',
  '文件大小不超过10MB': 'File size should not exceed 10MB',
  
  // 消息提示
  '发送失败': 'Send failed',
  '文件上传失败': 'File upload failed',
  '文件格式不支持': 'File format not supported',
  '文件大小超出限制': 'File size exceeds limit',
  '网络错误': 'Network error',
  '请稍后重试': 'Please try again later',
  
  // 空状态
  '暂无对话记录': 'No conversation history',
  '开始新的对话': 'Start a new conversation',
  
  // 加载状态
  '正在思考...': 'Thinking...',
  '正在处理...': 'Processing...',
  '正在生成回复...': 'Generating response...',
  
  // 其他
  '今天': 'Today',
  '昨天': 'Yesterday',
  '更早': 'Earlier',
  '系统': 'System',
  '用户': 'User',
  '助手': 'Assistant'
};

// 执行翻译
Object.keys(translations).forEach(chinese => {
  const english = translations[chinese];
  const regex = new RegExp(chinese.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
  content = content.replace(regex, english);
});

// 写回文件
fs.writeFileSync('src/views/Chat.vue', content, 'utf8');

console.log('Chat page translation completed!');
