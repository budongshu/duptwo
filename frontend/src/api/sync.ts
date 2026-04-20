import request from './index'
import type { ResultData, PageResult } from './index'

// 站点相关
export interface SyncStation {
  id: number
  name: string
  code: string
  url: string
  status: string
  statusText: string
  description: string
  isCenter: boolean
  lastSyncAt: string | null
  syncCount: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface SyncStationCreateReq {
  name: string
  code: string
  url: string
  description?: string
  isCenter?: boolean
}

export interface SyncStationUpdateReq {
  id: number
  name: string
  code: string
  url: string
  status?: string
  description?: string
  isCenter?: boolean
}

export interface SyncStationListReq {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string
}

// 站点注册
export interface SyncRegisterReq {
  stationCode: string
  stationName: string
  url: string
  password: string
}

export interface SyncRegisterResp {
  stationId: number
  apiKey: string
  message: string
}

// 上传记录同步
export interface SyncRecordItem {
  serialNo: string
  projectName: string
  diskLabel?: string
  destPath?: string
  fileSize?: number
  uploader?: string
  status?: string
  remark?: string
  data?: Record<string, any>
}

export interface SyncUploadReq {
  records: SyncRecordItem[]
}

export interface SyncDetailItem {
  serialNo: string
  projectName: string
  action: string
  result: string
  errorMsg?: string
  newSerialNo?: string
  oldSerialNo?: string
}

export interface SyncUploadResp {
  totalRecords: number
  successCount: number
  failCount: number
  conflictCount: number
  details: SyncDetailItem[]
  historyId: number
}

// 同步历史
export interface SyncHistory {
  id: number
  stationId: number
  stationName: string
  direction: string
  directionText: string
  status: string
  statusText: string
  totalRecords: number
  successCount: number
  failCount: number
  conflictCount: number
  startedAt: string | null
  completedAt: string | null
  errorMsg: string
  remark: string
  createdAt: string
}

export interface SyncHistoryDetail {
  id: number
  serialNo: string
  projectName: string
  action: string
  actionText: string
  result: string
  errorMsg?: string
  oldSerialNo?: string
  newSerialNo?: string
  createdAt: string
}

export interface SyncHistoryListReq {
  page?: number
  pageSize?: number
  stationId?: number
  direction?: string
  status?: string
  startDate?: string
  endDate?: string
}

// 同步状态
export interface SyncStatus {
  enabled: boolean
  mode: string
  isCenter: boolean
  stationId: string
  stationName: string
  centerUrl: string
  lastSyncAt: string | null
  syncQueueCount: number
}

// API 函数
export const syncApi = {
  // 站点管理
  createStation: (data: SyncStationCreateReq) =>
    request.post<ResultData<SyncStation>>('/sync/stations', data),

  listStations: (params?: SyncStationListReq) =>
    request.get<ResultData<PageResult<SyncStation>>>('/sync/stations', { params }),

  getAllStations: () =>
    request.get<ResultData<SyncStation[]>>('/sync/stations/all'),

  getStation: (id: number) =>
    request.get<ResultData<SyncStation>>(`/sync/stations/${id}`),

  updateStation: (data: SyncStationUpdateReq) =>
    request.put<ResultData<SyncStation>>('/sync/stations', data),

  deleteStation: (id: number) =>
    request.delete<ResultData<null>>(`/sync/stations/${id}`),

  // 站点注册
  register: (data: SyncRegisterReq) =>
    request.post<ResultData<SyncRegisterResp>>('/sync/register', data),

  // 上传记录同步
  uploadRecords: (data: SyncUploadReq) =>
    request.post<ResultData<SyncUploadResp>>('/sync/upload-records', data),

  // 同步历史
  getHistory: (params?: SyncHistoryListReq) =>
    request.get<ResultData<PageResult<SyncHistory>>>('/sync/history', { params }),

  getHistoryDetails: (id: number) =>
    request.get<ResultData<{ details: SyncHistoryDetail[] }>>(`/sync/history/${id}`),

  // 同步状态
  getStatus: () =>
    request.get<ResultData<SyncStatus>>('/sync/status'),
}
