<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">人员管理</h1>
        <span class="page-subtitle">项目人员库</span>
      </div>
      <div class="header-actions">
        <el-button type="success" size="small" @click="handleExport" :loading="exporting">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          导出
        </el-button>
        <el-button type="primary" size="small" @click="handleCreate">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          新增
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-bar__search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索姓名/手机/公司/职位/驻场地点"
          clearable
          size="small"
          @keyup.enter="handleSearch"
          style="width: 260px"
        />
      </div>
      <div class="filter-bar__selects">
        <el-select v-model="searchStatus" placeholder="状态" clearable size="small" style="width: 110px">
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="inactive" />
        </el-select>
        <el-select v-model="searchOnProject" placeholder="在项状态" clearable size="small" style="width: 110px">
          <el-option label="在项" value="在项" />
          <el-option label="离项" value="离项" />
        </el-select>
      </div>
      <div class="filter-bar__actions">
        <el-button type="primary" size="small" @click="handleSearch">查询</el-button>
        <el-button size="small" @click="handleReset">重置</el-button>
        <el-button v-if="selectedRows.length > 0" type="danger" size="small" plain @click="handleBatchDelete">
          批量删除 ({{ selectedRows.length }})
        </el-button>
      </div>
    </div>

    <!-- 表格卡片 -->
    <div class="content-card">
      <el-table
        ref="tableRef"
        v-model:selection="selectedRows"
        :data="tableData"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="38" fixed="left" />
        <el-table-column prop="name" label="姓名" min-width="80" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机号" min-width="110" />
        <el-table-column prop="company" label="所属公司" min-width="140" show-overflow-tooltip />
        <el-table-column prop="position" label="职位" min-width="90" show-overflow-tooltip />
        <el-table-column prop="workExperience" label="工作经验" min-width="80" align="center" />
        <el-table-column prop="entryDate" label="入项时间" min-width="100" align="center" />
        <el-table-column prop="projectStartDate" label="立项时间" min-width="100" align="center" />
        <el-table-column prop="onProjectStatus" label="在项状态" min-width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.onProjectStatus === '在项' ? 'success' : 'warning'" size="small" effect="light">
              {{ row.onProjectStatus || '离项' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="salary" label="薪资" min-width="80" align="center" />
        <el-table-column prop="location" label="驻场地点" min-width="110" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" min-width="68" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" min-width="60" align="center" />
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'edit', label: '编辑', type: 'primary' },
              { key: 'delete', label: '删除', type: 'danger' }
            ]" @action="(key) => handleAction(key, row)" />
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <span class="record-info">共 <strong>{{ pagination.total }}</strong> 条</span>
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

    <!-- 编辑/新增侧边栏 -->
    <el-drawer
      v-model="drawerVisible"
      direction="rtl"
      size="440px"
      :with-header="true"
      :destroy-on-close="true"
      class="personnel-drawer"
    >
      <template #header>
        <div class="drawer-title-inner">
          <span class="drawer-mode-tag" :class="isEdit ? 'tag--edit' : 'tag--new'">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="drawer-title-text">{{ isEdit ? form.name || '人员' : '新增人员' }}</span>
        </div>
      </template>

      <!-- 侧边栏内容 -->
      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top" class="edit-form">

          <!-- 用户关联选择 -->
          <div class="user-select-section">
            <div class="user-select-header">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              <span>关联系统用户</span>
            </div>
            <el-select
              v-model="selectedUserId"
              placeholder="从系统用户中选择（可选）"
              clearable
              filterable
              :loading="userSelectLoading"
              size="small"
              style="width: 100%"
            >
              <el-option
                v-for="user in allUsers"
                :key="user.id"
                :label="`${user.nickname || user.username}（${user.email}）`"
                :value="user.id"
              >
                <div class="user-option">
                  <span class="user-option-name">{{ user.nickname || user.username }}</span>
                  <span class="user-option-email">{{ user.email }}</span>
                </div>
              </el-option>
            </el-select>
            <div class="user-select-tip">选择已有用户可自动填充姓名和邮箱，也可直接手动输入</div>
          </div>

          <div class="form-divider"></div>

          <div class="form-grid">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="form.name" placeholder="请输入姓名" size="small" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号" size="small" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" size="small" />
            </el-form-item>
            <el-form-item label="所属公司" prop="company">
              <el-input v-model="form.company" placeholder="请输入所属公司" size="small" />
            </el-form-item>
            <el-form-item label="职位" prop="position">
              <el-input v-model="form.position" placeholder="请输入职位" size="small" />
            </el-form-item>
            <el-form-item label="工作经验" prop="workExperience">
              <el-input v-model="form.workExperience" placeholder="如：3年" size="small" />
            </el-form-item>
            <el-form-item label="入项时间" prop="entryDate">
              <el-date-picker v-model="form.entryDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" size="small" />
            </el-form-item>
            <el-form-item label="立项时间" prop="projectStartDate">
              <el-date-picker v-model="form.projectStartDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" size="small" />
            </el-form-item>
            <el-form-item label="在项状态" prop="onProjectStatus">
              <el-select v-model="form.onProjectStatus" placeholder="请选择在项状态" size="small" style="width: 100%">
                <el-option label="在项" value="在项" />
                <el-option label="离项" value="离项" />
              </el-select>
            </el-form-item>
            <el-form-item label="薪资" prop="salary">
              <el-input v-model="form.salary" placeholder="如：15-20K" size="small" />
            </el-form-item>
          </div>
          <el-form-item label="驻场地点" prop="location">
            <el-input v-model="form.location" placeholder="请输入人员驻场地点" size="small" />
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="请输入备注（可选）" size="small" />
          </el-form-item>
          <div class="form-row-2">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status" size="small">
                <el-radio value="active">启用</el-radio>
                <el-radio value="inactive">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="9999" size="small" />
              <span class="form-hint">数字越小越靠前</span>
            </el-form-item>
          </div>
        </el-form>
      </div>

      <!-- 侧边栏底部 -->
      <div class="drawer-foot">
        <el-button size="small" @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" size="small" :loading="submitting" @click="confirmSubmit">
          <svg v-if="!submitting" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, inject } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { PersonnelApi, type Personnel, type CreatePersonnelReq, type UpdatePersonnelReq } from '@/api/personnel'
