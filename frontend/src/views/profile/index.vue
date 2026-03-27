<template>
  <div class="profile-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">个人设置</h1>
        <span class="page-subtitle">管理个人信息与安全设置</span>
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
            {{ form.mfaEnabled ? 'MFA 已启用' : 'MFA 未启用' }}
          </span>
        </div>

        <div class="card-body">
          <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
            <div class="form-row-2">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="form.username" disabled>
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
                </el-input>
                <div class="field-tag field-tag--local" v-if="form.source === 'LOCAL'">本地账号</div>
                <div class="field-tag field-tag--ad" v-else>AD 域账号</div>
              </el-form-item>
              <el-form-item label="昵称" prop="nickname">
                <el-input v-model="form.nickname" placeholder="显示名称">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3"/></svg></template>
                </el-input>
              </el-form-item>
            </div>
            <div class="form-row-2">
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="form.email" :disabled="form.source === 'AD'" placeholder="邮箱地址">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg></template>
                </el-input>
              </el-form-item>
              <el-form-item label="手机号">
                <el-input v-model="form.phone" placeholder="手机号码">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg></template>
                </el-input>
              </el-form-item>
            </div>
          </el-form>
        </div>

        <div class="card-foot">
          <el-button type="primary" :loading="submitting" @click="handleSave">保存修改</el-button>
        </div>
      </div>

      <!-- 右侧：其他信息 -->
      <div class="profile-sidebar">

        <!-- 安全卡片 -->
        <div class="profile-card">
          <div class="card-header-simple">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            账号安全
          </div>
          <div class="card-body">
            <!-- 修改密码 -->
            <div class="info-row" v-if="form.source !== 'AD'">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <span>修改密码</span>
              </div>
              <el-button size="small" @click="showChangePwd = true">修改</el-button>
            </div>
            <div class="info-row info-row--disabled" v-else>
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <span>修改密码</span>
              </div>
              <span class="info-row-tip">AD域管理</span>
            </div>

            <!-- MFA 启用 -->
            <div class="info-row" v-if="form.source !== 'AD'">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                <span>MFA 认证</span>
              </div>
              <el-button v-if="!form.mfaEnabled" type="primary" size="small" plain @click="handleEnableMFA">启用 MFA</el-button>
              <el-button v-else type="warning" size="small" plain @click="handleDisableMFA">重置 MFA</el-button>
            </div>

            <!-- 角色信息 -->
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                <span>所属角色</span>
              </div>
              <span class="info-row-value">{{ currentUser?.roleName || '-' }}</span>
            </div>

            <!-- 最后登录 -->
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                <span>最后登录</span>
              </div>
              <span class="info-row-value">{{ currentUser?.lastLoginAt ? formatDate(currentUser.lastLoginAt) : '-' }}</span>
            </div>
          </div>
        </div>

        <!-- 账号信息 -->
        <div class="profile-card">
          <div class="card-header-simple">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
            账号信息
          </div>
          <div class="card-body">
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>账号来源</span>
              </div>
              <span class="info-row-value source-tag" :class="form.source === 'AD' ? 'source-tag--ad' : 'source-tag--local'">
                {{ form.source === 'AD' ? 'AD 域' : '本地' }}
              </span>
            </div>
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>账号状态</span>
              </div>
              <span class="info-row-value">
                <span class="status-dot" :class="form.status === 'active' ? 'status-dot--on' : 'status-dot--off'"></span>
                {{ form.status === 'active' ? '正常' : '禁用' }}
              </span>
            </div>
            <div class="info-row info-row--readonly">
              <div class="info-row-left">
                <span>注册时间</span>
              </div>
              <span class="info-row-value">{{ currentUser?.createdAt ? formatDate(currentUser.createdAt) : '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showChangePwd" width="420px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--warn">密码</span>
          <span class="dialog-title-text">修改密码</span>
        </div>
      </template>
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-position="top">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="pwdForm.oldPassword" type="password" placeholder="请输入当前密码" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="pwdForm.newPassword" type="password" placeholder="至少6位字符" show-password />
          <div class="password-strength" v-if="pwdForm.newPassword">
            <div class="strength-bar"><div class="strength-fill" :class="strengthClass" :style="{ width: strengthWidth }"></div></div>
          </div>
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="pwdForm.confirmPassword" type="password" placeholder="再次输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button size="small" @click="showChangePwd = false">取消</el-button>
        <el-button type="primary" size="small" :loading="pwdLoading" @click="handleChangePassword">确认修改</el-button>
      </template>
    </el-dialog>

    <!-- 启用 MFA 弹窗 -->
    <el-dialog v-model="showMFADialog" width="420px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--mfa">MFA</span>
          <span class="dialog-title-text">启用 MFA 认证</span>
        </div>
      </template>
      <div class="mfa-content">
        <p class="mfa-tip">使用身份验证器（如 Google Authenticator）扫描下方二维码，然后输入显示的 6 位验证码完成绑定。</p>
        <div class="qr-box" v-if="mfaSecret">
          <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="MFA QR" class="qr-img" />
          <div class="mfa-secret">密钥：{{ mfaSecret }}</div>
        </div>
        <el-form ref="mfaFormRef" :model="mfaForm" :rules="mfaRules" label-position="top">
          <el-form-item label="验证码" prop="code">
            <el-input v-model="mfaForm.code" placeholder="请输入6位验证码" maxlength="6" size="large" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button size="small" @click="showMFADialog = false">取消</el-button>
        <el-button type="primary" size="small" :loading="mfaLoading" @click="handleConfirmMFA">确认绑定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { UserApi, type User } from '@/api/user'
import { AuthApi } from '@/api/auth'

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
  if (value !== pwdForm.newPassword) callback(new Error('两次输入的密码不一致'))
  else callback()
}

const pwdRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认新密码', trigger: 'blur' }, { validator: validateConfirm, trigger: 'blur' }]
}

