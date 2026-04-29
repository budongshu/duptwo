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
        <el-button type="warning" @click="showImportDialog">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
          {{ t('common.import') }}
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('common.create') }}
        </el-button>
      </div>
    </header>

    <!-- 人员统计（基于所有数据，不受分页影响） -->
    <div class="person-stats" v-if="!loading && allPersonnelData.length > 0">
      <div class="person-stat person-stat--total" @click="filterByPosition('')">
        <span class="person-stat-num">{{ allPersonnelData.length }}</span>
        <span class="person-stat-label">全部人员</span>
      </div>
      <div class="stat-divider"></div>
      <div
        class="person-stat"
        v-for="s in positionStats"
        :key="s.position"
        @click="filterByPosition(s.position)"
        :style="{ background: s.color?.bg, borderTop: `2px solid ${s.color?.border}` }"
      >
        <span class="person-stat-num" :style="{ color: s.color?.text }">{{ s.count }}</span>
        <span class="person-stat-label" :style="{ color: s.color?.label }">{{ s.position }}</span>
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
        <!-- 字段显示控制 -->
        <el-popover placement="bottom-end" :width="220" trigger="click">
          <template #reference>
            <el-button>
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="12" x2="14" y2="12"/><line x1="4" y1="18" x2="8" y2="18"/></svg>
              字段
            </el-button>
          </template>
          <div class="column-settings">
            <div class="settings-header">
              <span class="settings-title">字段显示</span>
              <el-button type="primary" text size="small" @click="handleResetColumns">重置</el-button>
            </div>
            <div class="settings-list">
              <div v-for="col in visibleColumns" :key="col.key" class="settings-item">
                <el-checkbox v-model="col.visible" @change="saveColumnVisibility">
                  {{ col.label }}
                </el-checkbox>
              </div>
            </div>
            <div class="settings-footer">
              <span class="settings-hint">提示：拖拽表头可调整列顺序</span>
            </div>
          </div>
        </el-popover>
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
        <el-table-column v-if="isColumnVisible('name')" :label="t('personnel.list.form.nameLabel')" min-width="110" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="cell-name">
              <span class="name-avatar" :style="{ background: getAvatarColor(row.name).bg, color: getAvatarColor(row.name).text }">
                {{ row.name.charAt(0).toUpperCase() }}
              </span>
              <span class="name-text">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('phone')" :label="t('personnel.list.form.phone')" min-width="130">
          <template #default="{ row }">
            <span class="cell-with-copy">
              <span :class="row.phone ? '' : 'empty-text'">{{ row.phone || '—' }}</span>
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
        <el-table-column v-if="isColumnVisible('email')" :label="t('personnel.list.form.email')" min-width="180">
          <template #default="{ row }">
            <span class="cell-with-copy">
              <span class="email-cell" :class="row.email ? '' : 'empty-text'">{{ row.email || '—' }}</span>
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
        <el-table-column v-if="isColumnVisible('company')" prop="company" :label="t('personnel.list.form.company')" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="['cell-company', row.company ? '' : 'empty-text']">{{ row.company || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('position')" prop="position" :label="t('personnel.list.form.position')" min-width="90" align="center">
          <template #default="{ row }">
            <span v-if="row.position" class="cell-position" :style="{ background: positionColors[row.position]?.bg, color: positionColors[row.position]?.text }">{{ row.position }}</span>
            <span v-else class="cell-position cell-position--none">—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('workExperience')" prop="workExperience" :label="t('personnel.list.form.workExperience')" min-width="80" align="center" />
        <el-table-column v-if="isColumnVisible('entryDate')" prop="entryDate" :label="t('personnel.list.form.entryDate')" min-width="100" align="center" />
        <el-table-column v-if="isColumnVisible('projectStartDate')" prop="projectStartDate" :label="t('personnel.list.form.projectStartDate')" min-width="100" align="center" />
        <el-table-column v-if="isColumnVisible('onProjectStatus')" prop="onProjectStatus" :label="t('personnel.list.onProjectStatus')" min-width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.onProjectStatus === '在项' ? 'success' : 'warning'" size="small" effect="light">
              {{ row.onProjectStatus || t('personnel.list.offProject') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('salary')" prop="salary" :label="t('personnel.list.form.salary')" min-width="80" align="center" />
        <el-table-column v-if="isColumnVisible('location')" prop="location" :label="t('personnel.list.form.location')" min-width="110" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('status')" prop="status" :label="t('common.status')" min-width="68" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light">
              {{ row.status === 'active' ? t('common.enabled') : t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('sort')" prop="sort" :label="t('common.sort')" min-width="60" align="center" />
        <el-table-column :label="t('common.actions')" width="110" fixed="right" align="center">
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
          size="small"
        />
      </div>
    </div>

    <!-- 人员详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('personnel.list.detailTitle')"
      width="480px"
      :destroy-on-close="true"
      class="person-detail-dialog"
    >
      <div class="detail-content" v-if="currentDetail">
        <div class="detail-header">
          <div class="detail-avatar" :style="{ background: getAvatarColor(currentDetail.name) }">
            {{ currentDetail.name.charAt(0).toUpperCase() }}
          </div>
          <div class="detail-info">
            <div class="detail-name">{{ currentDetail.name }}</div>
            <div class="detail-position">
              <el-tag v-if="currentDetail.position" size="small" :style="{ background: positionColors[currentDetail.position]?.bg, color: positionColors[currentDetail.position]?.text }">
                {{ currentDetail.position }}
              </el-tag>
              <el-tag :type="currentDetail.status === 'active' ? 'success' : 'info'" size="small">
                {{ currentDetail.status === 'active' ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="detail-section-title">{{ t('personnel.list.form.basicInfo') }}</div>
          <div class="detail-grid">
            <div class="detail-item" v-if="currentDetail.phone">
              <span class="detail-label">{{ t('personnel.list.form.phone') }}</span>
              <span class="detail-value">{{ currentDetail.phone }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.email">
              <span class="detail-label">{{ t('personnel.list.form.email') }}</span>
              <span class="detail-value">{{ currentDetail.email }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.company">
              <span class="detail-label">{{ t('personnel.list.form.company') }}</span>
              <span class="detail-value">{{ currentDetail.company }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.location">
              <span class="detail-label">{{ t('personnel.list.form.location') }}</span>
              <span class="detail-value">{{ currentDetail.location }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="detail-section-title">{{ t('personnel.list.form.workInfo') }}</div>
          <div class="detail-grid">
            <div class="detail-item" v-if="currentDetail.workExperience">
              <span class="detail-label">{{ t('personnel.list.form.workExperience') }}</span>
              <span class="detail-value">{{ currentDetail.workExperience }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.entryDate">
              <span class="detail-label">{{ t('personnel.list.form.entryDate') }}</span>
              <span class="detail-value">{{ currentDetail.entryDate }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.projectStartDate">
              <span class="detail-label">{{ t('personnel.list.form.projectStartDate') }}</span>
              <span class="detail-value">{{ currentDetail.projectStartDate }}</span>
            </div>
            <div class="detail-item" v-if="currentDetail.onProjectStatus">
              <span class="detail-label">{{ t('personnel.list.onProjectStatus') }}</span>
              <span class="detail-value">
                <el-tag :type="currentDetail.onProjectStatus === '在项' ? 'success' : 'warning'" size="small">
                  {{ currentDetail.onProjectStatus }}
                </el-tag>
              </span>
            </div>
            <div class="detail-item" v-if="currentDetail.salary">
              <span class="detail-label">{{ t('personnel.list.form.salary') }}</span>
              <span class="detail-value">{{ currentDetail.salary }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section" v-if="currentDetail.remark">
          <div class="detail-section-title">{{ t('personnel.list.form.remark') }}</div>
          <div class="detail-remark">{{ currentDetail.remark }}</div>
        </div>

        <div class="detail-section detail-meta">
          <span class="meta-time">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {{ t('common.createdAt') }}: {{ formatDateTime(currentDetail.createdAt) }}
          </span>
          <span class="meta-time">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {{ t('common.updatedAt') }}: {{ formatDateTime(currentDetail.updatedAt) }}
          </span>
        </div>
      </div>
    </el-dialog>

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
                <el-option label="前端工程师" value="前端工程师" />
                <el-option label="算法工程师" value="算法工程师" />
                <el-option label="DBA数据库" value="DBA数据库" />
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

    <!-- 导入对话框 -->
    <el-dialog v-model="importDialogVisible" width="560px" destroy-on-close append-to-body class="import-dialog">
      <template #header>
        <div class="import-dlg-header">
          <span class="import-dlg-title">批量导入人员</span>
          <span class="import-dlg-sub">支持 xlsx 格式，可分批导入</span>
        </div>
      </template>

      <!-- 步骤条 -->
      <div class="import-steps">
        <div class="step-item" :class="{ 'step-item--active': activeImportStep === 1, 'step-item--done': activeImportStep > 1 }">
          <div class="step-circle">{{ activeImportStep > 1 ? '✓' : '1' }}</div>
          <span class="step-label">下载模板</span>
        </div>
        <div class="step-line" :class="{ 'step-line--active': activeImportStep > 1 }"></div>
        <div class="step-item" :class="{ 'step-item--active': activeImportStep === 2, 'step-item--done': activeImportStep > 2 }">
          <div class="step-circle">{{ activeImportStep > 2 ? '✓' : '2' }}</div>
          <span class="step-label">上传文件</span>
        </div>
        <div class="step-line" :class="{ 'step-line--active': activeImportStep > 2 }"></div>
        <div class="step-item" :class="{ 'step-item--active': activeImportStep === 3 }">
          <div class="step-circle">3</div>
          <span class="step-label">确认导入</span>
        </div>
      </div>

      <!-- Step 1: 下载模板 -->
      <div class="step-panel" v-if="activeImportStep === 1">
        <div class="step-hint">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <span>请先下载标准模板，按要求填写后上传</span>
        </div>
        <div class="field-guide" v-if="importTemplateFields && importTemplateFields.length > 0">
          <div class="fg-title">模板字段说明</div>
          <div class="fg-row" v-for="f in importTemplateFields" :key="f.code">
            <span class="fg-name" :class="{ 'fg-name--req': f.required }">{{ f.field }}</span>
            <span class="fg-type">{{ f.type === 'select' ? '下拉选择' : f.type === 'date' ? '日期' : '文本' }}</span>
            <span class="fg-example" v-if="f.example">示例：{{ f.example }}</span>
          </div>
        </div>
        <div class="step-actions">
          <button class="btn-outline" @click="downloadTemplate">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            下载模板
          </button>
          <button class="btn-primary" @click="activeImportStep = 2">
            已下载模板，继续
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
          </button>
        </div>
      </div>

      <!-- Step 2: 上传文件 -->
      <div class="step-panel" v-if="activeImportStep === 2">
        <div class="drop-zone" :class="{ 'drop-zone--selected': selectedFile }">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :show-file-list="false"
            :limit="1"
            accept=".xlsx"
            :on-change="handleFileChange"
          >
            <template #trigger>
              <div class="drop-zone__inner">
                <div class="drop-zone__icon" v-if="!selectedFile">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                </div>
                <div class="drop-zone__file" v-if="selectedFile">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                  <span class="drop-zone__fname">{{ selectedFile.name }}</span>
                </div>
                <span class="drop-zone__tip" v-if="!selectedFile">点击选择 Excel 文件，或拖拽到此处</span>
                <span class="drop-zone__tip drop-zone__tip--replace" v-else>点击可重新选择</span>
              </div>
            </template>
          </el-upload>
        </div>
        <div class="file-info" v-if="previewInfo">
          <span class="file-info__ok">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
            共 {{ previewInfo.dataRows }} 条有效数据
          </span>
          <span class="file-info__sheets" v-if="previewInfo.sheetName">工作表：{{ previewInfo.sheetName }}</span>
        </div>
        <div class="step-actions">
          <button class="btn-ghost" @click="activeImportStep = 1">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
            上一步
          </button>
          <button class="btn-primary" :disabled="!selectedFile" @click="startImport">
            开始导入
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
          </button>
        </div>
      </div>

      <!-- Step 3: 导入中 / 结果 -->
      <div class="step-panel" v-if="activeImportStep === 3">
        <!-- 导入中 -->
        <div class="importing-state" v-if="importing">
          <div class="progress-wrap">
            <div class="progress-row">
              <span class="progress-label">正在导入数据...</span>
              <span class="progress-pct">{{ importProgress }}%</span>
            </div>
            <div class="progress-bar"><div class="progress-bar__fill" :style="{ width: importProgress + '%' }"></div></div>
            <span class="progress-sub">请勿关闭页面</span>
          </div>
        </div>
        <!-- 导入结果 -->
        <div class="import-result" v-if="!importing && importResult">
          <div class="res-banner" :class="importResult && importResult.failed > 0 ? 'res-banner--warn' : 'res-banner--ok'">
            <svg v-if="importResult && importResult.failed === 0" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
            <svg v-else-if="importResult" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2.5"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            <div class="res-banner__text">
              <span class="res-title">导入完成，共 {{ importResult?.total ?? 0 }} 条</span>
              <span class="res-ok">成功 {{ importResult?.success ?? 0 }} 条</span>
              <span v-if="importResult && importResult.failed > 0" class="res-fail">失败 {{ importResult.failed }} 条</span>
            </div>
          </div>
          <div class="fail-list" v-if="!importing && importResult && importResult.failRows && importResult.failRows.length > 0">
            <div class="fail-list__title">失败明细</div>
            <div class="fail-item" v-for="(f, idx) in importResult.failRows" :key="idx">
              <span class="fail-item__row">第{{ f.row }}行</span>
              <span class="fail-item__data">{{ f.data }}</span>
              <span class="fail-item__reason">{{ f.reason }}</span>
            </div>
          </div>
        </div>
        <div class="step-actions" v-if="!importing">
          <button class="btn-primary" @click="closeImportDialog">完成</button>
        </div>
      </div>
    </el-dialog>
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
const importDialogVisible = ref(false)
const activeImportStep = ref(1)
const importing = ref(false)
const importProgress = ref(0)
const selectedFile = ref<File | null>(null)
const uploadRef = ref()
const previewInfo = ref<{ totalRows: number; dataRows: number; sheetName: string } | null>(null)
const importResult = ref<{ total: number; success: number; failed: number; failRows: { row: number; data: string; reason: string }[] } | null>(null)
const importTemplateFields = ref<{ field: string; code: string; required: boolean; type: string; options?: string; maxLength?: number; example?: string }[]>([])
const tableData = ref<Personnel[]>([])
// 所有人员数据（用于统计，不受当前页筛选影响）
const allPersonnelData = ref<Personnel[]>([])
const tableRef = ref()
const selectedRows = ref<Personnel[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 详情弹窗
const detailVisible = ref(false)
const currentDetail = ref<Personnel | null>(null)

// 用户选择相关
const allUsers = ref<User[]>([])
const selectedUserId = ref<number | undefined>()
const userSelectLoading = ref(false)

// 字段显示配置
const STORAGE_KEY = 'personnel_columns_visible'
interface ColumnConfig {
  key: string
  label: string
  visible: boolean
}
const visibleColumns = ref<ColumnConfig[]>([
  { key: 'name', label: '姓名', visible: true },
  { key: 'phone', label: '手机号', visible: true },
  { key: 'email', label: '邮箱', visible: true },
  { key: 'company', label: '公司', visible: false },
  { key: 'position', label: '职位', visible: true },
  { key: 'workExperience', label: '工作经验', visible: false },
  { key: 'entryDate', label: '入职日期', visible: false },
  { key: 'projectStartDate', label: '项目开始日期', visible: false },
  { key: 'onProjectStatus', label: '在项状态', visible: false },
  { key: 'salary', label: '薪资', visible: false },
  { key: 'location', label: '位置', visible: false },
  { key: 'status', label: '状态', visible: true },
  { key: 'sort', label: '排序', visible: false },
])

// 加载保存的字段显示配置
const loadColumnVisibility = () => {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved) {
    try {
      const config = JSON.parse(saved)
      visibleColumns.value.forEach(col => {
        if (config[col.key] !== undefined) {
          col.visible = config[col.key]
        }
      })
    } catch {}
  }
}

// 保存字段显示配置
const saveColumnVisibility = () => {
  const config: Record<string, boolean> = {}
  visibleColumns.value.forEach(col => { config[col.key] = col.visible })
  localStorage.setItem(STORAGE_KEY, JSON.stringify(config))
}

// 重置字段显示为默认配置
const handleResetColumns = () => {
  visibleColumns.value.forEach(col => {
    if (col.key === 'name') col.visible = true
    else if (col.key === 'phone') col.visible = true
    else if (col.key === 'email') col.visible = true
    else if (col.key === 'company') col.visible = false
    else if (col.key === 'position') col.visible = true
    else if (col.key === 'workExperience') col.visible = false
    else if (col.key === 'entryDate') col.visible = false
    else if (col.key === 'projectStartDate') col.visible = false
    else if (col.key === 'onProjectStatus') col.visible = false
    else if (col.key === 'salary') col.visible = false
    else if (col.key === 'location') col.visible = false
    else if (col.key === 'status') col.visible = true
    else if (col.key === 'sort') col.visible = false
  })
  saveColumnVisibility()
  ElMessage.success('已重置为默认配置')
}

// 检查列是否可见
const isColumnVisible = (key: string) => {
  const col = visibleColumns.value.find(c => c.key === key)
  return col ? col.visible : true
}

// 监听字段变化，保存配置
watch(visibleColumns, () => {
  saveColumnVisibility()
}, { deep: true })

const searchKeyword = ref('')
const searchStatus = ref('')
const searchOnProject = ref('')
const searchPosition = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 获取职位的 CSS class 名（用于颜色标签）
const getPositionClass = (position: string) => {
  return position.replace(/[^a-zA-Z0-9\u4e00-\u9fa5]/g, '-').toLowerCase()
}

// 人员头像颜色（基于姓名 hash，10种颜色）
const avatarColors = [
  { bg: '#e6f7ff', text: '#1890ff' },
  { bg: '#f6ffed', text: '#52c41a' },
  { bg: '#fff7e6', text: '#fa8c16' },
  { bg: '#f9f0ff', text: '#722ed1' },
  { bg: '#fff1f0', text: '#f5222d' },
  { bg: '#e6fffb', text: '#13c2c2' },
  { bg: '#f0f5ff', text: '#597ef7' },
  { bg: '#fff0f6', text: '#eb2f96' },
  { bg: '#fcffe6', text: '#d4b106' },
  { bg: '#f5f5f5', text: '#595959' },
]
const getAvatarColor = (name: string) => {
  if (!name) return avatarColors[0]
  let hash = 0
  for (let i = 0; i < name.length; i++) { hash = ((hash << 5) - hash) + name.charCodeAt(i) }
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

const formatDateTime = (timeStr: string) => {
  if (!timeStr) return '—'
  const d = new Date(timeStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  const ss = String(d.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}:${ss}`
}

// 职位颜色配置（浅色主题，实色背景 + 顶部彩边）
const positionColors: Record<string, { bg: string; border: string; text: string; label: string }> = {
  '测试工程师':    { bg: '#eff6ff', border: '#3b82f6', text: '#1d4ed8', label: '#3b82f6' },
  '前端工程师':    { bg: '#f0fdf4', border: '#22c55e', text: '#15803d', label: '#22c55e' },
  '算法工程师':    { bg: '#fdf4ff', border: '#a855f7', text: '#7e22ce', label: '#a855f7' },
  'DBA数据库':    { bg: '#fff7ed', border: '#f97316', text: '#c2410c', label: '#f97316' },
  '网络工程师':    { bg: '#ecfeff', border: '#06b6d4', text: '#0e7490', label: '#06b6d4' },
  '安全工程师':    { bg: '#f5f3ff', border: '#a855f7', text: '#7e22ce', label: '#a855f7' },
  '开发工程师':    { bg: '#fefce8', border: '#ca8a04', text: '#a16207', label: '#ca8a04' },
  '运维工程师':    { bg: '#f0fdf4', border: '#22c55e', text: '#15803d', label: '#22c55e' },
  '运营人员':      { bg: '#fff7ed', border: '#f97316', text: '#c2410c', label: '#f97316' },
  '合规专家':      { bg: '#f8fafc', border: '#64748b', text: '#475569', label: '#64748b' },
  '解决方案':      { bg: '#e0f2fe', border: '#0ea5e9', text: '#0369a1', label: '#0ea5e9' },
  '商务人员':      { bg: '#fdf2f8', border: '#ec4899', text: '#be185d', label: '#ec4899' },
  '成本人员':      { bg: '#ede9fe', border: '#8b5cf6', text: '#6d28d9', label: '#8b5cf6' },
  '驻场人员':      { bg: '#dcfce7', border: '#16a34a', text: '#166534', label: '#16a34a' },
  '驻场人员-ODC':  { bg: '#ccfbf1', border: '#14b8a6', text: '#0f766e', label: '#14b8a6' },
  '项目管理':      { bg: '#f7fee7', border: '#84cc16', text: '#4d7c0f', label: '#84cc16' },
  '合规负责人':    { bg: '#fef9c3', border: '#ca8a04', text: '#a16207', label: '#ca8a04' },
  '产品人员':      { bg: '#e0f2fe', border: '#38bdf8', text: '#0284c7', label: '#38bdf8' },
  '其他人员':      { bg: '#f5f5f5', border: '#737373', text: '#525252', label: '#737373' },
}

// 职位统计（基于所有数据，不受分页影响）
const positionStats = computed(() => {
  const positions = ['测试工程师', '前端工程师', '算法工程师', 'DBA数据库', '网络工程师', '安全工程师', '开发工程师', '运维工程师', '运营人员', '合规专家', '解决方案', '商务人员', '成本人员', '驻场人员', '驻场人员-ODC', '项目管理', '合规负责人', '产品人员', '其他人员']
  return positions
    .map(p => ({ position: p, count: allPersonnelData.value.filter(r => r.position === p).length, color: positionColors[p] }))
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
    // 先加载所有数据用于统计（不使用分页）
    const allRes = await PersonnelApi.list({
      page: 1, pageSize: 10000,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
      onProject: searchOnProject.value || undefined,
    })
    if (allRes.code === 200) {
      // 不受职位筛选影响，用于统计所有人员
      allPersonnelData.value = allRes.data.items || []
    }

    // 再加载当前页数据
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
  } catch (e: any) { ElMessage.error((e?.message) || t('common.exportError')) }
  finally { exporting.value = false }
}

const showImportDialog = async () => {
  importDialogVisible.value = true
  importResult.value = null
  previewInfo.value = null
  selectedFile.value = null
  activeImportStep.value = 1
  importing.value = false
  importProgress.value = 0
  try {
    const res = await PersonnelApi.getImportTemplate()
    if (res.code === 200) {
      importTemplateFields.value = res.data.fields || []
    }
  } catch {
    importTemplateFields.value = []
  }
}

const downloadTemplate = async () => {
  try {
    await PersonnelApi.downloadTemplate()
  } catch {
    ElMessage.error('下载模板失败')
  }
}

const handleFileChange = async (file: any) => {
  const rawFile = file.raw || file
  if (!rawFile) return
  selectedFile.value = rawFile
  previewInfo.value = null
  try {
    const res = await PersonnelApi.previewImport(rawFile)
    if (res.code === 200) {
      if (res.data.error) {
        ElMessage.warning('未找到有效表头，请确认使用的是标准模板')
      }
      previewInfo.value = res.data
    }
  } catch (e: any) {
    ElMessage.error(e.message || '预览失败')
  }
}

const startImport = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  activeImportStep.value = 3
  importing.value = true
  importProgress.value = 0
  importResult.value = null
  try {
    const res = await PersonnelApi.importPersonnel(selectedFile.value, (pct) => {
      importProgress.value = pct
    })
    importResult.value = res.data
    importProgress.value = 100
    if (res.data.failed === 0) {
      ElMessage.success(`导入成功，共 ${res.data.success} 条`)
    } else {
      ElMessage.warning(`部分导入成功 ${res.data.success} 条，失败 ${res.data.failed} 条`)
    }
    if (res.data.success > 0) {
      loadData()
    }
  } catch (e: any) {
    ElMessage.error(e.message || '导入失败')
    importResult.value = { total: 0, success: 0, failed: 0, failRows: [] }
  } finally {
    importing.value = false
  }
}

const closeImportDialog = () => {
  importDialogVisible.value = false
  importing.value = false
  importProgress.value = 0
  selectedFile.value = null
  previewInfo.value = null
  activeImportStep.value = 1
  if (importResult.value && importResult.value.success > 0) {
    loadData()
  }
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

// 查看详情
const handleView = async (row: Personnel) => {
  detailVisible.value = true
  currentDetail.value = row
}

const handleAction = (key: string, row: Personnel) => {
  if (key === 'view') handleView(row)
  else if (key === 'edit') handleEdit(row)
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
onMounted(() => { loadColumnVisibility(); loadData() })
</script>

<script lang="ts">
export default { name: 'PersonnelList' }
</script>

<style scoped lang="scss">
/* ==================== 通用 ==================== */
.empty-text { color: var(--el-text-color-placeholder); }

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
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: 14px 20px;
  margin-bottom: 12px;
  overflow-x: auto;
  gap: 0;
}

.person-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  padding: 8px 16px;
  cursor: pointer;
  border-radius: var(--radius-md);
  transition: all 0.18s ease;
  flex-shrink: 0;
  min-width: 64px;
  border: 1px solid transparent;
  border-top: 2px solid transparent;

  &:hover {
    filter: brightness(0.95);
  }

  &--total {
    padding: 8px 20px;
    background: #3b82f6;

    .person-stat-num { color: #ffffff; font-size: 20px; }
    .person-stat-label { color: rgba(255,255,255,0.8); font-size: 11px; }

    &:hover { background: #2563eb; }
  }
}

.person-stat-num {
  font-size: 18px;
  font-weight: 700;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.person-stat-label {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  letter-spacing: 0.2px;
}

.stat-divider {
  width: 1px;
  height: 36px;
  background: var(--color-border-light);
  flex-shrink: 0;
  margin: 0 4px;
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

/* 固定 drawer 高度，不受页面内容影响 */
.personnel-drawer {
  height: 100vh !important;
}

.personnel-drawer :deep(.el-drawer__body) {
  overflow-y: auto;
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

/* ==================== 人员详情弹窗 ==================== */
.detail-content {
  padding: 0 8px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-avatar {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.detail-info {
  flex: 1;
}

.detail-name {
  font-size: 20px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 8px;
}

.detail-position {
  display: flex;
  gap: 8px;
}

.detail-section {
  margin-bottom: 20px;
}

.detail-section-title {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #94a3b8;
}

.detail-value {
  font-size: 14px;
  color: #334155;
  font-weight: 500;
}

.detail-remark {
  background: var(--color-surface-2);
  border-radius: 8px;
  padding: 12px;
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.detail-meta {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--color-text-muted);
  padding-top: 16px;
  border-top: 1px solid var(--color-border-light);
}

.meta-time {
  display: inline-flex;
  align-items: center;
  gap: 4px;

  .el-icon { color: var(--color-text-muted); }
}

/* ==================== 单元格美化 ==================== */
.cell-name { display: flex; align-items: center; gap: 7px; }
.name-avatar { display: inline-flex; align-items: center; justify-content: center; width: 26px; height: 26px; border-radius: 50%; font-size: 12px; font-weight: 700; flex-shrink: 0; }
.name-text { font-weight: 500; color: var(--color-text-primary); }

.cell-company {
  color: var(--color-text-secondary);
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

/* 字段显示控制 */
.column-settings {
  .settings-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: var(--space-2);
    margin-bottom: var(--space-2);
    border-bottom: 1px solid var(--color-border-light);
  }

  .settings-title {
    font-family: 'Manrope', sans-serif;
    font-size: 14px;
    font-weight: 700;
    color: var(--color-primary);
  }

  .settings-list {
    max-height: 280px;
    overflow-y: auto;
    margin: 0 -12px;
    padding: 0 12px;
  }

  .settings-item {
    padding: 4px 10px;
    margin: 0 -10px;
    border-radius: var(--radius-sm);
    transition: background 0.15s ease;

    &:hover { background: var(--color-surface-2); }

    :deep(.el-checkbox) {
      width: 100%;
      .el-checkbox__label {
        font-size: 13px;
        color: var(--color-text-primary);
        font-weight: 500;
      }
      .el-checkbox__input.is-checked .el-checkbox__inner {
        background-color: var(--color-primary);
        border-color: var(--color-primary);
      }
      .el-checkbox__input.is-checked + .el-checkbox__label {
        color: var(--color-primary);
      }
    }
  }

  .settings-footer {
    margin-top: var(--space-2);
    padding-top: var(--space-2);
    border-top: 1px solid var(--color-border-light);
  }

  .settings-hint { font-size: 11px; color: var(--color-text-muted); }
}

/* ==================== 导入对话框 ==================== */
.import-dlg-header { display: flex; flex-direction: column; gap: 2px; padding: 4px 0; }
.import-dlg-title { font-size: 15px; font-weight: 700; color: #1c1917; }
.import-dlg-sub { font-size: 12px; color: #78716c; }

.import-steps { display: flex; align-items: center; justify-content: center; padding: 8px 0 20px; gap: 0; }
.step-item { display: flex; flex-direction: column; align-items: center; gap: 4px; flex-shrink: 0; }
.step-circle { width: 26px; height: 26px; border-radius: 50%; background: #e8e4de; color: #78716c; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; transition: all 0.2s; }
.step-item--active .step-circle { background: #409eff; color: #fff; }
.step-item--done .step-circle { background: #4a7c59; color: #fff; }
.step-label { font-size: 11px; color: #78716c; }
.step-item--active .step-label { color: #409eff; font-weight: 600; }
.step-item--done .step-label { color: #4a7c59; }
.step-line { flex: 1; height: 2px; background: #e8e4de; margin: 0 8px; margin-bottom: 18px; transition: background 0.2s; }
.step-line--active { background: #4a7c59; }

.step-panel { display: flex; flex-direction: column; gap: 14px; }
.step-hint { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #57534e; background: #f7f6f3; padding: 10px 12px; border-radius: 8px; border: 1px solid #e8e4de; }

.field-guide { background: #fafaf9; border: 1px solid #e8e4de; border-radius: 8px; padding: 12px; display: flex; flex-direction: column; gap: 6px; max-height: 240px; overflow-y: auto; }
.fg-title { font-size: 12px; font-weight: 600; color: #1c1917; margin-bottom: 4px; }
.fg-row { display: grid; grid-template-columns: 80px 64px 1fr; gap: 6px; align-items: baseline; }
.fg-name { font-size: 12px; color: #1c1917; }
.fg-name--req::after { content: ' *'; color: #ef4444; }
.fg-type { font-size: 11px; color: #78716c; }
.fg-example { font-size: 11px; color: #a8a29e; }

.step-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 4px; }

.btn-primary { display: flex; align-items: center; gap: 5px; padding: 8px 16px; background: #409eff; color: #fff; border: none; border-radius: 7px; font-size: 13px; font-weight: 600; cursor: pointer; transition: background 0.2s; &:hover { background: #66b1ff; } &:disabled { background: #a0cfff; cursor: not-allowed; } }
.btn-ghost { display: flex; align-items: center; gap: 5px; padding: 8px 14px; background: #fff; color: #57534e; border: 1px solid #e8e4de; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; &:hover { border-color: #a8a29e; } }
.btn-outline { display: flex; align-items: center; gap: 5px; padding: 8px 14px; background: #fff; color: #409eff; border: 1px solid #409eff; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; &:hover { background: #f0f7ff; } }

.drop-zone { border: 2px dashed #d4cfc6; border-radius: 10px; transition: all 0.2s; &:hover { border-color: #409eff; background: #f0f7ff; } &.drop-zone--selected { border-color: #4a7c59; border-style: solid; background: #f0fdf4; } }
.drop-zone__inner { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 32px 20px; cursor: pointer; gap: 10px; }
.drop-zone__icon { color: #a8a29e; }
.drop-zone__tip { font-size: 13px; color: #78716c; }
.drop-zone__tip--replace { color: #4a7c59; font-weight: 500; }
.drop-zone__file { display: flex; flex-direction: column; align-items: center; gap: 6px; }
.drop-zone__fname { font-size: 13px; color: #1c1917; font-weight: 500; max-width: 320px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.file-info { display: flex; align-items: center; gap: 12px; }
.file-info__ok { display: flex; align-items: center; gap: 5px; font-size: 12px; color: #4a7c59; font-weight: 600; }
.file-info__sheets { font-size: 11px; color: #a8a29e; }

.importing-state { display: flex; flex-direction: column; gap: 10px; }
.progress-wrap { display: flex; flex-direction: column; gap: 6px; }
.progress-row { display: flex; justify-content: space-between; align-items: center; }
.progress-label { font-size: 13px; color: #1c1917; font-weight: 500; }
.progress-pct { font-size: 22px; font-weight: 800; color: #1c1917; font-variant-numeric: tabular-nums; }
.progress-bar { height: 8px; background: #e8e4de; border-radius: 4px; overflow: hidden; }
.progress-bar__fill { height: 100%; background: linear-gradient(90deg, #409eff, #66b1ff); border-radius: 4px; transition: width 0.3s; }
.progress-sub { font-size: 11px; color: #a8a29e; }

.import-result { display: flex; flex-direction: column; gap: 12px; }
.res-banner { display: flex; align-items: flex-start; gap: 10px; padding: 14px 16px; border-radius: 10px; }
.res-banner--ok { background: #f0fdf4; border: 1px solid #bbf7d0; }
.res-banner--warn { background: #fffbeb; border: 1px solid #fde68a; }
.res-banner__text { display: flex; flex-direction: column; gap: 3px; }
.res-title { font-size: 13px; font-weight: 700; color: #1c1917; }
.res-ok { font-size: 12px; color: #4a7c59; font-weight: 600; }
.res-fail { font-size: 12px; color: #dc2626; font-weight: 600; }

.fail-list { display: flex; flex-direction: column; gap: 4px; max-height: 180px; overflow-y: auto; }
.fail-list__title { font-size: 12px; font-weight: 600; color: #1c1917; padding: 2px 0; }
.fail-item { display: grid; grid-template-columns: 50px 80px 1fr; gap: 8px; align-items: baseline; padding: 6px 8px; background: #fef2f2; border-radius: 6px; font-size: 11px; }
.fail-item__row { color: #78716c; }
.fail-item__data { color: #1c1917; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.fail-item__reason { color: #dc2626; }
</style>
