import fs from 'fs';

// 读取供应商管理页面文件
let content = fs.readFileSync('src/views/Supplier.vue', 'utf8');

// 定义翻译映射
const translations = {
  // 页面标题
  '供应商管理': 'Supplier Management',
  '供应商信息管理和维护': 'Supplier Information Management and Maintenance',
  
  // 查询表单
  '查询条件': 'Search Conditions',
  '供应商代码': 'Supplier Code',
  '请输入供应商代码': 'Enter supplier code',
  '供应商名称': 'Supplier Name',
  '请输入供应商名称': 'Enter supplier name',
  '类型': 'Type',
  '请选择类型': 'Select type',
  '状态': 'Status',
  '请选择状态': 'Select status',
  '查询': 'Search',
  '重置': 'Reset',
  
  // 类型选项
  '制造商': 'Manufacturer',
  '经销商': 'Distributor',
  '服务商': 'Service Provider',
  
  // 状态选项
  '启用': 'Active',
  '禁用': 'Inactive',
  '待审核': 'Pending',
  
  // 统计卡片
  '总数': 'Total',
  '启用': 'Active',
  '本月新增': 'Added This Month',
  '待审核': 'Pending',
  
  // 表格标题
  '供应商列表': 'Supplier List',
  '共 {{ pagination.total }} 条记录': 'Total {{ pagination.total }} records',
  
  // 表格列标题
  '供应商代码': 'Supplier Code',
  '供应商名称': 'Supplier Name',
  '类型': 'Type',
  '状态': 'Status',
  '联系人': 'Contact Person',
  '电话': 'Phone',
  '邮箱': 'Email',
  '地址': 'Address',
  '采购类别': 'Procurement Categories',
  '备注': 'Notes',
  '创建时间': 'Created Time',
  '操作': 'Actions',
  
  // 操作按钮
  '查看': 'View',
  '编辑': 'Edit',
  '删除': 'Delete',
  '添加': 'Add',
  '刷新': 'Refresh',
  '导出': 'Export',
  
  // 详情抽屉
  '供应商详情': 'Supplier Details',
  '基本信息': 'Basic Information',
  '联系信息': 'Contact Information',
  '采购信息': 'Procurement Information',
  '关闭': 'Close',
  
  // 创建/编辑表单
  '添加供应商': 'Add Supplier',
  '编辑供应商': 'Edit Supplier',
  '保存': 'Save',
  '取消': 'Cancel',
  '请选择': 'Please select',
  '请输入': 'Please enter',
  
  // 表单字段
  '供应商代码': 'Supplier Code',
  '供应商名称': 'Supplier Name',
  '类型': 'Type',
  '状态': 'Status',
  '联系人': 'Contact Person',
  '电话': 'Phone',
  '邮箱': 'Email',
  '地址': 'Address',
  '采购类别': 'Procurement Categories',
  '备注': 'Notes',
  
  // 表单验证
  '请输入供应商代码': 'Please enter supplier code',
  '请输入供应商名称': 'Please enter supplier name',
  '请选择类型': 'Please select type',
  '请选择状态': 'Please select status',
  '请输入联系人': 'Please enter contact person',
  '请输入电话': 'Please enter phone number',
  '请输入邮箱地址': 'Please enter email address',
  '请输入有效的邮箱地址': 'Please enter a valid email address',
  '请输入地址': 'Please enter address',
  '请选择采购类别': 'Please select procurement categories',
  
  // 消息提示
  '查询失败': 'Query failed',
  '操作成功': 'Operation successful',
  '操作失败': 'Operation failed',
  '确认删除': 'Confirm deletion',
  '确定要删除这条记录吗？': 'Are you sure you want to delete this record?',
  '添加成功': 'Added successfully',
  '编辑成功': 'Edited successfully',
  '删除成功': 'Deleted successfully',
  
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
fs.writeFileSync('src/views/Supplier.vue', content, 'utf8');

console.log('Supplier Management page translation completed!');
