<template>
  <Teleport to="body">
    <!-- Milestone Unlock Notification -->
    <transition name="badge-slide">
      <div v-if="visible && badge" class="milestone-toast" :style="{ '--badge-color': badge.color }">
        <div class="toast-glow"></div>
        <div class="toast-content">
          <div class="toast-emoji">{{ badge.emoji }}</div>
          <div class="toast-text">
            <div class="toast-title">🏆 里程碑解锁！</div>
            <div class="toast-name">{{ badge.name }}</div>
            <div class="toast-desc">{{ badge.description }}</div>
          </div>
          <button class="toast-close" @click="dismiss">×</button>
        </div>
        <div class="toast-progress">
          <div class="toast-progress-bar"></div>
        </div>
      </div>
    </transition>

    <!-- Milestone Stats Overlay (triggered by clicking the stats icon) -->
    <el-drawer
      v-model="drawerVisible"
      title="📊 数据里程碑"
      direction="rtl"
      size="360px"
      :append-to-body="true"
    >
      <div class="milestone-panel">
        <!-- Overall progress -->
        <div class="panel-header">
          <div class="panel-stat">
            <span class="stat-num">{{ exportCount }}</span>
            <span class="stat-label">累计导出</span>
          </div>
          <div class="panel-next" v-if="nextMilestone">
            <div class="next-label">下一个</div>
            <div class="next-info">
              <span>{{ nextMilestone.emoji }} {{ nextMilestone.name }}</span>
              <span class="next-count">还需 {{ nextMilestone.count - exportCount }} 次</span>
            </div>
          </div>
        </div>

        <!-- Progress bar -->
        <div class="progress-track" v-if="nextMilestone">
          <div class="progress-bar" :style="{ width: progressPercent + '%', background: nextMilestone.color }"></div>
          <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
        </div>

        <!-- Milestone list -->
        <div class="milestone-list">
          <div
            v-for="m in milestones"
            :key="m.key"
            class="milestone-item"
            :class="{ 'is-unlocked': m.unlocked }"
            :style="m.unlocked ? { '--mc': m.color } : {}"
          >
            <div class="mi-icon">{{ m.emoji }}</div>
            <div class="mi-info">
              <div class="mi-name">{{ m.name }} <span class="mi-name-en">({{ m.nameEn }})</span></div>
              <div class="mi-desc">{{ m.description }}</div>
            </div>
            <div class="mi-status">
              <span v-if="m.unlocked" class="mi-unlocked">✓</span>
              <span v-else class="mi-locked">{{ m.count }}次</span>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMilestones } from '@/composables/useMilestones'

defineProps<{ visible: boolean; badge: any }>()
const emit = defineEmits<{
  (e: 'dismiss'): void
}>()

const {
  exportCount,
  milestones,
  nextMilestone,
  progressPercent,
  dismissBadge,
} = useMilestones()

const drawerVisible = ref(false)

const dismiss = () => {
  dismissBadge()
  emit('dismiss')
}

const openDrawer = () => {
  drawerVisible.value = true
}

defineExpose({ drawerVisible, openDrawer })
</script>

<style scoped>
/* ========== Toast Notification ========== */
.milestone-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  width: 320px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border: 1px solid var(--badge-color, #f59e0b);
  border-radius: 16px;
  overflow: hidden;
  z-index: 99997;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 0 20px color-mix(in srgb, var(--badge-color) 30%, transparent);
}

.toast-glow {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, color-mix(in srgb, var(--badge-color) 8%, transparent), transparent 50%);
  pointer-events: none;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 16px 14px;
  position: relative;
}

.toast-emoji {
  font-size: 36px;
  flex-shrink: 0;
  filter: drop-shadow(0 0 8px color-mix(in srgb, var(--badge-color) 50%, transparent));
  animation: emojiPop 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes emojiPop {
  0% { transform: scale(0) rotate(-20deg); opacity: 0; }
  100% { transform: scale(1) rotate(0deg); opacity: 1; }
}

.toast-text {
  flex: 1;
}

.toast-title {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: var(--badge-color);
  margin-bottom: 2px;
}

.toast-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 18px;
  font-weight: 800;
  color: #fff;
  margin-bottom: 2px;
}

