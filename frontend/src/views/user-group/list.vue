<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('userGroup.list.title') }}</h1>
        <span class="page-subtitle">{{ t('userGroup.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('userGroup.list.create') }}
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input v-model="keyword" :placeholder="t('userGroup.list.searchPlaceholder')" clearable @keyup.enter="handleSearch" style="width: 240px">
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
        <el-table-column prop="name" :label="t('userGroup.list.table.name')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" :label="t('userGroup.list.table.code')" min-width="150">
          <template #default="{ row }">
            <code class="code-text">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="t('userGroup.list.table.description')" min-width="240" show-overflow-tooltip />
        <el-table-column prop="roleName" :label="t('userGroup.list.table.role')" min-width="140">
          <template #default="{ row }">
            <span v-if="row.roleName" class="role-badge">{{ row.roleName }}</span>
            <span v-else class="empty-text">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="t('common.sort')" width="70" align="center" />
        <el-table-column :label="t('common.actions')" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'view', label: t('common.view'), type: 'default' },
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
          small
        />
      </div>
    </div>

    <!-- 编辑/新增抽屉 -->
    <el-drawer v-model="drawerVisible" direction="rtl" size="460px" :destroy-on-close="true" class="personnel-drawer">
      <template #header>
        <div class="drawer-title-inner">
          <span class="drawer-mode-tag" :class="isEdit ? 'tag--edit' : 'tag--new'">{{ isEdit ? t('common.edit') : t('common.create') }}</span>
          <span class="drawer-title-text">{{ isEdit ? t('userGroup.list.drawer.editTitle') : t('userGroup.list.drawer.createTitle') }}</span>
        </div>
      </template>
      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top" class="edit-form">
          <el-form-item :label="t('userGroup.list.drawer.name')" prop="name">
            <el-input v-model="form.name" :placeholder="t('userGroup.list.drawer.namePlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('userGroup.list.drawer.code')" prop="code">
            <el-input v-model="form.code" :placeholder="t('userGroup.list.drawer.codePlaceholder')" :disabled="isEdit" />
          </el-form-item>
          <el-form-item :label="t('userGroup.list.drawer.description')">
            <el-input v-model="form.description" type="textarea" :rows="3" :placeholder="t('userGroup.list.drawer.descriptionPlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('userGroup.list.drawer.role')">
            <el-select v-model="form.roleId" clearable style="width: 100%">
              <el-option v-for="r in roleList" :key="r.id" :label="r.name" :value="r.id" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('common.sort')">
            <el-input-number v-model="form.sort" :min="0" :max="9999" style="width: 100%" />
          </el-form-item>
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

    <!-- 详情抽屉 -->
    <el-drawer v-model="detailVisible" direction="rtl" size="440px" class="personnel-drawer">
      <template #header>
        <div class="drawer-title-inner">
          <span class="drawer-mode-tag tag--view">{{ t('common.view') }}</span>
          <span class="drawer-title-text">{{ detailGroup?.name }}</span>
        </div>
      </template>
      <div class="drawer-body" v-if="detailGroup">
        <div class="detail-info">
          <div class="detail-name">{{ detailGroup.name }}</div>
          <code class="detail-code">{{ detailGroup.code }}</code>
        </div>
        <div class="detail-rows">
          <div class="detail-row">
            <span class="detail-label">{{ t('userGroup.list.detail.groupId') }}</span>
            <span class="detail-value">{{ detailGroup.id }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('userGroup.list.detail.description') }}</span>
            <span class="detail-value" :class="detailGroup.description ? '' : 'empty-text'">{{ detailGroup.description || '—' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('common.sort') }}</span>
            <span class="detail-value">{{ detailGroup.sort }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('userGroup.list.table.role') }}</span>
            <span class="detail-value" :class="detailGroup.roleName ? '' : 'empty-text'">{{ detailGroup.roleName || '—' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('common.createdAt') }}</span>
            <span class="detail-value">{{ formatDate(detailGroup.createdAt) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('common.updatedAt') }}</span>
            <span class="detail-value">{{ formatDate(detailGroup.updatedAt) }}</span>
          </div>
        </div>
      </div>
      <div class="drawer-foot">
        <el-button @click="detailVisible = false">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ t('common.close') }}
        </el-button>
        <el-button type="primary" @click="detailVisible = false; handleEdit(detailGroup!)">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          {{ t('common.edit') }}
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserGroupApi, type UserGroup, type CreateUserGroupReq, type UpdateUserGroupReq } from '@/api/user-group'
import { RoleApi } from '@/api/role'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<UserGroup[]>([])
const roleList = ref<{id: number, name: string}[]>([])
const tableRef = ref()
const selectedRows = ref<UserGroup[]>([])
const drawerVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const detailGroup = ref<UserGroup | null>(null)
const formRef = ref()
const keyword = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateUserGroupReq & {roleId?: number}>({ name: '', code: '', description: '', roleId: undefined, sort: 0 })

const formRules = {
  name: [{ required: true, message: t('userGroup.list.messages.nameRequired'), trigger: 'blur' }],
  code: [{ required: true, message: t('userGroup.list.messages.codeRequired'), trigger: 'blur' }],
}

