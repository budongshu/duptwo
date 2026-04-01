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
        <button class="view-tab" :class="{ active: viewMode === 'grid' }" @click="switchView('grid')">
          <el-icon><Grid /></el-icon>
          <span>{{ t('project.list.viewGrid') }}</span>
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'kanban' }" @click="switchView('kanban')">
          <el-icon><Menu /></el-icon>
          <span>{{ t('project.list.viewKanban') }}</span>
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'detail' }" @click="switchView('detail')">
          <el-icon><Tickets /></el-icon>
          <span>{{ t('project.list.viewBento') }}</span>
        </button>
      </div>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ t('project.list.addProject') }}
        </el-button>
      </div>
    </header>

    <!-- KPI 统计 -->
    <div class="kpi-cards" v-if="!loading && pagination.total > 0">
      <div class="kpi-card" @click="handleReset">
        <div class="kpi-num">{{ pagination.total }}</div>
        <div class="kpi-label">{{ t('project.list.allProjects') }}</div>
      </div>
      <div class="kpi-card kpi-card--green" @click="filterByStage('running')">
        <div class="kpi-num">{{ activeCount }}</div>
        <div class="kpi-label">{{ t('project.stage.running') }}</div>
      </div>
      <div class="kpi-card kpi-card--amber" @click="filterByRecords">
        <div class="kpi-num">{{ formatNumber(totalRecords) }}</div>
        <div class="kpi-label">{{ t('project.list.uploadRecords') }}</div>
      </div>
      <div class="kpi-card kpi-card--purple" @click="filterBySize">
        <div class="kpi-num">{{ formatBytes(totalDataSize) }}</div>
        <div class="kpi-label">{{ t('project.list.totalDataSize') }}</div>
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

    <!-- 网格视图 -->
    <div v-if="viewMode === 'grid'" class="project-grid" v-loading="loading">
      <div
        v-for="(project, idx) in tableData"
        :key="project.id"
        class="project-card"
        :style="{ animationDelay: `${idx * 0.03}s` }"
        @click="handleEdit(project)"
      >
        <!-- 卡片头部 -->
        <div class="card-head">
          <div class="card-avatar" :style="getProjectBgStyle(project.code)">
            {{ getInitials(project.name) }}
          </div>
          <div class="card-info">
            <div class="card-name">{{ project.name }}</div>
            <div class="card-code">{{ project.code }}</div>
          </div>
          <el-tag :type="project.status === 'active' ? 'success' : 'info'" size="small">
            {{ t('project.stage.' + project.stage) }}
          </el-tag>
        </div>

        <!-- 卡片描述 -->
        <div class="card-desc" v-if="project.description">{{ project.description }}</div>

        <!-- 卡片指标 -->
        <div class="card-metrics">
          <div class="metric-item">
            <el-icon><Document /></el-icon>
            <span>{{ formatNumber(project.recordCount || 0) }}</span>
            <small>{{ t('project.list.records') }}</small>
          </div>
          <div class="metric-item">
            <el-icon><Cpu /></el-icon>
            <span>{{ formatBytes(project.totalDataSize || 0) }}</span>
          </div>
          <div class="metric-item">
            <el-icon><User /></el-icon>
            <span>{{ getTeamCount(project) }}</span>
            <small>成员</small>
          </div>
        </div>

        <!-- 卡片底部 -->
        <div class="card-footer">
          <div class="card-avatar-list" v-if="getTeamMembers(project).length > 0">
            <div
              v-for="(member, mi) in getTeamMembers(project).slice(0, 5)"
              :key="mi"
              class="avatar-mini"
              :style="{ background: getAvatarColor(member) }"
              :title="member"
            >
              {{ member.charAt(0).toUpperCase() }}
            </div>
          </div>
          <span class="card-time">{{ formatDate(project.updatedAt) }}</span>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && tableData.length === 0" class="empty-state">
        <el-icon class="empty-icon"><Folder /></el-icon>
        <p class="empty-title">{{ t('project.list.noProjectData') }}</p>
        <p class="empty-hint">{{ t('project.list.noProjectHint') }}</p>
      </div>
    </div>

    <!-- 看板视图 -->
    <div v-if="viewMode === 'kanban' && !loading" class="kanban-board">
      <div v-for="col in kanbanColumns" :key="col.key" class="kanban-col">
        <div class="kanban-col-head">
          <span class="kanban-dot" :style="{ background: col.color }"></span>
          <span class="kanban-title">{{ t('project.stage.' + col.key) }}</span>
          <el-badge :value="kanbanData[col.key]?.length || 0" :max="99" />
        </div>
        <div class="kanban-col-body">
          <div
            v-for="p in (kanbanData[col.key] || [])"
            :key="p.id"
            class="kanban-card"
          >
            <!-- 卡片头部 -->
            <div class="kanban-card-head">
              <div class="kanban-avatar" :style="getProjectBgStyle(p.code)">
                {{ getInitials(p.name) }}
              </div>
              <div class="kanban-info">
                <div class="kanban-name" @click="showDetail(p)">{{ p.name }}</div>
                <div class="kanban-code">{{ p.code }}</div>
              </div>
              <el-tag :type="p.status === 'active' ? 'success' : 'info'" size="small">{{ p.status === 'active' ? t('common.enabled') : t('common.disabled') }}</el-tag>
            </div>

            <!-- 项目标签 -->
            <div class="kanban-tags" v-if="p.companyAddr || p.projectPeriod || getTeamMembers(p).length > 0">
              <el-tag type="warning" size="small" effect="plain">{{ t('project.stage.' + p.stage) }}</el-tag>
              <el-tag size="small" effect="plain" v-if="p.companyAddr">{{ p.companyAddr }}</el-tag>
              <el-tag size="small" effect="plain" v-if="p.projectPeriod">{{ p.projectPeriod }}</el-tag>
              <el-tag size="small" effect="plain" v-if="getTeamMembers(p).length > 0">
                <el-icon><User /></el-icon> {{ getTeamMembers(p).length }}人
              </el-tag>
              <el-tag size="small" effect="plain" v-if="p.onsiteStations && p.onsiteStations.length > 0">
                <el-icon><LocationInformation /></el-icon> {{ p.onsiteStations.length }}个驻场点
              </el-tag>
            </div>

            <!-- 团队成员列表 -->
            <div class="kanban-members" v-if="getTeamMembers(p).length > 0">
              <div class="kanban-member" v-for="member in getTeamMembers(p).slice(0, 5)" :key="member">
                <div class="kanban-member-avatar" :style="{ background: getAvatarColor(member) }">
                  {{ member.charAt(0).toUpperCase() }}
                </div>
                <span class="kanban-member-name">{{ member }}</span>
              </div>
              <span v-if="getTeamMembers(p).length > 5" class="kanban-member-more">+{{ getTeamMembers(p).length - 5 }}人</span>
            </div>

            <!-- 驻场点列表 -->
            <div class="kanban-stations" v-if="p.onsiteStations && p.onsiteStations.length > 0">
              <div class="kanban-station" v-for="(s, si) in p.onsiteStations.slice(0, 3)" :key="si">
                <el-icon><LocationInformation /></el-icon>
                <span>{{ s.location || s.person || '驻场点' }}</span>
                <span v-if="s.phone" class="station-phone">{{ s.phone }}</span>
              </div>
              <span v-if="p.onsiteStations.length > 3" class="kanban-station-more">+{{ p.onsiteStations.length - 3 }}个驻场点</span>
            </div>

            <!-- 底部指标和操作 -->
            <div class="kanban-card-foot">
              <span class="kanban-metric">
                <el-icon><Document /></el-icon> {{ formatNumber(p.recordCount || 0) }}
              </span>
              <span class="kanban-metric">
                <el-icon><Cpu /></el-icon> {{ formatBytes(p.totalDataSize || 0) }}
              </span>
              <div class="kanban-actions" @click.stop>
                <TableActions :actions="[
                  { key: 'switch-bento', label: t('common.view'), type: 'default' },
                  { key: 'delete', label: t('common.delete'), type: 'danger' }
                ]" @action="(key) => key === 'switch-bento' ? (switchView('detail'), detailVisible = false) : handleDelete(p)" />
              </div>
            </div>
          </div>
          <div v-if="!kanbanData[col.key]?.length" class="kanban-empty">{{ t('project.list.noData') }}</div>
        </div>
      </div>
    </div>

    <!-- 详情视图 -->
    <!-- 瀑布视图 -->
    <div v-if="viewMode === 'detail'" class="bento-view" v-loading="loading">
      <div class="bento-grid">
        <div
          v-for="(project, idx) in tableData"
          :key="project.id"
          class="bento-card"
          :style="{ animationDelay: `${idx * 0.04}s` }"
          @click="openEditFromDetail(project)"
        >
          <!-- 卡片顶部 -->
          <div class="bento-top">
            <div class="bento-avatar" :style="{ background: getProjectColor(project.code) }">
              {{ getInitials(project.name) }}
            </div>
            <div class="bento-top-badges">
              <el-tag :type="project.status === 'active' ? 'success' : 'info'" size="small">
                {{ project.status === 'active' ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
              <el-tag type="warning" size="small" effect="plain">
                {{ t('project.stage.' + project.stage) }}
              </el-tag>
            </div>
            <!-- 悬浮编辑按钮 -->
            <div class="bento-actions" @click.stop>
              <el-button size="small" text @click="openEditFromDetail(project)">
                <el-icon><Edit /></el-icon>
                {{ t('common.edit') }}
              </el-button>
            </div>
          </div>

          <!-- 卡片主体 -->
          <div class="bento-body">
            <!-- 项目名 + 编号 -->
            <div class="bento-header">
              <div class="bento-name">{{ project.name }}</div>
              <div class="bento-code">
                <el-icon><Link /></el-icon>
                {{ project.code }}
              </div>
            </div>

            <!-- 核心指标 -->
            <div class="bento-metrics">
              <div class="bento-metric">
                <div class="bento-metric-num">{{ formatNumber(project.recordCount || 0) }}</div>
                <div class="bento-metric-label">{{ t('project.list.records') }}</div>
              </div>
              <div class="bento-metric-sep"></div>
              <div class="bento-metric">
                <div class="bento-metric-num">{{ formatBytes(project.totalDataSize || 0) }}</div>
                <div class="bento-metric-label">{{ t('project.list.size') }}</div>
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
              <span class="bento-time">
                <el-icon><Timer /></el-icon>
                更新于 {{ formatDate(project.updatedAt) }}
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
    <div class="pagination-wrapper" v-if="pagination.total > 0 && viewMode === 'grid'">
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
        <el-button text @click="drawerVisible = false"><el-icon><Close /></el-icon></el-button>
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
              <el-button type="danger" size="small" text @click="removeStation(idx)">
                <el-icon><Delete /></el-icon>
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
        <el-button type="primary" :loading="submitting" @click="confirmSubmit">
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
        <el-button type="primary" @click="openEditFromDetail">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          {{ t('common.edit') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Grid, Menu, Plus, Search, Document, Cpu, User,
  Folder, Delete, Close, Location, Clock, Connection, LocationInformation, Tickets, Phone, Edit, Timer, Link
} from '@element-plus/icons-vue'
import { ProjectApi, type Project, type CreateProjectReq, type UpdateProjectReq, type OnSiteStation } from '@/api/project'
import { PersonnelApi, type Personnel } from '@/api/personnel'
import TableActions from '@/components/TableActions.vue'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Project[]>([])
const drawerVisible = ref(false)
const dialogPreviewInitials = ref('?')
const viewMode = ref<'kanban' | 'grid' | 'detail'>('grid')
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

const activeCount = computed(() => tableData.value.filter(p => p.status === 'active').length)
const totalRecords = computed(() => tableData.value.reduce((sum, p) => sum + (p.recordCount || 0), 0))
const totalDataSize = computed(() => tableData.value.reduce((sum, p) => sum + (p.totalDataSize || 0), 0))

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

const switchView = (mode: 'kanban' | 'grid' | 'detail') => {
  viewMode.value = mode
  if (mode === 'kanban') loadKanbanData()
  else if (mode === 'grid') loadData()
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
  loadData()
}
const filterByRecords = () => {
  tableData.value = [...tableData.value].sort((a, b) => (b.recordCount || 0) - (a.recordCount || 0))
}
const filterBySize = () => {
  tableData.value = [...tableData.value].sort((a, b) => (b.totalDataSize || 0) - (a.totalDataSize || 0))
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
onMounted(() => { loadData() })
</script>

<script lang="ts">
export default { name: 'ProjectList' }
</script>

<style scoped lang="scss">
$primary: #409eff;
$success: #67c23a;
$warning: #e6a23c;
$danger: #f56c6c;

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
.kpi-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.kpi-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  border: 1px solid #ebeef5;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;

  &:hover {
    box-shadow: 0 2px 12px rgba(0,0,0,0.08);
    transform: translateY(-1px);
  }
}

.kpi-num {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 6px;
}

.kpi-label {
  font-size: 13px;
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
    border-color: $primary;
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
  width: 280px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #ebeef5;
}

.kanban-col-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.kanban-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.kanban-title {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.kanban-col-body {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: calc(100vh - 400px);
  overflow-y: auto;
}

.kanban-card {
  background: #f9fafb;
  border-radius: 8px;
  padding: 14px;
  border: 1px solid #ebeef5;
  transition: border-color 0.2s, box-shadow 0.2s;

  &:hover {
    border-color: $primary;
    box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  }
}

.kanban-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.kanban-avatar {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
  font-family: 'Manrope', sans-serif;
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
  background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
  background-size: 200% 100%;
  border-radius: 10px;
  animation: shimmer 1.2s infinite;
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
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 18px;
  align-items: start;
}

.bento-card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e8ecf0;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.22s ease, box-shadow 0.22s ease, border-color 0.22s ease;
  animation: fadeSlideIn 0.4s ease both;
  position: relative;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.09);
    border-color: #d0d7de;

    .bento-actions { opacity: 1; }
    .bento-avatar { transform: scale(1.05); }
  }
}

/* 卡片顶部：白底 + 左侧大彩色头像 */
.bento-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 18px 14px 12px;
  background: #fff;
  position: relative;
  border-bottom: 1px solid #f0f2f5;
}

