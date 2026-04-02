<template>
  <div class="project-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('project.list.title') }}</h1>
        <span class="page-subtitle">{{ t('project.list.subtitle') }}</span>
      </div>
      <div class="header-actions">
      <div class="view-tabs">
        <button class="view-tab" :class="{ active: viewMode === 'detail' }" @click="switchView('detail')">
          <el-icon><Tickets /></el-icon>
          <span>项目详情</span>
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'kanban' }" @click="switchView('kanban')">
          <el-icon><Menu /></el-icon>
          <span>阶段看板</span>
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'network' }" @click="switchView('network')">
          <el-icon><Connection /></el-icon>
          <span>协作网络</span>
        </button>
      </div>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ t('project.list.addProject') }}
        </el-button>
      </div>
    </header>

    <!-- KPI 统计 -->
    <div class="kpi-cards" v-if="pagination.total > 0">
      <div class="kpi-card" :class="{ 'kpi-card--active': !searchKeyword && !searchStatus && !searchStage }" @click="handleReset">
        <div class="kpi-num">{{ pagination.total }}</div>
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
    <div class="stage-distribution" v-if="!loading && pagination.total > 0">
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
    <div class="filter-bar">
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
            <span class="kanban-records">{{ getStageRecords(col.key) }} 条记录</span>
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
              <div class="kanban-card-left" :style="{ background: col.color }"></div>
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

    <!-- 协作网络视图 -->
    <div v-show="viewMode === 'network' && !loading" class="network-view">
      <!-- 网络图容器 -->
      <div
        class="network-canvas"
        ref="networkCanvas"
        @wheel.prevent="handleZoom"
        @mousedown="startPan"
        @mousemove="doPan"
        @mouseup="endPan"
        @mouseleave="endPan"
      >
        <!-- SVG 层 -->
        <svg class="network-svg" :viewBox="`0 0 ${svgWidth} ${svgHeight}`" :style="svgTransform">
          <defs>
            <filter id="glow">
              <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
              <feMerge>
                <feMergeNode in="coloredBlur"/>
                <feMergeNode in="SourceGraphic"/>
              </feMerge>
            </filter>
            <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
              <polygon points="0 0, 10 3.5, 0 7" fill="#94a3b8"/>
            </marker>
          </defs>

          <!-- 关系连线 -->
          <g class="links-layer">
            <path
              v-for="(link, idx) in meshLinks"
              :key="'link-' + idx"
              :d="link.path"
              class="mesh-link"
              :class="{
                'mesh-link--active': link.isActive,
                'mesh-link--dim': hoveredNode && !link.isActive && !link.involved
              }"
              :stroke="link.isActive ? link.personColor : '#dde3ed'"
              :stroke-width="link.isActive ? 2.5 : 1"
              fill="none"
              :filter="link.isActive ? 'url(#glow)' : 'none'"
            />
          </g>

          <!-- 人员节点 -->
          <g
            v-for="person in meshPersons"
            :key="'person-' + person.id"
            class="node-group"
            :transform="`translate(${person.x}, ${person.y})`"
            @mouseenter="handleNodeHover(person, 'person')"
            @mouseleave="handleNodeLeave"
            @mousedown.stop="startDragNode($event, person, 'person')"
          >
            <!-- 光晕效果 -->
            <circle
              v-if="hoveredNode?.id === person.id || person.isActive"
              :r="person.radius + 8"
              :fill="person.color"
              opacity="0.2"
              class="node-glow"
            />
            <!-- 主圆 -->
            <circle
              :r="person.radius"
              :fill="person.color"
              class="node-circle"
              :class="{ 'node-circle--active': hoveredNode?.id === person.id || person.isActive }"
            />
            <!-- 首字母 -->
            <text
              :y="1"
              text-anchor="middle"
              dominant-baseline="middle"
              fill="#fff"
              font-size="14"
              font-weight="600"
            >{{ person.name.charAt(0).toUpperCase() }}</text>
            <!-- 名字标签 -->
            <text
              :y="person.radius + 16"
              text-anchor="middle"
              fill="#374151"
              font-size="11"
              font-weight="500"
            >{{ person.name }}</text>
            <!-- 项目数量 -->
            <text
              :y="person.radius + 28"
              text-anchor="middle"
              fill="#9ca3af"
              font-size="10"
            >{{ person.projectCount }} 项目</text>
          </g>

          <!-- 项目节点 -->
          <g
            v-for="project in meshProjects"
            :key="'project-' + project.id"
            class="node-group"
            :transform="`translate(${project.x}, ${project.y})`"
            @mouseenter="handleNodeHover(project, 'project')"
            @mouseleave="handleNodeLeave"
            @mousedown.stop="startDragNode($event, project, 'project')"
          >
            <!-- 光晕效果 -->
            <circle
              v-if="hoveredNode?.id === project.id || project.isActive"
              :r="project.radius + 8"
              fill="#f56c6c"
              opacity="0.2"
              class="node-glow"
            />
            <!-- 主圆 -->
            <rect
              :x="-project.radius"
              :y="-project.radius"
              :width="project.radius * 2"
              :height="project.radius * 2"
              rx="8"
              :fill="'#f56c6c'"
              class="node-rect"
              :class="{ 'node-rect--active': hoveredNode?.id === project.id || project.isActive }"
            />
            <!-- 图标 -->
            <text
              :y="1"
              text-anchor="middle"
              dominant-baseline="middle"
              fill="#fff"
              font-size="14"
            >📁</text>
            <!-- 项目名 -->
            <text
              :y="project.radius + 16"
              text-anchor="middle"
              fill="#374151"
              font-size="11"
              font-weight="500"
            >{{ project.name.length > 10 ? project.name.slice(0, 10) + '...' : project.name }}</text>
            <!-- 人数 -->
            <text
              :y="project.radius + 28"
              text-anchor="middle"
              fill="#9ca3af"
              font-size="10"
            >{{ project.personCount }} 人</text>
          </g>
        </svg>

        <!-- 缩放控制 -->
        <div class="zoom-controls">
          <button class="zoom-btn" @click="zoomIn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <span class="zoom-level">{{ Math.round(zoomLevel * 100) }}%</span>
          <button class="zoom-btn" @click="zoomOut">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <button class="zoom-btn" @click="resetView">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
          </button>
        </div>

        <!-- 提示信息 -->
        <div class="network-hint">
          <span>🖱️ 拖拽画布移动</span>
          <span>滚轮缩放</span>
          <span>拖拽节点重排</span>
        </div>
      </div>

      <!-- 关系详情面板 -->
      <div class="network-sidebar">
        <!-- 当前选中/悬停节点信息 -->
        <div class="sidebar-card" v-if="hoveredNode || selectedMeshNode">
          <div class="sidebar-header">
            <span class="sidebar-icon" :style="{ background: (hoveredNode || selectedMeshNode)?.color || '#409eff' }">
              {{ (hoveredNode || selectedMeshNode)?.type === 'person' ? '👤' : '📁' }}
            </span>
            <div class="sidebar-title">
              <span class="sidebar-name">{{ (hoveredNode || selectedMeshNode)?.name }}</span>
              <span class="sidebar-type">{{ (hoveredNode || selectedMeshNode)?.type === 'person' ? '人员' : '项目' }}</span>
            </div>
          </div>

          <div class="sidebar-stats">
            <div class="sidebar-stat">
              <span class="stat-num">{{ (hoveredNode || selectedMeshNode)?.type === 'person' ? (hoveredNode || selectedMeshNode)?.projectCount : (hoveredNode || selectedMeshNode)?.personCount }}</span>
              <span class="stat-label">{{ (hoveredNode || selectedMeshNode)?.type === 'person' ? '参与项目' : '项目成员' }}</span>
            </div>
          </div>

          <!-- 关联关系列表 -->
          <div class="related-list" v-if="relatedItems.length > 0">
            <div class="related-title">关联关系</div>
            <div
              v-for="item in relatedItems"
              :key="item.id"
              class="related-item"
              :class="{ 'related-item--dim': hoveredNode && item.id !== hoveredNode.id }"
            >
              <div class="related-avatar" :style="{ background: item.color }">
                {{ item.name.charAt(0).toUpperCase() }}
              </div>
              <div class="related-info">
                <span class="related-name">{{ item.name }}</span>
                <span class="related-type">{{ item.type === 'person' ? '人员' : '项目' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 网络统计 -->
        <div class="sidebar-card">
          <div class="sidebar-header">
            <span class="sidebar-title-text">网络概览</span>
          </div>
          <div class="overview-stats">
            <div class="overview-item">
              <span class="overview-num">{{ meshPersons.length }}</span>
              <span class="overview-label">人员节点</span>
            </div>
            <div class="overview-item">
              <span class="overview-num">{{ meshProjects.length }}</span>
              <span class="overview-label">项目节点</span>
            </div>
            <div class="overview-item">
              <span class="overview-num">{{ meshLinks.length }}</span>
              <span class="overview-label">关系连线</span>
            </div>
          </div>
        </div>

        <!-- 核心人员 -->
        <div class="sidebar-card">
          <div class="sidebar-header">
            <span class="sidebar-title-text">核心人员</span>
          </div>
          <div class="core-list">
            <div
              v-for="person in meshPersons.filter(p => p.projectCount >= 3).slice(0, 6)"
              :key="person.id"
              class="core-item"
            >
              <div class="core-avatar" :style="{ background: person.color }">
                {{ person.name.charAt(0).toUpperCase() }}
              </div>
              <div class="core-info">
                <span class="core-name">{{ person.name }}</span>
                <span class="core-count">{{ person.projectCount }} 个项目</span>
              </div>
              <div class="core-bar">
                <div class="core-fill" :style="{ width: Math.min(person.projectCount / 5 * 100, 100) + '%', background: person.color }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 瀑布视图 -->
    <div v-show="viewMode === 'detail'" class="bento-view">
      <div class="bento-grid">
        <div
          v-for="(project, idx) in tableData"
          :key="project.id"
          class="bento-card"
          :class="{ 'bento-card--active': project.status === 'active' }"
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
                  class="bento-member-chip"
                  :title="member + ' · ' + getMemberRole(project, member)"
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

            <!-- 驻场点 -->
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
                  <span class="station-person" v-if="s.person">{{ s.person }}</span>
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
        <el-icon class="empty-icon"><Folder /></el-icon>
        <p class="empty-title">{{ t('project.list.noProjectData') }}</p>
        <p class="empty-hint">{{ t('project.list.noProjectHint') }}</p>
      </div>
    </div>

    <!-- 加载骨架 -->
    <div v-if="loading" class="project-grid">
      <div v-for="i in 8" :key="i" class="skeleton-card"></div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="pagination.total > 0">
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
            <div class="form-section-title">{{ t('project.form.basicInfo') }}</div>
            <el-form-item :label="t('project.form.nameLabel')" prop="name">
              <el-input v-model="form.name" :placeholder="t('project.form.namePlaceholder')" @input="dialogPreviewInitials = getInitials(form.name)" />
            </el-form-item>
            <el-form-item :label="t('project.form.codeLabel')" prop="code">
              <el-input v-model="form.code" :placeholder="t('project.form.codePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('project.form.descriptionLabel')">
              <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('project.form.descriptionPlaceholder')" />
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item :label="t('project.form.stageLabel')">
                  <el-select v-model="form.stage" style="width: 100%">
                    <el-option v-for="s in stageOptions" :key="s.value" :label="t('project.stage.' + s.value)" :value="s.value" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.form.statusLabel')">
                  <el-select v-model="form.status" style="width: 100%">
                    <el-option :label="t('common.enabled')" value="active" />
                    <el-option :label="t('common.disabled')" value="inactive" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <!-- 团队成员 -->
          <div class="form-section">
            <div class="form-section-title">{{ t('project.form.teamMembers') }}</div>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item :label="t('project.role.projectPerson')">
                  <el-select v-model="form.projectPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.role.opsPerson')">
                  <el-select v-model="form.opsPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.role.opsStaffPerson')">
                  <el-select v-model="form.opsStaffPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.role.developerPerson')">
                  <el-select v-model="form.developerPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.role.testerPerson')">
                  <el-select v-model="form.testerPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.role.businessPerson')">
                  <el-select v-model="form.businessPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                    <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <!-- 项目信息 -->
          <div class="form-section">
            <div class="form-section-title">{{ t('project.form.projectInfo') }}</div>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item :label="t('project.form.companyAddr')">
                  <el-input v-model="form.companyAddr" :placeholder="t('project.form.companyAddrPlaceholder')" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('project.form.projectPeriod')">
                  <el-date-picker
                    v-model="projectPeriodRange"
                    type="daterange"
                    range-separator="~"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="YYYY-MM-DD"
                    format="YYYY-MM-DD"
                    style="width: 100%"
                    :shortcuts="dateShortcuts"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item :label="t('project.form.solutionLabel')">
              <el-input v-model="form.solution" type="textarea" :rows="2" :placeholder="t('project.form.solutionPlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('project.form.solutionPersonLabel')">
              <el-select v-model="form.solutionPerson" :placeholder="t('project.form.selectPerson')" :loading="personnelLoading" filterable clearable style="width: 100%">
                <el-option v-for="p in personnelList" :key="p.id" :label="p.name" :value="p.name" />
              </el-select>
            </el-form-item>
          </div>

          <!-- 驻场点 -->
          <div class="form-section">
            <div class="form-section-title">
              {{ t('project.list.onsiteStations') }}
              <el-badge v-if="form.onsiteStations.length > 0" :value="form.onsiteStations.length" />
            </div>
            <div v-for="(station, idx) in form.onsiteStations" :key="idx" class="station-row">
              <span class="station-num">{{ idx + 1 }}</span>
              <el-select
                v-model="station.person"
                :placeholder="t('project.form.selectPerson')"
                size="small"
                :loading="personnelLoading"
                filterable
                clearable
                style="width: 160px"
                @change="(val) => fillStationFromPersonnel(station, val)"
              >
                <el-option
                  v-for="p in personnelList"
                  :key="p.id"
                  :label="p.name"
                  :value="p.name"
                >
                  <div class="station-person-opt">
                    <span>{{ p.name }}</span>
                    <small v-if="p.company">{{ p.company }}</small>
                  </div>
                </el-option>
              </el-select>
              <el-input v-model="station.location" :placeholder="t('project.form.stationLocationPlaceholder')" size="small" style="width: 140px" />
              <el-input v-model="station.phone" :placeholder="t('project.form.stationContactPlaceholder')" size="small" style="width: 120px" />
              <el-button size="small" text @click="removeStation(idx)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
              </el-button>
            </div>
            <el-button type="primary" plain size="small" @click="addStation">
              <el-icon><Plus /></el-icon>
              {{ t('project.form.addStation') }}
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

        <!-- 驻场点 -->
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
import {
  Grid, Menu, Plus, Search, Document, Cpu, User,
  Folder, Location, Clock, Connection, LocationInformation, Tickets, Phone, Timer, Loading, ArrowDown
} from '@element-plus/icons-vue'
import { ProjectApi, type Project, type CreateProjectReq, type UpdateProjectReq, type OnSiteStation } from '@/api/project'
import { PersonnelApi, type Personnel } from '@/api/personnel'
import { UploadRecordApi } from '@/api/upload-record'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Project[]>([])
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

const stageOptions = [
  { value: 'planning' },
  { value: 'designing' },
  { value: 'deploying' },
  { value: 'running' },
  { value: 'paused' },
]

const kanbanColumns = [
  { key: 'planning', color: '#94a3b8' },
  { key: 'designing', color: '#6366f1' },
  { key: 'deploying', color: '#f59e0b' },
  { key: 'running', color: '#22c55e' },
  { key: 'paused', color: '#ef4444' },
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

// 各阶段项目数量
const stageCount = computed(() => ({
  planning: tableData.value.filter(p => p.stage === 'planning').length,
  designing: tableData.value.filter(p => p.stage === 'designing').length,
  deploying: tableData.value.filter(p => p.stage === 'deploying').length,
  running: tableData.value.filter(p => p.stage === 'running').length,
  paused: tableData.value.filter(p => p.stage === 'paused').length
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
    const fields = ['projectPerson', 'opsPerson', 'developerPerson', 'testerPerson', 'businessPerson', 'compliancePerson', 'opsStaffPerson', 'solutionPerson']
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

// ==================== 协作网络图 ====================
// ============ 网格网络视图接口 ============
interface MeshNode {
  id: string
  name: string
  x: number
  y: number
  vx: number
  vy: number
  radius: number
  color: string
  projectCount: number
  personCount: number
  projectIds?: number[]
  personNames?: string[]
  type: 'person' | 'project'
  isActive: boolean
}

interface MeshLink {
  sourceId: string
  targetId: string
  path: string
  personColor: string
  isActive: boolean
  involved: boolean
}

const networkCanvas = ref<HTMLElement>()
const svgWidth = ref(1000)
const svgHeight = ref(600)
const zoomLevel = ref(1)
const panX = ref(0)
const panY = ref(0)
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })

// 网格网络数据
const meshPersons = ref<MeshNode[]>([])
const meshProjects = ref<MeshNode[]>([])
const meshLinks = ref<MeshLink[]>([])
const hoveredNode = ref<MeshNode | null>(null)
const selectedMeshNode = ref<MeshNode | null>(null)
const draggingMeshNode = ref<{ node: MeshNode; type: string } | null>(null)
const dragStartOffset = ref({ x: 0, y: 0 })

const svgTransform = computed(() => ({
  transform: `translate(${panX.value}px, ${panY.value}px) scale(${zoomLevel.value})`
}))

const relatedItems = computed(() => {
  if (!hoveredNode.value && !selectedMeshNode.value) return []
  const node = hoveredNode.value || selectedMeshNode.value
  if (!node) return []

  const items: MeshNode[] = []
  if (node.type === 'person' && node.projectIds) {
    node.projectIds.forEach(pid => {
      const proj = meshProjects.value.find(p => p.id === pid)
      if (proj) items.push(proj)
    })
  } else if (node.type === 'project' && node.personNames) {
    node.personNames.forEach(name => {
      const person = meshPersons.value.find(p => p.name === name)
      if (person) items.push(person)
    })
  }
  return items
})

// 初始化网格网络数据
const initMeshNetwork = () => {
  if (!networkCanvas.value) return

  const container = networkCanvas.value
  svgWidth.value = container.offsetWidth
  svgHeight.value = container.offsetHeight

  const persons: Record<string, { name: string; projectIds: number[] }> = {}
  const projects: Record<number, { name: string; personNames: string[] }> = {}
  const personFields = ['projectPerson', 'opsPerson', 'developerPerson', 'testerPerson', 'businessPerson', 'compliancePerson', 'opsStaffPerson', 'solutionPerson']

  tableData.value.forEach(p => {
    projects[p.id] = { name: p.name, personNames: [] }
    personFields.forEach(field => {
      const val = (p as any)[field]
      if (val) {
        val.split(/[,，、]/).filter((n: string) => n.trim()).forEach((name: string) => {
          const trimmedName = name.trim()
          if (!persons[trimmedName]) {
            persons[trimmedName] = { name: trimmedName, projectIds: [] }
          }
          if (!persons[trimmedName].projectIds.includes(p.id)) {
            persons[trimmedName].projectIds.push(p.id)
          }
          if (!projects[p.id].personNames.includes(trimmedName)) {
            projects[p.id].personNames.push(trimmedName)
          }
        })
      }
    })
  })

  const centerX = svgWidth.value / 2
  const centerY = svgHeight.value / 2
  const personCount = Object.keys(persons).length
  const projectCount = Object.keys(projects).length
  const radius = Math.min(svgWidth.value, svgHeight.value) * 0.35

  // 圆形布局 - 人员在左半圆，项目在右半圆
  const personColors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#c71585', '#00b050', '#005eeb', '#ff6b6b', '#4ecdc4']

  meshPersons.value = Object.entries(persons).map(([name, data], idx) => {
    const angle = (Math.PI * (idx + 1)) / (personCount + 1) - Math.PI / 2
    const x = centerX - radius * Math.cos(angle) * 0.7
    const y = centerY + radius * Math.sin(angle) * 0.5
    return {
      id: 'person_' + name,
      name,
      x,
      y,
      vx: 0,
      vy: 0,
      radius: 24 + Math.min(data.projectIds.length * 2, 12),
      color: personColors[idx % personColors.length],
      projectCount: data.projectIds.length,
      personCount: 0,
      projectIds: data.projectIds,
      type: 'person' as const,
      isActive: false
    }
  })

  meshProjects.value = Object.entries(projects).map(([id, data], idx) => {
    const angle = (Math.PI * (idx + 1)) / (projectCount + 1) + Math.PI / 2
    const x = centerX + radius * Math.cos(angle) * 0.7
    const y = centerY + radius * Math.sin(angle) * 0.5
    return {
      id: 'project_' + id,
      name: data.name,
      x,
      y,
      vx: 0,
      vy: 0,
      radius: 22 + Math.min(data.personNames.length * 2, 14),
      color: '#f56c6c',
      projectCount: 0,
      personCount: data.personNames.length,
      personNames: data.personNames,
      type: 'project' as const,
      isActive: false
    }
  })

  updateMeshLinks()
}

const updateMeshLinks = () => {
  const links: MeshLink[] = []

  meshPersons.value.forEach(person => {
    person.projectIds?.forEach(projId => {
      const project = meshProjects.value.find(p => p.id === 'project_' + projId)
      if (project) {
        // 创建贝塞尔曲线
        const midX = (person.x + project.x) / 2
        const midY = (person.y + project.y) / 2 - 30
        const path = `M ${person.x} ${person.y} Q ${midX} ${midY} ${project.x} ${project.y}`

        links.push({
          sourceId: person.id,
          targetId: project.id,
          path,
          personColor: person.color,
          isActive: false,
          involved: false
        })
      }
    })
  })

  meshLinks.value = links
}

// 节点悬停
const handleNodeHover = (node: MeshNode, type: string) => {
  hoveredNode.value = node

  // 高亮相关连线
  meshLinks.value.forEach(link => {
    link.isActive = link.sourceId === node.id || link.targetId === node.id
    link.involved = link.sourceId === node.id || link.targetId === node.id
  })

  // 高亮相关节点
  if (node.type === 'person') {
    meshProjects.value.forEach(p => {
      p.isActive = node.projectIds?.includes(parseInt(p.id.replace('project_', '')))
    })
    meshPersons.value.forEach(p => {
      p.isActive = p.id === node.id
    })
  } else {
    meshPersons.value.forEach(p => {
      p.isActive = node.personNames?.includes(p.name)
    })
    meshProjects.value.forEach(p => {
      p.isActive = p.id === node.id
    })
  }
}

const handleNodeLeave = () => {
  hoveredNode.value = null
  meshLinks.value.forEach(link => {
    link.isActive = false
    link.involved = false
  })
  meshPersons.value.forEach(p => p.isActive = false)
  meshProjects.value.forEach(p => p.isActive = false)
}

// 拖拽节点
const startDragNode = (e: MouseEvent, node: MeshNode, type: string) => {
  draggingMeshNode.value = { node, type }
  dragStartOffset.value = { x: e.offsetX, y: e.offsetY }

  const onMove = (ev: MouseEvent) => {
    if (!draggingMeshNode.value || !networkCanvas.value) return
    const rect = networkCanvas.value.getBoundingClientRect()
    const newX = (ev.clientX - rect.left) / zoomLevel.value - dragStartOffset.value.x
    const newY = (ev.clientY - rect.top) / zoomLevel.value - dragStartOffset.value.y

    if (type === 'person') {
      const idx = meshPersons.value.findIndex(p => p.id === node.id)
      if (idx > -1) {
        meshPersons.value[idx].x = newX
        meshPersons.value[idx].y = newY
      }
    } else {
      const idx = meshProjects.value.findIndex(p => p.id === node.id)
      if (idx > -1) {
        meshProjects.value[idx].x = newX
        meshProjects.value[idx].y = newY
      }
    }

    updateMeshLinks()
  }

  const onUp = () => {
    draggingMeshNode.value = null
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }

  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

// 缩放控制
const handleZoom = (e: WheelEvent) => {
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  zoomLevel.value = Math.max(0.3, Math.min(3, zoomLevel.value + delta))
}

const zoomIn = () => { zoomLevel.value = Math.min(3, zoomLevel.value + 0.2) }
const zoomOut = () => { zoomLevel.value = Math.max(0.3, zoomLevel.value - 0.2) }
const resetView = () => {
  zoomLevel.value = 1
  panX.value = 0
  panY.value = 0
}

// 画布拖拽平移
const startPan = (e: MouseEvent) => {
  if (e.target === networkCanvas.value || (e.target as HTMLElement).classList.contains('network-svg')) {
    isPanning.value = true
    panStart.value = { x: e.clientX - panX.value, y: e.clientY - panY.value }
  }
}

const doPan = (e: MouseEvent) => {
  if (isPanning.value) {
    panX.value = e.clientX - panStart.value.x
    panY.value = e.clientY - panStart.value.y
  }
}

const endPan = () => {
  isPanning.value = false
}

const form = reactive<CreateProjectReq & { id?: number }>({
  name: '', code: '', description: '', status: 'active', stage: 'planning', sort: 0,
  projectPerson: '', opsPerson: '', opsStaffPerson: '', developerPerson: '', testerPerson: '', businessPerson: '', compliancePerson: '',
  solution: '', solutionPerson: '', companyAddr: '', projectPeriod: '', onsiteStations: [],
})

// 项目周期：双向绑定日期范围选择器
const projectPeriodRange = computed({
  get: () => {
    if (!form.projectPeriod) return null
    const parts = form.projectPeriod.split(' ~ ')
    if (parts.length === 2) return [parts[0], parts[1]]
    return null
  },
  set: (val: [string, string] | null) => {
    form.projectPeriod = val ? `${val[0]} ~ ${val[1]}` : ''
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
  '#3b82f6', // blue
  '#10b981', // emerald
  '#f59e0b', // amber
  '#ef4444', // red
  '#8b5cf6', // violet
  '#06b6d4', // cyan
  '#ec4899', // pink
  '#14b8a6', // teal
  '#f97316', // orange
  '#84cc16', // lime
  '#6366f1', // indigo
  '#a855f7', // purple
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

const avatarColors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#8b5cf6', '#06b6d4', '#f59e0b', '#ec4899']
const getAvatarColor = (name: string) => {
  if (!name) return avatarColors[0]
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = ((hash << 5) - hash) + name.charCodeAt(i)
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

const getTeamMembers = (project: Project) => {
  const members: string[] = []
  if (project.projectPerson) members.push(project.projectPerson)
  if (project.opsPerson) members.push(project.opsPerson)
  if (project.opsStaffPerson) members.push(project.opsStaffPerson)
  if (project.developerPerson) members.push(project.developerPerson)
  if (project.testerPerson) members.push(project.testerPerson)
  if (project.businessPerson) members.push(project.businessPerson)
  if (project.compliancePerson) members.push(project.compliancePerson)
  if (project.solutionPerson) members.push(project.solutionPerson)
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
    ['solutionPerson', '方案人员'],
  ]
  for (const [key, label] of roleMap) {
    if ((project as any)[key] === member) return label
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

const loadKanbanData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.getKanbanList()
    const list: any[] = Array.isArray(res) ? res : (res?.data ?? [])
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
  if (mode === 'kanban') loadKanbanData()
  else if (mode === 'network') {
    loadData()
    nextTick(() => initMeshNetwork())
  }
  else loadData()
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
  searchStage.value = searchStage.value === stage ? '' : stage
  pagination.page = 1
  loadData()
}

// 复制到剪贴板
const copyCode = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
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
    developerPerson: row.developerPerson || '', testerPerson: row.testerPerson || '',
    businessPerson: row.businessPerson || '', compliancePerson: row.compliancePerson || '',
    solution: row.solution || '', solutionPerson: row.solutionPerson || '',
    companyAddr: row.companyAddr || '', projectPeriod: row.projectPeriod || '',
    onsiteStations: row.onsiteStations ? [...row.onsiteStations.map(s => ({ ...s }))] : [],
  })
  dialogPreviewInitials.value = getInitials(row.name)
  await loadPersonnelList()
  drawerVisible.value = true
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
    projectPerson: '', opsPerson: '', opsStaffPerson: '', developerPerson: '', testerPerson: '', businessPerson: '', compliancePerson: '',
    solution: '', solutionPerson: '', companyAddr: '', projectPeriod: '', onsiteStations: [],
  })
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

// 从人员列表自动填入驻场点
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
onMounted(() => { loadData(); loadUploadStats() })
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

/* 视图切换 */
.view-tabs {
  display: flex;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
}

.view-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: none;
  background: transparent;
  font-size: 13px;
  font-weight: 500;
  color: #909399;
  cursor: pointer;
  transition: all 0.2s;
  border-right: 1px solid #ebeef5;

  &:last-child { border-right: none; }
  &:hover { background: #f5f7fa; color: #606266; }
  &.active {
    background: #409eff;
    color: #fff;
    .el-icon { color: #fff; }
  }
}

/* 阶段分布迷你图 */
.stage-distribution {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.stage-dist-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.dist-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stage-dist-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
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
  font-size: 14px;
  font-weight: 700;
  color: #303133;
}

.dist-stat-label {
  font-size: 11px;
  color: #909399;
}

.dist-sep {
  color: #dcdfe6;
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
  background: #f0f0f0;
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
  color: #909399;
}

.workload--normal .workload-fill { background: #409eff; }
.workload--warning .workload-fill { background: #e6a23c; }
.workload--idle .workload-fill { background: #67c23a; }

.stage-dist-chart {
  display: flex;
  gap: 20px;
}

.stage-dist-bar {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  height: 40px;
  flex: 1;
}

.stage-dist-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  min-width: 36px;
}

.stage-dist-fill {
  width: 100%;
  min-height: 3px;
  border-radius: 3px 3px 0 0;
  transition: height 0.3s;
}

.stage-dist-count {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.stage-dist-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: var(--color-text-secondary);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.kpi-cards {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.kpi-card {
  background: #fff;
  border-radius: 10px;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-left: 3px solid #64748b;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s, border-color 0.2s, border-left-color 0.2s;

  &:hover {
    box-shadow: 0 4px 16px rgba(0,0,0,0.08);
    transform: translateY(-2px);
  }

  &--active {
    border-color: #409eff;
    border-left-color: #22c55e;
    background: #ecf5ff;

    &.kpi-card--deploying { border-left-color: #f59e0b; }
    &.kpi-card--planning { border-left-color: #6366f1; }
    &.kpi-card--data { border-left-color: #409eff; }
    &.kpi-card--records { border-left-color: #8b5cf6; }
  }

  &--running {
    .kpi-num { color: #22c55e; }
    border-left: 3px solid #22c55e;
  }

  &--deploying {
    .kpi-num { color: #f59e0b; }
    border-left: 3px solid #f59e0b;
  }

  &--planning {
    .kpi-num { color: #6366f1; }
    border-left: 3px solid #6366f1;
  }

  &--data {
    .kpi-num { color: #409eff; }
    border-left: 3px solid #409eff;
  }

  &--records {
    .kpi-num { color: #8b5cf6; }
    border-left: 3px solid #8b5cf6;
  }
}

.kpi-num {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 6px;
}

.kpi-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

/* ==================== 协作网络视图 ==================== */
.network-view {
  display: flex;
  gap: 16px;
  height: calc(100vh - 280px);
}

.network-canvas {
  flex: 1;
  background: #f8fafc;
  border-radius: 16px;
  position: relative;
  overflow: hidden;
  cursor: grab;
}

.network-canvas:active {
  cursor: grabbing;
}

.network-svg {
  width: 100%;
  height: 100%;
  transform-origin: center center;
  transition: transform 0.1s ease-out;
}

/* 连线样式 */
.mesh-link {
  stroke-linecap: round;
  transition: stroke 0.3s, stroke-width 0.3s, opacity 0.3s;
}

.mesh-link--active {
  stroke-dasharray: 8 4;
  animation: dash 0.5s linear infinite;
}

.mesh-link--dim {
  opacity: 0.15;
}

@keyframes dash {
  to { stroke-dashoffset: -12; }
}

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

/* 看板 */
.kanban-board {
  display: flex;
  gap: 14px;
  overflow-x: auto;
  padding-bottom: 16px;
}

.kanban-col {
  width: 260px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #ebeef5;
  transition: all 0.2s;

  &--collapsed {
    width: 160px;
    .kanban-col-head { cursor: pointer; }
    .kanban-stats, .collapse-icon { display: none; }
  }

  &--dragover {
    border-color: #409eff;
    background-color: #ecf5ff;
  }
}

.kanban-col-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  border-bottom: 1px solid #f0f0f0;
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
  color: #303133;
}

.kanban-stats {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
}

.kanban-count {
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
  color: #606266;
}

.kanban-records {
  color: #909399;
}

.collapse-icon {
  color: #c0c4cc;
  transition: transform 0.2s;

  &.is-collapsed {
    transform: rotate(-90deg);
  }
}

.kanban-col-body {
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: calc(100vh - 380px);
  overflow-y: auto;
}

.kanban-card {
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  transition: all 0.2s;
  cursor: grab;

  &:active { cursor: grabbing; }

  &:hover {
    border-color: #409eff;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
  }

  &--dragging {
    opacity: 0.4;
    border-color: #409eff;
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  }
}

.kanban-card-inner {
  display: flex;
}

.kanban-card-left {
  width: 4px;
  border-radius: 8px 0 0 8px;
  flex-shrink: 0;
}

.kanban-card-content {
  flex: 1;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.kanban-card-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.kanban-avatar {
  width: 30px;
  height: 30px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
}

.kanban-info {
  flex: 1;
  min-width: 0;
}

.kanban-name {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.kanban-code {
  font-size: 11px;
  color: #c0c4cc;
}

.kanban-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 10px;
}

.kanban-members {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 10px;
  padding: 8px;
  background: #fff;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
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
  color: #606266;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kanban-member-more {
  font-size: 11px;
  color: #c0c4cc;
  padding: 2px 6px;
  background: #f5f5f5;
  border-radius: 4px;
}

.kanban-stations {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
  padding: 8px;
  background: #fff;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
}

.kanban-station {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #606266;

  .el-icon { color: #409eff; font-size: 12px; flex-shrink: 0; }
  span:first-of-type { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
}

.station-phone {
  font-size: 11px;
  color: #c0c4cc;
  flex-shrink: 0;
}

.kanban-station-more {
  font-size: 11px;
  color: #c0c4cc;
  padding-left: 20px;
}

.station-person-opt {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 0;

  small { font-size: 11px; color: #c0c4cc; }
}

.kanban-card-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
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
  color: #909399;

  .el-icon { font-size: 12px; }
}

.kanban-empty {
  text-align: center;
  padding: 24px;
  color: #c0c4cc;
  font-size: 13px;
}

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
}

.empty-icon {
  font-size: 48px;
  color: #dcdfe6;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: #606266;
  margin: 0 0 8px 0;
}

.empty-hint {
  font-size: 13px;
  color: #c0c4cc;
  margin: 0;
}

/* 骨架屏 */
.skeleton-card {
  height: 200px;
  background: #f0f0f0;
  border-radius: 10px;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* 抽屉 */
:deep(.project-drawer) {
  .el-drawer__body { padding: 0; }
}

.drawer-head {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 24px;
  border-bottom: 1px solid #ebeef5;
}

.drawer-avatar {
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
  height: calc(100vh - 140px);
}

.form-section {
  margin-bottom: 28px;
}

.form-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.station-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.station-num {
  width: 22px;
  height: 22px;
  background: $primary;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  flex-shrink: 0;
}

.drawer-foot {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #ebeef5;
  background: #fafafa;
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
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e8ecf0;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.22s ease;
  animation: fadeSlideIn 0.4s ease both;
  position: relative;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
    border-color: #d0d7de;

    .bento-hover-actions { opacity: 1; }
    .bento-avatar { transform: scale(1.05); }
  }

  &--active {
    border-left: 3px solid #22c55e;
  }
}

/* 左侧阶段色条 */
.bento-stage-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
}

/* 卡片顶部 */
.bento-top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 14px 12px;
  background: #fff;
  position: relative;
  border-bottom: 1px solid #f0f2f5;
}

.bento-top-info {
  flex: 1;
  min-width: 0;
}

.bento-avatar {
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
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
  transition: transform 0.2s;
}

.bento-hover-actions {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;

  .el-button {
    color: #9ca3af;
    padding: 4px 6px;

    &:hover { color: #3b82f6; }
  }

  .btn-delete:hover { color: #ef4444; }
}

/* 卡片主体 */
.bento-body {
  padding: 12px 14px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 项目名称 */
.bento-name {
  font-size: 14px;
  font-weight: 700;
  color: #111827;
  line-height: 1.3;
  letter-spacing: -0.2px;
}

.bento-code {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #9ca3af;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.2s;

  .el-icon { color: #d1d5db; transition: color 0.2s; }

  &:hover {
    color: #409eff;
    background: #ecf5ff;
    .el-icon { color: #409eff; }
  }
}

/* 核心指标 */
.bento-metrics {
  display: flex;
  align-items: center;
  background: #f9fafb;
  border-radius: 10px;
  padding: 10px 0;
  border: 1px solid #f0f2f5;
}

.bento-metric {
  flex: 1;
  text-align: center;
}

.bento-metric-num {
  font-size: 18px;
  font-weight: 800;
  color: #111827;
  line-height: 1.1;
  font-family: 'Manrope', sans-serif;
}

.bento-metric-label {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 2px;
}

.bento-metric-sep {
  width: 1px;
  height: 28px;
  background: #e5e7eb;
}

/* 描述 */
.bento-desc {
  font-size: 13px;
  color: #6b7280;
  line-height: 1.6;
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
  color: #6b7280;
  background: #f3f4f6;
  border-radius: 6px;
  padding: 4px 8px;

  .el-icon { color: #9ca3af; }
}

/* 解决方案 */
.bento-solution {
  border-left: 3px solid #3b82f6;
  padding-left: 12px;
}

.bento-solution-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: #9ca3af;
  margin-bottom: 4px;

  .el-icon { color: #3b82f6; }
}

.solution-person {
  color: #3b82f6;
  font-weight: 600;
}

.bento-solution-text {
  font-size: 13px;
  color: #374151;
  line-height: 1.5;
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
  color: #9ca3af;
  margin-bottom: 8px;

  .el-icon { color: #6b7280; }
}

.bento-team-count {
  background: #e5e7eb;
  color: #6b7280;
  border-radius: 10px;
  padding: 0 6px;
  font-size: 10px;
  font-weight: 700;
  line-height: 16px;
}

.bento-team-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.bento-member-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 3px 10px 3px 4px;
  cursor: default;
  transition: background 0.15s, border-color 0.15s;

  &:hover {
    background: #f0f4ff;
    border-color: #c7d2fe;
  }
}

.bento-member-avatar {
  width: 22px;
  height: 22px;
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
  color: #374151;
  font-weight: 500;
  white-space: nowrap;
}

.bento-member-overflow {
  display: flex;
  align-items: center;
  padding: 3px 10px;
  background: #f3f4f6;
  border: 1px dashed #d1d5db;
  border-radius: 20px;
  font-size: 11px;
  color: #9ca3af;
  font-weight: 600;
}

/* 驻场点 */
.bento-stations { }

.bento-stations-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: #9ca3af;
  margin-bottom: 8px;

  .el-icon { color: #6b7280; }
}

.bento-station-count {
  background: #e5e7eb;
  color: #6b7280;
  border-radius: 10px;
  padding: 0 6px;
  font-size: 10px;
  font-weight: 700;
  line-height: 16px;
}

.bento-station-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.bento-station-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 5px 10px;
  font-size: 12px;

  .el-icon { color: #6b7280; font-size: 12px; }
}

.station-location {
  color: #374151;
  font-weight: 500;
  flex: 1;
}

.station-person {
  color: #3b82f6;
  font-weight: 600;
  font-size: 11px;
}

.station-phone {
  color: #9ca3af;
  font-size: 11px;
}

/* 底部 */
.bento-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 8px;
  border-top: 1px solid #f3f4f6;
}

.bento-stage-tag {
  font-size: 10px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 4px;
}

.bento-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #d1d5db;

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
  .kpi-cards { grid-template-columns: repeat(3, 1fr); }
  .modal-info-grid { grid-template-columns: 1fr; }
  .modal-info-item--full { grid-column: 1; }
}

@media (max-width: 768px) {
  .kpi-cards { grid-template-columns: repeat(2, 1fr); }
  .filter-bar { flex-direction: column; align-items: stretch; }
  .filter-tip { margin-left: 0; }
}
</style>
