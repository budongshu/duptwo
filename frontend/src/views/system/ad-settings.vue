<template>
  <div class="ad-page">
    <!-- 页面标题栏 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">AD域配置</h1>
        <span class="page-subtitle">Active Directory 域用户认证</span>
      </div>
    </header>

    <!-- 总开关卡片 -->
    <div class="enable-card" :class="{ 'enable-card--active': form.enabled }">
      <div class="enable-left">
        <div class="enable-icon" :class="form.enabled ? 'icon-active' : 'icon-inactive'">
          <svg v-if="form.enabled" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
        </div>
        <div class="enable-text">
          <div class="enable-title">{{ form.enabled ? 'AD域认证已启用' : 'AD域认证未启用' }}</div>
          <div class="enable-desc">{{ form.enabled ? '系统登录页面将显示"域账号(AD)"登录选项' : '开启后支持 Active Directory 域用户登录' }}</div>
        </div>
      </div>
      <el-switch v-model="form.enabled" size="large" @change="onEnabledChange" />
    </div>

    <!-- 配置内容（disabled状态统一覆盖） -->
    <div class="config-content" :class="{ 'config-content--disabled': !form.enabled }">

      <!-- 连接配置区块 -->
      <div class="config-card">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #005eeb"></div>
          <span class="config-card-title">连接配置</span>
          <span class="config-card-badge">服务器</span>
        </div>
        <div class="config-card-body">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top" :disabled="loading">
            <div class="form-grid-2">
              <div class="form-field">
                <label class="field-label">服务器地址 <span class="required-mark">*</span></label>
                <el-input v-model="form.server" placeholder="ldap://192.168.1.100 或 192.168.1.100" clearable size="default" class="field-input">
                  <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg></template>
                </el-input>
              </div>
              <div class="form-field">
                <label class="field-label">端口 <span class="required-mark">*</span></label>
                <el-input-number v-model="form.port" :min="1" :max="65535" placeholder="389" style="width: 100%" size="default" />
                <div class="field-hint">标准LDAP: 389, LDAPS: 636</div>
              </div>
            </div>
            <div class="form-field">
              <label class="field-label">基准DN（Base DN） <span class="required-mark">*</span></label>
              <el-input v-model="form.base_dn" placeholder="OU=Users,DC=company,DC=com" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg></template>
              </el-input>
              <div class="field-hint">用户和组的搜索起点</div>
            </div>
            <div class="ssl-toggle-row">
              <el-switch v-model="form.use_ssl" size="small" />
              <span class="ssl-label">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                使用 LDAPS 加密连接（SSL/TLS）
              </span>
            </div>
          </el-form>
        </div>
      </div>

      <!-- 认证配置区块 -->
      <div class="config-card">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #06b6d4"></div>
          <span class="config-card-title">认证配置</span>
          <span class="config-card-badge">安全</span>
        </div>
        <div class="config-card-body">
          <el-form ref="formRef2" :model="form" :rules="rules" label-position="top" :disabled="loading">
            <div class="form-field">
              <label class="field-label">管理账号DN（Bind DN） <span class="required-mark">*</span></label>
              <el-input v-model="form.bind_dn" placeholder="CN=admin,DC=company,DC=com" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
              </el-input>
              <div class="field-hint">具有搜索用户权限的管理员账号DN</div>
            </div>
            <div class="form-field">
              <label class="field-label">管理账号密码</label>
              <el-input v-model="form.bind_password" type="password" placeholder="留空则保持原密码不变" show-password clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg></template>
              </el-input>
              <div class="field-hint">留空表示不修改密码</div>
            </div>
            <div class="form-field">
              <label class="field-label">用户搜索过滤器</label>
              <el-input v-model="form.user_filter" placeholder="(sAMAccountName=%s)" clearable size="default" class="field-input">
                <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></template>
              </el-input>
              <div class="field-hint">
                <code>%s</code> 替换为登录用户名。常用：
                <code>(sAMAccountName=%s)</code>（AD）、<code>(uid=%s)</code>（LDAP）
              </div>
            </div>
          </el-form>
        </div>
      </div>

      <!-- 自动注册区块 -->
      <div class="config-card" :class="{ 'config-card--warn': form.enabled && form.auto_register }">
        <div class="config-card-header">
          <div class="config-card-dot" style="background: #f59e0b"></div>
          <span class="config-card-title">自动注册</span>
          <span class="config-card-badge config-card-badge--warn">用户</span>
        </div>
        <div class="config-card-body">
          <div class="auto-reg-row">
            <div class="auto-reg-info">
              <div class="auto-reg-title">允许AD用户自动注册</div>
              <div class="auto-reg-desc">首次登录时自动创建本地账号（Source = AD）</div>
            </div>
            <el-switch v-model="form.auto_register" :disabled="!form.enabled" size="small" />
          </div>

          <div v-if="form.enabled && form.auto_register" class="auto-reg-expanded">
            <div class="policy-box">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              <div class="policy-content">
                <div class="policy-title">自动注册策略</div>
                <ul class="policy-list">
                  <li>首次AD登录时，系统自动在本地创建用户记录</li>
                  <li>若本地已有同名用户（Source=LOCAL），则拒绝注册</li>
                  <li>AD用户密码由域控制器管理，本地不存储</li>
                </ul>
              </div>
            </div>
            <div class="form-field" style="margin-top: 16px">
              <label class="field-label">新AD用户的默认角色</label>
              <el-select v-model="form.default_role_id" placeholder="请选择默认角色" size="default" style="width: 280px">
                <el-option v-for="role in roles" :key="role.id" :label="role.name + '（' + role.code + '）'" :value="role.id" />
              </el-select>
              <div class="field-hint">自动注册的AD用户将被授予该角色</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作区 -->
      <div class="actions-card">
        <div class="actions-left">
          <el-button @click="handleTest" :loading="testLoading" :disabled="!form.enabled || !form.server" size="default" class="test-btn">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="13 2 13 9 20 9"/><path d="M20 14.66V17a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h2"/></svg>
            测试连接
          </el-button>
          <!-- 测试结果 -->
          <transition name="result-fade">
            <span v-if="testResult !== null" class="test-result" :class="testResult ? 'result--ok' : 'result--fail'">
              <svg v-if="testResult" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
              <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              {{ testResult ? '连接成功' : '连接失败' }}
            </span>
          </transition>
        </div>
        <div class="actions-right">
          <el-button @click="loadConfig" :disabled="loading" size="default">重置</el-button>
          <el-button type="primary" @click="handleSave" :loading="saving" :disabled="!form.enabled" size="default" class="save-btn">
            <svg v-if="!saving" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
            保存配置
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { AdminApi, type ADConfig } from '@/api/admin'
import { RoleApi } from '@/api/role'

