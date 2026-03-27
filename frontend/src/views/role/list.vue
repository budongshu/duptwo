<template>
  <div class="page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">角色管理</h1>
        <p class="page-subtitle">管理系统角色和权限</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          新增角色
        </el-button>
      </div>
    </header>

    <div class="filter-card">
      <el-input v-model="keyword" placeholder="搜索角色名称/编码" clearable @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <div class="table-card">
      <div class="table-toolbar">
        <div class="toolbar-left">
          <span class="record-count">共 <strong>{{ pagination.total }}</strong> 条</span>
          <span v-if="selectedRows.length > 0" class="selection-count">
            已选 <strong>{{ selectedRows.length }}</strong> 项
          </span>
        </div>
        <div class="toolbar-right">
          <el-button v-if="selectedRows.length > 0" type="danger" class="batch-delete-btn" @click="handleBatchDelete">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            批量删除
          </el-button>
        </div>
      </div>

      <el-table ref="tableRef" v-model:selection="selectedRows" :data="tableData" v-loading="loading">
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column prop="name" label="角色名称" min-width="140">
          <template #default="{ row }">
            <div class="role-name-cell">
              <div class="role-avatar">{{ row.name?.charAt(0) || 'R' }}</div>
              <span class="role-name-text">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="角色编码" min-width="160">
          <template #default="{ row }">
            <code class="code-text">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="permissions" label="权限" min-width="300">
          <template #default="{ row }">
            <div class="perm-tags">
              <el-tag v-for="p in (row.permissions || []).slice(0, 3)" :key="p" type="info" size="small">{{ getPermLabel(p) }}</el-tag>
              <el-tag v-if="(row.permissions || []).length > 3" type="warning" size="small">+{{ row.permissions.length - 3 }}</el-tag>
              <span v-if="!row.permissions?.length" class="empty-text">—</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" min-width="80" align="center" />
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: '编辑', type: 'primary' },
              { key: 'delete', label: '删除', type: 'danger' }
            ]" @action="(key) => handleAction(key, row)" />
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize" :total="pagination.total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" background />
      </div>
    </div>

    <!-- 角色编辑弹窗 -->
    <el-dialog v-model="dialogVisible" width="640px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="dialog-title-text">{{ isEdit ? '编辑角色' : '新增角色' }}</span>
        </div>
      </template>
      <div class="dialog-body-inner">
        <!-- 基本信息区 -->
        <div class="info-bar">
          <div class="info-bar-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/>
            </svg>
          </div>
          <div class="info-bar-text">
            <span class="info-bar-title">基本信息</span>
            <span class="info-bar-sub">设置角色名称和编码</span>
          </div>
        </div>

        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="角色名称" prop="name">
                <el-input v-model="form.name" placeholder="请输入角色名称" clearable>
                  <template #prefix>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="角色编码" prop="code">
                <el-input v-model="form.code" placeholder="如：admin" :disabled="isEdit" clearable>
                  <template #prefix>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="描述">
                <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入角色描述信息" clearable />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="排序">
                <el-input-number v-model="form.sort" :min="0" :max="9999" style="width: 100%" />
                <span class="form-tip">数字越小排序越靠前</span>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>

        <!-- 权限设置区 -->
        <div class="info-bar">
          <div class="info-bar-icon shield">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          <div class="info-bar-text">
            <span class="info-bar-title">权限配置</span>
            <span class="info-bar-sub">勾选该角色拥有的操作权限</span>
          </div>
          <div class="perm-quick-actions">
            <el-button size="small" text @click="checkAllPerms">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 3px"><polyline points="20 6 9 17 4 12"/></svg>
              全选
            </el-button>
            <el-button size="small" text @click="clearAllPerms">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 3px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              清空
            </el-button>
          </div>
        </div>

        <!-- 权限分组展示 -->
        <div class="perm-groups">
          <div v-for="group in permGroups" :key="group.key" class="perm-group-card">
            <div class="perm-group-header">
              <div class="perm-group-icon" :style="{ background: group.bgColor, color: group.color }">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="group.iconPath"></svg>
              </div>
              <span class="perm-group-name">{{ group.name }}</span>
              <span class="perm-group-count">{{ getGroupCheckedCount(group.perms) }}/{{ group.perms.length }}</span>
            </div>
            <div class="perm-group-items">
              <div
                v-for="perm in group.perms"
                :key="perm.value"
                class="perm-item"
                :class="{ 'perm-checked': form.permissions.includes(perm.value) }"
                @click="togglePerm(perm.value)"
              >
                <div class="perm-checkbox" :class="{ checked: form.permissions.includes(perm.value) }">
                  <svg v-if="form.permissions.includes(perm.value)" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                </div>
                <div class="perm-item-info">
                  <span class="perm-item-label">{{ perm.label }}</span>
                  <span class="perm-item-value">{{ perm.value }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 5px">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { RoleApi, type Role, type CreateRoleReq, type UpdateRoleReq } from '@/api/role'
import TableActions from '@/components/TableActions.vue'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Role[]>([])
const tableRef = ref()
const selectedRows = ref<Role[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const keyword = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateRoleReq & { id?: number }>({ name: '', code: '', description: '', permissions: [], sort: 0 })

const formRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}

// 权限分组配置
const permGroups = [
  {
    key: 'upload',
    name: '上传记录',
    bgColor: '#e8f4ff',
    color: '#005eeb',
    iconPath: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/>',
    perms: [
      { label: '创建', value: 'upload:create' },
      { label: '读取', value: 'upload:read' },
      { label: '更新', value: 'upload:update' },
      { label: '删除', value: 'upload:delete' },
      { label: '导出', value: 'upload:export' },
    ]
  },
  {
    key: 'field',
    name: '字段配置',
    bgColor: '#f0f8e8',
    color: '#16a34a',
    iconPath: '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/>',
    perms: [
      { label: '创建', value: 'field-config:create' },
      { label: '读取', value: 'field-config:read' },
      { label: '更新', value: 'field-config:update' },
      { label: '删除', value: 'field-config:delete' },
    ]
  },
  {
    key: 'project',
    name: '项目管理',
    bgColor: '#fff7ed',
    color: '#c2410c',
    iconPath: '<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>',
    perms: [
      { label: '创建', value: 'project:create' },
      { label: '读取', value: 'project:read' },
      { label: '更新', value: 'project:update' },
      { label: '删除', value: 'project:delete' },
    ]
  },
  {
    key: 'personnel',
    name: '人员管理',
    bgColor: '#fdf4ff',
    color: '#a21caf',
    iconPath: '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>',
    perms: [
      { label: '创建', value: 'personnel:create' },
      { label: '读取', value: 'personnel:read' },
      { label: '更新', value: 'personnel:update' },
      { label: '删除', value: 'personnel:delete' },
      { label: '导出', value: 'personnel:export' },
    ]
  },
  {
    key: 'user',
    name: '用户管理',
    bgColor: '#fef3e8',
    color: '#d97706',
    iconPath: '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>',
    perms: [
      { label: '创建', value: 'user:create' },
      { label: '读取', value: 'user:read' },
      { label: '更新', value: 'user:update' },
      { label: '删除', value: 'user:delete' },
    ]
  },
  {
    key: 'role',
    name: '角色管理',
    bgColor: '#f3e8ff',
    color: '#7c3aed',
    iconPath: '<path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/>',
    perms: [
      { label: '创建', value: 'role:create' },
      { label: '读取', value: 'role:read' },
      { label: '更新', value: 'role:update' },
      { label: '删除', value: 'role:delete' },
    ]
  },
  {
    key: 'audit',
    name: '审计日志',
    bgColor: '#fef0f0',
    color: '#dc2626',
    iconPath: '<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/>',
    perms: [
      { label: '操作日志', value: 'audit:operation:read' },
      { label: '登录日志', value: 'audit:login:read' },
    ]
  },
  {
    key: 'system',
    name: '系统配置',
    bgColor: '#f0f8f8',
    color: '#0891b2',
    iconPath: '<circle cx="12" cy="12" r="3"/><path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/>',
    perms: [
      { label: '读取', value: 'config:read' },
      { label: '更新', value: 'config:update' },
      { label: '管理员', value: 'admin:all' },
    ]
  },
]

// 权限标签映射
const permLabelMap: Record<string, string> = {}
permGroups.forEach(g => g.perms.forEach(p => { permLabelMap[p.value] = p.label }))

const getPermLabel = (code: string) => permLabelMap[code] || code.split(':').pop() || code

const getGroupCheckedCount = (perms: { value: string }[]) => {
  return perms.filter(p => form.permissions.includes(p.value)).length
}

const togglePerm = (value: string) => {
  const idx = form.permissions.indexOf(value)
  if (idx > -1) form.permissions.splice(idx, 1)
  else form.permissions.push(value)
}

const checkAllPerms = () => {
  form.permissions = permGroups.flatMap(g => g.perms.map(p => p.value))
}

const clearAllPerms = () => {
  form.permissions = []
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await RoleApi.list({ page: pagination.page, pageSize: pagination.pageSize, keyword: keyword.value || undefined })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; handleSearch() }

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: undefined, name: '', code: '', description: '', permissions: [], sort: 0 })
  dialogVisible.value = true
}

