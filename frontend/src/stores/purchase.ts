import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { purchaseApi } from '@/api/purchase'
import type { PurchaseRequest, PurchaseRequestStatus, CreatePurchaseRequestRequest } from '@/types/purchase'

export const usePurchaseStore = defineStore('purchase', () => {
  const purchaseRequests = ref<PurchaseRequest[]>([])
  const currentRequest = ref<PurchaseRequest | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const pendingQuantity = computed(() => 
    purchaseRequests.value.filter(req => req.status === 'pending').length
  )

  const processingQuantity = computed(() => 
    purchaseRequests.value.filter(req => req.status === 'processing').length
  )

  const completedQuantity = computed(() => 
    purchaseRequests.value.filter(req => req.status === 'completed').length
  )

  const failedQuantity = computed(() => 
    purchaseRequests.value.filter(req => req.status === 'failed').length
  )

  const fetchPurchaseRequests = async (params?: {
    status?: PurchaseRequestStatus
    limit?: number
    offset?: number
  }) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await purchaseApi.getPurchaseRequests(params)
      purchaseRequests.value = response.requests
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取采购请求失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchPurchaseRequestById = async (requestId: string) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await purchaseApi.getPurchaseRequestById(requestId)
      currentRequest.value = response
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取采购请求Details失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createPurchaseRequest = async (data: CreatePurchaseRequestRequest) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await purchaseApi.createPurchaseRequest(data)
      // Refresh列表
      await fetchPurchaseRequests()
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'create采购请求失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updatePurchaseRequest = async (requestId: string, data: Partial<CreatePurchaseRequestRequest>) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await purchaseApi.updatePurchaseRequest(requestId, data)
      // 更新本地数据
      const index = purchaseRequests.value.findIndex(req => req.request_id === requestId)
      if (index !== -1) {
        purchaseRequests.value[index] = { ...purchaseRequests.value[index], ...data }
      }
      if (currentRequest.value?.request_id === requestId) {
        currentRequest.value = { ...currentRequest.value, ...data }
      }
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新采购请求失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deletePurchaseRequest = async (requestId: string) => {
    loading.value = true
    error.value = null
    
    try {
      await purchaseApi.deletePurchaseRequest(requestId)
      // 从本地列表Medium移除
      purchaseRequests.value = purchaseRequests.value.filter(req => req.request_id !== requestId)
      if (currentRequest.value?.request_id === requestId) {
        currentRequest.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Delete采购请求失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const getQueueStats = async () => {
    try {
      return await purchaseApi.getQueueStats()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取队列Status失败'
      throw err
    }
  }

  const clearQueue = async () => {
    loading.value = true
    error.value = null
    
    try {
      await purchaseApi.clearQueue()
      // Refresh列表
      await fetchPurchaseRequests()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '清空队列失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    purchaseRequests,
    currentRequest,
    loading,
    error,
    pendingQuantity,
    processingQuantity,
    completedQuantity,
    failedQuantity,
    fetchPurchaseRequests,
    fetchPurchaseRequestById,
    createPurchaseRequest,
    updatePurchaseRequest,
    deletePurchaseRequest,
    getQueueStats,
    clearQueue
  }
})
