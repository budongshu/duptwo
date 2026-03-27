<template>
  <div class="login-page">
    <!-- Dot-matrix background texture -->
    <div class="bg-texture" aria-hidden="true"></div>

    <!-- Top bar -->
    <header class="topbar">
      <div class="topbar-brand">
        <div class="brand-mark">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <rect x="2" y="2" width="9" height="9" rx="2" fill="#005eeb" opacity="0.9"/>
            <rect x="13" y="2" width="9" height="9" rx="2" fill="#005eeb" opacity="0.5"/>
            <rect x="2" y="13" width="9" height="9" rx="2" fill="#005eeb" opacity="0.5"/>
            <rect x="13" y="13" width="9" height="9" rx="2" fill="#005eeb" opacity="0.25"/>
          </svg>
        </div>
        <div class="brand-text">
          <span class="brand-name">数据登记平台</span>
          <span class="brand-sub">DataRegistry</span>
        </div>
      </div>
    </header>

    <!-- Centered card -->
    <main class="login-center">
      <div class="login-card" :class="{ 'is-loading': loading }">

        <!-- Corner brackets -->
        <span class="bracket bracket--tl" aria-hidden="true"></span>
        <span class="bracket bracket--tr" aria-hidden="true"></span>
        <span class="bracket bracket--bl" aria-hidden="true"></span>
        <span class="bracket bracket--br" aria-hidden="true"></span>

        <transition name="step-fade" mode="out-in">

          <!-- ===== Login Step ===== -->
          <div v-if="!mfaStep" key="login" class="card-body">
            <div class="card-head">
              <h1 class="card-title">欢迎回来</h1>
              <p class="card-desc">请登录以继续访问数据登记平台</p>
            </div>

            <!-- Auth type tabs -->
            <div class="auth-tabs" v-if="adEnabled">
              <button :class="['auth-tab', { active: loginType === 'local' }]" @click="loginType = 'local'">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                本地账号
              </button>
              <button :class="['auth-tab', { active: loginType === 'ad' }]" @click="loginType = 'ad'">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                AD 域账号
              </button>
            </div>

            <!-- Local form -->
            <el-form v-if="loginType === 'local'" ref="formRef" :model="form" :rules="rules" class="auth-form" @submit.prevent="handleLogin">
              <div class="field-group">
                <label class="field-label">用户名</label>
                <el-input
                  v-model="form.username"
                  placeholder="输入用户名"
                  size="large"
                  :prefix-icon="User"
                  class="form-input"
                  @keyup.enter="handleLogin"
                />
              </div>
              <div class="field-group">
                <label class="field-label">密码</label>
                <el-input
                  v-model="form.password"
                  type="password"
                  placeholder="输入密码"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  class="form-input"
                  @keyup.enter="handleLogin"
                />
              </div>
              <div class="field-group" v-if="captchaEnabled">
                <label class="field-label">验证码</label>
                <div class="captcha-row">
                  <el-input v-model="form.captcha" placeholder="图形验证码" size="large" class="form-input capt-input" @keyup.enter="handleLogin" />
                  <div class="captcha-img-box" @click="refreshCaptcha" title="点击刷新">
                    <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-img" />
                    <span v-else class="captcha-loading">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                    </span>
                  </div>
                </div>
              </div>

              <button type="submit" class="submit-btn" :disabled="loading" @click="handleLogin">
                <span v-if="!loading" class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                  登录系统
                </span>
                <span v-else class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                  验证中...
                </span>
              </button>
            </el-form>

            <!-- AD form -->
            <el-form v-else ref="adFormRef" :model="adForm" :rules="adRules" class="auth-form" @submit.prevent="handleADLogin">
              <div class="ad-notice">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
                AD 域账号，密码为您的网络密码
              </div>
              <div class="field-group">
                <label class="field-label">域用户名</label>
                <el-input v-model="adForm.username" placeholder="如 zhangsan" size="large" :prefix-icon="User" class="form-input" @keyup.enter="handleADLogin" />
              </div>
              <div class="field-group">
                <label class="field-label">密码</label>
                <el-input v-model="adForm.password" type="password" placeholder="AD 网络密码" size="large" :prefix-icon="Lock" show-password class="form-input" @keyup.enter="handleADLogin" />
              </div>
              <div class="field-group" v-if="captchaEnabled">
                <label class="field-label">验证码</label>
                <div class="captcha-row">
                  <el-input v-model="adForm.captcha" placeholder="图形验证码" size="large" class="form-input capt-input" @keyup.enter="handleADLogin" />
                  <div class="captcha-img-box" @click="refreshCaptcha" title="点击刷新">
                    <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-img" />
                    <span v-else class="captcha-loading">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                    </span>
                  </div>
                </div>
              </div>

              <button type="submit" class="submit-btn" :disabled="loading" @click="handleADLogin">
                <span v-if="!loading" class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                  登录系统
                </span>
                <span v-else class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                  验证中...
                </span>
              </button>
            </el-form>

            <div class="card-footer">
              <router-link to="/register" class="footer-link">没有账户？立即注册</router-link>
            </div>
          </div>

          <!-- ===== MFA Step ===== -->
          <div v-else key="mfa" class="card-body">
            <div class="card-head">
              <div class="mfa-badge">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                双因素认证
              </div>
              <h1 class="card-title">验证身份</h1>
              <p class="card-desc">请打开身份验证器，输入六位动态码</p>
            </div>

            <el-form ref="mfaFormRef" :model="mfaForm" :rules="mfaRules" class="auth-form" @submit.prevent="handleMFAVerify">
              <div class="mfa-user-row">
                <div class="mfa-avatar">{{ form.username ? form.username.slice(0, 1).toUpperCase() : 'U' }}</div>
                <div class="mfa-user-info">
                  <span class="mfa-user-name">{{ form.username }}</span>
                  <span class="mfa-user-label">待验证身份</span>
                </div>
              </div>

              <div class="field-group">
                <label class="field-label">动态验证码</label>
                <el-input
                  v-model="mfaForm.code"
                  placeholder="6 位数字验证码"
                  size="large"
                  maxlength="6"
                  :prefix-icon="Key"
                  class="form-input"
                  @keyup.enter="handleMFAVerify"
                />
              </div>

              <button type="submit" class="submit-btn submit-btn--accent" :disabled="loading" @click="handleMFAVerify">
                <span v-if="!loading" class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  验证并进入
                </span>
                <span v-else class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                  验证中...
                </span>
              </button>

              <button type="button" class="back-btn" @click="mfaStep = false; mfaForm.code = ''">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                返回重新登录
              </button>
            </el-form>
          </div>

        </transition>
      </div>

      <!-- Footer -->
      <footer class="login-footer">
        <span>数据登记平台</span>
        <span class="footer-sep">·</span>
        <span>duptwo · {{ new Date().getFullYear() }}</span>
      </footer>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Key } from '@element-plus/icons-vue'
