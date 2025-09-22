import fs from 'fs';

// 读取PR查询页面文件
let content = fs.readFileSync('src/views/PRQuery.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 表格标题
  'PR列表': 'PR List',
  '共 {{ pagination.total }} 条记录': 'Total {{ pagination.total }} records',
  
  // 表格列标题
  'PR编号': 'PR Number',
  '状态': 'Status',
  '申请人': 'Requester',
  '供应商': 'Supplier',
  '物料': 'Material',
  '数量': 'Quantity',
  '单价': 'Unit Price',
  '总金额': 'Total Amount',
  '优先级': 'Priority',
  '紧急程度': 'Urgency',
  '创建时间': 'Created Time',
  '操作': 'Actions',
  
  // 状态标签
  '待审核': 'Pending',
  '已批准': 'Approved',
  '处理中': 'Processing',
  '已完成': 'Completed',
  '已拒绝': 'Rejected',
  '失败': 'Failed',
  
  // 优先级
  '低': 'Low',
  '中': 'Medium',
  '高': 'High',
  '紧急': 'Urgent',
  
  // 紧急程度
  '正常': 'Normal',
  '紧急': 'Urgent',
  '非常紧急': 'Critical',
  
  // 操作按钮
  '查看': 'View',
  '编辑': 'Edit',
  '删除': 'Delete',
  '刷新': 'Refresh',
  '导出': 'Export',
  
  // 详情抽屉
  'PR详情': 'PR Details',
  '基本信息': 'Basic Information',
  '物料信息': 'Material Information',
  '供应商信息': 'Supplier Information',
  '审批信息': 'Approval Information',
  '关闭': 'Close',
  
  // 创建/编辑表单
  '创建PR': 'Create PR',
  '编辑PR': 'Edit PR',
  '保存': 'Save',
  '取消': 'Cancel',
  '请选择': 'Please select',
  '请输入': 'Please enter',
  
  // 表单字段
  '文档类型': 'Document Type',
  '工厂': 'Plant',
  '物料组': 'Material Group',
  '单位类型': 'Unit Type',
  '采购组织': 'Purchase Organization',
  '货币': 'Currency',
  '交货日期': 'Delivery Date',
  '备注': 'Comments',
  '重试次数': 'Retry Count',
  '更新时间': 'Updated Time',
  
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
fs.writeFileSync('src/views/PRQuery.vue', content, 'utf8');

console.log('PR Query page translation completed!');
