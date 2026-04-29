<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('user.list.title') }}</h1>
        <span class="page-subtitle">{{ t('user.list.subtitle') }}</span>
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
      <el-popover placement="bottom-end" :width="240" trigger="click">
        <template #reference>
          <el-button>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px">
              <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
              <rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>
            </svg>
            {{ t('user.list.columnDisplay') }}
          </el-button>
        </template>
        <div class="col-settings">
          <div class="col-settings__header">
            <span>{{ t('user.list.columnSettingsTitle') }}</span>
            <el-button type="primary" text size="small" @click="handleResetColumns">{{ t('user.list.columnReset') }}</el-button>
          </div>
          <div class="col-settings__list">
            <div v-for="col in allColumns" :key="col.prop" class="col-settings__item" :class="{ 'is-disabled': col.required }">
              <el-checkbox v-model="col.visible" :disabled="col.required" @change="saveColumnSettings">
                {{ col.label }}
              </el-checkbox>
            </div>
          </div>
        </div>
      </el-popover>
    </div>

    <!-- 表格卡片 -->
    <div class="content-card">
      <el-table ref="tableRef" v-model:selection="selectedRows" :data="tableData" v-loading="loading" stripe @selection-change="handleSelectionChange" @sort-change="handleSortChange">
        <el-table-column type="selection" width="40" fixed="left" />
        <el-table-column v-if="isColumnVisible('nickname')" prop="nickname" :label="t('user.list.table.name')" min-width="110" show-overflow-tooltip sortable="custom">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar" :style="{ background: getAvatarColor(row.nickname) }">{{ (row.nickname || row.username || 'U').charAt(0).toUpperCase() }}</div>
              <span class="user-name">{{ row.nickname || row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('username')" prop="username" :label="t('user.list.table.username')" min-width="100" show-overflow-tooltip sortable="custom" />
        <el-table-column v-if="isColumnVisible('email')" prop="email" :label="t('user.list.table.email')" min-width="160" show-overflow-tooltip sortable="custom" />
        <el-table-column v-if="isColumnVisible('roleName')" prop="roleName" :label="t('user.list.table.role')" min-width="100" show-overflow-tooltip sortable="custom">
          <template #default="{ row }">
            <el-tag v-if="row.roleName" type="info" size="small" effect="plain">{{ row.roleName }}</el-tag>
            <span v-else class="empty-text">—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('groupName')" prop="groupName" :label="t('user.list.table.group')" min-width="110" sortable="custom">
          <template #default="{ row }">
            <span v-if="row.groupName" class="group-link" @click="jumpToGroup(row.groupId)">{{ row.groupName }}</span>
            <span v-else class="empty-text">—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('status')" prop="status" :label="t('common.status')" width="90" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small" effect="light">
              {{ row.status === 'active' ? t('common.enabled') : t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('locked')" prop="locked" :label="t('user.list.table.locked')" width="80" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag v-if="row.locked" type="danger" size="small" effect="plain">
              <el-icon style="vertical-align: middle; margin-right: 2px"><Lock /></el-icon>
              {{ t('user.list.table.lockedStatus') }}
            </el-tag>
            <span v-else class="empty-text">—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('mfaEnabled')" prop="mfaEnabled" :label="t('user.list.table.mfa')" width="80" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="row.mfaEnabled ? 'success' : 'info'" size="small" effect="plain">
              {{ row.mfaEnabled ? t('user.list.table.mfaOn') : t('user.list.table.mfaOff') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('department')" prop="department" :label="t('user.list.table.department')" min-width="110" show-overflow-tooltip sortable="custom">
          <template #default="{ row }">
            <span :class="row.department ? '' : 'empty-text'">{{ row.department || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('source')" prop="source" :label="t('user.list.table.source')" width="90" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag v-if="row.source === 'AD'" type="warning" size="small" effect="plain">
              <el-icon style="vertical-align: middle; margin-right: 2px"><Key /></el-icon>
              AD域
            </el-tag>
            <el-tag v-else type="info" size="small" effect="plain">本地</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('createdAt')" prop="createdAt" :label="t('common.createdAt')" min-width="170" sortable="custom">
          <template #default="{ row }">
            <span class="time-text" v-html="formatDate(row.createdAt)"></span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="row.locked ? [
              { key: 'unlock', label: t('user.list.actions.unlock'), type: 'danger' },
              { key: 'edit', label: t('common.edit'), type: 'primary' },
              { key: 'resetPwd', label: t('user.list.form.resetPwd'), type: 'warning', disabled: row.source === 'AD', disabledTip: t('user.list.form.noResetPwdAD') },
              { key: 'delete', label: t('common.delete'), type: 'danger' }
            ] : [
              { key: 'edit', label: t('common.edit'), type: 'primary' },
              { key: 'resetPwd', label: t('user.list.form.resetPwd'), type: 'warning', disabled: row.source === 'AD', disabledTip: t('user.list.form.noResetPwdAD') },
              { key: 'delete', label: t('common.delete'), type: 'danger' }
            ]" :row="row" @action="(key, r) => handleAction(key, r)" />
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

          <!-- 密码要求（新增时显示） -->
          <div class="pwd-reqs" v-if="form.password && !isEdit">
            <span class="req-tag" :class="{ ok: formPwdChecks.length }">≥8位</span>
            <span class="req-tag" :class="{ ok: formPwdChecks.upper }">大写</span>
            <span class="req-tag" :class="{ ok: formPwdChecks.lower }">小写</span>
            <span class="req-tag" :class="{ ok: formPwdChecks.number }">数字</span>
            <span class="req-tag" :class="{ ok: formPwdChecks.special }">特殊</span>
          </div>

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
              <el-radio value="active">{{ t('common.enabled') }}</el-radio>
              <el-radio value="inactive">{{ t('common.disabled') }}</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item :label="t('common.sort')" prop="sort">
            <el-input-number v-model="form.sort" :min="0" :max="9999" />
          </el-form-item>
          <!-- 操作区 -->
          <div v-if="isEdit" class="action-section">
            <div class="action-divider"></div>
            <div class="action-row">
              <el-button type="warning" plain @click="handleResetPwdFromDrawer">
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
    <el-dialog v-model="showResetPwd" width="420px" destroy-on-close append-to-body class="reset-pwd-dialog">
      <template #header>
        <div class="pwd-dialog-head">
          <div class="pwd-dialog-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          </div>
          <div class="pwd-dialog-text">
            <span class="pwd-dialog-title">{{ t('user.list.form.resetPwd') }}</span>
            <span class="pwd-dialog-sub" v-if="resetPwdTarget">{{ resetPwdTarget.nickname || resetPwdTarget.username }}</span>
          </div>
        </div>
      </template>

      <el-form ref="resetPwdFormRef" :model="resetPwdForm" :rules="resetPwdRules" class="reset-pwd-form">
      <div class="pwd-dialog-body">
        <div class="pwd-field">
          <label class="pwd-label">{{ t('user.list.form.newPassword') }}</label>
          <el-form-item prop="newPassword" style="margin-bottom: 0">
            <el-input
              v-model="resetPwdForm.newPassword"
              type="password"
              show-password
              size="large"
              :placeholder="t('user.list.form.newPasswordPlaceholder')"
              class="pwd-input"
            />
          </el-form-item>
        </div>

        <!-- 密码强度条 + 要求标签 -->
        <div class="pwd-strength" v-if="resetPwdForm.newPassword">
          <div class="pwd-strength__bar">
            <div class="pwd-strength__fill" :class="'strength--' + passwordStrength" :style="{ width: passwordStrength === 'none' ? '0' : passwordStrength === 'weak' ? '25' : passwordStrength === 'fair' ? '50' : passwordStrength === 'good' ? '75' : '100' + '%' }"></div>
          </div>
          <div class="pwd-strength__tags">
            <span class="req-tag" :class="{ ok: resetPwdChecks.length }">{{ t('user.list.form.pwdLen') || '≥8位' }}</span>
            <span class="req-tag" :class="{ ok: resetPwdChecks.upper }">{{ t('user.list.form.pwdUpper') || '大写' }}</span>
            <span class="req-tag" :class="{ ok: resetPwdChecks.lower }">{{ t('user.list.form.pwdLower') || '小写' }}</span>
            <span class="req-tag" :class="{ ok: resetPwdChecks.number }">{{ t('user.list.form.pwdDigit') || '数字' }}</span>
            <span class="req-tag" :class="{ ok: resetPwdChecks.special }">{{ t('user.list.form.pwdSpecial') || '特殊' }}</span>
          </div>
        </div>

        <div class="pwd-field">
          <label class="pwd-label">{{ t('user.list.form.confirmPassword') }}</label>
          <el-form-item prop="confirmPassword" style="margin-bottom: 0">
            <el-input
              v-model="resetPwdForm.confirmPassword"
              type="password"
              show-password
              size="large"
              :placeholder="t('user.list.form.confirmPasswordPlaceholder')"
              class="pwd-input"
            />
          </el-form-item>
        </div>
      </div>
      </el-form>

      <template #footer>
        <div class="pwd-dialog-foot">
          <el-button size="default" @click="showResetPwd = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="resetPwdLoading" @click="confirmResetPwd">{{ t('common.confirm') }}</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 导入对话框 -->
    <el-dialog v-model="importDialogVisible" width="560px" destroy-on-close append-to-body class="import-dialog">
      <div class="import-dlg-header">
        <span class="import-dlg-title">{{ t('user.list.messages.importTitle') }}</span>
        <span class="import-dlg-sub">{{ t('user.list.messages.importSub') }}</span>
      </div>

      <div class="import-steps">
        <span class="step-num" :class="{ active: importStep >= 1, done: importStep > 1 }">1</span>
        <span class="step-line" :class="{ done: importStep > 1 }"></span>
        <span class="step-num" :class="{ active: importStep >= 2, done: importStep > 2 }">2</span>
        <span class="step-line" :class="{ done: importStep > 2 }"></span>
        <span class="step-num" :class="{ active: importStep >= 3 }">3</span>
      </div>

      <!-- Step 1: 上传文件 + 自动识别 -->
      <div v-if="importStep === 1" class="step-content">
        <!-- 下载模板按钮 -->
        <div class="step-toolbar">
          <el-button class="btn-outline" @click="downloadTemplate">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            {{ t('user.list.messages.downloadTemplate') }}
          </el-button>
        </div>
        <!-- 字段说明 -->
        <div class="field-guide" v-if="importTemplateFields && importTemplateFields.length > 0">
          <div class="fg-row" v-for="f in importTemplateFields" :key="f.code">
            <span class="fg-name">{{ f.field }}</span>
            <span class="fg-type">{{ f.required ? t('common.required') : t('common.optional') }}</span>
            <span class="fg-hint">{{ f.example ? `例: ${f.example}` : '' }}</span>
          </div>
        </div>
        <!-- 上传区域（选中文件后变为预览状态） -->
        <div class="upload-area" :class="{ 'upload-area--preview': selectedFile && !importing && importStep === 1 }" @click="triggerFileInput">
          <input ref="fileInputRef" type="file" accept=".xlsx,.xls" style="display:none" @change="handleFileChange" />
          <!-- 未选中文件 -->
          <template v-if="!selectedFile">
            <div class="upload-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
            </div>
            <div class="upload-text">
              <span class="upload-tip">{{ t('user.list.messages.uploadTip') }}</span>
              <span class="upload-hint">{{ t('user.list.messages.uploadHint') }}</span>
            </div>
          </template>
          <!-- 已选中文件：预览识别结果 -->
          <template v-else>
            <div class="preview-icon">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#409eff" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
            </div>
            <div class="preview-text">
              <span class="preview-title">{{ selectedFile.name }}</span>
              <span class="preview-count" v-if="previewing">{{ t('common.loading') }}</span>
              <span class="preview-count" v-else-if="previewTotal > 0">已识别 {{ previewTotal }} 条用户记录</span>
              <span class="preview-empty" v-else-if="previewError">{{ previewError }}</span>
              <span class="preview-empty" v-else>未识别到有效用户数据，请检查文件格式</span>
            </div>
          </template>
        </div>
        <div class="step-actions">
          <el-button @click="importDialogVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :disabled="!selectedFile || previewing" :loading="importing" @click="startImport">{{ importing ? t('user.list.messages.importProgress') : t('user.list.messages.startImport') }}</el-button>
        </div>
      </div>

      <!-- Step 2: 导入中 -->
      <div v-if="importStep === 2" class="step-content">
        <div class="importing-state">
          <div class="progress-wrap">
            <span class="progress-label">{{ t('user.list.messages.importProgress') }}</span>
            <span class="progress-pct">{{ importProgress }}%</span>
          </div>
          <div class="progress-bar"><div class="progress-bar__fill" :style="{ width: importProgress + '%' }"></div></div>
        </div>
      </div>

      <!-- Step 3: 导入结果 -->
      <div v-if="importStep === 3" class="step-content">
        <div class="import-result">
          <div class="res-banner" :class="importResult && importResult.failed > 0 ? 'res-banner--warn' : 'res-banner--ok'">
            <svg v-if="importResult && importResult.failed === 0" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
            <svg v-else-if="importResult" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2.5"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            <div class="res-text">
              <span class="res-title">导入完成，共 {{ importResult?.total ?? 0 }} 条</span>
              <span class="res-ok">成功 {{ importResult?.success ?? 0 }} 条</span>
              <span class="res-fail" v-if="importResult && importResult.failed > 0">失败 {{ importResult.failed }} 条</span>
            </div>
          </div>
          <div class="fail-list" v-if="importResult && importResult.failRows && importResult.failRows.length > 0">
            <div class="fail-item" v-for="(f, idx) in importResult.failRows" :key="idx">
              <span class="fail-row">第{{ f.row }}行</span>
              <span class="fail-reason">{{ f.reason }}</span>
            </div>
          </div>
        </div>
        <div class="step-actions">
          <el-button v-if="importResult && importResult.success > 0" type="primary" @click="importDialogVisible = false; loadData()">{{ t('common.done') }}</el-button>
          <el-button v-else @click="resetImport">{{ t('common.retry') }}</el-button>
        </div>
      </div>
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

// 密码策略 - 固定要求：8位以上，大小写字母、数字、特殊字符
const PASSWORD_MIN_LENGTH = 8
const PASSWORD_REQUIRE_UPPER = true
const PASSWORD_REQUIRE_LOWER = true
const PASSWORD_REQUIRE_DIGIT = true
const PASSWORD_REQUIRE_SPECIAL = true

const pwdPolicy = ref({
  passwordMinLength: PASSWORD_MIN_LENGTH,
  passwordRequireUppercase: PASSWORD_REQUIRE_UPPER,
  passwordRequireLowercase: PASSWORD_REQUIRE_LOWER,
  passwordRequireDigit: PASSWORD_REQUIRE_DIGIT,
  passwordRequireSpecial: PASSWORD_REQUIRE_SPECIAL
})

const resetPwdRules = computed<FormRules>(() => ({
  newPassword: [
    { required: true, message: t('user.list.form.passwordRequired'), trigger: 'blur' },
    { validator: (_: any, val: string, cb: any) => {
      if (!isPwdValid(val)) cb(new Error('至少8位，包含大小写字母、数字和特殊字符'))
      else cb()
    }, trigger: 'change' }
  ],
  confirmPassword: [
    { required: true, message: t('user.list.form.confirmPasswordRequired'), trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value && value !== resetPwdForm.newPassword) callback(new Error(t('user.list.form.passwordMismatch')))
        else callback()
      },
      trigger: 'change'
    }
  ]
}))

// 特殊字符检测：直接列出常用特殊字符，明确覆盖 %
const hasSpecialChar = (pwd: string): boolean => {
  const specialChars = '!@#$%^&*()_+-=[]{}|;:\'",.<>?/`~\\'
  for (const c of pwd) {
    if (specialChars.includes(c)) return true
  }
  return false
}

// 密码检查
const resetPwdChecks = computed(() => ({
  length: resetPwdForm.newPassword.length >= PASSWORD_MIN_LENGTH,
  upper: /[A-Z]/.test(resetPwdForm.newPassword),
  lower: /[a-z]/.test(resetPwdForm.newPassword),
  number: /[0-9]/.test(resetPwdForm.newPassword),
  special: hasSpecialChar(resetPwdForm.newPassword)
}))

// 密码强度
const passwordStrength = computed(() => {
  const p = resetPwdForm.newPassword
  if (!p) return 'none'
  const checks = [p.length >= 8, /[A-Z]/.test(p), /[a-z]/.test(p), /[0-9]/.test(p), hasSpecialChar(p)]
  const satisfied = checks.filter(Boolean).length
  if (satisfied <= 1) return 'weak'
  if (satisfied <= 2) return 'fair'
  if (satisfied <= 3) return 'good'
  return 'strong'
})

// 表单密码检查（新增用户）
const formPwdChecks = computed(() => ({
  length: form.password.length >= PASSWORD_MIN_LENGTH,
  upper: /[A-Z]/.test(form.password),
  lower: /[a-z]/.test(form.password),
  number: /[0-9]/.test(form.password),
  special: hasSpecialChar(form.password)
}))

// 所有密码要求是否满足
const isPwdValid = (pwd: string) => {
  return pwd.length >= PASSWORD_MIN_LENGTH &&
    /[A-Z]/.test(pwd) && /[a-z]/.test(pwd) &&
    /[0-9]/.test(pwd) && hasSpecialChar(pwd)
}

const roles = ref<any[]>([])
const groups = ref<any[]>([])

// 导入导出相关
const exporting = ref(false)
const importDialogVisible = ref(false)
const importing = ref(false)
const importProgress = ref(0)
const importResult = ref<{ total: number; success: number; failed: number; failRows: { row: number; data: string; reason: string }[] } | null>(null)
const importTemplateFields = ref<{ field: string; code: string; required: boolean; type: string; options?: string; maxLength?: number; example?: string }[]>([])
const importStep = ref(1)
const fileInputRef = ref()
const selectedFile = ref<File | null>(null)
const previewing = ref(false)
const previewTotal = ref(0)
const previewError = ref('')

const keyword = ref('')
const status = ref('')
const filterRoleId = ref<number | undefined>()
const filterGroupId = ref<number | undefined>()
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const sortField = ref('')
const sortOrder = ref('')

const form = reactive<CreateUserReq & { id?: number; roleId?: number; groupId?: number }>({
  username: '', nickname: '', email: '', password: '', roleId: undefined, groupId: undefined, status: 'active', sort: 0,
})

const formRules = computed(() => ({
  username: [{ required: true, message: t('user.list.messages.usernameRequired'), trigger: 'blur' }],
  nickname: [{ required: true, message: t('user.list.messages.nicknameRequired'), trigger: 'blur' }],
  email: [{ type: 'email', message: t('user.list.messages.emailFormatTip'), trigger: 'blur' }],
  password: isEdit.value ? [] : [
    { required: true, message: t('user.list.form.passwordRequired'), trigger: 'blur' },
    { validator: (_: any, val: string, cb: any) => {
      if (!isPwdValid(val)) cb(new Error('至少8位，包含大小写字母、数字和特殊字符'))
      else cb()
    }, trigger: 'blur' }
  ]
}))

// 格式化时间 YYYY-MM-DD HH:mm:ss
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '<span class="empty-text">—</span>'
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
      sortField: sortField.value || undefined,
      sortOrder: sortOrder.value || undefined,
    })
    if (res.code === 200) { tableData.value = res.data.items || []; pagination.total = res.data.total || 0 }
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; loadData() }
const handleReset = () => { keyword.value = ''; status.value = ''; filterRoleId.value = undefined; filterGroupId.value = undefined; sortField.value = ''; sortOrder.value = ''; pagination.page = 1; loadData() }
const handleSelectionChange = (rows: User[]) => { selectedRows.value = rows }
const handleSortChange = ({ prop, order }: { prop: string; order: string }) => {
  sortField.value = prop || ''
  sortOrder.value = order === 'ascending' ? 'asc' : order === 'descending' ? 'desc' : ''
  pagination.page = 1
  loadData()
}

// 字段显示配置
interface ColumnConfig { prop: string; label: string; visible: boolean; required: boolean }
const allColumns = ref<ColumnConfig[]>([
  { prop: 'nickname', label: '姓名', visible: true, required: true },
  { prop: 'username', label: '用户名', visible: true, required: true },
  { prop: 'email', label: '邮箱', visible: true, required: false },
  { prop: 'roleName', label: '角色', visible: true, required: false },
  { prop: 'groupName', label: '用户组', visible: true, required: false },
  { prop: 'status', label: '状态', visible: true, required: false },
  { prop: 'locked', label: '锁定', visible: true, required: false },
  { prop: 'mfaEnabled', label: 'MFA', visible: true, required: false },
  { prop: 'department', label: '部门', visible: true, required: false },
  { prop: 'source', label: '来源', visible: true, required: false },
  { prop: 'createdAt', label: '创建时间', visible: true, required: true },
])
const COL_STG_KEY = 'usr_list_cols'
const loadColumnSettings = () => {
  const saved = localStorage.getItem(COL_STG_KEY)
  if (saved) {
    try {
      const cfg: ColumnConfig[] = JSON.parse(saved)
      allColumns.value.forEach(col => {
        const c = cfg.find(x => x.prop === col.prop)
        if (c && !col.required) col.visible = c.visible
      })
    } catch {}
  }
}
const saveColumnSettings = () => {
  localStorage.setItem(COL_STG_KEY, JSON.stringify(allColumns.value.map(c => ({ prop: c.prop, visible: c.visible }))))
}
const handleResetColumns = () => {
  allColumns.value.forEach(c => { if (!c.required) c.visible = true })
  saveColumnSettings()
}
const isColumnVisible = (prop: string) => allColumns.value.find(c => c.prop === prop)?.visible !== false

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
  showResetPwd.value = true
}

// 从编辑抽屉调用重置密码
const handleResetPwdFromDrawer = async () => {
  if (!form.id) return
  resetPwdTarget.value = {
    id: form.id,
    username: form.username,
    nickname: form.nickname,
    email: form.email,
    phone: '',
    avatar: '',
    status: form.status,
    statusText: '',
    roleId: form.roleId || 0,
    roleName: '',
    groupId: form.groupId || 0,
    groupName: '',
    mfaEnabled: form.mfaEnabled || false,
    source: 'LOCAL',
    createdAt: ''
  } as User
  resetPwdForm.newPassword = ''
  resetPwdForm.confirmPassword = ''
  showResetPwd.value = true
}

const confirmResetPwd = async () => {
  // 先检查密码格式是否满足
  if (!isPwdValid(resetPwdForm.newPassword)) {
    ElMessage.warning('新密码不满足格式要求，请检查是否包含大小写字母、数字和特殊字符')
    return
  }
  // 再检查两次输入是否一致
  if (resetPwdForm.confirmPassword !== resetPwdForm.newPassword) {
    ElMessage.warning('两次输入的密码不一致，请重新确认')
    return
  }
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
  } catch (err: any) {
    // 统一错误消息提取
    let msg = t('user.list.messages.resetPwdFailed')
    if (err) {
      // axios 错误结构: err.response.data.message
      if (err.response?.data?.message) {
        msg = err.response.data.message
      } else if (err.data?.message) {
        msg = err.data.message
      } else if (err.message) {
        msg = err.message
      } else if (typeof err === 'string') {
        msg = err
      }
    }
    ElMessage.error(msg)
    console.error('Reset password error:', err)
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

// 导出
const handleExport = async () => {
  exporting.value = true
  try {
    const blob = await UserApi.exportExcel({ keyword: keyword.value || undefined, status: status.value || undefined })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `用户列表_${new Date().getTime()}.xlsx`
    link.click()
    URL.revokeObjectURL(link.href)
    ElMessage.success(t('common.exportSuccess'))
  } catch (e: any) { ElMessage.error(e.message || t('common.exportError')) }
  finally { exporting.value = false }
}

// 显示导入对话框
const showImportDialog = async () => {
  importDialogVisible.value = true
  importResult.value = null
  importing.value = false
  previewing.value = false
  previewTotal.value = 0
  previewError.value = ''
  importProgress.value = 0
  importStep.value = 1
  selectedFile.value = null
  try {
    const res = await UserApi.getImportTemplate()
    importTemplateFields.value = res.data.fields || []
  } catch { importTemplateFields.value = [] }
}

// 触发文件选择
const triggerFileInput = () => fileInputRef.value?.click()

// 下载模板
const downloadTemplate = async () => {
  try {
    await UserApi.downloadTemplate()
  } catch {
    ElMessage.error(t('user.list.messages.templateDownloadFailed'))
  }
}

// 选择文件并自动预览
const handleFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  selectedFile.value = file
  previewTotal.value = 0
  previewError.value = ''
  previewing.value = true
  try {
    const res = await UserApi.previewImport(file)
    previewTotal.value = res.data.total || 0
    if (previewTotal.value > 0) {
      ElMessage.success({ message: `已识别 ${previewTotal.value} 条用户记录`, grouping: true })
    } else {
      ElMessage.warning({ message: '未识别到有效用户数据，请检查文件格式' })
    }
  } catch (err: any) {
    previewTotal.value = 0
    previewError.value = err?.message || '预览识别失败'
    ElMessage.error({ message: previewError.value })
  } finally {
    previewing.value = false
  }
}

// 开始导入
const startImport = async () => {
  if (!selectedFile.value) return
  importing.value = true
  importProgress.value = 0
  importResult.value = null
  importStep.value = 2
  try {
    const res = await UserApi.importUsers(selectedFile.value, (pct) => { importProgress.value = pct })
    importResult.value = res.data
    importProgress.value = 100
    if (res.data.failed === 0) ElMessage.success(t('user.list.messages.importSuccessMsg', res.data.success))
    else ElMessage.warning(t('user.list.messages.importPartial', res.data.success, res.data.failed))
    importStep.value = 3
  } catch (e: any) { ElMessage.error(e.message || t('user.list.messages.importFailedMsg')); importResult.value = { total: 0, success: 0, failed: 0, failRows: [] } }
  finally { importing.value = false }
}

// 重置导入
const resetImport = () => {
  importStep.value = 1
  importResult.value = null
  previewTotal.value = 0
  previewError.value = ''
  previewing.value = false
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
onMounted(() => { loadData(); loadRolesAndGroups(); loadColumnSettings() })
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

.empty-text { color: var(--el-text-color-placeholder); }

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
  margin-bottom: 6px;
  .el-form-item__label { font-size: 12px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 4px; }
  .el-form-item__error { padding-top: 2px; }
}

.action-section { margin-top: var(--space-2); }
.action-divider { height: 1px; background: var(--color-border-light); margin-bottom: var(--space-3); }
.action-row { display: flex; gap: var(--space-2); flex-wrap: wrap; }

/* 重置密码弹窗 */
.reset-pwd-dialog {
  .el-dialog__header {
    padding: 20px 24px 16px !important;
    margin-right: 0 !important;
    border-bottom: none !important;
  }
  .el-dialog__body { padding: 0 24px 24px !important; }
  .el-dialog__footer { padding: 0 24px 20px !important; border-top: 1px solid #f0f0f0 !important; }
}

/* 弹窗头部 */
.pwd-dialog-head {
  display: flex;
  align-items: center;
  gap: 12px;
}
.pwd-dialog-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: #fff7ed;
  border: 1px solid #fed7aa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ea580c;
  flex-shrink: 0;
}
.pwd-dialog-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.pwd-dialog-title {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: #111827;
}
.pwd-dialog-sub {
  font-size: 12px;
  color: #9ca3af;
}

/* 弹窗内容 */
.pwd-dialog-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  animation: pwd-in 0.25s ease both;
}
@keyframes pwd-in {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 字段 */
.pwd-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.pwd-label {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  font-family: 'Manrope', sans-serif;
}
.pwd-input :deep(.el-input__wrapper) {
  border-radius: 8px !important;
  background: #fff !important;
  border: 1.5px solid #e5e7eb !important;
  box-shadow: none !important;
  padding: 4px 12px !important;
  transition: border-color 0.2s, box-shadow 0.2s !important;
  &:hover { border-color: #d1d5db !important; }
  &.is-focus {
    border-color: #005eeb !important;
    box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.08) !important;
  }
}
.pwd-input :deep(.el-input__inner) {
  font-size: 14px !important;
  color: #111827 !important;
  font-family: 'Manrope', 'DM Sans', sans-serif !important;
  &::placeholder { color: #d1d5db !important; }
}

/* 密码强度条 */
.pwd-strength {
  display: flex;
  flex-direction: column;
  gap: 8px;
  animation: strength-in 0.2s ease both;
}
@keyframes strength-in {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}
.pwd-strength__bar {
  height: 5px;
  background: #f3f4f6;
  border-radius: 3px;
  overflow: hidden;
}
.pwd-strength__fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s cubic-bezier(0.34, 1.56, 0.64, 1), background 0.3s;
  min-width: 4px;
}
.strength--none, .strength--weak { background: #f87171; }
.strength--fair  { background: #fbbf24; }
.strength--good  { background: #34d399; }
.strength--strong { background: #16a34a; }

.pwd-mismatch-tip {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-danger, #f56c6c);
  font-weight: 500;

  svg { flex-shrink: 0; }
}

.reset-pwd-form {
  :deep(.el-input__wrapper) {
    border-radius: var(--radius-md);
  }
}

.pwd-strength__tags {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

/* 密码要求标签 */
.req-tag {
  padding: 3px 9px;
  font-size: 11px;
  font-weight: 600;
  border-radius: 5px;
  background: #f9fafb;
  color: #d1d5db;
  border: 1px solid #e5e7eb;
  transition: all 0.2s;
  font-family: 'Manrope', sans-serif;

  &.ok {
    background: #f0fdf4;
    color: #16a34a;
    border-color: #bbf7d0;
  }
}

/* 弹窗底部按钮 */
.pwd-dialog-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
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

.pwd-preview { padding: 12px 14px; background: #f0fdf4; border: 1px solid #86efac; border-radius: var(--radius-md); }
.pwd-preview-label { display: flex; align-items: center; gap: 6px; font-size: 11px; color: #166534; font-weight: 600; margin-bottom: 6px; svg { stroke: #166534; } }
.pwd-preview-value { font-family: 'SF Mono', 'Consolas', monospace; font-size: 16px; font-weight: 700; color: #15803d; letter-spacing: 1px; padding: 8px 12px; background: rgba(255, 255, 255, 0.7); border-radius: var(--radius-sm); border: 1px dashed #86efac; }

@media (max-width: 1366px) {
  .page { padding: var(--space-3); gap: var(--space-2); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-bar { padding: var(--space-2) var(--space-3); gap: var(--space-2); }
}

/* ==================== 导入对话框 ==================== */
.import-dlg-header { display: flex; flex-direction: column; gap: 2px; padding: 4px 0; }
.import-dlg-title { font-size: 15px; font-weight: 700; color: #1c1917; }
.import-dlg-sub { font-size: 12px; color: #78716c; }

.step-toolbar { display: flex; justify-content: flex-end; padding: 0 0 12px; }

.import-steps { display: flex; align-items: center; justify-content: center; padding: 8px 0 20px; gap: 0; }
.step-num { width: 28px; height: 28px; border-radius: 50%; background: #e5e7eb; color: #9ca3af; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; transition: all 0.3s; }
.step-num.active { background: var(--color-primary); color: #fff; }
.step-num.done { background: var(--color-success); color: #fff; }
.step-line { flex: 1; height: 2px; background: #e5e7eb; max-width: 80px; transition: background 0.3s; }
.step-line.done { background: var(--color-success); }

.step-content { min-height: 200px; }

.field-guide { display: flex; flex-direction: column; gap: 6px; margin-bottom: 16px; background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 8px; padding: 12px; }
.fg-row { display: flex; align-items: center; gap: 8px; font-size: 12px; }
.fg-name { color: #374151; font-weight: 600; min-width: 60px; }
.fg-type { font-size: 10px; padding: 1px 6px; border-radius: 4px; background: #fee2e2; color: #dc2626; }
.fg-hint { color: #9ca3af; }

.upload-area { border: 2px dashed #e5e7eb; border-radius: 12px; padding: 32px; display: flex; flex-direction: column; align-items: center; gap: 8px; cursor: pointer; transition: border-color 0.2s; &:hover { border-color: var(--color-primary); } }
.upload-area--preview { border-style: solid; border-color: #dbeafe; background: #eff6ff; padding: 24px; }
.upload-icon { color: #94a3b8; }
.upload-text { display: flex; flex-direction: column; align-items: center; gap: 2px; }
.upload-tip { font-size: 13px; color: #374151; font-weight: 500; }
.upload-hint { font-size: 11px; color: #9ca3af; }

.preview-icon { color: #409eff; }
.preview-text { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.preview-title { font-size: 13px; color: #374151; font-weight: 500; }
.preview-count { font-size: 14px; color: #409eff; font-weight: 700; }
.preview-empty { font-size: 13px; color: #dc2626; }

.step-actions { display: flex; justify-content: flex-end; gap: 8px; padding-top: 16px; border-top: 1px solid #f0f0f0; margin-top: 16px; }

.importing-state { display: flex; flex-direction: column; gap: 10px; }
.progress-wrap { display: flex; justify-content: space-between; font-size: 13px; }
.progress-label { color: #374151; }
.progress-pct { color: var(--color-primary); font-weight: 600; }
.progress-bar { height: 6px; background: #f3f4f6; border-radius: 3px; overflow: hidden; }
.progress-bar__fill { height: 100%; background: var(--color-primary); border-radius: 3px; transition: width 0.3s; }

.import-result { display: flex; flex-direction: column; gap: 12px; }
.res-banner { display: flex; align-items: center; gap: 12px; padding: 14px; border-radius: 10px; }
.res-banner--ok { background: #f0fdf4; border: 1px solid #bbf7d0; }
.res-banner--warn { background: #fffbeb; border: 1px solid #fde68a; }
.res-text { display: flex; flex-direction: column; gap: 2px; }
.res-title { font-size: 14px; font-weight: 600; color: #111827; }
.res-ok { font-size: 12px; color: #16a34a; }
.res-fail { font-size: 12px; color: #dc2626; }

.fail-list { display: flex; flex-direction: column; gap: 6px; max-height: 150px; overflow-y: auto; }
.fail-item { display: flex; align-items: center; gap: 8px; padding: 8px 10px; background: #fef2f2; border-radius: 6px; font-size: 12px; }
.fail-row { color: #dc2626; font-weight: 600; }
.fail-reason { color: #7f1d1d; }

/* 字段显示配置 */
.col-settings { min-width: 220px; }
.col-settings__header { display: flex; justify-content: space-between; align-items: center; padding-bottom: 10px; margin-bottom: 8px; border-bottom: 1px solid #f0f0f0; font-size: 13px; font-weight: 600; color: #374151; }
.col-settings__list { display: flex; flex-direction: column; gap: 2px; }
.col-settings__item { padding: 4px 0; font-size: 13px; &.is-disabled { opacity: 0.5; } }
</style>
