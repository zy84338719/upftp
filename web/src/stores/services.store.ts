import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServicesConfig } from '@/types'
import { DEFAULT_PORTS } from '@/constants'
import {
  loadSettings as apiLoadSettings,
  saveServices as apiSaveServices,
} from '@/api'
import { generateMCPKey as generateMCPKeyUtil } from '@/utils'

export const useServicesStore = defineStore('services', () => {
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

  async function loadSettings() {
    try {
      const data = await apiLoadSettings()
      const s = data.settings as ServicesConfig || {}
      
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

  function generateMCPKey(): string {
    const key = generateMCPKeyUtil()
    services.value.mcpKey = key
    return key
  }

  return {
    services,
    loadSettings,
    saveServices,
    generateMCPKey
  }
})
