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
            <span v-if="searchDiskLabel" class="filter-tag">{{ searchDiskLabel }}<button @click="searchDiskLabel = ''; handleFilterChange()">×</button></span>
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
            <span class="panel__title">{{ t('uploadRecord.dashboard.dailyUpload') }}</span>
            <div class="disk-header-right">
              <div class="disk-time-group">
                <button
                  v-for="opt in trendTimeOptions"
                  :key="opt.key"
                  class="disk-time-btn"
                  :class="{ 'is-active': trendTimeKey === opt.key }"
                  @click="setTrendTime(opt.key)"
                >{{ opt.label }}</button>
                <button class="disk-time-btn disk-time-btn--custom" @click="showTrendCustom = true">
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                </button>
              </div>
              <div v-if="showTrendCustom" class="disk-custom-group">
                <el-date-picker
                  v-model="trendCustomRange"
                  type="daterange"
                  range-separator="~"
                  start-placeholder="开始"
                  end-placeholder="结束"
                  size="small"
                  format="MM-DD"
                  value-format="YYYY-MM-DD"
                  :clearable="true"
                  @change="handleTrendCustomChange"
                  class="disk-custom-picker"
                  :disabled-date="(date: Date) => date > new Date()"
                />
                <button class="disk-back-btn" @click="showTrendCustom = false; trendCustomRange = null; setTrendTime('today')">
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                </button>
              </div>
            </div>
          </div>
          <div class="panel__body panel__body--chart">
            <div ref="trendChartRef" class="chart-canvas"></div>
          </div>
        </div>

        <div class="panel panel--disk card-appear" style="animation-delay: 0.1s">
          <div class="panel__header">
            <span class="panel__title">{{ t('uploadRecord.dashboard.diskLabelStatus') }}</span>
            <div class="disk-header-right">
              <div class="disk-time-group">
                <button
                  v-for="opt in diskTimeOptions"
                  :key="opt.key"
                  class="disk-time-btn"
                  :class="{ 'is-active': diskTimeKey === opt.key }"
                  @click="setDiskTime(opt.key)"
                >{{ opt.label }}</button>
                <button class="disk-time-btn disk-time-btn--custom" @click="showDiskCustom = true">
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                </button>
              </div>
              <div v-if="showDiskCustom" class="disk-custom-group">
                <el-date-picker
                  v-model="diskCustomRange"
                  type="daterange"
                  range-separator="~"
                  start-placeholder="开始"
                  end-placeholder="结束"
                  size="small"
                  format="MM-DD"
                  value-format="YYYY-MM-DD"
                  :clearable="true"
                  @change="handleDiskCustomChange"
                  class="disk-custom-picker"
                  :disabled-date="(date: Date) => date > new Date()"
                />
                <button class="disk-back-btn" @click="showDiskCustom = false; diskCustomRange = null; setDiskTime('today')">
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                </button>
              </div>
              <span class="panel__meta">{{ diskLabelTotal }} {{ t('uploadRecord.dashboard.diskLabels') }}</span>
            </div>
          </div>
          <div class="panel__body panel__body--disk">
            <div v-if="diskLabelsLoading" class="disk-loading">
              <div class="disk-skeleton-grid">
                <div v-for="i in 8" :key="i" class="disk-skeleton"></div>
              </div>
            </div>
            <div v-else-if="diskLabels.length === 0" class="disk-empty">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#c4c9d4" stroke-width="1.5"><rect x="2" y="6" width="20" height="12" rx="2"/><line x1="6" y1="10" x2="6" y2="14"/><line x1="10" y1="10" x2="10" y2="14"/><line x1="14" y1="10" x2="14" y2="14"/></svg>
              <span>{{ t('uploadRecord.dashboard.noDiskLabels') }}</span>
            </div>
            <div v-else class="disk-grid">
              <div
                v-for="item in diskLabels"
                :key="item.diskLabel"
                class="disk-tile"
                :class="`disk-tile--${item.status}`"
                :title="`${item.diskLabel} - ${item.count} ${t('uploadRecord.dashboard.recordCount')} - ${t(getStatusText(item.status))}`"
              >
                <div class="disk-tile__inner">
                  <span class="disk-tile__name">{{ item.diskLabel }}</span>
                  <span class="disk-tile__count">{{ item.count }}</span>
                </div>
              </div>
            </div>
            <div class="disk-legend">
              <span class="disk-legend__item disk-legend__item--completed">
                <span class="disk-legend__dot"></span>{{ t('uploadRecord.dashboard.legendCompleted') }}
              </span>
              <span class="disk-legend__item disk-legend__item--failed">
                <span class="disk-legend__dot"></span>{{ t('uploadRecord.dashboard.legendFailed') }}
              </span>
              <span class="disk-legend__item disk-legend__item--mixed">
                <span class="disk-legend__dot"></span>{{ t('uploadRecord.dashboard.legendMixed') }}
              </span>
              <span class="disk-legend__item disk-legend__item--pending">
                <span class="disk-legend__dot"></span>{{ t('uploadRecord.dashboard.legendPending') }}
              </span>
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
            <span class="panel__title">{{ t('uploadRecord.dashboard.diskLabelDistribution') }}</span>
            <span class="panel__meta">{{ dataTypeTotal }} {{ t('uploadRecord.dashboard.diskLabelCount') }}</span>
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
                  <span class="col-disk" :title="item.diskLabel">{{ item.diskLabel }}</span>
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
              v-model="searchDiskLabel"
              :placeholder="t('uploadRecord.dashboard.filterAllDisk')"
              clearable
              filterable
              size="small"
              style="width: 100%"
              @change="handleFilterChange"
            >
              <el-option v-for="dt in diskLabelOptions" :key="dt" :label="dt" :value="dt" />
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
const activeTimePreset = ref<'today' | 'week' | 'month' | 'all'>('today')
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
const searchDiskLabel = ref('')
const searchStatus = ref('')
const searchUploader = ref('')
const projectList = ref<{ name: string; code: string }[]>([])
const diskLabelOptions = ref<string[]>([])
const uploaderOptions = ref<string[]>([])
const filterDrawerVisible = ref(false)

