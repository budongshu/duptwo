import request from './index'
import type { PageResult } from './index'

// 驻场点信息
export interface OnSiteStation {
  location: string  // 场地名称
  person: string    // 驻场人员
  phone: string     // 联系方式
}

// 项目请求参数
export interface ProjectReq {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string
  stage?: string
}

// 项目响应
export interface Project {
  id: number
  name: string
  code: string
  description: string
  status: string
  stage: string
  sort: number
  recordCount: number
  totalDataSize: number
  createdAt: string
  updatedAt: string
  projectPerson: string   // 项目人员
  opsPerson: string       // 运维人员
  developerPerson: string // 开发人员
  testerPerson: string   // 测试人员
  businessPerson: string // 商务人员
  compliancePerson: string // 合规专员
  opsStaffPerson: string  // 运营人员
  solution: string       // 解决方案描述
  solutionPerson: string // 解决方案人员
  companyAddr: string    // 公司地点
  projectPeriod: string   // 项目周期
  onsiteStations: OnSiteStation[]  // 驻场点列表
}

// 简单项目响应（用于下拉选择）
export interface ProjectSimple {
  id: number
  name: string
  code: string
}

// 创建项目请求
export interface CreateProjectReq {
  name: string
  code: string
  description?: string
  status?: string
  stage?: string
  sort?: number
  projectPerson?: string   // 项目人员
  opsPerson?: string       // 运维人员
  developerPerson?: string // 开发人员
  testerPerson?: string   // 测试人员
  businessPerson?: string // 商务人员
  compliancePerson?: string // 合规专员
  opsStaffPerson?: string  // 运营人员
  solution?: string       // 解决方案描述
  solutionPerson?: string // 解决方案人员
  companyAddr?: string     // 公司地点
  projectPeriod?: string   // 项目周期
  onsiteStations?: OnSiteStation[]  // 驻场点列表
}

// 更新项目请求
export interface UpdateProjectReq {
  id: number
  name: string
  code: string
  description?: string
  status?: string
  stage?: string
  sort?: number
  projectPerson?: string   // 项目人员
  opsPerson?: string       // 运维人员
  developerPerson?: string // 开发人员
  testerPerson?: string   // 测试人员
  businessPerson?: string // 商务人员
  compliancePerson?: string // 合规专员
  opsStaffPerson?: string  // 运营人员
  solution?: string       // 解决方案描述
  solutionPerson?: string // 解决方案人员
  companyAddr?: string    // 公司地点
  projectPeriod?: string  // 项目周期
  onsiteStations?: OnSiteStation[]  // 驻场点列表
}

// API 命名空间
export namespace ProjectApi {
  // 获取项目列表
  export const list = (params?: ProjectReq) => {
    return request.get<PageResult<Project[]>>('/projects', params)
  }

  // 获取所有项目（用于下拉选择）
  export const getSimpleList = () => {
    return request.get<ProjectSimple[]>('/projects/simple')
  }

  // 获取所有项目（看板视图）
  export const getKanbanList = () => {
    return request.get<Project[]>('/projects/kanban')
  }

  // 获取项目详情
  export const getById = (id: number) => {
    return request.get<Project>(`/projects/${id}`)
  }

  // 创建项目
  export const create = (data: CreateProjectReq) => {
    return request.post<Project>('/projects', data)
  }

  // 更新项目
  export const update = (data: UpdateProjectReq) => {
    return request.put<Project>('/projects', data)
  }

  // 删除项目
  export const del = (id: number) => {
    return request.delete<null>(`/projects/${id}`)
  }

  // 批量删除项目
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/projects/batch-delete', { ids })
  }
}
