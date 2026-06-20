<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-400 via-teal-500 to-cyan-600 dark:from-[#0d3320] dark:via-[#0a2e1c] dark:to-[#071f12] overflow-hidden relative py-8">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-32 h-32 bg-white/10 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-1/3 right-1/5 w-24 h-24 bg-cyan-300/15 rounded-full blur-lg animate-float-delayed"></div>
      <div class="absolute bottom-1/4 left-1/3 w-20 h-20 bg-emerald-300/10 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-1/2 right-1/4 w-16 h-16 bg-teal-200/10 rounded-full blur-lg animate-float-delayed"></div>
      <div class="absolute bottom-1/3 right-1/3 w-28 h-28 bg-white/5 rounded-full blur-xl animate-float"></div>
      <div class="absolute top-2/3 left-1/5 w-12 h-12 bg-cyan-200/10 rounded-full blur-md animate-float-delayed"></div>
    </div>
    <div class="bg-white dark:bg-[#1a1a1e] p-7 sm:p-9 rounded-2xl shadow-xl w-full max-w-md mx-4">
      <div class="text-center mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-800 dark:text-gray-200 mb-1">{{ isFirstUse ? '创建管理员' : '创建账号' }}</h1>
        <p v-if="isFirstUse" class="text-sm text-emerald-600 dark:text-emerald-400">首次使用，第一个注册的用户将成为管理员</p>
        <p v-else class="text-sm text-gray-400 dark:text-gray-500">注册你的备忘录账号</p>
      </div>

      <div class="flex flex-col gap-4">
        <n-input v-model:value="form.username" placeholder="用户名" size="large" clearable />
        <n-input v-model:value="form.password" type="password" placeholder="密码（至少6位）" size="large" show-password-on="click" />
        <n-input v-model:value="form.confirmPassword" type="password" placeholder="确认密码" size="large" show-password-on="click" />
        <n-input v-model:value="form.nickname" placeholder="昵称（可选，默认与用户名相同）" size="large" clearable />
        <n-input v-model:value="form.email" placeholder="邮箱（可选）" size="large" clearable />

        <n-divider class="!my-1" />
        <p class="text-xs text-gray-500 dark:text-gray-400 -mt-1">
          安全问题（可选，用于找回密码）。如需设置请完整填写三个，否则请全部留空。
        </p>
        <SecurityQuestionsForm v-model:questions="form.securityQuestions" v-model:answers="form.securityAnswers" />

        <n-button
          type="primary"
          size="large"
          block
          class="!mt-2"
          :loading="loading"
          :disabled="!allowRegister || configLoading"
          @click="handleRegister"
        >
          {{ !allowRegister ? '系统暂不允许注册' : (isFirstUse ? '创建管理员' : '注册') }}
        </n-button>
      </div>

      <div v-if="!isFirstUse" class="text-center mt-4">
        <n-button text type="primary" @click="$router.push('/login')">
          已有账号？立即登录
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NInput, NButton, NDivider, useMessage } from 'naive-ui'
import { register, checkHasUsers } from '@/api/auth'
import { getPublicConfig } from '@/api/config'
import { useAuthStore } from '@/stores/auth'
import SecurityQuestionsForm from '@/components/SecurityQuestionsForm.vue'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const allowRegister = ref(true)
const configLoading = ref(true)
const isFirstUse = ref(false)

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  email: '',
  securityQuestions: ['', '', ''] as string[],
  securityAnswers: ['', '', ''] as string[],
})

// 校验三个安全问题是否完整（全部填写且问题不重复）
function validateSecurity(questions: string[], answers: string[]): { complete: boolean; touched: boolean; error?: string } {
  const touched = questions.some((q) => q) || answers.some((a) => a)
  const allFilled = questions.every((q) => q && q.trim()) && answers.every((a) => a && a.trim())
  if (!allFilled) {
    return { complete: false, touched, error: touched ? '请完整填写三个安全问题及答案，或全部留空' : undefined }
  }
  const qs = questions.map((q) => q.trim())
  if (new Set(qs).size !== qs.length) {
    return { complete: false, touched, error: '三个安全问题不能重复' }
  }
  return { complete: true, touched }
}

const handleRegister = async () => {
  if (!allowRegister.value) return
  if (!form.username) { message.warning('请输入用户名'); return }
  if (!form.password) { message.warning('请输入密码'); return }
  if (form.password.length < 6) { message.warning('密码至少6位'); return }
  if (form.password !== form.confirmPassword) { message.warning('两次输入的密码不一致'); return }

  const sec = validateSecurity(form.securityQuestions, form.securityAnswers)
  if (sec.touched && !sec.complete) { message.warning(sec.error!); return }

  const nickname = form.nickname.trim() || form.username.trim()

  loading.value = true
  try {
    const res = await register(
      form.username,
      form.password,
      nickname,
      form.email,
      sec.complete ? form.securityQuestions : [],
      sec.complete ? form.securityAnswers : [],
    )

    if (isFirstUse.value) {
      // 首个用户即管理员：直接登录并进入首页（未设置安全问题时由首页弹窗提示）
      authStore.token = res.token
      authStore.user = res.user
      authStore.isAuthenticated = true
      localStorage.setItem('token', res.token)
      message.success('管理员创建成功')
      router.push('/')
    } else {
      // 普通注册：去登录页（登录后到首页会提示设置安全问题）
      message.success('注册成功，请登录')
      router.push('/login')
    }
  } catch (error: any) {
    message.error(error.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    // 首次使用（系统无用户）：允许注册以创建管理员，不受"允许注册"开关限制
    const usersRes = await checkHasUsers()
    isFirstUse.value = !usersRes.has_users
  } catch {
    isFirstUse.value = false
  }

  if (isFirstUse.value) {
    allowRegister.value = true
    configLoading.value = false
    return
  }

  try {
    const res = await getPublicConfig()
    allowRegister.value = res.configs?.allow_register === '1'
  } catch {
    allowRegister.value = true
  } finally {
    configLoading.value = false
  }
})
</script>
