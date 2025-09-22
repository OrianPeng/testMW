<template>
  <div class="test-api-page">
    <h1>API Test Page</h1>
    
    <div class="test-section">
      <h2>Backend Connection Test</h2>
      <a-button @click="testConnection" :loading="loading">
        Test Backend Connection
      </a-button>
      <div v-if="connectionResult" class="result">
        <pre>{{ connectionResult }}</pre>
      </div>
    </div>

    <div class="test-section">
      <h2>Purchase Requests API</h2>
      <a-button @click="testPurchaseRequests" :loading="loading">
        Test Purchase Requests
      </a-button>
      <div v-if="prResult" class="result">
        <pre>{{ prResult }}</pre>
      </div>
    </div>

    <div class="test-section">
      <h2>Create PR Test</h2>
      <a-button @click="testCreatePR" :loading="loading">
        Test Create PR
      </a-button>
      <div v-if="createResult" class="result">
        <pre>{{ createResult }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { purchaseApi } from '@/api/purchase'

const loading = ref(false)
const connectionResult = ref('')
const prResult = ref('')
const createResult = ref('')

const testConnection = async () => {
  loading.value = true
  try {
    const response = await purchaseApi.healthCheck()
    connectionResult.value = JSON.stringify(response, null, 2)
  } catch (error: any) {
    connectionResult.value = `Error: ${error.message}`
  } finally {
    loading.value = false
  }
}

const testPurchaseRequests = async () => {
  loading.value = true
  try {
    const response = await purchaseApi.getPurchaseRequests()
    prResult.value = JSON.stringify(response, null, 2)
  } catch (error: any) {
    prResult.value = `Error: ${error.message}`
  } finally {
    loading.value = false
  }
}

const testCreatePR = async () => {
  loading.value = true
  try {
    const testData = {
      request_id: `TEST-${Date.now()}`,
      doc_type: 'NB',
      plant: '1000',
      quantity: 10,
      unit_price: 50.0,
      material: 'TEST-MAT-001',
      delivery_date: '2024-01-15T00:00:00Z',
      vendor_code: 'TEST-VENDOR-001',
      short_text: 'Test Purchase Request',
      material_group: 'TEST-MG-001',
      unit_type: 'EA',
      requester: 'Test User',
      purchase_organization: 'TEST-PO-001',
      currency: 'CNY',
      priority: 2,
      urgency: 'normal',
      comments: 'Test comment'
    }
    
    const response = await purchaseApi.createPurchaseRequest(testData)
    createResult.value = JSON.stringify(response, null, 2)
  } catch (error: any) {
    createResult.value = `Error: ${error.message}`
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.test-api-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.test-section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
}

.result {
  margin-top: 10px;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
