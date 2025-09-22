<template>
  <div class="pr-query-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <FileSearchOutlined class="title-icon" />
            PR Query
          </h1>
          <p class="page-subtitle">Purchase Request Query and Management</p>
        </div>
        <div class="header-actions">
          <a-button
            type="primary"
            @click="showCreateModal = true"
            class="create-btn"
            size="large"
          >
            <PlusOutlined />
            Create PR
          </a-button>
        </div>
      </div>
    </div>

    <!-- Search Form -->
    <a-card
      class="query-card"
      :bordered="false"
      :body-style="{ padding: '24px' }"
    >
      <template #title>
        <div class="card-title">
          <SearchOutlined class="card-icon" />
          <span>Search Conditions</span>
        </div>
      </template>
      <a-form
        :model="queryForm"
        layout="inline"
        class="query-form"
        @finish="handleSearch"
      >
        <a-form-item label="PR Number">
          <a-input
            v-model:value="queryForm.request_id"
            placeholder="Enter PR number"
            allow-clear
          />
        </a-form-item>
        
        <a-form-item label="Status">
          <a-select
            v-model:value="queryForm.status"
            placeholder="Select status"
            allow-clear
            style="width: 120px"
          >
            <a-select-option value="pending">
              Pending
            </a-select-option>
            <a-select-option value="approved">
              Approved
            </a-select-option>
            <a-select-option value="rejected">
              Rejected
            </a-select-option>
            <a-select-option value="completed">
              Completed
            </a-select-option>
            <a-select-option value="processing">
              Processing
            </a-select-option>
            <a-select-option value="failed">
              Failed
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="Requester">
          <a-input
            v-model:value="queryForm.requester"
            placeholder="Enter requester"
            allow-clear
          />
        </a-form-item>
        
        <a-form-item label="Supplier">
          <a-input
            v-model:value="queryForm.vendor_code"
            placeholder="Enter supplier code"
            allow-clear
          />
        </a-form-item>
        
        <a-form-item label="Created Time">
          <a-range-picker
            v-model:value="queryForm.dateRange"
            format="YYYY-MM-DD"
          />
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
    <div class="stats-cards">
      <a-card class="stat-card pending" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon">
            <ClockCircleOutlined />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending }}</div>
            <div class="stat-title">Pending</div>
          </div>
        </div>
      </a-card>
      
      <a-card class="stat-card processing" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon">
            <SyncOutlined spin />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.processing }}</div>
            <div class="stat-title">Processing</div>
          </div>
        </div>
      </a-card>
      
      <a-card class="stat-card completed" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon">
            <CheckCircleOutlined />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.completed }}</div>
            <div class="stat-title">Completed</div>
          </div>
        </div>
      </a-card>
      
      <a-card class="stat-card failed" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon">
            <CloseCircleOutlined />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.failed }}</div>
            <div class="stat-title">Failed</div>
          </div>
        </div>
      </a-card>
    </div>

    <!-- Data Table -->
    <a-card class="table-card" :bordered="false">
      <template #title>
        <div class="card-title">
          <TableOutlined class="card-icon" />
          <span>PR List</span>
          <a-tag
            color="blue"
            class="record-count"
          >
            Total {{ pagination.total }} records
          </a-tag>
        </div>
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
        row-key="request_id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'priority'">
            <a-tag :color="getPriorityColor(record.priority)">
              {{ getPriorityText(record.priority) }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'urgency'">
            <a-tag :color="getUrgencyColor(record.urgency)">
              {{ getUrgencyText(record.urgency) }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'total_amount'">
            <span class="amount-text">
              {{ formatCurrency(record.total_amount, record.currency) }}
            </span>
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
                v-if="canApprove(record)" 
                type="link" 
                size="small" 
                @click="handleApprove(record)"
              >
                Approve
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
      title="PR Details"
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
          <a-descriptions-item label="PR Number">
            {{ currentRecord.request_id }}
          </a-descriptions-item>
          <a-descriptions-item label="Status">
            <a-tag :color="getStatusColor(currentRecord.status)">
              {{ getStatusText(currentRecord.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="Document Type">
            {{ currentRecord.doc_type }}
          </a-descriptions-item>
          <a-descriptions-item label="Plant">
            {{ currentRecord.plant }}
          </a-descriptions-item>
          <a-descriptions-item label="Material">
            {{ currentRecord.material }}
          </a-descriptions-item>
          <a-descriptions-item label="Quantity">
            {{ currentRecord.quantity }} {{ currentRecord.unit_type }}
          </a-descriptions-item>
          <a-descriptions-item label="Unit Price">
            {{ formatCurrency(currentRecord.unit_price, currentRecord.currency) }}
          </a-descriptions-item>
          <a-descriptions-item label="Total Amount">
            <span class="amount-text">
              {{ formatCurrency(currentRecord.total_amount, currentRecord.currency) }}
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="Supplier">
            {{ currentRecord.vendor_code }}
          </a-descriptions-item>
          <a-descriptions-item label="Delivery Date">
            {{ formatDate(currentRecord.delivery_date) }}
          </a-descriptions-item>
          <a-descriptions-item label="Created By">
            {{ currentRecord.requester }}
          </a-descriptions-item>
          <a-descriptions-item label="Purchase Organization">
            {{ currentRecord.purchase_organization }}
          </a-descriptions-item>
          <a-descriptions-item label="Priority">
            <a-tag :color="getPriorityColor(currentRecord.priority)">
              {{ getPriorityText(currentRecord.priority) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="Urgency">
            <a-tag :color="getUrgencyColor(currentRecord.urgency)">
              {{ getUrgencyText(currentRecord.urgency) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item
            label="Comments"
            :span="2"
          >
            {{ currentRecord.comments || 'None' }}
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-drawer>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="createModalVisible"
      :title="isEdit ? 'EditPR' : 'Create PR'"
      :width="800"
      @ok="handleCreateSubmit"
      @cancel="handleCreateCancel"
    >
      <a-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="PR Number"
              name="request_id"
            >
              <a-input
                v-model:value="createForm.request_id"
                placeholder="Please enter PR number"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Document Type"
              name="doc_type"
            >
              <a-select
                v-model:value="createForm.doc_type"
                placeholder="Please select document type"
              >
                <a-select-option value="NB">
                  Standard Purchase Request
                </a-select-option>
                <a-select-option value="UB">
                  Urgent Purchase Request
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="Plant"
              name="plant"
            >
              <a-input
                v-model:value="createForm.plant"
                placeholder="Please enter plant code"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Material"
              name="material"
            >
              <a-input
                v-model:value="createForm.material"
                placeholder="Please enter material number"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item
              label="Quantity"
              name="quantity"
            >
              <a-input-number
                v-model:value="createForm.quantity"
                :min="1"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
              label="Unit Price"
              name="unit_price"
            >
              <a-input-number
                v-model:value="createForm.unit_price"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
              label="Unit"
              name="unit_type"
            >
              <a-input
                v-model:value="createForm.unit_type"
                placeholder="e.g.: EA, KG"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="Supplier"
              name="vendor_code"
            >
              <a-input
                v-model:value="createForm.vendor_code"
                placeholder="Please enter supplier code"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Delivery Date"
              name="delivery_date"
            >
              <a-date-picker
                v-model:value="createForm.delivery_date"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item
          label="Material Description"
          name="short_text"
        >
          <a-textarea
            v-model:value="createForm.short_text"
            placeholder="Please enterMaterial Description"
          />
        </a-form-item>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
              label="Material Group"
              name="material_group"
            >
              <a-input
                v-model:value="createForm.material_group"
                placeholder="Please enterMaterial Group"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
              label="Purchase Organization"
              name="purchase_organization"
            >
              <a-input
                v-model:value="createForm.purchase_organization"
                placeholder="Please enter purchase organization"
              />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item
              label="Priority"
              name="priority"
            >
              <a-select
                v-model:value="createForm.priority"
                placeholder="Please selectPriority"
              >
                <a-select-option :value="1">
                  Low
                </a-select-option>
                <a-select-option :value="2">
                  Medium
                </a-select-option>
                <a-select-option :value="3">
                  High
                </a-select-option>
                <a-select-option :value="4">
                  Urgent
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
              label="Urgency"
              name="urgency"
            >
              <a-select
                v-model:value="createForm.urgency"
                placeholder="Please selectUrgency"
              >
                <a-select-option value="normal">
                  Normal
                </a-select-option>
                <a-select-option value="urgent">
                  Urgent
                </a-select-option>
                <a-select-option value="critical">
                  Urgent
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
              label="Currency"
              name="currency"
            >
              <a-select
                v-model:value="createForm.currency"
                placeholder="Please selectCurrency"
              >
                <a-select-option value="CNY">
                  Chinese Yuan
                </a-select-option>
                <a-select-option value="USD">
                  US Dollar
                </a-select-option>
                <a-select-option value="EUR">
                  Euro
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item
          label="Comments"
          name="comments"
        >
          <a-textarea
            v-model:value="createForm.comments"
            placeholder="Please enter comments"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Create PR Modal -->
    <CreatePRModal
      v-model:open="createModalVisible"
      @success="handlePRCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { usePurchaseStore } from '@/stores/purchase'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import {
  FileSearchOutlined,
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  ClockCircleOutlined,
  SyncOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  TableOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { PurchaseRequest, CreatePurchaseRequestRequest } from '@/types/purchase'
import { purchaseApi } from '@/api/purchase'
import CreatePRModal from '@/components/CreatePRModal.vue'

const purchaseStore = usePurchaseStore()
const userStore = useUserStore()

// Reactive data
const loading = ref(false)
const detailVisible = ref(false)
const createModalVisible = ref(false)
const isEdit = ref(false)
const currentRecord = ref<PurchaseRequest | null>(null)
const createFormRef = ref()
const dataSource = ref<PurchaseRequest[]>([])
const statistics = ref<any>({})

// Search Form
const queryForm = reactive({
  request_id: '',
  status: undefined,
  requester: '',
  vendor_code: '',
  dateRange: undefined
})

// Create form
const createForm = reactive<CreatePurchaseRequestRequest>({
  request_id: '',
  doc_type: '',
  plant: '',
  quantity: 1,
  unit_price: 0,
  material: '',
  delivery_date: '',
  vendor_code: '',
  short_text: '',
  material_group: '',
  unit_type: '',
  requester: userStore.user?.name || '',
  purchase_organization: '',
  currency: 'CNY',
  priority: 2,
  urgency: 'normal',
  comments: ''
})

// Form validation rules
const createRules = {
  request_id: [{ required: true, message: 'Please enter PR number' }],
  doc_type: [{ required: true, message: 'Please select document type' }],
  plant: [{ required: true, message: 'Please enter plant code' }],
  quantity: [{ required: true, message: 'Please enter quantity' }],
  unit_price: [{ required: true, message: 'Please enter unit price' }],
  material: [{ required: true, message: 'Please enter material number' }],
  delivery_date: [{ required: true, message: 'Please select delivery date' }],
  vendor_code: [{ required: true, message: 'Please enter supplier code' }],
  short_text: [{ required: true, message: 'Please enterMaterial Description' }],
  material_group: [{ required: true, message: 'Please enterMaterial Group' }],
  unit_type: [{ required: true, message: 'Please enterUnit' }],
  requester: [{ required: true, message: 'Please enter requester' }],
  purchase_organization: [{ required: true, message: 'Please enter purchase organization' }]
}

// Table column definitions
const columns = [
  {
    title: 'PR Number',
    dataIndex: 'request_id',
    key: 'request_id',
    width: 150,
    fixed: 'left'
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: 'Material',
    dataIndex: 'material',
    key: 'material',
    width: 120
  },
  {
    title: 'Quantity',
    dataIndex: 'quantity',
    key: 'quantity',
    width: 80
  },
  {
    title: 'Unit Price',
    dataIndex: 'unit_price',
    key: 'unit_price',
    width: 100
  },
  {
    title: 'Total Amount',
    dataIndex: 'total_amount',
    key: 'total_amount',
    width: 120
  },
  {
    title: 'Supplier',
    dataIndex: 'vendor_code',
    key: 'vendor_code',
    width: 120
  },
  {
    title: 'Created By',
    dataIndex: 'requester',
    key: 'requester',
    width: 100
  },
  {
    title: 'Priority',
    dataIndex: 'priority',
    key: 'priority',
    width: 80
  },
  {
    title: 'Urgency',
    dataIndex: 'urgency',
    key: 'urgency',
    width: 100
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

// Load data from API
const loadData = async () => {
  loading.value = true
  try {
    // Load purchase requests
    const prResponse = await purchaseApi.getPurchaseRequests({
      page: pagination.current,
      limit: pagination.pageSize,
      status: queryForm.status,
      requester: queryForm.requester,
      vendor_code: queryForm.vendor_code
    })
    
    dataSource.value = Array.isArray(prResponse.data?.items) ? prResponse.data.items : []
    pagination.total = prResponse.data?.total || prResponse.total || 0
    
    // Load statistics
    const statsResponse = await purchaseApi.getPRStatistics()
    statistics.value = statsResponse.data || {}
    
  } catch (error) {
    console.error('Failed to load data:', error)
    message.error('Failed to load data')
  } finally {
    loading.value = false
  }
}

// Statistics data
const stats = computed(() => ({
  pending: statistics.value.pending || 0,
  approved: statistics.value.approved || 0,
  processing: statistics.value.processing || 0,
  completed: statistics.value.completed || 0,
  rejected: statistics.value.rejected || 0,
  total: statistics.value.total || 0
}))

// Methods
const handleSearch = async () => {
  await loadData()
}

const handleReset = () => {
  Object.assign(queryForm, {
    request_id: '',
    status: undefined,
    requester: '',
    vendor_code: '',
    dateRange: undefined
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

const handleView = async (record: PurchaseRequest) => {
  try {
    const response = await purchaseApi.getPurchaseRequestById(record.request_id)
    currentRecord.value = response.data
    detailVisible.value = true
  } catch (error) {
    console.error('Failed to load PR details:', error)
    message.error('Failed to load PR details')
  }
}

const handleEdit = (record: PurchaseRequest) => {
  isEdit.value = true
  Object.assign(createForm, {
    ...record,
    delivery_date: dayjs(record.delivery_date)
  })
  createModalVisible.value = true
}

const handleApprove = (record: PurchaseRequest) => {
  message.info('Approval feature in development...')
}

const handleDelete = async (record: PurchaseRequest) => {
  try {
    await purchaseApi.deletePurchaseRequest(record.request_id)
    message.success('PR deleted successfully')
    loadData()
  } catch (error) {
    console.error('Failed to delete PR:', error)
    message.error('Failed to delete PR')
  }
}

// Create PR related methods
const handlePRCreated = () => {
  message.success('Purchase Request created successfully!')
  loadData()
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(createForm, {
    request_id: '',
    doc_type: '',
    plant: '',
    quantity: 1,
    unit_price: 0,
    material: '',
    delivery_date: '',
    vendor_code: '',
    short_text: '',
    material_group: '',
    unit_type: '',
    requester: userStore.user?.name || '',
    purchase_organization: '',
    currency: 'CNY',
    priority: 2,
    urgency: 'normal',
    comments: ''
  })
  createModalVisible.value = true
}

const handleCreateSubmit = async () => {
  try {
    await createFormRef.value.validate()
    
    const formData = {
      ...createForm,
      delivery_date: dayjs(createForm.delivery_date).format('YYYY-MM-DD')
    }
    
    await purchaseStore.createPurchaseRequest(formData)
    message.success(isEdit.value ? 'Update successful' : 'Create successful')
    createModalVisible.value = false
    handleSearch()
  } catch (error) {
    message.error(isEdit.value ? 'Update failed' : 'Create failed')
  }
}

const handleCreateCancel = () => {
  createModalVisible.value = false
}

// Helper methods
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'orange',
    approved: 'green',
    rejected: 'red',
    completed: 'green',
    processing: 'blue',
    failed: 'red',
    enough: 'green',
    lack: 'orange'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: 'Pending',
    approved: 'Approved',
    rejected: 'Rejected',
    completed: 'Completed',
    processing: 'Processing',
    failed: 'Failed',
    enough: 'Information complete',
    lack: 'Information incomplete'
  }
  return texts[status] || status
}

const getPriorityColor = (priority: number) => {
  const colors: Record<number, string> = {
    1: 'green',
    2: 'blue',
    3: 'orange',
    4: 'red'
  }
  return colors[priority] || 'default'
}

const getPriorityText = (priority: number) => {
  const texts: Record<number, string> = {
    1: 'Low',
    2: 'Medium',
    3: 'High',
    4: 'Urgent'
  }
  return texts[priority] || priority.toString()
}

const getUrgencyColor = (urgency: string) => {
  const colors: Record<string, string> = {
    normal: 'green',
    urgent: 'orange',
    critical: 'red'
  }
  return colors[urgency] || 'default'
}

const getUrgencyText = (urgency: string) => {
  const texts: Record<string, string> = {
    normal: 'Normal',
    urgent: 'Urgent',
    critical: 'Urgent'
  }
  return texts[urgency] || urgency
}

const formatCurrency = (amount: number, currency: string) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: currency || 'CNY'
  }).format(amount)
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const canEdit = (record: PurchaseRequest) => {
  return userStore.hasPermission('pr:create') && 
         (record.status === 'pending' || record.status === 'lack')
}

const canApprove = (record: PurchaseRequest) => {
  return userStore.hasPermission('pr:approve') && record.status === 'pending'
}

const canDelete = (record: PurchaseRequest) => {
  return userStore.hasPermission('pr:create') && 
         (record.status === 'pending' || record.status === 'lack')
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.pr-query-page {
  padding: 0;
  background: transparent;
}

/* Page header styles */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="10" cy="10" r="1" fill="rgba(255,255,255,0.1)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  opacity: 0.3;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 1;
}

.title-section {
  color: white;
}

.page-title {
  display: flex;
  align-items: center;
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.title-icon {
  margin-right: 12px;
  font-size: 28px;
  color: rgba(255, 255, 255, 0.9);
}

.page-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.create-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  font-weight: 600;
  height: 48px;
  padding: 0 24px;
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.create-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

/* Search card styles */
.query-card {
  margin-bottom: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  background: white;
}

.card-title {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.card-icon {
  margin-right: 8px;
  color: #3b82f6;
  font-size: 16px;
}

.record-count {
  margin-left: 12px;
  font-size: 12px;
  border-radius: 8px;
}

.query-form {
  margin-bottom: 0;
}

.query-form .ant-form-item {
  margin-bottom: 16px;
}

/* Statistics cards styles */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  background: white;
  transition: all 0.3s ease;
  overflow: hidden;
  position: relative;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
}

.stat-card.pending::before {
  background: linear-gradient(90deg, #f59e0b, #f97316);
}

.stat-card.processing::before {
  background: linear-gradient(90deg, #3b82f6, #1d4ed8);
}

.stat-card.completed::before {
  background: linear-gradient(90deg, #10b981, #059669);
}

.stat-card.failed::before {
  background: linear-gradient(90deg, #ef4444, #dc2626);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 24px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
  color: white;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
}

.stat-card.pending .stat-icon {
  background: linear-gradient(135deg, #f59e0b, #f97316);
}

.stat-card.processing .stat-icon {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
}

.stat-card.completed .stat-icon {
  background: linear-gradient(135deg, #10b981, #059669);
}

.stat-card.failed .stat-icon {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-title {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

/* Table card styles */
.table-card {
  margin-bottom: 0;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  background: white;
}

.table-card :deep(.ant-table) {
  border-radius: 12px;
  overflow: hidden;
}

.table-card :deep(.ant-table-thead > tr > th) {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
  font-weight: 600;
  color: #374151;
  padding: 16px 12px;
}

.table-card :deep(.ant-table-tbody > tr > td) {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
}

.table-card :deep(.ant-table-tbody > tr:hover > td) {
  background: #f8fafc;
}

.amount-text {
  font-weight: 600;
  color: #1f2937;
}

.detail-content {
  padding: 16px 0;
}

/* Status tag styles */
.status-tag {
  font-weight: 600;
  border-radius: 8px;
  padding: 4px 12px;
  font-size: 12px;
}

.priority-tag {
  font-weight: 600;
  border-radius: 8px;
  padding: 4px 12px;
  font-size: 12px;
}

.urgency-tag {
  font-weight: 600;
  border-radius: 8px;
  padding: 4px 12px;
  font-size: 12px;
}

/* Action button styles */
.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-buttons .ant-btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-buttons .ant-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Responsive design */
@media (max-width: 768px) {
  .page-header {
    padding: 24px 20px;
    margin-bottom: 16px;
  }
  
  .header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .title-icon {
    font-size: 20px;
  }
  
  .query-form .ant-form-item {
    margin-bottom: 16px;
  }
  
  .stats-cards {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-content {
    padding: 20px;
  }
  
  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>
