import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  username: string
  token: string
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const isAuthenticated = ref(false)
  const isLoading = ref(false)
  const error = ref('')
  const httpAuthEnabled = ref(false)

  // Initialize from storage
  const token = localStorage.getItem('auth_token') || sessionStorage.getItem('auth_token')
  const username = localStorage.getItem('auth_username') || sessionStorage.getItem('auth_username')
  if (token && username) {
    user.value = { username, token }
    isAuthenticated.value = true
  }

  // Actions
  async function checkHttpAuthEnabled() {
    try {
      const response = await fetch('/api/settings')
      const data = await response.json()
      httpAuthEnabled.value = data.httpAuthOn || false
      return httpAuthEnabled.value
    } catch (error) {
      console.error('Failed to check HTTP auth:', error)
      return false
    }
  }

  async function login(username: string, password: string, remember: boolean) {
    isLoading.value = true
    error.value = ''

    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password, remember })
      })

      const data = await response.json()

      if (response.ok) {
        const token = data.token
        if (remember) {
          localStorage.setItem('auth_token', token)
          localStorage.setItem('auth_username', username)
        } else {
          sessionStorage.setItem('auth_token', token)
          sessionStorage.setItem('auth_username', username)
        }
        user.value = { username, token }
        isAuthenticated.value = true
        return { success: true, redirect: data.redirect || '/' }
      } else {
        error.value = data.message || 'Login failed'
        return { success: false }
      }
    } catch (err) {
      error.value = 'Network error, please try again'
      return { success: false }
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_username')
    sessionStorage.removeItem('auth_token')
    sessionStorage.removeItem('auth_username')
    user.value = null
    isAuthenticated.value = false
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    httpAuthEnabled,
    checkHttpAuthEnabled,
    login,
    logout
  }
})
