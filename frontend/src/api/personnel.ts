import request from './index'
import type { PageResult } from './index'

// Authorization header helper
const getAuthHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` })
const getApiBase = () => (import.meta.env.VITE_API_URL as string || '/api').replace(/\/$/, '')

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
  position?: string
}

// 人员统计响应
export interface PersonnelStatistics {
  total: number
  byPosition: { position: string; count: number }[]
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

  // 获取人员统计（总数 + 职位分布）
  export const statistics = (params?: PersonnelReq) => {
    return request.get<PersonnelStatistics>('/personnels/statistics', params)
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

  // 导入相关类型
  export interface ImportTemplateField {
    field: string
    code: string
    required: boolean
    type: string
    options?: string
    maxLength?: number
    example?: string
  }
  export interface ImportTemplateResp {
    fields: ImportTemplateField[]
    sheetName: string
    title: string
  }
  export interface ImportFailRow {
    row: number
    data: string
    reason: string
  }
  export interface ImportResultResp {
    total: number
    success: number
    failed: number
    failRows: ImportFailRow[]
  }

  // 获取导入模板
  export const getImportTemplate = () => {
    return request.get<ImportTemplateResp>('/personnels/template')
  }

  // 预览上传文件
  export const previewImport = (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    const url = getApiBase() + '/personnels/preview'
    return fetch(url, {
      method: 'POST',
      credentials: 'include',
      headers: { ...getAuthHeader() },
      body: formData
    }).then(res => {
      if (!res.ok) throw new Error('预览请求失败')
      return res.json()
    }) as Promise<{ code: number; message: string; data: { totalRows: number; dataRows: number; sheetName: string; headers: string[]; error?: string } }>
  }

  // 下载导入模板 Excel
  export const downloadTemplate = () => {
    const url = getApiBase() + '/personnels/template'
    return fetch(url, {
      credentials: 'include',
      headers: { ...getAuthHeader() }
    }).then(res => {
      if (!res.ok) throw new Error('获取模板信息失败')
      return res.json()
    }).then(data => {
      generatePersonnelTemplate(data.data || data as ImportTemplateResp)
    })
  }

  // 导入人员（带进度回调）
  export const importPersonnel = (file: File, onProgress?: (pct: number) => void) => {
    return new Promise<import('.').ResultData<ImportResultResp>>((resolve, reject) => {
      const xhr = new XMLHttpRequest()
      const formData = new FormData()
      formData.append('file', file)
      xhr.open('POST', (import.meta.env.VITE_API_URL as string || '/api') + '/personnels/import')
      xhr.setRequestHeader('Authorization', `Bearer ${localStorage.getItem('token') || ''}`)
      xhr.withCredentials = true
      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          const resp = JSON.parse(xhr.responseText)
          if (resp.code === 200) {
            resolve({ code: 200, message: resp.message, data: resp.data })
          } else {
            reject(new Error(resp.message || '导入失败'))
          }
        } else {
          reject(new Error('导入请求失败'))
        }
        onProgress?.(100)
      }
      xhr.onerror = () => reject(new Error('网络错误'))
      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable && onProgress) {
          onProgress(Math.round((e.loaded / e.total) * 90))
        }
      }
      xhr.send(formData)
    })
  }
}

// 前端生成人员模板Excel
function generatePersonnelTemplate(template: ImportTemplateResp) {
  const script = document.createElement('script')
  script.src = 'https://cdn.sheetjs.com/xlsx-0.20.1/package/dist/xlsx.full.min.js'
  script.onload = () => {
    const XLSX = (window as any).XLSX

    // 第1行：标题
    const titleRow = [{ v: template.title || '人员导入模板', t: 's' }]
    // 第2行：必填说明
    const reqFields = template.fields.filter(f => f.required).map(f => f.field).join('、')
    const hintRow = [{ v: `* 必填字段：${reqFields || '姓名'}`, t: 's' }]
    // 第3行：表头
    const headerRow = template.fields.map(f => ({ v: f.required ? `${f.field} *` : f.field, t: 's' }))
    // 第4行：示例
    const demoRow = template.fields.map(f => ({ v: f.example || '', t: 's' }))

    const allRows = [titleRow, hintRow, headerRow, demoRow]
    const ws = XLSX.utils.aoa_to_sheet(allRows)

    const lastCol = template.fields.length - 1
    ws['!merges'] = [{ s: { r: 0, c: 0 }, e: { r: 0, c: lastCol } }]
    ws['!rows'] = [{ hpt: 36 }, { hpt: 22 }, { hpt: 30 }, { hpt: 28 }]
    ws['!cols'] = template.fields.map(f => ({ wch: Math.max(f.field.length, (f.example?.length || 0) + 2, 12) }))

    // 表头样式
    for (let c = 0; c < template.fields.length; c++) {
      const addr3 = XLSX.utils.encode_cell({ r: 2, c })
      const addr4 = XLSX.utils.encode_cell({ r: 3, c })
      ;(ws[addr3] as any).s = {
        fill: { fgColor: { rgb: '1F4E79' }, patternType: 'solid' },
        font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 11 },
        alignment: { horizontal: 'center', vertical: 'center' },
        border: { top: { style: 'thin', color: { rgb: '2E75B6' } }, bottom: { style: 'thin', color: { rgb: '2E75B6' } }, left: { style: 'thin', color: { rgb: '2E75B6' } }, right: { style: 'thin', color: { rgb: '2E75B6' } }
        }
      }
      ;(ws[addr4] as any).s = {
        fill: { fgColor: { rgb: 'EBF3FB' }, patternType: 'solid' },
        font: { color: { rgb: '1F4E79' }, sz: 10 },
        alignment: { horizontal: 'center', vertical: 'center' },
        border: { top: { style: 'thin', color: { rgb: 'BDD7EE' } }, bottom: { style: 'thin', color: { rgb: 'BDD7EE' } }, left: { style: 'thin', color: { rgb: 'BDD7EE' } }, right: { style: 'thin', color: { rgb: 'BDD7EE' } }
        }
      }
    }

    // 说明行样式
    const hintCell = XLSX.utils.encode_cell({ r: 1, c: 0 })
    ;(ws[hintCell] as any).s = { font: { color: { rgb: '595959' }, sz: 10 }, fill: { fgColor: { rgb: 'F2F2F2' }, patternType: 'solid' } }
    // 标题行样式
    const titleCell = XLSX.utils.encode_cell({ r: 0, c: 0 })
    ;(ws[titleCell] as any).s = {
      font: { bold: true, color: { rgb: 'FFFFFF' }, sz: 14 },
      fill: { fgColor: { rgb: '1F4E79' }, patternType: 'solid' },
      alignment: { horizontal: 'center', vertical: 'center' }
    }

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, template.sheetName || 'Sheet1')
    XLSX.writeFile(wb, `${template.title || '人员导入模板'}.xlsx`)
  }
  document.head.appendChild(script)
}
