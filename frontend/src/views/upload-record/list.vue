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
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          导出
        </el-button>
        <el-button type="warning" @click="showImportDialog">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
          导入
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建
        </el-button>
        <el-button @click="loadRecords">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px">
            <path d="M23 4v6h-6M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          刷新
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
        <div class="export-preview-card" :class="{ 'is-loading': exportPreviewLoading }">
          <div class="preview-stat">
            <div class="preview-num">{{ exportPreviewCount.toLocaleString() }}</div>
            <div class="preview-label">条记录</div>
          </div>
          <div class="preview-divider"></div>
          <div class="preview-right">
            <div class="preview-hint-text">
              <el-icon><InfoFilled /></el-icon>
              <span v-if="syncCurrentFilter">已同步当前筛选</span>
              <span v-else-if="!exportForm.diskLabel && !exportForm.projectName && !exportForm.status && !exportForm.uploader && !exportForm.keyword && !exportForm.dateRange?.length">全部记录</span>
              <span v-else>已应用筛选条件</span>
            </div>
            <div class="preview-filter-tags" v-if="!syncCurrentFilter">
              <span class="preview-tag" v-if="exportForm.diskLabel">标签</span>
              <span class="preview-tag" v-if="exportForm.projectName">项目</span>
              <span class="preview-tag" v-if="exportForm.status">状态</span>
              <span class="preview-tag" v-if="exportForm.uploader">上传人</span>
              <span class="preview-tag" v-if="exportForm.keyword">关键词</span>
              <span class="preview-tag" v-if="exportForm.dateRange?.length">日期</span>
            </div>
            <el-button class="preview-refresh-btn" size="small" :loading="exportPreviewLoading" @click="refreshExportPreview">
              <el-icon v-if="!exportPreviewLoading"><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>

        <!-- 快速同步 -->
        <div class="export-sync-section">
          <div class="sync-header">
            <span>快速同步当前筛选</span>
            <el-switch v-model="syncCurrentFilter" size="small" />
          </div>
          <div class="sync-active-tags" v-if="syncCurrentFilter">
            <el-tag v-if="searchDiskLabel" size="small" closable @close="searchDiskLabel = ''; syncFilterToExport()">
              标签: {{ searchDiskLabel }}
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
              <el-input v-model="exportForm.diskLabel" placeholder="输入标签筛选" clearable size="small" />
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
    <el-dialog v-model="importDialogVisible" width="640px" destroy-on-close append-to-body class="import-dialog">
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

        <!-- 横向步骤条 -->
        <div class="stepper">
          <div class="step-item" :class="{ 'step-item--active': activeImportStep === 1, 'step-item--done': activeImportStep > 1 }">
            <div class="step-dot">
              <svg v-if="activeImportStep > 1" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
              <span v-else>1</span>
            </div>
            <span class="step-label">{{ t('uploadRecord.list.step1DownloadTemplate') }}</span>
          </div>
          <div class="step-line" :class="{ 'step-line--active': activeImportStep > 1 }"></div>
          <div class="step-item" :class="{ 'step-item--active': activeImportStep === 2, 'step-item--done': activeImportStep > 2 }">
            <div class="step-dot">
              <svg v-if="activeImportStep > 2" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
              <span v-else>2</span>
            </div>
            <span class="step-label">{{ t('uploadRecord.list.step2UploadFile') }}</span>
          </div>
          <div class="step-line" :class="{ 'step-line--active': activeImportStep > 2 }"></div>
          <div class="step-item" :class="{ 'step-item--active': activeImportStep === 3 }">
            <div class="step-dot">
              <span>3</span>
            </div>
            <span class="step-label">{{ t('uploadRecord.list.step3ConfirmImport') }}</span>
          </div>
        </div>

        <!-- 步骤内容区 -->
        <div class="step-content">
          <!-- Step 1: 下载模板 -->
          <div class="step-panel" v-if="activeImportStep === 1">
            <div class="step-panel__icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
            </div>
            <p class="step-panel__desc">{{ t('uploadRecord.list.step1Desc') }}</p>
            <el-button class="step-panel__btn" :loading="downloadingTemplate" @click="handleDownloadTemplate">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:5px"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              {{ t('uploadRecord.list.downloadXlsx') }}
            </el-button>
            <button class="step-panel__next" @click="activeImportStep = 2">
              {{ t('uploadRecord.list.nextStep') || '下一步' }}
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
            </button>
          </div>

          <!-- Step 2: 上传文件 -->
          <div class="step-panel" v-if="activeImportStep === 2">
            <!-- 统一上传区：始终渲染一个 el-upload -->
            <el-upload
              ref="uploadRef"
              class="drop-zone"
              action="#"
              :auto-upload="false"
              :limit="1"
              accept=".xlsx"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              :file-list="fileList"
              :drag="!selectedFile"
            >
              <!-- 无文件：显示拖拽区 -->
              <div v-if="!selectedFile" class="drop-zone__inner">
                <div class="drop-zone__icon">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                </div>
                <p class="drop-zone__text">{{ t('uploadRecord.list.dropHint') || '将 xlsx 文件拖到此处，或' }}</p>
                <span class="drop-zone__link">{{ t('uploadRecord.list.selectFile') }}</span>
              </div>
            </el-upload>

            <!-- 有文件：显示文件药片 + 替换按钮 -->
            <div v-if="selectedFile" class="file-actions">
              <div class="file-chip">
                <div class="file-chip__icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                </div>
                <div class="file-chip__info">
                  <span class="file-chip__name">{{ selectedFile.name }}</span>
                  <span class="file-chip__size">{{ (selectedFile.size / 1024).toFixed(1) }} KB</span>
                </div>
                <button class="file-chip__remove" @click="handleFileRemove" :title="t('uploadRecord.list.removeFile')">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
              </div>
              <!-- 点击触发隐藏的 file input 替换文件 -->
              <button class="replace-btn" @click="triggerFileReplace">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                {{ t('uploadRecord.list.replaceFile') || '重新选择文件' }}
              </button>
            </div>

            <div class="step-panel__actions">
              <button class="step-panel__back" @click="activeImportStep = 1">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                {{ t('uploadRecord.list.prevStep') || '上一步' }}
              </button>
              <div class="step-panel__right">
                <div v-if="previewLoading" class="preview-pill preview-pill--loading">
                  <span class="preview-dot"></span>
                  {{ t('uploadRecord.list.parsingFile') }}
                </div>
                <div v-else-if="previewRowCount !== null" class="preview-pill preview-pill--ok">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  {{ t('uploadRecord.list.recognizedRows', [previewRowCount]) }}
                </div>
                <button class="step-panel__next" :disabled="!selectedFile" @click="activeImportStep = 3">
                  {{ t('uploadRecord.list.nextStep') || '下一步' }}
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Step 3: 确认导入 -->
          <div class="step-panel" v-if="activeImportStep === 3 && !importing && !importingDone">
            <div class="step-panel__icon step-panel__icon--confirm">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="20 6 9 17 4 12"/></svg>
            </div>
            <p class="step-panel__desc">{{ t('uploadRecord.list.step3Desc') }}</p>
            <p class="step-panel__file-name" v-if="selectedFile">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
              {{ selectedFile.name }}
            </p>
            <div class="step-panel__actions">
              <button class="step-panel__back" @click="activeImportStep = 2">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                {{ t('uploadRecord.list.prevStep') || '上一步' }}
              </button>
              <el-button type="primary" class="import-btn" @click="handleImport">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>
                {{ t('uploadRecord.list.startImport') }}
              </el-button>
            </div>
          </div>

          <!-- 导入中 — 真实进度条 -->
          <div class="importing-state" v-if="importing">
            <div class="progress-wrap">
              <div class="progress-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>
              </div>
              <div class="progress-content">
                <div class="progress-header">
                  <span class="progress-label">{{ t('uploadRecord.list.importingText') || '正在导入数据...' }}</span>
                  <span class="progress-pct">{{ importProgress }}%</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-bar__fill" :style="{ width: importProgress + '%' }"></div>
                </div>
                <span class="progress-sub">{{ t('uploadRecord.list.importProgressTip') || '正在上传并处理 Excel 数据，请稍候' }}</span>
              </div>
            </div>
          </div>

          <!-- 导入结果 -->
          <div class="import-result" v-if="importingDone && importResult">
            <div class="res-banner" :class="importResult.failed > 0 ? 'res-banner--warn' : 'res-banner--ok'">
              <div class="res-icon">
                <svg v-if="importResult.failed === 0" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2.5"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
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

        <!-- 字段填写说明 -->
        <div class="field-guide" v-if="importTemplateFields.length > 0 && !importingDone">
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

      </div>

      <template #footer>
        <el-button class="diag-close-btn" @click="closeImportDialog">{{ t('uploadRecord.list.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- 工具栏 -->
    <div class="toolbar">
      <!-- 搜索框 -->
      <div class="toolbar__search">
        <div class="search-box">
          <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <input
            v-model="searchKeyword"
            class="search-input"
            :placeholder="t('uploadRecord.list.searchPlaceholder') || '搜索流水号、项目、上传人...'"
            @keyup.enter="handleSearch"
          />
          <span class="search-kbd">↵</span>
        </div>
      </div>

      <!-- 激活筛选标签 -->
      <div class="toolbar__chips" v-if="hasActiveFilters">
        <span v-if="searchStatus" class="filter-chip">
          {{ t('common.status') }}: {{ statusLabels[searchStatus] || searchStatus }}
          <button class="chip-remove" @click="searchStatus = ''; handleSearch()">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </span>
        <span v-if="searchDiskLabel" class="filter-chip">
          {{ t('uploadRecord.list.diskLabel') }}: {{ searchDiskLabel }}
          <button class="chip-remove" @click="searchDiskLabel = ''; handleSearch()">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </span>
        <span v-if="searchProjectName" class="filter-chip">
          {{ t('uploadRecord.list.project') }}: {{ searchProjectName }}
          <button class="chip-remove" @click="searchProjectName = ''; handleSearch()">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </span>
        <span v-if="searchUploader" class="filter-chip">
          {{ t('uploadRecord.list.uploader') }}: {{ searchUploader }}
          <button class="chip-remove" @click="searchUploader = ''; handleSearch()">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </span>
        <button class="chips-clear" @click="handleReset">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ t('common.clear') }}
        </button>
      </div>

      <!-- 右侧操作 -->
      <div class="toolbar__actions">
        <!-- 筛选器 -->
        <div class="filter-group">
          <el-popover placement="bottom-start" :width="360" trigger="click" :hide-after="0" :show-after="0" popper-class="filter-panel-popover">
            <template #reference>
              <button class="action-btn action-btn--ghost" :class="{ 'is-active': hasActiveFilters }">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
                {{ t('common.filter') }}
                <span v-if="activeFilterCount > 0" class="filter-badge">{{ activeFilterCount }}</span>
              </button>
            </template>
            <div class="filter-panel">
              <div class="filter-panel__header">
                <span class="fp-title">{{ t('common.filter') }}</span>
                <button class="fp-reset" @click="handleReset">{{ t('common.reset') }}</button>
              </div>
              <div class="filter-panel__body">
                <div class="fp-group">
                  <label class="fp-label">{{ t('common.status') }}</label>
                  <div class="fp-tags fp-tags--wrap">
                    <button class="fp-tag" :class="{ 'is-selected': searchStatus === 'pending' }" @click="searchStatus = searchStatus === 'pending' ? '' : 'pending'; handleSearch()">{{ statusLabels['pending'] }}</button>
                    <button class="fp-tag" :class="{ 'is-selected': searchStatus === 'processing' }" @click="searchStatus = searchStatus === 'processing' ? '' : 'processing'; handleSearch()">{{ statusLabels['processing'] }}</button>
                    <button class="fp-tag" :class="{ 'is-selected': searchStatus === 'completed' }" @click="searchStatus = searchStatus === 'completed' ? '' : 'completed'; handleSearch()">{{ statusLabels['completed'] }}</button>
                    <button class="fp-tag" :class="{ 'is-selected': searchStatus === 'failed' }" @click="searchStatus = searchStatus === 'failed' ? '' : 'failed'; handleSearch()">{{ statusLabels['failed'] }}</button>
                  </div>
                </div>
                <div class="fp-group" v-if="diskLabelOptions && diskLabelOptions.length > 0">
                  <label class="fp-label">{{ t('uploadRecord.list.diskLabel') }}</label>
                  <div class="fp-tags fp-tags--wrap">
                    <button v-for="d in diskLabelOptions" :key="d" class="fp-tag" :class="{ 'is-selected': searchDiskLabel === d }" @click="searchDiskLabel = searchDiskLabel === d ? '' : d; handleSearch()">{{ d }}</button>
                  </div>
                </div>
                <div class="fp-group" v-if="projectNameOptions && projectNameOptions.length > 0">
                  <label class="fp-label">{{ t('uploadRecord.list.project') }}</label>
                  <div class="fp-tags fp-tags--wrap">
                    <button v-for="p in projectNameOptions" :key="p" class="fp-tag" :class="{ 'is-selected': searchProjectName === p }" @click="searchProjectName = searchProjectName === p ? '' : p; handleSearch()">{{ p }}</button>
                  </div>
                </div>
              </div>
            </div>
          </el-popover>

          <!-- 字段显示 -->
          <el-popover placement="bottom-end" :width="280" trigger="click" :hide-after="0" :show-after="0" popper-class="col-panel-popover">
            <template #reference>
              <button class="action-btn action-btn--ghost">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="12" x2="14" y2="12"/><line x1="4" y1="18" x2="8" y2="18"/></svg>
                {{ t('common.columns') }}
              </button>
            </template>
            <div class="filter-panel">
              <div class="filter-panel__header">
                <span class="fp-title">{{ t('common.columns') }}</span>
                <button class="fp-reset" @click="handleResetColumns">{{ t('common.reset') }}</button>
              </div>
              <div class="filter-panel__body">
                <div v-for="col in allColumns" :key="col.prop" class="col-toggle" :class="{ 'is-disabled': col.required }">
                  <el-checkbox v-model="col.visible" :disabled="col.required" @change="handleColumnToggle">
                    {{ col.label }}
                  </el-checkbox>
                </div>
                <div v-for="col in allDynamicColumns" :key="col.prop" class="col-toggle">
                  <el-checkbox v-model="col.visible" @change="handleColumnToggle">
                    {{ col.label }}
                  </el-checkbox>
                </div>
              </div>
            </div>
          </el-popover>
        </div>

        <!-- 批量操作 -->
        <template v-if="selectedRows.length > 0">
          <el-dropdown trigger="click" @command="handleBatchAction">
            <button class="action-btn action-btn--ghost">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/></svg>
              {{ t('uploadRecord.list.batchEdit') }}
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-for="s in statusOptions" :key="s.value" :command="s.value">
                  <span class="batch-status-dot" :class="`dot--${s.value}`"></span>
                  {{ t('common.status') }}: {{ s.label }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <button class="action-btn action-btn--danger" @click="handleBatchDelete">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
            {{ t('common.batchDelete') }} ({{ selectedRows.length }})
          </button>
        </template>

        <!-- 新建 -->
        <button class="action-btn action-btn--primary" @click="handleCreate">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('common.create') }}
        </button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-card">
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
          @sort-change="handleSortChange"
        >
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column v-if="isColumnVisible('serialNo')" prop="serialNo" :label="t('uploadRecord.list.colSerialNo')" min-width="120">
          <template #default="{ row }">
            <span class="serial-cell">
              <el-tooltip :content="row.serialNo" placement="top">
                <span class="serial-text">{{ row.serialNo }}</span>
              </el-tooltip>
              <el-tooltip :content="t('uploadRecord.list.copySerialNo') || '复制流水号'" placement="top">
                <button class="copy-btn copy-btn--sm serial-copy" @click.stop="copyPath(row.serialNo)">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('diskLabel')" prop="diskLabel" :label="t('uploadRecord.list.colDiskLabel')" min-width="100" align="center" sortable="custom">
          <template #default="{ row }">
            <span class="disk-cell">
              <el-tooltip :content="row.diskLabel" placement="top">
                <span class="disk-text">{{ row.diskLabel }}</span>
              </el-tooltip>
              <el-tooltip :content="t('uploadRecord.list.copyDiskLabel') || '复制磁盘标签'" placement="top">
                <button class="copy-btn copy-btn--sm disk-copy" @click.stop="copyPath(row.diskLabel)">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('projectName')" prop="projectName" :label="t('uploadRecord.list.colProjectName')" min-width="110" show-overflow-tooltip sortable="custom" />
        <!-- 动态字段列 -->
        <el-table-column
          v-for="col in visibleDynamicColumns"
          :key="col.code"
          :prop="`data.${col.code}`"
          :label="col.name"
          min-width="100"
          align="center"
          show-overflow-tooltip
          sortable="custom"
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
        <el-table-column v-if="isColumnVisible('fileSize')" prop="fileSize" :label="t('uploadRecord.list.colFileSize')" min-width="90" align="center" sortable="custom">
          <template #default="{ row }">
            <span>{{ formatBytes(row.fileSize) }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('uploader')" prop="uploader" :label="t('uploadRecord.list.colUploader')" min-width="80" align="center" />
        <el-table-column v-if="isColumnVisible('status')" prop="status" :label="t('uploadRecord.list.colStatus')" min-width="70" align="center" sortable="custom">
          <template #default="{ row }">
            <span class="status-badge" :class="getStatusClass(row.status)">
              {{ row.statusText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('remark')" prop="remark" :label="t('uploadRecord.list.colRemark')" min-width="100" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('createdAt')" prop="createdAt" :label="t('uploadRecord.list.colTime')" min-width="160" sortable="custom" />
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
    <el-dialog v-model="detailVisible" width="560px" append-to-body class="detail-dialog">
      <template #header>
        <div class="detail-head">
          <span class="detail-head__icon">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
          </span>
          <span class="detail-head__title">记录详情</span>
        </div>
      </template>

      <!-- 内容区 -->
      <div v-if="currentRecord" class="detail-body">
        <!-- 流水号 -->
        <div class="detail-row detail-row--serial">
          <span class="detail-row__label">流水号</span>
          <span class="detail-row__value detail-row__value--copy">
            <span class="detail-row__code">{{ currentRecord.serialNo }}</span>
            <el-tooltip content="复制流水号" placement="top">
              <button class="copy-btn" @click.stop="copyPath(currentRecord.serialNo || '')">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
              </button>
            </el-tooltip>
          </span>
        </div>

        <!-- 状态行 -->
        <div class="detail-row">
          <span class="detail-row__label">状态</span>
          <span class="status-badge" :class="`status-badge--${currentRecord.status}`">{{ currentRecord.statusText }}</span>
        </div>

        <!-- 主信息 -->
        <div class="detail-row">
          <span class="detail-row__label">磁盘标签</span>
          <span class="detail-row__value detail-row__value--copy">
            <span class="detail-row__text">{{ currentRecord.diskLabel }}</span>
            <el-tooltip content="复制" placement="top">
              <button class="copy-btn" @click.stop="copyPath(currentRecord.diskLabel)">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
              </button>
            </el-tooltip>
          </span>
        </div>

        <div class="detail-row">
          <span class="detail-row__label">项目名称</span>
          <span class="detail-row__text">{{ currentRecord.projectName || '—' }}</span>
        </div>

        <div class="detail-row">
          <span class="detail-row__label">上传人</span>
          <span class="detail-row__text">{{ currentRecord.uploader }}</span>
        </div>

        <div class="detail-row">
          <span class="detail-row__label">文件大小</span>
          <span class="detail-row__text">{{ formatBytes(currentRecord.fileSize) }}</span>
        </div>

        <div class="detail-row">
          <span class="detail-row__label">创建时间</span>
          <span class="detail-row__text">{{ currentRecord.createdAt }}</span>
        </div>

        <div class="detail-row">
          <span class="detail-row__label">更新时间</span>
          <span class="detail-row__text">{{ currentRecord.updatedAt }}</span>
        </div>

        <div class="detail-row detail-row--path">
          <span class="detail-row__label">目标路径</span>
          <span class="detail-row__value detail-row__value--path">
            <span class="detail-row__text detail-row__text--path">{{ currentRecord.destPath || '—' }}</span>
            <el-tooltip v-if="currentRecord.destPath" content="复制路径" placement="top">
              <button class="copy-btn" @click.stop="copyPath(currentRecord.destPath)">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
              </button>
            </el-tooltip>
          </span>
        </div>

        <div v-if="currentRecord.remark" class="detail-row">
          <span class="detail-row__label">备注</span>
          <span class="detail-row__text">{{ currentRecord.remark }}</span>
        </div>
      </div>
    </el-dialog>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" width="520px" destroy-on-close append-to-body class="edit-dialog">
      <template #header>
        <div class="edit-head">
          <span class="edit-head__icon">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
          </span>
          <span class="edit-head__title">编辑记录</span>
        </div>
      </template>
      <!-- 表单 -->
      <div class="edit-body">
        <div class="edit-row">
          <label class="edit-row__label">状态</label>
          <el-select v-model="editForm.status" placeholder="选择状态">
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
          </el-select>
        </div>
        <div class="edit-row">
          <label class="edit-row__label">文件大小</label>
          <div class="edit-size">
            <el-input-number v-model="editFileSizeVal" :min="0" :precision="3" controls-position="right" @change="syncEditFileSizeFromInput" />
            <el-select v-model="editFileSizeUnit" @change="syncEditFileSizeFromInput">
              <el-option label="B" value="B" />
              <el-option label="KB" value="KB" />
              <el-option label="MB" value="MB" />
              <el-option label="GB" value="GB" />
              <el-option label="TB" value="TB" />
            </el-select>
          </div>
          <span class="edit-hint" v-if="editFileSizeVal > 0">≈ {{ formatSizeInOtherUnits(editFileSizeVal, editFileSizeUnit) }}</span>
        </div>
        <div class="edit-row edit-row--full">
          <label class="edit-row__label">备注</label>
          <el-input v-model="editForm.remark" type="textarea" :rows="2" placeholder="添加备注" />
        </div>
        <template v-if="dynamicColumns.length > 0">
          <div class="edit-divider"></div>
          <div class="edit-row" v-for="col in dynamicColumns" :key="col.code">
            <label class="edit-row__label">{{ col.name }}</label>
            <el-select v-if="col.type === 'select'" v-model="editForm.data[col.code]" clearable :placeholder="col.placeholder || '请选择'">
              <el-option v-for="opt in col.options" :key="opt" :label="opt" :value="opt" />
            </el-select>
            <el-date-picker v-else-if="col.type === 'date'" v-model="editForm.data[col.code]" type="date" value-format="YYYY-MM-DD" :placeholder="col.placeholder || '选择日期'" popper-class="date-picker-popper" />
            <el-date-picker v-else-if="col.type === 'datetime'" v-model="editForm.data[col.code]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" :placeholder="col.placeholder || '选择时间'" popper-class="date-picker-popper" />
            <el-input-number v-else-if="col.type === 'number'" v-model="editForm.data[col.code]" :placeholder="col.placeholder || '输入数字'" />
            <el-input v-else v-model="editForm.data[col.code]" :placeholder="col.placeholder || '输入内容'" />
          </div>
        </template>
      </div>

      <!-- 底部 -->
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmEdit">保存</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新增上传记录弹窗 -->
    <el-dialog
      v-model="createVisible"
      width="560px"
      :close-on-click-modal="false"
      append-to-body
      class="create-dialog"
    >
      <template #header>
        <div class="dialog-head">
          <div class="dialog-head-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="12" y1="18" x2="12" y2="12"/>
              <line x1="9" y1="15" x2="15" y2="15"/>
            </svg>
          </div>
          <span class="dialog-head-title">新建上传记录</span>
        </div>
      </template>

      <!-- 表单区域 -->
      <div class="create-form">

        <!-- 提示卡片 -->
        <div class="form-hint">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
          <span><em>*</em> 为必填项</span>
        </div>

        <!-- 基础信息 -->
        <div class="form-section">
          <div class="form-section__label">基础信息</div>
          <div class="form-row">
            <div class="form-item">
              <label class="form-item__label">磁盘标签 <em>*</em></label>
              <el-input v-model="createForm.diskLabel" :placeholder="t('uploadRecord.list.createDiskLabelPlaceholder')" clearable />
            </div>
            <div class="form-item">
              <label class="form-item__label">项目名称 <em>*</em></label>
              <el-select v-model="createForm.projectName" :placeholder="t('uploadRecord.list.createProjectNamePlaceholder')" clearable filterable allow-create default-first-option :reserve-keyword="false">
                <el-option v-for="p in projectList" :key="p.id" :label="p.name" :value="p.name" />
              </el-select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-item">
              <label class="form-item__label">上传人 <em>*</em></label>
              <el-select v-model="createForm.uploader" :placeholder="t('uploadRecord.list.createUploaderPlaceholder')" filterable allow-create default-first-option :reserve-keyword="false">
                <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
              </el-select>
            </div>
            <div class="form-item">
              <label class="form-item__label">状态 <em>*</em></label>
              <el-select v-model="createForm.status">
                <el-option :label="t('status.pending')" value="pending"><span class="status-dot status-dot--pending"></span> {{ t('status.pending') }}</el-option>
                <el-option :label="t('status.processing')" value="processing"><span class="status-dot status-dot--processing"></span> {{ t('status.processing') }}</el-option>
                <el-option :label="t('status.completed')" value="completed"><span class="status-dot status-dot--completed"></span> {{ t('status.completed') }}</el-option>
                <el-option :label="t('status.failed')" value="failed"><span class="status-dot status-dot--failed"></span> {{ t('status.failed') }}</el-option>
              </el-select>
            </div>
          </div>
        </div>

        <!-- 文件信息 -->
        <div class="form-section">
          <div class="form-section__label">文件信息</div>
          <div class="form-row">
            <div class="form-item form-item--full">
              <label class="form-item__label">目标路径</label>
              <el-input v-model="createForm.destPath" :placeholder="t('uploadRecord.list.createDestPathPlaceholder')" clearable />
            </div>
          </div>
          <div class="form-row">
            <div class="form-item">
              <label class="form-item__label">文件大小</label>
              <div class="size-input">
                <el-input-number v-model="fileSizeInputVal" :min="0" :precision="3" controls-position="right" @change="syncFileSizeFromInput" />
                <el-select v-model="fileSizeUnit" @change="syncFileSizeFromInput">
                  <el-option label="B" value="B" /><el-option label="KB" value="KB" /><el-option label="MB" value="MB" /><el-option label="GB" value="GB" /><el-option label="TB" value="TB" />
                </el-select>
              </div>
              <span class="form-hint-inline" v-if="fileSizeInputVal > 0">≈ {{ formatSizeInOtherUnits(fileSizeInputVal, fileSizeUnit) }}</span>
            </div>
          </div>
          <div class="form-row">
            <div class="form-item form-item--date">
              <label class="form-item__label">创建时间</label>
              <el-date-picker
                v-model="createForm.createdAt"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                :placeholder="t('uploadRecord.list.datePlaceholder')"
                popper-class="date-picker-popper"
              />
            </div>
          </div>
        </div>

        <!-- 备注 -->
        <div class="form-section">
          <div class="form-section__label">备注</div>
          <el-input v-model="createForm.remark" type="textarea" :rows="2" :placeholder="t('uploadRecord.list.createRemarkPlaceholder')" show-word-limit maxlength="500" />
        </div>

      </div>

      <!-- 底部按钮 -->
      <template #footer>
        <div class="dialog-foot">
          <el-button @click="createVisible = false">{{ t('uploadRecord.list.createCancel') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmCreate">{{ t('uploadRecord.list.createConfirm') }}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, inject, nextTick, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { InfoFilled, Filter, Grid, Download } from '@element-plus/icons-vue'
import { UploadRecordApi, type UploadRecord, type ImportTemplateField, type ImportResultResp, type UploadRecordStatistics } from '@/api/upload-record'
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
  { prop: 'diskLabel', label: t('uploadRecord.list.colDiskLabel'), visible: true, required: true },
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
      // 恢复固定列设置
      allColumns.value.forEach(col => {
        const savedCol = savedConfig.find((c: ColumnConfig) => c.prop === col.prop)
        if (savedCol && !col.required) {
          col.visible = savedCol.visible
        }
      })
      // 恢复动态列设置
      allDynamicColumns.value.forEach(col => {
        const savedCol = savedConfig.find((c: ColumnConfig) => c.prop === col.prop)
        if (savedCol) {
          col.visible = savedCol.visible
        }
      })
    } catch (e) {
      console.error('Failed to load column settings:', e)
    }
  }
}

const saveColumnSettings = () => {
  const config = [
    ...allColumns.value.map(col => ({
      prop: col.prop,
      label: col.label,
      visible: col.visible,
      required: col.required
    })),
    ...allDynamicColumns.value.map(col => ({
      prop: col.prop,
      label: col.label,
      visible: col.visible,
      required: false
    }))
  ]
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
  allDynamicColumns.value.forEach(col => {
    col.visible = true
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
const importingDone = ref(false)
const importProgress = ref(0)
const activeImportStep = ref(1)
const downloadingTemplate = ref(false)
const selectedFile = ref<File | null>(null)
const fileList = ref<any[]>([])
const uploadRef = ref()
const importTemplateFields = ref<ImportTemplateField[]>([])
const importResult = ref<ImportResultResp | null>(null)
const previewRowCount = ref<number | null>(null)
const previewLoading = ref(false)

// 搜索激活状态
const searchActive = ref(false)

// options 字段是逗号分隔字符串，需要解析后展示
const getOptionsArr = (opts: string | undefined) => {
  if (!opts) return []
  return String(opts).split(',').map(s => s.trim()).filter(Boolean)
}
const searchDiskLabel = ref('')
const searchStatus = ref('')
const searchUploader = ref('')
const searchKeyword = ref('')
const searchDateRange = ref<string[]>([])

// 筛选相关
const statusLabels: Record<string, string> = {
  pending: '待处理',
  processing: '处理中',
  completed: '已完成',
  failed: '失败'
}

const diskLabelOptions = computed(() => {
  const labels = new Set<string>()
  tableData.value.forEach(r => { if (r.diskLabel) labels.add(r.diskLabel) })
  return Array.from(labels).slice(0, 20)
})

const projectNameOptions = computed(() => {
  const names = new Set<string>()
  tableData.value.forEach(r => { if (r.projectName) names.add(r.projectName) })
  return Array.from(names).slice(0, 20)
})

const hasActiveFilters = computed(() => !!(searchStatus.value || searchDiskLabel.value || searchProjectName.value || searchUploader.value || searchKeyword.value))
const activeFilterCount = computed(() => {
  let count = 0
  if (searchStatus.value) count++
  if (searchDiskLabel.value) count++
  if (searchProjectName.value) count++
  if (searchUploader.value) count++
  return count
})

const handleBatchAction = (status: string) => {
  handleBatchUpdateStatus(status)
}

// 排序
const sortField = ref('')
const sortOrder = ref<'asc' | 'desc'>('')

const tableData = ref<UploadRecord[]>([])
const dynamicColumns = ref<FieldConfig[]>([])
const projectList = ref<ProjectSimple[]>([])
const personnelList = ref<Personnel[]>([])
const searchProjectName = ref('')
const searchField = ref('')
const searchPlaceholder = computed(() => {
  const map: Record<string, string> = {
    serialNo: '搜索流水号',
    diskLabel: '搜索磁盘标签',
    projectName: '搜索项目名称',
    destPath: '搜索目标路径',
    uploader: '搜索上传人',
    status: '搜索状态',
    remark: '搜索备注'
  }
  return map[searchField.value] || t('uploadRecord.list.searchKeywordPlaceholder')
})
const detailVisible = ref(false)
const editVisible = ref(false)
const createVisible = ref(false)
const currentRecord = ref<UploadRecord | null>(null)
const tableRef = ref()

// 统计（全量数据，不受筛选条件影响）
const listStats = reactive<UploadRecordStatistics>({
  todayCount: 0, todaySize: 0, todaySizeStr: '0 B',
  weekCount: 0, weekSize: 0, weekSizeStr: '—',
  monthCount: 0, monthSize: 0, monthSizeStr: '—',
  totalCount: 0, totalSize: 0, totalSizeStr: '0 B',
  trend: [], byStatus: [], byDiskLabel: [], byProject: []
})
const loadListStats = async () => {
  try {
    const res = await UploadRecordApi.statistics() as { data: UploadRecordStatistics }
    const d = res.data
    Object.assign(listStats, d)
  } catch (e) { /* ignore */ }
}
const totalSize = computed(() => listStats.totalSize)
const statusCount = computed(() => {
  const s = { completed: 0, pending: 0, processing: 0, failed: 0 }
  for (const item of (listStats.byStatus || [])) {
    if (item.status in s) (s as any)[item.status] = item.count
  }
  return s
})
const filterByStatus = (status: string) => {
  if (searchStatus.value === status) {
    searchStatus.value = ''
  } else {
    searchStatus.value = status
    // Clear other filters to avoid conflicting filter conditions
    searchDiskLabel.value = ''
    searchUploader.value = ''
    searchProjectName.value = ''
    searchKeyword.value = ''
    searchField.value = ''
    searchDateRange.value = []
    searchActive.value = false
  }
  loadRecords()
}
const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[Math.min(i, units.length - 1)]
}

// 文件大小工具
const fileSizeUnit = ref('GB')
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
    fileSizeUnit.value = 'GB'
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

// 编辑弹窗文件大小
const editFileSizeVal = ref(0)
const editFileSizeUnit = ref('MB')

const syncEditFileSizeFromInput = () => {
  editForm.fileSize = editFileSizeVal.value * BYTE_UNITS[editFileSizeUnit.value]
}

// 初始化编辑弹窗的文件大小显示
const initEditFileSize = (bytes: number) => {
  if (bytes === 0) {
    editFileSizeVal.value = 0
    editFileSizeUnit.value = 'MB'
    return
  }
  const units = ['TB', 'GB', 'MB', 'KB', 'B']
  for (const unit of units) {
    const quotient = bytes / BYTE_UNITS[unit]
    if (quotient >= 1) {
      editFileSizeVal.value = Math.round(quotient * 1000) / 1000
      editFileSizeUnit.value = unit
      return
    }
  }
  editFileSizeVal.value = bytes
  editFileSizeUnit.value = 'B'
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

const editForm = reactive<{ id: number; serialNo: string; status: string; remark: string; fileSize: number; data: Record<string, any> }>({
  id: 0,
  serialNo: '',
  status: '',
  remark: '',
  fileSize: 0,
  data: {}
})

const exportForm = reactive({
  diskLabel: '',
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
const exportPreviewCount = ref(0)
const exportPreviewLoading = ref(false)

const getExportParams = () => {
  const p: Record<string, any> = {}
  if (exportForm.diskLabel) p.diskLabel = exportForm.diskLabel
  if (exportForm.projectName) p.projectName = exportForm.projectName
  if (exportForm.status) p.status = exportForm.status
  if (exportForm.uploader) p.uploader = exportForm.uploader
  if (exportForm.keyword) p.keyword = exportForm.keyword
  if (exportForm.dateRange && exportForm.dateRange.length === 2) {
    p.startDate = exportForm.dateRange[0]
    p.endDate = exportForm.dateRange[1]
  }
  return p
}

// 刷新导出预览计数
const refreshExportPreview = async () => {
  exportPreviewLoading.value = true
  try {
    const stats = await UploadRecordApi.statistics(getExportParams())
    exportPreviewCount.value = (stats as { data: { totalCount: number } }).data.totalCount ?? 0
  } catch {
    // ignore
  } finally {
    exportPreviewLoading.value = false
  }
}

const getStatusText = (status: string) => {
  return t(`status.${status}`)
}

// 状态选项列表
const statusOptions = computed(() => [
  { value: 'pending', label: t('status.pending') },
  { value: 'processing', label: t('status.processing') },
  { value: 'completed', label: t('status.completed') },
  { value: 'failed', label: t('status.failed') }
])

// 同步当前筛选条件到导出表单
const syncFilterToExport = () => {
  if (syncCurrentFilter.value) {
    exportForm.diskLabel = searchDiskLabel.value
    exportForm.projectName = searchProjectName.value
    exportForm.status = searchStatus.value
    exportForm.uploader = searchUploader.value
    exportForm.dateRange = searchDateRange.value ? [...searchDateRange.value] : []
    if (!searchField.value) {
      exportForm.keyword = searchKeyword.value
    }
  }
  refreshExportPreview()
}

// 导出字段全选切换
const handleExportAllChange = (val: boolean) => {
  allColumns.value.forEach(c => { if (c.prop !== 'data') c.visible = val })
}

const showExportDialog = () => {
  // 重置导出表单
  exportForm.diskLabel = ''
  exportForm.projectName = ''
  exportForm.status = ''
  exportForm.uploader = ''
  exportForm.dateRange = []
  exportForm.keyword = ''
  exportAllFields.value = true
  syncCurrentFilter.value = true
  exportPreviewCount.value = pagination.total
  // 同步当前筛选
  syncFilterToExport()
  // 异步刷新预览计数
  refreshExportPreview()
  exportDialogVisible.value = true
}

// 自动刷新预览：当同步关闭且手动改了导出筛选时，延迟刷新计数
watch(() => [exportForm.diskLabel, exportForm.projectName, exportForm.status, exportForm.uploader, exportForm.keyword, exportForm.dateRange], () => {
  if (!syncCurrentFilter.value) {
    const t = setTimeout(() => refreshExportPreview(), 400)
    return () => clearTimeout(t)
  }
})

const handleExport = async () => {
  exporting.value = true
  try {
    const params: any = {}
    if (exportForm.diskLabel) params.diskLabel = exportForm.diskLabel
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
  diskLabel: '',
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
    // 优先使用 Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(path)
      ElMessage.success(t('uploadRecord.list.copyPathSuccess'))
      return
    }
    // Fallback: 使用 textarea 复制（兼容 HTTP 环境）
    const textarea = document.createElement('textarea')
    textarea.value = path
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.select()
    const success = document.execCommand('copy')
    document.body.removeChild(textarea)
    if (success) {
      ElMessage.success(t('uploadRecord.list.copyPathSuccess'))
    } else {
      ElMessage.error(t('uploadRecord.list.copyPathFailed'))
    }
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
    const params: Record<string, any> = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      startDate: searchDateRange.value?.[0] || undefined,
      endDate: searchDateRange.value?.[1] || undefined,
      sortField: sortField.value || undefined,
      sortOrder: sortOrder.value || undefined,
      diskLabel: searchDiskLabel.value || undefined,
      status: searchStatus.value || undefined,
      uploader: searchUploader.value || undefined,
      projectName: searchProjectName.value || undefined
    }
    // 字段精确搜索：searchField 指定字段，keyword 填入对应字段
    if (searchField.value && searchKeyword.value.trim()) {
      ;(params as any)[searchField.value] = searchKeyword.value.trim()
    } else if (!searchField.value && searchKeyword.value.trim()) {
      // 未选字段时，keyword 搜索所有字段
      ;(params as any).keyword = searchKeyword.value.trim()
    }

    const res = await UploadRecordApi.list(params)
    tableData.value = res.data.items || []
    pagination.total = res.data.total || 0

    // 更新数据类型列表
    const stats = await UploadRecordApi.statistics()
  } catch (error) {
    console.error('Failed to load records:', error)
    ElMessage.error(t('uploadRecord.list.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    searchActive.value = true
  }
  pagination.page = 1
  loadRecords()
}

const getDynamicFieldType = (fieldCode: string) => {
  const col = dynamicColumns.value.find(c => c.code === fieldCode)
  return col?.type || 'text'
}

const getDynamicFieldOptions = (fieldCode: string) => {
  const col = dynamicColumns.value.find(c => c.code === fieldCode)
  if (!col?.options) return []
  return String(col.options).split(',').map(s => s.trim()).filter(Boolean)
}

const handleReset = () => {
  searchActive.value = false
  searchStatus.value = ''
  searchUploader.value = ''
  searchProjectName.value = ''
  searchKeyword.value = ''
  searchField.value = ''
  searchDateRange.value = []
  sortField.value = ''
  sortOrder.value = ''
  pagination.page = 1
  loadRecords()
}

const handleSortChange = ({ prop, order }: { prop: string; order: string }) => {
  sortField.value = prop || ''
  sortOrder.value = order === 'ascending' ? 'asc' : order === 'descending' ? 'desc' : ''
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
  editForm.fileSize = row.fileSize
  editForm.data = { ...row.data } || {}
  initEditFileSize(row.fileSize)
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
    loadListStats()
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
    loadListStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('uploadRecord.list.batchDeleteFailed'))
    }
  }
}

// 批量修改状态
const handleBatchUpdateStatus = async (status: string) => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      t('uploadRecord.list.batchUpdateStatusConfirmMsg', [selectedRows.value.length, getStatusText(status)]),
      t('uploadRecord.list.batchUpdateStatusConfirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const ids = selectedRows.value.map(row => row.id)
    await UploadRecordApi.batchUpdateStatus(ids, status as any)
    ElMessage.success(t('uploadRecord.list.batchUpdateStatusSuccess', [selectedRows.value.length, getStatusText(status)]))
    selectedRows.value = []
    loadRecords()
    loadListStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(t('uploadRecord.list.batchUpdateStatusFailed'))
    }
  }
}

const handleCreate = () => {
  createForm.diskLabel = ''
  createForm.projectName = ''
  createForm.destPath = ''
  createForm.uploader = ''
  createForm.status = 'pending'
  createForm.remark = ''
  createForm.createdAt = ''
  createForm.data = {}
  fileSizeInputVal.value = 0
  fileSizeUnit.value = 'GB'
  createForm.fileSize = 0
  createVisible.value = true
}

const confirmCreate = async () => {
  if (!createForm.diskLabel) {
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
      diskLabel: createForm.diskLabel,
      projectName: createForm.projectName,
      destPath: createForm.destPath,
      fileSize: createForm.fileSize, // 直接发送浮点数，后端统一处理四舍五入
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
      loadListStats()
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
      fileSize: editForm.fileSize > 0 ? editForm.fileSize : undefined, // 直接发送浮点数，后端统一处理四舍五入
      data: editForm.data
    })
    ElMessage.success(t('uploadRecord.list.updateSuccess'))
    editVisible.value = false
    loadRecords()
    loadListStats()
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
  activeImportStep.value = 1
  importingDone.value = false
  importProgress.value = 0
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
  activeImportStep.value = 2
  // 清空 el-upload 内部 file list，确保同名文件也能触发 onChange
  nextTick(() => {
    uploadRef.value?.clearFiles()
  })
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
  uploadRef.value?.clearFiles()
  fileList.value = []
  activeImportStep.value = 2
}

// 触发 file input 点击以替换文件（清空后点击，清除内部状态确保重新触发 onChange）
const triggerFileReplace = () => {
  uploadRef.value?.clearFiles()
  selectedFile.value = null
  previewRowCount.value = null
  nextTick(() => {
    const el = uploadRef.value?.$el as HTMLElement
    const input = el?.querySelector?.('[type="file"]') as HTMLInputElement | undefined
    input?.click()
  })
}

const handleImport = async () => {
  if (!selectedFile.value) {
    ElMessage.warning(t('uploadRecord.list.selectImportFile'))
    return
  }
  activeImportStep.value = 3
  importing.value = true
  importingDone.value = false
  importResult.value = null
  importProgress.value = 0
  try {
    const res = await UploadRecordApi.importRecords(selectedFile.value, (pct) => {
      importProgress.value = pct
    })
    if (res.code === 200) {
      importResult.value = res.data
      importProgress.value = 100
      if (res.data.failed === 0) {
        ElMessage.success(t('uploadRecord.list.importSuccessMsg', [res.data.success]))
      } else {
        ElMessage.warning(t('uploadRecord.list.importPartial', [res.data.success, res.data.failed]))
      }
      // 成功后刷新列表和统计
      loadRecords()
      loadListStats()
    } else {
      ElMessage.error(res.message || t('uploadRecord.list.importFailedMsg'))
    }
  } catch (e: any) {
    ElMessage.error(e.message || t('uploadRecord.list.importFailedMsg'))
  } finally {
    importing.value = false
    importingDone.value = true
  }
}

const closeImportDialog = () => {
  importDialogVisible.value = false
  importingDone.value = false
  importProgress.value = 0
  activeImportStep.value = 1
  selectedFile.value = null
  fileList.value = []
  importResult.value = null
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
  loadListStats()
})
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.page {
  width: 100%;
  min-height: 100vh;
  background: var(--color-page-bg);
  padding: var(--space-4);
}

