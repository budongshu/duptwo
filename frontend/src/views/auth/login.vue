<template>
  <div class="login-page">
    <!-- Layered grid background -->
    <div class="bg-grid" aria-hidden="true"></div>
    <div class="bg-fade" aria-hidden="true"></div>

    <!-- Top bar -->
    <header class="topbar">
      <div class="topbar-brand">
        <div class="brand-mark">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <rect x="2" y="2" width="9" height="9" rx="2" fill="#3b82f6"/>
            <rect x="13" y="2" width="9" height="9" rx="2" fill="#93c5fd"/>
            <rect x="2" y="13" width="9" height="9" rx="2" fill="#93c5fd"/>
            <rect x="13" y="13" width="9" height="9" rx="2" fill="#dbeafe"/>
          </svg>
        </div>
        <div class="brand-text">
          <span class="brand-name">数据登记平台</span>
          <span class="brand-sub">duptwo</span>
        </div>
      </div>
      <div class="topbar-version">v2.0</div>
    </header>

    <!-- Centered card -->
    <main class="login-center">
      <!-- Ambient glow behind card -->
      <div class="card-glow" aria-hidden="true"></div>

      <div class="login-card" :class="{ 'is-loading': loading }">

        <!-- Heavy top accent bar -->
        <div class="card-accent-bar" aria-hidden="true"></div>

        <!-- Card header with logo -->
        <div class="card-header-block">
          <div class="card-logo-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <rect x="2" y="2" width="9" height="9" rx="2" fill="#3b82f6"/>
              <rect x="13" y="2" width="9" height="9" rx="2" fill="#3b82f6" opacity="0.55"/>
              <rect x="2" y="13" width="9" height="9" rx="2" fill="#3b82f6" opacity="0.55"/>
              <rect x="13" y="13" width="9" height="9" rx="2" fill="#3b82f6" opacity="0.3"/>
            </svg>
          </div>
          <div class="card-header-text">
            <h1 class="card-title">欢迎回来</h1>
            <p class="card-subtitle">登录到数据登记平台</p>
          </div>
        </div>

        <!-- Divider -->
        <div class="card-divider" aria-hidden="true"></div>

        <transition name="step-fade" mode="out-in">

          <!-- ===== Login Step ===== -->
          <div v-if="!mfaStep" key="login" class="card-body">

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

            <!-- Friendly Error Alert -->
            <transition name="error-fade">
              <div v-if="loginError" class="error-alert" :class="`error-alert--${loginError.type}`">
                <div class="error-icon">
                  <svg v-if="loginError.type === 'username'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                  <svg v-else-if="loginError.type === 'password'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                  <svg v-else-if="loginError.type === 'captcha'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M7 7h.01M7 12h10M7 17h6"/></svg>
                  <svg v-else-if="loginError.type === 'lockout'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
                  <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
                </div>
                <div class="error-content">
                  <div class="error-title">{{ loginError.title }}</div>
                  <div class="error-desc">{{ loginError.desc }}</div>
                  <div v-if="loginError.type === 'lockout'" class="error-countdown">
                    {{ formatCountdown(lockoutRemaining) }} 后自动解锁
                  </div>
                </div>
                <button class="error-close" @click="loginError = null">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
                </button>
              </div>
            </transition>

            <!-- Local form -->
            <el-form v-if="loginType === 'local'" ref="formRef" :model="form" :rules="rules" class="auth-form" @submit.prevent="handleLogin">
              <div class="field-group">
                <label class="field-label">
                  <span class="label-dot"></span>
                  用户名
                </label>
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
                <label class="field-label">
                  <span class="label-dot"></span>
                  密码
                </label>
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

              <button type="submit" class="submit-btn" :disabled="loading">
                <span v-if="!loading" class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                  登 录 系 统
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

              <button type="submit" class="submit-btn" :disabled="loading">
                <span v-if="!loading" class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                  登 录 系 统
                </span>
                <span v-else class="btn-content">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                  验证中...
                </span>
              </button>
            </el-form>

            <div class="card-footer-row">
              <router-link to="/register" class="footer-link">没有账户？立即注册</router-link>
            </div>
          </div>

          <!-- ===== MFA Step ===== -->
          <div v-else key="mfa" class="card-body">
            <div class="card-header-block">
              <div class="mfa-icon-wrap">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              </div>
              <div class="card-header-text">
                <h1 class="card-title">身份验证</h1>
                <p class="card-subtitle">双因素认证已启用，请输入动态码</p>
              </div>
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
                <label class="field-label">
                  <span class="label-dot label-dot--amber"></span>
                  动态验证码
                </label>
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

              <button type="submit" class="submit-btn submit-btn--amber" :disabled="loading" @click="handleMFAVerify">
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
const lockoutRemaining = ref(0)
let lockoutTimer: ReturnType<typeof setInterval> | null = null

