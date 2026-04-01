<template>
  <div class="profile-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('profile.title') }}</h1>
        <span class="page-subtitle">{{ t('profile.subtitle') }}</span>
      </div>
    </header>

    <div class="profile-layout" v-loading="loading">
      <!-- 左侧：个人信息卡片 -->
      <div class="profile-card profile-card--main">
        <div class="card-header">
          <div class="card-header-avatar">
            <div class="avatar-large" :style="{ background: avatarColor }">
              {{ avatarInitials }}
            </div>
            <div class="avatar-info">
              <span class="avatar-name">{{ form.nickname || form.username }}</span>
              <span class="avatar-username">@{{ form.username }}</span>
            </div>
          </div>
          <span class="mfa-badge" :class="form.mfaEnabled ? 'mfa-badge--on' : 'mfa-badge--off'">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
            {{ form.mfaEnabled ? t('profile.mfaEnabled') : t('profile.mfaDisabled') }}
          </span>
        </div>

        <div class="card-body">
          <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
            <div class="form-row-2">
              <el-form-item :label="t('profile.username')" prop="username">
                <el-input v-model="form.username" disabled>
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
                </el-input>
                <div class="field-tag field-tag--local" v-if="form.source === 'LOCAL'">{{ t('profile.localAccount') }}</div>
                <div class="field-tag field-tag--ad" v-else>{{ t('profile.adAccount') }}</div>
              </el-form-item>
              <el-form-item :label="t('profile.nickname')" prop="nickname">
                <el-input v-model="form.nickname" :placeholder="t('profile.nicknamePlaceholder')">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3"/></svg></template>
                </el-input>
              </el-form-item>
            </div>
            <div class="form-row-2">
              <el-form-item :label="t('profile.email')" prop="email">
                <el-input v-model="form.email" :disabled="form.source === 'AD'" :placeholder="t('profile.emailPlaceholder')">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg></template>
                </el-input>
              </el-form-item>
              <el-form-item :label="t('profile.phone')">
                <el-input v-model="form.phone" :placeholder="t('profile.phonePlaceholder')">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg></template>
                </el-input>
              </el-form-item>
            </div>
          </el-form>
        </div>

        <div class="card-foot">
          <el-button type="primary" :loading="submitting" @click="handleSave">{{ t('profile.saveChanges') }}</el-button>
        </div>
      </div>

      <!-- 右侧：其他信息 -->
      <div class="profile-sidebar">

        <!-- 安全卡片 -->
        <div class="profile-card">
          <div class="card-header-simple">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            {{ t('profile.accountSecurity') }}
          </div>
          <div class="card-body">
            <!-- 修改密码 -->
            <div class="info-row" v-if="form.source !== 'AD'">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <span>{{ t('profile.changePassword') }}</span>
              </div>
              <el-button size="small" @click="showChangePwd = true">{{ t('profile.modify') }}</el-button>
            </div>
            <div class="info-row info-row--disabled" v-else>
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <span>{{ t('profile.changePassword') }}</span>
              </div>
              <span class="info-row-tip">{{ t('profile.adManaged') }}</span>
            </div>

            <!-- MFA 启用 -->
            <div class="info-row" v-if="form.source !== 'AD'">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                <span>{{ t('profile.mfaAuth') }}</span>
              </div>
              <el-button v-if="!form.mfaEnabled" type="primary" size="small" plain @click="handleEnableMFA">{{ t('profile.enableMfa') }}</el-button>
              <el-button v-else type="warning" size="small" plain @click="handleDisableMFA">{{ t('profile.resetMfa') }}</el-button>
            </div>

            <!-- 角色信息 -->
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                <span>{{ t('profile.role') }}</span>
              </div>
              <span class="info-row-value">{{ currentUser?.roleName || '-' }}</span>
            </div>

            <!-- 最后登录 -->
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                <span>{{ t('profile.lastLogin') }}</span>
              </div>
              <span class="info-row-value">{{ currentUser?.lastLoginAt ? formatDate(currentUser.lastLoginAt) : '-' }}</span>
            </div>
          </div>
        </div>

        <!-- 账号信息 -->
        <div class="profile-card">
          <div class="card-header-simple">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
            {{ t('profile.accountInfo') }}
          </div>
          <div class="card-body">
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>{{ t('profile.accountSource') }}</span>
              </div>
              <span class="info-row-value source-tag" :class="form.source === 'AD' ? 'source-tag--ad' : 'source-tag--local'">
                {{ form.source === 'AD' ? t('profile.adDomain') : t('profile.local') }}
              </span>
            </div>
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>{{ t('profile.accountStatus') }}</span>
              </div>
              <span class="info-row-value">
                <span class="status-dot" :class="form.status === 'active' ? 'status-dot--on' : 'status-dot--off'"></span>
                {{ form.status === 'active' ? t('profile.active') : t('profile.disabled') }}
              </span>
            </div>
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>{{ t('profile.registerTime') }}</span>
              </div>
              <span class="info-row-value">{{ currentUser?.createdAt ? formatDate(currentUser.createdAt) : '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showChangePwd" width="460px" destroy-on-close class="pwd-dialog">
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--warn">{{ t('profile.password') }}</span>
          <span class="dialog-title-text">{{ t('profile.changePassword') }}</span>
        </div>
      </template>
      <div class="pwd-dialog-body">
        <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-position="top">
          <el-form-item :label="t('profile.currentPassword')" prop="oldPassword">
            <el-input v-model="pwdForm.oldPassword" type="password" :placeholder="t('profile.currentPasswordPlaceholder')" show-password size="large" />
          </el-form-item>
          <el-form-item :label="t('profile.newPassword')" prop="newPassword">
            <el-input v-model="pwdForm.newPassword" type="password" :placeholder="t('profile.newPasswordPlaceholder', { min: policy.passwordMinLength })" show-password size="large" />
          </el-form-item>

          <!-- 密码要求清单 -->
          <div class="pwd-requirements" v-if="pwdForm.newPassword">
            <div class="req-title">{{ t('profile.passwordRequirements') }}</div>
            <div class="req-list">
              <div class="req-item" :class="{ satisfied: pwdChecks.length }">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>至少 {{ policy.passwordMinLength }} 个字符</span>
              </div>
              <div class="req-item" :class="{ satisfied: pwdChecks.upper }" v-if="policy.passwordRequireUppercase">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含大写字母 (A-Z)</span>
              </div>
              <div class="req-item" :class="{ satisfied: pwdChecks.lower }" v-if="policy.passwordRequireLowercase">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含小写字母 (a-z)</span>
              </div>
              <div class="req-item" :class="{ satisfied: pwdChecks.number }" v-if="policy.passwordRequireDigit">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含数字 (0-9)</span>
              </div>
              <div class="req-item" :class="{ satisfied: pwdChecks.special }" v-if="policy.passwordRequireSpecial">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含特殊字符 (!@#$%...)</span>
              </div>
            </div>
            <!-- 密码强度指示 -->
            <div class="strength-section">
              <div class="strength-label">
                <span>密码强度</span>
                <span class="strength-text" :class="strengthClass">{{ strengthLabel }}</span>
              </div>
              <div class="strength-bar-wrap">
                <div class="strength-bar-segments">
                  <div class="segment" :class="{ active: strengthScore >= 1, [strengthClass]: strengthScore >= 1 }"></div>
                  <div class="segment" :class="{ active: strengthScore >= 2, [strengthClass]: strengthScore >= 2 }"></div>
                  <div class="segment" :class="{ active: strengthScore >= 3, [strengthClass]: strengthScore >= 3 }"></div>
                  <div class="segment" :class="{ active: strengthScore >= 4, [strengthClass]: strengthScore >= 4 }"></div>
                </div>
              </div>
            </div>
          </div>

          <el-form-item :label="t('profile.confirmNewPassword')" prop="confirmPassword">
            <el-input v-model="pwdForm.confirmPassword" type="password" :placeholder="t('profile.confirmPasswordPlaceholder')" show-password size="large" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button size="default" @click="showChangePwd = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" size="default" :loading="pwdLoading" @click="handleChangePassword">{{ t('profile.confirmChange') }}</el-button>
      </template>
    </el-dialog>

    <!-- 启用 MFA 弹窗 -->
    <el-dialog v-model="showMFADialog" width="420px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--mfa">MFA</span>
          <span class="dialog-title-text">{{ t('profile.enableMfaTitle') }}</span>
        </div>
      </template>
      <div class="mfa-content">
        <p class="mfa-tip">{{ t('profile.mfaTip') }}</p>
        <div class="qr-box" v-if="mfaSecret">
          <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="MFA QR" class="qr-img" />
          <div class="mfa-secret">{{ t('profile.mfaSecret') }}: {{ mfaSecret }}</div>
        </div>
        <el-form ref="mfaFormRef" :model="mfaForm" :rules="mfaRules" label-position="top">
          <el-form-item :label="t('profile.mfaCode')" prop="code">
            <el-input v-model="mfaForm.code" :placeholder="t('profile.mfaCodePlaceholder')" maxlength="6" size="large" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button size="small" @click="showMFADialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" size="small" :loading="mfaLoading" @click="handleConfirmMFA">{{ t('profile.confirmBind') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { UserApi, type User } from '@/api/user'
import { AuthApi } from '@/api/auth'
import { usePasswordPolicy } from '@/composables/usePasswordPolicy'

const { t } = useI18n()
const { loadPolicy, checkPassword, passwordStrength: calcStrength, policy } = usePasswordPolicy()

const loading = ref(false)
const submitting = ref(false)
const showChangePwd = ref(false)
const showMFADialog = ref(false)
const mfaLoading = ref(false)
const pwdLoading = ref(false)
const formRef = ref<FormInstance>()
const pwdFormRef = ref<FormInstance>()
const mfaFormRef = ref<FormInstance>()

const currentUser = ref<User | null>(null)

const form = reactive({
  id: undefined as number | undefined,
  username: '',
  nickname: '',
  email: '',
  phone: '',
  status: 'active',
  source: 'LOCAL',
  mfaEnabled: false,
  roleName: ''
})

const pwdForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })
const mfaForm = reactive({ code: '' })
const mfaSecret = ref('')
const qrCodeUrl = ref('')

const avatarColor = computed(() => {
  if (!form.username) return '#6366f1'
  const colors = ['#6366f1', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#ef4444', '#14b8a6']
  let hash = 0
  for (const c of form.username) hash = c.charCodeAt(0) + ((hash << 5) - hash)
  return colors[Math.abs(hash) % colors.length]
})

const avatarInitials = computed(() => {
  const name = form.nickname || form.username
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
})

const formRules: FormRules = {
  email: [
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ]
}

const validateConfirm = (rule: any, value: string, callback: any) => {
  if (value !== pwdForm.newPassword) callback(new Error(t('profile.passwordMismatch')))
  else callback()
}

// 密码验证规则（使用全局密码策略）
const pwdRules = computed<FormRules>(() => ({
  oldPassword: [{ required: true, message: t('profile.currentPasswordRequired'), trigger: 'blur' }],
  newPassword: [
    { required: true, message: t('profile.newPasswordRequired'), trigger: 'blur' },
    { min: policy.value.passwordMinLength, message: t('profile.passwordMinLength', { min: policy.value.passwordMinLength }), trigger: 'blur' }
  ],
  confirmPassword: [{ required: true, message: t('profile.confirmPasswordRequired'), trigger: 'blur' }, { validator: validateConfirm, trigger: 'blur' }]
}))

const mfaRules = {
  code: [{ required: true, message: t('profile.mfaCodeRequired'), trigger: 'blur' }, { len: 6, message: t('profile.mfaCodeLen'), trigger: 'blur' }]
}

const strengthClass = computed(() => {
  const score = calcStrength(pwdForm.newPassword)
  if (score <= 2) return 'weak'
  if (score <= 3) return 'medium'
  return 'strong'
})

const strengthWidth = computed(() => {
  const score = calcStrength(pwdForm.newPassword)
  return (score / 5 * 100) + '%'
})

// 密码检查（用于显示各条件满足情况）
const pwdChecks = computed(() => checkPassword(pwdForm.newPassword))

// 密码强度分数
const strengthScore = computed(() => calcStrength(pwdForm.newPassword))

// 密码强度标签
const strengthLabel = computed(() => {
  const score = strengthScore.value
  if (score <= 1) return '弱'
  if (score <= 2) return '较弱'
  if (score <= 3) return '中等'
  if (score <= 4) return '强'
  return '非常强'
})

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  const ss = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
}

const loadProfile = async () => {
  loading.value = true
  try {
    const res = await AuthApi.getCurrentUser()
    if (res.code === 200) {
      const data = res.data
      currentUser.value = data
      Object.assign(form, {
        id: data.id,
        username: data.username,
        nickname: data.nickname,
        email: data.email,
        phone: data.phone,
        status: data.status,
        source: data.source,
        mfaEnabled: data.mfaEnabled,
        roleName: data.roleName
      })
    }
  } finally { loading.value = false }
}

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const res = await AuthApi.updateProfile({
      nickname: form.nickname,
      email: form.email,
      phone: form.phone
    })
    if (res.code === 200) {
      ElMessage.success(t('profile.saveSuccess'))
      // 更新本地存储的用户信息
      const stored = localStorage.getItem('user')
      if (stored) {
        const user = JSON.parse(stored)
        Object.assign(user, { nickname: form.nickname, email: form.email, phone: form.phone })
        localStorage.setItem('user', JSON.stringify(user))
      }
      loadProfile()
    } else {
      ElMessage.error(res.message || t('profile.saveFailed'))
    }
  } finally { submitting.value = false }
}

const handleChangePassword = async () => {
  const valid = await pwdFormRef.value?.validate().catch(() => false)
  if (!valid) return
  pwdLoading.value = true
  try {
    const res = await AuthApi.changePassword({
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword
    })
    if (res.code === 200) {
      ElMessage.success(t('profile.passwordChangeSuccess'))
      showChangePwd.value = false
      pwdForm.oldPassword = ''
      pwdForm.newPassword = ''
      pwdForm.confirmPassword = ''
    } else {
      ElMessage.error(res.message || t('profile.passwordChangeFailed'))
    }
  } finally { pwdLoading.value = false }
}

const handleEnableMFA = async () => {
  try {
    const res = await UserApi.generateMFASecret(form.id!)
    if (res.code === 200) {
      mfaSecret.value = res.data.secret
      qrCodeUrl.value = res.data.qrCode
      showMFADialog.value = true
    }
  } catch { ElMessage.error(t('profile.mfaInfoFailed')) }
}

const handleConfirmMFA = async () => {
  const valid = await mfaFormRef.value?.validate().catch(() => false)
  if (!valid) return
  mfaLoading.value = true
  try {
    const res = await UserApi.adminEnableMFA({ userId: form.id!, code: mfaForm.code })
    if (res.code === 200) {
      ElMessage.success(t('profile.mfaEnableSuccess'))
      showMFADialog.value = false
      loadProfile()
    } else {
      ElMessage.error(res.message || t('profile.mfaEnableFailed'))
    }
  } catch (e: any) { ElMessage.error(e.message || t('profile.mfaEnableFailed')) }
  finally { mfaLoading.value = false }
}

const handleDisableMFA = async () => {
  try {
    await UserApi.resetMFA({ userId: form.id! })
    ElMessage.success(t('profile.mfaResetSuccess'))
    loadProfile()
  } catch (e: any) { ElMessage.error(e.message || t('profile.mfaResetFailed')) }
}

onMounted(() => { loadProfile(); loadPolicy() })
</script>

<style scoped lang="scss">
.profile-page {
  padding: var(--space-4);
  min-height: 100vh;
  background: var(--color-page-bg);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-5);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.page-title {
  font-family: 'Manrope', sans-serif;
  font-size: 20px;
  font-weight: 800;
  color: var(--color-text-primary);
  letter-spacing: -0.3px;
}

.page-subtitle {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 2px;
  display: block;
}

.profile-layout {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: var(--space-3);
  align-items: start;
}

.profile-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 20px 16px;
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface);
}

