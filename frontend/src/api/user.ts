import request from './index'
import type { PageResult, ResultData } from './index'

// 用户信息
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: string
  statusText: string
  roleId: number
  roleName: string
  groupId: number
  groupName: string
  mfaEnabled: boolean
  source: string   // LOCAL / AD
  lastLoginAt?: string
  lastLoginIP?: string
  createdAt: string
  updatedAt: string
}

// 创建用户请求
export interface CreateUserReq {
  username: string
  password: string
  nickname?: string
  email?: string
  phone?: string
  roleId?: number
  groupId?: number
  mfaEnabled?: boolean
}

// 更新用户请求
export interface UpdateUserReq {
  id: number
  nickname?: string
  email?: string
  phone?: string
  roleId?: number
  groupId?: number
  status?: string
  mfaEnabled?: boolean
}

// 用户列表请求
export interface UserListReq {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string
  roleId?: number
  groupId?: number
}

// 重置密码请求
export interface ResetPasswordReq {
  userId: number
  newPassword: string
}

// 重置MFA请求
export interface ResetMFAReq {
  userId: number
}

// MFA密钥响应
export interface MFASecretResp {
  secret: string
  qrCode: string
}

// 管理员启用MFA请求
export interface AdminEnableMFAReq {
  userId: number
  code: string
}

export namespace UserApi {
  // 获取用户列表
  export const list = (params?: UserListReq) => {
    return request.get<PageResult<User[]>>('/users', params)
  }

  // 获取所有用户
  export const getAll = () => {
    return request.get<User[]>('/users/all')
  }

  // 获取用户详情
  export const getById = (id: number) => {
    return request.get<User>(`/users/${id}`)
  }

  // 创建用户
  export const create = (data: CreateUserReq) => {
    return request.post<User>('/users', data)
  }

  // 更新用户
  export const update = (data: UpdateUserReq) => {
    return request.put<User>('/users', data)
  }

  // 获取当前用户信息（个人设置）
  export const getCurrentUser = () => {
    return request.get<User>('/users/current')
  }

  // 更新当前用户资料（个人设置）
  export const updateProfile = (data: { nickname?: string; email?: string; phone?: string }) => {
    return request.put<User>('/users/profile', data)
  }

  // 删除用户
  export const del = (id: number) => {
    return request.delete<null>(`/users/${id}`)
  }

  // 批量删除用户
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/users/batch-delete', { ids })
  }

  // 重置用户密码
  export const resetPassword = (data: ResetPasswordReq) => {
    return request.post<null>('/users/reset-password', data)
  }

  // 重置用户MFA
  export const resetMFA = (data: ResetMFAReq) => {
    return request.post<null>('/users/reset-mfa', data)
  }

  // 生成MFA密钥
  export const generateMFASecret = (userId: number) => {
    return request.get<MFASecretResp>('/users/generate-mfa-secret', { id: userId })
  }

  // 管理员启用用户MFA
  export const adminEnableMFA = (data: AdminEnableMFAReq) => {
    return request.post<null>('/users/admin-enable-mfa', data)
  }

  // 导出用户列表
  export const exportExcel = (params?: { keyword?: string; status?: string }) => {
    return request.get<Blob>('/users/export', params, { responseType: 'blob' })
  }

  // 获取导入模板
  export const getImportTemplate = () => {
    return request.get<{ fields: any[]; sheetName: string; title: string }>('/users/template')
  }

  // 导入用户
  export const importUsers = (file: File, onProgress?: (pct: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<{ total: number; success: number; failed: number; failRows: any[] }>('/users/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => { if (e.total) onProgress?.(Math.round((e.loaded * 100) / e.total)) },
    })
  }
}
