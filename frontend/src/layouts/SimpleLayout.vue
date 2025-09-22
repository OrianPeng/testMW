<template>
  <a-layout class="main-layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      :width="240"
      :collapsed-width="80"
      class="layout-sider"
    >
      <div class="logo">
        <img v-if="!collapsed" src="/logo.svg" alt="Logo" class="logo-img" />
        <img v-else src="/logo-mini.svg" alt="Logo" class="logo-mini" />
        <span v-if="!collapsed" class="logo-text">Order Center</span>
      </div>
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        :inline-collapsed="collapsed"
        class="sidebar-menu"
        @click="handleMenuClick"
      >
        <a-menu-item key="/chat">
          <template #icon>
            <MessageOutlined />
          </template>
          <span>AI Assistant</span>
        </a-menu-item>
        
        <a-menu-item key="/pr-query">
          <template #icon>
            <FileSearchOutlined />
          </template>
          <span>PR Query</span>
        </a-menu-item>
        
        <a-menu-item key="/po-query">
          <template #icon>
            <ShoppingOutlined />
          </template>
          <span>PO Query</span>
        </a-menu-item>
        
        <a-menu-item key="/process-monitor">
          <template #icon>
            <MonitorOutlined />
          </template>
          <span>Process Monitor</span>
        </a-menu-item>
        
        <a-menu-item key="/supplier">
          <template #icon>
            <TeamOutlined />
          </template>
          <span>Supplier Management</span>
        </a-menu-item>
        
        <a-menu-item key="/dashboard">
          <template #icon>
            <DashboardOutlined />
          </template>
          <span>Dashboard</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <!-- 主内容区域 -->
    <a-layout>
      <!-- 顶部导航栏 -->
      <a-layout-header class="layout-header">
        <div class="header-left">
          <a-button
            type="text"
            @click="toggleCollapsed"
            class="trigger"
          >
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </a-button>
          
          <!-- Global Search -->
          <a-input-search
            v-model:value="searchValue"
            placeholder="Search PR/PO numbers..."
            style="width: 300px; margin-left: 16px"
            @search="handleSearch"
            allow-clear
            class="header-search"
          />
        </div>
        
        <div class="header-right">
          <!-- User Info -->
          <div class="user-info">
            <a-avatar size="32">U</a-avatar>
            <span class="user-name">User</span>
          </div>
        </div>
      </a-layout-header>

      <!-- Page Content -->
      <a-layout-content class="layout-content">
        <div class="content-wrapper">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MessageOutlined,
  FileSearchOutlined,
  ShoppingOutlined,
  MonitorOutlined,
  TeamOutlined,
  DashboardOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()

// Reactive data
const collapsed = ref(false)
const selectedKeys = ref<string[]>([])
const searchValue = ref('')

// 监听路由变化，更新选Medium的菜单项
watch(
  () => route.path,
  (newPath) => {
    selectedKeys.value = [newPath]
  },
  { immediate: true }
)

// Methods
const toggleCollapsed = () => {
  collapsed.value = !collapsed.value
}

const handleMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}

const handleSearch = (value: string) => {
  if (!value.trim()) return
  
  // 简单的搜索逻辑
  if (value.startsWith('PR-') || value.startsWith('pr-')) {
    router.push(`/pr-query?search=${encodeURIComponent(value)}`)
  } else if (value.startsWith('PO-') || value.startsWith('po-')) {
    router.push(`/po-query?search=${encodeURIComponent(value)}`)
  } else {
    message.info('Please enter a valid PR or PO number')
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.layout-sider {
  background: linear-gradient(180deg, #1e3a8a 0%, #1e40af 50%, #1d4ed8 100%);
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.15);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
}

.logo-img {
  height: 32px;
  margin-right: 8px;
  filter: brightness(0) invert(1);
}

.logo-mini {
  height: 24px;
  filter: brightness(0) invert(1);
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.sidebar-menu {
  border-right: none;
  background: transparent;
  padding: 16px 8px;
}

.sidebar-menu :deep(.ant-menu-item) {
  color: rgba(255, 255, 255, 0.8);
  margin: 4px 0;
  border-radius: 12px;
  height: 48px;
  line-height: 48px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.sidebar-menu :deep(.ant-menu-item::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, #3b82f6, #8b5cf6);
  opacity: 0;
  transition: opacity 0.3s ease;
  border-radius: 12px;
}

.sidebar-menu :deep(.ant-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.sidebar-menu :deep(.ant-menu-item:hover::before) {
  opacity: 0.1;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background: linear-gradient(45deg, #3b82f6, #8b5cf6);
  color: #fff;
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
  transform: translateX(4px);
}

.sidebar-menu :deep(.ant-menu-item-selected::after) {
  display: none;
}

.sidebar-menu :deep(.ant-menu-item .anticon) {
  font-size: 16px;
  margin-right: 12px;
}

.layout-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  z-index: 10;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  height: 64px;
  min-height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #64748b;
  padding: 8px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 40px;
  width: 40px;
}

.trigger:hover {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  transform: scale(1.05);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 12px;
  transition: all 0.3s ease;
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.user-info:hover {
  background: rgba(59, 130, 246, 0.1);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
}

.user-name {
  margin-left: 8px;
  font-weight: 600;
  color: #1e293b;
}

.layout-content {
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  overflow: auto;
  position: relative;
}

.layout-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%23e2e8f0" opacity="0.3"/><circle cx="75" cy="75" r="1" fill="%23cbd5e1" opacity="0.2"/><circle cx="50" cy="10" r="0.5" fill="%23f1f5f9" opacity="0.4"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  opacity: 0.3;
  pointer-events: none;
}

.content-wrapper {
  padding: 24px;
  min-height: calc(100vh - 64px);
  position: relative;
  z-index: 1;
}

/* 搜索框样式优化 */
.header-search {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  transition: all 0.3s ease;
}

.header-search:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
}

.header-search:focus-within {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.header-search .ant-input {
  height: 40px;
  border: none;
  box-shadow: none;
}

.header-search .ant-input-search-button {
  height: 40px;
  border: none;
  box-shadow: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .layout-sider {
    position: fixed !important;
    height: 100vh;
    left: 0;
    top: 0;
    z-index: 1000;
  }
  
  .layout-content {
    margin-left: 0 !important;
  }
  
  .header-search {
    width: 200px !important;
  }
  
  .logo-text {
    font-size: 16px;
  }
  
  .sidebar-menu :deep(.ant-menu-item) {
    height: 44px;
    line-height: 44px;
    font-size: 14px;
  }
}

/* 滚动条样式 */
.layout-content::-webkit-scrollbar {
  width: 6px;
}

.layout-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
}

.layout-content::-webkit-scrollbar-thumb {
  background: rgba(59, 130, 246, 0.3);
  border-radius: 3px;
}

.layout-content::-webkit-scrollbar-thumb:hover {
  background: rgba(59, 130, 246, 0.5);
}
</style>
