<template>
  <div class="login-page">
    <!-- 顶部导航栏 -->
    <header class="top-bar">
      <div class="top-bar-inner">
        <div class="brand">
          <div class="brand-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <path d="M12 2L2 7l10 5 10-5-10-5z" fill="rgba(255,255,255,0.9)"/>
              <path d="M2 17l10 5 10-5M2 12l10 5 10-5" stroke="rgba(255,255,255,0.6)" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <span class="brand-name">数据登记管理平台</span>
        </div>
        <div class="top-bar-right">
          <span class="version-tag">v1.0.0</span>
        </div>
      </div>
    </header>

    <!-- 背景：点阵网格 -->
    <div class="bg-grid"></div>

    <!-- 注册卡片 -->
    <main class="login-main">
      <div class="login-card">
        <div class="card-title">
          <h2>用户注册</h2>
          <p>创建您的账号</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          class="login-form"
          @submit.prevent="handleRegister"
        >
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="用户名"
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>
          <el-form-item prop="nickname">
            <el-input
              v-model="form.nickname"
              placeholder="昵称（选填）"
              size="large"
              :prefix-icon="UserFilled"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码"
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="确认密码"
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleRegister"
          >
            注 册
          </el-button>
        </el-form>

        <div class="login-footer">
          <span>已有账号？</span>
          <router-link to="/login" class="link">立即登录</router-link>
        </div>
      </div>
    </main>

    <!-- 底部信息 -->
    <footer class="page-footer">
      <span>DataRegistry</span>
      <span class="sep">·</span>
      <span>2024</span>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, UserFilled } from '@element-plus/icons-vue'
import { AuthApi } from '@/api/auth'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

const form = reactive({
  username: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度为 3-32 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 64, message: '密码长度为 6-64 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await AuthApi.register({
      username: form.username,
      password: form.password,
      nickname: form.nickname || undefined
    })
    if (res.code === 200) {
      ElMessage.success('注册成功，请登录')
      router.push('/login')
    } else {
      ElMessage.error(res.message || '注册失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
// ========== 页面布局（全屏铺满） ==========
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #0f1623;
  position: relative;
  overflow: hidden;
}

// ========== 顶部导航栏 ==========
.top-bar {
  position: relative;
  z-index: 10;
  height: 56px;
  background: rgba(15, 22, 35, 0.95);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(12px);
  flex-shrink: 0;
}

.top-bar-inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-name {
  font-size: 15px;
  font-weight: 600;
  color: #e2e8f0;
  letter-spacing: 0.5px;
}

.top-bar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-tag {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 8px;
  border-radius: 4px;
}

// ========== 背景点阵网格 ==========
.bg-grid {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background-image: radial-gradient(circle, rgba(255, 255, 255, 0.06) 1px, transparent 1px);
  background-size: 28px 28px;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse 80% 60% at 50% 50%, transparent 30%, #0f1623 100%);
  }
}

// ========== 主内容区 ==========
.login-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 16px;
  position: relative;
  z-index: 1;
}

// ========== 登录卡片 ==========
.login-card {
  width: 400px;
  max-width: 100%;
  background: rgba(255, 255, 255, 0.97);
  border-radius: 16px;
  box-shadow:
    0 4px 24px rgba(0, 0, 0, 0.3),
    0 1px 4px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

// ========== 卡片标题 ==========
.card-title {
  padding: 32px 36px 24px;
  border-bottom: 1px solid #f1f5f9;

  h2 {
    font-size: 20px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 6px 0;
    line-height: 1.2;
  }

  p {
    font-size: 13px;
    color: #94a3b8;
    margin: 0;
  }
}

// ========== 表单 ==========
.login-form {
  padding: 24px 36px 8px;

  :deep(.el-form-item) {
    margin-bottom: 18px;
  }

  :deep(.el-input__wrapper) {
    border-radius: 10px;
    box-shadow: 0 0 0 1px #e2e8f0;
    padding: 4px 12px;
    transition: box-shadow 0.2s;

    &:hover {
      box-shadow: 0 0 0 1px #cbd5e1;
    }

    &.is-focus {
      box-shadow: 0 0 0 2px #3b82f6;
    }
  }

  :deep(.el-input__inner) {
    font-size: 14px;
    color: #1e293b;

    &::placeholder {
      color: #94a3b8;
    }
  }

  :deep(.el-input__prefix .el-icon) {
    color: #94a3b8;
    margin-right: 2px;
  }
}

// ========== 登录按钮 ==========
.login-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 2px;
  border-radius: 10px;
  border: none;
  background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%);
  color: #fff;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
    background: linear-gradient(135deg, #2563eb 0%, #4f46e5 100%);
  }

  &:active {
    transform: translateY(0);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    transform: none;
  }
}

// ========== 注册入口 ==========
.login-footer {
  text-align: center;
  padding: 16px 36px 28px;
  font-size: 13px;
  color: #64748b;

  .link {
    color: #3b82f6;
    text-decoration: none;
    margin-left: 4px;
    font-weight: 500;

    &:hover {
      color: #2563eb;
    }
  }
}

// ========== 底部信息 ==========
.page-footer {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 12px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.2);
  flex-shrink: 0;

  .sep {
    margin: 0 6px;
  }
}
</style>
