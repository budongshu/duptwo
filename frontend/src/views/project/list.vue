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
            <div class="kanban-card-inner">
              <div class="kanban-card-content">
                <div class="kanban-card-head">
                  <div class="kanban-avatar" :style="getProjectBgStyle(p.code)">
                    {{ getInitials(p.name) }}
                  </div>
                  <div class="kanban-info">
                    <div class="kanban-name">{{ p.name }}</div>
                    <div class="kanban-code">{{ p.code }}</div>
                  </div>
                </div>
                <div class="kanban-card-foot">
                  <span class="kanban-metric"><el-icon><Document /></el-icon> {{ formatNumber(p.recordCount || 0) }}</span>
                  <span class="kanban-metric"><el-icon><User /></el-icon> {{ getTeamCount(p) }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-if="!kanbanData[col.key]?.length" class="kanban-empty">{{ t('project.list.noData') }}</div>
        </div>
      </div>
    </div>

    <!-- 人员负荷矩阵视图 -->
    <div v-if="viewMode === 'network'" class="matrix-view">
      <!-- 加载状态 -->
      <div v-if="loading" class="matrix-loading">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="matrixData.length === 0" class="matrix-empty">
        <div class="empty-icon">📊</div>
        <div class="empty-text">暂无人员负荷数据</div>
        <div class="empty-hint">请先在项目中添加人员信息</div>
      </div>

      <!-- 矩阵图主区域 -->
      <div v-else class="matrix-main">
        <!-- 顶部标题栏 -->
        <div class="matrix-header">
          <div class="matrix-header-left">
            <h3 class="matrix-title">人员负荷分布矩阵</h3>
            <span class="matrix-desc">基于项目数量和重要程度分析人员工作负载</span>
          </div>
          <div class="matrix-date">{{ currentDate }}</div>
        </div>

        <!-- 统计卡片 -->
        <div class="matrix-stats">
          <div class="stat-card">
            <div class="stat-icon stat-icon--blue">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ matrixData.length }}</div>
              <div class="stat-label">人员总数</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--green">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ matrixStats.totalProjects }}</div>
              <div class="stat-label">项目总数</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--orange">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/></svg>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ matrixStats.avgWorkload }}%</div>
              <div class="stat-label">平均负荷</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--red">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ matrixStats.overloadCount }}</div>
              <div class="stat-label">超负荷人员</div>
            </div>
          </div>
        </div>

        <!-- 矩阵内容区 -->
        <div class="matrix-content">
          <!-- Y轴标签 -->
          <div class="axis-y-container">
            <span class="axis-label-top">高</span>
            <span class="axis-label-text">项目重要程度</span>
            <span class="axis-label-bottom">低</span>
          </div>

          <!-- 四象限网格 -->
          <div class="matrix-grid-wrapper">
            <div class="matrix-grid">
              <!-- 左上：需关注（高负荷 + 高重要项目） -->
              <div class="quadrant quadrant--tl">
                <div class="quadrant-header">
                  <span class="quadrant-icon">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                  </span>
                  <span class="quadrant-title">需关注</span>
                  <span class="quadrant-count">{{ getCriticalPerson().length }}</span>
                </div>
                <div class="quadrant-body">
                  <div
                    v-for="person in getCriticalPerson()"
                    :key="person.name"
                    class="person-node"
                    @mouseenter="hoveredPerson = person"
                    @mouseleave="hoveredPerson = null"
                  >
                    <div class="node-circle" :style="{ borderColor: '#ea580c' }">
                      <span class="node-avatar" :style="{ background: '#ea580c' }">{{ person.name.substring(0, 1) }}</span>
                      <span class="node-badge">{{ person.projectCount }}</span>
                    </div>
                    <span class="node-name">{{ person.name }}</span>
                  </div>
                  <div v-if="getCriticalPerson().length === 0" class="quadrant-empty">暂无</div>
                </div>
              </div>

              <!-- 右上：可分配（低负荷 + 高重要项目） -->
              <div class="quadrant quadrant--tr">
                <div class="quadrant-header">
                  <span class="quadrant-icon">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                  </span>
                  <span class="quadrant-title">可分配</span>
                  <span class="quadrant-count">{{ getAvailablePerson().length }}</span>
                </div>
                <div class="quadrant-body">
                  <div
                    v-for="person in getAvailablePerson()"
                    :key="person.name"
                    class="person-node"
                    @mouseenter="hoveredPerson = person"
                    @mouseleave="hoveredPerson = null"
                  >
                    <div class="node-circle" :style="{ borderColor: '#16a34a' }">
                      <span class="node-avatar" :style="{ background: '#16a34a' }">{{ person.name.substring(0, 1) }}</span>
                      <span class="node-badge">{{ person.projectCount }}</span>
                    </div>
                    <span class="node-name">{{ person.name }}</span>
                  </div>
                  <div v-if="getAvailablePerson().length === 0" class="quadrant-empty">暂无</div>
                </div>
              </div>

              <!-- 左下：繁忙中（高负荷 + 普通项目） -->
              <div class="quadrant quadrant--bl">
                <div class="quadrant-header">
                  <span class="quadrant-icon">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  </span>
                  <span class="quadrant-title">繁忙中</span>
                  <span class="quadrant-count">{{ getBusyPerson().length }}</span>
                </div>
                <div class="quadrant-body">
                  <div
                    v-for="person in getBusyPerson()"
                    :key="person.name"
                    class="person-node"
                    @mouseenter="hoveredPerson = person"
                    @mouseleave="hoveredPerson = null"
                  >
                    <div class="node-circle" :style="{ borderColor: '#dc2626' }">
                      <span class="node-avatar" :style="{ background: '#dc2626' }">{{ person.name.substring(0, 1) }}</span>
                      <span class="node-badge">{{ person.projectCount }}</span>
                    </div>
                    <span class="node-name">{{ person.name }}</span>
                  </div>
                  <div v-if="getBusyPerson().length === 0" class="quadrant-empty">暂无</div>
                </div>
              </div>

              <!-- 右下：正常（低负荷 + 普通项目） -->
              <div class="quadrant quadrant--br">
                <div class="quadrant-header">
                  <span class="quadrant-icon">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                  </span>
                  <span class="quadrant-title">正常</span>
                  <span class="quadrant-count">{{ getNormalPerson().length }}</span>
                </div>
                <div class="quadrant-body">
                  <div
                    v-for="person in getNormalPerson()"
                    :key="person.name"
                    class="person-node"
                    @mouseenter="hoveredPerson = person"
                    @mouseleave="hoveredPerson = null"
                  >
                    <div class="node-circle" :style="{ borderColor: '#9333ea' }">
                      <span class="node-avatar" :style="{ background: '#9333ea' }">{{ person.name.substring(0, 1) }}</span>
                      <span class="node-badge">{{ person.projectCount }}</span>
                    </div>
                    <span class="node-name">{{ person.name }}</span>
                  </div>
                  <div v-if="getNormalPerson().length === 0" class="quadrant-empty">暂无</div>
                </div>
              </div>
            </div>

            <!-- X轴标签 -->
            <div class="axis-x-container">
              <span class="axis-label-left">低</span>
              <div class="axis-line"></div>
              <span class="axis-label-center">工作负载（项目数量）</span>
              <div class="axis-line"></div>
              <span class="axis-label-right">高</span>
            </div>
          </div>
        </div>

        <!-- 图例说明 -->
        <div class="matrix-legend">
          <div class="legend-section">
            <div class="legend-title">矩阵说明</div>
            <div class="legend-desc">基于人员参与项目数量和项目重要程度进行双维度分析</div>
          </div>
          <div class="legend-rules">
            <div class="legend-item">
              <span class="legend-label">X轴计算：</span>
              <span class="legend-text">工作负载 = 参与项目数 / 标准容量(8) × 100%，项目数 > 3 为高负荷</span>
            </div>
            <div class="legend-item">
              <span class="legend-label">Y轴计算：</span>
              <span class="legend-text">重要程度 = Σ(项目阶段权重)，运营中=3分，部署中=2分，其他=0分</span>
            </div>
          </div>
          <div class="legend-quadrants">
            <div class="legend-quadrant">
              <span class="legend-dot" style="background: #ea580c"></span>
              <span class="legend-name">需关注</span>
              <span class="legend-tip">高负荷 + 高重要项目，人手紧张需关注</span>
            </div>
            <div class="legend-quadrant">
              <span class="legend-dot" style="background: #16a34a"></span>
              <span class="legend-name">可分配</span>
              <span class="legend-tip">低负荷 + 高重要项目，Capacity 充足可分配</span>
            </div>
            <div class="legend-quadrant">
              <span class="legend-dot" style="background: #dc2626"></span>
              <span class="legend-name">繁忙中</span>
              <span class="legend-tip">高负荷 + 普通项目，项目较多需留意</span>
            </div>
            <div class="legend-quadrant">
              <span class="legend-dot" style="background: #9333ea"></span>
              <span class="legend-name">正常</span>
              <span class="legend-tip">低负荷 + 普通项目，负荷正常</span>
            </div>
          </div>
        </div>

        <!-- 悬停详情 -->
        <transition name="fade">
          <div v-if="hoveredPerson" class="person-detail">
            <div class="detail-header">
              <div class="detail-avatar" :style="{ background: hoveredPerson.color }">
                {{ hoveredPerson.name.substring(0, 1) }}
              </div>
              <div class="detail-info">
                <div class="detail-name">{{ hoveredPerson.name }}</div>
                <div class="detail-meta">{{ hoveredPerson.projectCount }} 个项目 · 负荷 {{ hoveredPerson.allocationRate }}%</div>
              </div>
            </div>
            <div class="detail-progress">
              <div class="progress-label">
                <span>工作负荷</span>
                <span>{{ hoveredPerson.allocationRate }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: hoveredPerson.allocationRate + '%', background: hoveredPerson.color }"></div>
              </div>
            </div>
            <div class="detail-roles" v-if="hoveredPerson.roles.length">
              <span v-for="role in hoveredPerson.roles" :key="role" class="role-tag">{{ role }}</span>
            </div>
            <div class="detail-projects" v-if="hoveredPerson.projects.length">
              <div class="projects-title">参与项目</div>
              <div v-for="project in hoveredPerson.projects" :key="project" class="project-item">
                <span class="project-dot" :style="{ background: hoveredPerson.color }"></span>
                {{ project }}
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>

    <!-- 瀑布视图 -->
    <div v-show="viewMode === 'detail'" class="bento-view">
      <div class="bento-grid">
        <div
          v-for="(project, idx) in tableData"
          :key="project.id"
          class="bento-card"
          :style="{ animationDelay: `${idx * 0.04}s` }"
        >
          <!-- 左侧阶段色条 -->
          <div class="bento-stage-bar" :style="{ background: getStageColor(project.stage) }"></div>

          <!-- 卡片顶部 -->
          <div class="bento-top">
            <div class="bento-avatar" :style="{ background: getProjectColor(project.code) }">
              {{ getInitials(project.name) }}
            </div>
            <div class="bento-top-info">
              <div class="bento-name">{{ project.name }}</div>
              <div class="bento-code" @click.stop="copyCode(project.code)">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
                {{ project.code }}
              </div>
            </div>
            <!-- 悬停显示的操作按钮 -->
            <div class="bento-hover-actions">
              <el-tooltip content="复制代码" placement="top">
                <el-button size="small" text @click.stop="copyCode(project.code)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
                </el-button>
              </el-tooltip>
              <el-tooltip content="编辑" placement="top">
                <el-button size="small" text @click.stop="openEditFromDetail(project)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button size="small" text class="btn-delete" @click.stop="handleDelete(project)">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                </el-button>
              </el-tooltip>
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

            <!-- 项目描述 -->
            <div class="bento-desc" v-if="project.description">
              {{ project.description }}
            </div>

            <!-- 项目信息 -->
            <div class="bento-info-row" v-if="project.companyAddr || project.projectPeriod">
              <div class="bento-info-item" v-if="project.companyAddr">
                <el-icon><Location /></el-icon>
                <span>{{ project.companyAddr }}</span>
              </div>
              <div class="bento-info-item" v-if="project.projectPeriod">
                <el-icon><Clock /></el-icon>
                <span class="info-label">项目时间</span>
                <span>{{ project.projectPeriod }}</span>
              </div>
            </div>

            <!-- 解决方案 -->
            <div class="bento-solution" v-if="project.solution">
              <div class="bento-solution-label">
                <el-icon><Connection /></el-icon>
                {{ t('project.form.solutionLabel') }}
                <span v-if="project.solutionPerson" class="solution-person">
                  · {{ project.solutionPerson }}
                </span>
              </div>
              <div class="bento-solution-text">{{ project.solution }}</div>
            </div>

            <!-- 团队成员 -->
            <div class="bento-team" v-if="getTeamMembers(project).length > 0">
              <div class="bento-team-label">
                <el-icon><User /></el-icon>
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
                <el-icon><LocationInformation /></el-icon>
                {{ t('project.list.onsiteStations') }}
                <span class="bento-station-count">{{ project.onsiteStations.length }}</span>
              </div>
              <div class="bento-station-list">
                <div
                  v-for="(s, si) in project.onsiteStations"
                  :key="si"
                  class="bento-station-chip"
                >
                  <el-icon><LocationInformation /></el-icon>
                  <span class="station-location">{{ s.location || '—' }}</span>
                  <span class="station-person bento-member-chip--clickable" v-if="s.person" @click.stop="handlePersonnelClick(s.person)">{{ s.person }}</span>
                  <span class="station-phone" v-if="s.phone">{{ s.phone }}</span>
                </div>
              </div>
            </div>

            <!-- 底部时间 -->
            <div class="bento-footer">
              <span class="bento-stage-tag" :style="{ background: getStageColor(project.stage) + '20', color: getStageColor(project.stage) }">
                {{ t('project.stage.' + project.stage) }}
              </span>
              <span class="bento-time">
                <el-icon><Timer /></el-icon>
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
          <div class="modal-section-content">{{ currentDetail.description }}</div>
        </div>

        <!-- 项目信息 -->
        <div class="detail-modal-section" v-if="currentDetail.companyAddr || currentDetail.projectPeriod || currentDetail.solution">
          <div class="modal-section-title">{{ t('project.form.projectInfo') }}</div>
          <div class="modal-info-grid">
            <div class="modal-info-item" v-if="currentDetail.companyAddr">
              <el-icon><Location /></el-icon>
              <span>{{ t('project.form.companyAddr') }}</span>
              <strong>{{ currentDetail.companyAddr }}</strong>
            </div>
            <div class="modal-info-item" v-if="currentDetail.projectPeriod">
              <el-icon><Clock /></el-icon>
              <span>{{ t('project.form.projectPeriod') }}</span>
              <strong>{{ currentDetail.projectPeriod }}</strong>
            </div>
            <div class="modal-info-item modal-info-item--full" v-if="currentDetail.solution">
              <el-icon><Connection /></el-icon>
              <span>{{ t('project.form.solutionLabel') }}</span>
              <strong>{{ currentDetail.solution }}</strong>
            </div>
          </div>
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
  Grid, Menu, Plus, Search, Document, Cpu, User,
  Folder, Location, Clock, Connection, LocationInformation, Tickets, Phone, Timer, Loading, ArrowDown,
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
let matrixChart: any = null
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

