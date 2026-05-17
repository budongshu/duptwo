import request from './index'
import type { PageResult } from './index'

// 上传记录请求参数
export interface UploadRecordReq {
  page?: number
  pageSize?: number
  diskLabel?: string
  projectName?: string
  status?: string
  uploader?: string
  startDate?: string
  endDate?: string
  keyword?: string
  serialNo?: string
}

// 上传记录响应
export interface UploadRecord {
  id: number
  serialNo: string
  diskLabel: string
  projectName: string
  destPath: string
  fileSize: number
  fileSizeStr: string
  uploader: string
  status: string
  statusText: string
  remark: string
  data: Record<string, any>
  createdAt: string
  updatedAt: string
}

// 创建上传记录请求
export interface CreateUploadRecordReq {
  diskLabel: string
  projectName?: string
  destPath: string
  fileSize: number
  uploader: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  remark?: string
  data?: Record<string, any>
  createdAt?: string  // 可选，格式：2006-01-02 或 2006-01-02T15:04:05Z
}

// 更新上传记录请求
export interface UpdateUploadRecordReq {
  id: number
  status: 'pending' | 'processing' | 'completed' | 'failed'
  remark?: string
  fileSize?: number  // 上传完成后补充文件大小（字节）
}

// 统计数据
export interface UploadRecordStatistics {
  todayCount: number
  todaySize: number
  todaySizeStr: string
  todayByStatus: StatusCount[]
  weekCount: number
  weekSize: number
  weekSizeStr: string
  monthCount: number
  monthSize: number
  monthSizeStr: string
  totalCount: number
  totalSize: number
  totalSizeStr: string
  trend: DailyTrend[]
  byStatus: StatusCount[]
  byDiskLabel: DiskLabelCount[]
  byProject: ProjectCount[]
}

export interface DailyTrend {
  date: string
  count: number
  totalSize: number
}

export interface StatusCount {
  status: string
  count: number
}

export interface DiskLabelCount {
  diskLabel: string
  count: number
}

export interface ProjectCount {
  projectName: string
  count: number
  totalSize: number
}

// 磁盘标签状态
export interface DiskLabelStatus {
  diskLabel: string
  count: number
  totalSize: number
  status: 'completed' | 'failed' | 'mixed' | 'pending'
}

// ============ 批量导入相关 ============

// 导入模板字段定义
export interface ImportTemplateField {
  field: string      // Excel列名（中文）
  code: string        // 字段代码
  required: boolean   // 是否必填
  type: string        // text/number/select/date
  options?: string    // 下拉选项
  maxLength?: number  // 最大长度
  example?: string    // 填写示例
}

// 导入模板信息
export interface ImportTemplateResp {
  fields: ImportTemplateField[]
  sheetName: string
  title: string
}

// 导入失败行
export interface ImportFailRow {
  row: number     // Excel行号（从2开始）
  data: string    // 该行原始数据
  reason: string  // 失败原因
}

// 导入结果
export interface ImportResultResp {
  total: number
  success: number
  failed: number
  failRows: ImportFailRow[]
}

// API 命名空间
export namespace UploadRecordApi {
  // 获取上传记录列表
  export const list = (params: UploadRecordReq) => {
    return request.get<PageResult<UploadRecord[]>>('/upload-records', params)
  }

  // 创建上传记录
  export const create = (data: CreateUploadRecordReq) => {
    return request.post<UploadRecord>('/upload-records', data)
  }

  // 获取上传记录详情
  export const getById = (id: number) => {
    return request.get<UploadRecord>(`/upload-records/${id}`)
  }

  // 更新上传记录
  export const update = (data: UpdateUploadRecordReq) => {
    return request.put<UploadRecord>('/upload-records', data)
  }

  // 删除上传记录
  export const del = (id: number) => {
    return request.delete<null>(`/upload-records/${id}`)
  }

