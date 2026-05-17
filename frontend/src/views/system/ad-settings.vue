<template>
  <div class="ad-page">
    <!-- 页面标题栏 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('adSettings.title') }}</h1>
        <span class="page-subtitle">{{ t('adSettings.subtitle') }}</span>
      </div>
      <div class="header-right">
        <el-button
          v-if="form.enabled"
          type="primary"
          size="default"
          class="sync-btn"
          @click="handleSync"
          :loading="syncing"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M23 4v6h-6M1 20v-6h6M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          {{ t('adSettings.actions.sync') }}
        </el-button>
      </div>
    </header>

    <!-- 状态面板 -->
    <div class="status-bar">
      <div class="status-bar-left">
        <div class="status-item" :class="form.enabled ? 'status--on' : 'status--off'">
          <span class="status-dot"></span>
          <span class="status-label">{{ form.enabled ? t('adSettings.enabled') : t('adSettings.disabled') }}</span>
        </div>
        <template v-if="form.enabled && form.server">
          <div class="status-sep"></div>
          <div class="status-info">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
            <span>{{ form.server }}:{{ form.port }}</span>
          </div>
          <div class="status-info" v-if="form.use_ssl">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            <span>LDAPS</span>
          </div>
        </template>
      </div>
      <div class="status-bar-right">
        <div class="sync-stat" v-if="lastSyncTime">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          <span>{{ t('adSettings.sync.lastSync') }}: {{ formatTime(lastSyncTime) }}</span>
        </div>
      </div>
    </div>

    <!-- 总开关卡片 -->
    <div class="enable-card" :class="{ 'enable-card--active': form.enabled }">
      <div class="enable-left">
        <div class="enable-icon" :class="form.enabled ? 'icon-active' : 'icon-inactive'">
          <svg v-if="form.enabled" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
        </div>
        <div class="enable-text">
          <div class="enable-title">{{ form.enabled ? t('adSettings.enabled') : t('adSettings.disabled') }}</div>
          <div class="enable-desc">{{ form.enabled ? t('adSettings.enabledDesc') : t('adSettings.disabledDesc') }}</div>
        </div>
      </div>
      <el-switch v-model="form.enabled" size="large" @change="onEnabledChange" />
    </div>

    <!-- 配置内容 -->
    <div class="config-content" :class="{ 'config-content--disabled': !form.enabled }">

      <!-- 连接配置 -->
      <div class="config-card">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #005eeb"></div>
          <span class="config-card-title">{{ t('adSettings.connection.title') }}</span>
        </div>
        <div class="config-card-body">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top" :disabled="loading">
            <div class="form-grid-2">
              <div class="form-field">
                <label class="field-label">{{ t('adSettings.connection.server') }} <span class="required-mark">*</span></label>
                <el-input v-model="form.server" :placeholder="t('adSettings.connection.serverPlaceholder')" clearable size="default" class="field-input">
                  <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg></template>
                </el-input>
              </div>
              <div class="form-field">
                <label class="field-label">{{ t('adSettings.connection.port') }} <span class="required-mark">*</span></label>
                <el-input-number v-model="form.port" :min="1" :max="65535" :placeholder="t('adSettings.connection.portPlaceholder')" style="width: 100%" size="default" />
                <div class="field-hint">{{ t('adSettings.connection.portHint') }}</div>
              </div>
            </div>
            <div class="form-field">
              <label class="field-label">{{ t('adSettings.connection.baseDn') }} <span class="required-mark">*</span></label>
              <el-input v-model="form.base_dn" :placeholder="t('adSettings.connection.baseDnPlaceholder')" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg></template>
              </el-input>
              <div class="field-hint">{{ t('adSettings.connection.baseDnHint') }}</div>
            </div>
            <div class="ssl-toggle-row">
              <el-switch v-model="form.use_ssl" size="small" />
              <span class="ssl-label">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                {{ t('adSettings.connection.useSsl') }}
              </span>
            </div>
          </el-form>
        </div>
      </div>

      <!-- 认证配置 -->
      <div class="config-card">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #06b6d4"></div>
          <span class="config-card-title">{{ t('adSettings.auth.title') }}</span>
        </div>
        <div class="config-card-body">
          <el-form ref="formRef2" :model="form" :rules="rules" label-position="top" :disabled="loading">
            <div class="form-field">
              <label class="field-label">{{ t('adSettings.auth.bindDn') }} <span class="required-mark">*</span></label>
              <el-input v-model="form.bind_dn" :placeholder="t('adSettings.auth.bindDnPlaceholder')" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
              </el-input>
              <div class="field-hint">{{ t('adSettings.auth.bindDnHint') }}</div>
            </div>
            <div class="form-field">
              <label class="field-label">{{ t('adSettings.auth.bindPassword') }}</label>
              <el-input v-model="form.bind_password" type="password" :placeholder="t('adSettings.auth.bindPasswordPlaceholder')" show-password clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg></template>
              </el-input>
              <div class="field-hint">{{ t('adSettings.auth.bindPasswordHint') }}</div>
            </div>
            <div class="form-field">
              <label class="field-label">{{ t('adSettings.auth.userFilter') }}</label>
              <el-input v-model="form.user_filter" :placeholder="t('adSettings.auth.userFilterPlaceholder')" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></template>
              </el-input>
              <div class="field-hint">{{ t('adSettings.auth.userFilterHint') }}</div>
            </div>
          </el-form>
        </div>
      </div>

      <!-- 自动注册 -->
      <div class="config-card" :class="{ 'config-card--warn': form.enabled && form.auto_register }">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #f59e0b"></div>
          <span class="config-card-title">{{ t('adSettings.autoReg.title') }}</span>
          <span class="config-card-badge config-card-badge--warn" v-if="form.enabled && form.auto_register">{{ t('adSettings.autoReg.badge') }}</span>
        </div>
        <div class="config-card-body">
          <div class="auto-reg-row">
            <div class="auto-reg-info">
              <div class="auto-reg-title">{{ t('adSettings.autoReg.allowAuto') }}</div>
              <div class="auto-reg-desc">{{ t('adSettings.autoReg.allowAutoDesc') }}</div>
            </div>
            <el-switch v-model="form.auto_register" :disabled="!form.enabled" size="small" />
          </div>

          <div v-if="form.enabled && form.auto_register" class="auto-reg-expanded">
            <div class="policy-box">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              <div class="policy-content">
                <div class="policy-title">{{ t('adSettings.autoReg.policyTitle') }}</div>
                <ul class="policy-list">
                  <li>{{ t('adSettings.autoReg.policy1') }}</li>
                  <li>{{ t('adSettings.autoReg.policy2') }}</li>
                  <li>{{ t('adSettings.autoReg.policy3') }}</li>
                </ul>
              </div>
            </div>
            <div class="form-field" style="margin-top: 16px">
              <label class="field-label">{{ t('adSettings.autoReg.defaultRole') }}</label>
              <el-select v-model="form.default_role_id" :placeholder="t('adSettings.autoReg.defaultRolePlaceholder')" size="default" style="width: 280px">
                <el-option v-for="role in roles" :key="role.id" :label="role.name + '（' + role.code + '）'" :value="role.id" />
              </el-select>
              <div class="field-hint">{{ t('adSettings.autoReg.defaultRoleHint') }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 同步配置 -->
      <div class="config-card">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #8b5cf6"></div>
          <span class="config-card-title">{{ t('adSettings.sync.title') }}</span>
        </div>
        <div class="config-card-body">
          <div class="sync-info">
            <div class="sync-info-item">
              <span class="sync-info-label">{{ t('adSettings.sync.totalAD') }}</span>
              <span class="sync-info-value">{{ syncStats.total }}</span>
            </div>
            <div class="sync-info-item">
              <span class="sync-info-label">{{ t('adSettings.sync.synced') }}</span>
              <span class="sync-info-value sync-info-value--success">{{ syncStats.synced }}</span>
            </div>
            <div class="sync-info-item" v-if="lastSyncResult">
              <span class="sync-info-label">{{ t('adSettings.sync.lastResult') }}</span>
              <div class="sync-result-tags">
                <span class="sync-tag sync-tag--created" v-if="lastSyncResult.created > 0">+{{ lastSyncResult.created }} {{ t('adSettings.sync.created') }}</span>
                <span class="sync-tag sync-tag--updated" v-if="lastSyncResult.updated > 0">~{{ lastSyncResult.updated }} {{ t('adSettings.sync.updated') }}</span>
                <span class="sync-tag sync-tag--skipped" v-if="lastSyncResult.skipped > 0">{{ lastSyncResult.skipped }} {{ t('adSettings.sync.skipped') }}</span>
                <span class="sync-tag sync-tag--disabled" v-if="lastSyncResult.disabled > 0">-{{ lastSyncResult.disabled }} {{ t('adSettings.sync.disabled') }}</span>
                <span class="sync-tag sync-tag--deleted" v-if="lastSyncResult.deleted > 0">x{{ lastSyncResult.deleted }} {{ t('adSettings.sync.deleted') }}</span>
              </div>
            </div>
          </div>
          <div class="sync-actions">
            <el-button @click="handlePreview" :loading="previewLoading" :disabled="!form.enabled || !form.server" size="default">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              {{ t('adSettings.sync.preview') }}
            </el-button>
            <el-button @click="handleSync('incremental')" :loading="syncing" :disabled="!form.enabled || !form.server" size="default">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12l7-7 7 7"/></svg>
              {{ t('adSettings.sync.incrementalSync') }}
            </el-button>
            <el-button type="primary" @click="handleSync('full')" :loading="syncing" :disabled="!form.enabled || !form.server" size="default" class="sync-submit-btn">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 4v6h-6M1 20v-6h6M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
              {{ t('adSettings.sync.fullSync') }}
            </el-button>
            <el-button type="danger" plain @click="handleResetAll" :loading="resetting" :disabled="!form.enabled || !form.server" size="default">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4"/></svg>
              {{ t('adSettings.sync.resetAll') }}
            </el-button>
          </div>
        </div>
      </div>

      <!-- 操作区 -->
      <div class="actions-card">
        <div class="actions-left">
          <el-button @click="handleTest" :loading="testLoading" :disabled="!form.enabled || !form.server" size="default" class="test-btn">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="13 2 13 9 20 9"/><path d="M20 14.66V17a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h2"/></svg>
            {{ t('adSettings.actions.test') }}
          </el-button>
          <transition name="result-fade">
            <span v-if="testResult !== null" class="test-result" :class="testResult ? 'result--ok' : 'result--fail'">
              <svg v-if="testResult" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
              <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              {{ testResult ? t('adSettings.actions.testSuccess') : t('adSettings.actions.testFail') }}
            </span>
          </transition>
        </div>
        <div class="actions-right">
          <el-button @click="loadConfig" :disabled="loading" size="default">{{ t('common.reset') }}</el-button>
          <el-button type="primary" @click="handleSave" :loading="saving" :disabled="!form.enabled" size="default" class="save-btn">
            <svg v-if="!saving" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
            {{ t('adSettings.actions.save') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 用户预览弹窗 -->
    <el-dialog v-model="previewDialogVisible" :title="t('adSettings.sync.previewTitle')" width="700px" class="preview-dialog">
      <div class="preview-header">
        <div class="preview-stats">
          <span class="preview-stat">
            <span class="preview-stat-num">{{ previewData.total || 0 }}</span>
            <span class="preview-stat-label">{{ t('adSettings.sync.totalAD') }}</span>
          </span>
          <span class="preview-stat preview-stat--synced">
            <span class="preview-stat-num">{{ previewData.synced || 0 }}</span>
            <span class="preview-stat-label">{{ t('adSettings.sync.synced') }}</span>
          </span>
        </div>
        <el-input v-model="previewSearch" :placeholder="t('adSettings.sync.searchPlaceholder')" clearable size="default" style="width: 200px">
          <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></template>
        </el-input>
      </div>
      <el-table :data="filteredPreviewUsers" max-height="400" class="preview-table" stripe>
        <el-table-column prop="username" :label="t('adSettings.sync.columns.username')" width="150" />
        <el-table-column prop="nickname" :label="t('adSettings.sync.columns.nickname')" width="150" />
        <el-table-column prop="email" :label="t('adSettings.sync.columns.email')" min-width="180" />
        <el-table-column :label="t('adSettings.sync.columns.status')" width="100" align="center">
          <template #default="{ row }">
            <span class="sync-status-tag" :class="row.synced ? 'synced' : 'not-synced'">
              {{ row.synced ? t('adSettings.sync.syncedTag') : t('adSettings.sync.notSyncedTag') }}
            </span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { AdminApi, type ADConfig, type ADUserPreview } from '@/api/admin'
import { RoleApi } from '@/api/role'

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const testLoading = ref(false)
const testResult = ref<boolean | null>(null)
const syncing = ref(false)
const resetting = ref(false)
const previewLoading = ref(false)
const previewDialogVisible = ref(false)
const previewSearch = ref('')
const previewData = ref<{ total: number; synced: number; users: ADUserPreview[] }>({ total: 0, synced: 0, users: [] })
const lastSyncTime = ref<string | null>(null)
const lastSyncResult = ref<{ created: number; updated: number; disabled: number; skipped: number } | null>(null)
const syncStats = reactive({ total: 0, synced: 0 })
const formRef = ref()
const formRef2 = ref()

const roles = ref<{ id: number; name: string; code: string }[]>([])

const form = reactive<ADConfig>({
  enabled: false,
  server: '',
  port: 389,
  use_ssl: false,
  base_dn: '',
  bind_dn: '',
  bind_password: '',
  user_filter: '(sAMAccountName=%s)',
  auto_register: false,
  default_role_id: 3
})

const filteredPreviewUsers = computed(() => {
  if (!previewSearch.value) return previewData.value.users
  const keyword = previewSearch.value.toLowerCase()
  return previewData.value.users.filter(u =>
    u.username.toLowerCase().includes(keyword) ||
    u.nickname.toLowerCase().includes(keyword) ||
    u.email?.toLowerCase().includes(keyword)
  )
})

const rules = {
  server: [{ required: true, message: t('adSettings.rules.serverRequired'), trigger: 'blur' }],
  base_dn: [{ required: true, message: t('adSettings.rules.baseDnRequired'), trigger: 'blur' }],
  bind_dn: [{ required: true, message: t('adSettings.rules.bindDnRequired'), trigger: 'blur' }]
}

const formatTime = (time: string) => {
  const d = new Date(time)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await AdminApi.getADConfig()
    if (res.code === 200) {
      Object.assign(form, res.data)
      // 保存旧密码到全局变量，用于后续操作（如测试连接、同步），留空保持不变
      if (res.data.bind_password) {
        ;(globalThis as any).__adBindPassword__ = res.data.bind_password
        form.bind_password = '' // 前端显示留空，提示"留空保持密码不变"
      }
    }
  } catch (e) {
    console.error(t('adSettings.messages.loadFailed'), e)
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  try {
    const res = await RoleApi.getAll()
    if (res.code === 200) {
      roles.value = res.data || []
    }
  } catch (e) {
    console.error(t('adSettings.messages.rolesFailed'), e)
  }
}

const loadSyncStats = async () => {
  try {
    const res = await AdminApi.getADUsers()
    if (res.code === 200) {
      syncStats.total = res.data.total
      syncStats.synced = res.data.synced
    }
  } catch (e) {
    // ignore
  }
}

const onEnabledChange = (val: boolean) => {
  if (!val) {
    form.auto_register = false
  }
}

const handleSave = async () => {
  const valid1 = await formRef.value?.validate().catch(() => false)
  const valid2 = await formRef2.value?.validate().catch(() => false)
  if (!valid1 || !valid2) return

  saving.value = true
  try {
    const res = await AdminApi.updateADConfig(form)
    if (res.code === 200) {
      ElMessage.success(t('adSettings.messages.saveSuccess'))
    } else {
      ElMessage.error(res.message || t('adSettings.messages.saveFailed'))
    }
  } catch (e: any) {
    ElMessage.error(e.message || t('adSettings.messages.saveFailed'))
  } finally {
    saving.value = false
  }
}

const handleTest = async () => {
  testResult.value = null
  testLoading.value = true
  try {
    const password = form.bind_password || (globalThis as any).__adBindPassword__ || ''
    const res = await AdminApi.testADConnection({
      server: form.server,
      port: form.port,
      use_ssl: form.use_ssl,
      base_dn: form.base_dn,
      bind_dn: form.bind_dn,
      bind_password: password,
      user_filter: form.user_filter
    })
    testResult.value = res.code === 200
    if (res.code === 200) {
      ElMessage.success(res.message || t('adSettings.messages.testSuccess'))
    } else {
      ElMessage.error(res.message || t('adSettings.messages.testFailed'))
    }
  } catch (e: any) {
    testResult.value = false
    ElMessage.error(e.message || t('adSettings.messages.testFailed'))
  } finally {
    testLoading.value = false
  }
}

const handlePreview = async () => {
  previewLoading.value = true
  previewDialogVisible.value = true
  previewSearch.value = ''
  try {
    const res = await AdminApi.getADUsers()
    if (res.code === 200) {
      previewData.value = res.data
    } else {
      ElMessage.error(res.message || t('adSettings.messages.previewFailed'))
    }
  } catch (e: any) {
    ElMessage.error(e.message || t('adSettings.messages.previewFailed'))
  } finally {
    previewLoading.value = false
  }
}

const handleSync = async (mode: 'full' | 'incremental' = 'incremental') => {
  syncing.value = true
  try {
    const res = await AdminApi.syncADUsers(mode)
    if (res.code === 200) {
      lastSyncTime.value = res.data.lastSyncAt
      lastSyncResult.value = res.data
      syncStats.total = res.data.total
      syncStats.synced = res.data.total
      const modeText = mode === 'full' ? t('adSettings.sync.fullSync') : t('adSettings.sync.incrementalSync')
      const disabledStr = res.data.deleted > 0 ? ` -${res.data.disabled} ~-${res.data.deleted}` : ` -${res.data.disabled}`
      ElMessage.success(t('adSettings.messages.syncSuccess') + ` [${modeText}]: +${res.data.created} ~${res.data.updated} ${res.data.skipped}${disabledStr}`)
      // 同步完成后刷新统计数据
      loadSyncStats()
    } else {
      ElMessage.error(res.message || t('adSettings.messages.syncFailed'))
    }
  } catch (e: any) {
    ElMessage.error(e.message || t('adSettings.messages.syncFailed'))
  } finally {
    syncing.value = false
  }
}

// 重置所有AD同步用户（删除所有Source=AD的用户）
const handleResetAll = async () => {
  await ElMessageBox.confirm(
    t('adSettings.sync.resetAllConfirm'),
    t('adSettings.sync.resetAllTitle'),
    { type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
  )
  resetting.value = true
  try {
    const res = await AdminApi.resetAllADUsers()
    if (res.code === 200) {
      ElMessage.success(t('adSettings.sync.resetAllSuccess') + `: ${res.data.count}`)
      lastSyncResult.value = null
      syncStats.total = 0
      syncStats.synced = 0
      // 重置完成后刷新统计数据
      loadSyncStats()
    } else {
      ElMessage.error(res.message || t('adSettings.messages.resetFailed'))
    }
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || t('adSettings.messages.resetFailed'))
    }
  } finally {
    resetting.value = false
  }
}

onMounted(() => {
  testResult.value = null
  loadConfig()
  loadRoles()
  loadSyncStats()
})
</script>

<style scoped lang="scss">
.ad-page {
  padding: var(--space-4);
  min-height: 100vh;
  background: var(--color-page-bg);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  overflow: visible;
}

/* 页面标题栏 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
}

.page-title {
  font-family: 'Manrope', sans-serif;
  font-size: 17px;
  font-weight: 800;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.3px;
}

.page-subtitle {
  font-size: 12px;
  color: var(--color-text-muted);
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.sync-btn {
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 状态栏 */
.status-bar {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: 12px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-xs);
}

.status-bar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-bar-right {
  display: flex;
  align-items: center;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;

  &.status--on {
    background: var(--color-success-bg);
    color: var(--color-success);
    .status-dot { background: var(--color-success); }
  }

  &.status--off {
    background: var(--gray-100);
    color: var(--gray-500);
    .status-dot { background: var(--gray-500); }
  }
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: dot-pulse 2s ease-in-out infinite;
}

@keyframes dot-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-sep {
  width: 1px;
  height: 16px;
  background: var(--color-border-light);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: var(--color-text-secondary);
  svg { color: var(--color-text-muted); }
}

.sync-stat {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--color-text-muted);
  svg { color: var(--color-text-muted); }
}

