<template>
  <div class="project-page">
    <!-- 页面标题栏 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">项目管理</h1>
        <span class="page-subtitle">项目信息库</span>
      </div>
      <div class="header-actions">
        <!-- 视图切换 -->
        <div class="view-toggle">
          <button class="view-btn" :class="{ active: viewMode === 'dashboard' }" @click="switchView('dashboard')" title="仪表盘视图">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="7" height="9" rx="1"/><rect x="14" y="3" width="7" height="5" rx="1"/>
              <rect x="14" y="12" width="7" height="9" rx="1"/><rect x="3" y="16" width="7" height="5" rx="1"/>
            </svg>
          </button>
          <button class="view-btn" :class="{ active: viewMode === 'kanban' }" @click="switchView('kanban')" title="看板视图">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="5" height="18" rx="1"/><rect x="10" y="3" width="5" height="11" rx="1"/><rect x="17" y="3" width="5" height="15" rx="1"/>
            </svg>
          </button>
          <button class="view-btn" :class="{ active: viewMode === 'grid' }" @click="switchView('grid')" title="卡片视图">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/>
              <rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/>
            </svg>
          </button>
        </div>
        <el-button type="primary" size="small" @click="handleCreate">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          新增项目
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索项目名称/编号"
          clearable
          size="small"
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <el-select v-model="searchStatus" placeholder="状态" clearable size="small" style="width: 110px">
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="inactive" />
        </el-select>
        <el-select v-model="searchStage" placeholder="阶段" clearable size="small" style="width: 120px">
          <el-option v-for="s in stageOptions" :key="s.value" :label="s.label" :value="s.value" />
        </el-select>
      </div>
      <div class="filter-right">
        <span class="result-count">共 <strong>{{ pagination.total }}</strong> 个项目<span v-if="activeKpi" class="filter-active-hint">（已筛选）</span></span>
        <el-button type="primary" size="small" @click="handleSearch">查询</el-button>
        <el-button size="small" @click="handleReset">重置</el-button>
      </div>
    </div>

    <!-- KPI 统计栏（仅仪表盘视图） -->
    <div class="kpi-stats-bar" v-if="viewMode === 'dashboard' && !loading && pagination.total > 0">
      <div class="kpi-card" :class="{ 'kpi-card--active': activeKpi === 'all' }" @click="filterByKpi('all')">
        <div class="kpi-card-inner">
          <div class="kpi-card-icon kpi-card-icon--blue">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
          </div>
          <div class="kpi-card-body">
            <div class="kpi-card-value">{{ pagination.total }}</div>
            <div class="kpi-card-label">全部项目</div>
          </div>
        </div>
      </div>

      <div class="kpi-card" :class="{ 'kpi-card--active': activeKpi === 'active' }" @click="filterByKpi('active')">
        <div class="kpi-card-inner">
          <div class="kpi-card-icon kpi-card-icon--green">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          </div>
          <div class="kpi-card-body">
            <div class="kpi-card-value">{{ activeCount }}</div>
            <div class="kpi-card-label">运营中</div>
          </div>
        </div>
      </div>

      <div class="kpi-card" :class="{ 'kpi-card--active': activeKpi === 'size' }" @click="filterByKpi('size')">
        <div class="kpi-card-inner">
          <div class="kpi-card-icon kpi-card-icon--purple">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>
          </div>
          <div class="kpi-card-body">
            <div class="kpi-card-value">{{ formatBytes(totalDataSize) }}</div>
            <div class="kpi-card-label">上传数据总量</div>
          </div>
        </div>
      </div>

      <div class="kpi-card" :class="{ 'kpi-card--active': activeKpi === 'records' }" @click="filterByKpi('records')">
        <div class="kpi-card-inner">
          <div class="kpi-card-icon kpi-card-icon--amber">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
          </div>
          <div class="kpi-card-body">
            <div class="kpi-card-value">{{ totalRecords.toLocaleString() }}</div>
            <div class="kpi-card-label">上传记录</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 项目详情横幅 -->
    <div v-if="selectedProject" class="project-detail-banner">
      <div class="banner-inner">
        <!-- 左侧: 头像 + 核心信息 -->
        <div class="banner-left">
          <div class="banner-avatar" :style="{ background: getProjectColor(selectedProject.code) }">
            {{ getInitials(selectedProject.name) }}
          </div>
          <div class="banner-core">
            <div class="banner-name">{{ selectedProject.name }}</div>
            <div class="banner-code">
              <code>{{ selectedProject.code }}</code>
            </div>
          </div>
          <span class="banner-stage-badge" :style="{ background: getStageColor(selectedProject.stage).bg, color: getStageColor(selectedProject.stage).text }">
            {{ getStageLabel(selectedProject.stage) }}
          </span>
          <span :class="['banner-status-badge', selectedProject.status === 'active' ? 'badge--on' : 'badge--off']">
            {{ selectedProject.status === 'active' ? '启用' : '禁用' }}
          </span>
        </div>

        <!-- 中间: 关键指标 -->
        <div class="banner-metrics">
          <div class="banner-metric">
            <span class="metric-num">{{ selectedProject.recordCount || 0 }}</span>
            <span class="metric-lbl">数据记录</span>
          </div>
          <div class="metric-sep"></div>
          <div class="banner-metric">
            <span class="metric-num">{{ formatFileSize(selectedProject.totalDataSize || 0) }}</span>
            <span class="metric-lbl">总数据量</span>
          </div>
          <div class="metric-sep"></div>
          <div class="banner-metric">
            <span class="metric-num">{{ selectedProject.onsiteStations?.length || 0 }}</span>
            <span class="metric-lbl">驻场点</span>
          </div>
        </div>

        <!-- 右侧: 团队成员 -->
        <div class="banner-team">
          <div class="team-row" v-if="selectedProject.projectPerson">
            <span class="team-role">👤</span>
            <span class="team-name">{{ selectedProject.projectPerson }}</span>
          </div>
          <div class="team-row" v-if="selectedProject.opsPerson">
            <span class="team-role">💻</span>
            <span class="team-name">{{ selectedProject.opsPerson }}</span>
          </div>
          <div class="team-row" v-if="selectedProject.developerPerson">
            <span class="team-role">🛠</span>
            <span class="team-name">{{ selectedProject.developerPerson }}</span>
          </div>
          <div class="team-row empty" v-if="!selectedProject.projectPerson && !selectedProject.opsPerson && !selectedProject.developerPerson">
            暂无团队成员
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="banner-actions">
          <el-button size="small" @click="selectedProject = null">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            收起
          </el-button>
          <el-button type="primary" size="small" @click="handleEdit(selectedProject)">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
            编辑
          </el-button>
        </div>
      </div>
    </div>

    <!-- ==================== 仪表盘视图 ==================== -->
    <div v-if="viewMode === 'dashboard' && !loading && tableData.length > 0" class="dashboard-container" v-loading="loading">

      <!-- 第一行：大数字 + 阶段分布 -->
      <div class="db-row db-row--hero">
        <!-- 左侧：核心数字 4宫格 -->
        <div class="db-hero-kpis">
          <div class="db-kpi-large" @click="filterByKpi('all')">
            <div class="kpi-ring kpi-ring--total">
              <svg width="72" height="72" viewBox="0 0 72 72">
                <circle cx="36" cy="36" r="30" fill="none" stroke="var(--color-border-light)" stroke-width="5"/>
                <circle cx="36" cy="36" r="30" fill="none" stroke="var(--color-primary)" stroke-width="5"
                  stroke-dasharray="188.5" :stroke-dashoffset="188.5 - (188.5 * Math.min(pagination.total / (pagination.total || 1), 1))"
                  stroke-linecap="round" transform="rotate(-90 36 36)" class="kpi-arc"/>
              </svg>
              <div class="kpi-ring-text">
                <span class="kpi-ring-num">{{ pagination.total }}</span>
                <span class="kpi-ring-label">项目</span>
              </div>
            </div>
            <div class="kpi-meta">
              <span class="kpi-main-label">全部项目</span>
              <span class="kpi-sub-label">当前全部项目</span>
            </div>
          </div>

          <div class="db-kpi-divider"></div>

          <div class="db-kpi-stat-row">
            <div class="db-kpi-stat" @click="filterByKpi('active')">
              <span class="stat-dot" style="background:#22c55e"></span>
              <span class="stat-num">{{ activeCount }}</span>
              <span class="stat-label">运营中</span>
            </div>
            <div class="db-kpi-stat" @click="filterByKpi('size')">
              <span class="stat-dot" style="background:#a855f7"></span>
              <span class="stat-num stat-num--sm">{{ formatBytes(totalDataSize) }}</span>
              <span class="stat-label">数据总量</span>
            </div>
            <div class="db-kpi-stat" @click="filterByKpi('records')">
              <span class="stat-dot" style="background:#f59e0b"></span>
              <span class="stat-num">{{ totalRecords.toLocaleString() }}</span>
              <span class="stat-label">上传记录</span>
            </div>
          </div>
        </div>

        <!-- 右侧：阶段分布 -->
        <div class="db-stage-dist">
          <div class="db-panel-header">
            <span class="db-panel-title">阶段分布</span>
            <span class="db-panel-sub">各阶段项目数量</span>
          </div>
          <div class="stage-bars">
            <div v-for="col in kanbanColumns" :key="col.key" class="stage-bar-item">
              <div class="stage-bar-label">
                <span class="stage-bar-dot" :style="{ background: col.color }"></span>
                <span class="stage-bar-name">{{ col.label }}</span>
              </div>
              <div class="stage-bar-track">
                <div class="stage-bar-fill"
                  :style="{
                    width: `${Math.max((tableData.filter(p => p.stage === col.key).length / (tableData.length || 1)) * 100, 4)}%`,
                    background: col.color
                  }">
                </div>
              </div>
              <span class="stage-bar-count">{{ tableData.filter(p => p.stage === col.key).length }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 第二行：最近项目 + 项目角色分布 -->
      <div class="db-row">
        <!-- 最近项目 -->
        <div class="db-panel db-panel--recent">
          <div class="db-panel-header">
            <span class="db-panel-title">最近项目</span>
            <span class="db-panel-sub">按更新时间排序</span>
          </div>
          <div class="recent-list">
            <div v-for="project in [...tableData].sort((a,b) => new Date(b.updatedAt||0).getTime() - new Date(a.updatedAt||0).getTime()).slice(0, 6)" :key="project.id" class="recent-item" @click="handleCardClick(project)">
              <div class="recent-avatar" :style="{ background: getProjectColor(project.code) }">{{ getInitials(project.name) }}</div>
              <div class="recent-info">
                <span class="recent-name">{{ project.name }}</span>
                <span class="recent-code">{{ project.code }}</span>
              </div>
              <span class="recent-stage" :style="{ background: getStageColor(project.stage).bg, color: getStageColor(project.stage).text }">{{ getStageLabel(project.stage) }}</span>
              <span class="recent-date">{{ formatDate(project.updatedAt) }}</span>
            </div>
            <div v-if="tableData.length === 0" class="recent-empty">暂无项目数据</div>
          </div>
        </div>

        <!-- 项目角色分布 -->
        <div class="db-panel db-panel--roles">
          <div class="db-panel-header">
            <span class="db-panel-title">团队覆盖</span>
            <span class="db-panel-sub">有团队人员的项目</span>
          </div>
          <div class="role-dist-list">
            <div class="role-dist-item">
              <span class="role-icon">👤</span>
              <span class="role-label">项目人员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.projectPerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.projectPerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
            <div class="role-dist-item">
              <span class="role-icon">💻</span>
              <span class="role-label">运维人员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.opsPerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.opsPerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
            <div class="role-dist-item">
              <span class="role-icon">🛠</span>
              <span class="role-label">开发人员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.developerPerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.developerPerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
            <div class="role-dist-item">
              <span class="role-icon">🧪</span>
              <span class="role-label">测试人员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.testerPerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.testerPerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
            <div class="role-dist-item">
              <span class="role-icon">💰</span>
              <span class="role-label">商务人员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.businessPerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.businessPerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
            <div class="role-dist-item">
              <span class="role-icon">🛡</span>
              <span class="role-label">合规专员</span>
              <div class="role-bar-track">
                <div class="role-bar-fill" :style="{ width: `${Math.round((tableData.filter(p => p.compliancePerson).length / (tableData.length||1)) * 100)}%` }"></div>
              </div>
              <span class="role-pct">{{ Math.round((tableData.filter(p => p.compliancePerson).length / (tableData.length||1)) * 100) }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 仪表盘空状态 -->
    <div v-if="viewMode === 'dashboard' && !loading && tableData.length === 0" class="dashboard-empty">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="7" height="9" rx="1"/><rect x="14" y="3" width="7" height="5" rx="1"/><rect x="14" y="12" width="7" height="9" rx="1"/><rect x="3" y="16" width="7" height="5" rx="1"/></svg>
      <p>暂无项目数据，切换其他视图或新增项目</p>
    </div>

    <!-- 项目看板视图 -->
    <div v-if="viewMode === 'kanban' && !loading" class="kanban-container">
      <div class="kanban-board">
        <div
          v-for="col in kanbanColumns"
          :key="col.key"
          class="kanban-column"
        >
          <!-- 栏头 -->
          <div class="kanban-col-header" :style="{ '--col-color': col.color }">
            <span class="col-dot" :style="{ background: col.color }"></span>
            <span class="col-title">{{ col.label }}</span>
            <span class="col-count">{{ kanbanData[col.key]?.length || 0 }}</span>
          </div>
          <!-- 栏内容 -->
          <div class="kanban-col-body">
            <div
              v-for="project in (kanbanData[col.key] || [])"
              :key="project.id"
              class="kanban-card"
              @click="handleKanbanCardClick(project)"
            >
              <!-- 项目名 + 状态 -->
              <div class="kanban-card-top">
                <div class="kanban-card-avatar" :style="{ background: getProjectColor(project.code) }">
                  {{ getInitials(project.name) }}
                </div>
                <div class="kanban-card-title-area">
                  <div class="kanban-card-name">{{ project.name }}</div>
                  <div class="kanban-card-code">{{ project.code }}</div>
                </div>
                <span class="kanban-status-pill" :class="project.status === 'active' ? 'pill--on' : 'pill--off'">
                  {{ project.status === 'active' ? '启用' : '禁用' }}
                </span>
              </div>

              <!-- 描述 -->
              <div class="kanban-card-desc" v-if="project.description">{{ project.description }}</div>

              <!-- 核心指标 -->
              <div class="kanban-metrics">
                <div class="metric-item">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                  <span class="metric-val">{{ project.recordCount || 0 }}</span>
                  <span class="metric-label">条记录</span>
                </div>
                <div class="metric-sep"></div>
                <div class="metric-item">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>
                  <span class="metric-val">{{ formatFileSize(project.totalDataSize || 0) }}</span>
                  <span class="metric-label">数据量</span>
                </div>
              </div>

              <!-- 底部操作 -->
              <div class="kanban-card-footer">
                <div class="kanban-actions">
                  <button class="kanban-action-btn" @click.stop="handleEdit(project)" title="编辑">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                  </button>
                  <button class="kanban-action-btn kanban-action-btn--danger" @click.stop="handleDelete(project)" title="删除">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/></svg>
                  </button>
                </div>
                <span class="kanban-card-date">{{ formatDate(project.createdAt) }}</span>
              </div>
            </div>

            <!-- 空栏 -->
            <div v-if="!kanbanData[col.key]?.length" class="kanban-empty-col">
              暂无数
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 项目卡片网格 -->
    <div v-if="viewMode === 'grid'" class="project-grid" v-loading="loading">
      <div
        v-for="(project, idx) in tableData"
        :key="project.id"
        class="project-card"
        :class="{ 'card-selected': selectedIds.includes(project.id) }"
        :style="{ animationDelay: `${idx * 0.05}s` }"
        @click="handleCardClick(project)"
      >
        <!-- 左侧彩色 accent 条 -->
        <div class="card-accent-bar" :style="{ background: getProjectColor(project.code) }"></div>

        <!-- 卡片主体 -->
        <div class="card-inner">
          <!-- 头部：头像 + 文字 -->
          <div class="card-header">
            <!-- 项目头像 -->
            <div class="card-avatar" :style="{ background: getProjectColor(project.code) }">
              {{ getInitials(project.name) }}
            </div>

            <div class="card-title-area">
              <div class="card-name">{{ project.name }}</div>
              <div class="card-meta-tags">
                <span class="card-code-chip">
                  <svg width="8" height="8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  {{ project.code }}
                </span>
                <span class="card-status-pill" :class="project.status === 'active' ? 'pill--active' : 'pill--inactive'">
                  <span class="pill-dot"></span>
                  {{ project.status === 'active' ? '启用' : '禁用' }}
                </span>
                <span class="record-chip">
                  <svg width="9" height="9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                  {{ project.recordCount || 0 }} 条
                </span>
                <span class="stage-chip" :style="{ background: getStageColor(project.stage).bg, color: getStageColor(project.stage).text }">
                  {{ getStageLabel(project.stage) }}
                </span>
              </div>
            </div>

            <!-- 选择框 -->
            <div
              class="card-checkbox"
              :class="{ checked: selectedIds.includes(project.id) }"
              @click.stop="toggleSelect(project.id)"
            >
              <svg v-if="selectedIds.includes(project.id)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
            </div>
          </div>

          <!-- 描述（无则不显示） -->
          <div class="card-desc" v-if="project.description">{{ project.description }}</div>

          <!-- 信息网格（2列 x 4行） -->
          <div class="card-info-grid">
            <div class="info-cell">
              <div class="info-icon icon-person">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">项目人员</span>
                <span class="info-value">{{ project.projectPerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell">
              <div class="info-icon icon-ops">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">运维人员</span>
                <span class="info-value">{{ project.opsPerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell">
              <div class="info-icon icon-dev">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">开发人员</span>
                <span class="info-value">{{ project.developerPerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell">
              <div class="info-icon icon-test">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">测试人员</span>
                <span class="info-value">{{ project.testerPerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell">
              <div class="info-icon icon-biz">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">商务人员</span>
                <span class="info-value">{{ project.businessPerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell">
              <div class="info-icon icon-compliance">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><polyline points="9 12 11 14 15 10"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">合规专员</span>
                <span class="info-value">{{ project.compliancePerson || '-' }}</span>
              </div>
            </div>
            <div class="info-cell info-cell--full">
              <div class="info-icon icon-location">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              </div>
              <div class="info-text">
                <span class="info-label">公司地点</span>
                <span class="info-value">{{ project.companyAddr || '-' }}</span>
              </div>
            </div>
          </div>

          <!-- 解决方案（全宽） -->
          <div class="card-solution" v-if="project.solution">
            <div class="solution-label">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              解决方案
            </div>
            <div class="solution-text">{{ project.solution }}</div>
          </div>

          <!-- 驻场点标签云 -->
          <div class="card-stations" v-if="project.onsiteStations && project.onsiteStations.length > 0">
            <div class="stations-label">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              驻场点
            </div>
            <div class="stations-cloud">
              <span v-for="(station, sidx) in project.onsiteStations" :key="sidx" class="station-tag">{{ station.location }}</span>
            </div>
          </div>

          <!-- 卡片底部 -->
          <div class="card-footer">
            <div class="card-actions">
              <TableActions :actions="[
                { key: 'edit', label: '编辑', type: 'primary' },
                { key: 'delete', label: '删除', type: 'danger' }
              ]" @action="(key) => handleAction(key, project)" />
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && tableData.length === 0" class="empty-state" style="grid-column: 1 / -1">
        <svg class="empty-icon" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
        <p class="empty-text">暂无项目，点击上方按钮新增</p>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="pagination.total > 0 && viewMode !== 'dashboard'">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[12, 24, 48, 96]"
        layout="sizes, prev, pager, next"
        background
      />
    </div>

    <!-- 编辑/新增侧边栏 -->
    <el-drawer
      v-model="drawerVisible"
      direction="rtl"
      size="580px"
      :with-header="false"
      :destroy-on-close="true"
      class="project-drawer"
    >
      <!-- 侧边栏头部 -->
      <div class="drawer-head">
        <div class="drawer-head-tag">
          <span class="drawer-mode-tag" :class="isEdit ? 'tag--edit' : 'tag--new'">{{ isEdit ? '编辑' : '新增' }}</span>
          <span class="drawer-entity-label">{{ isEdit ? form.name || '项目' : '新增项目' }}</span>
        </div>
        <button class="drawer-close-btn" @click="drawerVisible = false">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <!-- 侧边栏内容 -->
      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top" class="edit-form">
        <!-- 基本信息区块 -->
        <div class="form-section">
          <div class="form-section-header">
            <div class="form-section-dot" style="background: #005eeb"></div>
            基本信息
          </div>
          <div class="form-body">
            <div class="form-grid-2">
              <div class="form-field">
                <label class="field-label">项目名称 <span class="required-mark">*</span></label>
                <el-input v-model="form.name" placeholder="请输入项目名称" size="default" class="field-input" @input="updateDialogPreview" />
              </div>
              <div class="form-field">
                <label class="field-label">项目编号 <span class="required-mark">*</span></label>
                <el-input v-model="form.code" placeholder="如：PRJ-001" size="default" class="field-input" @input="updateDialogPreviewColor" />
              </div>
            </div>
            <div class="form-field">
              <label class="field-label">项目描述</label>
              <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入项目描述（可选）" size="default" class="field-input" />
            </div>
            <div class="form-grid-2">
              <div class="form-field">
                <label class="field-label">项目阶段</label>
                <el-select v-model="form.stage" placeholder="选择项目阶段" size="default" class="field-input">
                  <el-option label="待定中" value="planning" />
                  <el-option label="方案中" value="designing" />
                  <el-option label="部署中" value="deploying" />
                  <el-option label="运营中" value="running" />
                  <el-option label="暂定中" value="paused" />
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">启用状态</label>
                <el-select v-model="form.status" placeholder="选择状态" size="default" class="field-input">
                  <el-option label="启用" value="active" />
                  <el-option label="禁用" value="inactive" />
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 团队成员区块 -->
        <div class="form-section">
          <div class="form-section-header">
            <div class="form-section-dot" style="background: #06b6d4"></div>
            团队成员
            <span class="section-count" v-if="personnelList.length > 0">{{ personnelList.length }} 人可选</span>
          </div>
          <div class="form-body">
            <div class="form-grid-3">
              <div class="form-field">
                <label class="field-label">👤 项目人员</label>
                <el-select v-model="form.projectPerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">💻 运维人员</label>
                <el-select v-model="form.opsPerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">🛠 开发人员</label>
                <el-select v-model="form.developerPerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">🧪 测试人员</label>
                <el-select v-model="form.testerPerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">💰 商务人员</label>
                <el-select v-model="form.businessPerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
              <div class="form-field">
                <label class="field-label">🛡 合规专员</label>
                <el-select v-model="form.compliancePerson" placeholder="选择人员" size="default" class="field-input" :loading="personnelLoading" filterable clearable>
                  <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name">
                    <div class="personnel-opt"><span class="opt-name">{{ p.name }}</span><span class="opt-phone">{{ p.phone || '' }}</span></div>
                  </el-option>
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 项目信息区块 -->
        <div class="form-section">
          <div class="form-section-header">
            <div class="form-section-dot" style="background: #8b5cf6"></div>
            项目信息
          </div>
          <div class="form-body">
            <div class="form-grid-2">
              <div class="form-field">
                <label class="field-label">📍 公司地点</label>
                <el-input v-model="form.companyAddr" placeholder="如：北京市朝阳区" size="default" class="field-input" />
              </div>
              <div class="form-field">
                <label class="field-label">📅 项目周期</label>
                <el-input v-model="form.projectPeriod" placeholder="如：2024.01 - 2024.12" size="default" class="field-input" />
              </div>
              <div class="form-field" style="grid-column: 1 / -1">
                <label class="field-label">💡 解决方案</label>
                <el-input v-model="form.solution" placeholder="简要描述采用的解决方案" size="default" class="field-input" />
              </div>
            </div>
          </div>
        </div>

        <!-- 状态 + 排序 -->
        <div class="form-status-row">
          <div class="form-field">
            <label class="field-label">状态</label>
            <el-radio-group v-model="form.status" size="small">
              <el-radio-button value="active">启用</el-radio-button>
              <el-radio-button value="inactive">禁用</el-radio-button>
            </el-radio-group>
          </div>
          <div class="form-field">
            <label class="field-label">排序</label>
            <el-input-number v-model="form.sort" :min="0" :max="9999" size="small" controls-position="right" style="width: 120px" />
          </div>
        </div>

        <!-- 驻场点区块 -->
        <div class="form-section">
          <div class="form-section-header">
            <div class="form-section-dot" style="background: #8b5cf6"></div>
            驻场点
            <span class="section-count" v-if="form.onsiteStations.length > 0">{{ form.onsiteStations.length }} 个</span>
          </div>
          <div class="form-body">
            <div class="stations-list">
              <div v-if="form.onsiteStations.length === 0" class="stations-empty">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                <span>暂无驻场点，填写下方信息可添加</span>
              </div>
              <div v-for="(station, idx) in form.onsiteStations" :key="idx" class="station-card">
                <div class="station-number">{{ idx + 1 }}</div>
                <div class="station-fields">
                  <div class="station-field">
                    <span class="station-field-label">📍 场地</span>
                    <el-input v-model="station.location" placeholder="场地名称" size="small" class="station-input" />
                  </div>
                  <div class="station-field">
                    <span class="station-field-label">👤 人员</span>
                    <el-select
                      v-model="station.person"
                      placeholder="选择人员"
                      size="small"
                      class="station-input"
                      :loading="personnelLoading"
                      filterable
                      clearable
                      @change="(val) => handlePersonnelSelect(station, val)"
                    >
                      <el-option
                        v-for="p in personnelList"
                        :key="p.id"
                        :label="p.name"
                        :value="p.name"
                      >
                        <div class="personnel-option">
                          <span class="personnel-name">{{ p.name }}</span>
                          <span class="personnel-phone">{{ p.phone || '无电话' }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </div>
                  <div class="station-field">
                    <span class="station-field-label">📞 联系</span>
                    <el-input v-model="station.phone" placeholder="手机/电话" size="small" class="station-input" />
                  </div>
                </div>
                <button type="button" class="station-remove-btn" @click="removeStation(idx)" title="删除此驻场点">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
              </div>
            </div>
            <button type="button" class="add-station-btn" @click="addStation">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              添加驻场点
            </button>
          </div>
        </div>
        </el-form>
      </div>

      <!-- 侧边栏底部 -->
      <div class="drawer-foot">
        <el-button size="default" @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" size="default" :loading="submitting" @click="confirmSubmit" class="submit-btn">
          <svg v-if="!submitting" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          {{ isEdit ? '保存修改' : '创建项目' }}
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ProjectApi, type Project, type CreateProjectReq, type UpdateProjectReq, type OnSiteStation } from '@/api/project'
import { PersonnelApi, type Personnel } from '@/api/personnel'
import TableActions from '@/components/TableActions.vue'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Project[]>([])
const drawerVisible = ref(false)
const dialogPreviewColor = ref('#005eeb')
const dialogPreviewInitials = ref('??')

const updateDialogPreview = () => {
  dialogPreviewInitials.value = getInitials(form.name)
}
const updateDialogPreviewColor = () => {
  dialogPreviewColor.value = getProjectColor(form.code)
}
const isEdit = ref(false)
const formRef = ref()
const selectedIds = ref<number[]>([])

const searchKeyword = ref('')
const searchStatus = ref('')
const searchStage = ref('')
const pagination = reactive({ page: 1, pageSize: 24, total: 0 })
const activeKpi = ref<string | null>(null)
const selectedProject = ref<Project | null>(null)

const stageOptions = [
  { value: 'planning', label: '待定中' },
  { value: 'designing', label: '方案中' },
  { value: 'deploying', label: '部署中' },
  { value: 'running', label: '运营中' },
  { value: 'paused', label: '暂定中' },
]

// KPI Stats
const activeCount = computed(() => tableData.value.filter(p => p.status === 'active').length)
const totalStations = computed(() => tableData.value.reduce((sum, p) => sum + (p.onsiteStations?.length || 0), 0))
const totalRecords = computed(() => tableData.value.reduce((sum, p) => sum + (p.recordCount || 0), 0))
const totalDataSize = computed(() => tableData.value.reduce((sum, p) => sum + (p.totalDataSize || 0), 0))

// 自动单位换算
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(i > 0 ? 1 : 0)} ${units[i]}`
}

// ========== 看板视图 ==========
const viewMode = ref<'kanban' | 'grid' | 'dashboard'>('dashboard')

const kanbanColumns = [
  { key: 'planning', label: '待定中', color: '#94a3b8' },
  { key: 'designing', label: '方案中', color: '#6366f1' },
  { key: 'deploying', label: '部署中', color: '#f59e0b' },
  { key: 'running', label: '运营中', color: '#22c55e' },
  { key: 'paused', label: '暂定中', color: '#ef4444' },
]

const kanbanData = ref<Record<string, Project[]>>({})

const switchView = (mode: 'kanban' | 'grid' | 'dashboard') => {
  viewMode.value = mode
  if (mode === 'kanban') loadKanbanData()
  if (mode === 'grid') loadData()
  if (mode === 'dashboard') loadData()
}

const loadKanbanData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.getKanbanList()
    const raw: any = res
    const list: any[] = Array.isArray(raw) ? raw : (raw?.data ?? [])
    const data: Record<string, Project[]> = {}
    for (const col of kanbanColumns) data[col.key] = []
    for (const p of list) {
      if (!p || typeof p !== 'object') continue
      const stage = (p as Project).stage || 'planning'
      if (!data[stage]) data[stage] = []
      data[stage].push(p as Project)
    }
    kanbanData.value = data
  } catch (e) {
    kanbanData.value = {}
  } finally { loading.value = false }
}

const handleKanbanCardClick = (_project: Project) => {}

const formatFileSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, '0')}/${String(d.getDate()).padStart(2, '0')}`
}

const stageColorMap: Record<string, { bg: string; text: string }> = {
  planning: { bg: '#f1f5f9', text: '#64748b' },
  designing: { bg: '#eef2ff', text: '#6366f1' },
  deploying: { bg: '#fef3c7', text: '#d97706' },
  running: { bg: '#f0fdf4', text: '#16a34a' },
  paused: { bg: '#fef2f2', text: '#dc2626' },
}

const getStageColor = (stage: string) => stageColorMap[stage] || { bg: '#f1f5f9', text: '#64748b' }
const getStageLabel = (stage: string) => ({ planning: '待定中', designing: '方案中', deploying: '部署中', running: '运营中', paused: '暂定中' })[stage] || stage

const form = reactive<CreateProjectReq & { id?: number }>({
  name: '', code: '', description: '', status: 'active', sort: 0,
  projectPerson: '', opsPerson: '', developerPerson: '', testerPerson: '', businessPerson: '', compliancePerson: '',
  solution: '', companyAddr: '', projectPeriod: '',
  onsiteStations: [],
})

const formRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入项目编号', trigger: 'blur' }],
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
    })
    let list: Project[] = res.data?.items || res.data || []
    // 前端过滤：阶段
    if (searchStage.value) {
      list = list.filter(p => p.stage === searchStage.value)
    }
    // 前端排序：按记录数
    if (activeKpi.value === 'records') {
      list = [...list].sort((a, b) => (b.recordCount || 0) - (a.recordCount || 0))
    }
    tableData.value = list
    pagination.total = res.data?.total || 0
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; activeKpi.value = null; loadData() }
const handleReset = () => { searchKeyword.value = ''; searchStatus.value = ''; searchStage.value = ''; activeKpi.value = null; handleSearch() }

const filterByKpi = (type: string) => {
  if (activeKpi.value === type) {
    activeKpi.value = null
    searchStage.value = ''
    searchStatus.value = ''
  } else {
    activeKpi.value = type
    if (type === 'all') { searchStage.value = ''; searchStatus.value = '' }
    else if (type === 'active') { searchStage.value = ''; searchStatus.value = 'active' }
    else if (type === 'size') { searchStage.value = ''; searchStatus.value = '' }
    else if (type === 'records') { searchStage.value = ''; searchStatus.value = '' }
    else if (type === 'stations') { searchStage.value = ''; searchStatus.value = '' }
  }
  pagination.page = 1
  loadData()
}
const getInitials = (name: string) => {
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

// 预定义颜色盘，每个项目代码对应一个颜色
const projectColors = [
  '#005eeb', '#06b6d4', '#8b5cf6', '#ec4899',
  '#f59e0b', '#10b981', '#3b82f6', '#ef4444',
  '#6366f1', '#14b8a6', '#f97316', '#84cc16'
]

const getProjectColor = (code: string) => {
  if (!code) return projectColors[0]
  let hash = 0
  for (let i = 0; i < code.length; i++) {
    hash = ((hash << 5) - hash) + code.charCodeAt(i)
    hash |= 0
  }
  return projectColors[Math.abs(hash) % projectColors.length]
}

const handleCardClick = (project: Project) => {
  if (selectedProject.value?.id === project.id) {
    selectedProject.value = null
  } else {
    selectedProject.value = project
  }
}
const toggleSelect = (id: number) => {
  const idx = selectedIds.value.indexOf(id)
  if (idx === -1) selectedIds.value.push(id)
  else selectedIds.value.splice(idx, 1)
}

const personnelList = ref<Personnel[]>([])
const personnelLoading = ref(false)

const loadPersonnelList = async () => {
  personnelLoading.value = true
  try {
    const res = await PersonnelApi.getAll()
    if (res.code === 200) {
      personnelList.value = res.data || []
    }
  } finally {
    personnelLoading.value = false
  }
}

const handlePersonnelSelect = (station: OnSiteStation, name: string) => {
  if (!name) return
  const p = personnelList.value.find(x => x.name === name)
  if (p && p.phone) {
    station.phone = p.phone
  }
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    id: undefined, name: '', code: '', description: '', status: 'active', stage: 'planning', sort: 0,
    projectPerson: '', opsPerson: '', developerPerson: '', testerPerson: '', businessPerson: '', compliancePerson: '',
    solution: '', companyAddr: '', projectPeriod: '', onsiteStations: [],
  })
  dialogPreviewColor.value = '#005eeb'
  dialogPreviewInitials.value = '??'
  loadPersonnelList()
  drawerVisible.value = true
}

const handleEdit = (row: Project) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id, name: row.name, code: row.code,
    description: row.description || '', status: row.status, stage: row.stage || 'planning', sort: row.sort || 0,
    projectPerson: row.projectPerson || '', opsPerson: row.opsPerson || '',
    developerPerson: row.developerPerson || '', testerPerson: row.testerPerson || '',
    businessPerson: row.businessPerson || '', compliancePerson: row.compliancePerson || '',
    solution: row.solution || '',
    companyAddr: row.companyAddr || '', projectPeriod: row.projectPeriod || '',
    onsiteStations: row.onsiteStations ? [...row.onsiteStations.map(s => ({ ...s }))] : [],
  })
  dialogPreviewColor.value = getProjectColor(row.code)
  dialogPreviewInitials.value = getInitials(row.name)
  loadPersonnelList()
  drawerVisible.value = true
}

const handleDelete = async (row: Project) => {
  try {
    await ElMessageBox.confirm(`确定要删除项目"${row.name}"吗？`, '删除确认', {
      confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning'
    })
    await ProjectApi.del(row.id)
    ElMessage.success('删除成功'); loadKanbanData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || '删除失败') }
}

const addStation = () => {
  if (!form.onsiteStations) form.onsiteStations = []
  form.onsiteStations.push({ location: '', person: '', phone: '' })
}

const removeStation = (idx: number) => {
  form.onsiteStations.splice(idx, 1)
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      await ProjectApi.update(form as UpdateProjectReq)
      ElMessage.success('更新成功')
    } else {
      await ProjectApi.create(form as CreateProjectReq)
      ElMessage.success('创建成功')
    }
    drawerVisible.value = false
    loadKanbanData()
  } catch (e: any) { ElMessage.error(e.message || '操作失败') }
  finally { submitting.value = false }
}

