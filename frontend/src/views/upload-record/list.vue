<template>
  <div class="page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">上传记录</h1>
        <p class="page-subtitle">管理所有数据上传记录</p>
      </div>
      <div class="header-actions">
        <el-button type="success" :loading="exporting" @click="showExportDialog">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          导出 Excel
        </el-button>
        <el-button type="warning" @click="showImportDialog">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
          批量导入
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          新增上传记录
        </el-button>
        <el-button type="primary" @click="loadRecords">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 6px">
            <path d="M23 4v6h-6M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          刷新
        </el-button>
      </div>
    </header>

    <!-- 导出弹窗 -->
    <el-drawer
      v-model="exportDialogVisible"
      direction="rtl"
      size="440px"
      :with-header="true"
      class="export-drawer"
    >
      <template #header>
        <div class="drawer-title-inner">
          <div class="drawer-title-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
          </div>
          <span class="drawer-title-text">导出上传记录</span>
        </div>
      </template>

      <div class="export-form">
        <!-- 摘要预览 -->
        <div class="export-summary-card" v-if="exportPreviewCount > 0">
          <div class="summary-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/>
            </svg>
          </div>
          <div class="summary-info">
            <span class="summary-num">{{ exportPreviewCount }}</span>
            <span class="summary-label">条记录将导出</span>
          </div>
          <div class="summary-tip">按当前筛选条件</div>
        </div>

        <!-- 筛选条件 -->
        <div class="export-filter-section">
          <div class="filter-section-title">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
            筛选条件
          </div>

          <div class="filter-field">
            <label class="filter-label">磁盘标签</label>
            <el-input v-model="exportForm.dataType" placeholder="全部标签" clearable size="small" />
          </div>

          <div class="filter-field">
            <label class="filter-label">状态</label>
            <el-select v-model="exportForm.status" placeholder="全部状态" clearable size="small" style="width: 100%">
              <el-option label="待处理" value="pending" />
              <el-option label="处理中" value="processing" />
              <el-option label="已完成" value="completed" />
              <el-option label="失败" value="failed" />
            </el-select>
          </div>

          <div class="filter-field">
            <label class="filter-label">上传人</label>
            <el-input v-model="exportForm.uploader" placeholder="全部人员" clearable size="small" />
          </div>

          <div class="filter-field">
            <label class="filter-label">日期范围</label>
            <el-date-picker
              v-model="exportForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始"
              end-placeholder="结束"
              value-format="YYYY-MM-DD"
              size="small"
              style="width: 100%"
            />
          </div>

          <div class="filter-field">
            <label class="filter-label">关键词</label>
            <el-input v-model="exportForm.keyword" placeholder="路径/备注" clearable size="small" />
          </div>
        </div>
      </div>

      <div class="export-drawer-foot">
        <el-button size="small" @click="exportDialogVisible = false">取消</el-button>
        <el-button type="primary" size="small" :loading="exporting" @click="handleExport">
          <svg v-if="!exporting" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          导出 Excel
        </el-button>
      </div>
    </el-drawer>

    <!-- 批量导入弹窗 -->
    <el-dialog v-model="importDialogVisible" width="680px" destroy-on-close>
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag">导入</span>
          <span class="dialog-title-text">批量导入上传记录</span>
        </div>
      </template>
      <div class="import-dialog-body">
        <!-- 操作区 -->
        <div class="import-cards">
          <!-- 步骤1: 下载模板 -->
          <div class="import-card import-card--blue">
            <div class="import-card__badge">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
            </div>
            <div class="import-card__body">
              <div class="import-card__title">下载模板</div>
              <div class="import-card__desc">获取空白 Excel 模板</div>
            </div>
            <el-button class="import-card__action" type="primary" plain size="small" :loading="downloadingTemplate" @click="handleDownloadTemplate">
              下载 xlsx
            </el-button>
          </div>

          <!-- 步骤2: 上传文件 -->
          <div class="import-card import-card--teal">
            <div class="import-card__badge">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <div class="import-card__body">
              <div class="import-card__title">上传文件</div>
              <div class="import-card__desc">选择已填写的 xlsx 文件</div>
            </div>
            <el-upload
              ref="uploadRef"
              class="import-card__action"
              action="#"
              :auto-upload="false"
              :limit="1"
              accept=".xlsx"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              :file-list="fileList"
            >
              <el-button size="small" plain>选择文件</el-button>
            </el-upload>
          </div>

          <!-- 步骤3: 开始导入 -->
          <div class="import-card import-card--green">
            <div class="import-card__badge">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 11 12 14 22 4"/>
                <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
              </svg>
            </div>
            <div class="import-card__body">
              <div class="import-card__title">确认导入</div>
              <div class="import-card__desc">数据将批量写入系统</div>
            </div>
            <el-button class="import-card__action" type="success" plain size="small" :loading="importing" :disabled="!selectedFile" @click="handleImport">
              开始导入
            </el-button>
          </div>
        </div>

        <!-- 字段说明 -->
        <div class="import-field-guide" v-if="importTemplateFields.length > 0">
          <div class="guide-title">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
            字段填写说明
          </div>
          <div class="guide-table">
            <div class="guide-table-header">
              <span>字段名称</span>
              <span>是否必填</span>
              <span>填写示例</span>
              <span>说明</span>
            </div>
            <div class="guide-table-row" v-for="f in importTemplateFields" :key="f.code">
              <span>{{ f.field }}</span>
              <span :class="f.required ? 'required' : 'optional'">{{ f.required ? '必填 *' : '选填' }}</span>
              <span class="example">{{ f.example || '-' }}</span>
              <span class="field-type">{{ f.type === 'select' ? `下拉选项: ${f.options}` : f.type === 'number' ? '整数（字节数）' : '文本' }}</span>
            </div>
          </div>
        </div>

        <!-- 导入结果 -->
        <div class="import-result" v-if="importResult">
          <div class="result-summary" :class="importResult.failed > 0 ? 'has-error' : 'all-success'">
            <div class="result-icon">
              <svg v-if="importResult.failed === 0" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <div class="result-text">
              <div class="result-title">导入完成</div>
              <div class="result-stats">
                共 {{ importResult.total }} 行，成功 <strong style="color:#16a34a">{{ importResult.success }}</strong> 行，
                <span v-if="importResult.failed > 0">失败 <strong style="color:#dc2626">{{ importResult.failed }}</strong> 行</span>
                <span v-else>失败 0 行</span>
              </div>
            </div>
          </div>
          <div class="fail-rows" v-if="importResult.failRows.length > 0">
            <div class="fail-rows-title">失败行明细：</div>
            <div class="fail-rows-list">
              <div class="fail-row-item" v-for="(f, idx) in importResult.failRows" :key="idx">
                <span class="fail-row-num">第{{ f.row }}行</span>
                <span class="fail-row-reason">{{ f.reason }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="closeImportDialog">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <!-- 磁盘标签 -->
      <el-input v-model="searchDataType" placeholder="磁盘标签" clearable size="default" />

      <!-- 项目名称 -->
      <el-select v-model="searchProjectName" placeholder="项目名称" clearable size="default" filterable>
        <el-option v-for="p in projectList" :key="p.id" :label="p.name" :value="p.name" />
      </el-select>

      <!-- 状态 -->
      <el-select v-model="searchStatus" placeholder="状态" clearable size="default">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="失败" value="failed" />
      </el-select>

      <!-- 上传人 -->
      <el-input v-model="searchUploader" placeholder="上传人" clearable size="default" />

      <!-- 日期范围 -->
      <el-date-picker
        v-model="searchDateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始"
        end-placeholder="结束"
        value-format="YYYY-MM-DD"
        size="default"
      />

      <!-- 关键词 -->
      <el-input v-model="searchKeyword" placeholder="搜索路径/备注" clearable size="default">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="M21 21l-4.35-4.35"/>
          </svg>
        </template>
      </el-input>

      <!-- 操作按钮 -->
      <div class="filter-actions">
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <!-- 表格工具栏 -->
      <div class="table-toolbar">
        <div class="toolbar-left">
          <span class="record-count">共 <strong>{{ pagination.total }}</strong> 条记录</span>
          <span v-if="selectedRows.length > 0" class="selection-count">
            已选择 <strong>{{ selectedRows.length }}</strong> 项
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
            批量删除
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
                字段显示
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-left: 4px">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </el-button>
            </template>
            <div class="column-settings">
              <div class="settings-header">
                <span class="settings-title">字段显示配置</span>
                <el-button type="primary" text size="small" @click="handleResetColumns">重置</el-button>
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
                <span class="settings-hint">提示：拖拽表头可调整列顺序</span>
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
        <el-table-column v-if="isColumnVisible('serialNo')" prop="serialNo" label="流水号" min-width="110" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('dataType')" prop="dataType" label="磁盘标签" min-width="90" align="center">
          <template #default="{ row }">
            <span class="type-tag">{{ row.dataType }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('projectName')" prop="projectName" label="项目名称" min-width="110" show-overflow-tooltip />
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
        <el-table-column v-if="isColumnVisible('destPath')" label="目标路径" min-width="160" align="center">
          <template #default="{ row }">
            <div v-if="row.destPath" class="path-cell">
              <el-tooltip :content="row.destPath" placement="top" :show-after="300">
                <span class="path-text">{{ row.destPath }}</span>
              </el-tooltip>
              <el-tooltip content="复制路径" placement="top">
                <button class="copy-btn" @click.stop="copyPath(row.destPath)" title="复制路径">
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
        <el-table-column v-if="isColumnVisible('fileSize')" prop="fileSizeStr" label="文件大小" min-width="90" align="center" />
        <el-table-column v-if="isColumnVisible('uploader')" prop="uploader" label="上传人" min-width="80" align="center" />
        <el-table-column v-if="isColumnVisible('status')" prop="status" label="状态" min-width="70" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="getStatusClass(row.status)">
              {{ row.statusText }}
            </span>
          </template>
        </el-table-column>
        <el-table-column v-if="isColumnVisible('remark')" prop="remark" label="备注" min-width="100" show-overflow-tooltip />
        <el-table-column v-if="isColumnVisible('createdAt')" prop="createdAt" label="时间" min-width="130" />
        <el-table-column label="操作" width="110" fixed="right" align="center">
          <template #default="{ row }">
            <TableActions :actions="[
              { key: 'detail', label: '详情', type: 'primary' },
              { key: 'edit', label: '编辑', type: 'primary' },
              { key: 'delete', label: '删除', type: 'danger' }
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
        />
      </div>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" width="640px">
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag dialog-mode-tag--info">详情</span>
          <span class="dialog-title-text">记录详情</span>
        </div>
      </template>
      <div v-if="currentRecord" class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="流水号">
            <code class="serial-no">{{ currentRecord.serialNo }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="磁盘标签">{{ currentRecord.dataType }}</el-descriptions-item>
          <el-descriptions-item label="项目名称">{{ currentRecord.projectName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="上传人">{{ currentRecord.uploader }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <span class="status-chip" :class="getStatusClass(currentRecord.status)">
              {{ currentRecord.statusText }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="文件大小">{{ currentRecord.fileSizeStr }}</el-descriptions-item>
          <el-descriptions-item label="目标路径" :span="2">
            <div class="detail-path-cell">
              <el-tooltip :content="currentRecord.destPath" placement="top" :show-after="200" v-if="currentRecord.destPath">
                <span class="detail-path-text">{{ currentRecord.destPath }}</span>
              </el-tooltip>
              <span v-else class="no-data">-</span>
              <el-tooltip content="复制路径" placement="top" v-if="currentRecord.destPath">
                <button class="copy-btn copy-btn-lg" @click="copyPath(currentRecord.destPath)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ currentRecord.remark || '无' }}</el-descriptions-item>
          <el-descriptions-item label="上传时间">{{ currentRecord.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ currentRecord.updatedAt }}</el-descriptions-item>

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
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" width="520px">
      <template #header>
        <div class="dialog-head">
          <span class="dialog-mode-tag">编辑</span>
          <span class="dialog-title-text">编辑记录</span>
        </div>
      </template>
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>

        <!-- 动态字段编辑 -->
        <el-form-item
          v-for="col in dynamicColumns"
          :key="col.code"
          :label="col.name"
        >
          <el-select
            v-if="col.type === 'select'"
            v-model="editForm.data[col.code]"
            style="width: 100%"
            clearable
            :placeholder="col.placeholder"
          >
            <el-option v-for="opt in col.options" :key="opt" :label="opt" :value="opt" />
          </el-select>
          <el-date-picker
            v-else-if="col.type === 'date'"
            v-model="editForm.data[col.code]"
            type="date"
            value-format="YYYY-MM-DD"
            style="width: 100%"
            :placeholder="col.placeholder"
          />
          <el-date-picker
            v-else-if="col.type === 'datetime'"
            v-model="editForm.data[col.code]"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
            :placeholder="col.placeholder"
          />
          <el-input-number
            v-else-if="col.type === 'number'"
            v-model="editForm.data[col.code]"
            style="width: 100%"
            :placeholder="col.placeholder"
          />
          <el-input
            v-else
            v-model="editForm.data[col.code]"
            :placeholder="col.placeholder"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="editForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmEdit">保存</el-button>
      </template>
    </el-dialog>

    <!-- 新增上传记录抽屉 -->
    <el-drawer
      v-model="createVisible"
      :title="null"
      direction="rtl"
      size="520px"
      :before-close="() => createVisible = false"
      class="create-drawer"
    >
      <!-- 抽屉头部 -->
      <template #header>
        <div class="drawer-title">
          <div class="drawer-title-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="12" y1="18" x2="12" y2="12"/>
              <line x1="9" y1="15" x2="15" y2="15"/>
            </svg>
          </div>
          <span class="drawer-title-text">新增上传记录</span>
        </div>
      </template>

      <div class="create-drawer-content">
        <!-- 基本信息 -->
        <div class="form-section">
          <div class="form-section-label">
            <span class="label-dot label-dot--required"></span>
            基本信息
          </div>
          <el-form ref="createFormRef" :model="createForm" label-position="top">
            <el-form-item label="磁盘标签" prop="dataType" class="is-required-field">
              <el-input v-model="createForm.dataType" placeholder="输入磁盘标签，如：DataDisk-01" clearable>
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <ellipse cx="12" cy="5" rx="9" ry="3"/>
                    <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
                    <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
                  </svg>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="项目名称" prop="projectName" class="is-required-field">
              <el-input v-model="createForm.projectName" placeholder="输入项目名称" clearable>
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 3h7v7H3zM14 3h7v7h-7zM14 14h7v7h-7zM3 14h7v7H3z"/>
                  </svg>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="上传人" prop="uploader" class="is-required-field">
              <el-input v-model="createForm.uploader" placeholder="输入上传人姓名" clearable>
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </template>
              </el-input>
            </el-form-item>
          </el-form>
        </div>

        <!-- 文件信息 -->
        <div class="form-section">
          <div class="form-section-label">
            <span class="label-dot label-dot--required"></span>
            文件信息
          </div>
          <el-form :model="createForm" label-position="top">
            <el-form-item label="目标路径" prop="destPath" class="is-required-field">
              <el-input v-model="createForm.destPath" placeholder="/data/uploads/file.csv" clearable>
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
                    <polyline points="13 2 13 9 20 9"/>
                  </svg>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="文件大小" class="is-required-field">
              <div class="size-input-row">
                <el-input-number
                  v-model="fileSizeInputVal"
                  :min="0"
                  :precision="3"
                  placeholder="输入数值"
                  controls-position="right"
                  class="size-num-input"
                  @change="syncFileSizeFromInput"
                />
                <el-select v-model="fileSizeUnit" class="size-unit-select" @change="syncFileSizeFromInput">
                  <el-option label="B" value="B" />
                  <el-option label="KB" value="KB" />
                  <el-option label="MB" value="MB" />
                  <el-option label="GB" value="GB" />
                  <el-option label="TB" value="TB" />
                </el-select>
              </div>
              <div class="size-hint" v-if="fileSizeInputVal > 0">
                <span class="size-converted">≈ {{ formatSizeInOtherUnits(fileSizeInputVal, fileSizeUnit) }}</span>
              </div>
            </el-form-item>

            <el-form-item label="状态" prop="status" class="is-required-field">
              <el-select v-model="createForm.status" style="width: 100%">
                <el-option label="待处理" value="pending" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="失败" value="failed" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>

        <!-- 扩展字段 -->
        <div class="form-section" v-if="dynamicColumns.length > 0">
          <div class="form-section-label">
            <span class="label-dot label-dot--optional"></span>
            扩展字段
            <span class="section-count">{{ dynamicColumns.length }} 项</span>
          </div>
          <el-form :model="createForm" label-position="top">
            <el-form-item
              v-for="col in dynamicColumns"
              :key="'create-' + col.code"
              :label="col.name"
            >
              <el-input v-if="isIpField(col)" v-model="createForm.data[col.code]" placeholder="如：192.168.1.1" clearable>
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="2" y1="12" x2="22" y2="12"/>
                    <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
                  </svg>
                </template>
              </el-input>
              <el-input-number v-else-if="col.type === 'number'" v-model="createForm.data[col.code]" style="width: 100%" controls-position="right" :placeholder="col.placeholder" />
              <el-date-picker v-else-if="col.type === 'date'" v-model="createForm.data[col.code]" type="date" value-format="YYYY-MM-DD" style="width: 100%" :placeholder="col.placeholder" />
              <el-date-picker v-else-if="col.type === 'datetime'" v-model="createForm.data[col.code]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" :placeholder="col.placeholder" />
              <el-select v-else-if="col.type === 'select'" v-model="createForm.data[col.code]" style="width: 100%" clearable :placeholder="col.placeholder">
                <el-option v-for="opt in col.options" :key="opt" :label="opt" :value="opt" />
              </el-select>
              <el-input v-else v-model="createForm.data[col.code]" :placeholder="col.placeholder" clearable />
            </el-form-item>
          </el-form>
        </div>

        <!-- 备注信息 -->
        <div class="form-section">
          <div class="form-section-label">
            <span class="label-dot label-dot--optional"></span>
            备注信息
          </div>
          <el-form :model="createForm" label-position="top">
            <el-form-item label="备注">
              <el-input v-model="createForm.remark" type="textarea" :rows="3" placeholder="可选填写备注信息，如数据来源、特殊说明等" show-word-limit maxlength="500" />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="createVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmCreate">确认创建</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, inject } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadRecordApi, type UploadRecord, type ImportTemplateField, type ImportResultResp } from '@/api/upload-record'
import { FieldConfigApi, type FieldConfig } from '@/api/field-config'
import { ProjectApi, type ProjectSimple } from '@/api/project'
import TableActions from '@/components/TableActions.vue'

const trackExport = inject<(action?: string) => void>('trackExport')

// ==================== 字段配置 ====================
interface ColumnConfig {
  prop: string
  label: string
  visible: boolean
  required?: boolean
}

const allColumns = ref<ColumnConfig[]>([
  { prop: 'serialNo', label: '流水号', visible: true, required: true },
  { prop: 'dataType', label: '磁盘标签', visible: true, required: true },
  { prop: 'projectName', label: '项目名称', visible: true, required: false },
  { prop: 'destPath', label: '目标路径', visible: false, required: false },
  { prop: 'fileSize', label: '文件大小', visible: false, required: false },
  { prop: 'uploader', label: '上传人', visible: true, required: false },
  { prop: 'status', label: '状态', visible: true, required: true },
  { prop: 'remark', label: '备注', visible: false, required: false },
  { prop: 'createdAt', label: '时间', visible: true, required: false },
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
  ElMessage.success('已重置为默认配置')
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
const exporting = ref(false)
const selectedRows = ref<UploadRecord[]>([])
const exportDialogVisible = ref(false)
const exportPreviewCount = ref(0)

// ==================== 批量导入状态 ====================
const importDialogVisible = ref(false)
const importing = ref(false)
const downloadingTemplate = ref(false)
const selectedFile = ref<File | null>(null)
const fileList = ref<any[]>([])
const uploadRef = ref()
const importTemplateFields = ref<ImportTemplateField[]>([])
const importResult = ref<ImportResultResp | null>(null)
const searchDataType = ref('')
const searchStatus = ref('')
const searchUploader = ref('')
const searchKeyword = ref('')
const searchDateRange = ref<string[]>([])
const tableData = ref<UploadRecord[]>([])
const dataTypes = ref<string[]>([])
const dynamicColumns = ref<FieldConfig[]>([])
const projectList = ref<ProjectSimple[]>([])
const searchProjectName = ref('')
const detailVisible = ref(false)
const editVisible = ref(false)
const createVisible = ref(false)
const createFormRef = ref()
const currentRecord = ref<UploadRecord | null>(null)
const tableRef = ref()

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

const editForm = reactive<{ id: number; status: string; remark: string; data: Record<string, any> }>({
  id: 0,
  status: '',
  remark: '',
  data: {}
})

const exportForm = reactive({
  dataType: '',
  status: '',
  uploader: '',
  dateRange: [] as string[],
  keyword: ''
})

const createForm = reactive({
  dataType: '',
  projectName: '',
  destPath: '',
  fileSize: 0,
  uploader: '',
  status: 'pending',
  remark: '',
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
    ElMessage.success('路径已复制')
  } catch {
    ElMessage.error('复制失败')
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
    ElMessage.error('加载数据失败')
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
  editForm.status = row.status
  editForm.remark = row.remark
  editForm.data = { ...row.data } || {}
  editVisible.value = true
}

const handleDelete = async (row: UploadRecord) => {
  try {
    await ElMessageBox.confirm(`确定要删除流水号"${row.serialNo}"的记录吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await UploadRecordApi.del(row.id)
    ElMessage.success('删除成功')
    loadRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
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
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 条记录吗？`, '批量删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const ids = selectedRows.value.map(row => row.id)
    await UploadRecordApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 条记录`)
    selectedRows.value = []
    loadRecords()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
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
  createForm.data = {}
  fileSizeInputVal.value = 0
  fileSizeUnit.value = 'MB'
  createForm.fileSize = 0
  createVisible.value = true
}

const confirmCreate = async () => {
  if (!createForm.dataType) {
    ElMessage.error('请选择或输入磁盘标签')
    return
  }
  if (!createForm.destPath) {
    ElMessage.error('请输入目标路径')
    return
  }
  if (!createForm.fileSize || fileSizeInputVal.value <= 0) {
    ElMessage.error('请输入文件大小')
    return
  }
  if (!createForm.projectName) {
    ElMessage.error('请输入项目名称')
    return
  }
  if (!createForm.uploader) {
    ElMessage.error('请输入上传人')
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
      data: Object.keys(createForm.data).length > 0 ? createForm.data : undefined
    })
    if (res.code === 200) {
      ElMessage.success('创建成功')
      createVisible.value = false
      loadRecords()
    } else {
      ElMessage.error(res.message || '创建失败')
    }
  } catch (error) {
    ElMessage.error('创建失败')
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
    ElMessage.success('更新成功')
    editVisible.value = false
    loadRecords()
  } catch (error) {
    ElMessage.error('更新失败')
  } finally {
    submitting.value = false
  }
}

const showExportDialog = () => {
  exportForm.dataType = searchDataType.value
  exportForm.status = searchStatus.value
  exportForm.uploader = searchUploader.value
  exportForm.dateRange = searchDateRange.value || []
  exportForm.keyword = searchKeyword.value
  exportPreviewCount.value = pagination.total
  exportDialogVisible.value = true
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params: any = {}
    if (exportForm.dataType) params.dataType = exportForm.dataType
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
    // Track export milestone
    trackExport?.('upload-record')
  } catch (error) {
    console.error('Export failed:', error)
  } finally {
    exporting.value = false
  }
}

// ==================== 批量导入 ====================
const showImportDialog = async () => {
  importDialogVisible.value = true
  importResult.value = null
  selectedFile.value = null
  fileList.value = []
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
    ElMessage.error('模板下载失败')
  } finally {
    downloadingTemplate.value = false
  }
}

const handleFileChange = (file: any) => {
  selectedFile.value = file.raw as File
}

const handleFileRemove = () => {
  selectedFile.value = null
}

const handleImport = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择要导入的 Excel 文件')
    return
  }
  importing.value = true
  importResult.value = null
  try {
    const res = await UploadRecordApi.importRecords(selectedFile.value)
    if (res.code === 200) {
      importResult.value = res.data
      if (res.data.failed === 0) {
        ElMessage.success(`导入成功！共导入 ${res.data.success} 条记录`)
      } else {
        ElMessage.warning(`导入完成：成功 ${res.data.success} 行，失败 ${res.data.failed} 行，请查看下方失败明细`)
      }
    } else {
      ElMessage.error(res.message || '导入失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '导入失败，请检查文件格式是否正确')
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

/* 导出表单 */
.export-form {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--color-surface-2);
  display: flex;
  flex-direction: column;
  gap: 12px;

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
}

/* 摘要卡片 */
.export-summary-card {
  background: linear-gradient(135deg, rgba(34,197,94,0.08) 0%, rgba(22,163,74,0.04) 100%);
  border: 1px solid rgba(34,197,94,0.15);
  border-radius: var(--radius-md);
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 3px;
    height: 100%;
    background: var(--color-success);
    border-radius: 3px 0 0 3px;
  }
}

.summary-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: rgba(34,197,94,0.12);
  color: var(--color-success);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.summary-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.summary-num {
  font-family: 'Manrope', monospace;
  font-size: 22px;
  font-weight: 900;
  color: var(--color-success);
  line-height: 1;
}

.summary-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.summary-tip {
  margin-left: auto;
  font-size: 11px;
  color: var(--color-text-muted);
  align-self: flex-start;
}

/* 筛选区块 */
.export-filter-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.filter-section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--color-border-light);

  svg { color: var(--color-primary); }
}

.filter-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filter-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

/* 底部 */
.export-drawer-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

/* ==================== 批量导入弹窗 ==================== */
.import-dialog-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.import-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.import-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px 14px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  border-top-width: 3px;
  position: relative;
}

.import-card--blue { border-top-color: #005eeb; }
.import-card--teal { border-top-color: #0d9488; }
.import-card--green { border-top-color: #16a34a; }

.import-card__badge {
  width: 34px;
  height: 34px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 2px;
}
.import-card--blue .import-card__badge { background: #eff6ff; color: #005eeb; }
.import-card--teal .import-card__badge { background: #f0fdfa; color: #0d9488; }
.import-card--green .import-card__badge { background: #f0fdf4; color: #16a34a; }

.import-card__body { flex: 1; }
.import-card__title { font-size: 13px; font-weight: 700; color: var(--color-text-primary); margin-bottom: 2px; }
.import-card__desc { font-size: 11px; color: var(--color-text-muted); line-height: 1.4; }
.import-card__action { width: 100%; margin-top: 4px; }

.import-field-guide {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
  margin-top: var(--space-3);
}

.guide-title {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: var(--color-surface-3);
  border-bottom: 1px solid var(--color-border-light);
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  letter-spacing: 0.2px;
}

.guide-table {
  width: 100%;
}

.guide-table-header {
  display: grid;
  grid-template-columns: 100px 64px 120px 1fr;
  gap: 8px;
  padding: 8px 14px;
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  background: var(--color-surface-2);
  border-bottom: 1px solid var(--color-border-light);
}

.guide-table-row {
  display: grid;
  grid-template-columns: 100px 64px 120px 1fr;
  gap: 8px;
  padding: 8px 14px;
  align-items: center;
  border-bottom: 1px solid var(--color-border-light);
  font-size: 12px;
  color: var(--color-text-primary);
  &:last-child { border-bottom: none; }
  &:hover { background: rgba(0,94,235,0.02); }
}

.required { color: var(--color-danger); font-weight: 700; }
.optional { color: var(--color-text-muted); font-weight: 500; }
.example { font-family: 'SF Mono', 'DM Sans', monospace; font-size: 11px; color: var(--color-primary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.field-type { color: var(--color-text-muted); font-size: 11px; }

.import-result {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
  margin-top: var(--space-3);
}

.result-summary {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border-light);
  &.all-success { background: rgba(34,197,94,0.04); }
  &.has-error { background: rgba(217,119,6,0.04); }
}

.result-icon { flex-shrink: 0; display: flex; align-items: center; }
.result-text { flex: 1; }
.result-title { font-size: 14px; font-weight: 700; color: var(--color-text-primary); margin-bottom: 2px; font-family: 'Manrope', sans-serif; }
.result-stats { font-size: 12px; color: var(--color-text-secondary); }

.fail-rows { padding: 12px 14px; }
.fail-rows-title { font-size: 11px; font-weight: 700; color: var(--color-text-secondary); text-transform: uppercase; letter-spacing: 0.3px; margin-bottom: 8px; }
.fail-rows-list { display: flex; flex-direction: column; gap: 4px; }
.fail-row-item { display: flex; align-items: center; gap: 8px; padding: 5px 8px; background: var(--color-surface-2); border-radius: var(--radius-xs); font-size: 12px; }
.fail-row-num { color: var(--color-danger); font-weight: 700; font-family: 'SF Mono', monospace; min-width: 40px; }
.fail-row-reason { color: var(--color-text-secondary); }

/* ==================== 新增抽屉 ==================== */
.drawer-title {
  display: flex;
  align-items: center;
  gap: 10px;
}
.drawer-title-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.drawer-title-text {
  font-family: 'Manrope', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.create-drawer-content {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.form-section-label {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  background: var(--color-surface-3);
  border-bottom: 1px solid var(--color-border-light);
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.4px;
}

.label-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
  &--required { background: var(--color-danger); }
  &--optional { background: var(--color-text-muted); }
}

.section-count {
  margin-left: auto;
  font-size: 10px;
  font-weight: 600;
  color: var(--color-primary);
  background: var(--color-primary-light-9);
  padding: 1px 6px;
  border-radius: var(--radius-full);
  text-transform: none;
  letter-spacing: 0;
}

.form-section :deep(.el-form-item) {
  margin-bottom: 0;
  border-bottom: 1px solid var(--color-border-light);
  &:last-child { border-bottom: none; }
  .el-form-item__label {
    font-weight: 600;
    color: var(--color-text-primary);
    font-size: 13px;
    padding: 10px 14px 0;
    &::before { color: var(--color-danger); }
  }
  .el-form-item__content {
    padding: 4px 14px 10px;
  }
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
</style>