import { AuthApi } from '@/api/auth'

const router = useRouter()
const formRef = ref()
const adFormRef = ref()
const mfaFormRef = ref()
const loading = ref(false)
const mfaStep = ref(false)
const currentUserId = ref<number>(0)
const loginType = ref<'local' | 'ad'>('local')
const adEnabled = ref(false)
const captchaEnabled = ref(true)
const captchaId = ref('')
const captchaImage = ref('')

const form = reactive({ username: '', password: '', captcha: '' })
const adForm = reactive({ username: '', password: '', captcha: '' })
const mfaForm = reactive({ code: '' })

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const adRules = {
  username: [{ required: true, message: '请输入域用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const mfaRules = {
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

const refreshCaptcha = async () => {
  try {
    const res = await AuthApi.getCaptcha()
    if (res.code === 200) {
      captchaId.value = res.data.captchaId
      captchaImage.value = 'data:image/png;base64,' + res.data.captchaImage
      captchaEnabled.value = res.data.enabled
    }
  } catch {}
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    const res = await AuthApi.login({ username: form.username, password: form.password, captchaId: captchaId.value, captcha: form.captcha })
    if (res.code === 200) {
      if (res.data.mfaRequired) {
        currentUserId.value = res.data.user.id
        mfaStep.value = true
      } else {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        ElMessage.success('登录成功')
        router.push('/')
      }
    } else {
      ElMessage.error(res.message || '登录失败')
      if (captchaEnabled.value) refreshCaptcha()
    }
  } catch { ElMessage.error('登录失败'); if (captchaEnabled.value) refreshCaptcha() }
  finally { loading.value = false }
}

const handleADLogin = async () => {
  const valid = await adFormRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    const res = await AuthApi.adLogin({ username: adForm.username, password: adForm.password, captchaId: captchaId.value, captcha: adForm.captcha })
    if (res.code === 200) {
      if (res.data.mfaRequired) {
        currentUserId.value = res.data.user.id
        form.username = adForm.username
        mfaStep.value = true
      } else {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        ElMessage.success('登录成功')
        router.push('/')
      }
    } else { ElMessage.error(res.message || '登录失败'); if (captchaEnabled.value) refreshCaptcha() }
  } catch { ElMessage.error('登录失败'); if (captchaEnabled.value) refreshCaptcha() }
  finally { loading.value = false }
}

const handleMFAVerify = async () => {
  const valid = await mfaFormRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    const res = await AuthApi.mfaVerify({ userId: currentUserId.value, code: mfaForm.code })
    if (res.code === 200) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      ElMessage.success('验证通过，欢迎回来')
      router.push('/')
    } else {
      ElMessage.error(res.message || '验证码错误')
    }
  } catch { ElMessage.error('验证失败') }
  finally { loading.value = false }
}