/* ==================== 页面标题栏 ==================== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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
  align-items: center;
  flex-wrap: wrap;
  gap: var(--space-2);

  :deep(.el-button) {
    font-size: 13px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    padding: 7px 14px;
    white-space: nowrap;
  }
}

/* ==================== 工具栏 ==================== */
.toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: 12px 16px;
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  flex-wrap: wrap;
  min-height: 52px;
  margin-bottom: 0;

  &__search { display: flex; align-items: center; }
  &__chips { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; flex: 1; }
  &__actions { display: flex; align-items: center; gap: 6px; margin-left: auto; flex-shrink: 0; }
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--gray-50);
  border: 1.5px solid var(--gray-200);
  border-radius: 8px;
  padding: 0 10px;
  height: 34px;
  width: 280px;
  transition: all 0.15s;
  gap: 6px;

  &:focus-within {
    border-color: var(--color-primary);
    background: #fff;
    box-shadow: 0 0 0 3px rgba(27, 58, 138, 0.08);
  }

  .search-icon { color: var(--gray-400); flex-shrink: 0; }
}

.search-input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 13px;
  color: var(--color-text-primary);
  flex: 1;
  font-family: inherit;
  &::placeholder { color: var(--gray-400); }
}

.search-kbd {
  font-size: 10px;
  color: var(--gray-400);
  background: var(--gray-100);
  border: 1px solid var(--gray-200);
  border-radius: 4px;
  padding: 1px 5px;
  font-family: 'Manrope', monospace;
  flex-shrink: 0;
  line-height: 1.4;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px 3px 10px;
  background: var(--color-primary-light-9);
  border: 1px solid rgba(27, 58, 138, 0.15);
  border-radius: 20px;
  font-size: 12px;
  color: var(--color-primary);
  font-weight: 500;
  white-space: nowrap;
  animation: chipIn 0.15s ease;
}

