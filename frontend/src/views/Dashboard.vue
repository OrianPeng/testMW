<template>
  <div class="dashboard-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-title">
        <DashboardOutlined />
        <span>Dashboard</span>
      </div>
      <div class="page-actions">
        <a-range-picker
          v-model:value="dateRange"
          format="YYYY-MM-DD"
          @change="handleDateChange"
        />
        <a-button
          :loading="loading"
          @click="handleRefresh"
        >
          <ReloadOutlined />
          Refresh
        </a-button>
      </div>
    </div>

    <!-- Key Metrics Cards -->
    <a-row
      :gutter="16"
      class="metrics-row"
    >
      <a-col :span="6">
        <a-card class="metric-card">
          <a-statistic
            title="Today PR Created"
            :value="metrics.todayPR"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <FileTextOutlined />
            </template>
            <template #suffix>
              <span class="metric-trend up">+12%</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="metric-card">
          <a-statistic
            title="Today PO Sent"
            :value="metrics.todayPO"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <SendOutlined />
            </template>
            <template #suffix>
              <span class="metric-trend up">+8%</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="metric-card">
          <a-statistic
            title="Pending ApprovalQuantity"
            :value="metrics.pendingApproval"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
            <template #suffix>
              <span class="metric-trend down">-5%</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="metric-card">
          <a-statistic
            title="Average Approval Time"
            :value="metrics.avgApprovalTime"
            suffix="hours"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
            <template #suffix>
              <span class="metric-trend down">-2hours</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Charts Area -->
    <a-row
      :gutter="16"
      class="charts-row"
    >
      <a-col :span="12">
        <a-card title="PR Creation Trend">
          <div
            ref="prTrendChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="PO Status Distribution">
          <div
            ref="poStatusChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-row
      :gutter="16"
      class="charts-row"
    >
      <a-col :span="12">
        <a-card title="Supplier Procurement Amount Ranking">
          <div
            ref="supplierChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Department Procurement Statistics">
          <div
            ref="departmentChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- Data Table Area -->
    <a-row
      :gutter="16"
      class="tables-row"
    >
      <a-col :span="12">
        <a-card
          title="Recent PR"
          size="small"
        >
          <a-table
            :columns="recentPRColumns"
            :data-source="recentPRData"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusText(record.status) }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                {{ formatCurrency(record.amount) }}
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card
          title="Recent PO"
          size="small"
        >
          <a-table
            :columns="recentPOColumns"
            :data-source="recentPOData"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="getPOStatusColor(record.status)">
                  {{ getPOStatusText(record.status) }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                {{ formatCurrency(record.amount) }}
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { purchaseApi } from '@/api/purchase'
import {
  DashboardOutlined,
  ReloadOutlined,
  FileTextOutlined,
  SendOutlined,
  ClockCircleOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

// Reactive data
const loading = ref(false)
const dateRange = ref([dayjs().subtract(7, 'day'), dayjs()])

// Chart references
const prTrendChartRef = ref<HTMLElement>()
const poStatusChartRef = ref<HTMLElement>()
const supplierChartRef = ref<HTMLElement>()
const departmentChartRef = ref<HTMLElement>()

// Key metrics
const metrics = reactive({
  todayPR: 0,
  todayPO: 0,
  pendingApproval: 0,
  avgApprovalTime: 0
})

// Recent PR data
const recentPRData = ref([
  {
    id: '1',
    request_id: 'PR-20231201-001',
    requester: 'John Smith',
    amount: 50000,
    status: 'pending',
    created_at: '2023-12-01 10:30'
  },
  {
    id: '2',
    request_id: 'PR-20231201-002',
    requester: 'Jane Doe',
    amount: 30000,
    status: 'approved',
    created_at: '2023-12-01 14:20'
  },
  {
    id: '3',
    request_id: 'PR-20231201-003',
    requester: 'Mike Johnson',
    amount: 80000,
    status: 'processing',
    created_at: '2023-12-01 16:45'
  },
  {
    id: '4',
    request_id: 'PR-20231201-004',
    requester: 'Sarah Wilson',
    amount: 15000,
    status: 'rejected',
    created_at: '2023-12-01 09:15'
  },
  {
    id: '5',
    request_id: 'PR-20231201-005',
    requester: 'David Brown',
    amount: 25000,
    status: 'completed',
    created_at: '2023-11-30 15:20'
  }
])

// Recent PO data
const recentPOData = ref([
  {
    id: '1',
    po_number: 'PO-20231201-001',
    supplier: 'ABC Company',
    amount: 50000,
    status: 'sent',
    created_at: '2023-12-01 11:00'
  },
  {
    id: '2',
    po_number: 'PO-20231201-002',
    supplier: 'XYZ Company',
    amount: 30000,
    status: 'confirmed',
    created_at: '2023-12-01 15:30'
  },
  {
    id: '3',
    po_number: 'PO-20231201-003',
    supplier: 'DEF Company',
    amount: 80000,
    status: 'delivered',
    created_at: '2023-12-01 17:20'
  },
  {
    id: '4',
    po_number: 'PO-20231201-004',
    supplier: 'GHI Company',
    amount: 15000,
    status: 'draft',
    created_at: '2023-12-02 09:15'
  },
  {
    id: '5',
    po_number: 'PO-20231201-005',
    supplier: 'JKL Company',
    amount: 25000,
    status: 'cancelled',
    created_at: '2023-11-29 14:30'
  }
])

// Recent PR table columns
const recentPRColumns = [
  {
    title: 'PR Number',
    dataIndex: 'request_id',
    key: 'request_id',
    width: 120
  },
  {
    title: 'Created By',
    dataIndex: 'requester',
    key: 'requester',
    width: 80
  },
  {
    title: 'Amount',
    dataIndex: 'amount',
    key: 'amount',
    width: 100
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 80
  }
]

// Recent PO table columns
const recentPOColumns = [
  {
    title: 'PO Number',
    dataIndex: 'po_number',
    key: 'po_number',
    width: 120
  },
  {
    title: 'Supplier',
    dataIndex: 'supplier',
    key: 'supplier',
    width: 100
  },
  {
    title: 'Amount',
    dataIndex: 'amount',
    key: 'amount',
    width: 100
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: 80
  }
]

// Load data from API
const loadData = async () => {
  try {
    // Load dashboard metrics
    const metricsResponse = await purchaseApi.getDashboardMetrics()
    Object.assign(metrics, metricsResponse.data || {})
    
    // Load recent PR data
    const prResponse = await purchaseApi.getPurchaseRequests({ limit: 5 })
    recentPRData.value = Array.isArray(prResponse.data?.items) ? prResponse.data.items : []
    
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
    message.error('Failed to load dashboard data')
  }
}

// Methods
const handleDateChange = () => {
  handleRefresh()
}

const handleRefresh = async () => {
  loading.value = true
  try {
    await loadData()
    initCharts()
    message.success('Data refreshed')
  } catch (error) {
    console.error('Refresh failed:', error)
    message.error('Refresh failed')
  } finally {
    loading.value = false
  }
}

// Initialize charts
const initCharts = () => {
  nextTick(() => {
    // PR Creation Trend chart
    if (prTrendChartRef.value) {
      const prTrendChart = echarts.init(prTrendChartRef.value)
      const prTrendOption = {
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: ['12-01', '12-02', '12-03', '12-04', '12-05', '12-06', '12-07']
        },
        yAxis: {
          type: 'value',
          name: 'Quantity'
        },
        series: [
          {
            name: 'PR Creation Quantity',
            type: 'line',
            data: [20, 25, 18, 30, 22, 28, 25],
            smooth: true,
            itemStyle: {
              color: '#1890ff'
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
                  { offset: 1, color: 'rgba(24, 144, 255, 0.1)' }
                ]
              }
            }
          }
        ]
      }
      prTrendChart.setOption(prTrendOption)
    }

    // PO Status Distribution pie chart
    if (poStatusChartRef.value) {
      const poStatusChart = echarts.init(poStatusChartRef.value)
      const poStatusOption = {
        tooltip: {
          trigger: 'item'
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'POStatus',
            type: 'pie',
            radius: '50%',
            data: [
              { value: 35, name: 'Sent' },
              { value: 25, name: 'Confirmed' },
              { value: 20, name: 'Delivered' },
              { value: 15, name: 'Draft' },
              { value: 5, name: 'Cancelled' }
            ],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      }
      poStatusChart.setOption(poStatusOption)
    }

    // Supplier Procurement Amount Ranking
    if (supplierChartRef.value) {
      const supplierChart = echarts.init(supplierChartRef.value)
      const supplierOption = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        xAxis: {
          type: 'value',
          name: 'Amount (10K)'
        },
        yAxis: {
          type: 'category',
          data: ['ABC Company', 'XYZ Company', 'DEF Company', 'GHI Company', 'JKL Company']
        },
        series: [
          {
            name: 'Procurement Amount',
            type: 'bar',
            data: [120, 80, 60, 40, 30],
            itemStyle: {
              color: '#52c41a'
            }
          }
        ]
      }
      supplierChart.setOption(supplierOption)
    }

    // Department Procurement Statistics
    if (departmentChartRef.value) {
      const departmentChart = echarts.init(departmentChartRef.value)
      const departmentOption = {
        tooltip: {
          trigger: 'item'
        },
        series: [
          {
            name: 'Department Procurement',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '18',
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: [
              { value: 35, name: 'Procurement Department' },
              { value: 25, name: 'Production Department' },
              { value: 20, name: 'IT Department' },
              { value: 15, name: 'Administration Department' },
              { value: 5, name: 'Others' }
            ]
          }
        ]
      }
      departmentChart.setOption(departmentOption)
    }
  })
}

