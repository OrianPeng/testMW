import fs from 'fs';

// 读取流程监控页面文件
let content = fs.readFileSync('src/views/ProcessMonitor.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 页面标题
  '流程监控': 'Process Monitor',
  '流程状态监控和瓶颈分析': 'Process Status Monitoring and Bottleneck Analysis',
  
  // 概览统计
  '今日待处理': 'Today Pending',
  '逾期流程': 'Overdue Processes',
  '平均审批时长': 'Average Approval Time',
  '完成率': 'Completion Rate',
  '小时': 'hours',
  '%': '%',
  
  // 流程步骤
  'PR流程': 'PR Process',
  'PO流程': 'PO Process',
  '提交申请': 'Submit Application',
  '部门审批': 'Department Approval',
  '财务审批': 'Finance Approval',
  '采购执行': 'Procurement Execution',
  '供应商确认': 'Supplier Confirmation',
  '质量检验': 'Quality Inspection',
  '交付完成': 'Delivery Complete',
  
  // 状态
  '已完成': 'Completed',
  '进行中': 'Current',
  '待处理': 'Pending',
  '阻塞': 'Blocked',
  
  // 负责人
  '负责人': 'Responsible',
  '申请人': 'Applicant',
  '部门经理': 'Department Manager',
  '财务经理': 'Finance Manager',
  '采购员': 'Purchaser',
  '供应商': 'Supplier',
  '质检员': 'Quality Inspector',
  
  // 时间信息
  '处理时间': 'Processing Time',
  '待处理': 'Pending',
  '件': 'items',
  
  // 阻塞流程表格
  '阻塞流程': 'Blocked Processes',
  '流程编号': 'Process Number',
  '当前步骤': 'Current Step',
  '负责人': 'Responsible',
  '部门': 'Department',
  '阻塞时长': 'Blocked Duration',
  '优先级': 'Priority',
  '阻塞原因': 'Block Reason',
  '操作': 'Actions',
  
  // 操作按钮
  '联系': 'Contact',
  '刷新': 'Refresh',
  '导出': 'Export',
  
  // 优先级
  '低': 'Low',
  '中': 'Medium',
  '高': 'High',
  '紧急': 'Urgent',
  
  // 阻塞原因
  '等待部门经理审批': 'Waiting for department manager approval',
  '预算超限，需要额外审批': 'Budget exceeded, additional approval required',
  '供应商未响应': 'Supplier not responding',
  '检验设备故障': 'Inspection equipment failure',
  
  // 消息提示
  '查询失败': 'Query failed',
  '操作成功': 'Operation successful',
  '操作失败': 'Operation failed',
  
  // 分页
  '条/页': 'items/page',
  '跳至': 'Go to',
  '页': 'page'
};

// 执行翻译
Object.keys(translations).forEach(chinese => {
  const english = translations[chinese];
  const regex = new RegExp(chinese.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
  content = content.replace(regex, english);
});

// 写回文件
fs.writeFileSync('src/views/ProcessMonitor.vue', content, 'utf8');

console.log('Process Monitor page translation completed!');