@keyframes chipIn {
  from { opacity: 0; transform: scale(0.9); }
  to { opacity: 1; transform: scale(1); }
}

.chip-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: none;
  background: rgba(27, 58, 138, 0.12);
  border-radius: 50%;
  cursor: pointer;
  color: var(--color-primary);
  padding: 0;
  transition: background 0.12s;

  &:hover { background: rgba(27, 58, 138, 0.25); }
}

.chips-clear {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border: none;
  background: transparent;
  color: var(--gray-500);
  font-size: 12px;
  cursor: pointer;
  border-radius: 20px;
  font-family: inherit;
  transition: all 0.12s;

  &:hover { background: var(--gray-100); color: var(--color-text-primary); }
}

.filter-group { display: flex; align-items: center; gap: 4px; }

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  border: 1.5px solid transparent;
  font-family: inherit;
  white-space: nowrap;

  &--ghost {
    background: transparent;
    color: var(--color-text-secondary);
    border-color: var(--gray-200);

    &:hover { background: var(--gray-50); border-color: var(--gray-300); color: var(--color-text-primary); }
    &.is-active { background: var(--color-primary-light-9); border-color: rgba(27, 58, 138, 0.25); color: var(--color-primary); }
  }

  &--primary {
    background: var(--color-primary);
    color: #fff;
    border-color: var(--color-primary);
    &:hover { background: var(--color-primary-hover); border-color: var(--color-primary-hover); }
    &:active { transform: scale(0.98); }
  }

  &--danger {
    background: var(--color-danger-bg);
    color: var(--color-danger);
    border-color: rgba(153, 27, 27, 0.2);
    &:hover { background: rgba(153, 27, 27, 0.12); }
  }
}

