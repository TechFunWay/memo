import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Memo, MemoCreateRequest, MemoUpdateRequest } from '@/api/memo'
import * as memoApi from '@/api/memo'

export const useMemoStore = defineStore('memo', () => {
  const memos = ref<Memo[]>([])
  const loading = ref(false)
  const currentMemo = ref<Memo | null>(null)

  async function fetchMemos(params?: { search?: string; filter?: string; folder_id?: string }) {
    loading.value = true
    try {
      const response: any = await memoApi.getMemos(params)
      memos.value = response.memos || []
    } finally {
      loading.value = false
    }
  }

  async function createMemo(data: MemoCreateRequest) {
    const response: any = await memoApi.createMemo(data)
    await fetchMemos()
    return response
  }

  async function updateMemoFn(id: number, data: MemoUpdateRequest) {
    await memoApi.updateMemo(id, data)
    await fetchMemos()
  }

  async function deleteMemoFn(id: number) {
    await memoApi.deleteMemo(id)
    await fetchMemos()
  }

  async function fetchMemoDetail(id: number) {
    currentMemo.value = await memoApi.getMemo(id)
  }

  return {
    memos,
    loading,
    currentMemo,
    fetchMemos,
    createMemo,
    updateMemo: updateMemoFn,
    deleteMemo: deleteMemoFn,
    fetchMemoDetail,
  }
})