const mfaRules = {
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }, { len: 6, message: '验证码为6位数字', trigger: 'blur' }]
}

const strengthClass = computed(() => {
  const score = (pwdForm.newPassword.length >= 6 ? 1 : 0) + (/[A-Z]/.test(pwdForm.newPassword) ? 1 : 0) + (/[a-z]/.test(pwdForm.newPassword) ? 1 : 0) + (/[0-9]/.test(pwdForm.newPassword) ? 1 : 0) + (/[!@#$%^&*]/.test(pwdForm.newPassword) ? 1 : 0)
  if (score <= 2) return 'weak'
  if (score <= 3) return 'medium'
  return 'strong'
})

const strengthWidth = computed(() => {
  const score = (pwdForm.newPassword.length >= 6 ? 1 : 0) + (/[A-Z]/.test(pwdForm.newPassword) ? 1 : 0) + (/[a-z]/.test(pwdForm.newPassword) ? 1 : 0) + (/[0-9]/.test(pwdForm.newPassword) ? 1 : 0) + (/[!@#$%^&*]/.test(pwdForm.newPassword) ? 1 : 0)
  return (score / 5 * 100) + '%'
})

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}`
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
      ElMessage.success('保存成功')
      // 更新本地存储的用户信息
      const stored = localStorage.getItem('user')
      if (stored) {
        const user = JSON.parse(stored)
        Object.assign(user, { nickname: form.nickname, email: form.email, phone: form.phone })
        localStorage.setItem('user', JSON.stringify(user))
      }
      loadProfile()
    } else {
      ElMessage.error(res.message || '保存失败')
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
      ElMessage.success('密码修改成功')
      showChangePwd.value = false
      pwdForm.oldPassword = ''
      pwdForm.newPassword = ''
      pwdForm.confirmPassword = ''
    } else {
      ElMessage.error(res.message || '修改失败')
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
  } catch { ElMessage.error('获取MFA信息失败') }
}

const handleConfirmMFA = async () => {
  const valid = await mfaFormRef.value?.validate().catch(() => false)
  if (!valid) return
  mfaLoading.value = true
  try {
    const res = await UserApi.adminEnableMFA({ userId: form.id!, code: mfaForm.code })
    if (res.code === 200) {
      ElMessage.success('MFA 启用成功')
      showMFADialog.value = false
      loadProfile()
    } else {
      ElMessage.error(res.message || '启用失败')
    }
  } catch (e: any) { ElMessage.error(e.message || '启用失败') }
  finally { mfaLoading.value = false }
}

const handleDisableMFA = async () => {
  try {
    await UserApi.resetMFA({ userId: form.id! })
    ElMessage.success('MFA 已重置')
    loadProfile()
  } catch (e: any) { ElMessage.error(e.message || '重置失败') }
}

onMounted(() => { loadProfile() })
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
