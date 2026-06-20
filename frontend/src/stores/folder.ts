import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Folder, FolderCreateRequest, FolderUpdateRequest } from '@/api/folder'
import * as folderApi from '@/api/folder'

export const useFolderStore = defineStore('folder', () => {
  const folders = ref<Folder[]>([])
  const loading = ref(false)

  async function fetchFolders() {
    loading.value = true
    try {
      const response: any = await folderApi.getFolders()
      folders.value = response.folders || []
    } finally {
      loading.value = false
    }
  }

  async function createFolder(data: FolderCreateRequest) {
    await folderApi.createFolder(data)
    await fetchFolders()
  }

  async function updateFolderFn(id: number, data: FolderUpdateRequest) {
    await folderApi.updateFolder(id, data)
    await fetchFolders()
  }

  async function deleteFolderFn(id: number) {
    await folderApi.deleteFolder(id)
    await fetchFolders()
  }

  return {
    folders,
    loading,
    fetchFolders,
    createFolder,
    updateFolder: updateFolderFn,
    deleteFolder: deleteFolderFn,
  }
})
