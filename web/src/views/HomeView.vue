<template>
  <div class="home">
    <aside class="sidebar">
      <div class="logo">
        <div class="logo-mark">U</div>
        <span class="logo-text">UPFTP</span>
      </div>
      <nav class="nav-section">
        <a href="javascript:void(0)" @click="navigateTo('/')" class="nav-item active" id="navFiles">
          <span>📁</span> <span>{{ $t('all_files') }}</span>
        </a>
      </nav>
      <div class="tree-section">
        <div class="tree-header">
          <span>{{ $t('explorer') }}</span>
        </div>
        <div class="tree-body" id="fileTree">
          <TreeNodeComponent
            v-if="treeData"
            :node="treeData"
            :current-path="currentPath"
            @navigate="navigateTo"
          />
          <div v-else class="tree-loading">...</div>
        </div>
      </div>
      <div class="sidebar-bottom">
        <button class="settings-btn" @click="openSettings">
          <span>⚙️</span> <span>{{ $t('settings') }}</span>
        </button>
        <div class="lang-switch">
          <button
            class="lang-btn"
            :class="{ active: language === 'en' }"
            @click="setLang('en')"
          >
            EN
          </button>
          <button
            class="lang-btn"
            :class="{ active: language === 'zh' }"
            @click="setLang('zh')"
          >
            中文
          </button>
        </div>
      </div>
    </aside>

    <main class="main">
      <div class="header">
        <div class="header-left">
          <div class="breadcrumb" id="breadcrumb">
            <template v-if="currentPath === '/'">
              <a href="javascript:void(0)" @click="navigateTo('/')">/{{ $t('root') }}</a>
            </template>
            <template v-else>
              <a href="javascript:void(0)" @click="navigateTo('/')">/</a>
              <template v-for="(part, index) in breadcrumbParts" :key="index">
                <span v-if="index < breadcrumbParts.length - 1">
                  <a href="javascript:void(0)" @click="navigateTo(part.path)">{{ part.name }}</a> /
                </span>
                <span v-else>{{ part.name }}</span>
              </template>
            </template>
          </div>
          <h1 class="page-title" id="pageTitle">{{ pageTitle }}</h1>
        </div>
        <div class="header-right">
          <input
            type="text"
            class="search"
            :placeholder="$t('search')"
            v-model="searchQuery"
            @input="onSearch"
          />
          <button v-if="canUpload" class="btn-upload" @click="doUpload">
            ↑ <span>{{ $t('upload') }}</span>
          </button>
        </div>
      </div>
      <div class="content" id="content">
        <div class="table-head">
          <div>{{ $t('name') }}</div>
          <div>{{ $t('size') }}</div>
          <div>{{ $t('modified') }}</div>
          <div style="text-align: right">{{ $t('actions') }}</div>
        </div>
        <div id="fileListBody">
          <div v-if="filteredFiles.length === 0" class="empty">
            <p>{{ $t('empty') }}</p>
          </div>
          <div
            v-for="file in filteredFiles"
            :key="file.path"
            class="file-row"
            v-show="!searchQuery || file.name.toLowerCase().includes(searchQuery.toLowerCase())"
          >
            <div class="file-name-cell">
              <span class="file-icon">{{ file.icon || getFileIcon(file.name, file.isDir) }}</span>
              <span
                class="file-name"
                @click="file.isDir ? navigateTo(file.path) : showPreview(file)"
              >
                {{ file.name }}
              </span>
            </div>
            <div class="file-size">{{ file.size }}</div>
            <div class="file-date">{{ file.modTime }}</div>
            <div class="file-actions">
              <button class="act-btn" @click="copyLink(file.path)">{{ $t('copy_link') }}</button>
              <button class="act-btn" @click="dlFile(file.path)">
                {{ file.isDir ? $t('zip_download') : $t('download') }}
              </button>
              <button class="act-btn" @click="showQR(file.path)">{{ $t('qr_download') }}</button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Preview Modal -->
    <div
      class="modal-overlay"
      :class="{ show: showPreviewModal }"
      @click.self="closePreview"
    >
      <div class="modal">
        <div class="modal-header">
          <div class="modal-title">{{ previewTitle }}</div>
          <button class="modal-close" @click="closePreview">&times;</button>
        </div>
        <div class="modal-body">
          <div v-if="previewLoading" class="loading">{{ $t('loading') }}</div>
          <div v-else-if="previewError" class="error">{{ $t('preview_err') }}</div>
          <img
            v-else-if="previewType === 1"
            :src="previewUrl"
            alt=""
            @error="previewErr"
          />
          <video
            v-else-if="previewType === 3"
            controls
            :src="previewUrl"
            @error="previewErr"
          ></video>
          <audio
            v-else-if="previewType === 4"
            controls
            :src="previewUrl"
            @error="previewErr"
            style="width: 100%"
          ></audio>
          <pre v-else-if="previewType === 2 || previewType === 8">{{ previewContent }}</pre>
          <div v-else class="preview-unsupported">
            <p>{{ $t('preview_err') }}</p>
            <a :href="previewUrl" target="_blank" class="preview-download-link">
              {{ $t('download') }}
            </a>
          </div>
        </div>
      </div>
    </div>

    <!-- QR Modal -->
    <div class="modal-overlay" :class="{ show: showQRModal }" @click.self="closeQR">
      <div class="modal" style="max-width: 360px">
        <div class="modal-header">
          <div class="modal-title">{{ $t('qr_title') }}</div>
          <button class="modal-close" @click="closeQR">&times;</button>
        </div>
        <div class="modal-body" style="text-align: center; padding: 24px">
          <img
            :src="qrImageUrl"
            alt="QR"
            style="max-width: 240px; margin: 0 auto 16px; display: block"
          />
          <div style="font-size: 12px; color: #888; word-break: break-all">{{ qrLink }}</div>
        </div>
      </div>
    </div>

    <!-- Settings Modal -->
    <div
      class="modal-overlay"
      :class="{ show: showSettingsModal }"
      @click.self="closeSettings"
    >
      <div class="modal settings-modal">
        <div class="modal-header">
          <div class="modal-title">{{ $t('settings') }}</div>
          <button class="modal-close" @click="closeSettings">&times;</button>
        </div>
        <div class="modal-body">
          <!-- Language Section -->
          <div class="settings-section">
            <div class="settings-section-title">{{ $t('st_lang_title') }}</div>
            <div class="settings-row">
              <span class="settings-label">{{ $t('st_lang_label') }}</span>
              <div class="lang-switch" style="width: 120px">
                <button
                  class="lang-btn"
                  :class="{ active: language === 'en' }"
                  @click="setLang('en')"
                >
                  EN
                </button>
                <button
                  class="lang-btn"
                  :class="{ active: language === 'zh' }"
                  @click="setLang('zh')"
                >
                  中文
                </button>
              </div>
            </div>
          </div>

          <!-- Services Section -->
          <div class="settings-section">
            <div class="settings-section-title">{{ $t('st_services_title') }}</div>
            
            <!-- FTP -->
            <div class="settings-row">
              <span class="settings-label">{{ $t('st_enable_ftp') }}</span>
              <button
                class="settings-toggle"
                :class="{ on: services.enableFTP }"
                @click="toggleService('ftp')"
              ></button>
            </div>
            <div v-show="services.enableFTP" class="settings-row">
              <span class="settings-label">{{ $t('st_ftp_port') }}</span>
              <input
                type="text"
                class="settings-input"
                v-model="services.ftpPort"
                placeholder="2121"
              />
            </div>

            <!-- WebDAV -->
            <div class="settings-row">
              <span class="settings-label">{{ $t('st_enable_webdav') }}</span>
              <button
                class="settings-toggle"
                :class="{ on: services.enableWebDAV }"
                @click="toggleService('webdav')"
              ></button>
            </div>
            <div v-show="services.enableWebDAV" class="settings-row">
              <span class="settings-label">{{ $t('st_webdav_port') }}</span>
              <input
                type="text"
                class="settings-input"
                v-model="services.webDAVPort"
                placeholder="8080"
              />
            </div>

            <!-- NFS -->
            <div class="settings-row">
              <span class="settings-label">{{ $t('st_enable_nfs') }}</span>
              <button
                class="settings-toggle"
                :class="{ on: services.enableNFS }"
                @click="toggleService('nfs')"
              ></button>
            </div>
            <div v-show="services.enableNFS" class="settings-row">
              <span class="settings-label">{{ $t('st_nfs_port') }}</span>
              <input
                type="text"
                class="settings-input"
                v-model="services.nfsPort"
                placeholder="2049"
              />
            </div>

            <!-- MCP -->
            <div class="settings-row">
              <span class="settings-label">{{ $t('st_enable_mcp') }}</span>
              <button
                class="settings-toggle"
                :class="{ on: services.enableMCP }"
                @click="toggleService('mcp')"
              ></button>
            </div>

            <div class="settings-row">
              <span></span>
              <button class="settings-save" @click="saveServices">{{ $t('save') }}</button>
            </div>
            <div
              class="settings-status"
              :class="{ success: servicesStatus === 'success', error: servicesStatus === 'error' }"
            >
              {{ servicesMessage }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <div class="toast" :class="{ show: toastVisible }">{{ toastMessage }}</div>

    <input
      type="file"
      ref="fileInput"
      multiple
      style="display: none"
      @change="handleFileChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppStore, type FileInfo, type TreeNode, type ServicesConfig } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import TreeNodeComponent from '@/components/TreeNode.vue'

const router = useRouter()
const route = useRoute()
const store = useAppStore()
const authStore = useAuthStore()

// State
const treeData = ref<TreeNode | null>(null)
const showPreviewModal = ref(false)
const showQRModal = ref(false)
const showSettingsModal = ref(false)
const previewTitle = ref('')
const previewUrl = ref('')
const previewType = ref(0)
const previewContent = ref('')
const previewLoading = ref(false)
const previewError = ref(false)
const qrImageUrl = ref('')
const qrLink = ref('')
const toastVisible = ref(false)
const toastMessage = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

// Settings state
const servicesStatus = ref('')
const servicesMessage = ref('')

// Computed
const language = computed(() => store.language)
const currentPath = computed(() => store.currentPath)
const files = computed(() => store.files)
const searchQuery = computed({
  get: () => store.searchQuery,
  set: (val) => (store.searchQuery = val)
})
const filteredFiles = computed(() => store.filteredFiles)
const services = computed(() => store.services)

// 检查是否可以上传（如果启用了HTTP认证，则需要登录）
const canUpload = computed(() => {
  // 如果HTTP认证未启用，允许上传
  if (!store.httpAuthOn) {
    return true
  }
  // 如果启用了HTTP认证，需要登录
  return authStore.isAuthenticated
})

const pageTitle = computed(() => {
  if (currentPath.value === '/') return $t('all_files')
  const parts = currentPath.value.split('/').filter(Boolean)
  return parts[parts.length - 1] || $t('all_files')
})

const breadcrumbParts = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  let acc = ''
  return parts.map((part) => {
    acc += '/' + part
    return { name: part, path: acc }
  })
})

