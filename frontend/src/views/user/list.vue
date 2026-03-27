<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">用户管理</h1>
        <p class="page-subtitle">管理系统用户账号</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          新增用户
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <el-input v-model="keyword" placeholder="搜索用户名/昵称" clearable @keyup.enter="handleSearch" />
      <el-select v-model="status" placeholder="状态" clearable style="width: 120px" @change="handleSearch">
        <el-option label="正常" value="active" />
        <el-option label="禁用" value="inactive" />
      </el-select>
      <el-select v-model="filterRoleId" placeholder="角色" clearable style="width: 140px" @change="handleSearch">
        <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
      </el-select>
      <el-select v-model="filterGroupId" placeholder="用户组" clearable style="width: 140px" @change="handleSearch">
        <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <!-- 表格工具栏 -->
      <div class="table-toolbar">
        <div class="toolbar-left">
          <span class="record-count">共 <strong>{{ pagination.total }}</strong> 条</span>
          <span v-if="selectedRows.length > 0" class="selection-count">
            已选 <strong>{{ selectedRows.length }}</strong> 项
          </span>
        </div>
        <div class="toolbar-right">
          <el-button
            v-if="selectedRows.length > 0"
            type="danger"
            @click="handleBatchDelete"
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            批量删除
          </el-button>
        </div>
      </div>

      <el-table ref="tableRef" v-model:selection="selectedRows" :data="tableData" v-loading="loading" stripe>
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="100" />
        <el-table-column prop="email" label="邮箱" min-width="160" />
        <el-table-column prop="roleName" label="角色" min-width="100" align="center">
          <template #default="{ row }">
            <span class="tag-badge">{{ row.roleName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="groupName" label="用户组" min-width="130" align="center">
          <template #default="{ row }">
            <template v-if="row.groupName">
              <button class="group-badge" :title="`查看「${row.groupName}」详情`" @click.stop="navigateToGroup(row.groupId, row.groupName)">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                {{ row.groupName }}
                <span class="group-count">{{ groupMemberCount(row.groupId) }}</span>
              </button>
            </template>
            <span v-else class="no-group">未分配</span>
          </template>
        </el-table-column>
        <el-table-column prop="mfaEnabled" label="MFA" min-width="80" align="center">
          <template #default="{ row }">
            <span v-if="row.mfaEnabled" class="mfa-badge mfa-badge--enabled">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              已启用
            </span>
            <span v-else class="mfa-badge mfa-badge--disabled">未启用</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="80" align="center">
          <template #default="{ row }">
            <span :class="['status-badge', row.status === 'active' ? 'status-badge--success' : 'status-badge--danger']">
              {{ row.statusText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" min-width="80" align="center">
          <template #default="{ row }">
            <span v-if="row.source === 'AD'" class="source-badge source-badge--ad">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              AD
            </span>
            <span v-else class="source-badge source-badge--local">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              本地
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="lastLoginAt" label="最后登录" min-width="140">
          <template #default="{ row }">
            {{ formatDate(row.lastLoginAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="90" align="center">
          <template #default="{ row }">
            <div class="row-actions">
              <button class="row-action-btn row-action-btn--edit" title="编辑" @click="handleEdit(row)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="row-action-btn row-action-btn--delete" title="删除" @click="handleDelete(row)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
                </svg>
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
          @current-change="loadData"
          @size-change="loadData"
        />
      </div>
    </div>

    <!-- 用户编辑弹窗 -->
    <el-dialog v-model="drawerVisible" width="600px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="dialog-title-text">{{ isEdit ? '编辑用户' : '新增用户' }}</span>
        </div>
      </template>
      <div class="dialog-body-inner">
        <!-- AD用户提示 -->
        <div v-if="isEdit && form.source === 'AD'" class="ad-user-notice">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          <div class="ad-notice-text">
            <span class="ad-notice-title">AD域用户</span>
            <span class="ad-notice-desc">用户名、邮箱由AD域同步，密码由AD域控制器管理，仅可编辑本地扩展信息</span>
          </div>
        </div>

        <!-- 基本设置 -->
        <div class="info-bar">
          <div class="info-bar-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
          </div>
          <div class="info-bar-text">
            <span class="info-bar-title">基本信息</span>
            <span class="info-bar-sub">设置用户名、昵称和邮箱</span>
          </div>
        </div>
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <div class="form-row">
            <el-form-item label="用户名" prop="username" v-if="!isEdit" class="is-required-field">
              <el-input v-model="form.username" placeholder="3-32个字符">
                <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
              </el-input>
            </el-form-item>
            <el-form-item v-else label="用户名">
              <el-input :model-value="form.username" disabled>
                <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></template>
              </el-input>
              <div class="field-tag field-tag--ad" v-if="form.source === 'AD'">AD同步</div>
            </el-form-item>
            <el-form-item label="昵称">
              <el-input v-model="form.nickname" :placeholder="isEdit && form.source === 'AD' ? 'AD同步' : '用户显示名称'">
                <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3"/></svg></template>
              </el-input>
            </el-form-item>
          </div>
          <div class="form-row">
            <el-form-item label="邮箱" prop="email" :class="isEdit && form.source === 'AD' ? '' : 'is-required-field'">
              <el-input v-model="form.email" :placeholder="isEdit && form.source === 'AD' ? 'AD同步' : '邮箱地址'" :disabled="isEdit && form.source === 'AD'">
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

        <!-- 安全设置（仅本地用户，编辑时） -->
        <template v-if="isEdit && form.source !== 'AD'">
          <div class="info-bar">
            <div class="info-bar-icon shield">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            </div>
            <div class="info-bar-text">
              <span class="info-bar-title">安全设置</span>
              <span class="info-bar-sub">重置密码或 MFA 认证</span>
            </div>
          </div>
          <div class="security-actions">
            <!-- 密码重置 -->
            <div class="security-card">
              <div class="security-card-icon">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </div>
              <div class="security-card-info">
                <span class="security-card-title">密码重置</span>
                <span class="security-card-desc">生成新密码并发送至用户</span>
              </div>
              <el-button type="warning" size="small" @click="handleResetPassword">重置密码</el-button>
            </div>

            <!-- MFA 重置 -->
            <div class="security-card">
              <div class="security-card-icon" :class="{ 'mfa-on': form.mfaEnabled }">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              </div>
              <div class="security-card-info">
                <span class="security-card-title">MFA 重置</span>
                <span class="security-card-desc">{{ form.mfaEnabled ? '强制解除绑定，下次登录需重新启用' : '当前未启用 MFA' }}</span>
              </div>
              <el-button :type="form.mfaEnabled ? 'warning' : 'default'" size="small" :disabled="!form.mfaEnabled" @click="handleDisableMFA">重置 MFA</el-button>
            </div>
          </div>
        </template>

        <!-- 认证设置 -->
        <div class="info-bar">
          <div class="info-bar-icon shield">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          </div>
          <div class="info-bar-text">
            <span class="info-bar-title">认证配置</span>
            <span class="info-bar-sub">{{ isEdit && form.source === 'AD' ? 'AD用户MFA由域控制器管理' : '分配角色和 MFA 安全设置' }}</span>
          </div>
        </div>
        <el-form label-position="top">
          <div class="form-row">
            <el-form-item label="角色">
              <el-select v-model="form.roleId" placeholder="选择角色" style="width: 100%">
                <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg></template>
                <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="用户组">
              <el-select v-model="form.groupId" placeholder="选择用户组" style="width: 100%">
                <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg></template>
                <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
              </el-select>
            </el-form-item>
          </div>
          <div class="form-row">
            <el-form-item label="MFA 认证">
              <div class="mfa-toggle">
                <span class="mfa-status" :class="form.mfaEnabled ? 'enabled' : 'disabled'">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                  {{ form.mfaEnabled ? '已启用' : '未启用' }}
                </span>
                <span v-if="isEdit && form.source === 'AD'" class="mfa-ad-tag">AD域管理</span>
              </div>
            </el-form-item>
            <el-form-item label="账号状态" v-if="isEdit">
              <div class="status-options">
                <div :class="['status-opt', { active: form.status === 'active' }]" @click="form.status = 'active'">
                  <span class="status-dot status-dot--success"></span>
                  正常
                </div>
                <div :class="['status-opt', { active: form.status === 'inactive' }]" @click="form.status = 'inactive'">
                  <span class="status-dot status-dot--danger"></span>
                  禁用
                </div>
              </div>
            </el-form-item>
          </div>
        </el-form>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="drawerVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">保存</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="resetPwdVisible" width="400px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--warn">密码</span>
          <span class="dialog-title-text">重置密码</span>
        </div>
      </template>
      <el-form ref="resetPwdFormRef" :model="resetPwdForm" :rules="resetPwdRules" label-position="top">
        <el-form-item label="用户">
          <el-input :model-value="resetPwdTarget?.username" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword" class="is-required-field">
          <el-input v-model="resetPwdForm.newPassword" type="password" placeholder="请输入新密码" show-password />
          <div class="password-strength" v-if="resetPwdForm.newPassword">
            <div class="strength-bar">
              <div class="strength-fill" :class="resetPwdStrengthClass" :style="{ width: resetPwdStrengthWidth }"></div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPwdVisible = false">取消</el-button>
        <el-button type="primary" :loading="resetPwdLoading" @click="confirmResetPassword">确认重置</el-button>
      </template>
    </el-dialog>

    <!-- 启用/重置MFA对话框 -->
    <el-dialog v-model="enableMFAVisible" width="420px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--mfa">MFA</span>
          <span class="dialog-title-text">{{ form.mfaEnabled ? '重置MFA' : '启用MFA' }}</span>
        </div>
      </template>
      <div class="enable-mfa-content">
        <p class="enable-mfa-tip">请使用身份验证器（如Google Authenticator）扫描下方二维码，然后输入显示的6位验证码进行验证。</p>
        <div class="qr-code-box" v-if="mfaSecret">
          <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="MFA QR Code" class="qr-code-img" />
          <div class="secret-text">密钥: {{ mfaSecret }}</div>
        </div>
        <el-form ref="enableMFAFormRef" :model="enableMFAForm" :rules="enableMFARules" label-position="top">
          <el-form-item label="验证码" prop="code">
            <el-input v-model="enableMFAForm.code" placeholder="请输入6位验证码" maxlength="6" size="large" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="enableMFAVisible = false">取消</el-button>
        <el-button type="primary" :loading="enableMFALoading" @click="confirmEnableMFA">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserApi, type User, type CreateUserReq, type UpdateUserReq } from '@/api/user'
import { RoleApi, type Role } from '@/api/role'
import { UserGroupApi, type UserGroup } from '@/api/user-group'

const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<User[]>([])
const tableRef = ref()
const selectedRows = ref<User[]>([])
const roles = ref<Role[]>([])
const groups = ref<UserGroup[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 重置密码相关
const resetPwdVisible = ref(false)
const resetPwdTarget = ref<User | null>(null)
const resetPwdLoading = ref(false)
const resetPwdFormRef = ref()
const resetPwdForm = reactive({ newPassword: '' })

// 启用MFA相关
const enableMFAVisible = ref(false)
const enableMFALoading = ref(false)
const enableMFAFormRef = ref()
const mfaSecret = ref('')
const qrCodeUrl = ref('')
const enableMFAForm = reactive({ code: '' })
const enableMFATargetId = ref<number>()

const keyword = ref('')
const status = ref('')
const filterRoleId = ref<number>()
const filterGroupId = ref<number>()

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateUserReq & { id?: number; status?: string; source?: string }>({
  username: '', password: '', nickname: '', email: '', phone: '', roleId: undefined, groupId: undefined, status: 'active', mfaEnabled: false, source: 'LOCAL'
})

// 密码规则检查
const passwordRules = computed(() => {
  const pwd = form.password || ''
  return {
    length: pwd.length >= 6,
    upper: /[A-Z]/.test(pwd),
    lower: /[a-z]/.test(pwd),
    number: /[0-9]/.test(pwd),
    special: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd)
  }
})

// 密码强度计算
const passwordStrength = computed(() => {
  const rules = passwordRules.value
  let score = 0
  if (rules.length) score++
  if (rules.upper) score++
  if (rules.lower) score++
  if (rules.number) score++
  if (rules.special) score++
  return score
})

const strengthWidth = computed(() => {
  return (passwordStrength.value / 5 * 100) + '%'
})

const strengthClass = computed(() => {
  const score = passwordStrength.value
  if (score <= 2) return 'weak'
  if (score <= 3) return 'medium'
  return 'strong'
})

// 重置密码的密码强度
const resetPwdStrength = computed(() => {
  const pwd = resetPwdForm.newPassword || ''
  let score = 0
  if (pwd.length >= 6) score++
  if (/[A-Z]/.test(pwd)) score++
  if (/[a-z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd)) score++
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd)) score++
  return score
})

const resetPwdStrengthWidth = computed(() => (resetPwdStrength.value / 5 * 100) + '%')
const resetPwdStrengthClass = computed(() => {
  const score = resetPwdStrength.value
  if (score <= 2) return 'weak'
  if (score <= 3) return 'medium'
  return 'strong'
})

const validateEmail = (rule: any, value: any, callback: any) => {
  if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
    callback(new Error('请输入有效的邮箱地址'))
  } else {
    callback()
  }
}

const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度为3-32个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { validator: validateEmail, trigger: 'blur' }
  ]
}

const resetPwdRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ]
}