.card-header-simple {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border-light);
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;

  svg { color: var(--color-primary); }
}

.card-header-avatar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-large {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Manrope', sans-serif;
  font-size: 18px;
  font-weight: 900;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.avatar-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.avatar-name {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 800;
  color: var(--color-text-primary);
}

.avatar-username {
  font-size: 12px;
  color: var(--color-text-muted);
}

.mfa-badge {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: var(--radius-full);

  &--on {
    background: rgba(34,197,94,0.1);
    color: var(--color-success);
  }
  &--off {
    background: var(--color-surface-3);
    color: var(--color-text-muted);
  }
}

.card-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-foot {
  padding: 12px 20px;
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
  display: flex;
  justify-content: flex-end;
}

.form-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.field-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  margin-top: 4px;
  display: inline-block;

  &--local { background: rgba(99,102,241,0.1); color: #6366f1; }
  &--ad { background: rgba(245,158,11,0.1); color: #f59e0b; }
}

/* 右侧信息行 */
.profile-sidebar {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-border-light);
  gap: 8px;

  &:last-child { border-bottom: none; }

  &--readonly { cursor: default; }
  &--disabled { opacity: 0.5; }
}

.info-row-left {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);

  svg { color: var(--color-text-muted); flex-shrink: 0; }
}

.info-row-value {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 5px;
}

