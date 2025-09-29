export interface ApiResponse<T = any> {
  data: T
  message?: string
  success: boolean
}

export interface ApiError {
  message: string
  code?: string
  details?: any
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

export interface QueryParams {
  [key: string]: string | number | boolean | undefined
}

export interface ApiClientConfig {
  baseURL: string
  timeout: number
  headers?: { [key: string]: string }
}