// Methods
function $t(key: string): string {
  // This will be replaced by vue-i18n
  const messages: Record<string, Record<string, string>> = {
    en: {
      all_files: 'All Files',
      name: 'Name',
      size: 'Size',
      modified: 'Modified',
      actions: 'Actions',
      search: 'Search files...',
      upload: 'Upload',
      copy_link: 'Copy Link',
      download: 'Download',
      preview: 'Preview',
      empty: 'This folder is empty',
      root: 'Root',
      loading: 'Loading...',
      preview_err: 'Preview not available',
      copied: 'Link copied!',
      explorer: 'EXPLORER',
      qr_download: 'QR Download',
      qr_title: 'Scan to Download',
      zip_download: 'ZIP Download',
      settings: 'Settings',
      save_success: 'Saved!',
      save_error: 'Save failed',
      st_lang_title: 'Language',
      st_lang_label: 'Interface Language',
      st_http_auth_title: 'HTTP Authentication',
      st_enable: 'Enable',
      st_username: 'Username',
      st_password: 'Password',
      st_username_placeholder: 'Enter username',
      st_password_placeholder: 'Enter password',
      st_generate: 'Generate',
      st_services_title: 'Services',
      st_enable_ftp: 'Enable FTP',
      st_ftp_port: 'FTP Port',
      st_enable_webdav: 'Enable WebDAV',
      st_webdav_port: 'WebDAV Port',
      st_enable_nfs: 'Enable NFS',
      st_nfs_port: 'NFS Port',
      st_enable_mcp: 'Enable MCP',
      st_mcp_key: 'MCP Key',
      st_mcp_key_placeholder: 'Auto-generated or enter custom key',
      save: 'Save'
    },
    zh: {
      all_files: '所有文件',
      name: '文件名',
      size: '大小',
      modified: '修改时间',
      actions: '操作',
      search: '搜索文件...',
      upload: '上传',
      copy_link: '复制链接',
      download: '下载',
      preview: '预览',
      empty: '此文件夹为空',
      root: '根目录',
      loading: '加载中...',
      preview_err: '无法预览此文件',
      copied: '链接已复制！',
      explorer: '文件树',
      qr_download: '二维码下载',
      qr_title: '扫码下载',
      zip_download: '打包下载',
      settings: '设置',
      save_success: '已保存！',
      save_error: '保存失败',
      st_lang_title: '语言',
      st_lang_label: '界面语言',
      st_http_auth_title: 'HTTP 认证',
      st_enable: '启用',
      st_username: '用户名',
      st_password: '密码',
      st_username_placeholder: '输入用户名',
      st_password_placeholder: '输入密码',
      st_generate: '生成',
      st_services_title: '服务',
      st_enable_ftp: '启用 FTP',
      st_ftp_port: 'FTP 端口',
      st_enable_webdav: '启用 WebDAV',
      st_webdav_port: 'WebDAV 端口',
      st_enable_nfs: '启用 NFS',
      st_nfs_port: 'NFS 端口',
      st_enable_mcp: '启用 MCP',
      st_mcp_key: 'MCP 密钥',
      st_mcp_key_placeholder: '自动生成或输入自定义密钥',
      save: '保存'
    }
  }
  return messages[language.value]?.[key] || key
}

