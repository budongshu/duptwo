<template>
  <div class="network-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('project.network.title') }}</h1>
        <span class="page-subtitle">{{ t('project.network.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <el-button @click="resetGraph">
          <el-icon><Refresh /></el-icon>
          {{ t('project.network.resetView') }}
        </el-button>
        <el-select v-model="filterRole" :placeholder="t('project.network.filterRole')" clearable size="small" style="width: 140px">
          <el-option :label="t('project.network.allRoles')" value="" />
          <el-option :label="t('project.role.projectPerson')" value="projectPerson" />
          <el-option :label="t('project.role.opsPerson')" value="opsPerson" />
          <el-option :label="t('project.role.opsStaffPerson')" value="opsStaffPerson" />
          <el-option :label="t('project.role.developerPerson')" value="developerPerson" />
          <el-option :label="t('project.role.testerPerson')" value="testerPerson" />
          <el-option :label="t('project.role.businessPerson')" value="businessPerson" />
          <el-option :label="t('project.role.compliancePerson')" value="compliancePerson" />
          <el-option :label="t('project.role.solutionPerson')" value="solutionPerson" />
        </el-select>
        <el-select v-model="filterStage" :placeholder="t('project.network.filterStage')" clearable size="small" style="width: 140px">
          <el-option :label="t('project.network.allStages')" value="" />
          <el-option v-for="s in stageOptions" :key="s.value" :label="t('project.stage.' + s.value)" :value="s.value" />
        </el-select>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ t('project.list.addProject') }}
        </el-button>
      </div>
    </header>

    <!-- 关系网络图 -->
    <div class="graph-container" ref="graphContainer">
      <svg ref="svgRef" class="network-svg"></svg>

      <!-- 加载中 -->
      <div v-if="loading" class="graph-loading">
        <el-icon class="is-loading"><Loading /></el-icon>
        <span>{{ t('project.network.loading') }}</span>
      </div>

      <!-- 图例 -->
      <div class="graph-legend">
        <div class="legend-item">
          <span class="legend-node legend-node--person"></span>
          <span>{{ t('project.network.legendPerson') }}</span>
        </div>
        <div class="legend-item">
          <span class="legend-node legend-node--project"></span>
          <span>{{ t('project.network.legendProject') }}</span>
        </div>
        <div class="legend-item">
          <span class="legend-line"></span>
          <span>{{ t('project.network.legendRelation') }}</span>
        </div>
      </div>

      <!-- 工作量统计 -->
      <div class="workload-panel">
        <div class="workload-title">{{ t('project.network.workloadTitle') }}</div>
        <div class="workload-list">
          <div
            v-for="item in workloadList"
            :key="item.name"
            class="workload-item"
            :class="{ 'workload-item--highlight': highlightedNode === item.name }"
            @click="highlightNode(item.name)"
          >
            <div class="workload-avatar" :style="{ background: getAvatarColor(item.name) }">
              {{ item.name.charAt(0).toUpperCase() }}
            </div>
            <div class="workload-info">
              <div class="workload-name">{{ item.name }}</div>
              <div class="workload-count">
                <span class="workload-num">{{ item.count }}</span>
                <span class="workload-unit">{{ t('project.network.projects') }}</span>
              </div>
            </div>
            <div class="workload-bar">
              <div class="workload-bar-fill" :style="{ width: (item.count / maxWorkload * 100) + '%', background: getAvatarColor(item.name) }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 节点详情弹窗 -->
    <el-dialog v-model="nodeDialogVisible" :title="nodeDialogData.title" width="480px" :destroy-on-close="true">
      <div class="node-detail" v-if="nodeDialogData.type === 'person'">
        <div class="node-detail-header">
          <div class="node-detail-avatar" :style="{ background: getAvatarColor(nodeDialogData.name) }">
            {{ nodeDialogData.name?.charAt(0).toUpperCase() }}
          </div>
          <div class="node-detail-info">
            <h3>{{ nodeDialogData.name }}</h3>
            <div class="node-detail-meta">{{ t('project.network.memberOfProjects', { count: nodeDialogData.links?.length || 0 }) }}</div>
          </div>
        </div>
        <div class="node-detail-projects">
          <div class="node-detail-project" v-for="p in nodeDialogData.projects" :key="p.id" @click="handleEdit(p)">
            <div class="ndp-avatar" :style="{ background: getProjectColor(p.code) }">
              {{ getInitials(p.name) }}
            </div>
            <div class="ndp-info">
              <div class="ndp-name">{{ p.name }}</div>
              <div class="ndp-meta">
                <el-tag size="small" :type="p.status === 'active' ? 'success' : 'info'">{{ p.status === 'active' ? t('common.enabled') : t('common.disabled') }}</el-tag>
                <el-tag size="small" type="warning">{{ t('project.stage.' + p.stage) }}</el-tag>
              </div>
            </div>
            <el-icon class="ndp-arrow"><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
      <div class="node-detail" v-if="nodeDialogData.type === 'project'">
        <div class="node-detail-header">
          <div class="node-detail-avatar node-detail-avatar--project" :style="{ background: getProjectColor(nodeDialogData.code) }">
            {{ getInitials(nodeDialogData.name || '') }}
          </div>
          <div class="node-detail-info">
            <h3>{{ nodeDialogData.name }}</h3>
            <div class="node-detail-meta">{{ nodeDialogData.code }}</div>
          </div>
        </div>
        <div class="node-detail-stats">
          <div class="nds-item">
            <span class="nds-num">{{ nodeDialogData.recordCount || 0 }}</span>
            <span class="nds-label">{{ t('project.list.records') }}</span>
          </div>
          <div class="nds-item">
            <span class="nds-num">{{ formatBytes(nodeDialogData.totalDataSize || 0) }}</span>
            <span class="nds-label">{{ t('project.list.size') }}</span>
          </div>
          <div class="nds-item">
            <span class="nds-num">{{ nodeDialogData.links?.length || 0 }}</span>
            <span class="nds-label">{{ t('project.network.teamMembers') }}</span>
          </div>
        </div>
        <div class="node-detail-members">
          <div class="node-detail-member" v-for="m in nodeDialogData.members" :key="m" @click="highlightNode(m)">
            <div class="ndm-avatar" :style="{ background: getAvatarColor(m) }">{{ m.charAt(0).toUpperCase() }}</div>
            <span>{{ m }}</span>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="nodeDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="openProjectFromNode">{{ t('project.network.viewProject') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Plus, Loading, ArrowRight } from '@element-plus/icons-vue'
import * as d3 from 'd3'
import { ProjectApi, type Project } from '@/api/project'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const svgRef = ref<SVGSVGElement>()
const graphContainer = ref<HTMLDivElement>()
const filterRole = ref('')
const filterStage = ref('')
const nodeDialogVisible = ref(false)
const nodeDialogData = reactive<any>({
  type: '', name: '', code: '', title: '',
  links: [] as any[], projects: [] as Project[],
  members: [] as string[], recordCount: 0, totalDataSize: 0,
})
const highlightedNode = ref('')
const allProjects = ref<Project[]>([])

// D3 state
let simulation: any = null
let svg: any = null
let g: any = null
let zoom: any = null

const stageOptions = [
  { value: 'planning' }, { value: 'designing' }, { value: 'deploying' },
  { value: 'running' }, { value: 'paused' },
]

// Graph data
interface GraphNode {
  id: string
  name: string
  type: 'person' | 'project'
  code?: string
  role?: string
  stage?: string
  status?: string
  recordCount?: number
  totalDataSize?: number
  count: number
  fx?: number | null
  fy?: number | null
  x?: number
  y?: number
  vx?: number
  vy?: number
  index?: number
}

interface GraphLink {
  source: string | GraphNode
  target: string | GraphNode
  role?: string
}

const nodes = ref<GraphNode[]>([])
const links = ref<GraphLink[]>([])

const workloadList = computed(() => {
  return nodes.value
    .filter(n => n.type === 'person' && n.count > 0)
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

const maxWorkload = computed(() => {
  return Math.max(...workloadList.value.map(w => w.count), 1)
})

const getInitials = (name: string) => {
  if (!name) return '??'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

const projectColors = ['#4a6fa5', '#6b5b95', '#4a7c59', '#b87333', '#8b6f5b', '#5b7b8c', '#6b7c8c', '#7a8471']
const getProjectColor = (code: string) => {
  if (!code) return projectColors[0]
  let hash = 0
  for (let i = 0; i < code.length; i++) hash = ((hash << 5) - hash) + code.charCodeAt(i)
  return projectColors[Math.abs(hash) % projectColors.length]
}

const avatarColors = ['#6b5b95', '#4a7c59', '#b87333', '#8b6f5b', '#5b7b8c', '#6b7c8c', '#7a8471', '#a67c52']
const getAvatarColor = (name: string) => {
  if (!name) return avatarColors[0]
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = ((hash << 5) - hash) + name.charCodeAt(i)
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

const formatBytes = (bytes: number): string => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

const buildGraph = (projects: Project[]) => {
  const nodeMap = new Map<string, GraphNode>()
  const linkList: GraphLink[] = []

  // Build person nodes
  const roleFields: (keyof Project)[] = [
    'projectPerson', 'opsPerson', 'opsStaffPerson', 'developerPerson',
    'testerPerson', 'businessPerson', 'compliancePerson', 'solutionPerson'
  ]

  for (const p of projects) {
    // Filter by stage
    if (filterStage.value && p.stage !== filterStage.value) continue

    const personCount: Record<string, number> = {}
    for (const role of roleFields) {
      const name = p[role] as string
      if (!name) continue
      personCount[name] = (personCount[name] || 0) + 1
    }

    for (const [name, count] of Object.entries(personCount)) {
      if (!nodeMap.has(name)) {
        nodeMap.set(name, { id: name, name, type: 'person', count: 0 })
      }
      nodeMap.get(name)!.count += count
    }
  }

  // Filter person nodes by role
  if (filterRole.value) {
    const filtered = new Map<string, GraphNode>()
    for (const p of projects) {
      if (filterStage.value && p.stage !== filterStage.value) continue
      const name = p[filterRole.value as keyof Project] as string
      if (name && nodeMap.has(name)) {
        filtered.set(name, nodeMap.get(name)!)
      }
    }
    nodeMap.clear()
    filtered.forEach((v, k) => nodeMap.set(k, v))
  }

  // Build project nodes and links
  for (const p of projects) {
    if (filterStage.value && p.stage !== filterStage.value) continue

    if (!nodeMap.has('proj_' + p.id)) {
      nodeMap.set('proj_' + p.id, {
        id: 'proj_' + p.id,
        name: p.name,
        type: 'project',
        code: p.code,
        stage: p.stage,
        status: p.status,
        recordCount: p.recordCount,
        totalDataSize: p.totalDataSize,
        count: 0,
      })
    }

    for (const role of roleFields) {
      const name = p[role] as string
      if (!name) continue
      if (filterRole.value && role !== filterRole.value) continue
      if (!nodeMap.has(name)) continue
      linkList.push({ source: name, target: 'proj_' + p.id, role })
    }
  }

  nodes.value = Array.from(nodeMap.values())
  links.value = linkList
}

const renderGraph = () => {
  if (!svgRef.value || !graphContainer.value) return

  const container = graphContainer.value
  const width = container.clientWidth
  const height = container.clientHeight

  // Clear previous
  d3.select(svgRef.value).selectAll('*').remove()

  svg = d3.select(svgRef.value)
    .attr('width', width)
    .attr('height', height)

  // Defs for arrow markers
  const defs = svg.append('defs')
  defs.append('marker')
    .attr('id', 'arrow')
    .attr('viewBox', '0 -5 10 10')
    .attr('refX', 20)
    .attr('refY', 0)
    .attr('markerWidth', 6)
    .attr('markerHeight', 6)
    .attr('orient', 'auto')
    .append('path')
    .attr('d', 'M0,-5L10,0L0,5')
    .attr('fill', '#c0c4cc')

  // Zoom
  zoom = d3.zoom()
    .scaleExtent([0.2, 4])
    .on('zoom', (event: any) => {
      g.attr('transform', event.transform)
    })

  svg.call(zoom as any)

  g = svg.append('g')

  // Simulation
  simulation = d3.forceSimulation(nodes.value as any)
    .force('link', d3.forceLink(links.value as any)
      .id((d: any) => d.id)
      .distance(d => d.source.type === 'person' ? 120 : 180)
      .strength(0.5))
    .force('charge', d3.forceManyBody().strength(-400))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collision', d3.forceCollide().radius(d => getNodeRadius(d) + 10))

  // Links
  const link = g.append('g')
    .attr('class', 'links')
    .selectAll('line')
    .data(links.value)
    .join('line')
    .attr('stroke', '#dcdfe6')
    .attr('stroke-width', 1.5)
    .attr('stroke-opacity', 0.6)

  // Node groups
  const node = g.append('g')
    .attr('class', 'nodes')
    .selectAll('g')
    .data(nodes.value)
    .join('g')
    .attr('class', 'node-group')
    .call(d3.drag<any, any>()
      .on('start', dragStarted)
      .on('drag', dragged)
      .on('end', dragEnded) as any)

  // Node shapes
  node.each(function(d: any) {
    const el = d3.select(this)
    if (d.type === 'person') {
      // Circle for person
      el.append('circle')
        .attr('r', getNodeRadius(d))
        .attr('fill', getAvatarColor(d.name))
        .attr('stroke', '#fff')
        .attr('stroke-width', 3)
        .attr('cursor', 'pointer')

      el.append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', getNodeRadius(d) + 16)
        .attr('font-size', '12px')
        .attr('fill', '#606266')
        .attr('pointer-events', 'none')
        .text(d.name)
    } else {
      // Square for project
      const size = getNodeRadius(d) * 2
      el.append('rect')
        .attr('x', -size / 2)
        .attr('y', -size / 2)
        .attr('width', size)
        .attr('height', size)
        .attr('rx', 8)
        .attr('fill', getProjectColor(d.code || ''))
        .attr('stroke', '#fff')
        .attr('stroke-width', 3)
        .attr('cursor', 'pointer')

      el.append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', size / 2 + 16)
        .attr('font-size', '12px')
        .attr('fill', '#606266')
        .attr('pointer-events', 'none')
        .text(d.name)
    }
  })

  // Click handler
  node.on('click', (event: any, d: any) => {
    event.stopPropagation()
    showNodeDetail(d)
  })

  // Highlight on hover
  node.on('mouseenter', (event: any, d: any) => {
    highlightConnections(d, link, node)
  })
  node.on('mouseleave', () => {
    link.attr('stroke', '#dcdfe6').attr('stroke-opacity', 0.6)
    link.attr('stroke-width', 1.5)
    node.select('circle,rect').attr('opacity', 1)
  })

  simulation.on('tick', () => {
    link
      .attr('x1', (d: any) => d.source.x)
      .attr('y1', (d: any) => d.source.y)
      .attr('x2', (d: any) => d.target.x)
      .attr('y2', (d: any) => d.target.y)

    node.attr('transform', (d: any) => `translate(${d.x},${d.y})`)
  })
}

const getNodeRadius = (d: any) => {
  if (d.type === 'person') {
    return Math.max(22, Math.min(40, 18 + d.count * 5))
  }
  // Project node: size by team
  return Math.max(28, Math.min(50, 25 + d.count * 3))
}

const highlightConnections = (d: any, link: any, node: any) => {
  const connected = new Set<string>()
  connected.add(d.id)
  link.each((l: any) => {
    if (l.source.id === d.id) connected.add(l.target.id)
    if (l.target.id === d.id) connected.add(l.source.id)
  })

  link.attr('stroke', (l: any) =>
    l.source.id === d.id || l.target.id === d.id ? '#409eff' : '#dcdfe6')
    .attr('stroke-width', (l: any) =>
      l.source.id === d.id || l.target.id === d.id ? 2.5 : 1.5)
    .attr('stroke-opacity', (l: any) =>
      l.source.id === d.id || l.target.id === d.id ? 1 : 0.3)

  node.select('circle,rect').attr('opacity', (n: any) => connected.has(n.id) ? 1 : 0.3)
}

const dragStarted = (event: any, d: any) => {
  if (!event.active) simulation.alphaTarget(0.3).restart()
  d.fx = d.x
  d.fy = d.y
}

const dragged = (event: any, d: any) => {
  d.fx = event.x
  d.fy = event.y
}

const dragEnded = (event: any, d: any) => {
  if (!event.active) simulation.alphaTarget(0)
  d.fx = null
  d.fy = null
}

const showNodeDetail = (d: any) => {
  nodeDialogData.type = d.type
  nodeDialogData.name = d.name
  nodeDialogData.code = d.code
  nodeDialogData.title = d.type === 'person' ? d.name : d.name
  nodeDialogData.stage = d.stage
  nodeDialogData.status = d.status
  nodeDialogData.recordCount = d.recordCount
  nodeDialogData.totalDataSize = d.totalDataSize

  if (d.type === 'person') {
    const projects = allProjects.value.filter(p => {
      if (filterStage.value && p.stage !== filterStage.value) return false
      const roleFields = ['projectPerson', 'opsPerson', 'opsStaffPerson', 'developerPerson', 'testerPerson', 'businessPerson', 'compliancePerson', 'solutionPerson']
      for (const role of roleFields) {
        if (p[role] === d.name) return true
      }
      return false
    })
    nodeDialogData.projects = projects
    nodeDialogData.links = projects.map(p => p.id)
  } else {
    const project = allProjects.value.find(p => 'proj_' + p.id === d.id)
    if (project) {
      const members: string[] = []
      const roleFields = ['projectPerson', 'opsPerson', 'opsStaffPerson', 'developerPerson', 'testerPerson', 'businessPerson', 'compliancePerson', 'solutionPerson']
      for (const role of roleFields) {
        const name = project[role] as string
        if (name) members.push(name)
      }
      nodeDialogData.members = members
      nodeDialogData.links = members
    }
  }
  nodeDialogVisible.value = true
}

const highlightNode = (name: string) => {
  highlightedNode.value = highlightedNode.value === name ? '' : name
  if (!simulation) return

  if (highlightedNode.value) {
    const d = nodes.value.find(n => n.name === name)
    if (d) {
      const connected = new Set<string>([d.id])
      links.value.forEach(l => {
        const sid = typeof l.source === 'string' ? l.source : (l.source as any).id
        const tid = typeof l.target === 'string' ? l.target : (l.target as any).id
        if (sid === d.id) connected.add(tid)
        if (tid === d.id) connected.add(sid)
      })
      svg.selectAll('.node-group').attr('opacity', (n: any) => connected.has(n.id) ? 1 : 0.15)
      svg.selectAll('.links line').attr('stroke-opacity', (l: any) => {
        const sid = typeof l.source === 'string' ? l.source : l.source.id
        const tid = typeof l.target === 'string' ? l.target : l.target.id
        return sid === d.id || tid === d.id ? 1 : 0.1
      }).attr('stroke-width', (l: any) => {
        const sid = typeof l.source === 'string' ? l.source : l.source.id
        const tid = typeof l.target === 'string' ? l.target : l.target.id
        return sid === d.id || tid === d.id ? 3 : 1
      })
    }
  } else {
    svg.selectAll('.node-group').attr('opacity', 1)
    svg.selectAll('.links line').attr('stroke-opacity', 0.6).attr('stroke-width', 1.5)
  }
}

const resetGraph = () => {
  highlightedNode.value = ''
  if (!svg || !zoom) return
  svg.transition().duration(500).call(zoom.transform as any, d3.zoomIdentity)
  svg.selectAll('.node-group').attr('opacity', 1)
  svg.selectAll('.links line').attr('stroke-opacity', 0.6).attr('stroke-width', 1.5)
}

const openProjectFromNode = () => {
  if (nodeDialogData.type === 'project') {
    const p = allProjects.value.find(p => 'proj_' + p.id === nodeDialogData.code)
    if (p) handleEdit(p)
  } else {
    router.push('/projects')
  }
}

const handleEdit = (p: Project) => {
  nodeDialogVisible.value = false
  router.push('/projects')
}

const handleCreate = () => {
  router.push('/projects')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await ProjectApi.list({ pageSize: 500 })
    const projects: Project[] = res.data?.items || res.data || []
    allProjects.value = projects
    buildGraph(projects)
    await nextTick()
    renderGraph()
  } catch {
    ElMessage.error(t('common.loadError'))
  } finally {
    loading.value = false
  }
}

// Re-render on filter change
const rebuildGraph = () => {
  buildGraph(allProjects.value)
  simulation?.stop()
  renderGraph()
}

const handleResize = () => {
  if (!graphContainer.value || !svgRef.value) return
  const width = graphContainer.value.clientWidth
  const height = graphContainer.value.clientHeight
  svg.attr('width', width).attr('height', height)
  simulation?.force('center', d3.forceCenter(width / 2, height / 2))
  simulation?.alpha(0.3).restart()
}

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  loadData()
  resizeObserver = new ResizeObserver(handleResize)
  if (graphContainer.value) resizeObserver.observe(graphContainer.value)
})

onUnmounted(() => {
  simulation?.stop()
  resizeObserver?.disconnect()
})

// Watch filters
watch([filterRole, filterStage], rebuildGraph)
</script>

<script lang="ts">
export default { name: 'ProjectNetwork' }
</script>

<style scoped lang="scss">
$primary: #6b5b95;

.network-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #1c1917;
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #292524;
  border-bottom: 1px solid #44403c;
  flex-shrink: 0;
  gap: 16px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #fafaf9;
  margin: 0;
}

.page-subtitle {
  font-size: 12px;
  color: #78716c;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

/* Graph */
.graph-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.network-svg {
  width: 100%;
  height: 100%;
  display: block;
  background: radial-gradient(ellipse at center, #292524 0%, #1c1917 60%);
}

.graph-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #78716c;
  font-size: 14px;
}

/* Legend */
.graph-legend {
  position: absolute;
  top: 16px;
  left: 16px;
  background: rgba(41, 37, 36, 0.9);
  border: 1px solid #44403c;
  border-radius: 12px;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  backdrop-filter: blur(8px);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: #a8a29e;
}

.legend-node {
  display: inline-block;
  flex-shrink: 0;
}

.legend-node--person {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #6b5b95;
}

.legend-node--project {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  background: #4a7c59;
}

.legend-line {
  width: 20px;
  height: 2px;
  background: #78716c;
  flex-shrink: 0;
}

/* Workload Panel */
.workload-panel {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 240px;
  background: rgba(41, 37, 36, 0.9);
  border: 1px solid #44403c;
  border-radius: 12px;
  padding: 14px;
  backdrop-filter: blur(8px);
  max-height: calc(100vh - 140px);
  overflow-y: auto;
}

.workload-title {
  font-size: 11px;
  font-weight: 600;
  color: #a8a29e;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
}

.workload-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.workload-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover { background: rgba(255,255,255,0.05); }
  &--highlight { background: rgba(107, 91, 149, 0.2); }
}

.workload-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.workload-info {
  flex: 1;
  min-width: 0;
}

.workload-name {
  font-size: 12px;
  color: #e7e5e4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}

.workload-count {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.workload-num {
  font-size: 14px;
  font-weight: 700;
  color: #f1f5f9;
}

.workload-unit {
  font-size: 10px;
  color: #78716c;
}

.workload-bar {
  width: 40px;
  height: 4px;
  background: #44403c;
  border-radius: 2px;
  flex-shrink: 0;
  overflow: hidden;
}

.workload-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s;
}

/* Node Detail Dialog */
.node-detail {
  padding: 4px 0;
}

.node-detail-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 20px;
}

.node-detail-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

  &--project { border-radius: 12px; }
}

