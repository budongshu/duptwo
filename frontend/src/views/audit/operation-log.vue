<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">{{ t('audit.operation.title') }}</h1>
        <p class="page-subtitle">{{ t('audit.operation.subtitle') }}</p>
      </div>
      <div class="header-actions">
        <el-button type="success" @click="handleExport">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          {{ t('audit.operation.exportExcel') }}
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <el-input v-model="searchKeyword" :placeholder="t('audit.operation.searchPlaceholder')" clearable @keyup.enter="handleSearch" />
      <el-select v-model="searchAction" :placeholder="t('audit.operation.actionType')" clearable style="width: 120px" @change="handleSearch">
        <el-option :label="t('audit.operation.action.view')" value="view" />
        <el-option :label="t('audit.operation.action.create')" value="create" />
        <el-option :label="t('audit.operation.action.update')" value="update" />
        <el-option :label="t('audit.operation.action.delete')" value="delete" />
        <el-option :label="t('audit.operation.action.export')" value="export" />
      </el-select>
      <el-select v-model="searchResourceType" :placeholder="t('audit.operation.resourceType')" clearable style="width: 140px" @change="handleSearch">
        <el-option :label="t('audit.operation.resource.User')" value="User" />
        <el-option :label="t('audit.operation.resource.Role')" value="Role" />
        <el-option :label="t('audit.operation.resource.UserGroup')" value="UserGroup" />
        <el-option :label="t('audit.operation.resource.UploadRecord')" value="UploadRecord" />
        <el-option :label="t('audit.operation.resource.FieldConfig')" value="FieldConfig" />
      </el-select>
      <el-date-picker
        v-model="searchDateRange"
        type="daterange"
        :range-separator="t('common.to')"
        :start-placeholder="t('common.startDate')"
        :end-placeholder="t('common.endDate')"
        value-format="YYYY-MM-DD"
        style="width: 300px"
      />
      <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
      <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="username" :label="t('audit.operation.username')" min-width="100" />
        <el-table-column prop="menuName" :label="t('audit.operation.menu')" min-width="140">
          <template #default="{ row }">
            <span class="menu-name">{{ row.menuName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="action" :label="t('audit.operation.action')" min-width="80" align="center">
          <template #default="{ row }">
            <span class="action-badge" :class="getActionClass(row.action)">
              {{ row.actionText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="resourceType" :label="t('audit.operation.resourceType')" min-width="100" align="center">
          <template #default="{ row }">
            <span class="type-tag">{{ row.resourceType || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="resourceName" :label="t('audit.operation.resourceName')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="ipAddress" :label="t('audit.operation.ipAddress')" min-width="120" align="center">
          <template #default="{ row }">
            <span class="ip-text">{{ row.ipAddress || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('audit.operation.actionTime')" min-width="140">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="sizes, prev, pager, next"
          background
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, inject, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { OperationLogApi, type OperationLog } from '@/api/audit'

const { t } = useI18n()
const trackExport = inject<(success?: boolean) => void>('trackExport')

const loading = ref(false)
const tableData = ref<OperationLog[]>([])
const searchKeyword = ref('')
const searchAction = ref('')
const searchResourceType = ref('')
const searchDateRange = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const getActionClass = (action: string) => {
  const map: Record<string, string> = {
    view: 'action-badge--info',
    create: 'action-badge--success',
    update: 'action-badge--warning',
    delete: 'action-badge--danger',
    export: 'action-badge--primary'
  }
  return map[action] || ''
}

// 格式化日期
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await OperationLogApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      action: searchAction.value || undefined,
      resourceType: searchResourceType.value || undefined,
      startDate: searchDateRange.value?.[0] || undefined,
      endDate: searchDateRange.value?.[1] || undefined
    })
    if (res.code === 200) {
      tableData.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  searchKeyword.value = ''
  searchAction.value = ''
  searchResourceType.value = ''
  searchDateRange.value = []
  handleSearch()
}

const handleExport = async () => {
  try {
    const res = await OperationLogApi.exportExcel({
      keyword: searchKeyword.value || undefined,
      action: searchAction.value || undefined,
      resourceType: searchResourceType.value || undefined,
      startDate: searchDateRange.value?.[0] || undefined,
      endDate: searchDateRange.value?.[1] || undefined
    })
    const blob = new Blob([res as unknown as BlobPart], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 10)
    link.href = url
    link.download = `${t('audit.operation.fileName')}_${timestamp}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success(t('audit.operation.exportSuccess'))
    trackExport?.(true)
  } catch (error) {
    ElMessage.error(t('audit.operation.exportFailed'))
  }
}

onMounted(() => {
  loadData()
})

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
</script>

<style scoped lang="scss">
.page { min-height: 100vh; background: var(--color-page-bg); padding: 20px; }

.page-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 16px;
  padding: 20px 24px; background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
}

.page-title { font-family: 'Manrope', sans-serif; font-size: 20px; font-weight: 800; color: var(--color-text-primary); margin-bottom: 2px; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); }
.header-actions { display: flex; gap: 8px; }

.filter-card {
  display: flex; flex-wrap: wrap; gap: 10px; align-items: center;
  background: var(--color-surface); border-radius: 12px; padding: 16px 20px; margin-bottom: 14px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
  > * { flex: 0 0 auto; }
  .el-input { width: 180px; }
}

.table-card { background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light); overflow: hidden; }

:deep(.el-table) {
  --el-table-border-color: var(--color-border-light); width: 100% !important;
  th.el-table__cell {
    background-color: var(--color-surface-2) !important; border-bottom: 1px solid var(--color-border-light) !important;
    padding: 14px 16px !important; color: var(--color-text-secondary) !important; font-weight: 700 !important; font-size: 12px !important;
  }
  td.el-table__cell { padding: 14px 16px !important; font-size: 14px; color: var(--color-text-primary); border-bottom: 1px solid var(--color-border-light) !important; }
  .el-table__body tr:hover > td.el-table__cell { background-color: var(--color-surface-2) !important; }
}

.menu-name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.action-badge {
  display: inline-block; padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 600;
  &--info { background: var(--gray-100); color: var(--color-text-secondary); }
  &--success { background: var(--color-success-bg); color: var(--color-success); }
  &--warning { background: var(--color-warning-bg); color: var(--color-warning-text); }
  &--danger { background: var(--color-danger-bg); color: var(--color-danger); }
  &--primary { background: var(--color-primary-light-9); color: var(--color-primary); }
}

.type-tag {
  display: inline-block; padding: 4px 10px; background: var(--color-primary-light-9); color: var(--color-primary);
  border-radius: 6px; font-size: 12px; font-weight: 600;
}

.ip-text {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination-wrapper { padding: 16px 20px; display: flex; justify-content: flex-end; border-top: 1px solid var(--color-border-light); }
</style>
