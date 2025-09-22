<template>
  <div class="process-monitor-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-title">
        <MonitorOutlined />
        <span>Process Monitor</span>
      </div>
      <div class="page-actions">
        <a-button
          :loading="loading"
          @click="handleRefresh"
        >
          <ReloadOutlined />
          Refresh
        </a-button>
      </div>
    </div>

    <!-- 统计概览 -->
    <a-row
      :gutter="16"
      class="overview-row"
    >
      <a-col :span="6">
        <a-card class="overview-card">
          <a-statistic
            title="Today Pending"
            :value="overview.todayPending"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="overview-card">
          <a-statistic
            title="Overdue Processes"
            :value="overview.overdue"
            :value-style="{ color: '#ff4d4f' }"
          >
            <template #prefix>
              <ExclamationCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="overview-card">
          <a-statistic
            title="Average Approval Time"
            :value="overview.avgApprovalTime"
            suffix="hours"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="overview-card">
          <a-statistic
            title="ProcessCompletion Rate"
            :value="overview.completionRate"
            suffix="%"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- BlockedProcess列表 -->
    <a-card
      class="blocked-processes-card"
      title="BlockedProcess"
    >
      <template #extra>
        <a-tag color="red">
          Urgent
        </a-tag>
      </template>
      
      <a-table
        :columns="blockedColumns"
        :data-source="blockedProcesses"
        :loading="loading"
        :pagination="false"
        size="small"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'blocked_duration'">
            <a-tag :color="getDurationColor(record.blocked_duration)">
              {{ record.blocked_duration }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'priority'">
            <a-tag :color="getPriorityColor(record.priority)">
              {{ getPriorityText(record.priority) }}
            </a-tag>
          </template>
          
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button
                type="link"
                size="small"
                @click="handleContact(record)"
              >
                <PhoneOutlined />
                ContactResponsible
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="handleEscalate(record)"
              >
                <ExclamationOutlined />
                升级处理
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Process统计图表 -->
    <a-row
      :gutter="16"
      class="charts-row"
    >
      <a-col :span="12">
        <a-card title="ProcessStatus分布">
          <div
            ref="statusChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Approve时长趋势">
          <div
            ref="timeChartRef"
            style="height: 300px;"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- Process StepsDetails -->
    <a-card title="Process StepsDetails">
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane
          key="pr-process"
          tab="PR Process"
        >
          <div class="process-steps">
            <div
              v-for="(step, index) in prProcessSteps"
              :key="index"
              class="process-step"
              :class="{ 
                active: step.status === 'current',
                completed: step.status === 'completed',
                blocked: step.status === 'blocked'
              }"
            >
              <div class="process-step-icon">
                <CheckOutlined v-if="step.status === 'completed'" />
                <ClockCircleOutlined v-else-if="step.status === 'current'" />
                <ExclamationCircleOutlined v-else-if="step.status === 'blocked'" />
                <MinusOutlined v-else />
              </div>
              <div class="process-step-content">
                <div class="process-step-title">
                  {{ step.title }}
                </div>
                <div class="process-step-responsible">
                  Responsible: {{ step.responsible }}
                </div>
                <div
                  v-if="step.time"
                  class="process-step-time"
                >
                  Processing Time: {{ step.time }}
                </div>
                <div
                  v-if="step.count"
                  class="process-step-count"
                >
                  Pending: {{ step.count }} items
                </div>
              </div>
            </div>
          </div>
        </a-tab-pane>
        
        <a-tab-pane
          key="po-process"
          tab="PO Process"
        >
          <div class="process-steps">
            <div
              v-for="(step, index) in poProcessSteps"
              :key="index"
              class="process-step"
              :class="{ 
                active: step.status === 'current',
                completed: step.status === 'completed',
                blocked: step.status === 'blocked'
              }"
            >
              <div class="process-step-icon">
                <CheckOutlined v-if="step.status === 'completed'" />
                <ClockCircleOutlined v-else-if="step.status === 'current'" />
                <ExclamationCircleOutlined v-else-if="step.status === 'blocked'" />
                <MinusOutlined v-else />
              </div>
              <div class="process-step-content">
                <div class="process-step-title">
                  {{ step.title }}
                </div>
                <div class="process-step-responsible">
                  Responsible: {{ step.responsible }}
                </div>
                <div
                  v-if="step.time"
                  class="process-step-time"
                >
                  Processing Time: {{ step.time }}
                </div>
                <div
                  v-if="step.count"
                  class="process-step-count"
                >
                  Pending: {{ step.count }} items
                </div>
              </div>
            </div>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { purchaseApi } from '@/api/purchase'
