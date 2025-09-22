<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <img
          src="/logo.svg"
          alt="Logo"
          class="login-logo"
        >
        <h1 class="login-title">
          订单Medium心
        </h1>
        <p class="login-subtitle">
          智能采购管理System
        </p>
      </div>
      
      <a-form
        :model="loginForm"
        :rules="rules"
        class="login-form"
        layout="vertical"
        @finish="handleLogin"
      >
        <a-form-item
          label="Username"
          name="username"
        >
          <a-input
            v-model:value="loginForm.username"
            placeholder="Please enterUsername"
            size="large"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item
          label="Password"
          name="password"
        >
          <a-input-password
            v-model:value="loginForm.password"
            placeholder="Please enterPassword"
            size="large"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item>
          <a-checkbox v-model:checked="rememberMe">
            Remember me
          </a-checkbox>
        </a-form-item>
        
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            :loading="loading"
            block
          >
            Login
          </a-button>
        </a-form-item>
      </a-form>
      
      <div class="login-footer">
        <p class="demo-accounts">
          <strong>演示账号：</strong><br>
          管理员：admin / admin123<br>
          Purchaser：purchaser / purchaser123<br>
          NormalUser：user / user123
        </p>
      </div>
    </div>
    
    <!-- 背景装饰 -->
    <div class="login-background">
      <div class="bg-shape shape-1" />
      <div class="bg-shape shape-2" />
      <div class="bg-shape shape-3" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// Reactive data
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

// Form validation rules
const rules = {
  username: [
    { required: true, message: 'Please enterUsername', trigger: 'blur' },
    { min: 3, max: 20, message: 'Username长度在3到20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enterPassword', trigger: 'blur' },
    { min: 6, max: 20, message: 'Password长度在6到20个字符', trigger: 'blur' }
  ]
}

// Login处理
const handleLogin = async () => {
  loading.value = true
  
  try {
    const result = await userStore.login(loginForm.username, loginForm.password)
    
    if (result.success) {
      message.success('Login成功')
      router.push('/')
    } else {
      message.error(result.error || 'Login失败')
    }
  } catch (error) {
    message.error('Login失败，Please try again later')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-box {
  width: 400px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  position: relative;
  z-index: 10;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  height: 48px;
  margin-bottom: 16px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 8px 0;
}

.login-subtitle {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.login-form {
  margin-bottom: 24px;
}

.login-footer {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.demo-accounts {
  font-size: 12px;
  color: #666;
  margin: 0;
  line-height: 1.6;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

/* Responsive design */
@media (max-width: 480px) {
  .login-box {
    width: 90%;
    padding: 24px;
  }
  
  .login-title {
    font-size: 24px;
  }
}
</style>