.filter-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 16px;
  height: 16px;
  background: var(--color-primary);
  color: #fff;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 700;
  padding: 0 4px;
  line-height: 1;
}

.filter-panel {
  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--gray-100);
    margin-bottom: 10px;
  }
  &__body {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-height: 380px;
    overflow-y: auto;
    &::-webkit-scrollbar { width: 4px; }
    &::-webkit-scrollbar-track { background: transparent; }
    &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
  }
}

.fp-title { font-size: 13px; font-weight: 700; color: var(--color-text-primary); letter-spacing: -0.1px; }
.fp-reset {
  border: none;
  background: transparent;
  color: var(--color-primary);
  font-size: 12px;
  cursor: pointer;
  font-family: inherit;
  padding: 2px 6px;
  border-radius: 4px;
  &:hover { background: var(--color-primary-light-9); }
}

.fp-group { display: flex; flex-direction: column; gap: 6px; }
.fp-label { font-size: 11px; font-weight: 600; color: var(--gray-500); text-transform: uppercase; letter-spacing: 0.5px; }

.fp-tags {
  display: flex;
  flex-direction: column;
  gap: 4px;
  &--wrap { flex-direction: row; flex-wrap: wrap; }
}

.fp-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 6px;
  border: 1.5px solid var(--gray-200);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.12s;
  font-family: inherit;
  white-space: nowrap;

  &:hover { border-color: var(--gray-400); color: var(--color-text-primary); }
  &.is-selected { background: var(--color-primary); border-color: var(--color-primary); color: #fff; }
}

