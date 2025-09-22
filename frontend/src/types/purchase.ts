export type PurchaseRequestStatus = 
  | 'pending' 
  | 'approved' 
  | 'rejected' 
  | 'completed' 
  | 'cancelled' 
  | 'processing' 
  | 'failed' 
  | 'enough' 
  | 'lack'

export interface PurchaseRequest {
  id: number
  request_id: string
  doc_type: string
  plant: string
  quantity: number
  unit_price: number
  material: string
  delivery_date: string
  vendor_code: string
  short_text: string
  material_group: string
  unit_type: string
  requester: string
  purchase_organization: string
  currency: string
  total_amount: number
  priority: number
  status: PurchaseRequestStatus
  urgency: string
  approver_id?: string
  approver_name?: string
  approved_at?: string
  rejection_reason?: string
  comments: string
  retry_count: number
  error_msg?: string
  processed_at?: string
  created_at: string
  updated_at: string
}

export interface CreatePurchaseRequestRequest {
  request_id: string
  doc_type: string
  plant: string
  quantity: number
  unit_price: number
  material: string
  delivery_date: string
  vendor_code: string
  short_text: string
  material_group: string
  unit_type: string
  requester: string
  purchase_organization: string
  currency?: string
  priority: number
  urgency: string
  comments?: string
}

export interface UpdatePurchaseRequestRequest {
  doc_type?: string
  plant?: string
  quantity?: number
  unit_price?: number
  material?: string
  delivery_date?: string
  vendor_code?: string
  short_text?: string
  material_group?: string
  unit_type?: string
  requester?: string
  purchase_organization?: string
  currency?: string
  priority?: number
  urgency?: string
  comments?: string
  status?: PurchaseRequestStatus
  approver_id?: string
  approver_name?: string
  rejection_reason?: string
  retry_count?: number
  error_msg?: string
  processed_at?: string
}

export interface PurchaseRequestListResponse {
  requests: PurchaseRequest[]
  pending_count: number
  limit: number
  offset: number
}

export interface QueueStats {
  queue_stats: {
    pending_count: number
  }
  rpa_status: {
    available: boolean
  }
}

export interface PurchaseRequestQuery {
  status?: PurchaseRequestStatus
  requester?: string
  plant?: string
  vendor_code?: string
  material_group?: string
  priority?: number
  urgency?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface QueuedPurchaseRequest extends PurchaseRequest {
  queue_position: number
  enqueued_at: string
}

export interface RPARequest {
  id: string
  request_id: string
  data: Record<string, any>
  metadata?: Record<string, string>
}

export interface RPAResponse {
  id: string
  request_id: string
  success: boolean
  data?: Record<string, any>
  error?: string
}

export interface RPAStatus {
  is_available: boolean
  last_check: string
  current_task?: string
}

export interface PurchaseRequestCallbackResponse {
  request_id: string
  status: PurchaseRequestStatus
  message: string
  data?: any
}
