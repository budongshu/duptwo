<template>
  <div class="security-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('security.title') }}</h1>
        <span class="page-subtitle">{{ t('security.subtitle') }}</span>
      </div>
      <el-button type="primary" :loading="saving" @click="handleSave">
        <el-icon v-if="!saving"><Lock /></el-icon>
        {{ t('security.saveAll') }}
      </el-button>
    </header>

    <!-- 安全总览卡片 -->
    <div class="overview-cards" v-loading="loadingOverview">
      <div class="overview-card overview-card--warn" @click="activeTab = 'login'; showLockDrawer('users')">
        <el-icon class="overview-icon"><Warning /></el-icon>
        <div class="overview-info">
          <div class="overview-num">{{ overview.totalLockedUsers }}</div>
          <div class="overview-label">{{ t('security.lockedUsers') }}</div>
        </div>
      </div>
      <div class="overview-card overview-card--danger" @click="activeTab = 'ip'; showLockDrawer('ips')">
        <el-icon class="overview-icon"><CircleClose /></el-icon>
        <div class="overview-info">
          <div class="overview-num">{{ overview.totalLockedIPs }}</div>
          <div class="overview-label">{{ t('security.blockedIPs') }}</div>
        </div>
      </div>
      <div class="overview-card overview-card--success" @click="activeTab = 'ip'">
        <el-icon class="overview-icon"><Guide /></el-icon>
        <div class="overview-info">
          <div class="overview-num">{{ overview.whitelistCount || 0 }}</div>
          <div class="overview-label">{{ t('security.ipWhitelist') }}</div>
        </div>
      </div>
      <div class="overview-card overview-card--info" @click="activeTab = 'ip'">
        <el-icon class="overview-icon"><Lock /></el-icon>
        <div class="overview-info">
          <div class="overview-num">{{ overview.blacklistCount || 0 }}</div>
          <div class="overview-label">{{ t('security.ipBlacklist') }}</div>
        </div>
      </div>
    </div>

    <!-- 当前配置状态 -->
    <div class="status-panel">
      <div class="status-title">
        <el-icon><Monitor /></el-icon>
        <span>{{ t('security.currentStatus') }}</span>
      </div>
      <div class="status-tags">
        <el-tag :type="form.captchaEnabled ? 'success' : 'info'" size="small" effect="plain">
          <el-icon><Picture /></el-icon>
          验证码 {{ form.captchaEnabled ? '已启用' : '已禁用' }}
        </el-tag>
        <el-tag :type="form.inactiveAutoDisable ? 'warning' : 'info'" size="small" effect="plain">
          <el-icon><Timer /></el-icon>
          不活跃 {{ form.inactiveAutoDisable ? form.inactiveDaysThreshold + '天自动禁用' : '未启用' }}
        </el-tag>
        <el-tag type="info" size="small" effect="plain">
          <el-icon><User /></el-icon>
          登录锁定 {{ form.userLoginMaxAttempts }}次 / {{ form.userLoginLockMinutes }}分钟
        </el-tag>
        <el-tag type="info" size="small" effect="plain">
          <el-icon><Key /></el-icon>
          密码 {{ form.passwordMinLength }}位
          <template v-if="form.passwordRequireUppercase || form.passwordRequireLowercase || form.passwordRequireDigit || form.passwordRequireSpecial">
            ({{
              (form.passwordRequireUppercase ? '大写 ' : '') +
              (form.passwordRequireLowercase ? '小写 ' : '') +
              (form.passwordRequireDigit ? '数字 ' : '') +
              (form.passwordRequireSpecial ? '特殊' : '')
            }})
          </template>
        </el-tag>
        <el-tag type="info" size="small" effect="plain">
          <el-icon><Clock /></el-icon>
          会话 {{ form.sessionTimeoutHours }}小时
        </el-tag>
        <el-tag type="info" size="small" effect="plain">
          <el-icon><Connection /></el-icon>
          IP限制 {{ form.ipLoginMaxAttempts }}次 / {{ form.ipLoginLockMinutes }}分钟
        </el-tag>
      </div>
    </div>

    <!-- 配置标签页 -->
    <div class="security-tabs" v-loading="loadingSettings">
      <el-tabs v-model="activeTab">
      <!-- 登录安全 -->
      <el-tab-pane :label="t('security.loginBlock.title')" name="login">
        <div class="config-section">
          <!-- 验证码 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><Picture /></el-icon>
              <span>{{ t('security.loginBlock.captchaLabel') }}</span>
            </div>
            <div class="config-control">
              <el-switch v-model="form.captchaEnabled" />
              <span class="control-hint" v-if="form.captchaEnabled">
                {{ t('security.loginBlock.enableCaptchaDesc') }}
                <el-input-number v-model="form.captchaMinLen" :min="1" :max="10" size="small" controls-position="right" />
                {{ t('security.loginBlock.captchaThreshold') }}
              </span>
            </div>
          </div>

          <!-- 不活跃自动禁用 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><Timer /></el-icon>
              <span>{{ t('security.loginBlock.inactiveLabel') }}</span>
            </div>
            <div class="config-control">
              <el-switch v-model="form.inactiveAutoDisable" />
              <span class="control-hint" v-if="form.inactiveAutoDisable">
                {{ t('security.loginBlock.inactiveDaysDesc') }}
                <el-input-number v-model="form.inactiveDaysThreshold" :min="1" :max="365" size="small" controls-position="right" />
                {{ t('common.days') }}
              </span>
            </div>
          </div>

          <!-- 用户登录限制 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><User /></el-icon>
              <span>{{ t('security.loginBlock.userLimitLabel') }}</span>
            </div>
            <div class="config-control config-control--inline">
              <div class="inline-item">
                <span class="inline-label">{{ t('security.loginBlock.maxAttempts') }}</span>
                <el-input-number v-model="form.userLoginMaxAttempts" :min="1" :max="20" size="small" controls-position="right" />
              </div>
              <div class="inline-item">
                <span class="inline-label">{{ t('security.loginBlock.lockDuration') }}</span>
                <el-input-number v-model="form.userLoginLockMinutes" :min="1" :max="10080" size="small" controls-position="right" />
                <span class="unit">{{ t('common.minutes') }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- IP访问控制 -->
      <el-tab-pane :label="t('security.ipBlock.title')" name="ip">
        <div class="config-section">
          <!-- IP登录限制 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><Monitor /></el-icon>
              <span>{{ t('security.ipBlock.ipLimitLabel') }}</span>
            </div>
            <div class="config-control config-control--inline">
              <div class="inline-item">
                <span class="inline-label">{{ t('security.ipBlock.maxAttempts') }}</span>
                <el-input-number v-model="form.ipLoginMaxAttempts" :min="1" :max="100" size="small" controls-position="right" />
              </div>
              <div class="inline-item">
                <span class="inline-label">{{ t('security.ipBlock.banDuration') }}</span>
                <el-input-number v-model="form.ipLoginLockMinutes" :min="1" :max="10080" size="small" controls-position="right" />
                <span class="unit">{{ t('common.minutes') }}</span>
              </div>
            </div>
          </div>

          <!-- IP白名单 -->
          <div class="config-row config-row--vertical">
            <div class="config-label">
              <el-icon><Guide /></el-icon>
              <span>{{ t('security.ipBlock.whitelistLabel') }}</span>
            </div>
            <div class="config-control config-control--full">
              <el-input
                v-model="form.ipWhitelist"
                type="textarea"
                :rows="2"
                :placeholder="t('security.ipBlock.whitelistPlaceholder')"
                clearable
              />
              <span class="input-hint">{{ t('security.ipBlock.whitelistDesc') }}</span>
            </div>
          </div>

          <!-- IP黑名单 -->
          <div class="config-row config-row--vertical">
            <div class="config-label">
              <el-icon><CircleClose /></el-icon>
              <span>{{ t('security.ipBlock.blacklistLabel') }}</span>
            </div>
            <div class="config-control config-control--full">
              <el-input
                v-model="form.ipBlacklist"
                type="textarea"
                :rows="2"
                :placeholder="t('security.ipBlock.blacklistPlaceholder')"
                clearable
              />
              <span class="input-hint">{{ t('security.ipBlock.blacklistDesc') }}</span>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 密码策略 -->
      <el-tab-pane :label="t('security.passwordBlock.title')" name="password" disabled>
        <div class="config-section">
          <div class="pwd-fixed-notice">
            <el-icon><InfoFilled /></el-icon>
            <span>密码要求已固定：至少8位，必须包含大小写字母、数字和特殊字符</span>
          </div>
        </div>
      </el-tab-pane>

      <!-- 会话管理 -->
      <el-tab-pane :label="t('security.sessionBlock.title')" name="session">
        <div class="config-section">
          <div class="config-row">
            <div class="config-label">
              <el-icon><Clock /></el-icon>
              <span>{{ t('security.sessionBlock.timeoutLabel') }}</span>
            </div>
            <div class="config-control config-control--inline">
              <el-input-number v-model="form.sessionTimeoutHours" :min="1" :max="168" size="default" controls-position="right" />
              <span class="unit">{{ t('common.hours') }}</span>
              <span class="input-hint">{{ t('security.sessionBlock.timeoutDesc') }}</span>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    </div>

    <!-- 锁定列表抽屉 -->
    <el-drawer v-model="lockDrawerVisible" :title="lockDrawerTitle" size="420px" direction="rtl">
      <div class="lock-drawer-body">
        <el-empty v-if="lockList.length === 0" :description="lockDrawerType === 'users' ? t('security.noLockedUsers') : t('security.noBlockedIPs')" />
        <div v-else class="lock-list">
          <div v-for="item in lockList" :key="item.id" class="lock-item">
            <div class="lock-item-info">
              <el-tag :type="lockDrawerType === 'users' ? 'warning' : 'danger'" size="small">
                {{ lockDrawerType === 'users' ? t('security.lockTag') : 'Banned' }}
              </el-tag>
              <span class="lock-target">{{ item.target }}</span>
              <span class="lock-count">{{ t('security.failCount', { count: item.failCount }) }}</span>
            </div>
            <el-button type="danger" size="small" plain :title="t('security.unlock')" @click="handleUnlock(item)">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 9.9-1"/></svg>
            </el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Lock, Warning, CircleClose, Guide, Picture, Timer, User, Monitor, Key, Connection, Clock, InfoFilled } from '@element-plus/icons-vue'
