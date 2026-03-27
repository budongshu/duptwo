<template>
  <div class="security-page">
    <!-- 页面标题栏 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">登录安全</h1>
        <span class="page-subtitle">密码策略 · 登录限制 · 会话管理</span>
      </div>
      <div class="header-actions">
        <el-button type="primary" :loading="saving" @click="handleSave" size="default" class="save-btn">
          <svg v-if="!saving" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          保存全部设置
        </el-button>
      </div>
    </header>

    <!-- 安全态势总览 -->
    <div class="security-overview" v-if="!loadingOverview">
      <div class="overview-item overview-item--warn" @click="scrollTo('login')">
        <div class="overview-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
        </div>
        <div class="overview-data">
          <div class="overview-value">{{ overview.totalLockedUsers }}</div>
          <div class="overview-label">锁定用户</div>
        </div>
        <div class="overview-arrow">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
        </div>
      </div>
      <div class="overview-item overview-item--danger" @click="scrollTo('ip')">
        <div class="overview-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
        </div>
        <div class="overview-data">
          <div class="overview-value">{{ overview.totalLockedIPs }}</div>
          <div class="overview-label">封禁IP</div>
        </div>
        <div class="overview-arrow">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
        </div>
      </div>
      <div class="overview-item overview-item--success" @click="scrollTo('ip')">
        <div class="overview-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        </div>
        <div class="overview-data">
          <div class="overview-value">{{ overview.whitelistCount || 0 }}</div>
          <div class="overview-label">IP白名单</div>
        </div>
        <div class="overview-arrow">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
        </div>
      </div>
      <div class="overview-item overview-item--info">
        <div class="overview-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
        </div>
        <div class="overview-data">
          <div class="overview-value">{{ overview.blacklistCount || 0 }}</div>
          <div class="overview-label">IP黑名单</div>
        </div>
      </div>
    </div>

    <!-- 主体内容：4个配置区块横排 -->
    <div class="config-grid">

      <!-- 区块1：登录安全 -->
      <div class="config-block" id="block-login">
        <div class="config-block-header">
          <div class="block-icon block-icon--blue">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          </div>
          <div class="block-titles">
            <div class="block-title">登录安全</div>
            <div class="block-subtitle">验证码 · 不活跃用户 · 登录限制</div>
          </div>
          <div class="block-badge" :class="form.captchaEnabled ? 'badge--on' : 'badge--off'">
            {{ form.captchaEnabled ? '验证码已启用' : '验证码未启用' }}
          </div>
        </div>
        <div class="config-block-body">
          <!-- 验证码 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              登录验证码
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">启用登录验证码</div>
                <div class="field-desc">连续登录失败N次后自动显示图形验证码</div>
              </div>
              <el-switch v-model="form.captchaEnabled" size="small" />
            </div>
            <div class="field-row field-row--indent" v-if="form.captchaEnabled">
              <div class="field-info">
                <div class="field-name">触发阈值（次）</div>
                <div class="field-desc">连续失败此次数后启用验证码</div>
              </div>
              <el-input-number v-model="form.captchaMinLen" :min="1" :max="10" size="small" controls-position="right" />
            </div>
          </div>

          <!-- 不活跃用户 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              不活跃自动禁用
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">启用不活跃自动禁用</div>
                <div class="field-desc">超过指定天数未登录的用户自动标记为禁用</div>
              </div>
              <el-switch v-model="form.inactiveAutoDisable" size="small" />
            </div>
            <div class="field-row field-row--indent" v-if="form.inactiveAutoDisable">
              <div class="field-info">
                <div class="field-name">不活跃天数阈值</div>
                <div class="field-desc">超过此天数未登录将被自动禁用</div>
              </div>
              <el-input-number v-model="form.inactiveDaysThreshold" :min="1" :max="365" size="small" controls-position="right" />
            </div>
          </div>

          <!-- 用户登录限制 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/></svg>
              用户级别登录限制
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">最大连续失败次数</div>
                <div class="field-desc">同一用户名连续失败达到此次数后账号被锁定</div>
              </div>
              <el-input-number v-model="form.userLoginMaxAttempts" :min="1" :max="20" size="small" controls-position="right" />
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">账号锁定时长</div>
                <div class="field-desc">锁定后自动解锁的等待时间</div>
              </div>
              <div class="input-with-unit">
                <el-input-number v-model="form.userLoginLockMinutes" :min="1" :max="10080" size="small" controls-position="right" style="width: 100px" />
                <span class="unit-label">分钟</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 区块2：IP安全 -->
      <div class="config-block" id="block-ip">
        <div class="config-block-header">
          <div class="block-icon block-icon--amber">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>
          </div>
          <div class="block-titles">
            <div class="block-title">IP访问控制</div>
            <div class="block-subtitle">登录限制 · 黑白名单</div>
          </div>
          <div class="block-badge" :class="form.ipLoginMaxAttempts > 0 ? 'badge--on' : 'badge--off'">
            {{ form.ipLoginMaxAttempts > 0 ? 'IP限流已启用' : '未限制' }}
          </div>
        </div>
        <div class="config-block-body">
          <!-- IP登录限制 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/></svg>
              IP级别登录限制
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">最大连续失败次数</div>
                <div class="field-desc">同一IP连续失败达到此次数后IP被封禁</div>
              </div>
              <el-input-number v-model="form.ipLoginMaxAttempts" :min="1" :max="100" size="small" controls-position="right" />
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">IP封禁时长</div>
                <div class="field-desc">封禁后自动解封的等待时间</div>
              </div>
              <div class="input-with-unit">
                <el-input-number v-model="form.ipLoginLockMinutes" :min="1" :max="10080" size="small" controls-position="right" style="width: 100px" />
                <span class="unit-label">分钟</span>
              </div>
            </div>
          </div>

          <!-- IP白名单 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
              IP白名单
            </div>
            <div class="field-row field-row--vertical">
              <div class="field-info">
                <div class="field-name">允许访问的IP / 网段</div>
                <div class="field-desc">多个用英文逗号分隔，支持CIDR格式。白名单非空时仅白名单IP可登录。</div>
              </div>
              <el-input v-model="form.ipWhitelist" type="textarea" :rows="3" placeholder="如：192.168.1.100, 10.0.0.0/8" clearable size="default" />
            </div>
          </div>

          <!-- IP黑名单 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
              IP黑名单
            </div>
            <div class="field-row field-row--vertical">
              <div class="field-info">
                <div class="field-name">禁止访问的IP / 网段</div>
                <div class="field-desc">黑名单优先级高于白名单，可一次性封禁整个网段。</div>
              </div>
              <el-input v-model="form.ipBlacklist" type="textarea" :rows="3" placeholder="如：1.2.3.4, 5.6.7.0/24" clearable size="default" />
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- 第二行：密码安全 + 会话管理 -->
    <div class="config-grid config-grid--2">

      <!-- 区块3：密码策略 -->
      <div class="config-block" id="block-password">
        <div class="config-block-header">
          <div class="block-icon block-icon--green">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          </div>
          <div class="block-titles">
            <div class="block-title">密码策略</div>
            <div class="block-subtitle">强度要求 · 有效期</div>
          </div>
          <div class="block-badge" :class="form.passwordMinLength >= 8 ? 'badge--on' : 'badge--warn'">
            最小{{ form.passwordMinLength }}位
          </div>
        </div>
        <div class="config-block-body">
          <!-- 密码强度 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              密码强度要求
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">密码最小长度</div>
                <div class="field-desc">密码至少需要包含的字符数，建议不少于8位</div>
              </div>
              <el-input-number v-model="form.passwordMinLength" :min="6" :max="128" size="small" controls-position="right" />
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">密码过期时间</div>
                <div class="field-desc">0表示永不过期，单位为天</div>
              </div>
              <div class="input-with-unit">
                <el-input-number v-model="form.passwordExpiryDays" :min="0" :max="365" size="small" controls-position="right" style="width: 100px" />
                <span class="unit-label">天</span>
              </div>
            </div>
          </div>

          <!-- 密码组成 -->
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              密码组成规则
            </div>
            <div class="toggle-grid">
              <div class="toggle-item" :class="{ active: form.passwordRequireUppercase }" @click="form.passwordRequireUppercase = !form.passwordRequireUppercase">
                <div class="toggle-check">
                  <svg v-if="form.passwordRequireUppercase" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                </div>
                <div class="toggle-text">
                  <div class="toggle-name">大写字母</div>
                  <div class="toggle-desc">必须包含 A-Z</div>
                </div>
              </div>
              <div class="toggle-item" :class="{ active: form.passwordRequireLowercase }" @click="form.passwordRequireLowercase = !form.passwordRequireLowercase">
                <div class="toggle-check">
                  <svg v-if="form.passwordRequireLowercase" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                </div>
                <div class="toggle-text">
                  <div class="toggle-name">小写字母</div>
                  <div class="toggle-desc">必须包含 a-z</div>
                </div>
              </div>
              <div class="toggle-item" :class="{ active: form.passwordRequireDigit }" @click="form.passwordRequireDigit = !form.passwordRequireDigit">
                <div class="toggle-check">
                  <svg v-if="form.passwordRequireDigit" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                </div>
                <div class="toggle-text">
                  <div class="toggle-name">数字</div>
                  <div class="toggle-desc">必须包含 0-9</div>
                </div>
              </div>
              <div class="toggle-item" :class="{ active: form.passwordRequireSpecial }" @click="form.passwordRequireSpecial = !form.passwordRequireSpecial">
                <div class="toggle-check">
                  <svg v-if="form.passwordRequireSpecial" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                </div>
                <div class="toggle-text">
                  <div class="toggle-name">特殊字符</div>
                  <div class="toggle-desc">必须包含 !@#$%^&*</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 区块4：会话管理 -->
      <div class="config-block" id="block-session">
        <div class="config-block-header">
          <div class="block-icon block-icon--purple">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          </div>
          <div class="block-titles">
            <div class="block-title">会话管理</div>
            <div class="block-subtitle">超时配置</div>
          </div>
        </div>
        <div class="config-block-body">
          <div class="field-group">
            <div class="field-group-label">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              会话超时
            </div>
            <div class="field-row">
              <div class="field-info">
                <div class="field-name">会话超时时长</div>
                <div class="field-desc">用户无操作后自动退出登录的时长</div>
              </div>
              <div class="input-with-unit">
                <el-input-number v-model="form.sessionTimeoutHours" :min="1" :max="168" size="small" controls-position="right" style="width: 100px" />
                <span class="unit-label">小时</span>
              </div>
            </div>
          </div>

          <!-- 当前锁定情况 -->
          <div class="lock-widgets">
            <div class="lock-widget" @click="loadLockedUsers(); showLockPanel = 'users'">
              <div class="lock-widget-icon lock-widget-icon--warn">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </div>
              <div class="lock-widget-data">
                <div class="lock-widget-num">{{ overview.totalLockedUsers }}</div>
                <div class="lock-widget-label">锁定用户</div>
              </div>
            </div>
            <div class="lock-widget" @click="loadLockedIPs(); showLockPanel = 'ips'">
              <div class="lock-widget-icon lock-widget-icon--danger">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
              </div>
              <div class="lock-widget-data">
                <div class="lock-widget-num">{{ overview.totalLockedIPs }}</div>
                <div class="lock-widget-label">封禁IP</div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- 锁定面板（弹出一个drawer） -->
    <el-drawer v-model="lockDrawerVisible" size="480px" direction="rtl">
      <template #header>
        <div class="lock-drawer-head">
          <span class="lock-drawer-tag">锁定</span>
          <span class="lock-drawer-title">{{ lockPanelTitle }}</span>
        </div>
      </template>
      <template v-if="showLockPanel === 'users'">
        <div class="lock-list">
          <div v-if="lockedUsers.length === 0" class="lock-empty">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
            <span>暂无被锁定的用户</span>
          </div>
          <div v-for="user in lockedUsers" :key="user.target" class="lock-record">
            <div class="lock-record-left">
              <div class="lock-record-name">{{ user.target }}</div>
              <div class="lock-record-meta">
                <span>失败 {{ user.failCount }} 次</span>
                <span>{{ formatTime(user.lockedAt) }}</span>
              </div>
            </div>
            <el-button type="danger" size="small" plain @click="handleUnlockUser(user.target)">解锁</el-button>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="lock-list">
          <div v-if="lockedIPs.length === 0" class="lock-empty">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/></svg>
            <span>暂无被封禁的IP</span>
          </div>
          <div v-for="ip in lockedIPs" :key="ip.target" class="lock-record">
            <div class="lock-record-left">
              <div class="lock-record-name">{{ ip.target }}</div>
              <div class="lock-record-meta">
                <span>失败 {{ ip.failCount }} 次</span>
                <span>{{ formatTime(ip.lockedAt) }}</span>
              </div>
            </div>
            <el-button type="danger" size="small" plain @click="handleUnlockIP(ip.target)">解锁</el-button>
          </div>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { AdminApi, type SecuritySettings, type LockoutRecord } from '@/api/admin'

