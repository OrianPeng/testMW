<template>
  <a-modal
    v-model:open="visible"
    title="Create Purchase Request"
    :width="800"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Request ID" name="request_id">
            <a-input
              v-model:value="formData.request_id"
              placeholder="Auto-generated if empty"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Document Type" name="doc_type">
            <a-select
              v-model:value="formData.doc_type"
              placeholder="Select document type"
              :disabled="loading"
            >
              <a-select-option value="NB">NB</a-select-option>
              <a-select-option value="UB">UB</a-select-option>
              <a-select-option value="KB">KB</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Plant" name="plant">
            <a-input
              v-model:value="formData.plant"
              placeholder="Enter plant code"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Material Group" name="material_group">
            <a-input
              v-model:value="formData.material_group"
              placeholder="Enter material group"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Material" name="material">
            <a-input
              v-model:value="formData.material"
              placeholder="Enter material code"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Unit Type" name="unit_type">
            <a-select
              v-model:value="formData.unit_type"
              placeholder="Select unit type"
              :disabled="loading"
            >
              <a-select-option value="EA">EA</a-select-option>
              <a-select-option value="KG">KG</a-select-option>
              <a-select-option value="L">L</a-select-option>
              <a-select-option value="M">M</a-select-option>
              <a-select-option value="PC">PC</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Quantity" name="quantity">
            <a-input-number
              v-model:value="formData.quantity"
              :min="1"
              :precision="0"
              placeholder="Enter quantity"
              style="width: 100%"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Unit Price" name="unit_price">
            <a-input-number
              v-model:value="formData.unit_price"
              :min="0"
              :precision="2"
              placeholder="Enter unit price"
              style="width: 100%"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Total Amount" name="total_amount">
            <a-input-number
              v-model:value="formData.total_amount"
              :min="0"
              :precision="2"
              placeholder="Auto-calculated"
              style="width: 100%"
              :disabled="true"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Currency" name="currency">
            <a-select
              v-model:value="formData.currency"
              placeholder="Select currency"
              :disabled="loading"
            >
              <a-select-option value="CNY">CNY</a-select-option>
              <a-select-option value="USD">USD</a-select-option>
              <a-select-option value="EUR">EUR</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Vendor Code" name="vendor_code">
            <a-input
              v-model:value="formData.vendor_code"
              placeholder="Enter vendor code"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Purchase Organization" name="purchase_organization">
            <a-input
              v-model:value="formData.purchase_organization"
              placeholder="Enter purchase organization"
              :disabled="loading"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="Priority" name="priority">
            <a-select
              v-model:value="formData.priority"
              placeholder="Select priority"
              :disabled="loading"
            >
              <a-select-option :value="1">1 - Critical</a-select-option>
              <a-select-option :value="2">2 - High</a-select-option>
              <a-select-option :value="3">3 - Medium</a-select-option>
              <a-select-option :value="4">4 - Low</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="Urgency" name="urgency">
            <a-select
              v-model:value="formData.urgency"
              placeholder="Select urgency"
              :disabled="loading"
            >
              <a-select-option value="critical">Critical</a-select-option>
              <a-select-option value="urgent">Urgent</a-select-option>
              <a-select-option value="normal">Normal</a-select-option>
              <a-select-option value="low">Low</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>

      <a-form-item label="Delivery Date" name="delivery_date">
        <a-date-picker
          v-model:value="formData.delivery_date"
          show-time
          format="YYYY-MM-DD HH:mm:ss"
          placeholder="Select delivery date"
          style="width: 100%"
          :disabled="loading"
        />
      </a-form-item>

      <a-form-item label="Requester" name="requester">
        <a-input
          v-model:value="formData.requester"
          placeholder="Enter requester name"
          :disabled="loading"
        />
      </a-form-item>

      <a-form-item label="Short Text" name="short_text">
        <a-textarea
          v-model:value="formData.short_text"
          placeholder="Enter short description"
          :rows="3"
          :disabled="loading"
        />
      </a-form-item>

      <a-form-item label="Comments" name="comments">
        <a-textarea
          v-model:value="formData.comments"
          placeholder="Enter additional comments"
          :rows="2"
          :disabled="loading"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { purchaseApi } from '@/api/purchase'
