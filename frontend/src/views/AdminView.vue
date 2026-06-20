<template>
  <div class="h-screen flex flex-col bg-gray-50 dark:bg-[#101014]">
    <header class="bg-white dark:bg-[#1a1a1e] shadow-sm border-b dark:border-[#2d2d33] flex-shrink-0">
      <div class="px-3 sm:px-4 py-2 flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <n-button quaternary circle @click="$router.push('/')">
            <template #icon>
              <n-icon :component="ArrowBackOutline" />
            </template>
          </n-button>
          <h1 class="text-lg font-bold text-gray-800 dark:text-gray-200">管理员</h1>
        </div>
      </div>
    </header>

    <div class="flex-1 overflow-auto px-3 sm:px-4 py-4">
      <n-tabs type="line" animated>
        <n-tab-pane name="config" tab="系统配置">
          <n-card class="mt-4" title="配置项">
            <n-list>
              <n-list-item>
                <n-thing title="允许注册">
                  <template #description>
                    关闭后，新用户将无法注册账号
                  </template>
                  <template #header-extra>
                    <n-switch v-model:value="configForm.allow_register" :loading="configLoading" />
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
            <div class="mt-4 flex justify-end">
              <n-button type="primary" :loading="configSaving" @click="saveConfig">保存配置</n-button>
            </div>
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="users" tab="用户管理">
          <div class="mt-4">
            <div class="flex justify-end mb-3">
              <n-button type="primary" size="small" @click="openCreateUser">
                <template #icon>
                  <n-icon :component="PersonAddOutline" />
                </template>
                新增用户
              </n-button>
            </div>
            <n-data-table
              :columns="userColumns"
              :data="userStore.users"
              :loading="userLoading"
              :pagination="{ pageSize: 10 }"
              :scroll-x="640"
              size="small"
            />
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- 编辑用户弹窗 -->
    <n-modal v-model:show="showEditUserModal" preset="card" title="编辑用户" class="w-[90vw] max-w-[400px]">
      <n-form :model="editUserForm">
        <n-form-item label="昵称">
          <n-input v-model:value="editUserForm.nickname" placeholder="请输入昵称" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="editUserForm.email" placeholder="请输入邮箱" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showEditUserModal = false">取消</n-button>
          <n-button type="primary" :loading="editUserLoading" @click="submitEditUser">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 新增用户弹窗 -->
    <n-modal v-model:show="showCreateUserModal" preset="card" title="新增用户" class="w-[90vw] max-w-[400px]">
      <n-form :model="createUserForm">
        <n-form-item label="用户名" required>
          <n-input v-model:value="createUserForm.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码" required>
          <n-input v-model:value="createUserForm.password" type="password" placeholder="请输入密码（至少6位）" show-password-on="click" />
        </n-form-item>
        <n-form-item label="昵称">
          <n-input v-model:value="createUserForm.nickname" placeholder="可选，默认与用户名相同" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="createUserForm.email" placeholder="可选" />
        </n-form-item>
        <p class="text-xs text-gray-400 dark:text-gray-500 -mt-1">系统仅允许一个管理员，新增用户均为普通用户。</p>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateUserModal = false">取消</n-button>
          <n-button type="primary" :loading="createUserLoading" @click="submitCreateUser">创建</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 重置密码弹窗 -->
    <n-modal v-model:show="showResetPwdModal" preset="card" title="重置密码" class="w-[90vw] max-w-[400px]">
      <n-form>
        <n-form-item label="新密码">
          <n-input v-model:value="resetPwdForm.password" type="password" placeholder="请输入新密码" show-password-on="click" />
        </n-form-item>
        <n-form-item label="确认密码">
          <n-input v-model:value="resetPwdForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password-on="click" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showResetPwdModal = false">取消</n-button>
          <n-button type="primary" :loading="resetPwdLoading" @click="submitResetPassword">确认重置</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  NIcon,
  NSwitch,
  NList,
  NListItem,
  NThing,
  NTabs,
  NTabPane,
  NCard,
  NDataTable,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  NTag,
  useMessage,
  useDialog,
} from 'naive-ui'
import { ArrowBackOutline, PersonAddOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { getConfigs, updateConfig } from '@/api/config'
import { hashPassword } from '@/utils/crypto'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const authStore = useAuthStore()
const userStore = useUserStore()

const configLoading = ref(false)
const configSaving = ref(false)
const configForm = ref({ allow_register: true })

const userLoading = ref(false)
const showEditUserModal = ref(false)
const editUserLoading = ref(false)
const editUserForm = ref({ id: 0, nickname: '', email: '' })

const showResetPwdModal = ref(false)
const resetPwdLoading = ref(false)
const resetPwdForm = ref({ id: 0, password: '', confirmPassword: '' })

const showCreateUserModal = ref(false)
const createUserLoading = ref(false)
const createUserForm = ref({ username: '', password: '', nickname: '', email: '' })

const userColumns = computed(() => [
  { title: '用户名', key: 'username', minWidth: 100, resizable: true },
  { title: '昵称', key: 'nickname', minWidth: 80, resizable: true },
  { title: '邮箱', key: 'email', ellipsis: { tooltip: true } },
  {
    title: '角色',
    key: 'is_admin',
    width: 80,
    render: (row: any) => h(NTag, { type: row.is_admin ? 'error' : 'info', size: 'small' }, { default: () => row.is_admin ? '管理员' : '普通' }),
  },
  {
    title: '状态',
    key: 'is_active',
    width: 70,
    render: (row: any) => h(NTag, { type: row.is_active ? 'success' : 'default', size: 'small' }, { default: () => row.is_active ? '启用' : '禁用' }),
  },
  {
    title: '操作',
    key: 'actions',
    width: 230,
    fixed: 'right' as const,
    render: (row: any) => h('div', { class: 'flex items-center gap-1' }, [
      h(NButton, { size: 'small', secondary: true, onClick: () => openEditUser(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', secondary: true, type: 'warning', onClick: () => openResetPwd(row) }, { default: () => '重置密码' }),
      h(
        NButton,
        {
          size: 'small',
          type: row.is_active ? 'error' : 'success',
          secondary: true,
          disabled: row.is_admin,
          onClick: () => handleToggleUser(row),
        },
        { default: () => row.is_active ? '禁用' : '启用' }
      ),
      h(
        NButton,
        {
          size: 'small',
          type: 'error',
          secondary: true,
          disabled: row.is_admin || row.id === authStore.user?.id,
          onClick: () => handleDeleteUser(row),
        },
        { default: () => '删除' }
      ),
    ]),
  },
])

async function loadConfigs() {
  configLoading.value = true
  try {
    const res = await getConfigs()
    const configs = res.configs || {}
    configForm.value.allow_register = configs.allow_register !== '0'
  } catch {
    // ignore
  } finally {
    configLoading.value = false
  }
}

async function saveConfig() {
  configSaving.value = true
  try {
    await updateConfig('allow_register', configForm.value.allow_register ? '1' : '0')
    message.success('配置已保存')
  } finally {
    configSaving.value = false
  }
}

async function loadUsers() {
  userLoading.value = true
  try {
    await userStore.fetchUsers()
  } finally {
    userLoading.value = false
  }
}

function openEditUser(row: any) {
  editUserForm.value = { id: row.id, nickname: row.nickname || '', email: row.email || '' }
  showEditUserModal.value = true
}

async function submitEditUser() {
  editUserLoading.value = true
  try {
    await userStore.updateUserById(editUserForm.value.id, {
      nickname: editUserForm.value.nickname,
      email: editUserForm.value.email,
    })
    message.success('更新成功')
    showEditUserModal.value = false
  } finally {
    editUserLoading.value = false
  }
}

function openResetPwd(row: any) {
  resetPwdForm.value = { id: row.id, password: '', confirmPassword: '' }
  showResetPwdModal.value = true
}

async function submitResetPassword() {
  if (!resetPwdForm.value.password) { message.warning('请输入新密码'); return }
  if (resetPwdForm.value.password.length < 6) { message.warning('密码至少6位'); return }
  if (resetPwdForm.value.password !== resetPwdForm.value.confirmPassword) { message.warning('两次密码不一致'); return }

  resetPwdLoading.value = true
  try {
    await userStore.resetPasswordById(resetPwdForm.value.id, resetPwdForm.value.password)
    message.success('密码重置成功')
    showResetPwdModal.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || '重置失败')
  } finally {
    resetPwdLoading.value = false
  }
}

async function handleToggleUser(row: any) {
  try {
    await userStore.toggleUserStatusById(row.id, !row.is_active)
    message.success(row.is_active ? '已禁用' : '已启用')
  } catch (e: any) {
    message.error(e.response?.data?.error || '操作失败')
  }
}

function openCreateUser() {
  createUserForm.value = { username: '', password: '', nickname: '', email: '' }
  showCreateUserModal.value = true
}

async function submitCreateUser() {
  const f = createUserForm.value
  if (!f.username.trim()) { message.warning('请输入用户名'); return }
  if (!f.password) { message.warning('请输入密码'); return }
  if (f.password.length < 6) { message.warning('密码至少6位'); return }

  createUserLoading.value = true
  try {
    await userStore.createNewUser({
      username: f.username.trim(),
      password: f.password,
      nickname: f.nickname.trim(),
      email: f.email.trim(),
    })
    message.success('用户创建成功')
    showCreateUserModal.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || '创建失败')
  } finally {
    createUserLoading.value = false
  }
}

function handleDeleteUser(row: any) {
  dialog.warning({
    title: '删除用户',
    content: `确定删除用户「${row.username}」吗？该用户的所有备忘录将一并删除，操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await userStore.deleteUserById(row.id)
        message.success('删除成功')
      } catch (e: any) {
        message.error(e.response?.data?.error || '删除失败')
      }
    },
  })
}

onMounted(async () => {
  await Promise.all([loadConfigs(), loadUsers()])
})
</script>
