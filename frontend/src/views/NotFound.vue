<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="not-found-icon">
        <ExclamationCircleOutlined />
      </div>
      <h1 class="not-found-title">
        404
      </h1>
      <p class="not-found-description">
        Sorry, the page you visited does not exist
      </p>
      <a-button
        type="primary"
        @click="goHome"
      >
        Back to Home
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.not-found {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.not-found-content {
  text-align: center;
  padding: 40px;
}

.not-found-icon {
  font-size: 120px;
  color: #d9d9d9;
  margin-bottom: 24px;
}

.not-found-title {
  font-size: 72px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 16px 0;
}

.not-found-description {
  font-size: 16px;
  color: #666;
  margin: 0 0 32px 0;
}
</style>
