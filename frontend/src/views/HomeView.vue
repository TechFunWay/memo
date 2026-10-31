<template>
  <div class="h-screen flex flex-col w-screen max-w-[100vw] overflow-hidden bg-gray-50 dark:bg-[#101014]">
    <!-- 顶部导航 -->
    <header class="bg-white dark:bg-[#1a1a1e] shadow-sm border-b dark:border-[#2d2d33] flex-shrink-0">
      <div class="px-3 py-2 sm:px-4 sm:py-3 flex items-center justify-between">
        <div class="flex items-center space-x-3">
            <n-button quaternary circle @click="toggleSidebar">
              <template #icon>
                <n-icon :component="MenuOutline" />
              </template>
            </n-button>
            <div class="flex items-center gap-1.5">
              <h1 class="text-lg sm:text-xl font-bold text-gray-800 dark:text-gray-200">备忘录</h1>
              <span class="text-[10px] text-gray-400 dark:text-gray-500 leading-none">{{ appVersion?.replace('v', '') }}</span>
            </div>
          </div>
        <div class="flex items-center space-x-2">
          <n-button quaternary circle @click="toggleDarkMode">
            <template #icon>
              <n-icon :component="darkMode ? SunnyOutline : MoonOutline" />
            </template>
          </n-button>
          <n-dropdown :options="userMenuOptions" trigger="click">
            <n-button quaternary circle>
              <template #icon>
                <n-icon :component="PersonOutline" />
              </template>
            </n-button>
          </n-dropdown>
        </div>
      </div>
    </header>

    <div class="flex flex-1 overflow-hidden relative">
      <!-- 移动端遮罩 -->
      <div
        v-if="sidebarOpen && isMobile"
        class="fixed inset-0 bg-black/30 z-30"
        @click="sidebarOpen = false"
      />

      <!-- 侧边栏 -->
      <aside
        :style="!isMobile && sidebarOpen ? { width: `${sidebarWidth}px` } : undefined"
        :class="[
          'bg-gray-50 dark:bg-[#151518] border-r dark:border-[#2d2d33] flex-shrink-0 overflow-hidden',
          !isSidebarResizing && 'transition-all duration-300',
          sidebarOpen
            ? (isMobile ? 'w-72 fixed left-0 top-0 bottom-0 z-40' : 'sidebar-expanded')
            : (isMobile ? 'w-0' : 'w-14'),
        ]"
      >
        <div class="flex flex-col h-full overflow-hidden" style="min-width: 0;">
          <!-- 可滚动内容区 -->
          <div class="flex-1 overflow-y-auto overflow-x-hidden p-2">
            <!-- 新建按钮 -->
            <n-tooltip v-if="!sidebarOpen" placement="right" trigger="hover">
              <template #trigger>
                <n-button type="primary" block class="h-10" @click="goToEditor('new')">
                  <template #icon>
                    <n-icon size="20" :component="AddCircleOutline" />
                  </template>
                </n-button>
              </template>
              新建备忘录
            </n-tooltip>
            <n-button v-else type="primary" block @click="goToEditor('new')">
              <template #icon>
                <n-icon :component="AddCircleOutline" />
              </template>
              新建备忘录
            </n-button>

            <!-- 展开模式：用 NMenu -->
            <n-menu
              v-if="sidebarOpen || isMobile"
              class="mt-2"
              :options="menuOptions"
              :indent="18"
              :value="activeMenuKey"
              @update:value="handleMenuSelect"
            />
            <!-- 收缩模式：用自定义图标按钮 -->
            <div v-else class="mt-2 flex flex-col items-center space-y-1 px-1">
              <n-tooltip v-for="opt in menuOptions" :key="opt.key" placement="right" trigger="hover">
                <template #trigger>
                  <div
                    class="w-full h-10 flex items-center justify-center rounded-lg cursor-pointer transition-colors"
                    :class="activeMenuKey === opt.key
                      ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                      : 'hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'"
                    @click="handleMenuSelect(opt.key)"
                  >
                    <component :is="opt.icon" />
                  </div>
                </template>
                {{ opt.label }}
              </n-tooltip>
            </div>

            <!-- 文件夹区域 -->
            <n-divider v-if="sidebarOpen || isMobile" style="margin: 12px 0 4px" title-placement="left">文件夹</n-divider>
            <div v-else class="border-t dark:border-[#2d2d33] my-2" />

            <!-- 文件夹操作栏 -->
            <div v-if="sidebarOpen || isMobile" class="flex items-center justify-between mb-1 px-2">
              <span class="text-xs text-gray-500 dark:text-gray-400 truncate">文件夹</span>
              <n-tooltip placement="top" trigger="hover">
                <template #trigger>
                  <n-button text size="tiny" type="primary" @click="showCreateFolderModal = true">
                    <template #icon>
                      <n-icon :component="AddCircleOutline" />
                    </template>
                    新建
                  </n-button>
                </template>
                新建文件夹
              </n-tooltip>
            </div>
            <n-tooltip v-else placement="right" trigger="hover">
              <template #trigger>
                <div
                  class="w-10 h-10 flex items-center justify-center rounded-lg cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400 transition-colors mx-auto mb-1"
                  @click="showCreateFolderModal = true"
                >
                  <n-icon size="20" :component="AddCircleOutline" />
                </div>
              </template>
              新建文件夹
            </n-tooltip>

            <!-- 展开模式：文件夹自定义列表 -->
            <div v-if="sidebarOpen || isMobile" class="space-y-0.5 px-1">
              <div
                v-for="opt in folderMenuOptions"
                :key="opt.key"
                class="folder-item group flex items-center rounded-lg cursor-pointer transition-colors"
                :class="activeMenuKey === opt.key
                  ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                  : 'hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'"
                @click="handleMenuSelect(opt.key)"
              >
                <div class="flex items-center flex-1 min-w-0 px-2 py-1.5">
                  <n-icon size="16" class="flex-shrink-0 mr-2"><FolderOutline /></n-icon>
                  <span class="text-sm truncate">{{ opt.label }}</span>
                </div>
                <div class="flex items-center flex-shrink-0 mr-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <n-button text size="tiny" @click.stop="openRenameFolder(opt.folderId, opt.label)">
                    <template #icon><n-icon size="14"><PencilOutline /></n-icon></template>
                  </n-button>
                  <n-button text size="tiny" type="error" @click.stop="confirmDeleteFolder(opt.folderId, opt.label)">
                    <template #icon><n-icon size="14"><TrashOutline /></n-icon></template>
                  </n-button>
                </div>
              </div>
            </div>
            <!-- 收缩模式：用自定义图标按钮 -->
            <div v-else class="flex flex-col items-center space-y-1 px-1">
              <n-tooltip v-for="opt in folderMenuOptions" :key="opt.key" placement="right" trigger="hover">
                <template #trigger>
                  <div
                    class="w-full h-10 flex items-center justify-center rounded-lg cursor-pointer transition-colors"
                    :class="activeMenuKey === opt.key
                      ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                      : 'hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'"
                    @click="handleMenuSelect(opt.key)"
                  >
                    <component :is="opt.icon" />
                  </div>
                </template>
                {{ opt.label }}
              </n-tooltip>
            </div>
          </div>

          <!-- 底部折叠按钮（桌面端） -->
          <div v-if="!isMobile" class="flex-shrink-0 border-t dark:border-[#2d2d33] p-2">
            <n-tooltip :content="sidebarOpen ? '收起侧边栏' : '展开侧边栏'" placement="right" trigger="hover">
              <template #trigger>
                <n-button
                  quaternary
                  block
                  :aria-label="sidebarOpen ? '收起侧边栏' : '展开侧边栏'"
                  @click="toggleSidebar"
                  class="flex items-center justify-center"
                >
                  <template #icon>
                    <n-icon :component="sidebarOpen ? ArrowBackOutline : ArrowForwardOutline" />
                  </template>
                  <span v-if="sidebarOpen" class="ml-2 text-xs">收起</span>
                </n-button>
              </template>
            </n-tooltip>
          </div>
        </div>
      </aside>

      <!-- 桌面端侧边栏宽度调节条 -->
      <div
        v-if="sidebarOpen && !isMobile"
        class="sidebar-resizer"
        :class="{ 'is-resizing': isSidebarResizing }"
        role="separator"
        aria-orientation="vertical"
        aria-label="调整侧边栏宽度"
        :aria-valuemin="SIDEBAR_MIN_WIDTH"
        :aria-valuemax="SIDEBAR_MAX_WIDTH"
        :aria-valuenow="Math.round(sidebarWidth)"
        :aria-valuetext="`${Math.round(sidebarWidth)} 像素`"
        title="拖动调整侧边栏宽度，双击恢复默认宽度"
        tabindex="0"
        @pointerdown="startSidebarResize"
        @pointercancel="stopSidebarResize"
        @dblclick="resetSidebarWidth"
        @keydown="handleSidebarResizeKeydown"
      >
        <span class="sidebar-resizer-handle" aria-hidden="true" />
      </div>

      <!-- 主内容区 -->
      <main class="flex-1 overflow-auto bg-gray-100 dark:bg-[#101014] p-3 sm:p-4 w-full">
        <div class="w-full">
          <!-- 移动端顶部操作栏 -->
          <div v-if="isMobile" class="mb-3 flex items-center justify-between">
            <n-button text type="primary" @click="goToEditor('new')" class="flex-shrink-0">
              <template #icon>
                <n-icon :component="AddCircleOutline" />
              </template>
              新建
            </n-button>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              共 {{ filteredMemos.length }} 条
            </span>
          </div>

          <!-- 搜索栏 -->
          <div class="mb-3 sm:mb-4">
            <n-input
              v-model:value="searchQuery"
              placeholder="搜索备忘录..."
              clearable
              size="large"
              @update:value="handleSearch"
            >
              <template #prefix>
                <n-icon :component="SearchOutline" />
              </template>
            </n-input>
          </div>

          <!-- 备忘录列表 -->
          <n-spin :show="memoStore.loading">
            <div v-if="filteredMemos.length === 0" class="text-center py-20 text-gray-400 dark:text-gray-500">
              <n-icon size="48" :component="SearchOutline" class="mx-auto mb-4" />
              <p>暂无备忘录</p>
              <n-button text type="primary" class="mt-2" @click="goToEditor('new')">
                创建第一条备忘录
              </n-button>
            </div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
              <div
                v-for="memo in pagedMemos"
                :key="memo.id"
                class="bg-white dark:bg-[#1a1a1e] rounded-lg shadow-sm hover:shadow-md transition-shadow p-3 sm:p-4 cursor-pointer border-l-4"
                :class="memo.is_pinned ? 'border-l-yellow-400' : 'border-l-blue-400'"
                @click="goToEditor(memo.id)"
              >
                <div class="flex items-center gap-2 mb-2">
                  <n-icon v-if="memo.is_pinned" :component="Pin" color="#fbbf24" class="flex-shrink-0" />
                  <h3 class="font-semibold text-gray-800 dark:text-gray-200 truncate text-sm sm:text-base flex-1 min-w-0">{{ memo.title || '无标题' }}</h3>
                  <span v-if="getFolderName(memo.folder_id)" class="text-[10px] text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-[#252529] px-1.5 py-0.5 rounded flex-shrink-0">{{ getFolderName(memo.folder_id) }}</span>
                </div>
                <div class="text-gray-600 dark:text-gray-400 text-xs sm:text-sm line-clamp-3">{{ stripHtml(memo.content || '无内容') }}</div>
                <div class="mt-2 sm:mt-3 flex items-end justify-between text-xs text-gray-500 dark:text-gray-500">
                  <div class="flex flex-col min-w-0 flex-1 mr-2 leading-tight gap-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                    <span class="truncate">创建 {{ formatDateTime(memo.created_at) }}</span>
                    <span class="truncate">更新 {{ formatDateTime(memo.updated_at) }}</span>
                  </div>
                  <div class="flex space-x-2 flex-shrink-0">
                    <n-icon
                      :component="Pin"
                      :color="memo.is_pinned ? '#fbbf24' : undefined"
                      @click.stop="togglePin(memo)"
                    />
                    <n-icon
                      :component="StarOutline"
                      :color="memo.is_favorite ? '#f59e0b' : undefined"
                      @click.stop="toggleFavorite(memo)"
                    />
                    <n-icon :component="TrashOutline" color="#ef4444" @click.stop="handleDelete(memo.id)" />
                  </div>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <div v-if="filteredMemos.length > pageSize" class="mt-4 sm:mt-6 flex justify-center">
              <n-pagination
                v-model:page="currentPage"
                :page-size="pageSize"
                :item-count="filteredMemos.length"
                :page-slot="isMobile ? 5 : 9"
              />
            </div>
          </n-spin>
        </div>

        <!-- 移动端悬浮新建按钮 -->
        <div v-if="isMobile" class="fixed bottom-6 right-6 z-20 sm:hidden">
          <n-button type="primary" circle size="large" @click="goToEditor('new')" class="shadow-lg">
            <template #icon>
              <n-icon size="24" :component="AddCircleOutline" />
            </template>
          </n-button>
        </div>
      </main>
    </div>

    <!-- 修改密码弹窗 -->
    <n-modal v-model:show="showChangePasswordModal" preset="card" title="修改密码" class="w-[90vw] max-w-[400px]">
      <n-form :model="changePasswordForm">
        <n-form-item label="旧密码">
          <n-input v-model:value="changePasswordForm.oldPassword" type="password" placeholder="请输入旧密码" />
        </n-form-item>
        <n-form-item label="新密码">
          <n-input v-model:value="changePasswordForm.newPassword" type="password" placeholder="请输入新密码" />
        </n-form-item>
        <n-form-item label="确认新密码">
          <n-input v-model:value="changePasswordForm.confirmPassword" type="password" placeholder="请再次输入新密码" @keyup.enter="submitChangePassword" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showChangePasswordModal = false">取消</n-button>
          <n-button type="primary" :loading="changingPassword" @click="submitChangePassword">确认修改</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 新建文件夹弹窗 -->
    <n-modal v-model:show="showCreateFolderModal" preset="card" title="新建文件夹" class="w-[90vw] max-w-[400px]">
      <n-form :model="folderForm">
        <n-form-item label="文件夹名称">
          <n-input v-model:value="folderForm.name" placeholder="请输入文件夹名称" @keyup.enter="saveFolder" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateFolderModal = false">取消</n-button>
          <n-button type="primary" :loading="folderSaving" @click="saveFolder">创建</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 重命名文件夹弹窗 -->
    <n-modal v-model:show="showRenameFolderModal" preset="card" title="重命名文件夹" class="w-[90vw] max-w-[400px]">
      <n-form :model="renameFolderForm">
        <n-form-item label="文件夹名称">
          <n-input v-model:value="renameFolderForm.name" placeholder="请输入新名称" @keyup.enter="submitRenameFolder" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showRenameFolderModal = false">取消</n-button>
          <n-button type="primary" :loading="folderRenaming" @click="submitRenameFolder">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 安全设置弹窗 -->
    <n-modal v-model:show="showSecurityModal" preset="card" title="安全设置" class="w-[90vw] max-w-[440px]">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">设置三个安全问题用于找回密码（请完整填写）</p>
      <SecurityQuestionsForm v-model:questions="securityForm.questions" v-model:answers="securityForm.answers" @enter="submitSecurityQuestion" />
      <template #footer>
        <n-space justify="end">
          <n-button @click="showSecurityModal = false">取消</n-button>
          <n-button type="primary" :loading="securityLoading" @click="submitSecurityQuestion">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, h, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NIcon, NInput, NMenu, NDropdown, NModal, NForm, NFormItem, NSpace, NSelect, NSpin, NDivider, NPagination, useMessage, useDialog } from 'naive-ui'