.toast-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.toast-close {
  position: absolute;
  top: 8px;
  right: 10px;
  background: none;
  border: none;
  color: rgba(255,255,255,0.3);
  font-size: 20px;
  cursor: pointer;
  line-height: 1;
  transition: color 0.2s;
}
.toast-close:hover { color: rgba(255,255,255,0.7); }

.toast-progress {
  height: 3px;
  background: rgba(255,255,255,0.1);
}
.toast-progress-bar {
  height: 100%;
  background: var(--badge-color);
  animation: progressShrink 5s linear forwards;
  box-shadow: 0 0 8px var(--badge-color);
}
@keyframes progressShrink {
  from { width: 100%; }
  to { width: 0%; }
}

/* ========== Drawer Panel ========== */
.milestone-panel {
  padding: 0 4px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(135deg, rgba(102,126,234,0.1), rgba(118,75,162,0.1));
  border-radius: 12px;
  margin-bottom: 16px;
}

.panel-stat {
  display: flex;
  flex-direction: column;
}

.stat-num {
  font-family: 'DM Sans', sans-serif;
  font-size: 36px;
  font-weight: 800;
  color: #667eea;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: rgba(255,255,255,0.4);
  margin-top: 4px;
}

.panel-next {
  text-align: right;
}

.next-label {
  font-size: 10px;
  color: rgba(255,255,255,0.3);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.next-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: 13px;
  font-weight: 600;
  color: rgba(255,255,255,0.7);
}

.next-count {
  font-size: 11px;
  color: rgba(255,255,255,0.3);
  font-weight: 400;
}

.progress-track {
  height: 4px;
  background: rgba(255,255,255,0.08);
  border-radius: 2px;
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
}

.progress-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 2px;
  transition: width 0.5s ease;
}

/* ========== Milestone List ========== */
.milestone-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.milestone-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 12px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  transition: all 0.2s ease;
  opacity: 0.45;
}

.milestone-item.is-unlocked {
  opacity: 1;
  background: rgba(255,255,255,0.06);
  border-color: color-mix(in srgb, var(--mc) 30%, transparent);
  box-shadow: 0 0 12px color-mix(in srgb, var(--mc) 10%, transparent);
}

.mi-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.mi-info {
  flex: 1;
  min-width: 0;
}

.mi-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: rgba(255,255,255,0.8);
}

.is-unlocked .mi-name {
  color: #fff;
}

.mi-name-en {
  font-size: 11px;
  color: rgba(255,255,255,0.3);
  font-weight: 400;
}

.mi-desc {
  font-size: 11px;
  color: rgba(255,255,255,0.3);
  margin-top: 2px;
}

.mi-status {
  flex-shrink: 0;
}

.mi-unlocked {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--mc, #22c55e);
  color: white;
  font-size: 12px;
  font-weight: bold;
  box-shadow: 0 2px 8px color-mix(in srgb, var(--mc) 50%, transparent);
}

.mi-locked {
  font-size: 11px;
  color: rgba(255,255,255,0.2);
  font-family: 'JetBrains Mono', monospace;
}

/* ========== Transitions ========== */
.badge-slide-enter-active {
  animation: badgeIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}
.badge-slide-leave-active {
  animation: badgeOut 0.3s ease both;
}
@keyframes badgeIn {
  0% { opacity: 0; transform: translateX(100px) scale(0.8); }
  100% { opacity: 1; transform: translateX(0) scale(1); }
}
@keyframes badgeOut {
  0% { opacity: 1; transform: translateX(0); }
  100% { opacity: 0; transform: translateX(50px); }
}
</style>