// 矩阵数据
const matrixData = ref<MatrixPersonData[]>([])

// 悬停的人员详情
const hoveredPerson = ref<MatrixPersonData | null>(null)

// 当前日期
const currentDate = new Date().toLocaleDateString('zh-CN', {
  year: 'numeric',
  month: '2-digit',
  day: '2-digit'
}).replace(/\//g, '-')

// 矩阵统计数据
const matrixStats = computed(() => {
  const totalProjects = new Set(matrixData.value.flatMap(p => p.projects)).size
  const avgWorkload = matrixData.value.length > 0
    ? Math.round(matrixData.value.reduce((sum, p) => sum + p.allocationRate, 0) / matrixData.value.length)
    : 0
  const overloadCount = matrixData.value.filter(p => p.allocationRate > 75).length
  return { totalProjects, avgWorkload, overloadCount }
})

// 需关注：高负荷 + 高重要项目
const getCriticalPerson = () => {
  return matrixData.value.filter(p => p.projectCount > 3 && p.importanceScore >= 30)
}

// 可分配：低负荷 + 有高重要项目
const getAvailablePerson = () => {
  return matrixData.value.filter(p => p.projectCount <= 3 && p.importanceScore >= 30)
}

// 繁忙中：高负荷 + 普通项目
const getBusyPerson = () => {
  return matrixData.value.filter(p => p.projectCount > 3 && p.importanceScore < 30)
}

// 正常：低负荷 + 普通项目
const getNormalPerson = () => {
  return matrixData.value.filter(p => p.projectCount <= 3 && p.importanceScore < 30)
}

// 按状态统计
const getCountByStatus = (status: string) => {
  switch (status) {
    case 'overload': return matrixData.value.filter(p => p.projectCount > 5).length
    case 'busy': return matrixData.value.filter(p => p.projectCount > 3 && p.projectCount <= 5).length
    case 'normal': return matrixData.value.filter(p => p.projectCount <= 3).length
    default: return 0
  }
}

// 视图切换配置
const viewTabs = [
  { key: 'detail', label: '项目详情', icon: 'Tickets' },
  { key: 'kanban', label: '阶段看板', icon: 'Menu' },
  { key: 'network', label: '人员矩阵', icon: 'Grid' },
]

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
  const projectIndex = kanbanData.value[fromCol]?.findIndex(p => p.id === project.id)
  if (projectIndex !== undefined && projectIndex > -1) {
    kanbanData.value[fromCol].splice(projectIndex, 1)
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
  } catch {}
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
  '业务人员', '合规人员', '安全人员', '网络人员', '方案人员'
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
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
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

const avatarColors = ['#6b5b95', '#4a7c59', '#b87333', '#8b6f5b', '#5b7b8c', '#6b7c8c', '#7a8471', '#a67c52']
const getAvatarColor = (name: string) => {
  if (!name) return avatarColors[0]
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = ((hash << 5) - hash) + name.charCodeAt(i)
  return avatarColors[Math.abs(hash) % avatarColors.length]
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

const switchView = (mode: 'detail' | 'kanban' | 'network') => {
  viewMode.value = mode
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

// 加载矩阵视图数据
const loadNetworkData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.getKanbanList()
    const resp = res as any
    const list: any[] = Array.isArray(resp) ? resp : (Array.isArray(resp?.data) ? resp.data : [])
    tableData.value = list
    pagination.total = list.length

    // 计算矩阵数据
    calculateMatrixData(list)
  } finally { loading.value = false }
}

// 计算人员负荷矩阵数据
const calculateMatrixData = (projects: Project[]) => {
  // 收集所有人员及其参与的项目
  const personMap = new Map<string, {
    projectCount: number
    projects: Set<string>
    highImportance: number  // 参与高重要项目的数量
    roles: Set<string>
  }>()

  // 人员角色映射
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

  // 项目重要程度权重（根据阶段）
  const importanceWeight: Record<string, number> = {
    running: 3,    // 运营中 - 最重要
    deploying: 2,  // 部署中 - 重要
    designing: 1,  // 方案中 - 一般
    planning: 0,   // 待定中 - 低
    paused: 0,     // 暂定中 - 低
  }

  // 遍历所有项目，收集人员信息
  projects.forEach(project => {
    const weight = importanceWeight[project.stage || 'planning'] || 0

    roleFieldMap.forEach(([field, roleName]) => {
      const val = (project as any)[field]
      if (val && val.trim()) {
        // 支持逗号分隔的多个人员
        const persons = val.split(/[,，、]/).map((n: string) => n.trim()).filter((n: string) => n)
        persons.forEach(personName => {
          if (!personMap.has(personName)) {
            personMap.set(personName, {
              projectCount: 0,
              projects: new Set(),
              highImportance: 0,
              roles: new Set(),
            })
          }
          const info = personMap.get(personName)!
          info.projectCount++
          info.projects.add(project.name)
          info.roles.add(roleName)
          if (weight >= 2) {
            info.highImportance += weight
          }
        })
      }
    })
  })

  // 转换为矩阵数据
  const maxProjects = 8 // 标准容量
  matrixData.value = Array.from(personMap.entries()).map(([name, info]) => {
    const projectCount = info.projectCount
    // 分配率 = 项目数 / 标准容量 * 100
    const allocationRate = Math.min((projectCount / maxProjects) * 100, 100)

    // 重要程度分数 = 高重要项目加权分
    const importanceScore = Math.min(info.highImportance * 10, 100)

    // 重要程度等级
    let importanceLevel = 'low'
    if (importanceScore >= 60) {
      importanceLevel = 'high'
    } else if (importanceScore >= 30) {
      importanceLevel = 'middle'
    }

    // 颜色：繁忙程度
    let color = '#22c55e' // 正常 - 绿色
    if (projectCount > 5) {
      color = '#ef4444' // 过载 - 红色
    } else if (projectCount > 3) {
      color = '#f59e0b' // 偏忙 - 黄色
    }

    return {
      name,
      projectCount,
      allocationRate,
      importanceScore,
      importanceLevel,
      color,
      projects: Array.from(info.projects),
      roles: Array.from(info.roles),
    }
  })

  // 更新矩阵图
  nextTick(() => {
    if (viewMode.value === 'network') {
      initMatrixChart()
    }
  })
}

// 初始化矩阵图
const initMatrixChart = async () => {
  if (!matrixChartRef.value) return

  if (!echarts) {
    echarts = await import('echarts')
  }

  if (matrixChart) {
    matrixChart.dispose()
  }
  matrixChart = echarts.init(matrixChartRef.value)

  const data = matrixData.value.map(d => ({
    name: d.name,
    value: [d.allocationRate, d.importanceScore, d.projectCount],
    projects: d.projects,
    roles: d.roles,
    color: d.color,
  }))

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item',
      backgroundColor: '#fff',
      borderColor: '#e0e2ec',
      textStyle: { color: '#191c23', fontSize: 12 },
      formatter: (params: any) => {
        const d = params.data
        const rolesStr = d.roles.join('、')
        const projectsStr = d.projects.map((p: string) => `<br/>&nbsp;&nbsp;• ${p}`).join('')
        return `
          <div style="font-weight: 600; margin-bottom: 4px;">${d.name}</div>
          <div style="color: #666;">角色: ${rolesStr}</div>
          <div style="color: #666;">参与项目: ${d.projectCount} 项</div>
          <div style="color: #666; margin-top: 4px;">项目列表: ${projectsStr}</div>
        `
      }
    },
    grid: {
      left: 60,
      right: 40,
      top: 40,
      bottom: 60,
    },
    xAxis: {
      type: 'value',
      name: '分配率 (%)',
      nameLocation: 'middle',
      nameGap: 30,
      nameTextStyle: { color: '#6b7280', fontSize: 12 },
      min: 0,
      max: 100,
      splitNumber: 5,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisTick: { show: false },
      axisLabel: { color: '#9ca3af', fontSize: 11, formatter: '{value}%' },
      splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } },
    },
    yAxis: {
      type: 'value',
      name: '重要程度',
      nameLocation: 'middle',
      nameGap: 40,
      nameTextStyle: { color: '#6b7280', fontSize: 12 },
      min: 0,
      max: 100,
      splitNumber: 5,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisTick: { show: false },
      axisLabel: { color: '#9ca3af', fontSize: 11 },
      splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } },
    },
    series: [{
      type: 'scatter',
      symbolSize: (val: number[]) => {
        // 根据项目数量调整圆圈大小 (最小20, 最大60)
        return Math.min(Math.max(val[2] * 8, 20), 60)
      },
      data,
      itemStyle: {
        color: (params: any) => params.data.color,
        opacity: 0.85,
        borderColor: '#fff',
        borderWidth: 2,
      },
      emphasis: {
        scale: 1.3,
        itemStyle: {
          opacity: 1,
          shadowBlur: 10,
          shadowColor: 'rgba(0,0,0,0.2)',
        }
      },
      label: {
        show: true,
        position: 'inside',
        formatter: (params: any) => params.data.name.charAt(0),
        color: '#fff',
        fontSize: 12,
        fontWeight: 'bold',
      }
    }],
    // 四个象限的参考线
    MarkLine: {
      silent: true,
      symbol: 'none',
      lineStyle: { type: 'solid', color: '#d1d5db', width: 1 },
      data: [{
        xAxis: 50
      }, {
        yAxis: 50
      }]
    }
  }

  matrixChart.setOption(option)

  // 响应式
  window.addEventListener('resize', () => {
    matrixChart?.resize()
  })
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
onMounted(() => { loadData(); loadAllProjectsForStats(); loadUploadStats() })
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
.network-view {
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
  height: calc(100vh - 220px);
  display: flex;
  flex-direction: column;
}

