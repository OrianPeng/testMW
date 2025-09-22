<template>
  <div class="supplier-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-title">
        <TeamOutlined />
        <span>Supplier Management</span>
      </div>
      <div class="page-actions">
        <a-button
          type="primary"
          @click="showAddModal"
        >
          <PlusOutlined />
          AddSupplier
        </a-button>
      </div>
    </div>

    <!-- Search Form -->
    <a-card
      class="query-card"
      title="Search Conditions"
    >
      <a-form
        :model="queryForm"
        layout="inline"
        class="query-form"
        @finish="handleSearch"
      >
        <a-form-item label="Supplier Name">
          <a-input
            v-model:value="queryForm.name"
            placeholder="Please enterSupplier Name"
            allow-clear
          />
        </a-form-item>
        
        <a-form-item label="Supplier Code">
          <a-input
            v-model:value="queryForm.code"
            placeholder="Please enterSupplier Code"
            allow-clear
          />
        </a-form-item>
        
        <a-form-item label="Status">
          <a-select
            v-model:value="queryForm.status"
            placeholder="Please selectStatus"
            allow-clear
            style="width: 120px"
          >
            <a-select-option value="active">
              Activated
            </a-select-option>
            <a-select-option value="inactive">
              Deactivated
            </a-select-option>
            <a-select-option value="pending">
              Pending
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item>
          <a-space>
            <a-button
              type="primary"
              html-type="submit"
              :loading="loading"
            >
              <SearchOutlined />
              Search
            </a-button>
            <a-button @click="handleReset">
              <ReloadOutlined />
              Reset
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- Statistics Cards -->
    <a-row
      :gutter="16"
      class="stats-row"
    >
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Total Suppliers"
            :value="stats.total"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <TeamOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="ActivatedSupplier"
            :value="stats.active"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Added This Month"
            :value="stats.newThisMonth"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <PlusOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Pending"
            :value="stats.pending"
            :value-style="{ color: '#ff4d4f' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Data Table -->
    <a-card class="table-card">
      <template #title>
        <span>Supplier List</span>
        <a-tag
          color="blue"
          style="margin-left: 8px"
        >
          Total {{ pagination.total }} records
        </a-tag>
      </template>
      
      <template #extra>
        <a-space>
          <a-button
            :loading="loading"
            @click="handleRefresh"
          >
            <ReloadOutlined />
            Refresh
          </a-button>
          <a-button @click="handleExport">
            <DownloadOutlined />
            Export
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1200 }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'contact_info'">
            <div class="contact-info">
              <div>{{ record.contact_person }}</div>
              <div class="contact-detail">
                {{ record.phone }}
              </div>
              <div class="contact-detail">
                {{ record.email }}
              </div>
            </div>
          </template>
          
          <template v-else-if="column.key === 'created_at'">
            {{ formatDate(record.created_at) }}
          </template>
          
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button
                type="link"
                size="small"
                @click="handleView(record)"
              >
                View
              </a-button>
              <a-button 
                v-if="canEdit(record)" 
                type="link" 
                size="small" 
                @click="handleEdit(record)"
              >
                Edit
              </a-button>
              <a-button 
                v-if="canToggleStatus(record)" 
                type="link" 
                size="small" 
                @click="handleToggleStatus(record)"
              >
                {{ record.status === 'active' ? 'Deactivated' : 'Activated' }}
              </a-button>
              <a-button 
                v-if="canDelete(record)" 
                type="link" 
                size="small" 
                danger
                @click="handleDelete(record)"
              >
                Delete
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Detail Drawer -->
    <a-drawer
      v-model:open="detailVisible"
      title="Supplier Details"
      placement="right"
      :width="600"
    >
      <div
        v-if="currentRecord"
        class="detail-content"
      >
        <a-descriptions
          :column="2"
          bordered
        >
          <a-descriptions-item label="Supplier Code">
            {{ currentRecord.code }}
          </a-descriptions-item>
          <a-descriptions-item label="Supplier Name">
            {{ currentRecord.name }}
          </a-descriptions-item>
          <a-descriptions-item label="Status">
            <a-tag :color="getStatusColor(currentRecord.status)">
              {{ getStatusText(currentRecord.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="SupplierType">
            {{ currentRecord.type }}
          </a-descriptions-item>
          <a-descriptions-item label="Contact Person">
            {{ currentRecord.contact_person }}
          </a-descriptions-item>
          <a-descriptions-item label="ContactPhone">
            {{ currentRecord.phone }}
          </a-descriptions-item>
          <a-descriptions-item label="Email">
            {{ currentRecord.email }}
          </a-descriptions-item>
          <a-descriptions-item label="Address">
            {{ currentRecord.address }}
          </a-descriptions-item>
          <a-descriptions-item label="Main Procurement Categories">
            {{ currentRecord.categories ? currentRecord.categories.join(', ') : '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Created Time">
            {{ formatDate(currentRecord.created_at) }}
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-drawer>

    <!-- Add/Edit Modal -->
    <a-modal
      v-model:open="addModalVisible"
      :title="isEdit ? 'EditSupplier' : 'AddSupplier'"
      :width="600"
      @ok="handleAddSubmit"
      @cancel="handleAddCancel"
    >
      <a-form
        ref="addFormRef"
        :model="addForm"
        :rules="addRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="Supplier Code"
              name="code"
            >
              <a-input
                v-model:value="addForm.code"
                placeholder="Please enterSupplier Code"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Supplier Name"
              name="name"
            >
              <a-input
                v-model:value="addForm.name"
                placeholder="Please enterSupplier Name"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="SupplierType"
              name="type"
            >
              <a-select
                v-model:value="addForm.type"
                placeholder="Please selectSupplierType"
              >
                <a-select-option value="manufacturer">
                  Manufacturer
                </a-select-option>
                <a-select-option value="distributor">
                  Distributor
                </a-select-option>
                <a-select-option value="service">
                  Service Provider
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Status"
              name="status"
            >
              <a-select
                v-model:value="addForm.status"
                placeholder="Please selectStatus"
              >
                <a-select-option value="active">
                  Activated
                </a-select-option>
                <a-select-option value="inactive">
                  Deactivated
                </a-select-option>
                <a-select-option value="pending">
                  Pending
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item
          label="Contact Person"
          name="contact_person"
        >
          <a-input
            v-model:value="addForm.contact_person"
            placeholder="Please enter contact person name"
          />
        </a-form-item>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="ContactPhone"
              name="phone"
            >
              <a-input
                v-model:value="addForm.phone"
                placeholder="Please enterContactPhone"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Email"
              name="email"
            >
              <a-input
                v-model:value="addForm.email"
                placeholder="Please enter email address"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item
          label="Address"
          name="address"
        >
          <a-textarea
            v-model:value="addForm.address"
            placeholder="Please enter detailed address"
          />
        </a-form-item>
        
        <a-form-item
          label="Main Procurement Categories"
          name="categories"
        >
          <a-select
            v-model:value="addForm.categories"
            mode="multiple"
            placeholder="Please selectMain Procurement Categories"
          >
            <a-select-option value="Raw Materials">
              Raw Materials
            </a-select-option>
            <a-select-option value="Equipment">
              Equipment
            </a-select-option>
            <a-select-option value="Services">
              Services
            </a-select-option>
            <a-select-option value="Office Supplies">
              Office Supplies
            </a-select-option>
            <a-select-option value="ITEquipment">
              ITEquipment
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item
          label="Notes"
          name="notes"
        >
          <a-textarea
            v-model:value="addForm.notes"
            placeholder="Please enterNotesInformation"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { message, Modal } from 'ant-design-vue'
import {
  TeamOutlined,
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { purchaseApi } from '@/api/purchase'

// 供应商类型定义
interface Supplier {
  id: string
  code: string
  name: string
  type: string
  status: string
  contact_person: string
  phone: string
  email: string
  address: string
  categories: string[]
  notes: string
  created_at: string
}

const userStore = useUserStore()

// Reactive data
const loading = ref(false)
const detailVisible = ref(false)
const addModalVisible = ref(false)
const isEdit = ref(false)
const currentRecord = ref<any>(null)
const addFormRef = ref()

// Search Form
const queryForm = reactive({
  name: '',
  code: '',
  status: undefined
})

// Add Form
const addForm = reactive({
  code: '',
  name: '',
  type: '',
  status: 'active',
  contact_person: '',
  phone: '',
  email: '',
  address: '',
  categories: [],
  notes: ''
})

// Form validation rules
const addRules = {
  code: [{ required: true, message: 'Please enterSupplier Code' }],
  name: [{ required: true, message: 'Please enterSupplier Name' }],
  type: [{ required: true, message: 'Please selectSupplierType' }],
  contact_person: [{ required: true, message: 'Please enterContact Person' }],
  phone: [{ required: true, message: 'Please enterContactPhone' }],
  email: [
    { required: true, message: 'Please enter email address' },
    { type: 'email', message: 'Please enter valid email address' }
  ]
}

// Data source
const dataSource = ref<Supplier[]>([])
const statistics = ref<any>({})

// Table column definitions
const columns = [
  {
    title: 'Supplier Code',
    dataIndex: 'code',
    key: 'code',
    width: 120,
    fixed: 'left'
  },
  {
    title: 'Supplier Name',
    dataIndex: 'name',
    key: 'name',
    width: 200
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: 'Contact PersonInformation',
    key: 'contact_info',
    width: 200
  },
  {
    title: 'Main Procurement Categories',
    dataIndex: 'categories',
    key: 'categories',
    width: 150,
    customRender: ({ record }: any) => record.categories ? record.categories.join(', ') : '-'
  },
  {
    title: 'Created Time',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 150
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 200,
    fixed: 'right'
  }
]

// Pagination Config
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `Total ${total} records`
})

// Statistics data
const stats = computed(() => ({
  total: statistics.value.total || 0,
  active: statistics.value.active || 0,
  newThisMonth: statistics.value.new_this_month || 0,
  pending: statistics.value.pending || 0
}))

// Load data from API
const loadData = async () => {
  loading.value = true
  try {
    // Load suppliers
    const response = await purchaseApi.getSuppliers({
      page: pagination.current,
      limit: pagination.pageSize,
      status: queryForm.status,
      type: queryForm.type
    })
    
    // 处理响应数据
    if (response && response.success) {
      const items = Array.isArray(response.data?.items) ? response.data.items : []
      // 确保每个供应商都有categories字段
      dataSource.value = items.map((item: any) => ({
        ...item,
        categories: item.categories || []
      }))
      pagination.total = response.data?.pagination?.total || 0
    } else {
      dataSource.value = []
      pagination.total = 0
    }
    
    // Load statistics
    const statsResponse = await purchaseApi.getSupplierStatistics()
    if (statsResponse && statsResponse.success) {
      statistics.value = statsResponse.data || {}
    } else {
      statistics.value = {}
    }
    
  } catch (error) {
    console.error('Failed to load data:', error)
    message.error('Failed to load data')
    dataSource.value = []
    pagination.total = 0
    statistics.value = {}
  } finally {
    loading.value = false
  }
}

// Methods
const handleSearch = async () => {
  await loadData()
}

const handleReset = () => {
  Object.assign(queryForm, {
    name: '',
    code: '',
    status: undefined
  })
  handleSearch()
}

const handleRefresh = () => {
  loadData()
}

const handleExport = () => {
  message.info('Export feature in development...')
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleView = (record: any) => {
  currentRecord.value = record
  detailVisible.value = true
}

const handleEdit = (record: any) => {
  isEdit.value = true
  Object.assign(addForm, record)
  addModalVisible.value = true
}

const handleToggleStatus = async (record: any) => {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    const response = await purchaseApi.updateSupplier(record.id, { ...record, status: newStatus })
    if (response && response.success) {
      message.success(`Supplier has been ${newStatus === 'active' ? 'activated' : 'deactivated'}`)
      loadData()
    } else {
      message.error('Failed to toggle status')
    }
  } catch (error) {
    console.error('Failed to toggle status:', error)
    message.error('Failed to toggle status')
  }
}

const handleDelete = async (record: any) => {
  Modal.confirm({
    title: 'Confirm Delete',
    content: `Are you sure you want to delete supplier "${record.name}"?`,
    onOk: async () => {
      try {
        const response = await purchaseApi.deleteSupplier(record.id)
        if (response && response.success) {
          message.success('Supplier deleted successfully')
          loadData()
        } else {
          message.error('Failed to delete supplier')
        }
      } catch (error) {
        console.error('Failed to delete supplier:', error)
        message.error('Failed to delete supplier')
      }
    }
  })
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(addForm, {
    code: '',
    name: '',
    type: '',
    status: 'active',
    contact_person: '',
    phone: '',
    email: '',
    address: '',
    categories: [],
    notes: ''
  })
  addModalVisible.value = true
}

const handleAddSubmit = async () => {
  try {
    await addFormRef.value.validate()
    
    let response
    if (isEdit.value) {
      response = await purchaseApi.updateSupplier(addForm.id, addForm)
    } else {
      response = await purchaseApi.createSupplier(addForm)
    }
    
    if (response && response.success) {
      message.success(isEdit.value ? 'Update successful' : 'Added successfully')
      addModalVisible.value = false
      loadData()
    } else {
      message.error(isEdit.value ? 'Update failed' : 'Add failed')
    }
  } catch (error) {
    console.error('Failed to save supplier:', error)
    message.error(isEdit.value ? 'Update failed' : 'Add failed')
  }
}

const handleAddCancel = () => {
  addModalVisible.value = false
}

// Helper methods
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'green',
    inactive: 'red',
    pending: 'orange'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    active: 'Activated',
    inactive: 'Deactivated',
    pending: 'Pending'
  }
  return texts[status] || status
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const canEdit = (record: any) => {
  return userStore.hasPermission('supplier:create')
}

const canToggleStatus = (record: any) => {
  return userStore.hasPermission('supplier:create')
}

const canDelete = (record: any) => {
  return userStore.hasPermission('supplier:create')
}

// Lifecycle
onMounted(() => {
  handleSearch()
})
</script>

<style scoped>
.supplier-page {
  padding: 0;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-title {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
}

.page-title .anticon {
  margin-right: 8px;
  color: #1890ff;
}

.query-card {
  margin-bottom: 16px;
}

.query-form {
  margin-bottom: 0;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.table-card {
  margin-bottom: 0;
}

.contact-info {
  font-size: 12px;
}

.contact-detail {
  color: #666;
  margin-top: 2px;
}

.detail-content {
  padding: 16px 0;
}

/* Responsive design */
@media (max-width: 768px) {
  .query-form .ant-form-item {
    margin-bottom: 16px;
  }
  
  .stats-row .ant-col {
    margin-bottom: 16px;
  }
}
</style>
