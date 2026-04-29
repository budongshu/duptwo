<template>
  <div class="project-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('project.list.title') }}</h1>
        <span class="page-subtitle">{{ t('project.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <!-- 视图切换器 -->
        <div class="view-switcher">
          <div class="view-switcher-card">
            <button
              v-for="(tab, index) in viewTabs"
              :key="tab.key"
              class="view-switcher-btn"
              :class="{ active: viewMode === tab.key }"
              :style="{ '--delay': index * 0.05 + 's' }"
              @click="switchView(tab.key)"
            >
              <span class="view-btn-bg"></span>
              <el-icon class="view-btn-icon"><component :is="tab.icon" /></el-icon>
              <span class="view-btn-label">{{ tab.label }}</span>
              <span class="view-btn-indicator"></span>
            </button>
          </div>
        </div>
        <el-button type="primary" @click="handleCreate" class="add-btn">
          <el-icon><Plus /></el-icon>
          {{ t('project.list.addProject') }}
        </el-button>
      </div>
    </header>

    <!-- KPI 统计 -->
    <div class="kpi-cards" v-if="totalProjectCount > 0 && viewMode !== 'network'">
      <div class="kpi-card" :class="{ 'kpi-card--active': !searchKeyword && !searchStatus && !searchStage }" @click="handleReset">
        <div class="kpi-num">{{ totalProjectCount }}</div>
        <div class="kpi-label">总项目</div>
      </div>
      <div class="kpi-card kpi-card--running" :class="{ 'kpi-card--active': searchStage === 'running' }" @click="filterByStage('running')">
        <div class="kpi-num">{{ stageCount.running }}</div>
        <div class="kpi-label">运营中</div>
      </div>
      <div class="kpi-card kpi-card--deploying" :class="{ 'kpi-card--active': searchStage === 'deploying' }" @click="filterByStage('deploying')">
        <div class="kpi-num">{{ stageCount.deploying }}</div>
        <div class="kpi-label">部署中</div>
      </div>
      <div class="kpi-card kpi-card--planning" :class="{ 'kpi-card--active': searchStage === 'planning' }" @click="filterByStage('planning')">
        <div class="kpi-num">{{ stageCount.planning }}</div>
        <div class="kpi-label">待定中</div>
      </div>
      <div class="kpi-card kpi-card--designing" :class="{ 'kpi-card--active': searchStage === 'designing' }" @click="filterByStage('designing')">
        <div class="kpi-num">{{ stageCount.designing }}</div>
        <div class="kpi-label">方案中</div>
      </div>
      <div class="kpi-card kpi-card--data" @click="jumpToUpload">
        <div class="kpi-num">{{ formatBytes(totalUploadSize) }}</div>
        <div class="kpi-label">总数据大小</div>
      </div>
      <div class="kpi-card kpi-card--records" @click="jumpToUpload">
        <div class="kpi-num">{{ totalUploadCount }}</div>
        <div class="kpi-label">总上传条数</div>
      </div>
    </div>

    <!-- 阶段分布迷你图 -->
    <div class="stage-distribution" v-if="!loading && totalProjectCount > 0 && viewMode !== 'network'">
      <div class="stage-dist-header">
        <div class="dist-header-left">
          <span class="stage-dist-title">人员分布</span>
          <div class="dist-stats">
            <span class="dist-stat">
              <span class="dist-stat-num">{{ workLoadStats.personCount }}</span>
              <span class="dist-stat-label">人</span>
            </span>
            <span class="dist-sep">/</span>
            <span class="dist-stat">
              <span class="dist-stat-num">{{ workLoadStats.projectCount }}</span>
              <span class="dist-stat-label">项目</span>
            </span>
          </div>
        </div>
        <div class="dist-header-right">
          <div class="workload-indicator" :class="'workload--' + workLoadStats.status">
            <span class="workload-bar">
              <span class="workload-fill" :style="{ width: Math.min(workLoadStats.laborRate, 100) + '%' }"></span>
            </span>
            <span class="workload-text">{{ workLoadStats.avgProjectsPerPerson }} 人均</span>
          </div>
          <el-switch v-model="showActiveOnly" size="small" active-text="活跃" inactive-text="全部" />
        </div>
      </div>
      <div class="stage-dist-chart">
        <div class="stage-dist-bar">
          <div
            v-for="col in kanbanColumns"
            :key="col.key"
            class="stage-dist-item"
            :style="{ flex: Math.max(getStageCount(col.key), 0.5) }"
          >
            <div class="stage-dist-fill" :style="{ background: col.color, height: Math.max(getStageHeight(col.key), 5) + '%' }"></div>
            <span class="stage-dist-count">{{ getStageCount(col.key) }}</span>
          </div>
        </div>
        <div class="stage-dist-legend">
          <div v-for="col in kanbanColumns" :key="col.key" class="legend-item">
            <span class="legend-dot" :style="{ background: col.color }"></span>
            <span class="legend-label">{{ t('project.stage.' + col.key) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar" v-if="viewMode !== 'network'">
      <el-input
        v-model="searchKeyword"
        :placeholder="t('project.list.searchPlaceholder')"
        style="width: 240px"
        clearable
        @input="debouncedSearch"
        @keyup.enter="handleSearch"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-select v-model="searchStage" :placeholder="t('project.list.stage')" clearable style="width: 140px" @change="debouncedSearch">
        <el-option v-for="s in stageOptions" :key="s.value" :label="t('project.stage.' + s.value)" :value="s.value" />
      </el-select>
      <el-select v-model="searchStatus" :placeholder="t('project.form.statusLabel')" clearable style="width: 140px" @change="debouncedSearch">
        <el-option :label="t('common.enabled')" value="active" />
        <el-option :label="t('common.disabled')" value="inactive" />
      </el-select>
      <el-button @click="handleReset" text>重置</el-button>
      <!-- 拼图风格选择 -->
      <div class="masonry-style-picker" v-if="viewMode === 'detail'">
        <button
          v-for="s in masonryStyles"
          :key="s.key"
          class="masonry-style-btn"
          :class="{ active: masonryStyle === s.key }"
          :title="s.desc"
          @click="masonryStyle = s.key"
        >
          {{ s.label }}
        </button>
      </div>
      <div class="filter-tip">
        共 {{ pagination.total }} 个项目
      </div>
    </div>

    <!-- 看板视图 -->
    <div v-show="viewMode === 'kanban' && !loading" class="kanban-board">
      <div
        v-for="col in kanbanColumns"
        :key="col.key"
        class="kanban-col"
        :class="{
          'kanban-col--collapsed': collapsedStages.includes(col.key),
          'kanban-col--dragover': dragOverCol === col.key
        }"
        @dragover.prevent="onDragOver(col.key)"
        @dragenter.prevent="onDragEnter(col.key)"
        @dragleave="onDragLeave(col.key)"
        @drop="onDrop(col.key)"
      >
        <div class="kanban-col-head" @click="toggleStageCollapse(col.key)">
          <span class="kanban-dot" :style="{ background: col.color }"></span>
          <span class="kanban-title">{{ t('project.stage.' + col.key) }}</span>
          <span class="kanban-stats">
            <span class="kanban-count">{{ kanbanData[col.key]?.length || 0 }}</span>
            <span class="kanban-records">{{ getStageRecords(col.key) }} {{ t('project.list.records') }}</span>
            <span class="kanban-size">{{ formatBytes(getStageSize(col.key)) }}</span>
          </span>
          <el-icon class="collapse-icon" :class="{ 'is-collapsed': collapsedStages.includes(col.key) }">
            <ArrowDown />
          </el-icon>
        </div>
        <div class="kanban-col-body" v-show="!collapsedStages.includes(col.key)">
          <div
            v-for="p in (kanbanData[col.key] || [])"
            :key="p.id"
            class="kanban-card"
            :class="{ 'kanban-card--dragging': draggingProject?.id === p.id }"
            draggable="true"
            @dragstart="onDragStart(p, col.key)"
            @dragend="onDragEnd"
          >
            <div class="kanban-card-name" :title="p.name">{{ p.name }}</div>
            <div class="kanban-card-metrics">
              <span class="kanban-metric"><el-icon><User /></el-icon>{{ getTeamCount(p) }}</span>
              <span class="kanban-metric"><el-icon><Document /></el-icon>{{ formatNumber(p.recordCount || 0) }}</span>
              <span class="kanban-metric"><el-icon><Folder /></el-icon>{{ formatBytes(p.totalDataSize || 0) }}</span>
            </div>
          </div>
          <div v-if="!kanbanData[col.key]?.length" class="kanban-empty">{{ t('project.list.noData') }}</div>
        </div>
      </div>
    </div>

    <!-- 人员矩阵视图 -->
    <div v-if="viewMode === 'network'" class="matrix-view">
      <!-- 加载状态 -->
      <div v-if="loading" class="matrix-loading">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>
      <!-- 空状态 -->
      <div v-else-if="matrixPersons.length === 0" class="matrix-empty">
        <div class="empty-icon">👥</div>
        <div class="empty-text">暂无人员负荷数据</div>
        <div class="empty-hint">请先在项目中添加人员信息</div>
      </div>
      <!-- 主内容 -->
      <template v-else>
        <!-- 顶部标题栏 -->
        <div class="matrix-header">
          <div class="matrix-header-left">
            <h3 class="matrix-title">人员负荷矩阵</h3>
            <span class="matrix-desc">基于项目数量和重要程度分析人员工作负载</span>
          </div>
          <div class="matrix-header-right">
            <span class="matrix-date">{{ currentDate }}</span>
            <!-- 视图切换 -->
            <div class="matrix-view-toggle">
              <button :class="['toggle-btn', { active: matrixViewType === 'card' }]" @click="switchMatrixView('card')">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
                卡片矩阵
              </button>
              <button :class="['toggle-btn', { active: matrixViewType === 'scatter' }]" @click="switchMatrixView('scatter')">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="7" cy="15" r="2"/><circle cx="12" cy="9" r="2"/><circle cx="18" cy="18" r="2"/><circle cx="5" cy="5" r="2"/><circle cx="17" cy="7" r="2"/></svg>
                散点图
              </button>
            </div>
          </div>
        </div>

        <!-- 统计栏 -->
        <div class="matrix-stats">
          <div class="mstat-card">
            <div class="mstat-icon mstat-icon--blue">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            </div>
            <div class="mstat-info">
              <div class="mstat-value">{{ matrixPersons.length }}</div>
              <div class="mstat-label">人员</div>
            </div>
          </div>
          <div class="mstat-card">
            <div class="mstat-icon mstat-icon--green">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
            </div>
            <div class="mstat-info">
              <div class="mstat-value">{{ matrixStats.totalProjects }}</div>
              <div class="mstat-label">参与项目</div>
            </div>
          </div>
          <div class="mstat-card">
            <div class="mstat-icon mstat-icon--orange">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/></svg>
            </div>
            <div class="mstat-info">
              <div class="mstat-value">{{ matrixStats.avgWorkload }}%</div>
              <div class="mstat-label">平均负荷</div>
            </div>
          </div>
          <div class="mstat-card">
            <div class="mstat-icon mstat-icon--red">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <div class="mstat-info">
              <div class="mstat-value">{{ matrixStats.overloadCount }}</div>
              <div class="mstat-label">超负荷</div>
            </div>
          </div>
          <div class="mstat-card">
            <div class="mstat-icon mstat-icon--purple">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            </div>
            <div class="mstat-info">
              <div class="mstat-value">{{ matrixStats.idleCount }}</div>
              <div class="mstat-label">空闲</div>
            </div>
          </div>
        </div>

        <!-- 卡片矩阵视图 -->
        <div v-show="matrixViewType === 'card'" class="matrix-card-grid">
          <!-- 四象限 -->
          <div class="matrix-quadrant" :class="'quadrant--' + q.key" v-for="q in matrixQuadrants" :key="q.key">
            <div class="quadrant-header">
              <span class="quadrant-dot" :style="{ background: q.color }"></span>
              <span class="quadrant-name">{{ q.name }}</span>
              <span class="quadrant-count">{{ q.persons.length }}人</span>
            </div>
            <div class="quadrant-desc">{{ q.desc }}</div>
            <div class="quadrant-persons" v-if="q.persons.length > 0">
              <div
                v-for="person in q.persons"
                :key="person.name"
                class="person-bubble"
                :style="{ background: person.color }"
                :title="person.name + ' · ' + person.projectCount + '个项目'"
                @mouseenter="hoveredPerson = person"
                @mouseleave="hoveredPerson = null"
              >
                <span class="person-bubble__name">{{ person.name.substring(0, 1) }}</span>
                <div class="person-bubble__tooltip">
                  <div class="tooltip-name">{{ person.name }}</div>
                  <div class="tooltip-meta">{{ person.projectCount }}个项目 · 负荷{{ person.loadPercent }}%</div>
                  <div class="tooltip-bar"><div class="tooltip-bar-fill" :style="{ width: person.loadPercent + '%', background: person.color }"></div></div>
                  <div class="tooltip-role" v-if="person.roles.length">{{ person.roles[0] }}</div>
                  <div class="tooltip-projects" v-if="person.projects.length">
                    <div class="tooltip-proj-title">参与项目</div>
                    <div class="tooltip-proj-item" v-for="p in person.projects" :key="p.name">{{ p.name }}</div>
                  </div>
                </div>
              </div>
            </div>
            <div class="quadrant-empty" v-else>—</div>
          </div>
        </div>

        <!-- 散点图视图 -->
        <div v-show="matrixViewType === 'scatter'" class="matrix-scatter-wrap">
          <div ref="matrixChartRef" class="matrix-scatter-canvas"></div>
        </div>

        <!-- 图例 -->
        <div class="matrix-legend">
          <div class="legend-item" v-for="q in matrixQuadrants" :key="q.key">
            <span class="legend-dot" :style="{ background: q.color }"></span>
            <span class="legend-name">{{ q.name }}</span>
            <span class="legend-desc">{{ q.desc }}</span>
          </div>
        </div>
      </template>
    </div>

    <!-- 项目拼图视图 (Treemap) -->
    <div v-show="viewMode === 'detail' && masonryStyle === 'treemap'" class="treemap-view">
      <div ref="projectTreemapRef" class="treemap-canvas"></div>
    </div>

    <!-- 项目详情视图 -->
    <div v-show="viewMode === 'detail' && masonryStyle !== 'treemap'" class="bento-view">
      <div
        class="bento-grid"
        :style="{
          gridTemplateColumns: gridColumns,
          gap: currentMasonryStyle.gap + 'px'
        }"
      >
        <div
          v-for="(project, idx) in tableData"
          :key="project.id"
          class="bento-card"
          :class="{ 'is-dragging': draggingCard === idx, 'drag-over': dragOverCard === idx }"
          :style="{ animationDelay: `${idx * 0.04}s` }"
          draggable="true"
          @dragstart="onCardDragStart($event, idx)"
          @dragover="onCardDragOver($event, idx)"
          @drop="onCardDrop($event, idx)"
          @dragend="onCardDragEnd"
        >
          <!-- 卡片顶部 — 项目身份区 -->
          <div class="bento-top" :style="{ borderTop: `3px solid ${getStageColor(project.stage)}` }">
            <div class="bento-top-left">
              <div class="bento-avatar" :style="{ background: getProjectColor(project.code) }">
                {{ getInitials(project.name) }}
              </div>
            </div>
            <div class="bento-top-right">
              <div class="bento-name-row">
                <span class="bento-name">{{ project.name }}</span>
                <div class="bento-hover-actions">
                  <el-tooltip content="编辑" placement="top">
                    <el-button size="small" text @click.stop="openEditFromDetail(project)">
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                    <el-button size="small" text class="btn-delete" @click.stop="handleDelete(project)">
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                    </el-button>
                  </el-tooltip>
                </div>
              </div>
              <div class="bento-code" @click.stop="copyCode(project.code)">{{ project.code }}</div>
            </div>
          </div>

          <!-- 卡片主体 -->
          <div class="bento-body">
            <!-- 核心指标 -->
            <div class="bento-metrics">
              <div class="bento-metric">
                <div class="bento-metric-num">{{ formatNumber(project.recordCount || 0) }}</div>
                <div class="bento-metric-label">{{ t('project.list.records') }}</div>
              </div>
              <div class="bento-metric-sep"></div>
              <div class="bento-metric">
                <div class="bento-metric-num">{{ getTeamCount(project) }}</div>
                <div class="bento-metric-label">成员</div>
              </div>
              <div class="bento-metric-sep"></div>
              <div class="bento-metric">
                <div class="bento-metric-num bento-metric-num--size">{{ formatBytes(project.totalDataSize || 0) }}</div>
                <div class="bento-metric-label">{{ t('project.list.size') }}</div>
              </div>
            </div>

            <div class="bento-divider"></div>

            <!-- 项目描述 -->
            <div class="bento-desc-section" v-if="project.description">
              <div class="bento-desc-label">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                项目描述
              </div>
              <div class="bento-desc">{{ project.description }}</div>
            </div>

            <!-- 项目周期 & 公司 -->
            <div class="bento-period-section" v-if="project.projectPeriod || project.companyAddr">
              <div class="bento-period-label">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                项目周期
              </div>
              <div class="bento-period-row">
                <span class="bento-period-item" v-if="project.projectPeriod">{{ project.projectPeriod }}</span>
                <span class="bento-period-item bento-period-item--addr" v-if="project.companyAddr">
                  <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                  {{ project.companyAddr }}
                </span>
              </div>
            </div>

            <!-- 解决方案 -->
            <div class="bento-solution" v-if="project.solution">
              <div class="bento-solution-label">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>
                {{ t('project.form.solutionLabel') }}
                <span v-if="project.solutionPerson" class="solution-person">{{ project.solutionPerson }}</span>
              </div>
              <div class="bento-solution-text">{{ project.solution }}</div>
            </div>

            <!-- 团队成员 -->
            <div class="bento-team" v-if="getTeamMembers(project).length > 0">
              <div class="bento-team-label">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                团队成员
                <span class="bento-team-count">{{ getTeamMembers(project).length }}</span>
              </div>
              <div class="bento-team-list">
                <div
                  v-for="member in getDisplayMembers(project, 6).visible"
                  :key="member"
                  class="bento-member-chip bento-member-chip--clickable"
                  :title="member + ' · ' + getMemberRole(project, member)"
                  @click.stop="handlePersonnelClick(member)"
                >
                  <div class="bento-member-avatar" :style="{ background: getAvatarColor(member) }">
                    {{ member.charAt(0).toUpperCase() }}
                  </div>
                  <span class="bento-member-name">{{ member }}</span>
                </div>
                <div class="bento-member-overflow" v-if="getDisplayMembers(project, 6).overflow > 0">
                  +{{ getDisplayMembers(project, 6).overflow }}
                </div>
              </div>
            </div>

            <!-- 驻场地点 -->
            <div class="bento-stations" v-if="project.onsiteStations && project.onsiteStations.length > 0">
              <div class="bento-stations-label">
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                驻场地点
                <span class="bento-station-count">{{ project.onsiteStations.length }}</span>
              </div>
              <div class="bento-station-list">
                <div
                  v-for="(s, si) in project.onsiteStations"
                  :key="si"
                  class="bento-station-chip"
                >
                  <span class="station-location">{{ s.location || '—' }}</span>
                  <span class="station-person" v-if="s.person">{{ s.person }}</span>
                  <span class="station-phone" v-if="s.phone">{{ s.phone }}</span>
                </div>
              </div>
            </div>

            <!-- 底部 -->
            <div class="bento-footer">
              <span class="bento-stage-tag" :style="{ background: getStageColor(project.stage) + '18', color: getStageColor(project.stage) }">
                {{ t('project.stage.' + project.stage) }}
              </span>
              <span class="bento-time">
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                {{ formatDate(project.updatedAt) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && tableData.length === 0" class="empty-state">
        <div class="empty-icon-wrap">
          <el-icon class="empty-icon" :style="{ color: searchStage ? '#f59e0b' : '#dcdfe6' }">
            <Folder v-if="!searchStage" />
            <Filter v-else />
          </el-icon>
        </div>
        <p class="empty-title">
          <template v-if="searchStage">
            {{ stageDisplayName(searchStage) }} 暂无项目
          </template>
          <template v-else-if="searchKeyword || searchStatus">
            未找到匹配的项目
          </template>
          <template v-else>
            {{ t('project.list.noProjectData') }}
          </template>
        </p>
        <p class="empty-hint">
          <template v-if="searchStage">
            该阶段暂未分配项目，试试查看其他阶段
          </template>
          <template v-else-if="searchKeyword || searchStatus">
            {{ t('project.list.noProjectHint') }}
          </template>
          <template v-else>
            {{ t('project.list.noProjectHint') }}
          </template>
        </p>
        <div class="empty-actions">
          <el-button v-if="searchStage || searchKeyword || searchStatus" type="primary" plain size="default" @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            清除筛选
          </el-button>
          <el-button type="default" size="default" @click="switchView('kanban')">
            <el-icon><Grid /></el-icon>
            看板视图
          </el-button>
        </div>
      </div>
    </div>

    <!-- 加载骨架 -->
    <div v-if="loading" class="project-grid">
      <div v-for="i in 8" :key="i" class="skeleton-card"></div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="pagination.total > 0 && viewMode !== 'network'">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[12, 24, 48]"
        layout="sizes, prev, pager, next"
        background
      />
    </div>

    <!-- 新增/编辑抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      direction="rtl"
      size="560px"
      :with-header="false"
      :destroy-on-close="true"
      class="project-drawer"
    >
      <div class="drawer-head">
        <div class="drawer-avatar" :style="{ background: getProjectColor(form.code) }">
          {{ dialogPreviewInitials }}
        </div>
        <div class="drawer-title">
          <span class="drawer-mode" :class="isEdit ? 'mode--edit' : 'mode--new'">
            {{ isEdit ? t('common.edit') : t('project.form.newProject') }}
          </span>
          <span class="drawer-name">{{ isEdit ? form.name || t('project.form.project') : t('project.form.project') }}</span>
        </div>
        <el-button text @click="drawerVisible = false">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </el-button>
      </div>

      <div class="drawer-body">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <!-- 基本信息 -->
          <div class="form-section">
            <div class="form-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              {{ t('project.form.basicInfo') }}
            </div>
            <div class="form-row-2">
              <el-form-item :label="t('project.form.nameLabel')" prop="name" class="flex-1">
                <el-input v-model="form.name" :placeholder="t('project.form.namePlaceholder')" @input="dialogPreviewInitials = getInitials(form.name)" />
              </el-form-item>
              <el-form-item :label="t('project.form.codeLabel')" prop="code" class="flex-1">
                <el-input v-model="form.code" :placeholder="t('project.form.codePlaceholder')" />
              </el-form-item>
            </div>
            <div class="form-row-2">
              <el-form-item :label="t('project.form.stageLabel')" class="flex-1">
                <el-select v-model="form.stage" style="width: 100%">
                  <el-option v-for="s in stageOptions" :key="s.value" :label="t('project.stage.' + s.value)" :value="s.value" />
                </el-select>
              </el-form-item>
              <el-form-item :label="t('project.form.statusLabel')" class="flex-1">
                <el-select v-model="form.status" style="width: 100%">
                  <el-option :label="t('common.enabled')" value="active" />
                  <el-option :label="t('common.disabled')" value="inactive" />
                </el-select>
              </el-form-item>
            </div>
            <el-form-item :label="t('project.form.descriptionLabel')">
              <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('project.form.descriptionPlaceholder')" />
            </el-form-item>
          </div>

          <!-- 项目信息 -->
          <div class="form-section">
            <div class="form-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              {{ t('project.form.projectInfo') }}
            </div>
            <el-form-item :label="t('project.form.companyAddr')">
              <el-input v-model="form.companyAddr" :placeholder="t('project.form.companyAddrPlaceholder')" />
            </el-form-item>
            <div class="form-row-2">
              <el-form-item label="开始日期" class="flex-1">
                <el-date-picker
                  v-model="projectStartDate"
                  type="date"
                  placeholder="选择开始日期"
                  value-format="YYYY-MM-DD"
                  format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="结束日期" class="flex-1">
                <el-date-picker
                  v-model="projectEndDate"
                  type="date"
                  placeholder="选择结束日期"
                  value-format="YYYY-MM-DD"
                  format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
            </div>
          </div>

          <!-- 团队成员 -->
          <div class="form-section">
            <div class="form-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              {{ t('project.form.teamMembers') }}
            </div>
            <div class="member-list">
              <div v-for="(member, idx) in teamMembers" :key="idx" class="member-row">
                <span class="member-num">{{ idx + 1 }}</span>
                <div class="member-fields">
                  <el-select v-model="member.role" placeholder="岗位" filterable style="width: 140px">
                    <el-option v-for="role in availableRoles" :key="role" :label="role" :value="role" />
                  </el-select>
                  <el-select v-model="member.members" placeholder="人员" multiple filterable collapse-tags style="flex: 1">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </div>
                <el-button text type="danger" @click="removeTeamMember(idx)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                </el-button>
              </div>
            </div>
            <el-button @click="addTeamMember" type="default" class="add-btn-inline">
              <el-icon><Plus /></el-icon> 添加成员
            </el-button>
          </div>

          <!-- 驻场地点 -->
          <div class="form-section">
            <div class="form-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              {{ t('project.list.onsiteStations') }}
              <el-badge v-if="form.onsiteStations.length > 0" :value="form.onsiteStations.length" />
            </div>
            <div class="station-list">
              <div v-for="(station, idx) in form.onsiteStations" :key="idx" class="station-row">
                <span class="station-num">{{ idx + 1 }}</span>
                <div class="station-fields">
                  <el-select
                    v-model="station.person"
                    placeholder="人员"
                    filterable
                    clearable
                    :loading="personnelLoading"
                    style="flex: 1"
                    @change="(val) => fillStationFromPersonnel(station, val)"
                  >
                    <el-option
                      v-for="p in personnelList"
                      :key="p.id"
                      :label="p.name"
                      :value="p.name"
                    >
                      <div class="person-opt">
                        <span>{{ p.name }}</span>
                        <small v-if="p.company">{{ p.company }}</small>
                      </div>
                    </el-option>
                  </el-select>
                  <el-input v-model="station.location" placeholder="地点" size="default" style="flex: 1" />
                  <el-input v-model="station.phone" placeholder="联系方式" size="default" style="flex: 1" />
                </div>
                <el-button text type="danger" @click="removeStation(idx)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                </el-button>
              </div>
            </div>
            <el-button @click="addStation" type="default" class="add-btn-inline">
              <el-icon><Plus /></el-icon> {{ t('project.form.addStation') }}
            </el-button>
          </div>
        </el-form>
      </div>

      <div class="drawer-foot">
        <el-button @click="drawerVisible = false">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ t('common.cancel') }}
        </el-button>
        <el-button @click="confirmSubmit" :loading="submitting">
          <svg v-if="!submitting" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
          {{ isEdit ? t('common.save') : t('common.create') }}
        </el-button>
      </div>
    </el-drawer>

    <!-- 项目详情弹窗 -->
    <el-dialog v-model="detailVisible" :title="currentDetail?.name || '项目详情'" width="720px" :destroy-on-close="true" class="project-detail-dialog">
      <div class="detail-modal" v-if="currentDetail">
        <!-- 顶部信息 -->
        <div class="detail-modal-head">
          <div class="detail-modal-avatar" :style="getProjectBgStyle(currentDetail.code)">
            {{ getInitials(currentDetail.name) }}
          </div>
          <div class="detail-modal-title">
            <h2 class="modal-project-name">{{ currentDetail.name }}</h2>
            <div class="modal-project-meta">
              <el-tag :type="currentDetail.status === 'active' ? 'success' : 'info'" size="small">{{ currentDetail.status === 'active' ? t('common.enabled') : t('common.disabled') }}</el-tag>
              <el-tag type="warning" size="small">{{ t('project.stage.' + currentDetail.stage) }}</el-tag>
              <span class="modal-code">{{ currentDetail.code }}</span>
            </div>
          </div>
          <div class="detail-modal-stats">
            <div class="modal-stat">
              <span class="modal-stat-num">{{ formatNumber(currentDetail.recordCount || 0) }}</span>
              <span class="modal-stat-label">{{ t('project.list.records') }}</span>
            </div>
            <div class="modal-stat">
              <span class="modal-stat-num">{{ formatBytes(currentDetail.totalDataSize || 0) }}</span>
              <span class="modal-stat-label">{{ t('project.list.size') }}</span>
            </div>
          </div>
        </div>

        <!-- 基本信息 -->
        <div class="detail-modal-section" v-if="currentDetail.description">
          <div class="modal-section-title">项目描述</div>
          <div class="modal-desc">{{ currentDetail.description }}</div>
        </div>

        <!-- 项目信息 -->
        <div class="detail-modal-section" v-if="currentDetail.companyAddr || currentDetail.projectPeriod">
          <div class="modal-section-title">项目信息</div>
          <div class="modal-kv-grid">
            <div class="modal-kv-row" v-if="currentDetail.companyAddr">
              <span class="modal-kv-label"><el-icon><Location /></el-icon> {{ t('project.form.companyAddr') }}</span>
              <span class="modal-kv-value">{{ currentDetail.companyAddr }}</span>
            </div>
            <div class="modal-kv-row" v-if="currentDetail.projectPeriod">
              <span class="modal-kv-label"><el-icon><Clock /></el-icon> {{ t('project.form.projectPeriod') }}</span>
              <span class="modal-kv-value">{{ currentDetail.projectPeriod }}</span>
            </div>
          </div>
        </div>

        <!-- 解决方案 -->
        <div class="detail-modal-section" v-if="currentDetail.solution">
          <div class="modal-section-title">{{ t('project.form.solutionLabel') }}</div>
          <div class="modal-desc">{{ currentDetail.solution }}</div>
        </div>

        <!-- 团队成员 -->
        <div class="detail-modal-section" v-if="getTeamMembers(currentDetail).length > 0">
          <div class="modal-section-title">
            团队成员
            <el-badge :value="getTeamMembers(currentDetail).length" />
          </div>
          <div class="modal-team-grid">
            <div class="modal-team-item" v-for="member in getTeamMembers(currentDetail)" :key="member">
              <div class="modal-team-avatar" :style="{ background: getAvatarColor(member) }">
                {{ member.charAt(0).toUpperCase() }}
              </div>
              <span class="modal-team-name">{{ member }}</span>
            </div>
          </div>
        </div>

        <!-- 驻场地点 -->
        <div class="detail-modal-section" v-if="currentDetail.onsiteStations && currentDetail.onsiteStations.length > 0">
          <div class="modal-section-title">
            {{ t('project.list.onsiteStations') }}
            <el-badge :value="currentDetail.onsiteStations.length" />
          </div>
          <div class="modal-stations">
            <div class="modal-station" v-for="(s, si) in currentDetail.onsiteStations" :key="si">
              <div class="modal-station-icon"><el-icon><LocationInformation /></el-icon></div>
              <div class="modal-station-body">
                <div class="modal-station-location">{{ s.location || '—' }}</div>
                <div class="modal-station-meta">
                  <span v-if="s.person"><el-icon><User /></el-icon> {{ s.person }}</span>
                  <span v-if="s.phone"><el-icon><Phone /></el-icon> {{ s.phone }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ t('common.close') }}
        </el-button>
        <el-button @click="openEditFromDetail">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          {{ t('common.edit') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUser } from '@/composables/useUser'
import {
  Grid, Menu, Plus, Search, Cpu, User,
  Folder, Location, Clock, LocationInformation, Tickets, Phone, Timer, Loading, ArrowDown,
  Filter, RefreshLeft
} from '@element-plus/icons-vue'
import { ProjectApi, type Project, type CreateProjectReq, type UpdateProjectReq, type OnSiteStation } from '@/api/project'
import { PersonnelApi, type Personnel } from '@/api/personnel'
import { UploadRecordApi } from '@/api/upload-record'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()
const router = useRouter()
const { hasPermission } = useUser()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Project[]>([])
const allProjectsData = ref<Project[]>([]) // 所有项目数据（不受筛选影响，用于KPI统计）
const totalProjectCount = ref(0) // 全局项目总数
const drawerVisible = ref(false)
const dialogPreviewInitials = ref('?')
const viewMode = ref<'detail' | 'kanban' | 'network'>('detail')
const detailVisible = ref(false)
const currentDetail = ref<Project | null>(null)
const searchKeyword = ref('')
const searchStatus = ref('')
const searchStage = ref('')
const pagination = reactive({ page: 1, pageSize: 24, total: 0 })
const isEdit = ref(false)
const formRef = ref()
const personnelList = ref<Personnel[]>([])
const personnelLoading = ref(false)

// 矩阵图相关
const matrixChartRef = ref<HTMLElement>()
const projectTreemapRef = ref<HTMLElement>()

// 矩阵视图类型
const matrixViewType = ref<'card' | 'scatter'>('card')

// 矩阵数据
interface MatrixPerson {
  name: string
  projectCount: number
  loadPercent: number       // 负荷百分比
  hasHighImportance: boolean // 参与高重要项目
  color: string
  projects: { name: string; stage: string }[]
  roles: string[]
}
const matrixPersons = ref<MatrixPerson[]>([])

// 矩阵统计数据
const matrixStats = computed(() => {
  const persons = matrixPersons.value
  const totalProjects = new Set(persons.flatMap(p => p.projects.map(pj => pj.name))).size
  const avgWorkload = persons.length > 0 ? Math.round(persons.reduce((s, p) => s + p.loadPercent, 0) / persons.length) : 0
  const overloadCount = persons.filter(p => p.loadPercent > 75).length
  const idleCount = persons.filter(p => p.projectCount === 0).length
  return { totalProjects, avgWorkload, overloadCount, idleCount }
})

// 四象限定义
const matrixQuadrants = computed(() => {
  const persons = matrixPersons.value
  return [
    {
      key: 'attention',
      name: '需关注',
      desc: '高负荷 + 高重要项目',
      color: '#ea580c',
      persons: persons.filter(p => p.projectCount > 3 && p.hasHighImportance),
    },
    {
      key: 'assignable',
      name: '可分配',
      desc: '低负荷 + 高重要项目',
      color: '#22c55e',
      persons: persons.filter(p => p.projectCount <= 3 && p.hasHighImportance && p.projectCount > 0),
    },
    {
      key: 'busy',
      name: '繁忙中',
      desc: '高负荷 + 普通项目',
      color: '#dc2626',
      persons: persons.filter(p => p.projectCount > 3 && !p.hasHighImportance),
    },
    {
      key: 'normal',
      name: '正常',
      desc: '低负荷 + 普通项目',
      color: '#9333ea',
      persons: persons.filter(p => p.projectCount <= 3 && p.projectCount > 0 && !p.hasHighImportance),
    },
    {
      key: 'idle',
      name: '空闲',
      desc: '无参与项目',
      color: '#94a3b8',
      persons: persons.filter(p => p.projectCount === 0),
    },
  ]
})

// 悬停人员
const hoveredPerson = ref<MatrixPerson | null>(null)

// 视图切换
const switchMatrixView = (type: 'card' | 'scatter') => {
  matrixViewType.value = type
  if (type === 'scatter') {
    nextTick(() => initMatrixScatterChart())
  }
}

// 头像颜色池
const avatarColors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#00BFA5', '#7C3AED', '#DB2777', '#0ea5e9', '#f97316', '#8b5cf6', '#ec4899']
const getPersonColor = (name: string) => avatarColors[(name || '').charCodeAt(0) % avatarColors.length]

// 计算矩阵数据
const calculateMatrixData = (projects: Project[]) => {
  const personMap = new Map<string, {
    projectCount: number
    projects: { name: string; stage: string }[]
    roles: Set<string>
    hasHighImportance: boolean
  }>()

  const roleFieldMap: [string, string][] = [
    ['projectPerson', '项目负责人'],
    ['opsPerson', '运维负责人'],
    ['opsStaffPerson', '运维人员'],
    ['developerPerson', '开发人员'],
    ['testerPerson', '测试人员'],
    ['businessPerson', '业务人员'],
    ['compliancePerson', '合规人员'],
    ['securityPerson', '安全人员'],
    ['networkPerson', '网络人员'],
    ['solutionPerson', '方案人员'],
  ]

  projects.forEach(project => {
    const isHighImportance = project.stage === 'running' || project.stage === 'deploying'
    roleFieldMap.forEach(([field, roleName]) => {
      const val = (project as any)[field]
      if (val && val.trim()) {
        val.split(/[,，、]/).map((n: string) => n.trim()).filter(Boolean).forEach(personName => {
          if (!personMap.has(personName)) {
            personMap.set(personName, { projectCount: 0, projects: [], roles: new Set(), hasHighImportance: false })
          }
          const info = personMap.get(personName)!
          info.projectCount++
          info.projects.push({ name: project.name, stage: project.stage || '' })
          info.roles.add(roleName)
          if (isHighImportance) info.hasHighImportance = true
        })
      }
    })
  })

  matrixPersons.value = Array.from(personMap.entries()).map(([name, info]) => ({
    name,
    projectCount: info.projectCount,
    loadPercent: Math.min(Math.round((info.projectCount / 6) * 100), 100),
    hasHighImportance: info.hasHighImportance,
    color: getPersonColor(name),
    projects: info.projects,
    roles: Array.from(info.roles),
  }))
}

// 初始化散点图
const initMatrixScatterChart = async () => {
  if (!matrixChartRef.value) return
  if (!echarts) { echarts = await import('echarts') }
  if (matrixChart) { matrixChart.dispose() }
  matrixChart = echarts.init(matrixChartRef.value)
  matrixChart.resize()

  const persons = matrixPersons.value
  // 坐标：x = 项目数/6*100，y = 高重要?100:0
  const scatterData = persons.map(p => ({
    name: p.name,
    value: [
      p.projectCount > 0 ? Math.min((p.projectCount / 6) * 100, 100) : 0,
      p.hasHighImportance ? 80 : p.projectCount === 0 ? -10 : 30,
      p.projectCount,
      p.projects,
    ],
    itemStyle: { color: p.color, opacity: 0.85, borderColor: '#fff', borderWidth: 2 },
  }))

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      borderWidth: 1,
      borderRadius: 8,
      padding: [10, 14],
      textStyle: { color: '#191c23', fontSize: 12 },
      formatter: (params: any) => {
        const d = params.data
        const [x, y, count, projects] = d.value
        const name = d.name
        const person = persons.find(p => p.name === name)
        const rolesStr = person?.roles.join('、') || ''
        const pjs = (projects as any[]).map((p: any) => `<br/>&nbsp;&nbsp;• ${p.name}`).join('')
        return `<div style="font-weight:600;margin-bottom:4px">${name}</div><div style="color:#666">角色: ${rolesStr}</div><div style="color:#666">项目数: ${count} 项</div><div style="margin-top:4px;color:#666">项目列表: ${pjs}</div>`
      }
    },
    grid: { left: 64, right: 40, top: 24, bottom: 48 },
    xAxis: {
      type: 'value', name: '项目数量',
      nameLocation: 'middle', nameGap: 32,
      nameTextStyle: { color: '#6b7280', fontSize: 12 },
      min: 0, max: 100,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisTick: { show: false },
      axisLabel: { color: '#9ca3af', fontSize: 11, formatter: (v: number) => Math.round(v / 100 * 6) + '' },
      splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } },
    },
    yAxis: {
      type: 'value', name: '重要程度',
      nameLocation: 'middle', nameGap: 44,
      nameTextStyle: { color: '#6b7280', fontSize: 12 },
      min: -20, max: 100,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisTick: { show: false },
      axisLabel: { show: false },
      splitLine: { show: false },
    },
    series: [{
      type: 'scatter',
      symbolSize: (val: number[]) => Math.min(Math.max(val[2] * 8, 20), 56),
      data: scatterData,
      emphasis: {
        scale: 1.25,
        itemStyle: { opacity: 1, shadowBlur: 12, shadowColor: 'rgba(0,0,0,0.18)' }
      },
      label: {
        show: true, position: 'inside',
        formatter: (params: any) => params.data.name.charAt(0),
        color: '#fff', fontSize: 11, fontWeight: 'bold'
      }
    }],
    // 四象限参考线
    markLine: {
      silent: true, symbol: 'none',
      lineStyle: { type: 'solid', color: '#d1d5db', width: 1.5, type: 'dashed' as any },
      label: { show: false },
      data: [{ xAxis: 50 }, { yAxis: 50 }]
    },
    // 象限背景色
    visualMap: { show: false, type: 'continuous' as any }
  }

  matrixChart.setOption(option)
  window.addEventListener('resize', () => matrixChart?.resize())
}