.col-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  &.is-disabled { opacity: 0.5; pointer-events: none; }
  :deep(.el-checkbox) { font-size: 13px; color: var(--color-text-primary); }
}

/* 批量状态 dot */
.batch-status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-right: 4px;

  &.dot--pending { background: #f59e0b; }
  &.dot--processing { background: #3b82f6; }
  &.dot--completed { background: #22c55e; }
  &.dot--failed { background: #ef4444; }
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

.serial-text {
  display: inline-block;
  padding: 3px 10px;
  background: rgba(22, 163, 74, 0.12);
  color: #15803d;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.disk-text {
  display: inline-block;
  padding: 3px 10px;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.disk-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  position: relative;
  padding-right: 24px;
  min-width: 0;
}

.serial-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  position: relative;
  padding-right: 24px;
  min-width: 0;
}

.serial-copy {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
}

.disk-copy {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
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
.serial-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}
.disk-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.copy-btn--sm {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  padding: 0;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 3px;
  transition: color 0.15s, background-color 0.15s;
  flex-shrink: 0;
}
.copy-btn--sm:hover {
  color: #409eff;
  background-color: #f0f9ff;
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
}

.dynamic-field { font-size: 13px; color: var(--color-text-primary); font-weight: 500; }

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  background-clip: padding-box;

  &::after { display: none !important; }

  &--success { background: rgba(34,197,94,0.1); color: var(--color-success); }
  &--warning { background: rgba(245,158,11,0.1); color: var(--color-warning); }
  &--danger { background: rgba(239,68,68,0.08); color: var(--color-danger); }
  &--info { background: rgba(59,130,246,0.1); color: var(--chart-blue); }
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

.preview-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.preview-filter-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.preview-tag {
  font-size: 10px;
  font-weight: 600;
  color: #15803d;
  background: rgba(22, 163, 74, 0.1);
  border: 1px solid rgba(22, 163, 74, 0.25);
  border-radius: 4px;
  padding: 1px 5px;
  line-height: 1.4;
}

.preview-refresh-btn {
  align-self: flex-start;
  font-size: 11px;
  height: 24px;
  padding: 0 8px;
  border-radius: 5px;
  border-color: #86efac;
  color: #15803d;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 3px;

  &:hover {
    background: rgba(22, 163, 74, 0.06);
    border-color: #22c55e;
    color: #16a34a;
  }

  .el-icon {
    font-size: 12px;
    transition: transform 0.3s;
  }

  &:hover .el-icon {
    transform: rotate(180deg);
  }
}

.export-preview-card.is-loading .preview-num {
  opacity: 0.5;
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
  gap: 16px;
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

/* 横向步骤条 */
.stepper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 0 20px;
}
.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  position: relative;
  z-index: 1;
}
.step-dot {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #f3f4f6;
  border: 2px solid #d1d5db;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #9ca3af;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  font-family: 'Manrope', sans-serif;
}
.step-item--active .step-dot {
  background: #eff6ff;
  border-color: #005eeb;
  color: #005eeb;
  box-shadow: 0 0 0 4px rgba(0, 94, 235, 0.12);
  transform: scale(1.1);
}
.step-item--done .step-dot {
  background: #dcfce7;
  border-color: #16a34a;
  color: #16a34a;
}
.step-label {
  font-size: 11px;
  font-weight: 600;
  color: #9ca3af;
  text-align: center;
  white-space: nowrap;
  transition: color 0.2s;
}
.step-item--active .step-label { color: #005eeb; }
.step-item--done .step-label { color: #16a34a; }

.step-line {
  flex: 1;
  height: 2px;
  background: #e5e7eb;
  margin: 0 4px;
  margin-bottom: 22px;
  transition: background 0.4s;
  max-width: 80px;
}
.step-line--active {
  background: #16a34a;
}

/* 步骤内容区 */
.step-content {
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 10px;
  min-height: 160px;
  display: flex;
  flex-direction: column;
  animation: content-in 0.25s ease-out;
}
@keyframes content-in {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 步骤面板 */
.step-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 32px 20px;
  gap: 8px;
}
.step-panel__icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #005eeb;
  margin-bottom: 4px;
}
.step-panel__icon--confirm {
  background: #f0fdf4;
  border-color: #bbf7d0;
  color: #16a34a;
}
.step-panel__desc {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
  text-align: center;
}
.step-panel__file-name {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #005eeb;
  font-family: 'SF Mono', monospace;
  background: #eff6ff;
  padding: 4px 10px;
  border-radius: 5px;
  margin: 0;
}
.step-panel__btn {
  margin-top: 6px;
  display: flex;
  align-items: center;
}
.step-panel__next {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border: none;
  background: transparent;
  color: #005eeb;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s, color 0.15s;
  margin-top: 4px;
  &:hover { background: #eff6ff; }
  &:disabled { color: #d1d5db; cursor: not-allowed; &:hover { background: transparent; } }
}
.step-panel__back {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border: none;
  background: transparent;
  color: #9ca3af;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s, color 0.15s;
  &:hover { background: #f3f4f6; color: #6b7280; }
}
.step-panel__actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  margin-top: 4px;
}
.step-panel__right {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 拖拽上传区 */
.drop-zone {
  width: 100%;
  display: flex;
  justify-content: center;
}
.drop-zone :deep(.el-upload) {
  width: 100%;
}
.drop-zone :deep(.el-upload-drag) {
  width: 100%;
  height: 110px;
  border: 1.5px dashed #d1d5db;
  border-radius: 10px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s, background 0.2s;
  &:hover {
    border-color: #005eeb;
    background: #f0f7ff;
  }
}
.drop-zone__inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  pointer-events: none;
}
.drop-zone__icon {
  color: #9ca3af;
  transition: color 0.2s;
}
.drop-zone :deep(.el-upload-drag:hover) .drop-zone__icon {
  color: #005eeb;
}
.drop-zone__text {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}
.drop-zone__link {
  font-size: 13px;
  color: #005eeb;
  font-weight: 600;
  pointer-events: auto;
  &:hover { text-decoration: underline; }
}

/* 替换文件按钮 */
.drop-zone--replace {
  margin-top: 6px;
  display: flex;
  justify-content: center;
}
.drop-zone--replace :deep(.el-upload) {
  display: flex;
  justify-content: center;
}
.replace-btn {
  font-size: 12px;
  color: #9ca3af;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 5px;
  transition: color 0.15s, background 0.15s;
  &:hover { color: #005eeb; background: #eff6ff; }
}

/* 文件药片 */
.file-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: #fff;
  border: 1.5px solid #d1d5db;
  border-radius: 10px;
  width: 100%;
  animation: chip-in 0.2s ease-out;
}
@keyframes chip-in {
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
}
.file-chip__icon {
  color: #005eeb;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  background: #eff6ff;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.file-chip__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
}
.file-chip__name {
  font-size: 13px;
  font-weight: 600;
  color: #1f2937;
  font-family: 'SF Mono', monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.file-chip__size {
  font-size: 11px;
  color: #9ca3af;
}
.file-chip__remove {
  width: 26px;
  height: 26px;
  border: none;
  background: #f9fafb;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  flex-shrink: 0;
  transition: background 0.15s, color 0.15s;
  &:hover { background: #fef2f2; color: #dc2626; }
}

/* 预览药片 */
.preview-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}
.preview-pill--loading {
  background: #f9fafb;
  color: #9ca3af;
}
.preview-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #9ca3af;
  animation: dot-pulse 1.2s ease-in-out infinite;
}
@keyframes dot-pulse {
  0%, 100% { opacity: 0.4; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1); }
}
.preview-pill--ok {
  background: #f0fdf4;
  color: #16a34a;
}

/* 导入中 — 真实进度条 */
.importing-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 40px;
}
.progress-wrap {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
  max-width: 380px;
}
.progress-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #005eeb;
  flex-shrink: 0;
  animation: icon-pulse 1.5s ease-in-out infinite;
}
@keyframes icon-pulse {
  0%, 100% { opacity: 0.7; transform: scale(0.95); }
  50% { opacity: 1; transform: scale(1); }
}
.progress-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.progress-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.progress-label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}
.progress-pct {
  font-size: 13px;
  font-weight: 700;
  color: #005eeb;
  font-family: 'Manrope', monospace;
  min-width: 36px;
  text-align: right;
}
.progress-bar {
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}
.progress-bar__fill {
  height: 100%;
  background: #005eeb;
  border-radius: 3px;
  transition: width 0.3s ease-out;
  min-width: 4px;
}
.progress-sub {
  font-size: 11px;
  color: #9ca3af;
}

/* 导入按钮 */
.import-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

/* 导入结果 */
.import-result {
  padding: 0 4px 4px;
  animation: result-in 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes result-in {
  from { opacity: 0; transform: scale(0.96) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
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
  animation: fail-in 0.3s ease-out both;
}
.fail-item:nth-child(1) { animation-delay: 0.05s; }
.fail-item:nth-child(2) { animation-delay: 0.1s; }
.fail-item:nth-child(3) { animation-delay: 0.15s; }
.fail-item:nth-child(4) { animation-delay: 0.2s; }
.fail-item:nth-child(5) { animation-delay: 0.25s; }
@keyframes fail-in {
  from { opacity: 0; transform: translateX(-6px); }
  to { opacity: 1; transform: translateX(0); }
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

/* ==================== 抽屉头部通用 ==================== */
.drawer-head {
  display: flex;
  align-items: center;
  gap: 12px;
}
.drawer-head-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(0, 94, 235, 0.08);
  border: 1px solid rgba(0, 94, 235, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
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

/* ==================== 新增弹窗样式 ==================== */
.create-dialog {
  max-width: calc(100vw - 40px);

  :deep(.el-dialog__body) {
    padding: 0;
  }

  .dialog-head {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .dialog-head-icon {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    background: #eff6ff;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: #3b82f6;
  }

  .dialog-head-title {
    font-family: 'Manrope', sans-serif;
    font-size: 15px;
    font-weight: 600;
    color: #111827;
  }

  /* 表单区域 */
  .create-form {
    padding: 14px 20px 0;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  /* 提示卡片 */
  .form-hint {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: #f0fdf4;
    border-radius: 6px;
    font-size: 12px;
    color: #166534;

    svg {
      color: #22c55e;
      flex-shrink: 0;
    }

    em {
      font-style: normal;
      color: #dc2626;
      font-weight: 600;
    }
  }

  .form-hint-inline {
    font-size: 11px;
    color: #9ca3af;
    margin-top: 4px;
  }

  /* 表单分区 */
  .form-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .form-section__label {
    font-size: 11px;
    font-weight: 600;
    color: #6b7280;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 2px;
  }

  /* 表单行 */
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }

  .form-item {
    display: flex;
    flex-direction: column;
    gap: 4px;

    &--full {
      grid-column: 1 / -1;
    }

    &--date {
      margin-bottom: 4px;
    }
  }

  .form-item__label {
    font-size: 12px;
    font-weight: 500;
    color: #374151;

    em {
      font-style: normal;
      color: #dc2626;
    }
  }

  /* Element Plus 组件样式 */
  :deep(.el-input__wrapper),
  :deep(.el-select__wrapper) {
    border-radius: 6px;
  }

  :deep(.el-select) {
    width: 100%;
  }

  .form-item--date :deep(.el-date-editor) {
    width: 100% !important;

    .el-input__wrapper {
      width: 100%;
    }
  }

  /* 文件大小输入 */
  .size-input {
    display: flex;
    gap: 6px;

    :deep(.el-input-number) {
      flex: 1;
    }

    :deep(.el-select) {
      width: 80px;
      flex-shrink: 0;
    }
  }

  /* 底部按钮 */
  :deep(.el-dialog__footer) {
    padding: 12px 0 0;
    border: none;
  }

  :deep(.dialog-footer) {
    display: flex;
    justify-content: flex-end;
    gap: 10px;

    .el-button {
      border-radius: 6px;
      padding: 8px 16px;
      font-weight: 500;
    }
  }
}

/* ==================== 详情弹窗 ==================== */
.detail-dialog {
  max-width: calc(100vw - 40px);

  :deep(.el-dialog__header) {
    padding: 0 !important;
    margin: 0 !important;
    border-bottom: none !important;
  }

  :deep(.el-dialog__body) {
    padding: 0 !important;
  }

  :deep(.el-dialog__headerbtn) {
    top: 12px !important;
    right: 12px !important;
    width: 28px !important;
    height: 28px !important;
  }

  :deep(.el-dialog__headerbtn .el-icon) {
    font-size: 14px !important;
    color: #9ca3af;
  }

  :deep(.el-dialog__headerbtn:hover .el-icon) {
    color: #ef4444 !important;
    background: #fee2e2 !important;
    border-radius: 6px;
  }

  .detail-head {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 52px 10px 20px;

    &__icon {
      width: 28px;
      height: 28px;
      border-radius: 6px;
      background: #eff6ff;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #3b82f6;
      flex-shrink: 0;
    }

    &__title {
      font-size: 14px;
      font-weight: 600;
      color: #111827;
    }
  }

  .detail-body {
    display: flex;
    flex-direction: column;
  }

  .detail-row {
    display: flex;
    align-items: center;
    padding: 6px 20px;
    gap: 20px;
    border-bottom: 1px solid #f3f4f6;
    min-height: 32px;

    &:last-child {
      border-bottom: none;
    }

    &--path, &--serial {
      background: #fafafa;
    }

    &--remark {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
      padding: 8px 20px;
    }
  }

  .detail-row__label {
    width: 84px;
    flex-shrink: 0;
    font-size: 12px;
    font-weight: 500;
    color: #9ca3af;
    white-space: nowrap;
  }

  .detail-row__value {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
    gap: 6px;
    overflow: hidden;

    &--copy {
      display: flex;
      align-items: center;
      gap: 6px;
    }

    &--path {
      flex: 1;
      background: #f1f5f9;
      padding: 5px 10px;
      border-radius: 4px;
      border: 1px solid #e5e7eb;
      overflow: hidden;
      white-space: nowrap !important;
    }
  }

  .detail-row__text {
    font-size: 13px;
    color: #374151;
    white-space: nowrap !important;
    overflow: hidden;
    text-overflow: ellipsis;

    &--path {
      white-space: normal !important;
      word-break: break-all;
      overflow: visible;
      text-overflow: unset;
      font-family: 'SF Mono', Monaco, 'Courier New', monospace;
      font-size: 12px;
      color: var(--color-text-secondary);
      line-height: 1.5;
    }
  }

  .detail-row__code {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 12px;
    color: #1f2937;
    white-space: nowrap !important;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .detail-row__remark {
    width: 100%;
    background: #f8fafc;
    padding: 8px 10px;
    border-radius: 4px;
    border: 1px solid #e5e7eb;
    font-size: 13px;
    line-height: 1.5;
    color: #374151;
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;

    &::after { display: none !important; }

    &--pending { background: #fffbeb; color: #d97706; }
    &--processing { background: #eff6ff; color: #2563eb; }
    &--completed { background: #f0fdf4; color: #16a34a; }
    &--failed { background: #fef2f2; color: #dc2626; }
  }

  .copy-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    border: none;
    background: transparent;
    border-radius: 4px;
    cursor: pointer;
    color: #9ca3af;
    transition: all 0.15s;
    flex-shrink: 0;

    &:hover {
      background: #f3f4f6;
      color: #3b82f6;
    }
  }
}

/* ==================== 编辑弹窗 ==================== */
.edit-dialog {
  max-width: calc(100vw - 40px);

  :deep(.el-dialog__header) {
    padding: 0 !important;
    margin: 0 !important;
    border-bottom: none !important;
  }

  :deep(.el-dialog__body) {
    padding: 0 !important;
  }

  :deep(.el-dialog__footer) {
    padding: 12px 16px !important;
    border-top: 1px solid #f3f4f6;
  }

  :deep(.el-dialog__headerbtn) {
    top: 12px !important;
    right: 12px !important;
    width: 28px !important;
    height: 28px !important;
  }

  :deep(.el-dialog__headerbtn .el-icon) {
    font-size: 14px !important;
    color: #9ca3af;
  }

  :deep(.el-dialog__headerbtn:hover .el-icon) {
    color: #ef4444 !important;
    background: #fee2e2 !important;
    border-radius: 6px;
  }

  .edit-head {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 52px 10px 20px;

    &__icon {
      width: 28px;
      height: 28px;
      border-radius: 6px;
      background: #eff6ff;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #3b82f6;
      flex-shrink: 0;
    }

    &__title {
      font-size: 14px;
      font-weight: 600;
      color: #111827;
    }
  }

  .edit-body {
    padding: 12px 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .edit-row {
    display: flex;
    flex-direction: column;
    gap: 4px;

    &--full {
      flex-direction: column;
    }

    :deep(.el-select) {
      width: 100%;
    }
  }

  .edit-row__label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
  }

  .edit-size {
    display: flex;
    gap: 8px;

    :deep(.el-input-number) {
      flex: 1;
    }

    :deep(.el-select) {
      width: 80px;
      flex-shrink: 0;
    }
  }

  .edit-hint {
    font-size: 12px;
    color: #9ca3af;
  }

  .edit-divider {
    height: 1px;
    background: #f0f0f2;
  }
}

/* ==================== 通用复制按钮 ==================== */
.copy-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  border-radius: 5px;
  cursor: pointer;
  color: #9ca3af;
  transition: all 0.15s;

  &:hover {
    background: #f3f4f6;
    color: #3b82f6;
  }
}

/* ==================== 通用底部 ==================== */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;

  :deep(.el-button) {
    border-radius: 6px;
    padding: 8px 16px;
    font-weight: 500;
  }
}

/* ==================== 抽屉相关（旧样式保留） ==================== */
.drawer-head-icon--create {
  background: rgba(0, 94, 235, 0.1) !important;
  border-color: rgba(0, 94, 235, 0.15) !important;
  color: #005eeb !important;
}

/* 填写提示条 */
.create-tip {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 13px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  font-size: 12px;
  color: #1e40af;
  line-height: 1.4;

  em {
    font-style: normal;
    color: #dc2626;
    font-weight: 700;
  }

  .create-tip__icon {
    flex-shrink: 0;
    color: #3b82f6;
  }
}

/* 表单主体 */
.create-body {
  padding: 18px 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
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

/* ========== 现代字段布局 ========== */
.create-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.create-section__head {
  display: flex;
  align-items: center;
  gap: 6px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;

  &::before {
    content: '';
    display: inline-block;
    width: 3px;
    height: 13px;
    background: #005eeb;
    border-radius: 2px;
    flex-shrink: 0;
  }
}

.create-section__label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

.create-fields {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  align-items: start;
}

.create-field {
  display: flex;
  flex-direction: column;
  gap: 5px;

  &--wide {
    grid-column: 1 / -1;
  }
}

.create-field__label {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #4b5563;
  line-height: 1;

  em {
    font-style: normal;
    color: #dc2626;
    font-weight: 700;
    margin-left: 1px;
  }
}

/* 表单控件 */
.dr-input, .dr-select, .dr-textarea {
  width: 100%;
}

:deep(.dr-input .el-input__wrapper),
:deep(.dr-select .el-select__wrapper),
:deep(.el-date-editor.dr-input .el-input__wrapper) {
  background: #ffffff !important;
  border: 1px solid #d1d5db !important;
  border-radius: 7px !important;
  box-shadow: 0 0 0 1px #d1d5db inset !important;
  padding: 5px 11px !important;
  transition: border-color 0.2s, box-shadow 0.2s !important;
  font-family: 'DM Sans', sans-serif !important;
  font-size: 13px !important;
  height: 34px !important;
  box-sizing: border-box !important;
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
:deep(.dr-select .el-select__inner) {
  color: #1f2329 !important;
  font-size: 13px !important;
  font-weight: 500;
  &::placeholder { color: #c4c9d4 !important; }
}
:deep(.dr-input .el-input__prefix .el-icon),
:deep(.dr-select .el-select__prefix .el-icon) {
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
  width: 100%;
  box-sizing: border-box;
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
    height: 34px !important;
    box-sizing: border-box !important;
    &:hover { border-color: #9ca3af !important; box-shadow: 0 0 0 1px #9ca3af inset !important; }
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
    height: 34px !important;
    box-sizing: border-box !important;
    &:hover { border-color: #9ca3af !important; box-shadow: 0 0 0 1px #9ca3af inset !important; }
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
  .filter-card { padding: 8px 12px; }
}

/* ==================== 弹窗通用样式 ==================== */
:deep(.detail-modal),
:deep(.edit-modal) {
  .el-dialog__header {
    padding: 16px 20px;
    margin-right: 0;
    border-bottom: 1px solid var(--color-border-light);
  }
  .el-dialog__body {
    padding: 0;
  }
  .el-dialog__footer {
    padding: 12px 20px;
    border-top: 1px solid var(--color-border-light);
  }
}

.modal-head {
  display: flex;
  align-items: center;
  gap: 10px;
}
.modal-head__icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  &--info { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
  &--edit { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
}
.modal-head__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.modal-head__title {
  font-family: 'Manrope', 'Inter', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
  -webkit-font-smoothing: antialiased;
}
.modal-head__sub {
  display: flex;
  align-items: center;
  gap: 6px;
}
.modal-serial {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  color: var(--color-text-secondary);
  background: var(--color-page-bg);
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--color-border-light);
}
.modal-copy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  padding: 0;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 5px;
  transition: all 0.15s;
  &:hover {
    background: rgba(59, 130, 246, 0.08);
    border-color: rgba(59, 130, 246, 0.3);
    color: #3b82f6;
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* ==================== 详情弹窗 ==================== */
.detail-body {
  padding: 16px 20px;
}
.detail-grid {
  display: flex;
  flex-direction: column;
  gap: 0;
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  overflow: hidden;
}
.detail-row {
  display: flex;
  align-items: center;
  min-height: 40px;
  border-bottom: 1px solid var(--color-border-light);
  &:last-child { border-bottom: none; }
  &--path,
  &--full {
    align-items: flex-start;
  }
}
.detail-row__label {
  width: 100px;
  flex-shrink: 0;
  padding: 0 14px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  background: var(--color-page-bg);
  height: 100%;
  display: flex;
  align-items: center;
  border-right: 1px solid var(--color-border-light);
  min-height: 40px;
}
.detail-row__value {
  flex: 1;
  padding: 8px 14px;
  font-size: 13px;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  &--path {
    flex-wrap: wrap;
    gap: 8px;
    align-items: flex-start;
  }
  &--remark {
    color: var(--color-text-secondary);
    font-size: 12px;
    line-height: 1.6;
  }
  &--copy {
    gap: 6px;
    display: flex;
    align-items: center;
  }
}
.detail-status {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  &--completed { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
  &--pending { background: rgba(245, 158, 11, 0.1); color: #d97706; }
  &--processing { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
  &--failed { background: rgba(239, 68, 68, 0.1); color: #dc2626; }
}
.detail-path {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
  color: var(--color-text-secondary);
  word-break: break-all;
  line-height: 1.6;
  background: var(--color-page-bg);
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--color-border-light);
}
.detail-copy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  padding: 0;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.15s;
  flex-shrink: 0;
  &:hover {
    background: rgba(59, 130, 246, 0.08);
    border-color: rgba(59, 130, 246, 0.3);
    color: #3b82f6;
  }
}

/* ==================== 编辑弹窗 ==================== */
.edit-tip {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 14px;
  margin: 0 20px 4px;
  background: #f0f9ff;
  border: 1px solid #bae0ff;
  border-radius: 6px;
  font-size: 12px;
  color: #1d7ec7;
  line-height: 1.5;
}
.edit-tip__icon {
  flex-shrink: 0;
  margin-top: 1px;
  color: #409eff;
}
.edit-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.edit-divider {
  height: 1px;
  background: var(--color-border-light);
  margin: 2px 0;
}
.edit-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.edit-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.edit-select {
  width: 100%;
}

.edit-size-row {
  display: flex;
  gap: 6px;
  align-items: center;
}

.edit-size-num {
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

.edit-size-unit {
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

.edit-size-hint {
  font-size: 11px;
  color: #9ca3af;
  font-family: 'SF Mono', 'Fira Code', monospace;
  margin-top: 3px;
}

.edit-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
