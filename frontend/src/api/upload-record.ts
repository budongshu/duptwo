import request from './index'
import type { PageResult } from './index'

// 上传记录请求参数
export interface UploadRecordReq {
  page?: number
  pageSize?: number
  dataType?: string
  projectName?: string
  status?: string
  uploader?: string
  startDate?: string
  endDate?: string
  keyword?: string
}

// 上传记录响应
export interface UploadRecord {
  id: number
  serialNo: string
  dataType: string
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
  dataType: string
  projectName?: string
  destPath: string
  fileSize: number
  uploader: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  remark?: string
  data?: Record<string, any>
}

// 更新上传记录请求
export interface UpdateUploadRecordReq {
  id: number
  status: 'pending' | 'processing' | 'completed' | 'failed'
  remark?: string
}

// 统计数据
export interface UploadRecordStatistics {
  todayCount: number
  todaySize: number
  todaySizeStr: string
  weekCount: number
  monthCount: number
  totalCount: number
  totalSize: number
  totalSizeStr: string
  trend: DailyTrend[]
  byStatus: StatusCount[]
  byDataType: DataTypeCount[]
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

export interface DataTypeCount {
  dataType: string
  count: number
}

export interface ProjectCount {
  projectName: string
  count: number
  totalSize: number
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

  // 获取统计数据
  export const statistics = (params?: { startDate?: string; endDate?: string; projectName?: string; dataType?: string; status?: string; uploader?: string }) => {
    return request.get<UploadRecordStatistics>('/upload-records/statistics', params)
  }

  // 获取最近上传记录
  export const recent = (params?: { limit?: number; projectName?: string; dataType?: string; status?: string; uploader?: string }) => {
    return request.get<UploadRecord[]>('/upload-records/recent', params)
  }

  // 获取上传者列表
  export const uploaderList = () => {
    return request.get<string[]>('/upload-records/uploaders')
  }

  // 导出上传记录为 Excel（使用原生 fetch 避免拦截器干扰）
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
    const url = (import.meta.env.VITE_API_URL as string || '/api') + '/upload-records/export' + (qs ? '?' + qs : '')
    return fetch(url, {
      credentials: 'include'
    }).then(res => {
      if (!res.ok) {
        throw new Error('导出失败')
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
    const url = (import.meta.env.VITE_API_URL as string || '/api') + '/upload-records/template'
    return fetch(url, {
      credentials: 'include'
    }).then(res => {
      if (!res.ok) throw new Error('获取模板信息失败')
      return res.json()
    }).then(data => {
      // 使用前端生成模板（根据字段定义生成xlsx）
      generateTemplateExcel(data.data || data as ImportTemplateResp)
    })
  }

  // 导入上传记录 Excel
  export const importRecords = (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    const url = (import.meta.env.VITE_API_URL as string || '/api') + '/upload-records/import'
    return fetch(url, {
      method: 'POST',
      credentials: 'include',
      body: formData
    }).then(res => {
      if (!res.ok) throw new Error('导入请求失败')
      return res.json()
    }) as Promise<{ code: number; message: string; data: ImportResultResp }>
  }
}

// 前端生成模板Excel（使用 SheetJS）
function generateTemplateExcel(template: ImportTemplateResp) {
  // 动态加载 xlsx CDN
  const script = document.createElement('script')
  script.src = 'https://cdn.sheetjs.com/xlsx-0.20.1/package/dist/xlsx.full.min.js'
  script.onload = () => {
    const XLSX = (window as any).XLSX
    const worksheetData: any[] = []

    // 第1行：标题行
    worksheetData.push([template.title || '上传记录批量导入模板'])

    // 第2行：字段说明行
    worksheetData.push(['字段说明：* 为必填项'])

    // 第3行：表头
    const headers = template.fields.map(f => {
      const label = f.required ? `* ${f.field}` : f.field
      return {
        v: label,
        t: 's',
        l: {
          fill: { fgColor: { rgb: '005bbf' } },
          font: { bold: true, color: { rgb: 'FFFFFF' } }
        }
      }
    })
    worksheetData.push(headers as any)

    // 第4行起：示例数据
    template.fields.forEach(f => {
      worksheetData.push([f.example || ''])
    })

    const ws = XLSX.utils.aoa_to_sheet(worksheetData)

    // 设置表头行高和样式
    ws['!rows'] = [
      { hpt: 30 }, // 标题行
      { hpt: 20 }, // 说明行
      { hpt: 28 }, // 表头行
    ]

    // 设置列宽
    const colWidths = template.fields.map(f => ({
      wch: Math.max(f.field.length, f.example?.length || 0, 15) + 4
    }))
    ws['!cols'] = colWidths

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, template.sheetName || '上传记录导入')

    // 下载
    const filename = `上传记录导入模板_${new Date().toISOString().slice(0, 10).replace(/-/g, '')}.xlsx`
    XLSX.writeFile(wb, filename)
    document.body.removeChild(script)
  }
  script.onerror = () => {
    // CDN加载失败，提示用户
    alert('在线模板生成失败，请联系管理员')
    document.body.removeChild(script)
  }
  document.body.appendChild(script)
}