.info-row-tip {
  font-size: 11px;
  color: var(--color-text-muted);
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
  &--on { background: var(--color-success); }
  &--off { background: var(--color-danger); }
}

.source-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: var(--radius-sm);

  &--ad { background: rgba(245,158,11,0.1); color: #f59e0b; }
  &--local { background: rgba(99,102,241,0.1); color: #6366f1; }
}

/* 密码强度 */
.password-strength {
  margin-top: 6px;
  .strength-bar { height: 3px; background: var(--color-surface-3); border-radius: 2px; overflow: hidden; }
  .strength-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.3s ease, background 0.3s ease;
    &.weak { background: var(--color-danger); }
    &.medium { background: var(--color-warning); }
    &.strong { background: var(--color-success); }
  }
}

/* 密码要求清单 */
.pwd-requirements {
  margin-top: 8px;
  padding: 14px 16px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.req-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}

.req-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.req-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  color: var(--color-text-muted);
  transition: color 0.2s;

  svg {
    stroke: var(--color-text-muted);
    flex-shrink: 0;
    opacity: 0.4;
    transition: all 0.2s;
  }

  &.satisfied {
    color: var(--color-success);
    svg {
      stroke: var(--color-success);
      opacity: 1;
    }
  }
}

.strength-section {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border-light);
}

