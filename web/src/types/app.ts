export interface FileInfo {
  name: string
  size: string
  modTime: string
  isDir: boolean
  path: string
  canPreview?: boolean
  fileType?: number
  icon?: string
  mimeType?: string
}

export interface TreeNode {
  name: string
  path: string
  isDir: boolean
  children?: TreeNode[]
}

export interface AppConfig {
  language: string
  httpAuthOn: boolean
  httpAuthUser: string
  httpAuthPass: string
  files: FileInfo[]
  currentPath: string
}

export interface ServicesConfig {
  enableFTP: boolean
  enableWebDAV: boolean
  enableNFS: boolean
  enableMCP: boolean
  ftpPort: string
  webDAVPort: string
  nfsPort: string
  mcpKey: string
}

export type LocaleType = 'en' | 'zh'
