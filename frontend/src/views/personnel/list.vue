<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('personnel.list.title') }}</h1>
        <span class="page-subtitle">{{ t('personnel.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button type="success" @click="handleExport" :loading="exporting">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          {{ t('common.export') }}
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('common.create') }}
        </el-button>
      </div>
    </header>

    <!-- 人员统计 -->
    <div class="person-stats" v-if="!loading && pagination.total > 0">
      <div class="person-stat" @click="filterByPosition('')">
        <span class="person-stat-num">{{ pagination.total }}</span>
        <span class="person-stat-label">全部人员</span>
      </div>
      <div class="stat-divider"></div>
      <div class="person-stat" v-for="s in positionStats" :key="s.position" @click="filterByPosition(s.position)" :style="{ background: s.color.bg }">
        <span class="person-stat-num" :style="{ color: s.color.text }">{{ s.count }}</span>
        <span class="person-stat-label" :style="{ color: s.color.text }">{{ s.position }}</span>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-bar__search">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('personnel.list.searchPlaceholder')"
          clearable
          @keyup.enter="handleSearch"
          style="width: 260px"
        >
          <template #prefix>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          </template>
        </el-input>
      </div>
      <div class="filter-bar__selects">
        <el-select v-model="searchStatus" :placeholder="t('common.status')" clearable style="width: 110px">
          <el-option :label="t('common.enabled')" value="active" />
          <el-option :label="t('common.disabled')" value="inactive" />
        </el-select>
        <el-select v-model="searchOnProject" :placeholder="t('personnel.list.onProjectStatus')" clearable style="width: 110px">
          <el-option :label="t('personnel.list.onProject')" value="在项" />
          <el-option :label="t('personnel.list.offProject')" value="离项" />
        </el-select>
        <el-select v-model="searchPosition" :placeholder="t('personnel.list.form.position')" clearable style="width: 130px" @change="handleSearch">
          <el-option label="测试工程师" value="测试工程师" />
          <el-option label="网络工程师" value="网络工程师" />
          <el-option label="安全工程师" value="安全工程师" />
          <el-option label="开发工程师" value="开发工程师" />
          <el-option label="运维工程师" value="运维工程师" />
          <el-option label="运营人员" value="运营人员" />
          <el-option label="合规专家" value="合规专家" />
          <el-option label="解决方案" value="解决方案" />
          <el-option label="商务人员" value="商务人员" />
          <el-option label="成本人员" value="成本人员" />
          <el-option label="驻场人员" value="驻场人员" />
          <el-option label="驻场人员-ODC" value="驻场人员-ODC" />
          <el-option label="项目管理" value="项目管理" />
          <el-option label="合规负责人" value="合规负责人" />
          <el-option label="产品人员" value="产品人员" />
          <el-option label="其他人员" value="其他人员" />
        </el-select>
      </div>
      <div class="filter-bar__actions">
        <el-button type="primary" @click="handleSearch">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          {{ t('common.search') }}
        </el-button>
        <el-button @click="handleReset">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4"/></svg>
          {{ t('common.reset') }}
        </el-button>
        <el-button v-if="selectedRows.length > 0" type="danger" plain @click="handleBatchDelete">
          {{ t('common.batchDelete') }} ({{ selectedRows.length }})
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
        <el-table-column :label="t('personnel.list.form.nameLabel')" min-width="100" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="cell-name">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('personnel.list.form.phone')" min-width="130">
          <template #default="{ row }">
            <span class="cell-with-copy">
              <span>{{ row.phone || '—' }}</span>
              <el-tooltip v-if="row.phone" content="复制手机号" placement="top">
                <button class="copy-btn copy-btn--sm" @click.stop="copyText(row.phone)">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="t('personnel.list.form.email')" min-width="180">
          <template #default="{ row }">
            <span class="cell-with-copy">
              <span class="email-cell">{{ row.email || '—' }}</span>
              <el-tooltip v-if="row.email" content="复制邮箱" placement="top">
                <button class="copy-btn copy-btn--sm" @click.stop="copyText(row.email)">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="company" :label="t('personnel.list.form.company')" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="cell-company">{{ row.company || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="position" :label="t('personnel.list.form.position')" min-width="90" align="center">
          <template #default="{ row }">
            <span v-if="row.position" class="cell-position" :style="{ background: positionColors[row.position]?.bg, color: positionColors[row.position]?.text }">{{ row.position }}</span>
            <span v-else class="cell-position cell-position--none">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="workExperience" :label="t('personnel.list.form.workExperience')" min-width="80" align="center" />
        <el-table-column prop="entryDate" :label="t('personnel.list.form.entryDate')" min-width="100" align="center" />
        <el-table-column prop="projectStartDate" :label="t('personnel.list.form.projectStartDate')" min-width="100" align="center" />
        <el-table-column prop="onProjectStatus" :label="t('personnel.list.onProjectStatus')" min-width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.onProjectStatus === '在项' ? 'success' : 'warning'" size="small" effect="light">
              {{ row.onProjectStatus || t('personnel.list.offProject') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="salary" :label="t('personnel.list.form.salary')" min-width="80" align="center" />
        <el-table-column prop="location" :label="t('personnel.list.form.location')" min-width="110" show-overflow-tooltip />
        <el-table-column prop="status" :label="t('common.status')" min-width="68" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light">
              {{ row.status === 'active' ? t('common.enabled') : t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="t('common.sort')" min-width="60" align="center" />
        <el-table-column :label="t('common.actions')" width="80" fixed="right" align="center">
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
          small
        />
      </div>
    </div>

    <!-- 编辑/新增侧边栏 -->
    <el-drawer
      v-model="drawerVisible"
      direction="rtl"
      size="520px"
      :with-header="true"
      :destroy-on-close="true"
      class="personnel-drawer"
    >
      <template #header>
        <div class="drawer-title-inner">
          <span class="drawer-mode-tag" :class="isEdit ? 'tag--edit' : 'tag--new'">{{ isEdit ? t('common.edit') : t('common.create') }}</span>
          <span class="drawer-title-text">{{ isEdit ? form.name || t('personnel.list.form.person') : t('personnel.list.form.newPerson') }}</span>
        </div>
      </template>

      <!-- 侧边栏内容 -->
      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top" class="edit-form">

          <!-- 用户关联选择 -->
          <div class="user-select-section">
            <div class="user-select-header">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              <span>{{ t('personnel.list.form.linkSystemUser') }}</span>
            </div>
            <el-select
              v-model="selectedUserId"
              :placeholder="t('personnel.list.form.selectUserPlaceholder')"
              clearable
              filterable
              :loading="userSelectLoading"
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
            <div class="user-select-tip">{{ t('personnel.list.form.userSelectTip') }}</div>
          </div>

          <div class="form-divider"></div>

          <div class="form-grid">
            <el-form-item :label="t('personnel.list.form.nameLabel')" prop="name">
              <el-input v-model="form.name" :placeholder="t('personnel.list.form.namePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.phone')" prop="phone">
              <el-input v-model="form.phone" :placeholder="t('personnel.list.form.phonePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.email')" prop="email">
              <el-input v-model="form.email" :placeholder="t('personnel.list.form.emailPlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.company')" prop="company">
              <el-input v-model="form.company" :placeholder="t('personnel.list.form.companyPlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.position')" prop="position">
              <el-select v-model="form.position" :placeholder="t('personnel.list.form.positionPlaceholder')" style="width: 100%" clearable>
                <el-option label="测试工程师" value="测试工程师" />
                <el-option label="网络工程师" value="网络工程师" />
                <el-option label="安全工程师" value="安全工程师" />
                <el-option label="开发工程师" value="开发工程师" />
                <el-option label="运维工程师" value="运维工程师" />
                <el-option label="运营人员" value="运营人员" />
                <el-option label="合规专家" value="合规专家" />
                <el-option label="解决方案" value="解决方案" />
                <el-option label="商务人员" value="商务人员" />
                <el-option label="成本人员" value="成本人员" />
                <el-option label="驻场人员" value="驻场人员" />
                <el-option label="驻场人员-ODC" value="驻场人员-ODC" />
                <el-option label="项目管理" value="项目管理" />
                <el-option label="合规负责人" value="合规负责人" />
                <el-option label="产品人员" value="产品人员" />
                <el-option label="其他人员" value="其他人员" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.workExperience')" prop="workExperience">
              <el-input v-model="form.workExperience" :placeholder="t('personnel.list.form.workExperiencePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.entryDate')" prop="entryDate">
              <el-date-picker v-model="form.entryDate" type="date" :placeholder="t('personnel.list.form.selectDate')" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.projectStartDate')" prop="projectStartDate">
              <el-date-picker v-model="form.projectStartDate" type="date" :placeholder="t('personnel.list.form.selectDate')" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
            <el-form-item :label="t('personnel.list.onProjectStatus')" prop="onProjectStatus">
              <el-select v-model="form.onProjectStatus" :placeholder="t('personnel.list.form.selectOnProjectStatus')" style="width: 100%">
                <el-option :label="t('personnel.list.onProject')" value="在项" />
                <el-option :label="t('personnel.list.offProject')" value="离项" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('personnel.list.form.salary')" prop="salary">
              <el-input v-model="form.salary" :placeholder="t('personnel.list.form.salaryPlaceholder')" />
            </el-form-item>
          </div>
          <el-form-item :label="t('personnel.list.form.location')" prop="location">
            <el-input v-model="form.location" :placeholder="t('personnel.list.form.locationPlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('personnel.list.form.remark')" prop="remark">
            <el-input v-model="form.remark" type="textarea" :rows="2" :placeholder="t('personnel.list.form.remarkPlaceholder')" />
          </el-form-item>
          <div class="form-row-2">
            <el-form-item :label="t('common.status')" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio value="active">{{ t('common.enabled') }}</el-radio>
                <el-radio value="inactive">{{ t('common.disabled') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="t('common.sort')" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="9999" />
              <span class="form-hint">{{ t('personnel.list.form.sortHint') }}</span>
            </el-form-item>
          </div>
        </el-form>
      </div>

      <!-- 侧边栏底部 -->
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { PersonnelApi, type Personnel, type CreatePersonnelReq, type UpdatePersonnelReq } from '@/api/personnel'
import { UserApi, type User } from '@/api/user'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()

const trackExport = inject<(success?: boolean) => void>('trackExport')

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
const searchPosition = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 获取职位的 CSS class 名（用于颜色标签）
const getPositionClass = (position: string) => {
  return position.replace(/[^a-zA-Z0-9\u4e00-\u9fa5]/g, '-').toLowerCase()
}

// 职位颜色配置（16种，足够区分所有岗位）
const positionColors: Record<string, { bg: string; text: string }> = {
  '测试工程师':   { bg: '#e6f7ff', text: '#1890ff' },
  '网络工程师':   { bg: '#f6ffed', text: '#52c41a' },
  '安全工程师':   { bg: '#fff7e6', text: '#fa8c16' },
  '开发工程师':   { bg: '#f9f0ff', text: '#722ed1' },
  '运维工程师':   { bg: '#fff1f0', text: '#f5222d' },
  '运营人员':     { bg: '#e6fffb', text: '#13c2c2' },
  '合规专家':     { bg: '#f0f5ff', text: '#597ef7' },
  '解决方案':     { bg: '#fff0f6', text: '#eb2f96' },
  '商务人员':     { bg: '#fcffe6', text: '#d4b106' },
  '成本人员':     { bg: '#e6f7ff', text: '#096dd9' },
  '驻场人员':     { bg: '#f6ffed', text: '#389e0d' },
  '驻场人员-ODC': { bg: '#fff7e6', text: '#d46b08' },
  '项目管理':     { bg: '#f9f0ff', text: '#531dab' },
  '合规负责人':   { bg: '#fcffe6', text: '#ad8b00' },
  '产品人员':     { bg: '#e6fffb', text: '#08979c' },
  '其他人员':     { bg: '#f5f5f5', text: '#595959' },
}

// 职位统计
const positionStats = computed(() => {
  const positions = ['测试工程师', '网络工程师', '安全工程师', '开发工程师', '运维工程师', '运营人员', '合规专家', '解决方案', '商务人员', '成本人员', '驻场人员', '驻场人员-ODC', '项目管理', '合规负责人', '产品人员', '其他人员']
  return positions
    .map(p => ({
      position: p,
      count: tableData.value.filter(r => r.position === p).length,
      color: positionColors[p] || { bg: '#f0f0f0', text: '#666' }
    }))
    .filter(s => s.count > 0)
})

const filterByPosition = (position: string) => {
  searchPosition.value = searchPosition.value === position ? '' : position
  pagination.page = 1
  loadData()
}

const form = reactive<CreatePersonnelReq & { id?: number }>({
  name: '', phone: '', email: '', company: '', position: '',
  workExperience: '', entryDate: '', projectStartDate: '',
  onProjectStatus: '在项', salary: '', location: '',
  remark: '', status: 'active', sort: 0,
})

const formRules = computed(() => ({
  name: [{
    validator: (_rule: any, value: string, callback: Function) => {
      // 如果 name 有值，直接通过
      if (value && value.trim()) {
        callback()
        return
      }
      // 如果没有值，但选择了关联用户，等 watch 自动填充，暂时通过
      if (selectedUserId.value) {
        callback()
        return
      }
      // 都没有才报错
      callback(new Error(t('personnel.list.form.nameRequired')))
    },
    trigger: 'blur'
  }],
}))

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
      // 如果有职位筛选，客户端过滤（后端不支持position过滤则使用前端过滤）
      if (searchPosition.value) {
        tableData.value = tableData.value.filter(r => r.position === searchPosition.value)
        pagination.total = tableData.value.length
      }
    }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { searchKeyword.value = ''; searchStatus.value = ''; searchOnProject.value = ''; searchPosition.value = ''; pagination.page = 1; loadData() }
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
    link.href = url; link.download = `${t('personnel.list.exportFileName')}_${new Date().getTime()}.xlsx`
    document.body.appendChild(link); link.click(); document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success(t('common.exportSuccess'))
    trackExport?.(true)
  } catch (e: any) { ElMessage.error(e.message || t('common.exportError')) }
  finally { exporting.value = false }
}

const handleCreate = async () => {
  isEdit.value = false
  selectedUserId.value = undefined
  await loadUsers()
  Object.assign(form, { id: undefined, name: '', phone: '', email: '', company: '', position: '', workExperience: '', entryDate: '', projectStartDate: '', onProjectStatus: '在项', salary: '', location: '', remark: '', status: 'active', sort: 0 })
  drawerVisible.value = true
}

const handleEdit = async (row: Personnel) => {
  isEdit.value = true
  await loadUsers()
  // 尝试根据邮箱匹配已有用户
  selectedUserId.value = allUsers.value.find(u => u.email === row.email)?.id
  Object.assign(form, { id: row.id, name: row.name, phone: row.phone || '', email: row.email || '', company: row.company || '', position: row.position || '', workExperience: row.workExperience || '', entryDate: row.entryDate || '', projectStartDate: row.projectStartDate || '', onProjectStatus: row.onProjectStatus || '在项', salary: row.salary || '', location: row.location || '', remark: row.remark || '', status: row.status, sort: row.sort || 0 })
  drawerVisible.value = true
}

const handleDelete = async (row: Personnel) => {
  try {
    await ElMessageBox.confirm(t('personnel.list.form.deleteConfirm', { name: row.name }), t('personnel.list.form.deleteConfirmTitle'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    await PersonnelApi.del(row.id)
    ElMessage.success(t('common.deleteSuccess')); loadData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || t('common.deleteError')) }
}

const handleAction = (key: string, row: Personnel) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  const ids = selectedRows.value.map(r => r.id)
  try {
    await ElMessageBox.confirm(t('personnel.list.form.batchDeleteConfirm', { count: selectedRows.value.length }), t('personnel.list.form.batchDeleteTitle'), { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' })
    await PersonnelApi.batchDelete(ids)
    ElMessage.success(t('personnel.list.form.batchDeleteSuccess', { count: selectedRows.value.length }))
    selectedRows.value = []; loadData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || t('common.batchDeleteError')) }
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) { await PersonnelApi.update(form as UpdatePersonnelReq); ElMessage.success(t('common.updateSuccess')) }
    else { await PersonnelApi.create(form as CreatePersonnelReq); ElMessage.success(t('common.createSuccess')) }
    drawerVisible.value = false; loadData()
  } catch (e: any) { ElMessage.error(e.message || t('common.operationError')) }
  finally { submitting.value = false }
}

// 复制到剪贴板
const copyText = async (text: string) => {
  if (!text) {
    ElMessage.warning('无可复制内容')
    return
  }
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      ElMessage.success('已复制到剪贴板')
      return
    }
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.select()
    const success = document.execCommand('copy')
    document.body.removeChild(textarea)
    if (success) {
      ElMessage.success('已复制到剪贴板')
    } else {
      ElMessage.error('复制失败')
    }
  } catch {
    ElMessage.error('复制失败')
  }
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => loadData())
</script>

<script lang="ts">
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

/* 人员统计 */
.person-stats {
  display: flex;
  align-items: center;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  padding: 12px 20px;
  margin-bottom: 12px;
  overflow-x: auto;
  gap: 0;
}

.person-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
  flex-shrink: 0;
  min-width: 60px;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
}

.person-stat-num {
  font-size: 18px;
  font-weight: 700;
  line-height: 1;
}

.person-stat-label {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: #ebeef5;
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

/* ==================== 单元格美化 ==================== */
.cell-name {
  font-weight: 600;
  color: var(--color-text-primary);
  font-size: 13px;
}

.cell-company {
  color: var(--color-text-secondary);
  font-size: 13px;
}

.cell-position {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;

  &--none {
    background: var(--color-surface-2);
    color: var(--color-text-muted);
  }
}

/* ==================== 复制按钮样式 ==================== */
.cell-with-copy {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.email-cell {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text-muted);
  transition: all 0.15s ease;
  flex-shrink: 0;
  padding: 0;

  svg {
    width: 12px;
    height: 12px;
  }

  &:hover {
    background: var(--color-primary-light-9);
    color: var(--color-primary);
  }

  &:active {
    transform: scale(0.95);
  }
}

.copy-btn--sm {
  width: 20px;
  height: 20px;

  svg {
    width: 11px;
    height: 11px;
  }
}
</style>
