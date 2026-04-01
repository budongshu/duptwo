<template>
  <div class="app-wrapper">
    <!-- ========== 侧边栏 ========== -->
    <aside class="sidebar" :class="{ 'is-collapsed': isCollapsed }">

      <!-- 头部 -->
      <div class="sidebar-header">
        <router-link to="/" class="logo" :title="isCollapsed ? '数据管理平台' : undefined">
          <div class="logo-icon" @click="easterEggs?.handleLogoClick">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <rect x="1" y="1" width="6" height="6" rx="1.5" fill="currentColor" opacity="0.9"/>
              <rect x="9" y="1" width="6" height="6" rx="1.5" fill="currentColor" opacity="0.6"/>
              <rect x="1" y="9" width="6" height="6" rx="1.5" fill="currentColor" opacity="0.6"/>
              <rect x="9" y="9" width="6" height="6" rx="1.5" fill="currentColor" opacity="0.3"/>
            </svg>
          </div>
          <span v-if="!isCollapsed" class="logo-text">数据管理平台</span>
        </router-link>
      </div>

      <!-- 折叠按钮（浮动在底部导航右侧） -->
      <button class="collapse-btn" @click="toggleCollapse" :title="isCollapsed ? '展开侧边栏' : '收起侧边栏'">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline v-if="isCollapsed" points="9 6 15 12 9 18" />
          <polyline v-else points="15 18 9 12 15 6" />
        </svg>
      </button>

      <!-- 导航区域 -->
      <nav class="sidebar-nav" role="navigation">

        <!-- 分组导航 -->
        <div v-for="group in menuGroups" :key="group.key" class="nav-group">
          <template v-if="getVisibleItems(group.items).length > 0">
            <!-- 分组标题 -->
            <div v-if="!isCollapsed" class="nav-group-header" @click="toggleGroup(group.key)">
              <span class="nav-group-title">{{ t(group.labelKey) }}</span>
              <svg class="nav-group-chevron" :class="{ 'is-up': collapsedGroups.has(group.key) }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="6 15 12 9 18 15" />
              </svg>
            </div>
            <div v-else class="nav-group-header nav-group-header--icon">
              <span class="nav-group-line"></span>
            </div>

            <!-- 菜单项 -->
            <div v-if="!isCollapsed && !collapsedGroups.has(group.key)" class="nav-group-body">
              <template v-for="item in getVisibleItems(group.items)" :key="item.path || item.nameKey">

                <!-- 顶级父菜单（含子项） -->
                <div v-if="item.children" class="nav-parent">
                  <div
                    class="nav-item nav-item--parent"
                    :class="{ 'is-active': isChildActive(item) }"
                    @click="toggleGroup('sub_' + item.nameKey)"
                  >
                    <div class="nav-item__icon-wrap">
                      <component :is="item.icon" class="nav-icon" />
                    </div>
                    <span class="nav-item__label">{{ t(item.nameKey) }}</span>
                    <svg class="nav-parent-chevron" :class="{ 'is-open': !collapsedGroups.has('sub_' + item.nameKey) }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <polyline points="6 9 12 15 18 9" />
                    </svg>
                  </div>
                  <!-- 子菜单 -->
                  <div v-if="!collapsedGroups.has('sub_' + item.nameKey)" class="nav-children">
                    <router-link
                      v-for="child in item.children.filter(c => hasPermission(c.perm))"
                      :key="child.path"
                      :to="child.path!"
                      class="nav-item nav-item--child"
                      :class="{ 'is-active': isActive(child.path!) }"
                    >
                      <div class="nav-item__icon-wrap">
                        <component :is="child.icon" class="nav-icon" />
                      </div>
                      <span class="nav-item__label">{{ t(child.nameKey!) }}</span>
                      <span v-if="isActive(child.path!)" class="nav-item__active-dot"></span>
                    </router-link>
                  </div>
                </div>

                <!-- 普通菜单项 -->
                <router-link
                  v-else
                  :to="item.path!"
                  class="nav-item"
                  :class="{ 'is-active': isActive(item.path!) }"
                >
                  <div class="nav-item__icon-wrap">
                    <component :is="item.icon" class="nav-icon" />
                  </div>
                  <span class="nav-item__label">{{ t(item.nameKey) }}</span>
                  <span v-if="item.isNew" class="nav-tag-new">NEW</span>
                  <span v-else-if="item.badge" class="nav-badge">{{ item.badge }}</span>
                  <span v-if="isActive(item.path!)" class="nav-item__active-dot"></span>
                </router-link>

              </template>
            </div>

            <!-- 收起时只显示图标列表 -->
            <div v-if="isCollapsed" class="nav-group-body nav-group-body--collapsed">
              <router-link
                v-for="item in getVisibleItems(group.items)"
                :key="item.path"
                :to="item.path"
                class="nav-item nav-item--icon-only"
                :class="{ 'is-active': isActive(item.path) }"
                :title="t(item.nameKey) + (item.isNew ? ' (NEW)' : '')"
              >
                <component :is="item.icon" class="nav-icon" />
                <span v-if="item.isNew" class="nav-tag-new nav-tag-new--dot"></span>
              </router-link>
            </div>
          </template>
        </div>
      </nav>

      <!-- 底部用户区 -->
      <div class="sidebar-footer">
        <div class="footer-divider"></div>

        <!-- 服务健康状态 -->
        <div class="health-indicator" :class="healthStatus" :title="healthTooltip" @click="checkHealth">
          <span class="health-dot"></span>
          <span v-if="!isCollapsed" class="health-label">{{ healthLabel }}</span>
        </div>

        <div class="user-info" @click="showUserMenu = !showUserMenu" :title="isCollapsed ? (currentUser?.nickname || currentUser?.username) : undefined">
          <div class="user-avatar-wrap">
            <div class="user-avatar" :class="{ 'is-online': true }">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <span v-if="currentUser?.source === 'AD'" class="source-badge AD">AD</span>
            <span v-else class="source-badge LOCAL">本地</span>
          </div>

          <div v-if="!isCollapsed" class="user-details">
            <span class="user-name">{{ currentUser?.nickname || currentUser?.username || '未登录' }}</span>
            <span class="user-role">{{ currentUser?.roleName || '普通用户' }}</span>
          </div>

          <svg v-if="!isCollapsed" class="user-chevron" :class="{ 'is-open': showUserMenu }" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9"/>
          </svg>
        </div>

        <!-- 用户下拉菜单 -->
        <transition name="dropdown-anim">
          <div v-if="showUserMenu && !isCollapsed" class="user-dropdown">
            <!-- 深色模式 + 语言切换 -->
            <div class="dropdown-theme-row">
              <div class="dropdown-mode-btn" :class="{ active: isDark }" @click="toggleDark" :title="isDark ? t('theme.switchLight') : t('theme.switchDark')">
                <!-- Sun icon (light mode active → click to go dark) -->
                <svg v-if="isDark" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
                </svg>
                <!-- Moon icon (dark mode active → click to go light) -->
                <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
                </svg>
              </div>
              <div class="dropdown-lang-btns">
                <span class="dropdown-lang-btn" :class="{ active: locale === 'zh' }" @click="switchLocale('zh')">中</span>
                <span class="dropdown-lang-btn" :class="{ active: locale === 'en' }" @click="switchLocale('en')">EN</span>
              </div>
            </div>
            <div class="dropdown-divider"></div>
            <div class="dropdown-item" @click="router.push('/profile'); showUserMenu = false">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              {{ t('nav.profile') }}
            </div>
            <div class="dropdown-divider"></div>
            <div class="dropdown-item dropdown-item--danger" @click="handleLogout">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
              {{ t('nav.logout') }}
            </div>
          </div>
        </transition>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main-container" :class="{ 'sidebar-collapsed': isCollapsed }">
      <router-view v-slot="{ Component }">
        <transition name="page-transition" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, h, watch, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import type { UserInfo } from '@/api/auth'