const loadRoles = async () => { const res = await RoleApi.getAll(); if (res.code === 200) roles.value = res.data || [] }
const loadGroups = async () => { const res = await UserGroupApi.getAll(); if (res.code === 200) groups.value = res.data || [] }

// 格式化日期
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 用户组内成员数量
const groupMemberCount = (groupId: number | undefined) => {
  if (!groupId) return 0
  return tableData.value.filter(u => u.groupId === groupId).length
}

// 点击用户组 → 跳转到用户组页面并筛选
const navigateToGroup = (groupId: number | undefined, groupName: string) => {
  if (!groupId) return
  router.push({ path: '/system/user-group', query: { id: String(groupId) } })
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await UserApi.list({ page: pagination.page, pageSize: pagination.pageSize, keyword: keyword.value || undefined, status: status.value || undefined, roleId: filterRoleId.value, groupId: filterGroupId.value })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; status.value = ''; filterRoleId.value = undefined; filterGroupId.value = undefined; handleSearch() }
const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { username: '', password: '', nickname: '', email: '', phone: '', roleId: undefined, groupId: undefined, status: 'active', mfaEnabled: false, source: 'LOCAL' })
  drawerVisible.value = true
}
const handleEdit = (row: User) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, username: row.username, nickname: row.nickname, email: row.email, phone: row.phone, roleId: row.roleId || undefined, groupId: row.groupId || undefined, status: row.status, mfaEnabled: row.mfaEnabled, source: row.source || 'LOCAL' })
  drawerVisible.value = true
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const data: UpdateUserReq = { id: form.id!, nickname: form.nickname, email: form.email, phone: form.phone, roleId: form.roleId, groupId: form.groupId, status: form.status }
      const res = await UserApi.update(data)
      if (res.code === 200) { ElMessage.success('更新成功'); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || '更新失败')
    } else {
      const data: CreateUserReq = { username: form.username, password: form.password, nickname: form.nickname, email: form.email, phone: form.phone, roleId: form.roleId, groupId: form.groupId, mfaEnabled: form.mfaEnabled }
      const res = await UserApi.create(data)
      if (res.code === 200) { ElMessage.success('创建成功'); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || '创建失败')
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户"${row.username}"吗？`, '删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const res = await UserApi.del(row.id)
    if (res.code === 200) { ElMessage.success('删除成功'); loadData() }
    else ElMessage.error(res.message || '删除失败')
  } catch {}
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个用户吗？`, '批量删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const ids = selectedRows.value.map(row => row.id)
    await UserApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个用户`)
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') { ElMessage.error('批量删除失败') }
  }
}

// 重置密码
const handleResetPassword = () => {
  resetPwdTarget.value = { id: form.id!, username: form.username } as User
  resetPwdForm.newPassword = ''
  resetPwdVisible.value = true
}

const confirmResetPassword = async () => {
  const valid = await resetPwdFormRef.value?.validate().catch(() => false)
  if (!valid) return

  resetPwdLoading.value = true
  try {
    const res = await UserApi.resetPassword({
      userId: resetPwdTarget.value!.id,
      newPassword: resetPwdForm.newPassword
    })
    if (res.code === 200) {
      ElMessage.success('密码重置成功')
      resetPwdVisible.value = false
    } else {
      ElMessage.error(res.message || '重置失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '重置失败')
  } finally {
    resetPwdLoading.value = false
  }
}

// 禁用MFA（从编辑表单）
const handleDisableMFA = async () => {
  try {
    const res = await UserApi.resetMFA({ userId: form.id! })
    if (res.code === 200) {
      ElMessage.success('MFA 已重置')
      form.mfaEnabled = false
      drawerVisible.value = false
      loadData()
    } else {
      ElMessage.error(res.message || '重置失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '重置MFA失败')
  }
}

// 显示启用MFA对话框
const showEnableMFADialog = async () => {
  enableMFATargetId.value = form.id
  mfaSecret.value = ''
  qrCodeUrl.value = ''
  enableMFAForm.code = ''
  enableMFAVisible.value = true
  // 生成MFA密钥
  try {
    const res = await UserApi.generateMFASecret(form.id!)
    if (res.code === 200) {
      mfaSecret.value = res.data.secret
      qrCodeUrl.value = res.data.qrCode
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取MFA信息失败')
  }
}

// 确认启用/重置MFA
const enableMFARules = {
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

const confirmEnableMFA = async () => {
  const valid = await enableMFAFormRef.value?.validate().catch(() => false)
  if (!valid) return

  enableMFALoading.value = true
  try {
    const res = await UserApi.adminEnableMFA({
      userId: enableMFATargetId.value!,
      code: enableMFAForm.code
    })
    if (res.code === 200) {
      ElMessage.success('MFA启用成功')
      enableMFAVisible.value = false
      form.mfaEnabled = true
      loadData()
    } else {
      ElMessage.error(res.message || '启用失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '启用失败')
  } finally {
    enableMFALoading.value = false
  }
}

onMounted(() => { loadRoles(); loadGroups(); loadData() })
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.page {
  min-height: 100vh;
  background: var(--color-page-bg);
  padding: var(--space-4);
}

/* ==================== 页面标题栏 ==================== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
  flex-wrap: wrap;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.page-title {
  font-family: 'Manrope', sans-serif;
  font-size: 17px;
  font-weight: 800;
  color: var(--color-text-primary);
  margin: 0 0 2px;
  letter-spacing: -0.3px;
}

.page-subtitle { font-size: 12px; color: var(--color-text-muted); font-weight: 500; }
.header-actions { display: flex; gap: var(--space-2); }

/* ==================== 筛选栏 ==================== */
.filter-card {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
  align-items: center;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  margin-bottom: var(--space-3);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);

  > * { flex: 0 0 auto; }
  .el-input, .el-select { width: 150px; }
}

/* ==================== 内容卡片 ==================== */
.table-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  overflow: visible;
}

/* ==================== 表格工具栏 ==================== */
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

.toolbar-left {
  .record-count { font-size: 13px; color: var(--color-text-secondary); strong { color: var(--color-text-primary); font-weight: 700; } }
  .selection-count { margin-left: var(--space-4); font-size: 13px; color: var(--color-danger); strong { font-weight: 700; } }
}

.toolbar-right { display: flex; gap: var(--space-2); }

/* ==================== 表格 ==================== */
:deep(.el-table) {
  --el-table-border-color: var(--color-border-light);
  width: 100% !important;

  th.el-table__cell {
    background-color: var(--color-surface-3) !important;
    border-bottom: 1px solid var(--color-border) !important;
    padding: 12px 14px !important;
    color: var(--color-text-secondary) !important;
    font-weight: 700 !important;
    font-size: 11px !important;
    text-transform: uppercase;
    letter-spacing: 0.4px;
  }

  td.el-table__cell {
    padding: 12px 14px !important;
    font-size: 13px;
    color: var(--color-text-primary);
    border-bottom: 1px solid var(--color-border-light) !important;
  }

  .el-table__body tr:hover > td.el-table__cell {
    background-color: var(--color-primary-light-9) !important;
  }
}

/* ==================== 标签徽章 ==================== */
.tag-badge {
  display: inline-block;
  padding: 4px 10px;
  background: var(--color-primary-light-9);
  color: var(--chart-blue);
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 500;
}

.group-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.15s ease;

  svg { color: var(--color-primary); flex-shrink: 0; }

  &:hover {
    background: var(--color-primary-light-9);
    border-color: var(--color-primary);
    color: var(--color-primary);
    svg { color: var(--color-primary); }
  }
}

.group-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-full);
  font-size: 10px;
  font-weight: 700;
  font-family: 'Manrope', sans-serif;
  line-height: 1;
}

.no-group {
  color: var(--color-text-muted);
  font-size: 12px;
  font-style: italic;
}

.mfa-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;

  &--enabled { background: rgba(34,197,94,0.1); color: var(--color-success); }
  &--disabled { background: var(--color-surface-3); color: var(--color-text-muted); }
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
}

.status-badge--success { background: rgba(34,197,94,0.1); color: var(--color-success); }
.status-badge--danger { background: rgba(239,68,68,0.08); color: var(--color-danger); }

/* ==================== 来源标签 ==================== */
.source-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 700;

  &--ad {
    background: rgba(0, 94, 235, 0.1);
    color: var(--color-primary);
    border: 1px solid rgba(0, 94, 235, 0.2);
  }

  &--local {
    background: var(--color-surface-2);
    color: var(--color-text-muted);
    border: 1px solid var(--color-border-light);
  }
}

/* ==================== AD用户提示 ==================== */
.ad-user-notice {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3);
  background: rgba(0, 94, 235, 0.05);
  border: 1px solid rgba(0, 94, 235, 0.2);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-3);

  svg {
    color: var(--color-primary);
    flex-shrink: 0;
    margin-top: 1px;
  }
}

.ad-notice-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.ad-notice-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-primary);
}

.ad-notice-desc {
  font-size: 12px;
  color: var(--color-text-secondary);
  line-height: 1.4;
}

/* ==================== 字段标签 ==================== */
.field-tag {
  display: inline-block;
  margin-top: 4px;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: var(--radius-full);
  text-transform: uppercase;
  letter-spacing: 0.3px;

  &--ad {
    background: rgba(0, 94, 235, 0.1);
    color: var(--color-primary);
  }
}

/* ==================== MFA AD标签 ==================== */
.mfa-ad-tag {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
}

/* ==================== 分页 ==================== */
.pagination-wrapper {
  padding: var(--space-3) var(--space-4);
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

/* ==================== 弹窗内容 ==================== */
.dialog-body-inner {
  padding: 0 4px;
  max-height: 65vh;
  overflow-y: auto;
  &::-webkit-scrollbar { width: 5px; }
  &::-webkit-scrollbar-thumb { background: var(--color-border); border-radius: 3px; }
}

/* ==================== 信息提示条 ==================== */
.info-bar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-3);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-3);
  margin-top: var(--space-1);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
}

.info-bar-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.shield { background: rgba(217,119,6,0.1); color: var(--color-warning); }
}

.info-bar-text { display: flex; flex-direction: column; flex: 1; }
.info-bar-title { font-size: 14px; font-weight: 700; color: var(--color-text-primary); font-family: 'Manrope', sans-serif; }
.info-bar-sub { font-size: 12px; color: var(--color-text-muted); margin-top: 1px; }

/* ==================== 表单 ==================== */
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 0 var(--space-4); }

:deep(.el-form-item) {
  margin-bottom: var(--space-3);
  .el-form-item__label {
    font-weight: 600;
    color: var(--color-text-primary);
    font-size: 13px;
    margin-bottom: 6px;
    &::before { color: var(--color-danger); }
  }
}

:deep(.el-input__wrapper) {
  border-radius: var(--radius-sm) !important;
  box-shadow: none !important;
  border: 1.5px solid var(--color-border-light) !important;
  &:hover { border-color: var(--color-primary) !important; }
  &.is-focus { border-color: var(--color-primary) !important; box-shadow: 0 0 0 3px rgba(0,94,235,0.08) !important; }
}

:deep(.el-select__wrapper) {
  border-radius: var(--radius-sm) !important;
  box-shadow: none !important;
  border: 1.5px solid var(--color-border-light) !important;
  &:hover { border-color: var(--color-primary) !important; }
  &.is-focused { border-color: var(--color-primary) !important; box-shadow: 0 0 0 3px rgba(0,94,235,0.08) !important; }
}

:deep(.el-radio-group) { display: flex; gap: var(--space-4); }
:deep(.el-radio__label) { font-size: 13px; }

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

/* ==================== 密码强度 ==================== */
.password-strength {
  margin-top: var(--space-3);
  .strength-bar { height: 4px; background: var(--color-surface-3); border-radius: 2px; overflow: hidden; margin-bottom: var(--space-2); }
  .strength-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.3s ease, background 0.3s ease;
    &.weak { background: var(--color-danger); }
    &.medium { background: var(--color-warning); }
    &.strong { background: var(--color-success); }
  }
  .strength-labels { display: flex; gap: var(--space-3); justify-content: space-between; }
  .strength-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 3px;
    svg { opacity: 0.2; color: var(--color-text-muted); transition: all 0.2s; }
    &.active svg { opacity: 1; color: var(--color-success); }
  }
}

/* ==================== MFA 状态 ==================== */
.mfa-toggle { display: flex; align-items: center; gap: var(--space-3); }
.mfa-status {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 500;

  &.enabled { background: rgba(34,197,94,0.1); color: var(--color-success); }
  &.disabled { background: var(--color-surface-3); color: var(--color-text-muted); }
}

/* ==================== 状态选择 ==================== */
.status-options { display: flex; gap: var(--space-2); }
.status-opt {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-muted);
  transition: all 0.15s;

  &:hover { border-color: var(--color-border); color: var(--color-text-primary); }
  &.active { border-color: var(--color-primary); background: var(--color-primary-light-9); color: var(--color-primary); }
}

.status-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.status-dot--success { background: var(--color-success); }
.status-dot--danger { background: var(--color-danger); }

/* ==================== 行操作按钮 ==================== */
.row-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.row-action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  cursor: pointer;
  transition: all 0.15s ease;

  &--edit {
    color: var(--gray-400);
    &:hover {
      color: var(--color-primary);
      background: var(--color-primary-light-9);
    }
  }

  &--delete {
    color: var(--gray-400);
    &:hover {
      color: var(--color-danger);
      background: var(--color-danger-bg);
    }
  }
}

/* ==================== 安全设置卡片 ==================== */
.security-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  margin-bottom: 8px;
}

.security-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  transition: border-color 0.15s ease;

  &:hover { border-color: var(--color-primary); }
}

.security-card-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  background: rgba(217, 119, 6, 0.1);
  color: var(--color-warning);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.mfa-on {
    background: rgba(34, 197, 94, 0.1);
    color: var(--color-success);
  }
}

.security-card-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.security-card-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.security-card-desc {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
}

/* ==================== 弹窗底部 ==================== */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
  .el-button { min-width: 80px; border-radius: var(--radius-sm); font-weight: 600; }
}

/* ==================== QR码 ==================== */
.enable-mfa-content {
  .enable-mfa-tip { font-size: 13px; color: var(--color-text-muted); margin-bottom: var(--space-4); line-height: 1.5; }
  .qr-code-box {
    text-align: center;
    padding: var(--space-4);
    background: var(--color-surface-2);
    border-radius: var(--radius-md);
    margin-bottom: var(--space-4);
    .qr-code-img { width: 160px; height: 160px; margin-bottom: var(--space-2); }
    .secret-text { font-size: 11px; color: var(--color-text-muted); font-family: 'SF Mono', monospace; }
  }
}

/* ==================== 响应式 ==================== */
@media (max-width: 1366px) {
  .page { padding: var(--space-3); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-card { padding: var(--space-2) var(--space-3); margin-bottom: var(--space-2); }
}
</style>
