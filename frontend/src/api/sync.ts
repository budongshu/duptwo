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
  apiKey?: string
  lastSyncAt: number | null  // Unix ms
  syncCount: number
  remark: string
  createdAt: number           // Unix ms
  updatedAt: number           // Unix ms
  lastHeartbeatAt: number | null  // Unix ms
  lastConnectedAt: number | null // Unix ms，最后探测成功时间
  isConnected: boolean
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
  stationCode: string
  direction: string
  directionText: string
  status: string
  statusText: string
  totalRecords: number
  successCount: number
  failCount: number
  conflictCount: number
  startedAt: number    // Unix ms
  completedAt: number  // Unix ms
  errorMsg: string
  remark: string
  createdAt: number    // Unix ms
}

export interface SyncHistoryDetail {
  id: number
  historyId: number
  serialNo: string
  projectName: string
  action: string
  actionText: string
  result: string
  errorMsg: string
  oldSerialNo: string
  newSerialNo: string
  createdAt: number // unix ms
}

export interface SyncStationSummary {
  id: number
  name: string
  code: string
  status: string
  lastSyncAt: number | null  // Unix ms
  totalSyncs: number
  totalRecords: number
  successCount: number
  failCount: number
  conflictCount: number
}

export interface SyncHistoryListReq {
  page?: number
  pageSize?: number
  stationId?: number
  direction?: string
  status?: string
  startDate?: string
  endDate?: string
  keyword?: string
}

export interface SyncHistoryDetailResp {
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
  startedAt: number    // Unix ms
  completedAt: number  // Unix ms
  errorMsg: string
  total: number
  details: SyncHistoryDetail[]
}

// 同步状态
export interface SyncStatus {
  enabled: boolean
  mode: string
  isCenter: boolean
  stationId: string
  stationName: string
  centerUrl: string
  registered: boolean
  interval: string
  heartbeatInterval: string
  batchSize: number
  filter: {
    projectNames: string[]
  } | null
  lastSyncAt: number | null  // Unix ms
  lastSerialNo: string
  syncQueueCount: number
  queueTotal: number
  queuePending: number
  queueCompleted: number
  queueFailed: number
  lastErrorAt: number | null  // Unix ms
  lastError: string
}

export interface SyncResetKeyResp {
  id: number
  apiKey: string
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

  // 重置API Key
  resetApiKey: (id: number) =>
    request.post<ResultData<SyncResetKeyResp>>(`/sync/stations/${id}/reset-key`),

  // 站点注册
  register: (data: SyncRegisterReq) =>
    request.post<ResultData<SyncRegisterResp>>('/sync/register', data),

  // 上传记录同步
  uploadRecords: (data: SyncUploadReq) =>
    request.post<ResultData<SyncUploadResp>>('/sync/upload-records', data),

  // 同步历史
  getHistory: (params?: SyncHistoryListReq) =>
    request.get<ResultData<PageResult<SyncHistory>>>('/sync/history', { params }),

  getHistoryDetails: (id: number, result?: string, page = 1, pageSize = 50) =>
    request.get<ResultData<SyncHistoryDetailResp>>(`/sync/history/${id}`, {
      params: { result, page, pageSize }
    }),

  // 站点同步汇总
  getStationSummaries: (params?: { stationId?: number; startDate?: string; endDate?: string }) =>
    request.get<ResultData<SyncStationSummary[]>>('/sync/station-summaries', { params }),

  // 同步状态
  getStatus: () =>
    request.get<ResultData<SyncStatus>>('/sync/status'),
}