import type { CreatePurchaseRequestRequest } from '@/types/purchase'
import dayjs from 'dayjs'

interface Props {
  open: boolean
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref()
const loading = ref(false)

// Form data
const formData = reactive<CreatePurchaseRequestRequest>({
  request_id: '',
  doc_type: 'NB',
  plant: '',
  quantity: 1,
  unit_price: 0,
  material: '',
  delivery_date: '',
  vendor_code: '',
  short_text: '',
  material_group: '',
  unit_type: 'EA',
  requester: '',
  purchase_organization: '',
  currency: 'CNY',
  priority: 2,
  urgency: 'normal',
  comments: ''
})

// Computed total amount
const totalAmount = computed(() => {
  return formData.quantity * formData.unit_price
})

// Watch for total amount calculation
watch([() => formData.quantity, () => formData.unit_price], () => {
  formData.total_amount = totalAmount.value
})

// Form validation rules
const rules = {
  doc_type: [{ required: true, message: 'Please select document type' }],
  plant: [{ required: true, message: 'Please enter plant code' }],
  quantity: [
    { required: true, message: 'Please enter quantity' },
    { type: 'number', min: 1, message: 'Quantity must be at least 1' }
  ],
  unit_price: [
    { required: true, message: 'Please enter unit price' },
    { type: 'number', min: 0, message: 'Unit price must be non-negative' }
  ],
  material: [{ required: true, message: 'Please enter material code' }],
  delivery_date: [{ required: true, message: 'Please select delivery date' }],
  vendor_code: [{ required: true, message: 'Please enter vendor code' }],
  short_text: [{ required: true, message: 'Please enter short description' }],
  material_group: [{ required: true, message: 'Please enter material group' }],
  unit_type: [{ required: true, message: 'Please select unit type' }],
  requester: [{ required: true, message: 'Please enter requester name' }],
  purchase_organization: [{ required: true, message: 'Please enter purchase organization' }],
  priority: [{ required: true, message: 'Please select priority' }],
  urgency: [{ required: true, message: 'Please select urgency' }]
}

// Modal visibility
const visible = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})

// Handle form submission
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    // Prepare data for submission
    const submitData = {
      ...formData,
      delivery_date: formData.delivery_date ? dayjs(formData.delivery_date).format('YYYY-MM-DDTHH:mm:ssZ') : '',
      total_amount: totalAmount.value
    }
    
    // Remove empty request_id to let backend generate it
    if (!submitData.request_id) {
      delete submitData.request_id
    }
    
    await purchaseApi.createPurchaseRequest(submitData)
    
    message.success('Purchase request created successfully!')
    emit('success')
    handleCancel()
    
  } catch (error: any) {
    console.error('Create PR failed:', error)
    message.error(error.response?.data?.error?.message || 'Failed to create purchase request')
  } finally {
    loading.value = false
  }
}

// Handle cancel
const handleCancel = () => {
  // Reset form
  Object.assign(formData, {
    request_id: '',
    doc_type: 'NB',
    plant: '',
    quantity: 1,
    unit_price: 0,
    material: '',
    delivery_date: '',
    vendor_code: '',
    short_text: '',
    material_group: '',
    unit_type: 'EA',
    requester: '',
    purchase_organization: '',
    currency: 'CNY',
    priority: 2,
    urgency: 'normal',
    comments: ''
  })
  
  formRef.value?.resetFields()
  visible.value = false
}
</script>

<style scoped>
.ant-form-item {
  margin-bottom: 16px;
}

.ant-input-number {
  width: 100%;
}
</style>