interface LoginError {
  type: 'username' | 'password' | 'captcha' | 'lockout' | 'general'
  title: string
  desc: string
}

const loginError = ref<LoginError | null>(null)

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

const clearLockout = () => {
  lockoutRemaining.value = 0
  if (lockoutTimer) { clearInterval(lockoutTimer); lockoutTimer = null }
}

const showLockout = (minutes: number) => {
  loginError.value = {
    type: 'lockout',
    title: minutes > 0 ? '账号已被锁定' : '登录被拒绝',
    desc: minutes > 0 ? `登录失败次数过多，账号已临时锁定` : '请稍后再试'
  }
  lockoutRemaining.value = minutes * 60
  if (lockoutTimer) clearInterval(lockoutTimer)
  lockoutTimer = setInterval(() => {
    if (lockoutRemaining.value > 0) lockoutRemaining.value--
    else {
      clearInterval(lockoutTimer!)
      lockoutTimer = null
      loginError.value = null
    }
  }, 1000)
}

const formatCountdown = (seconds: number) => {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

// 分析错误类型并显示友好提示
const parseLoginError = (msg: string): LoginError => {
  if (msg.includes('锁定')) {
    const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '0')
    return {
      type: 'lockout',
      title: '账号已被锁定',
      desc: `登录失败次数过多，为保护账户安全已临时锁定，请 ${minutes} 分钟后再试`
    }
  }
  if (msg.includes('验证码已过期') || msg.includes('验证码错误')) {
    return {
      type: 'captcha',
      title: '验证码错误',
      desc: '验证码填写错误或已过期，请重新获取并填写'
    }
  }
  if (msg.includes('密码')) {
    return {
      type: 'password',
      title: '密码错误',
      desc: '您输入的密码不正确，请检查后重新输入'
    }
  }
  if (msg.includes('用户') || msg.includes('账号') || msg.includes('用户名')) {
    if (msg.includes('不存在') || msg.includes('未找到')) {
      return {
        type: 'username',
        title: '用户不存在',
        desc: '该用户名不存在，请检查或联系管理员'
      }
    }
    if (msg.includes('禁用')) {
      return {
        type: 'general',
        title: '账号已被禁用',
        desc: '您的账号已被禁用，请联系系统管理员'
      }
    }
  }
  return {
    type: 'general',
    title: '登录失败',
    desc: msg || '登录失败，请稍后重试'
  }
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  loginError.value = null
  clearLockout()
  try {
    const res = await AuthApi.login({ username: form.username, password: form.password, captchaId: captchaId.value, captcha: form.captcha })
    if (res.code === 200) {
      if (res.data.mfaRequired) {
        currentUserId.value = res.data.user.id
        mfaStep.value = true
      } else {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        router.push('/')
      }
    } else {
      const msg = res.message || '登录失败'
      if (msg.includes('锁定')) {
        const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
        showLockout(minutes)
      } else {
        loginError.value = parseLoginError(msg)
      }
      if (captchaEnabled.value) refreshCaptcha()
    }
  } catch (e: any) {
    const msg = e?.message || '登录失败'
    if (msg.includes('锁定')) {
      const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
      showLockout(minutes)
    } else {
      loginError.value = parseLoginError(msg)
    }
    if (captchaEnabled.value) refreshCaptcha()
  }
  finally { loading.value = false }
}

const handleADLogin = async () => {
  const valid = await adFormRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  loginError.value = null
  clearLockout()
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
        router.push('/')
      }
    } else {
      const msg = res.message || '登录失败'
      if (msg.includes('锁定')) {
        const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
        showLockout(minutes)
      } else {
        loginError.value = parseLoginError(msg)
      }
      if (captchaEnabled.value) refreshCaptcha()
    }
  } catch (e: any) {
    const msg = e?.message || '登录失败'
    if (msg.includes('锁定')) {
      const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
      showLockout(minutes)
    } else {
      loginError.value = parseLoginError(msg)
    }
    if (captchaEnabled.value) refreshCaptcha()
  }
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
      router.push('/')
    } else {
      loginError.value = {
        type: 'captcha',
        title: '验证码错误',
        desc: res.message || '动态验证码不正确，请重新输入'
      }
    }
  } catch {
    loginError.value = {
      type: 'general',
      title: '验证失败',
      desc: '验证过程出错，请稍后重试'
    }
  }
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
  background-color: #f1f5f9;
  position: relative;
  overflow: hidden;
}