import { AdminApi, SecuritySettings } from '@/api/admin'
import { usePasswordPolicy } from '@/composables/usePasswordPolicy'

const { clearCache } = usePasswordPolicy()

const { t } = useI18n()

// 状态
const activeTab = ref('login')
const saving = ref(false)
const loadingOverview = ref(false)
const loadingSettings = ref(false)
const lockDrawerVisible = ref(false)
const lockDrawerType = ref<'users' | 'ips'>('users')
const lockList = ref<any[]>([])

// 总览数据
const overview = reactive({
  totalLockedUsers: 0,
  totalLockedIPs: 0,
  whitelistCount: 0,
  blacklistCount: 0,
})

// 表单数据
const form = reactive<SecuritySettings>({
  id: 0,
  captchaEnabled: true,
  captchaMinLen: 3,
  inactiveAutoDisable: false,
  inactiveDaysThreshold: 90,
  userLoginMaxAttempts: 5,
  userLoginLockMinutes: 30,
  ipLoginMaxAttempts: 20,
  ipLoginLockMinutes: 60,
  ipWhitelist: '',
  ipBlacklist: '',
  passwordExpiryDays: 0,
  passwordMinLength: 8,
  passwordRequireUppercase: false,
  passwordRequireLowercase: false,
  passwordRequireDigit: false,
  passwordRequireSpecial: false,
  sessionTimeoutHours: 24,
})