/* 总开关卡片 */
.enable-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  padding: var(--space-4) var(--space-5);
  box-shadow: var(--shadow-xs);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  transition: all 0.25s ease;

  &--active {
    border-color: rgba(0, 94, 235, 0.35);
    box-shadow: 0 4px 20px rgba(0, 94, 235, 0.08), var(--shadow-xs);
  }
}

.enable-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.enable-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.25s ease;

  &.icon-inactive {
    background: var(--gray-100);
    color: var(--gray-500);
  }

  &.icon-active {
    background: var(--color-primary-light-9);
    color: var(--color-primary);
    box-shadow: 0 4px 12px rgba(0, 94, 235, 0.25);
  }
}

.enable-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.enable-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--color-text-primary);
}

.enable-desc {
  font-size: 12px;
  color: var(--color-text-muted);
}

/* 配置内容 */
.config-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  transition: opacity 0.25s ease;

  &--disabled {
    opacity: 0.55;
    pointer-events: none;
  }
}

/* 配置卡片 */
.config-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-xs);
  overflow: hidden;

  &--warn {
    border-color: rgba(245, 158, 11, 0.3);
  }
}

.config-card-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
}

.config-card-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.config-card-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.config-card-badge {
  margin-left: auto;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  background: var(--color-warning-bg);
  color: var(--color-warning);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.config-card-body {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* 表单 */
.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 var(--space-4);
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.field-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  letter-spacing: 0.2px;
}

.required-mark {
  color: var(--color-danger);
  margin-left: 2px;
}

.field-hint {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;

  code {
    background: var(--gray-100);
    padding: 1px 5px;
    border-radius: 3px;
    font-family: monospace;
    color: var(--color-primary);
    font-size: 10.5px;
  }
}

.ssl-toggle-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.ssl-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  svg { color: var(--color-text-muted); }
}