const handleEdit = (row: Role) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, name: row.name, code: row.code, description: row.description, permissions: [...(row.permissions || [])], sort: row.sort })
  dialogVisible.value = true
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const res = await RoleApi.update({ id: form.id as number, name: form.name, code: form.code, description: form.description, permissions: form.permissions, sort: form.sort })
      if (res.code === 200) { ElMessage.success('更新成功'); dialogVisible.value = false; loadData() }
      else ElMessage.error(res.message || '更新失败')
    } else {
      const res = await RoleApi.create(form)
      if (res.code === 200) { ElMessage.success('创建成功'); dialogVisible.value = false; loadData() }
      else ElMessage.error(res.message || '创建失败')
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色"${row.name}"吗？`, '删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const res = await RoleApi.del(row.id)
    if (res.code === 200) { ElMessage.success('删除成功'); loadData() }
    else ElMessage.error(res.message || '删除失败')
  } catch {}
}

const handleAction = (key: string, row: Role) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个角色吗？`, '批量删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const ids = selectedRows.value.map(row => row.id)
    await RoleApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个角色`)
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') { ElMessage.error('批量删除失败') }
  }
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })

onMounted(() => loadData())
</script>

<script lang="ts">
import { watch } from 'vue'
export default { name: 'RoleList' }
</script>

<style scoped lang="scss">
.page { min-height: 100vh; background: var(--color-page-bg); padding: 20px; }

.page-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 16px;
  padding: 20px 24px; background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
}

.header-left {
  display: flex; align-items: baseline; gap: 10px;
}

.page-title { font-family: 'Manrope', sans-serif; font-size: 20px; font-weight: 800; color: var(--color-text-primary); margin: 0; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); margin: 0; }
.header-actions { display: flex; gap: 8px; }

.filter-card {
  display: flex; flex-wrap: wrap; gap: 10px; align-items: center;
  background: var(--color-surface); border-radius: 12px; padding: 16px 20px; margin-bottom: 14px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
  > * { flex: 0 0 auto; }
  .el-input { width: 200px; }
}

.table-card { background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light); overflow: visible; }

.table-toolbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px 16px; border-bottom: 1px solid var(--color-border-light); background: var(--color-surface-2);
}

.toolbar-left {
  .record-count { font-size: 13px; color: var(--color-text-secondary); strong { color: var(--color-text-primary); font-weight: 700; } }
  .selection-count { margin-left: 16px; font-size: 13px; color: var(--color-primary); strong { font-weight: 700; } }
}

.toolbar-right { display: flex; gap: 10px; align-items: center; }
.batch-delete-btn { display: inline-flex !important; align-items: center; flex-shrink: 0; }

:deep(.el-table) {
  th.el-table__cell {
    background-color: var(--color-surface-2) !important; border-bottom: 1px solid var(--color-border-light) !important;
    padding: 12px 14px !important; color: var(--color-text-secondary) !important; font-weight: 700 !important; font-size: 12px !important;
  }
  td.el-table__cell { padding: 12px 14px !important; font-size: 14px; color: var(--color-text-primary); border-bottom: 1px solid var(--color-border-light) !important; }
  .el-table__body tr:hover > td.el-table__cell { background-color: var(--color-surface-2) !important; }
}

.role-name-cell { display: flex; align-items: center; gap: 10px; }
.role-avatar {
  width: 30px; height: 30px; border-radius: 50%;
  background: var(--color-primary); color: var(--color-surface);
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 13px; flex-shrink: 0;
}
.role-name-text { font-weight: 600; }

.code-text { font-size: 12px; color: var(--color-primary); background: var(--color-primary-light-9); padding: 3px 8px; border-radius: 5px; font-weight: 500; }
.perm-tags { display: flex; flex-wrap: wrap; gap: 5px; }
.empty-text { color: var(--color-text-secondary); font-size: 14px; }

.pagination-wrapper { padding: 14px 20px; display: flex; justify-content: flex-end; border-top: 1px solid var(--color-border-light); }

// ========== 弹窗内容 ==========
.dialog-body-inner {
  padding: 0 4px;
}

// 信息提示条
.info-bar {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 14px; border-radius: 10px; margin-bottom: 16px;
  background: var(--color-surface-2); border: 1px solid var(--color-border-light);
}

.info-bar-icon {
  width: 36px; height: 36px; border-radius: 9px;
  background: var(--color-primary-light-9); color: var(--color-primary);
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
  &.shield { background: var(--color-warning-bg); color: var(--color-warning); }
}

.info-bar-text {
  display: flex; flex-direction: column; flex: 1;
}

.info-bar-title {
  font-size: 14px; font-weight: 700; color: var(--color-text-primary);
  font-family: 'Manrope', sans-serif;
}

.info-bar-sub {
  font-size: 12px; color: var(--color-text-muted); margin-top: 1px;
}

.perm-quick-actions {
  display: flex; gap: 4px;
  .el-button {
    font-size: 12px; color: var(--color-text-secondary);
    &:hover { color: var(--color-primary); }
  }
}

// 表单
:deep(.el-form-item) {
  margin-bottom: 14px;
  .el-form-item__label {
    font-weight: 600; color: var(--color-text-primary); font-size: 13px; margin-bottom: 6px;
    &::before { color: var(--color-danger); }
  }
}

:deep(.el-input__wrapper) {
  border-radius: 8px; box-shadow: none !important;
  border: 1.5px solid var(--color-border) !important;
  &:hover { border-color: var(--color-primary) !important; }
  &.is-focus { border-color: var(--color-primary) !important; box-shadow: 0 0 0 3px var(--color-primary-light-9) !important; }
}

:deep(.el-textarea__inner) {
  border-radius: 8px; box-shadow: none !important;
  border: 1.5px solid var(--color-border) !important;
  &:hover { border-color: var(--color-primary) !important; }
  &:focus { border-color: var(--color-primary) !important; box-shadow: 0 0 0 3px var(--color-primary-light-9) !important; }
}

:deep(.el-input-number) {
  .el-input__wrapper { border-radius: 8px; }
}

.form-tip {
  display: block; font-size: 11px; color: var(--color-text-muted); margin-top: 4px;
}

// ========== 权限分组 ==========
.perm-groups {
  display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 4px;
}

.perm-group-card {
  border: 1px solid var(--color-border-light); border-radius: 10px; overflow: hidden;
  transition: border-color 0.15s;
  &:hover { border-color: var(--color-info); }
}

.perm-group-header {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 12px; background: var(--color-surface-2); border-bottom: 1px solid var(--color-border-light);
}

.perm-group-icon {
  width: 26px; height: 26px; border-radius: 7px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}

.perm-group-name {
  font-size: 13px; font-weight: 700; color: var(--color-text-primary); flex: 1;
  font-family: 'Manrope', sans-serif;
}

.perm-group-count {
  font-size: 11px; color: var(--color-text-muted); background: var(--gray-100);
  padding: 2px 7px; border-radius: 10px; font-weight: 500;
}

.perm-group-items {
  padding: 8px 10px; display: flex; flex-direction: column; gap: 3px;
}

.perm-item {
  display: flex; align-items: center; gap: 8px;
  padding: 7px 8px; border-radius: 7px; cursor: pointer;
  transition: all 0.15s;
  border: 1px solid transparent;

  &:hover {
    background: var(--color-primary-light-9);
    border-color: var(--color-info);
  }

  &.perm-checked {
    background: var(--color-primary-light-9);
    border-color: var(--color-primary);
  }
}

.perm-checkbox {
  width: 18px; height: 18px; border-radius: 5px; flex-shrink: 0;
  border: 1.5px solid var(--color-border);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.15s;
  background: var(--color-surface);

  &.checked {
    background: var(--color-primary); border-color: var(--color-primary); color: var(--color-surface);
  }
}

.perm-item-info {
  display: flex; flex-direction: column; gap: 1px;
}

.perm-item-label {
  font-size: 13px; font-weight: 600; color: var(--color-text-primary);
}

.perm-item-value {
  font-size: 10px; color: var(--color-text-muted); font-family: 'SF Mono', monospace;
}

// ========== 弹窗底部 ==========
.dialog-footer {
  display: flex; justify-content: flex-end; gap: 10px;
  .el-button {
    min-width: 80px; border-radius: 8px; font-weight: 600;
  }
}

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
  background: rgba(124, 58, 237, 0.1);
  color: #7c3aed;
  border: 1px solid rgba(124, 58, 237, 0.2);
}

.dialog-title-text {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

:deep(.el-dialog__body) {
  max-height: 70vh; overflow-y: auto;
  &::-webkit-scrollbar { width: 5px; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 3px; }
}
</style>
