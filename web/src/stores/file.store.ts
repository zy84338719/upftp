import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { FileInfo, TreeNode } from '@/types'
import { DEFAULT_PATH } from '@/constants'
import {
  loadFiles as apiLoadFiles,
  loadTree as apiLoadTree,
  uploadFiles as apiUploadFiles,
} from '@/api'
import { getFileIcon, getFileType, getDownloadUrl } from '@/utils'
import { useUIStore } from './ui.store'

export const useFileStore = defineStore('file', () => {
  const files = ref<FileInfo[]>([])
  const currentPath = ref(DEFAULT_PATH)
  const treeData = ref<TreeNode | null>(null)

  const filteredFiles = computed(() => {
    const uiStore = useUIStore()
    if (!uiStore.searchQuery) return files.value
    const query = uiStore.searchQuery.toLowerCase()
    return files.value.filter(f => f.name.toLowerCase().includes(query))
  })

  async function loadFiles(path: string = '/') {
    const uiStore = useUIStore()
    uiStore.isLoading = true
    try {
      const data = await apiLoadFiles(path)
      files.value = data.files || []
      currentPath.value = data.path || path
      return data
    } finally {
      uiStore.isLoading = false
    }
  }

  async function loadTree() {
    try {
      const tree = await apiLoadTree()
      treeData.value = tree
      return tree
    } catch (error) {
      console.error('Failed to load tree:', error)
      return null
    }
  }

  async function uploadFiles(filesToUpload: FileList, path: string) {
    const success = await apiUploadFiles(filesToUpload, path)
    if (success) {
      await loadFiles(currentPath.value)
    }
    return success
  }

  function initConfig(filesList: FileInfo[], path: string) {
    files.value = filesList || []
    currentPath.value = path || DEFAULT_PATH
  }

  return {
    files,
    currentPath,
    treeData,
    filteredFiles,
    loadFiles,
    loadTree,
    uploadFiles,
    getFileIcon,
    getFileType,
    getDownloadUrl,
    initConfig
  }
})