/* 自动注册 */
.auto-reg-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-warning-bg);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: var(--radius-md);
}

.auto-reg-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.auto-reg-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-warning-text);
}

.auto-reg-desc {
  font-size: 11.5px;
  color: rgba(161, 98, 7, 0.75);
}

.auto-reg-expanded {
  padding: var(--space-3);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  margin-top: var(--space-2);
  animation: expand-in 0.25s ease both;
}

@keyframes expand-in {
  from { opacity: 0; transform: translateY(-6px); }
  to { opacity: 1; transform: translateY(0); }
}

.policy-box {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-info-bg);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: var(--radius-md);

  svg {
    color: var(--color-info);
    flex-shrink: 0;
    margin-top: 2px;
  }
}

.policy-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.policy-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--color-info);
}

.policy-list {
  margin: 0;
  padding-left: 16px;
  font-size: 12px;
  color: #1e40af;
  line-height: 1.8;

  li { margin-bottom: 1px; }
}

/* 同步区块 */
.sync-info {
  display: flex;
  align-items: center;
  gap: var(--space-5);
  flex-wrap: wrap;
}

.sync-info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sync-info-label {
  font-size: 12px;
  color: var(--color-text-muted);
}

.sync-info-value {
  font-size: 14px;
  font-weight: 800;
  color: var(--color-text-primary);

  &--success {
    color: var(--color-success);
  }
}