const saving = ref(false)
const loadingOverview = ref(false)
const loadingLockUsers = ref(false)
const loadingLockIPs = ref(false)
const lockedUsers = ref<LockoutRecord[]>([])
const lockedIPs = ref<LockoutRecord[]>([])
const lockDrawerVisible = ref(false)
const showLockPanel = ref<'users' | 'ips'>('users')

const lockPanelTitle = computed(() => showLockPanel.value === 'users' ? '锁定用户列表' : '封禁IP列表')

const overview = reactive({
  totalLockedUsers: 0,
  totalLockedIPs: 0,
  whitelistCount: 0,
  blacklistCount: 0,
})

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
  passwordMinLength: 6,
  passwordRequireUppercase: false,
  passwordRequireLowercase: false,
  passwordRequireDigit: false,
  passwordRequireSpecial: false,
  sessionTimeoutHours: 24,
})

const loadSettings = async () => {
  try {
    const res = await AdminApi.getSecuritySettings()
    if (res.code === 200 && res.data) {
      Object.assign(form, res.data)
    }
  } catch (e) {
    ElMessage.error('加载安全设置失败')
  }
}

const loadOverview = async () => {
  try {
    const res = await AdminApi.getSecurityOverview()
    if (res.code === 200 && res.data) {
      Object.assign(overview, res.data)
    }
  } catch (e) {
    // ignore
  }
}

