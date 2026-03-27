import request from './index'
import type { PageResult } from './index'

// 用户组信息
export interface UserGroup {
  id: number
  name: string
  code: string
  description: string
  sort: number
  createdAt: string
  updatedAt: string
}

// 创建用户组请求
export interface CreateUserGroupReq {
  name: string
  code: string
  description?: string
  sort?: number
}

// 更新用户组请求
export interface UpdateUserGroupReq {
  id: number
  name: string
  code: string
  description?: string
  sort?: number
}

// 用户组列表请求
export interface UserGroupListReq {
  page?: number
  pageSize?: number
  keyword?: string
}

export namespace UserGroupApi {
  // 获取用户组列表
  export const list = (params?: UserGroupListReq) => {
    return request.get<PageResult<UserGroup[]>>('/user-groups', params)
  }

  // 获取所有用户组
  export const getAll = () => {
    return request.get<UserGroup[]>('/user-groups/all')
  }

  // 获取用户组详情
  export const getById = (id: number) => {
    return request.get<UserGroup>(`/user-groups/${id}`)
  }

  // 创建用户组
  export const create = (data: CreateUserGroupReq) => {
    return request.post<UserGroup>('/user-groups', data)
  }

  // 更新用户组
  export const update = (data: UpdateUserGroupReq) => {
    return request.put<UserGroup>('/user-groups', data)
  }

  // 删除用户组
  export const del = (id: number) => {
    return request.delete<null>(`/user-groups/${id}`)
  }

  // 批量删除用户组
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/user-groups/batch-delete', { ids })
  }
}