.matrix-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
  min-height: 0;
}

/* 矩阵主区域 */
.matrix-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 24px;
  min-height: 0;
}

/* 标题栏 */
.matrix-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.matrix-header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.matrix-title {
  font-size: 20px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.matrix-desc {
  font-size: 13px;
  color: #94a3b8;
}

.matrix-date {
  font-size: 13px;
  color: #94a3b8;
}

/* 统计卡片 */
.matrix-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid #f1f5f9;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;

  &--blue { background: #eff6ff; color: #3b82f6; }
  &--green { background: #f0fdf4; color: #22c55e; }
  &--orange { background: #fff7ed; color: #f59e0b; }
  &--red { background: #fef2f2; color: #dc2626; }
  &--purple { background: #faf5ff; color: #9333ea; }
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #94a3b8;
}

/* 矩阵内容区 */
.matrix-content {
  flex: 1;
  display: flex;
  gap: 16px;
  min-height: 0;
}

/* Y轴标签 */
.axis-y-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 8px;
}

.axis-label-top,
.axis-label-bottom {
  font-size: 12px;
  color: #94a3b8;
  font-weight: 500;
}

.axis-label-text {
  font-size: 13px;
  color: #64748b;
  font-weight: 600;
  writing-mode: vertical-rl;
  text-orientation: mixed;
  letter-spacing: 2px;
}

/* 四象限网格容器 */
.matrix-grid-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 四象限网格 */
.matrix-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 12px;
  min-height: 0;
  background: #fff;
  border-radius: 12px;
  padding: 16px;
}

/* X轴标签 */
.axis-x-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 8px 20px;
}

.axis-label-left,
.axis-label-right {
  font-size: 12px;
  color: #94a3b8;
  font-weight: 500;
}

.axis-label-center {
  font-size: 13px;
  color: #64748b;
  font-weight: 600;
}

.axis-line {
  flex: 1;
  height: 1px;
  background: #e2e8f0;
  max-width: 80px;
}

/* 象限 */
.quadrant {
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #f1f5f9;
  transition: all 0.2s;

  &:hover {
    border-color: #e2e8f0;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  }

  &--tl {
    background: #fff7ed;
    .quadrant-header { background: #ea580c; }
  }

  &--tr {
    background: #f0fdf4;
    .quadrant-header { background: #16a34a; }
  }

  &--bl {
    background: #fef2f2;
    .quadrant-header { background: #dc2626; }
  }

  &--br {
    background: #faf5ff;
    .quadrant-header { background: #9333ea; }
  }
}

.quadrant-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  color: #fff;
}

.quadrant-icon {
  display: flex;
  align-items: center;
}

.quadrant-title {
  font-size: 14px;
  font-weight: 700;
  flex: 1;
}

.quadrant-count {
  background: rgba(255, 255, 255, 0.25);
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
}

.quadrant-body {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  gap: 12px;
  overflow-y: auto;
}

.quadrant-empty {
  width: 100%;
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
  padding: 16px;
}

/* 人员节点 */
.person-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-2px);
  }
}

.node-circle {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  border: 3px solid;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

.node-avatar {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.node-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 9px;
  background: #1e293b;
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #fff;
}

.node-name {
  font-size: 11px;
  font-weight: 600;
  color: #334155;
  text-align: center;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 图例说明 */
.matrix-legend {
  background: #fff;
  border-radius: 12px;
  padding: 16px 20px;
  border: 1px solid #f1f5f9;
}

.legend-section {
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f1f5f9;
}

.legend-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 4px;
}

.legend-desc {
  font-size: 12px;
  color: #64748b;
}

.legend-rules {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f1f5f9;
}

.legend-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 11px;
}

.legend-label {
  font-weight: 600;
  color: #475569;
  white-space: nowrap;
}

.legend-text {
  color: #64748b;
}

.legend-quadrants {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.legend-quadrant {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 3px;
}

.legend-name {
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  white-space: nowrap;
}

.legend-tip {
  font-size: 11px;
  color: #94a3b8;
}

/* 悬停详情卡片 */
.person-detail {
  position: absolute;
  top: 160px;
  right: 36px;
  width: 260px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  padding: 16px;
  z-index: 100;
  animation: slideIn 0.2s ease-out;
  border: 1px solid #f1f5f9;
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(-8px); }
  to { opacity: 1; transform: translateY(0); }
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.detail-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.detail-name {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
}

.detail-meta {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 2px;
}

.detail-progress {
  margin-bottom: 12px;
}

.progress-label {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #64748b;
  margin-bottom: 4px;
}

.progress-bar {
  height: 4px;
  background: #f1f5f9;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s;
}

.detail-roles {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 10px;
}

.role-tag {
  background: #f1f5f9;
  color: #475569;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
}

.detail-projects {
  background: #f8fafc;
  border-radius: 8px;
  padding: 10px;
}

.projects-title {
  font-size: 10px;
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.project-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 0;
  font-size: 12px;
  color: #475569;

  &:not(:last-child) {
    border-bottom: 1px solid #e2e8f0;
  }
}

.project-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* 空状态 */
.matrix-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 400px;
  background: #fff;
  border-radius: 16px;
}

.matrix-loading {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 400px;
  background: #fff;
  border-radius: 16px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 64px;
  opacity: 0.4;
}

.empty-text {
  font-size: 18px;
  font-weight: 600;
  color: #64748b;
}

.empty-hint {
  font-size: 14px;
  color: #94a3b8;
}

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
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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
}

.kanban-col {
  width: 260px;
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
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: calc(100vh - 380px);
  overflow-y: auto;
}

.kanban-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e8e5e1;
  transition: all 0.2s ease;
  cursor: grab;

  &:active { cursor: grabbing; }

  &:hover {
    border-color: #a8a29e;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
    transform: translateY(-2px);
  }

  &--dragging {
    opacity: 0.4;
    border-color: #a8a29e;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.kanban-card-inner {
  display: flex;
}

.kanban-card-content {
  flex: 1;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.kanban-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
}

.kanban-avatar {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.kanban-info {
  flex: 1;
  min-width: 0;
}

.kanban-name {
  font-size: 13px;
  font-weight: 600;
  color: #1c1917;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.kanban-code {
  font-size: 11px;
  color: #a8a29e;
}

.kanban-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 10px;
}

.kanban-members {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 10px;
  padding: 10px;
  background: #fafaf9;
  border-radius: 10px;
  border: 1px solid #f0ede8;
}

.kanban-member {
  display: flex;
  align-items: center;
  gap: 4px;
}

.kanban-member-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.kanban-member-name {
  font-size: 11px;
  color: #57534e;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kanban-member-more {
  font-size: 11px;
  color: #a8a29e;
  padding: 2px 8px;
  background: #f5f5f4;
  border-radius: 8px;
}

.kanban-stations {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 10px;
  padding: 10px;
  background: #fafaf9;
  border-radius: 10px;
  border: 1px solid #f0ede8;
}

.kanban-station {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #57534e;

  .el-icon { color: #a8a29e; font-size: 12px; flex-shrink: 0; }
  span:first-of-type { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
}

.station-phone {
  font-size: 11px;
  color: #a8a29e;
  flex-shrink: 0;
}

.kanban-station-more {
  font-size: 11px;
  color: #a8a29e;
  padding-left: 20px;
}

.station-person-opt {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 0;

  small { font-size: 11px; color: #a8a29e; }
}

.kanban-card-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid #f0ede8;
}

.kanban-actions {
  display: flex;
  gap: 4px;
  margin-left: auto;
}

.kanban-metric {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #a8a29e;

  .el-icon { font-size: 12px; }
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

.bento-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
  align-items: start;
}

.bento-card {
  background: #fafaf9;
  border-radius: 16px;
  border: 1px solid #e8e5e1;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: fadeSlideIn 0.4s ease both;
  position: relative;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.08);
    border-color: #d4d0c8;

    .bento-hover-actions { opacity: 1; }
    .bento-avatar { transform: scale(1.04); }
  }
}

/* 左侧阶段色条 - 日式细线 */
.bento-stage-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 3px;
  height: 100%;
  opacity: 0.7;
}

/* 卡片顶部 */
.bento-top {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 20px 18px 16px;
  background: #fff;
  position: relative;
  border-bottom: 1px solid #f0ede8;
}

.bento-top-info {
  flex: 1;
  min-width: 0;
  padding-left: 4px;
}

.bento-avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transition: transform 0.25s ease;
}