const handleAction = (key: string, row: Project) => {
  if (key === 'edit') handleEdit(row)
  else if (key === 'delete') handleDelete(row)
}

watch(() => pagination.page, () => { if (viewMode.value === 'grid') loadData() })
watch(() => pagination.pageSize, () => { pagination.page = 1; if (viewMode.value === 'grid') loadData() })
onMounted(() => { loadData(); loadKanbanData() })
</script>

<script lang="ts">
export default { name: 'ProjectList' }
</script>

<style scoped lang="scss">
/* ==================== 页面布局 ==================== */
.project-page {
  padding: var(--space-4);
  min-height: 100vh;
  background: var(--color-page-bg);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  overflow: visible;
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

/* ==================== 筛选栏 ==================== */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
  backdrop-filter: blur(8px);
  flex-wrap: wrap;
}

.filter-left {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.filter-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.search-input { width: 220px; }

.result-count {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
  strong { color: var(--color-text-primary); font-weight: 700; }
}

/* ==================== 项目详情横幅 ==================== */
.project-detail-banner {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-3);
  overflow: hidden;
  box-shadow: var(--shadow-xs);
  animation: card-rise 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.banner-inner {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 18px 24px;
  flex-wrap: wrap;
}

.banner-left {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-shrink: 0;
}

.banner-avatar {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-family: 'Manrope', sans-serif;
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.5px;
  flex-shrink: 0;
}

.banner-core { flex: 0 0 auto; }
.banner-name {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 800;
  color: var(--color-text-primary);
  margin-bottom: 4px;
  letter-spacing: -0.3px;
}
.banner-code {
  code {
    font-family: 'SF Mono', monospace;
    font-size: 11px;
    color: var(--color-primary);
    background: var(--color-primary-light-9);
    padding: 2px 7px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(0,94,235,0.15);
  }
}

.banner-stage-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: var(--radius-full);
  font-size: 11px;
  font-weight: 700;
}