import {
  MonitorOutlined,
  ReloadOutlined,
  ClockCircleOutlined,
  ExclamationCircleOutlined,
  CheckCircleOutlined,
  PhoneOutlined,
  ExclamationOutlined,
  CheckOutlined,
  MinusOutlined
} from '@ant-design/icons-vue'
import * as echarts from 'echarts'

// Reactive data
const loading = ref(false)
const activeTab = ref('pr-process')
const statusChartRef = ref<HTMLElement>()
const timeChartRef = ref<HTMLElement>()

// 概览数据
const overview = reactive({
  todayPending: 15,
  overdue: 3,
  avgApprovalTime: 24,
  completionRate: 85
})

// BlockedProcess数据
const blockedProcesses = ref([
  {
    id: '1',
    request_id: 'PR-20231201-001',
    current_step: 'Department Approval',
    responsible: 'Mike Johnson',
    department: 'Procurement Department',
    blocked_duration: '2 day',
    priority: 4,
    reason: 'waiting Department ManagerApprove'
  },
  {
    id: '2',
    request_id: 'PR-20231201-002',
    current_step: 'Finance Approval',
    responsible: 'Sarah Wilson',
    department: 'Finance Department',
    blocked_duration: '1 day',
    priority: 3,
    reason: 'Budget exceeded, additional approval required'
  },
  {
    id: '3',
    request_id: 'PR-20231201-003',
    current_step: 'Supplier Confirmation',
    responsible: 'Jane Doe',
    department: 'Procurement Department',
    blocked_duration: '3 day',
    priority: 2,
    reason: 'Supplier'
  },
  {
    id: '4',
    request_id: 'PO-20231201-004',
    current_step: 'Quality Inspection',
    responsible: 'David Brown',
    department: 'Quality Department',
    blocked_duration: '1 day',
    priority: 3,
    reason: 'Inspection equipment failure'
  }
])

// BlockedProcess表格列
const blockedColumns = [
  {
    title: 'Process Number',
    dataIndex: 'request_id',
    key: 'request_id',
    width: 150
  },
  {
    title: 'Current Step',
    dataIndex: 'current_step',
    key: 'current_step',
    width: 120
  },
  {
    title: 'Responsible',
    dataIndex: 'responsible',
    key: 'responsible',
    width: 100
  },
  {
    title: 'Department',
    dataIndex: 'department',
    key: 'department',
    width: 100
  },
  {
    title: 'Blocked duration',
    dataIndex: 'blocked_duration',
    key: 'blocked_duration',
    width: 100
  },
  {
    title: 'Priority',
    dataIndex: 'priority',
    key: 'priority',
    width: 80
  },
  {
    title: 'reason',
    dataIndex: 'reason',
    key: 'reason',
    width: 200
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 150
  }
]

// PR Process步骤
const prProcessSteps = ref([
  {
    title: 'Submit Application',
    status: 'completed',
    responsible: 'Applicant',
    time: 'Average 0.5hours',
    count: 0
  },
  {
    title: 'Department Approval',
    status: 'current',
    responsible: 'Department Manager',
    time: 'Average 4hours',
    count: 8
  },
  {
    title: 'Finance Approval',
    status: 'pending',
    responsible: 'Finance Manager',
    time: 'Average 2hours',
    count: 5
  },
  {
    title: 'Procurement Execution',
    status: 'pending',
    responsible: 'Purchaser',
    time: 'Average 1hours',
    count: 2
  }
])

// PO Process步骤
const poProcessSteps = ref([
  {
    title: 'Generate PO',
    status: 'completed',
    responsible: 'System',
    time: 'Average 0.1hours',
    count: 0
  },
  {
    title: 'SendSupplier',
    status: 'current',
    responsible: 'Purchaser',
    time: 'Average 0.5hours',
    count: 3
  },
  {
    title: 'Supplier Confirmation',
    status: 'pending',
    responsible: 'Supplier',
    time: 'Average 24hours',
    count: 10
  },
  {
    title: 'Order Execution',
    status: 'pending',
    responsible: 'Supplier',
    time: 'Average 7 day',
    count: 15
  }
])

