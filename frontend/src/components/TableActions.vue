<template>
  <div class="ta-wrapper">
    <button
      v-for="action in visibleActions"
      :key="action.key"
      class="ta-btn"
      :class="`ta-btn--${action.type || 'default'}`"
      :title="action.label"
      @click="$emit('action', action.key, currentRow)"
    >
      <!-- 查看 -->
      <svg v-if="action.key === 'detail' || action.key === 'view'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
        <circle cx="12" cy="12" r="3"/>
      </svg>
      <!-- 编辑 -->
      <svg v-else-if="action.key === 'edit'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
      </svg>
      <!-- 删除 -->
      <svg v-else-if="action.key === 'delete'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="3 6 5 6 21 6"/>
        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
      </svg>
      <!-- 重置密码 -->
      <svg v-else-if="action.key === 'resetPwd'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
        <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
      </svg>
      <!-- 导出 -->
      <svg v-else-if="action.key === 'export'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
      <!-- 解锁 -->
      <svg v-else-if="action.key === 'unlock'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
        <path d="M7 11V7a5 5 0 0 1 9.9-1"/>
      </svg>
      <!-- 复制 -->
      <svg v-else-if="action.key === 'copy'" class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
      </svg>
      <!-- 其他默认 -->
      <svg v-else class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/>
      </svg>
    </button>

    <!-- 更多下拉 -->
    <el-dropdown v-if="overflowActions.length > 0" trigger="click" @command="(key:string) => $emit('action', key, currentRow)" popper-class="ta-dropdown-popper">
      <button class="ta-btn ta-btn--more" title="更多">
        <svg class="ta-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="5" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="19" r="1"/>
        </svg>
      </button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-for="a in overflowActions" :key="a.key" :command="a.key">
            <span class="dropdown-item-label">{{ a.label }}</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Action {
  key: string
  label: string
  type?: 'primary' | 'danger' | 'warning' | 'default' | 'info'
}

const props = withDefaults(defineProps<{
  actions: Action[]
  maxVisible?: number
  row?: any
}>(), {
  maxVisible: 5
})

const visibleActions = computed(() => props.actions.slice(0, props.maxVisible))
const overflowActions = computed(() => props.actions.slice(props.maxVisible))

// Use computed so row is always reactive when props.row changes
const currentRow = computed(() => props.row)

defineEmits<{
  action: [key: string, row?: any]
}>()
</script>

<style scoped lang="scss">
.ta-wrapper {
  display: flex;
  align-items: center;
  gap: 2px;
  justify-content: center;
}

.ta-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 1px solid transparent;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: var(--color-text-muted);
  transition: all 0.15s ease;
  flex-shrink: 0;

  &:hover { transform: translateY(-1px); }
  &:active { transform: scale(0.95); }

  // 查看/默认 — 蓝色
  &--default {
    &:hover {
      color: var(--color-primary);
      background: var(--color-primary-light-9);
      border-color: var(--color-primary-light-7);
    }
  }

  // 编辑 — 蓝色
  &--primary {
    &:hover {
      color: var(--color-primary);
      background: var(--color-primary-light-9);
      border-color: var(--color-primary-light-7);
    }
  }

  // 删除 — 红色
  &--danger {
    &:hover {
      color: var(--color-danger);
      background: rgba(220, 38, 38, 0.06);
      border-color: rgba(220, 38, 38, 0.2);
    }
  }

  // 警告 — 橙色
  &--warning {
    &:hover {
      color: var(--color-warning);
      background: rgba(245, 158, 11, 0.06);
      border-color: rgba(245, 158, 11, 0.2);
    }
  }

  // 更多
  &--more {
    &:hover {
      color: var(--color-text-primary);
      background: var(--gray-100);
    }
  }
}

.ta-icon {
  flex-shrink: 0;
  pointer-events: none;
}

.dropdown-item-label {
  font-size: 13px;
  color: var(--color-text-primary);
}
</style>