// 加载矩阵视图数据
const loadNetworkData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.getKanbanList()
    const resp = res as any
    const list: any[] = Array.isArray(resp) ? resp : (Array.isArray(resp?.data) ? resp.data : [])
    tableData.value = list
    pagination.total = list.length
    calculateMatrixData(list)
    nextTick(() => {
      if (viewMode.value === 'network' && matrixViewType.value === 'scatter') {
        initMatrixScatterChart()
      }
    })
  } finally { loading.value = false }
}
let matrixChart: any = null
let projectTreemap: any = null
let echarts: any = null

// 矩阵数据接口
interface MatrixPersonData {
  name: string
  projectCount: number  // 参与项目数
  allocationRate: number // 分配率 (0-100)
  importanceScore: number // 重要程度分数 (0-100)
  importanceLevel: string // 重要程度等级: low/middle/high
  color: string // 颜色
  projects: string[] // 参与的项目名称列表
  roles: string[] // 担任的角色列表
}

// 矩阵数据（旧冗余，由 matrixPersons 替代，保留兼容）
const matrixData = ref<MatrixPersonData[]>([])

// 当前日期
const currentDate = new Date().toLocaleDateString('zh-CN', {
  year: 'numeric',
  month: '2-digit',
  day: '2-digit'
}).replace(/\//g, '-')

// 视图切换配置
const viewTabs = [
  { key: 'detail', label: '项目详情', icon: 'Tickets' },
  { key: 'kanban', label: '阶段看板', icon: 'Menu' },
  { key: 'network', label: '人员矩阵', icon: 'Cpu' },
]

// 拼图风格配置
const masonryStyle = ref('standard')
const masonryStyles = [
  { key: 'tight', label: '紧紧', cols: 5, gap: 4, minWidth: 180, desc: '每行5列，适合大量项目快速浏览' },
  { key: 'compact', label: '紧凑', cols: 4, gap: 6, minWidth: 200, desc: '每行4列，平衡信息密度' },
  { key: 'standard', label: '标准', cols: 3, gap: 12, minWidth: 240, desc: '每行3列，平衡美观与信息' },
  { key: 'relaxed', label: '宽松', cols: 2, gap: 20, minWidth: 320, desc: '每行2列，大气舒适展示' },
  { key: 'info', label: '信息', cols: 3, gap: 8, minWidth: 220, desc: '每行3列，展示最全信息' },
  { key: 'treemap', label: '拼图', cols: 0, gap: 0, minWidth: 0, desc: '无缝拼接，矩形块瀑布流' },
]
const currentMasonryStyle = computed(() => masonryStyles.find(s => s.key === masonryStyle.value) || masonryStyles[2])

// 卡片拖拽排序
const draggingCard = ref<number | null>(null)
const dragOverCard = ref<number | null>(null)
const isCardDragging = ref(false)

// 窗口宽度（响应式重算用）
const windowWidth = ref(window.innerWidth)

// 响应式卡片最小宽度（防止14寸屏溢出）
const cardMinWidth = computed(() => {
  const style = currentMasonryStyle.value
  const ww = windowWidth.value
  if (ww < 768) return Math.max(140, style.minWidth - 60)
  if (ww < 1024) return Math.max(160, style.minWidth - 40)
  if (ww < 1280) return Math.max(170, style.minWidth - 20)
  return style.minWidth
})

// 网格列配置（响应式）
const gridColumns = computed(() => {
  const style = currentMasonryStyle.value
  const mw = cardMinWidth.value
  return `repeat(${style.cols}, minmax(${mw}px, 1fr))`
})

const onCardDragStart = (e: DragEvent, idx: number) => {
  draggingCard.value = idx
  isCardDragging.value = true
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', String(idx))
  }
}