const loading = ref(false)
const saving = ref(false)
const testLoading = ref(false)
const testResult = ref<boolean | null>(null)
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

const rules = {
  server: [{ required: true, message: '请输入AD服务器地址', trigger: 'blur' }],
  base_dn: [{ required: true, message: '请输入Base DN', trigger: 'blur' }],
  bind_dn: [{ required: true, message: '请输入管理账号DN', trigger: 'blur' }]
}

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await AdminApi.getADConfig()
    if (res.code === 200) {
      Object.assign(form, res.data)
      form.bind_password = '' // 清空密码显示
    }
  } catch (e) {
    console.error('加载配置失败', e)
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
    console.error('加载角色列表失败', e)
  }
}

const onEnabledChange = (val: boolean) => {
  if (!val) {
    form.auto_register = false
  }
}

const handleSave = async () => {
  // 合并两个form的校验
  const valid1 = await formRef.value?.validate().catch(() => false)
  const valid2 = await formRef2.value?.validate().catch(() => false)
  if (!valid1 || !valid2) return

  saving.value = true
  try {
    const res = await AdminApi.updateADConfig(form)
    if (res.code === 200) {
      ElMessage.success('AD配置已保存')
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
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
      ElMessage.success(res.message || '连接测试成功')
    } else {
      ElMessage.error(res.message || '连接测试失败')
    }
  } catch (e: any) {
    testResult.value = false
    ElMessage.error(e.message || '连接测试失败')
  } finally {
    testLoading.value = false
  }
}

onMounted(() => {
  testResult.value = null
  loadConfig()
  loadRoles()
})
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.ad-page {
  padding: var(--space-4);
  min-height: 100vh;
  background: var(--color-page-bg);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  overflow: visible;
}

/* ==================== 页面标题栏 ==================== */
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

/* ==================== 总开关卡片 ==================== */
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
  animation: card-rise 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both;

  &--active {
    border-color: rgba(0, 94, 235, 0.3);
    box-shadow: 0 4px 20px rgba(0, 94, 235, 0.1), var(--shadow-xs);
    background: linear-gradient(135deg, var(--color-surface) 0%, rgba(0, 94, 235, 0.03) 100%);
  }
}

@keyframes card-rise {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.enable-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.enable-icon {
  width: 52px;
  height: 52px;
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
  letter-spacing: -0.2px;
}

.enable-desc {
  font-size: 12px;
  color: var(--color-text-muted);
}

/* ==================== 配置内容 ==================== */
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

/* ==================== 配置卡片区块 ==================== */
.config-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-xs);
  overflow: hidden;
  animation: card-rise 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both;
  animation-delay: 0.05s;
  animation-fill-mode: both;

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
  background: var(--color-surface-2);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.3px;

  &--warn {
    background: var(--color-warning-bg);
    color: var(--color-warning);
  }
}

.config-card-body {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* ==================== 表单样式 ==================== */
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

/* ==================== 自动注册 ==================== */
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

/* ==================== 操作区 ==================== */
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
  animation: card-rise 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both 0.1s;
  animation-fill-mode: both;
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
}
</style>