.banner-status-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 600;
  &.badge--on { background: rgba(34,197,94,0.1); color: var(--color-success); }
  &.badge--off { background: rgba(239,68,68,0.08); color: var(--color-danger); }
}

.banner-metrics {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
  padding: 0 20px;
  border-left: 1px solid var(--color-border-light);
  border-right: 1px solid var(--color-border-light);
}

.banner-metric {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  min-width: 60px;
}

.metric-num {
  font-family: 'Manrope', sans-serif;
  font-size: 20px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1;
  letter-spacing: -0.5px;
}

.metric-lbl {
  font-size: 11px;
  color: var(--color-text-muted);
  font-weight: 500;
  white-space: nowrap;
}

.metric-sep {
  width: 1px;
  height: 36px;
  background: var(--color-border-light);
  flex-shrink: 0;
}

.banner-team {
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex: 0 0 160px;
  padding-left: 20px;
  border-left: 1px solid var(--color-border-light);
}

.team-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--color-text-secondary);
  &.empty { color: var(--color-text-muted); font-size: 11px; font-style: italic; }
}
.team-role { font-size: 12px; }
.team-name { font-weight: 500; color: var(--color-text-primary); }

.banner-actions {
  display: flex;
  gap: 6px;
  margin-left: auto;
  flex-shrink: 0;
  .el-button { border-radius: var(--radius-sm); }
}