// Helper methods
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'orange',
    approved: 'green',
    processing: 'blue',
    rejected: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: 'Pending Approval',
    approved: 'Approved',
    processing: 'Processing',
    rejected: 'Rejected'
  }
  return texts[status] || status
}

const getPOStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    draft: 'orange',
    sent: 'blue',
    confirmed: 'green',
    delivered: 'green',
    cancelled: 'red'
  }
  return colors[status] || 'default'
}

const getPOStatusText = (status: string) => {
  const texts: Record<string, string> = {
    draft: 'Draft',
    sent: 'Sent',
    confirmed: 'Confirmed',
    delivered: 'Delivered',
    cancelled: 'Cancelled'
  }
  return texts[status] || status
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY'
  }).format(amount)
}

// Lifecycle
onMounted(async () => {
  await loadData()
  initCharts()
})
</script>

<style scoped>
.dashboard-page {
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

.page-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.metrics-row {
  margin-bottom: 16px;
}

.metric-card {
  text-align: center;
}

.metric-trend {
  font-size: 12px;
  margin-left: 8px;
}

.metric-trend.up {
  color: #52c41a;
}

.metric-trend.down {
  color: #ff4d4f;
}

.charts-row {
  margin-bottom: 16px;
}

.tables-row {
  margin-bottom: 0;
}

/* Responsive design */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .page-actions {
    width: 100%;
    justify-content: space-between;
  }
  
  .metrics-row .ant-col {
    margin-bottom: 16px;
  }
  
  .charts-row .ant-col {
    margin-bottom: 16px;
  }
  
  .tables-row .ant-col {
    margin-bottom: 16px;
  }
}
</style>