const lockDrawerTitle = computed(() =>
  lockDrawerType.value === 'users' ? t('security.lockedUsersList') : t('security.blockedIPsList')
)

// 加载安全设置
async function loadSettings() {
  loadingSettings.value = true
  try {
    const data = await AdminApi.getSecuritySettings()
    // 使用单独赋值确保 Vue 响应式正确更新
    form.id = data.id ?? 0
    form.captchaEnabled = data.captchaEnabled ?? true
    form.captchaMinLen = data.captchaMinLen ?? 3
    form.inactiveAutoDisable = data.inactiveAutoDisable ?? false
    form.inactiveDaysThreshold = data.inactiveDaysThreshold ?? 90
    form.userLoginMaxAttempts = data.userLoginMaxAttempts ?? 5
    form.userLoginLockMinutes = data.userLoginLockMinutes ?? 30
    form.ipLoginMaxAttempts = data.ipLoginMaxAttempts ?? 20
    form.ipLoginLockMinutes = data.ipLoginLockMinutes ?? 60
    form.ipWhitelist = data.ipWhitelist ?? ''
    form.ipBlacklist = data.ipBlacklist ?? ''
    form.passwordExpiryDays = data.passwordExpiryDays ?? 0
    form.passwordMinLength = data.passwordMinLength ?? 8
    form.passwordRequireUppercase = data.passwordRequireUppercase ?? false
    form.passwordRequireLowercase = data.passwordRequireLowercase ?? false
    form.passwordRequireDigit = data.passwordRequireDigit ?? false
    form.passwordRequireSpecial = data.passwordRequireSpecial ?? false
    form.sessionTimeoutHours = data.sessionTimeoutHours ?? 24
  } catch {
    ElMessage.error(t('security.messages.loadSettingsFailed'))
  } finally {
    loadingSettings.value = false
  }
}

// 加载总览
async function loadOverview() {
  loadingOverview.value = true
  try {
    const data = await AdminApi.getSecurityOverview()
    Object.assign(overview, data)
  } catch {
    // 总览加载失败不影响主功能
  } finally {
    loadingOverview.value = false
  }
}

// 保存设置
async function handleSave() {
  saving.value = true
  try {
    await AdminApi.updateSecuritySettings(form)
    ElMessage.success(t('security.messages.saveSuccess'))
    loadOverview()
    loadSettings()
    clearCache() // 清除密码策略缓存
  } catch {
    ElMessage.error(t('security.messages.saveFailed'))
  } finally {
    saving.value = false
  }
}