async function navigateTo(path: string) {
  await store.loadFiles(path)
  window.history.pushState({ path }, '', path === '/' ? '/' : path)
  await loadTree()
}

async function loadTree() {
  const data = await store.loadTree()
  if (data) {
    treeData.value = data
  }
}

function getFileIcon(name: string, isDir: boolean): string {
  return store.getFileIcon(name, isDir)
}

function onSearch() {
  // Computed property handles filtering
}

function doUpload() {
  fileInput.value?.click()
}

async function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    await store.uploadFiles(target.files, currentPath.value)
    target.value = ''
  }
}

function copyLink(path: string) {
  store.copyLink(path)
  showToast($t('copied'))
}

function dlFile(path: string) {
  // 移除路径开头的 /，避免双斜杠
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  // 使用 window.open 确保正确代理到后端
  window.open('/download/' + cleanPath, '_blank')
}

function showQR(path: string) {
  // 移除路径开头的 /，避免双斜杠
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  const url = window.location.origin + '/download/' + cleanPath
  qrLink.value = url
  qrImageUrl.value = '/api/qrcode?url=' + encodeURIComponent(url)
  showQRModal.value = true
}

function closeQR() {
  showQRModal.value = false
}

async function showPreview(file: FileInfo) {
  previewTitle.value = file.name
  // 移除路径开头的 /，避免双斜杠
  const cleanPath = file.path.startsWith('/') ? file.path.slice(1) : file.path
  previewUrl.value = '/download/' + cleanPath
  previewType.value = file.fileType || store.getFileType(file.name)
  previewLoading.value = true
  previewError.value = false
  previewContent.value = ''
  showPreviewModal.value = true

  // FileType: 0=Unknown, 1=Image, 2=Text, 3=Video, 4=Audio, 5=PDF, 6=Document, 7=Archive, 8=Code
  if (previewType.value === 2 || previewType.value === 8) {
    // Text/Code file
    try {
      const response = await fetch(previewUrl.value)
      if (!response.ok) throw new Error('fetch failed')
      previewContent.value = await response.text()
    } catch {
      previewErr()
    }
  }
  previewLoading.value = false
}