.strength-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 11px;
  color: var(--color-text-secondary);
}

.strength-text {
  font-weight: 700;
  &.weak { color: var(--color-danger); }
  &.medium { color: var(--color-warning); }
  &.strong { color: var(--color-success); }
}

.strength-bar-segments {
  display: flex;
  gap: 4px;
}

.segment {
  flex: 1;
  height: 4px;
  background: var(--color-border-light);
  border-radius: 2px;
  transition: background 0.3s;

  &.weak.active { background: var(--color-danger); }
  &.medium.active { background: var(--color-warning); }
  &.strong.active { background: var(--color-success); }
}

/* 密码弹窗 */
.pwd-dialog-body {
  padding: 16px 20px;
}

/* MFA 弹窗 */
.mfa-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.mfa-tip {
  font-size: 13px;
  color: var(--color-text-muted);
  line-height: 1.5;
}

.qr-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.qr-img {
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: 8px;
  background: white;
}

.mfa-secret {
  font-size: 11px;
  color: var(--color-text-muted);
  font-family: monospace;
}

@media (max-width: 900px) {
  .profile-layout { grid-template-columns: 1fr; }
  .form-row-2 { grid-template-columns: 1fr; }
}

/* ==================== 对话框头部 ==================== */
:deep(.el-dialog__header) {
  padding: 16px 20px;
  margin-right: 0;
  border-bottom: 1px solid var(--color-border-light);
}

.dialog-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dialog-mode-tag {
  font-size: 10px;
  font-weight: 800;
  font-family: 'DM Sans', sans-serif;
  padding: 2px 8px;
  border-radius: 4px;
  letter-spacing: 0.5px;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border: 1px solid rgba(0, 94, 235, 0.2);

  &--warn {
    background: rgba(245, 158, 11, 0.1);
    color: #d97706;
    border-color: rgba(245, 158, 11, 0.2);
  }

  &--mfa {
    background: rgba(139, 92, 246, 0.1);
    color: #7c3aed;
    border-color: rgba(139, 92, 246, 0.2);
  }
}

.dialog-title-text {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}
</style>