const route = useRoute()
const router = useRouter()
const { locale, t } = useI18n()
const easterEggs = inject<{ handleLogoClick: () => void }>('easterEggs')
const isCollapsed = ref(false)
const showUserMenu = ref(false)
const currentUser = ref<UserInfo | null>(null)
const collapsedGroups = ref(new Set<string>(['audit']))
const isDark = ref(localStorage.getItem('isDark') === 'true')

const toggleDark = () => {
  isDark.value = !isDark.value
  localStorage.setItem('isDark', String(isDark.value))
  if (isDark.value) {
    document.documentElement.classList.add('el-theme-dark')
  } else {
    document.documentElement.classList.remove('el-theme-dark')
  }
}

const switchLocale = (l: 'zh' | 'en') => {
  locale.value = l
  localStorage.setItem('locale', l)
}

// 健康状态
const healthStatus = ref<'healthy' | 'unhealthy' | 'checking'>('checking')
const healthLabel = ref('检查中')
const healthTooltip = ref('正在检查服务状态...')
let healthTimer: ReturnType<typeof setInterval> | null = null

const checkHealth = async () => {
  healthStatus.value = 'checking'
  healthLabel.value = '检查中'
  healthTooltip.value = '正在检查服务状态...'
  const token = localStorage.getItem('token') || ''
  const baseUrl = (import.meta.env.VITE_API_URL as string || '').replace('/api', '')
  const start = Date.now()
  try {
    const res = await fetch(`${baseUrl}/health`, {
      credentials: 'include',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const ms = Date.now() - start
    if (res.ok) {
      healthStatus.value = 'healthy'
      healthLabel.value = `正常 ${ms}ms`
      healthTooltip.value = `服务正常 · 响应 ${ms}ms`
    } else {
      healthStatus.value = 'unhealthy'
      healthLabel.value = '异常'
      healthTooltip.value = `服务异常 · HTTP ${res.status}`
    }
  } catch {
    healthStatus.value = 'unhealthy'
    healthLabel.value = '离线'
    healthTooltip.value = '无法连接到服务'
  }
}

// ========== 图标定义（Lucide 风格，唯一语义化）==========
const iconProps = { width: 16, height: 16, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }

// 数据概览 - 仪表盘/图表
const DashboardIcon = () => h('svg', iconProps, [
  h('rect', { x: 3, y: 3, width: 7, height: 7, rx: 1 }),
  h('rect', { x: 14, y: 3, width: 7, height: 7, rx: 1 }),
  h('rect', { x: 14, y: 14, width: 7, height: 7, rx: 1 }),
  h('rect', { x: 3, y: 14, width: 7, height: 7, rx: 1 })
])

// 上传记录 - 列表/文档
const UploadRecordIcon = () => h('svg', iconProps, [
  h('path', { d: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z' }),
  h('polyline', { points: '14 2 14 8 20 8' }),
  h('line', { x1: 16, y1: 13, x2: 8, y2: 13 }),
  h('line', { x1: 16, y1: 17, x2: 8, y2: 17 })
])

// 项目列表 - 文件夹/目录树
const ProjectIcon = () => h('svg', iconProps, [
  h('path', { d: 'M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z' })
])

// 项目管理 - 格子/看板
const ProjectManageIcon = () => h('svg', iconProps, [
  h('rect', { x: 3, y: 3, width: 7, height: 9, rx: 1 }),
  h('rect', { x: 14, y: 3, width: 7, height: 5, rx: 1 }),
  h('rect', { x: 14, y: 12, width: 7, height: 9, rx: 1 }),
  h('rect', { x: 3, y: 16, width: 7, height: 5, rx: 1 })
])

// 人员管理 - 用户id卡
const PersonnelIcon = () => h('svg', iconProps, [
  h('rect', { x: 2, y: 5, width: 20, height: 14, rx: 2 }),
  h('circle', { cx: 12, cy: 10, r: 2 }),
  h('path', { d: 'M6 18h.01M18 18h.01' }),
  h('line', { x1: 2, y1: 10, x2: 4, y2: 10 })
])

// 用户管理 - 单用户
const UserIcon = () => h('svg', iconProps, [
  h('path', { d: 'M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2' }),
  h('circle', { cx: 12, cy: 7, r: 4 })
])

// 角色管理 - 盾牌
const RoleIcon = () => h('svg', iconProps, [
  h('path', { d: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z' })
])

// 用户组 - 多个用户/人堆
const GroupIcon = () => h('svg', iconProps, [
  h('circle', { cx: 9, cy: 7, r: 4 }),
  h('path', { d: 'M3 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2' }),
  h('circle', { cx: 17, cy: 8, r: 3 }),
  h('path', { d: 'M21 21v-1.5a3 3 0 0 0-2-2.83' })
])

// 字段配置 - 滑动条/设置
const FieldConfigIcon = () => h('svg', iconProps, [
  h('line', { x1: 4, y1: 21, x2: 4, y2: 14 }),
  h('line', { x1: 4, y1: 10, x2: 4, y2: 3 }),
  h('line', { x1: 12, y1: 21, x2: 12, y2: 12 }),
  h('line', { x1: 12, y1: 8, x2: 12, y2: 3 }),
  h('line', { x1: 20, y1: 21, x2: 20, y2: 16 }),
  h('line', { x1: 20, y1: 12, x2: 20, y2: 3 }),
  h('line', { x1: 1, y1: 14, x2: 7, y2: 14 }),
  h('line', { x1: 9, y1: 8, x2: 15, y2: 8 }),
  h('line', { x1: 17, y1: 16, x2: 23, y2: 16 })
])

// AD域配置 - 服务器/网络
const ADAuthorityIcon = () => h('svg', iconProps, [
  h('rect', { x: 2, y: 2, width: 20, height: 8, rx: 2 }),
  h('rect', { x: 2, y: 14, width: 20, height: 8, rx: 2 }),
  h('line', { x1: 6, y1: 6, x2: 6.01, y2: 6 }),
  h('line', { x1: 6, y1: 18, x2: 6.01, y2: 18 })
])

// 安全设置 - 锁+钥匙
const SecurityIcon = () => h('svg', iconProps, [
  h('rect', { x: 3, y: 11, width: 18, height: 11, rx: 2 }),
  h('path', { d: 'M7 11V7a5 5 0 0 1 10 0v4' }),
  h('circle', { cx: 12, cy: 16, r: 1 })
])

// 操作日志 - 文档+编辑
const OperationLogIcon = () => h('svg', iconProps, [
  h('path', { d: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z' }),
  h('polyline', { points: '14 2 14 8 20 8' }),
  h('line', { x1: 16, y1: 13, x2: 8, y2: 13 }),
  h('line', { x1: 16, y1: 17, x2: 8, y2: 17 }),
  h('polyline', { points: '10 9 9 9 8 9' })
])

// 登录日志 - 进出箭头
const LoginLogIcon = () => h('svg', iconProps, [
  h('path', { d: 'M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4' }),
  h('polyline', { points: '10 17 15 12 10 7' }),
  h('line', { x1: 15, y1: 12, x2: 3, y2: 12 })
])

// ========== 菜单配置（使用 i18n key）==========
// 支持嵌套子菜单的类型
interface NavItem {
  path?: string
  name: string
  icon: () => ReturnType<typeof h>
  perm?: string
  isNew?: boolean
  badge?: string
  children?: NavItem[]
}

// i18n key 常量
const menuGroups = reactive([
  {
    key: 'data',
    labelKey: 'sidebar.dataCenter',
    items: [
      { path: '/upload-record', nameKey: 'nav.dashboard', icon: DashboardIcon, perm: 'upload:read' },
      { path: '/upload-record/list', nameKey: 'nav.uploadRecords', icon: UploadRecordIcon, perm: 'upload:read' }
    ]
  },
  {
    key: 'center',
    labelKey: 'sidebar.projectCenter',
    items: [
      {
        nameKey: 'nav.projectList', icon: ProjectIcon, perm: 'project:read',
        children: [
          { path: '/projects', nameKey: 'nav.projectManage', icon: ProjectManageIcon, perm: 'project:read' },
        ]
      },
      { path: '/personnel', nameKey: 'nav.personnel', icon: PersonnelIcon, perm: 'personnel:read' },
    ]
  },
  {
    key: 'person',
    labelKey: 'sidebar.orgPersonnel',
    items: [
      { path: '/users', nameKey: 'nav.userManage', icon: UserIcon, perm: 'user:read' },
      { path: '/roles', nameKey: 'nav.roleManage', icon: RoleIcon, perm: 'role:read' },
      { path: '/user-groups', nameKey: 'nav.userGroups', icon: GroupIcon, perm: 'role:read' }
    ]
  },
  {
    key: 'system',
    labelKey: 'sidebar.system',
    items: [
      { path: '/field-config', nameKey: 'nav.fieldConfig', icon: FieldConfigIcon, perm: 'field-config:read' },
      { path: '/system/ad-settings', nameKey: 'nav.adSettings', icon: ADAuthorityIcon, perm: 'config:read' },
      { path: '/system/security-settings', nameKey: 'nav.securitySettings', icon: SecurityIcon, perm: 'config:read' }
    ]
  },
  {
    key: 'audit',
    labelKey: 'sidebar.audit',
    items: [
      { path: '/audit/operation-log', nameKey: 'nav.operationLog', icon: OperationLogIcon, perm: 'audit:operation:read' },
      { path: '/audit/login-log', nameKey: 'nav.loginLog', icon: LoginLogIcon, perm: 'audit:login:read' }
    ]
  }
])

// 兼容类型：nameKey 用于 i18n，name 用于普通文本
interface NavItem {
  path?: string
  name?: string
  nameKey?: string
  icon: () => ReturnType<typeof h>
  perm?: string
  isNew?: boolean
  badge?: string
  children?: NavItem[]
}

const hasPermission = (perm: string | undefined) => {
  if (!perm) return true
  if (!currentUser.value?.permissions || currentUser.value.permissions.length === 0) return true
  if (currentUser.value.permissions.includes('admin:all')) return true
  return currentUser.value.permissions.includes(perm)
}

// 扁平化所有可显示的菜单项（含嵌套）
const flattenItems = (items: NavItem[]): NavItem[] => {
  const result: NavItem[] = []
  for (const item of items) {
    if (item.path) result.push(item)
    if (item.children) {
      for (const child of item.children) {
        if (child.path) result.push(child)
      }
    }
  }
  return result
}

const getVisibleItems = (items: NavItem[]) => {
  return flattenItems(items).filter(item => hasPermission(item.perm))
}

// 检查是否有可见子项
const hasVisibleChildren = (item: NavItem): boolean => {
  if (!item.children) return false
  return flattenItems([item]).some(i => hasPermission(i.perm))
}

// 判断任意子项是否激活
const isChildActive = (item: NavItem): boolean => {
  if (!item.children) return false
  return item.children.some(child => child.path && isActive(child.path))
}

const toggleGroup = (key: string) => {
  if (collapsedGroups.value.has(key)) {
    collapsedGroups.value.delete(key)
  } else {
    collapsedGroups.value.add(key)
  }
}

const isActive = (path: string) => route.path === path

const loadUser = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try { currentUser.value = JSON.parse(userStr) } catch { currentUser.value = null }
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '退出确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    showUserMenu.value = false
    router.push('/login')
  } catch {}
}

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  if (isCollapsed.value) showUserMenu.value = false
}

const handleClickOutside = (e: MouseEvent) => {
  if (!(e.target as HTMLElement).closest('.sidebar-footer')) {
    showUserMenu.value = false
  }
}

watch(showUserMenu, (val) => {
  if (val) { document.addEventListener('click', handleClickOutside) }
  else { document.removeEventListener('click', handleClickOutside) }
})

onMounted(() => {
  loadUser()
  checkHealth()
  healthTimer = setInterval(checkHealth, 30000)
  // 初始化深色模式
  if (localStorage.getItem('isDark') === 'true') {
    document.documentElement.classList.add('el-theme-dark')
    isDark.value = true
  }
  // 初始化语言
  const savedLocale = localStorage.getItem('locale') as 'zh' | 'en'
  if (savedLocale) locale.value = savedLocale
})

onUnmounted(() => {
  if (healthTimer) clearInterval(healthTimer)
})
</script>

<style scoped lang="scss">
/* ==================== 基础布局 ==================== */
.app-wrapper {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background-color: var(--theme-main-bg);
}

/* ==================== 侧边栏 ==================== */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 210px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  background-color: var(--theme-sidebar-bg-start);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.12);
  border-right: 1px solid var(--theme-sidebar-border);

  &.is-collapsed {
    width: 56px;
  }
}

/* ==================== 头部 ==================== */
.sidebar-header {
  display: flex;
  align-items: center;
  padding: 0 14px;
  height: 52px;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  flex: 0 0 auto;
}

/* ==================== 导航 ==================== */
.logo {
  display: flex;
  align-items: center;
  gap: 9px;
  text-decoration: none;
  overflow: hidden;
  flex: 1;
  min-width: 0;
}

.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 10px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  position: relative;
  z-index: 1;

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.12);
    border-radius: 2px;
    &:hover { background: rgba(255, 255, 255, 0.2); }
  }
}

