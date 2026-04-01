<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('role.list.title') }}</h1>
        <span class="page-subtitle">{{ t('role.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('role.list.create') }}
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input v-model="keyword" :placeholder="t('role.list.searchPlaceholder')" clearable @keyup.enter="handleSearch" style="width: 240px">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </template>
      </el-input>
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

    <!-- 内容卡片 -->
    <div class="content-card">
      <el-table ref="tableRef" v-model:selection="selectedRows" :data="tableData" v-loading="loading" stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="38" fixed="left" />
        <el-table-column prop="name" :label="t('role.list.table.name')" min-width="120">
          <template #default="{ row }">
            <div class="role-cell">
              <div class="role-avatar">{{ (row.name || 'R').charAt(0).toUpperCase() }}</div>
              <span class="role-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="code" :label="t('role.list.table.code')" min-width="120">
          <template #default="{ row }">
            <code class="code-text">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="t('role.list.table.description')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="permissions" :label="t('role.list.table.permissions')" min-width="260">
          <template #default="{ row }">
            <div class="perm-tags">
              <el-tag v-for="p in (row.permissions || []).slice(0, 3)" :key="p" type="info" size="small" effect="plain">{{ getPermLabel(p) }}</el-tag>
              <el-tag v-if="(row.permissions || []).length > 3" type="warning" size="small" effect="plain">+{{ row.permissions.length - 3 }}</el-tag>
              <span v-if="!row.permissions?.length" class="empty-text">—</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="t('common.sort')" width="70" align="center" />
        <el-table-column :label="t('common.actions')" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: t('common.edit'), type: 'primary' },
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

    <!-- 编辑/新增弹窗 -->
    <el-dialog v-model="dialogVisible" width="600px" destroy-on-close :close-on-click-modal="false">
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag">{{ isEdit ? t('common.edit') : t('common.create') }}</span>
          <span class="dialog-title-text">{{ isEdit ? t('role.list.form.editTitle') : t('role.list.form.createTitle') }}</span>
        </div>
      </template>

      <div class="dialog-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <div class="form-grid">
            <el-form-item :label="t('role.list.form.name')" prop="name">
              <el-input v-model="form.name" :placeholder="t('role.list.form.namePlaceholder')" size="small" />
            </el-form-item>
            <el-form-item :label="t('role.list.form.code')" prop="code">
              <el-input v-model="form.code" :placeholder="t('role.list.form.codePlaceholder')" :disabled="isEdit" size="small" />
            </el-form-item>
          </div>
          <el-form-item :label="t('role.list.form.description')" prop="description">
            <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('role.list.form.descriptionPlaceholder')" size="small" />
          </el-form-item>
          <el-form-item :label="t('common.sort')" prop="sort">
            <el-input-number v-model="form.sort" :min="0" :max="9999" size="small" />
          </el-form-item>
        </el-form>

        <!-- 权限设置 -->
        <div class="perm-section">
          <div class="perm-header">
            <span class="perm-title">{{ t('role.list.form.permConfig') }}</span>
            <div class="perm-actions">
              <el-button size="small" text @click="checkAllPerms">{{ t('role.list.form.selectAll') }}</el-button>
              <el-button size="small" text @click="clearAllPerms">{{ t('role.list.form.clearAll') }}</el-button>
            </div>
          </div>
          <div class="perm-groups">
            <div v-for="group in translatedPermGroups" :key="group.key" class="perm-group">
              <div class="perm-group-head">
                <span class="perm-group-name">{{ group.name }}</span>
                <span class="perm-group-count">{{ getGroupCheckedCount(group.perms) }}/{{ group.perms.length }}</span>
              </div>
              <div class="perm-items">
                <div v-for="perm in group.perms" :key="perm.value" class="perm-item" :class="{ active: form.permissions.includes(perm.value) }" @click="togglePerm(perm.value)">
                  <div class="perm-check">
                    <svg v-if="form.permissions.includes(perm.value)" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                  </div>
                  <span class="perm-label">{{ perm.label }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-foot">
          <el-button @click="dialogVisible = false">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            {{ t('common.cancel') }}
          </el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">
            <svg v-if="!submitting" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
            {{ t('common.save') }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, inject, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { RoleApi, type Role, type CreateRoleReq, type UpdateRoleReq } from '@/api/role'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()
const refreshUser = inject<() => Promise<any>>('refreshUser', () => Promise.resolve())

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
  name: [{ required: true, message: t('role.list.messages.nameRequired'), trigger: 'blur' }],
  code: [{ required: true, message: t('role.list.messages.codeRequired'), trigger: 'blur' }],
}

// 权限分组
const permGroups = [
  { key: 'upload', name: '上传记录', perms: [{ label: '创建', value: 'upload:create' }, { label: '读取', value: 'upload:read' }, { label: '更新', value: 'upload:update' }, { label: '删除', value: 'upload:delete' }, { label: '导出', value: 'upload:export' }] },
  { key: 'field', name: '字段配置', perms: [{ label: '创建', value: 'field-config:create' }, { label: '读取', value: 'field-config:read' }, { label: '更新', value: 'field-config:update' }, { label: '删除', value: 'field-config:delete' }] },
  { key: 'project', name: '项目管理', perms: [{ label: '创建', value: 'project:create' }, { label: '读取', value: 'project:read' }, { label: '更新', value: 'project:update' }, { label: '删除', value: 'project:delete' }] },
  { key: 'personnel', name: '人员管理', perms: [{ label: '创建', value: 'personnel:create' }, { label: '读取', value: 'personnel:read' }, { label: '更新', value: 'personnel:update' }, { label: '删除', value: 'personnel:delete' }, { label: '导出', value: 'personnel:export' }] },
  { key: 'user', name: '用户管理', perms: [{ label: '创建', value: 'user:create' }, { label: '读取', value: 'user:read' }, { label: '更新', value: 'user:update' }, { label: '删除', value: 'user:delete' }] },
  { key: 'role', name: '角色管理', perms: [{ label: '创建', value: 'role:create' }, { label: '读取', value: 'role:read' }, { label: '更新', value: 'role:update' }, { label: '删除', value: 'role:delete' }] },
  { key: 'audit', name: '审计日志', perms: [{ label: '操作日志', value: 'audit:operation:read' }, { label: '登录日志', value: 'audit:login:read' }] },
  { key: 'system', name: '系统配置', perms: [{ label: '读取', value: 'config:read' }, { label: '更新', value: 'config:update' }, { label: '管理员', value: 'admin:all' }] },
]

const permLabelMap: Record<string, string> = {}
permGroups.forEach(g => g.perms.forEach(p => { permLabelMap[p.value] = p.label }))

const translatedPermGroups = computed(() =>
  permGroups.map(g => ({ ...g, perms: g.perms.map(p => ({ ...p, label: t(`role.perms.${p.value.replace(':', '_')}`) || p.label })) }))
)

const getPermLabel = (code: string) => permLabelMap[code] || code.split(':').pop() || code
const getGroupCheckedCount = (perms: { value: string }[]) => perms.filter(p => form.permissions.includes(p.value)).length

const togglePerm = (value: string) => {
  const idx = form.permissions.indexOf(value)
  if (idx > -1) form.permissions.splice(idx, 1)
  else form.permissions.push(value)
}

const checkAllPerms = () => { form.permissions = permGroups.flatMap(g => g.perms.map(p => p.value)) }
const clearAllPerms = () => { form.permissions = [] }

const loadData = async () => {
  loading.value = true
  try {
    const res = await RoleApi.list({ page: pagination.page, pageSize: pagination.pageSize, keyword: keyword.value || undefined })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; pagination.page = 1; loadData() }
const handleSelectionChange = (rows: Role[]) => { selectedRows.value = rows }

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
      if (res.code === 200) { ElMessage.success(t('common.updateSuccess')); dialogVisible.value = false; loadData(); refreshUser() }
      else ElMessage.error(res.message || t('role.list.messages.updateFailed'))
    } else {
      const res = await RoleApi.create(form)
      if (res.code === 200) { ElMessage.success(t('common.createSuccess')); dialogVisible.value = false; loadData(); refreshUser() }
      else ElMessage.error(res.message || t('role.list.messages.createFailed'))
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(t('role.list.messages.deleteConfirm', { name: row.name }), t('common.confirm'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    const res = await RoleApi.del(row.id)
    if (res.code === 200) { ElMessage.success(t('common.deleteSuccess')); loadData() }
    else ElMessage.error(res.message || t('role.list.messages.deleteFailed'))
  } catch {}
}

const handleBatchDelete = async () => {
  if (!selectedRows.value.length) return
  try {
    await ElMessageBox.confirm(t('role.list.messages.batchDeleteConfirm', { count: selectedRows.value.length }), t('common.batchConfirm'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    await RoleApi.batchDelete(selectedRows.value.map(r => r.id))
    ElMessage.success(t('role.list.messages.batchDeleteSuccess', { count: selectedRows.value.length }))
    selectedRows.value = []; loadData()
  } catch {}
}

const handleAction = (key: string, row: Role) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => loadData())
</script>

<script lang="ts">
export default { name: 'RoleList' }
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

.header-actions {
  display: flex;
  gap: var(--space-2);
}

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

  td.el-table__cell {
    padding: 9px 12px !important;
    border-bottom: 1px solid var(--color-border-light) !important;
    color: var(--color-text-primary);
  }

  .el-table__body tr:hover > td.el-table__cell {
    background-color: var(--color-primary-light-9) !important;
  }
}

.role-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 12px;
  flex-shrink: 0;
}

.role-name {
  font-weight: 600;
  font-size: 13px;
}

.code-text {
  font-size: 12px;
  color: var(--color-primary);
  background: var(--color-primary-light-9);
  padding: 2px 8px;
  border-radius: 5px;
}

.perm-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.empty-text {
  color: var(--color-text-secondary);
  font-size: 13px;
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

.record-info {
  font-size: 12px;
  color: var(--color-text-secondary);
}

/* ==================== 弹窗 ==================== */
.dialog-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dialog-mode-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  background: rgba(124, 58, 237, 0.1);
  color: #7c3aed;
}

.dialog-title-text {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.dialog-body {
  max-height: 65vh;
  overflow-y: auto;
  padding-right: 4px;

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 var(--space-4);
}

:deep(.el-form-item) {
  margin-bottom: var(--space-3);
  .el-form-item__label {
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-secondary);
    margin-bottom: 4px;
  }
}

/* 权限区域 */
.perm-section {
  margin-top: var(--space-2);
}

.perm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.perm-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.perm-actions {
  display: flex;
  gap: 4px;
  .el-button { font-size: 12px; color: var(--color-text-secondary); }
}

.perm-groups {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
}

.perm-group {
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.perm-group-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--color-surface-2);
  border-bottom: 1px solid var(--color-border-light);
}

.perm-group-name {
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.perm-group-count {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--gray-100);
  padding: 1px 6px;
  border-radius: 10px;
}

.perm-items {
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.perm-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
  border: 1px solid transparent;

  &:hover {
    background: var(--color-primary-light-9);
  }

  &.active {
    background: var(--color-primary-light-9);
    border-color: var(--color-primary);
  }
}

.perm-check {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 1.5px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s;

  .perm-item.active & {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: #fff;
  }
}

.perm-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.dialog-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 1366px) {
  .page { padding: var(--space-3); gap: var(--space-2); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-bar { padding: var(--space-2) var(--space-3); gap: var(--space-2); }
}

@media (max-width: 1024px) {
  .form-grid { grid-template-columns: 1fr; }
  .perm-groups { grid-template-columns: 1fr; }
}
</style>