const loadLockedUsers = async () => {
  loadingLockUsers.value = true
  try {
    const res = await AdminApi.getLockedUsers()
    if (res.code === 200) {
      lockedUsers.value = res.data || []
    }
  } catch (e) {
    ElMessage.error('加载锁定用户失败')
  } finally {
    loadingLockUsers.value = false
  }
}

const loadLockedIPs = async () => {
  loadingLockIPs.value = true
  try {
    const res = await AdminApi.getLockedIPs()
    if (res.code === 200) {
      lockedIPs.value = res.data || []
    }
  } catch (e) {
    ElMessage.error('加载锁定IP失败')
  } finally {
    loadingLockIPs.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    const res = await AdminApi.updateSecuritySettings(form)
    if (res.code === 200) {
      ElMessage.success('安全设置已保存')
      loadOverview()
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleUnlockUser = async (username: string) => {
  try {
    await ElMessageBox.confirm(`确定要解锁用户 "${username}" 吗？`, '解锁确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const res = await AdminApi.unlockUser(username)
    if (res.code === 200) {
      ElMessage.success('用户已解锁')
      loadLockedUsers()
      loadOverview()
    } else {
      ElMessage.error(res.message || '解锁失败')
    }
  } catch {}
}

const handleUnlockIP = async (ip: string) => {
  try {
    await ElMessageBox.confirm(`确定要解锁IP "${ip}" 吗？`, '解锁确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const res = await AdminApi.unlockIP(ip)
    if (res.code === 200) {
      ElMessage.success('IP已解锁')
      loadLockedIPs()
      loadOverview()
    } else {
      ElMessage.error(res.message || '解锁失败')
    }
  } catch {}
}

const scrollTo = (id: string) => {
  const el = document.getElementById('block-' + id)
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const formatTime = (timeStr: string) => {
  if (!timeStr) return '-'
  const d = new Date(timeStr)
  return d.toLocaleString('zh-CN')
}

onMounted(() => {
  loadSettings()
  loadingOverview.value = true
  loadOverview().finally(() => { loadingOverview.value = false })
  loadLockedUsers()
  loadLockedIPs()
})
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.security-page {
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
  animation: card-rise 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes card-rise {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
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

.header-actions { display: flex; gap: var(--space-2); }
.save-btn { font-weight: 700; }

/* ==================== 安全态势总览 ==================== */
.security-overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-3);
  animation: card-rise 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both 0.05s;
  animation-fill-mode: both;
}

.overview-item {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  padding: var(--space-4);
  display: flex;
  align-items: center;
  gap: var(--space-3);
  box-shadow: var(--shadow-xs);
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: 4px;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-sm);
  }

  &--warn::before { background: var(--color-warning); }
  &--danger::before { background: var(--color-danger); }
  &--success::before { background: var(--color-success); }
  &--info::before { background: var(--color-primary); }
}

.overview-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .overview-item--warn & { background: var(--color-warning-bg); color: var(--color-warning); }
  .overview-item--danger & { background: var(--color-danger-bg); color: var(--color-danger); }
  .overview-item--success & { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
  .overview-item--info & { background: var(--color-primary-light-9); color: var(--color-primary); }
}

.overview-data { flex: 1; }

.overview-value {
  font-family: 'Manrope', sans-serif;
  font-size: 28px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1;
  letter-spacing: -0.5px;
}

.overview-label {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 4px;
  font-weight: 600;
}

.overview-arrow {
  color: var(--color-text-muted);
  opacity: 0;
  transition: opacity 0.2s ease;
  .overview-item:hover & { opacity: 1; }
}

/* ==================== 配置网格 ==================== */
.config-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
  animation: card-rise 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both 0.1s;
  animation-fill-mode: both;
}

.config-grid--2 {
  grid-template-columns: 1fr 1fr;
  animation-delay: 0.15s;
}

/* ==================== 配置区块 ==================== */
.config-block {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-xs);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.config-block-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: relative;
}

.block-icon {
  width: 42px;
  height: 42px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &--blue { background: rgba(0, 94, 235, 0.1); color: var(--color-primary); }
  &--amber { background: rgba(245, 158, 11, 0.1); color: var(--chart-amber); }
  &--green { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
  &--purple { background: rgba(139, 92, 246, 0.1); color: var(--chart-purple); }
}

.block-titles {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.block-title {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 800;
  color: var(--color-text-primary);
  letter-spacing: -0.2px;
}

.block-subtitle {
  font-size: 11px;
  color: var(--color-text-muted);
  font-weight: 500;
}

.block-badge {
  font-size: 10.5px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: var(--radius-full);
  white-space: nowrap;

  &.badge--on { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
  &.badge--off { background: var(--gray-100); color: var(--color-text-muted); }
  &.badge--warn { background: var(--color-warning-bg); color: var(--color-warning); }
}

.config-block-body {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* ==================== 字段分组 ==================== */
.field-group {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: var(--space-3);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field-group-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10.5px;
  font-weight: 700;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 2px;

  svg { color: var(--color-primary); flex-shrink: 0; }
}

.field-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: 6px 0;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  &:last-child { border-bottom: none; }

  &--indent { padding-left: var(--space-3); }
  &--vertical { flex-direction: column; align-items: stretch; gap: var(--space-2); }
}

.field-info { flex: 1; min-width: 0; }

.field-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.field-desc {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
}

.input-with-unit {
  display: flex;
  align-items: center;
  gap: 6px;
}

.unit-label {
  font-size: 12px;
  color: var(--color-text-muted);
  font-weight: 600;
  white-space: nowrap;
}

/* ==================== 密码开关网格 ==================== */
.toggle-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-2);
  margin-top: var(--space-1);
}

.toggle-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover { border-color: var(--color-primary); }

  &.active {
    border-color: var(--color-success);
    background: rgba(22, 163, 74, 0.05);
    .toggle-check { background: var(--color-success); border-color: var(--color-success); color: white; }
    .toggle-name { color: var(--color-success); }
  }
}

.toggle-check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1.5px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.toggle-text { display: flex; flex-direction: column; gap: 1px; }

.toggle-name {
  font-size: 12.5px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.toggle-desc {
  font-size: 10px;
  color: var(--color-text-muted);
}

/* ==================== 锁定widgets ==================== */
.lock-widgets {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-2);
}

.lock-widget {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    border-color: var(--color-primary);
    transform: scale(1.01);
  }
}

.lock-widget-icon {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &--warn { background: var(--color-warning-bg); color: var(--color-warning); }
  &--danger { background: var(--color-danger-bg); color: var(--color-danger); }
}

.lock-widget-data { display: flex; flex-direction: column; gap: 2px; }

.lock-widget-num {
  font-family: 'Manrope', sans-serif;
  font-size: 20px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1;
}

.lock-widget-label {
  font-size: 11px;
  color: var(--color-text-muted);
  font-weight: 600;
}

/* ==================== 锁定列表（drawer） ==================== */
.lock-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  padding: var(--space-3);
}

