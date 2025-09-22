import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: string
  username: string
  name: string
  email: string
  role: 'admin' | 'purchaser' | 'requester' | 'approver'
  department: string
  avatar?: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  
  const userRole = computed(() => user.value?.role || 'requester')

  const initializeUser = () => {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (error) {
        console.error('Failed to parse saved user data:', error)
        localStorage.removeItem('user')
        localStorage.removeItem('token')
      }
    }
  }

  const login = async (username: string, password: string) => {
    try {
      // 模拟LoginAPI调用
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      if (!response.ok) {
        throw new Error('Login失败')
      }

      const data = await response.json()
      
      // 模拟User数据
      const mockUser: User = {
        id: '1',
        username,
        name: username === 'admin' ? 'System管理员' : username === 'purchaser' ? 'Purchaser' : 'NormalUser',
        email: `${username}@company.com`,
        role: username === 'admin' ? 'admin' : username === 'purchaser' ? 'purchaser' : 'requester',
        department: username === 'admin' ? 'ITDepartment' : username === 'purchaser' ? '采购部' : '业务Department',
        avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
      }

      user.value = mockUser
      token.value = data.token || 'mock-token-' + Date.now()
      
      localStorage.setItem('user', JSON.stringify(mockUser))
      localStorage.setItem('token', token.value)
      
      return { success: true }
    } catch (error) {
      console.error('Login error:', error)
      return { success: false, error: error instanceof Error ? error.message : 'Login失败' }
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('user')
    localStorage.removeItem('token')
  }

  const hasPermission = (permission: string) => {
    if (!user.value) return false
    
    const rolePermissions: Record<string, string[]> = {
      admin: ['*'],
      purchaser: ['pr:read', 'pr:create', 'po:read', 'po:create', 'supplier:read', 'supplier:create', 'process:read'],
      approver: ['pr:read', 'pr:approve', 'po:read', 'process:read'],
      requester: ['pr:read', 'pr:create', 'po:read', 'process:read']
    }

    const permissions = rolePermissions[user.value.role] || []
    return permissions.includes('*') || permissions.includes(permission)
  }

  return {
    user,
    token,
    isAuthenticated,
    userRole,
    initializeUser,
    login,
    logout,
    hasPermission
  }
})
