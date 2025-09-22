import axios from 'axios'
import type { 
  CreatePurchaseRequestRequest, 
  UpdatePurchaseRequestRequest
} from '@/types/purchase'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      // 处理未授权
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const purchaseApi = {
  // 获取采购请求列表
  getPurchaseRequests: async (params?: {
    status?: string
    requester?: string
    plant?: string
    vendor_code?: string
    material_group?: string
    priority?: number
    urgency?: string
    start_date?: string
    end_date?: string
    page?: number
    limit?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  }): Promise<any> => {
    return api.get('/purchase-requests', { params })
  },

  // 获取采购请求详情
  getPurchaseRequestById: async (requestId: string): Promise<any> => {
    return api.get(`/purchase-requests/${requestId}`)
  },

  // 创建采购请求
  createPurchaseRequest: async (data: CreatePurchaseRequestRequest): Promise<any> => {
    return api.post('/purchase-requests', data)
  },

  // 更新采购请求
  updatePurchaseRequest: async (requestId: string, data: UpdatePurchaseRequestRequest): Promise<any> => {
    return api.put(`/purchase-requests/${requestId}`, data)
  },

  // 删除采购请求
  deletePurchaseRequest: async (requestId: string): Promise<any> => {
    return api.delete(`/purchase-requests/${requestId}`)
  },

  // 获取采购请求统计
  getPRStatistics: async (): Promise<any> => {
    return api.get('/purchase-requests/statistics')
  },

  // 获取采购订单列表
  getPurchaseOrders: async (params?: {
    status?: string
    supplier_name?: string
    supplier_code?: string
    page?: number
    limit?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  }): Promise<any> => {
    return api.get('/purchase-orders', { params })
  },

  // 获取采购订单详情
  getPurchaseOrderById: async (poNumber: string): Promise<any> => {
    return api.get(`/purchase-orders/${poNumber}`)
  },

  // 创建采购订单
  createPurchaseOrder: async (data: any): Promise<any> => {
    return api.post('/purchase-orders', data)
  },

  // 更新采购订单
  updatePurchaseOrder: async (poNumber: string, data: any): Promise<any> => {
    return api.put(`/purchase-orders/${poNumber}`, data)
  },

  // 删除采购订单
  deletePurchaseOrder: async (poNumber: string): Promise<any> => {
    return api.delete(`/purchase-orders/${poNumber}`)
  },

  // 获取采购订单统计
  getPOStatistics: async (): Promise<any> => {
    return api.get('/purchase-orders/statistics')
  },

  // 获取供应商列表
  getSuppliers: async (params?: {
    status?: string
    type?: string
    page?: number
    limit?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  }): Promise<any> => {
    return api.get('/suppliers', { params })
  },

  // 获取供应商详情
  getSupplierById: async (id: string): Promise<any> => {
    return api.get(`/suppliers/${id}`)
  },

  // 创建供应商
  createSupplier: async (data: any): Promise<any> => {
    return api.post('/suppliers', data)
  },

  // 更新供应商
  updateSupplier: async (id: string, data: any): Promise<any> => {
    return api.put(`/suppliers/${id}`, data)
  },

  // 删除供应商
  deleteSupplier: async (id: string): Promise<any> => {
    return api.delete(`/suppliers/${id}`)
  },

  // 获取供应商统计
  getSupplierStatistics: async (): Promise<any> => {
    return api.get('/suppliers/statistics')
  },

  // 获取流程监控概览
  getProcessOverview: async (): Promise<any> => {
    return api.get('/process-monitor/overview')
  },

  // 获取被阻塞的流程
  getBlockedProcesses: async (): Promise<any> => {
    return api.get('/process-monitor/blocked')
  },

  // 获取流程步骤
  getProcessSteps: async (processType?: string): Promise<any> => {
    return api.get('/process-monitor/steps', { 
      params: processType ? { process_type: processType } : {} 
    })
  },

  // 获取仪表板指标
  getDashboardMetrics: async (): Promise<any> => {
    return api.get('/dashboard/metrics')
  },

  // 获取PR趋势图表
  getPRTrendChart: async (): Promise<any> => {
    return api.get('/dashboard/charts/pr-trend')
  },

  // 获取PO状态分布
  getPOStatusDistribution: async (): Promise<any> => {
    return api.get('/dashboard/charts/po-status')
  },

  // 获取供应商排名
  getSupplierRanking: async (): Promise<any> => {
    return api.get('/dashboard/charts/supplier-ranking')
  },

  // 获取部门统计
  getDepartmentStats: async (): Promise<any> => {
    return api.get('/dashboard/charts/department-stats')
  },

  // AI助手发送消息
  sendAIMessage: async (data: {
    message: string
    conversation_id?: string
  }): Promise<any> => {
    return api.post('/ai-assistant/message', data)
  },

  // 获取队列统计（保持兼容性）
  getQueueStats: async (): Promise<any> => {
    return api.get('/status')
  },

  // 清空队列（保持兼容性）
  clearQueue: async (): Promise<any> => {
    return api.post('/purchase-requests/queue/clear-all')
  },

  // 健康检查（保持兼容性）
  healthCheck: async (): Promise<any> => {
    return api.get('/health')
  }
}

export default api
