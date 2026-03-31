import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true }
    },
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/:path(.*)',
      name: 'files',
      component: HomeView
    }
  ]
})

// 检查 HTTP 认证是否启用的标志
let checkedHttpAuth = false
let httpAuthEnabled = false

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // 如果是公开页面，直接放行
  if (to.meta.public) {
    next()
    return
  }

  // 首次访问时检查 HTTP 认证是否启用
  if (!checkedHttpAuth) {
    try {
      httpAuthEnabled = await authStore.checkHttpAuthEnabled()
      checkedHttpAuth = true
    } catch (error) {
      console.error('Failed to check HTTP auth:', error)
      // 如果检查失败，假设需要认证
      httpAuthEnabled = true
      checkedHttpAuth = true
    }
  }

  // 如果没有启用 HTTP 认证，直接放行
  if (!httpAuthEnabled) {
    next()
    return
  }

  // 如果启用了 HTTP 认证，检查用户是否已登录
  if (!authStore.isAuthenticated) {
    const token = localStorage.getItem('auth_token') || sessionStorage.getItem('auth_token')
    if (!token) {
      next('/login')
      return
    }
  }

  next()
})

export default router
