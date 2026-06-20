<template>
  <div class="h-screen flex flex-col bg-gray-50 dark:bg-[#101014]">
    <!-- 顶部工具栏 -->
    <header class="bg-white dark:bg-[#1a1a1e] shadow-sm border-b dark:border-[#2d2d33] flex-shrink-0">
      <div class="px-3 sm:px-4 py-2 flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <n-button quaternary circle @click="goBack">
            <template #icon>
              <n-icon :component="ArrowBackOutline" />
            </template>
          </n-button>
          <h1 class="text-lg font-bold text-gray-800 dark:text-gray-200">
            {{ isEdit ? '编辑备忘录' : '新建备忘录' }}
          </h1>
        </div>
        <div class="flex items-center space-x-2">
          <span v-if="saving" class="text-xs text-gray-400">保存中...</span>
          <span v-else-if="lastSaved" class="text-xs text-gray-400">已保存</span>
          <n-button type="primary" :loading="saving" @click="saveMemo">保存</n-button>
        </div>
      </div>
    </header>

    <!-- 标题输入 -->
    <div class="bg-white dark:bg-[#1a1a1e] border-b dark:border-[#2d2d33] flex-shrink-0 px-3 sm:px-4 py-3">
      <n-input
        v-model:value="title"
        placeholder="请输入标题..."
        size="large"
        class="!text-xl !font-semibold"
        :input-props="{ style: 'background: transparent;' }"
      />
    </div>

    <!-- 文件夹选择 -->
    <div class="bg-white dark:bg-[#1a1a1e] border-b dark:border-[#2d2d33] flex-shrink-0 px-3 sm:px-4 py-2">
      <div class="flex items-center space-x-2">
        <span class="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">文件夹</span>
        <n-select
          v-model:value="folderId"
          :options="folderOptions"
          placeholder="选择文件夹（可选）"
          clearable
          class="flex-1"
          size="small"
        />
      </div>
    </div>

    <!-- 格式工具栏 -->
    <div class="bg-white dark:bg-[#1a1a1e] border-b dark:border-[#2d2d33] flex-shrink-0 px-3 sm:px-4 py-2">
      <div class="flex items-center space-x-1 flex-wrap gap-1">
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('bold')" :disabled="saving">
              <span class="font-bold text-sm">B</span>
            </n-button>
          </template>
          粗体 (Ctrl+B)
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('italic')" :disabled="saving">
              <span class="italic text-sm">I</span>
            </n-button>
          </template>
          斜体 (Ctrl+I)
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('underline')" :disabled="saving">
              <span class="underline text-sm">U</span>
            </n-button>
          </template>
          下划线 (Ctrl+U)
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('strikeThrough')" :disabled="saving">
              <span class="line-through text-sm">S</span>
            </n-button>
          </template>
          删除线
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="triggerColorPicker('foreColor')" :disabled="saving">
              <span class="text-sm" style="background: linear-gradient(to right, red, blue); -webkit-background-clip: text; -webkit-text-fill-color: transparent;">A</span>
            </n-button>
          </template>
          文字颜色
        </n-tooltip>
        <input ref="colorInputRef" type="color" class="hidden" @change="handleColorChange" />
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="triggerColorPicker('hiliteColor')" :disabled="saving">
              <span class="text-sm" style="background: #FFD700; padding: 0 4px; border-radius: 2px;">A</span>
            </n-button>
          </template>
          背景色
        </n-tooltip>
        <n-divider vertical style="height: 24px" />
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('insertUnorderedList')" :disabled="saving">
              <template #icon>
                <n-icon :component="List" />
              </template>
            </n-button>
          </template>
          无序列表
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('insertOrderedList')" :disabled="saving">
              <template #icon>
                <n-icon :component="ListOutline" />
              </template>
            </n-button>
          </template>
          有序列表
        </n-tooltip>
        <n-divider vertical style="height: 24px" />
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('formatBlock', '<h2>')" :disabled="saving">
              <template #icon>
                <n-icon :component="Text" />
              </template>
            </n-button>
          </template>
          标题
        </n-tooltip>
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="execCmd('removeFormat')" :disabled="saving">
              <template #icon>
                <n-icon :component="BrushOutline" />
              </template>
            </n-button>
          </template>
          清除格式
        </n-tooltip>
        <n-divider vertical style="height: 24px" />
        <n-tooltip placement="top" trigger="hover">
          <template #trigger>
            <n-button quaternary size="small" @click="triggerImageUpload" :disabled="saving">
              <template #icon>
                <n-icon :component="ImageOutline" />
              </template>
            </n-button>
          </template>
          上传图片
        </n-tooltip>
        <input
          ref="imageInputRef"
          type="file"
          accept="image/*"
          class="hidden"
          @change="handleImageUpload"
        />
      </div>
    </div>

    <!-- 编辑器区域 -->
    <div
      ref="editorRef"
      class="flex-1 overflow-auto px-3 sm:px-4 py-4 outline-none text-gray-800 dark:text-gray-200"
      contenteditable
      @input="handleInput"
      @paste="handlePaste"
      @keydown="handleKeydown"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NIcon, NInput, NTooltip, NSelect, NDivider, useMessage, useDialog } from 'naive-ui'
import {
  ArrowBackOutline,
  List,
  ListOutline,
  Text,
  BrushOutline,
  ImageOutline,
} from '@vicons/ionicons5'
import { useMemoStore } from '@/stores/memo'
import { useFolderStore } from '@/stores/folder'
import { uploadImage } from '@/api/upload'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const memoStore = useMemoStore()
const folderStore = useFolderStore()