// Load data from API
const loadData = async () => {
  try {
    // Load process overview
    const overviewResponse = await purchaseApi.getProcessOverview()
    Object.assign(overview, overviewResponse.data || {})
    
    // Load blocked processes
    const blockedResponse = await purchaseApi.getBlockedProcesses()
    blockedProcesses.value = Array.isArray(blockedResponse.data?.items) ? blockedResponse.data.items : []
    
  } catch (error) {
    console.error('Failed to load process monitor data:', error)
    message.error('Failed to load process monitor data')
  }
}

// Methods
const handleRefresh = async () => {
  loading.value = true
  try {
    await loadData()
    message.success('Data refreshed')
  } catch (error) {
    console.error('Refresh failed:', error)
    message.error('Refresh failed')
  } finally {
    loading.value = false
  }
}

const handleContact = (record: any) => {
  message.info(` Contacting ${record.responsible}...`)
}

const handleEscalate = (record: any) => {
  message.info(`Escalating process ${record.request_id}...`)
}

const getDurationColor = (duration: string) => {
  if (duration.includes('day')) {
    const days = parseInt(duration)
    if (days >= 3) return 'red'
    if (days >= 1) return 'orange'
  }
  return 'green'
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

// Initialize charts
const initCharts = () => {
  nextTick(() => {
    // Status分布饼图
    if (statusChartRef.value) {
      const statusChart = echarts.init(statusChartRef.value)
      const statusOption = {
        tooltip: {
          trigger: 'item'
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'ProcessStatus',
            type: 'pie',
            radius: '50%',
            data: [
              { value: 35, name: 'Completed' },
              { value: 25, name: 'Current' },
              { value: 15, name: 'Approving' },
              { value: 10, name: 'Blocked' },
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
      statusChart.setOption(statusOption)
    }

    // Approve时长趋势图
    if (timeChartRef.value) {
      const timeChart = echarts.init(timeChartRef.value)
      const timeOption = {
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']
        },
        yAxis: {
          type: 'value',
          name: 'hours'
        },
        series: [
          {
            name: 'Average Approval Time',
            type: 'line',
            data: [20, 25, 18, 22, 28, 15, 12],
            smooth: true,
            itemStyle: {
              color: '#1890ff'
            }
          }
        ]
      }
      timeChart.setOption(timeOption)
    }
  })
}

// Lifecycle
onMounted(async () => {
  await loadData()
  initCharts()
})
</script>

<style scoped>
.process-monitor-page {
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

.overview-row {
  margin-bottom: 16px;
}

.overview-card {
  text-align: center;
}

.blocked-processes-card {
  margin-bottom: 16px;
}

.charts-row {
  margin-bottom: 16px;
}

.process-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.process-step {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  border-radius: 8px;
  background: #fafafa;
  border: 1px solid #e8e8e8;
  transition: all 0.3s;
}

.process-step.active {
  background: #e6f7ff;
  border-color: #1890ff;
}

.process-step.completed {
  background: #f6ffed;
  border-color: #52c41a;
}

.process-step.blocked {
  background: #fff2f0;
  border-color: #ff4d4f;
}

.process-step-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e8e8e8;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  flex-shrink: 0;
}

.process-step.active .process-step-icon {
  background: #1890ff;
  color: #fff;
}

.process-step.completed .process-step-icon {
  background: #52c41a;
  color: #fff;
}

.process-step.blocked .process-step-icon {
  background: #ff4d4f;
  color: #fff;
}

.process-step-content {
  flex: 1;
}

.process-step-title {
  font-size: 16px;
  font-weight: 500;
  color: #262626;
  margin-bottom: 4px;
}

.process-step-responsible {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.process-step-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.process-step-count {
  font-size: 12px;
  color: #1890ff;
  font-weight: 500;
}

/* Responsive design */
@media (max-width: 768px) {
  .overview-row .ant-col {
    margin-bottom: 16px;
  }
  
  .charts-row .ant-col {
    margin-bottom: 16px;
  }
  
  .process-step {
    flex-direction: column;
    text-align: center;
  }
  
  .process-step-icon {
    margin-right: 0;
    margin-bottom: 8px;
  }
}
</style>
