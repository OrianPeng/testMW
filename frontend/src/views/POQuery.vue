<template>
  <div class="po-query-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-title">
        <ShoppingOutlined />
        <span>PO Query</span>
      </div>
      <div class="page-actions">
        <a-button
          type="primary"
          @click="handleSendToSupplier"
        >
          <SendOutlined />
          Send to Supplier
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
        <a-form-item label="PO Number">
          <a-input
            v-model:value="queryForm.po_number"
            placeholder="Please enterPO Number"
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
            <a-select-option value="draft">
              Draft
            </a-select-option>
            <a-select-option value="sent">
              Sent
            </a-select-option>
            <a-select-option value="confirmed">
              Confirmed
            </a-select-option>
            <a-select-option value="delivered">
              Delivered
            </a-select-option>
            <a-select-option value="cancelled">
              Cancelled
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="Supplier">
          <a-input
            v-model:value="queryForm.supplier"
            placeholder="Please enter supplier name"
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
    <a-row
      :gutter="16"
      class="stats-row"
    >
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Pending Send"
            :value="stats.draft"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <EditOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Sent"
            :value="stats.sent"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <SendOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic
            title="Confirmed"
            :value="stats.confirmed"
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
            title="Delivered"
            :value="stats.delivered"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CarOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Data Table -->
    <a-card class="table-card">
      <template #title>
        <span>PO List</span>
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
        row-key="po_number"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
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
                v-if="canSend(record)" 
                type="link" 
                size="small" 
                @click="handleSend(record)"
              >
                Send
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
                v-if="canCancel(record)" 
                type="link" 
                size="small" 
                danger
                @click="handleCancel(record)"
              >
                Cancel
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Detail Drawer -->
    <a-drawer
      v-model:open="detailVisible"
      title="PO Details"
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
          <a-descriptions-item label="PO Number">
            {{ currentRecord.po_number }}
          </a-descriptions-item>
          <a-descriptions-item label="Status">
            <a-tag :color="getStatusColor(currentRecord.status)">
              {{ getStatusText(currentRecord.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="Supplier">
            {{ currentRecord.supplier_name }}
          </a-descriptions-item>
          <a-descriptions-item label="Supplier Code">
            {{ currentRecord.supplier_code }}
          </a-descriptions-item>
          <a-descriptions-item label="Total Amount">
            <span class="amount-text">
              {{ formatCurrency(currentRecord.total_amount, currentRecord.currency) }}
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="Currency">
            {{ currentRecord.currency }}
          </a-descriptions-item>
          <a-descriptions-item label="Created By">
            {{ currentRecord.created_by }}
          </a-descriptions-item>
          <a-descriptions-item label="Created Time">
            {{ formatDate(currentRecord.created_at) }}
          </a-descriptions-item>
          <a-descriptions-item label="Sent Time">
            {{ currentRecord.sent_at ? formatDate(currentRecord.sent_at) : 'Not Sent' }}
          </a-descriptions-item>
          <a-descriptions-item label="Confirmed Time">
            {{ currentRecord.confirmed_at ? formatDate(currentRecord.confirmed_at) : 'Not Confirmed' }}
          </a-descriptions-item>
          <a-descriptions-item
            label="Notes"
            :span="2"
          >
            {{ currentRecord.notes || 'None' }}
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import {
  ShoppingOutlined,
  SendOutlined,
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  EditOutlined,
  CheckCircleOutlined,
  CarOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { purchaseApi } from '@/api/purchase'

const userStore = useUserStore()

// Reactive data
const loading = ref(false)
const detailVisible = ref(false)
const currentRecord = ref<any>(null)
const dataSource = ref<any[]>([])
const statistics = ref<any>({})

// Search Form
const queryForm = reactive({
  po_number: '',
  status: undefined,
  supplier: '',
  dateRange: undefined
})

// Load data from API
const loadData = async () => {
  loading.value = true
  try {
    // Load purchase orders
    const poResponse = await purchaseApi.getPurchaseOrders({
      page: pagination.current,
      limit: pagination.pageSize,
      status: queryForm.status,
      supplier_name: queryForm.supplier,
      po_number: queryForm.po_number
    })
    
    dataSource.value = Array.isArray(poResponse.data?.items) ? poResponse.data.items : []
    pagination.total = poResponse.data?.total || poResponse.total || 0
    
    // Load statistics
    const statsResponse = await purchaseApi.getPOStatistics()
    statistics.value = statsResponse.data || {}
    
  } catch (error) {
    console.error('Failed to load data:', error)
    message.error('Failed to load data')
  } finally {
    loading.value = false
  }
}

// Table column definitions
const columns = [
  {
    title: 'PO Number',
    dataIndex: 'po_number',
    key: 'po_number',
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
    title: 'Supplier',
    dataIndex: 'supplier_name',
    key: 'supplier_name',
    width: 150
  },
  {
    title: 'Total Amount',
    dataIndex: 'total_amount',
    key: 'total_amount',
    width: 120
  },
  {
    title: 'Created By',
    dataIndex: 'created_by',
    key: 'created_by',
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

// Statistics data
const stats = computed(() => ({
  draft: statistics.value.draft || 0,
  sent: statistics.value.sent || 0,
  confirmed: statistics.value.confirmed || 0,
  delivered: statistics.value.delivered || 0
}))

// Methods
const handleSearch = async () => {
  await loadData()
}

const handleReset = () => {
  Object.assign(queryForm, {
    po_number: '',
    status: undefined,
    supplier: '',
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

const handleView = async (record: any) => {
  try {
    const response = await purchaseApi.getPurchaseOrderById(record.po_number)
    currentRecord.value = response.data
    detailVisible.value = true
  } catch (error) {
    console.error('Failed to load PO details:', error)
    message.error('Failed to load PO details')
  }
}

const handleSend = (record: any) => {
  message.info('Send feature in development...')
}

const handleEdit = (record: any) => {
  message.info('Edit feature in development...')
}

const handleCancel = (record: any) => {
  message.info('Cancel feature in development...')
}

const handleSendToSupplier = () => {
  message.info('批量Send feature in development...')
}

// Helper methods
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    draft: 'orange',
    sent: 'blue',
    confirmed: 'green',
    delivered: 'green',
    cancelled: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    draft: 'Draft',
    sent: 'Sent',
    confirmed: 'Confirmed',
    delivered: 'Delivered',
    cancelled: 'Cancelled'
  }
  return texts[status] || status
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

const canSend = (record: any) => {
  return userStore.hasPermission('po:create') && record.status === 'draft'
}

const canEdit = (record: any) => {
  return userStore.hasPermission('po:create') && record.status === 'draft'
}

const canCancel = (record: any) => {
  return userStore.hasPermission('po:create') && 
         (record.status === 'draft' || record.status === 'sent')
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.po-query-page {
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

.amount-text {
  font-weight: 500;
  color: #262626;
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
