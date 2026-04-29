<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">{{ t('fieldConfig.list.title') }}</h1>
        <p class="page-subtitle">{{ t('fieldConfig.list.subtitle') }}</p>
      </div>
      <el-button type="primary" @click="handleCreate">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        {{ t('fieldConfig.list.addField') }}
      </el-button>
    </header>

    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar__search">
        <div class="search-box">
          <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <input v-model="searchKeyword" class="search-input" :placeholder="t('fieldConfig.list.searchPlaceholder')" @keyup.enter="handleSearch" />
          <span class="search-kbd">↵</span>
        </div>
      </div>

      <div class="toolbar__actions">
        <template v-if="selectedRows.length > 0">
          <button class="action-btn action-btn--danger" @click="handleBatchDelete">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
            {{ t('fieldConfig.list.batchDelete') }} ({{ selectedRows.length }})
          </button>
        </template>
        <button class="action-btn action-btn--primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('fieldConfig.list.addField') }}
        </button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="content-card">
      <el-table
        ref="tableRef"
        v-model:selection="selectedRows"
        :data="tableData"
        v-loading="loading"
        stripe
      >
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column prop="name" :label="t('fieldConfig.list.fieldName')" min-width="140">
          <template #default="{ row }">
            <span class="field-name">{{ row.name }}</span>
            <span v-if="row.required" class="required-mark">*</span>
          </template>
        </el-table-column>
        <el-table-column prop="code" :label="t('fieldConfig.list.fieldCode')" min-width="150">
          <template #default="{ row }">
            <code class="field-code">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('fieldConfig.list.type')" min-width="100" align="center">
          <template #default="{ row }">
            <span class="type-badge">{{ getTypeText(row.type) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="options" :label="t('fieldConfig.list.options')" min-width="200">
          <template #default="{ row }">
            <span v-if="row.options && row.options.length" class="options-text">
              {{ row.options.join(', ') }}
            </span>
            <span v-else class="options-empty">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="placeholder" :label="t('fieldConfig.list.placeholder')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sort" :label="t('fieldConfig.list.sort')" min-width="80" align="center" />
        <el-table-column prop="enabled" :label="t('fieldConfig.list.status')" min-width="90" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="row.enabled ? 'status-badge--success' : 'status-badge--disabled'">
              {{ row.enabled ? t('fieldConfig.list.enabled') : t('fieldConfig.list.disabled') }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: t('common.edit'), type: 'primary' },
              { key: 'delete', label: t('common.delete'), type: 'danger' }
            ]" @action="(key) => handleAction(key, row)" />
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
        />
      </div>
    </div>

    <!-- 侧边抽屉 -->
    <el-drawer v-model="drawerVisible" size="520px" direction="rtl">
      <template #header>
        <div class="drawer-head">
          <span class="drawer-mode-tag">{{ isEdit ? t('common.edit') : t('fieldConfig.list.add') }}</span>
          <span class="drawer-title-text">{{ isEdit ? t('fieldConfig.list.editTitle') : t('fieldConfig.list.addTitle') }}</span>
        </div>
      </template>
      <div class="drawer-content">
        <!-- 基本设置 -->
        <div class="form-section">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3h7a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-7m0-18H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7m0-18v18"/></svg>
            <span>{{ t('fieldConfig.list.fieldSettings') }}</span>
          </div>
          <el-form :model="form" label-position="top" :rules="formRules" ref="formRef">
            <el-form-item :label="t('fieldConfig.list.fieldName')" prop="name" class="is-required">
              <el-input v-model="form.name" :placeholder="t('fieldConfig.list.fieldNamePlaceholder')" maxlength="64" size="small">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3h7a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-7m0-18H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7m0-18v18"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.fieldCode')" prop="code" class="is-required">
              <el-input v-model="form.code" :placeholder="t('fieldConfig.list.fieldCodePlaceholder')" maxlength="64" :disabled="isEdit" size="small">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.fieldType')" prop="type" class="is-required">
              <el-select v-model="form.type" style="width: 100%" size="small">
                <el-option :label="t('fieldConfig.list.typeText')" value="text" />
                <el-option :label="t('fieldConfig.list.typeNumber')" value="number" />
                <el-option :label="t('fieldConfig.list.typeSelect')" value="select" />
                <el-option :label="t('fieldConfig.list.typeMultiselect')" value="multiselect" />
                <el-option :label="t('fieldConfig.list.typeDate')" value="date" />
                <el-option :label="t('fieldConfig.list.typeDatetime')" value="datetime" />
                <el-option :label="t('fieldConfig.list.typeTextarea')" value="textarea" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.required')">
              <el-switch v-model="form.required" />
            </el-form-item>
          </el-form>
        </div>

        <!-- 选项配置 -->
        <div class="form-section" v-if="form.type === 'select' || form.type === 'multiselect'">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
            <span>{{ t('fieldConfig.list.optionsConfig') }}</span>
          </div>
          <el-form label-position="top">
            <el-form-item :label="t('fieldConfig.list.optionsList')">
              <el-input
                v-model="optionsInput"
                type="textarea"
                :rows="3"
                size="small"
                :placeholder="t('fieldConfig.list.optionsPlaceholder')"
              />
            </el-form-item>
          </el-form>
        </div>

        <!-- 界面配置 -->
        <div class="form-section">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
            <span>{{ t('fieldConfig.list.uiConfig') }}</span>
          </div>
          <el-form label-position="top">
            <el-form-item :label="t('fieldConfig.list.defaultValue')">
              <el-input v-model="form.defaultValue" :placeholder="t('fieldConfig.list.defaultValuePlaceholder')" maxlength="512" size="small">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.placeholder')">
              <el-input v-model="form.placeholder" :placeholder="t('fieldConfig.list.placeholderHint')" maxlength="256" size="small">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.sort')">
              <el-input-number v-model="form.sort" :min="0" :max="999" style="width: 100%" size="small" />
            </el-form-item>
            <el-form-item :label="t('fieldConfig.list.enabledStatus')">
              <el-switch v-model="form.enabled" />
              <span class="switch-hint">{{ form.enabled ? t('fieldConfig.list.enabledHint') : t('fieldConfig.list.disabledHint') }}</span>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="drawerVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">{{ t('common.save') }}</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { FieldConfigApi, type FieldConfig, type CreateFieldConfigReq } from '@/api/field-config'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const searchKeyword = ref('')
const tableData = ref<FieldConfig[]>([])
const tableRef = ref()
const selectedRows = ref<FieldConfig[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<CreateFieldConfigReq & { id?: number }>({
  name: '',
  code: '',
  type: 'text',
  required: false,
  options: [],
  defaultValue: '',
  placeholder: '',
  sort: 0,
  enabled: true
})

const optionsInput = ref('')

const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: t('fieldConfig.list.formRules.nameRequired'), trigger: 'blur' }],
  code: [
    { required: true, message: t('fieldConfig.list.formRules.codeRequired'), trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: t('fieldConfig.list.formRules.codePattern'), trigger: 'blur' }
  ],
  type: [{ required: true, message: t('fieldConfig.list.fieldType'), trigger: 'change' }]
}))

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    text: t('fieldConfig.list.typeText'),
    number: t('fieldConfig.list.typeNumber'),
    select: t('fieldConfig.list.typeSelect'),
    multiselect: t('fieldConfig.list.typeMultiselect'),
    date: t('fieldConfig.list.typeDate'),
    datetime: t('fieldConfig.list.typeDatetime'),
    textarea: t('fieldConfig.list.typeTextarea')
  }
  return map[type] || type
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await FieldConfigApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined
    })
    tableData.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (error) {
    console.error('Failed to load data:', error)
    ElMessage.error(t('fieldConfig.list.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  searchKeyword.value = ''
  pagination.page = 1
  loadData()
}

const resetForm = () => {
  form.name = ''
  form.code = ''
  form.type = 'text'
  form.required = false
  form.options = []
  form.defaultValue = ''
  form.placeholder = ''
  form.sort = 0
  form.enabled = true
  optionsInput.value = ''
}

const handleCreate = () => {
  isEdit.value = false
  delete form.id
  resetForm()
  drawerVisible.value = true
}

const handleEdit = (row: FieldConfig) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.code = row.code
  form.type = row.type
  form.required = row.required
  form.options = row.options || []
  form.defaultValue = row.defaultValue
  form.placeholder = row.placeholder
  form.sort = row.sort
  form.enabled = row.enabled
  optionsInput.value = row.options ? row.options.join(',') : ''
  drawerVisible.value = true
}

const confirmSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    if (form.type === 'select' && optionsInput.value) {
      form.options = optionsInput.value.split(',').map(s => s.trim()).filter(Boolean)
    } else {
      form.options = []
    }

    if (isEdit.value) {
      await FieldConfigApi.update(form as any)
      ElMessage.success(t(isEdit.value ? 'fieldConfig.list.updateSuccess' : 'fieldConfig.list.createSuccess'))
    } else {
      await FieldConfigApi.create(form)
      ElMessage.success(t('fieldConfig.list.createSuccess'))
    }
    drawerVisible.value = false
    loadData()
  } catch (error: any) {
    ElMessage.error(error?.message || t('fieldConfig.list.opFailed'))
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: FieldConfig) => {
  try {
    await ElMessageBox.confirm(t('fieldConfig.list.deleteConfirm', { name: row.name }), t('fieldConfig.list.deleteTitle'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await FieldConfigApi.del(row.id)
    ElMessage.success(t('fieldConfig.list.deleteSuccess'))
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('fieldConfig.list.deleteFailed'))
    }
  }
}

