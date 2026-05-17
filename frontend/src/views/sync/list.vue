<template>
  <div class="sync-page">
    <div class="page-container">

    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('sync.title') }}</h1>
        <span class="page-subtitle">{{ t('sync.subtitle') }}</span>
      </div>
      <div class="header-right">
        <el-button class="btn-flat" @click="loadSyncStatus" :loading="statusLoading">
          <el-icon><Refresh /></el-icon>
          {{ t('common.refresh') }}
        </el-button>
      </div>
    </header>

    <!-- 模式切换：紧凑胶囊式标签 -->
    <div class="mode-switcher">
      <button
        class="mode-pill"
        :class="{ 'is-active': activeTab === 'agent', 'mode-agent': true }"
        @click="switchTab('agent')"
      >
        <el-icon><Monitor /></el-icon>
        <span>{{ t('sync.tabAgent') }}</span>
        <span v-if="syncStatus.mode === 'agent'" class="pill-dot"></span>
      </button>
      <button
        class="mode-pill"
        :class="{ 'is-active': activeTab === 'center', 'mode-center': true }"
        @click="switchTab('center')"
      >
        <el-icon><OfficeBuilding /></el-icon>
        <span>{{ t('sync.tabCenter') }}</span>
        <span v-if="syncStatus.mode === 'center'" class="pill-dot"></span>
      </button>
    </div>

    <!-- 内容区 -->
    <div class="page-content">
      <!-- ==================== AGENT TAB ==================== -->
      <div v-show="activeTab === 'agent'" class="tab-content">
        <!-- 未配置 Agent -->
        <div v-if="!syncStatus.enabled || syncStatus.mode !== 'agent'" class="not-configured">
          <div class="not-configured-icon"><el-icon><Monitor /></el-icon></div>
          <div class="not-configured-title">{{ t('sync.agentNotConfigured') }}</div>
          <div class="not-configured-hint">{{ t('sync.agentNotConfiguredHint') }}</div>
        </div>

        <!-- Agent 已配置 -->
        <template v-else>
          <!-- 连接状态条 -->
          <div class="conn-status-bar" :class="connectionStatus">
            <div class="conn-left">
              <span class="conn-dot"></span>
              <span class="conn-label">{{ connectionStatusText }}</span>
            </div>
            <div class="conn-right">
              <span class="conn-tag" :class="syncStatus.registered ? 'registered' : 'unregistered'">
                {{ syncStatus.registered ? '已注册到 Center' : '未注册' }}
              </span>
            </div>
          </div>

          <!-- 基本信息 + 队列状态 合并行 -->
          <div class="agent-grid">
            <div class="info-card">
              <div class="info-card-title">{{ t('sync.basicInfo') }}</div>
              <div class="info-row"><span class="info-lbl">Agent ID</span><span class="info-val mono">#{{ syncStatus.stationId || '-' }}</span></div>
              <div class="info-row"><span class="info-lbl">Agent 名称</span><span class="info-val">{{ syncStatus.stationName || '-' }}</span></div>
              <div class="info-row"><span class="info-lbl">{{ t('sync.centerUrl') }}</span><a v-if="syncStatus.centerUrl" class="info-val url" :href="syncStatus.centerUrl" target="_blank" :title="syncStatus.centerUrl">{{ syncStatus.centerUrl }}</a><span v-else class="info-val muted">-</span></div>
              <div class="info-row"><span class="info-lbl">{{ t('sync.interval') }}</span><span class="info-val">{{ syncStatus.interval || '-' }}</span></div>
              <div class="info-row"><span class="info-lbl">{{ t('sync.batchSize') }}</span><span class="info-val">{{ syncStatus.batchSize || 0 }}</span></div>
              <div class="info-row"><span class="info-lbl">{{ t('sync.filterProjects') }}</span><span class="info-val">{{ (syncStatus.filter?.projectNames?.length) ? syncStatus.filter.projectNames.join(', ') : t('sync.allProjects') }}</span></div>
            </div>

            <!-- 队列状态 -->
            <div class="info-card" v-if="syncStatus.registered">
              <div class="info-card-title">{{ t('sync.queueStatus') }}</div>
              <div class="queue-grid">
                <div class="queue-item"><span class="qnum pending">{{ syncStatus.queuePending || 0 }}</span><span class="qlbl">{{ t('sync.queuePending') }}</span></div>
                <div class="queue-item"><span class="qnum completed">{{ syncStatus.queueCompleted || 0 }}</span><span class="qlbl">{{ t('sync.queueCompleted') }}</span></div>
                <div class="queue-item"><span class="qnum failed">{{ syncStatus.queueFailed || 0 }}</span><span class="qlbl">{{ t('sync.queueFailed') }}</span></div>
                <div class="queue-item"><span class="qnum total">{{ syncStatus.queueTotal || 0 }}</span><span class="qlbl">{{ t('sync.queueTotal') }}</span></div>
              </div>
            </div>
          </div>

          <!-- 最后同步 -->
          <div class="info-card" v-if="syncStatus.lastSyncAt">
            <div class="info-card-title">{{ t('sync.lastSync') }}</div>
            <div class="info-row"><span class="info-lbl">{{ t('sync.lastSync') }}</span><span class="info-val">{{ formatTime(syncStatus.lastSyncAt) }}</span></div>
            <div class="info-row" v-if="syncStatus.lastSerialNo"><span class="info-lbl">{{ t('sync.lastSerialNo') }}</span><span class="info-val mono">{{ syncStatus.lastSerialNo }}</span></div>
          </div>

          <!-- 错误 -->
          <div class="error-card" v-if="syncStatus.lastError">
            <el-icon><Warning /></el-icon><span>{{ syncStatus.lastError }}</span>
          </div>
        </template>
      </div>

      <!-- ==================== CENTER TAB ==================== -->
      <div v-show="activeTab === 'center'" class="tab-content">
        <!-- 未配置 Center -->
        <div v-if="!syncStatus.enabled || syncStatus.mode !== 'center'" class="not-configured">
          <div class="not-configured-icon"><el-icon><OfficeBuilding /></el-icon></div>
          <div class="not-configured-title">{{ t('sync.centerNotConfigured') }}</div>
          <div class="not-configured-hint">{{ t('sync.centerNotConfiguredHint') }}</div>
        </div>

        <!-- Center 已配置 -->
        <template v-else>
          <!-- 概览统计 -->
          <div class="overview-bar">
            <div class="overview-stat-card">
              <span class="ovnum">{{ stationList.length }}</span>
              <span class="ovlbl">{{ t('sync.totalStations') }}</span>
            </div>
            <div class="overview-stat-card">
              <span class="ovnum">{{ totalSyncRecords }}</span>
              <span class="ovlbl">{{ t('sync.totalSyncRecords') }}</span>
            </div>
            <div class="overview-stat-card stat-success">
              <span class="ovnum">{{ totalSyncSuccess }}</span>
              <span class="ovlbl">{{ t('sync.successCount') }}</span>
            </div>
            <div class="overview-stat-card stat-fail">
              <span class="ovnum">{{ totalSyncFail }}</span>
              <span class="ovlbl">{{ t('sync.failCount') }}</span>
            </div>
            <div class="overview-stat-card stat-conflict">
              <span class="ovnum">{{ totalSyncConflict }}</span>
              <span class="ovlbl">{{ t('sync.conflictCount') }}</span>
            </div>
          </div>

          <!-- 子标签：站点管理 | 同步历史 -->
          <div class="center-sub-tabs">
            <button
              class="sub-tab-btn"
              :class="{ active: centerSubTab === 'stations' }"
              @click="centerSubTab = 'stations'"
            >
              <el-icon><OfficeBuilding /></el-icon>
              {{ t('sync.stations') }}
              <span class="sub-tab-count">{{ filteredStations.length }}</span>
            </button>
            <button
              class="sub-tab-btn"
              :class="{ active: centerSubTab === 'history' }"
              @click="centerSubTab = 'history'; loadCenterData(); loadHistory()"
            >
              <el-icon><Clock /></el-icon>
              {{ t('sync.syncHistory') }}
              <span v-if="historyList.length > 0" class="sub-tab-count">{{ historyList.length }}</span>
            </button>
          </div>

          <!-- 站点管理面板 -->
          <div v-show="centerSubTab === 'stations'" class="sub-panel stations-panel">
            <div class="sub-panel-toolbar">
              <el-input v-model="stationKeyword" :placeholder="t('common.search') + '...'" clearable style="width: 260px">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-button class="btn-primary" @click="openStationDrawer('create')">
                <el-icon><Plus /></el-icon>
                {{ t('sync.createStation') }}
              </el-button>
            </div>

            <!-- 站点卡片 -->
            <div class="card-grid">
              <div v-for="station in filteredStations" :key="station.id" class="station-card">
                <div class="station-bar-color" :style="{ background: getStationColor(station.code) }"></div>
                <div class="station-body-inner">
                  <div class="station-top-row">
                    <div class="station-avatar" :style="{ background: getStationColor(station.code) }">
                      {{ station.name.charAt(0).toUpperCase() }}
                    </div>
                    <div class="station-meta-info">
                      <div class="station-card-name">{{ station.name }}</div>
                      <div class="station-card-code"><span class="code-label">编号</span>{{ station.code }}</div>
                    </div>
                    <div class="station-id-badge">#{{ station.id }}</div>
                    <div class="station-status-dot" :class="station.isConnected ? 'connected' : 'disconnected'" :title="station.isConnected ? '在线' : '离线'"></div>
                  </div>
                  <div class="station-info-rows">
                    <div class="info-row" v-if="station.url"><span class="info-lbl">Agent 地址</span><a class="info-val url" :href="station.url" target="_blank" :title="station.url">{{ station.url }}</a></div>
                    <div class="info-row"><span class="info-lbl">{{ t('sync.lastSync') }}</span><span class="info-val" :class="station.lastSyncAt ? '' : 'muted'">{{ formatTime(station.lastSyncAt) }}</span></div>
                    <div class="info-row"><span class="info-lbl">{{ t('sync.lastConnected') }}</span><span class="info-val" :class="station.lastConnectedAt ? '' : 'muted'">{{ formatTime(station.lastConnectedAt) }}</span></div>
                    <div class="info-row"><span class="info-lbl">最后心跳</span><span class="info-val" :class="station.isConnected ? 'connected' : 'disconnected'">{{ formatTime(station.lastHeartbeatAt) }}</span></div>
                    <div class="info-row"><span class="info-lbl">同步次数</span><span class="info-val count">{{ station.syncCount }}</span></div>
                  </div>
                  <div class="station-actions">
                    <el-button class="btn-icon" @click="testConnection(station)" :loading="testingConn === station.id" :type="station.isConnected ? 'success' : 'warning'" plain :title="station.isConnected ? 'Agent 在线' : 'Agent 离线'">
                      <el-icon><Connection /></el-icon>
                    </el-button>
                    <el-button class="btn-icon" @click="openStationDrawer('edit', station)"><el-icon><Edit /></el-icon></el-button>
                    <el-button class="btn-icon" @click="showApiKey(station)" plain :title="'查看 API Key'"><el-icon><Key /></el-icon></el-button>
                    <el-button class="btn-icon" @click="resetStationKey(station)" plain :title="'重新注册，生成新 Key'" :loading="resettingKey === station.id"><el-icon><RefreshRight /></el-icon></el-button>
                    <el-button class="btn-icon btn-icon-danger" @click="handleDeleteStation(station)" :title="'删除站点'"><el-icon><Delete /></el-icon></el-button>
                  </div>
                </div>
              </div>
              <div v-if="filteredStations.length === 0" class="empty-state">
                <div class="empty-icon-wrap"><el-icon class="empty-icon"><OfficeBuilding /></el-icon></div>
                <div class="empty-title">{{ t('sync.noStations') }}</div>
                <div class="empty-hint">{{ t('sync.createFirstStation') }}</div>
              </div>
            </div>
          </div>

          <!-- 同步历史面板（独立满宽） -->
          <div v-show="centerSubTab === 'history'" class="sub-panel history-panel">
            <div class="sub-panel-toolbar">
              <el-input v-model="historyKeyword" :placeholder="t('common.search') + '...'" clearable style="width: 200px">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="historyStationId" :placeholder="t('sync.selectStation')" clearable style="width: 180px">
                <el-option :label="t('sync.allStations')" :value="0" />
                <el-option v-for="s in stationList" :key="s.id" :label="s.name" :value="s.id" />
              </el-select>
              <el-date-picker v-model="historyDateRange" type="daterange" range-separator="-" :start-placeholder="t('common.startDate')" :end-placeholder="t('common.endDate')" value-format="YYYY-MM-DD" style="width: 260px" />
              <el-button class="btn-icon" @click="loadHistory" :title="'刷新'"><el-icon><Refresh /></el-icon></el-button>
            </div>

            <!-- 历史表格 -->
            <div class="history-table-wrap" v-loading="historyLoading">
              <table v-if="historyList.length > 0" class="history-table">
                <thead>
                  <tr>
                    <th>{{ t('common.createdAt') }}</th>
                    <th>{{ t('sync.stations') }}</th>
                    <th>{{ t('sync.direction') }}</th>
                    <th>{{ t('common.status') }}</th>
                    <th style="text-align:right">{{ t('sync.totalRecords') }}</th>
                    <th style="text-align:right">{{ t('sync.successCount') }}</th>
                    <th style="text-align:right">{{ t('sync.failCount') }}</th>
                    <th style="text-align:right">{{ t('sync.conflictCount') }}</th>
                    <th>{{ t('common.remark') }}</th>
                    <th>{{ t('common.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in historyList" :key="row.id" class="history-row">
                    <td class="td-time">{{ formatTime(row.createdAt) }}</td>
                    <td><span class="station-chip" :style="{ background: getStationColor(row.stationCode || row.stationName) }">{{ row.stationName || row.stationCode || '-' }}</span></td>
                    <td><span class="dir-tag">{{ row.directionText }}</span></td>
                    <td><span :class="['status-tag', row.status]">{{ row.statusText }}</span></td>
                    <td class="td-num">{{ row.totalRecords }}</td>
                    <td class="td-num ok">{{ row.successCount }}</td>
                    <td class="td-num" :class="{ err: row.failCount > 0 }">{{ row.failCount }}</td>
                    <td class="td-num" :class="{ warn: row.conflictCount > 0 }">{{ row.conflictCount }}</td>
                    <td class="td-remark">{{ row.remark || '-' }}</td>
                    <td><el-button class="btn-link" link @click="viewBatchDetail(row)">{{ t('sync.viewRecords') }}</el-button></td>
                  </tr>
                </tbody>
              </table>
              <div v-else class="history-empty">
                <div class="empty-icon-wrap"><el-icon class="empty-icon"><Document /></el-icon></div>
                <div class="empty-title">{{ t('sync.noSyncData') }}</div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 站点编辑抽屉 -->
    <el-drawer v-model="stationDrawerVisible" :title="stationDrawerTitle" size="520px" class="sync-drawer">
      <div class="drawer-content">
        <div class="drawer-tip">
          <el-icon><InfoFilled /></el-icon>
          <span v-if="stationEditMode === 'create'">{{ t('sync.createTip') }}</span>
          <span v-else>{{ t('sync.editTip') }}</span>
        </div>
        <el-form :model="stationForm" label-position="top">
          <div class="form-section">
            <div class="form-section-title">{{ t('sync.basicInfo') }}</div>
            <el-form-item label="Agent 名称" required>
              <el-input v-model="stationForm.name" placeholder="例如：上海节点" />
            </el-form-item>
            <el-form-item label="Agent 编号" required>
              <el-input v-model="stationForm.code" placeholder="唯一标识，如 sh-local1">
                <template #suffix><el-tooltip :content="t('sync.codeTip')" placement="top"><el-icon><QuestionFilled /></el-icon></el-tooltip></template>
              </el-input>
            </el-form-item>
            <el-form-item label="Agent 地址" required>
              <el-input v-model="stationForm.url" placeholder="https://agent.example.com:8080">
                <template #prefix><el-icon><Link /></el-icon></template>
              </el-input>
            </el-form-item>
          </div>
          <div class="form-section">
            <div class="form-section-title">{{ t('sync.stationStatus') }}</div>
            <el-form-item :label="t('sync.stationStatus')">
              <el-radio-group v-model="stationForm.status">
                <el-radio value="active">{{ t('sync.enabled') }}</el-radio>
                <el-radio value="inactive">{{ t('sync.disabled') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </div>
          <div class="form-section">
            <div class="form-section-title">{{ t('common.remark') }}</div>
            <el-form-item :label="t('sync.stationDescription')">
              <el-input v-model="stationForm.description" type="textarea" :rows="3" />
            </el-form-item>
          </div>
        </el-form>
      </div>
      <template #footer>
        <div class="drawer-footer">
          <el-button class="btn-flat" @click="stationDrawerVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button class="btn-primary" @click="saveStation" :loading="stationSaving">{{ t('common.save') }}</el-button>
        </div>
      </template>
    </el-drawer>

    <!-- API Key 弹窗 -->
    <el-dialog v-model="apiKeyDialogVisible" :title="t('sync.apiKey')" width="480px" class="apikey-dialog">
      <div class="apikey-content" v-if="currentStationForKey">
        <div class="apikey-station">
          <div class="station-avatar-sm" :style="{ background: getStationColor(currentStationForKey.code) }">
            {{ currentStationForKey.name.charAt(0).toUpperCase() }}
          </div>
          <div>
            <div class="apikey-station-name">{{ currentStationForKey.name }}</div>
            <div class="apikey-station-code">{{ currentStationForKey.code }}</div>
          </div>
        </div>
        <div class="apikey-tip"><el-icon><Warning /></el-icon><span>{{ t('sync.apiKeyWarning') }}</span></div>
        <div class="apikey-display">
          <el-input v-model="displayApiKey" readonly />
          <el-button class="btn-primary" @click="copyApiKey" :disabled="!displayApiKey"><el-icon><DocumentCopy /></el-icon></el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 同步详情抽屉 -->
    <el-drawer v-model="stationHistoryVisible" :title="currentStationName + ' - ' + t('sync.recordDetails')" size="880px" direction="rtl">
      <div class="station-history-drawer">
        <!-- 批次筛选 -->
        <div class="batch-toolbar">
          <el-date-picker v-model="batchDateRange" type="daterange" range-separator="-" :start-placeholder="t('common.startDate')" :end-placeholder="t('common.endDate')" value-format="YYYY-MM-DD" style="width: 240px" @change="loadStationBatches" />
          <el-select v-model="detailResultFilter" style="width: 130px" @change="loadBatchDetails">
            <el-option :label="t('common.all')" value="" />
            <el-option label="Success" value="success" />
            <el-option label="Failed" value="failed" />
            <el-option label="Conflict" value="conflict" />
          </el-select>
          <el-button class="btn-icon" @click="loadStationBatches" :title="'刷新'"><el-icon><Refresh /></el-icon></el-button>
        </div>

        <!-- 批次列表 -->
        <div class="batch-section">
          <div class="section-header">
            <span class="section-title-sm">{{ t('sync.batchList') }}</span>
          </div>
          <el-table :data="batchList" stripe v-loading="batchLoading" @row-click="viewBatchDetail">
            <el-table-column :label="t('common.createdAt')" width="160"><template #default="{ row }">{{ formatTime(row.createdAt) }}</template></el-table-column>
            <el-table-column prop="directionText" :label="t('sync.direction')" width="90" align="center" />
            <el-table-column :label="t('common.status')" width="110" align="center"><template #default="{ row }"><span :class="['status-tag', row.status]">{{ row.statusText }}</span></template></el-table-column>
            <el-table-column :label="t('sync.totalRecords')" width="90" align="center"><template #default="{ row }"><span class="num">{{ row.totalRecords }}</span></template></el-table-column>
            <el-table-column :label="t('sync.successCount')" width="90" align="center"><template #default="{ row }"><span class="num ok">{{ row.successCount }}</span></template></el-table-column>
            <el-table-column :label="t('sync.failCount')" width="90" align="center"><template #default="{ row }"><span class="num err">{{ row.failCount }}</span></template></el-table-column>
            <el-table-column :label="t('sync.conflictCount')" width="90" align="center"><template #default="{ row }"><span class="num warn">{{ row.conflictCount }}</span></template></el-table-column>
            <el-table-column :label="t('common.actions')" width="90" align="center" fixed="right"><template #default="{ row }"><el-button class="btn-link" link type="primary" @click.stop="viewBatchDetail(row)">{{ t('sync.viewRecords') }}</el-button></template></el-table-column>
          </el-table>
          <div class="pagination-wrapper">
            <el-pagination v-model:current-page="batchPage" v-model:page-size="batchPageSize" :total="batchTotal" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @change="loadStationBatches" />
          </div>
        </div>

        <!-- 详情列表 -->
        <div v-if="activeBatchId" class="detail-section">
          <div class="section-header">
            <span class="section-title-sm">{{ t('sync.recordDetails') }}</span>
          </div>
          <el-table :data="detailList" stripe max-height="400">
            <el-table-column prop="serialNo" :label="t('sync.originalSerialNo')" min-width="140" />
            <el-table-column prop="projectName" :label="t('sync.projectName')" min-width="100" />
            <el-table-column :label="t('sync.action')" width="80" align="center"><template #default="{ row }"><span :class="['action-tag', row.action]">{{ row.actionText }}</span></template></el-table-column>
            <el-table-column :label="t('sync.result')" width="80" align="center"><template #default="{ row }"><span :class="['result-tag', row.result]">{{ row.result }}</span></template></el-table-column>
            <el-table-column :label="t('sync.errorReason')" min-width="160"><template #default="{ row }"><span v-if="row.errorMsg" class="error-msg">{{ row.errorMsg }}</span><span v-else class="none">-</span></template></el-table-column>
            <el-table-column :label="t('sync.newSerialNo')" min-width="140"><template #default="{ row }"><span v-if="row.newSerialNo" class="new-serial">{{ row.newSerialNo }}</span><span v-else class="none">-</span></template></el-table-column>
            <el-table-column :label="t('common.createdAt')" width="140"><template #default="{ row }">{{ formatTime(row.createdAt) }}</template></el-table-column>
          </el-table>
          <div class="pagination-wrapper">
            <el-pagination v-model:current-page="detailPage" v-model:page-size="detailPageSize" :total="detailTotal" :page-sizes="[20, 50, 100]" layout="total, sizes, prev, pager, next" @change="loadBatchDetails" />
          </div>
        </div>
      </div>
    </el-drawer>
  </div><!-- end page-container -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Refresh, Monitor, OfficeBuilding, Link, Clock, Guide, Edit, Delete, Key, QuestionFilled, InfoFilled, Warning, DocumentCopy, Right, Document, Connection, RefreshRight } from '@element-plus/icons-vue'
import { syncApi, type SyncStation, type SyncHistory, type SyncStationSummary } from '@/api/sync'

const { t } = useI18n()

// ==================== 通用 ====================
const statusLoading = ref(false)
const syncStatus = ref<any>({})

const activeTab = ref('agent')
const centerSubTab = ref('stations')

const switchTab = (tab: 'agent' | 'center') => {
  activeTab.value = tab
  if (tab === 'center') loadCenterData()
}

// formatTime handles Unix ms timestamps (number), ISO strings, and null
const formatTime = (time: string | null | number | undefined) => {
  if (!time) return '-'
  if (typeof time === 'number' || (typeof time === 'string' && /^\d{13}$/.test(time))) {
    const d = new Date(typeof time === 'number' ? time : Number(time))
    if (!isNaN(d.getTime()) && d.getTime() > 0) {
      // 使用本地时区格式化，避免 toISOString() 转为 UTC 少 8 小时
      return d.toLocaleString('zh-CN', { hour12: false }).replace(/\//g, '-')
    }
    return '-'
  }
  if (typeof time === 'string') {
    const clean = time.replace('T', ' ').slice(0, 16)
    if (clean !== '1969-12-31 23:59' && clean !== '1970-01-01 00:00') return clean
  }
  return '-'
}

const getStationColor = (code: string) => {
  const colors = ['#5b7b8c', '#4a7c59', '#8b6f5b', '#6b5b95', '#b87333', '#4a6fa5', '#7a8471', '#a67c52']
  let hash = 0
  for (let i = 0; i < code.length; i++) hash = ((hash << 5) - hash) + code.charCodeAt(i)
  return colors[Math.abs(hash) % colors.length]
}

const parseHeartbeatInterval = (interval: string) => {
  if (!interval) return '60秒'
  return interval.replace('s', '秒').replace('m', '分钟')
}

const loadSyncStatus = async () => {
  statusLoading.value = true
  try {
    const res = await syncApi.getStatus()
    syncStatus.value = res.data || {}
    if (res.data?.enabled) {
      activeTab.value = res.data?.mode === 'agent' ? 'agent' : 'center'
      if (activeTab.value === 'center') loadCenterData()
    }
  } catch {
    syncStatus.value = {}
  } finally {
    statusLoading.value = false
  }
}

// ==================== Agent ====================
const connectionStatus = computed(() => {
  if (!syncStatus.value.enabled || syncStatus.value.mode !== 'agent') return 'offline'
  return syncStatus.value.registered ? 'connected' : 'unregistered'
})
const connectionStatusText = computed(() => {
  const map: Record<string, string> = {
    connected: t('sync.statusConnected'),
    unregistered: t('sync.statusUnregistered'),
    offline: t('sync.statusOffline')
  }
  return map[connectionStatus.value] || t('sync.statusOffline')
})

// ==================== Center ====================
const stationList = ref<SyncStation[]>([])
const stationKeyword = ref('')
let stationKeywordTimer: ReturnType<typeof setTimeout> | null = null
const testingConn = ref<number | null>(null)
const resettingKey = ref<number | null>(null)
const stationDrawerVisible = ref(false)
const stationDrawerTitle = ref('')
const stationForm = ref({ id: 0, name: '', code: '', url: '', status: 'active', description: '' })
const stationSaving = ref(false)
const stationEditMode = ref<'create' | 'edit'>('create')
const apiKeyDialogVisible = ref(false)
const currentStationForKey = ref<SyncStation | null>(null)
const displayApiKey = ref('')

// 同步历史相关
const historyList = ref<SyncHistory[]>([])
const historyKeyword = ref('')
const historyStationId = ref<number>(0)
const historyDateRange = ref<[string, string] | null>(null)
const historyPage = ref(1)
const historyPageSize = ref(10)
const historyTotal = ref(0)
const historyLoading = ref(false)

// 监听筛选条件变化，自动重新加载
watch([historyStationId, historyKeyword, historyDateRange], () => {
  historyPage.value = 1
  loadHistory()
})

// 站点搜索（防抖 300ms）
let stationKeywordTimer2: ReturnType<typeof setTimeout> | null = null
watch(stationKeyword, () => {
  if (stationKeywordTimer2) clearTimeout(stationKeywordTimer2)
  stationKeywordTimer2 = setTimeout(() => {
    loadCenterData()
  }, 300)
})

const stationHistoryVisible = ref(false)
const currentStationName = ref('')
const batchDateRange = ref<[string, string] | null>(null)
const batchList = ref<SyncHistory[]>([])
const batchLoading = ref(false)
const batchPage = ref(1)
const batchPageSize = ref(10)
const batchTotal = ref(0)
const activeBatchId = ref<number | null>(null)
const detailResultFilter = ref('')
const detailList = ref<any[]>([])
const detailPage = ref(1)
const detailPageSize = ref(50)
const detailTotal = ref(0)

const totalSyncRecords = computed(() => stationSummaries.value.reduce((sum, s) => sum + s.totalRecords, 0))
const totalSyncSuccess = computed(() => stationSummaries.value.reduce((sum, s) => sum + s.successCount, 0))
const totalSyncFail = computed(() => stationSummaries.value.reduce((sum, s) => sum + s.failCount, 0))
const totalSyncConflict = computed(() => stationSummaries.value.reduce((sum, s) => sum + s.conflictCount, 0))
const stationSummaries = ref<SyncStationSummary[]>([])

const filteredStations = computed(() => {
  if (!stationKeyword.value) return stationList.value
  const kw = stationKeyword.value.toLowerCase()
  return stationList.value.filter(s => s.name.toLowerCase().includes(kw) || s.code.toLowerCase().includes(kw))
})

// Helper maps for station lookup
const stationMap = computed(() => {
  const m = new Map<number, SyncStation>()
  stationList.value.forEach(s => m.set(s.id, s))
  return m
})
const stationSummariesMap = computed(() => {
  const m = new Map<number, SyncStationSummary>()
  stationSummaries.value.forEach(s => m.set(s.id, s))
  return m
})
const getStationCode = (id: number) => stationMap.value.get(id)?.code || ''
const getStationName = (id: number) => stationMap.value.get(id)?.name || stationSummariesMap.value.get(id)?.name || `#${id}`

const loadCenterData = async () => {
  try {
    const [stationsRes, summariesRes] = await Promise.all([
      syncApi.listStations({ page: 1, pageSize: 100 }),
      syncApi.getStationSummaries({})
    ])
    stationList.value = stationsRes.data.items
    stationSummaries.value = summariesRes.data || []
  } catch (e) { console.error(e) }
}

// 加载同步历史（表格模式，支持搜索+日期过滤）
const loadHistory = async () => {
  historyLoading.value = true
  try {
    const res = await syncApi.getHistory({
      page: historyPage.value,
      pageSize: historyPageSize.value,
      stationId: historyStationId.value || undefined,
      startDate: historyDateRange.value?.[0],
      endDate: historyDateRange.value?.[1],
      keyword: historyKeyword.value || undefined,
    })
    historyList.value = res.data.items || []
    historyTotal.value = res.data.total || 0
    // 同时刷新站点汇总数据（概览栏统计）
    await loadStationSummaries()
  } catch (e: any) {
    console.error('[loadHistory] failed:', e)
  } finally { historyLoading.value = false }
}

const testConnection = async (station: SyncStation) => {
  testingConn.value = station.id
  try {
    await loadCenterData()
    const updated = stationList.value.find(s => s.id === station.id)
    if (updated) {
      if (updated.isConnected) {
        ElMessage.success(`Agent "${station.name}" 在线`)
      } else {
        ElMessage.error(`Agent "${station.name}" 离线`)
      }
    }
  } finally {
    testingConn.value = null
  }
}

const loadStationSummaries = async () => {
  try {
    const res = await syncApi.getStationSummaries({
      stationId: historyStationId.value || undefined,
      startDate: historyDateRange.value?.[0],
      endDate: historyDateRange.value?.[1],
    })
    stationSummaries.value = res.data || []
  } catch (e) { console.error(e) }
}

const openStationHistory = async (station: SyncStationSummary) => {
  currentStationName.value = station.name
  batchPage.value = 1
  batchDateRange.value = null
  activeBatchId.value = null
  stationHistoryVisible.value = true
  await loadStationBatches()
}

const loadStationBatches = async () => {
  if (!currentStationSummary.value) return
  batchLoading.value = true
  try {
    const res = await syncApi.getHistory({
      page: batchPage.value, pageSize: batchPageSize.value,
      stationId: currentStationSummary.value.id,
      startDate: batchDateRange.value?.[0], endDate: batchDateRange.value?.[1],
    })
    batchList.value = res.data.items
    batchTotal.value = res.data.total
  } catch (e) { console.error(e) }
  finally { batchLoading.value = false }
}

const currentStationSummary = ref<SyncStationSummary | null>(null)

const viewBatchDetail = async (row: SyncHistory) => {
  if (activeBatchId.value === row.id) {
    activeBatchId.value = null
    detailList.value = []
    detailTotal.value = 0
    return
  }
  currentStationSummary.value = stationSummariesMap.value.get(row.stationId) || null
  currentStationName.value = getStationName(row.stationId)
  activeBatchId.value = row.id
  detailPage.value = 1
  detailResultFilter.value = ''
  stationHistoryVisible.value = true
  await loadBatchDetails()
}

const loadBatchDetails = async () => {
  if (!activeBatchId.value) return
  try {
    const res = await syncApi.getHistoryDetails(activeBatchId.value, detailResultFilter.value || undefined, detailPage.value, detailPageSize.value)
    detailList.value = res.data.details || []
    detailTotal.value = res.data.total || 0
  } catch (e) { console.error(e) }
}

const openStationDrawer = (mode: 'create' | 'edit', station?: SyncStation) => {
  stationEditMode.value = mode
  stationDrawerTitle.value = mode === 'create' ? t('sync.createStation') : t('sync.editStation')
  if (mode === 'edit' && station) {
    stationForm.value = { id: station.id, name: station.name, code: station.code, url: station.url, status: station.status, description: station.description || '' }
  } else {
    stationForm.value = { id: 0, name: '', code: '', url: '', status: 'active', description: '' }
  }
  stationDrawerVisible.value = true
}

const saveStation = async () => {
  if (!stationForm.value.name) { ElMessage.warning(t('sync.messages.nameRequired')); return }
  if (!stationForm.value.code) { ElMessage.warning(t('sync.messages.codeRequired')); return }
  if (!stationForm.value.url) { ElMessage.warning(t('sync.messages.urlRequired')); return }
  stationSaving.value = true
  try {
    if (stationEditMode.value === 'create') {
      const res = await syncApi.createStation(stationForm.value)
      ElMessage.success(t('sync.messages.createSuccess'))
      currentStationForKey.value = res.data
      displayApiKey.value = res.data.apiKey || t('sync.noApiKey')
      apiKeyDialogVisible.value = true
    } else {
      await syncApi.updateStation(stationForm.value)
      ElMessage.success(t('sync.messages.updateSuccess'))
    }
    stationDrawerVisible.value = false
    loadCenterData()
  } catch (e: any) { ElMessage.error(e.message || t('sync.messages.createFailed')) }
  finally { stationSaving.value = false }
}

const handleDeleteStation = async (station: SyncStation) => {
  try {
    await syncApi.deleteStation(station.id)
    ElMessage.success(t('sync.messages.deleteSuccess'))
    loadCenterData()
  } catch (e: any) { ElMessage.error(t('sync.messages.deleteFailed')) }
}

const showApiKey = (station: SyncStation) => {
  currentStationForKey.value = station
  displayApiKey.value = station.apiKey || t('sync.noApiKey')
  apiKeyDialogVisible.value = true
}

const copyApiKey = async () => {
  try {
    await navigator.clipboard.writeText(displayApiKey.value)
    ElMessage.success(t('sync.apiKeyCopied'))
  } catch { ElMessage.error('复制失败') }
}

const resetStationKey = async (station: SyncStation) => {
  try {
    await ElMessageBox.confirm(
      `将为站点「${station.name}」重新生成 API Key，Agent 重新同步时会自动用新 Key 重新注册。`,
      '重新注册',
      { confirmButtonText: '确定重置', cancelButtonText: '取消', type: 'warning' }
    )
  } catch { return }

  resettingKey.value = station.id
  try {
    const res = await syncApi.resetApiKey(station.id)
    currentStationForKey.value = station
    displayApiKey.value = res.data?.data?.apiKey || res.data?.apiKey || ''
    apiKeyDialogVisible.value = true
    ElMessage.success('新 API Key 已生成，请告知 Agent 重新同步')
  } catch (e: any) {
    console.error('[resetStationKey] 错误:', e)
    const msg = e?.response?.data?.message || e?.message || '重置失败，请刷新页面后重试'
    ElMessage.error(msg)
  } finally {
    resettingKey.value = null
  }
}

onMounted(() => { loadSyncStatus() })
</script>

<style scoped lang="scss">
// ========= 设计体系 =========
// 字号: 12辅助 / 13正文 / 14副标题 / 15卡片标题 / 22大标题 / 26数字
// 颜色: #1c1917正文 / #57534e次要 / #78716c辅助 / #a8a29e禁用
// max-width: 1280px 居中
// 原则: 无渐变，纯色背景，扁平设计

.sync-page { padding: 20px 0; min-height: 100vh; background: #f5f4f0; }
.page-container { max-width: 1280px; margin: 0 auto; padding: 0 20px; }

.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { font-size: 22px; font-weight: 700; color: #1c1917; margin: 0; }
.page-subtitle { font-size: 13px; color: #78716c; display: block; margin-top: 3px; }
.header-right { display: flex; gap: 8px; }

/* ============ 按钮体系（扁平无渐变） ============ */
.btn-primary {
  background: #409eff !important;
  border-color: #409eff !important;
  color: #fff !important;
  border-radius: 4px !important;
  &:hover { background: #66b1ff !important; border-color: #66b1ff !important; }
}
.btn-flat {
  background: #fff !important;
  border: 1px solid #dcdfe6 !important;
  color: #606266 !important;
  border-radius: 4px !important;
  &:hover { border-color: #409eff !important; color: #409eff !important; background: #ecf5ff !important; }
}
.btn-icon {
  background: #fff !important;
  border: 1px solid #e8e4de !important;
  color: #57534e !important;
  border-radius: 4px !important;
  &:hover { border-color: #409eff !important; color: #409eff !important; background: #ecf5ff !important; }
}
.btn-icon-danger {
  color: #dc2626 !important;
  &:hover { border-color: #dc2626 !important; background: #fef2f2 !important; color: #dc2626 !important; }
}
.btn-link { color: #409eff !important; &:hover { color: #66b1ff !important; } }

/* 模式切换胶囊 */
.mode-switcher {
  display: flex;
  gap: 6px;
  padding: 14px 16px 14px;
}

.mode-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  border-radius: 8px;
  border: 1.5px solid #e8e4de;
  background: #fff;
  font-size: 13px;
  font-weight: 600;
  color: #78716c;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  .el-icon { font-size: 14px; transition: color 0.2s; }

  .pill-dot {
    position: absolute;
    top: 6px;
    right: 6px;
    width: 7px;
    height: 7px;
    border-radius: 50%;
  }

  &:hover { border-color: #d4cfc6; color: #57534e; }

  /* Agent 激活态 */
  &.mode-agent.is-active {
    background: #fff7ed;
    border-color: #ea580c;
    color: #c2410c;
    box-shadow: 0 2px 8px rgba(234, 88, 12, 0.12);
    .pill-dot { background: #ea580c; }
    .el-icon { color: #ea580c; }
  }

  /* Center 激活态 */
  &.mode-center.is-active {
    background: #eff6ff;
    border-color: #2563eb;
    color: #1d4ed8;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.12);
    .pill-dot { background: #2563eb; }
    .el-icon { color: #2563eb; }
  }
}

/* Center 子标签页 */
.center-sub-tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid #e8e4de;
  border-top: 1px solid #e8e4de;
  padding: 4px 16px 0;
}

.sub-tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  font-size: 13px;
  font-weight: 600;
  color: #78716c;
  cursor: pointer;
  transition: all 0.2s;
  .el-icon { font-size: 14px; }
  &:hover { color: #1c1917; }
  &.active {
    color: #1c1917;
    border-bottom-color: #1c1917;
  }
}

.sub-tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 18px;
  padding: 0 6px;
  background: #f0ede8;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 700;
  color: #57534e;
  line-height: 1;
  margin-left: 2px;
}

.sub-panel {
  width: 100%;
}

.sub-panel-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

/* 旧 bento 布局移除 */
.center-bento { display: none; }
.bento-left { display: none; }
.bento-right { display: none; }
.bento-panel { display: none; }
.bento-panel-header { display: none; }
.panel-title { display: none; }
.panel-count { display: none; }
.station-scroll-area { display: none; }
.history-toolbar { display: none; }
.history-scroll-area { display: none; }

/* 内容区 */
.page-content { background: #fff; border-radius: 16px; border: 1px solid #e8e4de; overflow: hidden; }
.tab-content { padding: 0; }

/* 未配置提示 */
.not-configured { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; gap: 10px; }
.not-configured-icon { width: 56px; height: 56px; border-radius: 50%; background: #f5f4f0; display: flex; align-items: center; justify-content: center; .el-icon { font-size: 26px; color: #c8c3bb; } }
.not-configured-title { font-size: 14px; font-weight: 600; color: #78716c; }
.not-configured-hint { font-size: 12px; color: #a8a29e; text-align: center; max-width: 280px; }

/* 通用区块标题 */
.section-title { font-size: 14px; font-weight: 600; color: #1c1917; }
.section-title-sm { font-size: 13px; font-weight: 600; color: #57534e; }
.section-block { margin-bottom: 20px; }
.section-block-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; gap: 12px; }
.toolbar-right { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }

/* ==================== Agent ==================== */
.agent-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 12px; }

// 连接状态条
.conn-status-bar { display: flex; align-items: center; justify-content: space-between; padding: 12px 18px; border-radius: 10px; border: 1px solid; margin-bottom: 16px;
  &.connected { background: #f0fdf4; border-color: #bbf7d0; .conn-dot { background: #4a7c59; } .conn-label { color: #4a7c59; } }
  &.unregistered { background: #fffbeb; border-color: #fde68a; .conn-dot { background: #b87333; } .conn-label { color: #b87333; } }
  &.offline { background: #f5f4f0; border-color: #e8e4de; .conn-dot { background: #a8a29e; } .conn-label { color: #78716c; } }
}
.conn-left { display: flex; align-items: center; gap: 10px; }
.conn-dot { width: 9px; height: 9px; border-radius: 50%; }
.conn-label { font-size: 14px; font-weight: 600; }
.conn-tag { font-size: 12px; padding: 3px 10px; border-radius: 5px; font-weight: 500;
  &.registered { background: #dcfce7; color: #4a7c59; }
  &.unregistered { background: #fef3c7; color: #b87333; }
}

// 信息卡片
.info-card { background: #fafaf9; border: 1px solid #e8e4de; border-radius: 12px; padding: 16px; }
.info-card-title { font-size: 11px; font-weight: 600; color: #a8a29e; text-transform: uppercase; letter-spacing: 0.6px; margin-bottom: 12px; }
.info-row { display: flex; align-items: center; justify-content: space-between; padding: 9px 0; border-bottom: 1px solid #f0ede8; &:last-child { border-bottom: none; padding-bottom: 0; } &:first-child { padding-top: 0; } }
.info-lbl { font-size: 13px; color: #78716c; min-width: 100px; }
.info-val { font-size: 13px; color: #1c1917; text-align: right; &.mono { font-family: 'Monaco', monospace; color: #57534e; } &.url { max-width: 280px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #409eff; text-decoration: none; &:hover { text-decoration: underline; } } &.count { font-variant-numeric: tabular-nums; } &.connected { color: #4a7c59; } &.disconnected { color: #ef4444; } &.muted { color: #d4cfc6; } }

// 队列状态
.queue-card { background: #fafaf9; border: 1px solid #e8e4de; border-radius: 12px; padding: 16px; }
.queue-title { font-size: 11px; font-weight: 600; color: #a8a29e; text-transform: uppercase; letter-spacing: 0.6px; margin-bottom: 12px; }
.queue-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; }
.queue-item { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 14px 10px; background: #fff; border: 1px solid #e8e4de; border-radius: 8px; }
.qnum { font-size: 22px; font-weight: 700; color: #1c1917; line-height: 1;
  &.pending { color: #b87333; }
  &.completed { color: #4a7c59; }
  &.failed { color: #dc2626; }
  &.total { color: #409eff; }
}
.qlbl { font-size: 11px; color: #78716c; }

// 错误提示
.error-card { background: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 12px 14px; display: flex; align-items: flex-start; gap: 8px; font-size: 12px; color: #dc2626; margin-bottom: 16px; .el-icon { font-size: 14px; flex-shrink: 0; margin-top: 1px; } }

/* ==================== Center ==================== */
// 概览统计条
.overview-bar { display: grid; grid-template-columns: repeat(5, 1fr); gap: 10px; padding: 14px 16px; background: #fafaf9; border-bottom: 1px solid #f0ede8; }
.overview-stat-card { display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 12px 10px; background: #fff; border: 1px solid #e8e4de; border-radius: 10px;
  .ovnum { font-size: 22px; font-weight: 700; color: #1c1917; line-height: 1; }
  .ovlbl { font-size: 11px; color: #78716c; }
  &.stat-success { border-color: #bbf7d0; background: #f0fdf4; .ovnum { color: #4a7c59; } .ovlbl { color: #4a7c59; } }
  &.stat-fail { border-color: #fecaca; background: #fef2f2; .ovnum { color: #dc2626; } .ovlbl { color: #dc2626; } }
  &.stat-conflict { border-color: #fde68a; background: #fffbeb; .ovnum { color: #b87333; } .ovlbl { color: #b87333; } }
}

.sub-panel { padding: 16px; }

.stations-panel {
  padding: 16px 20px;
  .sub-panel-toolbar { padding-bottom: 14px; }
}

.history-panel {
  padding: 16px 20px;
  .sub-panel-toolbar { padding-bottom: 14px; }
}

.sub-panel-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

/* 站点卡片网格 */
.card-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 14px; }
.station-card { background: #fafaf9; border: 1px solid #e8e4de; border-radius: 12px; overflow: hidden; transition: all 0.2s; display: flex;
  &:hover { border-color: #d4cfc6; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
}
.station-bar-color { width: 4px; flex-shrink: 0; }
.station-body-inner { flex: 1; padding: 14px; display: flex; flex-direction: column; gap: 10px; }
.station-top-row { display: flex; align-items: center; gap: 10px; }
.station-avatar { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 15px; font-weight: 700; color: #fff; flex-shrink: 0; }
.station-meta-info { flex: 1; min-width: 0; }
.station-card-name { font-size: 14px; font-weight: 600; color: #1c1917; }
.station-card-code { display: flex; align-items: center; gap: 4px; font-size: 11px; color: #a8a29e; .code-label { font-size: 11px; color: #78716c; } }
.station-status-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; &.connected { background: #4a7c59; } &.disconnected { background: #ef4444; } }
.station-id-badge { font-size: 11px; color: #a8a29e; background: #f0ede8; padding: 2px 6px; border-radius: 4px; font-family: 'Monaco', monospace; flex-shrink: 0; }
.station-info-rows { display: flex; flex-direction: column; gap: 4px; }
.station-actions { display: flex; gap: 6px; justify-content: flex-end; }
.conn-hint { display: flex; align-items: center; gap: 4px; font-size: 11px; color: #a8a29e; margin-top: 4px; .el-icon { font-size: 12px; } }

/* ==================== 同步历史表格 ==================== */
.history-table-wrap {
  border: 1px solid #e8e4de;
  border-radius: 12px;
  overflow: auto;
  max-height: 600px;
}
.history-table { width: 100%; border-collapse: collapse; font-size: 14px; min-width: 800px;
  thead { background: #fafaf9; position: sticky; top: 0; z-index: 1;
    th { padding: 12px 16px; text-align: left; font-size: 12px; font-weight: 700; color: #78716c; letter-spacing: 0.5px; border-bottom: 1px solid #e8e4de; white-space: nowrap; }
  }
  tbody tr { border-bottom: 1px solid #f0ede8; cursor: pointer; transition: background 0.15s;
    &:last-child { border-bottom: none; }
    &:hover { background: #fafaf9; }
  }
  td { padding: 12px 16px; color: #1c1917; vertical-align: middle; }
}
.td-time { font-family: 'Monaco', monospace; font-size: 12px; color: #57534e; white-space: nowrap; }
.td-center { text-align: center; }
.td-num { text-align: right; font-variant-numeric: tabular-nums; font-weight: 600; &.ok { color: #4a7c59; } &.err { color: #dc2626; } &.warn { color: #b87333; } }
.td-remark { font-size: 12px; color: #78716c; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.station-chip { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; color: #fff; }
.dir-tag { display: inline-block; padding: 2px 7px; border-radius: 4px; font-size: 11px; font-weight: 500; background: #e8f0fe; color: #3b5bdb; }

.history-empty { display: flex; flex-direction: column; align-items: center; padding: 40px 20px; }

/* 表格内标签 */
.status-tag { font-size: 11px; font-weight: 500; padding: 2px 7px; border-radius: 4px; display: inline-block;
  &.pending { background: #fef9c3; color: #92400e; }
  &.processing { background: #dbeafe; color: #1e40af; }
  &.completed { background: #dcfce7; color: #4a7c59; }
  &.failed { background: #fee2e2; color: #dc2626; }
}
.action-tag { font-size: 11px; padding: 2px 7px; border-radius: 4px;
  &.create { background: #dcfce7; color: #4a7c59; }
  &.update { background: #dbeafe; color: #1e40af; }
  &.conflict { background: #fef3c7; color: #b87333; }
}
.result-tag { font-size: 11px; padding: 2px 7px; border-radius: 4px;
  &.success { background: #dcfce7; color: #4a7c59; }
  &.failed { background: #fee2e2; color: #dc2626; }
  &.conflict { background: #fef3c7; color: #b87333; }
}
.num { font-size: 13px; font-weight: 600; &.ok { color: #4a7c59; } &.err { color: #dc2626; } &.warn { color: #b87333; } }
.error-msg { font-size: 11px; color: #dc2626; line-height: 1.4; }
.new-serial { font-size: 11px; color: #409eff; font-weight: 500; }
.none { font-size: 11px; color: #d4cfc6; }

/* 空状态 */
.empty-state { grid-column: 1 / -1; display: flex; flex-direction: column; align-items: center; padding: 40px 20px; }
.empty-icon-wrap { width: 56px; height: 56px; border-radius: 50%; background: #f5f4f0; display: flex; align-items: center; justify-content: center; margin-bottom: 12px; }
.empty-icon { font-size: 26px; color: #c8c3bb; }
.empty-title { font-size: 13px; font-weight: 600; color: #78716c; margin-bottom: 4px; }
.empty-hint { font-size: 12px; color: #a8a29e; }

/* 抽屉 */
.station-history-drawer { display: flex; flex-direction: column; gap: 16px; }
.batch-toolbar { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.batch-section { }
.detail-section { border-top: 1px solid #f0ede8; padding-top: 16px; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; gap: 12px; }
.pagination-wrapper { display: flex; justify-content: center; margin-top: 14px; }

.drawer-content { padding: 0 4px; }
.drawer-tip { display: flex; align-items: center; gap: 8px; padding: 11px 14px; background: #fef9c3; border-radius: 8px; font-size: 12px; color: #92400e; margin-bottom: 16px; .el-icon { font-size: 14px; } }
.form-section { margin-bottom: 16px; &:last-child { margin-bottom: 0; } }
.form-section-title { font-size: 11px; font-weight: 600; color: #a8a29e; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 10px; }
.drawer-footer { display: flex; justify-content: flex-end; gap: 8px; padding: 12px 16px; border-top: 1px solid #f0ede8; }

.apikey-content { padding: 4px; }
.apikey-station { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.station-avatar-sm { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0; }
.apikey-station-name { font-size: 14px; font-weight: 600; color: #1c1917; }
.apikey-station-code { font-size: 12px; color: #78716c; }
.apikey-tip { display: flex; align-items: center; gap: 6px; padding: 9px 12px; background: #fef2f2; border-radius: 6px; font-size: 12px; color: #dc2626; margin-bottom: 12px; .el-icon { font-size: 13px; } }
.apikey-display { display: flex; gap: 8px; }
</style>
