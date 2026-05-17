import request from './index'
import type { PageResult } from './index'

// 角色信息
export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string[]
  sort: number
  createdAt: string
  updatedAt: string
}

// 创建角色请求
export interface CreateRoleReq {
  name: string
  code: string
  description?: string
  permissions?: string[]
  sort?: number
}

// 更新角色请求
export interface UpdateRoleReq {
  id: number
  name: string
  code: string
  description?: string
  permissions?: string[]
  sort?: number
}

// 角色列表请求
export interface RoleListReq {
  page?: number
  pageSize?: number
  keyword?: string
}

export namespace RoleApi {
  // 获取角色列表
  export const list = (params?: RoleListReq) => {
    return request.get<PageResult<Role[]>>('/roles', params)
  }

  // 获取所有角色
  export const getAll = () => {
    return request.get<Role[]>('/roles/all')
  }

  // 获取角色详情
  export const getById = (id: number) => {
    return request.get<Role>(`/roles/${id}`)
  }

  // 创建角色
  export const create = (data: CreateRoleReq) => {
    return request.post<Role>('/roles', data)
  }

  // 更新角色
  export const update = (data: UpdateRoleReq) => {
    return request.put<Role>('/roles', data)
  }

  // 删除角色
  export const del = (id: number) => {
    return request.delete<null>(`/roles/${id}`)
  }

  // 批量删除角色
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/roles/batch-delete', { ids })
  }
}