.node-detail-info {
  flex: 1;

  h3 {
    font-size: 16px;
    font-weight: 700;
    color: #1c1917;
    margin: 0 0 4px 0;
  }
}

.node-detail-meta {
  font-size: 12px;
  color: #78716c;
}

.node-detail-projects {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.node-detail-project {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #fafaf9;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid #e8e5e1;

  &:hover { background: #f5f5f4; border-color: #d4d0c8; }
}

.ndp-avatar {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ndp-info { flex: 1; }

.ndp-name {
  font-size: 13px;
  font-weight: 600;
  color: #1c1917;
  margin-bottom: 4px;
}

.ndp-meta { display: flex; gap: 4px; }

.ndp-arrow { color: #d4d0c8; }

.node-detail-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  padding: 14px;
  background: #fafaf9;
  border-radius: 12px;
}

.nds-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.nds-num {
  font-size: 18px;
  font-weight: 700;
  color: #1c1917;
}

.nds-label {
  font-size: 11px;
  color: #78716c;
}

.node-detail-members {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.node-detail-member {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px 4px 4px;
  background: #f5f5f4;
  border: 1px solid #e8e5e1;
  border-radius: 20px;
  cursor: pointer;
  font-size: 12px;
  color: #44403c;
  transition: all 0.2s;

  &:hover { background: #f0ede8; border-color: #d4d0c8; }
}

.ndm-avatar {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: #fff;
}
</style>