function previewErr() {
  previewError.value = true
}

function closePreview() {
  showPreviewModal.value = false
  previewUrl.value = ''
  previewContent.value = ''
}

function openSettings() {
  servicesStatus.value = ''
  servicesMessage.value = ''
  store.loadSettings()
  showSettingsModal.value = true
}

function closeSettings() {
  showSettingsModal.value = false
}

function toggleService(service: 'ftp' | 'webdav' | 'nfs' | 'mcp') {
  store.services[service === 'ftp' ? 'enableFTP' : service === 'webdav' ? 'enableWebDAV' : service === 'nfs' ? 'enableNFS' : 'enableMCP'] = 
    !store.services[service === 'ftp' ? 'enableFTP' : service === 'webdav' ? 'enableWebDAV' : service === 'nfs' ? 'enableNFS' : 'enableMCP']
}



async function saveServices() {
  servicesStatus.value = ''
  servicesMessage.value = ''
  const result = await store.saveServices(store.services)
  if (result.success) {
    servicesStatus.value = 'success'
    servicesMessage.value = $t('save_success')
  } else {
    servicesStatus.value = 'error'
    servicesMessage.value = result.error || $t('save_error')
  }
}

function setLang(lang: string) {
  store.setLang(lang)
}

function showToast(msg: string) {
  toastMessage.value = msg
  toastVisible.value = true
  setTimeout(() => {
    toastVisible.value = false
  }, 2500)
}

