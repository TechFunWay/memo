import api from '@/utils/request'

export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  is_active: boolean
  is_admin: boolean
  created_at: string
  updated_at: string
}

export interface UpdateUserRequest {
  nickname?: string
  email?: string
  avatar?: string
}

export interface CreateUserRequest {
  username: string
  password: string
  nickname?: string
  email?: string
  is_admin?: boolean
}

export async function getAllUsers() {
  return api.get('/users')
}

export async function createUser(data: CreateUserRequest) {
  return api.post('/users', data)
}

export async function deleteUser(id: number) {
  return api.delete(`/users/${id}`)
}

export async function updateUser(id: number, data: UpdateUserRequest) {
  return api.put(`/users/${id}`, data)
}

export async function toggleUserStatus(id: number, isActive: boolean) {
  return api.post(`/users/${id}/toggle`, { is_active: isActive ? 1 : 0 })
}

export async function resetUserPassword(id: number, newPassword: string) {
  return api.post(`/users/${id}/reset-password`, { new_password: newPassword })
}