.logo-icon {
  width: 30px;
  height: 30px;
  border-radius: 7px;
  background: var(--theme-sidebar-logo-icon-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--theme-sidebar-logo-icon-color);
  flex-shrink: 0;
}

.logo-text {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 800;
  color: var(--theme-logo-text-color, #ffffff);
  white-space: nowrap;
  letter-spacing: -0.3px;
  line-height: 1.2;
  transition: opacity 0.2s ease;
}

/* 收起/展开按钮 — 浮动于侧边栏右下角 */
.collapse-btn {
  position: absolute;
  bottom: 90px;
  right: 10px;
  z-index: 10;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 8px;
  background: var(--theme-sidebar-collapse-btn-bg);
  color: var(--theme-sidebar-collapse-btn-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);

  &:hover {
    background: var(--theme-sidebar-collapse-btn-hover-bg);
    color: var(--theme-sidebar-collapse-btn-hover-color);
    transform: scale(1.08);
  }
}

.sidebar.is-collapsed .collapse-btn {
  bottom: 86px;
  right: 14px;
}

/* ==================== 导航 ==================== */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 10px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  position: relative;
  z-index: 1;

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.12);
    border-radius: 2px;
    &:hover { background: rgba(255, 255, 255, 0.2); }
  }
}

/* 首页入口 */
.nav-home {
  margin-bottom: 6px;
}