const hasActiveFilters = computed(() => !!(searchProjectName.value || searchDiskLabel.value || searchStatus.value || searchUploader.value || dateRange.value?.length))

const statistics = reactive<UploadRecordStatistics>({
  todayCount: 0, todaySize: 0, todaySizeStr: '0 B',
  weekCount: 0, weekSize: 0, weekSizeStr: '',
  monthCount: 0, monthSize: 0, monthSizeStr: '',
  totalCount: 0, totalSize: 0, totalSizeStr: '0 B',
  trend: [], byStatus: [], byDiskLabel: [], byProject: []
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
const dataTypeTotal = computed(() => (statistics.byDiskLabel || []).length)
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

let trendChart: any = null, dataTypeChart: any = null, projectChart: any = null

const trendChartRef = ref<HTMLElement>()
const dataTypeChartRef = ref<HTMLElement>()
const projectChartRef = ref<HTMLElement>()

// 每日上传统计日期范围（独立控制）
const trendTimeKey = ref('today')
const trendCustomRange = ref<[string, string] | null>(null)
const showTrendCustom = ref(false)
const trendTimeOptions = [
  { key: 'today', label: '今天' },
  { key: '3d', label: '近3天' },
  { key: '7d', label: '近7天' },
  { key: '30d', label: '近30天' },
]
const getTrendDateRange = (): [string, string] => {
  const today = new Date()
  const fmt = (d: Date) => d.toISOString().split('T')[0]
  switch (trendTimeKey.value) {
    case 'today': return [fmt(today), fmt(today)]
    case '3d': { const d = new Date(today); d.setDate(d.getDate() - 2); return [fmt(d), fmt(today)] }
    case '7d': { const d = new Date(today); d.setDate(d.getDate() - 6); return [fmt(d), fmt(today)] }
    case '30d': { const d = new Date(today); d.setDate(d.getDate() - 29); return [fmt(d), fmt(today)] }
    default: return ['', '']
  }
}
const loadTrendData = async () => {
  try {
    const params: Record<string, any> = {}
    if (trendCustomRange.value?.length === 2) {
      params.startDate = trendCustomRange.value[0]
      params.endDate = trendCustomRange.value[1]
    } else {
      const range = getTrendDateRange()
      if (range[0]) params.startDate = range[0]
      if (range[1]) params.endDate = range[1]
    }
    const res = await UploadRecordApi.statistics(params)
    const data = res.data || {}
    statistics.trend = data.trend || []
    statistics.byDiskLabel = data.byDiskLabel || []
    statistics.byProject = data.byProject || []
    updateAllCharts()
  } catch (error) { console.error('Failed to load trend data:', error) }
}
const setTrendTime = (key: string) => {
  trendTimeKey.value = key
  showTrendCustom.value = false
  trendCustomRange.value = null
  loadTrendData()
}
const handleTrendCustomChange = (val: any) => { if (val) { trendTimeKey.value = 'custom'; trendCustomRange.value = val; loadTrendData() } }

// 磁盘标签状态面板（独立时间控制）
const diskLabels = ref<{ diskLabel: string; count: number; status: 'completed' | 'failed' | 'mixed' | 'pending' }[]>([])
const diskLabelsLoading = ref(false)
const diskTimeKey = ref('today')
const diskCustomRange = ref<[string, string] | null>(null)
const showDiskCustom = ref(false)
const diskTimeOptions = [
  { key: 'today', label: '今天' },
  { key: '3d', label: '近3天' },
  { key: '7d', label: '近7天' },
  { key: '30d', label: '近30天' },
]
const getDiskDateRange = (): [string, string] => {
  const today = new Date()
  const fmt = (d: Date) => d.toISOString().split('T')[0]
  switch (diskTimeKey.value) {
    case 'today': return [fmt(today), fmt(today)]
    case '3d': { const d = new Date(today); d.setDate(d.getDate() - 2); return [fmt(d), fmt(today)] }
    case '7d': { const d = new Date(today); d.setDate(d.getDate() - 6); return [fmt(d), fmt(today)] }
    case '30d': { const d = new Date(today); d.setDate(d.getDate() - 29); return [fmt(d), fmt(today)] }
    default: return ['', '']
  }
}
const setDiskTime = (key: string) => { diskTimeKey.value = key; showDiskCustom.value = false; diskCustomRange.value = null; loadDiskLabels() }
const handleDiskCustomChange = (val: any) => { if (val) { diskTimeKey.value = 'custom'; loadDiskLabels() } }
const loadDiskLabels = async () => {
  diskLabelsLoading.value = true
  try {
    const [s, e] = getDiskDateRange()
    const params: Record<string, string> = {}
    if (s) params.startDate = s
    if (e) params.endDate = e
    const res = await UploadRecordApi.diskLabels(params)
    diskLabels.value = res.data || []
  } catch { diskLabels.value = [] }
  finally { diskLabelsLoading.value = false }
}
const diskLabelTotal = computed(() => diskLabels.value.length)

const initCharts = async () => {
  if (!trendChartRef.value) return
  echarts = await import('echarts')
  trendChart = echarts.init(trendChartRef.value)
  dataTypeChart = echarts.init(dataTypeChartRef.value)
  projectChart = echarts.init(projectChartRef.value)
  updateAllCharts()
  window.addEventListener('resize', handleResize)
}

const updateTrendChart = () => {
  if (!trendChart || !echarts) return
  const arr = trendDaily.value || []
  const dates = arr.map((t: any) => t.date)
  const counts = arr.map((t: any) => t.count)
  const sizes = arr.map((t: any) => t.totalSize)
  const maxCount = Math.max(...counts, 1)
  const maxSize = Math.max(...sizes, 1)
  trendChart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      textStyle: { color: '#191c23', fontSize: 12 },
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const d = arr[params[0].dataIndex]
        if (!d) return ''
        return `${params[0].name}<br/>${params[0].marker} 记录数：${d.count} 条<br/>${params[1].marker} 数据量：${formatFileSizeStatic(d.totalSize)}`
      }
    },
    legend: {
      data: ['记录数', '数据量'],
      top: 2,
      right: 8,
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { color: '#6b7280', fontSize: 11, fontFamily: 'Manrope, Inter, sans-serif', fontWeight: '500' }
    },
    grid: { top: 38, right: 56, bottom: 24, left: 48, containLabel: false },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: '#e5e7eb', width: 1 } },
      axisTick: { show: false },
      axisLabel: { color: '#9ca3af', fontSize: 10, fontFamily: 'Manrope, Inter, sans-serif', formatter: (v: string) => v.slice(5) },
      splitLine: { show: false }
    },
    yAxis: [
      {
        type: 'value',
        name: '',
        max: Math.ceil(maxCount * 1.15),
        show: true,
        position: 'left',
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: '#3b82f6', fontSize: 9, fontFamily: 'Manrope, Inter, sans-serif', formatter: (v: number) => v >= 1000 ? (v / 1000).toFixed(0) + 'k' : v },
        splitLine: { lineStyle: { color: '#f3f4f6', width: 1, type: 'dashed' } }
      },
      {
        type: 'value',
        name: '',
        max: Math.ceil(maxSize * 1.15),
        show: true,
        position: 'right',
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: '#22c55e', fontSize: 9, fontFamily: 'Manrope, Inter, sans-serif', formatter: (v: number) => formatFileSizeStatic(v) },
        splitLine: { show: false }
      }
    ],
    series: [
      {
        name: '记录数',
        type: 'bar',
        data: counts,
        barWidth: '30%',
        barGap: '20%',
        yAxisIndex: 0,
        itemStyle: { color: '#3b82f6', borderRadius: [0, 0, 0, 0] },
        tooltip: { formatter: (params: any) => `${params.name}<br/>记录数：${params.value} 条` }
      },
      {
        name: '数据量',
        type: 'bar',
        data: sizes,
        barWidth: '30%',
        barGap: '20%',
        yAxisIndex: 1,
        itemStyle: { color: '#22c55e', borderRadius: [0, 0, 0, 0] },
        tooltip: { formatter: (params: any) => `${params.name}<br/>数据量：${formatFileSizeStatic(params.value)}` }
      }
    ]
  })
}