/* ==================== 仪表盘视图 ==================== */
.dashboard-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  animation: card-rise 0.4s ease both;
}

.db-row {
  display: grid;
  gap: 16px;

  &.db-row--hero {
    grid-template-columns: 340px 1fr;
  }
}

.db-row:not(.db-row--hero) {
  grid-template-columns: 1fr 380px;
}

/* 英雄数字区 */
.db-hero-kpis {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.db-kpi-divider {
  width: 80%;
  height: 1px;
  background: var(--color-border-light);
}

.db-kpi-large {
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover { transform: scale(1.02); }
}

.kpi-ring {
  position: relative;
  flex-shrink: 0;
}

.kpi-ring-text {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.kpi-ring-num {
  font-family: 'Manrope', monospace;
  font-size: 24px;
  font-weight: 900;
  color: var(--color-text-primary);
  line-height: 1;
}

.kpi-ring-label {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-muted);
  margin-top: 2px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.kpi-arc {
  transition: stroke-dashoffset 1s cubic-bezier(0.4, 0, 0.2, 1);
}

.kpi-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kpi-main-label {
  font-family: 'Manrope', monospace;
  font-size: 18px;
  font-weight: 800;
  color: var(--color-text-primary);
}

.kpi-sub-label {
  font-size: 12px;
  color: var(--color-text-muted);
}

.db-kpi-stat-row {
  display: flex;
  gap: 16px;
  width: 100%;
  justify-content: center;
}

.db-kpi-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  transition: background 0.15s ease;

  &:hover { background: var(--color-surface-2); }
}

.stat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.stat-num {
  font-family: 'Manrope', monospace;
  font-size: 20px;
  font-weight: 900;
  color: var(--color-text-primary);
  line-height: 1;

  &--sm {
    font-size: 15px;
    font-weight: 800;
    letter-spacing: -0.3px;
  }
}

.stat-label {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

/* 阶段分布面板 */
.db-stage-dist {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.db-panel-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-border-light);
}

.db-panel-title {
  font-family: 'Manrope', sans-serif;
  font-size: 14px;
  font-weight: 800;
  color: var(--color-text-primary);
}

.db-panel-sub {
  font-size: 11px;
  color: var(--color-text-muted);
}

.stage-bars {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.stage-bar-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stage-bar-label {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 80px;
  flex-shrink: 0;
}

.stage-bar-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.stage-bar-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.stage-bar-track {
  flex: 1;
  height: 8px;
  background: var(--color-surface-2);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.stage-bar-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 4%;
}

.stage-bar-count {
  font-family: 'Manrope', monospace;
  font-size: 13px;
  font-weight: 800;
  color: var(--color-text-primary);
  width: 20px;
  text-align: right;
  flex-shrink: 0;
}

/* 通用面板 */
.db-panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: 20px;
}

/* 最近项目 */
.recent-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 4px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background 0.15s ease;

  &:hover { background: var(--color-surface-2); }
}