.bento-avatar {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 800;
  color: #fff;
  font-family: 'Manrope', sans-serif;
  letter-spacing: -0.5px;
  flex-shrink: 0;
  box-shadow: 0 3px 10px rgba(0,0,0,0.18);
  transition: transform 0.2s;
}

.bento-top-badges {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  padding-left: 14px;
  flex: 1;
}

.bento-actions {
  position: absolute;
  top: 12px;
  right: 10px;
  opacity: 0;
  transition: opacity 0.2s;

  .el-button {
    font-size: 12px;
    color: #6b7280;
    padding: 4px 8px;

    &:hover { color: #3b82f6; background: #f0f4ff; }
  }
}

/* 卡片主体 */
.bento-body {
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 项目名称 */
.bento-name {
  font-size: 15px;
  font-weight: 700;
  color: #111827;
  line-height: 1.3;
  margin-bottom: 3px;
  letter-spacing: -0.3px;
}

.bento-code {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;

  .el-icon { color: #d1d5db; }
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
  padding-top: 4px;
  border-top: 1px solid #f3f4f6;
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
  .kpi-cards { grid-template-columns: repeat(2, 1fr); }
  .modal-info-grid { grid-template-columns: 1fr; }
  .modal-info-item--full { grid-column: 1; }
}

@media (max-width: 768px) {
  .kpi-cards { grid-template-columns: 1fr 1fr; }
  .filter-bar { flex-direction: column; align-items: stretch; }
  .filter-tip { margin-left: 0; }
}
</style>