const onCardDragOver = (e: DragEvent, idx: number) => {
  e.preventDefault()
  if (draggingCard.value !== null && draggingCard.value !== idx) {
    dragOverCard.value = idx
  }
}

const onCardDrop = (e: DragEvent, idx: number) => {
  e.preventDefault()
  if (draggingCard.value !== null && draggingCard.value !== idx) {
    const data = [...tableData.value]
    const [removed] = data.splice(draggingCard.value, 1)
    data.splice(idx, 0, removed)
    tableData.value = data
  }
  draggingCard.value = null
  dragOverCard.value = null
  isCardDragging.value = false
}

const onCardDragEnd = () => {
  draggingCard.value = null
  dragOverCard.value = null
  isCardDragging.value = false
}

const stageOptions = [
  { value: 'planning' },
  { value: 'designing' },
  { value: 'deploying' },
  { value: 'running' },
  { value: 'paused' },
]

const kanbanColumns = [
  { key: 'running', color: '#4a7c59' },
  { key: 'deploying', color: '#b87333' },
  { key: 'designing', color: '#6b5b95' },
  { key: 'planning', color: '#4a6fa5' },
]

const kanbanData = ref<Record<string, Project[]>>({})

// 看板折叠状态
const collapsedStages = ref<string[]>([])

