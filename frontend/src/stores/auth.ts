import { defineStore } from 'pinia'
import { ref } from 'vue'
import router from '@/router'
import { login as apiLogin, logout as apiLogout, changePassword as apiChangePassword, getCurrentUser } from '@/api/auth'

export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  is_admin: boolean
  is_active: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const isAuthenticated = ref(!!token.value)

  async function login(username: string, password: string, keepLogin = false) {
    const result = await apiLogin(username, password, keepLogin)
    token.value = result.token
    user.value = result.user
    isAuthenticated.value = true
    localStorage.setItem('token', result.token)
  }

  async function logout() {
    await apiLogout()
    token.value = ''
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem('token')
    router.push('/login')
  }

  async function checkAuth() {
    if (!token.value) {
      isAuthenticated.value = false
      user.value = null
      return
    }
    try {
      const userData = await getCurrentUser()
      user.value = userData
      isAuthenticated.value = true
    } catch {
      token.value = ''
      user.value = null
      isAuthenticated.value = false
      localStorage.removeItem('token')
    }
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    return apiChangePassword(oldPassword, newPassword)
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    checkAuth,
    changePassword,
  }
})
