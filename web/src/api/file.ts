import type { FileInfo, TreeNode } from '@/types'
import { apiClient } from './client'

export interface LoadFilesResponse {
  files: FileInfo[]
  path: string
}

export interface LoadTreeResponse {
  tree: TreeNode
}

export async function loadFiles(path: string = '/'): Promise<LoadFilesResponse> {
  const url = `/api/files?path=${encodeURIComponent(path)}`
  return await apiClient.get<LoadFilesResponse>(url)
}

export async function loadTree(): Promise<TreeNode> {
  const response = await apiClient.get<LoadTreeResponse>('/api/tree')
  return response.tree
}

export async function uploadFiles(filesToUpload: FileList, path: string): Promise<boolean> {
  const formData = new FormData()
  formData.append('path', path)
  for (let i = 0; i < filesToUpload.length; i++) {
    const file = filesToUpload[i]
    if (file) {
      formData.append('files', file)
    }
  }

  try {
    await apiClient.postFormData('/api/upload', formData)
    return true
  } catch (error) {
    return false
  }
}

export function getDownloadUrl(path: string): string {
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  return `/download/${cleanPath}`
}

export function copyLink(
  path: string,
  httpAuthOn: boolean,
  httpAuthUser: string,
  httpAuthPass: string
): string {
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  let url = `${window.location.origin}/download/${cleanPath}`
  if (httpAuthOn && httpAuthUser && httpAuthPass) {
    url = `${window.location.origin.replace('://', `://${encodeURIComponent(httpAuthUser)}:${encodeURIComponent(httpAuthPass)}@`)}/download/${cleanPath}`
  }

  if (navigator.clipboard) {
    navigator.clipboard.writeText(url)
  } else {
    const ta = document.createElement('textarea')
    ta.value = url
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
  }
  return url
}