// 格式化时间 YYYY-MM-DD HH:mm:ss
const formatDate = (dateStr: string | undefined | null) => {
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

const loadData = async () => {
  loading.value = true
  try {
    const res = await UserGroupApi.list({ page: pagination.page, pageSize: pagination.pageSize, keyword: keyword.value || undefined })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const loadRoles = async () => {
  const res = await RoleApi.list({ page: 1, pageSize: 100 })
  if (res.code === 200) { roleList.value = res.data.items || [] }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; pagination.page = 1; loadData() }
const handleSelectionChange = (rows: UserGroup[]) => { selectedRows.value = rows }

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', code: '', description: '', roleId: undefined, sort: 0 })
  drawerVisible.value = true
}

const handleAction = (key: string, row: UserGroup) => {
  if (key === 'view') handleView(row)
  else if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleEdit = (row: UserGroup) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, name: row.name, code: row.code, description: row.description, roleId: row.roleId || undefined, sort: row.sort })
  drawerVisible.value = true
}

const handleView = (row: UserGroup) => {
  detailGroup.value = row
  detailVisible.value = true
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const res = await UserGroupApi.update({ id: form.id as number, name: form.name, code: form.code, description: form.description, roleId: form.roleId || 0, sort: form.sort })
      if (res.code === 200) { ElMessage.success(t('common.updateSuccess')); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || t('userGroup.list.messages.updateFailed'))
    } else {
      const res = await UserGroupApi.create(form)
      if (res.code === 200) { ElMessage.success(t('common.createSuccess')); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || t('userGroup.list.messages.createFailed'))
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: UserGroup) => {
  try {
    await ElMessageBox.confirm(t('userGroup.list.messages.deleteConfirm', { name: row.name }), t('common.confirm'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    const res = await UserGroupApi.del(row.id)
    if (res.code === 200) { ElMessage.success(t('common.deleteSuccess')); loadData() }
    else ElMessage.error(res.message || t('common.deleteError'))
  } catch {}
}

const handleBatchDelete = async () => {
  if (!selectedRows.value.length) return
  try {
    await ElMessageBox.confirm(t('userGroup.list.messages.batchDeleteConfirm', { count: selectedRows.value.length }), t('common.batchConfirm'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    await UserGroupApi.batchDelete(selectedRows.value.map(r => r.id))
    ElMessage.success(t('userGroup.list.messages.batchDeleteSuccess', { count: selectedRows.value.length }))
    selectedRows.value = []; loadData()
  } catch {}
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => { loadData(); loadRoles() })
</script>

<script lang="ts">
export default { name: 'UserGroupList' }
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

.code-text {
  font-size: 12px;
  color: var(--color-primary);
  background: var(--color-primary-light-9);
  padding: 2px 8px;
  border-radius: 5px;
}

.role-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  background: rgba(0, 94, 235, 0.08);
  color: var(--color-primary);
  border: 1px solid rgba(0, 94, 235, 0.15);
}

.empty-text { color: var(--el-text-color-placeholder); }

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

/* ==================== 侧边栏 ==================== */
:deep(.personnel-drawer) {
  .el-drawer__header {
    padding: 10px 16px;
    margin-bottom: 0;
    border-bottom: 1px solid var(--color-border-light);
    background: var(--color-surface);
    align-items: center;
  }

  .el-drawer__body {
    padding: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}

.drawer-title-inner {
  display: flex;
  align-items: center;
  gap: 8px;
}

.drawer-mode-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: var(--radius-full);

  &.tag--edit { background: rgba(0, 94, 235, 0.1); color: var(--color-primary); }
  &.tag--new { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
  &.tag--view { background: rgba(0, 94, 235, 0.1); color: var(--color-primary); }
}

.drawer-title-text {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--color-surface-2);

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
}

.drawer-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

/* ==================== 表单 ==================== */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  animation: form-in 0.3s ease both 0.05s;
}

@keyframes form-in {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.edit-form :deep(.el-form-item) {
  margin-bottom: 6px;
  .el-form-item__label {
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-secondary);
    margin-bottom: 4px;
  }
  .el-form-item__error {
    padding-top: 2px;
  }
}

/* 详情页 */
.detail-info {
  background: var(--color-primary-light-9);
  border: 1px solid rgba(0, 94, 235, 0.15);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-bottom: 16px;
  text-align: center;
}

.detail-name {
  font-size: 17px;
  font-weight: 800;
  color: var(--color-text-primary);
  font-family: 'Manrope', sans-serif;
  margin-bottom: 6px;
}

.detail-code {
  font-size: 12px;
  color: var(--color-primary);
  background: rgba(0, 94, 235, 0.08);
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(0, 94, 235, 0.15);
}

.detail-rows {
  display: flex;
  flex-direction: column;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-border-light);
  gap: 12px;
  &:last-child { border-bottom: none; }
}

.detail-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.detail-value {
  font-size: 12px;
  color: var(--color-text-primary);
  font-weight: 500;
  text-align: right;
  word-break: break-all;
}

@media (max-width: 1366px) {
  .page { padding: var(--space-3); gap: var(--space-2); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-bar { padding: var(--space-2) var(--space-3); gap: var(--space-2); }
}
</style>
