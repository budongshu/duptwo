<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">登录日志</h1>
        <p class="page-subtitle">记录用户登录系统的情况</p>
      </div>
      <div class="header-actions">
        <el-button type="success" @click="handleExport">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          导出 Excel
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <el-input v-model="searchKeyword" placeholder="搜索用户名" clearable @keyup.enter="handleSearch" />
      <el-select v-model="searchStatus" placeholder="登录状态" clearable style="width: 120px" @change="handleSearch">
        <el-option label="成功" value="success" />
        <el-option label="失败" value="failed" />
      </el-select>
      <el-date-picker
        v-model="searchDateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        style="width: 300px"
      />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="status" label="状态" min-width="90" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="row.status === 'success' ? 'status-badge--success' : 'status-badge--danger'">
              <svg v-if="row.status === 'success'" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
              <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              {{ row.statusText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="loginMethod" label="登录方式" min-width="100" align="center">
          <template #default="{ row }">
            <span class="method-tag">{{ row.loginMethod === 'mfa' ? 'MFA' : '密码' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="mfaUsed" label="MFA验证" min-width="90" align="center">
          <template #default="{ row }">
            <span v-if="row.mfaUsed" class="mfa-enabled">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              已验证
            </span>
            <span v-else class="mfa-disabled">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="failReason" label="失败原因" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="fail-reason" v-if="row.failReason">{{ row.failReason }}</span>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="ipAddress" label="IP地址" min-width="130" align="center">
          <template #default="{ row }">
            <span class="ip-text">{{ row.ipAddress || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="登录时间" min-width="140">
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
          layout="total, sizes, prev, pager, next"
          background
          @current-change="loadData"
          @size-change="loadData"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, inject } from 'vue'
import { ElMessage } from 'element-plus'
import { LoginLogApi, type LoginLog } from '@/api/audit'

const trackExport = inject<(action?: string) => void>('trackExport')

const loading = ref(false)
const tableData = ref<LoginLog[]>([])
const searchKeyword = ref('')
const searchStatus = ref('')
const searchDateRange = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 格式化日期
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await LoginLogApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
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
  searchStatus.value = ''
  searchDateRange.value = []
  handleSearch()
}

const handleExport = async () => {
  try {
    const res = await LoginLogApi.exportExcel({
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
      startDate: searchDateRange.value?.[0] || undefined,
      endDate: searchDateRange.value?.[1] || undefined
    })
    const blob = new Blob([res as unknown as BlobPart], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 10)
    link.href = url
    link.download = `登录日志_${timestamp}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
    trackExport?.('login-log')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

onMounted(() => {
  loadData()
})
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

.status-badge {
  display: inline-flex; align-items: center; gap: 4px; padding: 4px 10px; border-radius: 6px; font-size: 13px; font-weight: 600;
  &--success { background: var(--color-success-bg); color: var(--color-success); }
  &--danger { background: var(--color-danger-bg); color: var(--color-danger); }
}

.method-tag {
  display: inline-block; padding: 4px 10px; background: var(--color-primary-light-9); color: var(--color-primary);
  border-radius: 6px; font-size: 12px; font-weight: 600;
}

.mfa-enabled {
  display: inline-flex; align-items: center; gap: 4px;
  color: var(--color-success); font-size: 13px; font-weight: 500;
}

.mfa-disabled {
  color: var(--color-text-secondary); font-size: 14px;
}

.fail-reason {
  color: var(--color-danger); font-size: 13px;
}

.empty-text {
  color: var(--color-text-muted); font-size: 14px;
}

.ip-text {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination-wrapper { padding: 16px 20px; display: flex; justify-content: flex-end; border-top: 1px solid var(--color-border-light); }
</style>