.sync-result-tags {
  display: flex;
  align-items: center;
  gap: 6px;
}

.sync-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 10px;

  &--created {
    background: var(--color-success-bg);
    color: var(--color-success);
  }

  &--updated {
    background: var(--color-primary-light-9);
    color: var(--color-primary);
  }

  &--skipped {
    background: var(--gray-100);
    color: var(--color-text-muted);
  }

  &--disabled {
    background: var(--color-danger-bg);
    color: var(--color-danger);
  }

  &--deleted {
    background: var(--gray-200);
    color: var(--color-text-secondary);
  }
}

.sync-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-2);
}

.sync-submit-btn {
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 操作区 */
.actions-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  padding: var(--space-3) var(--space-4);
  box-shadow: var(--shadow-xs);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

.actions-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.actions-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.test-btn {
  font-weight: 600;
}

.save-btn {
  font-weight: 700;
  min-width: 110px;
}

.test-result {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 700;

  &.result--ok {
    background: var(--color-success-bg);
    color: var(--color-success);
  }

  &.result--fail {
    background: var(--color-danger-bg);
    color: var(--color-danger);
  }
}

.result-fade-enter-active,
.result-fade-leave-active {
  transition: all 0.2s ease;
}
.result-fade-enter-from,
.result-fade-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

/* 预览弹窗 */
.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.preview-stats {
  display: flex;
  gap: var(--space-4);
}

.preview-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 8px 16px;
  background: var(--gray-100);
  border-radius: var(--radius-md);

  &--synced {
    background: var(--color-success-bg);
  }
}

.preview-stat-num {
  font-size: 18px;
  font-weight: 800;
  color: var(--color-text-primary);
}

.preview-stat-label {
  font-size: 11px;
  color: var(--color-text-muted);
}

.sync-status-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 10px;

  &.synced {
    background: var(--color-success-bg);
    color: var(--color-success);
  }

  &.not-synced {
    background: var(--gray-100);
    color: var(--color-text-muted);
  }
}

/* Element Plus overrides */
:deep(.el-input-number .el-input__inner) {
  text-align: left;
}

/* 响应式 */
@media (max-width: 768px) {
  .form-grid-2 { grid-template-columns: 1fr; }
  .enable-card { flex-direction: column; align-items: flex-start; }
  .actions-card { flex-direction: column; align-items: stretch; }
  .actions-left, .actions-right { flex-wrap: wrap; }
  .sync-info { flex-direction: column; align-items: flex-start; }
  .preview-header { flex-direction: column; align-items: flex-start; gap: var(--space-2); }
}
</style>