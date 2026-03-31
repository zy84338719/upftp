<template>
  <div class="login-page">
    <div class="login-container">
      <div class="logo">
        <div class="logo-mark">U</div>
        <span class="logo-text">UPFTP</span>
      </div>

      <div v-if="error" class="error show">{{ error }}</div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">{{ $t('username') }}</label>
          <input
            type="text"
            id="username"
            v-model="username"
            :placeholder="$t('username')"
            required
            autofocus
            autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label for="password">{{ $t('password') }}</label>
          <input
            type="password"
            id="password"
            v-model="password"
            :placeholder="$t('password')"
            required
            autocomplete="current-password"
          />
        </div>

        <div class="remember-me">
          <input type="checkbox" id="remember" v-model="remember" />
          <label for="remember">{{ $t('remember_me') }}</label>
        </div>

        <button type="submit" :disabled="isLoading">
          <span v-if="!isLoading">{{ $t('login_btn') }}</span>
          <div v-else class="loading show">
            <div class="spinner"></div>
          </div>
        </button>
      </form>

      <div class="footer">
        <p>{{ $t('secure_connection') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const remember = ref(false)
const error = ref('')
const isLoading = ref(false)

onMounted(() => {
  const rememberedUsername = localStorage.getItem('auth_username')
  if (rememberedUsername) {
    username.value = rememberedUsername
    remember.value = true
  }
})

async function handleLogin() {
  isLoading.value = true
  error.value = ''

  const result = await authStore.login(username.value, password.value, remember.value)

  if (result.success) {
    router.push(result.redirect || '/')
  } else {
    error.value = authStore.error
  }

  isLoading.value = false
}
</script>

<style scoped>
.login-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  background: #f8f9fa;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
  padding: 20px;
}

.login-container {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 400px;
  border: 1px solid #e5e5e5;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 32px;
}

.logo-mark {
  width: 32px;
  height: 32px;
  background: #d97706;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 700;
  font-size: 16px;
  border-radius: 6px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: #1a1a1a;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 6px;
  color: #555;
  font-weight: 500;
  font-size: 13px;
}

input[type='text'],
input[type='password'] {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.15s;
  background: #fff;
  outline: none;
}

input[type='text']:focus,
input[type='password']:focus {
  border-color: #999;
}

.remember-me {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.remember-me input[type='checkbox'] {
  margin-right: 8px;
}

.remember-me label {
  margin: 0;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  font-weight: 400;
}

button {
  width: 100%;
  padding: 10px 16px;
  background: #d97706;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
  margin-top: 0;
}

button:hover:not(:disabled) {
  background: #b45309;
}

button:active:not(:disabled) {
  background: #92400e;
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error {
  background: #fee;
  color: #c33;
  padding: 10px 14px;
  border-radius: 6px;
  margin-bottom: 20px;
  display: none;
  border: 1px solid #fcc;
  font-size: 13px;
}

.error.show {
  display: block;
}

.loading {
  display: none;
  text-align: center;
  padding: 0;
}

.loading.show {
  display: block;
}

.spinner {
  border: 2px solid #f3f3f3;
  border-top: 2px solid #fff;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
  display: inline-block;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #eee;
  color: #999;
  font-size: 12px;
}

@media (max-width: 480px) {
  .login-container {
    padding: 30px 20px;
  }

  .logo-text {
    font-size: 18px;
  }

  button {
    padding: 10px 14px;
  }
}
</style>