import {
  MenuOutline,
  SunnyOutline,
  MoonOutline,
  PersonOutline,
  AddCircleOutline,
  SearchOutline,
  StarOutline,
  TrashOutline,
  FolderOutline,
  PencilOutline,
  Pin,
} from '@vicons/ionicons5'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { useMemoStore } from '@/stores/memo'
import { useFolderStore } from '@/stores/folder'
import type { Memo } from '@/api/memo'
import { getVersion } from '@/api/system'
import { setSecurityQuestions, checkSecurityQuestion } from '@/api/auth'
import SecurityQuestionsForm from '@/components/SecurityQuestionsForm.vue'

const appVersion = ref('v1.1.1')

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const themeStore = useThemeStore()
const authStore = useAuthStore()
const memoStore = useMemoStore()
const folderStore = useFolderStore()

const sidebarOpen = ref(false)
const searchQuery = ref('')
const darkMode = computed(() => themeStore.darkMode)
const activeMenuKey = ref<string>('all')
const isMobile = ref(false)
const SIDEBAR_WIDTH_STORAGE_KEY = 'memo-sidebar-width'
const SIDEBAR_DEFAULT_WIDTH = 240
const SIDEBAR_MIN_WIDTH = 200
const SIDEBAR_MAX_WIDTH = 360
const sidebarWidth = ref(SIDEBAR_DEFAULT_WIDTH)
const isSidebarResizing = ref(false)

