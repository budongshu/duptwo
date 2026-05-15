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
        <el-tag :type="form.registrationEnabled ? 'success' : 'danger'" size="small" effect="plain">
          <el-icon><UserFilled /></el-icon>
          注册 {{ form.registrationEnabled ? '已启用' : '已禁用' }}
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

          <!-- 自主注册开关 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><UserFilled /></el-icon>
              <span>允许自主注册</span>
            </div>
            <div class="config-control">
              <el-switch v-model="form.registrationEnabled" />
              <span class="control-hint" v-if="!form.registrationEnabled">
                <span style="color: #f56c6c; font-size: 12px;">关闭后，用户将无法自主注册账号</span>
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
      <el-tab-pane :label="t('security.passwordBlock.title')" name="password">
        <div class="config-section">
          <!-- 密码有效期 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><Timer /></el-icon>
              <span>{{ t('security.passwordBlock.expiryLabel') }}</span>
            </div>
            <div class="config-control config-control--inline">
              <el-switch v-model="form.passwordExpiryDaysEnabled" @change="handleExpiryToggle" />
              <span class="control-hint" v-if="form.passwordExpiryDaysEnabled">
                <el-input-number v-model="form.passwordExpiryDays" :min="1" :max="365" size="small" controls-position="right" />
                <span class="unit">{{ t('common.days') }}</span>
                <span class="input-hint">{{ t('security.passwordBlock.expiryDesc') }}</span>
              </span>
              <span class="control-hint" v-else>
                <el-tag type="info" size="small">{{ t('security.passwordBlock.noExpiry') }}</el-tag>
              </span>
            </div>
          </div>

          <!-- 密码最小长度 -->
          <div class="config-row">
            <div class="config-label">
              <el-icon><Minus /></el-icon>
              <span>{{ t('security.passwordBlock.minLengthLabel') }}</span>
            </div>
            <div class="config-control config-control--inline">
              <el-input-number v-model="form.passwordMinLength" :min="6" :max="32" size="default" controls-position="right" />
              <span class="unit">{{ t('common.characters') }}</span>
              <span class="input-hint">{{ t('security.passwordBlock.minLengthDesc') }}</span>
            </div>
          </div>

          <!-- 密码复杂度要求 -->
          <div class="config-row config-row--vertical">
            <div class="config-label">
              <el-icon><Key /></el-icon>
              <span>{{ t('security.passwordBlock.complexityLabel') }}</span>
            </div>
            <div class="config-control config-control--full">
              <div class="complexity-grid">
                <div class="complexity-item" :class="{ active: form.passwordRequireUppercase }" @click="form.passwordRequireUppercase = !form.passwordRequireUppercase">
                  <div class="complexity-check">
                    <el-icon v-if="form.passwordRequireUppercase"><Check /></el-icon>
                  </div>
                  <div class="complexity-info">
                    <span class="complexity-title">{{ t('security.passwordBlock.requireUppercase') }}</span>
                    <span class="complexity-desc">Aa Bb Cc</span>
                  </div>
                </div>
                <div class="complexity-item" :class="{ active: form.passwordRequireLowercase }" @click="form.passwordRequireLowercase = !form.passwordRequireLowercase">
                  <div class="complexity-check">
                    <el-icon v-if="form.passwordRequireLowercase"><Check /></el-icon>
                  </div>
                  <div class="complexity-info">
                    <span class="complexity-title">{{ t('security.passwordBlock.requireLowercase') }}</span>
                    <span class="complexity-desc">aa bb cc</span>
                  </div>
                </div>
                <div class="complexity-item" :class="{ active: form.passwordRequireDigit }" @click="form.passwordRequireDigit = !form.passwordRequireDigit">
                  <div class="complexity-check">
                    <el-icon v-if="form.passwordRequireDigit"><Check /></el-icon>
                  </div>
                  <div class="complexity-info">
                    <span class="complexity-title">{{ t('security.passwordBlock.requireDigit') }}</span>
                    <span class="complexity-desc">0 1 2 3 9</span>
                  </div>
                </div>
                <div class="complexity-item" :class="{ active: form.passwordRequireSpecial }" @click="form.passwordRequireSpecial = !form.passwordRequireSpecial">
                  <div class="complexity-check">
                    <el-icon v-if="form.passwordRequireSpecial"><Check /></el-icon>
                  </div>
                  <div class="complexity-info">
                    <span class="complexity-title">{{ t('security.passwordBlock.requireSpecial') }}</span>
                    <span class="complexity-desc">!@#$%^&*</span>
                  </div>
                </div>
              </div>
              <span class="input-hint">{{ t('security.passwordBlock.complexityDesc') }}</span>
            </div>
          </div>

          <!-- 密码强度预览 -->
          <div class="config-row config-row--vertical">
            <div class="config-label">
              <el-icon><View /></el-icon>
              <span>{{ t('security.passwordBlock.strengthPreview') }}</span>
            </div>
            <div class="config-control config-control--full">
              <div class="strength-preview">
                <div class="strength-input-wrapper">
                  <el-input
                    v-model="previewPassword"
                    type="password"
                    :placeholder="t('security.passwordBlock.previewPlaceholder')"
                    show-password
                    clearable
                  />
                </div>
                <div class="strength-meter">
                  <div class="strength-bar" :style="{ width: strengthScore.width + '%', background: strengthScore.color }"></div>
                </div>
                <div class="strength-labels">
                  <span class="strength-score" :style="{ color: strengthScore.color }">{{ strengthScore.label }}</span>
                  <span class="strength-rules">
                    <span v-for="(rule, idx) in strengthRules" :key="idx" :class="{ met: rule.met }">
                      <el-icon v-if="rule.met"><Check /></el-icon>
                      <el-icon v-else><Close /></el-icon>
                      {{ rule.text }}
                    </span>
                  </span>
                </div>
              </div>
              <span class="input-hint">{{ t('security.passwordBlock.previewDesc') }}</span>
            </div>
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
import { Lock, Warning, CircleClose, Guide, Picture, Timer, User, UserFilled, Monitor, Key, Connection, Clock, Minus, Check, Close, View } from '@element-plus/icons-vue'
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
const form = reactive<SecuritySettings & { passwordExpiryDaysEnabled: boolean }>({
  id: 0,
  captchaEnabled: true,
  captchaMinLen: 3,
  registrationEnabled: true,
  inactiveAutoDisable: false,
  inactiveDaysThreshold: 90,
  userLoginMaxAttempts: 5,
  userLoginLockMinutes: 30,
  ipLoginMaxAttempts: 20,
  ipLoginLockMinutes: 60,
  ipWhitelist: '',
  ipBlacklist: '',
  passwordExpiryDays: 0,
  passwordExpiryDaysEnabled: false,
  passwordMinLength: 8,
  passwordRequireUppercase: false,
  passwordRequireLowercase: false,
  passwordRequireDigit: false,
  passwordRequireSpecial: false,
  sessionTimeoutHours: 24,
})