const updateDataTypeChart = () => {
  if (!dataTypeChart || !echarts) return
  const all = statistics.byDiskLabel || []
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
        return `${d.diskLabel || '(空)'}<br/>${item.marker} ${d.count} 条<br/>数据量: ${formatFileSizeStatic(d.totalSize || 0)}`
      }
    },
    grid: { top: 10, right: 80, bottom: 10, left: 10, containLabel: true },
    xAxis: { type: 'value', show: false },
    yAxis: { type: 'category', data: reversed.map((d: any) => d.diskLabel || '(空)'), axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#414754', fontSize: 12 } },
    series: [{
      type: 'bar',
      data: reversed.map((d: any) => d.count),
      barWidth: 14,
      itemStyle: { borderRadius: [0, 6, 6, 0], color: '#3b82f6' },
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
      itemStyle: { borderRadius: [0, 6, 6, 0], color: '#8b5cf6' },
      label: { show: true, position: 'right', color: '#727785', fontSize: 11, formatter: (params: any) => formatFileSizeStatic(reversed[params.dataIndex]?.totalSize || 0) }
    }]
  })
}

const updateAllCharts = () => { updateTrendChart(); updateDataTypeChart(); updateProjectChart() }

const formatFileSizeStatic = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

const handleResize = () => { trendChart?.resize(); dataTypeChart?.resize(); projectChart?.resize() }