.lock-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-8) var(--space-4);
  color: var(--color-text-muted);
  font-size: 13px;
  font-weight: 500;
  svg { opacity: 0.4; }
}

.lock-record {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  transition: all 0.15s ease;

  &:hover {
    border-color: var(--color-primary);
    background: var(--color-primary-light-9);
  }
}

.lock-record-left { display: flex; flex-direction: column; gap: 3px; }

.lock-record-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  font-family: 'SF Mono', monospace;
}

.lock-record-meta {
  display: flex;
  gap: var(--space-3);
  font-size: 11px;
  color: var(--color-text-muted);
}

/* ==================== Element Plus覆写 ==================== */
:deep(.el-input-number .el-input__inner) { text-align: left; }
:deep(.el-input__wrapper) { border-radius: var(--radius-sm) !important; box-shadow: none !important; border: 1.5px solid var(--color-border-light) !important; &:focus { border-color: var(--color-primary) !important; } }
:deep(.el-textarea__inner) { border-radius: var(--radius-sm) !important; box-shadow: none !important; border: 1.5px solid var(--color-border-light) !important; &:focus { border-color: var(--color-primary) !important; } }
:deep(.el-drawer__header) { padding: 16px 20px; margin-bottom: 0; border-bottom: 1px solid var(--color-border-light); color: var(--color-text-primary); }
:deep(.el-drawer__body) { padding: 0; overflow-y: auto; }

.lock-drawer-head {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.lock-drawer-tag {
  font-size: 10px;
  font-weight: 800;
  font-family: 'DM Sans', sans-serif;
  padding: 2px 8px;
  border-radius: 4px;
  letter-spacing: 0.5px;
  background: var(--color-danger-bg);
  color: var(--color-danger);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.lock-drawer-title {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

/* ==================== 响应式 ==================== */
@media (max-width: 1200px) {
  .security-overview { grid-template-columns: repeat(2, 1fr); }
  .config-grid, .config-grid--2 { grid-template-columns: 1fr; }
}
@media (max-width: 768px) {
  .security-page { padding: var(--space-3); }
  .security-overview { grid-template-columns: 1fr 1fr; }
  .toggle-grid { grid-template-columns: 1fr; }
}
</style>