  // 批量删除上传记录
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/upload-records/batch-delete', { ids })
  }

  // 批量更新上传记录状态
  export const batchUpdateStatus = (ids: number[], status: 'pending' | 'processing' | 'completed' | 'failed') => {
    return request.post<null>('/upload-records/batch-update-status', { ids, status })
  }

  // 获取统计数据
  export const statistics = (params?: { startDate?: string; endDate?: string; projectName?: string; diskLabel?: string; status?: string; uploader?: string }) => {
    return request.get<UploadRecordStatistics>('/upload-records/statistics', params)
  }

  // 获取磁盘标签状态列表
  export const diskLabels = (params?: { projectName?: string; diskLabel?: string; startDate?: string; endDate?: string }) => {
    return request.get<DiskLabelStatus[]>('/upload-records/disk-labels', params)
  }

  // 获取最近上传记录
  export const recent = (params?: { limit?: number; projectName?: string; diskLabel?: string; status?: string; uploader?: string }) => {
    return request.get<UploadRecord[]>('/upload-records/recent', params)
  }

  // 获取上传者列表
  export const uploaderList = () => {
    return request.get<string[]>('/upload-records/uploaders')
  }

  // 导出上传记录为 Excel（使用原生 fetch 避免拦截器干扰）
const getApiBase = () => (import.meta.env.VITE_API_URL as string || '/api').replace(/\/$/, '')
const getAuthHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` })

  export const exportExcel = (params?: UploadRecordReq) => {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== null && v !== '') {
          query.append(k, String(v))
        }
      })
    }
    const qs = query.toString()
    // 直接使用 /api 前缀，不要重复
    const url = '/api/upload-records/export' + (qs ? '?' + qs : '')
    return fetch(url, {
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        'Accept': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      }
    }).then(res => {
      if (!res.ok) {
        return res.json().then(err => Promise.reject(err)).catch(() => Promise.reject(new Error('导出失败')))
      }
      return res.blob()
    })
  }

  // 获取导入模板字段定义
  export const getImportTemplate = () => {
    return request.get<ImportTemplateResp>('/upload-records/template')
  }

  // 下载导入模板 Excel
  export const downloadTemplate = () => {
    const url = getApiBase() + '/upload-records/template'
    return fetch(url, {
      credentials: 'include',
      headers: { ...getAuthHeader() }
    }).then(res => {
      if (!res.ok) throw new Error('获取模板信息失败')
      return res.json()
    }).then(data => {
      generateTemplateExcel(data.data || data as ImportTemplateResp)
    })
  }

  // 导入上传记录 Excel（支持进度回调）
  export const importRecords = (file: File, onProgress?: (pct: number) => void): Promise<{ code: number; message: string; data: ImportResultResp }> => {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()
      const formData = new FormData()
      formData.append('file', file)
      const url = getApiBase() + '/upload-records/import'
      xhr.open('POST', url, true)
      xhr.withCredentials = true
      const authHeader = getAuthHeader()
      for (const key in authHeader) {
        xhr.setRequestHeader(key, (authHeader as any)[key])
      }
      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable && onProgress) {
          onProgress(Math.round((e.loaded / e.total) * 100))
        }
      }
      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(JSON.parse(xhr.responseText))
        } else {
          reject(new Error('导入请求失败'))
        }
      }
      xhr.onerror = () => reject(new Error('导入请求失败'))
      xhr.send(formData)
    })
  }

  // 预览上传文件行数
  export const previewImport = (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    const url = getApiBase() + '/upload-records/preview'
    return fetch(url, {
      method: 'POST',
      credentials: 'include',
      headers: { ...getAuthHeader() },
      body: formData
    }).then(res => {
      if (!res.ok) throw new Error('预览请求失败')
      return res.json()
    }) as Promise<{ code: number; message: string; data: { totalRows: number; dataRows: number; sheetName: string } }>
  }
}

// Excel列名转换（支持 > 26 列，如 0→A, 25→Z, 26→AA, 27→AB）
function colName(n: number): string {
  let s = ''
  n++
  while (n > 0) {
    n--
    s = String.fromCharCode(65 + (n % 26)) + s
    n = Math.floor(n / 26)
  }
  return s
}

// 前端生成模板Excel（使用 SheetJS）
function generateTemplateExcel(template: ImportTemplateResp) {
  const script = document.createElement('script')
  script.src = 'https://cdn.sheetjs.com/xlsx-0.20.1/package/dist/xlsx.full.min.js'
  script.onload = () => {
    const XLSX = (window as any).XLSX

    // ---------- 构建工作表数据（纯数据，不用样式对象混入 aoa） ----------
    // 第1行：标题
    const titleRow: any[] = [{ v: template.title || '上传记录批量导入模板', t: 's', w: '标题' }]

    // 第2行：必填说明
    const reqFields = template.fields.filter(f => f.required).map(f => f.field).join('、')
    const hintRow: any[] = [{ v: `* 必填字段：${reqFields || '无'}`, t: 's' }]

    // 第3行：表头
    const headerRow: any[] = template.fields.map(f => ({
      v: f.required ? `${f.field} *` : f.field,
      t: 's'
    }))

    // 第4行：示例数据（一行完整的 demo）
    const demoRow: any[] = template.fields.map(f => ({
      v: f.example || '',
      t: 's'
    }))

    // 合并所有行
    const allRows = [titleRow, hintRow, headerRow, demoRow]
    const ws = XLSX.utils.aoa_to_sheet(allRows)

    // ---------- 设置行列样式 ----------
    // 第1行：标题（合并到最后一列）
    const lastCol = colName(template.fields.length - 1)
    ws['!merges'] = [{ s: { r: 0, c: 0 }, e: { r: 0, c: template.fields.length - 1 } }]
    ws['!rows'] = [
      { hpt: 36 },  // 标题行
      { hpt: 22 },  // 说明行
      { hpt: 30 },  // 表头行
      { hpt: 28 },  // 示例数据行
    ]

    // 设置列宽
    ws['!cols'] = template.fields.map(f => ({
      wch: Math.max(f.field.length, (f.example?.length || 0) + 2, 12)
    }))

    // 手动给第3行（表头）和第4行（示例）加样式
    for (let c = 0; c < template.fields.length; c++) {
      const addr3 = XLSX.utils.encode_cell({ r: 2, c })
      const addr4 = XLSX.utils.encode_cell({ r: 3, c })
      // 表头样式：深蓝底白字
      ;(ws[addr3] as any).s = {
        fill: { fgColor: { rgb: '1F4E79' }, patternType: 'solid' },
        font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 },
        alignment: { horizontal: 'center', vertical: 'center' },
        border: {
          top: { style: 'thin', color: { rgb: '2E75B6' } },
          bottom: { style: 'thin', color: { rgb: '2E75B6' } },
          left: { style: 'thin', color: { rgb: 'BDD7EE' } },
          right: { style: 'thin', color: { rgb: 'BDD7EE' } }
        }
      }
      // 示例行样式：浅蓝底
      ;(ws[addr4] as any).s = {
        fill: { fgColor: { rgb: 'DEEAF1' }, patternType: 'solid' },
        font: { color: { rgb: '2E75B6' }, sz: 11 },
        alignment: { horizontal: 'left', vertical: 'center' },
        border: {
          top: { style: 'thin', color: { rgb: 'BDD7EE' } },
          bottom: { style: 'thin', color: { rgb: 'BDD7EE' } },
          left: { style: 'thin', color: { rgb: 'BDD7EE' } },
          right: { style: 'thin', color: { rgb: 'BDD7EE' } }
        }
      }
    }

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, template.sheetName || '上传记录导入')

    const filename = `上传记录导入模板_${new Date().toISOString().slice(0, 10).replace(/-/g, '')}.xlsx`
    XLSX.writeFile(wb, filename)
    document.body.removeChild(script)
  }
  script.onerror = () => {
    document.body.removeChild(script)
    throw new Error('在线模板生成失败，请联系管理员')
  }
  document.body.appendChild(script)
}