// Lifecycle
onMounted(async () => {
  // Check for initial config from server
  const configEl = document.getElementById('app-config')
  if (configEl) {
    try {
      const config = JSON.parse(configEl.textContent || '{}')
      store.initConfig(config)
    } catch (e) {
      console.error('Failed to parse config:', e)
    }
  }

  // Load initial files
  const path = route.params.path ? '/' + (Array.isArray(route.params.path) ? route.params.path.join('/') : route.params.path) : '/'
  await store.loadFiles(path)
  await loadTree()
})

// Handle browser back/forward
window.addEventListener('popstate', (e) => {
  if (e.state && e.state.path) {
    store.loadFiles(e.state.path)
  }
})
</script>

<style scoped>
.home {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: #f8f9fa;
  color: #1a1a1a;
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #e5e5e5;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-mark {
  width: 28px;
  height: 28px;
  background: #d97706;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 700;
  font-size: 14px;
  border-radius: 6px;
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  text-decoration: none;
  color: #666;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.15s;
}

.nav-item:hover {
  background: #f5f5f5;
  color: #333;
}

.nav-item.active {
  background: #d97706;
  color: #fff;
}

.sidebar-bottom {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  font-size: 13px;
  font-weight: 500;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
  transition: all 0.15s;
}

.settings-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.lang-switch {
  display: flex;
  gap: 4px;
  background: #f0f0f0;
  border-radius: 6px;
  padding: 3px;
}

.lang-btn {
  flex: 1;
  padding: 6px 0;
  text-align: center;
  border: none;
  background: transparent;
  font-size: 12px;
  font-weight: 600;
  border-radius: 4px;
  cursor: pointer;
  color: #999;
  transition: all 0.15s;
}

.lang-btn.active {
  background: #fff;
  color: #333;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.tree-section {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.tree-header {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #aaa;
  padding: 0 12px 8px;
}

.tree-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.tree-body::-webkit-scrollbar {
  width: 4px;
}

.tree-body::-webkit-scrollbar-thumb {
  background: #ddd;
  border-radius: 2px;
}

.tree-loading {
  padding: 8px 12px;
  font-size: 11px;
  color: #ccc;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.header {
  padding: 24px 32px;
  border-bottom: 1px solid #e5e5e5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  background: #fff;
  flex-shrink: 0;
}

.header-left {
  min-width: 0;
}

.breadcrumb {
  font-size: 13px;
  color: #999;
  margin-bottom: 4px;
}

.breadcrumb a {
  color: #999;
  text-decoration: none;
}

.breadcrumb a:hover {
  color: #333;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.header-right {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-shrink: 0;
}

.search {
  padding: 8px 14px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 13px;
  width: 260px;
  outline: none;
  transition: border 0.15s;
}

.search:focus {
  border-color: #999;
}