// 阶段分布开关
const showActiveOnly = ref(false)

// 切换阶段折叠
const toggleStageCollapse = (stageKey: string) => {
  const idx = collapsedStages.value.indexOf(stageKey)
  if (idx > -1) {
    collapsedStages.value.splice(idx, 1)
  } else {
    collapsedStages.value.push(stageKey)
  }
}

// 获取阶段项目数
const getStageCount = (stageKey: string) => {
  return tableData.value.filter(p => p.stage === stageKey && (!showActiveOnly.value || p.status === 'active')).length
}

// 获取阶段记录总数
const getStageRecords = (stageKey: string) => {
  return (kanbanData.value[stageKey] || []).reduce((sum, p) => sum + (p.recordCount || 0), 0)
}

// 获取阶段数据总容量
const getStageSize = (stageKey: string) => {
  return (kanbanData.value[stageKey] || []).reduce((sum, p) => sum + (p.totalDataSize || 0), 0)
}

// 获取阶段条形图高度百分比
const getStageHeight = (stageKey: string) => {
  const counts = kanbanColumns.map(col => getStageCount(col.key))
  const max = Math.max(...counts, 1)
  return Math.round((getStageCount(stageKey) / max) * 100)
}

// 获取阶段颜色
const getStageColor = (stage: string) => {
  const col = kanbanColumns.find(c => c.key === stage)
  return col?.color || '#94a3b8'
}

// 拖拽相关状态
const draggingProject = ref<Project | null>(null)
const draggingFromCol = ref<string>('')
const dragOverCol = ref<string>('')

// 拖拽开始
const onDragStart = (project: Project, colKey: string) => {
  draggingProject.value = project
  draggingFromCol.value = colKey
}

// 拖拽悬停
const onDragOver = (colKey: string) => {
  dragOverCol.value = colKey
}

// 进入目标列
const onDragEnter = (colKey: string) => {
  dragOverCol.value = colKey
}

// 离开目标列
const onDragLeave = (colKey: string) => {
  if (dragOverCol.value === colKey) {
    dragOverCol.value = ''
  }
}

// 放下卡片
const onDrop = async (targetColKey: string) => {
  dragOverCol.value = ''

  if (!draggingProject.value || draggingFromCol.value === targetColKey) {
    draggingProject.value = null
    draggingFromCol.value = ''
    return
  }

  const project = draggingProject.value
  const fromCol = draggingFromCol.value

  // 乐观更新：先移动卡片到新列
  const colData = kanbanData.value[fromCol]
  if (colData) {
    const projectIndex = colData.findIndex(p => p.id === project.id)
    if (projectIndex > -1) {
      colData.splice(projectIndex, 1)
    }
  }

  if (!kanbanData.value[targetColKey]) {
    kanbanData.value[targetColKey] = []
  }
  kanbanData.value[targetColKey].push({ ...project, stage: targetColKey })

  // 调用 API 更新后端
  try {
    await ProjectApi.update({
      id: project.id,
      name: project.name,
      code: project.code,
      description: project.description,
      status: project.status,
      stage: targetColKey,
      sort: project.sort,
      projectPerson: project.projectPerson,
      opsPerson: project.opsPerson,
      developerPerson: project.developerPerson,
      testerPerson: project.testerPerson,
      businessPerson: project.businessPerson,
      compliancePerson: project.compliancePerson,
      opsStaffPerson: project.opsStaffPerson,
      solution: project.solution,
      solutionPerson: project.solutionPerson,
      companyAddr: project.companyAddr,
      projectPeriod: project.projectPeriod,
      onsiteStations: project.onsiteStations,
    })
  } catch (err: any) {
    // 回滚：恢复原位置
    const targetIndex = kanbanData.value[targetColKey]?.findIndex(p => p.id === project.id)
    if (targetIndex !== undefined && targetIndex > -1) {
      kanbanData.value[targetColKey].splice(targetIndex, 1)
    }
    if (!kanbanData.value[fromCol]) {
      kanbanData.value[fromCol] = []
    }
    kanbanData.value[fromCol].push(project)

    ElMessage.error(err?.message || '移动失败')
  }

  draggingProject.value = null
  draggingFromCol.value = ''
}

// 拖拽结束
const onDragEnd = () => {
  draggingProject.value = null
  draggingFromCol.value = ''
  dragOverCol.value = ''
}

const activeCount = computed(() => tableData.value.filter(p => p.status === 'active').length)
const totalRecords = computed(() => tableData.value.reduce((sum, p) => sum + (p.recordCount || 0), 0))
const totalDataSize = computed(() => tableData.value.reduce((sum, p) => sum + (p.totalDataSize || 0), 0))

// 各阶段项目数量（基于所有项目，不受筛选影响）
const stageCount = computed(() => ({
  planning: allProjectsData.value.filter(p => p.stage === 'planning').length,
  designing: allProjectsData.value.filter(p => p.stage === 'designing').length,
  deploying: allProjectsData.value.filter(p => p.stage === 'deploying').length,
  running: allProjectsData.value.filter(p => p.stage === 'running').length,
  paused: allProjectsData.value.filter(p => p.stage === 'paused').length
}))

// 上传记录统计数据
const totalUploadSize = ref(0)
const totalUploadCount = ref(0)

const loadUploadStats = async () => {
  try {
    const stats = await UploadRecordApi.statistics()
    if (stats.data) {
      totalUploadSize.value = stats.data.totalSize || 0
      totalUploadCount.value = stats.data.totalCount || 0
    }
  } catch (e) {
    console.error('Failed to load upload stats:', e)
  }
}

// 跳转到上传记录页面
const jumpToUpload = () => {
  router.push('/upload-record')
}

// 工作负荷统计
const workLoadStats = computed(() => {
  const projectCount = tableData.value.length
  // 收集所有项目相关人员
  const personSet = new Set<string>()
  tableData.value.forEach(p => {
    const fields = ['projectPerson', 'opsPerson', 'opsStaffPerson', 'developerPerson', 'testerPerson', 'businessPerson', 'compliancePerson', 'securityPerson', 'networkPerson', 'solutionPerson']
    fields.forEach(field => {
      const val = (p as any)[field]
      if (val) {
        val.split(/[,，、]/).filter((n: string) => n.trim()).forEach((name: string) => personSet.add(name.trim()))
      }
    })
  })
  const personCount = personSet.size

  // 计算公式：每5个项目 = 1人力
  const requiredLabor = projectCount / 5  // 所需人力
  const laborRate = personCount > 0 ? Math.round((requiredLabor / personCount) * 100) : 0  // 人力充足率
  const avgProjectsPerPerson = personCount > 0 ? (projectCount / personCount).toFixed(1) : '0'  // 人均项目数

  // 判断状态
  let status = 'normal'
  let statusText = '正常'
  if (laborRate > 120) {
    status = 'warning'
    statusText = '偏少'
  } else if (laborRate < 70) {
    status = 'idle'
    statusText = '富余'
  }

  return {
    projectCount,
    personCount,
    requiredLabor: requiredLabor.toFixed(1),
    laborRate,
    avgProjectsPerPerson,
    status,
    statusText
  }
})

// ==================== 协作网络视图 ====================
interface NetworkNode {
  id: string
  name: string
  type: 'person' | 'project'
  x: number
  y: number
  vx: number
  vy: number
  color: string
  linkCount: number
  isActive: boolean
  isDragging: boolean
  links: string[] // 连接的其他节点ID
}

interface NetworkLink {
  source: string
  target: string
  isActive: boolean
}

const networkCanvasRef = ref<HTMLElement>()
const canvasWidth = ref(800)
const canvasHeight = ref(600)
const netPeople = ref<NetworkNode[]>([])
const netProjects = ref<NetworkNode[]>([])
const netLinks = ref<NetworkLink[]>([])
const hoveredNode = ref<NetworkNode | null>(null)
const hoveredRelations = ref<NetworkNode[]>([])
const isAnyActive = ref(false)

const draggingNode = ref<NetworkNode | null>(null)
const dragOffset = ref({ x: 0, y: 0 })
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })
const viewOffset = ref({ x: 0, y: 0 })
const zoom = ref(1)

const hasNetworkData = computed(() => netPeople.value.length > 0 || netProjects.value.length > 0)

const getNodeX = (nodeId: string) => {
  const person = netPeople.value.find(n => n.id === nodeId)
  if (person) return person.x + viewOffset.value.x
  const project = netProjects.value.find(n => n.id === nodeId)
  return project ? project.x + viewOffset.value.x : 0
}

const getNodeY = (nodeId: string) => {
  const person = netPeople.value.find(n => n.id === nodeId)
  if (person) return person.y + viewOffset.value.y
  const project = netProjects.value.find(n => n.id === nodeId)
  return project ? project.y + viewOffset.value.y : 0
}

// 初始化网络数据
const initNetworkData = () => {
  console.log('[Network] tableData count:', tableData.value.length)
  console.log('[Network] tableData sample:', JSON.stringify(tableData.value[0]))
  if (!networkCanvasRef.value) {
    console.log('[Network] canvas ref is null')
    return
  }

  const rect = networkCanvasRef.value.getBoundingClientRect()
  canvasWidth.value = rect.width
  canvasHeight.value = rect.height

  // 收集关系数据
  const personMap = new Map<string, { name: string; projectIds: number[] }>()
  const projectMap = new Map<number, { name: string; personNames: string[] }>()
  const personFields = ['projectPerson', 'opsPerson', 'opsStaffPerson', 'solutionPerson']

  tableData.value.forEach(p => {
    projectMap.set(p.id, { name: p.name, personNames: [] })

    // 从人员字段收集
    personFields.forEach(field => {
      const val = (p as any)[field]
      if (val) {
        val.split(/[,，、]/).filter((n: string) => n.trim()).forEach((name: string) => {
          const trimmedName = name.trim()
          if (!trimmedName) return
          if (!personMap.has(trimmedName)) {
            personMap.set(trimmedName, { name: trimmedName, projectIds: [] })
          }
          if (!personMap.get(trimmedName)!.projectIds.includes(p.id)) {
            personMap.get(trimmedName)!.projectIds.push(p.id)
          }
          if (!projectMap.get(p.id)!.personNames.includes(trimmedName)) {
            projectMap.get(p.id)!.personNames.push(trimmedName)
          }
        })
      }
    })

    // 从驻场地点收集人员
    const stations: any[] = (p as any).onsiteStations || []
    stations.forEach((station: any) => {
      if (station.person) {
        const name = station.person.trim()
        if (!name) return
        if (!personMap.has(name)) {
          personMap.set(name, { name, projectIds: [] })
        }
        if (!personMap.get(name)!.projectIds.includes(p.id)) {
          personMap.get(name)!.projectIds.push(p.id)
        }
        if (!projectMap.get(p.id)!.personNames.includes(name)) {
          projectMap.get(p.id)!.personNames.push(name)
        }
      }
    })
  })

  // 创建人员节点
  const personColors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#c71585', '#00b050', '#005eeb', '#ff6b6b', '#4ecdc4']
  const people = Array.from(personMap.entries()).map(([name, data], idx) => ({
    id: 'p_' + name,
    name,
    type: 'person' as const,
    x: 0,
    y: 0,
    vx: 0,
    vy: 0,
    color: personColors[idx % personColors.length],
    linkCount: data.projectIds.length,
    isActive: false,
    isDragging: false,
    links: data.projectIds.map(id => 'proj_' + id)
  }))

  // 创建项目节点
  const projects = Array.from(projectMap.entries()).map(([id, data], idx) => ({
    id: 'proj_' + id,
    name: data.name,
    type: 'project' as const,
    x: 0,
    y: 0,
    vx: 0,
    vy: 0,
    color: '#f56c6c',
    linkCount: data.personNames.length,
    isActive: false,
    isDragging: false,
    links: data.personNames.map(n => 'p_' + n)
  }))

  // 创建连线
  const links: NetworkLink[] = []
  people.forEach(person => {
    person.links.forEach(targetId => {
      links.push({ source: person.id, target: targetId, isActive: false })
    })
  })

  netPeople.value = people
  netProjects.value = projects
  netLinks.value = links

  // 初始布局 - 环形分布
  setTimeout(() => applyForceLayout(), 50)
}

// 力导向布局
const applyForceLayout = () => {
  const centerX = canvasWidth.value / 2
  const centerY = canvasHeight.value / 2
  const radius = Math.min(canvasWidth.value, canvasHeight.value) * 0.35

  // 初始位置 - 人员在外圈，项目在内圈
  const personCount = netPeople.value.length
  const projectCount = netProjects.value.length

  netPeople.value.forEach((node, i) => {
    const angle = (i / personCount) * Math.PI * 2
    node.x = centerX + radius * 1.2 * Math.cos(angle)
    node.y = centerY + radius * 1.2 * Math.sin(angle)
  })

  netProjects.value.forEach((node, i) => {
    const angle = (i / projectCount) * Math.PI * 2 + Math.PI / projectCount
    node.x = centerX + radius * 0.6 * Math.cos(angle)
    node.y = centerY + radius * 0.6 * Math.sin(angle)
  })

  // 运行力导向模拟
  runForceSimulation()
}

