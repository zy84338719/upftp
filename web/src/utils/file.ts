import { ICON_MAP, FILE_TYPE_MAP, ICON_FOLDER, ICON_DEFAULT, FileType } from '@/constants'
import { getDownloadUrl as apiGetDownloadUrl } from '@/api'

export function getFileIcon(name: string, isDir: boolean): string {
  if (isDir) return ICON_FOLDER

  const ext = (name.split('.').pop() || '').toLowerCase()
  return ICON_MAP[ext] || ICON_DEFAULT
}

export function getFileType(name: string): number {
  const ext = (name.split('.').pop() || '').toLowerCase()
  return FILE_TYPE_MAP[ext] || FileType.Unknown
}

export function getDownloadUrl(path: string): string {
  return apiGetDownloadUrl(path)
}