let sidebarResizeStartX = 0
let sidebarResizeStartWidth = SIDEBAR_DEFAULT_WIDTH
let previousBodyUserSelect = ''
let previousBodyCursor = ''

const currentPage = ref(1)
const pageSize = 12

const showCreateFolderModal = ref(false)
const folderSaving = ref(false)
const folderForm = ref({ name: '' })

const showRenameFolderModal = ref(false)
const folderRenaming = ref(false)
const renameFolderForm = ref({ id: 0, name: '' })

const showChangePasswordModal = ref(false)
const changingPassword = ref(false)
const changePasswordForm = ref({ oldPassword: '', newPassword: '', confirmPassword: '' })

const showSecurityModal = ref(false)
const securityLoading = ref(false)
const securityForm = ref<{ questions: string[]; answers: string[] }>({ questions: ['', '', ''], answers: ['', '', ''] })

function checkMobile() {
  isMobile.value = window.innerWidth < 768
}

function clampSidebarWidth(width: number) {
  return Math.min(SIDEBAR_MAX_WIDTH, Math.max(SIDEBAR_MIN_WIDTH, width))
}

function loadSidebarWidth() {
  try {
    const storedWidth = Number(localStorage.getItem(SIDEBAR_WIDTH_STORAGE_KEY))
    if (Number.isFinite(storedWidth)) {
      sidebarWidth.value = clampSidebarWidth(storedWidth)
    }
  } catch {
    // 隐私模式下 localStorage 不可用时使用默认宽度
  }
}