/* Subtle dot texture on white-ish bg */
.bg-grid {
  position: fixed;
  inset: 0;
  background-image: radial-gradient(circle, rgba(15, 23, 42, 0.04) 1px, transparent 1px);
  background-size: 28px 28px;
  pointer-events: none;
  z-index: 0;
}

.bg-fade {
  position: fixed;
  inset: 0;
  background: radial-gradient(ellipse 70% 60% at 50% 40%, transparent 50%, rgba(15, 23, 42, 0.05) 100%);
  pointer-events: none;
  z-index: 0;
}

/* ==================== Topbar ==================== */
.topbar {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 56px;
  height: 52px;
  border-bottom: 1px solid #e2e8f0;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  flex-shrink: 0;
  box-shadow: 0 1px 0 rgba(15, 23, 42, 0.04);
}

.topbar-brand {
  display: flex;
  align-items: center;
  gap: 13px;
}

.brand-mark {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  border: 1.5px solid rgba(37, 99, 235, 0.2);
  background: rgba(37, 99, 235, 0.05);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.08);
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-name {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.3px;
  line-height: 1;
}

.brand-sub {
  font-family: 'DM Sans', 'Courier New', monospace;
  font-size: 10px;
  color: rgba(37, 99, 235, 0.6);
  letter-spacing: 2.5px;
  font-weight: 600;
}

.topbar-version {
  font-family: 'Courier New', monospace;
  font-size: 11px;
  color: #cbd5e1;
  letter-spacing: 1px;
}

/* ==================== Center layout ==================== */
.login-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 28px 20px;
  position: relative;
  z-index: 1;
  gap: 28px;
}