onMounted(() => { refreshCaptcha() })
</script>

<style scoped lang="scss">
/* ==================== Page ==================== */
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f7;
  position: relative;
  overflow: hidden;
}

/* Subtle dot texture */
.bg-texture {
  position: fixed;
  inset: 0;
  background-image: radial-gradient(circle, rgba(0,0,0,0.05) 1px, transparent 1px);
  background-size: 24px 24px;
  pointer-events: none;
  z-index: 0;
}

/* ==================== Topbar ==================== */
.topbar {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  padding: 0 48px;
  height: 60px;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  flex-shrink: 0;
}

.topbar-brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1.5px solid rgba(0, 94, 235, 0.25);
  background: rgba(0, 94, 235, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.brand-name {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 800;
  color: #1f2329;
  letter-spacing: -0.3px;
  line-height: 1;
}

.brand-sub {
  font-family: 'DM Sans', sans-serif;
  font-size: 10px;
  color: rgba(0, 94, 235, 0.65);
  letter-spacing: 2px;
  text-transform: uppercase;
  font-weight: 600;
}

/* ==================== Center layout ==================== */
.login-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  position: relative;
  z-index: 1;
  gap: 32px;
}

/* ==================== Card ==================== */
.login-card {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 44px 40px;
  position: relative;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04);
  animation: cardIn 0.5s cubic-bezier(0.4, 0, 0.2, 1) both;

  &.is-loading {
    pointer-events: none;
    opacity: 0.85;
  }
}

