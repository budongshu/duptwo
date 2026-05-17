import { ref, computed } from 'vue'

// 全局筛选状态（供所有页面共享）
export interface GlobalFilter {
  projectName: string
  diskLabel: string
  status: string
  uploader: string
  startDate: string
  endDate: string
}

const globalFilter = ref<GlobalFilter>({
  projectName: '',
  diskLabel: '',
  status: '',
  uploader: '',
  startDate: '',
  endDate: '',
})

// 是否有激活的筛选
export const hasActiveGlobalFilter = computed(() => {
  return !!(
    globalFilter.value.projectName ||
    globalFilter.value.diskLabel ||
    globalFilter.value.status ||
    globalFilter.value.uploader ||
    globalFilter.value.startDate ||
    globalFilter.value.endDate
  )
})

// 设置筛选条件（合并更新）
export const setGlobalFilter = (filter: Partial<GlobalFilter>) => {
  globalFilter.value = { ...globalFilter.value, ...filter }
}

// 清空筛选条件
export const clearGlobalFilter = () => {
  globalFilter.value = {
    projectName: '',
    diskLabel: '',
    status: '',
    uploader: '',
    startDate: '',
    endDate: '',
  }
}

// 获取 API 请求参数
export const getGlobalFilterParams = () => {
  const params: Record<string, string> = {}
  if (globalFilter.value.projectName) params.projectName = globalFilter.value.projectName
  if (globalFilter.value.diskLabel) params.diskLabel = globalFilter.value.diskLabel
  if (globalFilter.value.status) params.status = globalFilter.value.status
  if (globalFilter.value.uploader) params.uploader = globalFilter.value.uploader
  if (globalFilter.value.startDate) params.startDate = globalFilter.value.startDate
  if (globalFilter.value.endDate) params.endDate = globalFilter.value.endDate
  return params
}