const loadStatistics = async () => {
  try {
    // 1. 顶部卡片展示全局汇总（不带筛选条件）
    const res = await UploadRecordApi.statistics()
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
    statistics.byStatus = data.byStatus || []

    const trend = data.trend || []
    if (trend.length >= 2) {
      const last = trend[trend.length - 1]
      const prev2 = trend[Math.max(0, trend.length - 2)]
      const base = prev2.count || 1
      const diff = last.count - prev2.count
      const pct = Math.round((Math.abs(diff) / base) * 100)
      todayTrend.value = diff >= 0 ? pct : -pct
    }

    // 2. 项目分布、磁盘标签分布跟随页面整体筛选
    const params: Record<string, any> = {}
    if (dateRange.value?.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    } else if (dateRange.value?.length === 1) {
      params.startDate = dateRange.value[0]
    }
    if (searchProjectName.value) params.projectName = searchProjectName.value
    if (searchDiskLabel.value) params.diskLabel = searchDiskLabel.value
    if (searchStatus.value) params.status = searchStatus.value
    if (searchUploader.value) params.uploader = searchUploader.value

    const chartRes = await UploadRecordApi.statistics(params)
    const chartData = chartRes.data || {}
    statistics.trend = chartData.trend || []
    statistics.byDiskLabel = chartData.byDiskLabel || []
    statistics.byProject = chartData.byProject || []

    updateAllCharts()
  } catch (error) { console.error('Failed to load statistics:', error) }
}