@keyframes cardIn {
  from { opacity: 0; transform: translateY(24px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Corner brackets */
.bracket {
  position: absolute;
  width: 14px;
  height: 14px;
  border-color: rgba(0, 94, 235, 0.4);
  border-style: solid;

  &--tl { top: -1px; left: -1px; border-width: 2px 0 0 2px; border-radius: 16px 0 0 0; }
  &--tr { top: -1px; right: -1px; border-width: 2px 2px 0 0; border-radius: 0 16px 0 0; }
  &--bl { bottom: -1px; left: -1px; border-width: 0 0 2px 2px; border-radius: 0 0 0 16px; }
  &--br { bottom: -1px; right: -1px; border-width: 0 2px 2px 0; border-radius: 0 0 16px 0; }
}

/* ==================== Card content ==================== */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.card-head {
  margin-bottom: 32px;
}

.card-title {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 26px;
  font-weight: 800;
  color: #1f2329;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
  line-height: 1.2;
}

.card-desc {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  color: #5f6368;
  margin: 0;
  line-height: 1.5;
}

/* MFA badge */
.mfa-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: #92400e;
  letter-spacing: 1px;
  text-transform: uppercase;
  padding: 4px 10px;
  border: 1px solid rgba(217, 119, 6, 0.25);
  border-radius: 4px;
  background: rgba(217, 119, 6, 0.08);
  margin-bottom: 14px;

  svg { stroke: #d97706; }
}

/* ==================== Auth tabs ==================== */
.auth-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 28px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  background: #f8f9fb;
  width: fit-content;
}

.auth-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 22px;
  background: none;
  border: none;
  color: #9ca3af;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border-right: 1px solid #e5e7eb;

  &:last-child { border-right: none; }

  svg { opacity: 0.5; transition: opacity 0.2s; }

  &.active {
    background: #ffffff;
    color: #005eeb;
    box-shadow: inset 0 -2px 0 #005eeb;

    svg { opacity: 1; stroke: #005eeb; }
  }

  &:hover:not(.active) {
    color: #4b5563;
    svg { opacity: 0.7; }
  }
}

/* ==================== Form ==================== */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.field-group {
  margin-bottom: 18px;
}

.field-label {
  display: block;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  letter-spacing: 0.3px;
  margin-bottom: 8px;
}

.ad-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: rgba(217, 119, 6, 0.06);
  border: 1px solid rgba(217, 119, 6, 0.2);
  border-radius: 8px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #92400e;
  margin-bottom: 18px;

  svg { flex-shrink: 0; stroke: #d97706; }
}

/* Input styling */
:deep(.form-input) {
  .el-input__wrapper {
    background: #ffffff !important;
    border-radius: 8px !important;
    border: 1px solid #d1d5db !important;
    box-shadow: 0 0 0 1px #d1d5db inset !important;
    padding: 6px 14px !important;
    transition: border-color 0.2s ease, box-shadow 0.2s ease !important;
    width: 100%;

    &:hover {
      border-color: #9ca3af !important;
      box-shadow: 0 0 0 1px #9ca3af inset !important;
    }
    &.is-focus {
      border-color: #005eeb !important;
      background: rgba(0, 94, 235, 0.02) !important;
      box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.1) inset, 0 0 0 1px #005eeb inset !important;
    }
  }

  .el-input__inner {
    color: #1f2329 !important;
    font-size: 14px;
    font-family: 'DM Sans', sans-serif;
    font-weight: 500;
    &::placeholder { color: #9ca3af !important; }
  }

  .el-input__prefix .el-icon { color: #9ca3af; margin-right: 8px; }
}

/* Captcha */
.captcha-row {
  display: flex;
  gap: 10px;
  align-items: stretch;
}

.capt-input { flex: 1; }

.captcha-img-box {
  width: 128px;
  height: 42px;
  border-radius: 8px;
  background: #ffffff;
  border: 1px solid #d1d5db;
  cursor: pointer;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s, box-shadow 0.2s;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px #d1d5db inset;

  &:hover {
    border-color: #005eeb;
    box-shadow: 0 0 0 1px #005eeb inset;
  }
}

.captcha-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  image-rendering: crisp-edges;
}
.captcha-loading { color: #9ca3af; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ==================== Submit button ==================== */
.submit-btn {
  width: 100%;
  height: 48px;
  border-radius: 10px;
  background: #005eeb;
  border: none;
  color: #ffffff;
  font-family: 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.5px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 8px;
  transition: background 0.2s ease, transform 0.15s ease, box-shadow 0.25s ease;

  &:hover:not(:disabled) {
    background: #1a73e8;
    transform: translateY(-2px);
    box-shadow: 0 6px 28px rgba(0, 94, 235, 0.35);
  }

  &:active:not(:disabled) { transform: scale(0.985); }
  &:disabled { opacity: 0.5; cursor: not-allowed; }

  &--accent {
    background: #d97706;
    &:hover:not(:disabled) {
      background: #b45309;
      box-shadow: 0 6px 28px rgba(217, 119, 6, 0.35);
    }
  }
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* MFA user row */
.mfa-user-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(217, 119, 6, 0.06);
  border: 1px solid rgba(217, 119, 6, 0.2);
  border-radius: 10px;
  margin-bottom: 20px;
}

.mfa-avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(217, 119, 6, 0.1);
  border: 1px solid rgba(217, 119, 6, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Manrope', sans-serif;
  font-size: 15px;
  font-weight: 800;
  color: #d97706;
  flex-shrink: 0;
}

.mfa-user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mfa-user-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #1f2329;
}

.mfa-user-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  color: #92400e;
  opacity: 0.8;
}

/* Back button */
.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  margin-top: 14px;
  padding: 10px;
  background: none;
  border: none;
  color: #9ca3af;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;

  &:hover { color: #005eeb; background: rgba(0, 94, 235, 0.04); }
}

/* Card footer */
.card-footer {
  margin-top: 24px;
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #f0f1f3;
}

.footer-link {
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  color: #005eeb;
  text-decoration: none;
  transition: color 0.2s;

  &:hover { color: #1a73e8; }
}

/* ==================== Footer ==================== */
.login-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #9ca3af;
  letter-spacing: 0.5px;

  .footer-sep { opacity: 0.5; }
}

/* ==================== Transitions ==================== */
.step-fade-enter-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.step-fade-leave-active {
  transition: opacity 0.15s ease;
}
.step-fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.step-fade-leave-to {
  opacity: 0;
}

/* ==================== Responsive ==================== */
@media (max-width: 520px) {
  .topbar { padding: 0 20px; }
  .login-card { padding: 32px 24px; border-radius: 12px; }
  .card-title { font-size: 22px; }
  .auth-tab { padding: 9px 16px; font-size: 12px; }
}
</style>
