<template>
  <div class="login-page">
    <!-- 动态背景 -->
    <div class="bg-animate">
      <div class="blob blob--1"></div>
      <div class="blob blob--2"></div>
      <div class="blob blob--3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-wrapper">
      <div class="login-card">

        <!-- Logo + 标题 -->
        <div class="card-head">
          <div class="logo-wrap">
            <div class="logo-bounce">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                <rect x="2" y="2" width="9" height="9" rx="2" fill="#818cf8"/>
                <rect x="13" y="2" width="9" height="9" rx="2" fill="#a78bfa"/>
                <rect x="2" y="13" width="9" height="9" rx="2" fill="#c4b5fd"/>
                <rect x="13" y="13" width="9" height="9" rx="2" fill="#e0e7ff"/>
              </svg>
            </div>
            <div class="logo-ring"></div>
          </div>
          <h1 class="title">数据登记平台</h1>
          <p class="subtitle">duptwo</p>
        </div>

        <!-- 表单区域 -->
        <div class="card-body">
          <div class="greet" :class="{ 'greet--error': loginError }">
            <span v-if="!loginError">欢迎登录</span>
            <span v-else>{{ loginError.title }}</span>
          </div>

          <!-- 错误提示 -->
          <transition name="err-pop">
            <div v-if="loginError" class="error-box" :class="`error-box--${loginError.type}`">
              <div class="error-icon">
                <svg v-if="loginError.type === 'password' || loginError.type === 'username'" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
                <svg v-else-if="loginError.type === 'lockout'" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/><line x1="12" y1="16" x2="12" y2="18"/></svg>
                <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
              </div>
              <div class="error-text">
                <div class="error-main">{{ loginError.desc }}</div>
                <div v-if="loginError.remainingAttempts !== undefined && loginError.remainingAttempts > 0 && loginError.type !== 'lockout'" class="error-remaining">
                  剩余 {{ loginError.remainingAttempts }} / {{ loginError.maxAttempts }} 次
                </div>
                <div v-if="loginError.type === 'lockout'" class="error-countdown">
                  {{ formatCountdown(lockoutRemaining) }} 后自动解锁
                </div>
              </div>
              <button class="error-close" @click="loginError = null">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
          </transition>

          <!-- 登录表单 -->
          <el-form ref="formRef" :model="form" :rules="rules" class="login-form" @submit.prevent="handleLogin">
            <div class="field">
              <label>用户名</label>
              <el-input
                v-model="form.username"
                placeholder="输入用户名"
                size="large"
                :prefix-icon="User"
                @keyup.enter="handleLogin"
              />
            </div>
            <div class="field">
              <label>密码</label>
              <el-input
                v-model="form.password"
                type="password"
                placeholder="输入密码"
                size="large"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </div>
            <div class="field" v-if="captchaEnabled">
              <label>验证码</label>
              <div class="captcha-row">
                <el-input
                  v-model="form.captcha"
                  placeholder="输入验证码"
                  size="large"
                  @keyup.enter="handleLogin"
                />
                <div class="captcha-img" @click="refreshCaptcha">
                  <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
                  <span v-else class="captcha-loading">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                  </span>
                </div>
              </div>
            </div>

            <button type="submit" class="login-btn" :disabled="loading" :class="{ 'login-btn--loading': loading }">
              <span v-if="loading" class="btn-loading">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
                验证中...
              </span>
              <span v-else>
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" style="vertical-align: -2px; margin-right: 5px"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                登 录
              </span>
              <span class="btn-shine"></span>
            </button>
          </el-form>

          <div class="card-footer">
            <router-link v-if="registrationEnabled" to="/register">没有账号？立即注册</router-link>
          </div>
        </div>
      </div>

      <div class="login-footer">
        <span>duptwo · {{ new Date().getFullYear() }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { AuthApi } from '@/api/auth'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const mfaStep = ref(false)
const currentUserId = ref<number>(0)
const loginType = ref<'local' | 'ad'>('local')
const adEnabled = ref(false)
const captchaEnabled = ref(true)
const captchaId = ref('')
const captchaImage = ref('')
const lockoutRemaining = ref(0)
const registrationEnabled = ref(true)
let lockoutTimer: ReturnType<typeof setInterval> | null = null

interface LoginError {
  type: 'username' | 'password' | 'captcha' | 'lockout' | 'general'
  title: string
  desc: string
  remainingAttempts?: number
  maxAttempts?: number
}

const loginError = ref<LoginError | null>(null)

const form = reactive({ username: '', password: '', captcha: '' })

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
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
  if (lockoutTimer) { clearInterval(lockoutTimer); lockoutTimer = null }
  lockoutRemaining.value = minutes * 60
  const lockoutEndTime = Date.now() + minutes * 60 * 1000

  loginError.value = {
    type: 'lockout',
    title: '账号已被锁定',
    desc: `由于连续登录失败，为了保护您的账户安全，系统已临时锁定。请等待 ${minutes} 分钟后自动解锁，或联系管理员协助。`,
    lockoutEndTime: lockoutEndTime
  }

  lockoutTimer = setInterval(() => {
    const remaining = Math.max(0, Math.ceil((lockoutEndTime - Date.now()) / 1000))
    lockoutRemaining.value = remaining
    if (remaining <= 0) {
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

const parseLoginError = (msg: string): LoginError => {
  if (msg.includes('锁定')) {
    const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
    return {
      type: 'lockout', title: '账号已被锁定',
      desc: `由于连续登录失败，为了保护您的账户安全，系统已临时锁定。请等待 ${minutes} 分钟后自动解锁，或联系管理员协助。`
    }
  }
  if (msg.includes('验证码')) {
    return { type: 'captcha', title: '验证码错误', desc: '验证码填写错误或已过期，请点击图片刷新后重新输入' }
  }
  if (msg.includes('请填写验证码')) {
    return { type: 'captcha', title: '请填写验证码', desc: '请先获取并填写图形验证码后再点击登录' }
  }
  if (msg.includes('用户不存在') || msg.includes('不存在')) {
    return { type: 'username', title: '用户不存在', desc: '该用户名不存在，请检查输入是否有误' }
  }
  if (msg.includes('密码错误') || msg.includes('密码不正确')) {
    return { type: 'password', title: '密码输入错误', desc: '您输入的密码不正确，请检查大小写是否正确' }
  }
  if (msg.includes('用户名或密码错误') || msg.includes('不正确')) {
    return { type: 'password', title: '密码输入错误', desc: '用户名或密码不正确，请仔细核对后重新输入' }
  }
  if (msg.includes('禁用') || msg.includes('停用')) {
    return { type: 'general', title: '账号已被停用', desc: '您的账号已被管理员停用，如需恢复使用，请联系系统管理员' }
  }
  if (msg.includes('未激活') || msg.includes('待审核')) {
    return { type: 'general', title: '账号待审核', desc: '您的账号正在等待管理员审核，审核通过后即可正常登录' }
  }
  return { type: 'general', title: '登录失败', desc: msg ? `提示：${msg}` : '登录遇到问题，请稍后重试' }
}

const handleLogin = async () => {
  if (!form.username.trim()) {
    loginError.value = { type: 'username', title: '请输入用户名', desc: '用户名不能为空' }
    return
  }
  if (!form.password) {
    loginError.value = { type: 'password', title: '请输入密码', desc: '密码不能为空' }
    return
  }
  if (captchaEnabled.value && !form.captcha.trim()) {
    loginError.value = { type: 'captcha', title: '请输入验证码', desc: '请输入图片中的验证码' }
    return
  }

  loginError.value = null
  loading.value = true

  try {
    const res = await AuthApi.login({
      username: form.username,
      password: form.password,
      captchaId: captchaId.value,
      captcha: form.captcha
    })

    if (res.data.mfaRequired) {
      currentUserId.value = res.data.user.id
      mfaStep.value = true
    } else {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      router.push('/')
    }
  } catch (err: any) {
    let msg = ''
    let remainingAttempts = 0
    let maxAttempts = 5

    if (err && err.message) {
      msg = err.message
      if (err.data?.remainingAttempts !== undefined) {
        remainingAttempts = err.data.remainingAttempts
        maxAttempts = err.data.maxAttempts || 5
      }
    } else if (err?.response?.data) {
      const respData = err.response.data
      msg = respData.message || '登录失败'
      if (respData.data?.remainingAttempts !== undefined) {
        remainingAttempts = respData.data.remainingAttempts
        maxAttempts = respData.data.maxAttempts || 5
      }
    } else {
      msg = '登录失败，请稍后重试'
    }

    const errorInfo = parseLoginError(msg)
    loginError.value = { ...errorInfo, remainingAttempts, maxAttempts }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  refreshCaptcha()
  try {
    const res = await AuthApi.getRegistrationStatus()
    if (res.code === 200 && res.data) {
      registrationEnabled.value = res.data.registrationEnabled
    }
  } catch {}
})
</script>

<style scoped lang="scss">
/* ==================== 背景动画 ==================== */
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f1ff;
  position: relative;
  overflow: hidden;
}

.bg-animate {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(72px);
  opacity: 0.55;
  animation: blob-float 8s ease-in-out infinite;

  &--1 {
    width: 380px;
    height: 380px;
    background: #c7d2fe;
    top: -120px;
    right: -80px;
    animation-duration: 9s;
  }

  &--2 {
    width: 280px;
    height: 280px;
    background: #ddd6fe;
    bottom: -80px;
    left: -60px;
    animation-duration: 11s;
    animation-delay: -3s;
  }

  &--3 {
    width: 200px;
    height: 200px;
    background: #fae8ff;
    top: 40%;
    left: 20%;
    animation-duration: 13s;
    animation-delay: -6s;
  }
}

@keyframes blob-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(20px, -25px) scale(1.05); }
  66% { transform: translate(-15px, 15px) scale(0.95); }
}

/* ==================== 卡片容器 ==================== */
.login-wrapper {
  width: 100%;
  max-width: 480px;
  padding: 16px;
  position: relative;
  z-index: 1;
}

.login-card {
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 8px 40px rgba(99, 102, 241, 0.10), 0 2px 8px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  animation: card-in 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes card-in {
  from { opacity: 0; transform: translateY(24px) scale(0.96); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

/* ==================== 头部 Logo ==================== */
.card-head {
  padding: 18px 56px 16px;
  text-align: center;
  border-bottom: 1px solid #f1f0ff;
  background: linear-gradient(180deg, #fafafa, #fff);
}

.logo-wrap {
  position: relative;
  width: 46px;
  height: 46px;
  margin: 0 auto 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-bounce {
  animation: logo-wobble 3s ease-in-out infinite;
  position: relative;
  z-index: 1;
}

.logo-ring {
  position: absolute;
  inset: 0;
  border-radius: 16px;
  border: 2px solid #c7d2fe;
  animation: ring-pulse 3s ease-in-out infinite;
}

@keyframes logo-wobble {
  0%, 100% { transform: rotate(-3deg) scale(1); }
  50% { transform: rotate(3deg) scale(1.03); }
}

@keyframes ring-pulse {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.08); opacity: 0.3; }
}

.title {
  font-size: 16px;
  font-weight: 700;
  color: #3730a3;
  margin: 0 0 2px;
  letter-spacing: -0.2px;
}

.subtitle {
  font-size: 10px;
  color: #c7d2fe;
  margin: 0;
  font-family: 'SF Mono', 'Monaco', monospace;
  font-weight: 600;
  letter-spacing: 1.5px;
}

/* ==================== 表单区域 ==================== */
.card-body {
  padding: 16px 56px 24px;
}

.greet {
  font-size: 16px;
  font-weight: 700;
  color: #3730a3;
  margin-bottom: 18px;
  letter-spacing: -0.2px;

  &--error { color: #dc2626; }
}

/* 错误提示 */
.error-box {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 12px;
  margin-bottom: 16px;
  position: relative;

  &--password,
  &--username {
    background: #fef2f2;
    border: 1px solid #fecaca;
    .error-icon { color: #ef4444; }
    .error-main { color: #991b1b; }
  }

  &--captcha {
    background: #fefce8;
    border: 1px solid #fde68a;
    .error-icon { color: #ca8a04; }
    .error-main { color: #78350f; }
  }

  &--lockout {
    background: #fff1f2;
    border: 1px solid #fecdd3;
    .error-icon { color: #e11d48; }
    .error-main { color: #881337; }
  }

  &--general {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    .error-icon { color: #64748b; }
    .error-main { color: #334155; }
  }
}

.error-icon { flex-shrink: 0; margin-top: 1px; }
.error-text { flex: 1; min-width: 0; }
.error-main { font-size: 12.5px; line-height: 1.5; }
.error-countdown {
  font-size: 13px;
  font-weight: 700;
  font-family: 'SF Mono', monospace;
  color: #e11d48;
  margin-top: 5px;
}
.error-remaining { font-size: 11px; color: #94a3b8; margin-top: 4px; font-weight: 500; }
.error-close {
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  transition: all 0.15s;
  &:hover { background: rgba(0, 0, 0, 0.06); color: #64748b; }
}

/* 错误动画 */
.err-pop-enter-active { animation: err-pop-in 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
.err-pop-leave-active { animation: err-pop-in 0.2s ease reverse; }
@keyframes err-pop-in {
  from { opacity: 0; transform: translateY(-6px) scale(0.96); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

/* 表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field {
  label {
    display: block;
    font-size: 12.5px;
    font-weight: 600;
    color: #6366f1;
    margin-bottom: 6px;
    letter-spacing: 0.2px;
  }
}

/* 输入框 */
:deep(.el-input__wrapper) {
  border-radius: 10px !important;
  box-shadow: 0 0 0 1px #e0e7ff inset !important;
  padding: 3px 10px !important;
  background: #f8f9ff !important;
  transition: all 0.15s ease !important;

  &:hover {
    background: #fff !important;
    box-shadow: 0 0 0 1px #c7d2fe inset !important;
  }

  &.is-focus {
    background: #fff !important;
    box-shadow: 0 0 0 1.5px #818cf8 inset !important;
  }
}

:deep(.el-input__inner) {
  font-size: 13px !important;
  color: #3730a3 !important;
  &::placeholder { color: #c7d2fe !important; }
}

:deep(.el-input__prefix .el-input__prefix-inner .el-icon) {
  color: #c7d2fe;
  transition: color 0.18s;
}

:deep(.el-input__wrapper.is-focus .el-input__prefix .el-input__prefix-inner .el-icon) {
  color: #6366f1;
}

/* 验证码 */
.captcha-row {
  display: flex;
  gap: 8px;
  .el-input { flex: 1; }
}
.captcha-img {
  width: 94px;
  height: 38px;
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  background: #f1f5f9;
  border: 1px dashed #c7d2fe;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.18s;
  &:hover { border-color: #818cf8; background: #eef2ff; }
  img { width: 100%; height: 100%; object-fit: contain; }
}
.captcha-loading { color: #c7d2fe; }

/* ==================== 登录按钮 ==================== */
.login-btn {
  width: 100%;
  height: 46px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 6px;
  position: relative;
  overflow: hidden;
  letter-spacing: 1px;

  &:hover:not(:disabled) {
    background: linear-gradient(135deg, #4f46e5, #7c3aed);
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(99, 102, 241, 0.35);
  }

  &:active:not(:disabled) {
    transform: scale(0.98) translateY(0);
    box-shadow: 0 2px 8px rgba(99, 102, 241, 0.25);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}

.btn-shine {
  position: absolute;
  top: 0;
  left: -75%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.25), transparent);
  transform: skewX(-20deg);
  animation: btn-shine 3s ease-in-out infinite;

  .login-btn--loading & { display: none; }
}

@keyframes btn-shine {
  0%, 100% { left: -75%; }
  50% { left: 130%; }
}

.btn-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

/* ==================== 底部 ==================== */
.card-footer {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #f1f0ff;
  text-align: center;
  font-size: 12px;
  color: #6366f1;

  a {
    color: #6366f1;
    text-decoration: none;
    font-weight: 600;
    &:hover { color: #4f46e5; text-decoration: underline; }
  }
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 11px;
  color: #c7d2fe;
  font-weight: 500;
}

/* 旋转动画 */
.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