// 密码预览
const previewPassword = ref('')

// 密码有效期开关
function handleExpiryToggle(val: boolean) {
  if (val && form.passwordExpiryDays === 0) {
    form.passwordExpiryDays = 90
  }
}

// 计算密码强度
const strengthScore = computed(() => {
  if (!previewPassword.value) {
    return { label: t('security.passwordBlock.strengthEmpty'), width: 0, color: '#c0c4cc' }
  }
  let score = 0
  const pwd = previewPassword.value

  // 长度得分
  if (pwd.length >= form.passwordMinLength) score += 20
  else if (pwd.length >= form.passwordMinLength - 2) score += 10

  // 复杂度得分
  if (/[A-Z]/.test(pwd) && form.passwordRequireUppercase) score += 20
  else if (/[A-Z]/.test(pwd)) score += 5

  if (/[a-z]/.test(pwd) && form.passwordRequireLowercase) score += 20
  else if (/[a-z]/.test(pwd)) score += 5

  if (/[0-9]/.test(pwd) && form.passwordRequireDigit) score += 20
  else if (/[0-9]/.test(pwd)) score += 5

  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd) && form.passwordRequireSpecial) score += 20
  else if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd)) score += 5

  // 长度额外加分
  if (pwd.length >= 12) score += 10
  if (pwd.length >= 16) score += 5

  if (score >= 80) return { label: t('security.passwordBlock.strengthStrong'), width: 100, color: '#67c23a' }
  if (score >= 50) return { label: t('security.passwordBlock.strengthMedium'), width: 60, color: '#e6a23c' }
  return { label: t('security.passwordBlock.strengthWeak'), width: 30, color: '#f56c6c' }
})