function persistSidebarWidth() {
  try {
    localStorage.setItem(SIDEBAR_WIDTH_STORAGE_KEY, String(Math.round(sidebarWidth.value)))
  } catch {
    // 不影响侧边栏当前会话的缩放
  }
}

checkMobile()
loadSidebarWidth()
window.addEventListener('resize', checkMobile)
onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  stopSidebarResize()
})

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

function startSidebarResize(event: PointerEvent) {
  if (isMobile.value || !sidebarOpen.value) return

  event.preventDefault()
  sidebarResizeStartX = event.clientX
  sidebarResizeStartWidth = sidebarWidth.value
  isSidebarResizing.value = true

  previousBodyUserSelect = document.body.style.userSelect
  previousBodyCursor = document.body.style.cursor
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'col-resize'

  window.addEventListener('pointermove', handleSidebarResize)
  window.addEventListener('pointerup', stopSidebarResize)
  window.addEventListener('pointercancel', stopSidebarResize)
}

function handleSidebarResize(event: PointerEvent) {
  if (!isSidebarResizing.value) return
  sidebarWidth.value = clampSidebarWidth(sidebarResizeStartWidth + event.clientX - sidebarResizeStartX)
}

function stopSidebarResize() {
  if (!isSidebarResizing.value) return

  isSidebarResizing.value = false
  persistSidebarWidth()
  document.body.style.userSelect = previousBodyUserSelect
  document.body.style.cursor = previousBodyCursor

  window.removeEventListener('pointermove', handleSidebarResize)
  window.removeEventListener('pointerup', stopSidebarResize)
  window.removeEventListener('pointercancel', stopSidebarResize)
}