const runForceSimulation = () => {
  const iterations = 100
  const repulsion = 5000
  const attraction = 0.05
  const damping = 0.9

  for (let iter = 0; iter < iterations; iter++) {
    // 计算排斥力（所有节点互相排斥）
    netPeople.value.forEach(p1 => {
      netProjects.value.forEach(p2 => {
        const dx = p1.x - p2.x
        const dy = p1.y - p2.y
        const dist = Math.sqrt(dx * dx + dy * dy) || 1
        const force = repulsion / (dist * dist)
        p1.vx += (dx / dist) * force
        p1.vy += (dy / dist) * force
        p2.vx -= (dx / dist) * force
        p2.vy -= (dy / dist) * force
      })
    })

    netProjects.value.forEach(p1 => {
      netProjects.value.forEach(p2 => {
        if (p1.id === p2.id) return
        const dx = p1.x - p2.x
        const dy = p1.y - p2.y
        const dist = Math.sqrt(dx * dx + dy * dy) || 1
        const force = repulsion * 0.3 / (dist * dist)
        p1.vx += (dx / dist) * force
        p1.vy += (dy / dist) * force
        p2.vx -= (dx / dist) * force
        p2.vy -= (dy / dist) * force
      })
    })

    netPeople.value.forEach(p1 => {
      netPeople.value.forEach(p2 => {
        if (p1.id === p2.id) return
        const dx = p1.x - p2.x
        const dy = p1.y - p2.y
        const dist = Math.sqrt(dx * dx + dy * dy) || 1
        const force = repulsion * 0.3 / (dist * dist)
        p1.vx += (dx / dist) * force
        p1.vy += (dy / dist) * force
        p2.vx -= (dx / dist) * force
        p2.vy -= (dy / dist) * force
      })
    })

    // 计算吸引力（连接的节点互相吸引）
    netLinks.value.forEach(link => {
      const source = netPeople.value.find(n => n.id === link.source) || netProjects.value.find(n => n.id === link.source)
      const target = netPeople.value.find(n => n.id === link.target) || netProjects.value.find(n => n.id === link.target)
      if (!source || !target) return

      const dx = target.x - source.x
      const dy = target.y - source.y
      const dist = Math.sqrt(dx * dx + dy * dy) || 1
      const force = dist * attraction

      source.vx += (dx / dist) * force
      source.vy += (dy / dist) * force
      target.vx -= (dx / dist) * force
      target.vy -= (dy / dist) * force
    })

    // 应用速度和阻尼
    netPeople.value.forEach(node => {
      if (node.isDragging) return
      node.x += node.vx
      node.y += node.vy
      node.vx *= damping
      node.vy *= damping
      const margin = 80
      node.x = Math.max(margin, Math.min(canvasWidth.value - margin, node.x))
      node.y = Math.max(margin, Math.min(canvasHeight.value - margin, node.y))
    })
    netProjects.value.forEach(node => {
      if (node.isDragging) return
      node.x += node.vx
      node.y += node.vy
      node.vx *= damping
      node.vy *= damping
      const margin = 80
      node.x = Math.max(margin, Math.min(canvasWidth.value - margin, node.x))
      node.y = Math.max(margin, Math.min(canvasHeight.value - margin, node.y))
    })
  }
}

const resetNetworkLayout = () => {
  viewOffset.value = { x: 0, y: 0 }
  zoom.value = 1
  applyForceLayout()
}

// 节点拖拽
const onNodeMouseDown = (e: MouseEvent, node: NetworkNode, type: string) => {
  draggingNode.value = node
  node.isDragging = true
  const rect = networkCanvasRef.value!.getBoundingClientRect()
  dragOffset.value = {
    x: e.clientX - rect.left - node.x,
    y: e.clientY - rect.top - node.y
  }
}

const onCanvasMouseDown = (e: MouseEvent) => {
  if (e.target === networkCanvasRef.value || (e.target as HTMLElement).classList.contains('network-svg')) {
    isPanning.value = true
    panStart.value = { x: e.clientX - viewOffset.value.x, y: e.clientY - viewOffset.value.y }
  }
}

const onCanvasMouseMove = (e: MouseEvent) => {
  if (draggingNode.value) {
    const rect = networkCanvasRef.value!.getBoundingClientRect()
    draggingNode.value.x = e.clientX - rect.left - dragOffset.value.x
    draggingNode.value.y = e.clientY - rect.top - dragOffset.value.y
  } else if (isPanning.value) {
    viewOffset.value = {
      x: e.clientX - panStart.value.x,
      y: e.clientY - panStart.value.y
    }
  }
}

const onCanvasMouseUp = () => {
  if (draggingNode.value) {
    draggingNode.value.isDragging = false
    draggingNode.value = null
  }
  isPanning.value = false
}

const onCanvasWheel = (e: WheelEvent) => {
  const delta = e.deltaY > 0 ? 0.9 : 1.1
  zoom.value = Math.max(0.5, Math.min(2, zoom.value * delta))
}

// 节点悬停
const onNodeHover = (node: NetworkNode, type: string) => {
  hoveredNode.value = node
  isAnyActive.value = true

  // 高亮相关节点和连线
  node.isActive = true
  node.links.forEach(linkId => {
    const linkedNode = netPeople.value.find(n => n.id === linkId) || netProjects.value.find(n => n.id === linkId)
    if (linkedNode) linkedNode.isActive = true
  })

  netLinks.value.forEach(link => {
    link.isActive = link.source === node.id || link.target === node.id
  })

  // 更新关联关系列表
  hoveredRelations.value = node.links
    .map(linkId => netPeople.value.find(n => n.id === linkId) || netProjects.value.find(n => n.id === linkId))
    .filter((n): n is NetworkNode => n !== undefined)
}

const onNodeLeave = () => {
  hoveredNode.value = null
  isAnyActive.value = false
  netPeople.value.forEach(n => n.isActive = false)
  netProjects.value.forEach(n => n.isActive = false)
  netLinks.value.forEach(l => l.isActive = false)
  hoveredRelations.value = []
}

const onLinkHover = (link: NetworkLink) => {
  link.isActive = true
  const sourceNode = netPeople.value.find(n => n.id === link.source) || netProjects.value.find(n => n.id === link.source)
  const targetNode = netPeople.value.find(n => n.id === link.target) || netProjects.value.find(n => n.id === link.target)
  if (sourceNode) sourceNode.isActive = true
  if (targetNode) targetNode.isActive = true
}

const onLinkLeave = () => {
  netLinks.value.forEach(l => l.isActive = false)
  netPeople.value.forEach(n => n.isActive = false)
  netProjects.value.forEach(n => n.isActive = false)
}

const form = reactive<CreateProjectReq & { id?: number }>({
  name: '', code: '', description: '', status: 'active', stage: 'planning', sort: 0,
  projectPerson: '', opsPerson: '', opsStaffPerson: '', developerPerson: '', testerPerson: '', businessPerson: '',
  compliancePerson: '', securityPerson: '', networkPerson: '',
  solution: '', solutionPerson: '', companyAddr: '', projectPeriod: '', onsiteStations: [],
})

// 团队成员数据结构
interface TeamMemberRow {
  role: string
  members: string[]
}

// 可用岗位列表
const availableRoles = [
  '项目负责人', '运维负责人', '运维人员', '开发人员', '测试人员',
  '业务人员', '合规人员', '安全人员', '网络人员', '方案人员', '产品人员'
]

// 团队成员（动态表格）
const teamMembers = ref<TeamMemberRow[]>([])

// 添加团队成员行
const addTeamMember = () => {
  teamMembers.value.push({ role: '', members: [] })
}

// 删除团队成员行
const removeTeamMember = (index: number) => {
  teamMembers.value.splice(index, 1)
}

// 将 teamMembers 转换为表单字段（保存时调用）
const syncTeamMembersToForm = () => {
  // 重置所有人员字段
  form.projectPerson = ''
  form.opsPerson = ''
  form.opsStaffPerson = ''
  form.developerPerson = ''
  form.testerPerson = ''
  form.businessPerson = ''
  form.compliancePerson = ''
  form.securityPerson = ''
  form.networkPerson = ''
  form.solutionPerson = ''

  // 根据岗位映射填充
  teamMembers.value.forEach(row => {
    if (!row.role || row.members.length === 0) return
    const membersStr = row.members.join('，')
    switch (row.role) {
      case '项目负责人': form.projectPerson = membersStr; break
      case '运维负责人': form.opsPerson = membersStr; break
      case '运维人员': form.opsStaffPerson = membersStr; break
      case '开发人员': form.developerPerson = membersStr; break
      case '测试人员': form.testerPerson = membersStr; break
      case '业务人员': form.businessPerson = membersStr; break
      case '合规人员': form.compliancePerson = membersStr; break
      case '安全人员': form.securityPerson = membersStr; break
      case '网络人员': form.networkPerson = membersStr; break
      case '方案人员': form.solutionPerson = membersStr; break
    }
  })
}

// 将表单字段转换为 teamMembers（加载时调用）
const syncFormToTeamMembers = () => {
  teamMembers.value = []
  const roleFieldMap: [string, string][] = [
    ['项目负责人', 'projectPerson'],
    ['运维负责人', 'opsPerson'],
    ['运维人员', 'opsStaffPerson'],
    ['开发人员', 'developerPerson'],
    ['测试人员', 'testerPerson'],
    ['业务人员', 'businessPerson'],
    ['合规人员', 'compliancePerson'],
    ['安全人员', 'securityPerson'],
    ['网络人员', 'networkPerson'],
    ['方案人员', 'solutionPerson'],
  ]

  roleFieldMap.forEach(([role, field]) => {
    const val = (form as any)[field]
    if (val && val.trim()) {
      const members = val.split(/[,，、]/).map((n: string) => n.trim()).filter((n: string) => n)
      if (members.length > 0) {
        teamMembers.value.push({ role, members })
      }
    }
  })
}

// 项目周期：分别处理开始和结束日期
const projectStartDate = computed({
  get: () => {
    if (!form.projectPeriod) return null
    const parts = form.projectPeriod.split(' ~ ')
    return parts[0] || null
  },
  set: (val: string | null) => {
    const endParts = form.projectPeriod?.split(' ~ ')?.[1] || ''
    form.projectPeriod = val && endParts ? `${val} ~ ${endParts}` : val || ''
  }
})

const projectEndDate = computed({
  get: () => {
    if (!form.projectPeriod) return null
    const parts = form.projectPeriod.split(' ~ ')
    return parts[1] || null
  },
  set: (val: string | null) => {
    const startParts = form.projectPeriod?.split(' ~ ')?.[0] || ''
    form.projectPeriod = startParts && val ? `${startParts} ~ ${val}` : val || ''
  }
})

// 日期快捷选项
const dateShortcuts = [
  { text: '本周', value: () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - start.getDay() * 86400000)
    return [start, end]
  }},
  { text: '本月', value: () => {
    const end = new Date()
    const start = new Date(end.getFullYear(), end.getMonth(), 1)
    return [start, end]
  }},
  { text: '本季度', value: () => {
    const now = new Date()
    const qStart = new Date(now.getFullYear(), Math.floor(now.getMonth() / 3) * 3, 1)
    return [qStart, now]
  }},
  { text: '今年', value: () => {
    const end = new Date()
    const start = new Date(end.getFullYear(), 0, 1)
    return [start, end]
  }},
]

const formRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入项目编号', trigger: 'blur' }],
}

const formatNumber = (num: number): string => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num.toString()
}

