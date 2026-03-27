import axios, { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'

export interface ResultData<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResult<T = any> {
  total: number
  items: T
}

// API配置
const config: AxiosRequestConfig = {
  baseURL: import.meta.env.VITE_API_URL as string || '/api',
  timeout: 20000,
  withCredentials: true
}

class RequestHttp {
  service: AxiosInstance

  constructor(cfg: AxiosRequestConfig) {
    this.service = axios.create(cfg)

    // 请求拦截器：添加 JWT Token
    this.service.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.service.interceptors.response.use(
      (response: AxiosResponse<ResultData>) => {
        // 如果是 blob 类型，直接返回
        if (response.config.responseType === 'blob') {
          return response.data
        }
        const { data } = response
        if (data.code !== 200) {
          return Promise.reject(data)
        }
        return data
      },
      (error: any) => {
        // 如果是 401 未授权，跳转到登录页
        if (error.response?.status === 401) {
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          return Promise.reject(error)
        }
        // 处理 400/500 等错误，提取后端返回的错误消息
        if (error.response?.data) {
          return Promise.reject(error.response.data)
        }
        return Promise.reject(error)
      }
    )
  }

  get<T>(url: string, params?: object, config = {}): Promise<ResultData<T>> {
    return this.service.get(url, { params, ...config })
  }

  post<T>(url: string, params?: object, timeout?: number): Promise<ResultData<T>> {
    return this.service.post(url, params, { timeout })
  }

  put<T>(url: string, params?: object, config = {}): Promise<ResultData<T>> {
    return this.service.put(url, params, config)
  }

  delete<T>(url: string, params?: any, config = {}): Promise<ResultData<T>> {
    return this.service.delete(url, { params, ...config })
  }

  // 导出文件（返回原始响应，不经过响应拦截器）
  exportFile(url: string, params?: object): Promise<any> {
    return this.service.get(url, { params, responseType: 'blob' }).then((res: any) => res.data)
  }
}

export default new RequestHttp(config)
