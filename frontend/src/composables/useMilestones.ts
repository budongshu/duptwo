import { ref, computed, onMounted } from 'vue'

// ========== Milestone Definitions ==========
export interface Milestone {
  key: string
  count: number
  name: string
  nameEn: string
  emoji: string
  description: string
  color: string
  unlocked: boolean
}

const STORAGE_KEY = 'datareg_milestones'

const MILESTONE_DEFS = [
  { key: 'export_1',  count: 1,  name: '初出茅庐',   nameEn: 'Data Pioneer',   emoji: '🌱', description: '完成了首次数据导出', color: '#22c55e' },
  { key: 'export_5',  count: 5,  name: '小有成就',   nameEn: 'Data Collector', emoji: '📦', description: '累计导出5次数据', color: '#3b82f6' },
  { key: 'export_10', count: 10, name: '数据搬运工', nameEn: 'Data Porter',    emoji: '🏋️', description: '累计导出10次数据', color: '#f59e0b' },
  { key: 'export_20', count: 20, name: '数据达人',   nameEn: 'Data Pro',       emoji: '💫', description: '累计导出20次数据', color: '#8b5cf6' },
  { key: 'export_50', count: 50, name: '数据大师',   nameEn: 'Data Master',    emoji: '🏆', description: '累计导出50次数据', color: '#ef4444' },
  { key: 'export_100',count: 100,name: '数据领主',   nameEn: 'Data Overlord',  emoji: '👑', description: '累计导出100次数据', color: '#ec4899' },
]

// ========== State ==========
const exportCount = ref(0)
const unlockedMilestones = ref<string[]>([])
const showBadge = ref(false)
const currentBadge = ref<Milestone | null>(null)
let badgeTimer: ReturnType<typeof setTimeout> | null = null

// ========== Load from localStorage ==========
const loadState = () => {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const data = JSON.parse(saved)
      exportCount.value = data.exportCount || 0
      unlockedMilestones.value = data.unlockedMilestones || []
    }
  } catch {}
}

const saveState = () => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({
    exportCount: exportCount.value,
    unlockedMilestones: unlockedMilestones.value,
  }))
}

// ========== Track Export ==========
const trackExport = () => {
  exportCount.value++

  // Check for newly unlocked milestones
  for (const def of MILESTONE_DEFS) {
    if (exportCount.value === def.count && !unlockedMilestones.value.includes(def.key)) {
      unlockedMilestones.value.push(def.key)
      showBadgeNotification(def)
    }
  }

  saveState()
}

// ========== Show Badge Notification ==========
const showBadgeNotification = (def: typeof MILESTONE_DEFS[0]) => {
  if (badgeTimer) clearTimeout(badgeTimer)

  currentBadge.value = { ...def, unlocked: true }
  showBadge.value = true

  badgeTimer = setTimeout(() => {
    showBadge.value = false
  }, 5000)
}

const dismissBadge = () => {
  showBadge.value = false
  if (badgeTimer) clearTimeout(badgeTimer)
}

// ========== Get All Milestones (with unlock status) ==========
const getMilestones = computed<Milestone[]>(() => {
  return MILESTONE_DEFS.map(def => ({
    ...def,
    unlocked: unlockedMilestones.value.includes(def.key),
  }))
})

// ========== Next milestone to unlock ==========
const nextMilestone = computed(() => {
  return MILESTONE_DEFS.find(def => !unlockedMilestones.value.includes(def.key)) || null
})

// ========== Progress percentage ==========
const progressPercent = computed(() => {
  const next = nextMilestone.value
  if (!next) return 100
  const prev = MILESTONE_DEFS[MILESTONE_DEFS.indexOf(next) - 1]
  const prevCount = prev ? prev.count : 0
  const range = next.count - prevCount
  const progress = exportCount.value - prevCount
  return Math.min(100, Math.round((progress / range) * 100))
})

// ========== Init ==========
onMounted(() => {
  loadState()
})

export function useMilestones() {
  return {
    exportCount,
    unlockedMilestones,
    showBadge,
    currentBadge,
    milestones: getMilestones,
    nextMilestone,
    progressPercent,
    trackExport,
    dismissBadge,
    getMilestones,
  }
}
