import fs from 'fs';

// 读取PO查询页面文件
let content = fs.readFileSync('src/views/POQuery.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 页面标题
  'PO查询': 'PO Query',
  '采购订单查询和管理': 'Purchase Order Query and Management',
  
  // 查询表单
  '查询条件': 'Search Conditions',
  'PO编号': 'PO Number',
  '请输入PO编号': 'Enter PO number',
  '状态': 'Status',
  '请选择状态': 'Select status',
  '供应商': 'Supplier',
  '请输入供应商': 'Enter supplier',
  '日期范围': 'Date Range',
  '查询': 'Search',
  '重置': 'Reset',
  
  // 状态选项
  '草稿': 'Draft',
  '已发送': 'Sent',
  '已确认': 'Confirmed',
  '已交付': 'Delivered',
  '已取消': 'Cancelled',
  
  // 统计卡片
  '草稿': 'Draft',
  '已发送': 'Sent',
  '已确认': 'Confirmed',
  '已交付': 'Delivered',
  
  // 表格标题
  'PO列表': 'PO List',
  '共 {{ pagination.total }} 条记录': 'Total {{ pagination.total }} records',
  
  // 表格列标题
  'PO编号': 'PO Number',
  '状态': 'Status',
  '供应商名称': 'Supplier Name',
  '供应商代码': 'Supplier Code',
  '总金额': 'Total Amount',
  '货币': 'Currency',
  '创建人': 'Created By',
  '创建时间': 'Created Time',
  '发送时间': 'Sent Time',
  '确认时间': 'Confirmed Time',
  '备注': 'Notes',
  '操作': 'Actions',
  
  // 操作按钮
  '查看': 'View',
  '编辑': 'Edit',
  '删除': 'Delete',
  '发送': 'Send',
  '确认': 'Confirm',
  '刷新': 'Refresh',
  '导出': 'Export',
  
  // 详情抽屉
  'PO详情': 'PO Details',
  '基本信息': 'Basic Information',
  '供应商信息': 'Supplier Information',
  '订单信息': 'Order Information',
  '关闭': 'Close',
  
  // 创建/编辑表单
  '创建PO': 'Create PO',
  '编辑PO': 'Edit PO',
  '保存': 'Save',
  '取消': 'Cancel',
  '请选择': 'Please select',
  '请输入': 'Please enter',
  
  // 表单字段
  '供应商名称': 'Supplier Name',
  '供应商代码': 'Supplier Code',
  '总金额': 'Total Amount',
  '货币': 'Currency',
  '备注': 'Notes',
  
  // 消息提示
  '查询失败': 'Query failed',
  '操作成功': 'Operation successful',
  '操作失败': 'Operation failed',
  '确认删除': 'Confirm deletion',
  '确定要删除这条记录吗？': 'Are you sure you want to delete this record?',
  
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
fs.writeFileSync('src/views/POQuery.vue', content, 'utf8');

console.log('PO Query page translation completed!');