.recent-avatar {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Manrope', monospace;
  font-size: 12px;
  font-weight: 800;
  color: white;
  flex-shrink: 0;
}

.recent-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.recent-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recent-code {
  font-size: 11px;
  color: var(--color-text-muted);
}

.recent-stage {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: var(--radius-full);
  white-space: nowrap;
}

.recent-date {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.recent-empty {
  text-align: center;
  padding: 24px;
  color: var(--color-text-muted);
  font-size: 13px;
}

/* 角色分布 */
.role-dist-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 4px;
}

.role-dist-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-icon {
  font-size: 14px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.role-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  width: 70px;
  flex-shrink: 0;
}

.role-bar-track {
  flex: 1;
  height: 6px;
  background: var(--color-surface-2);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.role-bar-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: var(--radius-full);
  transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.role-pct {
  font-family: 'Manrope', monospace;
  font-size: 12px;
  font-weight: 800;
  color: var(--color-text-primary);
  width: 32px;
  text-align: right;
  flex-shrink: 0;
}

/* 仪表盘空状态 */
.dashboard-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 60px;
  color: var(--color-text-muted);
  font-size: 14px;
  animation: card-rise 0.4s ease both;

  svg { opacity: 0.4; }
}


.kanban-container {
  overflow-x: auto;
  overflow-y: visible;
  padding-bottom: var(--space-3);
  animation: card-rise 0.4s ease both;
}

