<template>
  <div class="page">
    <header class="page-header">
      <div>
        <h1 class="page-title">用户组</h1>
        <p class="page-subtitle">管理系统用户组</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          新增用户组
        </el-button>
      </div>
    </header>

    <div class="filter-card">
      <el-input v-model="keyword" placeholder="搜索用户组名称/编码" clearable @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

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
        <el-table-column prop="name" label="用户组名称" min-width="150" />
        <el-table-column prop="code" label="用户组编码" min-width="150">
          <template #default="{ row }">
            <code class="code-text">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="280" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" min-width="80" align="center" />
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'view', label: '详情', type: 'primary' },
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

    <!-- 侧边抽屉 -->
    <el-drawer v-model="drawerVisible" size="480px" direction="rtl">
      <template #header>
        <div class="drawer-head">
          <span class="drawer-mode-tag">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="drawer-title-text">{{ isEdit ? '编辑用户组' : '新增用户组' }}</span>
        </div>
      </template>
      <div class="drawer-content">
        <!-- 基本设置 -->
        <div class="form-section">
          <div class="section-header">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            <span>基本设置</span>
          </div>
          <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
            <el-form-item label="用户组名称" prop="name" class="is-required">
              <el-input v-model="form.name" placeholder="请输入用户组名称">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="用户组编码" prop="code" class="is-required">
              <el-input v-model="form.code" placeholder="请输入用户组编码，如：developers" :disabled="isEdit">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入用户组描述信息" />
            </el-form-item>
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
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

    <!-- 用户组详情抽屉 -->
    <el-drawer v-model="detailVisible" size="480px" direction="rtl">
      <template #header>
        <div class="drawer-head">
          <span class="drawer-mode-tag drawer-mode-tag--info">详情</span>
          <span class="drawer-title-text">用户组详情</span>
        </div>
      </template>
      <div class="detail-content" v-if="detailGroup">
        <!-- 概览 -->
        <div class="detail-overview">
          <div class="detail-avatar">{{ detailGroup.name.slice(0, 2).toUpperCase() }}</div>
          <div class="detail-info">
            <div class="detail-name">{{ detailGroup.name }}</div>
            <div class="detail-code">
              <code>{{ detailGroup.code }}</code>
            </div>
          </div>
        </div>

        <!-- 基本信息 -->
        <div class="detail-section">
          <div class="detail-section-label">基本信息</div>
          <div class="detail-row">
            <span class="detail-label">用户组ID</span>
            <span class="detail-value mono">{{ detailGroup.id }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">用户组名称</span>
            <span class="detail-value">{{ detailGroup.name }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">用户组编码</span>
            <span class="detail-value mono">{{ detailGroup.code }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">描述</span>
            <span class="detail-value">{{ detailGroup.description || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">排序权重</span>
            <span class="detail-value">{{ detailGroup.sort }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">创建时间</span>
            <span class="detail-value">{{ detailGroup.createdAt || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">更新时间</span>
            <span class="detail-value">{{ detailGroup.updatedAt || '-' }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="detailVisible = false">关闭</el-button>
          <el-button type="primary" @click="detailVisible = false; handleEdit(detailGroup!)">编辑用户组</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserGroupApi, type UserGroup, type CreateUserGroupReq, type UpdateUserGroupReq } from '@/api/user-group'
import TableActions from '@/components/TableActions.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<UserGroup[]>([])
const tableRef = ref()
const selectedRows = ref<UserGroup[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const detailVisible = ref(false)
const detailGroup = ref<UserGroup | null>(null)
const formRef = ref()
const keyword = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateUserGroupReq>({ name: '', code: '', description: '', sort: 0 })

const formRules = {
  name: [{ required: true, message: '请输入用户组名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入用户组编码', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await UserGroupApi.list({ page: pagination.page, pageSize: pagination.pageSize, keyword: keyword.value || undefined })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; handleSearch() }
const handleCreate = () => { isEdit.value = false; Object.assign(form, { name: '', code: '', description: '', sort: 0 }); drawerVisible.value = true }
const handleEdit = (row: UserGroup) => { isEdit.value = true; Object.assign(form, { id: row.id, name: row.name, code: row.code, description: row.description, sort: row.sort }); drawerVisible.value = true }
const handleView = (row: UserGroup) => { detailGroup.value = row; detailVisible.value = true }
const openDetailById = async (id: number) => {
  // API 获取详情或从列表中找到对应项
  const found = tableData.value.find(g => g.id === id)
  if (found) { detailGroup.value = found; detailVisible.value = true }
  else {
    // 尝试按ID加载
    try {
      const res = await UserGroupApi.getById(id)
      if (res.code === 200) { detailGroup.value = res.data; detailVisible.value = true }
    } catch {}
  }
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const res = await UserGroupApi.update({ id: form.id as number, name: form.name, code: form.code, description: form.description, sort: form.sort })
      if (res.code === 200) { ElMessage.success('更新成功'); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || '更新失败')
    } else {
      const res = await UserGroupApi.create(form)
      if (res.code === 200) { ElMessage.success('创建成功'); drawerVisible.value = false; loadData() }
      else ElMessage.error(res.message || '创建失败')
    }
  } finally { submitting.value = false }
}

const handleDelete = async (row: UserGroup) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户组"${row.name}"吗？`, '删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const res = await UserGroupApi.del(row.id)
    if (res.code === 200) { ElMessage.success('删除成功'); loadData() }
    else ElMessage.error(res.message || '删除失败')
  } catch {}
}

const handleAction = (key: string, row: UserGroup) => {
  if (key === 'view') handleView(row)
  else if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个用户组吗？`, '批量删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    const ids = selectedRows.value.map(row => row.id)
    await UserGroupApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个用户组`)
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') { ElMessage.error('批量删除失败') }
  }
}

onMounted(() => {
  loadData()
  // 从用户管理跳转过来的情况：带 id 参数时自动打开详情
  const idParam = route.query.id
  if (idParam) {
    const id = Number(idParam)
    if (!isNaN(id)) {
      // 等数据加载后再打开详情
      setTimeout(() => openDetailById(id), 300)
    }
  }
})
</script>

<style scoped lang="scss">
.page { min-height: 100vh; background: var(--color-page-bg); padding: 20px; }

.page-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 16px;
  padding: 20px 24px; background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
}

.page-title { font-family: 'Manrope', sans-serif; font-size: 20px; font-weight: 800; color: var(--color-text-primary); margin-bottom: 2px; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); }
.header-actions { display: flex; gap: 8px; }

.filter-card {
  display: flex; flex-wrap: wrap; gap: 10px; align-items: center;
  background: var(--color-surface); border-radius: 12px; padding: 16px 20px; margin-bottom: 14px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light);
  > * { flex: 0 0 auto; }
  .el-input { width: 200px; }
}

.table-card { background: var(--color-surface); border-radius: 12px; box-shadow: var(--shadow-xs); border: 1px solid var(--color-border-light); overflow: hidden; }

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

.toolbar-left {
  .record-count { font-size: 13px; color: var(--color-text-secondary); strong { color: var(--color-text-primary); font-weight: 700; } }
  .selection-count { margin-left: 16px; font-size: 13px; color: var(--color-danger); strong { font-weight: 700; } }
}

.toolbar-right { display: flex; gap: 10px; }

:deep(.el-table) {
  --el-table-border-color: var(--color-border-light); width: 100% !important;
  th.el-table__cell {
    background-color: var(--color-surface-2) !important; border-bottom: 1px solid var(--color-border-light) !important;
    padding: 14px 16px !important; color: var(--color-text-secondary) !important; font-weight: 700 !important; font-size: 12px !important;
  }
  td.el-table__cell { padding: 14px 16px !important; font-size: 14px; color: var(--color-text-primary); border-bottom: 1px solid var(--color-border-light) !important; }
  .el-table__body tr:hover > td.el-table__cell { background-color: var(--color-surface-2) !important; }
}

.code-text { font-size: 13px; color: var(--color-primary); background: var(--color-primary-light-9); padding: 4px 10px; border-radius: 6px; font-weight: 500; }

.pagination-wrapper { padding: 16px 20px; display: flex; justify-content: flex-end; border-top: 1px solid var(--color-border-light); }

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

  &--info {
    background: var(--color-success-light-9);
    color: var(--color-success);
    border-color: rgba(0, 176, 80, 0.2);
  }
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
  display: flex; align-items: center; gap: 8px;
  margin-bottom: 16px; padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-text-primary); font-weight: 600; font-size: 15px;
  svg { color: var(--color-text-secondary); }
}

/* ==================== 用户组详情抽屉 ==================== */
.detail-content {
  padding: 20px 24px;
  height: calc(100vh - 140px);
  overflow-y: auto;
}

.detail-overview {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--color-primary-light-9);
  border: 1px solid rgba(0,94,235,0.15);
  border-radius: var(--radius-md);
  margin-bottom: 24px;
}

.detail-avatar {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-md);
  background: var(--color-primary);
  color: #fff;
  font-family: 'Manrope', sans-serif;
  font-size: 18px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  letter-spacing: -0.5px;
}

.detail-info { flex: 1; min-width: 0; }
.detail-name { font-size: 17px; font-weight: 700; color: var(--color-text-primary); font-family: 'Manrope', sans-serif; margin-bottom: 4px; }
.detail-code {
  code {
    font-family: 'SF Mono', 'DM Mono', monospace;
    font-size: 12px;
    color: var(--color-primary);
    background: rgba(0,94,235,0.08);
    padding: 2px 8px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(0,94,235,0.15);
  }
}

.detail-section { margin-bottom: 20px; }
.detail-section-label {
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding-bottom: 8px;
  margin-bottom: 8px;
  border-bottom: 1px solid var(--color-border-light);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 7px 0;
  border-bottom: 1px solid var(--color-border-light);
  gap: 12px;
  &:last-child { border-bottom: none; }
}
.detail-label { font-size: 13px; color: var(--color-text-secondary); flex-shrink: 0; }
.detail-value { font-size: 13px; color: var(--color-text-primary); font-weight: 500; text-align: right; }
.detail-value.mono { font-family: 'SF Mono', 'DM Mono', monospace; font-size: 12px; color: var(--color-primary); }

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 14px 20px;
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface);
  .el-button { min-width: 80px; border-radius: var(--radius-sm); font-weight: 600; }
}

:deep(.el-form-item) {
  margin-bottom: 18px;
  .el-form-item__label {
    font-weight: 500; color: var(--color-text-primary); font-size: 14px;
    &::before { color: var(--color-danger); }
  }
}

:deep(.el-input) {
  .el-input__wrapper {
    border-radius: 8px; padding: 4px 12px;
    box-shadow: 0 0 0 1px var(--color-border) inset;
    &:hover { box-shadow: 0 0 0 1px var(--color-info) inset; }
    &.is-focus { box-shadow: 0 0 0 2px var(--color-primary-light-8) inset, 0 0 0 1px var(--color-primary) inset; }
  }
  .el-input__prefix { color: var(--color-text-secondary); margin-right: 8px; }
}

:deep(.el-input-number) {
  .el-input__wrapper { border-radius: 8px; }
}

:deep(.el-textarea) {
  .el-textarea__inner {
    border-radius: 8px;
    &:hover { box-shadow: 0 0 0 1px var(--color-info) inset; }
    &:focus { box-shadow: 0 0 0 2px var(--color-primary-light-8) inset, 0 0 0 1px var(--color-primary) inset; }
  }
}
</style>