/* 分组 */
.nav-group {
  margin-bottom: 4px;
}

.nav-group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 10px 5px;
  cursor: pointer;
  user-select: none;
  border-radius: 6px;
  transition: all 0.15s ease;

  &:hover {
    .nav-group-title { color: var(--theme-sidebar-group-title-hover); }
  }

  &--icon {
    justify-content: center;
    padding: 10px 0 5px;
    cursor: default;
  }
}

.nav-group-line {
  display: block;
  width: 24px;
  height: 1px;
  background: var(--theme-sidebar-border);
  margin: 0 auto;
}

.nav-group-title {
  font-family: 'DM Sans', 'Manrope', sans-serif;
  font-size: 12px;
  font-weight: 700;
  color: var(--theme-sidebar-group-title);
  letter-spacing: 0.3px;
  transition: color 0.15s;
}

.nav-group-chevron {
  color: var(--theme-sidebar-chevron);
  transition: transform 0.2s ease;
  flex-shrink: 0;

  &.is-up { transform: rotate(-90deg); }
}

.nav-group-body {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 2px 0;

  &--collapsed {
    gap: 2px;
    padding: 2px 0;
  }
}

/* 菜单项 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 8px;
  color: var(--theme-sidebar-text);
  text-decoration: none;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 450;
  transition: all 0.18s ease;
  position: relative;
  cursor: pointer;
  overflow: hidden;
  white-space: nowrap;

  /* 入场动画延迟 */
  animation: nav-slide-in 0.3s ease both;
  animation-delay: calc(var(--stagger, 0) * 40ms);

  @keyframes nav-slide-in {
    from { opacity: 0; transform: translateX(-8px); }
    to { opacity: 1; transform: translateX(0); }
  }

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%) scaleY(0);
    width: 3px;
    height: 60%;
    background: var(--theme-sidebar-active-color);
    border-radius: 0 2px 2px 0;
    transition: transform 0.2s ease;
  }

  &:hover {
    background: var(--theme-sidebar-hover-bg);
    color: var(--theme-sidebar-text-hover);

    .nav-icon { color: var(--theme-sidebar-text-hover); }
  }

  &.is-active {
    background: var(--theme-sidebar-active-bg);
    color: var(--theme-sidebar-active-color);
    font-weight: 600;

    &::before { transform: translateY(-50%) scaleY(1); }

    .nav-icon { color: var(--theme-sidebar-active-icon-color); }
    .nav-item__label { color: var(--theme-sidebar-active-color); }
  }

  /* 收起时的图标模式 */
  &--icon-only {
    justify-content: center;
    padding: 9px;

    .nav-icon { color: var(--theme-sidebar-nav-icon); }

    &:hover {
      background: var(--theme-sidebar-hover-bg);
      .nav-icon { color: var(--theme-sidebar-text-hover); }
    }

    &.is-active {
      background: var(--theme-sidebar-active-bg);
      .nav-icon { color: var(--theme-sidebar-active-icon-color); }
    }
  }

  /* 首页特殊样式 */
  &--home {
    margin-bottom: 4px;
    background: var(--theme-sidebar-section-bg);

    &:hover { background: var(--theme-sidebar-hover-bg); }
    &.is-active {
      background: var(--theme-sidebar-active-bg);
      .nav-icon { color: var(--theme-sidebar-active-icon-color); }
    }
  }

  /* 父级菜单（可折叠） */
  &--parent {
    font-weight: 500;
    color: var(--theme-sidebar-text);

    &:hover {
      background: var(--theme-sidebar-hover-bg);
      color: var(--theme-sidebar-text-hover);
      .nav-icon { color: var(--theme-sidebar-text-hover); }
    }

    &.is-active {
      background: var(--theme-sidebar-active-bg);
      color: var(--theme-sidebar-active-color);
      font-weight: 600;
      .nav-icon { color: var(--theme-sidebar-active-icon-color); }
      .nav-item__label { color: var(--theme-sidebar-active-color); }
      &::before { transform: translateY(-50%) scaleY(1); }
    }
  }

  /* 子菜单项 */
  &--child {
    padding: 6px 10px 6px 38px;
    font-size: 12.5px;

    &::before { display: none; }

    &:hover {
      background: var(--theme-sidebar-hover-bg);
      color: var(--theme-sidebar-text-hover);
      .nav-icon { color: var(--theme-sidebar-text-hover); }
    }

    &.is-active {
      background: var(--theme-sidebar-active-bg);
      color: var(--theme-sidebar-active-color);
      font-weight: 600;
      .nav-icon { color: var(--theme-sidebar-active-icon-color); }
      .nav-item__label { color: var(--theme-sidebar-active-color); }
    }
  }
}