.kanban-board {
  display: flex;
  gap: 12px;
  min-width: max-content;
  height: calc(100vh - 240px);
}

.kanban-column {
  width: 280px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-surface-2);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
}

.kanban-col-header {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 11px 14px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.col-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.col-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  flex: 1;
}

.col-count {
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-muted);
  background: var(--color-surface-3);
  padding: 2px 7px;
  border-radius: var(--radius-full);
  font-family: 'Manrope', sans-serif;
}

.kanban-col-body {
  flex: 1;
  overflow-y: auto;
  padding: 10px 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.1); border-radius: 2px; }
}

.kanban-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: 13px;
  cursor: pointer;
  transition: all 0.18s ease;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);

  &:hover {
    border-color: var(--color-primary);
    box-shadow: 0 4px 16px rgba(0,0,0,0.08);
    transform: translateY(-1px);
  }
}

.kanban-card-top {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 9px;
}

.kanban-card-avatar {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 800;
  flex-shrink: 0;
  font-family: 'Manrope', sans-serif;
}

.kanban-card-title-area {
  flex: 1;
  min-width: 0;
}

.kanban-card-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

.kanban-card-code {
  font-size: 11px;
  color: var(--color-text-muted);
  font-family: 'Manrope', monospace;
  margin-top: 1px;
}

