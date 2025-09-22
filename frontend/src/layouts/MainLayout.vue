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
        <img
          v-if="!collapsed"
          src="/logo.svg"
          alt="Logo"
          class="logo-img"
        >
        <img
          v-else
          src="/logo-mini.svg"
          alt="Logo"
          class="logo-mini"
        >
        <span
          v-if="!collapsed"
          class="logo-text"
        >订单Medium心</span>
      </div>
      
      <a-menu
        v-model:selected-keys="selectedKeys"
        mode="inline"
        :inline-collapsed="collapsed"
        class="sidebar-menu"
      >
        <a-menu-item key="/chat">
          <template #icon>
            <MessageOutlined />
          </template>
          <span>智能Assistant</span>
        </a-menu-item>
        
        <a-menu-item key="/pr-query">
          <template #icon>
            <FileSearchOutlined />
          </template>
          <span>PRSearch</span>
        </a-menu-item>
        
        <a-menu-item key="/po-query">
          <template #icon>
            <ShoppingOutlined />
          </template>
          <span>POSearch</span>
        </a-menu-item>
        
        <a-menu-item key="/process-monitor">
          <template #icon>
            <MonitorOutlined />
          </template>
          <span>Process监控</span>
        </a-menu-item>
        
        <a-menu-item key="/supplier">
          <template #icon>
            <TeamOutlined />
          </template>
          <span>Supplier管理</span>
        </a-menu-item>
        
        <a-menu-item key="/dashboard">
          <template #icon>
            <DashboardOutlined />
          </template>
          <span>数据看板</span>
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
            class="trigger"
            @click="toggleCollapsed"
          >
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </a-button>
          
          <!-- 全局搜索 -->
          <a-input-search
            v-model:value="searchValue"
            placeholder="搜索PR/PO Number..."
            style="width: 300px; margin-left: 16px"
            allow-clear
            @search="handleSearch"
          />
        </div>
        
        <div class="header-right">
          <!-- 通知铃铛 -->
          <a-badge
            :count="notificationQuantity"
            :offset="[10, 0]"
          >
            <a-button
              type="text"
              @click="showNotifications"
            >
              <BellOutlined style="font-size: 16px" />
            </a-button>
          </a-badge>
          
          <!-- User头像和下拉菜单 -->
          <a-dropdown placement="bottomRight">
            <div class="user-info">
              <a-avatar
                :src="userStore.user?.avatar"
                :size="32"
              >
                {{ userStore.user?.name?.charAt(0) }}
              </a-avatar>
              <span class="user-name">{{ userStore.user?.name }}</span>
              <DownOutlined style="font-size: 12px; margin-left: 4px" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile">
                  <UserOutlined />
                  个人资料
                </a-menu-item>
                <a-menu-item key="settings">
                  <SettingOutlined />
                  System设置
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item
                  key="logout"
                  @click="handleLogout"
                >
                  <LogoutOutlined />
                  退出Login
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- page面内容 -->
      <a-layout-content class="layout-content">
        <div class="content-wrapper">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>

    <!-- 通知抽屉 -->
    <a-drawer
      v-model:open="notificationVisible"
      title="System通知"
      placement="right"
      :width="400"
    >
      <div class="notification-list">
        <a-empty
          v-if="notifications.length === 0"
          description="暂None通知"
        />
        <div v-else>
          <div
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item"
            :class="{ unread: !notification.read }"
          >
            <div class="notification-content">
              <div class="notification-title">
                {{ notification.title }}
              </div>
              <div class="notification-message">
                {{ notification.message }}
              </div>
              <div class="notification-time">
                {{ formatTime(notification.time) }}
              </div>
            </div>
            <a-tag
              v-if="!notification.read"
              color="blue"
            >
              新
            </a-tag>
          </div>
        </div>
      </div>
    </a-drawer>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MessageOutlined,
  FileSearchOutlined,
  ShoppingOutlined,
  MonitorOutlined,
  TeamOutlined,
  DashboardOutlined,
  BellOutlined,
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// Reactive data
const collapsed = ref(false)
const selectedKeys = ref<string[]>([])
const searchValue = ref('')
const notificationVisible = ref(false)
const notificationQuantity = ref(3)

// 通知数据
const notifications = ref([
  {
    id: '1',
    title: 'PRApprove提醒',
    message: '您有3个PR待Approve',
    time: new Date(),
    read: false
  },
  {
    id: '2',
    title: 'POStatus更新',
    message: 'PO-20231201-001 已Send to Supplier',
    time: new Date(Date.now() - 1000 * 60 * 30),
    read: false
  },
  {
    id: '3',
    title: 'System维护通知',
    message: 'System将于今晚22:00-24:00进行维护',
    time: new Date(Date.now() - 1000 * 60 * 60 * 2),
    read: true
  }
])

// 计算属性
const unreadQuantity = computed(() => 
  notifications.value.filter(n => !n.read).length
)

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

const handleSearch = (value: string) => {
  if (!value.trim()) return
  
  // 简单的搜索逻辑，可以根据实际需求扩展
  if (value.startsWith('PR-') || value.startsWith('pr-')) {
    router.push(`/pr-query?search=${encodeURIComponent(value)}`)
  } else if (value.startsWith('PO-') || value.startsWith('po-')) {
    router.push(`/po-query?search=${encodeURIComponent(value)}`)
  } else {
    message.info('Please enter有效的PR或PO Number')
  }
}

const showNotifications = () => {
  notificationVisible.value = true
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
  message.success('已退出Login')
}

const formatTime = (time: Date) => {
  return dayjs(time).fromNow()
}

// Lifecycle
onMounted(() => {
  // 初始化通知Quantity
  notificationQuantity.value = unreadQuantity.value
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.layout-sider {
  background: #001529;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid #002140;
}

.logo-img {
  height: 32px;
  margin-right: 8px;
}

.logo-mini {
  height: 24px;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
}

.sidebar-menu {
  border-right: none;
  background: #001529;
}

.sidebar-menu :deep(.ant-menu-item) {
  color: rgba(255, 255, 255, 0.65);
  margin: 4px 8px;
  border-radius: 6px;
}

.sidebar-menu :deep(.ant-menu-item:hover) {
  background-color: #1890ff;
  color: #fff;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background-color: #1890ff;
  color: #fff;
}

.sidebar-menu :deep(.ant-menu-item-selected::after) {
  display: none;
}

.layout-header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  line-height: 64px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
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
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.user-name {
  margin-left: 8px;
  font-weight: 500;
}

.layout-content {
  background: #f5f5f5;
  overflow: auto;
}

.content-wrapper {
  padding: 24px;
  min-height: calc(100vh - 64px);
}

.notification-list {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.notification-item.unread {
  background-color: #f6ffed;
  margin: 0 -16px;
  padding: 12px 16px;
}

.notification-content {
  flex: 1;
}

.notification-title {
  font-weight: 500;
  color: #262626;
  margin-bottom: 4px;
}

.notification-message {
  color: #666;
  font-size: 12px;
  margin-bottom: 4px;
}

.notification-time {
  color: #999;
  font-size: 11px;
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
  
  .header-left .ant-input-search {
    width: 200px !important;
  }
}
</style>