// 显示锁定抽屉
async function showLockDrawer(type: 'users' | 'ips') {
  lockDrawerType.value = type
  lockDrawerVisible.value = true
  try {
    if (type === 'users') {
      lockList.value = await AdminApi.getLockedUsers()
    } else {
      lockList.value = await AdminApi.getLockedIPs()
    }
  } catch {
    ElMessage.error(type === 'users' ? t('security.messages.loadLockedUsersFailed') : t('security.messages.loadLockedIPsFailed'))
    lockList.value = []
  }
}

// 解锁
async function handleUnlock(item: any) {
  const target = item.target
  const confirmMsg = lockDrawerType.value === 'users'
    ? t('security.messages.unlockUserConfirm', { username: target })
    : t('security.messages.unlockIPConfirm', { ip: target })

  try {
    if (lockDrawerType.value === 'users') {
      await AdminApi.unlockUser(target)
      ElMessage.success(t('security.messages.userUnlocked'))
    } else {
      await AdminApi.unlockIP(target)
      ElMessage.success(t('security.messages.ipUnlocked'))
    }
    lockList.value = lockList.value.filter(i => i.target !== target)
    loadOverview()
  } catch {
    ElMessage.error(t('security.messages.unlockFailed'))
  }
}

onMounted(() => {
  loadSettings()
  loadOverview()
})
</script>

<style scoped>
.security-page {
  padding: 24px;
  min-height: 100vh;
  background: #f5f7fa;
}

/* 页面标题 */
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
  font-weight: 600;
  color: #303133;
  margin: 0;
}
.page-subtitle {
  font-size: 13px;
  color: #909399;
}

/* 总览卡片 */
.overview-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}
.overview-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;
  border: 1px solid #ebeef5;
}
.overview-card:hover {
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  transform: translateY(-1px);
}
.overview-icon {
  font-size: 28px;
  flex-shrink: 0;
}
.overview-card--warn .overview-icon { color: #e6a23c; }
.overview-card--danger .overview-icon { color: #f56c6c; }
.overview-card--success .overview-icon { color: #67c23a; }
.overview-card--info .overview-icon { color: #409eff; }
.overview-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.overview-num {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}
.overview-label {
  font-size: 13px;
  color: #909399;
}

/* 当前配置状态 */
.status-panel {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  padding: 14px 20px;
  margin-bottom: 16px;
}

.status-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;

  .el-icon { color: #409eff; }
}

.status-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.el-tag) {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  cursor: default;

  .el-icon { font-size: 12px; }
}

/* 标签页 */
.security-tabs {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  border: 1px solid #ebeef5;
}

/* 配置区块 */
.config-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.config-row {
  display: flex;
  align-items: flex-start;
  gap: 20px;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}
.config-row:last-child {
  border-bottom: none;
}
.config-row--vertical {
  flex-direction: column;
  gap: 10px;
}
.config-label {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 160px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  padding-top: 4px;
}
.config-label .el-icon {
  color: #409eff;
}
.config-control {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.config-control--full {
  width: 100%;
}
.config-control--inline {
  flex-wrap: wrap;
  gap: 20px;
}
.control-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #909399;
}
.inline-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.inline-label {
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
}
.unit {
  font-size: 13px;
  color: #909399;
}
.input-hint {
  font-size: 12px;
  color: #c0c4cc;
  width: 100%;
}

/* 密码策略固定提示 */
.pwd-fixed-notice {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 18px;
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  font-size: 13px;
  color: #0369a1;

  .el-icon { font-size: 18px; color: #0ea5e9; }
}

/* 复选框网格 */
.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px 24px;
  width: 100%;
}
:deep(.el-checkbox) {
  display: flex;
  align-items: flex-start;
  white-space: normal;
  line-height: 1.6;
}
:deep(.el-checkbox__label) {
  font-size: 13px;
  color: #606266;
}
.checkbox-sub {
  display: block;
  font-size: 12px;
  color: #c0c4cc;
  line-height: 1.4;
  margin-top: 2px;
}

/* 锁定列表 */
.lock-drawer-body {
  padding: 0 4px;
}
.lock-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.lock-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f9fafb;
  border-radius: 6px;
  border: 1px solid #ebeef5;
}
.lock-item-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.lock-target {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  font-family: monospace;
}
.lock-count {
  font-size: 12px;
  color: #909399;
}

/* 响应式 */
@media (max-width: 1024px) {
  .overview-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .checkbox-grid {
    grid-template-columns: 1fr;
  }
}
@media (max-width: 768px) {
  .overview-cards {
    grid-template-columns: 1fr 1fr;
  }
  .config-row {
    flex-direction: column;
    gap: 10px;
  }
  .config-label {
    min-width: unset;
  }
}
</style>