.nav-item__icon-wrap {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon {
  width: 16px;
  height: 16px;
  color: var(--theme-sidebar-nav-icon);
  transition: color 0.15s;
  flex-shrink: 0;
}

.nav-item__label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  color: inherit;
  transition: color 0.15s;
}

/* 父级菜单指示器 */
.nav-parent-chevron {
  color: var(--theme-sidebar-chevron);
  flex-shrink: 0;
  transition: transform 0.2s ease;
  &.is-open { transform: rotate(180deg); }
}

/* 子菜单容器 */
.nav-children {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 2px 0;
  border-left: 1px solid var(--theme-sidebar-border);
  margin-left: 22px;
  margin-top: 2px;
  margin-bottom: 2px;
  padding-left: 8px;
}

/* 活动指示点 */
.nav-item__active-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--theme-sidebar-active-color);
  flex-shrink: 0;
  box-shadow: 0 0 6px var(--theme-sidebar-active-color);
}

/* 标签 */
.nav-tag-new {
  font-family: 'DM Sans', sans-serif;
  font-size: 9px;
  font-weight: 800;
  background: var(--theme-new-tag-bg);
  color: #fff;
  padding: 1px 5px;
  border-radius: 3px;
  letter-spacing: 0.5px;
  flex-shrink: 0;
  box-shadow: 0 1px 4px rgba(239, 68, 68, 0.4);

  &--dot {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 6px;
    height: 6px;
    padding: 0;
    border-radius: 50%;
  }
}

