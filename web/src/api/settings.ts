import type { ServicesConfig } from '@/types'
import { apiClient } from './client'

export interface SetLangRequest {
  language: string
}

export interface LoadSettingsResponse {
  settings: {
    enableFTP?: boolean
    enableWebDAV?: boolean
    enableNFS?: boolean
    enableMCP?: boolean
    ftpPort?: string | number
    webDAVPort?: string | number
    nfsPort?: string | number
  }
}

export interface SaveServicesRequest {
  enableFTP: boolean
  enableWebDAV: boolean
  enableNFS: boolean
  enableMCP: boolean
  ftpPort: string
  webDAVPort: string
  nfsPort: string
}

export async function setLang(lang: string): Promise<void> {
  await apiClient.post('/api/settings/language', { language: lang })
}

export async function loadSettings(): Promise<LoadSettingsResponse> {
  return await apiClient.get<LoadSettingsResponse>('/api/settings')
}

export async function saveServices(
  config: ServicesConfig
): Promise<any> {
  return await apiClient.post('/api/settings/services', {
    enableFTP: config.enableFTP,
    enableWebDAV: config.enableWebDAV,
    enableNFS: config.enableNFS,
    enableMCP: config.enableMCP,
    ftpPort: config.ftpPort,
    webDAVPort: config.webDAVPort,
    nfsPort: config.nfsPort,
  })
}