import { UserApi, type User } from '@/api/user'
import TableActions from '@/components/TableActions.vue'

const trackExport = inject<(action?: string) => void>('trackExport')

const loading = ref(false)
const exporting = ref(false)
const submitting = ref(false)
const tableData = ref<Personnel[]>([])
const tableRef = ref()
const selectedRows = ref<Personnel[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 用户选择相关
const allUsers = ref<User[]>([])
const selectedUserId = ref<number | undefined>()
const userSelectLoading = ref(false)

const searchKeyword = ref('')
const searchStatus = ref('')
const searchOnProject = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreatePersonnelReq & { id?: number }>({
  name: '', phone: '', email: '', company: '', position: '',
  workExperience: '', entryDate: '', projectStartDate: '',
  onProjectStatus: '离项', salary: '', location: '',
  remark: '', status: 'active', sort: 0,
})

const formRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
}

// 加载可选用户列表（用于关联选择）
const loadUsers = async () => {
  userSelectLoading.value = true
  try {
    const res = await UserApi.getAll()
    if (res.code === 200) {
      allUsers.value = res.data || []
    }
  } finally {
    userSelectLoading.value = false
  }
}

// 关联用户变化时，自动填充姓名和邮箱
watch(selectedUserId, (uid) => {
  if (!uid) return
  const user = allUsers.value.find(u => u.id === uid)
  if (user) {
    form.name = user.nickname || user.username
    form.email = user.email || ''
  }
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await PersonnelApi.list({
      page: pagination.page, pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
      onProject: searchOnProject.value || undefined,
    })
    if (res.code === 200) {
      tableData.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { searchKeyword.value = ''; searchStatus.value = ''; searchOnProject.value = ''; pagination.page = 1; loadData() }
const handleSelectionChange = (rows: Personnel[]) => { selectedRows.value = rows }

const handleExport = async () => {
  try {
    exporting.value = true
    const blob = await PersonnelApi.exportExcel({
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
      onProject: searchOnProject.value || undefined,
    })
    const url = URL.createObjectURL(blob as Blob)
    const link = document.createElement('a')
    link.href = url; link.download = `人员列表_${new Date().getTime()}.xlsx`
    document.body.appendChild(link); link.click(); document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
    trackExport?.('personnel')
  } catch (e: any) { ElMessage.error(e.message || '导出失败') }
  finally { exporting.value = false }
}

const handleCreate = async () => {
  isEdit.value = false
  selectedUserId.value = undefined
  await loadUsers()
  Object.assign(form, { id: undefined, name: '', phone: '', email: '', company: '', position: '', workExperience: '', entryDate: '', projectStartDate: '', onProjectStatus: '离项', salary: '', location: '', remark: '', status: 'active', sort: 0 })
  drawerVisible.value = true
}

const handleEdit = async (row: Personnel) => {
  isEdit.value = true
  await loadUsers()
  // 尝试根据邮箱匹配已有用户
  selectedUserId.value = allUsers.value.find(u => u.email === row.email)?.id
  Object.assign(form, { id: row.id, name: row.name, phone: row.phone || '', email: row.email || '', company: row.company || '', position: row.position || '', workExperience: row.workExperience || '', entryDate: row.entryDate || '', projectStartDate: row.projectStartDate || '', onProjectStatus: row.onProjectStatus || '离项', salary: row.salary || '', location: row.location || '', remark: row.remark || '', status: row.status, sort: row.sort || 0 })
  drawerVisible.value = true
}

const handleDelete = async (row: Personnel) => {
  try {
    await ElMessageBox.confirm(`确定要删除人员"${row.name}"吗？`, '删除确认', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    await PersonnelApi.del(row.id)
    ElMessage.success('删除成功'); loadData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || '删除失败') }
}

const handleAction = (key: string, row: Personnel) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  const ids = selectedRows.value.map(r => r.id)
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个人员吗？`, '批量删除', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    await PersonnelApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个人员`)
    selectedRows.value = []; loadData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || '批量删除失败') }
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) { await PersonnelApi.update(form as UpdatePersonnelReq); ElMessage.success('更新成功') }
    else { await PersonnelApi.create(form as CreatePersonnelReq); ElMessage.success('创建成功') }
    drawerVisible.value = false; loadData()
  } catch (e: any) { ElMessage.error(e.message || '操作失败') }
  finally { submitting.value = false }
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => loadData())
</script>

