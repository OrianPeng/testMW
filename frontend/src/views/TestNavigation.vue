<template>
  <div class="test-navigation">
    <div class="test-header">
      <h1>🎉 Secondary Interface Test Successful!</h1>
      <p>If you can see this page, it means the route navigation is working properly</p>
    </div>
    
    <div class="test-content">
      <a-row :gutter="16">
        <a-col :span="8">
          <a-card class="test-card" :bordered="false">
            <template #title>
              <CheckCircleOutlined style="color: #52c41a" />
              <span>Route Navigation</span>
            </template>
            <p>✅ Sidebar menu clicks work properly</p>
            <p>✅ Page switching works properly</p>
            <p>✅ Route parameter passing works properly</p>
          </a-card>
        </a-col>
        
        <a-col :span="8">
          <a-card class="test-card" :bordered="false">
            <template #title>
              <CheckCircleOutlined style="color: #52c41a" />
              <span>Layout System</span>
            </template>
            <p>✅ Main layout works properly</p>
            <p>✅ Sidebar works properly</p>
            <p>✅ Content area works properly</p>
          </a-card>
        </a-col>
        
        <a-col :span="8">
          <a-card class="test-card" :bordered="false">
            <template #title>
              <CheckCircleOutlined style="color: #52c41a" />
              <span>Design System</span>
            </template>
            <p>✅ Modern design</p>
            <p>✅ Responsive layout</p>
            <p>✅ Animation effects</p>
          </a-card>
        </a-col>
      </a-row>
      
      <div class="test-actions">
        <a-space size="large">
          <a-button type="primary" @click="goToPRQuery">
            <FileSearchOutlined />
            View PR Query
          </a-button>
          <a-button @click="goToDashboard">
            <DashboardOutlined />
            View Dashboard
          </a-button>
          <a-button @click="goToChat">
            <MessageOutlined />
            View AI Assistant
          </a-button>
        </a-space>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import {
  CheckCircleOutlined,
  FileSearchOutlined,
  DashboardOutlined,
  MessageOutlined
} from '@ant-design/icons-vue'

const router = useRouter()

const goToPRQuery = () => {
  router.push('/pr-query')
}

const goToDashboard = () => {
  router.push('/dashboard')
}

const goToChat = () => {
  router.push('/chat')
}
</script>

<style scoped>
.test-navigation {
  padding: 24px;
  background: transparent;
}

.test-header {
  text-align: center;
  margin-bottom: 32px;
  padding: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
}

.test-header h1 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 16px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.test-header p {
  font-size: 18px;
  margin: 0;
  opacity: 0.9;
}

.test-content {
  margin-bottom: 32px;
}

.test-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  background: white;
  transition: all 0.3s ease;
  height: 200px;
}

.test-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.test-card :deep(.ant-card-head) {
  border-bottom: 2px solid #f0f0f0;
  background: #fafafa;
  border-radius: 16px 16px 0 0;
}

.test-card :deep(.ant-card-head-title) {
  display: flex;
  align-items: center;
  font-weight: 600;
  color: #1f2937;
}

.test-card :deep(.ant-card-head-title .anticon) {
  margin-right: 8px;
  font-size: 16px;
}

.test-card :deep(.ant-card-body) {
  padding: 24px;
}

.test-card p {
  margin: 8px 0;
  color: #6b7280;
  font-size: 14px;
}

.test-actions {
  text-align: center;
  padding: 32px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.test-actions .ant-btn {
  height: 48px;
  padding: 0 24px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s ease;
}

.test-actions .ant-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.test-actions .ant-btn-primary {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  border: none;
}

.test-actions .ant-btn-primary:hover {
  background: linear-gradient(135deg, #1d4ed8, #1e40af);
}

/* Responsive design */
@media (max-width: 768px) {
  .test-navigation {
    padding: 16px;
  }
  
  .test-header {
    padding: 24px 20px;
  }
  
  .test-header h1 {
    font-size: 24px;
  }
  
  .test-header p {
    font-size: 16px;
  }
  
  .test-actions {
    padding: 24px;
  }
  
  .test-actions .ant-space {
    flex-direction: column;
    width: 100%;
  }
  
  .test-actions .ant-btn {
    width: 100%;
    margin-bottom: 12px;
  }
}
</style>