function resetSidebarWidth() {
  sidebarWidth.value = SIDEBAR_DEFAULT_WIDTH
  persistSidebarWidth()
}

function handleSidebarResizeKeydown(event: KeyboardEvent) {
  let nextWidth: number | null = null

  if (event.key === 'ArrowLeft') {
    nextWidth = sidebarWidth.value - 16
  } else if (event.key === 'ArrowRight') {
    nextWidth = sidebarWidth.value + 16
  } else if (event.key === 'Home') {
    nextWidth = SIDEBAR_MIN_WIDTH
  } else if (event.key === 'End') {
    nextWidth = SIDEBAR_MAX_WIDTH
  }

  if (nextWidth !== null) {
    event.preventDefault()
    sidebarWidth.value = clampSidebarWidth(nextWidth)
    persistSidebarWidth()
  }
}

const toggleDarkMode = () => {
  themeStore.toggleDarkMode()
}

function goToEditor(id: number | string) {
  router.push(`/editor/${id}`)
}

const menuOptions = computed(() => [
  {
    label: '全部备忘录',
    key: 'all',
    icon: () => h(NIcon, null, { default: () => h(AddCircleOutline) }),
  },
  {
    label: '收藏夹',
    key: 'favorite',
    icon: () => h(NIcon, { color: '#f59e0b' }, { default: () => h(StarOutline) }),
  },
  {
    label: '置顶',
    key: 'pinned',
    icon: () => h(NIcon, { color: '#fbbf24' }, { default: () => h(AddCircleOutline) }),
  },
])