// 密码规则检查
const strengthRules = computed(() => {
  const pwd = previewPassword.value
  return [
    { text: t('security.passwordBlock.ruleLength', { n: form.passwordMinLength }), met: pwd.length >= form.passwordMinLength },
    { text: t('security.passwordBlock.ruleUppercase'), met: form.passwordRequireUppercase ? /[A-Z]/.test(pwd) : true },
    { text: t('security.passwordBlock.ruleLowercase'), met: form.passwordRequireLowercase ? /[a-z]/.test(pwd) : true },
    { text: t('security.passwordBlock.ruleDigit'), met: form.passwordRequireDigit ? /[0-9]/.test(pwd) : true },
    { text: t('security.passwordBlock.ruleSpecial'), met: form.passwordRequireSpecial ? /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd) : true },
  ]
})

const lockDrawerTitle = computed(() =>
  lockDrawerType.value === 'users' ? t('security.lockedUsersList') : t('security.blockedIPsList')
)

// 加载安全设置
async function loadSettings() {
  loadingSettings.value = true
  try {
    const resp = await AdminApi.getSecuritySettings()
    const data = resp.data as SecuritySettings
    // 使用单独赋值确保 Vue 响应式正确更新
    form.id = data.id ?? 0
    form.captchaEnabled = data.captchaEnabled ?? true
    form.captchaMinLen = data.captchaMinLen ?? 3
    form.registrationEnabled = data.registrationEnabled ?? true
    form.inactiveAutoDisable = data.inactiveAutoDisable ?? false
    form.inactiveDaysThreshold = data.inactiveDaysThreshold ?? 90
    form.userLoginMaxAttempts = data.userLoginMaxAttempts ?? 5
    form.userLoginLockMinutes = data.userLoginLockMinutes ?? 30
    form.ipLoginMaxAttempts = data.ipLoginMaxAttempts ?? 20
    form.ipLoginLockMinutes = data.ipLoginLockMinutes ?? 60
    form.ipWhitelist = data.ipWhitelist ?? ''
    form.ipBlacklist = data.ipBlacklist ?? ''
    form.passwordExpiryDays = data.passwordExpiryDays ?? 0
    form.passwordExpiryDaysEnabled = (data.passwordExpiryDays ?? 0) > 0
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
    const resp = await AdminApi.getSecurityOverview()
    Object.assign(overview, resp.data)
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
    const saveData = {
      ...form,
      passwordExpiryDays: form.passwordExpiryDaysEnabled ? form.passwordExpiryDays : 0,
    }
    await AdminApi.updateSecuritySettings(saveData)
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
      const resp = await AdminApi.getLockedUsers()
      lockList.value = resp.data ?? []
    } else {
      const resp = await AdminApi.getLockedIPs()
      lockList.value = resp.data ?? []
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

/* 复杂度选择器 */
.complexity-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  width: 100%;
}
.complexity-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: #f5f7fa;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}
.complexity-item:hover {
  border-color: #c0c4cc;
}
.complexity-item.active {
  background: #ecf5ff;
  border-color: #409eff;
}
.complexity-check {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 2px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
}
.complexity-item.active .complexity-check {
  background: #409eff;
  border-color: #409eff;
  color: #fff;
}
.complexity-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.complexity-title {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
}
.complexity-desc {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
}

/* 密码强度预览 */
.strength-preview {
  width: 100%;
  max-width: 480px;
}
.strength-input-wrapper {
  margin-bottom: 12px;
}
.strength-meter {
  height: 6px;
  background: #e4e7ed;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 10px;
}
.strength-bar {
  height: 100%;
  border-radius: 3px;
  transition: all 0.3s ease;
}
.strength-labels {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
}
.strength-score {
  font-size: 13px;
  font-weight: 600;
}
.strength-rules {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.strength-rules span {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #c0c4cc;
  transition: color 0.2s;
}
.strength-rules span.met {
  color: #67c23a;
}
.strength-rules .el-icon {
  font-size: 12px;
}

/* 响应式 */
@media (max-width: 1024px) {
  .overview-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .checkbox-grid {
    grid-template-columns: 1fr;
  }
  .complexity-grid {
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