const loadRecentRecords = async () => {
  try {
    const params: any = { limit: 20 }
    if (searchProjectName.value) params.projectName = searchProjectName.value
    if (searchDiskLabel.value) params.diskLabel = searchDiskLabel.value
    if (searchStatus.value) params.status = searchStatus.value
    if (searchUploader.value) params.uploader = searchUploader.value
    const res = await UploadRecordApi.recent(params)
    recentRecords.value = res.data || []
  } catch (error) { console.error('Failed to load recent records:', error) }
}

const loadUploaderList = async () => { try { const res = await UploadRecordApi.uploaderList(); uploaderOptions.value = res.data || [] } catch { uploaderOptions.value = [] } }
const loadProjects = async () => { try { const res = await ProjectApi.getSimpleList(); projectList.value = res.data || [] } catch (error) { console.error('Failed to load projects:', error) } }
const handleFilterChange = () => { loadStatistics(); loadRecentRecords() }
const resetFilters = () => { searchProjectName.value = ''; searchDiskLabel.value = ''; searchStatus.value = ''; searchUploader.value = ''; setTimePreset('today'); loadStatistics(); loadRecentRecords() }

onMounted(() => {
  setTimeRange('month')
  initCharts()
  isInitialLoading.value = false
  loadProjects()
  loadUploaderList()
  loadDiskLabels()
  loadStatistics()
  loadTrendData()
  loadRecentRecords()
  refreshTimer = setInterval(() => {
    nextRefresh.value = 30
    const tick = setInterval(() => { nextRefresh.value--; if (nextRefresh.value <= 0) clearInterval(tick) }, 1000)
    loadDiskLabels()
    loadStatistics()
    loadTrendData()
    loadRecentRecords()
  }, 30000)
})

