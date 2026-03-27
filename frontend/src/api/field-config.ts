import request from './index'
import type { PageResult } from './index'

// 字段配置请求参数
export interface FieldConfigReq {
  page?: number
  pageSize?: number
  keyword?: string
  enabled?: boolean
}

// 字段配置响应
export interface FieldConfig {
  id: number
  name: string
  code: string
  type: 'text' | 'number' | 'select' | 'date' | 'datetime'
  required: boolean
  options: string[]
  defaultValue: string
  placeholder: string
  sort: number
  enabled: boolean
  createdAt: string
  updatedAt: string
}

// 创建字段配置请求
export interface CreateFieldConfigReq {
  name: string
  code: string
  type: 'text' | 'number' | 'select' | 'date' | 'datetime'
  required?: boolean
  options?: string[]
  defaultValue?: string
  placeholder?: string
  sort?: number
  enabled?: boolean
}

// 更新字段配置请求
export interface UpdateFieldConfigReq {
  id: number
  name: string
  code: string
  type: 'text' | 'number' | 'select' | 'date' | 'datetime'
  required?: boolean
  options?: string[]
  defaultValue?: string
  placeholder?: string
  sort?: number
  enabled?: boolean
}

// API 命名空间
export namespace FieldConfigApi {
  // 获取字段配置列表
  export const list = (params?: FieldConfigReq) => {
    return request.get<PageResult<FieldConfig[]>>('/field-configs', params)
  }

  // 创建字段配置
  export const create = (data: CreateFieldConfigReq) => {
    return request.post<FieldConfig>('/field-configs', data)
  }

  // 获取字段配置详情
  export const getById = (id: number) => {
    return request.get<FieldConfig>(`/field-configs/${id}`)
  }

  // 更新字段配置
  export const update = (data: UpdateFieldConfigReq) => {
    return request.put<FieldConfig>('/field-configs', data)
  }

  // 删除字段配置
  export const del = (id: number) => {
    return request.delete<null>(`/field-configs/${id}`)
  }

  // 批量删除字段配置
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/field-configs/batch-delete', { ids })
  }

  // 获取所有启用的字段配置
  export const getAllEnabled = () => {
    return request.get<FieldConfig[]>('/field-configs/all')
  }
}