const editorRef = ref<HTMLElement | null>(null)
const imageInputRef = ref<HTMLInputElement | null>(null)
const colorInputRef = ref<HTMLInputElement | null>(null)
const title = ref('')
const folderId = ref<number | null>(null)
const saving = ref(false)
const lastSaved = ref(false)
const isEdit = ref(false)
const memoId = ref<number | null>(null)
const currentColorCmd = ref('')

const htmlContent = ref('')
const dirty = ref(false)
let autoSaveTimer: ReturnType<typeof setTimeout> | null = null

const folderOptions = computed(() =>
  folderStore.folders.map((f) => ({ label: f.name, value: f.id }))
)

function execCmd(command: string, value?: string) {
  editorRef.value?.focus()
  document.execCommand(command, false, value)
  handleInput()
}

function triggerColorPicker(command: string) {
  currentColorCmd.value = command
  colorInputRef.value?.click()
}

function handleColorChange(event: Event) {
  const input = event.target as HTMLInputElement
  const color = input.value
  if (color && currentColorCmd.value) {
    editorRef.value?.focus()
    document.execCommand(currentColorCmd.value, false, color)
    handleInput()
  }
  input.value = ''
}

function triggerImageUpload() {
  imageInputRef.value?.click()
}

async function handleImageUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  await insertImage(file)
  input.value = ''
}

async function insertImage(file: File) {
  if (!file.type.startsWith('image/')) {
    message.warning('请选择图片文件')
    return
  }

  try {
    // 上传到服务器并插入返回的 URL
    const url = await uploadImage(file)
    editorRef.value?.focus()
    document.execCommand('insertImage', false, url)
    handleInput()
  } catch {
    message.error('图片上传失败')
  }
}

async function handlePaste(event: ClipboardEvent) {
  const items = event.clipboardData?.items
  if (!items) return

  for (const item of items) {
    if (item.type.startsWith('image/')) {
      event.preventDefault()
      const file = item.getAsFile()
      if (file) {
        await insertImage(file)
      }
      return
    }
  }
}

function handleInput() {
  htmlContent.value = editorRef.value?.innerHTML || ''
  dirty.value = true
  triggerAutoSave()
}

function triggerAutoSave() {
  if (!isEdit.value || !memoId.value) return
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    doAutoSave()
  }, 2000)
}

async function doAutoSave() {
  const content = editorRef.value?.innerHTML || ''
  const plainText = editorRef.value?.innerText?.trim() || ''
  if (!title.value && !plainText) return
  if (!memoId.value) return

  saving.value = true
  try {
    await memoStore.updateMemo(memoId.value, {
      title: title.value,
      content,
      folder_id: folderId.value,
    })
    dirty.value = false
    lastSaved.value = true
    setTimeout(() => { lastSaved.value = false }, 3000)
  } catch {
    // 静默失败，避免干扰用户
  } finally {
    saving.value = false
  }
}

function handleKeydown(event: KeyboardEvent) {
  // Ctrl+S / Cmd+S 保存
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    saveMemo()
  }
}

async function saveMemo() {
  const content = editorRef.value?.innerHTML || ''
  const plainText = editorRef.value?.innerText?.trim() || ''

  if (!title.value && !plainText) {
    message.warning('请输入标题或内容')
    return
  }

  saving.value = true
  try {
    if (isEdit.value && memoId.value) {
      await memoStore.updateMemo(memoId.value, {
        title: title.value,
        content,
        folder_id: folderId.value,
      })
      dirty.value = false
      message.success('更新成功')
    } else {
      const res = await memoStore.createMemo({
        title: title.value,
        content,
        folder_id: folderId.value,
      })
      // 新创建后转为编辑模式，开启自动保存
      isEdit.value = true
      memoId.value = res.id
      dirty.value = false
      message.success('保存成功')
    }
    lastSaved.value = true
    setTimeout(() => { lastSaved.value = false }, 3000)
  } finally {
    saving.value = false
  }
}

function goBack() {
  if (dirty.value) {
    dialog.warning({
      title: '提示',
      content: '内容未保存，确定要离开吗？',
      positiveText: '确定离开',
      negativeText: '继续编辑',
      onPositiveClick: () => {
        router.push('/')
      },
    })
    return
  }
  router.push('/')
}

onMounted(async () => {
  await folderStore.fetchFolders()

  const id = route.params.id as string
  if (id && id !== 'new') {
    isEdit.value = true
    memoId.value = parseInt(id, 10)
    // 先从 store 中找，如果找不到再从 API 加载
    let memo = memoStore.memos.find((m) => m.id === memoId.value)
    if (!memo) {
      try {
        await memoStore.fetchMemoDetail(memoId.value)
        memo = memoStore.currentMemo
      } catch {
        message.error('备忘录不存在')
        router.push('/')
        return
      }
    }
    if (memo) {
      title.value = memo.title
      folderId.value = memo.folder_id
      if (editorRef.value) {
        editorRef.value.innerHTML = memo.content || ''
        htmlContent.value = memo.content || ''
      }
    }
  }
})

onUnmounted(() => {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
})
</script>

<style scoped>
.editor-content:empty::before {
  content: attr(placeholder);
  color: #9ca3af;
  pointer-events: none;
}

.editor-content {
  min-height: 200px;
  line-height: 1.6;
  font-size: 15px;
  color: #374151;
}

.dark .editor-content {
  color: #d1d5db;
}

.editor-content img {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 8px;
  margin: 8px 0;
  display: block;
}

.editor-content ul,
.editor-content ol {
  padding-left: 1.5em;
}

.editor-content h2 {
  font-size: 1.25em;
  font-weight: 600;
  margin: 0.5em 0;
}
</style>