.nav-badge {
  font-size: 10px;
  font-weight: 700;
  background: var(--theme-sidebar-badge-bg);
  color: var(--theme-sidebar-text-hover);
  padding: 1px 6px;
  border-radius: 10px;
  flex-shrink: 0;
}

/* ==================== 底部用户区 ==================== */
.sidebar-footer {
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  padding: 0 8px 12px;
}

.footer-divider {
  height: 1px;
  background: var(--theme-sidebar-border);
  margin: 0 6px 10px;
}

/* 服务健康指示器 */
.health-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 10px;
  margin: 0 0 4px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s ease;
  user-select: none;

  &:hover { background: var(--theme-sidebar-hover-bg); }
}

.health-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);

  .healthy & {
    background: #22c55e;
    box-shadow: 0 0 6px rgba(34, 197, 94, 0.7);
  }
  .unhealthy & {
    background: #ef4444;
    box-shadow: 0 0 6px rgba(239, 68, 68, 0.7);
  }
  .checking & {
    background: #f59e0b;
    animation: health-pulse 1s ease-in-out infinite;
  }
}

@keyframes health-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.85); }
}

.health-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 500;
  color: var(--theme-sidebar-text);
  opacity: 0.7;
  white-space: nowrap;
  transition: color 0.2s ease;

  .healthy & { color: #86efac; opacity: 1; }
  .unhealthy & { color: #fca5a5; opacity: 1; }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;

  &:hover { background: var(--theme-sidebar-hover-bg); }
}

.user-avatar-wrap {
  position: relative;
  flex-shrink: 0;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--theme-sidebar-user-avatar-bg);
  border: 1.5px solid var(--theme-sidebar-user-avatar-border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: all 0.2s ease;
  position: relative;

  /* 在线状态点 */
  &.is-online::after {
    content: '';
    position: absolute;
    bottom: -1px;
    right: -1px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #22c55e;
    border: 1.5px solid var(--theme-sidebar-bg-end);
  }

  .is-collapsed & {
    width: 36px;
    height: 36px;
  }
}

.source-badge {
  position: absolute;
  bottom: -2px;
  right: -5px;
  font-size: 7.5px;
  font-weight: 800;
  padding: 1px 3px;
  border-radius: 3px;
  letter-spacing: 0.3px;
  border: 1px solid rgba(0,0,0,0.25);
  font-family: 'DM Sans', monospace;

  &.AD {
    background: var(--theme-badge-ad-bg);
    color: white;
  }
  &.LOCAL {
    background: var(--theme-badge-local-bg);
    color: white;
  }
}

.user-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.user-name {
  font-family: 'DM Sans', 'Manrope', sans-serif;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--theme-sidebar-user-name);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

.user-role {
  font-size: 10px;
  color: var(--theme-sidebar-user-role);
  margin-top: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-chevron {
  color: var(--theme-sidebar-chevron);
  flex-shrink: 0;
  transition: transform 0.2s ease;
  &.is-open { transform: rotate(180deg); }
}

/* 用户下拉菜单 */
.user-dropdown {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 6px;
  right: 6px;
  background: var(--theme-sidebar-dropdown-bg);
  border: 1px solid var(--theme-sidebar-dropdown-border);
  border-radius: 10px;
  padding: 4px;
  box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(12px);
  z-index: 200;
}

.dropdown-theme-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px 8px;
}

.dropdown-mode-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--theme-sidebar-dropdown-item);
  transition: all 0.15s ease;
  background: transparent;

  &:hover {
    background: var(--theme-sidebar-hover-bg);
    color: var(--theme-sidebar-dropdown-item-hover);
  }

  &.active {
    background: rgba(245, 158, 11, 0.15);
    color: #fbbf24;
  }
}

