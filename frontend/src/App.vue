<template>
  <n-config-provider :theme="theme">
    <n-message-provider>
      <n-dialog-provider>
        <router-view />

        <!-- 赞赏横幅：仅管理员可见，未支持当前版本前常驻展示（关闭不记忆，
             刷新后再次出现，升级到新版本后也会再次出现） -->
        <div
          v-if="showSupportUI && support.bannerVisible && !support.donateSupported"
          class="fixed bottom-20 left-3 right-3 sm:left-auto sm:right-20 sm:w-[380px] z-40 flex items-center justify-between gap-3 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 px-4 py-2.5 text-sm text-white shadow-lg"
        >
          <span>☕ 如果备忘录对你有帮助，欢迎请作者喝杯咖啡——1 元也是心意，完全自愿。</span>
          <span class="flex items-center gap-2 shrink-0">
            <button class="underline underline-offset-2 hover:opacity-80" @click="support.open()">去赞赏</button>
            <button class="hover:opacity-80" title="关闭" aria-label="关闭赞赏横幅" @click="support.dismissBanner()">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
          </span>
        </div>

        <!-- 支持按钮：常驻入口，始终可以打开赞赏弹窗 -->
        <button
          v-if="showSupportUI"
          class="fixed bottom-5 right-5 z-40 flex h-11 w-11 items-center justify-center rounded-full bg-gradient-to-r from-amber-500 to-orange-500 text-lg text-white shadow-lg transition-transform hover:scale-105"
          aria-label="支持作者"
          title="支持作者"
          @click="support.open()"
        >
          ☕
        </button>

        <SupportModal />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, watch, onMounted } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import { useRoute } from 'vue-router'
import { darkTheme } from 'naive-ui'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { useSupportStore } from '@/stores/support'
import SupportModal from '@/components/SupportModal.vue'

const route = useRoute()
const themeStore = useThemeStore()
const authStore = useAuthStore()
const support = useSupportStore()
const theme = computed(() => (themeStore.darkMode ? darkTheme : null))

watch(
  () => themeStore.darkMode,
  (dark) => {
    document.documentElement.classList.toggle('dark', dark)
  },
  { immediate: true },
)

// 赞赏提示仅对已登录的管理员显示，登录/注册/找回密码页不出现
const AUTH_PATHS = ['/login', '/register', '/forgot-password']
const isAuthPage = computed(() => AUTH_PATHS.includes(route.path))
const showSupportUI = computed(
  () => isAuthPage.value === false && authStore.isAuthenticated && authStore.user?.is_admin === true,
)

// 登录后（含刷新时已持有会话）初始化一次：拉取版本并调度弹窗/横幅
let supportInited = false
watch(showSupportUI, (visible) => {
  if (visible && !supportInited) {
    supportInited = true
    void support.init(true)
  }
})

onMounted(() => {
  authStore.checkAuth()
})
</script>
