import fs from 'fs';

// 读取数据看板页面文件
let content = fs.readFileSync('src/views/Dashboard.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 页面标题
  '数据看板': 'Dashboard',
  '采购数据统计和分析': 'Procurement Data Statistics and Analysis',
  
  // 关键指标
  '今日PR创建': 'Today PR Created',
  '今日PO发送': 'Today PO Sent',
  '待审批': 'Pending Approval',
  '平均审批时长': 'Average Approval Time',
  '小时': 'hours',
  
  // 图表标题
  'PR创建趋势': 'PR Creation Trend',
  'PO状态分布': 'PO Status Distribution',
  '供应商排行': 'Supplier Ranking',
  '部门统计': 'Department Statistics',
  
  // 最近记录
  '最近PR': 'Recent PR',
  '最近PO': 'Recent PO',
  'PR编号': 'PR Number',
  '申请人': 'Requester',
  '金额': 'Amount',
  '状态': 'Status',
  '创建时间': 'Created Time',
  'PO编号': 'PO Number',
  '供应商': 'Supplier',
  '发送时间': 'Sent Time',
  
  // 状态标签
  '待审核': 'Pending',
  '已批准': 'Approved',
  '处理中': 'Processing',
  '已拒绝': 'Rejected',
  '已完成': 'Completed',
  '已发送': 'Sent',
  '已确认': 'Confirmed',
  '已交付': 'Delivered',
  '已取消': 'Cancelled',
  
  // 操作按钮
  '查看详情': 'View Details',
  '刷新': 'Refresh',
  '导出': 'Export',
  
  // 时间范围
  '最近7天': 'Last 7 Days',
  '最近30天': 'Last 30 Days',
  '最近90天': 'Last 90 Days',
  
  // 图表标签
  '日期': 'Date',
  '数量': 'Quantity',
  '状态': 'Status',
  '数量': 'Count',
  '供应商': 'Supplier',
  '订单数': 'Order Count',
  '部门': 'Department',
  'PR数量': 'PR Count',
  
  // 消息提示
  '查询失败': 'Query failed',
  '操作成功': 'Operation successful',
  '操作失败': 'Operation failed',
  
  // 空状态
  '暂无数据': 'No Data',
  '暂无PR记录': 'No PR records',
  '暂无PO记录': 'No PO records',
  
  // 加载状态
  '加载中...': 'Loading...',
  '数据加载中': 'Data loading',
  
  // 其他
  '总计': 'Total',
  '平均': 'Average',
  '最高': 'Highest',
  '最低': 'Lowest'
};

// 执行翻译
Object.keys(translations).forEach(chinese => {
  const english = translations[chinese];
  const regex = new RegExp(chinese.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
  content = content.replace(regex, english);
});

// 写回文件
fs.writeFileSync('src/views/Dashboard.vue', content, 'utf8');

console.log('Dashboard page translation completed!');