.kanban-status-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
  letter-spacing: 0.2px;
}
.pill--on { background: rgba(34,197,94,0.1); color: #22c55e; }
.pill--off { background: rgba(239,68,68,0.1); color: #ef4444; }

.kanban-card-desc {
  font-size: 11.5px;
  color: var(--color-text-muted);
  line-height: 1.4;
  margin-bottom: 9px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.kanban-metrics {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 9px;
  background: var(--color-surface-2);
  border-radius: var(--radius-sm);
  margin-bottom: 9px;
}

.metric-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--color-text-muted);

  svg { flex-shrink: 0; }
}

.metric-val {
  font-size: 13px;
  font-weight: 800;
  color: var(--color-text-primary);
  font-family: 'Manrope', sans-serif;
}

.metric-label {
  font-size: 10px;
  color: var(--color-text-muted);
}

.metric-sep {
  width: 1px;
  height: 14px;
  background: var(--color-border);
  flex-shrink: 0;
}

.kanban-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.kanban-actions {
  display: flex;
  gap: 3px;
}

.kanban-action-btn {
  width: 26px;
  height: 26px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;

  &:hover { background: var(--color-surface-3); color: var(--color-primary); }
  &--danger:hover { background: rgba(239,68,68,0.1); color: #ef4444; }
}

.kanban-card-date {
  font-size: 10px;
  color: var(--color-text-muted);
  font-family: 'Manrope', monospace;
}

.kanban-empty-col {
  padding: 20px 0;
  text-align: center;
  font-size: 12px;
  color: var(--color-text-muted);
}

/* ==================== 视图切换 ==================== */
.view-toggle {
  display: flex;
  gap: 2px;
  padding: 3px;
  background: var(--color-surface-3);
  border-radius: var(--radius-md);
}

.view-btn {
  width: 30px;
  height: 26px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;

  &:hover { color: var(--color-text-primary); }
  &.active { background: var(--color-surface); color: var(--color-primary); box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
}

/* ==================== KPI 统计栏 ==================== */
.kpi-stats-bar {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-3);
  animation: card-rise 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.kpi-card {
  background: var(--color-surface);
  border-radius: 14px;
  border: 1px solid var(--color-border-light);
  padding: 18px 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;

  &:hover {
    border-color: transparent;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    transform: translateY(-2px);
  }

  &--active {
    border-color: transparent;
    box-shadow: 0 0 0 2px var(--color-primary), 0 4px 20px rgba(0, 94, 235, 0.15);
    transform: translateY(-2px);
  }
}

.kpi-card-inner {
  display: flex;
  align-items: center;
  gap: 14px;
}

.kpi-card-icon {
  width: 42px;
  height: 42px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &--blue { background: rgba(0, 94, 235, 0.1); color: #005eeb; }
  &--green { background: rgba(22, 163, 74, 0.1); color: #16a34a; }
  &--purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
  &--amber { background: rgba(245, 158, 11, 0.1); color: #d97706; }
}

.kpi-card-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.kpi-card-value {
  font-family: 'Manrope', 'DM Sans', sans-serif;
  font-size: 22px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1.1;
  letter-spacing: -0.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.kpi-card-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-muted);
  white-space: nowrap;
}

/* ==================== 项目卡片网格 ==================== */
.project-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(440px, 1fr));
  gap: var(--space-3);
}

/* ==================== 项目卡片 ==================== */
.project-card {
  position: relative;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
  display: flex;
  flex-direction: column;
  animation: card-rise 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;

  &:hover {
    border-color: transparent;
    box-shadow: 0 8px 32px rgba(0, 94, 235, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06);
    transform: translateY(-3px) scale(1.01);
  }

  &.card-selected {
    border-color: transparent;
    box-shadow: 0 0 0 2px var(--color-primary), 0 8px 32px rgba(0, 94, 235, 0.15);
    background: var(--color-primary-light-9);
  }
}

@keyframes card-rise {
  from {
    opacity: 0;
    transform: translateY(16px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 左侧彩色 accent 条 */
.card-accent-bar {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  border-radius: 4px 0 0 4px;
  transition: width 0.25s ease;
}

.project-card:hover .card-accent-bar {
  width: 5px;
}

/* 卡片主体（右侧偏移，留给 accent bar） */
.card-inner {
  padding: var(--space-3) var(--space-3) var(--space-3) calc(var(--space-3) + 6px);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  flex: 1;
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  position: relative;
}

.card-checkbox {
  position: absolute;
  top: 0;
  right: 0;
  width: 20px;
  height: 20px;
  border-radius: var(--radius-sm);
  border: 1.5px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  background: var(--color-surface);
  cursor: pointer;
  flex-shrink: 0;

  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: white;
    transform: scale(1.1);
  }

  &:hover:not(.checked) {
    border-color: var(--color-primary);
    transform: scale(1.05);
  }
}

/* 项目头像 */
.card-avatar {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Manrope', sans-serif;
  font-weight: 800;
  font-size: 18px;
  color: white;
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(255,255,255,0.25) 0%, transparent 60%);
  }
}

/* 项目文字区 */
.card-title-area {
  flex: 1;
  min-width: 0;
}

.card-code {
  display: inline-block;
  font-family: 'SF Mono', monospace;
  font-size: 10.5px;
  font-weight: 700;
  color: var(--color-primary);
  background: var(--color-primary-light-9);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  letter-spacing: 0.5px;
  margin-bottom: 5px;
}

.card-name {
  font-size: 17px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1.35;
  margin-bottom: 4px;
  letter-spacing: -0.2px;
}

.card-desc {
  font-size: 12px;
  color: var(--color-text-muted);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 项目代码 + 状态 + 记录数标签行 */
.card-meta-tags {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 5px;
  margin-top: 3px;
}

.card-code-chip {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 7px;
  background: var(--color-primary-light-9);
  border: 1px solid rgba(0, 94, 235, 0.12);
  border-radius: var(--radius-sm);
  font-family: 'SF Mono', monospace;
  font-size: 10px;
  font-weight: 700;
  color: var(--color-primary);
  letter-spacing: 0.3px;
  svg { color: var(--color-primary); opacity: 0.7; }
}

.card-status-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-size: 10.5px;
  font-weight: 700;

  &.pill--active {
    background: var(--color-success-bg);
    color: var(--color-success);
    .pill-dot {
      background: var(--color-success);
      box-shadow: 0 0 4px var(--color-success);
      animation: pulse-dot 2s ease-in-out infinite;
    }
  }

  &.pill--inactive {
    background: var(--gray-100);
    color: var(--color-text-muted);
    .pill-dot { background: var(--gray-400); }
  }
}

.pill-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  flex-shrink: 0;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.85); }
}

.record-chip {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 7px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-full);
  font-size: 10.5px;
  font-weight: 600;
  color: var(--color-text-secondary);
  svg { color: var(--color-text-muted); }
}

.stage-chip {
  display: inline-flex;
  align-items: center;
  padding: 2px 7px;
  border-radius: var(--radius-full);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

/* 状态 + 记录数标签 - removed, replaced by .card-meta-tags above */

/* 信息网格（2列） */
.card-info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-1) var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface-2);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
}

.info-cell {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  min-width: 0;
}

.info-cell--full {
  grid-column: 1 / -1;
  padding-top: var(--space-2);
  border-top: 1px dashed var(--color-border-light);
  margin-top: var(--space-1);
}

.info-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.icon-person    { background: rgba(59, 130, 246, 0.12); color: var(--chart-blue); }
  &.icon-ops       { background: rgba(22, 163, 74, 0.12);  color: var(--chart-green); }
  &.icon-dev      { background: rgba(139, 92, 246, 0.12); color: var(--chart-purple); }
  &.icon-test     { background: rgba(245, 158, 11, 0.12);  color: var(--chart-amber); }
  &.icon-biz      { background: rgba(236, 72, 153, 0.12);  color: #ec4899; }
  &.icon-compliance { background: rgba(0, 94, 235, 0.12); color: var(--chart-blue); }
  &.icon-location  { background: rgba(249, 115, 22, 0.12); color: var(--chart-amber); }
}

.info-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.info-label {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--color-text-muted);
}

.info-value {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 公司地点（全宽） */
.card-location {
  padding: 0 var(--space-3);
}

.location-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: var(--color-text-muted);
  margin-bottom: 4px;

  svg { color: var(--chart-amber); }
}

.location-text {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 解决方案 */
.card-solution {
  padding: 0 var(--space-3);
}

.solution-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: var(--color-text-muted);
  margin-bottom: 5px;

  svg { color: var(--color-info); }
}