const folderMenuOptions = computed(() =>
  folderStore.folders.map((f) => ({
    label: f.name,
    key: `folder_${f.id}`,
    folderId: f.id,
    icon: () => h(NIcon, null, { default: () => h(FolderOutline) }),
  }))
)

const filteredMemos = computed(() => {
  let result = memoStore.memos
  if (activeMenuKey.value === 'favorite') {
    result = result.filter((m) => m.is_favorite)
  } else if (activeMenuKey.value === 'pinned') {
    result = result.filter((m) => m.is_pinned)
  } else if (activeMenuKey.value?.startsWith('folder_')) {
    const folderId = parseInt(activeMenuKey.value.replace('folder_', ''), 10)
    result = result.filter((m) => m.folder_id === folderId)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    result = result.filter((m) => {
      const title = (m.title || '').toLowerCase()
      const text = stripHtml(m.content || '').toLowerCase()
      return title.includes(q) || text.includes(q)
    })
  }
  return result
})

// 分页：仅展示当前页
const pagedMemos = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredMemos.value.slice(start, start + pageSize)
})

// 切换筛选/搜索时回到第一页
watch([searchQuery, activeMenuKey], () => {
  currentPage.value = 1
})

// 数据变化后若当前页超出范围则收敛
watch(() => filteredMemos.value.length, (len) => {
  const maxPage = Math.max(1, Math.ceil(len / pageSize))
  if (currentPage.value > maxPage) currentPage.value = maxPage
})

