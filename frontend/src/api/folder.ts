import api from '@/utils/request'

export interface Folder {
  id: number
  user_id: number
  parent_id: number | null
  name: string
  position: number
  created_at: string
  updated_at: string
}

export interface FolderCreateRequest {
  name: string
  parent_id?: number | null
}

export interface FolderUpdateRequest {
  name?: string
  parent_id?: number | null
  position?: number
}

export async function getFolders() {
  return api.get('/folders')
}

export async function createFolder(data: FolderCreateRequest) {
  return api.post('/folders', data)
}

export async function updateFolder(id: number, data: FolderUpdateRequest) {
  return api.put(`/folders/${id}`, data)
}

export async function deleteFolder(id: number) {
  return api.delete(`/folders/${id}`)
}