.btn-upload {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #d97706;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}

.btn-upload:hover {
  background: #b45309;
}

.content {
  flex: 1;
  overflow-y: auto;
  background: #fff;
}

.table-head {
  display: grid;
  grid-template-columns: 1fr 100px 160px 260px;
  padding: 10px 32px;
  border-bottom: 1px solid #eee;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #999;
  position: sticky;
  top: 0;
  background: #fff;
  z-index: 1;
}

.file-row {
  display: grid;
  grid-template-columns: 1fr 100px 160px 260px;
  padding: 12px 32px;
  border-bottom: 1px solid #f5f5f5;
  align-items: center;
  font-size: 13px;
  transition: background 0.1s;
}

.file-row:hover {
  background: #fafafa;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.file-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.file-name {
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-name:hover {
  color: #d97706;
}

.file-size {
  color: #888;
}

.file-date {
  color: #888;
}

.file-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
}

.act-btn {
  padding: 4px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  color: #555;
}

.act-btn:hover {
  background: #333;
  color: #fff;
  border-color: #333;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #bbb;
}

.empty p {
  font-size: 14px;
}

.modal-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  align-items: center;
  justify-content: center;
}

.modal-overlay.show {
  display: flex;
}

.modal {
  background: #fff;
  border-radius: 12px;
  width: 90vw;
  max-width: 900px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
}

.modal-title {
  font-size: 15px;
  font-weight: 600;
}

.modal-close {
  width: 32px;
  height: 32px;
  border: none;
  background: #f5f5f5;
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
}

.modal-close:hover {
  background: #e0e0e0;
}

.modal-body {
  flex: 1;
  overflow: auto;
  padding: 20px;
  min-height: 200px;
}

.modal-body pre {
  white-space: pre-wrap;
  word-break: break-all;
  font-family: 'SF Mono', 'Menlo', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #333;
}

.modal-body img {
  max-width: 100%;
  max-height: 70vh;
  display: block;
  margin: 0 auto;
}

.modal-body video,
.modal-body audio {
  width: 100%;
}

.modal-body .loading {
  text-align: center;
  padding: 40px;
  color: #bbb;
}

.modal-body .error {
  text-align: center;
  padding: 40px;
  color: #d97706;
}

.preview-unsupported {
  text-align: center;
  padding: 40px;
  color: #888;
}

.preview-download-link {
  display: inline-block;
  margin-top: 16px;
  padding: 10px 20px;
  background: #d97706;
  color: #fff;
  text-decoration: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
}

.preview-download-link:hover {
  background: #b45309;
}

.settings-modal {
  max-width: 480px;
}

.settings-section {
  margin-bottom: 24px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section-title {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}

.settings-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.settings-label {
  font-size: 13px;
  color: #555;
}

.settings-input {
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 13px;
  width: 200px;
  outline: none;
}

.settings-input:focus {
  border-color: #999;
}

.settings-input-group {
  display: flex;
  gap: 8px;
  width: 200px;
}

.settings-input-group .settings-input {
  flex: 1;
  width: auto;
}

.settings-btn-small {
  padding: 6px 12px;
  background: #f0f0f0;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}

.settings-btn-small:hover {
  background: #e0e0e0;
}

.settings-toggle {
  position: relative;
  width: 44px;
  height: 24px;
  background: #ddd;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.15s;
  border: none;
}

.settings-toggle.on {
  background: #d97706;
}

.settings-toggle::after {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  width: 18px;
  height: 18px;
  background: #fff;
  border-radius: 50%;
  transition: transform 0.15s;
}

.settings-toggle.on::after {
  transform: translateX(20px);
}

.settings-save {
  padding: 8px 16px;
  background: #d97706;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}

.settings-save:hover {
  background: #b45309;
}

.settings-status {
  font-size: 12px;
  color: #999;
  margin-top: 6px;
}

.settings-status.success {
  color: #22c55e;
}

.settings-status.error {
  color: #ef4444;
}

.toast {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: #333;
  color: #fff;
  padding: 10px 18px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  display: none;
  z-index: 2000;
}

.toast.show {
  display: block;
}

@media (max-width: 768px) {
  .sidebar {
    display: none;
  }
  .table-head {
    display: none;
  }
  .file-row {
    grid-template-columns: 1fr auto;
  }
  .file-size,
  .file-date {
    display: none;
  }
  .header {
    flex-direction: column;
    align-items: stretch;
  }
  .header-right {
    flex-direction: column;
  }
  .search {
    width: 100%;
  }
}
</style>
