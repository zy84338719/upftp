export class ApiError extends Error {
  status?: number
  data?: any

  constructor(message: string, status?: number, data?: any) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

function getErrorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    return error.message
  }
  if (error instanceof Error) {
    return error.message
  }
  if (typeof error === 'string') {
    return error
  }
  return '未知错误'
}

// 获取 token
function getAuthToken(): string | null {
  return localStorage.getItem('auth_token') || sessionStorage.getItem('auth_token')
}

async function handleResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type')
  
  if (!response.ok) {
    let errorData: any = null
    let errorMessage = `请求失败: ${response.status} ${response.statusText}`
    
    if (contentType?.includes('application/json')) {
      try {
        errorData = await response.json()
        if (errorData.message) {
          errorMessage = errorData.message
        } else if (errorData.error) {
          errorMessage = errorData.error
        }
      } catch {
      }
    } else {
      try {
        errorMessage = await response.text()
      } catch {
      }
    }
    
    throw new ApiError(errorMessage, response.status, errorData)
  }

  if (contentType?.includes('application/json')) {
    return await response.json() as T
  }
  
  return undefined as T
}

export interface RequestOptions {
  headers?: Record<string, string>
  timeout?: number
}

export const apiClient = {
  async get<T>(url: string, options?: RequestOptions): Promise<T> {
    const controller = new AbortController()
    const timeoutId = options?.timeout 
      ? setTimeout(() => controller.abort(), options.timeout)
      : undefined

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options?.headers,
    }

    // 添加 Authorization header
    const token = getAuthToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    try {
      const response = await fetch(url, {
        method: 'GET',
        headers,
        signal: controller.signal,
      })
      return await handleResponse<T>(response)
    } finally {
      if (timeoutId) clearTimeout(timeoutId)
    }
  },

  async post<T>(url: string, data?: any, options?: RequestOptions): Promise<T> {
    const controller = new AbortController()
    const timeoutId = options?.timeout 
      ? setTimeout(() => controller.abort(), options.timeout)
      : undefined

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options?.headers,
    }

    // 添加 Authorization header
    const token = getAuthToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers,
        body: data ? JSON.stringify(data) : undefined,
        signal: controller.signal,
      })
      return await handleResponse<T>(response)
    } finally {
      if (timeoutId) clearTimeout(timeoutId)
    }
  },

  async postFormData<T>(url: string, formData: FormData, options?: RequestOptions): Promise<T> {
    const controller = new AbortController()
    const timeoutId = options?.timeout 
      ? setTimeout(() => controller.abort(), options.timeout)
      : undefined

    const headers: Record<string, string> = {
      ...options?.headers,
    }

    // 添加 Authorization header
    const token = getAuthToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers,
        body: formData,
        signal: controller.signal,
      })
      return await handleResponse<T>(response)
    } finally {
      if (timeoutId) clearTimeout(timeoutId)
    }
  },

  getErrorMessage,
  ApiError,
}
