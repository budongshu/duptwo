import request from './index'
import type { PageResult } from './index'

// Authorization header helper
const getAuthHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` })

// 人员响应
export interface Personnel {
  id: number
  name: string
  phone: string
  email: string
  company: string
  position: string
  workExperience: string
  entryDate: string
  projectStartDate: string
  onProjectStatus: string
  salary: string
  location: string
  remark: string
  status: string
  sort: number
  createdAt: string
  updatedAt: string
}

// 人员请求参数
export interface PersonnelReq {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string
  onProject?: string
}

// 创建人员请求
export interface CreatePersonnelReq {
  name: string
  phone?: string
  email?: string
  company?: string
  position?: string
  workExperience?: string
  entryDate?: string
  projectStartDate?: string
  onProjectStatus?: string
  salary?: string
  location?: string
  remark?: string
  status?: string
  sort?: number
}

// 更新人员请求
export interface UpdatePersonnelReq {
  id: number
  name: string
  phone?: string
  email?: string
  company?: string
  position?: string
  workExperience?: string
  entryDate?: string
  projectStartDate?: string
  onProjectStatus?: string
  salary?: string
  location?: string
  remark?: string
  status?: string
  sort?: number
}

// API 命名空间
export namespace PersonnelApi {
  // 获取人员列表
  export const list = (params?: PersonnelReq) => {
    return request.get<PageResult<Personnel[]>>('/personnels', params)
  }

  // 获取所有人员（用于下拉选择）
  export const getAll = () => {
    return request.get<Personnel[]>('/personnels/all')
  }

  // 获取人员详情
  export const getById = (id: number) => {
    return request.get<Personnel>(`/personnels/${id}`)
  }

  // 创建人员
  export const create = (data: CreatePersonnelReq) => {
    return request.post<Personnel>('/personnels', data)
  }

  // 更新人员
  export const update = (data: UpdatePersonnelReq) => {
    return request.put<Personnel>('/personnels', data)
  }

  // 删除人员
  export const del = (id: number) => {
    return request.delete<null>(`/personnels/${id}`)
  }

  // 批量删除人员
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/personnels/batch-delete', { ids })
  }

  // 导出人员Excel（使用原生 fetch 避免拦截器干扰）
  export const exportExcel = (params?: PersonnelReq) => {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== null && v !== '') {
          query.append(k, String(v))
        }
      })
    }
    const qs = query.toString()
    const url = (import.meta.env.VITE_API_URL as string || '/api') + '/personnels/export' + (qs ? '?' + qs : '')
    return fetch(url, {
      credentials: 'include',
      headers: { ...getAuthHeader() }
    }).then(res => {
      if (!res.ok) {
        throw new Error('导出失败')
      }
      return res.blob()
    })
  }
}