const formatBytes = (bytes: number): string => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[Math.min(i, units.length - 1)]
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  const ss = String(d.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}:${ss}`
}

// 项目首字符（单个更美观）
const getInitials = (name: string) => {
  if (!name) return '?'
  const first = name.trim().charAt(0)
  // 中文取第一个字，英文取第一个字母
  if (/[\u4e00-\u9fa5]/.test(first)) return first
  return first.toUpperCase()
}

// 项目渐变色（基于项目名生成，每项目独特）
const projectColorPalette = [
  '#4a6fa5', // indigo blue
  '#6b5b95', // muted purple
  '#4a7c59', // moss green
  '#b87333', // copper
  '#8b6f5b', // taupe
  '#5b7b8c', // slate blue
  '#a67c52', // warm brown
  '#7a8471', // olive gray
  '#6b7c8c', // blue gray
  '#9c7a6b', // clay
  '#5c6b7a', // dark slate
  '#7a6b8c', // dusty purple
]

const getProjectColor = (code: string) => {
  if (!code) return projectColorPalette[0]
  let hash = 0
  for (let i = 0; i < code.length; i++) hash = ((hash << 5) - hash) + code.charCodeAt(i)
  return projectColorPalette[Math.abs(hash) % projectColorPalette.length]
}

const getProjectBgStyle = (code: string) => {
  return { background: getProjectColor(code) }
}

const getTeamMembers = (project: Project) => {
  const members: string[] = []
  const fields: (keyof Project)[] = [
    'projectPerson', 'opsPerson', 'opsStaffPerson', 'developerPerson',
    'testerPerson', 'businessPerson', 'compliancePerson', 'securityPerson',
    'networkPerson', 'solutionPerson'
  ]
  fields.forEach(field => {
    const val = (project as any)[field]
    if (val) {
      // 支持逗号分隔的多个人员
      val.split(/[,，、]/).filter((n: string) => n.trim()).forEach((name: string) => {
        if (!members.includes(name.trim())) {
          members.push(name.trim())
        }
      })
    }
  })
  return members
}

const getTeamCount = (project: Project) => getTeamMembers(project).length

// 返回最多 displayLimit 个成员，及溢出数量
const getDisplayMembers = (project: Project, displayLimit = 5) => {
  const members = getTeamMembers(project)
  return {
    visible: members.slice(0, displayLimit),
    overflow: Math.max(0, members.length - displayLimit),
  }
}

const getMemberRole = (project: Project, member: string) => {
  const roleMap: [keyof Project, string][] = [
    ['projectPerson', '项目负责人'],
    ['opsPerson', '运维负责人'],
    ['opsStaffPerson', '运维人员'],
    ['developerPerson', '开发人员'],
    ['testerPerson', '测试人员'],
    ['businessPerson', '业务人员'],
    ['compliancePerson', '合规人员'],
    ['securityPerson', '安全人员'],
    ['networkPerson', '网络人员'],
    ['solutionPerson', '方案人员'],
  ]
  for (const [key, label] of roleMap) {
    const val = (project as any)[key]
    if (val && val.split(/[,，、]/).map((n: string) => n.trim()).includes(member)) {
      return label
    }
  }
  return ''
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      status: searchStatus.value || undefined,
      stage: searchStage.value || undefined,
    })
    tableData.value = res.data?.items || res.data || []
    pagination.total = res.data?.total || 0
    // 如果是拼图风格，初始化 Treemap
    if (masonryStyle.value === 'treemap') {
      nextTick(() => initProjectTreemap())
    }
  } finally { loading.value = false }
}

// 加载所有项目用于KPI统计（不分页，不受筛选影响）
const loadAllProjectsForStats = async () => {
  try {
    const res = await ProjectApi.getKanbanList()
    const list: any[] = Array.isArray(res) ? res : (Array.isArray(res?.data) ? res.data : [])
    allProjectsData.value = list
    totalProjectCount.value = list.length
  } catch { allProjectsData.value = []; totalProjectCount.value = 0 }
}

const loadKanbanData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.getKanbanList()
    const list: any[] = Array.isArray(res) ? res : (Array.isArray(res?.data) ? res.data : [])
    const data: Record<string, Project[]> = {}
    for (const col of kanbanColumns) data[col.key] = []
    for (const p of list) {
      if (!p || typeof p !== 'object') continue
      const stage = (p as Project).stage || 'planning'
      if (!data[stage]) data[stage] = []
      data[stage].push(p as Project)
    }
    kanbanData.value = data
  } catch { kanbanData.value = {} }
  finally { loading.value = false }
}

const switchView = (mode: string) => {
  viewMode.value = mode as any
  if (mode === 'kanban') {
    loadKanbanData()
    loadAllProjectsForStats()
  }
  else if (mode === 'network') {
    loadNetworkData().then(() => {
      nextTick(() => initNetworkData())
    })
    loadAllProjectsForStats()
  }
  else loadData()
}

// ==================== 项目拼图（Treemap） ====================
const projectTreemapData = computed(() => {
  if (!tableData.value.length) return []
  // 使用平方根压缩数值差距，让小块也能看清
  return tableData.value.map(p => {
    const raw = Math.max(p.recordCount || 1, 1)
    return {
      name: p.name,
      value: Math.pow(raw, 0.55), // 指数 < 1 压缩极差，让小项目也有足够面积
      rawValue: raw,
      id: p.id,
      code: p.code,
      stage: p.stage,
      stageName: stageDisplayName(p.stage),
      recordCount: p.recordCount || 0,
      totalSize: p.totalDataSize || 0,
      teamCount: getTeamCount(p),
      stageColor: getStageColor(p.stage),
    }
  })
})

const initProjectTreemap = async () => {
  if (!projectTreemapRef.value) return
  if (!echarts) echarts = await import('echarts')

  if (projectTreemap) projectTreemap.dispose()
  projectTreemap = echarts.init(projectTreemapRef.value)
  projectTreemap.resize()

  const data = projectTreemapData.value
  if (!data.length) return

  const option = {
    backgroundColor: '#f8f7f5',
    tooltip: {
      trigger: 'item',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      borderWidth: 1,
      textStyle: { color: '#191c23', fontSize: 13, fontFamily: 'Manrope, Inter, sans-serif' },
      padding: [12, 16],
      formatter: (params: any) => {
        const d = params.data
        return `<div style="font-weight:700;margin-bottom:8px;color:#1c1917;font-size:14px">${d.name}</div>
          <div style="color:#6b7280;font-size:12px">编号: ${d.code}</div>
          <div style="color:#6b7280;font-size:12px">阶段: ${d.stageName}</div>
          <div style="color:#6b7280;font-size:12px">记录: ${d.recordCount} 条 | 大小: ${formatBytes(d.totalSize)}</div>
          <div style="color:#6b7280;font-size:12px">成员: ${d.teamCount} 人</div>`
      }
    },
    series: [{
      type: 'treemap',
      width: '100%',
      height: '100%',
      roam: false,
      nodeClick: false,
      breadcrumb: { show: false },
      // 让方块尽量保持正方形，避免细长条看不清
      squareRatio: 1,
      label: {
        show: true,
        fontSize: 14,
        fontWeight: 'bold',
        formatter: (params: any) => {
          const d = params.data
          // 名称最多显示 10 个字（中文），过长自动截断
          const name = d.name.length > 10 ? d.name.substring(0, 9) + '…' : d.name
          const line2 = `{count|${d.recordCount}条 | ${d.stageName || ''}}`
          return `{name|${name}}\n${line2}`
        },
        rich: {
          name: { fontSize: 15, fontWeight: 'bold', color: '#fff', lineHeight: 24 },
          count: { fontSize: 11, color: 'rgba(255,255,255,0.85)', lineHeight: 16 },
        },
      },
      upperLabel: {
        show: true,
        height: 28,
        fontSize: 12,
        fontWeight: 'bold',
        formatter: (params: any) => {
          const d = params.data
          return `{stage|${d.stageName}}  {size|${formatBytes(d.totalSize)}}`
        },
        rich: {
          stage: { fontSize: 12, fontWeight: 'bold', color: 'rgba(255,255,255,0.95)', backgroundColor: 'rgba(0,0,0,0.25)', borderRadius: 4, padding: [2, 8, 2, 8] },
          size: { fontSize: 11, color: 'rgba(255,255,255,0.85)' },
        }
      },
      itemStyle: {
        borderColor: 'rgba(255,255,255,0.25)',
        borderWidth: 4,
        gapWidth: 4,
        borderRadius: 10,
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 24,
          shadowColor: 'rgba(0,0,0,0.3)',
          borderWidth: 5,
          borderColor: '#fff',
        },
        label: {
          rich: {
            name: { fontSize: 16, fontWeight: 'bold' },
            count: { fontSize: 12 },
          }
        }
      },
      levels: [{
        // 最外层容器
        itemStyle: {
          borderColor: '#f8f7f5',
          borderWidth: 8,
          gapWidth: 8,
        },
        upperLabel: { show: false },
        label: { show: false },
      }, {
        // 所有子块
        colorSaturation: [0.35, 0.55],
        itemStyle: {
          borderColor: 'rgba(255,255,255,0.25)',
          borderWidth: 4,
          gapWidth: 4,
          borderRadius: 10,
        },
        label: { show: true },
        upperLabel: { show: true },
      }],
      data: data.map(d => {
        const baseColor = d.stageColor || '#6366f1'
        return {
          name: d.name,
          value: d.value,
          id: d.id,
          code: d.code,
          stageName: d.stageName,
          recordCount: d.recordCount,
          totalSize: d.totalSize,
          teamCount: d.teamCount,
          itemStyle: { color: baseColor },
        }
      }),
    }]
  }

  projectTreemap.setOption(option)

  projectTreemap.on('click', (params: any) => {
    if (params.data?.id) {
      const project = tableData.value.find(p => p.id === params.data.id)
      if (project) openEditFromDetail(project)
    }
  })

  window.addEventListener('resize', () => { projectTreemap?.resize() })
}

const handleSearch = () => { pagination.page = 1; loadData() }
const debouncedSearch = (() => {
  let timer: ReturnType<typeof setTimeout>
  return () => {
    clearTimeout(timer)
    timer = setTimeout(() => { pagination.page = 1; loadData() }, 400)
  }
})()
const handleReset = () => {
  searchKeyword.value = ''
  searchStatus.value = ''
  searchStage.value = ''
  pagination.page = 1
  loadData()
}
const filterByStage = (stage: string) => {
  if (viewMode.value !== 'detail') {
    switchView('detail')
  }
  searchStage.value = searchStage.value === stage ? '' : stage
  pagination.page = 1
  loadData()
}

// 阶段显示名称
const stageDisplayName = (stage: string) => {
  const map: Record<string, string> = {
    running: '运营中',
    deploying: '部署中',
    designing: '方案中',
    planning: '待定中',
    paused: '暂定中'
  }
  return map[stage] || stage
}

// 复制到剪贴板
const copyCode = (text: string) => {
  // 优先使用 Clipboard API
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success('已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制失败')
    })
    return
  }
  // Fallback: 使用 textarea 复制（兼容 HTTP 环境）
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
}

// 快捷操作
const handleQuickAction = (cmd: string, project: Project) => {
  switch (cmd) {
    case 'copy-name':
      copyCode(project.name)
      break
    case 'copy-code':
      copyCode(project.code)
      break
    case 'copy-link':
      copyCode(`${window.location.origin}/projects/${project.code}`)
      break
  }
}

const handleEdit = async (row: Project) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id, name: row.name, code: row.code,
    description: row.description || '', status: row.status, stage: row.stage || 'planning', sort: row.sort || 0,
    projectPerson: row.projectPerson || '', opsPerson: row.opsPerson || '',
    opsStaffPerson: row.opsStaffPerson || '',
    developerPerson: row.developerPerson || '', testerPerson: row.testerPerson || '',
    businessPerson: row.businessPerson || '', compliancePerson: row.compliancePerson || '',
    securityPerson: row.securityPerson || '', networkPerson: row.networkPerson || '',
    solution: row.solution || '', solutionPerson: row.solutionPerson || '',
    companyAddr: row.companyAddr || '', projectPeriod: row.projectPeriod || '',
    onsiteStations: row.onsiteStations ? [...row.onsiteStations.map(s => ({ ...s }))] : [],
  })
  syncFormToTeamMembers() // 同步到团队成员表格
  dialogPreviewInitials.value = getInitials(row.name)
  await loadPersonnelList()
  drawerVisible.value = true
}

// 点击人员姓名，跳转到人员管理详情
const handlePersonnelClick = (personName: string) => {
  if (!hasPermission('personnel:read')) {
    ElMessage.warning('您没有查看人员详情的权限')
    return
  }
  router.push({ name: 'PersonnelList', query: { search: personName } })
}

const handleDelete = async (row: Project) => {
  try {
    await ElMessageBox.confirm(`确定删除项目「${row.name}」吗？`, '删除确认', {
      confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning'
    })
    await ProjectApi.del(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    loadData()
    loadKanbanData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e.message || t('common.deleteError')) }
}

const handleCreate = async () => {
  isEdit.value = false
  Object.assign(form, {
    id: undefined, name: '', code: '', description: '', status: 'active', stage: 'planning', sort: 0,
    projectPerson: '', opsPerson: '', opsStaffPerson: '', developerPerson: '', testerPerson: '', businessPerson: '',
    compliancePerson: '', securityPerson: '', networkPerson: '',
    solution: '', solutionPerson: '', companyAddr: '', projectPeriod: '', onsiteStations: [],
  })
  teamMembers.value = [] // 重置团队成员
  dialogPreviewInitials.value = '?'
  await loadPersonnelList()
  drawerVisible.value = true
}

const loadPersonnelList = async () => {
  if (personnelList.value.length > 0) return
  personnelLoading.value = true
  try {
    const res = await PersonnelApi.getAll()
    if (res.code === 200) personnelList.value = res.data || []
  } finally { personnelLoading.value = false }
}

const handlePersonnelSelect = (station: OnSiteStation, name: string) => {
  if (!name) return
  const p = personnelList.value.find(x => x.name === name)
  if (p?.phone) station.phone = p.phone
}

// 从人员列表自动填入驻场地点
const fillStationFromPersonnel = (station: OnSiteStation, name: string) => {
  if (!name) {
    station.location = ''
    station.phone = ''
    return
  }
  const p = personnelList.value.find(x => x.name === name)
  if (p) {
    // 自动填入驻场地点
    if (!station.location && p.location) station.location = p.location
    // 自动填入手机
    if (!station.phone && p.phone) station.phone = p.phone
  }
}

const addStation = () => {
  if (!form.onsiteStations) form.onsiteStations = []
  form.onsiteStations.push({ location: '', person: '', phone: '' })
}
const removeStation = (idx: number) => form.onsiteStations.splice(idx, 1)

const showDetail = (project: Project) => {
  currentDetail.value = project
  detailVisible.value = true
}

const openEditFromDetail = (project?: Project) => {
  const target = project || currentDetail.value
  if (!target) return
  detailVisible.value = false
  handleEdit(target)
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    syncTeamMembersToForm() // 同步团队成员到表单字段
    if (isEdit.value) {
      await ProjectApi.update(form as UpdateProjectReq)
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await ProjectApi.create(form as CreateProjectReq)
      ElMessage.success(t('common.createSuccess'))
    }
    drawerVisible.value = false
    loadData()
    loadKanbanData()
  } catch (e: any) { ElMessage.error(e.message || t('common.operationError')) }
  finally { submitting.value = false }
}

watch(() => pagination.page, () => loadData())
watch(() => pagination.pageSize, () => { pagination.page = 1; loadData() })
watch(masonryStyle, (val) => {
  if (val === 'treemap' && tableData.value.length > 0) {
    nextTick(() => initProjectTreemap())
  }
})
onMounted(() => {
  loadData()
  loadAllProjectsForStats()
  loadUploadStats()
  window.addEventListener('resize', () => { windowWidth.value = window.innerWidth })
})
</script>

<script lang="ts">
export default { name: 'ProjectList' }
</script>

<style scoped lang="scss">
$primary: #409eff;
$success: #67c23a;
$warning: #e6a23c;
$danger: #f56c6c;

/* 视图切换过渡效果 */
.bento-view,
.kanban-board,
.network-view,
.treemap-view {
  transition: opacity 0.2s ease;
}

/* 视图显示动画 */
.bento-view,
.kanban-board,
.network-view {
  animation: fadeIn 0.25s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.project-page {
  padding: 24px;
  min-height: 100vh;
  background: #f5f7fa;
}

/* 页面标题 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.page-subtitle {
  font-size: 13px;
  color: #909399;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

/* 视图切换器 */
.view-switcher {
  display: flex;
  align-items: center;
}

.view-switcher-card {
  display: flex;
  background: #f1f5f9;
  padding: 4px;
  border-radius: 12px;
  gap: 2px;
}

.view-switcher-btn {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border: none;
  background: transparent;
  font-size: 13px;
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  border-radius: 10px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;

  .view-btn-icon {
    font-size: 16px;
    transition: transform 0.3s ease;
  }

  .view-btn-label {
    position: relative;
    z-index: 1;
  }

  .view-btn-bg {
    position: absolute;
    inset: 0;
    background: #fff;
    border-radius: 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    opacity: 0;
    transform: scale(0.8);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .view-btn-indicator {
    position: absolute;
    bottom: 4px;
    left: 50%;
    transform: translateX(-50%) scaleX(0);
    width: 20px;
    height: 3px;
    background: linear-gradient(90deg, #3b82f6, #6366f1);
    border-radius: 2px;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  &:hover:not(.active) {
    color: #1e293b;

    .view-btn-icon {
      transform: scale(1.1);
    }
  }

  &.active {
    color: #3b82f6;

    .view-btn-bg {
      opacity: 1;
      transform: scale(1);
      box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
    }

    .view-btn-icon {
      color: #3b82f6;
      transform: scale(1.1);
    }

    .view-btn-indicator {
      transform: translateX(-50%) scaleX(1);
    }
  }
}

/* 添加按钮 */
.add-btn {
  border-radius: 10px !important;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
    transform: translateY(-1px);
  }
}

/* 阶段分布迷你图 - 日式美学 */
.stage-distribution {
  background: #fafaf9;
  border: 1px solid #e8e5e1;
  border-radius: 14px;
  padding: 14px 18px;
  margin-bottom: 16px;
}

.stage-dist-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.dist-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stage-dist-title {
  font-size: 13px;
  font-weight: 600;
  color: #44403c;
}

.dist-stats {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.dist-stat {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.dist-stat-num {
  font-size: 16px;
  font-weight: 700;
  color: #1c1917;
}

.dist-stat-label {
  font-size: 11px;
  color: #a8a29e;
}

.dist-sep {
  color: #d4d0c8;
}

.dist-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.workload-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.workload-bar {
  width: 60px;
  height: 6px;
  background: #f5f5f4;
  border-radius: 3px;
  overflow: hidden;
}

.workload-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s;
}

.workload-text {
  font-size: 11px;
  color: #a8a29e;
}

.workload--normal .workload-fill { background: #6b5b95; }
.workload--warning .workload-fill { background: #b87333; }
.workload--idle .workload-fill { background: #4a7c59; }

.stage-dist-chart {
  display: flex;
  gap: 20px;
}

.stage-dist-bar {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  height: 44px;
  flex: 1;
}

.stage-dist-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  min-width: 40px;
}

.stage-dist-fill {
  width: 100%;
  min-height: 4px;
  border-radius: 4px 4px 0 0;
  transition: height 0.3s;
}

.stage-dist-count {
  font-size: 11px;
  font-weight: 600;
  color: #78716c;
}

.stage-dist-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #78716c;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

/* KPI 卡片 - 日式美学 */
.kpi-cards {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 12px;
  margin-bottom: 0;
}

.kpi-card {
  background: #fafaf9;
  border-radius: 14px;
  padding: 16px 14px;
  border: 1px solid #e8e5e1;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 0;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: #d4d0c8;
    border-radius: 14px 14px 0 0;
  }

  &:hover {
    transform: translateY(-3px);
    border-color: #d4d0c8;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  }

  &--active {
    border-color: #a8a29e;
    background: #fff;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);

    &::before { background: #6b5b95; }

    &.kpi-card--running::before { background: #4a7c59; }
    &.kpi-card--deploying::before { background: #b87333; }
    &.kpi-card--planning::before { background: #4a6fa5; }
    &.kpi-card--designing::before { background: #6b5b95; }
    &.kpi-card--paused::before { background: #a8a29e; }
    &.kpi-card--data::before { background: #6b5b95; }
    &.kpi-card--records::before { background: #6b5b95; }
  }

  &--running::before { background: #4a7c59; }
  &--deploying::before { background: #b87333; }
  &--planning::before { background: #4a6fa5; }
  &--designing::before { background: #6b5b95; }
  &--paused::before { background: #a8a29e; }
  &--data::before { background: #6b5b95; }
  &--records::before { background: #6b5b95; }
}

.kpi-num {
  font-size: 22px;
  font-weight: 800;
  color: #1c1917;
  line-height: 1;
  margin-bottom: 6px;
  letter-spacing: -0.5px;
  font-family: 'Noto Sans SC', 'Manrope', sans-serif;
}

.kpi-label {
  font-size: 11px;
  color: #78716c;
  font-weight: 500;
  letter-spacing: 0.3px;
}


/* ==================== 人员负荷矩阵视图 ==================== */
.matrix-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px 24px;
}

.matrix-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.matrix-view-toggle {
  display: flex;
  gap: 6px;
}

.toggle-btn {
  padding: 6px 16px;
  border-radius: 6px;
  border: 1.5px solid #e2e8f0;
  background: #fff;
  color: #64748b;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;

  &:hover {
    border-color: #94a3b8;
    color: #334155;
  }

  &.active {
    background: #1e293b;
    border-color: #1e293b;
    color: #fff;
  }

  svg { width: 14px; height: 14px; }
}

/* 统计卡片 */
.matrix-stats {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.mstat-card {
  flex: 1;
  background: #fff;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);

  .mstat-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    flex-shrink: 0;
  }

  .mstat-body { flex: 1; min-width: 0; }
  .mstat-num {
    font-size: 20px;
    font-weight: 800;
    color: #1e293b;
    line-height: 1.1;
    font-family: 'Manrope', sans-serif;
  }
  .mstat-label { font-size: 11px; color: #94a3b8; margin-top: 2px; }
}

/* 卡片矩阵 */
.matrix-card-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.matrix-quadrant {
  background: #fff;
  border: 1px solid #f1f5f9;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);

  .quadrant-header {
    padding: 10px 14px;
    border-bottom: 1px solid #f1f5f9;
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #fafafa;

    .quadrant-title {
      font-size: 12px;
      font-weight: 600;
      color: #475569;
      letter-spacing: 0.3px;
    }
    .quadrant-count {
      font-size: 18px;
      font-weight: 800;
      color: #1e293b;
      font-family: 'Manrope', sans-serif;
    }
  }

  .quadrant-body {
    padding: 10px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-height: 60px;
    max-height: 160px;
    overflow-y: auto;

    &::-webkit-scrollbar { width: 4px; }
    &::-webkit-scrollbar-track { background: transparent; }
    &::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 2px; }
  }
}

/* 人员气泡 */
.person-bubble {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  cursor: default;
  transition: all 0.15s;
  font-size: 12px;
  color: #334155;
  position: relative;

  &:hover {
    background: #f1f5f9;
    border-color: #e2e8f0;
    transform: translateX(2px);
  }

  .person-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
    font-weight: 700;
    color: #fff;
    flex-shrink: 0;
    text-shadow: 0 1px 1px rgba(0,0,0,0.1);
  }

  .person-info {
    flex: 1;
    min-width: 0;
    .person-name {
      font-size: 12px;
      font-weight: 600;
      color: #1e293b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .person-meta {
      font-size: 10px;
      color: #94a3b8;
      margin-top: 1px;
    }
  }
}

/* 散点图 */
.matrix-scatter-wrap {
  flex: 1;
  background: #fff;
  border: 1px solid #f1f5f9;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  min-height: 280px;
  position: relative;

  .matrix-scatter-canvas {
    width: 100%;
    height: 100%;
    min-height: 280px;
  }
}

/* 图例 */
.matrix-legend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  flex-shrink: 0;

  .legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: #64748b;

    .legend-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
    }
  }
}

/* 空状态 / 加载 */
.matrix-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
  color: #94a3b8;

  svg { width: 48px; height: 48px; opacity: 0.5; }
  p { font-size: 13px; }
}

.matrix-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
  color: #94a3b8;

  .loading-spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #e2e8f0;
    border-top-color: #94a3b8;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  p { font-size: 13px; }
}

@keyframes spin { to { transform: rotate(360deg); } }


/* 看板 */

/* 节点样式 */
.node-group {
  cursor: pointer;
  transition: transform 0.2s;
}

.node-group:hover {
  transform: scale(1.1);
}

.node-circle {
  transition: fill 0.2s, stroke 0.2s;
  stroke: transparent;
  stroke-width: 3;
}

.node-circle--active {
  stroke: #fff;
  stroke-width: 4;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.2));
}

.node-rect {
  transition: fill 0.2s, stroke 0.2s;
  stroke: transparent;
  stroke-width: 3;
}

.node-rect--active {
  stroke: #fff;
  stroke-width: 4;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.2));
}

.node-glow {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.2; r: attr(r); }
  50% { opacity: 0.35; }
}

/* 缩放控制 */
.zoom-controls {
  position: absolute;
  bottom: 20px;
  left: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  padding: 8px 12px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.zoom-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: #f1f5f9;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  transition: all 0.2s;
}

.zoom-btn:hover {
  background: #e2e8f0;
  color: #334155;
}

.zoom-level {
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  min-width: 40px;
  text-align: center;
}

/* 提示信息 */
.network-hint {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  gap: 16px;
  background: rgba(255, 255, 255, 0.9);
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 11px;
  color: #64748b;
}

/* 侧边栏 */
.network-sidebar {
  width: 280px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.sidebar-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.sidebar-title {
  display: flex;
  flex-direction: column;
}

.sidebar-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.sidebar-type {
  font-size: 11px;
  color: #94a3b8;
}

.sidebar-title-text {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
}

.sidebar-stats {
  display: flex;
  gap: 16px;
  padding: 12px 0;
  border-top: 1px solid #f1f5f9;
  border-bottom: 1px solid #f1f5f9;
}

.sidebar-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.stat-num {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
}

.stat-label {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 2px;
}

.related-list {
  margin-top: 12px;
}

.related-title {
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  margin-bottom: 8px;
}

.related-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px;
  border-radius: 8px;
  transition: background 0.2s;
}

.related-item:hover {
  background: #f8fafc;
}

.related-item--dim {
  opacity: 0.4;
}

.related-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
}

.related-info {
  display: flex;
  flex-direction: column;
}

.related-name {
  font-size: 12px;
  font-weight: 500;
  color: #1e293b;
}

.related-type {
  font-size: 10px;
  color: #94a3b8;
}

.overview-stats {
  display: flex;
  justify-content: space-around;
  padding: 8px 0;
}

.overview-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.overview-num {
  font-size: 20px;
  font-weight: 700;
  color: #409eff;
}

.overview-label {
  font-size: 10px;
  color: #94a3b8;
  margin-top: 2px;
}

.core-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.core-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.core-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  flex-shrink: 0;
}

.core-info {
  flex: 1;
  min-width: 0;
}

.core-name {
  font-size: 12px;
  font-weight: 500;
  color: #1e293b;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.core-count {
  font-size: 10px;
  color: #94a3b8;
}

.core-bar {
  width: 60px;
  height: 4px;
  background: #f1f5f9;
  border-radius: 2px;
  overflow: hidden;
}

.core-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s;
}

/* 旧样式保留兼容性 */
.network-container {
  flex: 1;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 12px;
  position: relative;
  overflow: hidden;
}

.network-lines {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.network-link {
  stroke: #e4e7ed;
  stroke-width: 1.5;
  transition: stroke 0.2s;

  &--highlight {
    stroke: #409eff;
    stroke-width: 2;
  }
}

.network-persons,
.network-projects {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.network-person-node,
.network-project-node {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: grab;
  user-select: none;
  transition: box-shadow 0.2s, border-color 0.2s;

  &:active { cursor: grabbing; }

  &:hover {
    border-color: #409eff;
    box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
  }

  &.network-node--selected {
    border-color: #409eff;
    box-shadow: 0 4px 16px rgba(64, 158, 255, 0.3);
  }
}

.network-project-node {
  border-left: 3px solid #f56c6c;
}

.node-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.node-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #fef0f0;
  color: #f56c6c;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.node-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.node-name {
  font-size: 12px;
  font-weight: 600;
  color: #303133;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-count {
  font-size: 10px;
  color: #909399;
}

.node-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: #fff;

  &--heavy { background: #f56c6c; }
  &--light { background: #67c23a; }
}

/* 工作负荷统计面板 */
.network-stats {
  width: 220px;
  flex-shrink: 0;
}

.stats-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 12px;
  overflow: hidden;
}

.stats-header {
  padding: 12px 14px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.stats-title {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
}

.stats-list {
  padding: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.stats-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
  border-radius: 6px;
  margin-bottom: 4px;

  &:last-child { margin-bottom: 0; }

  &:hover { background: #f5f7fa; }

  &--heavy {
    background: #fef0f0;
    .stats-count { color: #f56c6c; }
  }

  &--light {
    background: #f0f9eb;
    .stats-count { color: #67c23a; }
  }
}

.stats-item-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stats-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 700;
  color: #fff;
}

.stats-name {
  font-size: 12px;
  color: #606266;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stats-item-right {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.stats-count {
  font-size: 14px;
  font-weight: 700;
  color: #303133;
}

.stats-unit {
  font-size: 10px;
  color: #909399;
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #fff;
  border-radius: 8px;
  padding: 12px 16px;
  border: 1px solid #ebeef5;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-tip {
  margin-left: auto;
  font-size: 13px;
  color: #909399;
}

/* 网格 */
.project-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  min-height: 200px;
  position: relative;

  &--loading {
    opacity: 0.6;
  }
}

.grid-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 10;
  color: #409eff;
}

.project-card {
  background: #fff;
  border-radius: 10px;
  border: 1px solid #ebeef5;
  padding: 20px;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;
  animation: fadeIn 0.3s ease both;

  &:hover {
    box-shadow: 0 4px 16px rgba(0,0,0,0.1);
    transform: translateY(-2px);
    border-color: #409eff;
  }
}

.card-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.card-avatar {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  font-family: 'Manrope', sans-serif;
  letter-spacing: -0.5px;
}

.card-info {
  flex: 1;
  min-width: 0;
}

.card-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-code {
  font-size: 12px;
  color: #c0c4cc;
  margin-top: 2px;
}

.card-desc {
  font-size: 13px;
  color: #909399;
  margin-bottom: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-metrics {
  display: flex;
  gap: 16px;
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.metric-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  font-weight: 600;
  color: #606266;

  .el-icon { color: #c0c4cc; }

  small {
    font-size: 11px;
    color: #c0c4cc;
    font-weight: 400;
  }
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-avatar-list {
  display: flex;
  gap: 4px;
}

.avatar-mini {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 700;
  color: #fff;
}

.card-time {
  font-size: 12px;
  color: #c0c4cc;
}

/* 看板 - 日式美学 */
.kanban-board {
  display: flex;
  gap: 14px;
  overflow-x: auto;
  padding-bottom: 16px;
  align-items: flex-start;
}

.kanban-col {
  width: 280px;
  flex-shrink: 0;
  background: #fafaf9;
  border-radius: 14px;
  border: 1px solid #e8e5e1;
  transition: all 0.2s;

  &--collapsed {
    width: 160px;
    .kanban-col-head { cursor: pointer; }
    .kanban-stats, .collapse-icon { display: none; }
  }

  &--dragover {
    border-color: #a8a29e;
    background-color: #f5f5f4;
  }
}

.kanban-col-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  border-bottom: 1px solid #f0ede8;
  cursor: pointer;
  user-select: none;
}

.kanban-dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  flex-shrink: 0;
}

.kanban-title {
  flex: 1;
  font-size: 13px;
  font-weight: 600;
  color: #44403c;
}

.kanban-stats {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
}

.kanban-count {
  background: #f5f5f4;
  padding: 2px 8px;
  border-radius: 6px;
  font-weight: 600;
  color: #78716c;
}

.kanban-records {
  color: #a8a29e;
}

.kanban-size {
  color: #a8a29e;
  font-weight: 500;
  font-size: 11px;
}

.collapse-icon {
  color: #d4d0c8;
  transition: transform 0.2s;

  &.is-collapsed {
    transform: rotate(-90deg);
  }
}

.kanban-col-body {
  padding: 10px 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-content: flex-start;
  min-height: 60px;
  max-height: calc(100vh - 320px);
  overflow-y: auto;
}

.kanban-card {
  width: 140px;
  min-height: 56px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e8e5e1;
  padding: 8px 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  transition: all 0.2s ease;
  cursor: grab;

  &:active { cursor: grabbing; }

  &:hover {
    border-color: #a8a29e;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }

  &--dragging {
    opacity: 0.4;
    border-color: #a8a29e;
  }
}

.kanban-card-name {
  font-size: 12px;
  font-weight: 600;
  color: #1c1917;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.kanban-card-metrics {
  display: flex;
  gap: 6px;
  font-size: 10px;
  color: #78716c;
}

.kanban-metric {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  white-space: nowrap;

  .el-icon {
    font-size: 11px;
    color: #a8a29e;
  }
}

.kanban-empty {
  text-align: center;
  padding: 24px;
  color: #d4d0c8;
  font-size: 13px;
}

/* 空状态 - 日式美学 */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
}

.empty-icon-wrap {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #f5f5f4;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.empty-icon {
  font-size: 36px;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: #78716c;
  margin: 0 0 8px 0;
}

.empty-hint {
  font-size: 13px;
  color: #a8a29e;
  margin: 0 0 20px 0;
}

.empty-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 骨架屏 */
.skeleton-card {
  height: 200px;
  background: #f5f5f4;
  border-radius: 14px;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* 抽屉 */
:deep(.project-drawer) {
  .el-drawer__body {
    padding: 0;
    display: flex;
    flex-direction: column;
    height: 100%;
  }
}

.drawer-head {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 24px;
  border-bottom: 1px solid #f0ede8;
  background: #fafaf9;
  flex-shrink: 0;
}

.drawer-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  font-family: 'Manrope', sans-serif;
  letter-spacing: -0.5px;
}

.drawer-title {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.drawer-mode {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;

  &.mode--edit { color: $primary; }
  &.mode--new { color: $success; }
}

.drawer-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.drawer-body {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

.form-section {
  margin-bottom: 28px;
  padding-bottom: 28px;
  border-bottom: 1px dashed #e8ecef;

  &:last-child {
    border-bottom: none;
    margin-bottom: 0;
  }
}

.form-section-title {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  text-transform: uppercase;
  letter-spacing: 0.3px;

  svg {
    color: $primary;
  }
}

.form-row-2 {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;

  .flex-1 { flex: 1; }
  .flex-2 { flex: 2; }
}

.form-row-3 {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;

  > * { flex: 1; }
}

.empty-hint {
  padding: 20px;
  color: #9ca3af;
  text-align: center;
  font-size: 13px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px dashed #e5e7eb;
}

.station-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

.station-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.station-fields {
  display: flex;
  gap: 10px;
  flex: 1;

  .el-input {
    flex: 1;
  }
}

.station-num {
  width: 28px;
  height: 28px;
  background: $primary;
  color: #fff;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

/* 成员列表 */
.member-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

.member-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.member-fields {
  display: flex;
  gap: 10px;
  flex: 1;

  .el-select {
    flex: 1;
  }
}

.member-num {
  width: 28px;
  height: 28px;
  background: $primary;
  color: #fff;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

/* 统一添加按钮 */
.add-btn-inline {
  width: 100%;
  justify-content: center;
  border-style: dashed !important;
  color: #64748b;

  &:hover {
    color: $primary;
    border-color: $primary;
  }
}

/* 人员下拉选项 */
.person-opt {
  display: flex;
  align-items: center;
  gap: 8px;

  small {
    color: #94a3b8;
    font-size: 11px;
  }
}

.drawer-foot {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #ebeef5;
  background: #fafbfc;
  flex-shrink: 0;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

/* 瀑布视图 - Bento Cards */
.bento-view { width: 100%; }

/* 项目拼图视图 */
.treemap-view {
  width: 100%;
  height: calc(100vh - 200px);
  min-height: 500px;
  padding: 16px 20px;
}

.treemap-canvas {
  width: 100%;
  height: 100%;
}

.bento-grid {
  display: grid;
  align-items: start;
}

.bento-card {
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  animation: fadeSlideIn 0.4s ease both;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.07);
    border-color: #d4cfc8;

    .bento-hover-actions { opacity: 1; }
    .bento-avatar { transform: scale(1.04); }
  }

  &.is-dragging {
    opacity: 0.4;
    transform: scale(0.97);
    border-style: dashed;
    border-color: #9ca3af;
  }

  &.drag-over {
    border-color: #3b82f6;
    border-style: dashed;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
  }
}

/* 拼图风格选择器 */
.masonry-style-picker {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 10px;
  flex-shrink: 0;
}

.masonry-style-btn {
  padding: 4px 12px;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  color: #78716c;
  cursor: pointer;
  transition: all 0.2s;

  &:hover { background: #f5f4f2; color: #1c1917; }
  &.active { background: #1c1917; color: #fff; }
}

/* 卡片顶部 */
.bento-top {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px 14px;
  background: #fff;
  position: relative;
}

.bento-top-left {
  flex-shrink: 0;
}

.bento-top-right {
  flex: 1;
  min-width: 0;
}

.bento-name-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.bento-hover-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;
  flex-shrink: 0;
}

.bento-name {
  font-size: 14px;
  font-weight: 700;
  color: #1c1917;
  line-height: 1.4;
  letter-spacing: -0.2px;
  word-break: break-word;
}

.bento-avatar {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  transition: transform 0.2s ease;
}

.bento-hover-actions {
  display: flex;
  align-items: center;
  gap: 1px;
  opacity: 0;
  transition: opacity 0.2s;

  .el-button {
    color: #a8a29e;
    padding: 4px 6px;

    &:hover { color: #57534e; background: #f5f5f4; }
  }

  .btn-delete:hover { color: #dc2626; }
}

/* 项目名称 */
.bento-name {
  font-size: 14px;
  font-weight: 700;
  color: #1c1917;
  line-height: 1.4;
  letter-spacing: -0.2px;
  word-break: break-word;
}

.bento-code {
  font-size: 11px;
  color: #a8a29e;
  cursor: pointer;
  margin-top: 3px;
  display: inline-block;
  transition: color 0.2s;

  &:hover { color: #78716c; }
}

/* 卡片主体 */
.bento-body {
  padding: 14px 16px 14px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #fafaf8;
}

.bento-divider {
  height: 1px;
  background: #e8e5e1;
  margin: 0 -2px;
}

/* 核心指标 */
.bento-metrics {
  display: flex;
  align-items: center;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 10px;
  padding: 10px 0;
}

.bento-metric {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
}

.bento-metric-num {
  font-size: 16px;
  font-weight: 700;
  color: #1c1917;
  font-variant-numeric: tabular-nums;
  line-height: 1;

  &--size { font-size: 14px; }
}

.bento-metric-label {
  font-size: 10px;
  color: #a8a29e;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.bento-metric-sep {
  width: 1px;
  height: 28px;
  background: #e8e5e1;
}

/* 描述区块 */
.bento-desc-section { }

.bento-desc-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.bento-desc {
  font-size: 12px;
  color: #78716c;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 项目周期区块 */
.bento-period-section { }

.bento-period-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.bento-period-row {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.bento-period-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #57534e;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 6px;
  padding: 3px 10px;

  &--addr { color: #78716c; }
}

/* 解决方案 */
.bento-solution { }

.bento-solution-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.solution-person {
  color: #78716c;
  font-weight: 600;
  font-size: 11px;
  text-transform: none;
  letter-spacing: 0;
}

.bento-solution-text {
  font-size: 12px;
  color: #57534e;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 团队成员 */
.bento-team { }

.bento-team-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #a8a29e;
  margin-bottom: 7px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.bento-team-count {
  background: #e8e5e1;
  color: #78716c;
  border-radius: 10px;
  padding: 0 7px;
  font-size: 10px;
  font-weight: 700;
  line-height: 16px;
}

.bento-team-list {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.bento-member-chip {
  display: flex;
  align-items: center;
  gap: 5px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 20px;
  padding: 3px 10px 3px 3px;
  cursor: default;
  transition: all 0.18s ease;

  &:hover {
    border-color: #d4cfc8;
    background: #fafaf8;
  }

  &--clickable {
    cursor: pointer;

    &:hover {
      border-color: #a8a29e;
      background: #f5f5f4;
    }
  }
}

.bento-member-avatar {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.bento-member-name {
  font-size: 11px;
  color: #44403c;
  font-weight: 500;
  white-space: nowrap;
}

.bento-member-overflow {
  display: flex;
  align-items: center;
  padding: 3px 10px;
  background: #f5f5f4;
  border: 1px dashed #d4cfc8;
  border-radius: 20px;
  font-size: 11px;
  color: #a8a29e;
  font-weight: 600;
}

/* 驻场地点 */
.bento-stations { }

.bento-stations-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #a8a29e;
  margin-bottom: 6px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.bento-station-count {
  background: #e8e5e1;
  color: #78716c;
  border-radius: 10px;
  padding: 0 7px;
  font-size: 10px;
  font-weight: 700;
  line-height: 16px;
}

.bento-station-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bento-station-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 8px;
  padding: 5px 10px;
  font-size: 11px;
}

.station-location { color: #44403c; font-weight: 500; flex: 1; }
.station-person { color: #78716c; padding-left: 6px; border-left: 1px solid #e8e5e1; cursor: pointer; &:hover { color: #409eff; } }
.station-phone { color: #a8a29e; padding-left: 6px; border-left: 1px solid #e8e5e1; font-size: 10px; }

/* 底部 */
.bento-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 2px;
}

.bento-stage-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 12px;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.bento-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: #a8a29e;
}

/* 核心指标 - 日式简约风格 */
.bento-metrics {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 14px;
  padding: 16px 0;
  border: 1px solid #f0ede8;
}

.bento-metric {
  flex: 1;
  text-align: center;
}

.bento-metric-num {
  font-size: 22px;
  font-weight: 800;
  color: #1c1917;
  line-height: 1;
  font-family: 'Noto Sans SC', 'Manrope', sans-serif;
  letter-spacing: -0.5px;
}

.bento-metric-num--size {
  font-size: 16px;
  font-weight: 700;
  color: #78716c;
}

.bento-metric-label {
  font-size: 11px;
  color: #a8a29e;
  margin-top: 6px;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.bento-metric-sep {
  width: 1px;
  height: 36px;
  background: #e8e5e1;
}

/* 描述区块 */
.bento-desc-section {
  margin-bottom: 10px;
}

.bento-desc-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 500;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.3px;
}

.bento-desc {
  font-size: 13px;
  color: #78716c;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 项目周期区块 */
.bento-period-section {
  margin-bottom: 8px;
}

.bento-period-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 500;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.3px;
}

.bento-period-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.bento-period-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #78716c;
  background: #f5f5f4;
  border-radius: 6px;
  padding: 4px 10px;
}

/* 解决方案 */
.bento-solution {
  margin-bottom: 8px;
}

.bento-solution-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 500;
  color: #a8a29e;
  margin-bottom: 4px;
  letter-spacing: 0.3px;
}

.solution-person {
  color: #78716c;
  font-weight: 600;
}

.bento-solution-text {
  font-size: 13px;
  color: #57534e;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 团队成员 */
.bento-team { }

.bento-team-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 500;
  color: #a8a29e;
  margin-bottom: 8px;
  letter-spacing: 0.3px;
}

.bento-team-count {
  background: #f5f5f4;
  color: #78716c;
  border-radius: 10px;
  padding: 0 8px;
  font-size: 10px;
  font-weight: 700;
  line-height: 18px;
}

.bento-team-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.bento-member-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 24px;
  padding: 4px 12px 4px 4px;
  cursor: default;
  transition: all 0.2s ease;

  &:hover {
    background: #fafaf9;
    border-color: #d4d0c8;
  }

  &--clickable {
    cursor: pointer;

    &:hover {
      background: #f5f5f4;
      border-color: #a8a29e;
      transform: translateY(-1px);
    }
  }
}

.bento-member-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.bento-member-name {
  font-size: 12px;
  color: #44403c;
  font-weight: 500;
  white-space: nowrap;
}

.bento-member-overflow {
  display: flex;
  align-items: center;
  padding: 4px 12px;
  background: #f5f5f4;
  border: 1px dashed #d4d0c8;
  border-radius: 24px;
  font-size: 11px;
  color: #a8a29e;
  font-weight: 600;
}

/* 驻场地点 */
.bento-stations { }

.bento-stations-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #a8a29e;
  margin-bottom: 10px;
  font-weight: 500;
  letter-spacing: 0.3px;

  .el-icon { color: #a8a29e; }
}

.bento-station-count {
  background: #f5f5f4;
  color: #78716c;
  border-radius: 10px;
  padding: 0 8px;
  font-size: 10px;
  font-weight: 700;
  line-height: 18px;
}

.bento-station-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bento-station-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  border: 1px solid #e8e5e1;
  border-radius: 10px;
  padding: 8px 12px;
  font-size: 12px;

  .el-icon { color: #a8a29e; font-size: 13px; }
}

.station-location {
  color: #44403c;
  font-weight: 500;
  flex: 1;
}

.station-person {
  color: #78716c;
  font-weight: 600;
  font-size: 11px;
  padding: 2px 8px;
  background: #f5f5f4;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #e8e5e1;
    color: #57534e;
  }
}

.station-phone {
  color: #a8a29e;
  font-size: 11px;
}

/* 底部 */
.bento-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid #f0ede8;
}

.bento-stage-tag {
  font-size: 10px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  letter-spacing: 0.3px;
}

.bento-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #d4d0c8;

  .el-icon { font-size: 12px; }
}

@keyframes fadeSlideIn {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 详情弹窗 */
:deep(.project-detail-dialog) {
  .el-dialog__body { padding: 0; }
}

.detail-modal {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-modal-head {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 24px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-modal-avatar {
  width: 60px;
  height: 60px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.18);
  font-family: 'Manrope', sans-serif;
  letter-spacing: -1px;
}

.detail-modal-title {
  flex: 1;
}

.modal-project-name {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 8px 0;
}

.modal-project-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-code {
  font-size: 12px;
  color: #c0c4cc;
  font-family: monospace;
}

.detail-modal-stats {
  display: flex;
  gap: 20px;
}

.modal-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.modal-stat-num {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.modal-stat-label {
  font-size: 12px;
  color: #909399;
}

.detail-modal-section {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;

  &:last-child { border-bottom: none; }
}

.modal-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.7;
}

.modal-kv-grid {
  display: flex;
  flex-direction: column;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
}

.modal-kv-row {
  display: grid;
  grid-template-columns: 120px 1fr;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
  gap: 8px;
  font-size: 13px;

  &:last-child { border-bottom: none; }
}

.modal-kv-label {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #909399;
  font-size: 12px;
  white-space: nowrap;

  .el-icon { color: #c0c4cc; font-size: 13px; }
}

.modal-kv-value {
  color: #303133;
  font-weight: 450;
  word-break: break-all;
}

.modal-team-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  min-width: 60px;
}

.modal-team-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
}

.modal-team-name {
  font-size: 12px;
  color: #606266;
  text-align: center;
  max-width: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.modal-stations {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.modal-station {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.modal-station-icon {
  width: 32px;
  height: 32px;
  background: #e6f4ff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
  flex-shrink: 0;
}

.modal-station-body { flex: 1; }

.modal-station-location {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.modal-station-meta {
  display: flex;
  gap: 12px;

  span {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: #909399;

    .el-icon { font-size: 12px; }
  }
}

@media (max-width: 1024px) {
  .kpi-cards { grid-template-columns: repeat(4, 1fr); }
  .modal-kv-grid { grid-template-columns: 1fr; }
}

@media (max-width: 768px) {
  .kpi-cards { grid-template-columns: repeat(3, 1fr); }
  .filter-bar { flex-direction: column; align-items: stretch; }
  .filter-tip { margin-left: 0; }
}
</style>
