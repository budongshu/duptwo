<template>
  <div class="login-page">
    <!-- 背景 -->
    <div class="bg-pattern"></div>

    <!-- 登录卡片 -->
    <div class="login-wrapper">
      <div class="login-card">
        <!-- Logo 区域 -->
        <div class="card-head">
          <div class="logo">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
              <rect x="2" y="2" width="9" height="9" rx="2" fill="#3b82f6"/>
              <rect x="13" y="2" width="9" height="9" rx="2" fill="#60a5fa"/>
              <rect x="2" y="13" width="9" height="9" rx="2" fill="#60a5fa"/>
              <rect x="13" y="13" width="9" height="9" rx="2" fill="#93c5fd"/>
            </svg>
          </div>
          <h1 class="title">数据登记平台</h1>
          <p class="subtitle">duptwo</p>
        </div>

        <!-- 表单区域 -->
        <div class="card-body">
          <div class="greet">
            <span v-if="!loginError">欢迎登录</span>
            <span v-else class="greet-error">{{ loginError.title }}</span>
          </div>

          <!-- 错误提示 -->
          <transition name="error-slide">
            <div v-if="loginError" class="error-box" :class="`error-box--${loginError.type}`">
              <div class="error-icon">
                <svg v-if="loginError.type === 'password'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <svg v-else-if="loginError.type === 'lockout'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                  <line x1="12" y1="16" x2="12" y2="18"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 8v4M12 16h.01"/>
                </svg>
              </div>
              <div class="error-text">
                <div class="error-main">{{ loginError.desc }}</div>
                <div v-if="loginError.remainingAttempts !== undefined && loginError.remainingAttempts > 0 && loginError.type !== 'lockout'" class="error-remaining">
                  剩余 {{ loginError.remainingAttempts }} / {{ loginError.maxAttempts }} 次尝试机会
                </div>
                <div v-if="loginError.type === 'lockout'" class="error-countdown">
                  {{ formatCountdown(lockoutRemaining) }} 后自动解锁
                </div>
              </div>
              <button class="error-close" @click="loginError = null">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
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
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin">
                      <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                    </svg>
                  </span>
                </div>
              </div>
            </div>

            <button type="submit" class="login-btn" :disabled="loading">
              <span v-if="loading" class="btn-loading">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin">
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                </svg>
                验证中...
              </span>
              <span v-else>登 录</span>
            </button>
          </el-form>

          <div class="card-footer">
            <router-link v-if="registrationEnabled" to="/register">没有账号？立即注册</router-link>
          </div>
        </div>
      </div>

      <!-- 底部 -->
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
  remainingAttempts?: number  // 剩余尝试次数
  maxAttempts?: number       // 最大尝试次数
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
  // 如果已经有定时器在运行，先清除
  if (lockoutTimer) {
    clearInterval(lockoutTimer)
    lockoutTimer = null
  }

  // 设置初始值（分钟转为秒）
  lockoutRemaining.value = minutes * 60

  // 锁定结束时间 = 现在 + 分钟数
  const lockoutEndTime = Date.now() + minutes * 60 * 1000

  loginError.value = {
    type: 'lockout',
    title: '账号已被锁定',
    desc: `由于连续登录失败，为了保护您的账户安全，系统已临时锁定。请等待 ${minutes} 分钟后自动解锁，或联系管理员协助。`,
    lockoutEndTime: lockoutEndTime  // 保存结束时间用于精确计算
  }

  lockoutTimer = setInterval(() => {
    // 实时计算剩余秒数，确保和后端同步
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

// 解析登录错误
const parseLoginError = (msg: string): LoginError => {
  if (msg.includes('锁定')) {
    const minutes = parseInt(msg.match(/(\d+)\s*分钟/)?.[1] || '5')
    return {
      type: 'lockout',
      title: '账号已被锁定',
      desc: `由于连续登录失败，为了保护您的账户安全，系统已临时锁定。请等待 ${minutes} 分钟后自动解锁，或联系管理员协助。`
    }
  }
  if (msg.includes('验证码')) {
    return {
      type: 'captcha',
      title: '验证码错误',
      desc: '验证码填写错误或已过期，请点击图片刷新后重新输入'
    }
  }
  if (msg.includes('请填写验证码')) {
    return {
      type: 'captcha',
      title: '请填写验证码',
      desc: '请先获取并填写图形验证码后再点击登录'
    }
  }
  if (msg.includes('用户不存在') || msg.includes('不存在')) {
    return {
      type: 'username',
      title: '用户不存在',
      desc: '该用户名不存在，请检查输入是否有误'
    }
  }
  if (msg.includes('密码错误') || msg.includes('密码不正确')) {
    return {
      type: 'password',
      title: '密码输入错误',
      desc: '您输入的密码不正确，请检查大小写是否正确'
    }
  }
  if (msg.includes('用户名或密码错误') || msg.includes('不正确')) {
    return {
      type: 'password',
      title: '密码输入错误',
      desc: '用户名或密码不正确，请仔细核对后重新输入'
    }
  }
  if (msg.includes('禁用') || msg.includes('停用')) {
    return {
      type: 'general',
      title: '账号已被停用',
      desc: '您的账号已被管理员停用，如需恢复使用，请联系系统管理员'
    }
  }
  if (msg.includes('未激活') || msg.includes('待审核')) {
    return {
      type: 'general',
      title: '账号待审核',
      desc: '您的账号正在等待管理员审核，审核通过后即可正常登录'
    }
  }
  return {
    type: 'general',
    title: '登录失败',
    desc: msg ? `提示：${msg}` : '登录遇到问题，请稍后重试'
  }
}

const handleLogin = async () => {
  // 1. 前端基础验证（不清空任何数据）
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

  // 2. 开始登录请求
  loginError.value = null
  loading.value = true

  try {
    const res = await AuthApi.login({
      username: form.username,
      password: form.password,
      captchaId: captchaId.value,
      captcha: form.captcha
    })

    // 3. 登录成功
    if (res.data.mfaRequired) {
      currentUserId.value = res.data.user.id
      mfaStep.value = true
    } else {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      router.push('/')
    }
  } catch (err: any) {
    // 4. 登录失败 - 只显示错误，不清空任何输入
    let msg = ''
    let remainingAttempts = 0
    let maxAttempts = 5

    // 拦截器 reject 的格式: { code, message, data }
    if (err && err.message) {
      msg = err.message
      // 尝试从 err.data 中提取剩余尝试次数
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

    // 解析错误类型，并带上剩余尝试次数
    const errorInfo = parseLoginError(msg)
    loginError.value = {
      ...errorInfo,
      remainingAttempts,
      maxAttempts
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  refreshCaptcha()
  // 获取注册功能状态
  try {
    const res = await AuthApi.getRegistrationStatus()
    if (res.code === 200 && res.data) {
      registrationEnabled.value = res.data.registrationEnabled
    }
  } catch {}
})
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
  position: relative;
}

.bg-pattern {
  position: fixed;
  inset: 0;
  background-image: radial-gradient(circle at 1px 1px, #e2e8f0 1px, transparent 0);
  background-size: 24px 24px;
  pointer-events: none;
}

.login-wrapper {
  width: 100%;
  max-width: 480px;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.login-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.card-head {
  padding: 36px 40px 28px;
  text-align: center;
  border-bottom: 1px solid #f1f5f9;
}

.logo {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  background: #eff6ff;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 4px;
}

.subtitle {
  font-size: 13px;
  color: #94a3b8;
  margin: 0;
  font-family: 'SF Mono', monospace;
}

.card-body {
  padding: 28px 40px 36px;
}

.greet {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 20px;

  &-error {
    color: #ef4444;
  }
}

/* 错误提示框 */
.error-box {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 10px;
  margin-bottom: 20px;
  position: relative;

  &--password,
  &--username {
    background: #fef2f2;
    border: 1px solid #fecaca;

    .error-icon { color: #ef4444; }
    .error-main { color: #b91c1c; }
  }

  &--captcha {
    background: #fefce8;
    border: 1px solid #fef08a;

    .error-icon { color: #ca8a04; }
    .error-main { color: #854d0e; }
  }

  &--lockout {
    background: #fff1f2;
    border: 1px solid #fecdd3;

    .error-icon { color: #e11d48; }
    .error-main { color: #9f1239; }
  }

  &--general {
    background: #f8fafc;
    border: 1px solid #e2e8f0;

    .error-icon { color: #64748b; }
    .error-main { color: #475569; }
  }
}

.error-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.error-text {
  flex: 1;
  min-width: 0;
}

.error-main {
  font-size: 13px;
  line-height: 1.5;
}

.error-countdown {
  font-size: 14px;
  font-weight: 700;
  font-family: 'SF Mono', monospace;
  color: #e11d48;
  margin-top: 6px;
}

.error-remaining {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
  font-weight: 500;
}

.error-close {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  transition: all 0.2s;

  &:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #64748b;
  }
}

.error-slide-enter-active,
.error-slide-leave-active {
  transition: all 0.3s ease;
}

.error-slide-enter-from,
.error-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #374151;
    margin-bottom: 8px;
  }
}

:deep(.el-input__wrapper) {
  border-radius: 10px !important;
  box-shadow: 0 0 0 1px #e2e8f0 !important;
  padding: 4px 12px !important;

  &:hover {
    box-shadow: 0 0 0 1px #cbd5e1 !important;
  }

  &.is-focus {
    box-shadow: 0 0 0 2px #3b82f6 !important;
  }
}

:deep(.el-input__inner) {
  font-size: 14px !important;
  color: #1e293b !important;

  &::placeholder {
    color: #94a3b8 !important;
  }
}

/* 验证码 */
.captcha-row {
  display: flex;
  gap: 10px;

  .el-input {
    flex: 1;
  }
}

.captcha-img {
  width: 100px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s;

  &:hover {
    border-color: #3b82f6;
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
}

.captcha-loading {
  color: #94a3b8;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 48px;
  border: none;
  border-radius: 10px;
  background: #3b82f6;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 8px;

  &:hover:not(:disabled) {
    background: #2563eb;
  }

  &:active:not(:disabled) {
    transform: scale(0.98);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}

.btn-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

/* 底部 */
.card-footer {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f1f5f9;
  text-align: center;
  font-size: 13px;
  color: #64748b;

  a {
    color: #3b82f6;
    text-decoration: none;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 12px;
  color: #94a3b8;
}
</style>