onUnmounted(() => { if (refreshTimer) clearInterval(refreshTimer); trendChart?.dispose(); dataTypeChart?.dispose(); projectChart?.dispose(); window.removeEventListener('resize', handleResize) })
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
.panel--disk { grid-column: 2 / 4; grid-row: 1 / 2; }
.panel--project { grid-column: 1 / 2; grid-row: 2 / 3; }
.panel--datatype { grid-column: 2 / 3; grid-row: 2 / 3; }
.panel--health { grid-column: 3 / 4; grid-row: 2 / 3; }
.panel--activity { grid-column: 1 / 4; grid-row: 3 / 4; }
@media (max-width: 1100px) {
  .charts-grid { grid-template-columns: 1fr 1fr; }
  .panel--trend { grid-column: 1 / 2; grid-row: 1 / 2; }
  .panel--disk { grid-column: 2 / 3; grid-row: 1 / 2; }
  .panel--project { grid-column: 1 / 2; grid-row: 2 / 3; }
  .panel--datatype { grid-column: 2 / 3; grid-row: 2 / 3; }
  .panel--health { grid-column: 1 / 3; }
  .panel--activity { grid-column: 1 / 3; }
}
@media (max-width: 700px) {
  .charts-grid { grid-template-columns: 1fr; }
  .panel--trend, .panel--disk, .panel--project, .panel--datatype, .panel--health, .panel--activity { grid-column: 1 / 2; grid-row: auto; }
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

/* 磁盘标签状态面板 */
.panel__body--disk { padding: 0; }
.disk-grid { display: flex; flex-wrap: wrap; gap: 8px; padding: 14px 16px; }
.disk-tile { width: 64px; height: 64px; border-radius: 10px; cursor: default; transition: transform 0.15s; }
.disk-tile:hover { transform: scale(1.05); }
.disk-tile__inner { width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; border-radius: 10px; padding: 4px; gap: 2px; }
.disk-tile--completed .disk-tile__inner { background: rgba(34,197,94,0.1); border: 1px solid rgba(34,197,94,0.3); }
.disk-tile--failed .disk-tile__inner { background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.3); }
.disk-tile--mixed .disk-tile__inner { background: rgba(245,158,11,0.1); border: 1px solid rgba(245,158,11,0.3); }
.disk-tile--pending .disk-tile__inner { background: rgba(156,163,175,0.1); border: 1px solid rgba(156,163,175,0.3); }
.disk-tile__name { font-size: 10px; font-weight: 600; color: var(--color-text-primary); text-align: center; word-break: break-all; line-height: 1.2; max-height: 28px; overflow: hidden; }
.disk-tile__count { font-size: 15px; font-weight: 700; color: var(--color-text-primary); font-family: 'Manrope', 'Inter', sans-serif; -webkit-font-smoothing: antialiased; }
.disk-legend { display: flex; gap: 12px; padding: 8px 16px; border-top: 1px solid var(--color-border-light); flex-wrap: wrap; }
.disk-legend__item { display: flex; align-items: center; gap: 5px; font-size: 11px; color: var(--color-text-secondary); font-weight: 500; }
.disk-legend__dot { width: 8px; height: 8px; border-radius: 3px; }
.disk-legend__item--completed .disk-legend__dot { background: #22c55e; }
.disk-legend__item--failed .disk-legend__dot { background: #ef4444; }
.disk-legend__item--mixed .disk-legend__dot { background: #f59e0b; }
.disk-legend__item--pending .disk-legend__dot { background: #9ca3af; }
.disk-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; padding: 40px 20px; color: var(--color-text-muted); font-size: 13px; }
.disk-loading { padding: 14px 16px; }
.disk-skeleton-grid { display: flex; flex-wrap: wrap; gap: 8px; }
.disk-skeleton { width: 64px; height: 64px; border-radius: 10px; background: linear-gradient(90deg, var(--gray-100) 25%, var(--gray-50) 50%, var(--gray-100) 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
.disk-header-right { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.disk-time-group {
  display: inline-flex;
  align-items: center;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 3px;
  gap: 1px;
}
.disk-time-btn {
  position: relative;
  padding: 5px 13px;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s;
  font-family: 'DM Sans', 'Manrope', sans-serif;
  white-space: nowrap;
  outline: none;
  -webkit-font-smoothing: antialiased;
}
.disk-time-btn:hover:not(.is-active) {
  color: #1e40af;
  background: #eff6ff;
}
.disk-time-btn.is-active {
  background: #1d4ed8;
  color: #ffffff;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(29, 78, 216, 0.3);
  border-radius: 6px;
}
/* 选中按钮下方小三角指示器 */
.disk-time-btn.is-active::after {
  content: '';
  position: absolute;
  bottom: -9px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid #1d4ed8;
}
.disk-time-btn--custom {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
}
.disk-custom-group { display: flex; align-items: center; gap: 6px; }
.disk-custom-picker { width: 200px; }
.disk-custom-picker :deep(.el-input__wrapper) {
  background: #ffffff !important;
  border: 1px solid #e5e7eb !important;
  border-radius: 6px !important;
  box-shadow: none !important;
  padding: 4px 8px !important;
  font-family: 'DM Sans', sans-serif !important;
  font-size: 12px !important;
}
.disk-back-btn {
  display: flex;
  align-items: center;
  padding: 5px 8px;
  border: none;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 12px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s;
  font-family: 'DM Sans', sans-serif;
}
.disk-back-btn:hover { background: #e5e7eb; color: #374151; }
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
  display: inline-flex;
  align-items: center;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 3px;
  gap: 1px;
}

.preset-btn {
  position: relative;
  padding: 5px 13px;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s;
  font-family: 'DM Sans', sans-serif;
  white-space: nowrap;
  -webkit-font-smoothing: antialiased;
}

.preset-btn:hover:not(.is-active) {
  color: #1e40af;
  background: #eff6ff;
}

.preset-btn.is-active {
  background: #1d4ed8;
  color: #ffffff;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(29, 78, 216, 0.3);
  border-radius: 6px;
}

.preset-btn.is-active::after {
  content: '';
  position: absolute;
  bottom: -9px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid #1d4ed8;
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
