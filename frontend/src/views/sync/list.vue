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
            <div class="station-header">
              <div class="station-avatar" :style="{ background: getStationColor(station.code) }">
                {{ station.name.charAt(0).toUpperCase() }}
              </div>
              <div class="station-info">
                <div class="station-name">{{ station.name }}</div>
                <div class="station-code">{{ station.code }}</div>
              </div>
              <div :class="['station-status', station.status]">
                {{ station.status === 'active' ? t('sync.enabled') : t('sync.disabled') }}
              </div>
            </div>

            <div class="station-body">
              <div class="station-stat">
                <span class="stat-label">{{ t('sync.stationUrl') }}</span>
                <span class="stat-value">{{ station.url || '-' }}</span>
              </div>
              <div class="station-stat">
                <span class="stat-label">{{ t('sync.lastSync') }}</span>
                <span class="stat-value">{{ station.lastSyncAt || '-' }}</span>
              </div>
              <div class="station-stat">
                <span class="stat-label">{{ t('sync.syncCount') }}</span>
                <span class="stat-value">{{ station.syncCount }}</span>
              </div>
            </div>

            <div class="station-footer">
              <el-button size="small" @click="openStationDrawer('edit', station)">
                {{ t('common.edit') }}
              </el-button>
              <el-button size="small" type="danger" @click="handleDeleteStation(station)">
                {{ t('common.delete') }}
              </el-button>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredStations.length === 0" class="empty-state">
            <el-icon class="empty-icon"><Folder /></el-icon>
            <div class="empty-title">{{ t('sync.noStations') }}</div>
          </div>
        </div>
      </div>

      <!-- 同步历史 -->
      <div v-show="activeTab === 'history'" class="tab-content">
        <div class="content-header">
          <div class="header-left">
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
          <el-table-column prop="stationName" :label="t('sync.stations')" min-width="120" />
          <el-table-column prop="directionText" :label="t('sync.direction')" width="100" />
          <el-table-column prop="statusText" :label="t('common.status')" width="100">
            <template #default="{ row }">
              <span :class="['status-tag', row.status]">{{ row.statusText }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.totalRecords')" width="100">
            <template #default="{ row }">
              <span class="num">{{ row.totalRecords }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.successCount')" width="100">
            <template #default="{ row }">
              <span class="num success">{{ row.successCount }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.failCount')" width="100">
            <template #default="{ row }">
              <span class="num fail">{{ row.failCount }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('sync.conflictCount')" width="100">
            <template #default="{ row }">
              <span class="num conflict">{{ row.conflictCount }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" :label="t('common.createTime')" width="160" />
          <el-table-column :label="t('common.actions')" width="120" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="viewHistoryDetail(row)">
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
    <el-drawer v-model="stationDrawerVisible" :title="stationDrawerTitle" size="440px">
      <el-form :model="stationForm" label-width="100px">
        <el-form-item :label="t('sync.stationName')" required>
          <el-input v-model="stationForm.name" />
        </el-form-item>
        <el-form-item :label="t('sync.stationCode')" required>
          <el-input v-model="stationForm.code" />
        </el-form-item>
        <el-form-item :label="t('sync.stationUrl')" required>
          <el-input v-model="stationForm.url" />
        </el-form-item>
        <el-form-item :label="t('sync.stationStatus')">
          <el-select v-model="stationForm.status" style="width: 100%">
            <el-option label="Active" value="active" />
            <el-option label="Inactive" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.remark')">
          <el-input v-model="stationForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="stationDrawerVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="saveStation" :loading="stationSaving">
            {{ t('common.save') }}
          </el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 同步详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" :title="t('sync.historyDetail')" width="640px">
      <el-table :data="detailList" stripe size="small">
        <el-table-column prop="serialNo" :label="t('sync.originalSerialNo')" min-width="140" />
        <el-table-column prop="projectName" :label="t('project.list.form.projectName')" min-width="100" />
        <el-table-column prop="actionText" :label="t('sync.action')" width="80">
          <template #default="{ row }">
            <span :class="['action-tag', row.action]">{{ row.actionText }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="result" :label="t('sync.result')" width="80">
          <template #default="{ row }">
            <span :class="['result-tag', row.result]">{{ row.result }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="newSerialNo" :label="t('sync.newSerialNo')" min-width="140">
          <template #default="{ row }">
            {{ row.newSerialNo || '-' }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Plus, Folder, Refresh, Grid, Document, Connection
} from '@element-plus/icons-vue'
import { syncApi, type SyncStation, type SyncHistory } from '@/api/sync'

const { t } = useI18n()

// 状态
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
  id: 0,
  name: '',
  code: '',
  url: '',
  status: 'active',
  description: '',
})
const stationSaving = ref(false)
const stationEditMode = ref<'create' | 'edit'>('create')

// 历史相关
const historyDateRange = ref<[string, string] | null>(null)
const historyList = ref<SyncHistory[]>([])
const historyLoading = ref(false)
const historyPage = ref(1)
const historyPageSize = ref(10)
const historyTotal = ref(0)

// 详情相关
const detailDialogVisible = ref(false)
const detailList = ref<any[]>([])

// 计算属性
const filteredStations = computed(() => {
  if (!stationKeyword.value) return stationList.value
  const kw = stationKeyword.value.toLowerCase()
  return stationList.value.filter(s =>
    s.name.toLowerCase().includes(kw) || s.code.toLowerCase().includes(kw)
  )
})

// 方法
const getStationColor = (code: string) => {
  const colors = ['#6b5b95', '#4a7c59', '#b87333', '#4a6fa5', '#8b6f5b', '#5b7b8c']
  let hash = 0
  for (let i = 0; i < code.length; i++) hash = ((hash << 5) - hash) + code.charCodeAt(i)
  return colors[Math.abs(hash) % colors.length]
}

const loadStations = async () => {
  stationLoading.value = true
  try {
    const res = await syncApi.listStations({ page: 1, pageSize: 100 })
    stationList.value = res.data.items
  } catch (e) {
    console.error(e)
  } finally {
    stationLoading.value = false
  }
}

const loadHistory = async () => {
  historyLoading.value = true
  try {
    const res = await syncApi.getHistory({
      page: historyPage.value,
      pageSize: historyPageSize.value,
      startDate: historyDateRange.value?.[0],
      endDate: historyDateRange.value?.[1],
    })
    historyList.value = res.data.items
    historyTotal.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    historyLoading.value = false
  }
}

const openStationDrawer = (mode: 'create' | 'edit', station?: SyncStation) => {
  stationEditMode.value = mode
  stationDrawerTitle.value = mode === 'create' ? t('sync.createStation') : t('sync.editStation')
  if (mode === 'edit' && station) {
    stationForm.value = {
      id: station.id,
      name: station.name,
      code: station.code,
      url: station.url,
      status: station.status,
      description: station.description || '',
    }
  } else {
    stationForm.value = {
      id: 0,
      name: '',
      code: '',
      url: '',
      status: 'active',
      description: '',
    }
  }
  stationDrawerVisible.value = true
}

const saveStation = async () => {
  if (!stationForm.value.name) {
    ElMessage.warning(t('sync.messages.nameRequired'))
    return
  }
  stationSaving.value = true
  try {
    if (stationEditMode.value === 'create') {
      await syncApi.createStation(stationForm.value)
      ElMessage.success(t('sync.messages.createSuccess'))
    } else {
      await syncApi.updateStation(stationForm.value)
      ElMessage.success(t('sync.messages.updateSuccess'))
    }
    stationDrawerVisible.value = false
    loadStations()
  } catch (e: any) {
    ElMessage.error(e.message || t('sync.messages.createFailed'))
  } finally {
    stationSaving.value = false
  }
}

const handleDeleteStation = async (station: SyncStation) => {
  try {
    await ElMessageBox.confirm(
      t('sync.deleteConfirm', { name: station.name }),
      t('common.confirm'),
      { type: 'warning' }
    )
    await syncApi.deleteStation(station.id)
    ElMessage.success(t('sync.messages.deleteSuccess'))
    loadStations()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(t('sync.messages.deleteFailed'))
  }
}

const viewHistoryDetail = async (row: SyncHistory) => {
  try {
    const res = await syncApi.getHistoryDetails(row.id)
    detailList.value = res.data.details || []
    detailDialogVisible.value = true
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadStations()
  loadHistory()
})
</script>

<style scoped lang="scss">
.sync-page {
  padding: 20px 24px;
  height: 100vh;
  overflow-y: auto;
  background: #fafaf9;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #1c1917;
  margin: 0;
}

.page-subtitle {
  font-size: 13px;
  color: #a8a29e;
}

/* 标签页 */
.view-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 10px;
  font-size: 13px;
  color: #78716c;
  cursor: pointer;
  transition: all 0.2s;

  .el-icon { font-size: 14px; }

  &:hover {
    border-color: #d4d0c8;
    color: #57534e;
  }

  &.active {
    background: #fff;
    border-color: #6b5b95;
    color: #6b5b95;
    box-shadow: 0 2px 8px rgba(107, 91, 149, 0.15);
  }
}

/* 内容区 */
.page-content {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
}

.tab-content {
  padding: 20px;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.header-right {
  display: flex;
  gap: 8px;
}

/* 站点卡片 */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.station-card {
  background: #fafaf9;
  border-radius: 14px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
  transition: all 0.2s;

  &:hover {
    border-color: #d4d0c8;
    transform: translateY(-2px);
  }
}

.station-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #f0ede8;
}

.station-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.station-info {
  flex: 1;
  min-width: 0;
}

.station-name {
  font-size: 14px;
  font-weight: 600;
  color: #1c1917;
}

.station-code {
  font-size: 11px;
  color: #a8a29e;
}

.station-status {
  font-size: 11px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 6px;

  &.active {
    background: #f0fdf4;
    color: #4a7c59;
  }

  &.inactive {
    background: #f5f5f4;
    color: #a8a29e;
  }
}

.station-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.station-stat {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-label {
  font-size: 12px;
  color: #a8a29e;
}

.stat-value {
  font-size: 12px;
  color: #57534e;
  font-weight: 500;
}

.station-footer {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #f0ede8;

  .el-button {
    flex: 1;
  }
}

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 48px;
  color: #d4d0c8;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 14px;
  color: #a8a29e;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 标签样式 */
.status-tag {
  font-size: 11px;
  font-weight: 500;
  padding: 3px 8px;
  border-radius: 6px;

  &.pending { background: #fef9c3; color: #854d0e; }
  &.processing { background: #dbeafe; color: #1e40af; }
  &.completed { background: #f0fdf4; color: #4a7c59; }
  &.failed { background: #fef2f2; color: #dc2626; }
}

.num {
  font-weight: 600;
  &.success { color: #4a7c59; }
  &.fail { color: #dc2626; }
  &.conflict { color: #b87333; }
}

.action-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;

  &.create { background: #f0fdf4; color: #4a7c59; }
  &.update { background: #dbeafe; color: #1e40af; }
  &.conflict { background: #fef3c7; color: #b87333; }
  &.skip { background: #f5f5f4; color: #78716c; }
}

.result-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;

  &.success { background: #f0fdf4; color: #4a7c59; }
  &.failed { background: #fef2f2; color: #dc2626; }
}

/* 抽屉 */
.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