.solution-text {
  font-size: 12.5px;
  font-weight: 500;
  color: var(--color-text-secondary);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 驻场点标签云 */
.card-stations {
  padding: 0 var(--space-3) var(--space-3);
}

.stations-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: var(--color-text-muted);
  margin-bottom: 6px;
}

.stations-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.station-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 9px;
  background: var(--color-primary-light-9);
  border: 1px solid rgba(0, 94, 235, 0.15);
  border-radius: var(--radius-full);
  font-size: 11px;
  font-weight: 600;
  color: var(--color-primary);
  transition: all 0.15s ease;

  &:hover {
    background: rgba(0, 94, 235, 0.15);
    transform: scale(1.03);
  }
}

/* 卡片底部操作条 */
.card-footer {
  padding: var(--space-2) var(--space-3);
  border-top: 1px solid var(--color-border-light);
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: var(--space-1);
  background: var(--color-surface-2);
  margin-top: auto;
  border-radius: 0 0 calc(var(--radius-lg) - 1px) calc(var(--radius-lg) - 1px);
}

.card-actions {
  display: flex;
  gap: var(--space-1);
}

/* ==================== 空状态 ==================== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  gap: var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  animation: card-rise 0.5s ease both;
}

.empty-icon {
  color: var(--gray-300);
  opacity: 0.5;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.empty-text {
  font-size: 15px;
  color: var(--color-text-muted);
  margin: 0;
  font-weight: 500;
}

/* ==================== 分页 ==================== */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
  box-shadow: var(--shadow-xs);
  border: 1px solid var(--color-border-light);
}

/* ==================== 响应式 ==================== */
@media (max-width: 640px) {
  :deep(.project-drawer) { width: 100% !important; }
  .form-grid-2 { grid-template-columns: 1fr; }
  .form-grid-3 { grid-template-columns: 1fr 1fr; }
  .form-status-row { flex-direction: column; align-items: flex-start; gap: 8px; }
  .station-fields { grid-template-columns: 1fr; }
}

.divider-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-3);
}

.form-hint {
  margin-left: var(--space-2);
  font-size: 11px;
  color: var(--color-text-muted);
}

.stations-form {
  padding: 0 4px;
}

.station-form-item {
  padding: var(--space-2) 0;
  border-bottom: 1px dashed var(--color-border-light);
  &:last-child { border-bottom: none; }
}

.stations-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-4);
  background: var(--color-surface-2);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 13px;
}

/* ==================== 侧边栏 ==================== */
:deep(.project-drawer) {
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
  padding: 12px;
  background: var(--color-surface-2);

  &::-webkit-scrollbar { width: 3px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb { background: var(--gray-200); border-radius: 2px; }
}

/* 表单通用区块 */
.form-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
  margin-bottom: 8px;
}

.form-section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.form-section-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.section-count {
  margin-left: auto;
  font-size: 10px;
  font-weight: 700;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  padding: 1px 6px;
  border-radius: var(--radius-full);
  text-transform: none;
  letter-spacing: 0;
}

.form-body {
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 8px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.field-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.required-mark {
  color: var(--color-danger);
}

/* 状态+排序行 */
.form-status-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  margin-bottom: 8px;
}

/* 驻场点卡片 */
.stations-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stations-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  background: var(--color-surface);
  border: 1.5px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 12px;
}

.station-card {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--color-primary);
    box-shadow: 0 2px 8px rgba(0, 94, 235, 0.08);
  }
}

.station-number {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  font-size: 10px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.station-fields {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 6px;
}

.station-field {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.station-field-label {
  font-size: 9px;
  font-weight: 700;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.station-remove-btn {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  flex-shrink: 0;

  &:hover {
    background: var(--color-danger-bg);
    color: var(--color-danger);
  }
}

.add-station-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  width: 100%;
  padding: 8px;
  border: 1.5px dashed var(--color-border);
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--color-primary);
    color: var(--color-primary);
    background: var(--color-primary-light-9);
  }
}

/* 人员选项 */
.personnel-opt {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.opt-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.opt-phone {
  font-size: 11px;
  color: var(--color-text-muted);
}

.personnel-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.personnel-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.personnel-phone {
  font-size: 11px;
  color: var(--color-text-muted);
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

/* ==================== 表单区块 ==================== */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  animation: drawer-form-in 0.4s ease both 0.05s;
}

@keyframes drawer-form-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.form-section {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.form-section-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.form-section-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.section-count {
  margin-left: auto;
  font-size: 10px;
  font-weight: 700;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  padding: 1px 8px;
  border-radius: var(--radius-full);
  text-transform: none;
  letter-spacing: 0;
}

.form-body {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;

  &.form-field--inline {
    flex-direction: row;
    align-items: center;
    gap: var(--space-3);
  }
}

.field-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  letter-spacing: 0.2px;
}

.required-mark {
  color: var(--color-danger);
  margin-left: 2px;
}

.field-prefix-icon {
  font-size: 13px;
  opacity: 0.6;
}

/* 状态切换按钮 */
.status-toggle {
  display: flex;
  gap: 4px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 3px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
  color: var(--color-text-secondary);

  &.active {
    background: var(--color-primary);
    color: white;
    box-shadow: 0 2px 8px rgba(0, 94, 235, 0.3);
  }

  &:not(.active):hover {
    background: var(--color-surface-2);
    color: var(--color-text-primary);
  }
}

.toggle-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--gray-300);

  &.active {
    background: #22c55e;
    box-shadow: 0 0 6px rgba(34, 197, 94, 0.6);
    animation: pulse-dot 2s ease-in-out infinite;
  }
}

/* 排序输入 */
.sort-input-wrap {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.sort-hint {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
}

/* 内联行 */
.form-row-inline {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: var(--space-4);
  align-items: center;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: var(--space-3) var(--space-4);
}

/* 驻场点卡片 */
.stations-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.stations-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-4);
  background: var(--color-surface);
  border: 1.5px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 13px;
}

.station-card {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
  animation: station-in 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) both;

  &:hover {
    border-color: var(--color-primary);
    box-shadow: 0 2px 8px rgba(0, 94, 235, 0.08);
  }
}

@keyframes station-in {
  from { opacity: 0; transform: scale(0.97) translateY(-4px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.station-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--color-primary-light-9);
  color: var(--color-primary);
  font-size: 11px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-family: 'Manrope', monospace;
}

.station-fields {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: var(--space-2);
}

.station-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.station-field-label {
  font-size: 10px;
  font-weight: 700;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.station-input {
  .el-input__wrapper {
    border-radius: var(--radius-sm);
  }
}

.personnel-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: 2px 0;
}

.personnel-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.personnel-phone {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 人员下拉选项 */
.personnel-opt {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  padding: 1px 0;
}

.opt-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.opt-phone {
  font-size: 11px;
  color: var(--color-text-muted);
}

.station-remove-btn {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  flex-shrink: 0;

  &:hover {
    background: var(--color-danger-bg);
    color: var(--color-danger);
    transform: scale(1.1);
  }
}

.add-station-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: var(--space-3);
  border: 1.5px dashed var(--color-border);
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: var(--space-2);

  &:hover {
    border-color: var(--color-primary);
    color: var(--color-primary);
    background: var(--color-primary-light-9);
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }
}

/* 响应式 */
@media (max-width: 1366px) {
  .project-page { padding: var(--space-3); }
  .project-grid { grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: var(--space-2); }
}

@media (max-width: 640px) {
  .form-grid-2 { grid-template-columns: 1fr; }
  .form-row-inline { grid-template-columns: 1fr; }
  .station-fields { grid-template-columns: 1fr; }
  :deep(.project-dialog .el-dialog) { width: 95% !important; }
}
</style>