function handleMenuSelect(key: string) {
  activeMenuKey.value = key
  if (isMobile.value) {
    sidebarOpen.value = false
  }
}

function handleSearch() {
  // filteredMemos is computed, no action needed
}

async function toggleFavorite(memo: Memo) {
  await memoStore.updateMemo(memo.id, { is_favorite: !memo.is_favorite })
  message.success(memo.is_favorite ? '已取消收藏' : '已收藏')
}

async function togglePin(memo: Memo) {
  await memoStore.updateMemo(memo.id, { is_pinned: !memo.is_pinned })
  message.success(memo.is_pinned ? '已取消置顶' : '已置顶')
}

async function handleDelete(id: number) {
  dialog.warning({
    title: '确认删除',
    content: '删除后无法恢复，确定要删除这条备忘录吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await memoStore.deleteMemo(id)
      message.success('已删除')
    },
  })
}

async function saveFolder() {
  if (!folderForm.value.name) {
    message.warning('请输入文件夹名称')
    return
  }
  folderSaving.value = true
  try {
    await folderStore.createFolder({ name: folderForm.value.name })
    message.success('创建成功')
    showCreateFolderModal.value = false
    folderForm.value = { name: '' }
  } finally {
    folderSaving.value = false
  }
}

function openRenameFolder(id: number, name: string) {
  renameFolderForm.value = { id, name }
  showRenameFolderModal.value = true
}

async function submitRenameFolder() {
  if (!renameFolderForm.value.name) {
    message.warning('请输入文件夹名称')
    return
  }
  folderRenaming.value = true
  try {
    await folderStore.updateFolder(renameFolderForm.value.id, { name: renameFolderForm.value.name })
    message.success('重命名成功')
    showRenameFolderModal.value = false
  } finally {
    folderRenaming.value = false
  }
}

