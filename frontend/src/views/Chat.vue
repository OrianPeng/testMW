<template>
  <div class="chat-page">
    <div class="chat-container">
      <!-- Chat Header -->
      <div class="chat-header">
        <div class="chat-title">
          <RobotOutlined />
          <span>AI Assistant</span>
        </div>
        <div class="chat-actions">
          <a-button
            type="primary"
            @click="showCreatePRModal"
          >
            <FileTextOutlined />
            Create PR
          </a-button>
        </div>
      </div>


      <!-- AI Chatbot iframe -->
      <div class="chat-iframe-container">
        <div v-if="iframeLoading" class="loading-overlay">
          <a-spin size="large" />
          <p>Loading AI Assistant...</p>
        </div>
        <iframe
          ref="chatIframe"
          src="http://localhost/chatbot/HHOn2vBJZJHHTSug"
          style="width: 100%; height: 100%; min-height: 700px"
          frameborder="0"
          allow="microphone"
          @load="onIframeLoad"
          @error="onIframeError"
        />
      </div>
    </div>

    <!-- Create PR Modal -->
    <CreatePRModal
      v-model:open="createPRVisible"
      @success="handlePRCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import {
  RobotOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import CreatePRModal from '@/components/CreatePRModal.vue'

// Reactive data
const createPRVisible = ref(false)
const chatIframe = ref<HTMLIFrameElement>()
const iframeLoading = ref(true)


// Methods


const onIframeLoad = () => {
  console.log('Chatbot iframe loaded successfully')
  iframeLoading.value = false
  antMessage.success('AI Assistant is ready!')
}

const onIframeError = () => {
  console.error('Failed to load chatbot iframe')
  iframeLoading.value = false
  antMessage.error('Failed to load AI Assistant. Please check the chatbot URL.')
}

// Create PR related methods
const showCreatePRModal = () => {
  createPRVisible.value = true
}

const handlePRCreated = () => {
  antMessage.success('Purchase Request created successfully! You can view it in the PR Query page.')
}

// Lifecycle
onMounted(() => {
  // Initialize chat iframe
  console.log('Chat page mounted, iframe will load automatically')
})
</script>

<style scoped>
.chat-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
  flex-shrink: 0;
  min-height: 60px;
}

.chat-title {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  flex-shrink: 0;
}

.chat-title .anticon {
  margin-right: 8px;
  color: #1890ff;
  font-size: 20px;
  line-height: 1;
}

.chat-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-shrink: 0;
}


.chat-iframe-container {
  flex: 1;
  position: relative;
  background: #f5f5f5;
  overflow: hidden;
}

.chat-iframe-container iframe {
  border: none;
  width: 100%;
  height: 100%;
  min-height: 600px;
  background: #fff;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.loading-overlay p {
  margin-top: 16px;
  color: #666;
  font-size: 16px;
}

/* Responsive design */
@media (max-width: 768px) {
  .chat-header {
    padding: 12px 16px;
    min-height: 56px;
  }
  
  .chat-title {
    font-size: 16px;
  }
  
  .chat-title .anticon {
    font-size: 18px;
  }
  
  .chat-iframe-container {
    min-height: 500px;
  }
  
  .chat-iframe-container iframe {
    min-height: 500px;
  }
}

@media (max-width: 600px) {
  .chat-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
    padding: 16px;
    min-height: auto;
  }
  
  .chat-title {
    justify-content: center;
  }
  
  .chat-actions {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .chat-header {
    padding: 12px;
  }
  
  .chat-title {
    font-size: 15px;
  }
  
  .chat-title .anticon {
    font-size: 16px;
  }
}
</style>