<script lang="ts">
import { watch } from 'vue'
export default { name: 'PersonnelList' }
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.page {
  padding: var(--space-4);
  min-height: 100vh;
  background: var(--color-page-bg);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
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

.btn-icon {
  margin-right: 4px;
  flex-shrink: 0;
}

/* ==================== 筛选栏 ==================== */
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

  &__search { display: flex; align-items: center; gap: var(--space-2); }
  &__selects { display: flex; align-items: center; gap: var(--space-2); }
  &__actions { display: flex; align-items: center; gap: var(--space-2); margin-left: auto; }
}

/* ==================== 内容卡片 ==================== */
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

/* ==================== 表格 ==================== */
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

/* ==================== 分页栏 ==================== */
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
  strong { color: var(--color-text-primary); font-weight: 700; }
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

/* 内联标题区 */
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
  letter-spacing: 0.3px;
  text-transform: uppercase;

  &.tag--edit {
    background: rgba(0, 94, 235, 0.1);
    color: var(--color-primary);
  }
  &.tag--new {
    background: rgba(22, 163, 74, 0.1);
    color: var(--color-success);
  }
}

.drawer-title-text {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: -0.2px;
}

/* 侧边栏内容 */
.drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--color-surface-2);

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
}

/* 侧边栏底部 */
.drawer-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

/* ==================== 表单样式 ==================== */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  animation: drawer-form-in 0.4s ease both 0.05s;
}

@keyframes drawer-form-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.edit-form :deep(.el-form-item) {
  margin-bottom: 0;
  .el-form-item__label {
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-secondary);
    margin-bottom: 4px;
  }
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 var(--space-4);
}

.form-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 var(--space-4);
}

.form-hint {
  margin-left: var(--space-2);
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 用户关联选择区 */
.user-select-section {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: var(--space-3);
  margin-bottom: var(--space-3);
}

.user-select-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: var(--space-2);
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);

  svg { color: var(--color-primary); }
}

.user-select-tip {
  margin-top: var(--space-1);
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
}

/* 用户下拉选项 */
.user-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 2px 0;
}

.user-option-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.user-option-email {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 分隔线 */
.form-divider {
  height: 1px;
  background: var(--color-border-light);
  margin-bottom: var(--space-3);
}

/* ==================== 响应式 14寸 ==================== */
@media (max-width: 1366px) {
  .page { padding: var(--space-3); gap: var(--space-2); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-bar { padding: var(--space-2) var(--space-3); gap: var(--space-2); }
  .filter-bar__actions { margin-left: 0; }
}

@media (max-width: 1024px) {
  .form-grid { grid-template-columns: 1fr; }
  .form-row-2 { grid-template-columns: 1fr; }
}
</style>