function confirmDeleteFolder(id: number, name: string) {
  const memoCount = memoStore.memos.filter((m) => m.folder_id === id).length
  dialog.warning({
    title: '确认删除',
    content: `文件夹「${name}」内有 ${memoCount} 条备忘录，删除后这些备忘录将变为未分类。确定删除吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await folderStore.deleteFolder(id)
        await memoStore.fetchMemos()
        message.success('已删除')
        if (activeMenuKey.value === `folder_${id}`) {
          activeMenuKey.value = 'all'
        }
      } catch {
        message.error('删除失败')
      }
    },
  })
}

function getFolderName(folderId: number | null): string {
  if (folderId == null) return ''
  const folder = folderStore.folders.find((f) => f.id === folderId)
  return folder?.name || ''
}

function stripHtml(html: string): string {
  const div = document.createElement('div')
  div.innerHTML = html
  return div.textContent || div.innerText || ''
}

async function submitSecurityQuestion() {
  const { questions, answers } = securityForm.value
  if (questions.some((q) => !q || !q.trim()) || answers.some((a) => !a || !a.trim())) {
    message.warning('请完整填写三个安全问题及答案')
    return
  }
  const qs = questions.map((q) => q.trim())
  if (new Set(qs).size !== qs.length) {
    message.warning('三个安全问题不能重复')
    return
  }

  securityLoading.value = true
  try {
    await setSecurityQuestions(questions, answers)
    message.success('安全设置已保存')
    showSecurityModal.value = false
    securityForm.value = { questions: ['', '', ''], answers: ['', '', ''] }
  } catch (e: any) {
    message.error(e.response?.data?.error || '保存失败')
  } finally {
    securityLoading.value = false
  }
}


async function submitChangePassword() {
  const { oldPassword, newPassword, confirmPassword } = changePasswordForm.value
  if (!oldPassword) {
    message.warning('请输入旧密码')
    return
  }
  if (!newPassword) {
    message.warning('请输入新密码')
    return
  }
  if (newPassword !== confirmPassword) {
    message.warning('两次输入的密码不一致')
    return
  }
  changingPassword.value = true
  try {
    await authStore.changePassword(oldPassword, newPassword)
    message.success('密码修改成功')
    showChangePasswordModal.value = false
    changePasswordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (e: any) {
    message.error(e.response?.data?.error || '修改失败')
  } finally {
    changingPassword.value = false
  }
}

const userMenuOptions = computed(() => {
  const opts: any[] = []
  if (authStore.user?.is_admin) {
    opts.push({
      label: '管理员',
      key: 'admin',
      props: {
        onClick: () => {
          router.push('/admin')
        },
      },
    })
  }
  opts.push({
    label: '修改密码',
    key: 'changePassword',
    props: {
      onClick: () => {
        showChangePasswordModal.value = true
      },
    },
  })
  opts.push({
    label: '安全设置',
    key: 'security',
    props: {
      onClick: () => {
        showSecurityModal.value = true
      },
    },
  })
  opts.push({
    label: '退出登录',
    key: 'logout',
    props: {
      onClick: () => {
        authStore.logout()
      },
    },
  })
  return opts
})

// 直接展示后端存储的北京时间，不做时区换算（避免二次偏移）
function formatDateTime(s?: string) {
  if (!s) return ''
  const m = s.replace('T', ' ').match(/(\d{4})-(\d{2})-(\d{2})[ ]?(\d{2}):(\d{2})/)
  if (!m) return s
  return `${m[1]}-${m[2]}-${m[3]} ${m[4]}:${m[5]}`
}

onMounted(async () => {
  await Promise.all([memoStore.fetchMemos(), folderStore.fetchFolders()])
  try {
    const res = await getVersion()
    appVersion.value = res.version || 'v1.1.1'
  } catch {
    // ignore version fetch errors
  }

  // 检查是否已设置安全问题，未设置则提示
  try {
    const res = await checkSecurityQuestion()
    if (!res.has_security_question) {
      showSecurityModal.value = true
    }
  } catch {
    // ignore check errors
  }
})
</script>

<style scoped>
.sidebar-resizer {
  position: relative;
  z-index: 5;
  display: flex;
  flex: 0 0 8px;
  width: 8px;
  height: 100%;
  align-items: center;
  justify-content: center;
  cursor: col-resize;
  touch-action: none;
  outline: none;
}

.sidebar-resizer::before {
  content: '';
  width: 1px;
  height: 100%;
  background: transparent;
  transition: background-color 150ms ease;
}

.sidebar-resizer-handle {
  position: absolute;
  width: 4px;
  height: 36px;
  border-radius: 9999px;
  background: #cbd5e1;
  opacity: 0;
  transition: opacity 150ms ease, background-color 150ms ease;
}

.sidebar-resizer:hover::before,
.sidebar-resizer:focus-visible::before,
.sidebar-resizer.is-resizing::before {
  background: #94a3b8;
}

.sidebar-resizer:hover .sidebar-resizer-handle,
.sidebar-resizer:focus-visible .sidebar-resizer-handle,
.sidebar-resizer.is-resizing .sidebar-resizer-handle {
  opacity: 1;
}

.sidebar-resizer:hover .sidebar-resizer-handle,
.sidebar-resizer:focus-visible .sidebar-resizer-handle {
  background: #64748b;
}

.dark .sidebar-resizer-handle {
  background: #475569;
}

.dark .sidebar-resizer:hover .sidebar-resizer-handle,
.dark .sidebar-resizer:focus-visible .sidebar-resizer-handle,
.dark .sidebar-resizer.is-resizing .sidebar-resizer-handle {
  background: #94a3b8;
}

@media (prefers-reduced-motion: reduce) {
  .sidebar-resizer::before,
  .sidebar-resizer-handle {
    transition: none;
  }
}
</style>
