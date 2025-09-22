import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/test',
    name: 'SimpleTest',
    component: () => import('@/views/SimpleTest.vue')
  },
  {
    path: '/test-nav',
    name: 'TestNavigation',
    component: () => import('@/views/TestNavigation.vue')
  },
  {
    path: '/test-api',
    name: 'TestAPI',
    component: () => import('@/views/TestAPI.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/SimpleLayout.vue'),
    children: [
      {
        path: '',
        redirect: '/chat'
      },
      {
        path: 'chat',
        name: 'Chat',
        component: () => import('@/views/Chat.vue'),
        meta: { title: 'AI Assistant', icon: 'MessageOutlined' }
      },
      {
        path: 'pr-query',
        name: 'PRQuery',
        component: () => import('@/views/PRQuery.vue'),
        meta: { title: 'PR Query', icon: 'FileSearchOutlined' }
      },
      {
        path: 'po-query',
        name: 'POQuery',
        component: () => import('@/views/POQuery.vue'),
        meta: { title: 'PO Query', icon: 'ShoppingOutlined' }
      },
      {
        path: 'process-monitor',
        name: 'ProcessMonitor',
        component: () => import('@/views/ProcessMonitor.vue'),
        meta: { title: 'Process Monitor', icon: 'MonitorOutlined' }
      },
      {
        path: 'supplier',
        name: 'Supplier',
        component: () => import('@/views/Supplier.vue'),
        meta: { title: 'Supplier Management', icon: 'TeamOutlined' }
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: 'Dashboard', icon: 'DashboardOutlined' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 简化路由守卫，暂时不进行权限检查
  next()
})

export default router
