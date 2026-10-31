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

      <n-alert v-else-if="suspectDataLoss" type="warning" :bordered="false">
        <p class="text-sm leading-relaxed">
          检测到数据异常：系统历史上已有账号，但当前数据库是空的。原数据文件可能暂时未被挂载（常见于
          NAS 重启时数据目录加载失败），重启 NAS 后再打开通常可恢复。
        </p>
        <p class="text-sm leading-relaxed mt-1.5">在数据找回之前，请先不要创建新账号。</p>
        <div class="flex gap-2 mt-3">
          <n-button size="small" type="primary" :loading="rechecking" @click="recheckUsers">重新检测</n-button>
          <n-button size="small" quaternary @click="$router.replace('/register')">仍要创建新账号</n-button>
        </div>
      </n-alert>

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
          <n-button
            v-if="fnosEnabled"
            secondary
            size="large"
            block
            :loading="loading"
            :disabled="loading"
            @click="handleFnOSLogin"
          >
            使用飞牛 NAS 登录
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
import { NInput, NButton, NCheckbox, NSpin, NAlert, useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { checkHasUsers, fnosLogin } from '@/api/auth'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const loadingInit = ref(true)
const suspectDataLoss = ref(false)
const rechecking = ref(false)
const rememberAccount = ref(false)
const keepLogin = ref(true)
const fnosEnabled = import.meta.env.VITE_FNOS_APP === 'true'

const form = reactive({
  username: '',
  password: '',
})

onMounted(async () => {
  // fnOS 入口由统一网关注入 NAS 身份：优先自动登录，未绑定时进入绑定流程。
  if (fnosEnabled) {
    loading.value = true
    try {
      await loginWithFnOS()
      return
    } catch (error: any) {
      message.error(error.response?.data?.error || '飞牛 NAS 自动登录失败，可重试或使用应用账号登录')
    } finally {
      loading.value = false
      loadingInit.value = false
    }
  }

  // 首次使用（系统无用户）→ 跳转注册页创建管理员；
  // 疑似数据丢失（曾有账号但现在是空库）→ 停留登录页告警，避免误导用户新建账号
  try {
    const res = await checkHasUsers()
    if (!res.has_users) {
      if (res.suspect_data_loss) {
        suspectDataLoss.value = true
        return
      }
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

const recheckUsers = async () => {
  rechecking.value = true
  try {
    const res = await checkHasUsers()
    if (res.has_users) {
      suspectDataLoss.value = false
      message.success('数据已恢复，请用原来的账号登录')
    } else {
      message.warning('数据库仍为空。建议先重启 NAS 再试；若仍未恢复，请联系开发者找回数据')
    }
  } catch {
    message.error('检测失败，请稍后重试')
  } finally {
    rechecking.value = false
  }
}

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

const handleFnOSLogin = async () => {
  loading.value = true
  try {
    await loginWithFnOS()
  } catch (error: any) {
    message.error(error.response?.data?.error || '飞牛 NAS 登录失败')
  } finally {
    loading.value = false
  }
}

const loginWithFnOS = async () => {
  const result = await fnosLogin()
  if (result.binding_required) {
    await router.replace({ path: '/register', query: { fnos: 'bind', fnos_username: result.fnos_username || '' } })
    return
  }
  authStore.token = result.token
  authStore.user = result.user
  authStore.isAuthenticated = true
  localStorage.setItem('token', result.token)
  await router.replace('/')
}
</script>
