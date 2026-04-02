<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('uploadRecord.list.title') }}</h1>
        <p class="page-subtitle">{{ t('uploadRecord.list.subtitle') }}</p>
      </div>
      <div class="header-actions">
        <div class="header-stat" @click="filterByStatus('')">
          <span class="header-stat-num">{{ pagination.total }}</span>
          <span class="header-stat-label">全部</span>
        </div>
        <div class="header-stat header-stat--success" @click="filterByStatus('completed')">
          <span class="header-stat-num">{{ statusCount.completed }}</span>
          <span class="header-stat-label">已完成</span>
        </div>
        <div class="header-stat header-stat--warning" @click="filterByStatus('pending')">
          <span class="header-stat-num">{{ statusCount.pending }}</span>
          <span class="header-stat-label">待处理</span>
        </div>
        <div class="header-stat header-stat--info" @click="filterByStatus('processing')">
          <span class="header-stat-num">{{ statusCount.processing }}</span>
          <span class="header-stat-label">处理中</span>
        </div>
        <div class="header-stat header-stat--danger" @click="filterByStatus('failed')">
          <span class="header-stat-num">{{ statusCount.failed }}</span>
          <span class="header-stat-label">失败</span>
        </div>
        <div class="header-stat header-stat--purple" @click="">
          <span class="header-stat-num">{{ formatBytes(totalSize) }}</span>
          <span class="header-stat-label">总大小</span>
        </div>
        <el-button type="success" :loading="exporting" @click="showExportDialog">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          {{ t('uploadRecord.list.exportExcel') }}
        </el-button>
        <el-button type="warning" @click="showImportDialog">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
          {{ t('uploadRecord.list.batchImport') }}
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          {{ t('uploadRecord.list.createRecord') }}
        </el-button>
        <el-button @click="loadRecords">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M23 4v6h-6M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          {{ t('uploadRecord.list.refresh') }}
        </el-button>
      </div>
    </header>

    <!-- 导出抽屉 -->
    <el-drawer
      v-model="exportDialogVisible"
      direction="rtl"
      size="420px"
      :with-header="true"
      append-to-body
      class="export-drawer"
    >
      <template #header>
        <div class="export-drawer-header">
          <div class="export-drawer-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
          </div>
          <div class="export-drawer-head-text">
            <span class="export-drawer-title">{{ t('uploadRecord.list.exportDialogTitle') }}</span>
            <span class="export-drawer-sub">选择筛选条件导出数据</span>
          </div>
        </div>
      </template>

      <div class="export-drawer-body">
        <!-- 预览统计 -->
        <div class="export-preview-card">
          <div class="preview-stat">
            <div class="preview-num">{{ pagination.total }}</div>
            <div class="preview-label">条记录</div>
          </div>
          <div class="preview-divider"></div>
          <div class="preview-hint-text">
            <el-icon><InfoFilled /></el-icon>
            基于当前筛选条件
          </div>
        </div>

        <!-- 快速同步 -->
        <div class="export-sync-section">
          <div class="sync-header">
            <span>快速同步当前筛选</span>
            <el-switch v-model="syncCurrentFilter" size="small" />
          </div>
          <div class="sync-active-tags" v-if="syncCurrentFilter">
            <el-tag v-if="searchDataType" size="small" closable @close="searchDataType = ''; syncFilterToExport()">
              标签: {{ searchDataType }}
            </el-tag>
            <el-tag v-if="searchProjectName" size="small" closable @close="searchProjectName = ''; syncFilterToExport()">
              项目: {{ searchProjectName }}
            </el-tag>
            <el-tag v-if="searchStatus" size="small" closable @close="searchStatus = ''; syncFilterToExport()">
              {{ getStatusText(searchStatus) }}
            </el-tag>
            <el-tag v-if="searchUploader" size="small" closable @close="searchUploader = ''; syncFilterToExport()">
              上传人: {{ searchUploader }}
            </el-tag>
            <el-tag v-if="searchDateRange?.length" size="small" closable @close="searchDateRange = []; syncFilterToExport()">
              {{ searchDateRange[0] }} ~ {{ searchDateRange[1] }}
            </el-tag>
          </div>
        </div>

        <!-- 高级筛选 -->
        <div class="export-filter-card">
          <div class="filter-card-header">
            <el-icon><Filter /></el-icon>
            <span>自定义筛选</span>
          </div>
          <div class="filter-card-body">
            <div class="filter-row">
              <label>磁盘标签</label>
              <el-input v-model="exportForm.dataType" placeholder="输入标签筛选" clearable size="small" />
            </div>
            <div class="filter-row">
              <label>项目名称</label>
              <el-select v-model="exportForm.projectName" placeholder="选择项目" clearable filterable size="small" style="width: 100%">
                <el-option v-for="p in projectList" :key="p.id" :label="p.name" :value="p.name" />
              </el-select>
            </div>
            <div class="filter-row">
              <label>状态</label>
              <el-select v-model="exportForm.status" placeholder="选择状态" clearable size="small" style="width: 100%">
                <el-option label="待处理" value="pending" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="失败" value="failed" />
              </el-select>
            </div>
            <div class="filter-row">
              <label>上传人</label>
              <el-input v-model="exportForm.uploader" placeholder="输入上传人" clearable size="small" />
            </div>
            <div class="filter-row">
              <label>关键词</label>
              <el-input v-model="exportForm.keyword" placeholder="搜索关键词" clearable size="small" />
            </div>
            <div class="filter-row">
              <label>日期范围</label>
              <el-date-picker
                v-model="exportForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
                size="small"
                style="width: 100%"
              />
            </div>
          </div>
        </div>

        <!-- 导出字段选择 -->
        <div class="export-field-card">
          <div class="filter-card-header">
            <el-icon><Grid /></el-icon>
            <span>导出字段</span>
          </div>
          <div class="field-check-list">
            <el-checkbox v-model="exportAllFields" @change="handleExportAllChange">全选</el-checkbox>
            <div class="field-check-grid">
              <el-checkbox v-for="col in allColumns.filter(c => c.prop !== 'data')" :key="col.prop" v-model="col.visible">
                {{ col.label }}
              </el-checkbox>
            </div>
          </div>
        </div>
      </div>

      <div class="export-drawer-footer">
        <el-button @click="exportDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="exporting" @click="handleExport">
          <el-icon v-if="!exporting"><Download /></el-icon>
          导出 Excel
        </el-button>
      </div>
    </el-drawer>

    <!-- 批量导入弹窗 -->
    <el-dialog v-model="importDialogVisible" width="700px" destroy-on-close append-to-body class="import-dialog">
      <template #header>
        <div class="diag-head">
          <div class="diag-head-icon">
            <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="16 16 12 12 8 16"/>
              <line x1="12" y1="12" x2="12" y2="21"/>
              <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
            </svg>
          </div>
          <div class="diag-head-text">
            <span class="diag-head-title">{{ t('uploadRecord.list.importDialogTitle') }}</span>
            <span class="diag-head-sub">{{ t('uploadRecord.list.importDialogSub') }}</span>
          </div>
        </div>
      </template>

      <div class="import-body">

        <!-- 三步骤卡片 -->
        <div class="step-row">
          <div class="step-card">
            <div class="step-header">
              <span class="step-num">1</span>
              <span class="step-name">{{ t('uploadRecord.list.step1DownloadTemplate') }}</span>
            </div>
            <div class="step-desc">{{ t('uploadRecord.list.step1Desc') }}</div>
            <el-button class="step-btn" size="small" :loading="downloadingTemplate" @click="handleDownloadTemplate">{{ t('uploadRecord.list.downloadXlsx') }}</el-button>
          </div>

          <div class="step-arrow">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
          </div>

          <div class="step-card">
            <div class="step-header">
              <span class="step-num">2</span>
              <span class="step-name">{{ t('uploadRecord.list.step2UploadFile') }}</span>
            </div>
            <div class="step-desc">{{ t('uploadRecord.list.step2Desc') }}</div>
            <el-upload ref="uploadRef" class="step-upload" action="#" :auto-upload="false" :limit="1" accept=".xlsx" :on-change="handleFileChange" :on-remove="handleFileRemove" :file-list="fileList">
              <el-button size="small" plain>{{ t('uploadRecord.list.selectFile') }}</el-button>
            </el-upload>
            <div v-if="previewLoading" class="preview-hint loading">{{ t('uploadRecord.list.parsingFile') }}</div>
            <div v-else-if="previewRowCount !== null" class="preview-hint">{{ t('uploadRecord.list.recognizedRows', [previewRowCount]) }}</div>
          </div>

          <div class="step-arrow">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
          </div>

          <div class="step-card">
            <div class="step-header">
              <span class="step-num">3</span>
              <span class="step-name">{{ t('uploadRecord.list.step3ConfirmImport') }}</span>
            </div>
            <div class="step-desc">{{ t('uploadRecord.list.step3Desc') }}</div>
            <el-button class="step-btn" type="primary" size="small" :loading="importing" :disabled="!selectedFile" @click="handleImport">{{ t('uploadRecord.list.startImport') }}</el-button>
          </div>
        </div>

        <!-- 字段填写说明 -->
        <div class="field-guide" v-if="importTemplateFields.length > 0">
          <div class="field-guide-head">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
            {{ t('uploadRecord.list.fieldGuideTitle') }}
          </div>
          <div class="field-guide-body">
            <div class="fg-col-header">
              <span>{{ t('uploadRecord.list.fieldName') }}</span>
              <span>{{ t('uploadRecord.list.required') }}</span>
              <span>{{ t('uploadRecord.list.fillExample') }}</span>
              <span>{{ t('uploadRecord.list.formatDesc') }}</span>
            </div>
            <div class="fg-row" v-for="f in importTemplateFields" :key="f.code">
              <span class="fg-name">{{ f.field }}</span>
              <span class="fg-required">
                <span class="badge" :class="f.required ? 'badge--danger' : 'badge--gray'">{{ f.required ? t('uploadRecord.list.required') : t('uploadRecord.list.optional') }}</span>
              </span>
              <span class="fg-example">{{ f.example || '—' }}</span>
              <span class="fg-type">
                <span class="type-chip" v-if="f.type === 'select'">{{ t('uploadRecord.list.selectOptions') }}{{ getOptionsArr(f.options).slice(0,3).join(' / ') }}{{ getOptionsArr(f.options).length > 3 ? '...' : '' }}</span>
                <span class="type-chip" v-else-if="f.type === 'number'">{{ t('uploadRecord.list.integerBytes') }}</span>
                <span class="type-chip" v-else>{{ t('uploadRecord.list.text') }}</span>
              </span>
            </div>
          </div>
        </div>

        <!-- 导入结果 -->
        <div class="import-result" v-if="importResult">
          <div class="res-banner" :class="importResult.failed > 0 ? 'res-banner--warn' : 'res-banner--ok'">
            <div class="res-icon">
              <svg v-if="importResult.failed === 0" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2.5"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <div class="res-info">
              <span class="res-title">{{ t('uploadRecord.list.importComplete', [importResult.total]) }}</span>
              <span class="res-detail">
                <span class="res-ok">{{ t('uploadRecord.list.importSuccess', [importResult.success]) }}</span>
                <span v-if="importResult.failed > 0" class="res-fail">{{ t('uploadRecord.list.importFailed', [importResult.failed]) }}</span>
              </span>
            </div>
          </div>
          <div class="fail-list" v-if="importResult.failRows.length > 0">
            <div class="fail-list-title">{{ t('uploadRecord.list.failedRows') }}</div>
            <div class="fail-list-body">
              <div class="fail-item" v-for="(f, idx) in importResult.failRows" :key="idx">
                <span class="fail-item-num">{{ t('uploadRecord.list.row', [f.row]) }}</span>
                <span class="fail-item-msg">{{ f.reason }}</span>
              </div>
            </div>
          </div>
        </div>

      </div>

      <template #footer>
        <el-button class="diag-close-btn" @click="closeImportDialog">{{ t('uploadRecord.list.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <!-- 磁盘标签 -->
      <el-input v-model="searchDataType" :placeholder="t('uploadRecord.list.searchDiskLabel')" clearable size="default" />

      <!-- 项目名称 -->
      <el-select v-model="searchProjectName" :placeholder="t('uploadRecord.list.searchProjectName')" clearable size="default" filterable>
        <el-option v-for="p in projectList" :key="p.id" :label="p.name" :value="p.name" />
      </el-select>

      <!-- 状态 -->
      <el-select v-model="searchStatus" :placeholder="t('uploadRecord.list.searchStatus')" clearable size="default">
        <el-option :label="t('status.pending')" value="pending" />
        <el-option :label="t('status.processing')" value="processing" />
        <el-option :label="t('status.completed')" value="completed" />
        <el-option :label="t('status.failed')" value="failed" />
      </el-select>

      <!-- 上传人 -->
      <el-input v-model="searchUploader" :placeholder="t('uploadRecord.list.searchUploader')" clearable size="default" />

      <!-- 日期范围 -->
      <el-date-picker
        v-model="searchDateRange"
        type="daterange"
        :range-separator="t('common.to')"
        :start-placeholder="t('common.startDate')"
        :end-placeholder="t('common.endDate')"
        value-format="YYYY-MM-DD"
        size="default"
      />

      <!-- 关键词 -->
      <el-input v-model="searchKeyword" :placeholder="t('uploadRecord.list.searchKeywordPlaceholder')" clearable size="default">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="M21 21l-4.35-4.35"/>
          </svg>
        </template>
      </el-input>

      <!-- 操作按钮 -->
      <div class="filter-actions">
        <el-button type="primary" @click="handleSearch">{{ t('uploadRecord.list.query') }}</el-button>
        <el-button @click="handleReset">{{ t('uploadRecord.list.reset') }}</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <!-- 表格工具栏 -->
      <div class="table-toolbar">
        <div class="toolbar-left">
          <span class="record-count">{{ t('uploadRecord.list.recordCount', [pagination.total]) }}</span>
          <span v-if="selectedRows.length > 0" class="selection-count">
            {{ t('uploadRecord.list.selectedItems', [selectedRows.length]) }}
          </span>
        </div>
        <div class="toolbar-right">
          <!-- 批量删除按钮 -->
          <el-button
            v-if="selectedRows.length > 0"
            type="danger"
            @click="handleBatchDelete"
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            {{ t('uploadRecord.list.batchDelete') }}
          </el-button>
          <!-- 字段显示配置 -->
          <el-popover
            placement="bottom-end"
            :width="280"
            trigger="click"
            popper-class="column-settings-popover"
          >
            <template #reference>
              <el-button class="column-settings-btn">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
                  <rect x="3" y="3" width="7" height="7"/>
                  <rect x="14" y="3" width="7" height="7"/>
                  <rect x="14" y="14" width="7" height="7"/>
                  <rect x="3" y="14" width="7" height="7"/>
                </svg>
                {{ t('uploadRecord.list.columnDisplay') }}
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-left: 4px">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </el-button>
            </template>
            <div class="column-settings">
              <div class="settings-header">
                <span class="settings-title">{{ t('uploadRecord.list.columnSettingsTitle') }}</span>
                <el-button type="primary" text size="small" @click="handleResetColumns">{{ t('uploadRecord.list.columnReset') }}</el-button>
              </div>
              <div class="settings-list">
                <div
                  v-for="col in allColumns"
                  :key="col.prop"
                  class="settings-item"
                  :class="{ 'is-disabled': col.required }"
                >
                  <el-checkbox
                    v-model="col.visible"
                    :disabled="col.required"
                    @change="handleColumnToggle"
                  >
                    {{ col.label }}
                  </el-checkbox>
                </div>
              </div>
              <div class="settings-footer">
                <span class="settings-hint">{{ t('uploadRecord.list.columnHint') }}</span>
              </div>
            </div>
          </el-popover>
        </div>
      </div>

      <!-- 表格 -->
      <div class="table-scroll-wrapper">
        <el-table
          ref="tableRef"
          v-model:selection="selectedRows"
          :data="tableData"
          v-loading="loading"
          stripe
          size="small"
          :row-class-name="tableRowClassName"
          @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column v-if="isColumnVisible('serialNo')" prop="serialNo" :label="t('uploadRecord.list.colSerialNo')" min-width="110" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('dataType')" prop="dataType" :label="t('uploadRecord.list.colDiskLabel')" min-width="90" align="center">
          <template #default="{ row }">
            <span class="type-tag">{{ row.dataType }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('projectName')" prop="projectName" :label="t('uploadRecord.list.colProjectName')" min-width="110" show-overflow-tooltip />
        <!-- 动态字段列 -->
        <el-table-column
          v-for="col in visibleDynamicColumns"
          :key="col.code"
          :prop="`data.${col.code}`"
          :label="col.name"
          min-width="100"
          align="center"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <span class="dynamic-field">{{ row.data?.[col.code] ?? '-' }}</span>
          </template>
        </el-table-column>
        <!-- 目标路径列：文字被裁剪，悬浮显示完整路径 -->
        <el-table-column v-if="isColumnVisible('destPath')" :label="t('uploadRecord.list.colDestPath')" min-width="160" align="center">
          <template #default="{ row }">
            <div v-if="row.destPath" class="path-cell">
              <el-tooltip :content="row.destPath" placement="top" :show-after="300">
                <span class="path-text">{{ row.destPath }}</span>
              </el-tooltip>
              <el-tooltip :content="t('uploadRecord.list.copyPath')" placement="top">
                <button class="copy-btn" @click.stop="copyPath(row.destPath)" :title="t('uploadRecord.list.copyPath')">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </div>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('fileSize')" prop="fileSizeStr" :label="t('uploadRecord.list.colFileSize')" min-width="90" align="center" />
        <el-table-column v-if="isColumnVisible('uploader')" prop="uploader" :label="t('uploadRecord.list.colUploader')" min-width="80" align="center" />
        <el-table-column v-if="isColumnVisible('status')" prop="status" :label="t('uploadRecord.list.colStatus')" min-width="70" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="getStatusClass(row.status)">
              {{ row.statusText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('remark')" prop="remark" :label="t('uploadRecord.list.colRemark')" min-width="100" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('createdAt')" prop="createdAt" :label="t('uploadRecord.list.colTime')" min-width="130" />
        <el-table-column :label="t('uploadRecord.list.colActions')" width="110" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'detail', label: t('uploadRecord.list.actionDetail'), type: 'primary' },
              { key: 'edit', label: t('uploadRecord.list.actionEdit'), type: 'primary' },
              { key: 'delete', label: t('uploadRecord.list.actionDelete'), type: 'danger' }
            ]" @action="(key) => handleAction(key, row)" />
          </template>
        </el-table-column>
      </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
          @current-change="loadRecords"
          @size-change="loadRecords"
        />
      </div>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" width="640px" append-to-body>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--info">{{ t('uploadRecord.list.detailDialogTag') }}</span>
          <span class="dialog-title-text">{{ t('uploadRecord.list.detailDialogTitle') }}</span>
        </div>
      </template>
      <div v-if="currentRecord" class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="t('uploadRecord.list.detailSerialNo')">
            <code class="serial-no">{{ currentRecord.serialNo }}</code>
          </el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailDiskLabel')">{{ currentRecord.dataType }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailProjectName')">{{ currentRecord.projectName || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailUploader')">{{ currentRecord.uploader }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailStatus')">
            <span class="status-chip" :class="getStatusClass(currentRecord.status)">
              {{ currentRecord.statusText }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailFileSize')">{{ currentRecord.fileSizeStr }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailDestPath')" :span="2">
            <div class="detail-path-cell">
              <el-tooltip :content="currentRecord.destPath" placement="top" :show-after="200" v-if="currentRecord.destPath">
                <span class="detail-path-text">{{ currentRecord.destPath }}</span>
              </el-tooltip>
              <span v-else class="no-data">-</span>
              <el-tooltip :content="t('uploadRecord.list.copyPath')" placement="top" v-if="currentRecord.destPath">
                <button class="copy-btn copy-btn-lg" @click="copyPath(currentRecord.destPath)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailRemark')" :span="2">{{ currentRecord.remark || t('uploadRecord.list.detailNoRemark') }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailUploadTime')">{{ currentRecord.createdAt }}</el-descriptions-item>
          <el-descriptions-item :label="t('uploadRecord.list.detailUpdateTime')">{{ currentRecord.updatedAt }}</el-descriptions-item>

          <!-- 动态字段 -->
          <template v-if="currentRecord.data && Object.keys(currentRecord.data).length > 0">
            <el-descriptions-item
              v-for="col in dynamicColumns"
              :key="col.code"
              :label="col.name"
            >
              {{ currentRecord.data[col.code] ?? '-' }}
            </el-descriptions-item>
          </template>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">{{ t('uploadRecord.list.detailClose') }}</el-button>
      </template>
    </el-dialog>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" width="560px" destroy-on-close append-to-body>
      <template #header>
        <div class="edit-dialog-header">
          <div class="edit-dialog-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
          </div>
          <div class="edit-dialog-head-text">
            <span class="edit-dialog-title">{{ t('uploadRecord.list.editDialogTitle') }}</span>
            <span class="edit-dialog-sub" v-if="editForm.serialNo">{{ editForm.serialNo }}</span>
          </div>
        </div>
      </template>

      <div class="edit-dialog-body">
        <!-- 基本信息卡片 -->
        <div class="edit-card">
          <div class="edit-card-header">
            <el-icon><Document /></el-icon>
            <span>基本信息</span>
          </div>
          <div class="edit-card-body">
            <div class="edit-field">
              <label>状态</label>
              <el-select v-model="editForm.status" placeholder="选择状态">
                <el-option label="待处理" value="pending" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="失败" value="failed" />
              </el-select>
            </div>
            <div class="edit-field">
              <label>备注</label>
              <el-input v-model="editForm.remark" type="textarea" :rows="2" placeholder="添加备注信息" />
            </div>
          </div>
        </div>

        <!-- 动态字段卡片 -->
        <div class="edit-card" v-if="dynamicColumns.length > 0">
          <div class="edit-card-header">
            <el-icon><Grid /></el-icon>
            <span>扩展信息</span>
          </div>
          <div class="edit-card-body">
            <div class="edit-field" v-for="col in dynamicColumns" :key="col.code">
              <label>{{ col.name }}</label>
              <el-select
                v-if="col.type === 'select'"
                v-model="editForm.data[col.code]"
                clearable
                :placeholder="col.placeholder || '请选择'"
              >
                <el-option v-for="opt in col.options" :key="opt" :label="opt" :value="opt" />
              </el-select>
              <el-date-picker
                v-else-if="col.type === 'date'"
                v-model="editForm.data[col.code]"
                type="date"
                value-format="YYYY-MM-DD"
                :placeholder="col.placeholder || '选择日期'"
              />
              <el-date-picker
                v-else-if="col.type === 'datetime'"
                v-model="editForm.data[col.code]"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                :placeholder="col.placeholder || '选择时间'"
              />
              <el-input-number
                v-else-if="col.type === 'number'"
                v-model="editForm.data[col.code]"
                :placeholder="col.placeholder || '输入数字'"
              />
              <el-input
                v-else
                v-model="editForm.data[col.code]"
                :placeholder="col.placeholder || '输入内容'"
              />
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="edit-dialog-footer">
          <el-button @click="editVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmEdit">
            <el-icon v-if="!submitting"><Select /></el-icon>
            保存修改
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新增上传记录抽屉 -->
    <el-drawer
      v-model="createVisible"
      :title="null"
      direction="rtl"
      size="600px"
      :before-close="() => createVisible = false"
      append-to-body
      class="create-drawer"
    >
      <!-- 抽屉头部 -->
      <template #header>
        <div class="drawer-head">
          <div class="drawer-head-icon">
            <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="12" y1="18" x2="12" y2="12"/>
              <line x1="9" y1="15" x2="15" y2="15"/>
            </svg>
          </div>
          <div class="drawer-head-text">
            <span class="drawer-head-title">{{ t('uploadRecord.list.createDialogTitle') }}</span>
            <span class="drawer-head-sub">{{ t('uploadRecord.list.createDialogSub') }}</span>
          </div>
        </div>
      </template>

      <!-- 表单区域 -->
      <div class="create-body">

        <!-- 第一组：基础信息 -->
        <div class="field-group">
          <div class="field-group-header">
            <span class="field-group-label">{{ t('uploadRecord.list.createBasicInfo') }}</span>
            <span class="field-group-hint">{{ t('uploadRecord.list.createBasicInfoRequired') }}</span>
          </div>
          <div class="field-group-body">
            <div class="field-row">
              <div class="field-cell">
                <label class="field-label">
                  {{ t('uploadRecord.list.createDiskLabel') }} <span class="field-required">*</span>
                </label>
                <el-input
                  v-model="createForm.dataType"
                  :placeholder="t('uploadRecord.list.createDiskLabelPlaceholder')"
                  clearable
                  class="dr-input"
                />
              </div>
              <div class="field-cell">
                <label class="field-label">
                  {{ t('uploadRecord.list.createProjectName') }} <span class="field-required">*</span>
                </label>
                <el-input
                  v-model="createForm.projectName"
                  :placeholder="t('uploadRecord.list.createProjectNamePlaceholder')"
                  clearable
                  class="dr-input"
                />
              </div>
            </div>
            <div class="field-row">
              <div class="field-cell">
                <label class="field-label">
                  {{ t('uploadRecord.list.createUploader') }} <span class="field-required">*</span>
                </label>
                <el-select
                  v-model="createForm.uploader"
                  :placeholder="t('uploadRecord.list.createUploaderPlaceholder')"
                  filterable
                  allow-create
                  default-first-option
                  :reserve-keyword="false"
                  class="dr-select"
                >
                  <el-option
                    v-for="p in personnelList"
                    :key="p.id"
                    :label="p.name"
                    :value="p.name"
                  />
                </el-select>
              </div>
              <div class="field-cell">
                <label class="field-label">{{ t('uploadRecord.list.createStatus') }} <span class="field-required">*</span></label>
                <el-select v-model="createForm.status" class="dr-select dr-select--status">
                  <el-option :label="t('status.pending')" value="pending">
                    <span class="status-dot status-dot--pending"></span> {{ t('status.pending') }}
                  </el-option>
                  <el-option :label="t('status.processing')" value="processing">
                    <span class="status-dot status-dot--processing"></span> {{ t('status.processing') }}
                  </el-option>
                  <el-option :label="t('status.completed')" value="completed">
                    <span class="status-dot status-dot--completed"></span> {{ t('status.completed') }}
                  </el-option>
                  <el-option :label="t('status.failed')" value="failed">
                    <span class="status-dot status-dot--failed"></span> {{ t('status.failed') }}
                  </el-option>
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 第二组：文件信息 -->
        <div class="field-group">
          <div class="field-group-header">
            <span class="field-group-label">{{ t('uploadRecord.list.createFileInfo') }}</span>
          </div>
          <div class="field-group-body">
            <div class="field-row">
              <div class="field-cell">
                <label class="field-label">{{ t('uploadRecord.list.createDestPath') }}</label>
                <el-input
                  v-model="createForm.destPath"
                  :placeholder="t('uploadRecord.list.createDestPathPlaceholder')"
                  clearable
                  class="dr-input"
                />
              </div>
              <div class="field-cell">
                <label class="field-label">{{ t('uploadRecord.list.createFileSize') }}</label>
                <div class="size-row">
                  <el-input-number
                    v-model="fileSizeInputVal"
                    :min="0"
                    :precision="3"
                    controls-position="right"
                    class="size-num"
                    @change="syncFileSizeFromInput"
                  />
                  <el-select v-model="fileSizeUnit" class="size-unit" @change="syncFileSizeFromInput">
                    <el-option label="B" value="B" />
                    <el-option label="KB" value="KB" />
                    <el-option label="MB" value="MB" />
                    <el-option label="GB" value="GB" />
                    <el-option label="TB" value="TB" />
                  </el-select>
                </div>
                <div class="size-hint" v-if="fileSizeInputVal > 0">
                  ≈ {{ formatSizeInOtherUnits(fileSizeInputVal, fileSizeUnit) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 第三组：时间信息 -->
        <div class="field-group">
          <div class="field-group-header">
            <span class="field-group-label">{{ t('createdAt') }}</span>
          </div>
          <div class="field-group-body">
            <div class="field-row">
              <div class="field-cell field-cell--full">
                <label class="field-label">{{ t('createdAt') }}</label>
                <el-date-picker
                  v-model="createForm.createdAt"
                  type="date"
                  value-format="YYYY-MM-DD"
                  :placeholder="t('uploadRecord.list.datePlaceholder')"
                  style="width:100%"
                  class="dr-input"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 第四组：扩展字段 -->
        <div class="field-group" v-if="dynamicColumns.length > 0">
          <div class="field-group-header">
            <span class="field-group-label">{{ t('uploadRecord.list.createExtendedFields') }}</span>
            <span class="field-group-badge">{{ dynamicColumns.length }}</span>
          </div>
          <div class="field-group-body">
            <div class="field-row" v-for="col in dynamicColumns" :key="'c-' + col.code">
              <div class="field-cell field-cell--full">
                <label class="field-label">{{ col.name }}</label>
                <el-input v-if="isIpField(col)" v-model="createForm.data[col.code]" :placeholder="t('uploadRecord.list.ipPlaceholder')" clearable class="dr-input" />
                <el-input-number v-else-if="col.type === 'number'" v-model="createForm.data[col.code]" style="width:100%" controls-position="right" :placeholder="col.placeholder" class="dr-input" />
                <el-date-picker v-else-if="col.type === 'date'" v-model="createForm.data[col.code]" type="date" value-format="YYYY-MM-DD" style="width:100%" :placeholder="col.placeholder" class="dr-input" />
                <el-date-picker v-else-if="col.type === 'datetime'" v-model="createForm.data[col.code]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width:100%" :placeholder="col.placeholder" class="dr-input" />
                <el-select v-else-if="col.type === 'select'" v-model="createForm.data[col.code]" style="width:100%" clearable :placeholder="col.placeholder" class="dr-select">
                  <el-option v-for="opt in col.options" :key="opt" :label="opt" :value="opt" />
                </el-select>
                <el-input v-else v-model="createForm.data[col.code]" :placeholder="col.placeholder" clearable class="dr-input" />
              </div>
            </div>
          </div>
        </div>

        <!-- 第五组：备注 -->
        <div class="field-group">
          <div class="field-group-header">
            <span class="field-group-label">{{ t('uploadRecord.list.createRemark') }}</span>
          </div>
          <div class="field-group-body">
            <el-input
              v-model="createForm.remark"
              type="textarea"
              :rows="3"
              :placeholder="t('uploadRecord.list.createRemarkPlaceholder')"
              show-word-limit
              maxlength="500"
              class="dr-textarea"
            />
          </div>
        </div>

      </div>

      <!-- 底部按钮 -->
      <template #footer>
        <div class="drawer-foot">
          <el-button class="btn-cancel" @click="createVisible = false">{{ t('uploadRecord.list.createCancel') }}</el-button>
          <el-button type="primary" class="btn-confirm" :loading="submitting" @click="confirmCreate">{{ t('uploadRecord.list.createConfirm') }}</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, inject } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { InfoFilled, Filter, Grid, Download } from '@element-plus/icons-vue'
import { UploadRecordApi, type UploadRecord, type ImportTemplateField, type ImportResultResp } from '@/api/upload-record'
import { FieldConfigApi, type FieldConfig } from '@/api/field-config'
import { ProjectApi, type ProjectSimple } from '@/api/project'
import { PersonnelApi, type Personnel } from '@/api/personnel'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()

const trackExport = inject<(success?: boolean) => void>('trackExport')

// ==================== 字段配置 ====================
interface ColumnConfig {
  prop: string
  label: string
  visible: boolean
  required?: boolean
}

const allColumns = ref<ColumnConfig[]>([
  { prop: 'serialNo', label: t('uploadRecord.list.colSerialNo'), visible: true, required: true },
  { prop: 'dataType', label: t('uploadRecord.list.colDiskLabel'), visible: true, required: true },
  { prop: 'projectName', label: t('uploadRecord.list.colProjectName'), visible: true, required: false },
  { prop: 'destPath', label: t('uploadRecord.list.colDestPath'), visible: false, required: false },
  { prop: 'fileSize', label: t('uploadRecord.list.colFileSize'), visible: false, required: false },
  { prop: 'uploader', label: t('uploadRecord.list.colUploader'), visible: true, required: false },
  { prop: 'status', label: t('uploadRecord.list.colStatus'), visible: true, required: true },
  { prop: 'remark', label: t('uploadRecord.list.colRemark'), visible: false, required: false },
  { prop: 'createdAt', label: t('uploadRecord.list.colTime'), visible: true, required: false },
])

const STORAGE_KEY = 'upload_record_columns'

const loadColumnSettings = () => {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved) {
    try {
      const savedConfig = JSON.parse(saved)
      allColumns.value.forEach(col => {
        const savedCol = savedConfig.find((c: ColumnConfig) => c.prop === col.prop)
        if (savedCol && !col.required) {
          col.visible = savedCol.visible
        }
      })
    } catch (e) {
      console.error('Failed to load column settings:', e)
    }
  }
}

const saveColumnSettings = () => {
  const config = allColumns.value.map(col => ({
    prop: col.prop,
    label: col.label,
    visible: col.visible,
    required: col.required
  }))
  localStorage.setItem(STORAGE_KEY, JSON.stringify(config))
}

const handleColumnToggle = () => {
  saveColumnSettings()
}

const handleResetColumns = () => {
  allColumns.value.forEach(col => {
    if (!col.required) {
      col.visible = false
    }
  })
  saveColumnSettings()
  ElMessage.success(t('uploadRecord.list.resetColumnSuccess'))
}

const isColumnVisible = (prop: string) => {
  const col = allColumns.value.find(c => c.prop === prop)
  return col?.visible ?? false
}

const visibleDynamicColumns = computed(() => {
  return dynamicColumns.value.filter(col => {
    const config = allDynamicColumns.value.find(c => c.prop === col.code)
    return config?.visible ?? false
  })
})

const allDynamicColumns = ref<ColumnConfig[]>([])

const updateDynamicColumnsConfig = () => {
  allDynamicColumns.value = dynamicColumns.value.map(col => ({
    prop: col.code,
    label: col.name,
    visible: true,
    required: false
  }))
}

// ==================== 其他状态 ====================
const loading = ref(false)
const submitting = ref(false)
const selectedRows = ref<UploadRecord[]>([])
const exportDialogVisible = ref(false)

// ==================== 批量导入状态 ====================
const importDialogVisible = ref(false)
const importing = ref(false)
const downloadingTemplate = ref(false)
const selectedFile = ref<File | null>(null)
const fileList = ref<any[]>([])
const uploadRef = ref()
const importTemplateFields = ref<ImportTemplateField[]>([])
const importResult = ref<ImportResultResp | null>(null)
const previewRowCount = ref<number | null>(null)
const previewLoading = ref(false)

// options 字段是逗号分隔字符串，需要解析后展示
const getOptionsArr = (opts: string | undefined) => {
  if (!opts) return []
  return String(opts).split(',').map(s => s.trim()).filter(Boolean)
}
const searchDataType = ref('')
const searchStatus = ref('')
const searchUploader = ref('')
const searchKeyword = ref('')
const searchDateRange = ref<string[]>([])
const tableData = ref<UploadRecord[]>([])
const dataTypes = ref<string[]>([])
const dynamicColumns = ref<FieldConfig[]>([])
const projectList = ref<ProjectSimple[]>([])
const personnelList = ref<Personnel[]>([])
const searchProjectName = ref('')
const detailVisible = ref(false)
const editVisible = ref(false)
const createVisible = ref(false)
const createFormRef = ref()
const currentRecord = ref<UploadRecord | null>(null)
const tableRef = ref()

// 统计
const totalSize = computed(() => tableData.value.reduce((sum, r) => sum + (r.fileSize || 0), 0))
const statusCount = computed(() => {
  const s = { completed: 0, pending: 0, processing: 0, failed: 0 }
  for (const r of tableData.value) {
    const st = r.status || 'pending'
    if (st in s) (s as any)[st]++
  }
  return s
})
const filterByStatus = (status: string) => {
  searchStatus.value = searchStatus.value === status ? '' : status
  loadRecords()
}
const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

// 文件大小工具
const fileSizeUnit = ref('MB')
const fileSizeInputVal = ref(0)

const BYTE_UNITS: Record<string, number> = {
  B: 1,
  KB: 1024,
  MB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
  TB: 1024 * 1024 * 1024 * 1024,
  PB: 1024 * 1024 * 1024 * 1024 * 1024,
}

// 将用户输入（数值+单位）转换为字节，更新到 createForm.fileSize
const syncFileSizeFromInput = () => {
  createForm.fileSize = fileSizeInputVal.value * BYTE_UNITS[fileSizeUnit.value]
}

// 将 bytes 转换为显示值（用于初始化已有数据）
const initFileSizeFromBytes = (bytes: number) => {
  if (bytes === 0) {
    fileSizeInputVal.value = 0
    fileSizeUnit.value = 'MB'
    return
  }
  // 找到最合适的单位
  const units = ['PB', 'TB', 'GB', 'MB', 'KB', 'B']
  for (const unit of units) {
    const quotient = bytes / BYTE_UNITS[unit]
    if (quotient >= 1) {
      fileSizeInputVal.value = Math.round(quotient * 1000) / 1000
      fileSizeUnit.value = unit
      return
    }
  }
  fileSizeInputVal.value = bytes
  fileSizeUnit.value = 'B'
}

// 格式化为其他单位显示（仅展示用，非精确转换）
const formatSizeInOtherUnits = (val: number, unit: string): string => {
  const bytes = val * BYTE_UNITS[unit]
  if (bytes === 0) return ''
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const parts: string[] = []
  for (const u of units) {
    if (u === unit) continue
    const q = bytes / BYTE_UNITS[u]
    if (q >= 1) {
      parts.push(`${q >= 1000 ? q.toFixed(0) : q.toFixed(2)} ${u}`)
    }
  }
  return parts.length > 0 ? `= ${parts.join(', ')}` : ''
}

// 根据字段名称判断类型，返回对应的 CSS class
const isIpField = (col: FieldConfig) => {
  const code = (col.code || '').toLowerCase()
  const name = (col.name || '').toLowerCase()
  return code.includes('ip') || name.includes('ip') || code.includes('addr') || name.includes('地址')
}

const isPathField = (col: FieldConfig) => {
  const code = (col.code || '').toLowerCase()
  const name = (col.name || '').toLowerCase()
  return code.includes('path') || code.includes('dir') || code.includes('folder') ||
         name.includes('路径') || name.includes('目录') || name.includes('文件夹')
}

const getFieldClass = (col: FieldConfig) => {
  if (isIpField(col)) return 'input-ip'
  if (col.type === 'number') return 'input-num'
  if (isPathField(col)) return 'input-path'
  return ''
}

const editForm = reactive<{ id: number; serialNo: string; status: string; remark: string; data: Record<string, any> }>({
  id: 0,
  serialNo: '',
  status: '',
  remark: '',
  data: {}
})

const exportForm = reactive({
  dataType: '',
  projectName: '',
  status: '',
  uploader: '',
  dateRange: [] as string[],
  keyword: ''
})

// 导出相关
const syncCurrentFilter = ref(true)
const exportAllFields = ref(true)
const exporting = ref(false)

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待处理', processing: '处理中', completed: '已完成', failed: '失败' }
  return map[status] || status
}

// 同步当前筛选条件到导出表单
const syncFilterToExport = () => {
  if (syncCurrentFilter.value) {
    exportForm.dataType = searchDataType.value
    exportForm.projectName = searchProjectName.value
    exportForm.status = searchStatus.value
    exportForm.uploader = searchUploader.value
    exportForm.dateRange = searchDateRange.value ? [...searchDateRange.value] : []
    exportForm.keyword = searchKeyword.value
  }
}

// 导出字段全选切换
const handleExportAllChange = (val: boolean) => {
  allColumns.value.forEach(c => { if (c.prop !== 'data') c.visible = val })
}

const showExportDialog = () => {
  // 重置导出表单
  exportForm.dataType = ''
  exportForm.projectName = ''
  exportForm.status = ''
  exportForm.uploader = ''
  exportForm.dateRange = []
  exportForm.keyword = ''
  exportAllFields.value = true
  // 同步当前筛选
  syncFilterToExport()
  exportDialogVisible.value = true
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params: any = {}
    if (exportForm.dataType) params.dataType = exportForm.dataType
    if (exportForm.projectName) params.projectName = exportForm.projectName
    if (exportForm.status) params.status = exportForm.status
    if (exportForm.uploader) params.uploader = exportForm.uploader
    if (exportForm.keyword) params.keyword = exportForm.keyword
    if (exportForm.dateRange && exportForm.dateRange.length === 2) {
      params.startDate = exportForm.dateRange[0]
      params.endDate = exportForm.dateRange[1]
    }

    const res = await UploadRecordApi.exportExcel(params)
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 10)
    link.href = url
    link.download = `上传记录_${timestamp}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    exportDialogVisible.value = false
    ElMessage.success('导出成功')
    trackExport?.(true)
  } catch (error: any) {
    ElMessage.error(error?.message || '导出失败')
  } finally {
    exporting.value = false
  }
}

const createForm = reactive({
  dataType: '',
  projectName: '',
  destPath: '',
  fileSize: 0,
  uploader: '',
  status: 'pending',
  remark: '',
  createdAt: '',
  data: {} as Record<string, any>
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const getStatusClass = (status: string) => {
  return `status-chip--${status === 'completed' ? 'success' : status === 'failed' ? 'danger' : status === 'pending' ? 'warning' : 'info'}`
}

// 复制路径到剪贴板
const copyPath = async (path: string) => {
  try {
    await navigator.clipboard.writeText(path)
    ElMessage.success(t('uploadRecord.list.copyPathSuccess'))
  } catch {
    ElMessage.error(t('uploadRecord.list.copyPathFailed'))
  }
}

const tableRowClassName = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 0 ? 'even-row' : 'odd-row'
}

const handleSelectionChange = (rows: UploadRecord[]) => {
  selectedRows.value = rows
}

const loadFieldConfigs = async () => {
  try {
    const res = await FieldConfigApi.getAllEnabled()
    dynamicColumns.value = res.data || []
    updateDynamicColumnsConfig()
  } catch (error) {
    console.error('Failed to load field configs:', error)
  }
}

const loadProjects = async () => {
  try {
    const res = await ProjectApi.getSimpleList()
    projectList.value = res.data || []
  } catch (error) {
    console.error('Failed to load projects:', error)
  }
}

const loadPersonnel = async () => {
  try {
    const res = await PersonnelApi.getAll()
    const raw: any = res.data
    personnelList.value = Array.isArray(raw) ? raw : []
  } catch (error) {
    console.error('Failed to load personnel:', error)
    personnelList.value = []
  }
}

const loadRecords = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      dataType: searchDataType.value || undefined,
      status: searchStatus.value || undefined,
      uploader: searchUploader.value || undefined,
      projectName: searchProjectName.value || undefined,
      keyword: searchKeyword.value || undefined,
      startDate: searchDateRange.value?.[0] || undefined,
      endDate: searchDateRange.value?.[1] || undefined
    }

    const res = await UploadRecordApi.list(params)
    tableData.value = res.data.items || []
    pagination.total = res.data.total || 0

    // 更新数据类型列表
    const stats = await UploadRecordApi.statistics()
    if (stats.data.byDataType) {
      dataTypes.value = stats.data.byDataType.map((d: { dataType: string }) => d.dataType)
    }
  } catch (error) {
    console.error('Failed to load records:', error)
    ElMessage.error(t('uploadRecord.list.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadRecords()
}

const handleReset = () => {
  searchDataType.value = ''
  searchStatus.value = ''
  searchUploader.value = ''
  searchProjectName.value = ''
  searchKeyword.value = ''
  searchDateRange.value = []
  pagination.page = 1
  loadRecords()
}

const handleDetail = (row: UploadRecord) => {
  currentRecord.value = row
  detailVisible.value = true
}

const handleEdit = (row: UploadRecord) => {
  editForm.id = row.id
  editForm.serialNo = row.serialNo
  editForm.status = row.status
  editForm.remark = row.remark
  editForm.data = { ...row.data } || {}
  editVisible.value = true
}

const handleDelete = async (row: UploadRecord) => {
  try {
    await ElMessageBox.confirm(
      t('uploadRecord.list.deleteConfirmMsg', [row.serialNo]),
      t('uploadRecord.list.deleteConfirmTitle'),
      {
        confirmButtonText: t('uploadRecord.list.deleteConfirmBtn'),
        cancelButtonText: t('uploadRecord.list.deleteCancelBtn'),
        type: 'warning'
      }
    )

    await UploadRecordApi.del(row.id)
    ElMessage.success(t('uploadRecord.list.deleteSuccess'))
    loadRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('uploadRecord.list.deleteFailed'))
    }
  }
}

const handleAction = (key: string, row: UploadRecord) => {
  if (key === 'detail') handleDetail(row)
  else if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      t('uploadRecord.list.batchDeleteConfirmMsg', [selectedRows.value.length]),
      t('uploadRecord.list.batchDeleteConfirmTitle'),
      {
        confirmButtonText: t('uploadRecord.list.deleteConfirmBtn'),
        cancelButtonText: t('uploadRecord.list.deleteCancelBtn'),
        type: 'warning'
      }
    )

    const ids = selectedRows.value.map(row => row.id)
    await UploadRecordApi.batchDelete(ids)
    ElMessage.success(t('uploadRecord.list.batchDeleteSuccess', [selectedRows.value.length]))
    selectedRows.value = []
    loadRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('uploadRecord.list.batchDeleteFailed'))
    }
  }
}

const handleCreate = () => {
  createForm.dataType = ''
  createForm.projectName = ''
  createForm.destPath = ''
  createForm.uploader = ''
  createForm.status = 'pending'
  createForm.remark = ''
  createForm.createdAt = ''
  createForm.data = {}
  fileSizeInputVal.value = 0
  fileSizeUnit.value = 'MB'
  createForm.fileSize = 0
  createVisible.value = true
}

const confirmCreate = async () => {
  if (!createForm.dataType) {
    ElMessage.error(t('uploadRecord.list.selectDiskLabel'))
    return
  }
  if (!createForm.destPath) {
    ElMessage.error(t('uploadRecord.list.enterDestPath'))
    return
  }
  if (!createForm.fileSize || fileSizeInputVal.value <= 0) {
    ElMessage.error(t('uploadRecord.list.enterFileSize'))
    return
  }
  if (!createForm.projectName) {
    ElMessage.error(t('uploadRecord.list.enterProjectName'))
    return
  }
  if (!createForm.uploader) {
    ElMessage.error(t('uploadRecord.list.enterUploader'))
    return
  }
  submitting.value = true
  try {
    const res = await UploadRecordApi.create({
      dataType: createForm.dataType,
      projectName: createForm.projectName,
      destPath: createForm.destPath,
      fileSize: createForm.fileSize,
      uploader: createForm.uploader,
      status: createForm.status,
      remark: createForm.remark || undefined,
      createdAt: createForm.createdAt || undefined,
      data: Object.keys(createForm.data).length > 0 ? createForm.data : undefined
    })
    if (res.code === 200) {
      ElMessage.success(t('uploadRecord.list.createSuccess'))
      createVisible.value = false
      loadRecords()
    } else {
      ElMessage.error(res.message || t('uploadRecord.list.createFailed'))
    }
  } catch (error) {
    ElMessage.error(t('uploadRecord.list.createFailed'))
  } finally {
    submitting.value = false
  }
}

const confirmEdit = async () => {
  try {
    submitting.value = true
    await UploadRecordApi.update({
      id: editForm.id,
      status: editForm.status as any,
      remark: editForm.remark,
      data: editForm.data
    })
    ElMessage.success(t('uploadRecord.list.updateSuccess'))
    editVisible.value = false
    loadRecords()
  } catch (error) {
    ElMessage.error(t('uploadRecord.list.updateFailed'))
  } finally {
    submitting.value = false
  }
}

// ==================== 批量导入 ====================
const showImportDialog = async () => {
  importDialogVisible.value = true
  importResult.value = null
  selectedFile.value = null
  fileList.value = []
  previewRowCount.value = null
  // 加载模板字段说明
  try {
    const res = await UploadRecordApi.getImportTemplate()
    if (res.code === 200 && res.data) {
      importTemplateFields.value = res.data.fields || []
    }
  } catch (e) {
    importTemplateFields.value = []
  }
}

const handleDownloadTemplate = async () => {
  downloadingTemplate.value = true
  try {
    await UploadRecordApi.downloadTemplate()
  } catch (e) {
    ElMessage.error(t('uploadRecord.list.templateDownloadFailed'))
  } finally {
    downloadingTemplate.value = false
  }
}

const handleFileChange = async (file: any) => {
  selectedFile.value = file.raw as File
  previewRowCount.value = null
  previewLoading.value = true
  try {
    const res = await UploadRecordApi.previewImport(selectedFile.value)
    if (res.code === 200) {
      previewRowCount.value = res.data.dataRows
    }
  } catch (e) {
    console.error('预览失败', e)
  } finally {
    previewLoading.value = false
  }
}

const handleFileRemove = () => {
  selectedFile.value = null
  previewRowCount.value = null
}

const handleImport = async () => {
  if (!selectedFile.value) {
    ElMessage.warning(t('uploadRecord.list.selectImportFile'))
    return
  }
  importing.value = true
  importResult.value = null
  try {
    const res = await UploadRecordApi.importRecords(selectedFile.value)
    if (res.code === 200) {
      importResult.value = res.data
      if (res.data.failed === 0) {
        ElMessage.success(t('uploadRecord.list.importSuccessMsg', [res.data.success]))
      } else {
        ElMessage.warning(t('uploadRecord.list.importPartial', [res.data.success, res.data.failed]))
      }
    } else {
      ElMessage.error(res.message || t('uploadRecord.list.importFailedMsg'))
    }
  } catch (e: any) {
    ElMessage.error(e.message || t('uploadRecord.list.importFailedMsg'))
  } finally {
    importing.value = false
  }
}

const closeImportDialog = () => {
  importDialogVisible.value = false
  if (importResult.value && importResult.value.success > 0) {
    loadRecords()
  }
}

onMounted(() => {
  loadColumnSettings()
  loadFieldConfigs()
  loadProjects()
  loadPersonnel()
  loadRecords()
})
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.page {
  width: 100%;
  min-height: 100vh;
  background: var(--color-page-bg);
  padding: var(--space-4);
  overflow-x: hidden;
}

/* ==================== 页面标题栏 ==================== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
  flex-wrap: wrap;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

.header-left {}

.page-title {
  font-family: 'Manrope', sans-serif;
  font-size: 17px;
  font-weight: 800;
  color: var(--color-text-primary);
  margin-bottom: 2px;
  letter-spacing: -0.3px;
}

.page-subtitle {
  font-size: 12px;
  color: var(--color-text-muted);
}

.header-stat {
  display: flex;
  align-items: baseline;
  gap: 4px;
  padding: 6px 14px;
  background: #f0f2f5;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  flex-shrink: 0;

  &:hover { background: #e6e8eb; }
}

.header-stat-num {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.header-stat-label {
  font-size: 12px;
  color: #909399;
}

.header-stat--success .header-stat-num { color: #67c23a; }
.header-stat--warning .header-stat-num { color: #e6a23c; }
.header-stat--info .header-stat-num { color: #409eff; }
.header-stat--danger .header-stat-num { color: #f56c6c; }
.header-stat--purple .header-stat-num { color: #8b5cf6; }

.header-actions {
  display: flex;
  gap: var(--space-2);

  :deep(.el-button) {
    font-size: 13px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    padding: 8px 16px;
  }
}

/* ==================== 筛选栏 ==================== */
.filter-card {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
  align-items: center;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  margin-bottom: var(--space-3);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);

  > * { flex: 0 0 auto; }
  .el-input,
  .el-select { width: 130px; }
  .el-date-editor { width: 220px !important; }
}

.filter-actions {
  margin-left: auto;
  display: flex;
  gap: var(--space-2);

  :deep(.el-button) {
    font-size: 13px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    padding: 6px 16px;
  }
}

/* ==================== 内容卡片 ==================== */
.table-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
}

.table-scroll-wrapper {
  overflow-x: auto;
  overflow-y: visible;

  &::-webkit-scrollbar { height: 6px; }
  &::-webkit-scrollbar-track { background: var(--color-surface-2); border-radius: 3px; }
  &::-webkit-scrollbar-thumb {
    background: var(--color-border);
    border-radius: 3px;
    &:hover { background: var(--color-border-hover, var(--color-border)); }
  }
}

/* ==================== 表格工具栏 ==================== */
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

.toolbar-left {
  .record-count { font-size: 13px; color: var(--color-text-secondary); strong { color: var(--color-text-primary); font-weight: 700; } }
  .selection-count { margin-left: var(--space-4); font-size: 13px; color: var(--color-danger); strong { font-weight: 700; } }
}

.toolbar-right { display: flex; gap: var(--space-2); }

.column-settings-btn {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-primary);
  background: var(--color-primary-light-9);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  padding: 6px 12px;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(0,94,235,0.12);
    color: var(--color-primary);
    border-color: var(--color-primary);
  }
}

/* ==================== 字段配置弹窗 ==================== */
.column-settings {
  .settings-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: var(--space-3);
    margin-bottom: var(--space-3);
    border-bottom: 1px solid var(--color-border-light);
  }

  .settings-title {
    font-family: 'Manrope', sans-serif;
    font-size: 14px;
    font-weight: 700;
    color: var(--color-primary);
  }

  .settings-list {
    max-height: 320px;
    overflow-y: auto;
    margin: 0 -12px;
    padding: 0 12px;
  }

  .settings-item {
    padding: 8px 10px;
    margin: 0 -10px;
    border-radius: var(--radius-sm);
    transition: background 0.15s ease;

    &:hover { background: var(--color-surface-2); }
    &.is-disabled { opacity: 0.7; }

    :deep(.el-checkbox) {
      width: 100%;
      .el-checkbox__label {
        font-size: 13px;
        color: var(--color-primary);
        font-weight: 500;
      }
      .el-checkbox__input.is-checked .el-checkbox__inner {
        background-color: var(--color-primary);
        border-color: var(--color-primary);
      }
      .el-checkbox__input.is-checked + .el-checkbox__label {
        color: var(--color-text-primary);
      }
    }
  }

  .settings-footer {
    margin-top: var(--space-3);
    padding-top: var(--space-3);
    border-top: 1px solid var(--color-border-light);
  }

  .settings-hint { font-size: 11px; color: var(--color-text-muted); }
}

/* ==================== 表格样式 ==================== */
:deep(.el-table) {
  --el-table-border-color: var(--color-border-light);
  width: 100% !important;
  table-layout: fixed;

  th.el-table__cell {
    background-color: var(--color-surface-3) !important;
    border-bottom: 1px solid var(--color-border) !important;
    padding: 12px 14px !important;
    color: var(--color-text-secondary) !important;
    font-weight: 700 !important;
    font-size: 11px !important;
    text-transform: uppercase;
    letter-spacing: 0.4px;
  }

  td.el-table__cell {
    padding: 10px 14px !important;
    font-size: 13px;
    color: var(--color-text-primary);
    border-bottom: 1px solid var(--color-border-light) !important;
  }

  .el-table__body tr:hover > td.el-table__cell {
    background-color: var(--color-primary-light-9) !important;
  }
}

.type-tag {
  display: inline-block;
  padding: 3px 10px;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid rgba(0,94,235,0.12);
}

.path-text {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  color: var(--color-primary);
  display: block;
  overflow: hidden;
  text-overflow: clip;
  white-space: nowrap;
  max-width: 280px;
}

.path-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  justify-content: center;
}

.no-data { color: var(--color-text-muted); font-style: italic; }

.copy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: all 0.2s;
  flex-shrink: 0;

  &:hover {
    background: var(--color-primary-light-9);
    color: var(--color-primary);
  }
  &:active { transform: scale(0.95); }
  &.copy-btn-lg { width: 26px; height: 26px; }
}

.detail-path-cell {
  display: flex;
  align-items: flex-start;
  gap: var(--space-2);
  width: 100%;
}

.detail-path-text {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  color: var(--color-text-secondary);
  word-break: break-all;
  flex: 1;
  line-height: 1.6;
}

.dynamic-field { font-size: 13px; color: var(--color-text-primary); font-weight: 500; }

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;

  &--success { background: rgba(34,197,94,0.1); color: var(--color-success); }
  &--warning { background: rgba(245,158,11,0.1); color: var(--color-warning); }
  &--danger { background: rgba(239,68,68,0.08); color: var(--color-danger); }
  &--info { background: rgba(59,130,246,0.1); color: var(--chart-blue); }
}

.serial-no {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  background: var(--color-primary-light-9);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  color: var(--color-primary);
  font-weight: 600;
  border: 1px solid rgba(0,94,235,0.12);
}

/* ==================== 分页 ==================== */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: var(--space-3) var(--space-4);
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface-2);
}

/* ==================== 导出侧边栏 ==================== */
.drawer-title-inner {
  display: flex;
  align-items: center;
  gap: 8px;
}

.drawer-title-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: rgba(34, 197, 94, 0.1);
  color: var(--color-success);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.drawer-title-text {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-text-primary);
}

:deep(.export-drawer) {
  .el-drawer__header {
    padding: 14px 16px;
    margin-bottom: 0;
    border-bottom: 1px solid var(--color-border-light);
    background: var(--color-surface);
  }

  .el-drawer__body {
    padding: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}

.export-drawer-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.export-drawer-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #22c55e;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.export-drawer-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.export-drawer-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.export-drawer-sub {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.export-drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;

  &::-webkit-scrollbar { width: 4px; }
  &::-webkit-scrollbar-thumb { background: #d1d5db; border-radius: 2px; }
}

/* 预览统计卡片 */
.export-preview-card {
  background: #f0fdf4;
  border: 1px solid #86efac;
  border-radius: 10px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.preview-stat {
  text-align: center;
}

.preview-num {
  font-family: 'Manrope', sans-serif;
  font-size: 28px;
  font-weight: 900;
  color: #16a34a;
  line-height: 1;
}

.preview-label {
  font-size: 12px;
  color: #16a34a;
  margin-top: 4px;
}

.preview-divider {
  width: 1px;
  height: 36px;
  background: rgba(22, 163, 74, 0.2);
}

.preview-hint-text {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #16a34a;

  .el-icon { font-size: 14px; }
}

/* 同步筛选区域 */
.export-sync-section {
  background: #fff;
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  padding: 12px 14px;
}

.sync-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 10px;
}

.sync-active-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

/* 筛选卡片 */
.export-filter-card,
.export-field-card {
  background: #fff;
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  overflow: hidden;
}

.filter-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  border-bottom: 1px solid var(--color-border-light);
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  background: var(--color-page-bg);

  .el-icon { color: var(--color-primary); }
}

.filter-card-body {
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filter-row {
  display: flex;
  flex-direction: column;
  gap: 4px;

  label {
    font-size: 11px;
    font-weight: 600;
    color: var(--color-text-secondary);
    letter-spacing: 0.3px;
  }
}

/* 字段选择 */
.field-check-list {
  padding: 12px 14px;

  .el-checkbox {
    margin-right: 0;
    margin-bottom: 8px;
  }
}

.field-check-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
  margin-top: 8px;
}

/* 底部 */
.export-drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 16px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

/* ==================== 批量导入弹窗 ==================== */
.import-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 弹窗头部 */
.diag-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-bottom: 14px;
  border-bottom: 1px solid #e5e7eb;
}
.diag-head-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #eff6ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #005eeb;
  flex-shrink: 0;
}
.diag-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.diag-head-title {
  font-size: 15px;
  font-weight: 700;
  color: #1f2937;
  font-family: 'Manrope', 'DM Sans', sans-serif;
  letter-spacing: -0.2px;
}
.diag-head-sub {
  font-size: 12px;
  color: #9ca3af;
}

/* 三步骤卡片 */
.step-row {
  display: grid;
  grid-template-columns: 1fr auto 1fr auto 1fr;
  gap: 0;
  align-items: stretch;
}
.step-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px 14px 12px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  position: relative;
}
.step-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 2px;
}
.step-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  font-family: 'Manrope', sans-serif;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  flex-shrink: 0;
}
.step-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.3;
}
.step-desc {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
  flex: 1;
}
.step-btn {
  width: 100%;
  margin-top: 4px;
}
.step-upload {
  width: 100%;
  margin-top: 4px;
}
.step-upload :deep(.el-upload) {
  width: 100%;
}
.step-upload :deep(.el-upload__trigger) {
  width: 100%;
}
.preview-hint {
  margin-top: 6px;
  font-size: 11px;
  color: var(--color-primary);
  font-weight: 500;
  padding: 4px 8px;
  background: var(--color-primary-light-9);
  border-radius: 4px;
}
.preview-hint.loading {
  color: var(--color-text-muted);
  background: var(--color-surface-3);
}

.step-arrow {
  display: flex;
  align-items: center;
  padding: 0 8px;
  color: var(--color-text-muted);
  flex-shrink: 0;
}

/* 字段填写说明 */
.field-guide {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}
.field-guide-head {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
  font-size: 12px;
  font-weight: 700;
  color: #6b7280;
  letter-spacing: 0.2px;
}
.field-guide-body { padding: 0; }

.fg-col-header {
  display: grid;
  grid-template-columns: 100px 60px 130px 1fr;
  gap: 8px;
  padding: 8px 14px;
  font-size: 11px;
  font-weight: 700;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.fg-row {
  display: grid;
  grid-template-columns: 100px 60px 130px 1fr;
  gap: 8px;
  padding: 9px 14px;
  align-items: center;
  border-bottom: 1px solid #f3f4f6;
  font-size: 12px;
  color: #374151;
  transition: background 0.15s;
  &:last-child { border-bottom: none; }
  &:hover { background: #f9fafb; }
}
.fg-name { font-weight: 600; color: #1f2937; font-size: 12px; }
.fg-example {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 11px;
  color: #005eeb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.fg-type { font-size: 11px; }

.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.2px;
}
.badge--danger { background: #fef2f2; color: #dc2626; }
.badge--gray { background: #f3f4f6; color: #9ca3af; }

.type-chip {
  display: inline-block;
  padding: 2px 7px;
  background: #f0f9ff;
  color: #0369a1;
  border-radius: 5px;
  font-size: 11px;
  font-weight: 500;
}

/* 导入结果 */
.import-result { }
.res-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 10px;
  border: 1px solid;
}
.res-banner--ok {
  background: #f0fdf4;
  border-color: #bbf7d0;
}
.res-banner--warn {
  background: #fffbeb;
  border-color: #fde68a;
}
.res-icon { flex-shrink: 0; display: flex; align-items: center; }
.res-info { display: flex; flex-direction: column; gap: 2px; }
.res-title { font-size: 14px; font-weight: 700; color: #1f2937; font-family: 'Manrope', sans-serif; }
.res-detail { font-size: 12px; color: #6b7280; display: flex; gap: 12px; }
.res-ok { color: #16a34a; font-weight: 600; }
.res-fail { color: #d97706; font-weight: 600; }

.fail-list { }
.fail-list-title {
  font-size: 11px;
  font-weight: 700;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  padding: 10px 14px 6px;
  border-top: 1px solid #e5e7eb;
  margin-top: 4px;
}
.fail-list-body { padding: 0 14px 12px; display: flex; flex-direction: column; gap: 4px; }
.fail-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 10px;
  background: #fff;
  border: 1px solid #f3f4f6;
  border-radius: 7px;
  font-size: 12px;
}
.fail-item-num {
  color: #dc2626;
  font-weight: 700;
  font-family: 'SF Mono', monospace;
  min-width: 44px;
  flex-shrink: 0;
}
.fail-item-msg { color: #6b7280; font-size: 12px; }

.diag-close-btn { }

/* ==================== 新增抽屉 ==================== */
/* 抽屉头部 */
.drawer-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px 0;
}
.drawer-head-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(0, 94, 235, 0.08);
  border: 1px solid rgba(0, 94, 235, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #005eeb;
  flex-shrink: 0;
}
.drawer-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.drawer-head-title {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 800;
  color: #1f2329;
  letter-spacing: -0.3px;
  line-height: 1.2;
}
.drawer-head-sub {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #9ca3af;
}

/* 表单主体 */
.create-body {
  padding: 20px 24px 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  flex: 1;
}

/* 字段分组 */
.field-group {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}

.field-group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: #f8f9fb;
  border-bottom: 1px solid #e5e7eb;
}

.field-group-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}

.field-group-hint {
  font-family: 'DM Sans', sans-serif;
  font-size: 10px;
  color: #9ca3af;
  font-weight: 500;
}

.field-group-badge {
  font-family: 'DM Sans', sans-serif;
  font-size: 10px;
  font-weight: 700;
  color: #005eeb;
  background: rgba(0, 94, 235, 0.08);
  padding: 1px 7px;
  border-radius: 99px;
}

.field-group-body {
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #ffffff;
}

/* 双列行 */
.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 16px;
}

.field-cell {
  display: flex;
  flex-direction: column;
  gap: 5px;
  &--full { grid-column: 1 / -1; }
}

.field-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #4b5563;
  line-height: 1;
}

.field-required {
  color: #dc2626;
  margin-left: 1px;
}

/* 表单控件 */
.dr-input, .dr-select, .dr-textarea {
  width: 100%;
}

:deep(.dr-input .el-input__wrapper),
:deep(.dr-select .el-input__wrapper) {
  background: #ffffff !important;
  border: 1px solid #d1d5db !important;
  border-radius: 7px !important;
  box-shadow: 0 0 0 1px #d1d5db inset !important;
  padding: 5px 11px !important;
  transition: border-color 0.2s, box-shadow 0.2s !important;
  font-family: 'DM Sans', sans-serif !important;
  font-size: 13px !important;
  &:hover {
    border-color: #9ca3af !important;
    box-shadow: 0 0 0 1px #9ca3af inset !important;
  }
  &.is-focus {
    border-color: #005eeb !important;
    box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.1) inset, 0 0 0 1px #005eeb inset !important;
  }
}
:deep(.dr-input .el-input__inner),
:deep(.dr-select .el-input__inner) {
  color: #1f2329 !important;
  font-size: 13px !important;
  font-weight: 500;
  &::placeholder { color: #c4c9d4 !important; }
}
:deep(.dr-input .el-input__prefix .el-icon),
:deep(.dr-select .el-input__prefix .el-icon) {
  color: #9ca3af;
}

/* 状态 select 样式 */
.dr-select--status {
  :deep(.el-select__rendered .el-select__placeholder) {
    color: #4b5563;
    font-weight: 500;
  }
}

/* 状态圆点 */
.status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-right: 6px;
  flex-shrink: 0;
  &--pending { background: #9ca3af; }
  &--processing { background: #005eeb; }
  &--completed { background: #16a34a; }
  &--failed { background: #dc2626; }
}

/* 文件大小行 */
.size-row {
  display: flex;
  gap: 6px;
  align-items: center;
}

.size-num {
  flex: 1;
  :deep(.el-input__wrapper) {
    background: #ffffff !important;
    border: 1px solid #d1d5db !important;
    border-radius: 7px !important;
    box-shadow: 0 0 0 1px #d1d5db inset !important;
    padding: 5px 8px !important;
    font-family: 'DM Sans', sans-serif !important;
    font-size: 13px !important;
    &:hover { border-color: #9ca3af !important; }
    &.is-focus { border-color: #005eeb !important; box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.1) inset !important; }
  }
  :deep(.el-input__inner) { color: #1f2329 !important; font-size: 13px !important; }
}

.size-unit {
  width: 68px;
  flex-shrink: 0;
  :deep(.el-input__wrapper) {
    background: #ffffff !important;
    border: 1px solid #d1d5db !important;
    border-radius: 7px !important;
    box-shadow: 0 0 0 1px #d1d5db inset !important;
    padding: 5px 8px !important;
    &:hover { border-color: #9ca3af !important; }
    &.is-focus { border-color: #005eeb !important; box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.1) inset !important; }
  }
  :deep(.el-input__inner) { color: #1f2329 !important; font-size: 12px !important; }
}

.size-hint {
  font-size: 11px;
  color: #9ca3af;
  font-family: 'SF Mono', 'Fira Code', monospace;
  margin-top: 3px;
}

/* 备注 */
.dr-textarea {
  :deep(.el-textarea__inner) {
    background: #ffffff !important;
    border: 1px solid #d1d5db !important;
    border-radius: 7px !important;
    box-shadow: 0 0 0 1px #d1d5db inset !important;
    padding: 8px 11px !important;
    font-family: 'DM Sans', sans-serif !important;
    font-size: 13px !important;
    color: #1f2329 !important;
    resize: none;
    transition: border-color 0.2s, box-shadow 0.2s !important;
    &:hover { border-color: #9ca3af !important; }
    &:focus {
      border-color: #005eeb !important;
      box-shadow: 0 0 0 3px rgba(0, 94, 235, 0.1) inset !important;
    }
    &::placeholder { color: #c4c9d4 !important; }
  }
}

/* 底部按钮 */
.drawer-foot {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 24px;
  border-top: 1px solid #f0f1f3;
  background: #ffffff;
  flex-shrink: 0;
}

.btn-cancel {
  height: 38px;
  padding: 0 20px;
  border: 1px solid #e5e7eb !important;
  border-radius: 8px !important;
  color: #4b5563 !important;
  font-family: 'DM Sans', sans-serif !important;
  font-size: 13px !important;
  font-weight: 600 !important;
  background: #ffffff !important;
  transition: all 0.2s;
  &:hover {
    border-color: #9ca3af !important;
    color: #1f2329 !important;
    background: #f9fafb !important;
  }
}

.btn-confirm {
  height: 38px;
  padding: 0 22px;
  border-radius: 8px !important;
  background: #005eeb !important;
  border: none !important;
  color: #ffffff !important;
  font-family: 'DM Sans', sans-serif !important;
  font-size: 13px !important;
  font-weight: 700 !important;
  letter-spacing: 0.3px;
  transition: all 0.2s;
  &:hover {
    background: #1a73e8 !important;
    transform: translateY(-1px);
    box-shadow: 0 4px 14px rgba(0, 94, 235, 0.3);
  }
  &:active { transform: scale(0.98); }
  &:disabled { opacity: 0.5; }
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 14px 20px;
  border-top: 1px solid var(--color-border-light);
  background: var(--color-surface);
  .el-button {
    border-radius: var(--radius-sm);
    font-weight: 600;
    font-size: 13px;
    min-width: 80px;
  }
}

.size-input-row { display: flex; gap: 6px; align-items: center; }
.size-num-input { flex: 1; }
.size-unit-select { width: 72px; flex-shrink: 0; }
.size-hint { margin-top: 4px; }
.size-converted { font-size: 11px; color: var(--color-text-muted); font-family: 'SF Mono', monospace; }

/* ==================== Element Plus 覆盖 ==================== */
:deep(.el-date-editor) {
  .el-input__wrapper {
    background: var(--color-surface) !important;
    border: 1px solid var(--color-border-light) !important;
    border-radius: var(--radius-sm) !important;
    box-shadow: none !important;
    padding: 5px 12px !important;
  }
  .el-range-separator { color: var(--color-text-muted) !important; font-size: 12px !important; }
  .el-range-input {
    color: var(--color-text-primary) !important;
    font-size: 12px !important;
    -webkit-font-smoothing: antialiased;
  }
}

:deep(.el-select) {
  .el-input__wrapper {
    background: var(--color-surface) !important;
    border: 1px solid var(--color-border-light) !important;
    border-radius: var(--radius-sm) !important;
    box-shadow: none !important;
    padding: 5px 12px !important;
  }
  .el-input__inner { font-size: 12px !important; -webkit-font-smoothing: antialiased; }
}

:deep(.el-button--primary) {
  background: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
  font-size: 12px !important;
  &:hover {
    background: var(--chart-blue) !important;
    border-color: var(--chart-blue) !important;
  }
}

:deep(.el-drawer) {
  .el-drawer__header {
    padding: 0 !important;
    margin-bottom: 0 !important;
    border-bottom: none !important;
    display: none !important;
  }
  .el-drawer__body { padding: 0 !important; }
}

/* ==================== 响应式 ==================== */
@media (max-width: 1366px) {
  .page { padding: var(--space-3); }
  .page-header { padding: var(--space-3) var(--space-4); }
  .filter-card { padding: var(--space-2) var(--space-3); gap: var(--space-2); margin-bottom: var(--space-2); }
  .filter-card .el-input,
  .filter-card .el-select { width: 120px; }
  .filter-card .el-date-editor { width: 200px !important; }
}

/* ==================== 对话框头部 ==================== */
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
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border: 1px solid rgba(0, 94, 235, 0.2);

  &--info {
    background: rgba(0, 176, 80, 0.1);
    color: #00b050;
    border-color: rgba(0, 176, 80, 0.2);
  }
}

.dialog-title-text {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

/* ==================== 编辑弹窗样式 ==================== */
.edit-dialog-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.edit-dialog-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.edit-dialog-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.edit-dialog-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.edit-dialog-sub {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-family: 'SF Mono', Monaco, monospace;
}

.edit-dialog-body {
  padding: 8px 4px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.edit-card {
  background: var(--color-page-bg);
  border-radius: 10px;
  overflow: hidden;
}

.edit-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid var(--color-border-light);
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);

  .el-icon {
    color: var(--color-primary);
  }
}

.edit-card-body {
  padding: 16px;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.edit-field {
  display: flex;
  flex-direction: column;
  gap: 6px;

  label {
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-secondary);
    letter-spacing: 0.3px;
  }

  .el-select,
  .el-input,
  .el-date-editor,
  .el-input-number,
  .el-textarea {
    width: 100%;
  }
}

.edit-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border-light);
}
</style>
