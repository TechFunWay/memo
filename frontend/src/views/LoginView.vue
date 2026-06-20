<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-400 via-teal-500 to-cyan-600 dark:from-[#0d3320] dark:via-[#0a2e1c] dark:to-[#071f12] overflow-hidden relative">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-32 h-32 bg-white/10 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-1/3 right-1/5 w-24 h-24 bg-cyan-300/15 rounded-full blur-lg animate-float-delayed"></div>
      <div class="absolute bottom-1/4 left-1/3 w-20 h-20 bg-emerald-300/10 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-1/2 right-1/4 w-16 h-16 bg-teal-200/10 rounded-full blur-lg animate-float-delayed"></div>
      <div class="absolute bottom-1/3 right-1/3 w-28 h-28 bg-white/5 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-2/3 left-1/5 w-12 h-12 bg-cyan-200/10 rounded-full blur-md animate-float-delayed"></div>
    </div>
    <div class="bg-white dark:bg-[#1a1a1e] p-7 sm:p-9 rounded-2xl shadow-xl w-full max-w-md mx-4">
      <div class="text-center mb-7">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-800 dark:text-gray-200">备忘录</h1>
        <p class="text-sm text-gray-400 dark:text-gray-500 mt-1.5">欢迎回来，请登录</p>
      </div>

      <div v-if="loadingInit" class="flex justify-center py-6">
        <n-spin size="medium" />
      </div>

      <template v-else>
        <div class="flex flex-col gap-5">
          <n-input v-model:value="form.username" placeholder="用户名" size="large" clearable @keyup.enter="handleLogin" />
          <n-input v-model:value="form.password" type="password" placeholder="密码" size="large" show-password-on="click" @keyup.enter="handleLogin" />
          <div class="flex items-center justify-between -mt-1 px-0.5">
            <n-checkbox v-model:checked="rememberAccount">记住账号</n-checkbox>
            <n-checkbox v-model:checked="keepLogin">保持登录7天</n-checkbox>
          </div>
          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </n-button>
        </div>

        <div class="text-center mt-5 flex justify-center gap-4">
          <n-button text size="small" @click="$router.push('/forgot-password')">
            忘记密码？
          </n-button>
          <n-button text size="small" type="primary" @click="$router.push('/register')">
            没有账号？立即注册
          </n-button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NInput, NButton, NCheckbox, NSpin, useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { checkHasUsers } from '@/api/auth'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const loadingInit = ref(true)
const rememberAccount = ref(false)
const keepLogin = ref(true)

const form = reactive({
  username: '',
  password: '',
})

onMounted(async () => {
  // 首次使用（系统无用户）→ 跳转注册页创建管理员
  try {
    const res = await checkHasUsers()
    if (!res.has_users) {
      router.replace('/register')
      return
    }
  } catch {
    // 检查失败时按正常登录处理
  } finally {
    loadingInit.value = false
  }

  // 恢复记住的账号
  const savedUsername = localStorage.getItem('remembered_username')
  if (savedUsername) {
    form.username = savedUsername
    rememberAccount.value = true
  }
  keepLogin.value = localStorage.getItem('keep_login') !== 'false'
})

const handleLogin = async () => {
  if (!form.username) { message.warning('请输入用户名'); return }
  if (!form.password) { message.warning('请输入密码'); return }

  loading.value = true
  try {
    await authStore.login(form.username, form.password, keepLogin.value)

    if (rememberAccount.value) {
      localStorage.setItem('remembered_username', form.username)
    } else {
      localStorage.removeItem('remembered_username')
    }
    localStorage.setItem('keep_login', keepLogin.value ? 'true' : 'false')

    router.push('/')
  } catch (error: any) {
    message.error(error.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>
