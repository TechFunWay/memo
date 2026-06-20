import api from '@/utils/request'

export interface Memo {
  id: number
  user_id: number
  folder_id: number | null
  title: string
  content: string
  content_html: string
  is_pinned: boolean
  is_favorite: boolean
  tags: string
  position: number
  created_at: string
  updated_at: string
}

export interface MemoCreateRequest {
  title: string
  content: string
  folder_id?: number | null
  tags?: string
}

export interface MemoUpdateRequest {
  title?: string
  content?: string
  folder_id?: number | null
  tags?: string
  is_pinned?: boolean
  is_favorite?: boolean
  position?: number
}

export async function getMemos(params?: { search?: string; filter?: string; folder_id?: string }) {
  return api.get('/memos', { params })
}

export async function getMemo(id: number) {
  return api.get(`/memos/${id}`)
}

export async function createMemo(data: MemoCreateRequest) {
  return api.post('/memos', data)
}

export async function updateMemo(id: number, data: MemoUpdateRequest) {
  return api.put(`/memos/${id}`, data)
}

export async function deleteMemo(id: number) {
  return api.delete(`/memos/${id}`)
}
