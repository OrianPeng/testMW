<template>
  <div class="test-page">
    <h1>Order Center Test Page</h1>
    <p>If you can see this page, the Vue application has started successfully!</p>
    <a-button type="primary" @click="showMessage">点击测试</a-button>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue'

const showMessage = () => {
  message.success('Vue应用运行正常！')
}
</script>

<style scoped>
.test-page {
  padding: 50px;
  text-align: center;
}

h1 {
  color: #1890ff;
  margin-bottom: 20px;
}

p {
  font-size: 16px;
  color: #666;
  margin-bottom: 30px;
}
</style>
