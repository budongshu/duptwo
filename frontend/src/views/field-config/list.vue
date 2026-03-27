<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">字段配置</h1>
        <p class="page-subtitle">管理上传记录的动态字段</p>
      </div>
      <el-button type="primary" @click="handleCreate">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        新增字段
      </el-button>
    </header>

    <!-- 搜索栏 -->
    <div class="filter-card">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索字段名称/编码"
        clearable
        @keyup.enter="handleSearch"
      />
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

      <el-table
        ref="tableRef"
        v-model:selection="selectedRows"
        :data="tableData"
        v-loading="loading"
        stripe
      >
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column prop="name" label="字段名称" min-width="140">
          <template #default="{ row }">
            <span class="field-name">{{ row.name }}</span>
            <span v-if="row.required" class="required-mark">*</span>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="字段编码" min-width="150">
          <template #default="{ row }">
            <code class="field-code">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" min-width="100" align="center">
          <template #default="{ row }">
            <span class="type-badge">{{ getTypeText(row.type) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="options" label="选项" min-width="200">
          <template #default="{ row }">
            <span v-if="row.options && row.options.length" class="options-text">
              {{ row.options.join(', ') }}
            </span>
            <span v-else class="options-empty">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="placeholder" label="占位提示" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" min-width="80" align="center" />
        <el-table-column prop="enabled" label="状态" min-width="90" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="row.enabled ? 'status-badge--success' : 'status-badge--disabled'">
              {{ row.enabled ? '启用' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: '编辑', type: 'primary' },
              { key: 'delete', label: '删除', type: 'danger' }
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
          <span class="drawer-mode-tag">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="drawer-title-text">{{ isEdit ? '编辑字段' : '新增字段' }}</span>
        </div>
      </template>
      <div class="drawer-content">
        <!-- 基本设置 -->
        <div class="form-section">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3h7a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-7m0-18H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7m0-18v18"/></svg>
            <span>字段设置</span>
          </div>
          <el-form :model="form" label-position="top" :rules="formRules" ref="formRef">
            <el-form-item label="字段名称" prop="name" class="is-required">
              <el-input v-model="form.name" placeholder="如：目标路径" maxlength="64">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3h7a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-7m0-18H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7m0-18v18"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="字段编码" prop="code" class="is-required">
              <el-input v-model="form.code" placeholder="如：targetPath，只能是英文和数字" maxlength="64" :disabled="isEdit">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="字段类型" prop="type" class="is-required">
              <el-select v-model="form.type" style="width: 100%">
                <el-option label="文本" value="text" />
                <el-option label="数字" value="number" />
                <el-option label="下拉选择" value="select" />
                <el-option label="多选" value="multiselect" />
                <el-option label="日期" value="date" />
                <el-option label="日期时间" value="datetime" />
                <el-option label="文本域" value="textarea" />
              </el-select>
            </el-form-item>
            <el-form-item label="必填">
              <el-switch v-model="form.required" />
            </el-form-item>
          </el-form>
        </div>

        <!-- 选项配置 -->
        <div class="form-section" v-if="form.type === 'select' || form.type === 'multiselect'">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
            <span>选项配置</span>
          </div>
          <el-form label-position="top">
            <el-form-item label="选项列表">
              <el-input
                v-model="optionsInput"
                type="textarea"
                :rows="3"
                placeholder="多个选项用英文逗号分隔，如：选项1,选项2,选项3"
              />
            </el-form-item>
          </el-form>
        </div>

        <!-- 界面配置 -->
        <div class="form-section">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
            <span>界面配置</span>
          </div>
          <el-form label-position="top">
            <el-form-item label="默认值">
              <el-input v-model="form.defaultValue" placeholder="默认填充的值" maxlength="512">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="占位提示">
              <el-input v-model="form.placeholder" placeholder="如：请输入目标路径" maxlength="256">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" :min="0" :max="999" style="width: 100%" />
            </el-form-item>
            <el-form-item label="启用状态">
              <el-switch v-model="form.enabled" />
              <span class="switch-hint">{{ form.enabled ? '启用后该字段将在表单中显示' : '禁用后该字段将在表单中隐藏' }}</span>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="drawerVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">保存</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { FieldConfigApi, type FieldConfig, type CreateFieldConfigReq } from '@/api/field-config'
import TableActions from '@/components/TableActions.vue'

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

const formRules: FormRules = {
  name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入字段编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '编码必须以字母开头，只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择字段类型', trigger: 'change' }]
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    text: '文本',
    number: '数字',
    select: '下拉',
    multiselect: '多选',
    date: '日期',
    datetime: '日期时间',
    textarea: '文本域'
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
    ElMessage.error('加载数据失败')
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
      ElMessage.success('更新成功')
    } else {
      await FieldConfigApi.create(form)
      ElMessage.success('创建成功')
    }
    drawerVisible.value = false
    loadData()
  } catch (error: any) {
    ElMessage.error(error?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: FieldConfig) => {
  try {
    await ElMessageBox.confirm(`确定要删除字段"${row.name}"吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await FieldConfigApi.del(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
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
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个字段吗？`, '批量删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const ids = selectedRows.value.map(row => row.id)
    await FieldConfigApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个字段`)
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
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
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 16px;
  padding: 20px 24px;
  background: var(--color-surface);
  border-radius: 12px;
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.page-title {
  font-family: 'Manrope', sans-serif;
  font-size: 20px;
  font-weight: 800;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.page-subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.filter-card {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  background: var(--color-surface);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 14px;
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);

  > * { flex: 0 0 auto; }
  .el-input { width: 200px; }
}

.table-card {
  background: var(--color-surface);
  border-radius: 12px;
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

.toolbar-left {
  .record-count {
    font-size: 13px;
    color: var(--color-text-secondary);
    strong { color: var(--color-text-primary); font-weight: 700; }
  }
  .selection-count {
    margin-left: 16px;
    font-size: 13px;
    color: var(--color-danger);
    strong { font-weight: 700; }
  }
}

.toolbar-right {
  display: flex;
  gap: 10px;
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
  padding: 20px 24px;
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
  font-size: 15px;
  svg { color: var(--color-text-secondary); }
}

:deep(.el-form-item) {
  margin-bottom: 18px;
  .el-form-item__label {
    font-weight: 500;
    color: var(--color-text-primary);
    font-size: 14px;
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