const handleAction = (key: string, row: FieldConfig) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(t('fieldConfig.list.batchDeleteConfirm', { count: selectedRows.value.length }), t('fieldConfig.list.batchDeleteTitle'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    const ids = selectedRows.value.map(row => row.id)
    await FieldConfigApi.batchDelete(ids)
    ElMessage.success(t('fieldConfig.list.batchDeleteSuccess', { count: selectedRows.value.length }))
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('fieldConfig.list.batchDeleteFailed'))
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: var(--color-page-bg);
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
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
  margin: 0;
  letter-spacing: -0.3px;
}

.page-subtitle {
  font-size: 12px;
  color: var(--color-text-muted);
}

/* ==================== 工具栏 ==================== */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  gap: var(--space-3);
  flex-wrap: wrap;
}

.toolbar__search {
  flex: 1;
  min-width: 200px;
  max-width: 360px;
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 0 10px;
  gap: 8px;
  transition: all 0.2s;

  &:focus-within {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.08);
  }
}

.search-icon { color: var(--color-text-muted); flex-shrink: 0; }

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 13px;
  color: var(--color-text-primary);
  padding: 7px 0;
  min-width: 0;

  &::placeholder { color: var(--color-text-muted); }
}

.search-kbd {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--gray-100);
  padding: 1px 6px;
  border-radius: 4px;
  font-family: monospace;
  border: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.toolbar__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-shrink: 0;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 7px 14px;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 600;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
  background: transparent;
  color: var(--color-text-secondary);
  border-color: var(--color-border);

  svg { flex-shrink: 0; }

  &:hover {
    background: var(--color-surface-2);
    border-color: var(--color-border);
    color: var(--color-text-primary);
  }

  &--primary {
    background: var(--color-primary);
    color: #fff;
    border-color: var(--color-primary);
    &:hover { background: #005eea; border-color: #005eea; }
    svg { stroke: #fff; }
  }

  &--danger {
    background: var(--color-danger);
    color: #fff;
    border-color: var(--color-danger);
    &:hover { background: #dc2626; border-color: #dc2626; }
    svg { stroke: #fff; }
  }

  &--ghost {
    background: transparent;
    color: var(--color-text-secondary);
    border-color: transparent;
    &:hover { background: var(--color-surface-2); color: var(--color-text-primary); border-color: var(--color-border); }
  }
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
  --el-table-border-color: var(--color-border-light);
  width: 100% !important;

  th.el-table__cell {
    background-color: var(--color-surface-2) !important;
    border-bottom: 1px solid var(--color-border-light) !important;
    padding: 14px 16px !important;
    color: var(--color-text-secondary) !important;
    font-weight: 700 !important;
    font-size: 12px !important;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  td.el-table__cell {
    border-bottom: 1px solid var(--color-border-light) !important;
    padding: 14px 16px !important;
    font-size: 14px;
    color: var(--color-text-primary);
  }

  .el-table__body tr:hover > td.el-table__cell {
    background-color: var(--color-surface-2) !important;
  }
}

.field-name {
  font-weight: 600;
  color: var(--color-text-primary);
  font-size: 14px;
}

.required-mark {
  color: var(--color-danger);
  margin-left: 2px;
  font-weight: bold;
}

.field-code {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  background: var(--color-primary-light-9);
  padding: 4px 10px;
  border-radius: 6px;
  color: var(--color-primary);
  font-weight: 500;
}

.type-badge {
  display: inline-block;
  padding: 4px 12px;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
}

.options-text {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.options-empty {
  color: var(--color-text-muted);
  font-size: 14px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;

  &--success { background: var(--color-success-bg); color: var(--color-success); }
  &--disabled { background: var(--gray-100); color: var(--color-text-secondary); }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--color-border-light);
}

// 抽屉样式
:deep(.el-drawer) {
  .el-drawer__header {
    padding: 16px 20px;
    margin-bottom: 0;
    border-bottom: 1px solid var(--color-border-light);
    color: var(--color-text-primary);
    font-weight: 700;
    font-size: 16px;
  }
  .el-drawer__body { padding: 0; }
}

.drawer-head {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.drawer-mode-tag {
  font-size: 10px;
  font-weight: 800;
  font-family: 'DM Sans', sans-serif;
  padding: 2px 8px;
  border-radius: 4px;
  letter-spacing: 0.5px;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border: 1px solid rgba(0, 94, 235, 0.2);
}

.drawer-title-text {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.drawer-content {
  padding: 16px 20px;
  height: calc(100vh - 140px);
  overflow-y: auto;
}

.form-section {
  margin-bottom: 28px;
  &:last-child { margin-bottom: 0; }
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-text-primary);
  font-weight: 600;
  font-size: 13px;
  svg { color: var(--color-text-secondary); }
}

:deep(.el-form-item) {
  margin-bottom: 18px;
  .el-form-item__label {
    font-weight: 500;
    color: var(--color-text-primary);
    font-size: 12px;
    &::before { color: var(--color-danger); }
  }
}

:deep(.el-input) {
  .el-input__wrapper {
    border-radius: 8px;
    padding: 4px 12px;
    box-shadow: 0 0 0 1px var(--color-border) inset;
    &:hover { box-shadow: 0 0 0 1px var(--color-info) inset; }
    &.is-focus { box-shadow: 0 0 0 2px var(--color-primary-light-8) inset, 0 0 0 1px var(--color-primary) inset; }
  }
  .el-input__prefix { color: var(--color-text-secondary); margin-right: 8px; }
}

:deep(.el-select) {
  .el-select__wrapper {
    border-radius: 8px;
    min-height: 36px;
    box-shadow: 0 0 0 1px var(--color-border) inset;
    &:hover { box-shadow: 0 0 0 1px var(--color-info) inset; }
  }
}

:deep(.el-textarea) {
  .el-textarea__inner {
    border-radius: 8px;
    &:hover { box-shadow: 0 0 0 1px var(--color-info) inset; }
    &:focus { box-shadow: 0 0 0 2px var(--color-primary-light-8) inset, 0 0 0 1px var(--color-primary) inset; }
  }
}

:deep(.el-input-number) {
  .el-input__wrapper { border-radius: 8px; }
}

.switch-hint {
  margin-left: 12px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 24px;
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface);
}
</style>
