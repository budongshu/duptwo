<template>
  <div class="sync-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('sync.title') }}</h1>
        <span class="page-subtitle">{{ t('sync.subtitle') }}</span>
      </div>
    </header>

    <!-- 标签页 -->
    <div class="view-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-btn', { active: activeTab === tab.key }]"
        @click="activeTab = tab.key"
      >
        <el-icon><component :is="tab.icon" /></el-icon>
        {{ tab.label }}
      </button>
    </div>

    <!-- 内容区 -->
    <div class="page-content">
      <!-- 站点管理 -->
      <div v-show="activeTab === 'stations'" class="tab-content">
        <div class="content-header">
          <div class="header-left">
            <el-input
              v-model="stationKeyword"
              :placeholder="t('common.search') + '...'"
              clearable
              style="width: 200px"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="openStationDrawer('create')">
              <el-icon><Plus /></el-icon>
              {{ t('sync.createStation') }}
            </el-button>
          </div>
        </div>

        <!-- 站点列表 -->
        <div class="card-grid">
          <div
            v-for="station in filteredStations"
            :key="station.id"
            class="station-card"
          >
            <!-- 左侧状态条 -->
            <div class="station-bar" :style="{ background: getStationColor(station.code) }"></div>

            <div class="station-inner">
              <div class="station-top">
                <div class="station-avatar" :style="{ background: getStationColor(station.code) }">
                  {{ station.name.charAt(0).toUpperCase() }}
                </div>
                <div class="station-info">
                  <div class="station-name">{{ station.name }}</div>
                  <div class="station-meta">
                    <span class="station-code">
                      <el-icon><Link /></el-icon>
                      {{ station.code }}
                    </span>
                    <span class="station-status-tag" :class="station.status">
                      <span class="status-dot"></span>
                      {{ station.status === 'active' ? t('sync.enabled') : t('sync.disabled') }}
                    </span>
                  </div>
                </div>
              </div>

              <div class="station-body">
                <div class="info-row">
                  <span class="info-label">
                    <el-icon><Location /></el-icon>
                    {{ t('sync.stationUrl') }}
                  </span>
                  <span class="info-value url" :title="station.url || '-'">
                    {{ station.url || '-' }}
                  </span>
                </div>
                <div class="info-row">
                  <span class="info-label">
                    <el-icon><Clock /></el-icon>
                    {{ t('sync.lastSync') }}
                  </span>
                  <span class="info-value">{{ formatTime(station.lastSyncAt) }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">
                    <el-icon><Guide /></el-icon>
                    {{ t('sync.syncCount') }}
                  </span>
                  <span class="info-value count">{{ station.syncCount }}</span>
                </div>
              </div>

              <div class="station-footer">
                <el-button size="small" @click="openStationDrawer('edit', station)">
                  <el-icon><Edit /></el-icon>
                  {{ t('common.edit') }}
                </el-button>
                <el-button size="small" @click="showApiKey(station)" type="info" plain>
                  <el-icon><Key /></el-icon>
                  {{ t('sync.apiKey') }}
                </el-button>
                <el-button size="small" type="danger" @click="handleDeleteStation(station)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredStations.length === 0" class="empty-state">
            <div class="empty-icon-wrap">
              <el-icon class="empty-icon"><Connection /></el-icon>
            </div>
            <div class="empty-title">{{ t('sync.noStations') }}</div>
            <div class="empty-hint">点击上方按钮添加第一个同步站点</div>
          </div>
        </div>
      </div>

      <!-- 同步历史 -->
      <div v-show="activeTab === 'history'" class="tab-content">
        <div class="content-header">
          <div class="header-left">
            <el-select v-model="historyStationId" :placeholder="t('sync.stations')" clearable style="width: 150px">
              <el-option :label="t('sync.stationList')" :value="0" />
              <el-option v-for="s in stationList" :key="s.id" :label="s.name" :value="s.id" />
            </el-select>
            <el-date-picker
              v-model="historyDateRange"
              type="daterange"
              range-separator="-"
              :start-placeholder="t('common.startDate')"
              :end-placeholder="t('common.endDate')"
              value-format="YYYY-MM-DD"
              style="width: 240px"
            />
          </div>
          <div class="header-right">
            <el-button @click="loadHistory">
              <el-icon><Refresh /></el-icon>
              {{ t('common.refresh') }}
            </el-button>
          </div>
        </div>

        <!-- 历史记录表格 -->
        <el-table :data="historyList" stripe v-loading="historyLoading">
          <el-table-column prop="stationName" :label="t('sync.stations')" min-width="120">
            <template #default="{ row }">
              <div class="station-cell">
                <span class="station-dot" :style="{ background: getStationColor(row.stationId.toString()) }"></span>
                {{ row.stationName || '-' }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="directionText" :label="t('sync.direction')" width="90" align="center" />
          <el-table-column prop="statusText" :label="t('common.status')" width="100" align="center">
            <template #default="{ row }">
              <span :class="['status-tag', row.status]">{{ row.statusText }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.totalRecords')" width="80" align="center">
            <template #default="{ row }">
              <span class="num">{{ row.totalRecords }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.successCount')" width="80" align="center">
            <template #default="{ row }">
              <span class="num success">{{ row.successCount }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.failCount')" width="80" align="center">
            <template #default="{ row }">
              <span class="num fail">{{ row.failCount }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.conflictCount')" width="90" align="center">
            <template #default="{ row }">
              <span class="num conflict">{{ row.conflictCount }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" :label="t('common.createTime')" width="160" />
          <el-table-column :label="t('common.actions')" width="100" fixed="right" align="center">
            <template #default="{ row }">
              <el-button size="small" link @click="viewHistoryDetail(row)">
                {{ t('sync.syncDetail') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="historyPage"
            v-model:page-size="historyPageSize"
            :total="historyTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @change="loadHistory"
          />
        </div>
      </div>
    </div>

    <!-- 站点编辑抽屉 -->
    <el-drawer v-model="stationDrawerVisible" :title="stationDrawerTitle" size="520px" class="sync-drawer">
      <div class="drawer-content">
        <!-- 帮助提示 -->
        <div class="drawer-tip">
          <el-icon><InfoFilled /></el-icon>
          <span v-if="stationEditMode === 'create'">创建站点后，将自动生成 API Key，请妥善保管</span>
          <span v-else>修改站点信息不会影响已生成的 API Key</span>
        </div>

        <el-form :model="stationForm" label-position="top" class="station-form">
          <div class="form-section">
            <div class="form-section-title">{{ t('sync.stationName') }}</div>
            <el-form-item :label="t('sync.stationName')" required>
              <el-input v-model="stationForm.name" :placeholder="t('sync.messages.nameRequired')" />
            </el-form-item>
            <el-form-item :label="t('sync.stationCode')" required>
              <el-input v-model="stationForm.code" :placeholder="t('sync.messages.codeRequired')">
                <template #suffix>
                  <el-tooltip content="站点唯一标识，用于Agent识别" placement="top">
                    <el-icon><QuestionFilled /></el-icon>
                  </el-tooltip>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="t('sync.stationUrl')" required>
              <el-input v-model="stationForm.url" placeholder="https://center.example.com">
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
            <el-form-item :label="t('common.remark')">
              <el-input v-model="stationForm.description" type="textarea" :rows="3" />
            </el-form-item>
          </div>
        </el-form>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="stationDrawerVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="saveStation" :loading="stationSaving">
            {{ t('common.save') }}
          </el-button>
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

        <div class="apikey-tip">
          <el-icon><Warning /></el-icon>
          <span>API Key 只显示一次，请立即复制并妥善保管！</span>
        </div>

        <div class="apikey-display">
          <el-input v-model="displayApiKey" readonly :placeholder="'API Key 将显示在这里'" />
          <el-button type="primary" @click="copyApiKey" :disabled="!displayApiKey">
            <el-icon><DocumentCopy /></el-icon>
            {{ t('sync.copyApiKey') }}
          </el-button>
        </div>

        <div class="apikey-hint">
          <p>将此 API Key 配置到 Agent 站点的 <code>app.yaml</code> 中：</p>
          <pre class="config-example">sync:
  enabled: true
  mode: "agent"
  center_url: "{{ currentStationForKey.url }}"
  api_key: "{{ displayApiKey }}"
  station_id: "{{ currentStationForKey.code }}"
  station_name: "{{ currentStationForKey.name }}"</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 同步详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" :title="t('sync.historyDetail')" width="720px">
      <div class="detail-summary">
        <div class="summary-item">
          <span class="summary-num">{{ detailSummary.total }}</span>
          <span class="summary-label">{{ t('sync.totalRecords') }}</span>
        </div>
        <div class="summary-item success">
          <span class="summary-num">{{ detailSummary.success }}</span>
          <span class="summary-label">{{ t('sync.successCount') }}</span>
        </div>
        <div class="summary-item fail">
          <span class="summary-num">{{ detailSummary.fail }}</span>
          <span class="summary-label">{{ t('sync.failCount') }}</span>
        </div>
        <div class="summary-item conflict">
          <span class="summary-num">{{ detailSummary.conflict }}</span>
          <span class="summary-label">{{ t('sync.conflictCount') }}</span>
        </div>
      </div>

      <el-table :data="detailList" stripe size="small" max-height="400">
        <el-table-column prop="serialNo" :label="t('sync.originalSerialNo')" min-width="140" />
        <el-table-column prop="projectName" :label="t('project.list.form.projectName')" min-width="100" />
        <el-table-column prop="actionText" :label="t('sync.action')" width="80" align="center">
          <template #default="{ row }">
            <span :class="['action-tag', row.action]">{{ row.actionText }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="result" :label="t('sync.result')" width="80" align="center">
          <template #default="{ row }">
            <span :class="['result-tag', row.result]">{{ row.result }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="newSerialNo" :label="t('sync.newSerialNo')" min-width="140">
          <template #default="{ row }">
            <span v-if="row.newSerialNo" class="new-serial">{{ row.newSerialNo }}</span>
            <span v-else class="none">-</span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Search, Plus, Refresh, Grid, Document, Connection, Link, Clock, Guide, Edit, Delete, Key,
  Location, QuestionFilled, InfoFilled, Warning, DocumentCopy
} from '@element-plus/icons-vue'
import { syncApi, type SyncStation, type SyncHistory } from '@/api/sync'

const { t, locale } = useI18n()

const activeTab = ref('stations')
const tabs = computed(() => [
  { key: 'stations', label: t('sync.stations'), icon: Grid },
  { key: 'history', label: t('sync.history'), icon: Document },
])

// 站点相关
const stationKeyword = ref('')
const stationList = ref<SyncStation[]>([])
const stationLoading = ref(false)
const stationDrawerVisible = ref(false)
const stationDrawerTitle = ref('')
const stationForm = ref({
  id: 0, name: '', code: '', url: '', status: 'active', description: '',
})
const stationSaving = ref(false)
const stationEditMode = ref<'create' | 'edit'>('create')

// API Key 相关
const apiKeyDialogVisible = ref(false)
const currentStationForKey = ref<SyncStation | null>(null)
const displayApiKey = ref('')

// 历史相关
const historyStationId = ref<number>(0)
const historyDateRange = ref<[string, string] | null>(null)
const historyList = ref<SyncHistory[]>([])
const historyLoading = ref(false)
const historyPage = ref(1)
const historyPageSize = ref(10)
const historyTotal = ref(0)

// 详情相关
const detailDialogVisible = ref(false)
const detailList = ref<any[]>([])

const detailSummary = computed(() => {
  let total = 0, success = 0, fail = 0, conflict = 0
  detailList.value.forEach(d => {
    total++
    if (d.result === 'success') success++
    else if (d.result === 'failed') fail++
    if (d.action === 'conflict') conflict++
  })
  return { total, success, fail, conflict }
})

const filteredStations = computed(() => {
  if (!stationKeyword.value) return stationList.value
  const kw = stationKeyword.value.toLowerCase()
  return stationList.value.filter(s => s.name.toLowerCase().includes(kw) || s.code.toLowerCase().includes(kw))
})

const getStationColor = (code: string) => {
  const colors = ['#6b5b95', '#4a7c59', '#b87333', '#4a6fa5', '#8b6f5b', '#5b7b8c', '#7a8471', '#a67c52']
  let hash = 0
  for (let i = 0; i < code.length; i++) hash = ((hash << 5) - hash) + code.charCodeAt(i)
  return colors[Math.abs(hash) % colors.length]
}

const formatTime = (time: string | null) => {
  if (!time) return '-'
  return time.replace('T', ' ').slice(0, 16)
}

const loadStations = async () => {
  stationLoading.value = true
  try {
    const res = await syncApi.listStations({ page: 1, pageSize: 100 })
    stationList.value = res.data.items
  } catch (e) { console.error(e) }
  finally { stationLoading.value = false }
}

const loadHistory = async () => {
  historyLoading.value = true
  try {
    const res = await syncApi.getHistory({
      page: historyPage.value, pageSize: historyPageSize.value,
      stationId: historyStationId.value || undefined,
      startDate: historyDateRange.value?.[0],
      endDate: historyDateRange.value?.[1],
    })
    historyList.value = res.data.items
    historyTotal.value = res.data.total
  } catch (e) { console.error(e) }
  finally { historyLoading.value = false }
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
      // 自动显示 API Key
      currentStationForKey.value = res.data
      displayApiKey.value = res.data.apiKey || '请到站点详情查看'
      apiKeyDialogVisible.value = true
    } else {
      await syncApi.updateStation(stationForm.value)
      ElMessage.success(t('sync.messages.updateSuccess'))
    }
    stationDrawerVisible.value = false
    loadStations()
  } catch (e: any) { ElMessage.error(e.message || t('sync.messages.createFailed')) }
  finally { stationSaving.value = false }
}

const handleDeleteStation = async (station: SyncStation) => {
  try {
    await syncApi.deleteStation(station.id)
    ElMessage.success(t('sync.messages.deleteSuccess'))
    loadStations()
  } catch (e: any) { ElMessage.error(t('sync.messages.deleteFailed')) }
}

const showApiKey = (station: SyncStation) => {
  currentStationForKey.value = station
  // 实际项目中需要从后端获取真实 API Key，这里模拟
  displayApiKey.value = station.apiKey || '请重新保存站点获取'
  apiKeyDialogVisible.value = true
}

const copyApiKey = async () => {
  try {
    await navigator.clipboard.writeText(displayApiKey.value)
    ElMessage.success(t('sync.apiKeyCopied'))
  } catch { ElMessage.error('复制失败') }
}

const viewHistoryDetail = async (row: SyncHistory) => {
  try {
    const res = await syncApi.getHistoryDetails(row.id)
    detailList.value = res.data.details || []
    detailDialogVisible.value = true
  } catch (e) { console.error(e) }
}

onMounted(() => { loadStations(); loadHistory() })
</script>

<style scoped lang="scss">
.sync-page {
  padding: 20px 24px;
  min-height: 100vh;
  background: #fafaf9;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.header-left { display: flex; flex-direction: column; gap: 4px; }
.page-title { font-size: 20px; font-weight: 700; color: #1c1917; margin: 0; }
.page-subtitle { font-size: 13px; color: #a8a29e; }

/* 标签页 */
.view-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}
.tab-btn {
  display: flex; align-items: center; gap: 6px;
  padding: 10px 18px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 12px;
  font-size: 13px; font-weight: 500;
  color: #78716c;
  cursor: pointer;
  transition: all 0.2s;
  .el-icon { font-size: 15px; }
  &:hover { border-color: #d4d0c8; color: #57534e; }
  &.active {
    background: #fff;
    border-color: #6b5b95;
    color: #6b5b95;
    box-shadow: 0 2px 10px rgba(107, 91, 149, 0.15);
  }
}

/* 内容区 */
.page-content {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
}
.tab-content { padding: 20px; }
.content-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 20px; gap: 12px;
}
.header-right { display: flex; gap: 8px; }

/* 站点卡片 */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}
.station-card {
  background: #fafaf9;
  border-radius: 16px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
  transition: all 0.25s;
  display: flex;
  &:hover { border-color: #d4d0c8; transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.06); }
}
.station-bar { width: 4px; flex-shrink: 0; }
.station-inner { flex: 1; padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.station-top { display: flex; align-items: center; gap: 12px; }
.station-avatar {
  width: 44px; height: 44px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700; color: #fff; flex-shrink: 0;
}
.station-info { flex: 1; min-width: 0; }
.station-name { font-size: 15px; font-weight: 600; color: #1c1917; }
.station-meta { display: flex; align-items: center; gap: 8px; margin-top: 4px; }
.station-code {
  display: flex; align-items: center; gap: 4px;
  font-size: 11px; color: #a8a29e;
  .el-icon { font-size: 11px; }
}
.station-status-tag {
  display: flex; align-items: center; gap: 4px;
  font-size: 11px; font-weight: 500; padding: 2px 8px; border-radius: 6px;
  .status-dot { width: 6px; height: 6px; border-radius: 50%; }
  &.active { background: #f0fdf4; color: #4a7c59; .status-dot { background: #4a7c59; } }
  &.inactive { background: #f5f5f4; color: #a8a29e; .status-dot { background: #a8a29e; } }
}

.station-body { display: flex; flex-direction: column; gap: 8px; }
.info-row { display: flex; align-items: center; justify-content: space-between; padding: 6px 0; border-bottom: 1px solid #f0ede8; &:last-child { border-bottom: none; } }
.info-label { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #a8a29e; .el-icon { font-size: 13px; } }
.info-value { font-size: 12px; color: #57534e; font-weight: 500; &.url { max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; } &.count { font-size: 14px; font-weight: 700; color: #6b5b95; } }

.station-footer { display: flex; gap: 8px; margin-top: 4px; .el-button { flex: 1; } }

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  display: flex; flex-direction: column; align-items: center;
  padding: 80px 20px;
}
.empty-icon-wrap {
  width: 80px; height: 80px; border-radius: 50%;
  background: #f5f5f4;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 20px;
}
.empty-icon { font-size: 36px; color: #d4d0c8; }
.empty-title { font-size: 16px; font-weight: 600; color: #78716c; margin-bottom: 8px; }
.empty-hint { font-size: 13px; color: #a8a29e; }

/* 分页 */
.pagination-wrapper { display: flex; justify-content: center; margin-top: 20px; }

/* 表格内样式 */
.station-cell { display: flex; align-items: center; gap: 8px; }
.station-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.status-tag { font-size: 11px; font-weight: 500; padding: 3px 10px; border-radius: 6px; &.pending { background: #fef9c3; color: #854d0e; } &.processing { background: #dbeafe; color: #1e40af; } &.completed { background: #f0fdf4; color: #4a7c59; } &.failed { background: #fef2f2; color: #dc2626; } }
.num { font-weight: 600; &.success { color: #4a7c59; } &.fail { color: #dc2626; } &.conflict { color: #b87333; } }
.action-tag { font-size: 11px; padding: 2px 8px; border-radius: 4px; &.create { background: #f0fdf4; color: #4a7c59; } &.update { background: #dbeafe; color: #1e40af; } &.conflict { background: #fef3c7; color: #b87333; } &.skip { background: #f5f5f4; color: #78716c; } }
.result-tag { font-size: 11px; padding: 2px 8px; border-radius: 4px; &.success { background: #f0fdf4; color: #4a7c59; } &.failed { background: #fef2f2; color: #dc2626; } }
.new-serial { color: #6b5b95; font-weight: 500; } .none { color: #d4d0c8; }

/* 抽屉样式 */
.drawer-content { padding: 0 4px; }
.drawer-tip {
  display: flex; align-items: center; gap: 8px;
  padding: 12px 16px; background: #fef9c3; border-radius: 10px;
  font-size: 13px; color: #854d0e; margin-bottom: 20px;
  .el-icon { font-size: 16px; }
}
.station-form { }
.form-section { margin-bottom: 24px; &:last-child { margin-bottom: 0; } }
.form-section-title { font-size: 12px; font-weight: 600; color: #a8a29e; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px; }
.drawer-footer { display: flex; justify-content: flex-end; gap: 12px; padding: 16px 24px; border-top: 1px solid #f0ede8; }

/* API Key 弹窗 */
.apikey-content { display: flex; flex-direction: column; gap: 16px; }
.apikey-station { display: flex; align-items: center; gap: 12px; padding: 12px; background: #fafaf9; border-radius: 10px; }
.station-avatar-sm { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; color: #fff; }
.apikey-station-name { font-size: 14px; font-weight: 600; color: #1c1917; }
.apikey-station-code { font-size: 11px; color: #a8a29e; }
.apikey-tip {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 14px; background: #fef2f2; border-radius: 8px;
  font-size: 12px; color: #dc2626;
  .el-icon { font-size: 14px; }
}
.apikey-display { display: flex; gap: 8px; }
.apikey-hint {
  p { font-size: 12px; color: #78716c; margin: 0 0 8px 0; }
  code { background: #f5f5f4; padding: 2px 6px; border-radius: 4px; color: #6b5b95; }
}
.config-example {
  background: #1c1917; color: #e7e5e4;
  padding: 14px 16px; border-radius: 10px;
  font-size: 12px; font-family: 'Monaco', 'Menlo', monospace;
  overflow-x: auto; margin: 0;
}

/* 详情弹窗 */
.detail-summary { display: flex; gap: 16px; margin-bottom: 20px; }
.summary-item {
  flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px;
  padding: 16px; background: #fafaf9; border-radius: 12px; border: 1px solid #e8e5e1;
  &.success { border-color: #4a7c59; .summary-num { color: #4a7c59; } }
  &.fail { border-color: #dc2626; .summary-num { color: #dc2626; } }
  &.conflict { border-color: #b87333; .summary-num { color: #b87333; } }
}
.summary-num { font-size: 24px; font-weight: 700; color: #1c1917; }
.summary-label { font-size: 11px; color: #78716c; }
</style>