.bento-hover-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  align-items: center;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;

  .el-button {
    color: #a8a29e;
    padding: 5px 7px;
    background: rgba(250, 250, 249, 0.9);

    &:hover { color: #57534e; background: #f5f5f4; }
  }

  .btn-delete:hover { color: #dc2626; }
}

/* 卡片主体 */
.bento-body {
  padding: 16px 18px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 项目名称 */
.bento-name {
  font-size: 15px;
  font-weight: 700;
  color: #1c1917;
  line-height: 1.4;
  letter-spacing: -0.3px;
}

.bento-code {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #a8a29e;
  cursor: pointer;
  padding: 3px 8px;
  border-radius: 6px;
  margin-top: 4px;
  transition: all 0.2s;
  background: #f5f5f4;

  .el-icon { color: #d4d0c8; transition: color 0.2s; }

  &:hover {
    color: #78716c;
    background: #f0ede8;
    .el-icon { color: #78716c; }
  }
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

/* 描述 */
.bento-desc {
  font-size: 13px;
  color: #78716c;
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 项目信息行 */
.bento-info-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.bento-info-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #78716c;
  background: #f5f5f4;
  border-radius: 8px;
  padding: 5px 10px;

  .el-icon { color: #a8a29e; }

  // 标签样式统一，与团队成员标签对齐
  .info-label {
    color: #a8a29e;
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.3px;
    margin-right: 2px;
  }
}

/* 解决方案 */
.bento-solution {
  border-left: 3px solid #d6d3d1;
  padding-left: 14px;
}

.bento-solution-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: #a8a29e;
  margin-bottom: 6px;
  font-weight: 500;
  letter-spacing: 0.3px;

  .el-icon { color: #a8a29e; }
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
  gap: 6px;
  font-size: 11px;
  color: #a8a29e;
  margin-bottom: 10px;
  font-weight: 500;
  letter-spacing: 0.3px;

  .el-icon { color: #a8a29e; }
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
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-section-content {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.modal-info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.modal-info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #909399;

  .el-icon { color: #c0c4cc; }
  span { flex-shrink: 0; }
  strong { color: #303133; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

  &--full { grid-column: 1 / -1; }
}

.modal-team-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
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
  .modal-info-grid { grid-template-columns: 1fr; }
  .modal-info-item--full { grid-column: 1; }
}

@media (max-width: 768px) {
  .kpi-cards { grid-template-columns: repeat(3, 1fr); }
  .filter-bar { flex-direction: column; align-items: stretch; }
  .filter-tip { margin-left: 0; }
}
</style>
