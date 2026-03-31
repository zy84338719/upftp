export const FileType = {
  Unknown: 0,
  Image: 1,
  Text: 2,
  Video: 3,
  Audio: 4,
  PDF: 5,
  Document: 6,
  Archive: 7,
  Code: 8
} as const

export type FileTypeValue = typeof FileType[keyof typeof FileType]

export const ICON_MAP: Record<string, string> = {
  jpg: '🖼️', jpeg: '🖼️', png: '🖼️', gif: '🖼️', svg: '🖼️', webp: '🖼️',
  mp4: '🎬', avi: '🎬', mov: '🎬', mkv: '🎬', webm: '🎬',
  mp3: '🎵', wav: '🎵', flac: '🎵', ogg: '🎵', m4a: '🎵',
  pdf: '📄', doc: '📄', docx: '📄', xls: '📄', xlsx: '📄', ppt: '📄', pptx: '📄',
  zip: '📦', rar: '📦', '7z': '📦', tar: '📦', gz: '📦',
  js: '💻', ts: '💻', py: '💻', go: '💻', html: '💻', css: '💻',
  sh: '💻', java: '💻', c: '💻', cpp: '💻', rs: '💻',
  txt: '📝', md: '📝', json: '📝', yaml: '📝', yml: '📝',
  xml: '📝', csv: '📝', log: '📝', toml: '📝'
}

export const FILE_TYPE_MAP: Record<string, FileTypeValue> = {
  jpg: FileType.Image, jpeg: FileType.Image, png: FileType.Image, gif: FileType.Image,
  svg: FileType.Image, webp: FileType.Image, bmp: FileType.Image,
  mp4: FileType.Video, avi: FileType.Video, mov: FileType.Video, mkv: FileType.Video, webm: FileType.Video,
  mp3: FileType.Audio, wav: FileType.Audio, flac: FileType.Audio, ogg: FileType.Audio, m4a: FileType.Audio,
  pdf: FileType.PDF,
  txt: FileType.Text, md: FileType.Text, json: FileType.Text, yaml: FileType.Text, yml: FileType.Text,
  xml: FileType.Text, csv: FileType.Text, log: FileType.Text, toml: FileType.Text,
  js: FileType.Code, ts: FileType.Code, py: FileType.Code, go: FileType.Code, html: FileType.Code,
  css: FileType.Code, sh: FileType.Code, java: FileType.Code, c: FileType.Code, cpp: FileType.Code, rs: FileType.Code
}

export const DEFAULT_PORTS = {
  ftp: '2121',
  webDAV: '8080',
  nfs: '2049'
} as const

export const DEFAULT_LANGUAGE = 'en'
export const DEFAULT_PATH = '/'

export const LOCAL_STORAGE_KEYS = {
  language: 'upftp-lang'
} as const

export const ICON_FOLDER = '📁'
export const ICON_DEFAULT = '📄'

export const MCP_KEY_LENGTH = 43
export const MCP_KEY_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_@'