/* Ambient shadow under card */
.card-glow {
  position: absolute;
  width: 580px;
  height: 280px;
  background: radial-gradient(ellipse, rgba(15, 23, 42, 0.07) 0%, transparent 70%);
  border-radius: 50%;
  pointer-events: none;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

/* ==================== Card ==================== */
.login-card {
  width: 100%;
  max-width: 480px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 0;
  position: relative;
  box-shadow:
    0 0 0 1px rgba(37, 99, 235, 0.04),
    0 1px 2px rgba(15, 23, 42, 0.04),
    0 4px 8px -2px rgba(15, 23, 42, 0.06),
    0 12px 32px -8px rgba(15, 23, 42, 0.12);
  overflow: hidden;
  animation: cardIn 0.55s cubic-bezier(0.4, 0, 0.2, 1) both;

  &.is-loading {
    pointer-events: none;
    opacity: 0.92;
  }
}

@keyframes cardIn {
  from { opacity: 0; transform: translateY(22px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

/* Thick top accent bar */
.card-accent-bar {
  height: 3px;
  background: linear-gradient(90deg, #1e3a8a 0%, #1d4ed8 30%, #3b82f6 60%, #1d4ed8 100%);
  border-radius: 16px 16px 0 0;
}

/* Card header block */
.card-header-block {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 28px 0;
}

.card-logo-icon {
  width: 42px;
  height: 42px;
  border-radius: 11px;
  background: linear-gradient(145deg, #eff6ff 0%, #dbeafe 100%);
  border: 1px solid #bfdbfe;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.12);
}

.card-header-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.card-title {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 21px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.4px;
  line-height: 1.2;
}

.card-subtitle {
  font-family: 'DM Sans', sans-serif;
  font-size: 12.5px;
  color: #64748b;
  margin: 0;
  line-height: 1.5;
}

/* Heavy divider */
.card-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, #e2e8f0 15%, #e2e8f0 85%, transparent);
  margin: 18px 28px;
}

/* ==================== Card content ==================== */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0 28px 24px;
}

/* ==================== Auth tabs ==================== */
.auth-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 20px;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
  background: #f8fafc;
  width: fit-content;
}

.auth-tab {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 22px;
  background: none;
  border: none;
  color: #94a3b8;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border-right: 1.5px solid #e2e8f0;

  &:last-child { border-right: none; }
  svg { opacity: 0.45; transition: opacity 0.2s; }

  &.active {
    background: #ffffff;
    color: #1d4ed8;
    box-shadow: inset 0 -2.5px 0 #1d4ed8;
    svg { opacity: 1; stroke: #1d4ed8; }
  }

  &:hover:not(.active) {
    color: #475569;
    svg { opacity: 0.7; }
  }
}

/* ==================== Friendly Error Alert ==================== */
.error-alert {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 14px 16px;
  border-radius: 12px;
  margin-bottom: 18px;
  animation: errorIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) both;
  position: relative;

  &--username {
    background: linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%);
    border: 1.5px solid #fed7aa;
    .error-icon { background: rgba(249, 115, 22, 0.12); svg { stroke: #ea580c; } }
    .error-title { color: #9a3412; }
    .error-desc { color: #c2410c; }
  }

  &--password {
    background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
    border: 1.5px solid #fecaca;
    .error-icon { background: rgba(220, 38, 38, 0.12); svg { stroke: #dc2626; } }
    .error-title { color: #991b1b; }
    .error-desc { color: #b91c1c; }
  }

  &--captcha {
    background: linear-gradient(135deg, #fefce8 0%, #fef9c3 100%);
    border: 1.5px solid #fef08a;
    .error-icon { background: rgba(234, 179, 8, 0.12); svg { stroke: #ca8a04; } }
    .error-title { color: #854d0e; }
    .error-desc { color: #a16207; }
  }

  &--lockout {
    background: linear-gradient(135deg, #fff1f2 0%, #ffe4e6 100%);
    border: 1.5px solid #fecdd3;
    .error-icon { background: rgba(220, 38, 38, 0.12); svg { stroke: #dc2626; } }
    .error-title { color: #991b1b; }
    .error-desc { color: #b91c1c; }
    .error-countdown {
      color: #dc2626;
      font-weight: 700;
      font-size: 16px;
      margin-top: 8px;
      font-family: 'SF Mono', monospace;
    }
  }

  &--general {
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    border: 1.5px solid #e2e8f0;
    .error-icon { background: rgba(100, 116, 139, 0.12); svg { stroke: #64748b; } }
    .error-title { color: #475569; }
    .error-desc { color: #64748b; }
  }
}

@keyframes errorIn {
  from { opacity: 0; transform: scale(0.95) translateY(-6px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.error-fade-enter-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.error-fade-leave-active {
  transition: opacity 0.15s ease;
}
.error-fade-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(-6px);
}
.error-fade-leave-to {
  opacity: 0;
}

.error-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.error-content {
  flex: 1;
  min-width: 0;
}

.error-title {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 700;
  margin-bottom: 4px;
}

.error-desc {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  line-height: 1.5;
  opacity: 0.9;
}

.error-close {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.5;
  transition: all 0.2s;
  flex-shrink: 0;

  svg { stroke: #64748b; }

  &:hover {
    opacity: 1;
    background: rgba(0, 0, 0, 0.05);
  }
}

/* ==================== Form ==================== */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.field-group {
  margin-bottom: 14px;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 7px;
  font-family: 'DM Sans', 'PingFang SC', sans-serif;
  font-size: 11.5px;
  font-weight: 700;
  color: #334155;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  margin-bottom: 8px;
}

.label-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #1d4ed8;
  flex-shrink: 0;
  box-shadow: 0 0 6px rgba(29, 78, 216, 0.4);
}

.label-dot--amber {
  background: #d97706;
  box-shadow: 0 0 6px rgba(217, 119, 6, 0.4);
}

.ad-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 13px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 9px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #92400e;
  margin-bottom: 16px;

  svg { flex-shrink: 0; stroke: #d97706; }
}

/* Input styling */
:deep(.form-input) {
  .el-input__wrapper {
    background: #ffffff !important;
    border-radius: 10px !important;
    border: 1.5px solid #cbd5e1 !important;
    box-shadow: 0 0 0 1px #cbd5e1 inset !important;
    padding: 5px 14px !important;
    transition: border-color 0.2s ease, box-shadow 0.2s ease !important;
    width: 100%;

    &:hover {
      border-color: #94a3b8 !important;
      box-shadow: 0 0 0 1px #94a3b8 inset !important;
    }
    &.is-focus {
      border-color: #1d4ed8 !important;
      background: #f8faff !important;
      box-shadow: 0 0 0 3px rgba(29, 78, 216, 0.08) inset, 0 0 0 1px #1d4ed8 inset !important;
    }
  }

  .el-input__inner {
    color: #0f172a !important;
    font-size: 14px;
    font-family: 'DM Sans', 'PingFang SC', sans-serif;
    font-weight: 500;
    &::placeholder { color: #94a3b8 !important; }
  }

  .el-input__prefix .el-icon { color: #94a3b8; margin-right: 10px; }
}

/* Captcha */
.captcha-row {
  display: flex;
  gap: 12px;
  align-items: stretch;
}

.capt-input { flex: 1; }

.captcha-img-box {
  width: 126px;
  height: 42px;
  border-radius: 10px;
  background: #ffffff;
  border: 1.5px solid #cbd5e1;
  cursor: pointer;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s, box-shadow 0.2s;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px #cbd5e1 inset;

  &:hover {
    border-color: #1d4ed8;
    box-shadow: 0 0 0 1px #1d4ed8 inset, 0 0 10px rgba(29, 78, 216, 0.1);
  }
}

.captcha-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  image-rendering: crisp-edges;
}
.captcha-loading { color: #94a3b8; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ==================== Submit button ==================== */
.submit-btn {
  width: 100%;
  height: 46px;
  border-radius: 12px;
  background: linear-gradient(135deg, #1e3a8a 0%, #1d4ed8 45%, #2563eb 100%);
  border: none;
  color: #ffffff;
  font-family: 'DM Sans', 'PingFang SC', sans-serif;
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 2.5px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 6px;
  transition: all 0.25s ease;
  box-shadow:
    0 1px 2px rgba(29, 78, 216, 0.3),
    0 4px 12px rgba(29, 78, 216, 0.25),
    inset 0 1px 0 rgba(255,255,255,0.1);
  position: relative;
  overflow: hidden;

  &::after {
    content: '';
    position: absolute;
    top: 0; left: -100%;
    width: 60%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
    transition: left 0.5s ease;
  }

  &:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow:
      0 4px 6px rgba(29, 78, 216, 0.3),
      0 10px 24px rgba(29, 78, 216, 0.3),
      inset 0 1px 0 rgba(255,255,255,0.15);
    &::after { left: 150%; }
  }

  &:active:not(:disabled) { transform: scale(0.985); }
  &:disabled { opacity: 0.5; cursor: not-allowed; }

  &--amber {
    background: linear-gradient(135deg, #b45309 0%, #d97706 50%, #f59e0b 100%);
    box-shadow:
      0 1px 2px rgba(217, 119, 6, 0.3),
      0 4px 12px rgba(217, 119, 6, 0.25),
      inset 0 1px 0 rgba(255,255,255,0.1);

    &:hover:not(:disabled) {
      box-shadow:
        0 4px 6px rgba(217, 119, 6, 0.3),
        0 10px 24px rgba(217, 119, 6, 0.3),
        inset 0 1px 0 rgba(255,255,255,0.15);
    }
  }
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 10px;
  position: relative;
  z-index: 1;
}

/* MFA user row */
.mfa-user-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 13px 17px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
  margin-bottom: 18px;
}

.mfa-icon-wrap {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  svg { stroke: #d97706; }
}

.mfa-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.25);
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
  gap: 3px;
}

.mfa-user-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: #0f172a;
}

.mfa-user-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  color: #92400e;
  opacity: 0.85;
}

/* Back button */
.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  margin-top: 12px;
  padding: 10px;
  background: none;
  border: 1.5px solid #e2e8f0;
  color: #64748b;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 10px;
  transition: all 0.2s;

  &:hover {
    color: #1d4ed8;
    border-color: #93c5fd;
    background: #eff6ff;
  }
}

/* Card footer */
.card-footer-row {
  margin-top: 18px;
  padding-top: 16px;
  border-top: 1px solid #f1f5f9;
  text-align: center;
}

.footer-link {
  font-family: 'DM Sans', 'PingFang SC', sans-serif;
  font-size: 13px;
  color: #1d4ed8;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;

  &:hover { color: #1e3a8a; text-decoration: underline; }
}

/* ==================== Footer ==================== */
.login-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'DM Sans', 'PingFang SC', sans-serif;
  font-size: 12px;
  color: #94a3b8;
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
  transform: translateY(10px);
}
.step-fade-leave-to {
  opacity: 0;
}

/* ==================== Responsive ==================== */
@media (max-width: 520px) {
  .topbar { padding: 0 20px; }
  .login-card { max-width: 360px; border-radius: 16px; }
  .card-header-block, .card-body { padding-left: 24px; padding-right: 24px; }
  .card-divider { margin-left: 24px; margin-right: 24px; }
  .card-title { font-size: 19px; }
  .auth-tab { padding: 9px 16px; font-size: 12px; }
}
</style>