.dropdown-lang-btns {
  display: flex;
  gap: 2px;
  background: var(--theme-sidebar-hover-bg);
  border-radius: 6px;
  padding: 2px;
}

.dropdown-lang-btn {
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  font-family: 'DM Sans', sans-serif;
  cursor: pointer;
  color: var(--theme-sidebar-dropdown-item);
  opacity: 0.5;
  transition: all 0.15s ease;

  &:hover { opacity: 0.8; }
  &.active {
    background: var(--color-primary);
    color: #fff;
    opacity: 1;
  }
}

.dropdown-header {
  font-family: 'DM Sans', sans-serif;
  font-size: 10px;
  font-weight: 700;
  color: var(--theme-sidebar-dropdown-item);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 4px 10px 4px;
  opacity: 0.6;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 6px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12.5px;
  color: var(--theme-sidebar-dropdown-item);
  cursor: pointer;
  transition: all 0.15s;

  svg { color: var(--theme-sidebar-dropdown-svg); flex-shrink: 0; }

  &:hover {
    background: var(--theme-sidebar-hover-bg);
    color: var(--theme-sidebar-dropdown-item-hover);
    svg { color: var(--theme-sidebar-dropdown-item-hover); }
  }

  &--danger:hover {
    background: var(--theme-sidebar-danger-hover-bg);
    color: var(--theme-sidebar-danger-hover-color);
    svg { color: var(--theme-sidebar-danger-hover-svg); }
  }
}

