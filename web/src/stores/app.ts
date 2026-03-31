import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import i18n from '@/i18n'
import type { FileInfo, TreeNode, AppConfig, ServicesConfig, LocaleType } from '@/types'
import {
  DEFAULT_PORTS,
  DEFAULT_LANGUAGE,
  DEFAULT_PATH,
  LOCAL_STORAGE_KEYS
} from '@/constants'
import {
  loadFiles as apiLoadFiles,
  loadTree as apiLoadTree,
  uploadFiles as apiUploadFiles,
  setLang as apiSetLang,
  loadSettings as apiLoadSettings,
  saveServices as apiSaveServices,
} from '@/api'
import { getFileIcon, getFileType, getDownloadUrl, copyLink as copyLinkUtil, generateMCPKey as generateMCPKeyUtil } from '@/utils'

export type { FileInfo, TreeNode, AppConfig, ServicesConfig }

export const useAppStore = defineStore('app', () => {
  // State
  const language = ref<LocaleType>((localStorage.getItem(LOCAL_STORAGE_KEYS.language) as LocaleType) || DEFAULT_LANGUAGE)
  const httpAuthOn = ref(false)
  const files = ref<FileInfo[]>([])
  const currentPath = ref(DEFAULT_PATH)
  const treeData = ref<TreeNode | null>(null)
  const searchQuery = ref('')
  const isLoading = ref(false)

  // Services state
  const services = ref<ServicesConfig>({
    enableFTP: false,
    enableWebDAV: false,
    enableNFS: false,
    enableMCP: false,
    ftpPort: DEFAULT_PORTS.ftp,
    webDAVPort: DEFAULT_PORTS.webDAV,
    nfsPort: DEFAULT_PORTS.nfs,
    mcpKey: ''
  })

  // Getters
  const filteredFiles = computed(() => {
    if (!searchQuery.value) return files.value
    const query = searchQuery.value.toLowerCase()
    return files.value.filter(f => f.name.toLowerCase().includes(query))
  })

  // Actions
  async function setLang(lang: string) {
    const locale = lang as LocaleType
    language.value = locale
    localStorage.setItem(LOCAL_STORAGE_KEYS.language, lang)
    i18n.global.locale.value = locale
    await apiSetLang(lang)
  }

  async function loadFiles(path: string = '/') {
    isLoading.value = true
    try {
      const data = await apiLoadFiles(path)
      files.value = data.files || []
      currentPath.value = data.path || path
      return data
    } finally {
      isLoading.value = false
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

  function copyLink(path: string) {
    return copyLinkUtil(path, httpAuthOn.value, '', '')
  }

  async function loadSettings() {
    try {
      const data = await apiLoadSettings()
      const s = data.settings as ServicesConfig || {}
      
      // Update services config
      services.value = {
        enableFTP: s.enableFTP || false,
        enableWebDAV: s.enableWebDAV || false,
        enableNFS: s.enableNFS || false,
        enableMCP: s.enableMCP || false,
        ftpPort: s.ftpPort?.toString() || DEFAULT_PORTS.ftp,
        webDAVPort: s.webDAVPort?.toString() || DEFAULT_PORTS.webDAV,
        nfsPort: s.nfsPort?.toString() || DEFAULT_PORTS.nfs,
        mcpKey: s.mcpKey || ''
      }
      
      return data
    } catch (error) {
      console.error('Failed to load settings:', error)
      return null
    }
  }

  async function saveServices(config: ServicesConfig) {
    return await apiSaveServices(config)
  }



  function initConfig(config: AppConfig) {
    language.value = (config.language as LocaleType) || DEFAULT_LANGUAGE
    httpAuthOn.value = config.httpAuthOn || false
    files.value = config.files || []
    currentPath.value = config.currentPath || DEFAULT_PATH
    i18n.global.locale.value = language.value
  }

  return {
    language,
    httpAuthOn,
    files,
    currentPath,
    treeData,
    searchQuery,
    isLoading,
    services,
    filteredFiles,
    setLang,
    loadFiles,
    loadTree,
    uploadFiles,
    getFileIcon,
    getFileType,
    getDownloadUrl,
    copyLink,
    loadSettings,
    saveServices,
    initConfig
  }
})
