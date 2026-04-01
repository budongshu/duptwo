<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('user.list.title') }}</h1>
        <span class="page-subtitle">{{ t('user.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('user.list.create') }}
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input v-model="keyword" :placeholder="t('user.list.searchPlaceholder')" clearable @keyup.enter="handleSearch" style="width: 240px">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </template>
      </el-input>
      <el-select v-model="status" :placeholder="t('user.list.filter.status')" clearable style="width: 110px" @change="handleSearch">
        <el-option :label="t('common.enabled')" value="active" />
        <el-option :label="t('common.disabled')" value="inactive" />
      </el-select>
      <el-select v-model="filterRoleId" :placeholder="t('user.list.filter.role')" clearable style="width: 130px" @change="handleSearch">
        <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
      </el-select>
      <el-select v-model="filterGroupId" :placeholder="t('user.list.filter.group')" clearable style="width: 130px" @change="handleSearch">
        <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
      </el-select>
      <el-button type="primary" @click="handleSearch">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        {{ t('common.search') }}
      </el-button>
      <el-button @click="handleReset">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4"/></svg>
        {{ t('common.reset') }}
      </el-button>
      <el-button v-if="selectedRows.length > 0" type="danger" @click="handleBatchDelete">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
        {{ t('common.batchDelete') }} ({{ selectedRows.length }})
      </el-button>
    </div>

    <!-- 表格卡片 -->
    <div class="content-card">
      <el-table ref="tableRef" v-model:selection="selectedRows" :data="tableData" v-loading="loading" stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="40" fixed="left" />
        <el-table-column prop="nickname" :label="t('user.list.table.name')" min-width="110" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar" :style="{ background: getAvatarColor(row.nickname) }">{{ (row.nickname || row.username || 'U').charAt(0).toUpperCase() }}</div>
              <span class="user-name">{{ row.nickname || row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="username" :label="t('user.list.table.username')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="email" :label="t('user.list.table.email')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="roleName" :label="t('user.list.table.role')" min-width="100" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.roleName || '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="groupName" :label="t('user.list.table.group')" min-width="110">
          <template #default="{ row }">
            <span v-if="row.groupName" class="group-link" @click="jumpToGroup(row.groupId)">{{ row.groupName }}</span>
            <span v-else class="empty-text">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('common.status')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light">
              {{ row.status === 'active' ? t('common.enabled') : t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mfaEnabled" :label="t('user.list.table.mfa')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.mfaEnabled ? 'success' : 'info'" size="small" effect="plain">
              {{ row.mfaEnabled ? t('user.list.table.mfaOn') : t('user.list.table.mfaOff') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('common.createdAt')" min-width="170">
          <template #default="{ row }">
            <span class="time-text">{{ formatDate(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: t('common.edit'), type: 'primary' },
              { key: 'resetPwd', label: t('user.list.form.resetPwd'), type: 'warning' },
              { key: 'delete', label: t('common.delete'), type: 'danger' }
            ]" @action="(key) => handleAction(key, row)" />
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <span class="record-info">{{ t('common.totalRecords', { total: pagination.total }) }}</span>
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="sizes, prev, pager, next"
          background
        />
      </div>
    </div>

    <!-- 编辑/新增抽屉 -->
    <el-drawer v-model="drawerVisible" direction="rtl" size="460px" :destroy-on-close="true" class="personnel-drawer">
      <template #header>
        <div class="drawer-title-inner">
          <span class="drawer-mode-tag" :class="isEdit ? 'tag--edit' : 'tag--new'">{{ isEdit ? t('common.edit') : t('common.create') }}</span>
          <span class="drawer-title-text">{{ isEdit ? (form.nickname || form.username || t('user.list.form.user')) : t('user.list.form.newUser') }}</span>
        </div>
      </template>
      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top" class="edit-form">
          <el-form-item :label="t('user.list.form.nickname')" prop="nickname">
            <el-input v-model="form.nickname" :placeholder="t('user.list.form.nicknamePlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('user.list.form.username')" prop="username">
            <el-input v-model="form.username" :placeholder="t('user.list.form.usernamePlaceholder')" :disabled="isEdit" />
          </el-form-item>
          <el-form-item :label="t('user.list.form.email')" prop="email">
            <el-input v-model="form.email" :placeholder="t('user.list.form.emailPlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('user.list.form.password')" :prop="isEdit ? '' : 'password'">
            <el-input v-model="form.password" type="password" show-password :placeholder="isEdit ? t('user.list.form.passwordEditTip') : t('user.list.form.passwordPlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('user.list.form.role')" prop="roleId">
            <el-select v-model="form.roleId" :placeholder="t('user.list.form.rolePlaceholder')" style="width: 100%" clearable>
              <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('user.list.form.group')" prop="groupId">
            <el-select v-model="form.groupId" :placeholder="t('user.list.form.groupPlaceholder')" style="width: 100%" clearable>
              <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('common.status')" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="active">{{ t('common.enabled') }}</el-radio>
              <el-radio label="inactive">{{ t('common.disabled') }}</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item :label="t('common.sort')" prop="sort">
            <el-input-number v-model="form.sort" :min="0" :max="9999" />
          </el-form-item>
          <!-- 操作区 -->
          <div v-if="isEdit" class="action-section">
            <div class="action-divider"></div>
            <div class="action-row">
              <el-button type="warning" plain @click="handleResetPwd">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                {{ t('user.list.form.resetPwd') }}
              </el-button>
              <el-button v-if="!form.mfaEnabled" type="success" plain @click="handleEnableMfa">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>
                {{ t('user.list.form.enableMfa') }}
              </el-button>
              <el-button v-else type="info" plain @click="handleResetMfa">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                {{ t('user.list.form.resetMfa') }}
              </el-button>
            </div>
          </div>
        </el-form>
      </div>
      <div class="drawer-foot">
        <el-button @click="drawerVisible = false">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ t('common.cancel') }}
        </el-button>
        <el-button type="primary" :loading="submitting" @click="confirmSubmit">
          <svg v-if="!submitting" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          {{ isEdit ? t('common.save') : t('common.create') }}
        </el-button>
      </div>
    </el-drawer>

    <!-- 重置密码弹窗 -->
    <el-dialog v-model="showResetPwd" width="480px" destroy-on-close class="reset-pwd-dialog">
      <template #header>
        <div class="dialog-head">
          <div class="dialog-icon dialog-icon--warning">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          </div>
          <div class="dialog-head-text">
            <span class="dialog-title-text">{{ t('user.list.form.resetPwd') }}</span>
            <span class="dialog-subtitle">{{ resetPwdTarget ? (resetPwdTarget.nickname || resetPwdTarget.username) : '' }} 的密码</span>
          </div>
        </div>
      </template>
      <div class="reset-pwd-body">
        <div class="pwd-notice">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
          <span>请设置一个新密码，并当面或通过安全渠道告知用户</span>
        </div>

        <el-form ref="resetPwdFormRef" :model="resetPwdForm" :rules="resetPwdRules" label-position="top">
          <el-form-item :label="t('user.list.form.newPassword')" prop="newPassword">
            <el-input v-model="resetPwdForm.newPassword" type="password" show-password size="large" placeholder="输入新密码" />
          </el-form-item>

          <!-- 密码要求清单 -->
          <div class="pwd-requirements" v-if="resetPwdForm.newPassword">
            <div class="req-title">密码要求</div>
            <div class="req-list">
              <div class="req-item" :class="{ satisfied: resetPwdChecks.length }">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>至少 {{ pwdPolicy.passwordMinLength }} 个字符</span>
              </div>
              <div class="req-item" :class="{ satisfied: resetPwdChecks.upper }" v-if="pwdPolicy.passwordRequireUppercase">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含大写字母 (A-Z)</span>
              </div>
              <div class="req-item" :class="{ satisfied: resetPwdChecks.lower }" v-if="pwdPolicy.passwordRequireLowercase">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含小写字母 (a-z)</span>
              </div>
              <div class="req-item" :class="{ satisfied: resetPwdChecks.number }" v-if="pwdPolicy.passwordRequireDigit">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含数字 (0-9)</span>
              </div>
              <div class="req-item" :class="{ satisfied: resetPwdChecks.special }" v-if="pwdPolicy.passwordRequireSpecial">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <span>包含特殊字符 (!@#$%...)</span>
              </div>
            </div>
          </div>

          <el-form-item :label="t('user.list.form.confirmPassword')" prop="confirmPassword">
            <el-input v-model="resetPwdForm.confirmPassword" type="password" show-password size="large" placeholder="再次输入新密码" />
          </el-form-item>
        </el-form>

        <!-- 显示设置的密码（供管理员参考） -->
        <div class="pwd-preview" v-if="resetPwdForm.newPassword && resetPwdForm.newPassword === resetPwdForm.confirmPassword">
          <div class="pwd-preview-label">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
            密码预览（请当面告知用户）
          </div>
          <div class="pwd-preview-value">{{ resetPwdForm.newPassword }}</div>
        </div>
      </div>
      <template #footer>
        <el-button size="default" @click="showResetPwd = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="resetPwdLoading" @click="confirmResetPwd">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { UserApi, type User, type CreateUserReq, type UpdateUserReq } from '@/api/user'
import { RoleApi } from '@/api/role'
import { UserGroupApi } from '@/api/user-group'
import { AdminApi } from '@/api/admin'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<User[]>([])
const tableRef = ref()
const selectedRows = ref<User[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 重置密码相关
const showResetPwd = ref(false)
const resetPwdTarget = ref<User | null>(null)
const resetPwdFormRef = ref<FormInstance>()
const resetPwdLoading = ref(false)
const resetPwdForm = reactive({ newPassword: '', confirmPassword: '' })

// 密码策略
const pwdPolicy = ref({
  passwordMinLength: 8,
  passwordRequireUppercase: false,
  passwordRequireLowercase: false,
  passwordRequireDigit: false,
  passwordRequireSpecial: false
})

const resetPwdRules = computed<FormRules>(() => ({
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: pwdPolicy.value.passwordMinLength, message: `密码至少 ${pwdPolicy.value.passwordMinLength} 个字符`, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== resetPwdForm.newPassword) callback(new Error('两次输入的密码不一致'))
        else callback()
      },
      trigger: 'blur'
    }
  ]
}))

// 密码检查
const resetPwdChecks = computed(() => ({
  length: resetPwdForm.newPassword.length >= pwdPolicy.value.passwordMinLength,
  upper: !pwdPolicy.value.passwordRequireUppercase || /[A-Z]/.test(resetPwdForm.newPassword),
  lower: !pwdPolicy.value.passwordRequireLowercase || /[a-z]/.test(resetPwdForm.newPassword),
  number: !pwdPolicy.value.passwordRequireDigit || /[0-9]/.test(resetPwdForm.newPassword),
  special: !pwdPolicy.value.passwordRequireSpecial || /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(resetPwdForm.newPassword)
}))

const loadPwdPolicy = async () => {
  try {
    const res = await AdminApi.getSecuritySettings()
    if (res.code === 200 && res.data) {
      pwdPolicy.value = {
        passwordMinLength: res.data.passwordMinLength || 8,
        passwordRequireUppercase: res.data.passwordRequireUppercase || false,
        passwordRequireLowercase: res.data.passwordRequireLowercase || false,
        passwordRequireDigit: res.data.passwordRequireDigit || false,
        passwordRequireSpecial: res.data.passwordRequireSpecial || false
      }
    }
  } catch {}
}

const roles = ref<any[]>([])
const groups = ref<any[]>([])

const keyword = ref('')
const status = ref('')
const filterRoleId = ref<number | undefined>()
const filterGroupId = ref<number | undefined>()
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateUserReq & { id?: number; roleId?: number; groupId?: number }>({
  username: '', nickname: '', email: '', password: '', roleId: undefined, groupId: undefined, status: 'active', sort: 0,
})

const formRules = {
  username: [{ required: true, message: t('user.list.form.usernameRequired'), trigger: 'blur' }],
  nickname: [{ required: true, message: t('user.list.form.nicknameRequired'), trigger: 'blur' }],
  email: [{ type: 'email', message: t('user.list.form.emailFormatTip'), trigger: 'blur' }],
}

// 格式化时间 YYYY-MM-DD HH:mm:ss
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  const ss = String(d.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}:${ss}`
}

// 头像颜色
const avatarColors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#00BFA5', '#7C3AED', '#DB2777']
const getAvatarColor = (name: string) => {
  const idx = (name || '').charCodeAt(0) % avatarColors.length
  return avatarColors[idx]
}

// 跳转到用户组
const jumpToGroup = (groupId: number | undefined) => {
  if (groupId) router.push({ path: '/user-groups', query: { id: String(groupId) } })
}

const loadRolesAndGroups = async () => {
  try {
    const [rRes, gRes] = await Promise.all([RoleApi.getAll(), UserGroupApi.getAll()])
    if (rRes.code === 200) roles.value = rRes.data || []
    if (gRes.code === 200) groups.value = gRes.data || []
  } catch {}
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await UserApi.list({
      page: pagination.page, pageSize: pagination.pageSize,
      keyword: keyword.value || undefined,
      status: status.value || undefined,
      roleId: filterRoleId.value || undefined,
      groupId: filterGroupId.value || undefined,
    })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; status.value = ''; filterRoleId.value = undefined; filterGroupId.value = undefined; pagination.page = 1; loadData() }
const handleSelectionChange = (rows: User[]) => { selectedRows.value = rows }

const handleCreate = async () => {
  isEdit.value = false
  await loadRolesAndGroups()
  Object.assign(form, { id: undefined, username: '', nickname: '', email: '', password: '', roleId: undefined, groupId: undefined, status: 'active', sort: 0 })
  drawerVisible.value = true
}

const handleEdit = async (row: User) => {
  isEdit.value = true
  await loadRolesAndGroups()
  Object.assign(form, {
    id: row.id, username: row.username, nickname: row.nickname || '',
    email: row.email || '', password: '', roleId: row.roleId, groupId: row.groupId,
    status: row.status, sort: row.sort || 0, mfaEnabled: row.mfaEnabled,
  })
  drawerVisible.value = true
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const data = { id: form.id as number, username: form.username, nickname: form.nickname, email: form.email, roleId: form.roleId, groupId: form.groupId, status: form.status, sort: form.sort } as UpdateUserReq
      if (form.password) { (data as any).password = form.password }
      const res = await UserApi.update(data)
      if (res.code === 200) { ElMessage.success(t('common.updateSuccess')); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || t('user.list.messages.updateFailed'))
    } else {
      const res = await UserApi.create(form as CreateUserReq)
      if (res.code === 200) { ElMessage.success(t('common.createSuccess')); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || t('user.list.messages.createFailed'))
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm(t('user.list.messages.deleteConfirm', { name: row.nickname || row.username }), t('user.list.messages.deleteConfirmTitle'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    const res = await UserApi.del(row.id)
    if (res.code === 200) { ElMessage.success(t('common.deleteSuccess')); loadData() }
    else ElMessage.error(res.message || t('common.deleteError'))
  } catch {}
}

const handleBatchDelete = async () => {
  if (!selectedRows.value.length) return
  try {
    await ElMessageBox.confirm(t('user.list.messages.batchDeleteConfirm', { count: selectedRows.value.length }), t('common.batchConfirm'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    const ids = selectedRows.value.map(r => r.id)
    await UserApi.batchDelete(ids)
    ElMessage.success(t('user.list.messages.batchDeleteSuccess', { count: selectedRows.value.length }))
    selectedRows.value = []; loadData()
  } catch {}
}

const handleResetPwd = async (row: User) => {
  resetPwdTarget.value = row
  resetPwdForm.newPassword = ''
  resetPwdForm.confirmPassword = ''
  await loadPwdPolicy()
  showResetPwd.value = true
}

const confirmResetPwd = async () => {
  const valid = await resetPwdFormRef.value?.validate().catch(() => false)
  if (!valid) return
  resetPwdLoading.value = true
  try {
    const res = await UserApi.resetPassword({
      userId: resetPwdTarget.value!.id,
      newPassword: resetPwdForm.newPassword
    })
    if (res.code === 200) {
      ElMessage.success(t('user.list.messages.resetPwdSuccess'))
      showResetPwd.value = false
    } else {
      ElMessage.error(res.message || t('user.list.messages.resetPwdFailed'))
    }
  } finally {
    resetPwdLoading.value = false
  }
}

const handleEnableMfa = async () => {
  try {
    await ElMessageBox.confirm(t('user.list.messages.enableMfaConfirm', { name: form.nickname || form.username }), t('user.list.form.enableMfa'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'info' })
    const res = await UserApi.adminEnableMfa(form.id as number)
    if (res.code === 200) { ElMessage.success(t('user.list.messages.mfaEnableSuccess')); form.mfaEnabled = true }
    else ElMessage.error(res.message || t('user.list.messages.mfaEnableFailed'))
  } catch {}
}

const handleResetMfa = async () => {
  try {
    await ElMessageBox.confirm(t('user.list.messages.resetMfaConfirm', { name: form.nickname || form.username }), t('user.list.form.resetMfa'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    const res = await UserApi.resetMfa(form.id as number)
    if (res.code === 200) { ElMessage.success(t('user.list.messages.mfaResetSuccess')); form.mfaEnabled = false }
    else ElMessage.error(res.message || t('user.list.messages.mfaResetFailed'))
  } catch {}
}

const handleAction = (key: string, row: User) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
  else if (key === 'resetPwd') handleResetPwd(row)
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => loadData())
</script>

<script lang="ts">
export default { name: 'UserList' }
</script>

<style scoped lang="scss">
.page {
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
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.header-left { display: flex; align-items: baseline; gap: var(--space-3); }
.page-title { font-family: 'Manrope', sans-serif; font-size: 17px; font-weight: 800; color: var(--color-text-primary); margin: 0; letter-spacing: -0.3px; }
.page-subtitle { font-size: 12px; color: var(--color-text-muted); font-weight: 500; }
.header-actions { display: flex; gap: var(--space-2); }

.filter-bar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  flex-wrap: wrap;
}

.content-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  overflow: visible;
  display: flex;
  flex-direction: column;
  flex: 1;
}

:deep(.el-table) {
  font-size: 13px;
  th.el-table__cell {
    background-color: var(--color-surface-3) !important;
    color: var(--color-text-secondary) !important;
    font-weight: 600;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.4px;
    padding: 10px 12px !important;
    border-bottom: 1px solid var(--color-border) !important;
  }
  td.el-table__cell { padding: 9px 12px !important; border-bottom: 1px solid var(--color-border-light) !important; color: var(--color-text-primary); }
  .el-table__body tr:hover > td.el-table__cell { background-color: var(--color-primary-light-9) !important; }
}

.user-cell { display: flex; align-items: center; gap: 8px; }
.user-avatar { width: 28px; height: 28px; border-radius: 50%; color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 12px; flex-shrink: 0; }
.user-name { font-weight: 600; font-size: 13px; color: var(--color-text-primary); }

.group-link {
  color: var(--color-primary);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  text-decoration: none;
  &:hover { text-decoration: underline; }
}

.empty-text { color: var(--color-text-secondary); font-size: 13px; }

.time-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-family: 'SF Mono', 'DM Mono', monospace;
}

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
}
.record-info { font-size: 12px; color: var(--color-text-secondary); }

/* 抽屉 */
:deep(.personnel-drawer) {
  .el-drawer__header { padding: 10px 16px; margin-bottom: 0; border-bottom: 1px solid var(--color-border-light); background: var(--color-surface); align-items: center; }
  .el-drawer__body { padding: 0; display: flex; flex-direction: column; overflow: hidden; }
}
.drawer-title-inner { display: flex; align-items: center; gap: 8px; }
.drawer-mode-tag {
  font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: var(--radius-full);
  &.tag--edit { background: rgba(0, 94, 235, 0.1); color: var(--color-primary); }
  &.tag--new { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
}
.drawer-title-text { font-family: 'Manrope', sans-serif; font-size: 14px; font-weight: 700; color: var(--color-text-primary); }
.drawer-body { flex: 1; overflow-y: auto; padding: 16px; background: var(--color-surface-2); &::-webkit-scrollbar { width: 3px; } &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; } }
.drawer-foot { display: flex; justify-content: flex-end; gap: 8px; padding: 12px 16px; background: var(--color-surface); border-top: 1px solid var(--color-border-light); flex-shrink: 0; }

/* 表单 */
.edit-form { display: flex; flex-direction: column; gap: var(--space-3); animation: form-in 0.3s ease both 0.05s; }
@keyframes form-in { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }
.edit-form :deep(.el-form-item) {
  margin-bottom: 0;
  .el-form-item__label { font-size: 12px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 4px; }
}

.action-section { margin-top: var(--space-2); }
.action-divider { height: 1px; background: var(--color-border-light); margin-bottom: var(--space-3); }
.action-row { display: flex; gap: var(--space-2); flex-wrap: wrap; }

/* 重置密码弹窗 */
.reset-pwd-dialog {
  .el-dialog__header {
    padding: 16px 20px;
    margin-right: 0;
    border-bottom: 1px solid var(--color-border-light);
  }
}

.dialog-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.dialog-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &--warning {
    background: rgba(245, 158, 11, 0.1);
    border: 1px solid rgba(245, 158, 11, 0.2);
    svg { stroke: #f59e0b; }
  }
}

.dialog-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.dialog-subtitle {
  font-size: 12px;
  color: var(--color-text-muted);
}

.reset-pwd-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.pwd-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: var(--radius-md);
  font-size: 12.5px;
  color: #92400e;

  svg { stroke: #d97706; flex-shrink: 0; }
}

.pwd-requirements {
  padding: 12px 14px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.req-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.req-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.req-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--color-text-muted);
  transition: color 0.2s;

  svg {
    stroke: var(--color-text-muted);
    flex-shrink: 0;
    opacity: 0.4;
    transition: all 0.2s;
  }

  &.satisfied {
    color: var(--color-success);
    svg {
      stroke: var(--color-success);
      opacity: 1;
    }
  }
}

.pwd-preview {
  padding: 12px 14px;
  background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
  border: 1px solid #86efac;
  border-radius: var(--radius-md);
}

.pwd-preview-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #166534;
  font-weight: 600;
  margin-bottom: 6px;

  svg { stroke: #166534; }
}

.pwd-preview-value {
  font-family: 'SF Mono', 'Consolas', monospace;
  font-size: 16px;
  font-weight: 700;
  color: #15803d;
  letter-spacing: 1px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-sm);
  border: 1px dashed #86efac;
}

@media (max-width: 1366px) {
  .page { padding: var(--space-3); gap: var(--space-2); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-bar { padding: var(--space-2) var(--space-3); gap: var(--space-2); }
}
</style>