.dropdown-divider {
  height: 1px;
  background: var(--theme-sidebar-dropdown-border);
  margin: 4px 6px;
}

/* 下拉动画 */
.dropdown-anim-enter-active { transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1); }
.dropdown-anim-leave-active { transition: all 0.15s ease; }
.dropdown-anim-enter-from { opacity: 0; transform: translateY(8px) scale(0.96); }
.dropdown-anim-leave-to { opacity: 0; transform: translateY(4px); }

/* ==================== 主内容区 ==================== */
.main-container {
  flex: 1;
  margin-left: 210px;
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1), background-color 0.3s ease;
  background: var(--theme-main-bg);

  &::-webkit-scrollbar { width: 5px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb {
    background: var(--color-scrollbar);
    border-radius: 3px;
    &:hover { background: var(--color-scrollbar-hover); }
  }

  &.sidebar-collapsed { margin-left: 56px; }
}

/* 页面切换动画 — 纯 opacity 切换，无 transform 抖动 */
.page-transition-enter-active {
  transition: opacity 0.15s ease;
}
.page-transition-leave-active {
  transition: opacity 0.1s ease;
}
.page-transition-enter-from,
.page-transition-leave-to {
  opacity: 0;
}

/* ==================== 收起状态下的特殊处理 ==================== */
.sidebar.is-collapsed {
  .logo { justify-content: center; gap: 0; }
  .logo-text { opacity: 0; width: 0; }
  .user-details { display: none; }
  .user-chevron { display: none; }
  .user-info { justify-content: center; padding: 8px; }
  .user-avatar-wrap { margin: 0; }
  .nav-group-header { padding: 8px 0 4px; }
  .nav-group-title { display: none; }
  .nav-group-chevron { display: none; }
}
</style>
