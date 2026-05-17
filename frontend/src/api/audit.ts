import request from './index'
import type { PageResult } from './index'

// 操作日志
export interface OperationLog {
  id: number
  userId: number
  username: string
  menuName: string
  action: string
  actionText: string
  resourceType: string
  resourceId: number
  resourceName: string
  ipAddress: string
  detail: string
  createdAt: string
}

export interface OperationLogListReq {
  page?: number
  pageSize?: number
  keyword?: string
  userId?: number
  menuName?: string
  action?: string
  resourceType?: string
  startDate?: string
  endDate?: string
}

export namespace OperationLogApi {
  export const list = (params?: OperationLogListReq) => {
    return request.get<PageResult<OperationLog[]>>('/audit/operation-logs', params)
  }

  export const exportExcel = (params?: OperationLogListReq) => {
    return request.get('/audit/operation-logs/export', params, { responseType: 'blob' })
  }
}

// 登录日志
export interface LoginLog {
  id: number
  userId: number
  username: string
  status: string
  statusText: string
  ipAddress: string
  failReason: string
  mfaUsed: boolean
  loginMethod: string
  createdAt: string
}

export interface LoginLogListReq {
  page?: number
  pageSize?: number
  keyword?: string
  userId?: number
  status?: string
  startDate?: string
  endDate?: string
}

export namespace LoginLogApi {
  export const list = (params?: LoginLogListReq) => {
    return request.get<PageResult<LoginLog[]>>('/audit/login-logs', params)
  }

  export const exportExcel = (params?: LoginLogListReq) => {
    return request.get('/audit/login-logs/export', params, { responseType: 'blob' })
  }
}
