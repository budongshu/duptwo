<template>
  <div class="dashboard">
    <!-- ==================== 主内容区 ==================== -->
    <main class="dashboard-main">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="header-left">
          <h1 class="page-title">{{ t('uploadRecord.dashboard.title') }}</h1>
          <div class="auto-refresh" :class="{ refreshing: nextRefresh <= 5 }">
            <span class="refresh-dot"></span>
            <span class="refresh-text">{{ nextRefresh }}s</span>
          </div>
        </div>
        <div class="page-meta">
          <span v-if="hasActiveFilters" class="filter-tags">
            <span v-if="searchProjectName" class="filter-tag">{{ searchProjectName }}<button @click="searchProjectName = ''; handleFilterChange()">×</button></span>
            <span v-if="searchDataType" class="filter-tag">{{ searchDataType }}<button @click="searchDataType = ''; handleFilterChange()">×</button></span>
            <span v-if="searchStatus" class="filter-tag">{{ t(getStatusText(searchStatus)) }}<button @click="searchStatus = ''; handleFilterChange()">×</button></span>
            <span v-if="searchUploader" class="filter-tag">{{ searchUploader }}<button @click="searchUploader = ''; handleFilterChange()">×</button></span>
          </span>
          <el-button type="primary" @click="filterDrawerVisible = true" class="filter-btn" :class="{ 'has-badge': hasActiveFilters }">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 5px"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
            {{ t('uploadRecord.dashboard.filter') }}
          </el-button>
        </div>
      </div>

      <!-- ==================== 核心指标卡片 ==================== -->
      <div class="stats-row" v-if="isInitialLoading">
        <div v-for="i in 5" :key="i" class="stat-card skeleton-card">
          <div class="skeleton" style="height: 12px; width: 60%; margin-bottom: 12px;"></div>
          <div class="skeleton" style="height: 32px; width: 80%; margin-bottom: 8px;"></div>
          <div class="skeleton" style="height: 10px; width: 45%;"></div>
        </div>
      </div>
      <div class="stats-row" v-else>
        <div class="stat-card stat-card--primary">
          <div class="stat-card__top">
            <span class="stat-card__label">{{ t('uploadRecord.dashboard.todayUpload') }}</span>
            <div class="stat-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/></svg>
            </div>
          </div>
          <div class="stat-card__big-value">
            <span class="big-counter">{{ statistics.todayCount }}</span>
            <span class="big-unit">{{ t('uploadRecord.dashboard.recordCount') }}</span>
          </div>
          <div class="stat-card__meta-row">
            <span class="stat-card__size">{{ statistics.todaySizeStr }}</span>
            <span class="stat-card__trend" :class="todayTrend >= 0 ? 'trend--up' : 'trend--down'">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path v-if="todayTrend >= 0" d="M18 15l-6-6-6 6"/><path v-else d="M6 9l6 6 6-6"/></svg>
              {{ Math.abs(todayTrend) }}%
            </span>
          </div>
        </div>

        <div class="stat-card stat-card--week">
          <div class="stat-card__top">
            <span class="stat-card__label">{{ t('uploadRecord.dashboard.weekUpload') }}</span>
            <div class="stat-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            </div>
          </div>
          <div class="stat-card__big-value">
            <span class="big-counter">{{ statistics.weekCount }}</span>
            <span class="big-unit">{{ t('uploadRecord.dashboard.recordCount') }}</span>
          </div>
          <div class="stat-card__meta-row">
            <span class="stat-card__size">{{ statistics.weekSizeStr || '—' }}</span>
            <span class="stat-card__trend" :class="weekTrend >= 0 ? 'trend--up' : 'trend--down'">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path v-if="weekTrend >= 0" d="M18 15l-6-6-6 6"/><path v-else d="M6 9l6 6 6-6"/></svg>
              {{ Math.abs(weekTrend) }}%
            </span>
          </div>
        </div>

        <div class="stat-card stat-card--month">
          <div class="stat-card__top">
            <span class="stat-card__label">{{ t('uploadRecord.dashboard.monthUpload') }}</span>
            <div class="stat-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
            </div>
          </div>
          <div class="stat-card__big-value">
            <span class="big-counter">{{ statistics.monthCount }}</span>
            <span class="big-unit">{{ t('uploadRecord.dashboard.recordCount') }}</span>
          </div>
          <div class="stat-card__meta-row">
            <span class="stat-card__size">{{ statistics.monthSizeStr || '—' }}</span>
            <span class="stat-card__trend" :class="monthTrend >= 0 ? 'trend--up' : 'trend--down'">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path v-if="monthTrend >= 0" d="M18 15l-6-6-6 6"/><path v-else d="M6 9l6 6 6-6"/></svg>
              {{ Math.abs(monthTrend) }}%
            </span>
          </div>
        </div>

        <div class="stat-card stat-card--total">
          <div class="stat-card__top">
            <span class="stat-card__label">{{ t('uploadRecord.dashboard.totalCount') }}</span>
            <div class="stat-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>
            </div>
          </div>
          <div class="stat-card__big-value">
            <span class="big-counter">{{ statistics.totalCount }}</span>
            <span class="big-unit">{{ t('uploadRecord.dashboard.recordCount') }}</span>
          </div>
          <div class="stat-card__meta-row">
            <span class="stat-card__size">{{ statistics.totalSizeStr }}</span>
          </div>
        </div>

        <div class="stat-card stat-card--rate">
          <div class="stat-card__top">
            <span class="stat-card__label">{{ t('uploadRecord.dashboard.successRate') }}</span>
            <div class="stat-card__icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            </div>
          </div>
          <div class="stat-card__big-value">
            <span class="big-counter">{{ successRate }}</span>
            <span class="big-unit">%</span>
          </div>
          <div class="stat-card__meta-row">
            <div class="rate-bar"><div class="rate-bar__fill" :style="{ width: successRate + '%' }"></div></div>
          </div>
        </div>
      </div>

      <!-- ==================== 图表区域 ==================== -->
      <div class="charts-grid">
        <div class="panel panel--trend card-appear">
          <div class="panel__header">
            <div class="panel__tabs">
              <button v-for="tab in trendTabs" :key="tab.key" class="panel-tab" :class="{ 'is-active': activeTrendTab === tab.key }" @click="setTrendTab(tab.key)">
                <span class="tab-dot" :style="{ background: tab.color }"></span>
                {{ tab.label }}
              </button>
            </div>
            <span class="panel__meta">共 {{ totalTrendCount }} {{ t('uploadRecord.dashboard.recordCount') }}</span>
          </div>
          <div class="panel__body panel__body--chart">
            <div ref="trendChartRef" class="chart-canvas"></div>
          </div>
        </div>

        <div class="panel panel--status card-appear" style="animation-delay: 0.1s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.statusDistribution') }}</span>
            <span class="panel__meta">{{ (statistics.byStatus || []).length }} {{ t('uploadRecord.dashboard.statusTypes') }}</span>
          </div>
          <div class="panel__body panel__body--status">
            <div ref="donutChartRef" class="chart-canvas chart-canvas--donut"></div>
            <div class="status-legend">
              <div v-for="item in (statistics.byStatus || [])" :key="item.status" class="legend-row">
                <span class="legend-bar" :style="{ background: statusColors[item.status] }"></span>
                <span class="legend-name">{{ t(getStatusText(item.status)) }}</span>
                <span class="legend-count">{{ item.count }}</span>
                <span class="legend-pct">{{ getPercent(item.count) }}%</span>
              </div>
            </div>
          </div>
        </div>

        <div class="panel panel--project card-appear" style="animation-delay: 0.2s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.projectDistribution') }}</span>
            <span class="panel__meta">{{ projectTotal }} {{ t('uploadRecord.dashboard.projectCount') }}</span>
          </div>
          <div class="panel__body">
            <div ref="projectChartRef" class="chart-canvas chart-canvas--bar"></div>
          </div>
        </div>

        <div class="panel panel--datatype card-appear" style="animation-delay: 0.3s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.dataTypeDistribution') }}</span>
            <span class="panel__meta">{{ dataTypeTotal }} {{ t('uploadRecord.dashboard.dataTypeCount') }}</span>
          </div>
          <div class="panel__body">
            <div ref="dataTypeChartRef" class="chart-canvas chart-canvas--bar"></div>
          </div>
        </div>

        <div class="panel panel--health card-appear" style="animation-delay: 0.4s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.processingHealth') }}</span>
          </div>
          <div class="panel__body panel__body--health">
            <div class="health-gauge">
              <div class="gauge-ring" :style="{ '--pct': successRate }">
                <svg viewBox="0 0 200 200">
                  <circle cx="100" cy="100" r="80" fill="none" stroke="#f0f0f5" stroke-width="14"/>
                  <circle cx="100" cy="100" r="80" fill="none" :stroke="successRate >= 80 ? '#22c55e' : successRate >= 50 ? '#f59e0b' : '#ef4444'" stroke-width="14" stroke-linecap="round" :stroke-dasharray="`${successRate * 5.02} 502`" transform="rotate(-90 100 100)"/>
                </svg>
                <div class="gauge-center">
                  <span class="gauge-value">{{ successRate }}%</span>
                  <span class="gauge-label">{{ t('uploadRecord.dashboard.ratePercent') }}</span>
                </div>
              </div>
            </div>
            <div class="health-bars">
              <div v-for="item in (statistics.byStatus || [])" :key="item.status" class="health-bar-item">
                <div class="health-bar-item__header">
                  <span class="health-bar-item__name">{{ t(getStatusText(item.status)) }}</span>
                  <span class="health-bar-item__count">{{ item.count }}</span>
                </div>
                <div class="health-bar-item__track">
                  <div class="health-bar-item__fill" :style="{ width: getPercent(item.count) + '%', background: statusColors[item.status] }"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 近期活动 -->
        <div class="panel panel--activity card-appear" style="animation-delay: 0.5s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.recentActivity') }}</span>
            <router-link to="/upload-record/list" class="panel__link">{{ t('uploadRecord.dashboard.viewAll') }} →</router-link>
          </div>
          <div class="panel__body panel__body--activity">
            <div class="activity-table">
              <div class="table-head">
                <span>#</span>
                <span>{{ t('uploadRecord.dashboard.tableSerial') }}</span>
                <span>{{ t('uploadRecord.dashboard.tableDisk') }}</span>
                <span>{{ t('uploadRecord.dashboard.tableUploadSize') }}</span>
                <span>{{ t('uploadRecord.dashboard.tableStatus') }}</span>
                <span>{{ t('uploadRecord.dashboard.tableTime') }}</span>
              </div>
              <div class="table-body">
                <div v-for="(item, index) in recentRecords" :key="item.id" class="table-row">
                  <span class="row-num">{{ String(index + 1).padStart(2, '0') }}</span>
                  <span class="row-serial" :title="item.serialNo">{{ item.serialNo }}</span>
                  <span class="col-disk" :title="item.dataType">{{ item.dataType }}</span>
                  <span class="col-size">{{ item.fileSizeStr }}</span>
                  <span>
                    <span class="status-badge" :class="getStatusClass(item.status)">{{ item.statusText }}</span>
                  </span>
                  <span class="row-time">{{ formatTime(item.createdAt) }}</span>
                </div>
                <div v-if="recentRecords.length === 0" class="table-empty">{{ t('uploadRecord.dashboard.noUploadRecord') }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- ==================== 筛选抽屉 ==================== -->
    <el-drawer v-model="filterDrawerVisible" direction="rtl" size="320px" :with-header="true">
      <template #header>
        <div class="filter-drawer-header">
          <div class="filter-drawer-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
            </svg>
          </div>
          <span class="filter-drawer-title">{{ t('uploadRecord.dashboard.filterCondition') }}</span>
        </div>
      </template>

      <div class="filter-drawer-body">
        <!-- 时间范围 -->
        <div class="filter-card">
          <div class="filter-card-header">
            <span>时间范围</span>
          </div>
          <div class="filter-card-body">
            <div class="quick-presets">
              <button
                v-for="preset in timePresets"
                :key="preset.key"
                class="preset-btn"
                :class="{ 'is-active': activeTimePreset === preset.key }"
                @click="setTimePreset(preset.key)"
              >
                {{ preset.label }}
              </button>
            </div>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              :clearable="false"
              size="small"
              style="width: 100%"
              @change="loadStatistics"
            />
          </div>
        </div>

        <!-- 项目名称 -->
        <div class="filter-card">
          <div class="filter-card-header">
            <span>项目名称</span>
          </div>
          <div class="filter-card-body">
            <el-select
              v-model="searchProjectName"
              :placeholder="t('uploadRecord.dashboard.filterAllProjects')"
              clearable
              filterable
              size="small"
              style="width: 100%"
              @change="handleFilterChange"
            >
              <el-option v-for="p in projectList" :key="p.code" :label="p.name" :value="p.name" />
            </el-select>
          </div>
        </div>

        <!-- 磁盘标签 -->
        <div class="filter-card">
          <div class="filter-card-header">
            <span>磁盘标签</span>
          </div>
          <div class="filter-card-body">
            <el-select
              v-model="searchDataType"
              :placeholder="t('uploadRecord.dashboard.filterAllDisk')"
              clearable
              filterable
              size="small"
              style="width: 100%"
              @change="handleFilterChange"
            >
              <el-option v-for="dt in dataTypeOptions" :key="dt" :label="dt" :value="dt" />
            </el-select>
          </div>
        </div>

        <!-- 上传状态 -->
        <div class="filter-card">
          <div class="filter-card-header">
            <span>上传状态</span>
          </div>
          <div class="filter-card-body">
            <el-select
              v-model="searchStatus"
              :placeholder="t('uploadRecord.dashboard.filterAllStatus')"
              clearable
              size="small"
              style="width: 100%"
              @change="handleFilterChange"
            >
              <el-option :label="t('status.completed')" value="completed" />
              <el-option :label="t('status.processing')" value="processing" />
              <el-option :label="t('status.pending')" value="pending" />
              <el-option :label="t('status.failed')" value="failed" />
            </el-select>
          </div>
        </div>

        <!-- 上传人 -->
        <div class="filter-card">
          <div class="filter-card-header">
            <span>上传人</span>
          </div>
          <div class="filter-card-body">
            <el-select
              v-model="searchUploader"
              :placeholder="t('uploadRecord.dashboard.filterAllUploader')"
              clearable
              filterable
              size="small"
              style="width: 100%"
              @change="handleFilterChange"
            >
              <el-option v-for="u in uploaderOptions" :key="u" :label="u" :value="u" />
            </el-select>
          </div>
        </div>
      </div>

      <div class="filter-drawer-footer">
        <el-button @click="resetFilters" size="small">重置</el-button>
        <el-button type="primary" size="small" @click="filterDrawerVisible = false">完成</el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import type { UploadRecordStatistics, UploadRecord } from '@/api/upload-record'
import { UploadRecordApi } from '@/api/upload-record'
import { ProjectApi } from '@/api/project'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

let echarts: any = null
const isInitialLoading = ref(true)

const dateRange = ref<string[]>([])
const activeTimePreset = ref<'today' | 'week' | 'month' | 'all'>('month')
const nextRefresh = ref(30)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const timePresets = [
  { key: 'today' as const, label: t('uploadRecord.dashboard.timePresetToday') },
  { key: 'week' as const, label: t('uploadRecord.dashboard.timePresetWeek') },
  { key: 'month' as const, label: t('uploadRecord.dashboard.timePresetMonth') },
  { key: 'all' as const, label: t('uploadRecord.dashboard.timePresetAll') }
]

const setTimeRange = (key: typeof activeTimePreset.value) => {
  const today = new Date()
  const fmt = (d: Date) => d.toISOString().split('T')[0]
  switch (key) {
    case 'today': dateRange.value = [fmt(today), fmt(today)]; break
    case 'week': { const d = new Date(today); d.setDate(d.getDate() - 6); dateRange.value = [fmt(d), fmt(today)]; break }
    case 'month': { const d = new Date(today); d.setDate(d.getDate() - 29); dateRange.value = [fmt(d), fmt(today)]; break }
    case 'all': dateRange.value = []; break
  }
}

const setTimePreset = (key: typeof activeTimePreset.value) => {
  activeTimePreset.value = key
  setTimeRange(key)
  loadStatistics()
}

const searchProjectName = ref('')
const searchDataType = ref('')
const searchStatus = ref('')
const searchUploader = ref('')
const projectList = ref<{ name: string; code: string }[]>([])
const dataTypeOptions = ref<string[]>([])
const uploaderOptions = ref<string[]>([])
const filterDrawerVisible = ref(false)

const hasActiveFilters = computed(() => !!(searchProjectName.value || searchDataType.value || searchStatus.value || searchUploader.value || dateRange.value?.length))

const statistics = reactive<UploadRecordStatistics>({
  todayCount: 0, todaySize: 0, todaySizeStr: '0 B',
  weekCount: 0, weekSize: 0, weekSizeStr: '',
  monthCount: 0, monthSize: 0, monthSizeStr: '',
  totalCount: 0, totalSize: 0, totalSizeStr: '0 B',
  trend: [], byStatus: [], byDataType: [], byProject: []
})

const recentRecords = ref<UploadRecord[]>([])
const todayTrend = ref(0)
const weekTrend = ref(0)
const monthTrend = ref(0)

const statusColors: Record<string, string> = { completed: '#22c55e', pending: '#f59e0b', processing: '#3b82f6', failed: '#ef4444' }
const getStatusText = (status: string) => ({ completed: 'status.completed', pending: 'status.pending', processing: 'status.processing', failed: 'status.failed' })[status] || status
const getStatusClass = (status: string) => `badge--${({ completed: 'success', pending: 'warning', processing: 'info', failed: 'danger' })[status] || 'info'}`
const getPercent = (count: number) => { const arr = statistics.byStatus || []; const total = arr.reduce((s, i) => s + i.count, 0); return total ? Math.round((count / total) * 100) : 0 }
const successRate = computed(() => { if (statistics.totalCount === 0) return 0; const arr = statistics.byStatus || []; const completed = arr.find(s => s.status === 'completed'); return completed ? Math.round((completed.count / statistics.totalCount) * 100) : 0 })
const totalTrendCount = computed(() => (statistics.trend || []).reduce((s, t) => s + t.count, 0))
const dataTypeTotal = computed(() => (statistics.byDataType || []).length)
const projectTotal = computed(() => (statistics.byProject || []).length)

// 将后端返回的累计趋势数据转换为每日增量
const trendDaily = computed(() => {
  const arr = statistics.trend || []
  if (arr.length === 0) return []
  return arr.map((t: any, idx: number) => {
    const prevCount = idx > 0 ? (arr[idx - 1]?.count || 0) : 0
    const prevSize = idx > 0 ? (arr[idx - 1]?.totalSize || 0) : 0
    return {
      date: t.date,
      count: t.count - prevCount,
      totalSize: t.totalSize - prevSize,
    }
  })
})
const formatTime = (timeStr: string) => { const d = new Date(timeStr); return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}` }

let trendChart: any = null, donutChart: any = null, dataTypeChart: any = null, projectChart: any = null
const activeTrendTab = ref('count')
const trendTabs = [{ key: 'count', label: t('uploadRecord.dashboard.chartRecordCount'), color: '#3b82f6' }, { key: 'size', label: t('uploadRecord.dashboard.chartDataSize'), color: '#8b5cf6' }]
const setTrendTab = (key: string) => { activeTrendTab.value = key; nextTick(() => updateTrendChart()) }

const trendChartRef = ref<HTMLElement>()
const donutChartRef = ref<HTMLElement>()
const dataTypeChartRef = ref<HTMLElement>()
const projectChartRef = ref<HTMLElement>()

const initCharts = async () => {
  if (!trendChartRef.value) return
  echarts = await import('echarts')
  trendChart = echarts.init(trendChartRef.value)
  donutChart = echarts.init(donutChartRef.value)
  dataTypeChart = echarts.init(dataTypeChartRef.value)
  projectChart = echarts.init(projectChartRef.value)
  updateAllCharts()
  window.addEventListener('resize', handleResize)
}

const updateTrendChart = () => {
  if (!trendChart || !echarts) return
  const arr = trendDaily.value || []
  const field = activeTrendTab.value === 'count' ? 'count' : 'totalSize'
  const formatter = activeTrendTab.value === 'count' ? (v: number) => `${v} 条` : (v: number) => formatFileSizeStatic(v)
  trendChart.setOption({
    backgroundColor: 'transparent', tooltip: { trigger: 'axis', backgroundColor: '#fff', borderColor: '#e0e2ec', textStyle: { color: '#191c23', fontSize: 12 }, formatter: (params: any) => { const p = params[0]; return `${p.name}<br/>${p.marker} ${formatter(p.value)}` } },
    grid: { top: 10, right: 16, bottom: 10, left: 10, containLabel: true },
    xAxis: { type: 'category', data: arr.map((t: any) => t.date), axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#9ca3af', fontSize: 11, formatter: (v: string) => v.slice(5) } },
    yAxis: { type: 'value', show: false, boundaryGap: ['2%', '2%'] },
    series: [{ type: 'line', data: arr.map((t: any) => t[field]), smooth: 0.4, symbol: 'none', lineStyle: { color: '#3b82f6', width: 2 }, areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: 'rgba(59,130,246,0.15)' }, { offset: 1, color: 'rgba(59,130,246,0)' }]) } }]
  })
}

const updateDonutChart = () => {
  if (!donutChart || !echarts) return
  const arr = statistics.byStatus || []
  donutChart.setOption({ backgroundColor: 'transparent', tooltip: { trigger: 'item', backgroundColor: '#fff', borderColor: '#e0e2ec', textStyle: { color: '#191c23', fontSize: 12 } }, legend: { show: false }, series: [{ type: 'pie', radius: ['55%', '75%'], center: ['50%', '50%'], data: arr.map((s: any) => ({ name: s.status, value: s.count, itemStyle: { color: statusColors[s.status] } })), label: { show: false }, emphasis: { scale: true, scaleSize: 4 } }] })
}

const updateDataTypeChart = () => {
  if (!dataTypeChart || !echarts) return
  const all = statistics.byDataType || []
  const arr = all.slice(0, 10)
  // 反转数组用于显示（从下到上）
  const reversed = [...arr].reverse()
  dataTypeChart.setOption({
    backgroundColor: 'transparent',
    animation: true,
    animationDuration: 800,
    animationEasing: 'cubicOut',
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      textStyle: { color: '#191c23', fontSize: 12 },
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const item = params[0]
        const d = reversed[item.dataIndex]
        if (!d) return ''
        return `${d.dataType || '(空)'}<br/>${item.marker} ${d.count} 条<br/>数据量: ${formatFileSizeStatic(d.totalSize || 0)}`
      }
    },
    grid: { top: 10, right: 80, bottom: 10, left: 10, containLabel: true },
    xAxis: { type: 'value', show: false },
    yAxis: { type: 'category', data: reversed.map((d: any) => d.dataType || '(空)'), axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#414754', fontSize: 12 } },
    series: [{
      type: 'bar',
      data: reversed.map((d: any) => d.count),
      barWidth: 14,
      itemStyle: { borderRadius: [0, 6, 6, 0], color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [{ offset: 0, color: '#06b6d4' }, { offset: 1, color: '#3b82f6' }]) },
      label: { show: true, position: 'right', color: '#727785', fontSize: 11, formatter: (params: any) => formatFileSizeStatic(reversed[params.dataIndex]?.totalSize || 0) }
    }]
  })
}

const updateProjectChart = () => {
  if (!projectChart || !echarts) return
  const all = statistics.byProject || []
  const arr = all.slice(0, 10)
  const reversed = [...arr].reverse()
  projectChart.setOption({
    backgroundColor: 'transparent',
    animation: true,
    animationDuration: 800,
    animationEasing: 'cubicOut',
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      textStyle: { color: '#191c23', fontSize: 12 },
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const item = params[0]
        const p = reversed[item.dataIndex]
        if (!p) return ''
        return `${p.projectName || '(空项目)'}<br/>${item.marker} ${p.count} 条<br/>数据量: ${formatFileSizeStatic(p.totalSize || 0)}`
      }
    },
    grid: { top: 10, right: 80, bottom: 10, left: 10, containLabel: true },
    xAxis: { type: 'value', show: false },
    yAxis: { type: 'category', data: reversed.map((p: any) => p.projectName || '(空项目)'), axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#414754', fontSize: 12 } },
    series: [{
      type: 'bar',
      data: reversed.map((p: any) => p.count),
      barWidth: 14,
      itemStyle: { borderRadius: [0, 6, 6, 0], color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [{ offset: 0, color: '#8b5cf6' }, { offset: 1, color: '#a78bfa' }]) },
      label: { show: true, position: 'right', color: '#727785', fontSize: 11, formatter: (params: any) => formatFileSizeStatic(reversed[params.dataIndex]?.totalSize || 0) }
    }]
  })
}

const updateAllCharts = () => { updateTrendChart(); updateDonutChart(); updateDataTypeChart(); updateProjectChart() }

const formatFileSizeStatic = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

const handleResize = () => { trendChart?.resize(); donutChart?.resize(); dataTypeChart?.resize(); projectChart?.resize() }

const loadStatistics = async () => {
  try {
    const params: Record<string, any> = {}
    if (dateRange.value?.length === 2) { params.startDate = dateRange.value[0]; params.endDate = dateRange.value[1] }
    if (searchProjectName.value) params.projectName = searchProjectName.value
    if (searchDataType.value) params.dataType = searchDataType.value
    if (searchStatus.value) params.status = searchStatus.value
    if (searchUploader.value) params.uploader = searchUploader.value
    const res = await UploadRecordApi.statistics(params)
    const data = res.data || {}

    // 同步更新统计数据（确保立即显示）
    statistics.todayCount = data.todayCount ?? 0
    statistics.todaySize = data.todaySize ?? 0
    statistics.todaySizeStr = data.todaySizeStr ?? '0 B'
    statistics.weekCount = data.weekCount ?? 0
    statistics.weekSize = data.weekSize ?? 0
    statistics.weekSizeStr = data.weekSizeStr || '—'
    statistics.monthCount = data.monthCount ?? 0
    statistics.monthSize = data.monthSize ?? 0
    statistics.monthSizeStr = data.monthSizeStr || '—'
    statistics.totalCount = data.totalCount ?? 0
    statistics.totalSize = data.totalSize ?? 0
    statistics.totalSizeStr = data.totalSizeStr ?? '0 B'
    statistics.trend = data.trend || []
    statistics.byStatus = data.byStatus || []
    statistics.byDataType = data.byDataType || []
    statistics.byProject = data.byProject || []

    const trend = data.trend || []
    if (trend.length >= 2) {
      const last = trend[trend.length - 1]
      const prev2 = trend[Math.max(0, trend.length - 2)]
      const base = prev2.count || 1
      const diff = last.count - prev2.count
      const pct = Math.round((Math.abs(diff) / base) * 100)
      todayTrend.value = diff >= 0 ? pct : -pct
    }

    updateAllCharts()
  } catch (error) { console.error('Failed to load statistics:', error) }
}

const loadRecentRecords = async () => {
  try {
    const params: any = { limit: 20 }
    if (searchProjectName.value) params.projectName = searchProjectName.value
    if (searchDataType.value) params.dataType = searchDataType.value
    if (searchStatus.value) params.status = searchStatus.value
    if (searchUploader.value) params.uploader = searchUploader.value
    const res = await UploadRecordApi.recent(params)
    recentRecords.value = res.data || []
  } catch (error) { console.error('Failed to load recent records:', error) }
}

const loadUploaderList = async () => { try { const res = await UploadRecordApi.uploaderList(); uploaderOptions.value = res.data || [] } catch { uploaderOptions.value = [] } }
const loadProjects = async () => { try { const res = await ProjectApi.getSimpleList(); projectList.value = res.data || [] } catch (error) { console.error('Failed to load projects:', error) } }
const handleFilterChange = () => { loadStatistics(); loadRecentRecords() }
const resetFilters = () => { searchProjectName.value = ''; searchDataType.value = ''; searchStatus.value = ''; searchUploader.value = ''; setTimePreset('month'); loadStatistics(); loadRecentRecords() }

onMounted(() => {
  setTimeRange('month')
  initCharts()
  isInitialLoading.value = false
  loadProjects()
  loadUploaderList()
  loadStatistics()
  loadRecentRecords()
  refreshTimer = setInterval(() => {
    nextRefresh.value = 30
    const tick = setInterval(() => { nextRefresh.value--; if (nextRefresh.value <= 0) clearInterval(tick) }, 1000)
    loadStatistics()
    loadRecentRecords()
  }, 30000)
})

onUnmounted(() => { if (refreshTimer) clearInterval(refreshTimer); trendChart?.dispose(); donutChart?.dispose(); dataTypeChart?.dispose(); projectChart?.dispose(); window.removeEventListener('resize', handleResize) })
</script>

<style scoped lang="scss">
.dashboard { width: 100%; min-height: 100vh; background: var(--color-page-bg); }
.dashboard-main { padding: 24px 24px 32px; min-width: 0; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-wrap: wrap; gap: 10px; }
.header-left { display: flex; align-items: center; gap: 14px; }
.page-meta { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; position: relative; }
.page-title { font-family: 'Manrope', 'Inter', sans-serif; font-size: 24px; font-weight: 800; color: var(--color-text-primary); margin: 0; letter-spacing: -0.5px; -webkit-font-smoothing: antialiased; }
.auto-refresh { display: flex; align-items: center; gap: 6px; padding: 4px 10px; background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-xl); font-size: 12px; color: var(--color-text-secondary); }
.auto-refresh .refresh-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--color-success); animation: pulse 2s ease-in-out infinite; }
.auto-refresh.refreshing .refresh-dot { background: var(--color-warning); animation: none; }
.auto-refresh.refreshing .refresh-text { color: var(--color-warning); font-weight: 600; }
@keyframes pulse { 0%, 100% { opacity: 1; transform: scale(1); } 50% { opacity: 0.6; transform: scale(0.85); } }
.filter-btn { display: flex; align-items: center; font-weight: 600; }
.filter-btn.has-badge::after { content: ''; position: absolute; top: -2px; right: -2px; width: 8px; height: 8px; border-radius: 50%; background: var(--color-badge-dot); border: 2px solid var(--color-surface); }
.filter-tags { display: flex; gap: 6px; flex-wrap: wrap; }
.filter-tag { display: inline-flex; align-items: center; gap: 4px; padding: 3px 8px; background: var(--color-primary-light-9); border: 1px solid rgba(0,94,235,0.2); border-radius: var(--radius-full); font-size: 11.5px; color: var(--color-primary); font-weight: 500; }
.filter-tag button { border: none; background: none; color: var(--color-primary); cursor: pointer; font-size: 13px; padding: 0; line-height: 1; }
.filter-tag button:hover { color: var(--color-danger); }
.stats-row { display: grid; grid-template-columns: repeat(5, 1fr); gap: 14px; margin-bottom: 20px; }
@media (max-width: 1100px) { .stats-row { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 700px) { .stats-row { grid-template-columns: repeat(2, 1fr); } }
.stat-card { background: var(--color-surface); border: 1px solid var(--color-border-light); border-radius: 14px; padding: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); overflow: hidden; position: relative; transition: box-shadow 0.2s ease; display: flex; flex-direction: column; }
.stat-card:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.07); }
.skeleton-card .skeleton { background: linear-gradient(90deg, var(--gray-100) 25%, var(--gray-50) 50%, var(--gray-100) 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; border-radius: 4px; }
@keyframes shimmer { 0% { background-position: -200% 0; } 100% { background-position: 200% 0; } }
.stat-card__top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
.stat-card__label { font-size: 11.5px; font-weight: 600; color: var(--color-text-secondary); text-transform: uppercase; letter-spacing: 0.3px; }
.stat-card__icon { color: var(--color-text-muted); opacity: 0.6; }
.stat-card__big-value { display: flex; align-items: baseline; gap: 4px; margin-bottom: 10px; }
.big-counter { font-family: 'Manrope', 'Inter', sans-serif; font-size: 32px; font-weight: 800; color: var(--color-text-primary); line-height: 1; letter-spacing: -1.5px; -webkit-font-smoothing: antialiased; }
.big-unit { font-size: 13px; font-weight: 600; color: var(--color-text-muted); }
.stat-card__meta-row { display: flex; justify-content: space-between; align-items: center; margin-top: auto; }
.stat-card__size { font-family: 'Manrope', 'Inter', sans-serif; font-size: 13px; font-weight: 700; color: var(--color-text-primary); -webkit-font-smoothing: antialiased; }
.stat-card__trend { display: flex; align-items: center; gap: 2px; font-size: 11px; font-weight: 700; padding: 3px 8px; border-radius: 5px; }
.stat-card__trend.trend--up { color: var(--color-success); background: rgba(34,197,94,0.1); }
.stat-card__trend.trend--down { color: var(--color-danger); background: rgba(239,68,68,0.1); }
.rate-bar { flex: 1; height: 6px; background: var(--color-surface-3); border-radius: 3px; overflow: hidden; }
.rate-bar__fill { height: 100%; background: linear-gradient(90deg, #06b6d4, #3b82f6); border-radius: 3px; transition: width 0.6s ease; min-width: 2px; }
.charts-grid { display: grid; grid-template-columns: 2fr 1fr 1fr; grid-template-rows: auto auto auto; gap: 14px; }
.panel--trend { grid-column: 1 / 2; grid-row: 1 / 2; }
.panel--status { grid-column: 2 / 4; grid-row: 1 / 2; }
.panel--project { grid-column: 1 / 2; grid-row: 2 / 3; }
.panel--datatype { grid-column: 2 / 3; grid-row: 2 / 3; }
.panel--health { grid-column: 3 / 4; grid-row: 2 / 3; }
.panel--activity { grid-column: 1 / 4; grid-row: 3 / 4; }
@media (max-width: 1100px) {
  .charts-grid { grid-template-columns: 1fr 1fr; }
  .panel--trend { grid-column: 1 / 3; }
  .panel--status { grid-column: 1 / 3; }
  .panel--project { grid-column: 1 / 2; }
  .panel--datatype { grid-column: 2 / 3; }
  .panel--health { grid-column: 1 / 3; }
  .panel--activity { grid-column: 1 / 3; }
}
@media (max-width: 700px) {
  .charts-grid { grid-template-columns: 1fr; }
  .panel--trend, .panel--status, .panel--project, .panel--datatype, .panel--health, .panel--activity { grid-column: 1 / 2; grid-row: auto; }
}
.panel { background: var(--color-surface); border-radius: 14px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.04); border: 1px solid var(--color-border-light); transition: box-shadow 0.2s ease; }
.panel:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.07); }
.panel__header { display: flex; justify-content: space-between; align-items: center; padding: 14px 18px; border-bottom: 1px solid var(--color-border-light); flex-wrap: wrap; gap: 8px; }
.panel__title { font-family: 'Manrope', 'Inter', sans-serif; font-size: 14px; font-weight: 600; color: var(--color-text-primary); -webkit-font-smoothing: antialiased; }
.panel__meta { font-size: 12px; color: var(--color-text-muted); }
.panel__link { font-size: 12px; color: var(--color-primary); text-decoration: none; font-weight: 500; }
.panel__link:hover { text-decoration: underline; }
.panel__tabs { display: flex; gap: 4px; }
.panel-tab { display: flex; align-items: center; gap: 5px; padding: 5px 10px; border: none; background: var(--color-surface-2); border-radius: 6px; font-size: 11px; font-weight: 600; color: var(--color-text-secondary); cursor: pointer; transition: all 0.15s ease; }
.panel-tab:hover { color: var(--color-text-primary); }
.panel-tab.is-active { background: var(--color-surface); color: var(--color-primary); box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
.tab-dot { width: 6px; height: 6px; border-radius: 50%; }
.panel__body { padding: 14px 18px; }
.panel__body--chart { padding: 10px 14px 14px; }
.panel__body--status { display: flex; flex-direction: column; align-items: center; gap: 14px; }
.panel__body--health { display: flex; align-items: center; gap: 24px; }
.panel__body--activity { padding: 0; max-height: 280px; overflow-y: auto; }
.chart-canvas { width: 100%; height: 220px; }
.chart-canvas--donut { height: 180px; width: 180px; }
.chart-canvas--bar { height: 200px; }
.status-legend { width: 100%; display: flex; flex-direction: column; gap: 8px; }
.legend-row { display: flex; align-items: center; gap: 8px; font-size: 12px; }
.legend-bar { width: 3px; height: 16px; border-radius: 2px; flex-shrink: 0; }
.legend-name { flex: 1; color: var(--color-text-secondary); font-weight: 500; }
.legend-count { font-weight: 700; color: var(--color-text-primary); font-family: 'Manrope', 'Inter', sans-serif; min-width: 28px; text-align: right; -webkit-font-smoothing: antialiased; }
.legend-pct { font-size: 11px; color: var(--color-text-muted); min-width: 32px; text-align: right; }
.health-gauge { flex-shrink: 0; }
.gauge-ring { position: relative; width: 110px; height: 110px; }
.gauge-ring svg { width: 100%; height: 100%; }
.gauge-center { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; }
.gauge-value { display: block; font-family: 'Manrope', 'Inter', sans-serif; font-size: 22px; font-weight: 700; color: var(--color-text-primary); line-height: 1; -webkit-font-smoothing: antialiased; }
.gauge-label { display: block; font-size: 10px; color: var(--color-text-muted); margin-top: 4px; }
.health-bars { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.health-bar-item__header { display: flex; justify-content: space-between; margin-bottom: 4px; }
.health-bar-item__name { font-size: 11px; color: var(--color-text-secondary); font-weight: 500; }
.health-bar-item__count { font-size: 11px; font-weight: 600; color: var(--color-text-primary); }
.health-bar-item__track { width: 100%; height: 6px; background: var(--color-surface-2); border-radius: 3px; overflow: hidden; }
.health-bar-item__fill { height: 100%; border-radius: 3px; transition: width 0.5s ease; min-width: 2px; }
.activity-table { width: 100%; }
.table-head { display: grid; grid-template-columns: 40px 1.2fr 80px 70px 64px 76px; gap: 8px; padding: 10px 16px; font-size: 10px; font-weight: 700; color: var(--color-text-secondary); border-bottom: 1px solid var(--color-border-light); background: var(--gray-50); }
.table-row { display: grid; grid-template-columns: 40px 1.2fr 80px 70px 64px 76px; gap: 8px; padding: 9px 16px; align-items: center; border-bottom: 1px solid var(--color-border-light); font-size: 12px; color: var(--color-text-primary); }
.table-row:last-child { border-bottom: none; }
.table-row:hover { background: rgba(0,94,235,0.03); }
.row-num { font-family: 'Manrope', sans-serif; font-size: 11px; font-weight: 600; color: var(--color-text-muted); }
.row-serial { font-family: 'Manrope', 'DM Sans', monospace; font-size: 11.5px; font-weight: 500; color: var(--color-text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.col-disk { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 11.5px; color: var(--color-text-secondary); }
.col-size { font-family: 'Manrope', sans-serif; font-size: 11.5px; font-weight: 600; color: var(--color-text-primary); }
.row-time { font-size: 11.5px; color: var(--color-text-muted); white-space: nowrap; }
.status-badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }
.status-badge.badge--success { background: var(--color-success-bg); color: var(--color-success-text); }
.status-badge.badge--warning { background: var(--color-warning-bg); color: var(--color-warning-text); }
.status-badge.badge--info { background: var(--color-info-bg); color: var(--color-info-text); }
.status-badge.badge--danger { background: var(--color-danger-bg); color: var(--color-danger-text); }
.table-empty { padding: 24px; text-align: center; color: var(--color-text-muted); font-size: 12px; }

/* ==================== 筛选抽屉样式 ==================== */
.filter-drawer-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-drawer-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.filter-drawer-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.filter-drawer-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filter-card {
  background: #fff;
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  overflow: hidden;
}

.filter-card-header {
  padding: 10px 14px;
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-page-bg);
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.filter-card-body {
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-presets {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
}

.preset-btn {
  padding: 6px 4px;
  border: 1px solid var(--color-border);
  background: transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.15s;

  &:hover {
    border-color: var(--color-primary);
    color: var(--color-primary);
  }

  &.is-active {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: #fff;
    font-weight: 600;
  }
}

.filter-drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 12px 16px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
}
</style>
