import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAllUsers, updateUser, toggleUserStatus, resetUserPassword, createUser, deleteUser } from '@/api/user'
import type { User, CreateUserRequest } from '@/api/user'
import { hashPassword } from '@/utils/crypto'

export const useUserStore = defineStore('user', () => {
  const users = ref<User[]>([])

  async function fetchUsers() {
    const res = await getAllUsers()
    users.value = res.users || []
  }

  async function updateUserById(id: number, data: { nickname?: string; email?: string; avatar?: string }) {
    await updateUser(id, data)
    await fetchUsers()
  }

  async function toggleUserStatusById(id: number, isActive: boolean) {
    await toggleUserStatus(id, isActive)
    await fetchUsers()
  }

  async function resetPasswordById(id: number, newPassword: string) {
    await resetUserPassword(id, hashPassword(newPassword))
  }

  async function createNewUser(data: CreateUserRequest) {
    await createUser({ ...data, password: hashPassword(data.password) })
    await fetchUsers()
  }

  async function deleteUserById(id: number) {
    await deleteUser(id)
    await fetchUsers()
  }

  return {
    users,
    fetchUsers,
    updateUserById,
    toggleUserStatusById,
    resetPasswordById,
    createNewUser,
    deleteUserById,
  }
})
