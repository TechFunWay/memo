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
    <div class="bg-white dark:bg-[#1a1a1e] p-6 sm:p-8 rounded-2xl shadow-xl w-full max-w-md mx-4">
      <div class="text-center mb-6">
        <h1 class="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-200">找回密码</h1>
      </div>

      <!-- 步骤1：输入用户名 -->
      <div v-if="step === 1" class="flex flex-col gap-4">
        <p class="text-sm text-gray-500 dark:text-gray-400 text-center">输入用户名，验证安全问题</p>
        <n-input v-model:value="form.username" placeholder="用户名" size="large" clearable @keyup.enter="handleGetQuestion" />
        <n-button type="primary" size="large" block :loading="loading" @click="handleGetQuestion">
          下一步
        </n-button>
        <div class="text-center mt-1">
          <n-button text type="primary" @click="$router.push('/login')">返回登录</n-button>
        </div>
      </div>

      <!-- 步骤2：回答三个安全问题 + 设置新密码 -->
      <div v-else-if="step === 2" class="flex flex-col gap-4">
        <p class="text-sm text-gray-500 dark:text-gray-400 text-center">请回答以下三个安全问题</p>
        <div v-for="(q, i) in securityQuestions" :key="i" class="flex flex-col gap-1.5">
          <div class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ i + 1 }}. {{ q }}</div>
          <n-input v-model:value="form.answers[i]" :placeholder="`问题 ${i + 1} 的答案`" size="large" clearable />
        </div>
        <n-divider class="!my-1" />
        <n-input v-model:value="form.newPassword" type="password" placeholder="新密码" size="large" show-password-on="click" />
        <n-input v-model:value="form.confirmPassword" type="password" placeholder="确认新密码" size="large" show-password-on="click" @keyup.enter="handleResetPassword" />
        <n-button type="primary" size="large" block :loading="loading" @click="handleResetPassword">
          重置密码
        </n-button>
        <div class="text-center mt-1">
          <n-button text @click="step = 1">返回上一步</n-button>
        </div>
      </div>

      <!-- 完成 -->
      <div v-else class="text-center">
        <p class="text-green-600 dark:text-green-400 mb-3">密码重置成功！</p>
        <n-button type="primary" @click="$router.push('/login')">返回登录</n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { NInput, NButton, NDivider, useMessage } from 'naive-ui'
import { getSecurityQuestions, verifySecurityAnswers } from '@/api/auth'

const message = useMessage()

const loading = ref(false)
const step = ref(1)
const securityQuestions = ref<string[]>([])

const form = reactive({
  username: '',
  answers: ['', '', ''] as string[],
  newPassword: '',
  confirmPassword: '',
})

async function handleGetQuestion() {
  if (!form.username) { message.warning('请输入用户名'); return }

  loading.value = true
  try {
    const res = await getSecurityQuestions(form.username)
    securityQuestions.value = res.questions
    form.answers = ['', '', '']
    step.value = 2
  } catch (error: any) {
    message.error(error.response?.data?.error || '该用户未设置安全问题')
  } finally {
    loading.value = false
  }
}

async function handleResetPassword() {
  if (form.answers.some((a) => !a || !a.trim())) { message.warning('请回答全部三个安全问题'); return }
  if (!form.newPassword) { message.warning('请输入新密码'); return }
  if (form.newPassword !== form.confirmPassword) { message.warning('两次密码不一致'); return }
  if (form.newPassword.length < 6) { message.warning('密码至少6位'); return }

  loading.value = true
  try {
    await verifySecurityAnswers(form.username, form.answers, form.newPassword)
    step.value = 3
    message.success('密码重置成功')
  } catch (error: any) {
    message.error(error.response?.data?.error || '重置失败')
  } finally {
    loading.value = false
  }
}
</script>
